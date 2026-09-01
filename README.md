<h1 align="center">Grok Go SDK</h1>

<p align="center">
  <a href="https://github.com/lancekrogers/grok-go-sdk/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/lancekrogers/grok-go-sdk/ci.yml?branch=main&label=CI"></a>
  <a href="https://pkg.go.dev/github.com/lancekrogers/grok-go-sdk"><img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/lancekrogers/grok-go-sdk.svg"></a>
  <img alt="License: MIT" src="https://img.shields.io/badge/license-MIT-blue.svg">
</p>

A Go library for programmatically integrating the [grok](https://github.com/xai-org/grok)
CLI in Go applications. Wraps the headless and stdio-agent surfaces of the
`grok` binary so Go programs can drive Grok Build sessions, stream output,
manage sessions, attach MCP servers, and enforce permissions, budgets, and
retry policies.

> This SDK is in early development. Public API is unstable until `v1.0`.

## Highlights

- Idiomatic Go wrapper for every `grok` subcommand we ship today.
- Streaming, stdio-agent, sessions, MCP, worktree, admin commands.
- Plugin manager with logging, metrics, audit, and tool-filter hooks.
- Budget tracker with warning and exceeded callbacks.
- Retry policy with exponential backoff and rate-limit honor.
- `pkg/grok/dangerous` for opt-in unsafe operations behind an env guard.

## Features

### Core
- `RunPrompt` / `RunPromptCtx`: single-shot prompts with `json` or `plain`.
- `StreamPrompt`: newline-delimited event channel with goleak-clean cancel.
- `RunFromStdin` / `RunFromStdinCtx`: pipe arbitrary stdin into grok.
- Session helpers: `GenerateSessionID` (UUID v7), `Continue`, `Resume`.

### Advanced
- `StartStdioAgent`: full-duplex JSON-RPC-style session with `RequestResponse`.
- `RunHeadlessAgent`, `StartServeAgent`, `StartLeaderAgent`, `LeaderList`/`Info`/`Kill`/`Profile*`.
- Admin subcommand wrappers: `MCP*`, `Worktree*`, `Memory*`, `Update*`, `Setup`, `Login`, `Logout`, `Inspect`, `Models`, `Export`, `Trace`, `Import`, `Completions`, `SessionsList`/`Search`/`Delete`.
- `grok plugin` tree: `Plugin*` (install/uninstall/enable/…) and `Marketplace*` wrappers.
- Permission rule validation + built-in tool name constants.
- Sandbox-profile resolution + permissive heuristic.
- Subagent JSON marshaling.

### Developer Experience
- `PluginManager` with `OnBeforeRun` and `OnAfterRun` hooks.
- Four built-in plugins: `LoggingPlugin`, `MetricsPlugin`, `AuditPlugin`, `ToolFilterPlugin`.
- `BudgetTracker` with warning/exceeded callbacks.
- `RetryPolicy` with backoff, jitter, and rate-limit honor.

## Out of scope

`grok wrap` (a PTY/clipboard passthrough that replaced the old `grok ssh`) and
`grok dashboard` (an interactive TUI) are interactive, terminal-bound commands.
They are intentionally **not** wrapped by this headless SDK — run the `grok`
binary directly for those. Shell completion scripts are available via
`Completions(ctx, shell)`.

## Installation

```
go get github.com/lancekrogers/grok-go-sdk@latest
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func main() {
	c, err := grok.NewClientFromPath()
	if err != nil {
		log.Fatal(err)
	}
	res, err := c.RunPrompt("Write a one-line Go function that returns 42", &grok.RunOptions{
		Format: grok.JSONOutput,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Text)
}
```

## Prerequisites

- Go 1.24 or newer.
- The [grok CLI](https://github.com/xai-org/grok) installed and authenticated via `grok login`.
- xAI API access for the real-binary lane (mock binary works offline).
- [just](https://github.com/casey/just) for running project commands.

## Interactive Demos

Every example under `examples/` is runnable. See [docs/DEMOS.md](docs/DEMOS.md)
for the full list. Common starting points:

```
just demo basic
just demo streaming
just demo sessions
just demo agent-stdio
```

## Core Features

### Basic prompts

```go
res, err := c.RunPrompt("Hello", &grok.RunOptions{Format: grok.JSONOutput})
```

### Streaming

```go
events, errs := c.StreamPrompt(ctx, "Explain Go channels", &grok.RunOptions{})
for ev := range events {
	fmt.Print(ev.Text)
}
```

### Sessions

```go
first, _ := c.RunPromptCtx(ctx, "Remember the number 42.", opts)
again, _ := c.ResumeConversationCtx(ctx, "What number?", first.SessionID)
```

### MCP

```go
servers, _ := c.MCPList(ctx)
report, _ := c.MCPDoctor(ctx)
```

## Advanced Features

### Plugins

```go
pm := grok.NewPluginManager()
pm.Register(&grok.LoggingPlugin{SanitizeSecrets: true}, nil)
pm.Register(&grok.MetricsPlugin{}, nil)
opts := &grok.RunOptions{PluginManager: pm}
```

### Budget tracking

```go
bt := grok.NewBudgetTracker(&grok.BudgetConfig{MaxBudgetUSD: 5.00, WarningThreshold: 0.8})
opts := &grok.RunOptions{BudgetTracker: bt}
```

### Subagents

```go
opts := &grok.RunOptions{
	Agent:  "security",
	Agents: map[string]*grok.SubagentConfig{
		"security": {Description: "...", Prompt: "...", Tools: []string{grok.ToolRead}},
	},
}
```

### Retry

```go
res, err := c.RunPromptWithRetryCtx(ctx, "hello", opts, grok.DefaultRetryPolicy())
```

### Permissions

```go
opts := &grok.RunOptions{
	AllowRules: []string{grok.BuildToolGlob(grok.ToolBash, "git status*")},
	DenyRules:  []string{grok.BuildToolGlob(grok.ToolWebFetch, "*")},
}
```

## API Reference

### RunOptions (selected)

| Field | Purpose |
|-------|---------|
| `Format` | `JSONOutput`, `StreamingJSONOutput`, `PlainOutput` |
| `Model` | model identifier (e.g. `grok-build`) |
| `WorkingDirectory` | overrides the spawned `grok` process cwd |
| `Timeout` | per-call timeout (uses `context.WithTimeout`) |
| `Continue`, `ResumeID` | session controls |
| `Agent`, `Agents` | subagent selection + inline definitions |
| `AllowRules`, `DenyRules` | permission rule lists |
| `BudgetTracker`, `PluginManager` | run-time hooks |
| `MaxBudgetUSD` | per-call budget cap |
| `Worktree` | `-w` worktree enablement + name |
| `SandboxProfile` | sandbox profile name |
| `Check`, `BestOfN` | self-verification + n-way race |
| `AllowDangerousMode` | opt-in for bypass-permissions or always-approve |

### Core methods

| Method | Purpose |
|--------|---------|
| `NewClient`, `NewClientFromPath` | construct a `*GrokClient` |
| `RunPrompt`, `RunPromptCtx` | single-shot prompt |
| `StreamPrompt` | streaming events |
| `RunFromStdin`, `RunFromStdinCtx` | pipe stdin into grok |
| `StartStdioAgent` | long-lived agent session |
| `RunPromptWithRetry*` | retry-wrapped single-shot |
| `ContinueConversation*`, `ResumeConversation*` | session helpers |
| `MCP*`, `Worktree*`, `Memory*`, `Update*`, `Setup`, `Login`, `Logout`, `Inspect`, `Models`, `Export`, `Trace`, `Import`, `Completions`, `Sessions*`, `Plugin*`, `Marketplace*`, `Leader*` | admin wrappers |

## Security-Sensitive Features

The `pkg/grok/dangerous` subpackage exposes operations that weaken safety
controls (bypass permissions, disable sandbox, auto-approve all tool calls).
They require explicit opt-in:

```
export GROK_ENABLE_DANGEROUS="i-accept-all-risks"
```

`NewDangerousClient` returns an error if the env var is missing or if
`GO_ENV` or `NODE_ENV` is `production`. See
[pkg/grok/dangerous/README.md](pkg/grok/dangerous/README.md).

## Testing

```
just test all                # all unit tests
just test integration        # integration tests against the mock binary
just test integration-real   # real-binary lane
```

## Development

```
just build all      # library + examples + mock
just lint           # fmt + vet
just mock build     # rebuild the test mock binary
just demo basic     # run an example
```

## Documentation

- [docs/DEMOS.md](docs/DEMOS.md) - every example
- [docs/CLI_REFERENCE.md](docs/CLI_REFERENCE.md) - every flag the SDK emits
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - subprocess/IPC architecture
- [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) - dev environment + PR norms
- [docs/RELEASE_NOTES_v0.2.0.md](docs/RELEASE_NOTES_v0.2.0.md) - changelog

## Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md).

## Related

Go SDKs for other coding-agent CLIs:

- [claude-code-go](https://github.com/lancekrogers/claude-code-go) wraps `claude`
- [vercel-fx-go](https://github.com/Obedience-Corp/vercel-fx-go) wraps `fx`
- [cursor-agent-go](https://github.com/Obedience-Corp/cursor-agent-go) wraps `cursor-agent`

## License

MIT. See [LICENSE](LICENSE).

## Acknowledgments

Built on top of the [grok CLI](https://github.com/xai-org/grok) by xAI.
