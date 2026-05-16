package grok

import (
	"context"
	"sync"
	"time"
)

type AuditRecord struct {
	Time   time.Time
	Prompt string
	Args   []string
	Result *GrokResult
	Err    error
}

type AuditPlugin struct {
	mu         sync.Mutex
	maxRecords int
	records    []AuditRecord
}

func NewAuditPlugin(maxRecords int) *AuditPlugin {
	if maxRecords <= 0 {
		maxRecords = 100
	}
	return &AuditPlugin{maxRecords: maxRecords}
}

func (p *AuditPlugin) Name() string { return "audit" }

func (p *AuditPlugin) OnBeforeRun(_ context.Context, _ BeforeRunEvent) error { return nil }

func (p *AuditPlugin) OnAfterRun(_ context.Context, ev AfterRunEvent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	rec := AuditRecord{
		Time:   time.Now(),
		Prompt: ev.Prompt,
		Args:   ev.Args,
		Result: ev.Result,
		Err:    ev.Err,
	}
	p.records = append(p.records, rec)
	if len(p.records) > p.maxRecords {
		p.records = p.records[len(p.records)-p.maxRecords:]
	}
	return nil
}

func (p *AuditPlugin) GetRecords() []AuditRecord {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]AuditRecord, len(p.records))
	copy(out, p.records)
	return out
}
