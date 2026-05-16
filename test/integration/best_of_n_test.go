//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestBestOfN_Validation(t *testing.T) {
	c := grok.NewClient("/nonexistent")
	_, err := c.RunPromptCtx(context.Background(), "trivial", &grok.RunOptions{
		BestOfN: 4,
		Format:  grok.JSONOutput,
	})
	if err != nil {
		if _, ok := err.(*grok.GrokError); !ok {
			t.Fatalf("expected exec error not validation error, got: %v", err)
		}
	}
}

func TestBestOfN_RealRun(t *testing.T) {
	if os.Getenv("INTEGRATION_REAL") != "1" {
		t.Skip("requires real grok")
	}
	if os.Getenv("INTEGRATION_SLOW") != "1" {
		t.Skip("set INTEGRATION_SLOW=1 to run; best-of-n spawns N parallel grok calls")
	}
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	res, err := c.RunPromptCtx(ctx, "Reply with exactly one word: hi", &grok.RunOptions{
		BestOfN:      2,
		Format:       grok.JSONOutput,
		MaxBudgetUSD: 0.10,
	})
	if err != nil {
		t.Fatalf("best-of-n: %v", err)
	}
	if res.Text == "" {
		t.Fatal("empty result")
	}
}
