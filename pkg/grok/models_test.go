package grok

import (
	"context"
	"testing"
)

func TestModels_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	mm, err := c.Models(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mm) == 0 {
		t.Fatal("expected at least one model")
	}
	hasDefault := false
	for _, m := range mm {
		if m.Default {
			hasDefault = true
		}
	}
	if !hasDefault {
		t.Fatal("expected default-marked model")
	}
}

func TestParseModels_Captured(t *testing.T) {
	in := []byte("Default model: grok-build\n* grok-build (default)\n- grok-fast\n- grok-think\n")
	got := parseModels(in)
	if len(got) != 3 {
		t.Fatalf("got %d, want 3: %#v", len(got), got)
	}
	if got[0].ID != "grok-build" || !got[0].Default {
		t.Fatalf("row 0: %#v", got[0])
	}
}
