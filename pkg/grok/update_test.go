package grok

import (
	"context"
	"testing"
)

func TestUpdateCheck_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ui, err := c.UpdateCheck(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ui.Current == "" {
		t.Fatalf("missing current: %#v", ui)
	}
}
