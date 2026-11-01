package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/wlbr/cselo/elo"
)

func (s *Postgres) GetPlayerStatsForPeriod(ctx context.Context, playerID int64, start, end time.Time) (*elo.PlayerStats, error) {
	var matches, kills, deaths, assists, headshots int
	err := s.Db.QueryRow(ctx, `
WITH completed_matches AS (
	SELECT id FROM matches WHERE completed
),
participation AS (
	SELECT k.match FROM kills k JOIN completed_matches m ON m.id = k.match WHERE (k.actor = $1 OR k.victim = $1) AND k.timestmp >= $2 AND k.timestmp < $3
	UNION SELECT a.match FROM assists a JOIN completed_matches m ON m.id = a.match WHERE (a.actor = $1 OR a.victim = $1) AND a.timestmp >= $2 AND a.timestmp < $3
	UNION SELECT b.match FROM blindings b JOIN completed_matches m ON m.id = b.match WHERE (b.actor = $1 OR b.victim = $1) AND b.timestmp >= $2 AND b.timestmp < $3
	UNION SELECT g.match FROM grenadethrows g JOIN completed_matches m ON m.id = g.match WHERE g.actor = $1 AND g.timestmp >= $2 AND g.timestmp < $3
	UNION SELECT s.match FROM scoreaction s JOIN completed_matches m ON m.id = s.match WHERE s.actor = $1 AND s.timestmp >= $2 AND s.timestmp < $3
)
SELECT
	(SELECT COUNT(*) FROM participation),
	(SELECT COUNT(*) FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.actor = $1 AND k.timestmp >= $2 AND k.timestmp < $3),
	(SELECT COUNT(*) FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.victim = $1 AND k.timestmp >= $2 AND k.timestmp < $3),
	(SELECT COUNT(*) FROM assists a JOIN completed_matches m ON m.id = a.match WHERE a.actor = $1 AND a.timestmp >= $2 AND a.timestmp < $3),
	(SELECT COUNT(*) FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.actor = $1 AND k.headshot AND k.timestmp >= $2 AND k.timestmp < $3)
`, playerID, start, end).Scan(&matches, &kills, &deaths, &assists, &headshots)
	if err != nil {
		return nil, fmt.Errorf("get statistics for player %d from %s to %s: %w", playerID, start, end, err)
	}
	return elo.NewPlayerStats(matches, kills, deaths, assists, headshots), nil
}

func (s *Postgres) GetPlayerResultsForPeriod(ctx context.Context, start, end time.Time) ([]*elo.PlayerResult, error) {
	rows, err := s.Db.Query(ctx, `
WITH completed_matches AS (
	SELECT id FROM matches WHERE completed
),
participation AS (
	SELECT k.actor AS player_id, k.match FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.timestmp >= $1 AND k.timestmp < $2
	UNION SELECT k.victim, k.match FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.timestmp >= $1 AND k.timestmp < $2
	UNION SELECT a.actor, a.match FROM assists a JOIN completed_matches m ON m.id = a.match WHERE a.timestmp >= $1 AND a.timestmp < $2
	UNION SELECT a.victim, a.match FROM assists a JOIN completed_matches m ON m.id = a.match WHERE a.timestmp >= $1 AND a.timestmp < $2
	UNION SELECT b.actor, b.match FROM blindings b JOIN completed_matches m ON m.id = b.match WHERE b.timestmp >= $1 AND b.timestmp < $2
	UNION SELECT b.victim, b.match FROM blindings b JOIN completed_matches m ON m.id = b.match WHERE b.timestmp >= $1 AND b.timestmp < $2
	UNION SELECT g.actor, g.match FROM grenadethrows g JOIN completed_matches m ON m.id = g.match WHERE g.timestmp >= $1 AND g.timestmp < $2
	UNION SELECT s.actor, s.match FROM scoreaction s JOIN completed_matches m ON m.id = s.match WHERE s.timestmp >= $1 AND s.timestmp < $2
),
match_counts AS (
	SELECT player_id, COUNT(*) AS matches FROM participation GROUP BY player_id
),
kill_counts AS (
	SELECT k.actor AS player_id, COUNT(*) AS kills, COUNT(*) FILTER (WHERE k.headshot) AS headshots
	FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.timestmp >= $1 AND k.timestmp < $2 GROUP BY k.actor
),
death_counts AS (
	SELECT k.victim AS player_id, COUNT(*) AS deaths
	FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.timestmp >= $1 AND k.timestmp < $2 GROUP BY k.victim
),
assist_counts AS (
	SELECT a.actor AS player_id, COUNT(*) AS assists
	FROM assists a JOIN completed_matches m ON m.id = a.match WHERE a.timestmp >= $1 AND a.timestmp < $2 GROUP BY a.actor
)
SELECT p.id, p.initialname, p.steamid, COALESCE(p.profileid, ''),
	mc.matches, COALESCE(kc.kills, 0), COALESCE(dc.deaths, 0),
	COALESCE(ac.assists, 0), COALESCE(kc.headshots, 0)
FROM match_counts mc
JOIN players p ON p.id = mc.player_id
LEFT JOIN kill_counts kc ON kc.player_id = p.id
LEFT JOIN death_counts dc ON dc.player_id = p.id
LEFT JOIN assist_counts ac ON ac.player_id = p.id
ORDER BY COALESCE(kc.kills, 0) DESC, COALESCE(dc.deaths, 0), p.initialname`, start, end)
	if err != nil {
		return nil, fmt.Errorf("get player results from %s to %s: %w", start, end, err)
	}
	defer rows.Close()

	var results []*elo.PlayerResult
	for rows.Next() {
		player := &elo.Player{}
		var matches, kills, deaths, assists, headshots int
		if err := rows.Scan(
			&player.ID, &player.Name, &player.SteamID, &player.ProfileID,
			&matches, &kills, &deaths, &assists, &headshots,
		); err != nil {
			return nil, fmt.Errorf("scan player result: %w", err)
		}
		results = append(results, &elo.PlayerResult{
			Player: player,
			Stats:  elo.NewPlayerStats(matches, kills, deaths, assists, headshots),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate player results: %w", err)
	}
	return results, nil
}

func (s *Postgres) GetHeadToHeadForPeriod(ctx context.Context, firstPlayerID, secondPlayerID int64, start, end time.Time) (*elo.HeadToHeadStats, error) {
	stats := &elo.HeadToHeadStats{}
	err := s.Db.QueryRow(ctx, `
SELECT
	COUNT(*) FILTER (WHERE k.actor = $1 AND k.victim = $2),
	COUNT(*) FILTER (WHERE k.actor = $1 AND k.victim = $2 AND k.headshot),
	COUNT(*) FILTER (WHERE k.actor = $2 AND k.victim = $1),
	COUNT(*) FILTER (WHERE k.actor = $2 AND k.victim = $1 AND k.headshot)
FROM kills k
JOIN matches m ON m.id = k.match
WHERE m.completed
	AND ((k.actor = $1 AND k.victim = $2) OR (k.actor = $2 AND k.victim = $1))
	AND k.timestmp >= $3 AND k.timestmp < $4`,
		firstPlayerID, secondPlayerID, start, end,
	).Scan(&stats.FirstKills, &stats.FirstHeadshots, &stats.SecondKills, &stats.SecondHeadshots)
	if err != nil {
		return nil, fmt.Errorf("get head-to-head statistics for players %d and %d from %s to %s: %w", firstPlayerID, secondPlayerID, start, end, err)
	}
	return stats, nil
}

func (s *Postgres) GetWeaponStatsForPeriod(ctx context.Context, playerID int64, start, end time.Time) ([]*elo.WeaponStats, error) {
	rows, err := s.Db.Query(ctx, `
SELECT COALESCE(k.weapon, 'unknown'), COUNT(*), COUNT(*) FILTER (WHERE k.headshot)
FROM kills k
JOIN matches m ON m.id = k.match
WHERE m.completed AND k.actor = $1 AND k.timestmp >= $2 AND k.timestmp < $3
GROUP BY k.weapon
ORDER BY COUNT(*) DESC, k.weapon`, playerID, start, end)
	if err != nil {
		return nil, fmt.Errorf("get weapon statistics for player %d: %w", playerID, err)
	}
	defer rows.Close()

	var stats []*elo.WeaponStats
	for rows.Next() {
		stat := &elo.WeaponStats{}
		if err := rows.Scan(&stat.Weapon, &stat.Kills, &stat.Headshots); err != nil {
			return nil, fmt.Errorf("scan weapon statistics: %w", err)
		}
		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate weapon statistics: %w", err)
	}
	return stats, nil
}

func (s *Postgres) GetMapStatsForPeriod(ctx context.Context, playerID int64, start, end time.Time) ([]*elo.MapStats, error) {
	rows, err := s.Db.Query(ctx, `
WITH completed_matches AS (
	SELECT id FROM matches WHERE completed
),
participation AS (
	SELECT k.match FROM kills k JOIN completed_matches cm ON cm.id = k.match WHERE (k.actor = $1 OR k.victim = $1) AND k.timestmp >= $2 AND k.timestmp < $3
	UNION SELECT a.match FROM assists a JOIN completed_matches cm ON cm.id = a.match WHERE (a.actor = $1 OR a.victim = $1) AND a.timestmp >= $2 AND a.timestmp < $3
	UNION SELECT b.match FROM blindings b JOIN completed_matches cm ON cm.id = b.match WHERE (b.actor = $1 OR b.victim = $1) AND b.timestmp >= $2 AND b.timestmp < $3
	UNION SELECT g.match FROM grenadethrows g JOIN completed_matches cm ON cm.id = g.match WHERE g.actor = $1 AND g.timestmp >= $2 AND g.timestmp < $3
	UNION SELECT s.match FROM scoreaction s JOIN completed_matches cm ON cm.id = s.match WHERE s.actor = $1 AND s.timestmp >= $2 AND s.timestmp < $3
),
match_counts AS (
	SELECT COALESCE(m.mapname, '') AS map_name, COUNT(*) AS matches
	FROM participation p JOIN matches m ON m.id = p.match
	GROUP BY m.mapname
),
kill_counts AS (
	SELECT COALESCE(m.mapname, '') AS map_name, COUNT(*) AS kills, COUNT(*) FILTER (WHERE k.headshot) AS headshots
	FROM kills k JOIN matches m ON m.id = k.match
	WHERE m.completed AND k.actor = $1 AND k.timestmp >= $2 AND k.timestmp < $3
	GROUP BY m.mapname
),
death_counts AS (
	SELECT COALESCE(m.mapname, '') AS map_name, COUNT(*) AS deaths
	FROM kills k JOIN matches m ON m.id = k.match
	WHERE m.completed AND k.victim = $1 AND k.timestmp >= $2 AND k.timestmp < $3
	GROUP BY m.mapname
),
assist_counts AS (
	SELECT COALESCE(m.mapname, '') AS map_name, COUNT(*) AS assists
	FROM assists a JOIN matches m ON m.id = a.match
	WHERE m.completed AND a.actor = $1 AND a.timestmp >= $2 AND a.timestmp < $3
	GROUP BY m.mapname
)
SELECT mc.map_name, mc.matches, COALESCE(kc.kills, 0), COALESCE(dc.deaths, 0),
	COALESCE(ac.assists, 0), COALESCE(kc.headshots, 0)
FROM match_counts mc
LEFT JOIN kill_counts kc ON kc.map_name = mc.map_name
LEFT JOIN death_counts dc ON dc.map_name = mc.map_name
LEFT JOIN assist_counts ac ON ac.map_name = mc.map_name
ORDER BY COALESCE(kc.kills, 0) DESC, mc.map_name`, playerID, start, end)
	if err != nil {
		return nil, fmt.Errorf("get map statistics for player %d: %w", playerID, err)
	}
	defer rows.Close()

	var stats []*elo.MapStats
	for rows.Next() {
		stat := &elo.MapStats{}
		if err := rows.Scan(
			&stat.MapName, &stat.Matches, &stat.Kills, &stat.Deaths, &stat.Assists, &stat.Headshots,
		); err != nil {
			return nil, fmt.Errorf("scan map statistics: %w", err)
		}
		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate map statistics: %w", err)
	}
	return stats, nil
}
