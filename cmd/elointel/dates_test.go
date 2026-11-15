package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/analytics"
)

type rangeSource struct {
	fakeStatisticsSource
	mu         sync.Mutex
	operation  string
	start, end *time.Time
}

func (s *rangeSource) record(operation string, start, end *time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operation, s.start, s.end = operation, start, end
}

func (s *rangeSource) FindPlayerByID(_ context.Context, id int64) (*elo.Player, error) {
	return &elo.Player{ID: id, Name: "player"}, nil
}

func (s *rangeSource) GetPlayerStats(_ context.Context, _ int64) (*elo.PlayerStats, error) {
	s.record("stats", nil, nil)
	return elo.NewPlayerStats(0, 0, 0, 0, 0), nil
}

func (s *rangeSource) GetPlayerStatsForPeriod(_ context.Context, _ int64, start, end *time.Time) (*elo.PlayerStats, error) {
	s.record("stats", start, end)
	return elo.NewPlayerStats(0, 0, 0, 0, 0), nil
}

func (s *rangeSource) GetMatches(_ context.Context, _, _ int) ([]*elo.Match, error) {
	s.record("matches", nil, nil)
	return []*elo.Match{}, nil
}

func (s *rangeSource) GetMatchesForPeriod(_ context.Context, start, end *time.Time, _, _ int) ([]*elo.Match, error) {
	s.record("matches", start, end)
	return []*elo.Match{}, nil
}

func (s *rangeSource) GetPlayerResultsForPeriod(_ context.Context, start, end *time.Time) ([]*elo.PlayerResult, error) {
	s.record("results", start, end)
	return []*elo.PlayerResult{}, nil
}

func (s *rangeSource) GetHeadToHeadForPeriod(_ context.Context, _, _ int64, start, end *time.Time) (*elo.HeadToHeadStats, error) {
	s.record("head", start, end)
	return &elo.HeadToHeadStats{}, nil
}

func (s *rangeSource) GetWeaponStatsForPeriod(_ context.Context, _ int64, start, end *time.Time) ([]*elo.WeaponStats, error) {
	s.record("weapons", start, end)
	return []*elo.WeaponStats{}, nil
}

func (s *rangeSource) GetMapStatsForPeriod(_ context.Context, _ int64, start, end *time.Time) ([]*elo.MapStats, error) {
	s.record("maps", start, end)
	return []*elo.MapStats{}, nil
}

func (s *rangeSource) assertRange(t *testing.T, operation, start, end string) {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.operation != operation || formatRangeDate(s.start) != start || formatRangeEndDate(s.end) != end {
		t.Fatalf("got %s [%v, %v), want %s [%s, %s inclusive]", s.operation, s.start, s.end, operation, start, end)
	}
}

func TestOptionalDatesAcrossMCPAndGraphQL(t *testing.T) {
	source := &rangeSource{}
	handler, err := newHandler("test", analytics.New(source), "", "")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint: server.URL + "/mcp", DisableStandaloneSSE: true,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	tools := []struct {
		name, operation string
		arguments       map[string]any
	}{
		{"get_player_profile", "stats", map[string]any{"player_id": "1"}},
		{"get_player_period_stats", "stats", map[string]any{"player_id": "1"}},
		{"list_matches", "matches", map[string]any{}},
		{"get_player_results_for_period", "results", map[string]any{}},
		{"get_head_to_head", "head", map[string]any{"first_player_id": "1", "second_player_id": "2"}},
		{"get_player_weapon_stats", "weapons", map[string]any{"player_id": "1"}},
		{"get_player_map_stats", "maps", map[string]any{"player_id": "1"}},
	}
	ranges := []struct{ start, end string }{
		{"", ""}, {"2026-10-01", ""}, {"", "2026-10-03"}, {"2026-10-01", "2026-10-03"},
	}
	for _, tool := range tools {
		for _, bounds := range ranges {
			t.Run(tool.name+"/"+bounds.start+"/"+bounds.end, func(t *testing.T) {
				args := make(map[string]any)
				for key, value := range tool.arguments {
					args[key] = value
				}
				if bounds.start != "" {
					args["start_date"] = bounds.start
				}
				if bounds.end != "" {
					args["end_date"] = bounds.end
				}
				result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tool.name, Arguments: args})
				if err != nil || result.IsError {
					t.Fatalf("MCP call: result=%+v error=%v", result, err)
				}
				source.assertRange(t, tool.operation, bounds.start, bounds.end)
			})
		}
	}
	daily, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_daily_player_results", Arguments: map[string]any{"date": "2026-10-03"},
	})
	if err != nil || daily.IsError {
		t.Fatalf("daily results: result=%+v error=%v", daily, err)
	}
	source.assertRange(t, "results", "2026-10-03", "2026-10-03")
	daily, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_daily_player_results", Arguments: map[string]any{},
	})
	if err != nil || !daily.IsError {
		t.Fatalf("daily results accepted a missing date: result=%+v error=%v", daily, err)
	}

	for _, test := range []struct {
		query, operation, start, end string
	}{
		{`{ player(id: "1") { stats { kills } } }`, "stats", "", ""},
		{`{ player(id: "1") { stats(startDate: "2026-10-01") { kills } } }`, "stats", "2026-10-01", ""},
		{`{ player(id: "1") { stats(endDate: "2026-10-03") { kills } } }`, "stats", "", "2026-10-03"},
		{`{ matches { id } }`, "matches", "", ""},
		{`{ matches(startDate: "2026-10-01", endDate: "2026-10-03") { id } }`, "matches", "2026-10-01", "2026-10-03"},
	} {
		body, err := json.Marshal(map[string]string{"query": test.query})
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.Post(server.URL+"/graphql", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		var result struct {
			Errors []any `json:"errors"`
		}
		err = json.NewDecoder(response.Body).Decode(&result)
		response.Body.Close()
		if err != nil || response.StatusCode != http.StatusOK || len(result.Errors) != 0 {
			t.Fatalf("GraphQL query %s: response=%+v error=%v", test.query, result, err)
		}
		source.assertRange(t, test.operation, test.start, test.end)
	}

	for _, tool := range tools {
		args := make(map[string]any)
		for key, value := range tool.arguments {
			args[key] = value
		}
		args["start_date"], args["end_date"] = "2026-10-03", "2026-10-01"
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tool.name, Arguments: args})
		if err != nil || !result.IsError {
			t.Fatalf("%s accepted reversed range: result=%+v error=%v", tool.name, result, err)
		}
	}
	request := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"{ matches(startDate: \"bad\") { id } }"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if !strings.Contains(response.Body.String(), "start_date must be YYYY-MM-DD") {
		t.Fatalf("GraphQL accepted invalid date: %s", response.Body.String())
	}
}
