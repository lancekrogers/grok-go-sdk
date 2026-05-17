package grok

import (
	"strings"
	"testing"
)

func TestEnumConstants(t *testing.T) {
	cases := []struct {
		got, want string
	}{
		{string(PlainOutput), "plain"},
		{string(JSONOutput), "json"},
		{string(StreamingJSONOutput), "streaming-json"},
		{string(TextInput), "text"},
		{string(StreamJSONInput), "stream-json"},
		{string(EffortLow), "low"},
		{string(EffortMedium), "medium"},
		{string(EffortHigh), "high"},
		{string(EffortXHigh), "xhigh"},
		{string(EffortMax), "max"},
		{string(PermissionDefault), "default"},
		{string(PermissionAcceptEdits), "acceptEdits"},
		{string(PermissionAuto), "auto"},
		{string(PermissionDontAsk), "dontAsk"},
		{string(PermissionBypassPermissions), "bypassPermissions"},
		{string(PermissionPlan), "plan"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got %q want %q", c.got, c.want)
		}
	}
}

func TestPreprocessOptions_BestOfNWithoutPrompt(t *testing.T) {
	err := PreprocessOptions(&RunOptions{BestOfN: 4})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "BestOfN") {
		t.Errorf("expected BestOfN message, got %v", err)
	}
}

func TestPreprocessOptions_CheckWithoutPrompt(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Check: true})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPreprocessOptions_BothMemoryFlags(t *testing.T) {
	err := PreprocessOptions(&RunOptions{
		Prompt:             "x",
		ExperimentalMemory: true,
		NoMemory:           true,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("got %v", err)
	}
}

func TestPreprocessOptions_AlwaysApproveRequiresDangerous(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Prompt: "x", AlwaysApprove: true})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "AllowDangerousMode") {
		t.Errorf("got %v", err)
	}

	err = PreprocessOptions(&RunOptions{Prompt: "x", AlwaysApprove: true, AllowDangerousMode: true})
	if err != nil {
		t.Fatalf("AllowDangerousMode should permit AlwaysApprove: %v", err)
	}
}

func TestPreprocessOptions_BypassRequiresDangerous(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Prompt: "x", PermissionMode: PermissionBypassPermissions})
	if err == nil {
		t.Fatal("expected error")
	}
	err = PreprocessOptions(&RunOptions{Prompt: "x", PermissionMode: PermissionBypassPermissions, AllowDangerousMode: true})
	if err != nil {
		t.Fatalf("AllowDangerousMode should permit bypass: %v", err)
	}
}

func TestPreprocessOptions_PermissiveSandboxRequiresDangerous(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Prompt: "x", SandboxProfile: "off"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "AllowDangerousMode") {
		t.Errorf("got %v", err)
	}

	err = PreprocessOptions(&RunOptions{Prompt: "x", SandboxProfile: "off", AllowDangerousMode: true})
	if err != nil {
		t.Fatalf("AllowDangerousMode should permit permissive sandbox: %v", err)
	}
}

func TestPreprocessOptions_EnvPermissiveSandboxRequiresDangerous(t *testing.T) {
	t.Setenv("GROK_SANDBOX", "off")
	err := PreprocessOptions(&RunOptions{Prompt: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "AllowDangerousMode") {
		t.Errorf("got %v", err)
	}
}

func TestPreprocessOptions_ForkRejected(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Prompt: "x", ForkSession: true})
	if err == nil {
		t.Fatal("expected error (fork not supported yet)")
	}
}

func TestPreprocessOptions_RestoreCodeRequiresResume(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Prompt: "x", RestoreCode: true})
	if err == nil {
		t.Fatal("expected error (RestoreCode requires ResumeID)")
	}
}

func TestPreprocessOptions_EmptyAllowRule(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Prompt: "x", AllowRules: []string{""}})
	if err == nil {
		t.Fatal("expected error on empty rule")
	}
}

func TestPreprocessOptions_NullByteInRule(t *testing.T) {
	err := PreprocessOptions(&RunOptions{Prompt: "x", DenyRules: []string{"Bash(\x00)"}})
	if err == nil {
		t.Fatal("expected error on null byte")
	}
}

func TestPreprocessOptions_DefaultsAppliedWhenEmpty(t *testing.T) {
	t.Setenv("GROK_SANDBOX", "test-profile")
	opts := &RunOptions{Prompt: "x"}
	if err := PreprocessOptions(opts); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if opts.Format != PlainOutput {
		t.Errorf("Format default got %q", opts.Format)
	}
	if opts.SandboxProfile != "test-profile" {
		t.Errorf("SandboxProfile default got %q", opts.SandboxProfile)
	}
}

func TestCloneRunOptions_DeepCopy(t *testing.T) {
	orig := &RunOptions{
		AllowRules: []string{"a", "b"},
		Agents: map[string]*SubagentConfig{
			"x": {Description: "d", Tools: []string{"t1"}},
		},
		Env: []string{"K=V"},
	}

	clone := cloneRunOptions(orig)

	clone.AllowRules[0] = "MUTATED"
	clone.Agents["x"].Description = "MUTATED"
	clone.Agents["x"].Tools[0] = "MUTATED"
	clone.Env[0] = "MUTATED"
	clone.Agents["new"] = &SubagentConfig{Description: "added"}

	if orig.AllowRules[0] != "a" {
		t.Errorf("AllowRules shared: %v", orig.AllowRules)
	}
	if orig.Agents["x"].Description != "d" {
		t.Errorf("agent description shared: %v", orig.Agents["x"])
	}
	if orig.Agents["x"].Tools[0] != "t1" {
		t.Errorf("agent tools shared: %v", orig.Agents["x"].Tools)
	}
	if orig.Env[0] != "K=V" {
		t.Errorf("env shared: %v", orig.Env)
	}
	if _, exists := orig.Agents["new"]; exists {
		t.Errorf("agents map shared: %v", orig.Agents)
	}
}

func TestCloneRunOptions_Nil(t *testing.T) {
	clone := cloneRunOptions(nil)
	if clone == nil {
		t.Fatal("nil clone")
	}
}
