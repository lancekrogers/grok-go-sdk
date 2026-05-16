package grok

import (
	"context"
	"testing"
)

func TestStartServeAgent_PortBinding(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	sa, err := c.StartServeAgent(context.Background(), &ServeAgentConfig{Bind: "127.0.0.1:0"})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if sa.Addr() == "" {
		t.Fatal("expected non-empty Addr after start")
	}
	if err := sa.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
