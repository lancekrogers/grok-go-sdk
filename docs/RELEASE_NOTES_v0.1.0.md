# Release v0.1.0

Initial public release.

## Scope

- Generic Go SDK wrapping the `grok` CLI.
- MIT licensed.
- Captured grok version: `grok 0.1.211` (see CLI_REFERENCE.md for the full surface).
- Headless `-p`, `streaming-json`, and `agent stdio` transports, all verified end-to-end against real grok.
- Sessions: continue, resume, share, trace, import.
- Admin subcommands: MCP, worktree, memory, update, setup, login, inspect, models, leader.
- Ergonomics: plugin manager, budget tracker, retry policy, structured errors.
- 10 examples including `agent_runtime` (embedded long-running agent host pattern).
- `dangerous` subpackage with `GROK_ENABLE_DANGEROUS` env guard.

## Transport Notes

- `grok agent stdio` speaks the Agent Client Protocol (ACP) over JSON-RPC 2.0. `StartStdioAgent` exposes typed helpers (`Initialize`, `Authenticate`, `NewSession`, `Prompt`, `PromptText`) and an `Updates()` channel for streaming `session/update` notifications. `Call`/`Notify` are available as low-level escape hatches.
- `--output-format streaming-json` emits `thought`, `text`, `end`, and `error` events with `data`/`stopReason`/`sessionId` fields. The SDK exposes them as `EventText`, `EventThought`, `EventEnd`, `EventError`. Presumed-but-unobserved types (`assistant`, `delta`, `result`, etc.) remain defined for forward compatibility via the permissive `Event.Raw` passthrough.
- `grok mcp doctor` exits non-zero when one or more MCP servers fail their handshake but still emits a useful diagnostic report on stdout. `MCPDoctor` returns the report alongside the error so callers can surface both.

## Known Limitations

- Fork session (`--fork-session`) not yet supported by grok; SDK pre-plumbs the field with a validation error.
- Custom `--session-id` not honored by grok; SDK field present for forward compatibility.
- `Check` (`--check`) and `BestOfN` (`--best-of-n`) run multiple grok passes per call, so they are slow and expensive. The integration tests for these are gated behind `INTEGRATION_SLOW=1`.
- Plain-text admin command parsers (`SessionsList`, `MCPList`, etc.) tolerate missing columns and non-table outputs, but may need updates after grok releases that change the shape of these outputs.
- API is pre-v1.0; breaking changes are allowed until v1.0.

## Verification

- `go test -race ./pkg/grok/...` passes.
- `go test -tags integration ./test/integration` (mock lane) passes.
- `INTEGRATION_REAL=1 go test -tags integration ./test/integration` passes 11/11 against `grok 0.1.211`, covering basic, streaming, stdio agent (full ACP turn), sessions multi-turn, MCP list, MCP doctor, sessions list, subagents, worktree create/list/remove, and validation paths for `Check`/`BestOfN`.

## Acknowledgments

- xAI for the grok CLI.
- The Go community for tooling and `go.uber.org/goleak`, `github.com/google/uuid`.
