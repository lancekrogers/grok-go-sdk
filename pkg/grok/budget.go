package grok

import "sync"

type BudgetConfig struct {
	MaxBudgetUSD     float64
	WarningThreshold float64
	OnBudgetWarning  func(current, max float64)
	OnBudgetExceeded func(current, max float64)
}

type BudgetTracker struct {
	mu         sync.Mutex
	cfg        BudgetConfig
	totalSpent float64
	warned     bool
	exceeded   bool
}

func NewBudgetTracker(cfg *BudgetConfig) *BudgetTracker {
	c := BudgetConfig{}
	if cfg != nil {
		c = *cfg
	}
	return &BudgetTracker{cfg: c}
}

func (b *BudgetTracker) Add(costUSD float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.totalSpent += costUSD
	if b.cfg.MaxBudgetUSD <= 0 {
		return
	}
	if !b.warned && b.cfg.WarningThreshold > 0 &&
		b.totalSpent >= b.cfg.MaxBudgetUSD*b.cfg.WarningThreshold {
		b.warned = true
		if b.cfg.OnBudgetWarning != nil {
			b.cfg.OnBudgetWarning(b.totalSpent, b.cfg.MaxBudgetUSD)
		}
	}
	if !b.exceeded && b.totalSpent > b.cfg.MaxBudgetUSD {
		b.exceeded = true
		if b.cfg.OnBudgetExceeded != nil {
			b.cfg.OnBudgetExceeded(b.totalSpent, b.cfg.MaxBudgetUSD)
		}
	}
}

func (b *BudgetTracker) TotalSpent() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.totalSpent
}

func (b *BudgetTracker) RemainingBudget() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cfg.MaxBudgetUSD <= 0 {
		return 0
	}
	return b.cfg.MaxBudgetUSD - b.totalSpent
}

func (b *BudgetTracker) Check() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cfg.MaxBudgetUSD > 0 && b.totalSpent > b.cfg.MaxBudgetUSD {
		return &GrokError{Type: ErrorValidation, Message: "budget exceeded"}
	}
	return nil
}
