package grok

import (
	"testing"
	"time"
)

func TestGenerateSessionID_TimeOrdered(t *testing.T) {
	a := GenerateSessionID()
	time.Sleep(2 * time.Millisecond)
	b := GenerateSessionID()
	if a >= b {
		t.Fatalf("expected time-ordered, got %q >= %q", a, b)
	}
}

func TestContinueConversation_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	res, err := c.ContinueConversation("follow up")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("nil result")
	}
}
