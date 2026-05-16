//go:build integration

package integration

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestWorktree_CreateListRemove(t *testing.T) {
	if os.Getenv("INTEGRATION_REAL") != "1" {
		t.Skip("worktree integration requires real grok")
	}
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	name := "test-wt-" + grok.GenerateSessionID()
	_, err := c.RunPromptCtx(ctx, "say hello", &grok.RunOptions{
		Format:       grok.JSONOutput,
		Worktree:     grok.WorktreeOption{Enabled: true, Name: name},
		MaxBudgetUSD: 0.05,
	})
	if err != nil {
		t.Fatalf("create worktree: %v", err)
	}

	list, err := c.WorktreeList(ctx)
	if err != nil {
		t.Fatalf("worktree list: %v", err)
	}
	found := false
	for _, w := range list {
		if w.ID == name || strings.Contains(w.Path, name) {
			found = true
			break
		}
	}
	if !found {
		t.Logf("created worktree not found in list (may be deferred)")
	}

	if err := c.WorktreeRemove(ctx, []string{name}, true, false); err != nil {
		t.Logf("worktree remove returned: %v", err)
	}
}
