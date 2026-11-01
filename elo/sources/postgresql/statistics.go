package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/wlbr/cselo/elo"
)

const playerStatsQuery = `
WITH completed_matches AS (
	SELECT id FROM matches WHERE completed
),
participation AS (
	SELECT k.match FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.actor = $1 OR k.victim = $1
	UNION
	SELECT a.match FROM assists a JOIN completed_matches m ON m.id = a.match WHERE a.actor = $1 OR a.victim = $1
	UNION
	SELECT b.match FROM blindings b JOIN completed_matches m ON m.id = b.match WHERE b.actor = $1 OR b.victim = $1
	UNION
	SELECT g.match FROM grenadethrows g JOIN completed_matches m ON m.id = g.match WHERE g.actor = $1
	UNION
	SELECT s.match FROM scoreaction s JOIN completed_matches m ON m.id = s.match WHERE s.actor = $1
)
SELECT
	(SELECT COUNT(*) FROM participation),
	(SELECT COUNT(*) FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.actor = $1),
	(SELECT COUNT(*) FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.victim = $1),
	(SELECT COUNT(*) FROM assists a JOIN completed_matches m ON m.id = a.match WHERE a.actor = $1),
	(SELECT COUNT(*) FROM kills k JOIN completed_matches m ON m.id = k.match WHERE k.actor = $1 AND k.headshot)
`

func (s *Postgres) GetPlayerStats(ctx context.Context, playerID int64) (*elo.PlayerStats, error) {
	var matches, kills, deaths, assists, headshots int
	if err := s.Db.QueryRow(ctx, playerStatsQuery, playerID).Scan(
		&matches, &kills, &deaths, &assists, &headshots,
	); err != nil {
		return nil, fmt.Errorf("get statistics for player %d: %w", playerID, err)
	}
	return elo.NewPlayerStats(matches, kills, deaths, assists, headshots), nil
}

func (s *Postgres) GetMatchParticipants(ctx context.Context, matchID int64) ([]*elo.MatchParticipant, error) {
	rows, err := s.Db.Query(ctx, `
WITH completed_match AS (
	SELECT id FROM matches WHERE id = $1 AND completed
),
participants AS (
	SELECT k.actor AS player_id FROM kills k JOIN completed_match m ON m.id = k.match
	UNION SELECT k.victim FROM kills k JOIN completed_match m ON m.id = k.match
	UNION SELECT a.actor FROM assists a JOIN completed_match m ON m.id = a.match
	UNION SELECT a.victim FROM assists a JOIN completed_match m ON m.id = a.match
	UNION SELECT b.actor FROM blindings b JOIN completed_match m ON m.id = b.match
	UNION SELECT b.victim FROM blindings b JOIN completed_match m ON m.id = b.match
	UNION SELECT g.actor FROM grenadethrows g JOIN completed_match m ON m.id = g.match
	UNION SELECT s.actor FROM scoreaction s JOIN completed_match m ON m.id = s.match
)
SELECT
	p.id, p.initialname, p.steamid, COALESCE(p.profileid, ''),
	(SELECT COUNT(*) FROM kills k JOIN completed_match m ON m.id = k.match WHERE k.actor = p.id),
	(SELECT COUNT(*) FROM kills k JOIN completed_match m ON m.id = k.match WHERE k.victim = p.id),
	(SELECT COUNT(*) FROM assists a JOIN completed_match m ON m.id = a.match WHERE a.actor = p.id),
	(SELECT COUNT(*) FROM kills k JOIN completed_match m ON m.id = k.match WHERE k.actor = p.id AND k.headshot)
FROM participants
JOIN players p ON p.id = participants.player_id
ORDER BY p.initialname`, matchID)
	if err != nil {
		return nil, fmt.Errorf("get participants for match %d: %w", matchID, err)
	}
	defer rows.Close()

	var participants []*elo.MatchParticipant
	for rows.Next() {
		player := &elo.Player{}
		var kills, deaths, assists, headshots int
		if err := rows.Scan(
			&player.ID, &player.Name, &player.SteamID, &player.ProfileID,
			&kills, &deaths, &assists, &headshots,
		); err != nil {
			return nil, fmt.Errorf("scan participant for match %d: %w", matchID, err)
		}
		participants = append(participants, &elo.MatchParticipant{
			Player: player,
			Stats:  elo.NewPlayerStats(1, kills, deaths, assists, headshots),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate participants for match %d: %w", matchID, err)
	}
	return participants, nil
}

func (s *Postgres) FindPlayerByID(ctx context.Context, id int64) (*elo.Player, error) {
	player := &elo.Player{}
	err := s.Db.QueryRow(ctx,
		`SELECT id, initialname, steamid, COALESCE(profileid, '') FROM players WHERE id = $1`, id,
	).Scan(&player.ID, &player.Name, &player.SteamID, &player.ProfileID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get player %d: %w", id, err)
	}
	return player, nil
}
