package elo

import "testing"

func TestNewPlayerStatsCalculatesRates(t *testing.T) {
	stats := NewPlayerStats(4, 12, 8, 3, 6)

	if stats.HeadshotRate != 0.5 {
		t.Fatalf("HeadshotRate = %v, want 0.5", stats.HeadshotRate)
	}
	if stats.KillDeathRatio != 1.5 {
		t.Fatalf("KillDeathRatio = %v, want 1.5", stats.KillDeathRatio)
	}
	if stats.KillsPerMatch != 3 {
		t.Fatalf("KillsPerMatch = %v, want 3", stats.KillsPerMatch)
	}
}

func TestNewPlayerStatsAvoidsDivisionByZero(t *testing.T) {
	stats := NewPlayerStats(0, 0, 0, 0, 0)

	if stats.HeadshotRate != 0 || stats.KillDeathRatio != 0 || stats.KillsPerMatch != 0 {
		t.Fatalf("unexpected zero-denominator statistics: %+v", stats)
	}
}
