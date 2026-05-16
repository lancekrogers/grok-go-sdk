package grok

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLocateBinary_PathHit(t *testing.T) {
	dir := t.TempDir()
	fakeBin := filepath.Join(dir, "grok")
	if err := os.WriteFile(fakeBin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	got, err := LocateBinary()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != fakeBin {
		t.Fatalf("got %q want %q", got, fakeBin)
	}
}

func TestLocateBinary_NotFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	_, err := LocateBinary()
	if err == nil {
		t.Fatal("expected error")
	}
	var ge *GrokError
	if !errors.As(err, &ge) || ge.Type != ErrorValidation {
		t.Fatalf("want ErrorValidation, got %v", err)
	}
}
