//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestMCPList(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := c.MCPList(ctx); err != nil {
		t.Fatalf("mcp list: %v", err)
	}
}

func TestMCPDoctor(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := c.MCPDoctor(ctx); err != nil {
		t.Fatalf("mcp doctor: %v", err)
	}
}
