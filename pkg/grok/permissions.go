package grok

import (
	"fmt"
	"strings"
)

func ValidateAllowDenyRule(rule string) error {
	if rule == "" {
		return &GrokError{Type: ErrorValidation, Message: "permission rule is empty"}
	}
	if strings.ContainsRune(rule, 0) {
		return &GrokError{Type: ErrorValidation, Message: "permission rule contains null byte"}
	}
	return nil
}

func validateRules(rules []string, label string) error {
	for i, r := range rules {
		if err := ValidateAllowDenyRule(r); err != nil {
			return fmt.Errorf("validation: %s[%d]: %w", label, i, err)
		}
	}
	return nil
}
