package grok

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type StdioConfig struct {
	Env []string

	ClientName    string
	ClientVersion string
}

type StdioSession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser

	writeMu sync.Mutex
	nextID  atomic.Int64

	pending   sync.Map
	updatesCh chan SessionUpdate
	notifCh   chan RPCMessage
	errsCh    chan error
	closed    chan struct{}
	readDone  chan struct{}

	closeOnce   sync.Once
	closedOnce  sync.Once
	streamsOnce sync.Once
}

func (c *GrokClient) StartStdioAgent(ctx context.Context, cfg *StdioConfig) (*StdioSession, error) {
	if cfg == nil {
		cfg = &StdioConfig{}
	}
	cmd := c.command(ctx, "agent", "stdio")
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
		cmd:       cmd,
		stdin:     stdin,
		stdout:    stdout,
		updatesCh: make(chan SessionUpdate, 64),
		notifCh:   make(chan RPCMessage, 32),
		errsCh:    make(chan error, 4),
		closed:    make(chan struct{}),
		readDone:  make(chan struct{}),
	}
	go s.readLoop()
	return s, nil
}

func (s *StdioSession) Updates() <-chan SessionUpdate { return s.updatesCh }

func (s *StdioSession) Notifications() <-chan RPCMessage { return s.notifCh }

func (s *StdioSession) Errors() <-chan error { return s.errsCh }

func (s *StdioSession) Done() <-chan struct{} { return s.closed }

func (s *StdioSession) Call(ctx context.Context, method string, params any, result any) error {
	id := s.nextID.Add(1)
	idJSON, _ := json.Marshal(id)
	idStr := strconv.FormatInt(id, 10)

	ch := make(chan *RPCMessage, 1)
	s.pending.Store(idStr, ch)
	defer s.pending.Delete(idStr)

	var paramsRaw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("marshal params: %w", err)
		}
		paramsRaw = b
	}

	msg := RPCMessage{
		JSONRPC: "2.0",
		ID:      idJSON,
		Method:  method,
		Params:  paramsRaw,
	}
	if err := s.writeMessage(&msg); err != nil {
		return err
	}

	select {
	case resp := <-ch:
		if resp.Error != nil {
			return resp.Error
		}
		if result != nil && len(resp.Result) > 0 {
			if err := json.Unmarshal(resp.Result, result); err != nil {
				return fmt.Errorf("unmarshal result for %s: %w", method, err)
			}
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-s.closed:
		return errors.New("stdio session closed before response")
	}
}

func (s *StdioSession) Notify(method string, params any) error {
	var paramsRaw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return err
		}
		paramsRaw = b
	}
	return s.writeMessage(&RPCMessage{JSONRPC: "2.0", Method: method, Params: paramsRaw})
}

func (s *StdioSession) writeMessage(m *RPCMessage) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_, err = s.stdin.Write(b)
	return err
}

func (s *StdioSession) Initialize(ctx context.Context, clientName, clientVersion string) (*InitializeResult, error) {
	if clientName == "" {
		clientName = "grok-go-sdk"
	}
	if clientVersion == "" {
		clientVersion = "0.1.0"
	}
	var r InitializeResult
	err := s.Call(ctx, "initialize", InitializeParams{
		ProtocolVersion: acpProtocolVersion,
		ClientInfo:      ClientInfo{Name: clientName, Version: clientVersion},
	}, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *StdioSession) Authenticate(ctx context.Context, methodID string) (*AuthenticateResult, error) {
	var r AuthenticateResult
	if err := s.Call(ctx, "authenticate", AuthenticateParams{MethodID: methodID}, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *StdioSession) NewSession(ctx context.Context, cwd string, mcpServers []MCPServerSpec) (string, error) {
	if mcpServers == nil {
		mcpServers = []MCPServerSpec{}
	}
	var r NewSessionResult
	if err := s.Call(ctx, "session/new", NewSessionParams{CWD: cwd, MCPServers: mcpServers}, &r); err != nil {
		return "", err
	}
	return r.SessionID, nil
}

func (s *StdioSession) Prompt(ctx context.Context, sessionID string, blocks []PromptBlock) (*PromptResult, error) {
	var r PromptResult
	if err := s.Call(ctx, "session/prompt", PromptParams{SessionID: sessionID, Prompt: blocks}, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *StdioSession) PromptText(ctx context.Context, sessionID, text string) (*PromptResult, error) {
	return s.Prompt(ctx, sessionID, []PromptBlock{TextBlock(text)})
}

func (s *StdioSession) readLoop() {
	defer func() {
		s.signalClosed()
		s.closeStreams()
		close(s.readDone)
	}()
	sc := bufio.NewScanner(s.stdout)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		line := append([]byte(nil), sc.Bytes()...)
		if len(line) == 0 {
			continue
		}
		var m RPCMessage
		if err := json.Unmarshal(line, &m); err != nil {
			s.sendErr(fmt.Errorf("decode rpc message: %w (line=%s)", err, truncate(line, 200)))
			continue
		}
		m.Raw = line

		switch {
		case m.IsResponse():
			s.deliverResponse(&m)
		case m.IsRequest():
			s.handleServerRequest(&m)
		case m.IsNotification():
			s.handleNotification(&m)
		}
	}
	if err := sc.Err(); err != nil {
		s.sendErr(fmt.Errorf("scanner: %w", err))
	}
}

func (s *StdioSession) deliverResponse(m *RPCMessage) {
	var idStr string
	if err := json.Unmarshal(m.ID, &idStr); err != nil {
		var idNum int64
		if err := json.Unmarshal(m.ID, &idNum); err != nil {
			return
		}
		idStr = strconv.FormatInt(idNum, 10)
	}
	v, ok := s.pending.LoadAndDelete(idStr)
	if !ok {
		return
	}
	ch := v.(chan *RPCMessage)
	select {
	case ch <- m:
	default:
	}
}

func (s *StdioSession) handleServerRequest(m *RPCMessage) {
	resp := RPCMessage{
		JSONRPC: "2.0",
		ID:      m.ID,
		Error: &RPCError{
			Code:    -32601,
			Message: "Method not found",
		},
	}
	_ = s.writeMessage(&resp)
}

func (s *StdioSession) handleNotification(m *RPCMessage) {
	if m.Method == "session/update" && len(m.Params) > 0 {
		var su SessionUpdate
		if err := json.Unmarshal(m.Params, &su); err == nil {
			su.Raw = m.Params
			su.Update.Raw = m.Params
			select {
			case s.updatesCh <- su:
			case <-s.closed:
				return
			default:
			}
			return
		}
	}
	select {
	case s.notifCh <- *m:
	case <-s.closed:
	default:
	}
}

func (s *StdioSession) sendErr(err error) {
	select {
	case s.errsCh <- err:
	case <-s.closed:
	default:
	}
}

func (s *StdioSession) Close() error {
	var closeErr error
	s.closeOnce.Do(func() {
		s.signalClosed()
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
		<-s.readDone
	})
	return closeErr
}

func (s *StdioSession) signalClosed() {
	s.closedOnce.Do(func() {
		close(s.closed)
	})
}

func (s *StdioSession) closeStreams() {
	s.streamsOnce.Do(func() {
		close(s.updatesCh)
		close(s.notifCh)
		close(s.errsCh)
	})
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
