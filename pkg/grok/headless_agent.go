package grok

import (
	"context"
	"io"
	"os/exec"
)

type HeadlessAgentConfig struct {
	// Shared `grok agent` options (emitted before the `headless` subcommand).
	Reauth           bool   // --reauth
	Model            string // -m/--model
	ReasoningEffort  string // --reasoning-effort
	AlwaysApprove    bool   // --always-approve
	AgentProfile     string // --agent-profile
	Leader           bool   // --leader: connect to a shared leader process
	NoLeader         bool   // --no-leader: force a new agent
	CLIChatProxyBase string // --cli-chat-proxy-base-url
	XAIAPIBase       string // --xai-api-base-url

	// `agent headless` options.
	GrokWSOrigin string
	GrokWSURL    string

	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// headlessArgs builds the full argv for `grok agent [shared] headless [opts]`.
func headlessArgs(cfg *HeadlessAgentConfig) []string {
	args := []string{"agent"}
	if cfg.Reauth {
		args = append(args, "--reauth")
	}
	if cfg.Model != "" {
		args = append(args, "--model", cfg.Model)
	}
	if cfg.ReasoningEffort != "" {
		args = append(args, "--reasoning-effort", cfg.ReasoningEffort)
	}
	if cfg.AlwaysApprove {
		args = append(args, "--always-approve")
	}
	if cfg.AgentProfile != "" {
		args = append(args, "--agent-profile", cfg.AgentProfile)
	}
	if cfg.Leader {
		args = append(args, "--leader")
	}
	if cfg.NoLeader {
		args = append(args, "--no-leader")
	}
	if cfg.CLIChatProxyBase != "" {
		args = append(args, "--cli-chat-proxy-base-url", cfg.CLIChatProxyBase)
	}
	if cfg.XAIAPIBase != "" {
		args = append(args, "--xai-api-base-url", cfg.XAIAPIBase)
	}
	args = append(args, "headless")
	if cfg.GrokWSOrigin != "" {
		args = append(args, "--grok-ws-origin", cfg.GrokWSOrigin)
	}
	if cfg.GrokWSURL != "" {
		args = append(args, "--grok-ws-url", cfg.GrokWSURL)
	}
	return args
}

func (c *GrokClient) RunHeadlessAgent(ctx context.Context, cfg *HeadlessAgentConfig) error {
	if cfg == nil {
		cfg = &HeadlessAgentConfig{}
	}
	cmd := execCommand(ctx, c.BinPath, headlessArgs(cfg)...)
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
