package grok

import (
	"context"
	"testing"
)

func TestInspect_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ir, err := c.Inspect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ir.Version == "" {
		t.Fatalf("missing version: %#v", ir)
	}
}
