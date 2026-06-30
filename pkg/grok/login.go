package grok

import (
	"context"
	"os"
	"os/exec"
)

type LoginMode int

const (
	LoginInteractive LoginMode = iota
	LoginOAuth
	LoginDevice
)

func (c *GrokClient) Login(ctx context.Context, mode LoginMode) error {
	args := []string{"login"}
	switch mode {
	case LoginOAuth:
		args = append(args, "--oauth")
	case LoginDevice:
		args = append(args, "--device-auth")
	}
	cmd := execCommand(ctx, c.BinPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			ge := ParseError("", ee.ExitCode())
			ge.Original = err
			return ge
		}
		return err
	}
	return nil
}

// Logout signs out and clears cached credentials (grok logout).
func (c *GrokClient) Logout(ctx context.Context) error {
	_, err := c.runSubcommand(ctx, []string{"logout"})
	return err
}
