package grok

import (
	"strings"
	"testing"
)

func TestRunFromStdin_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	transcript := strings.NewReader(`{"role":"user","content":"summarize"}` + "\n")
	res, err := c.RunFromStdin(transcript, "", &RunOptions{
		Format: JSONOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Text == "" {
		t.Fatal("empty result")
	}
}
