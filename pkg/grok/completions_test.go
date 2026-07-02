package grok

import (
	"context"
	"reflect"
	"testing"
)

func TestCompletions_Valid(t *testing.T) {
	got := captureArgs(t)
	c := NewClient("grok")
	if _, err := c.Completions(context.Background(), "zsh"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"completions", "zsh"}
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("got %v want %v", *got, want)
	}
}

func TestCompletions_InvalidShell(t *testing.T) {
	c := NewClient("grok")
	_, err := c.Completions(context.Background(), "tcsh")
	if err == nil {
		t.Fatal("expected validation error for unsupported shell")
	}
	if ge, ok := err.(*GrokError); !ok || ge.Type != ErrorValidation {
		t.Fatalf("want validation *GrokError, got %T %v", err, err)
	}
}

func TestCompletions_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	out, err := c.Completions(context.Background(), "bash")
	if err != nil || len(out) == 0 {
		t.Fatalf("completions bash: out=%q err=%v", out, err)
	}
}
