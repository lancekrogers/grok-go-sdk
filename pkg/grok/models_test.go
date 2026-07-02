package grok

import (
	"context"
	"os"
	"path/filepath"
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

func TestParseModels_CapturedFixture(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "models.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("fixture not yet captured: %v", err)
	}
	got := parseModels(data)
	if len(got) == 0 {
		t.Fatal("no models parsed from fixture")
	}
	if got[0].ID == "" {
		t.Fatalf("first model missing id: %#v", got[0])
	}
}
