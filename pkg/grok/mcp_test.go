package grok

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestParseMCPList_Fixture(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "mcp-list.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip("fixture not yet captured")
	}
	got := parseMCPList(data)
	if len(got) == 0 {
		t.Fatal("no servers parsed")
	}
}

func TestParseMCPList_Inline(t *testing.T) {
	in := []byte("alpha\n  command: /usr/bin/server\n  transport: stdio\n\nbeta\n  url: https://example.invalid\n  transport: http\n")
	got := parseMCPList(in)
	if len(got) != 2 {
		t.Fatalf("got %d servers, want 2: %#v", len(got), got)
	}
	if got[0].Name != "alpha" || got[0].Command != "/usr/bin/server" || got[0].Transport != MCPTransportStdio {
		t.Fatalf("alpha mismatch: %#v", got[0])
	}
	if got[1].Name != "beta" || got[1].URL != "https://example.invalid" || got[1].Transport != MCPTransportHTTP {
		t.Fatalf("beta mismatch: %#v", got[1])
	}
}

func TestMCPRemove_EmptyName(t *testing.T) {
	c := NewClient("/nonexistent")
	if err := c.MCPRemove(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	}
}

func TestMCPAdd_EmptyName(t *testing.T) {
	c := NewClient("/nonexistent")
	if err := c.MCPAdd(context.Background(), MCPServerConfig{}); err == nil {
		t.Fatal("expected error")
	}
}
