package grok

import (
	"encoding/json"
)

type EventType string

const (
	EventAssistant         EventType = "assistant"
	EventDelta             EventType = "delta"
	EventThought           EventType = "thought"
	EventTool              EventType = "tool"
	EventToolResult        EventType = "tool_result"
	EventPermissionRequest EventType = "permission_request"
	EventResult            EventType = "result"
	EventError             EventType = "error"
)

type Event struct {
	Type      EventType       `json:"type"`
	Text      string          `json:"text,omitempty"`
	Tool      string          `json:"tool,omitempty"`
	Args      json.RawMessage `json:"args,omitempty"`
	SessionID string          `json:"sessionId,omitempty"`
	Raw       json.RawMessage `json:"-"`
}

func parseEventLine(b []byte) (Event, error) {
	var e Event
	if err := json.Unmarshal(b, &e); err != nil {
		return Event{}, err
	}
	e.Raw = append(e.Raw[:0], b...)
	return e, nil
}
