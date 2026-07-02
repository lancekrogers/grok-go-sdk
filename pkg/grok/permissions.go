package grok

import (
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
			ge := validationErrorf("%s[%d]: %s", label, i, err.Error())
			ge.Original = err
			return ge
		}
	}
	return nil
}
