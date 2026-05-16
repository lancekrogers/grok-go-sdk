package grok

import (
	"fmt"
	"os"
	"time"
)

type OutputFormat string

const (
	PlainOutput         OutputFormat = "plain"
	JSONOutput          OutputFormat = "json"
	StreamingJSONOutput OutputFormat = "streaming-json"
)

type InputFormat string

const (
	TextInput       InputFormat = "text"
	StreamJSONInput InputFormat = "stream-json"
)

type EffortLevel string

const (
	EffortLow    EffortLevel = "low"
	EffortMedium EffortLevel = "medium"
	EffortHigh   EffortLevel = "high"
	EffortXHigh  EffortLevel = "xhigh"
	EffortMax    EffortLevel = "max"
)

type PermissionMode string

const (
	PermissionDefault           PermissionMode = "default"
	PermissionAcceptEdits       PermissionMode = "acceptEdits"
	PermissionAuto              PermissionMode = "auto"
	PermissionDontAsk           PermissionMode = "dontAsk"
	PermissionBypassPermissions PermissionMode = "bypassPermissions"
	PermissionPlan              PermissionMode = "plan"
)

type WorktreeOption struct {
	Enabled bool
	Name    string
}

type SubagentConfig struct {
	Description string   `json:"description"`
	Prompt      string   `json:"prompt"`
	Tools       []string `json:"tools,omitempty"`
	Model       string   `json:"model,omitempty"`
	Effort      string   `json:"effort,omitempty"`
}

type RunOptions struct {
	Format      OutputFormat
	InputFormat InputFormat

	Prompt               string
	PromptFile           string
	PromptJSON           string
	Verbatim             bool
	SystemPromptOverride string
	Rules                string
	RulesFile            string

	SessionID          string
	ResumeID           string
	Continue           bool
	ForkSession        bool
	RestoreCode        bool
	NoMemory           bool
	ExperimentalMemory bool

	WorkingDirectory string
	Worktree         WorktreeOption

	Agent               string
	AgentDefinitionFile string
	Agents              map[string]*SubagentConfig
	AgentsJSON          string
	NoSubagents         bool
	AgentProfile        string

	AllowRules         []string
	DenyRules          []string
	AllowedTools       []string
	DisallowedTools    []string
	DisableWebSearch   bool
	PermissionMode     PermissionMode
	AlwaysApprove      bool
	AllowDangerousMode bool

	SandboxProfile string

	Model           string
	Effort          EffortLevel
	ReasoningEffort string

	MaxTurns    int
	BestOfN     int
	Check       bool
	NoPlan      bool
	NoAltScreen bool

	MCPConfigPath   string
	MCPConfigs      []string
	StrictMCPConfig bool

	OAuth bool

	Timeout       time.Duration
	Env           []string
	BudgetTracker *BudgetTracker
	PluginManager *PluginManager
	MaxBudgetUSD  float64
	RetryPolicy   *RetryPolicy
}

func PreprocessOptions(opts *RunOptions) error {
	if opts == nil {
		return fmt.Errorf("validation: RunOptions is nil")
	}

	hasPrompt := opts.Prompt != "" || opts.PromptFile != "" || opts.PromptJSON != ""

	if opts.BestOfN > 0 && !hasPrompt {
		return fmt.Errorf("validation: BestOfN requires a prompt source")
	}
	if opts.Check && !hasPrompt {
		return fmt.Errorf("validation: Check requires a prompt source")
	}
	if opts.ExperimentalMemory && opts.NoMemory {
		return fmt.Errorf("validation: ExperimentalMemory and NoMemory are mutually exclusive")
	}
	if opts.RestoreCode && opts.ResumeID == "" {
		return fmt.Errorf("validation: RestoreCode requires ResumeID")
	}
	if opts.ForkSession {
		return fmt.Errorf("validation: ForkSession not yet supported by grok CLI")
	}

	if err := validateRules(opts.AllowRules, "AllowRules"); err != nil {
		return err
	}
	if err := validateRules(opts.DenyRules, "DenyRules"); err != nil {
		return err
	}

	if !opts.AllowDangerousMode {
		if opts.AlwaysApprove {
			return fmt.Errorf("validation: AlwaysApprove requires AllowDangerousMode=true")
		}
		if opts.PermissionMode == PermissionBypassPermissions {
			return fmt.Errorf("validation: PermissionBypassPermissions requires AllowDangerousMode=true")
		}
	}

	if opts.Format == "" {
		opts.Format = PlainOutput
	}
	if opts.SandboxProfile == "" {
		opts.SandboxProfile = os.Getenv("GROK_SANDBOX")
	}

	return nil
}

func cloneRunOptions(opts *RunOptions) *RunOptions {
	if opts == nil {
		return &RunOptions{}
	}
	cp := *opts

	if opts.AllowRules != nil {
		cp.AllowRules = append([]string(nil), opts.AllowRules...)
	}
	if opts.DenyRules != nil {
		cp.DenyRules = append([]string(nil), opts.DenyRules...)
	}
	if opts.AllowedTools != nil {
		cp.AllowedTools = append([]string(nil), opts.AllowedTools...)
	}
	if opts.DisallowedTools != nil {
		cp.DisallowedTools = append([]string(nil), opts.DisallowedTools...)
	}
	if opts.MCPConfigs != nil {
		cp.MCPConfigs = append([]string(nil), opts.MCPConfigs...)
	}
	if opts.Env != nil {
		cp.Env = append([]string(nil), opts.Env...)
	}
	if opts.Agents != nil {
		cp.Agents = make(map[string]*SubagentConfig, len(opts.Agents))
		for k, v := range opts.Agents {
			if v == nil {
				cp.Agents[k] = nil
				continue
			}
			cv := *v
			if v.Tools != nil {
				cv.Tools = append([]string(nil), v.Tools...)
			}
			cp.Agents[k] = &cv
		}
	}

	return &cp
}

