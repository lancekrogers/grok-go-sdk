//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestBasicRunPrompt(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := c.RunPromptCtx(ctx, "say hello", &grok.RunOptions{
		Format:       grok.JSONOutput,
		MaxBudgetUSD: 0.05,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("nil result")
	}
	if res.Text == "" {
		t.Fatal("empty text")
	}
	if res.SessionID == "" {
		t.Fatal("missing session id")
	}
}
