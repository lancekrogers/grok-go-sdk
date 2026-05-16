package grok

import (
	"context"
	"testing"
)

func TestSetup_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if err := c.Setup(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
