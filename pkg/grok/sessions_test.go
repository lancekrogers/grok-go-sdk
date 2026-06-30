package grok

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSessionsList_AgainstFixture(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "sessions-list.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip("fixture not yet captured")
	}
	got := parseSessions(data)
	if len(got) == 0 {
		t.Fatal("no sessions parsed")
	}
	if got[0].ID == "" {
		t.Fatal("first session missing id")
	}
}

func TestSessionsDelete_EmptyID(t *testing.T) {
	c := NewClient("/nonexistent")
	if err := c.SessionsDelete(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestSessionsDelete_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if err := c.SessionsDelete(context.Background(), "mock-session-01HMOCK0000000000000000000"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
