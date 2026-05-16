//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestCheckLoop_Validation(t *testing.T) {
	c := grok.NewClient("/nonexistent")
	_, err := c.RunPromptCtx(context.Background(), "trivial", &grok.RunOptions{
		Check:  true,
		Format: grok.JSONOutput,
	})
	if err != nil {
		if _, ok := err.(*grok.GrokError); !ok {
			t.Fatalf("expected exec-related error, got: %v", err)
		}
	}
}

func TestCheckLoop_RealRun(t *testing.T) {
	if os.Getenv("INTEGRATION_REAL") != "1" {
		t.Skip("requires real grok (check loop runs multiple passes)")
	}
	if os.Getenv("INTEGRATION_SLOW") != "1" {
		t.Skip("set INTEGRATION_SLOW=1 to run; --check makes multiple grok passes")
	}
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	res, err := c.RunPromptCtx(ctx, "Reply with exactly 'OK'.", &grok.RunOptions{
		Check:        true,
		Format:       grok.JSONOutput,
		MaxBudgetUSD: 0.10,
	})
	if err != nil {
		t.Fatalf("check run: %v", err)
	}
	if res.Text == "" {
		t.Fatal("empty result")
	}
}

