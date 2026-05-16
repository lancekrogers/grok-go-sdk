package dangerous

import "testing"

func TestNewDangerousClient_RequiresEnv(t *testing.T) {
	t.Setenv("GROK_ENABLE_DANGEROUS", "")
	if _, err := NewDangerousClient("/bin/echo"); err != ErrNotEnabled {
		t.Fatalf("got %v want %v", err, ErrNotEnabled)
	}
}

func TestNewDangerousClient_RefusesProduction(t *testing.T) {
	t.Setenv("GROK_ENABLE_DANGEROUS", "i-accept-all-risks")
	t.Setenv("GO_ENV", "production")
	if _, err := NewDangerousClient("/bin/echo"); err != ErrProduction {
		t.Fatalf("got %v want %v", err, ErrProduction)
	}
}

func TestNewDangerousClient_OK(t *testing.T) {
	t.Setenv("GROK_ENABLE_DANGEROUS", "i-accept-all-risks")
	t.Setenv("GO_ENV", "")
	t.Setenv("NODE_ENV", "")
	if _, err := NewDangerousClient("/bin/echo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
