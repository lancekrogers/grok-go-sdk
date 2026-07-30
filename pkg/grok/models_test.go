package grok

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestModels_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	mm, err := c.Models(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mm) == 0 {
		t.Fatal("expected at least one model")
	}
	hasDefault := false
	for _, m := range mm {
		if m.Default {
			hasDefault = true
		}
	}
	if !hasDefault {
		t.Fatal("expected default-marked model")
	}
}

func TestParseModels_Captured(t *testing.T) {
	in := []byte("Default model: grok-build\n* grok-build (default)\n- grok-fast\n- grok-think\n")
	got := parseModels(in)
	if len(got) != 3 {
		t.Fatalf("got %d, want 3: %#v", len(got), got)
	}
	if got[0].ID != "grok-build" || !got[0].Default {
		t.Fatalf("row 0: %#v", got[0])
	}
}

func TestParseModels_CapturedFixture(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "models.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("fixture not yet captured: %v", err)
	}
	got := parseModels(data)
	if len(got) == 0 {
		t.Fatal("no models parsed from fixture")
	}
	if got[0].ID == "" {
		t.Fatalf("first model missing id: %#v", got[0])
	}
}

func TestEnrichModelsFromCache(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "models_cache.json")
	// Shape matches the CLI cache (info.context_window, info.name).
	const body = `{
  "fetched_at": "2026-07-30T00:00:00Z",
  "origin": "https://cli-chat-proxy.grok.com/v1/models",
  "models": {
    "grok-4.5": {
      "info": {
        "id": "grok-4.5",
        "name": "Grok 4.5",
        "context_window": 500000
      }
    },
    "grok-build-0.1": {
      "info": {
        "id": "grok-build-0.1",
        "name": "Grok Build 0.1",
        "context_window": 256000
      }
    }
  }
}`
	if err := os.WriteFile(cachePath, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	models := []ModelInfo{
		{ID: "grok-4.5", Default: true},
		{ID: "grok-build-0.1"},
		{ID: "unknown-local"},
	}
	enrichModelsFromCache(models, cachePath)

	if models[0].ContextWindow != 500000 || models[0].Name != "Grok 4.5" {
		t.Fatalf("grok-4.5 = %#v, want ContextWindow=500000 Name=Grok 4.5", models[0])
	}
	if models[1].ContextWindow != 256000 || models[1].Name != "Grok Build 0.1" {
		t.Fatalf("grok-build-0.1 = %#v, want ContextWindow=256000", models[1])
	}
	if models[2].ContextWindow != 0 || models[2].Name != "" {
		t.Fatalf("unknown model should stay unenriched: %#v", models[2])
	}
}

func TestEnrichModelsFromCache_MissingOrInvalid(t *testing.T) {
	models := []ModelInfo{{ID: "grok-4.5"}}
	enrichModelsFromCache(models, "")
	enrichModelsFromCache(models, filepath.Join(t.TempDir(), "nope.json"))
	bad := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(bad, []byte("{not json"), 0o600)
	enrichModelsFromCache(models, bad)
	if models[0].ContextWindow != 0 {
		t.Fatalf("expected no enrichment on missing/invalid cache, got %#v", models[0])
	}
}

func TestModels_EnrichesFromCache(t *testing.T) {
	mock := buildOrLocateMock(t)
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "models_cache.json")
	// Mock server lists grok-build as default (see mock models output).
	// Enrich whatever ID the mock returns by writing a matching entry after a
	// first Models() probe, or use a wildcard-free cache keyed after parse.
	c := NewClient(mock)
	listed, err := c.Models(context.Background())
	if err != nil {
		t.Fatalf("Models: %v", err)
	}
	if len(listed) == 0 {
		t.Fatal("expected models from mock")
	}
	id := listed[0].ID
	body := `{"models":{"` + id + `":{"info":{"id":"` + id + `","name":"Mock Model","context_window":424242}}}}`
	if err := os.WriteFile(cachePath, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	prev := modelsCachePathFn
	modelsCachePathFn = func() string { return cachePath }
	t.Cleanup(func() { modelsCachePathFn = prev })

	got, err := c.Models(context.Background())
	if err != nil {
		t.Fatalf("Models with cache: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected models")
	}
	var found *ModelInfo
	for i := range got {
		if got[i].ID == id {
			found = &got[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("model %q missing from %#v", id, got)
	}
	if found.ContextWindow != 424242 || found.Name != "Mock Model" {
		t.Fatalf("enriched model = %#v, want ContextWindow=424242 Name=Mock Model", *found)
	}
}
