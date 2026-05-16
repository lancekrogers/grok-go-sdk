package grok

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func LocateBinary() (string, error) {
	if p, err := exec.LookPath("grok"); err == nil {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err == nil {
		candidates := []string{
			filepath.Join(home, ".grok", "bin", "grok"),
		}
		switch runtime.GOOS {
		case "darwin":
			candidates = append(candidates,
				"/opt/homebrew/bin/grok",
				"/usr/local/bin/grok",
			)
		case "linux":
			candidates = append(candidates,
				"/usr/local/bin/grok",
				"/usr/bin/grok",
			)
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
	}
	return "", &GrokError{
		Type:    ErrorValidation,
		Message: "grok binary not found on PATH or in standard locations",
	}
}
