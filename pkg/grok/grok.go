package grok

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var execCommand = exec.CommandContext

type GrokClient struct {
	BinPath        string
	DefaultOptions *RunOptions
	Env            []string
	WorkingDir     string
	// LeaderSocket, when set, is passed as the global --leader-socket <path> so
	// spawned grok commands talk to a specific leader process.
	LeaderSocket string
}

// withLeaderSocket appends the global --leader-socket flag when configured.
func (c *GrokClient) withLeaderSocket(args []string) []string {
	if c.LeaderSocket == "" {
		return args
	}
	out := make([]string, 0, len(args)+2)
	out = append(out, args...)
	out = append(out, "--leader-socket", c.LeaderSocket)
	return out
}

func (c *GrokClient) command(ctx context.Context, args ...string) *exec.Cmd {
	return execCommand(ctx, c.BinPath, c.withLeaderSocket(args)...)
}

func NewClient(binPath string) *GrokClient {
	return &GrokClient{
		BinPath:        binPath,
		DefaultOptions: &RunOptions{Format: PlainOutput},
	}
}

func NewClientFromPath() (*GrokClient, error) {
	p, err := LocateBinary()
	if err != nil {
		return nil, err
	}
	return NewClient(p), nil
}

func (c *GrokClient) RunPrompt(prompt string, opts *RunOptions) (*GrokResult, error) {
	return c.RunPromptCtx(context.Background(), prompt, opts)
}

func (c *GrokClient) RunPromptCtx(ctx context.Context, prompt string, opts *RunOptions) (*GrokResult, error) {
	prepared, err := c.prepareOptionsWithPrompt(opts, prompt)
	if err != nil {
		return nil, err
	}
	if prepared.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, prepared.Timeout)
		defer cancel()
	}
	args := BuildArgs(prompt, prepared)
	if prepared.BudgetTracker != nil {
		if err := prepared.BudgetTracker.Check(); err != nil {
			return nil, err
		}
		if prepared.MaxBudgetUSD > 0 && prepared.BudgetTracker.TotalSpent() > prepared.MaxBudgetUSD {
			return nil, &GrokError{Type: ErrorValidation, Message: "per-call budget exceeded"}
		}
	}
	if prepared.PluginManager != nil {
		if err := prepared.PluginManager.fireBefore(ctx, BeforeRunEvent{Prompt: prompt, Opts: prepared, Args: args}); err != nil {
			return nil, err
		}
	}
	cmd := c.command(ctx, args...)
	cmd.Dir = c.workDir(prepared)
	cmd.Env = c.envFor(prepared)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		exitCode := 1
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		}
		ge := ParseError(stderr.String(), exitCode)
		ge.Original = err
		if prepared.PluginManager != nil {
			_ = prepared.PluginManager.fireAfter(ctx, AfterRunEvent{Prompt: prompt, Opts: prepared, Args: args, Err: ge})
		}
		return nil, ge
	}
	result, err := decodeOutput(prepared.Format, stdout.Bytes())
	if err != nil {
		return nil, err
	}
	if prepared.BudgetTracker != nil && result != nil {
		prepared.BudgetTracker.Add(result.CostUSD)
	}
	if prepared.PluginManager != nil {
		if err := prepared.PluginManager.fireAfter(ctx, AfterRunEvent{Prompt: prompt, Opts: prepared, Args: args, Result: result}); err != nil {
			return result, err
		}
	}
	return result, nil
}

func (c *GrokClient) runSubcommandTolerant(ctx context.Context, args []string) ([]byte, error) {
	cmd := c.command(ctx, args...)
	cmd.Dir = c.WorkingDir
	cmd.Env = c.envBase(nil)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return stdout.Bytes(), nil
	}
	if ee, ok := err.(*exec.ExitError); ok {
		ge := ParseError(stderr.String(), ee.ExitCode())
		ge.Original = err
		return stdout.Bytes(), ge
	}
	return stdout.Bytes(), err
}

func (c *GrokClient) runSubcommand(ctx context.Context, args []string) ([]byte, error) {
	cmd := c.command(ctx, args...)
	cmd.Dir = c.WorkingDir
	cmd.Env = c.envBase(nil)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		exitCode := 1
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		}
		ge := ParseError(stderr.String(), exitCode)
		ge.Original = err
		return nil, ge
	}
	return stdout.Bytes(), nil
}

func (c *GrokClient) prepareOptions(opts *RunOptions) (*RunOptions, error) {
	return c.prepareOptionsWithPrompt(opts, "")
}

func (c *GrokClient) prepareOptionsWithPrompt(opts *RunOptions, prompt string) (*RunOptions, error) {
	if opts == nil {
		opts = c.DefaultOptions
	}
	prepared := cloneRunOptions(opts)
	if prompt != "" && prepared.Prompt == "" {
		prepared.Prompt = prompt
	}
	if err := PreprocessOptions(prepared); err != nil {
		return nil, err
	}
	if prompt != "" && prepared.Prompt == prompt {
		prepared.Prompt = ""
	}
	return prepared, nil
}

func (c *GrokClient) workDir(opts *RunOptions) string {
	if opts != nil && opts.WorkingDirectory != "" {
		return opts.WorkingDirectory
	}
	return c.WorkingDir
}

func (c *GrokClient) envBase(extra []string) []string {
	base := append([]string(nil), os.Environ()...)
	base = append(base, c.Env...)
	base = append(base, extra...)
	return base
}

func (c *GrokClient) envFor(opts *RunOptions) []string {
	if opts == nil {
		return c.envBase(nil)
	}
	return c.envBase(opts.Env)
}

func decodeOutput(format OutputFormat, data []byte) (*GrokResult, error) {
	if format == JSONOutput {
		var r GrokResult
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, &GrokError{Type: ErrorValidation, Message: "decode: " + err.Error()}
		}
		return &r, nil
	}
	return &GrokResult{Text: string(data)}, nil
}

func BuildArgs(prompt string, opts *RunOptions) []string {
	args := []string{}

	switch {
	case opts.PromptFile != "":
		args = append(args, "--prompt-file", opts.PromptFile)
	case opts.PromptJSON != "":
		args = append(args, "--prompt-json", opts.PromptJSON)
	case prompt != "":
		args = append(args, "-p", prompt)
	case opts.Prompt != "":
		args = append(args, "-p", opts.Prompt)
	}

	if opts.Format != "" {
		args = append(args, "--output-format", string(opts.Format))
	}
	if opts.JSONSchema != "" {
		args = append(args, "--json-schema", opts.JSONSchema)
	}

	if opts.Agent != "" {
		args = append(args, "--agent", opts.Agent)
	} else if opts.AgentDefinitionFile != "" {
		args = append(args, "--agent", opts.AgentDefinitionFile)
	}
	if opts.AgentsJSON != "" {
		args = append(args, "--agents", opts.AgentsJSON)
	} else if len(opts.Agents) > 0 {
		if b, err := json.Marshal(opts.Agents); err == nil {
			args = append(args, "--agents", string(b))
		}
	}
	if opts.NoSubagents {
		args = append(args, "--no-subagents")
	}

	for _, r := range opts.AllowRules {
		args = append(args, "--allow", r)
	}
	for _, r := range opts.DenyRules {
		args = append(args, "--deny", r)
	}
	if len(opts.AllowedTools) > 0 {
		args = append(args, "--tools", strings.Join(opts.AllowedTools, ","))
	}
	if len(opts.DisallowedTools) > 0 {
		args = append(args, "--disallowed-tools", strings.Join(opts.DisallowedTools, ","))
	}
	if opts.DisableWebSearch {
		args = append(args, "--disable-web-search")
	}
	if opts.AlwaysApprove {
		args = append(args, "--always-approve")
	}
	if opts.PermissionMode != "" {
		args = append(args, "--permission-mode", string(opts.PermissionMode))
	}

	if opts.SandboxProfile != "" {
		args = append(args, "--sandbox", opts.SandboxProfile)
	}

	if opts.ResumeID != "" {
		args = append(args, "--resume", opts.ResumeID)
	} else if opts.Continue {
		args = append(args, "--continue")
	}
	if opts.RestoreCode {
		args = append(args, "--restore-code")
	}
	if opts.ForkSession {
		args = append(args, "--fork-session")
	}
	if opts.SessionID != "" {
		args = append(args, "-s", opts.SessionID)
	}

	if opts.NoMemory {
		args = append(args, "--no-memory")
	}
	if opts.ExperimentalMemory {
		args = append(args, "--experimental-memory")
	}

	if opts.WorkingDirectory != "" {
		args = append(args, "--cwd", opts.WorkingDirectory)
	}
	if opts.Worktree.Enabled {
		if opts.Worktree.Name != "" {
			args = append(args, "-w", opts.Worktree.Name)
		} else {
			args = append(args, "-w")
		}
		if opts.Worktree.Ref != "" {
			args = append(args, "--worktree-ref", opts.Worktree.Ref)
		}
	}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}
	if opts.Effort != "" {
		args = append(args, "--effort", string(opts.Effort))
	}
	if opts.ReasoningEffort != "" {
		args = append(args, "--reasoning-effort", opts.ReasoningEffort)
	}

	if opts.MaxTurns > 0 {
		args = append(args, "--max-turns", strconv.Itoa(opts.MaxTurns))
	}
	if opts.BestOfN > 0 {
		args = append(args, "--best-of-n", strconv.Itoa(opts.BestOfN))
	}
	if opts.Check {
		args = append(args, "--check")
	}
	if opts.NoPlan {
		args = append(args, "--no-plan")
	}
	if opts.NoAltScreen {
		args = append(args, "--no-alt-screen")
	}

	if opts.SystemPromptOverride != "" {
		args = append(args, "--system-prompt-override", opts.SystemPromptOverride)
	}
	if opts.RulesFile != "" {
		args = append(args, "--rules", "@"+opts.RulesFile)
	} else if opts.Rules != "" {
		args = append(args, "--rules", opts.Rules)
	}
	if opts.Verbatim {
		args = append(args, "--verbatim")
	}

	if opts.OAuth {
		args = append(args, "--oauth")
	}

	return args
}

func init() {
	_ = time.Now
}
