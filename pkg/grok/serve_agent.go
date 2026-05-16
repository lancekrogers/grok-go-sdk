package grok

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
	"sync"
	"syscall"
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
	waitDone chan struct{}
}

var listeningRe = regexp.MustCompile(`listening on ([0-9A-Fa-f:.\[\]]+:[0-9]+)`)

func (c *GrokClient) StartServeAgent(ctx context.Context, cfg *ServeAgentConfig) (*ServeAgent, error) {
	if cfg == nil {
		cfg = &ServeAgentConfig{}
	}
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

	cmd := execCommand(ctx, c.BinPath, args...)
	stderr := &syncBuffer{}
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	sa := &ServeAgent{cmd: cmd, waitDone: make(chan struct{})}
	go func() {
		_ = cmd.Wait()
		close(sa.waitDone)
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
		if s.cmd.Process != nil {
			_ = s.cmd.Process.Signal(syscall.SIGTERM)
		}
		select {
		case <-s.waitDone:
		case <-time.After(5 * time.Second):
			if s.cmd.Process != nil {
				_ = s.cmd.Process.Kill()
			}
			<-s.waitDone
		}
	})
	return s.stopErr
}
