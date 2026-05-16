package grok

import (
	"sync"
	"testing"
)

func TestBudgetTracker_WarningCallback(t *testing.T) {
	var warnings int
	bt := NewBudgetTracker(&BudgetConfig{
		MaxBudgetUSD:     1.0,
		WarningThreshold: 0.5,
		OnBudgetWarning:  func(_, _ float64) { warnings++ },
	})
	bt.Add(0.3)
	bt.Add(0.3)
	bt.Add(0.1)
	if warnings != 1 {
		t.Fatalf("warnings=%d want 1", warnings)
	}
}

func TestBudgetTracker_ExceededCallback(t *testing.T) {
	var exceeded int
	bt := NewBudgetTracker(&BudgetConfig{
		MaxBudgetUSD:     1.0,
		OnBudgetExceeded: func(_, _ float64) { exceeded++ },
	})
	bt.Add(0.5)
	bt.Add(0.6)
	bt.Add(0.5)
	if exceeded != 1 {
		t.Fatalf("exceeded=%d want 1", exceeded)
	}
}

func TestBudgetTracker_Concurrent_Add(t *testing.T) {
	bt := NewBudgetTracker(&BudgetConfig{MaxBudgetUSD: 100})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bt.Add(0.01)
		}()
	}
	wg.Wait()
	if got := bt.TotalSpent(); got < 0.99 || got > 1.01 {
		t.Fatalf("total=%f want ~1.0", got)
	}
}

func TestBudgetTracker_CheckRefusesOverBudget(t *testing.T) {
	bt := NewBudgetTracker(&BudgetConfig{MaxBudgetUSD: 0.5})
	bt.Add(0.6)
	if err := bt.Check(); err == nil {
		t.Fatal("expected error")
	}
}
