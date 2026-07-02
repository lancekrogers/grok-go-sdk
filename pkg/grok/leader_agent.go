package grok

import (
	"context"
	"io"
	"os/exec"
	"sync"
)

// LeaderAgentConfig configures StartLeaderAgent (grok agent leader) — the shared
// leader process other clients connect to.
type LeaderAgentConfig struct {
	NoExitOnDisconnect bool   // --no-exit-on-disconnect
	RelayOnDemand      bool   // --relay-on-demand
	NoAutoUpdate       bool   // --no-auto-update
	GrokWSOrigin       string // --grok-ws-origin
	GrokWSURL          string // --grok-ws-url
	Env                []string
	Stderr             io.Writer
}

// LeaderAgent is a handle to a running `grok agent leader` process.
type LeaderAgent struct {
	cmd      *exec.Cmd
	stopOnce sync.Once
	stopErr  error
	waitDone chan error
}

func leaderAgentArgs(cfg *LeaderAgentConfig) []string {
	args := []string{"agent", "leader"}
	if cfg.NoExitOnDisconnect {
		args = append(args, "--no-exit-on-disconnect")
	}
	if cfg.RelayOnDemand {
		args = append(args, "--relay-on-demand")
	}
	if cfg.NoAutoUpdate {
		args = append(args, "--no-auto-update")
	}
	if cfg.GrokWSOrigin != "" {
		args = append(args, "--grok-ws-origin", cfg.GrokWSOrigin)
	}
	if cfg.GrokWSURL != "" {
		args = append(args, "--grok-ws-url", cfg.GrokWSURL)
	}
	return args
}

// StartLeaderAgent launches `grok agent leader` and returns a handle. Call Stop
// to terminate it.
func (c *GrokClient) StartLeaderAgent(ctx context.Context, cfg *LeaderAgentConfig) (*LeaderAgent, error) {
	if cfg == nil {
		cfg = &LeaderAgentConfig{}
	}
	cmd := c.command(ctx, leaderAgentArgs(cfg)...)
	cmd.Env = c.envBase(cfg.Env)
	cmd.Stderr = cfg.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	la := &LeaderAgent{cmd: cmd, waitDone: make(chan error, 1)}
	go func() {
		la.waitDone <- cmd.Wait()
	}()
	return la, nil
}

// Stop signals the leader process to terminate, escalating to kill after a grace period.
func (l *LeaderAgent) Stop() error {
	l.stopOnce.Do(func() {
		l.stopErr = stopGracefully(l.cmd, l.waitDone)
	})
	return l.stopErr
}
