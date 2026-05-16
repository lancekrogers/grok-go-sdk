package grok

import (
	"context"
	"strings"
	"testing"
)

func TestShare_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	url, err := c.Share(context.Background(), "mock-session-01HMOCK0000000000000000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(url, "https://") {
		t.Fatalf("expected URL, got %q", url)
	}
}
