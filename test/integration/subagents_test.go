//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestSubagents_RealRun(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	agents := map[string]*grok.SubagentConfig{
		"security": {
			Description: "Security analysis",
			Prompt:      "You are a security expert. Be terse.",
			Tools:       []string{grok.ToolRead, grok.ToolGrep},
		},
	}
	res, err := c.RunPromptCtx(ctx, "Is `db.Exec(\"SELECT * FROM users WHERE id=\" + req.ID)` safe? Reply yes or no.", &grok.RunOptions{
		Format:       grok.JSONOutput,
		Agent:        "security",
		Agents:       agents,
		MaxBudgetUSD: 0.05,
	})
	if err != nil {
		t.Fatalf("subagents: %v", err)
	}
	if res.Text == "" {
		t.Fatal("empty response")
	}
}
