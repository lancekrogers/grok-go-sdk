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

func TestStdioAgent_FullTurn(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	s, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer s.Close()

	init, err := s.Initialize(ctx, "grok-go-sdk-integration", "0.1.0")
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	if init.ProtocolVersion == 0 {
		t.Fatalf("initialize missing protocolVersion: %#v", init)
	}
	if len(init.AuthMethods) == 0 {
		t.Fatal("initialize missing authMethods")
	}

	if _, err := s.Authenticate(ctx, "cached_token"); err != nil {
		t.Fatalf("authenticate: %v", err)
	}

	cwd, _ := os.Getwd()
	sessID, err := s.NewSession(ctx, cwd, nil)
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	if sessID == "" {
		t.Fatal("empty session id")
	}

	var assembled strings.Builder
	done := make(chan struct{})
	go func() {
		defer close(done)
		for u := range s.Updates() {
			if u.Update.SessionUpdate == grok.UpdateAgentMessageChunk {
				assembled.WriteString(u.Update.ContentText())
			}
		}
	}()

	res, err := s.PromptText(ctx, sessID, "Reply with exactly the string OK and nothing else.")
	if err != nil {
		t.Fatalf("prompt: %v", err)
	}
	if res.StopReason == "" {
		t.Fatalf("missing stopReason: %#v", res)
	}
	if res.Meta == nil {
		t.Fatalf("missing usage meta: %#v", res)
	}
	if res.Meta.TotalTokens == 0 {
		t.Fatalf("expected non-zero totalTokens: %#v", res.Meta)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	<-done

	if assembled.Len() == 0 {
		t.Fatalf("no agent_message_chunk content received")
	}
}
