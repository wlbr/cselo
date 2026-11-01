package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/wlbr/cselo/elo"
)

const playerStatsQuery = `
WITH participation AS (
	SELECT match FROM kills WHERE actor = $1 OR victim = $1
	UNION
	SELECT match FROM assists WHERE actor = $1 OR victim = $1
	UNION
	SELECT match FROM blindings WHERE actor = $1 OR victim = $1
	UNION
	SELECT match FROM grenadethrows WHERE actor = $1
	UNION
	SELECT match FROM scoreaction WHERE actor = $1
)
SELECT
	(SELECT COUNT(*) FROM participation),
	(SELECT COUNT(*) FROM kills WHERE actor = $1),
	(SELECT COUNT(*) FROM kills WHERE victim = $1),
	(SELECT COUNT(*) FROM assists WHERE actor = $1),
	(SELECT COUNT(*) FROM kills WHERE actor = $1 AND headshot)
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
WITH participants AS (
	SELECT actor AS player_id FROM kills WHERE match = $1
	UNION SELECT victim FROM kills WHERE match = $1
	UNION SELECT actor FROM assists WHERE match = $1
	UNION SELECT victim FROM assists WHERE match = $1
	UNION SELECT actor FROM blindings WHERE match = $1
	UNION SELECT victim FROM blindings WHERE match = $1
	UNION SELECT actor FROM grenadethrows WHERE match = $1
	UNION SELECT actor FROM scoreaction WHERE match = $1
)
SELECT
	p.id, p.initialname, p.steamid, COALESCE(p.profileid, ''),
	(SELECT COUNT(*) FROM kills WHERE match = $1 AND actor = p.id),
	(SELECT COUNT(*) FROM kills WHERE match = $1 AND victim = p.id),
	(SELECT COUNT(*) FROM assists WHERE match = $1 AND actor = p.id),
	(SELECT COUNT(*) FROM kills WHERE match = $1 AND actor = p.id AND headshot)
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
