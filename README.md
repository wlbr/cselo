# cselo

cselo collects Counter-Strike 2 (CS2) log data to build player and match
statistics. `elohttp` receives and processes game logs, storing events in
PostgreSQL. `elointel` reads that data and makes statistics available through
GraphQL and MCP.

## elohttp: log receiver

`elohttp` accepts CS2 log lines in HTTP request bodies, parses events such as
kills, assists, rounds, and matches, and stores them in PostgreSQL. It listens
on port `42820` by default; use `-port` or `[Elo] Port` to change it.

Configure PostgreSQL in an INI file based on `cselo.ini.TEMPLATE` and prepare
the database schema before starting the receiver. **`scripts/create-tables.sql`
drops existing tables and their data**, so use it only for a new database or
when intentionally resetting one.

```bash
go build -o elohttp ./cmd/elohttp
./elohttp -cfg cselo-local.ini
```

To process an existing log file instead of listening for HTTP requests:

```bash
./elohttp -cfg cselo-local.ini -import data/test.log
```

## elointel statistics server

`elointel` serves GraphQL at `POST /graphql` and Streamable HTTP MCP at `/mcp`
on the same listener. Both use a shared analytics service and PostgreSQL source.
Player statistics are calculated from imported event tables, including historical
matches, without a separate aggregation job.

### GraphQL

Send queries to `POST /graphql`. For example:

```bash
curl -X POST http://localhost:42821/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"{ players(name: \"Jagger\", limit: 10) { id name steamid } }"}'
```

```graphql
query PlayerProfile($id: ID!) {
  player(id: $id) {
    name
    steamid
    stats(startDate: "2026-10-01", endDate: "2026-10-31") {
      matches
      kills
      deaths
      assists
      headshots
      headshotRate
      killDeathRatio
      killsPerMatch
    }
  }
}
```

Variables (replace `1` with a player ID from `players`):

```json
{"id": "1"}
```

For example, the same query can be sent with curl:

```bash
curl -X POST http://localhost:42821/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"query PlayerProfile($id: ID!) { player(id: $id) { name steamid stats(startDate: \"2026-10-01\") { matches kills deaths } } }","variables":{"id":"1"}}'
```

For completed matches, reports include observed participants and their
match-specific statistics. Incomplete matches have no reported participants:

```graphql
query MatchReport($id: ID!) {
  match(id: $id) {
    mapname
    scorea
    scoreb
    durationSeconds
    participants {
      player { name steamid }
      stats { kills deaths assists headshots }
    }
  }
}
```

Variables (replace `6448` with a match ID from `matches`):

```json
{"id": "6448"}
```

`players` and `matches` accept `limit` (1-100, default 50) and `offset` arguments.
`players` also accepts an optional `name` substring filter. `matches` accepts
`startDate` and `endDate` to filter by match start date, and `Player.stats`
accepts them to filter events. `player(id: ID!)` and `match(id: ID!)` return
`null` when the ID is not found.

```graphql
query RecentMatches {
  matches(startDate: "2026-10-01") {
    id
    mapname
    start
  }
}
```

To browse all matches, request successive pages by increasing `offset` by the
page size until a page is empty:

```graphql
query MatchPage {
  matches(limit: 100, offset: 0) {
    id
    mapname
    scorea
    scoreb
    completed
  }
}
```

### MCP

The MCP endpoint replaces the former standalone `elomcp` program and exposes
read-only player and match statistics tools:

| Question | MCP tool |
| --- | --- |
| Find a player by name | `search_players` |
| Show all-time player statistics | `get_player_profile` |
| Compare who killed whom between two players for a date range | `get_head_to_head` |
| Show all players' results for one day | `get_daily_player_results` |
| Show the player-results table for a date range | `get_player_results_for_period` |
| Show one player's statistics over a date range | `get_player_period_stats` |
| Show one player's weapon or map performance | `get_player_weapon_stats`, `get_player_map_stats` |
| Browse matches or inspect a match | `list_matches`, `get_match_report` |

Tool arguments (IDs are positive numeric cselo IDs encoded as strings):

| Tool | Required arguments | Optional arguments |
| --- | --- | --- |
| `search_players` | `name` | `limit`, `offset` |
| `get_player_profile` | `player_id` | `start_date`, `end_date` |
| `get_head_to_head` | `first_player_id`, `second_player_id` | `start_date`, `end_date` |
| `get_daily_player_results` | `date` | — |
| `get_player_results_for_period` | — | `start_date`, `end_date` |
| `get_player_period_stats` | `player_id` | `start_date`, `end_date` |
| `get_player_weapon_stats` | `player_id` | `start_date`, `end_date` |
| `get_player_map_stats` | `player_id` | `start_date`, `end_date` |
| `list_matches` | — | `start_date`, `end_date`, `limit`, `offset` |
| `get_match_report` | `match_id` | — |

For `search_players` and `list_matches`, `limit` defaults to 20 and may be
1–100; `offset` defaults to 0 and must be non-negative. Increase `offset` to
retrieve subsequent pages.

Examples of MCP `tools/call` request payloads for a connected client (replace IDs
with values from your data):

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search_players","arguments":{"name":"Jagger","limit":10}}}
```

```json
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_player_profile","arguments":{"player_id":"1","start_date":"2026-10-01"}}}
```

```json
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_head_to_head","arguments":{"first_player_id":"1","second_player_id":"2","start_date":"2026-10-01","end_date":"2026-10-31"}}}
```

All dates are UTC `YYYY-MM-DD` values. Date ranges include both the start and
end date. For time-based statistics and match lists, either bound may be
omitted; omitting both selects all time. MCP uses `start_date` and `end_date`,
while GraphQL uses `startDate` and `endDate`. The daily-results tool still
requires its single `date`. Player searches and individual match reports are
not date-filtered. Match lists and reports can include incomplete matches, but
event-derived player statistics include only events from completed matches.

### Statistics

Player statistics are derived from recorded events in completed matches. The
`matches` count is the number of distinct completed matches with observed
participation; for a date range, the participation event must fall within that
range. Match-list date filters instead use the match start date.

The ratio fields are calculated as follows:

| Field | Formula | If the denominator is zero |
| --- | --- | --- |
| `headshotRate` | headshots / kills | 0 |
| `killDeathRatio` | kills / deaths | 0 |
| `killsPerMatch` | kills / matches | 0 |

`headshotRate` is a fraction, not a percentage. Match-participant statistics
use a match count of 1 for that completed match.

### Configuration and startup

Use the same `-cfg` INI format as `elohttp`. The `[MCP]` host and port configure
the listener for **both** APIs (defaults: `localhost` and `42821`):

```ini
[MCP]
Host = localhost
Port = 42821
Path = /mcp
; Optional: uncomment and set both values to protect MCP with HTTP Basic Auth.
; BasicAuthUser = mcp-user
; BasicAuthPassword = change-me
```

The old `[GraphQL] Port` and `-graphql-port` settings do not control elointel.
The MCP path is fixed at `/mcp`; stdio transport is not supported.
Build and run:

```bash
go build -o elointel ./cmd/elointel
./elointel -cfg cselo-local.ini
```

For example, connect an MCP client to the running server without Basic Auth:

```json
{
  "mcpServers": {
    "cselo-stats": {
      "type": "http",
      "url": "http://localhost:42821/mcp"
    }
  }
}
```

If Basic Auth is enabled, configure the client to send an `Authorization`
header on MCP requests. For a client supporting a `headers` option, add:

```json
"headers": {
  "Authorization": "Basic <base64(username:password)>"
}
```

Replace the placeholder with the Base64 encoding of the configured username,
a colon, and password. Client configuration syntax varies; store credentials
securely rather than committing them to a configuration file.

Basic Auth protects only `/mcp`, not GraphQL. To expose the server remotely,
configure `MCP.Host` and protect both endpoints at a reverse proxy; use HTTPS
to avoid transmitting credentials insecurely.

## Code

* Documentation: https://godoc.org/github.com/wlbr/cselo
* Lint: http://go-lint.appspot.com/github.com/wlbr/cselo
* Continuous Integration: [![GoBuild](https://github.com/wlbr/cselo/actions/workflows/go-build.yml/badge.svg)](https://github.com/wlbr/cselo/actions/workflows/go-build.yml)
* Test Coverage: [![Coverage Status](https://coveralls.io/repos/github/wlbr/cselo/badge.svg?branch=master)](https://coveralls.io/github/wlbr/cselo?branch=master)
* Metrics: [![GoReportCard](https://goreportcard.com/badge/github.com/wlbr/cselo)](https://goreportcard.com/report/github.com/wlbr/cselo)
