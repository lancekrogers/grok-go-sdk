package grok

import (
	"context"
	"testing"
	"time"
)

func TestWorktreeList_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	got, err := c.WorktreeList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one worktree from mock")
	}
	if got[0].ID == "" || got[0].Path == "" {
		t.Fatalf("missing fields: %#v", got[0])
	}
}

func TestWorktreeShow_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	w, err := c.WorktreeShow(context.Background(), "mock-wt-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.ID == "" {
		t.Fatal("empty id")
	}
}

func TestWorktreeRemove_NoIDs(t *testing.T) {
	c := NewClient("/nonexistent")
	if err := c.WorktreeRemove(context.Background(), nil, false, false); err == nil {
		t.Fatal("expected error")
	}
}

func TestWorktreeGC_Argv(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if err := c.WorktreeGC(context.Background(), 24*time.Hour, true, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorktreeDBPath_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	p, err := c.WorktreeDBPath(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == "" {
		t.Fatal("empty path")
	}
}

func TestParseWorktreeList_ShortRows(t *testing.T) {
	in := []byte("only-id\n# comment\nid1\tpath1\n\nid2\tpath2\tbranch2\n")
	got := parseWorktreeList(in)
	if len(got) != 3 {
		t.Fatalf("got %d, want 3: %#v", len(got), got)
	}
	if got[0].ID != "only-id" {
		t.Fatalf("row 0 id: %#v", got[0])
	}
	if got[2].Branch != "branch2" {
		t.Fatalf("row 2 branch: %#v", got[2])
	}
}
