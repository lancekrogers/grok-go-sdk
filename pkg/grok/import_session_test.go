package grok

import (
	"context"
	"testing"
)

func TestImport_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if _, err := c.Import(context.Background(), []string{"./fake.jsonl"}, &ImportOptions{JSON: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
