# Release v0.1.0

Initial private release.

## Scope

- Generic Go SDK wrapping the `grok` CLI
- MIT licensed
- Captured grok version: `grok 0.1.210 (8b63e9068c)` (see CLI_REFERENCE.md for the full surface)
- Headless `-p`, `streaming-json`, and `agent stdio` transports
- Sessions: continue, resume, share, trace, import
- Admin subcommands: MCP, worktree, memory, update, setup, login, inspect, models, leader
- Ergonomics: plugin manager, budget tracker, retry policy, structured errors
- 10 examples including `agent_runtime` (embedded long-running agent host pattern)
- `dangerous` subpackage with `GROK_ENABLE_DANGEROUS` env guard

## Known Limitations

- Fork session (`--fork-session`) not yet supported by grok; SDK pre-plumbs the field with a validation error.
- Custom `--session-id` not honored by grok; SDK field present for forward compatibility.
- Real `grok agent stdio` uses JSON-RPC 2.0; the SDK's `AgentMessage` envelope is a simpler shape that works against the mock binary today. Migration to JSON-RPC 2.0 envelope is a planned follow-up; see design doc 10 Q8.
- Streaming event taxonomy in `pkg/grok/stream.go` lists the presumed event types from design doc 04. Real captures show `thought`, `text`, `end`, `error`. The permissive `Raw` passthrough means unknown types still round-trip safely; constants are scheduled to be added in a follow-up.
- Some grok subcommand outputs are plain-text only; parsers tolerate missing columns but may need updates after grok releases new versions.
- API is pre-v1.0; breaking changes are allowed until v1.0.

## Acknowledgments

- xAI for the grok CLI.
- The Go community for tooling and `go.uber.org/goleak`, `github.com/google/uuid`.
