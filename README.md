# cselo



## Code
* Documentation: https://godoc.org/github.com/wlbr/cselo
* Lint: http://go-lint.appspot.com/github.com/wlbr/cselo
* Continuous Integration: [![GoBuild](https://github.com/wlbr/cselo/actions/workflows/go-build.yml/badge.svg)](https://github.com/wlbr/cselo/actions/workflows/go-build.yml)
* Test Coverage: [![Coverage Status](https://coveralls.io/repos/github/wlbr/cselo/badge.svg?branch=master)](https://coveralls.io/github/wlbr/cselo?branch=master)
* Metrics: [![GoReportCard](https://goreportcard.com/badge/github.com/wlbr/cselo)](https://goreportcard.com/report/github.com/wlbr/cselo)

## GraphQL statistics API

`elogql` serves GraphQL at `POST /query` on port `8080` by default. Player statistics
are calculated from the imported event tables, so they are available for historical
matches without a separate aggregation job.

```graphql
query PlayerProfile($id: ID!) {
  player(id: $id) {
    name
    steamid
    stats {
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

Match reports include each observed participant and their match-specific statistics:

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

`players` and `matches` accept `limit` (1-100, default 50) and `offset` arguments.
`players` also accepts an optional `name` substring filter.

## MCP statistics server

`elomcp` exposes read-only player and match statistics tools. It uses the same
`elo.Config` initialization and `-cfg` INI format as `elohttp`. `elomcp` and
`elogql` share the same PostgreSQL statistics source, so neither service depends
on the other process being available.

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

All dates are UTC `YYYY-MM-DD` values. Date ranges include both the start and
end date. `get_head_to_head` requires a date range; `get_player_profile` and
`list_matches` accept either no dates (all-time/recent) or both range dates.
All event-derived statistics include only events from completed matches; live
and incomplete matches do not contribute to any reported metric.

Add the following section to the same INI file used by `elohttp`:

```ini
[MCP]
Transport = stdio
Host = localhost
Port = 42821
Path = /mcp
```

`elogql` uses its own `[GraphQL] Port` setting (default `8080`):

```ini
[GraphQL]
Port = 8080
```

Build and run it as a local stdio MCP server:

```bash
go build -o elomcp ./cmd/elomcp
./elomcp -cfg cselo-local.ini
```

For example, a desktop MCP client configuration can launch it with:

```json
{
  "mcpServers": {
    "cselo-stats": {
      "command": "/absolute/path/to/elomcp",
      "args": ["-cfg", "/absolute/path/to/cselo-local.ini"]
    }
  }
}
```

For a remote MCP endpoint, set `Transport = http` and configure `MCP.Port`.
Its Streamable HTTP endpoint is served at `Path`
(default `/mcp`) and `/healthz` returns `204 No Content`. The HTTP transport is
unauthenticated; only expose it behind an authenticated HTTPS reverse proxy, or
use it exclusively for the public statistics data.
