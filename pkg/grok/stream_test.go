package grok

import "testing"

func TestParseEventLine_KnownTypes(t *testing.T) {
	cases := []struct {
		line string
		want EventType
	}{
		{`{"type":"assistant","text":"hi"}`, EventAssistant},
		{`{"type":"tool","tool":"Read","args":{"path":"a"}}`, EventTool},
		{`{"type":"result","sessionId":"abc"}`, EventResult},
	}
	for _, tc := range cases {
		ev, err := parseEventLine([]byte(tc.line))
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if ev.Type != tc.want {
			t.Fatalf("got %q want %q", ev.Type, tc.want)
		}
		if len(ev.Raw) == 0 {
			t.Fatal("Raw not populated")
		}
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
