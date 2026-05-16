## Module

- Module path: `github.com/lancekrogers/grok-go-sdk`
- Go version: pinned at the latest stable release available at scaffold time;
  re-verify in [09-execution-plan.md](09-execution-plan.md).
- License: MIT.
- Visibility: private at creation, flipped public when the API has been
  exercised by an internal consumer and at least one external user.

## Package Layout

```
grok-go-sdk/
├── go.mod
├── go.sum
├── LICENSE
├── README.md
├── justfile
├── justfiles/
│   ├── build.just
│   ├── test.just
│   ├── lint.just
│   └── demo.just
├── pkg/
│   └── grok/
│       ├── grok.go               # Client, RunPrompt, RunPromptCtx
│       ├── options.go            # RunOptions, PreprocessOptions, BuildArgs
│       ├── result.go             # GrokResult, JSON decoding
│       ├── stream.go             # StreamPrompt, streaming-json parser
│       ├── stdin.go              # RunFromStdin
│       ├── stdio_agent.go        # `grok agent stdio` long-running transport
│       ├── headless_agent.go     # `grok agent headless` wrapper
│       ├── serve_agent.go        # `grok agent serve` lifecycle controller
│       ├── leader.go             # `grok leader` and `agent --leader` mode
│       ├── session.go            # Session ID generation, resume, fork, continue
│       ├── sessions.go           # `grok sessions list|search` wrapper
│       ├── share.go              # `grok share`
│       ├── trace.go              # `grok trace`
│       ├── import_session.go     # `grok import`
│       ├── inspect.go            # `grok inspect`
│       ├── models.go             # `grok models`
│       ├── memory.go             # `grok memory clear`
│       ├── worktree.go           # `grok worktree` admin + -w session worktree
│       ├── update.go             # `grok update --check`
│       ├── setup.go              # `grok setup`
│       ├── login.go              # `grok login` launcher
│       ├── mcp.go                # `grok mcp list|add|remove|doctor`
│       ├── permissions.go        # allow/deny rules, permission mode enums
│       ├── tools.go              # built-in tool registry + filtering
│       ├── sandbox.go            # sandbox profile resolution
│       ├── subagent.go           # SubagentConfig, --agents JSON
│       ├── plugin.go             # PluginManager + lifecycle hooks
│       ├── plugin_logging.go
│       ├── plugin_metrics.go
│       ├── plugin_audit.go
│       ├── plugin_filter.go
│       ├── budget.go             # BudgetTracker, BudgetConfig
│       ├── retry.go              # RetryPolicy, backoff, jitter
│       ├── errors.go             # GrokError, ParseError, ErrorType enum
│       ├── locate.go             # LocateBinary, default search paths
│       ├── jsonrpc.go            # shared helpers for stdio/serve JSON-RPC
│       └── doc.go                # package doc
├── pkg/
│   └── grok/
│       └── dangerous/
│           ├── doc.go
│           ├── dangerous.go      # --always-approve, bypassPermissions, sandbox=off
│           └── README.md         # CLAUDE_ENABLE_DANGEROUS-style guard
├── examples/
│   ├── basic/
│   ├── streaming/
│   ├── sessions/
│   ├── mcp/
│   ├── agent_stdio/
│   ├── worktree/
│   ├── best_of_n/
│   ├── check_loop/
│   ├── subagents/
│   └── agent_runtime/            # embedded long-running agent host
├── test/
│   ├── integration/              # build tag: integration, real grok binary
│   ├── mockserver/               # fake grok binary used by unit tests
│   └── testdata/                 # captured JSON and streaming-json fixtures
└── docs/
    ├── DEMOS.md
    ├── CONTRIBUTING.md
    ├── CLI_REFERENCE.md          # generated mirror of design doc 02
    └── RELEASE_NOTES_v0.1.0.md
```

The split between `pkg/grok` and `pkg/grok/dangerous` mirrors
`claude-code-go`'s `pkg/claude` and `pkg/claude/dangerous`. Anything that
materially weakens safety (always-approve, bypass permissions, sandbox
disabled) lives in `dangerous` behind an env-var guard.

## Internal vs External Boundary

Public API lives entirely under `pkg/grok` and `pkg/grok/dangerous`.

There is **no** top-level `internal/` package in the first release. The
surface is small and there is no shared helper that genuinely needs to be
hidden. If that changes (for example, if the streaming parser grows enough
state machine that it warrants test-only export), introduce
`pkg/grok/internal/` at that point.

## Dependency Policy

Standard library only for the first release, with these named exceptions:

- `github.com/google/uuid` for session ID generation. Same dependency that
  `claude-code-go` uses, well-maintained, single-purpose.

Everything else (JSON decoding, exec, channels, context) is stdlib.

New dependencies require a written justification (per root CLAUDE.md). The
SDK favors small, durable surfaces over rich third-party helpers.

## Concurrency Model

- The `Client` struct is safe for concurrent use across goroutines. Every
  call constructs its own `exec.Cmd` and owns its own pipes.
- The `BudgetTracker` is a value type with an internal mutex. Multiple
  callers can share one tracker across many `RunPrompt` calls.
- The `PluginManager` is similarly mutex-guarded. Hooks fire in registration
  order.
- Streaming events: `StreamPrompt` returns receive-only channels. The SDK
  closes both channels when the underlying process exits; the caller must
  drain them or call `ctx.Cancel()` to avoid leaking the reader goroutine.
- `agent stdio` returns a `StdioSession` with its own goroutines for
  read/write pumps. Closing the session shuts both pumps down deterministically.

## Error Model

All errors returned by public methods are either `*GrokError` (process,
auth, rate-limit, validation) or wrapped context errors. `GrokError`
implements `IsRetryable()` and exposes `Type` (`auth`, `rate_limit`,
`process`, `validation`, `transport`, `unknown`). Mirrors the
`claude-code-go` `ClaudeError`. See [04-api-surface.md](04-api-surface.md)
for the exact types.

## Exec Surface

A single `execCommand = exec.CommandContext` variable is exported in
`grok.go` so tests can swap it out. Every call goes through this hook,
including the stdio agent transport. The mock server in
`test/mockserver/main.go` is a real Go binary the test harness compiles and
points the SDK at, so we exercise the real `exec.CommandContext` path.

## Configuration Discovery

- The Client struct accepts `BinPath`. `NewClient("grok")` resolves via
  `LocateBinary()` if the path is bare.
- `LocateBinary()` order: `PATH` lookup first, then `~/.grok/bin/grok`,
  then a list of platform-specific fallbacks.
- The SDK does **not** read or write grok's `config.toml` directly. If a
  consumer needs to know what grok itself thinks the config says, they
  invoke `Client.Inspect(ctx)` which shells out to `grok inspect --json`
  and decodes the result.

## Lifecycle of a Headless Call

1. `RunPromptCtx` is invoked.
2. Options are cloned and run through `PreprocessOptions` (validates flag
   combinations, expands `Rules` vs `RulesFile`, normalizes permission rules).
3. `PluginManager.OnBeforeRun` hook fires.
4. `BuildArgs` constructs the argv.
5. `BudgetTracker` enforces pre-call budget; aborts if exceeded.
6. `exec.CommandContext` runs the binary with cwd, env, stdin.
7. Stdout is parsed per `OutputFormat`. Stderr is captured for error
   classification.
8. `PluginManager.OnAfterRun` hook fires with the parsed result.
9. Result returned to caller; errors classified through `ParseError`.

The streaming path replaces step 7 with a goroutine that scans
newline-delimited JSON and emits typed events.

## Cross-Reference to claude-code-go

Where a concept exists in `claude-code-go` and grok-go-sdk, keep names
aligned so a developer who knows one is immediately productive:

| claude-code-go | grok-go-sdk | Equivalence |
| --- | --- | --- |
| `ClaudeClient` | `GrokClient` | Identical role. |
| `ClaudeResult` | `GrokResult` | Different fields, same semantics. |
| `ClaudeError` | `GrokError` | Same type taxonomy. |
| `RunOptions` | `RunOptions` | Same name, grok-specific fields. |
| `StreamPrompt` | `StreamPrompt` | Same signature shape, different event types. |
| `BuildArgs` | `BuildArgs` | Same purpose. |
| `PreprocessOptions` | `PreprocessOptions` | Same purpose. |
| `dangerous` subpackage | `dangerous` subpackage | Same guard env-var pattern; env var renamed `GROK_ENABLE_DANGEROUS`. |

Where grok has a capability claude does not (`--best-of-n`, `--check`,
`--worktree`, `--sandbox`, `--no-plan`, `--no-subagents`,
`--experimental-memory`, `--verbatim`, `agent stdio`, leader processes), the
SDK adds fields without trying to invent parallels in `claude-code-go`.
