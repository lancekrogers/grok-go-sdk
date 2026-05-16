package grok

import (
	"testing"
)

func TestResolveSandboxProfile_Env(t *testing.T) {
	t.Setenv("GROK_SANDBOX", "foo")
	if got := ResolveSandboxProfile(""); got != "foo" {
		t.Fatalf("got %q want foo", got)
	}
	if got := ResolveSandboxProfile("bar"); got != "bar" {
		t.Fatalf("explicit got %q want bar", got)
	}
}

func TestSandboxIsPermissive(t *testing.T) {
	for _, name := range []string{"off", "Off", "disabled", "none"} {
		if !SandboxIsPermissive(name) {
			t.Fatalf("%q should be permissive", name)
		}
	}
	for _, name := range []string{"strict", "default", ""} {
		if SandboxIsPermissive(name) {
			t.Fatalf("%q should not be permissive", name)
		}
	}
}
