package grok

import (
	"context"
	"strings"
)

type ModelInfo struct {
	ID      string
	Default bool
}

func (c *GrokClient) Models(ctx context.Context) ([]ModelInfo, error) {
	out, err := c.runSubcommand(ctx, []string{"models"})
	if err != nil {
		return nil, err
	}
	return parseModels(out), nil
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
