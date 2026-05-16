//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestStdioAgent_RoundTrip(t *testing.T) {
	if os.Getenv("INTEGRATION_REAL") == "1" {
		t.Skip("real grok agent stdio uses JSON-RPC 2.0; SDK envelope is incompatible (see design doc 10 Q8)")
	}
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	s, err := c.StartStdioAgent(ctx, &grok.StdioConfig{})
	if err != nil {
		t.Fatalf("start stdio: %v", err)
	}
	defer s.Close()

	resp, err := s.RequestResponse(ctx, grok.AgentMessage{
		Type: "user",
		ID:   "req-1",
		Text: "ping",
	})
	if err != nil {
		t.Fatalf("request response: %v", err)
	}
	if resp.ID != "req-1" {
		t.Fatalf("response id mismatch: %q", resp.ID)
	}
}
