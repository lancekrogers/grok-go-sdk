# Dogfood Findings: v0.1.0

Run date: 2026-05-16
Consumer: `test/integration/` (in-repo integration suite acts as the first consumer until an external Go application is identified)
SDK commit: see git log of this branch
Notes: a real external consumer has not yet been wired; the integration suite plays the role of dogfood by exercising the public API surface.

## Test Results

| Lane | Pass | Fail | Skipped |
|------|------|------|---------|
| Unit (`go test ./pkg/grok/...`) | all | 0 | 0 |
| Integration mock (`go test -tags integration ./test/integration`) | all | 0 | 2 real-only |

## Findings

### Finding 1: real `grok agent stdio` is JSON-RPC 2.0, SDK envelope is simpler

- Severity: major
- Description: Real captures (see `test/testdata/stdio-agent/handshake.jsonl`) confirm the wire protocol is JSON-RPC 2.0 (`jsonrpc`/`method`/`params`/`id`/`result`/`error`). The SDK's `AgentMessage` is `{type,id,text,payload}`. Mock binary matches the simpler envelope; the SDK works against the mock but cannot drive a real `grok agent stdio` session end-to-end yet.
- Resolution: Documented as a known limitation in [RELEASE_NOTES_v0.1.0.md](RELEASE_NOTES_v0.1.0.md). Migration to JSON-RPC 2.0 envelope is scheduled as a follow-up festival. `RequestResponse` already correlates by `id`, so the public method signature stays stable across the migration.

### Finding 2: streaming event-type constants do not match real grok output

- Severity: minor
- Description: `pkg/grok/stream.go` exposes `EventAssistant`, `EventDelta`, `EventResult`, etc. Real captures show `thought`, `text`, `end`, `error`. The permissive `Raw` passthrough still works for unknown types; no consumer code breaks today, but the typed constants are misleading.
- Resolution: Documented as a known limitation in [RELEASE_NOTES_v0.1.0.md](RELEASE_NOTES_v0.1.0.md). Adding `EventText` and `EventEnd` constants and adjusting `Event.Text` to `Event.Data` is a non-breaking follow-up (existing fields remain via `omitempty`).

### Finding 3: no external Go consumer is currently wired

- Severity: minor (process)
- Description: The dogfood phase originally targeted an external Go application. None of the obey-campaign Go projects currently shell out to `grok`. The integration suite serves as a stand-in consumer.
- Resolution: Document the gap in this file and re-open the dogfood phase when an external consumer is identified. v0.1.0 release can proceed since the public API surface is exercised by the integration suite and the 10 example programs build clean.

## Summary

No blocking findings. v0.1.0 is ready to tag once the sequence's commit gate
clears.
