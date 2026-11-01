package elo

type PlayerStats struct {
	Matches        int
	Kills          int
	Deaths         int
	Assists        int
	Headshots      int
	HeadshotRate   float64
	KillDeathRatio float64
	KillsPerMatch  float64
}

func NewPlayerStats(matches, kills, deaths, assists, headshots int) *PlayerStats {
	stats := &PlayerStats{
		Matches:   matches,
		Kills:     kills,
		Deaths:    deaths,
		Assists:   assists,
		Headshots: headshots,
	}

	if kills > 0 {
		stats.HeadshotRate = float64(headshots) / float64(kills)
	}
	if deaths > 0 {
		stats.KillDeathRatio = float64(kills) / float64(deaths)
	}
	if matches > 0 {
		stats.KillsPerMatch = float64(kills) / float64(matches)
	}

	return stats
}

type MatchParticipant struct {
	Player *Player
	Stats  *PlayerStats
}
