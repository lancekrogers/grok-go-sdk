package grok

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
	"sync"
	"time"
)

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

type ServeAgentConfig struct {
	Bind         string
	Secret       string
	Remote       string
	GrokWSOrigin string
	GrokWSURL    string
}

type ServeAgent struct {
	cmd      *exec.Cmd
	addr     string
	stopOnce sync.Once
	stopErr  error
	waitDone chan error
}

var listeningRe = regexp.MustCompile(`listening on ([0-9A-Fa-f:.\[\]]+:[0-9]+)`)

func serveArgs(cfg *ServeAgentConfig) []string {
	args := []string{"agent", "serve"}
	if cfg.Bind != "" {
		args = append(args, "--bind", cfg.Bind)
	}
	if cfg.Secret != "" {
		args = append(args, "--secret", cfg.Secret)
	}
	if cfg.Remote != "" {
		args = append(args, "--remote", cfg.Remote)
	}
	if cfg.GrokWSOrigin != "" {
		args = append(args, "--grok-ws-origin", cfg.GrokWSOrigin)
	}
	if cfg.GrokWSURL != "" {
		args = append(args, "--grok-ws-url", cfg.GrokWSURL)
	}
	return args
}

func (c *GrokClient) StartServeAgent(ctx context.Context, cfg *ServeAgentConfig) (*ServeAgent, error) {
	if cfg == nil {
		cfg = &ServeAgentConfig{}
	}
	cmd := c.command(ctx, serveArgs(cfg)...)
	stderr := &syncBuffer{}
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	sa := &ServeAgent{cmd: cmd, waitDone: make(chan error, 1)}
	go func() {
		sa.waitDone <- cmd.Wait()
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if m := listeningRe.FindStringSubmatch(stderr.String()); len(m) == 2 {
			sa.addr = m[1]
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if sa.addr == "" && cfg.Bind != "" {
		sa.addr = cfg.Bind
	}
	return sa, nil
}

func (s *ServeAgent) Addr() string { return s.addr }

func (s *ServeAgent) Stop() error {
	s.stopOnce.Do(func() {
		s.stopErr = stopGracefully(s.cmd, s.waitDone)
	})
	return s.stopErr
}
