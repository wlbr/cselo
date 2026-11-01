package main

import (
	"context"
	"testing"
	"time"

	"github.com/wlbr/cselo/elo"
)

func TestSearchPlayers(t *testing.T) {
	tools := &statsTools{source: fakeStatisticsSource{
		players: []*elo.Player{{ID: 2, Name: "Jagger", SteamID: "STEAM_1:0:1"}},
	}}
	_, result, err := tools.searchPlayers(context.Background(), nil, searchPlayersInput{Name: "Jagger"})
	if err != nil {
		t.Fatalf("searchPlayers returned error: %v", err)
	}
	if len(result.Players) != 1 || result.Players[0].Name != "Jagger" {
		t.Fatalf("unexpected players: %+v", result.Players)
	}
}

func TestParseOptionalDateRangeRequiresBothDates(t *testing.T) {
	if _, _, _, err := parseOptionalDateRange("2026-10-01", ""); err == nil {
		t.Fatal("parseOptionalDateRange returned nil error")
	}
}

func TestGetMatchReportRejectsInvalidID(t *testing.T) {
	tools := &statsTools{}
	_, _, err := tools.getMatchReport(context.Background(), nil, matchReportInput{MatchID: "invalid"})
	if err == nil {
		t.Fatal("getMatchReport returned nil error")
	}
}

func TestNewMCPServerRegistersTools(t *testing.T) {
	server := newMCPServer("test", fakeStatisticsSource{})
	if server == nil {
		t.Fatal("newMCPServer returned nil")
	}
}

func TestParseDateRangeIncludesBothDays(t *testing.T) {
	start, end, err := parseDateRange("2026-10-01", "2026-10-03")
	if err != nil {
		t.Fatalf("parseDateRange returned error: %v", err)
	}
	if start.Format(time.DateOnly) != "2026-10-01" {
		t.Errorf("start = %s, want 2026-10-01", start)
	}
	if end.Format(time.DateOnly) != "2026-10-04" {
		t.Errorf("end = %s, want 2026-10-04", end)
	}
}

func TestParseDateRangeRejectsReversedDates(t *testing.T) {
	if _, _, err := parseDateRange("2026-10-03", "2026-10-01"); err == nil {
		t.Fatal("parseDateRange returned nil error")
	}
}

type fakeStatisticsSource struct {
	players []*elo.Player
}

func (s fakeStatisticsSource) SearchPlayers(_ context.Context, _ string, _ int, _ int) ([]*elo.Player, error) {
	return s.players, nil
}

func (s fakeStatisticsSource) FindPlayerByID(_ context.Context, _ int64) (*elo.Player, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetPlayerStats(_ context.Context, _ int64) (*elo.PlayerStats, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetMatches(_ context.Context, _ int, _ int) ([]*elo.Match, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetMatchesForPeriod(_ context.Context, _, _ time.Time, _ int, _ int) ([]*elo.Match, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetMatchByID(_ context.Context, _ int64) (*elo.Match, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetMatchParticipants(_ context.Context, _ int64) ([]*elo.MatchParticipant, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetPlayerStatsForPeriod(_ context.Context, _ int64, _, _ time.Time) (*elo.PlayerStats, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetPlayerResultsForPeriod(_ context.Context, _, _ time.Time) ([]*elo.PlayerResult, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetHeadToHeadForPeriod(_ context.Context, _, _ int64, _, _ time.Time) (*elo.HeadToHeadStats, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetWeaponStatsForPeriod(_ context.Context, _ int64, _, _ time.Time) ([]*elo.WeaponStats, error) {
	return nil, nil
}

func (s fakeStatisticsSource) GetMapStatsForPeriod(_ context.Context, _ int64, _, _ time.Time) ([]*elo.MapStats, error) {
	return nil, nil
}
