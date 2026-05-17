package grok

import (
	"context"
	"errors"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"
)

// mockExitScript builds an execCommand that runs /bin/sh -c so the resulting
// exec failure is an *exec.ExitError, then RunPromptCtx parses stderr into a GrokError.
func mockExitScript(t *testing.T, stderr string, exitCode int) (cleanup func(), count *atomic.Int32) {
	t.Helper()
	prev := execCommand
	var c atomic.Int32
	execCommand = func(ctx context.Context, _ string, _ ...string) *exec.Cmd {
		c.Add(1)
		// shell command emits stderr and exits with code
		shCmd := "printf '%s' " + shQuote(stderr) + " 1>&2; exit " + itoa(exitCode)
		return exec.CommandContext(ctx, "/bin/sh", "-c", shCmd)
	}
	return func() { execCommand = prev }, &c
}

func shQuote(s string) string {
	return "'" + s + "'"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func TestRetry_NonRetryable_StopsAfterOne(t *testing.T) {
	cleanup, count := mockExitScript(t, "Error: authentication failed", 1)
	defer cleanup()
	c := NewClient("/nonexistent")
	_, err := c.RunPromptWithRetryCtx(context.Background(), "hi", nil, &RetryPolicy{MaxRetries: 3, BaseDelay: time.Millisecond})
	if err == nil {
		t.Fatal("expected error")
	}
	if got := count.Load(); got != 1 {
		t.Fatalf("attempts=%d want 1", got)
	}
}

func TestRetry_ProcessError_RetriesUpToMax(t *testing.T) {
	cleanup, count := mockExitScript(t, "Error: process crashed", 1)
	defer cleanup()
	c := NewClient("/nonexistent")
	_, err := c.RunPromptWithRetryCtx(context.Background(), "hi", nil, &RetryPolicy{MaxRetries: 2, BaseDelay: time.Millisecond})
	if err == nil {
		t.Fatal("expected error")
	}
	if got := count.Load(); got != 3 {
		t.Fatalf("attempts=%d want 3", got)
	}
}

func TestRetryPolicy_ZeroMaxDelayDoesNotClamp(t *testing.T) {
	p := &RetryPolicy{
		BaseDelay:     100 * time.Millisecond,
		BackoffFactor: 2,
		Jitter:        false,
	}
	if got, want := p.backoffDelay(1), 100*time.Millisecond; got != want {
		t.Fatalf("attempt 1 delay=%s want %s", got, want)
	}
	if got, want := p.backoffDelay(2), 200*time.Millisecond; got != want {
		t.Fatalf("attempt 2 delay=%s want %s", got, want)
	}
}

func TestRetry_RespectsContextCancel(t *testing.T) {
	cleanup, _ := mockExitScript(t, "Error: process crashed", 1)
	defer cleanup()
	c := NewClient("/nonexistent")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	_, err := c.RunPromptWithRetryCtx(ctx, "hi", nil, &RetryPolicy{MaxRetries: 5, BaseDelay: 100 * time.Millisecond})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v want context.Canceled", err)
	}
}
