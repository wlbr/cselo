package elo

import (
	"strings"
	"testing"

	"github.com/go-gcfg/gcfg"
)

func TestConfigReadsGraphQLAndMCPSections(t *testing.T) {
	config := new(Config)
	if err := gcfg.ReadStringInto(config, `
[GraphQL]
Port = 8080
[MCP]
Port = 42821
`); err != nil {
		t.Fatalf("ReadStringInto returned error: %v", err)
	}

	if config.GraphQL.Port != "8080" {
		t.Errorf("GraphQL.Port = %q, want %q", config.GraphQL.Port, "8080")
	}
	if config.MCP.Port != "42821" {
		t.Errorf("MCP.Port = %q, want %q", config.MCP.Port, "42821")
	}
}

func TestConfigReadsServiceSectionsFromINI(t *testing.T) {
	config := new(Config)
	err := gcfg.ReadStringInto(config, strings.Join([]string{
		"[GraphQL]",
		"Port = 8080",
		"[MCP]",
		"Transport = http",
		"Port = 42821",
		"Path = /mcp",
	}, "\n"))
	if err != nil {
		t.Fatalf("ReadStringInto returned error: %v", err)
	}
	if config.MCP.Transport != "http" || config.MCP.Path != "/mcp" {
		t.Fatalf("unexpected MCP config: %+v", config.MCP)
	}
}
