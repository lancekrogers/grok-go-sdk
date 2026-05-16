package grok

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestStdioAgent_RequestResponse(t *testing.T) {
	defer goleak.VerifyNone(t)
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sess, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	resp, err := sess.RequestResponse(ctx, AgentMessage{Type: "user", ID: "1", Text: "hi"})
	if err != nil {
		t.Fatalf("rr: %v", err)
	}
	if resp.ID != "1" {
		t.Fatalf("got id %q want %q", resp.ID, "1")
	}
	if err := sess.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestStdioAgent_ReceivePanicsOnSecondCall(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sess, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer sess.Close()
	_, _ = sess.Receive()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on second Receive")
		}
	}()
	_, _ = sess.Receive()
}

func TestRequestResponse_RoundTrip(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	s, err := c.StartStdioAgent(context.Background(), &StdioConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := s.RequestResponse(ctx, AgentMessage{Type: "echo", Text: "ping"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "ping" {
		t.Fatalf("got %q want %q", resp.Text, "ping")
	}
}

func TestRequestResponse_ContextCancel(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	s, err := c.StartStdioAgent(context.Background(), &StdioConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = s.RequestResponse(ctx, AgentMessage{Type: "echo", ID: "no-such-id-will-never-match", Text: "ping"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRequestResponse_ConcurrentInFlight(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	s, err := c.StartStdioAgent(context.Background(), &StdioConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, err := s.RequestResponse(ctx, AgentMessage{Type: "echo", Text: "ping"}); err != nil {
				t.Errorf("rr: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestStdioAgent_ConcurrentSend(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sess, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer sess.Close()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = sess.Send(AgentMessage{Type: "user", Text: "ping"})
		}(i)
	}
	wg.Wait()
}
