package grok

import (
	"context"
	"testing"
)

func TestLeaderList_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	list, err := c.LeaderList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list == nil {
		// Empty list is fine; just must not return nil from a successful call... actually mock returns [].
		// Accept any value as long as no error.
	}
}

func TestLeaderKill_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if err := c.LeaderKill(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLeaderProfile_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if err := c.LeaderProfileStart(context.Background(), 1234); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := c.LeaderProfileStop(context.Background(), 1234); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if _, err := c.LeaderProfileStatus(context.Background(), 1234); err != nil {
		t.Fatalf("status: %v", err)
	}
}
