package grok

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunHeadlessAgent_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	var out bytes.Buffer
	in := strings.NewReader(`{"type":"user","text":"hi"}` + "\n")
	err := c.RunHeadlessAgent(context.Background(), &HeadlessAgentConfig{
		Stdin:  in,
		Stdout: &out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
