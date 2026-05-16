//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func TestStreamPrompt_TextAndEnd(t *testing.T) {
	c := newClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	events, errs := c.StreamPrompt(ctx, "say hello", &grok.RunOptions{
		Format:       grok.StreamingJSONOutput,
		MaxBudgetUSD: 0.05,
	})
	var assembled strings.Builder
	var endEvent *grok.Event
	for {
		select {
		case ev, ok := <-events:
			if !ok {
				goto done
			}
			if ev.Type == grok.EventText {
				assembled.WriteString(ev.Content())
			}
			if ev.Type == grok.EventEnd {
				copy := ev
				endEvent = &copy
			}
		case err, ok := <-errs:
			if !ok {
				continue
			}
			if err != nil {
				t.Fatalf("stream error: %v", err)
			}
		case <-ctx.Done():
			t.Fatalf("timeout: %v", ctx.Err())
		}
	}
done:
	if assembled.Len() == 0 {
		t.Fatal("no text events received")
	}
	if endEvent == nil {
		t.Fatal("no end event received")
	}
	if endEvent.SessionID == "" {
		t.Fatalf("end event missing sessionId: %#v", *endEvent)
	}
	if endEvent.StopReason == "" {
		t.Fatalf("end event missing stopReason: %#v", *endEvent)
	}
}
