package postgresql

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestOptionalDateBoundsAgainstPostgres(t *testing.T) {
	url := os.Getenv("CSELO_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set CSELO_TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	config.MaxConns = 1
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	_, err = pool.Exec(ctx, `
CREATE TEMP TABLE players (id bigint, initialname text, steamid text, profileid text);
CREATE TEMP TABLE matches (
  id bigint, gamemode text, mapgroup text, mapfullname text, mapname text,
  scorea int, scoreb int, rounds int, matchstart timestamp, matchend timestamp,
  duration interval, completed bool
);
CREATE TEMP TABLE kills (actor bigint, victim bigint, match bigint, timestmp timestamp, headshot bool, weapon text);
CREATE TEMP TABLE assists (actor bigint, victim bigint, match bigint, timestmp timestamp);
CREATE TEMP TABLE blindings (actor bigint, victim bigint, match bigint, timestmp timestamp);
CREATE TEMP TABLE grenadethrows (actor bigint, match bigint, timestmp timestamp);
CREATE TEMP TABLE scoreaction (actor bigint, match bigint, timestmp timestamp);
INSERT INTO players VALUES (1, 'Alice', 'alice', ''), (2, 'Bob', 'bob', '');
INSERT INTO matches VALUES
  (1, '', '', '', 'de_a', 1, 0, 1, '2026-10-01', '2026-10-01', interval '1 minute', true),
  (2, '', '', '', 'de_b', 1, 0, 1, '2026-10-03', '2026-10-03', interval '1 minute', true),
  (3, '', '', '', 'de_c', 0, 0, 0, '2026-10-03', '2026-10-03', interval '0', false);
INSERT INTO kills VALUES
  (1, 2, 1, '2026-10-01 23:59:59', true, 'ak'),
  (1, 2, 1, NULL, false, 'ak'),
  (1, 2, 2, '2026-10-03 00:00:00', false, 'awp'),
  (2, 1, 2, '2026-10-03 00:00:00', false, 'ak'),
  (1, 2, 3, '2026-10-03 00:00:00', false, 'awp');
`)
	if err != nil {
		t.Fatal(err)
	}
	source := &Postgres{Db: pool}
	cutoff := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	firstDay := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name               string
		start, end         *time.Time
		kills, deaths      int
		matches, maps      int
		first, second      int
		weapons, matchList int
	}{
		{name: "all time", kills: 3, deaths: 1, matches: 2, maps: 2, first: 3, second: 1, weapons: 2, matchList: 3},
		{name: "start only", start: &cutoff, kills: 1, deaths: 1, matches: 1, maps: 1, first: 1, second: 1, weapons: 1, matchList: 2},
		{name: "end only", end: &cutoff, kills: 1, deaths: 0, matches: 1, maps: 1, first: 1, second: 0, weapons: 1, matchList: 1},
		{name: "both bounds", start: &firstDay, end: &cutoff, kills: 1, deaths: 0, matches: 1, maps: 1, first: 1, second: 0, weapons: 1, matchList: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := source.GetPlayerStatsForPeriod(ctx, 1, tt.start, tt.end)
			if err != nil {
				t.Fatal(err)
			}
			if stats.Kills != tt.kills || stats.Deaths != tt.deaths || stats.Matches != tt.matches {
				t.Errorf("player stats = %+v", stats)
			}
			if tt.start == nil && tt.end == nil {
				allTime, err := source.GetPlayerStats(ctx, 1)
				if err != nil {
					t.Fatal(err)
				}
				if *stats != *allTime {
					t.Errorf("unbounded stats = %+v, all-time stats = %+v", stats, allTime)
				}
			}
			results, err := source.GetPlayerResultsForPeriod(ctx, tt.start, tt.end)
			if err != nil {
				t.Fatal(err)
			}
			if len(results) != 2 || results[0].Player.ID != 1 || results[0].Stats.Kills != tt.kills {
				t.Errorf("player results = %+v", results)
			}
			head, err := source.GetHeadToHeadForPeriod(ctx, 1, 2, tt.start, tt.end)
			if err != nil {
				t.Fatal(err)
			}
			if head.FirstKills != tt.first || head.SecondKills != tt.second {
				t.Errorf("head-to-head = %+v", head)
			}
			weapons, err := source.GetWeaponStatsForPeriod(ctx, 1, tt.start, tt.end)
			if err != nil {
				t.Fatal(err)
			}
			if len(weapons) != tt.weapons {
				t.Errorf("weapon stats = %+v", weapons)
			}
			maps, err := source.GetMapStatsForPeriod(ctx, 1, tt.start, tt.end)
			if err != nil {
				t.Fatal(err)
			}
			if len(maps) != tt.maps {
				t.Errorf("map stats = %+v", maps)
			}
			matches, err := source.GetMatchesForPeriod(ctx, tt.start, tt.end, 20, 0)
			if err != nil {
				t.Fatal(err)
			}
			if len(matches) != tt.matchList {
				t.Errorf("matches = %d, want %d", len(matches), tt.matchList)
			}
		})
	}
}
