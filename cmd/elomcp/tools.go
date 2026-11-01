package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/wlbr/cselo/elo"
)

type statisticsSource interface {
	SearchPlayers(context.Context, string, int, int) ([]*elo.Player, error)
	FindPlayerByID(context.Context, int64) (*elo.Player, error)
	GetPlayerStats(context.Context, int64) (*elo.PlayerStats, error)
	GetMatches(context.Context, int, int) ([]*elo.Match, error)
	GetMatchesForPeriod(context.Context, time.Time, time.Time, int, int) ([]*elo.Match, error)
	GetMatchByID(context.Context, int64) (*elo.Match, error)
	GetMatchParticipants(context.Context, int64) ([]*elo.MatchParticipant, error)
	GetPlayerStatsForPeriod(context.Context, int64, time.Time, time.Time) (*elo.PlayerStats, error)
	GetPlayerResultsForPeriod(context.Context, time.Time, time.Time) ([]*elo.PlayerResult, error)
	GetHeadToHeadForPeriod(context.Context, int64, int64, time.Time, time.Time) (*elo.HeadToHeadStats, error)
	GetWeaponStatsForPeriod(context.Context, int64, time.Time, time.Time) ([]*elo.WeaponStats, error)
	GetMapStatsForPeriod(context.Context, int64, time.Time, time.Time) ([]*elo.MapStats, error)
}

type searchPlayersInput struct {
	Name   string `json:"name" jsonschema:"Player name or a distinctive fragment of it"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Maximum number of results from 1 to 100; defaults to 20"`
	Offset int    `json:"offset,omitempty" jsonschema:"Zero-based result offset; defaults to 0"`
}

type playerSummary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	SteamID string `json:"steamid"`
}

type searchPlayersOutput struct {
	Players []playerSummary `json:"players"`
}

type playerProfileInput struct {
	PlayerID  string `json:"player_id" jsonschema:"Numeric cselo player ID"`
	StartDate string `json:"start_date,omitempty" jsonschema:"Optional first included UTC date in YYYY-MM-DD format; must be supplied with end_date"`
	EndDate   string `json:"end_date,omitempty" jsonschema:"Optional last included UTC date in YYYY-MM-DD format; must be supplied with start_date"`
}

type playerStats struct {
	Matches        int     `json:"matches"`
	Kills          int     `json:"kills"`
	Deaths         int     `json:"deaths"`
	Assists        int     `json:"assists"`
	Headshots      int     `json:"headshots"`
	HeadshotRate   float64 `json:"headshotRate"`
	KillDeathRatio float64 `json:"killDeathRatio"`
	KillsPerMatch  float64 `json:"killsPerMatch"`
}

type playerProfileOutput struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Player    *struct {
		playerSummary
		Stats playerStats `json:"stats"`
	} `json:"player"`
}

type listMatchesInput struct {
	StartDate string `json:"start_date,omitempty" jsonschema:"Optional first included UTC date in YYYY-MM-DD format; must be supplied with end_date"`
	EndDate   string `json:"end_date,omitempty" jsonschema:"Optional last included UTC date in YYYY-MM-DD format; must be supplied with start_date"`
	Limit     int    `json:"limit,omitempty" jsonschema:"Maximum number of matches from 1 to 100; defaults to 20"`
	Offset    int    `json:"offset,omitempty" jsonschema:"Zero-based result offset; defaults to 0"`
}

type matchSummary struct {
	ID              string `json:"id"`
	MapName         string `json:"mapname"`
	ScoreA          int    `json:"scorea"`
	ScoreB          int    `json:"scoreb"`
	DurationSeconds int    `json:"durationSeconds"`
	Completed       bool   `json:"completed"`
}

type listMatchesOutput struct {
	StartDate string         `json:"start_date,omitempty"`
	EndDate   string         `json:"end_date,omitempty"`
	Matches   []matchSummary `json:"matches"`
}

type matchReportInput struct {
	MatchID string `json:"match_id" jsonschema:"Numeric cselo match ID"`
}

type matchReportOutput struct {
	Match *struct {
		matchSummary
		Participants []struct {
			Player playerSummary `json:"player"`
			Stats  playerStats   `json:"stats"`
		} `json:"participants"`
	} `json:"match"`
}

type playerPeriodInput struct {
	PlayerID  string `json:"player_id" jsonschema:"Numeric cselo player ID"`
	StartDate string `json:"start_date" jsonschema:"First included UTC date in YYYY-MM-DD format"`
	EndDate   string `json:"end_date" jsonschema:"Last included UTC date in YYYY-MM-DD format"`
}

type playerPeriodOutput struct {
	Player playerSummary `json:"player"`
	Stats  playerStats   `json:"stats"`
}

type playerResultsInput struct {
	StartDate string `json:"start_date" jsonschema:"First included UTC date in YYYY-MM-DD format"`
	EndDate   string `json:"end_date" jsonschema:"Last included UTC date in YYYY-MM-DD format"`
}

type dailyPlayerResultsInput struct {
	Date string `json:"date" jsonschema:"UTC day in YYYY-MM-DD format"`
}

type playerResult struct {
	Player playerSummary `json:"player"`
	Stats  playerStats   `json:"stats"`
}

type playerResultsOutput struct {
	StartDate string         `json:"start_date"`
	EndDate   string         `json:"end_date"`
	Results   []playerResult `json:"results"`
}

type headToHeadInput struct {
	FirstPlayerID  string `json:"first_player_id" jsonschema:"Numeric cselo player ID for the first player"`
	SecondPlayerID string `json:"second_player_id" jsonschema:"Numeric cselo player ID for the second player"`
	StartDate      string `json:"start_date" jsonschema:"First included UTC date in YYYY-MM-DD format"`
	EndDate        string `json:"end_date" jsonschema:"Last included UTC date in YYYY-MM-DD format"`
}

type headToHeadOutput struct {
	StartDate       string        `json:"start_date"`
	EndDate         string        `json:"end_date"`
	FirstPlayer     playerSummary `json:"first_player"`
	SecondPlayer    playerSummary `json:"second_player"`
	FirstKills      int           `json:"first_kills"`
	FirstHeadshots  int           `json:"first_headshots"`
	SecondKills     int           `json:"second_kills"`
	SecondHeadshots int           `json:"second_headshots"`
}

type weaponStatsOutput struct {
	Player    playerSummary `json:"player"`
	StartDate string        `json:"start_date"`
	EndDate   string        `json:"end_date"`
	Weapons   []weaponStat  `json:"weapons"`
}

type weaponStat struct {
	Weapon    string `json:"weapon"`
	Kills     int    `json:"kills"`
	Headshots int    `json:"headshots"`
}

type mapStatsOutput struct {
	Player    playerSummary `json:"player"`
	StartDate string        `json:"start_date"`
	EndDate   string        `json:"end_date"`
	Maps      []mapStat     `json:"maps"`
}

type mapStat struct {
	MapName   string `json:"map_name"`
	Matches   int    `json:"matches"`
	Kills     int    `json:"kills"`
	Deaths    int    `json:"deaths"`
	Assists   int    `json:"assists"`
	Headshots int    `json:"headshots"`
}

func registerTools(server *mcp.Server, source statisticsSource) {
	tools := &statsTools{source: source}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_players",
		Description: "Find Counter-Strike players by name before retrieving player statistics.",
	}, tools.searchPlayers)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_player_profile",
		Description: "Get a player's all-time statistics or statistics for an optional inclusive UTC date range, including kills, deaths, assists, headshots, K/D, and matches played.",
	}, tools.getPlayerProfile)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_matches",
		Description: "List recent Counter-Strike matches or matches in an optional inclusive UTC date range, with map, score, duration, and completion status.",
	}, tools.listMatches)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_match_report",
		Description: "Get a match report including its score, duration, and per-player statistics.",
	}, tools.getMatchReport)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_player_period_stats",
		Description: "Get one player's statistics for an inclusive date range. Use search_players first when the player ID is unknown.",
	}, tools.getPlayerPeriodStats)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_daily_player_results",
		Description: "Get the player-results table for one UTC day, ranked by kills. It includes matches, kills, deaths, assists, headshots, K/D, and kills per match.",
	}, tools.getDailyPlayerResults)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_player_results_for_period",
		Description: "Get the player-results table for an inclusive UTC date range, ranked by kills.",
	}, tools.getPlayerResultsForPeriod)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_head_to_head",
		Description: "Compare two players over an inclusive UTC date range: how many times the first killed the second and vice versa, including headshots. Use search_players first when IDs are unknown.",
	}, tools.getHeadToHead)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_player_weapon_stats",
		Description: "Get a player's weapon kill and headshot breakdown for an inclusive UTC date range.",
	}, tools.getPlayerWeaponStats)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_player_map_stats",
		Description: "Get a player's match count, kills, deaths, assists, and headshots by map for an inclusive UTC date range.",
	}, tools.getPlayerMapStats)
}

type statsTools struct {
	source statisticsSource
}

func (t *statsTools) searchPlayers(ctx context.Context, _ *mcp.CallToolRequest, input searchPlayersInput) (*mcp.CallToolResult, searchPlayersOutput, error) {
	if input.Name == "" {
		return nil, searchPlayersOutput{}, fmt.Errorf("name is required")
	}
	limit, err := validateLimit(input.Limit)
	if err != nil {
		return nil, searchPlayersOutput{}, err
	}
	if err := validateOffset(input.Offset); err != nil {
		return nil, searchPlayersOutput{}, err
	}
	players, err := t.source.SearchPlayers(ctx, input.Name, limit, input.Offset)
	if err != nil {
		return nil, searchPlayersOutput{}, err
	}
	result := make([]playerSummary, 0, len(players))
	for _, player := range players {
		result = append(result, newPlayerSummary(player))
	}
	return nil, searchPlayersOutput{Players: result}, nil
}

func (t *statsTools) getPlayerProfile(ctx context.Context, _ *mcp.CallToolRequest, input playerProfileInput) (*mcp.CallToolResult, playerProfileOutput, error) {
	if input.PlayerID == "" {
		return nil, playerProfileOutput{}, fmt.Errorf("player_id is required")
	}
	id, err := parseID(input.PlayerID, "player_id")
	if err != nil {
		return nil, playerProfileOutput{}, err
	}
	player, err := t.source.FindPlayerByID(ctx, id)
	if err != nil {
		return nil, playerProfileOutput{}, err
	}
	if player == nil {
		return nil, playerProfileOutput{}, fmt.Errorf("player %d was not found", id)
	}
	start, end, hasRange, err := parseOptionalDateRange(input.StartDate, input.EndDate)
	if err != nil {
		return nil, playerProfileOutput{}, err
	}
	var stats *elo.PlayerStats
	if hasRange {
		stats, err = t.source.GetPlayerStatsForPeriod(ctx, id, start, end)
	} else {
		stats, err = t.source.GetPlayerStats(ctx, id)
	}
	if err != nil {
		return nil, playerProfileOutput{}, err
	}
	return nil, playerProfileOutput{Player: &struct {
		playerSummary
		Stats playerStats `json:"stats"`
	}{playerSummary: newPlayerSummary(player), Stats: newPlayerStats(stats)},
		StartDate: formatRangeDate(start, hasRange),
		EndDate:   formatRangeEndDate(end, hasRange),
	}, nil
}

func (t *statsTools) listMatches(ctx context.Context, _ *mcp.CallToolRequest, input listMatchesInput) (*mcp.CallToolResult, listMatchesOutput, error) {
	limit, err := validateLimit(input.Limit)
	if err != nil {
		return nil, listMatchesOutput{}, err
	}
	if err := validateOffset(input.Offset); err != nil {
		return nil, listMatchesOutput{}, err
	}
	start, end, hasRange, err := parseOptionalDateRange(input.StartDate, input.EndDate)
	if err != nil {
		return nil, listMatchesOutput{}, err
	}
	var matches []*elo.Match
	if hasRange {
		matches, err = t.source.GetMatchesForPeriod(ctx, start, end, limit, input.Offset)
	} else {
		matches, err = t.source.GetMatches(ctx, limit, input.Offset)
	}
	if err != nil {
		return nil, listMatchesOutput{}, err
	}
	result := make([]matchSummary, 0, len(matches))
	for _, match := range matches {
		result = append(result, newMatchSummary(match))
	}
	return nil, listMatchesOutput{
		StartDate: formatRangeDate(start, hasRange),
		EndDate:   formatRangeEndDate(end, hasRange),
		Matches:   result,
	}, nil
}

func (t *statsTools) getMatchReport(ctx context.Context, _ *mcp.CallToolRequest, input matchReportInput) (*mcp.CallToolResult, matchReportOutput, error) {
	if input.MatchID == "" {
		return nil, matchReportOutput{}, fmt.Errorf("match_id is required")
	}
	id, err := parseID(input.MatchID, "match_id")
	if err != nil {
		return nil, matchReportOutput{}, err
	}
	match, err := t.source.GetMatchByID(ctx, id)
	if err != nil {
		return nil, matchReportOutput{}, err
	}
	if match == nil {
		return nil, matchReportOutput{}, fmt.Errorf("match %d was not found", id)
	}
	participants, err := t.source.GetMatchParticipants(ctx, id)
	if err != nil {
		return nil, matchReportOutput{}, err
	}
	result := matchReportOutput{Match: &struct {
		matchSummary
		Participants []struct {
			Player playerSummary `json:"player"`
			Stats  playerStats   `json:"stats"`
		} `json:"participants"`
	}{matchSummary: newMatchSummary(match)}}
	for _, participant := range participants {
		result.Match.Participants = append(result.Match.Participants, struct {
			Player playerSummary `json:"player"`
			Stats  playerStats   `json:"stats"`
		}{Player: newPlayerSummary(participant.Player), Stats: newPlayerStats(participant.Stats)})
	}
	return nil, result, nil
}

func (t *statsTools) getPlayerPeriodStats(ctx context.Context, _ *mcp.CallToolRequest, input playerPeriodInput) (*mcp.CallToolResult, playerPeriodOutput, error) {
	player, id, start, end, err := t.playerAndPeriod(ctx, input.PlayerID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, playerPeriodOutput{}, err
	}
	stats, err := t.source.GetPlayerStatsForPeriod(ctx, id, start, end)
	if err != nil {
		return nil, playerPeriodOutput{}, err
	}
	return nil, playerPeriodOutput{Player: newPlayerSummary(player), Stats: newPlayerStats(stats)}, nil
}

func (t *statsTools) getDailyPlayerResults(ctx context.Context, _ *mcp.CallToolRequest, input dailyPlayerResultsInput) (*mcp.CallToolResult, playerResultsOutput, error) {
	start, end, err := parseDateRange(input.Date, input.Date)
	if err != nil {
		return nil, playerResultsOutput{}, err
	}
	return t.playerResults(ctx, start, end)
}

func (t *statsTools) getPlayerResultsForPeriod(ctx context.Context, _ *mcp.CallToolRequest, input playerResultsInput) (*mcp.CallToolResult, playerResultsOutput, error) {
	start, end, err := parseDateRange(input.StartDate, input.EndDate)
	if err != nil {
		return nil, playerResultsOutput{}, err
	}
	return t.playerResults(ctx, start, end)
}

func (t *statsTools) getHeadToHead(ctx context.Context, _ *mcp.CallToolRequest, input headToHeadInput) (*mcp.CallToolResult, headToHeadOutput, error) {
	firstID, err := parseID(input.FirstPlayerID, "first_player_id")
	if err != nil {
		return nil, headToHeadOutput{}, err
	}
	secondID, err := parseID(input.SecondPlayerID, "second_player_id")
	if err != nil {
		return nil, headToHeadOutput{}, err
	}
	first, err := t.source.FindPlayerByID(ctx, firstID)
	if err != nil {
		return nil, headToHeadOutput{}, err
	}
	if first == nil {
		return nil, headToHeadOutput{}, fmt.Errorf("player %d was not found", firstID)
	}
	second, err := t.source.FindPlayerByID(ctx, secondID)
	if err != nil {
		return nil, headToHeadOutput{}, err
	}
	if second == nil {
		return nil, headToHeadOutput{}, fmt.Errorf("player %d was not found", secondID)
	}
	start, end, err := parseDateRange(input.StartDate, input.EndDate)
	if err != nil {
		return nil, headToHeadOutput{}, err
	}
	stats, err := t.source.GetHeadToHeadForPeriod(ctx, firstID, secondID, start, end)
	if err != nil {
		return nil, headToHeadOutput{}, err
	}
	return nil, headToHeadOutput{
		StartDate:       start.Format(time.DateOnly),
		EndDate:         end.AddDate(0, 0, -1).Format(time.DateOnly),
		FirstPlayer:     newPlayerSummary(first),
		SecondPlayer:    newPlayerSummary(second),
		FirstKills:      stats.FirstKills,
		FirstHeadshots:  stats.FirstHeadshots,
		SecondKills:     stats.SecondKills,
		SecondHeadshots: stats.SecondHeadshots,
	}, nil
}

func (t *statsTools) getPlayerWeaponStats(ctx context.Context, _ *mcp.CallToolRequest, input playerPeriodInput) (*mcp.CallToolResult, weaponStatsOutput, error) {
	player, id, start, end, err := t.playerAndPeriod(ctx, input.PlayerID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, weaponStatsOutput{}, err
	}
	stats, err := t.source.GetWeaponStatsForPeriod(ctx, id, start, end)
	if err != nil {
		return nil, weaponStatsOutput{}, err
	}
	result := weaponStatsOutput{
		Player:    newPlayerSummary(player),
		StartDate: start.Format(time.DateOnly),
		EndDate:   end.AddDate(0, 0, -1).Format(time.DateOnly),
		Weapons:   make([]weaponStat, 0, len(stats)),
	}
	for _, stat := range stats {
		result.Weapons = append(result.Weapons, weaponStat{
			Weapon: stat.Weapon, Kills: stat.Kills, Headshots: stat.Headshots,
		})
	}
	return nil, result, nil
}

func (t *statsTools) getPlayerMapStats(ctx context.Context, _ *mcp.CallToolRequest, input playerPeriodInput) (*mcp.CallToolResult, mapStatsOutput, error) {
	player, id, start, end, err := t.playerAndPeriod(ctx, input.PlayerID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, mapStatsOutput{}, err
	}
	stats, err := t.source.GetMapStatsForPeriod(ctx, id, start, end)
	if err != nil {
		return nil, mapStatsOutput{}, err
	}
	result := mapStatsOutput{
		Player:    newPlayerSummary(player),
		StartDate: start.Format(time.DateOnly),
		EndDate:   end.AddDate(0, 0, -1).Format(time.DateOnly),
		Maps:      make([]mapStat, 0, len(stats)),
	}
	for _, stat := range stats {
		result.Maps = append(result.Maps, mapStat{
			MapName: stat.MapName, Matches: stat.Matches, Kills: stat.Kills,
			Deaths: stat.Deaths, Assists: stat.Assists, Headshots: stat.Headshots,
		})
	}
	return nil, result, nil
}

func (t *statsTools) playerResults(ctx context.Context, start, end time.Time) (*mcp.CallToolResult, playerResultsOutput, error) {
	results, err := t.source.GetPlayerResultsForPeriod(ctx, start, end)
	if err != nil {
		return nil, playerResultsOutput{}, err
	}
	output := playerResultsOutput{
		StartDate: start.Format(time.DateOnly),
		EndDate:   end.AddDate(0, 0, -1).Format(time.DateOnly),
		Results:   make([]playerResult, 0, len(results)),
	}
	for _, result := range results {
		output.Results = append(output.Results, playerResult{
			Player: newPlayerSummary(result.Player),
			Stats:  newPlayerStats(result.Stats),
		})
	}
	return nil, output, nil
}

func (t *statsTools) playerAndPeriod(ctx context.Context, playerID, startDate, endDate string) (*elo.Player, int64, time.Time, time.Time, error) {
	id, err := parseID(playerID, "player_id")
	if err != nil {
		return nil, 0, time.Time{}, time.Time{}, err
	}
	player, err := t.source.FindPlayerByID(ctx, id)
	if err != nil {
		return nil, 0, time.Time{}, time.Time{}, err
	}
	if player == nil {
		return nil, 0, time.Time{}, time.Time{}, fmt.Errorf("player %d was not found", id)
	}
	start, end, err := parseDateRange(startDate, endDate)
	if err != nil {
		return nil, 0, time.Time{}, time.Time{}, err
	}
	return player, id, start, end, nil
}

func parseID(value, name string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return id, nil
}

func parseDateRange(startDate, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("start_date must be YYYY-MM-DD: %w", err)
	}
	end, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must be YYYY-MM-DD: %w", err)
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must not be before start_date")
	}
	return start, end.AddDate(0, 0, 1), nil
}

func parseOptionalDateRange(startDate, endDate string) (time.Time, time.Time, bool, error) {
	if startDate == "" && endDate == "" {
		return time.Time{}, time.Time{}, false, nil
	}
	start, end, err := parseDateRange(startDate, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, false, err
	}
	return start, end, true, nil
}

func formatRangeDate(start time.Time, hasRange bool) string {
	if !hasRange {
		return ""
	}
	return start.Format(time.DateOnly)
}

func formatRangeEndDate(end time.Time, hasRange bool) string {
	if !hasRange {
		return ""
	}
	return end.AddDate(0, 0, -1).Format(time.DateOnly)
}

func newPlayerSummary(player *elo.Player) playerSummary {
	return playerSummary{
		ID:      strconv.FormatInt(player.ID, 10),
		Name:    player.Name,
		SteamID: player.SteamID,
	}
}

func newPlayerStats(stats *elo.PlayerStats) playerStats {
	return playerStats{
		Matches:        stats.Matches,
		Kills:          stats.Kills,
		Deaths:         stats.Deaths,
		Assists:        stats.Assists,
		Headshots:      stats.Headshots,
		HeadshotRate:   stats.HeadshotRate,
		KillDeathRatio: stats.KillDeathRatio,
		KillsPerMatch:  stats.KillsPerMatch,
	}
}

func newMatchSummary(match *elo.Match) matchSummary {
	return matchSummary{
		ID:              strconv.FormatInt(match.ID(), 10),
		MapName:         match.MapName,
		ScoreA:          match.ScoreA,
		ScoreB:          match.ScoreB,
		DurationSeconds: int(match.Duration.Seconds()),
		Completed:       match.Completed,
	}
}
