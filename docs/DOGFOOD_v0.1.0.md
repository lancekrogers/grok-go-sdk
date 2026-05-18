# Dogfood Findings: v0.1.0

Run date: 2026-05-16
Grok CLI version: `grok 0.1.211`
SDK commit: pre-v0.1.0 release candidate
Consumer: `test/integration/` plus the example programs under `examples/`. Examples are the verification surface, following the model used by `claude-code-go`.

## Test Results

| Lane | Pass | Fail | Skipped |
|------|------|------|---------|
| Unit (`go test -race ./pkg/grok/...`) | all | 0 | 0 |
| Integration mock (`go test -tags integration ./test/integration`) | all | 0 | real-only and slow-only |
| Integration real (`INTEGRATION_REAL=1 go test -tags integration ./test/integration`) | 11/11 | 0 | slow-only (`INTEGRATION_SLOW=1`) |

## Resolved Findings

### Finding 1 (resolved): real `grok agent stdio` is JSON-RPC 2.0

Status: FIXED in this branch.

Real `grok agent stdio` speaks the Agent Client Protocol (ACP) over JSON-RPC 2.0. The previous SDK envelope `{type,id,text}` was fictional and would have returned `-32601 Method not found` against real grok. `pkg/grok/acp.go` and the rewritten `pkg/grok/stdio_agent.go` now expose typed `Initialize`/`Authenticate`/`NewSession`/`Prompt`/`PromptText` helpers and an `Updates()` channel for streaming `session/update` notifications. `Call`/`Notify` are available as low-level escape hatches. `TestStdioAgent_FullTurn` exercises a full handshake against grok 0.1.211 and asserts text chunks accumulate plus a final result with token usage.

### Finding 2 (resolved): streaming event-type constants do not match real grok output

Status: FIXED in PR #2 (`fix(stream): match real grok wire format`).

`EventText`, `EventEnd` are now exposed alongside `EventThought` and `EventError`. `Event` carries `Data`, `StopReason`, `RequestID`, `Message`. `Event.Content()` and `Event.IsTerminal()` helpers added. Presumed-but-unobserved types stay defined for forward compatibility via `Event.Raw` passthrough. Verified end-to-end via `TestStreamPrompt_TextAndEnd` against grok 0.1.211.

### Finding 3 (resolved): `--check` and `--best-of-n` validation rejected RunPromptCtx callers

Status: FIXED in this branch.

`PreprocessOptions` validated `opts.Prompt != ""` to satisfy the `--check` and `--best-of-n` prompt-source requirements, but `RunPromptCtx(ctx, prompt, opts)` passes the prompt as a function argument, not via `opts.Prompt`. Added `prepareOptionsWithPrompt` so RunPromptCtx, StreamPrompt, and RunFromStdinCtx pre-set `opts.Prompt` for validation and restore it afterward so BuildArgs continues to use the function-arg pathway. `TestCheckLoop_Validation` and `TestBestOfN_Validation` lock in the regression.

### Finding 4 (resolved): `grok mcp doctor` exits non-zero with a useful stdout report

Status: FIXED in this branch.

`grok mcp doctor` writes a per-server diagnostic to stdout and exits 1 when any server fails its handshake. The previous `MCPDoctor` discarded the report. Switched to `runSubcommandTolerant`, which returns stdout regardless of exit code. `MCPDoctor` now returns `(report, error)` so callers can render the report alongside the failure.

## Outstanding (not blocking v0.1.0)

- `Check` and `BestOfN` real-lane integration tests are gated behind `INTEGRATION_SLOW=1` since each spawns multiple full grok passes per call.
- No external Go consumer has wired the SDK yet. The integration suite + 10 example programs cover the public surface end-to-end. Re-open the dogfood phase once an external consumer is identified.
- Stdio agent server-initiated requests (`fs/read_text_file`, `terminal/*`, etc.) are auto-responded with `-32601 Method not found`. Clients that want to expose filesystem or terminal capabilities to the agent need a follow-up to register handlers.

## Summary

No blocking findings. The SDK drives every documented surface end-to-end against `grok 0.1.211`. Ready to re-tag `v0.1.0`.
