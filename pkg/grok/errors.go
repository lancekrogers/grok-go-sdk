package grok

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ErrorType string

const (
	ErrorAuth       ErrorType = "auth"
	ErrorRateLimit  ErrorType = "rate_limit"
	ErrorProcess    ErrorType = "process"
	ErrorValidation ErrorType = "validation"
	ErrorTransport  ErrorType = "transport"
	ErrorUnknown    ErrorType = "unknown"
)

type GrokError struct {
	Type       ErrorType
	Message    string
	Stderr     string
	ExitCode   int
	Original   error
	retryAfter time.Duration
}

func (e *GrokError) Error() string {
	if e.Message != "" {
		return string(e.Type) + ": " + e.Message
	}
	return string(e.Type)
}

func (e *GrokError) Unwrap() error { return e.Original }

func (e *GrokError) IsRetryable() bool {
	switch e.Type {
	case ErrorRateLimit, ErrorTransport, ErrorProcess:
		return true
	}
	return false
}

func (e *GrokError) RetryDelay() time.Duration { return e.retryAfter }

var retryAfterRe = regexp.MustCompile(`retry[- ]?after[^\d]*(\d+)`)

func ParseError(stderr string, exitCode int) *GrokError {
	s := strings.ToLower(stderr)
	e := &GrokError{Stderr: stderr, ExitCode: exitCode}
	switch {
	case strings.Contains(s, "auth") || strings.Contains(s, "unauthorized") || strings.Contains(s, "authorizationrequired"):
		e.Type = ErrorAuth
		e.Message = "authentication required"
	case strings.Contains(s, "rate") && strings.Contains(s, "limit"):
		e.Type = ErrorRateLimit
		e.Message = "rate limited"
		if m := retryAfterRe.FindStringSubmatch(s); len(m) == 2 {
			if n, err := strconv.Atoi(m[1]); err == nil {
				e.retryAfter = time.Duration(n) * time.Second
			}
		}
	case strings.Contains(s, "invalid") || strings.Contains(s, "validation"):
		e.Type = ErrorValidation
		e.Message = "invalid request"
	case strings.Contains(s, "transport"):
		e.Type = ErrorTransport
		e.Message = "transport error"
	case exitCode != 0:
		e.Type = ErrorProcess
		e.Message = "non-zero exit"
	default:
		e.Type = ErrorUnknown
	}
	return e
}

var ErrNotImplemented = errors.New("not implemented in this grok version")
