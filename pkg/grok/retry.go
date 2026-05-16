package grok

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type RetryPolicy struct {
	MaxRetries    int
	BaseDelay     time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
	Jitter        bool
}

func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:    3,
		BaseDelay:     100 * time.Millisecond,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	}
}

func (p *RetryPolicy) backoffDelay(attempt int) time.Duration {
	d := float64(p.BaseDelay)
	for i := 0; i < attempt-1; i++ {
		d *= p.BackoffFactor
	}
	if time.Duration(d) > p.MaxDelay {
		d = float64(p.MaxDelay)
	}
	if p.Jitter {
		d = d * (0.5 + rand.Float64()*0.5)
	}
	return time.Duration(d)
}

func (c *GrokClient) RunPromptWithRetry(prompt string, opts *RunOptions, policy *RetryPolicy) (*GrokResult, error) {
	return c.RunPromptWithRetryCtx(context.Background(), prompt, opts, policy)
}

func (c *GrokClient) RunPromptWithRetryCtx(ctx context.Context, prompt string, opts *RunOptions, policy *RetryPolicy) (*GrokResult, error) {
	if policy == nil {
		policy = DefaultRetryPolicy()
	}
	var lastErr error
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := policy.backoffDelay(attempt)
			if ge, ok := lastErr.(*GrokError); ok && ge.Type == ErrorRateLimit && ge.RetryDelay() > 0 {
				delay = ge.RetryDelay()
			}
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		result, err := c.RunPromptCtx(ctx, prompt, opts)
		if err == nil {
			return result, nil
		}
		lastErr = err
		ge, ok := err.(*GrokError)
		if !ok || !ge.IsRetryable() {
			return nil, err
		}
	}
	return nil, fmt.Errorf("max retries (%d) exceeded: %w", policy.MaxRetries, lastErr)
}
