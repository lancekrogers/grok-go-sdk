package grok

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildArgs_Coverage(t *testing.T) {
	cases := []struct {
		name     string
		prompt   string
		opts     *RunOptions
		contains []string
	}{
		{"plain prompt", "hello", &RunOptions{}, []string{"-p", "hello"}},
		{"prompt file", "", &RunOptions{PromptFile: "/tmp/p.txt"}, []string{"--prompt-file", "/tmp/p.txt"}},
		{"prompt json", "", &RunOptions{PromptJSON: `{"k":"v"}`}, []string{"--prompt-json", `{"k":"v"}`}},
		{"json format", "x", &RunOptions{Format: JSONOutput}, []string{"--output-format", "json"}},
		{"streaming format", "x", &RunOptions{Format: StreamingJSONOutput}, []string{"--output-format", "streaming-json"}},
		{"agent name", "x", &RunOptions{Agent: "security"}, []string{"--agent", "security"}},
		{"agent file fallback", "x", &RunOptions{AgentDefinitionFile: "/tmp/a.json"}, []string{"--agent", "/tmp/a.json"}},
		{"agents json raw", "x", &RunOptions{AgentsJSON: `{"x":1}`}, []string{"--agents", `{"x":1}`}},
		{"agents map", "x", &RunOptions{Agents: map[string]*SubagentConfig{"x": {Description: "d", Prompt: "p"}}}, []string{"--agents"}},
		{"no subagents", "x", &RunOptions{NoSubagents: true}, []string{"--no-subagents"}},
		{"allow rules", "x", &RunOptions{AllowRules: []string{"Read(*)", "Bash(git *)"}}, []string{"--allow", "Read(*)", "--allow", "Bash(git *)"}},
		{"deny rules", "x", &RunOptions{DenyRules: []string{"Bash(rm*)"}}, []string{"--deny", "Bash(rm*)"}},
		{"allowed tools csv", "x", &RunOptions{AllowedTools: []string{"Read", "Grep"}}, []string{"--tools", "Read,Grep"}},
		{"disallowed tools csv", "x", &RunOptions{DisallowedTools: []string{"Bash", "Write"}}, []string{"--disallowed-tools", "Bash,Write"}},
		{"disable web search", "x", &RunOptions{DisableWebSearch: true}, []string{"--disable-web-search"}},
		{"always approve", "x", &RunOptions{AlwaysApprove: true, AllowDangerousMode: true}, []string{"--always-approve"}},
		{"permission mode", "x", &RunOptions{PermissionMode: PermissionDontAsk}, []string{"--permission-mode", "dontAsk"}},
		{"sandbox", "x", &RunOptions{SandboxProfile: "strict"}, []string{"--sandbox", "strict"}},
		{"resume id", "x", &RunOptions{ResumeID: "abc"}, []string{"--resume", "abc"}},
		{"continue", "x", &RunOptions{Continue: true}, []string{"--continue"}},
		{"restore code", "x", &RunOptions{ResumeID: "abc", RestoreCode: true}, []string{"--restore-code"}},
		{"session id alone", "x", &RunOptions{SessionID: "uuid-1"}, []string{"-s", "uuid-1"}},
		{"fork with resume", "x", &RunOptions{ResumeID: "abc", ForkSession: true, SessionID: "uuid-2"}, []string{"--resume", "abc", "--fork-session", "-s", "uuid-2"}},
		{"no memory", "x", &RunOptions{NoMemory: true}, []string{"--no-memory"}},
		{"experimental memory", "x", &RunOptions{ExperimentalMemory: true}, []string{"--experimental-memory"}},
		{"cwd", "x", &RunOptions{WorkingDirectory: "/tmp"}, []string{"--cwd", "/tmp"}},
		{"worktree named", "x", &RunOptions{Worktree: WorktreeOption{Enabled: true, Name: "wt1"}}, []string{"-w", "wt1"}},
		{"worktree bare", "x", &RunOptions{Worktree: WorktreeOption{Enabled: true}}, []string{"-w"}},
		{"model", "x", &RunOptions{Model: "grok-build"}, []string{"--model", "grok-build"}},
		{"effort", "x", &RunOptions{Effort: EffortHigh}, []string{"--effort", "high"}},
		{"reasoning effort", "x", &RunOptions{ReasoningEffort: "high"}, []string{"--reasoning-effort", "high"}},
		{"max turns", "x", &RunOptions{MaxTurns: 5}, []string{"--max-turns", "5"}},
		{"best of n", "x", &RunOptions{BestOfN: 4}, []string{"--best-of-n", "4"}},
		{"check", "x", &RunOptions{Check: true}, []string{"--check"}},
		{"no plan", "x", &RunOptions{NoPlan: true}, []string{"--no-plan"}},
		{"no alt screen", "x", &RunOptions{NoAltScreen: true}, []string{"--no-alt-screen"}},
		{"system prompt override", "x", &RunOptions{SystemPromptOverride: "you are X"}, []string{"--system-prompt-override", "you are X"}},
		{"rules inline", "x", &RunOptions{Rules: "be terse"}, []string{"--rules", "be terse"}},
		{"rules file", "x", &RunOptions{RulesFile: "/tmp/r.txt"}, []string{"--rules", "@/tmp/r.txt"}},
		{"verbatim", "x", &RunOptions{Verbatim: true}, []string{"--verbatim"}},
		{"oauth", "x", &RunOptions{OAuth: true}, []string{"--oauth"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildArgs(tc.prompt, tc.opts)
			joined := strings.Join(got, " ")
			for i := 0; i < len(tc.contains); i++ {
				if !contains(got, tc.contains[i]) {
					t.Errorf("missing %q in %q", tc.contains[i], joined)
				}
			}
		})
	}
}

func contains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}

func TestBuildArgs_AllowDenyOrder(t *testing.T) {
	args := BuildArgs("hi", &RunOptions{
		AllowRules: []string{"a1", "a2"},
		DenyRules:  []string{"d1"},
	})
	want := []string{"-p", "hi", "--allow", "a1", "--allow", "a2", "--deny", "d1"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("got %v want %v", args, want)
	}
}

func TestBuildArgs_JSONSchema(t *testing.T) {
	schema := "{\n  \"type\": \"object\",\n  \"properties\": {\n    \"name\": {\"type\": \"string\"}\n  }\n}"
	opts := &RunOptions{Prompt: "x", JSONSchema: schema}
	if err := PreprocessOptions(opts); err != nil {
		t.Fatalf("preprocess: %v", err)
	}
	if opts.Format != JSONOutput {
		t.Fatalf("JSONSchema should imply JSONOutput, got %q", opts.Format)
	}
	args := BuildArgs("x", opts)
	if !contains(args, "--output-format") || !contains(args, "json") {
		t.Fatalf("missing implied --output-format json: %v", args)
	}
	// The schema must travel as a single argv element, never split on its newlines.
	pos := -1
	for i, a := range args {
		if a == "--json-schema" {
			pos = i
			break
		}
	}
	if pos < 0 || pos+1 >= len(args) || args[pos+1] != schema {
		t.Fatalf("schema not passed as a single argv element: %v", args)
	}
}

func TestBuildArgs_PromptSourcePrecedence(t *testing.T) {
	args := BuildArgs("from-arg", &RunOptions{Prompt: "from-opts", PromptFile: "/f"})
	if !contains(args, "--prompt-file") {
		t.Error("prompt-file should win over both")
	}
	if contains(args, "from-arg") || contains(args, "from-opts") {
		t.Errorf("prompt-file precedence violated: %v", args)
	}
}

func TestRunPromptCtx_Decode(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	res, err := c.RunPrompt("hi", &RunOptions{Format: JSONOutput})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if res.Text == "" {
		t.Fatal("empty text")
	}
	if res.SessionID == "" {
		t.Fatal("empty session id")
	}
}

func TestRunSubcommand_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	out, err := c.runSubcommand(t.Context(), []string{"models"})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if !strings.Contains(string(out), "grok-build") {
		t.Errorf("expected 'grok-build' in output, got %q", string(out))
	}
}

func TestPrepareOptionsWithPrompt_CheckAndBestOfN(t *testing.T) {
	c := NewClient("/nonexistent")

	// Check with prompt via arg should succeed (the new WithPrompt injects for hasPrompt validation in PreprocessOptions)
	_, err := c.prepareOptionsWithPrompt(&RunOptions{Check: true}, "Find the largest prime below 100.")
	if err != nil {
		t.Fatalf("Check with prompt arg should not error: %v", err)
	}

	// BestOfN too
	_, err = c.prepareOptionsWithPrompt(&RunOptions{BestOfN: 3}, "hello")
	if err != nil {
		t.Fatalf("BestOfN with prompt arg should not error: %v", err)
	}

	// When prompt already in opts, still works and keeps it
	_, err = c.prepareOptionsWithPrompt(&RunOptions{Check: true, Prompt: "pre-set"}, "")
	if err != nil {
		t.Fatalf("Check with pre-set Prompt should not error: %v", err)
	}
}
