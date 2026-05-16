package grok

import (
	"context"
	"errors"
	"os/exec"
	"sync/atomic"
	"testing"
)

type recordingPlugin struct {
	name    string
	order   *[]string
	before  func() error
}

func (r *recordingPlugin) Name() string { return r.name }
func (r *recordingPlugin) OnBeforeRun(_ context.Context, _ BeforeRunEvent) error {
	*r.order = append(*r.order, r.name+":before")
	if r.before != nil {
		return r.before()
	}
	return nil
}
func (r *recordingPlugin) OnAfterRun(_ context.Context, _ AfterRunEvent) error {
	*r.order = append(*r.order, r.name+":after")
	return nil
}

func TestPlugins_FireInRegistrationOrder(t *testing.T) {
	pm := NewPluginManager()
	var order []string
	pm.Register(&recordingPlugin{name: "a", order: &order}, nil)
	pm.Register(&recordingPlugin{name: "b", order: &order}, nil)
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	_, err := c.RunPromptCtx(context.Background(), "hi", &RunOptions{PluginManager: pm, Format: JSONOutput})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"a:before", "b:before", "a:after", "b:after"}
	if len(order) != len(want) {
		t.Fatalf("got %v want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("idx %d: got %s want %s", i, order[i], want[i])
		}
	}
}

func TestPlugins_BeforeRunErrorAbortsExec(t *testing.T) {
	pm := NewPluginManager()
	var order []string
	var execCount atomic.Int32
	pm.Register(&recordingPlugin{name: "blocker", order: &order, before: func() error { return errors.New("nope") }}, nil)
	prev := execCommand
	execCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		execCount.Add(1)
		return prev(ctx, name, args...)
	}
	t.Cleanup(func() { execCommand = prev })

	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	_, err := c.RunPromptCtx(context.Background(), "hi", &RunOptions{PluginManager: pm, Format: JSONOutput})
	if err == nil {
		t.Fatal("expected error")
	}
	if execCount.Load() != 0 {
		t.Fatalf("expected 0 execs, got %d", execCount.Load())
	}
}

func TestMetricsPlugin_Accumulates(t *testing.T) {
	pm := NewPluginManager()
	m := &MetricsPlugin{}
	pm.Register(m, nil)
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	for i := 0; i < 3; i++ {
		if _, err := c.RunPromptCtx(context.Background(), "hi", &RunOptions{PluginManager: pm, Format: JSONOutput}); err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
	}
	if got := m.GetMetrics()["runs"]; got != 3 {
		t.Fatalf("runs=%v want 3", got)
	}
}

func TestAuditPlugin_RingBuffer(t *testing.T) {
	a := NewAuditPlugin(2)
	pm := NewPluginManager()
	pm.Register(a, nil)
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	for i := 0; i < 3; i++ {
		_, _ = c.RunPromptCtx(context.Background(), "hi", &RunOptions{PluginManager: pm, Format: JSONOutput})
	}
	if got := len(a.GetRecords()); got != 2 {
		t.Fatalf("records=%d want 2", got)
	}
}
