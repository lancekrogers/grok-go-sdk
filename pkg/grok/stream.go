package grok

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type EventType string

const (
	EventText              EventType = "text"
	EventThought           EventType = "thought"
	EventEnd               EventType = "end"
	EventError             EventType = "error"
	EventAssistant         EventType = "assistant"
	EventDelta             EventType = "delta"
	EventTool              EventType = "tool"
	EventToolResult        EventType = "tool_result"
	EventPermissionRequest EventType = "permission_request"
	EventResult            EventType = "result"
)

type Event struct {
	Type EventType `json:"type"`

	Data       string `json:"data,omitempty"`
	StopReason string `json:"stopReason,omitempty"`
	SessionID  string `json:"sessionId,omitempty"`
	RequestID  string `json:"requestId,omitempty"`
	Message    string `json:"message,omitempty"`

	Text string          `json:"text,omitempty"`
	Tool string          `json:"tool,omitempty"`
	Args json.RawMessage `json:"args,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (e Event) Content() string {
	switch e.Type {
	case EventText, EventThought:
		if e.Data != "" {
			return e.Data
		}
		return e.Text
	case EventError:
		return e.Message
	case EventAssistant, EventDelta:
		if e.Text != "" {
			return e.Text
		}
		return e.Data
	}
	return ""
}

func (e Event) IsTerminal() bool {
	return e.Type == EventEnd || e.Type == EventResult || e.Type == EventError
}

func parseEventLine(b []byte) (Event, error) {
	var e Event
	if err := json.Unmarshal(b, &e); err != nil {
		return Event{}, err
	}
	e.Raw = append(e.Raw[:0], b...)
	return e, nil
}

func (c *GrokClient) StreamPrompt(ctx context.Context, prompt string, opts *RunOptions) (<-chan Event, <-chan error) {
	events := make(chan Event, 16)
	errs := make(chan error, 4)

	prepared, err := c.prepareOptionsWithPrompt(opts, prompt)
	if err != nil {
		errs <- err
		close(events)
		close(errs)
		return events, errs
	}
	prepared.Format = StreamingJSONOutput

	args := BuildArgs(prompt, prepared)
	cmd := execCommand(ctx, c.BinPath, args...)
	cmd.Dir = c.workDir(prepared)
	cmd.Env = c.envFor(prepared)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errs <- err
		close(events)
		close(errs)
		return events, errs
	}

	if err := cmd.Start(); err != nil {
		errs <- err
		close(events)
		close(errs)
		return events, errs
	}

	go pumpStream(ctx, cmd, stdout, &stderr, events, errs)
	return events, errs
}

func pumpStream(ctx context.Context, cmd *exec.Cmd, r io.ReadCloser, stderr *bytes.Buffer, events chan<- Event, errs chan<- error) {
	var once sync.Once
	closeAll := func() {
		once.Do(func() {
			close(events)
			close(errs)
		})
	}
	defer closeAll()

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for sc.Scan() {
			line := sc.Bytes()
			if len(line) == 0 {
				continue
			}
			ev, err := parseEventLine(line)
			if err != nil {
				select {
				case errs <- err:
				case <-ctx.Done():
					return
				}
				continue
			}
			select {
			case events <- ev:
			case <-ctx.Done():
				return
			}
		}
		if err := sc.Err(); err != nil {
			select {
			case errs <- err:
			case <-ctx.Done():
			}
		}
	}()

	select {
	case <-done:
		reportStreamWaitError(ctx, cmd.Wait(), stderr, errs)
	case <-ctx.Done():
		if cmd.Process != nil {
			_ = cmd.Process.Signal(syscall.SIGTERM)
		}
		timer := time.NewTimer(2 * time.Second)
		defer timer.Stop()
		select {
		case <-done:
		case <-timer.C:
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
		}
		waitErr := cmd.Wait()
		if ctx.Err() == nil {
			reportStreamWaitError(ctx, waitErr, stderr, errs)
		}
	}
}

func reportStreamWaitError(ctx context.Context, err error, stderr *bytes.Buffer, errs chan<- error) {
	if err == nil {
		return
	}
	exitCode := 1
	if ee, ok := err.(*exec.ExitError); ok {
		exitCode = ee.ExitCode()
	}
	ge := ParseError(stderr.String(), exitCode)
	ge.Original = err
	select {
	case errs <- ge:
	case <-ctx.Done():
	}
}
