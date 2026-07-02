# Release v0.2.0

CLI alignment release for the current Grok Build CLI surface.

## Breaking Changes

- Removed dead `RunOptions` fields that emitted CLI flags no longer accepted by `grok`: `InputFormat`, `MCPConfigPath`, `MCPConfigs`, and `StrictMCPConfig`.
- Removed the `InputFormat` type/constants because the CLI no longer has `--input-format`.
- Removed the `Share` wrapper because `grok share` is no longer present.
- Remodeled `MCPServerConfig` and `MCPAdd` for the positional `grok mcp add` interface.

## Migration

- Replace per-run MCP config fields with persistent CLI configuration:
  `grok mcp add -s project <name> <command-or-url> -- <args...>`.
- Replace `Share` usage with `Export(ctx, sessionID, opts)` for Markdown session artifacts.
- Use `OutputFormat` only for output selection: `plain`, `json`, or `streaming-json`.

## CLI Alignment

- Retagged `InspectReport` for camelCase `inspect --json` payloads and preserved structured `source` fields losslessly.
- Wired `--session-id`, `--fork-session`, `--json-schema`, and `--worktree-ref`.
- Added wrappers for `sessions delete`, `logout`, `export`, `completions`, and the `plugin` subcommand tree.
- Filled wrapper gaps for MCP remove/doctor, headless agent shared flags, agent serve remote binding, agent leader flags, and leader profile options.
- Updated built-in tool constants and sandbox profile handling against the live CLI.
- Fixed `worktree gc --max-age` formatting to emit single-unit durations such as `72h`.

## Drift Hardening

- Refreshed CLI reference docs and live output fixtures.
- Added fixture-backed decode/parser coverage for result, stream, inspect, models, MCP, sessions, leader, and version output.
- Added `just cli-drift` and `just cli-drift-update`, plus a non-blocking CI drift check.

## Verification

- P0 grep guards for removed flags and `share` are clean.
- `just test all` passes.
- `just lint` passes.
- `go test ./...` passes.
- `just cli-drift` passes.
- Live smoke checks passed with installed `grok 0.2.81 (d3b6135814b2) [stable]`: `grok version --json`, isolated `grok mcp add -s project`/`grok mcp list`, and SDK `Inspect` decode.
