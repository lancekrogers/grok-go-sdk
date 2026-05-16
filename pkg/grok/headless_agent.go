package grok

import (
	"context"
	"io"
	"os/exec"
)

type HeadlessAgentConfig struct {
	GrokWSOrigin string
	GrokWSURL    string
	Env          []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
}

func (c *GrokClient) RunHeadlessAgent(ctx context.Context, cfg *HeadlessAgentConfig) error {
	if cfg == nil {
		cfg = &HeadlessAgentConfig{}
	}
	args := []string{"agent", "headless"}
	if cfg.GrokWSOrigin != "" {
		args = append(args, "--grok-ws-origin", cfg.GrokWSOrigin)
	}
	if cfg.GrokWSURL != "" {
		args = append(args, "--grok-ws-url", cfg.GrokWSURL)
	}
	cmd := execCommand(ctx, c.BinPath, args...)
	cmd.Env = c.envBase(cfg.Env)
	cmd.Stdin = cfg.Stdin
	cmd.Stdout = cfg.Stdout
	cmd.Stderr = cfg.Stderr
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
