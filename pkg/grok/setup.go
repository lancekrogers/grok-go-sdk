package grok

import "context"

func (c *GrokClient) Setup(ctx context.Context) error {
	_, err := c.runSubcommand(ctx, []string{"setup"})
	return err
}
