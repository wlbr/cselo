package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/analytics"
)

type trackingSource struct {
	fakeStatisticsSource
	searches atomic.Int32
	err      error
}

func (s *trackingSource) SearchPlayers(_ context.Context, _ string, _, _ int) ([]*elo.Player, error) {
	s.searches.Add(1)
	return s.players, s.err
}

func TestBothProtocolsUseAnalytics(t *testing.T) {
	source := &trackingSource{fakeStatisticsSource: fakeStatisticsSource{
		players: []*elo.Player{{ID: 2, Name: "Jagger", SteamID: "STEAM_1:0:1"}},
	}}
	handler, err := newHandler("test", analytics.New(source), "", "")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	body := `{"query":"{ players(name: \"Jagger\") { id name steamid } }"}`
	response, err := http.Post(server.URL+"/graphql", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var graphql struct {
		Data struct {
			Players []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"players"`
		} `json:"data"`
		Errors []any `json:"errors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&graphql); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || len(graphql.Errors) != 0 ||
		len(graphql.Data.Players) != 1 || graphql.Data.Players[0].Name != "Jagger" {
		t.Fatalf("unexpected GraphQL response: status=%d body=%+v", response.StatusCode, graphql)
	}

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint: server.URL + "/mcp", DisableStandaloneSSE: true,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 10 {
		t.Fatalf("got %d MCP tools, want 10", len(tools.Tools))
	}
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "search_players", Arguments: map[string]any{"name": "Jagger"},
	})
	if err != nil || result.IsError {
		t.Fatalf("MCP search failed: result=%+v error=%v", result, err)
	}
	content, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(content, []byte(`"name":"Jagger"`)) {
		t.Fatalf("unexpected MCP result: %s", content)
	}
	if source.searches.Load() != 2 {
		t.Fatalf("search source called %d times, want 2", source.searches.Load())
	}
}

func TestHandlerAuthenticationAndRoutes(t *testing.T) {
	service := analytics.New(fakeStatisticsSource{})
	if _, err := newHandler("test", service, "user", ""); err == nil {
		t.Fatal("incomplete credentials should be rejected")
	}
	handler, err := newHandler("test", service, "user", "password")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		path string
		want int
	}{
		{"/mcp", http.StatusUnauthorized},
		{"/graphql", http.StatusOK},
		{"/grapql", http.StatusNotFound},
		{"/healthz", http.StatusNotFound},
		{"/query", http.StatusNotFound},
	} {
		t.Run(test.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			if test.path == "/graphql" {
				request = httptest.NewRequest(http.MethodPost, test.path, strings.NewReader(`{"query":"{ __typename }"}`))
				request.Header.Set("Content-Type", "application/json")
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.want {
				t.Fatalf("status = %d, want %d", response.Code, test.want)
			}
		})
	}
}

func TestGraphQLPropagatesAnalyticsErrors(t *testing.T) {
	source := &trackingSource{err: errors.New("statistics source unavailable")}
	handler, err := newHandler("test", analytics.New(source), "", "")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"{ players(name: \"test\") { name } }"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if !strings.Contains(response.Body.String(), "statistics source unavailable") {
		t.Fatalf("expected source error, got %s", response.Body.String())
	}
}

func TestAnalyticsRejectsInvalidSearchBeforeQueryingSource(t *testing.T) {
	source := &trackingSource{}
	service := analytics.New(source)
	for _, test := range []struct {
		name   string
		limit  int
		offset int
	}{
		{"", 20, 0},
		{"test", 101, 0},
		{"test", 20, -1},
	} {
		if _, err := service.SearchPlayers(context.Background(), test.name, test.limit, test.offset); err == nil {
			t.Errorf("SearchPlayers(%q, %d, %d) accepted invalid input", test.name, test.limit, test.offset)
		}
	}
	if source.searches.Load() != 0 {
		t.Fatal("invalid search reached the source")
	}
}

type profileSource struct {
	fakeStatisticsSource
	statsCalls atomic.Int32
}

func (s *profileSource) FindPlayerByID(_ context.Context, id int64) (*elo.Player, error) {
	if id != 2 {
		return nil, nil
	}
	return &elo.Player{ID: id, Name: "Jagger", SteamID: "STEAM_1:0:1"}, nil
}

func (s *profileSource) GetPlayerStats(_ context.Context, id int64) (*elo.PlayerStats, error) {
	s.statsCalls.Add(1)
	return elo.NewPlayerStats(3, 12, 6, 2, 4), nil
}

func TestPlayerProfileUsesSameStatistics(t *testing.T) {
	source := &profileSource{}
	handler, err := newHandler("test", analytics.New(source), "", "")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	response, err := http.Post(server.URL+"/graphql", "application/json",
		strings.NewReader(`{"query":"{ player(id: \"2\") { name stats { matches kills deaths } } }"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var graphql struct {
		Data struct {
			Player struct {
				Stats struct {
					Matches int `json:"matches"`
					Kills   int `json:"kills"`
					Deaths  int `json:"deaths"`
				} `json:"stats"`
			} `json:"player"`
		} `json:"data"`
		Errors []any `json:"errors"`
	}
	if err := json.NewDecoder(response.Body).Decode(&graphql); err != nil {
		t.Fatal(err)
	}
	if len(graphql.Errors) != 0 || graphql.Data.Player.Stats.Kills != 12 {
		t.Fatalf("unexpected GraphQL response: %+v", graphql)
	}

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint: server.URL + "/mcp", DisableStandaloneSSE: true,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_player_profile", Arguments: map[string]any{"player_id": "2"},
	})
	if err != nil || result.IsError {
		t.Fatalf("MCP profile failed: result=%+v error=%v", result, err)
	}
	content, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var profile struct {
		Player struct {
			Stats struct {
				Matches int `json:"matches"`
				Kills   int `json:"kills"`
				Deaths  int `json:"deaths"`
			} `json:"stats"`
		} `json:"player"`
	}
	if err := json.Unmarshal(content, &profile); err != nil {
		t.Fatal(err)
	}
	if profile.Player.Stats != graphql.Data.Player.Stats || source.statsCalls.Load() != 2 {
		t.Fatalf("MCP stats %+v, GraphQL stats %+v, source calls %d", profile.Player.Stats, graphql.Data.Player.Stats, source.statsCalls.Load())
	}
}
