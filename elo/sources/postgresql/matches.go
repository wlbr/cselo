package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/wlbr/cselo/elo"
)

const matchColumns = `
id, COALESCE(gamemode, ''), COALESCE(mapgroup, ''), COALESCE(mapfullname, ''),
COALESCE(mapname, ''), scorea, scoreb, rounds, matchstart, COALESCE(matchend, matchstart),
COALESCE(EXTRACT(EPOCH FROM duration)::bigint, 0), completed`

func (s *Postgres) GetMatches(ctx context.Context, limit, offset int) ([]*elo.Match, error) {
	return s.getMatches(ctx, `
SELECT `+matchColumns+`
FROM matches
ORDER BY matchstart DESC, id DESC
LIMIT $1 OFFSET $2`, limit, offset)
}

func (s *Postgres) GetMatchesForPeriod(ctx context.Context, start, end time.Time, limit, offset int) ([]*elo.Match, error) {
	return s.getMatches(ctx, `
SELECT `+matchColumns+`
FROM matches
WHERE matchstart >= $1 AND matchstart < $2
ORDER BY matchstart DESC, id DESC
LIMIT $3 OFFSET $4`, start, end, limit, offset)
}

func (s *Postgres) getMatches(ctx context.Context, query string, arguments ...interface{}) ([]*elo.Match, error) {
	if len(arguments) < 2 {
		return nil, fmt.Errorf("match query requires pagination arguments")
	}
	limit, ok := arguments[len(arguments)-2].(int)
	if !ok {
		return nil, fmt.Errorf("match query limit is invalid")
	}
	offset, ok := arguments[len(arguments)-1].(int)
	if !ok {
		return nil, fmt.Errorf("match query offset is invalid")
	}
	if err := validatePage(limit, offset); err != nil {
		return nil, err
	}
	rows, err := s.Db.Query(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("get matches: %w", err)
	}
	defer rows.Close()

	var matches []*elo.Match
	for rows.Next() {
		match, err := scanMatch(rows)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate matches: %w", err)
	}
	return matches, nil
}

func (s *Postgres) GetMatchByID(ctx context.Context, id int64) (*elo.Match, error) {
	match, err := scanMatch(s.Db.QueryRow(ctx, `SELECT `+matchColumns+` FROM matches WHERE id = $1`, id))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get match %d: %w", id, err)
	}
	return match, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanMatch(row rowScanner) (*elo.Match, error) {
	match := &elo.Match{}
	var matchID int64
	var matchFullName string
	var durationSeconds int64
	if err := row.Scan(
		&matchID, &match.GameMode, &match.MapGroup, &matchFullName,
		&match.MapName, &match.ScoreA, &match.ScoreB, &match.Rounds,
		&match.Start, &match.End, &durationSeconds, &match.Completed,
	); err != nil {
		return nil, err
	}
	match.SetID(matchID)
	match.Set(match.GameMode, match.MapGroup, matchFullName, match.MapName,
		match.ScoreA, match.ScoreB, match.Completed, match.End, time.Duration(durationSeconds)*time.Second)
	return match, nil
}

func validatePage(limit, offset int) error {
	if limit < 1 || limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	if offset < 0 {
		return fmt.Errorf("offset must not be negative")
	}
	return nil
}
