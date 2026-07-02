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
	for _, name := range []string{SandboxOff, "Off", SandboxNone} {
		if !SandboxIsPermissive(name) {
			t.Fatalf("%q should be permissive", name)
		}
	}
	for _, name := range []string{SandboxStrict, SandboxWorkspace, SandboxReadOnly, "disabled", "default", ""} {
		if SandboxIsPermissive(name) {
			t.Fatalf("%q should not be permissive", name)
		}
	}
}

func TestSandboxProfiles_VerifiedNames(t *testing.T) {
	got := []string{SandboxOff, SandboxNone, SandboxWorkspace, SandboxReadOnly, SandboxStrict}
	want := []string{"off", "none", "workspace", "read-only", "strict"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("profile %d got %q want %q", i, got[i], want[i])
		}
	}
}
