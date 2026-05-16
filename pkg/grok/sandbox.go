package grok

import (
	"os"
	"strings"
)

func ResolveSandboxProfile(field string) string {
	if field != "" {
		return field
	}
	return os.Getenv("GROK_SANDBOX")
}

func SandboxIsPermissive(profile string) bool {
	switch strings.ToLower(profile) {
	case "off", "none", "disabled":
		return true
	}
	return false
}
