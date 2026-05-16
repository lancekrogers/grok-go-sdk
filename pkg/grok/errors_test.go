package grok

import (
	"errors"
	"testing"
	"time"
)

func TestErrorType_Constants(t *testing.T) {
	cases := []struct {
		got, want string
	}{
		{string(ErrorAuth), "auth"},
		{string(ErrorRateLimit), "rate_limit"},
		{string(ErrorProcess), "process"},
		{string(ErrorValidation), "validation"},
		{string(ErrorTransport), "transport"},
		{string(ErrorUnknown), "unknown"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got %q want %q", c.got, c.want)
		}
	}
}

func TestGrokError_Error(t *testing.T) {
	e := &GrokError{Type: ErrorAuth, Message: "bad token"}
	want := "auth: bad token"
	if e.Error() != want {
		t.Errorf("got %q want %q", e.Error(), want)
	}

	e2 := &GrokError{Type: ErrorUnknown}
	if e2.Error() != "unknown" {
		t.Errorf("got %q", e2.Error())
	}
}

func TestGrokError_Unwrap(t *testing.T) {
	inner := errors.New("inner")
	e := &GrokError{Type: ErrorProcess, Original: inner}
	if !errors.Is(e, inner) {
		t.Fatal("Is failed")
	}
}

func TestGrokError_IsRetryable(t *testing.T) {
	for _, tt := range []struct {
		typ  ErrorType
		want bool
	}{
		{ErrorRateLimit, true},
		{ErrorTransport, true},
		{ErrorProcess, true},
		{ErrorAuth, false},
		{ErrorValidation, false},
		{ErrorUnknown, false},
	} {
		e := &GrokError{Type: tt.typ}
		if e.IsRetryable() != tt.want {
			t.Errorf("Type=%s got=%v want=%v", tt.typ, e.IsRetryable(), tt.want)
		}
	}
}

func TestParseError_Auth(t *testing.T) {
	e := ParseError("Error: unauthorized request", 1)
	if e.Type != ErrorAuth {
		t.Errorf("got %s want auth", e.Type)
	}
	if e.IsRetryable() {
		t.Error("auth should not be retryable")
	}
}

func TestParseError_RateLimit_WithRetryAfter(t *testing.T) {
	e := ParseError("rate limited; retry after 5 seconds", 429)
	if e.Type != ErrorRateLimit {
		t.Errorf("got %s want rate_limit", e.Type)
	}
	if !e.IsRetryable() {
		t.Error("should be retryable")
	}
	if e.RetryDelay() != 5*time.Second {
		t.Errorf("RetryDelay got %v want 5s", e.RetryDelay())
	}
}

func TestParseError_Validation(t *testing.T) {
	e := ParseError("invalid flag --foo", 2)
	if e.Type != ErrorValidation {
		t.Errorf("got %s want validation", e.Type)
	}
	if e.IsRetryable() {
		t.Error("validation should not be retryable")
	}
}

func TestParseError_Transport(t *testing.T) {
	e := ParseError("transport channel closed", 1)
	if e.Type != ErrorTransport {
		t.Errorf("got %s want transport", e.Type)
	}
}

func TestParseError_Process_NonZeroExit(t *testing.T) {
	e := ParseError("something went wrong", 1)
	if e.Type != ErrorProcess {
		t.Errorf("got %s want process", e.Type)
	}
}

func TestParseError_Unknown(t *testing.T) {
	e := ParseError("", 0)
	if e.Type != ErrorUnknown {
		t.Errorf("got %s want unknown", e.Type)
	}
}
