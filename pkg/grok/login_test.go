package grok

import (
	"context"
	"testing"
)

func TestLogin_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	if err := c.Login(context.Background(), LoginOAuth); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
