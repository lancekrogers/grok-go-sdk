package grok

import (
	"context"
	"testing"
)

func TestMemoryClear_ScopeRequired(t *testing.T) {
	c := NewClient("/nonexistent")
	if err := c.MemoryClear(context.Background(), 0); err == nil {
		t.Fatal("expected error for missing scope")
	}
}

func TestMemoryClear_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if err := c.MemoryClear(context.Background(), MemoryAll); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.MemoryClear(context.Background(), MemoryWorkspace); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
