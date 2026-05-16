//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

var (
	mockBinOnce sync.Once
	mockBinPath string
	mockBinErr  error
)

func newClient(t *testing.T) *grok.GrokClient {
	t.Helper()
	if os.Getenv("INTEGRATION_REAL") == "1" {
		path, err := exec.LookPath("grok")
		if err != nil {
			t.Skip("INTEGRATION_REAL=1 but grok not on PATH")
		}
		return grok.NewClient(path)
	}
	return grok.NewClient(locateMock(t))
}

func locateMock(t *testing.T) string {
	t.Helper()
	if env := os.Getenv("GROK_MOCK_BIN"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return env
		}
	}
	mockBinOnce.Do(func() {
		_, this, _, _ := runtime.Caller(0)
		repo := filepath.Join(filepath.Dir(this), "..", "..")
		outDir := filepath.Join(repo, "test", "mockserver", "bin")
		outBin := filepath.Join(outDir, "grok-mock")
		cmd := exec.Command("go", "build", "-o", outBin, "./test/mockserver")
		cmd.Dir = repo
		if out, err := cmd.CombinedOutput(); err != nil {
			mockBinErr = err
			t.Logf("go build:\n%s", out)
			return
		}
		mockBinPath = outBin
	})
	if mockBinErr != nil {
		t.Fatalf("build mock: %v", mockBinErr)
	}
	return mockBinPath
}
