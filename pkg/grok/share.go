package grok

import (
	"context"
	"strings"
)

func (c *GrokClient) Share(ctx context.Context, sessionID string) (string, error) {
	out, err := c.runSubcommand(ctx, []string{"share", sessionID})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
