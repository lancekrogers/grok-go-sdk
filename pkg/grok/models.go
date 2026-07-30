package grok

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ModelInfo describes a model advertised by the installed Grok CLI.
//
// The CLI's `grok models` text only lists IDs and which is default. Richer
// metadata (display name, context window) is filled from the CLI models cache
// when present (~/.grok/models_cache.json), which the binary writes after
// fetching https://cli-chat-proxy.grok.com/v1/models (or the configured
// origin). ContextWindow is 0 when unknown.
type ModelInfo struct {
	ID            string
	Default       bool
	Name          string // display name from catalog; empty when unknown
	ContextWindow int    // tokens; 0 when unknown
}

// modelsCachePath returns the path to the Grok CLI models cache. Overridable
// in tests via modelsCachePathFn.
func modelsCachePath() string {
	return modelsCachePathFn()
}

var modelsCachePathFn = defaultModelsCachePath

func defaultModelsCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".grok", "models_cache.json")
}

// Models lists models available to the installed Grok CLI and enriches them
// with catalog metadata (including context_window) from the CLI models cache
// when that cache is present.
func (c *GrokClient) Models(ctx context.Context) ([]ModelInfo, error) {
	out, err := c.runSubcommand(ctx, []string{"models"})
	if err != nil {
		return nil, err
	}
	models := parseModels(out)
	enrichModelsFromCache(models, modelsCachePath())
	return models, nil
}

func parseModels(b []byte) []ModelInfo {
	var out []ModelInfo
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Default model:") {
			continue
		}
		if strings.HasPrefix(line, "*") {
			id := strings.TrimSpace(strings.TrimPrefix(line, "*"))
			id = strings.TrimSuffix(id, " (default)")
			id = strings.TrimSpace(id)
			if id != "" {
				out = append(out, ModelInfo{ID: id, Default: true})
			}
			continue
		}
		if strings.HasPrefix(line, "- ") {
			id := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			if id != "" {
				out = append(out, ModelInfo{ID: id})
			}
		}
	}
	return out
}

// modelsCacheFile is the on-disk shape written by the Grok CLI after a live
// /v1/models fetch. Only the fields the SDK needs are decoded.
type modelsCacheFile struct {
	Models map[string]modelsCacheEntry `json:"models"`
}

type modelsCacheEntry struct {
	Info modelsCacheInfo `json:"info"`
}

type modelsCacheInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ContextWindow int    `json:"context_window"`
}

// enrichModelsFromCache fills Name and ContextWindow from the CLI models cache
// for matching IDs. Missing/unreadable cache or unknown IDs leave fields zero.
// The CLI list remains the source of which models are available.
func enrichModelsFromCache(models []ModelInfo, cachePath string) {
	if cachePath == "" || len(models) == 0 {
		return
	}
	data, err := os.ReadFile(cachePath)
	if err != nil || len(data) == 0 {
		return
	}
	var cache modelsCacheFile
	if err := json.Unmarshal(data, &cache); err != nil || len(cache.Models) == 0 {
		return
	}
	for i := range models {
		entry, ok := cache.Models[models[i].ID]
		if !ok {
			// Some caches key by info.model while listing uses a shorter alias.
			entry, ok = lookupCacheEntry(cache.Models, models[i].ID)
		}
		if !ok {
			continue
		}
		if entry.Info.Name != "" {
			models[i].Name = entry.Info.Name
		}
		if entry.Info.ContextWindow > 0 {
			models[i].ContextWindow = entry.Info.ContextWindow
		}
	}
}

func lookupCacheEntry(models map[string]modelsCacheEntry, id string) (modelsCacheEntry, bool) {
	if e, ok := models[id]; ok {
		return e, true
	}
	// Case-insensitive match as a last resort (IDs are normally lowercase).
	lower := strings.ToLower(id)
	for k, e := range models {
		if strings.ToLower(k) == lower {
			return e, true
		}
		if e.Info.ID != "" && strings.ToLower(e.Info.ID) == lower {
			return e, true
		}
	}
	return modelsCacheEntry{}, false
}
