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

func TestSessionsContinueAndResume(t *testing.T) {
	if os.Getenv("INTEGRATION_REAL") != "1" {
		t.Skip("requires real grok session memory")
	}
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	first, err := c.RunPromptCtx(ctx, "remember the number 42", &grok.RunOptions{
		Format:       grok.JSONOutput,
		MaxBudgetUSD: 0.05,
	})
	if err != nil {
		t.Fatalf("first prompt: %v", err)
	}
	if first.SessionID == "" {
		t.Fatal("no session id")
	}

	second, err := c.ResumeConversationCtx(ctx, "what number did I just tell you?", first.SessionID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if !strings.Contains(second.Text, "42") {
		t.Fatalf("resumed session did not recall: %q", second.Text)
	}
}

func TestSessionsList(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := c.SessionsList(ctx, 5)
	if err != nil {
		t.Fatalf("sessions list: %v", err)
	}
}
