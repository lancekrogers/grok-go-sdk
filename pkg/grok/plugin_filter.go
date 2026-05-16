package grok

import (
	"context"
	"fmt"
	"path/filepath"
)

type ToolFilterPlugin struct {
	blocked map[string]string
}

func NewToolFilterPlugin(blocked map[string]string) *ToolFilterPlugin {
	return &ToolFilterPlugin{blocked: blocked}
}

func (p *ToolFilterPlugin) Name() string { return "tool_filter" }

func (p *ToolFilterPlugin) OnBeforeRun(_ context.Context, ev BeforeRunEvent) error {
	for _, a := range ev.Args {
		for pattern, reason := range p.blocked {
			matched, _ := filepath.Match(pattern, a)
			if matched {
				return fmt.Errorf("blocked tool/arg %q: %s", a, reason)
			}
		}
	}
	return nil
}

func (p *ToolFilterPlugin) OnAfterRun(_ context.Context, ev AfterRunEvent) error {
	if ev.Result == nil {
		return nil
	}
	for pattern, reason := range p.blocked {
		matched, _ := filepath.Match(pattern, ev.Result.Subtype)
		if matched {
			return fmt.Errorf("blocked result subtype %q: %s", ev.Result.Subtype, reason)
		}
	}
	return nil
}
