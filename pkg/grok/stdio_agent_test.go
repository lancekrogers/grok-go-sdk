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
