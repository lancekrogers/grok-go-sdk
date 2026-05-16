package grok

import (
	"context"
	"sync"
)

type BeforeRunEvent struct {
	Prompt string
	Opts   *RunOptions
	Args   []string
}

type AfterRunEvent struct {
	Prompt string
	Opts   *RunOptions
	Args   []string
	Result *GrokResult
	Err    error
}

type Plugin interface {
	Name() string
	OnBeforeRun(ctx context.Context, ev BeforeRunEvent) error
	OnAfterRun(ctx context.Context, ev AfterRunEvent) error
}

type PluginManager struct {
	mu      sync.Mutex
	plugins []Plugin
	init    bool
}

func NewPluginManager() *PluginManager { return &PluginManager{} }

func (pm *PluginManager) Register(p Plugin, _ any) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins = append(pm.plugins, p)
}

func (pm *PluginManager) Initialize(ctx context.Context) error {
	pm.mu.Lock()
	pm.init = true
	pm.mu.Unlock()
	return nil
}

func (pm *PluginManager) Shutdown(ctx context.Context) error {
	pm.mu.Lock()
	pm.init = false
	pm.mu.Unlock()
	return nil
}

func (pm *PluginManager) fireBefore(ctx context.Context, ev BeforeRunEvent) error {
	pm.mu.Lock()
	list := append([]Plugin(nil), pm.plugins...)
	pm.mu.Unlock()
	for _, p := range list {
		if err := p.OnBeforeRun(ctx, ev); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PluginManager) fireAfter(ctx context.Context, ev AfterRunEvent) error {
	pm.mu.Lock()
	list := append([]Plugin(nil), pm.plugins...)
	pm.mu.Unlock()
	for _, p := range list {
		if err := p.OnAfterRun(ctx, ev); err != nil {
			return err
		}
	}
	return nil
}
