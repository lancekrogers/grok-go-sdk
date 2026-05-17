package grok

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestParseEventLine_KnownTypes(t *testing.T) {
	cases := []struct {
		line        string
		want        EventType
		wantContent string
	}{
		{`{"type":"thought","data":"The"}`, EventThought, "The"},
		{`{"type":"text","data":"Hello"}`, EventText, "Hello"},
		{`{"type":"end","stopReason":"EndTurn","sessionId":"abc","requestId":"r1"}`, EventEnd, ""},
		{`{"type":"error","message":"unknown model"}`, EventError, "unknown model"},
	}
	for _, tc := range cases {
		ev, err := parseEventLine([]byte(tc.line))
		if err != nil {
			t.Fatalf("parse %q: %v", tc.line, err)
		}
		if ev.Type != tc.want {
			t.Fatalf("got %q want %q", ev.Type, tc.want)
		}
		if got := ev.Content(); got != tc.wantContent {
			t.Fatalf("content for %q: got %q want %q", tc.line, got, tc.wantContent)
		}
		if len(ev.Raw) == 0 {
			t.Fatal("Raw not populated")
		}
	}
}

func TestEvent_EndCarriesSessionMetadata(t *testing.T) {
	ev, err := parseEventLine([]byte(`{"type":"end","stopReason":"EndTurn","sessionId":"sess-1","requestId":"req-1"}`))
	if err != nil {
		t.Fatal(err)
	}
	if !ev.IsTerminal() {
		t.Fatal("end should be terminal")
	}
	if ev.SessionID != "sess-1" || ev.RequestID != "req-1" || ev.StopReason != "EndTurn" {
		t.Fatalf("missing metadata: %#v", ev)
	}
}

func TestParseEventLine_UnknownType(t *testing.T) {
	ev, err := parseEventLine([]byte(`{"type":"future_event","custom":42}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev.Type != "future_event" {
		t.Fatalf("got %q", ev.Type)
	}
	if len(ev.Raw) == 0 {
		t.Fatal("Raw empty for unknown type")
	}
}

func TestParseEventLine_InvalidJSON(t *testing.T) {
	_, err := parseEventLine([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStreamPrompt_AgainstMock(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx := context.Background()
	events, errs := c.StreamPrompt(ctx, "hello", &RunOptions{})
	var got []Event
	done := false
	for !done {
		select {
		case ev, ok := <-events:
			if !ok {
				done = true
				break
			}
			got = append(got, ev)
		case err, ok := <-errs:
			if !ok {
				continue
			}
			if err != nil {
				t.Fatalf("stream error: %v", err)
			}
		}
	}
	if len(got) == 0 {
		t.Fatal("no events received")
	}
}

func TestStreamPrompt_ProcessErrorReported(t *testing.T) {
	mock := buildOrLocateMock(t)
	t.Setenv("GROK_MOCK_EXIT_CODE", "2")
	t.Setenv("GROK_MOCK_STDERR", "invalid stream request")
	c := NewClient(mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	events, errs := c.StreamPrompt(ctx, "hello", &RunOptions{})
	var gotErr error
	for events != nil || errs != nil {
		select {
		case _, ok := <-events:
			if !ok {
				events = nil
			}
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				gotErr = err
			}
		case <-ctx.Done():
			t.Fatalf("timeout waiting for stream close: %v", ctx.Err())
		}
	}
	if gotErr == nil {
		t.Fatal("expected stream process error")
	}
	var ge *GrokError
	if !errors.As(gotErr, &ge) {
		t.Fatalf("expected *GrokError, got %T: %v", gotErr, gotErr)
	}
	if ge.Type != ErrorValidation {
		t.Fatalf("got error type %q want %q", ge.Type, ErrorValidation)
	}
}

func TestStreamPrompt_Cancellation_NoGoroutineLeak(t *testing.T) {
	defer goleak.VerifyNone(t)
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx, cancel := context.WithCancel(context.Background())
	events, errs := c.StreamPrompt(ctx, "hello", &RunOptions{})
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	for range events {
	}
	for range errs {
	}
}
