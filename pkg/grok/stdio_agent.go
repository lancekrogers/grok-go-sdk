package grok

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type AgentMessage struct {
	Type    string          `json:"type"`
	ID      string          `json:"id,omitempty"`
	Text    string          `json:"text,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Raw     json.RawMessage `json:"-"`
}

type StdioConfig struct {
	Model               string
	ReasoningEffort     string
	AlwaysApprove       bool
	AgentProfile        string
	UseLeader           bool
	NoLeader            bool
	GrokWSOrigin        string
	GrokWSURL           string
	CLIChatProxyBaseURL string
	XAIAPIBaseURL       string
	Env                 []string
}

type StdioSession struct {
	cmd           *exec.Cmd
	stdin         io.WriteCloser
	stdout        io.ReadCloser
	writeMu       sync.Mutex
	messages      chan AgentMessage
	errs          chan error
	closed        chan struct{}
	receiveCalled chan struct{}
	closeOnce     sync.Once
	shutdownOnce  sync.Once
	pending       *pendingTable
}

type pendingTable struct {
	mu      sync.Mutex
	waiters map[string]chan AgentMessage
}

func newPendingTable() *pendingTable {
	return &pendingTable{waiters: make(map[string]chan AgentMessage)}
}

func (p *pendingTable) register(id string) chan AgentMessage {
	ch := make(chan AgentMessage, 1)
	p.mu.Lock()
	p.waiters[id] = ch
	p.mu.Unlock()
	return ch
}

func (p *pendingTable) deliver(m AgentMessage) bool {
	if m.ID == "" {
		return false
	}
	p.mu.Lock()
	ch, ok := p.waiters[m.ID]
	if ok {
		delete(p.waiters, m.ID)
	}
	p.mu.Unlock()
	if !ok {
		return false
	}
	ch <- m
	return true
}

func (p *pendingTable) drop(id string) {
	p.mu.Lock()
	delete(p.waiters, id)
	p.mu.Unlock()
}

func (c *GrokClient) StartStdioAgent(ctx context.Context, cfg *StdioConfig) (*StdioSession, error) {
	if cfg == nil {
		cfg = &StdioConfig{}
	}
	head := []string{}
	if cfg.Model != "" {
		head = append(head, "--model", cfg.Model)
	}
	if cfg.ReasoningEffort != "" {
		head = append(head, "--reasoning-effort", cfg.ReasoningEffort)
	}
	if cfg.AlwaysApprove {
		head = append(head, "--always-approve")
	}
	if cfg.AgentProfile != "" {
		head = append(head, "--agent-profile", cfg.AgentProfile)
	}
	if cfg.UseLeader {
		head = append(head, "--leader")
	}
	if cfg.NoLeader {
		head = append(head, "--no-leader")
	}
	if cfg.GrokWSOrigin != "" {
		head = append(head, "--grok-ws-origin", cfg.GrokWSOrigin)
	}
	if cfg.GrokWSURL != "" {
		head = append(head, "--grok-ws-url", cfg.GrokWSURL)
	}
	if cfg.CLIChatProxyBaseURL != "" {
		head = append(head, "--cli-chat-proxy-base-url", cfg.CLIChatProxyBaseURL)
	}
	if cfg.XAIAPIBaseURL != "" {
		head = append(head, "--xai-api-base-url", cfg.XAIAPIBaseURL)
	}
	args := append(head, "agent", "stdio")

	cmd := execCommand(ctx, c.BinPath, args...)
	cmd.Env = c.envBase(cfg.Env)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	s := &StdioSession{
		cmd:           cmd,
		stdin:         stdin,
		stdout:        stdout,
		messages:      make(chan AgentMessage, 16),
		errs:          make(chan error, 4),
		closed:        make(chan struct{}),
		receiveCalled: make(chan struct{}, 1),
		pending:       newPendingTable(),
	}
	go s.readLoop()
	return s, nil
}

func (s *StdioSession) readLoop() {
	defer s.shutdown()
	sc := bufio.NewScanner(s.stdout)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var m AgentMessage
		if err := json.Unmarshal(line, &m); err != nil {
			select {
			case s.errs <- err:
			case <-s.closed:
				return
			}
			continue
		}
		m.Raw = append([]byte(nil), line...)
		if s.pending.deliver(m) {
			continue
		}
		select {
		case s.messages <- m:
		case <-s.closed:
			return
		}
	}
	if err := sc.Err(); err != nil {
		select {
		case s.errs <- err:
		case <-s.closed:
		}
	}
}

func (s *StdioSession) Send(m AgentMessage) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = s.stdin.Write(b)
	return err
}

func (s *StdioSession) Receive() (<-chan AgentMessage, <-chan error) {
	select {
	case s.receiveCalled <- struct{}{}:
		return s.messages, s.errs
	default:
		panic("Receive called more than once")
	}
}

func (s *StdioSession) RequestResponse(ctx context.Context, msg AgentMessage) (AgentMessage, error) {
	if msg.ID == "" {
		msg.ID = GenerateSessionID()
	}
	ch := s.pending.register(msg.ID)
	if err := s.Send(msg); err != nil {
		s.pending.drop(msg.ID)
		return AgentMessage{}, err
	}
	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		s.pending.drop(msg.ID)
		return AgentMessage{}, ctx.Err()
	case <-s.closed:
		s.pending.drop(msg.ID)
		return AgentMessage{}, errors.New("stdio session closed before response")
	}
}

func (s *StdioSession) Close() error {
	var closeErr error
	s.closeOnce.Do(func() {
		_ = s.Send(AgentMessage{Type: "shutdown"})
		_ = s.stdin.Close()
		done := make(chan struct{})
		go func() {
			_ = s.cmd.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			if s.cmd.Process != nil {
				_ = s.cmd.Process.Signal(syscall.SIGTERM)
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				if s.cmd.Process != nil {
					_ = s.cmd.Process.Kill()
				}
				<-done
			}
		}
		s.shutdown()
	})
	return closeErr
}

func (s *StdioSession) shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.closed)
		close(s.messages)
		close(s.errs)
	})
}

var ErrReceiveAlreadyCalled = errors.New("Receive already called once on this session")
