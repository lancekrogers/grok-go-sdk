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
	if got[0].ID == "SESSION ID" {
		t.Fatalf("parsed table header as session: %#v", got[0])
	}
	if got[0].UpdatedAt == "" {
		t.Fatalf("first session missing updated date: %#v", got[0])
	}
}

func TestParseSessions_FutureUUIDPrefix(t *testing.T) {
	in := []byte("01a12345-6789-7abc-8def-0123456789ab  2026-09-01  2026-09-01  local  future session\n")
	got := parseSessions(in)
	if len(got) != 1 {
		t.Fatalf("got %d sessions, want 1: %#v", len(got), got)
	}
	if got[0].ID != "01a12345-6789-7abc-8def-0123456789ab" {
		t.Fatalf("session id mismatch: %#v", got[0])
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
