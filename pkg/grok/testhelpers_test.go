package grok

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

var (
	mockBinOnce sync.Once
	mockBinPath string
	mockBinErr  error
)

func buildOrLocateMock(t *testing.T) string {
	t.Helper()
	mockBinOnce.Do(func() {
		_, this, _, _ := runtime.Caller(0)
		repo := filepath.Join(filepath.Dir(this), "..", "..")
		outDir := filepath.Join(repo, "test", "mockserver", "bin")
		outBin := filepath.Join(outDir, "grok-mock")

		cmd := exec.Command("go", "build", "-o", outBin, "./test/mockserver")
		cmd.Dir = repo
		if out, err := cmd.CombinedOutput(); err != nil {
			mockBinErr = err
			t.Logf("go build output:\n%s", out)
			return
		}
		mockBinPath = outBin
	})
	if mockBinErr != nil {
		t.Fatalf("build mock: %v", mockBinErr)
	}
	return mockBinPath
}
