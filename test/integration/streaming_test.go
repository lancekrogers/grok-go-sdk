//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestStreamPrompt_AssistantOrResult(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	events, errs := c.StreamPrompt(ctx, "say hello", &grok.RunOptions{
		Format:       grok.StreamingJSONOutput,
		MaxBudgetUSD: 0.05,
	})
	sawDone := false
	for !sawDone {
		select {
		case ev, ok := <-events:
			if !ok {
				sawDone = true
				break
			}
			if ev.Type == grok.EventAssistant || ev.Type == grok.EventResult || ev.Type == "text" || ev.Type == "end" {
				sawDone = true
			}
		case err := <-errs:
			if err != nil {
				t.Fatalf("stream error: %v", err)
			}
		case <-ctx.Done():
			t.Fatalf("timeout: %v", ctx.Err())
		}
	}
	if !sawDone {
		t.Fatal("no terminating event seen")
	}
}
