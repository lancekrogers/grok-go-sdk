package grok

import (
	"os"
	"strings"
)

const (
	SandboxOff       = "off"
	SandboxNone      = "none"
	SandboxWorkspace = "workspace"
	SandboxReadOnly  = "read-only"
	SandboxStrict    = "strict"
)

func ResolveSandboxProfile(field string) string {
	if field != "" {
		return field
	}
	return os.Getenv("GROK_SANDBOX")
}

func SandboxIsPermissive(profile string) bool {
	switch strings.ToLower(profile) {
	case SandboxOff, SandboxNone:
		return true
	}
	return false
}
