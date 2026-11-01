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
