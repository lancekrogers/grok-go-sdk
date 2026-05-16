package dangerous

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

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

func mockBinPath(t *testing.T) string {
	t.Helper()
	_, this, _, _ := runtime.Caller(0)
	repo := filepath.Join(filepath.Dir(this), "..", "..", "..")
	outDir := filepath.Join(repo, "test", "mockserver", "bin")
	outBin := filepath.Join(outDir, "grok-mock")
	cmd := exec.Command("go", "build", "-o", outBin, "./test/mockserver")
	cmd.Dir = repo
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("go build:\n%s", out)
		t.Fatalf("build mock: %v", err)
	}
	return outBin
}

func TestBypassAllPermissions_AgainstMock(t *testing.T) {
	t.Setenv("GROK_ENABLE_DANGEROUS", "i-accept-all-risks")
	t.Setenv("GO_ENV", "")
	t.Setenv("NODE_ENV", "")
	c, err := NewDangerousClient(mockBinPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.BypassAllPermissions("hi", &grok.RunOptions{Format: grok.JSONOutput}); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestWithBypassPermissionsMode_AgainstMock(t *testing.T) {
	t.Setenv("GROK_ENABLE_DANGEROUS", "i-accept-all-risks")
	t.Setenv("GO_ENV", "")
	t.Setenv("NODE_ENV", "")
	c, err := NewDangerousClient(mockBinPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.WithBypassPermissionsMode("hi", &grok.RunOptions{Format: grok.JSONOutput}); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestDisableSandbox_AgainstMock(t *testing.T) {
	t.Setenv("GROK_ENABLE_DANGEROUS", "i-accept-all-risks")
	t.Setenv("GO_ENV", "")
	t.Setenv("NODE_ENV", "")
	c, err := NewDangerousClient(mockBinPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.DisableSandbox("hi", &grok.RunOptions{Format: grok.JSONOutput}); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestBypassAllPermissions_DoesNotMutateCallerOpts(t *testing.T) {
	t.Setenv("GROK_ENABLE_DANGEROUS", "i-accept-all-risks")
	t.Setenv("GO_ENV", "")
	t.Setenv("NODE_ENV", "")
	c, _ := NewDangerousClient(mockBinPath(t))
	in := &grok.RunOptions{Format: grok.JSONOutput}
	_, _ = c.BypassAllPermissions("hi", in)
	if in.AlwaysApprove {
		t.Fatal("caller opts were mutated")
	}
	if in.AllowDangerousMode {
		t.Fatal("caller opts were mutated")
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
