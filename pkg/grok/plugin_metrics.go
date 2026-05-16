package grok

import (
	"context"
	"sync"
)

type MetricsPlugin struct {
	mu            sync.Mutex
	runs          float64
	errors        float64
	totalCostUSD  float64
	totalDuration float64
}

func (p *MetricsPlugin) Name() string { return "metrics" }

func (p *MetricsPlugin) OnBeforeRun(_ context.Context, _ BeforeRunEvent) error { return nil }

func (p *MetricsPlugin) OnAfterRun(_ context.Context, ev AfterRunEvent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.runs++
	if ev.Err != nil {
		p.errors++
		return nil
	}
	if ev.Result != nil {
		p.totalCostUSD += ev.Result.CostUSD
		p.totalDuration += float64(ev.Result.DurationMS)
	}
	return nil
}

func (p *MetricsPlugin) GetMetrics() map[string]float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return map[string]float64{
		"runs":           p.runs,
		"errors":         p.errors,
		"total_cost_usd": p.totalCostUSD,
		"total_duration": p.totalDuration,
	}
}
