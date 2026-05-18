package grok

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestStdioSession_FullHandshakeAndPrompt(t *testing.T) {
	defer goleak.VerifyNone(t)
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	init, err := s.Initialize(ctx, "test", "0.0.0")
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	if init.ProtocolVersion == 0 {
		t.Fatalf("initialize missing protocolVersion: %#v", init)
	}
	if len(init.AuthMethods) == 0 {
		t.Fatalf("initialize missing authMethods: %#v", init)
	}

	if _, err := s.Authenticate(ctx, "cached_token"); err != nil {
		t.Fatalf("authenticate: %v", err)
	}

	sessID, err := s.NewSession(ctx, "/tmp", nil)
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	if sessID == "" {
		t.Fatal("empty session id")
	}

	var collected []string
	done := make(chan struct{})
	go func() {
		defer close(done)
		for u := range s.Updates() {
			if u.Update.SessionUpdate == UpdateAgentMessageChunk {
				collected = append(collected, u.Update.ContentText())
			}
		}
	}()

	res, err := s.PromptText(ctx, sessID, "hello world")
	if err != nil {
		t.Fatalf("prompt: %v", err)
	}
	if res.StopReason == "" {
		t.Fatalf("missing stopReason: %#v", res)
	}
	if res.Meta == nil || res.Meta.TotalTokens == 0 {
		t.Fatalf("missing usage meta: %#v", res)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	<-done
	if len(collected) == 0 {
		t.Fatal("no agent_message_chunk content received")
	}
}

func TestStdioSession_MethodNotFound(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var ignored map[string]any
	err = s.Call(ctx, "nonexistent/method", map[string]string{}, &ignored)
	if err == nil {
		t.Fatal("expected error for unknown method")
	}
	rerr, ok := err.(*RPCError)
	if !ok {
		t.Fatalf("expected *RPCError, got %T: %v", err, err)
	}
	if rerr.Code != -32601 {
		t.Fatalf("expected code -32601, got %d", rerr.Code)
	}
}

func TestStdioSession_ContextCancel(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	s, err := c.StartStdioAgent(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.Initialize(ctx, "test", "0"); err == nil {
		t.Fatal("expected ctx.Err()")
	}
}

func TestStdioSession_ConcurrentCalls(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err := s.Initialize(ctx, "test", "0"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Authenticate(ctx, "cached_token"); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := s.NewSession(ctx, "/tmp", nil); err != nil {
				t.Errorf("new session: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestStdioSession_CloseAfterPromptDoesNotRaceReadLoop(t *testing.T) {
	mock := buildOrLocateMock(t)
	for i := 0; i < 25; i++ {
		c := NewClient(mock)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		s, err := c.StartStdioAgent(ctx, nil)
		if err != nil {
			cancel()
			t.Fatalf("start %d: %v", i, err)
		}
		if _, err := s.Initialize(ctx, "test", "0"); err != nil {
			cancel()
			t.Fatalf("initialize %d: %v", i, err)
		}
		if _, err := s.Authenticate(ctx, "cached_token"); err != nil {
			cancel()
			t.Fatalf("authenticate %d: %v", i, err)
		}
		sessID, err := s.NewSession(ctx, "/tmp", nil)
		if err != nil {
			cancel()
			t.Fatalf("new session %d: %v", i, err)
		}
		if _, err := s.PromptText(ctx, sessID, "hello world"); err != nil {
			cancel()
			t.Fatalf("prompt %d: %v", i, err)
		}
		if err := s.Close(); err != nil {
			cancel()
			t.Fatalf("close %d: %v", i, err)
		}
		select {
		case <-s.Done():
		default:
			cancel()
			t.Fatalf("session %d was not marked closed", i)
		}
		cancel()
	}
}
