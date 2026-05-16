package grok

import "testing"

func TestValidateAllowDenyRule(t *testing.T) {
	if err := ValidateAllowDenyRule(""); err == nil {
		t.Fatal("expected error on empty")
	}
	if err := ValidateAllowDenyRule("Read(*)"); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if err := ValidateAllowDenyRule("Read(\x00)"); err == nil {
		t.Fatal("expected null-byte rejection")
	}
}

func TestBuildToolGlob(t *testing.T) {
	if got := BuildToolGlob("Bash", "git status*"); got != "Bash(git status*)" {
		t.Fatalf("got %q", got)
	}
	if got := BuildToolGlob("Read", ""); got != "Read" {
		t.Fatalf("got %q", got)
	}
}
