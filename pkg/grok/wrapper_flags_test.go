package grok

import (
	"context"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// captureArgs overrides execCommand to record the emitted argv and run a
// harmless no-op command instead of the real binary.
func captureArgs(t *testing.T) *[]string {
	t.Helper()
	prev := execCommand
	got := new([]string)
	execCommand = func(ctx context.Context, _ string, args ...string) *exec.Cmd {
		*got = args
		return exec.CommandContext(ctx, "true")
	}
	t.Cleanup(func() { execCommand = prev })
	return got
}

func TestMCPRemove_Scope(t *testing.T) {
	got := captureArgs(t)
	c := NewClient("grok")
	_ = c.MCPRemove(context.Background(), "srv", MCPScopeProject)
	want := []string{"mcp", "remove", "-s", "project", "srv"}
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("got %v want %v", *got, want)
	}
}

func TestMCPDoctor_NameJSON(t *testing.T) {
	got := captureArgs(t)
	c := NewClient("grok")
	_, _ = c.MCPDoctor(context.Background(), "srv", true)
	want := []string{"mcp", "doctor", "srv", "--json"}
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("got %v want %v", *got, want)
	}
}

func TestLeaderProfileStart_Options(t *testing.T) {
	got := captureArgs(t)
	c := NewClient("grok")
	_ = c.LeaderProfileStart(context.Background(), 42, LeaderProfileStartOptions{Output: "/tmp/p", FrequencyHz: 200, JSON: true})
	want := []string{"leader", "profile", "start", "--pid", "42", "--output", "/tmp/p", "--frequency-hz", "200", "--json"}
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("got %v want %v", *got, want)
	}
}

func TestLeaderSocket_Threaded(t *testing.T) {
	cases := []struct {
		name string
		call func(*GrokClient)
		want []string
	}{
		{
			name: "subcommand",
			call: func(c *GrokClient) { _, _ = c.MCPList(context.Background()) },
			want: []string{"mcp", "list", "--leader-socket", "/tmp/x.sock"},
		},
		{
			name: "run prompt",
			call: func(c *GrokClient) { _, _ = c.RunPrompt("hi", &RunOptions{}) },
			want: []string{"-p", "hi", "--output-format", "plain", "--leader-socket", "/tmp/x.sock"},
		},
		{
			name: "stdin",
			call: func(c *GrokClient) { _, _ = c.RunFromStdin(strings.NewReader("hi"), "", &RunOptions{}) },
			want: []string{"--output-format", "plain", "--leader-socket", "/tmp/x.sock"},
		},
		{
			name: "login",
			call: func(c *GrokClient) { _ = c.Login(context.Background(), LoginOAuth) },
			want: []string{"login", "--oauth", "--leader-socket", "/tmp/x.sock"},
		},
		{
			name: "headless",
			call: func(c *GrokClient) { _ = c.RunHeadlessAgent(context.Background(), &HeadlessAgentConfig{}) },
			want: []string{"agent", "headless", "--leader-socket", "/tmp/x.sock"},
		},
		{
			name: "stdio",
			call: func(c *GrokClient) {
				s, _ := c.StartStdioAgent(context.Background(), nil)
				if s != nil {
					_ = s.Close()
				}
			},
			want: []string{"agent", "stdio", "--leader-socket", "/tmp/x.sock"},
		},
		{
			name: "serve",
			call: func(c *GrokClient) {
				sa, _ := c.StartServeAgent(context.Background(), &ServeAgentConfig{Bind: "127.0.0.1:0"})
				if sa != nil {
					_ = sa.Stop()
				}
			},
			want: []string{"agent", "serve", "--bind", "127.0.0.1:0", "--leader-socket", "/tmp/x.sock"},
		},
		{
			name: "leader agent",
			call: func(c *GrokClient) {
				la, _ := c.StartLeaderAgent(context.Background(), nil)
				if la != nil {
					_ = la.Stop()
				}
			},
			want: []string{"agent", "leader", "--leader-socket", "/tmp/x.sock"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := captureArgs(t)
			c := NewClient("grok")
			c.LeaderSocket = "/tmp/x.sock"
			tc.call(c)
			if !reflect.DeepEqual(*got, tc.want) {
				t.Fatalf("got %v want %v", *got, tc.want)
			}
		})
	}
}

func TestHeadlessArgs_SharedFlags(t *testing.T) {
	got := headlessArgs(&HeadlessAgentConfig{
		Reauth: true, Model: "grok-build", ReasoningEffort: "high",
		AlwaysApprove: true, AgentProfile: "/p.json", Leader: true,
		CLIChatProxyBase: "https://proxy", XAIAPIBase: "https://api",
		GrokWSURL: "wss://ws",
	})
	want := []string{
		"agent", "--reauth", "--model", "grok-build", "--reasoning-effort", "high",
		"--always-approve", "--agent-profile", "/p.json", "--leader",
		"--cli-chat-proxy-base-url", "https://proxy", "--xai-api-base-url", "https://api",
		"headless", "--grok-ws-url", "wss://ws",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v\nwant %v", got, want)
	}
}

func TestServeArgs_BindSecretRemote(t *testing.T) {
	got := serveArgs(&ServeAgentConfig{Bind: "127.0.0.1:9", Secret: "s", Remote: "r"})
	want := []string{"agent", "serve", "--bind", "127.0.0.1:9", "--secret", "s", "--remote", "r"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestLeaderAgentArgs(t *testing.T) {
	got := leaderAgentArgs(&LeaderAgentConfig{NoExitOnDisconnect: true, RelayOnDemand: true, NoAutoUpdate: true})
	want := []string{"agent", "leader", "--no-exit-on-disconnect", "--relay-on-demand", "--no-auto-update"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
