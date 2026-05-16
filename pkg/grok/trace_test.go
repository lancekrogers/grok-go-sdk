package grok

import (
	"context"
	"testing"
)

func TestTrace_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	path, uploaded, err := c.Trace(context.Background(), "mock-session-01HMOCK0000000000000000000", &TraceOptions{LocalOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Fatal("empty path")
	}
	if uploaded {
		t.Fatal("expected uploaded=false from mock")
	}
}
