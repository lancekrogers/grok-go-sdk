package grok

import (
	"encoding/json"
	"testing"
)

func TestMarshalAgents_RoundTrip(t *testing.T) {
	in := map[string]*SubagentConfig{
		"summarizer": {Description: "summarize text", Prompt: "summarize."},
	}
	s, err := MarshalAgents(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if s == "" {
		t.Fatal("empty json")
	}
	var out map[string]*SubagentConfig
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out["summarizer"].Description != "summarize text" {
		t.Fatalf("round-trip mismatch: %#v", out)
	}
}

func TestMarshalAgents_Empty(t *testing.T) {
	s, err := MarshalAgents(nil)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if s != "" {
		t.Fatalf("expected empty string, got %q", s)
	}
}
