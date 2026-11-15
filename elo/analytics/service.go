// Package analytics provides the statistics operations shared by API transports.
package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/wlbr/cselo/elo"
)

// Source supplies the data used by the analytics service.
type Source interface {
	GetPlayers(context.Context, int, int) ([]*elo.Player, error)
	SearchPlayers(context.Context, string, int, int) ([]*elo.Player, error)
	FindPlayerByID(context.Context, int64) (*elo.Player, error)
	GetPlayerStats(context.Context, int64) (*elo.PlayerStats, error)
	GetMatches(context.Context, int, int) ([]*elo.Match, error)
	GetMatchesForPeriod(context.Context, *time.Time, *time.Time, int, int) ([]*elo.Match, error)
	GetMatchByID(context.Context, int64) (*elo.Match, error)
	GetMatchParticipants(context.Context, int64) ([]*elo.MatchParticipant, error)
	GetPlayerStatsForPeriod(context.Context, int64, *time.Time, *time.Time) (*elo.PlayerStats, error)
	GetPlayerResultsForPeriod(context.Context, *time.Time, *time.Time) ([]*elo.PlayerResult, error)
	GetHeadToHeadForPeriod(context.Context, int64, int64, *time.Time, *time.Time) (*elo.HeadToHeadStats, error)
	GetWeaponStatsForPeriod(context.Context, int64, *time.Time, *time.Time) ([]*elo.WeaponStats, error)
	GetMapStatsForPeriod(context.Context, int64, *time.Time, *time.Time) ([]*elo.MapStats, error)
}

// Service validates analytics requests and delegates them to a Source.
// Periods are half-open intervals: start is inclusive and end is exclusive.
// Either bound may be nil; both nil means all time.
// Paginated requests require a limit between 1 and 100 and a non-negative offset.
type Service struct {
	source Source
}

// New creates an analytics service. source supplies the players, matches,
// and statistics queried by the service.
func New(source Source) *Service { return &Service{source: source} }

func validatePage(limit, offset int) error {
	if limit < 1 || limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	if offset < 0 {
		return fmt.Errorf("offset must not be negative")
	}
	return nil
}

func validateID(id int64) error {
	if id < 1 {
		return fmt.Errorf("id must be a positive integer")
	}
	return nil
}

// Periods are expressed as a half-open interval [start, end).
func validatePeriod(start, end *time.Time) error {
	if start != nil && end != nil && !start.Before(*end) {
		return fmt.Errorf("period must have a start before its end")
	}
	return nil
}

// GetPlayers returns a page of players. ctx controls the request lifetime;
// limit is the maximum number of results and offset is the number to skip.
func (s *Service) GetPlayers(ctx context.Context, limit, offset int) ([]*elo.Player, error) {
	if err := validatePage(limit, offset); err != nil {
		return nil, err
	}
	return s.source.GetPlayers(ctx, limit, offset)
}

// SearchPlayers returns a page of players matching the non-empty name fragment.
// ctx controls the request lifetime; limit is the maximum number of results
// and offset is the number to skip.
func (s *Service) SearchPlayers(ctx context.Context, name string, limit, offset int) ([]*elo.Player, error) {
	if name == "" {
		return nil, fmt.Errorf("player name must not be empty")
	}
	if err := validatePage(limit, offset); err != nil {
		return nil, err
	}
	return s.source.SearchPlayers(ctx, name, limit, offset)
}

// FindPlayerByID looks up a player by positive id. ctx controls the request
// lifetime.
func (s *Service) FindPlayerByID(ctx context.Context, id int64) (*elo.Player, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	return s.source.FindPlayerByID(ctx, id)
}

// GetPlayerStats returns all-time statistics for the player identified by
// positive id. ctx controls the request lifetime.
func (s *Service) GetPlayerStats(ctx context.Context, id int64) (*elo.PlayerStats, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	return s.source.GetPlayerStats(ctx, id)
}

// GetMatches returns a page of matches. ctx controls the request lifetime;
// limit is the maximum number of results and offset is the number to skip.
func (s *Service) GetMatches(ctx context.Context, limit, offset int) ([]*elo.Match, error) {
	if err := validatePage(limit, offset); err != nil {
		return nil, err
	}
	return s.source.GetMatches(ctx, limit, offset)
}

// GetMatchesForPeriod returns a page of matches starting in [start, end).
// ctx controls the request lifetime; start is inclusive and end is exclusive.
// Either bound may be nil to leave that side unbounded.
// limit is the maximum number of results and offset is the number to skip.
func (s *Service) GetMatchesForPeriod(ctx context.Context, start, end *time.Time, limit, offset int) ([]*elo.Match, error) {
	if err := validatePeriod(start, end); err != nil {
		return nil, err
	}
	if err := validatePage(limit, offset); err != nil {
		return nil, err
	}
	return s.source.GetMatchesForPeriod(ctx, start, end, limit, offset)
}

// GetMatchByID looks up a match by positive id. ctx controls the request
// lifetime.
func (s *Service) GetMatchByID(ctx context.Context, id int64) (*elo.Match, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	return s.source.GetMatchByID(ctx, id)
}

// GetMatchParticipants returns participants and their statistics for the match
// identified by positive id. ctx controls the request lifetime.
func (s *Service) GetMatchParticipants(ctx context.Context, id int64) ([]*elo.MatchParticipant, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	return s.source.GetMatchParticipants(ctx, id)
}

// GetPlayerStatsForPeriod returns statistics for the player identified by
// positive id. ctx controls the request lifetime; start is inclusive and end
// is exclusive when selecting events. Either bound may be nil.
func (s *Service) GetPlayerStatsForPeriod(ctx context.Context, id int64, start, end *time.Time) (*elo.PlayerStats, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validatePeriod(start, end); err != nil {
		return nil, err
	}
	return s.source.GetPlayerStatsForPeriod(ctx, id, start, end)
}

// GetPlayerResultsForPeriod returns player results for events in [start, end).
// ctx controls the request lifetime; start is inclusive and end is exclusive.
// Either bound may be nil.
func (s *Service) GetPlayerResultsForPeriod(ctx context.Context, start, end *time.Time) ([]*elo.PlayerResult, error) {
	if err := validatePeriod(start, end); err != nil {
		return nil, err
	}
	return s.source.GetPlayerResultsForPeriod(ctx, start, end)
}

// GetHeadToHeadForPeriod returns kills and headshots between two players.
// first and second are positive player IDs; ctx controls the request lifetime.
// start is inclusive and end is exclusive when selecting events; either may be nil.
func (s *Service) GetHeadToHeadForPeriod(ctx context.Context, first, second int64, start, end *time.Time) (*elo.HeadToHeadStats, error) {
	if err := validateID(first); err != nil {
		return nil, err
	}
	if err := validateID(second); err != nil {
		return nil, err
	}
	if err := validatePeriod(start, end); err != nil {
		return nil, err
	}
	return s.source.GetHeadToHeadForPeriod(ctx, first, second, start, end)
}

// GetWeaponStatsForPeriod returns weapon statistics for the player identified
// by positive id. ctx controls the request lifetime; start is inclusive and
// end is exclusive when selecting events; either may be nil.
func (s *Service) GetWeaponStatsForPeriod(ctx context.Context, id int64, start, end *time.Time) ([]*elo.WeaponStats, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validatePeriod(start, end); err != nil {
		return nil, err
	}
	return s.source.GetWeaponStatsForPeriod(ctx, id, start, end)
}

// GetMapStatsForPeriod returns map statistics for the player identified by
// positive id. ctx controls the request lifetime; start is inclusive and end
// is exclusive when selecting events; either may be nil.
func (s *Service) GetMapStatsForPeriod(ctx context.Context, id int64, start, end *time.Time) ([]*elo.MapStats, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validatePeriod(start, end); err != nil {
		return nil, err
	}
	return s.source.GetMapStatsForPeriod(ctx, id, start, end)
}
