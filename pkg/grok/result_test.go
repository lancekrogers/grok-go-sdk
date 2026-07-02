package grok

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGrokResult_DecodeSayHello(t *testing.T) {
	_, this, _, _ := runtime.Caller(0)
	repo := filepath.Join(filepath.Dir(this), "..", "..")
	path := filepath.Join(repo, "test", "testdata", "json", "say-hello.json")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("fixture not yet captured: %v", err)
	}
	if len(data) == 0 {
		t.Skip("fixture is empty (scaffold placeholder)")
	}
	var r GrokResult
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if r.SessionID == "" {
		t.Error("session id missing")
	}
	if r.Text == "" {
		t.Error("text missing")
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("raw decode failed: %v", err)
	}
	known := map[string]bool{
		"text": true, "stopReason": true, "sessionId": true, "requestId": true, "thought": true,
		"costUsd": true, "durationMs": true, "durationApiMs": true, "numTurns": true,
		"isError": true, "subtype": true,
	}
	for key := range raw {
		if !known[key] {
			t.Fatalf("GrokResult fixture contains unmapped field %q", key)
		}
	}
	if _, ok := raw["text"]; !ok {
		t.Fatal("raw text field missing")
	}
}
