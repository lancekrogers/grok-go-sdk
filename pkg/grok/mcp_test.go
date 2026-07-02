package grok

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseMCPList_Fixture(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "mcp-list.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip("fixture not yet captured")
	}
	got := parseMCPList(data)
	if bytes.HasPrefix(bytes.TrimSpace(data), []byte("No MCP servers configured")) {
		if len(got) != 0 {
			t.Fatalf("empty-state fixture should parse to 0 servers, got %d: %#v", len(got), got)
		}
		return
	}
	if len(got) == 0 {
		t.Fatal("no servers parsed from populated fixture")
	}
}

func TestParseMCPList_Inline(t *testing.T) {
	// grok 0.2.77 lists one server per line as "<name>: <command-or-url>".
	in := []byte("  alpha: /usr/bin/server arg1\n  beta: https://example.invalid/mcp\n")
	got := parseMCPList(in)
	if len(got) != 2 {
		t.Fatalf("got %d servers, want 2: %#v", len(got), got)
	}
	if got[0].Name != "alpha" || got[0].CommandOrURL != "/usr/bin/server arg1" {
		t.Fatalf("alpha mismatch: %#v", got[0])
	}
	if got[1].Name != "beta" || got[1].CommandOrURL != "https://example.invalid/mcp" {
		t.Fatalf("beta mismatch: %#v", got[1])
	}
}

func TestMCPAddArgs(t *testing.T) {
	stdio := mcpAddArgs(MCPServerConfig{
		Name:         "alpha",
		CommandOrURL: "npx",
		Args:         []string{"-y", "pkg"},
		Env:          map[string]string{"FOO": "bar"},
	})
	wantStdio := []string{"mcp", "add", "-e", "FOO=bar", "alpha", "npx", "--", "-y", "pkg"}
	if !reflect.DeepEqual(stdio, wantStdio) {
		t.Fatalf("stdio argv:\n got %#v\nwant %#v", stdio, wantStdio)
	}

	http := mcpAddArgs(MCPServerConfig{
		Name:         "beta",
		CommandOrURL: "https://mcp.example.com/mcp",
		Transport:    MCPTransportHTTP,
		Scope:        MCPScopeProject,
		Headers:      map[string]string{"Authorization": "Bearer xyz"},
	})
	wantHTTP := []string{"mcp", "add", "-t", "http", "-s", "project", "-H", "Authorization: Bearer xyz", "beta", "https://mcp.example.com/mcp"}
	if !reflect.DeepEqual(http, wantHTTP) {
		t.Fatalf("http argv:\n got %#v\nwant %#v", http, wantHTTP)
	}
}

func TestParseMCPList_EmptyState(t *testing.T) {
	in := []byte("No MCP servers configured. Run `grok mcp add --help` to get started.\n")
	got := parseMCPList(in)
	if len(got) != 0 {
		t.Fatalf("got %d servers, want 0: %#v", len(got), got)
	}
}

func TestMCPRemove_EmptyName(t *testing.T) {
	c := NewClient("/nonexistent")
	if err := c.MCPRemove(context.Background(), "", ""); err == nil {
		t.Fatal("expected error")
	} else {
		requireValidationError(t, err)
	}
}

func TestMCPAdd_EmptyName(t *testing.T) {
	c := NewClient("/nonexistent")
	if err := c.MCPAdd(context.Background(), MCPServerConfig{}); err == nil {
		t.Fatal("expected error")
	} else {
		requireValidationError(t, err)
	}
}
