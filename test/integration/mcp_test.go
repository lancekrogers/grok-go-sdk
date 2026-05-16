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
	report, err := c.MCPDoctor(ctx)
	// doctor reports problems by writing a useful stdout report and exiting non-zero.
	// MCPDoctor returns both so callers can present the report alongside the error.
	if report == "" && err != nil {
		t.Fatalf("mcp doctor: empty report with error: %v", err)
	}
}
