This document is the **complete** captured `grok --help` surface. The SDK
public API under `pkg/grok` must cover everything listed here that is callable
in non-interactive mode. When the CLI ships new flags, refresh this file first,
then update the Go option structs.

The reference binary lived at the standard `~/.grok/bin/grok` install path on
the capture date. The default model reported by `grok models` was `grok-build`.

---

## 1. Top-Level Invocation

```
grok [OPTIONS] [COMMAND]
```

With no subcommand and no `-p`, grok launches the interactive TUI. With `-p`,
`--prompt-file`, or `--prompt-json` it runs single-turn headless and exits.

### Top-level options

| Flag | Argument | Notes |
| --- | --- | --- |
| `--agent <NAME>` | string | Agent name or path to a definition file. |
| `--agents <JSON>` | JSON string | Inline subagent definitions. |
| `--allow <RULE>` | string, repeatable | Permission allow rule. Repeat to add more. |
| `--always-approve` | flag | Auto-approve all tool executions. Dangerous. |
| `--best-of-n <N>` | int | Run the task N ways in parallel and pick the best (headless only). |
| `-c, --continue` | flag | Continue the most recent session for the current cwd. |
| `--check` | flag | Append a self-verification loop to the prompt (headless only). |
| `--cwd <CWD>` | path | Working directory. |
| `--deny <RULE>` | string, repeatable | Permission deny rule. |
| `--disable-web-search` | flag | Disable web search and web fetch tools. |
| `--disallowed-tools <TOOLS>` | comma list | Built-in tools to remove. |
| `--effort <LEVEL>` | enum | `low`, `medium`, `high`, `xhigh`, `max`. |
| `--experimental-memory` | flag | Enable cross-session memory. |
| `-h, --help` | flag | Print help. |
| `-m, --model <MODEL>` | string | Model ID to use. |
| `--max-turns <N>` | int | Maximum number of agent turns. |
| `--no-alt-screen` | flag | Run inline instead of using the terminal alternate screen. |
| `--no-memory` | flag | Disable cross-session memory for this session. |
| `--no-plan` | flag | Disable plan mode. |
| `--no-subagents` | flag | Disable subagent spawning. |
| `--oauth` | flag | Use OAuth when the welcome screen starts authentication. |
| `--output-format <FORMAT>` | enum | `plain` (default), `json`, `streaming-json`. Headless mode. |
| `-p, --single <PROMPT>` | string | Single-turn prompt. Prints response to stdout and exits. |
| `--permission-mode <MODE>` | enum | `default`, `acceptEdits`, `auto`, `dontAsk`, `bypassPermissions`, `plan`. |
| `--prompt-file <PATH>` | path | Single-turn prompt loaded from a file. |
| `--prompt-json <JSON>` | JSON | Single-turn prompt as JSON content blocks. |
| `-r, --resume [<SESSION_ID>]` | optional string | Resume a session by ID, or most recent if omitted. |
| `--reasoning-effort <EFFORT>` | string | Reasoning effort for reasoning models. |
| `--restore-code` | flag | Check out the original session's commit when resuming. |
| `--rules <RULES>` | string or `@file` | Extra rules appended to the system prompt. `@file` loads from file. |
| `--sandbox <PROFILE>` | string, also env `GROK_SANDBOX` | Sandbox profile for filesystem and network access. |
| `--system-prompt-override <PROMPT>` | string | Override the agent's system prompt. |
| `--tools <TOOLS>` | comma list | Built-in tools to allow. |
| `-v, --version` | flag | Print version. |
| `--verbatim` | flag | Send the prompt exactly as given (no auto-prefixing). |
| `-w, --worktree [<WORKTREE>]` | optional string | Start the session in a new git worktree, optionally named. |

### Notes on top-level behavior

- `--allow` and `--deny` are repeatable. The Go SDK exposes them as
  `[]string` and emits one CLI arg per element.
- `--prompt-file` and `--prompt-json` are mutually exclusive with `-p`. The
  Go SDK validates this in `PreprocessOptions`.
- `--best-of-n`, `--check`, `--output-format` are documented as
  headless-only. The SDK rejects them if `-p` / `--prompt-file` /
  `--prompt-json` is not set.
- `--worktree` integrates with `grok worktree` state. The SDK exposes a
  `WorktreeName` option and a separate `Worktree` subpackage for the
  management subcommands.
- `--rules @file` is a syntactic shortcut. The SDK accepts `Rules string`
  and `RulesFile string` and assembles the correct CLI form.

---

## 2. Subcommands

### 2.1 `grok agent`

Run grok without the interactive UI. Top-level shared options:

| Flag | Notes |
| --- | --- |
| `--reauth`, `----reauthenticate` | Run authentication before starting the agent. |
| `-m, --model <MODEL>` | Same as top-level. |
| `--reasoning-effort <EFFORT>` | Same. |
| `--always-approve` | Same. |
| `--agent-profile <PATH>` | Path to an agent profile file. |
| `--leader` | Connect to a shared leader process instead of starting a new agent. Default per `[cli] use_leader` in config.toml. |
| `--no-leader` | Force start a new agent even with leader enabled in config. |
| `--grok-ws-origin <ORIGIN>` | Override the WS origin. |
| `--grok-ws-url <URL>` | Override the WS URL. |
| `--cli-chat-proxy-base-url <URL>` | Override the CLI chat proxy base URL. |
| `--xai-api-base-url <URL>` | Override the public xAI API base URL. |

Subcommands:

| Subcommand | Args | Purpose |
| --- | --- | --- |
| `stdio` | none | Run the agent over stdio. Full duplex agent protocol. |
| `headless` | `--grok-ws-origin`, `--grok-ws-url` | Run the agent headlessly over the Grok WebSocket relay. |
| `serve` | see below | Run the agent as a local WebSocket server. |
| `leader` | (no listed args) | Run as the shared leader process for other clients. |
| `help` | | Print help. |

`agent serve` options:

| Flag | Notes |
| --- | --- |
| `--bind <BIND>` | Listen address. Default `127.0.0.1:2419`. |
| `--secret <SECRET>` | Auth secret. Auto-generated if omitted. Env `GROK_AGENT_SECRET`. |
| `--remote <REMOTE>` | Remote agent URL for proxy mode. |
| `--grok-ws-origin`, `--grok-ws-url` | Same as headless. |

### 2.2 `grok import`

```
grok import [OPTIONS] [TARGETS]...
```

Import sessions into Grok.

| Flag | Notes |
| --- | --- |
| `--list` | List available sessions without importing. |
| `--json` | NDJSON output to stdout. |

Positional `TARGETS` are session IDs or `.jsonl` paths.

### 2.3 `grok inspect`

```
grok inspect [OPTIONS]
```

Show the configuration grok discovers for the current directory.

| Flag | Notes |
| --- | --- |
| `--json` | Machine-readable JSON. |

### 2.4 `grok leader`

Manage running leader processes.

| Subcommand | Purpose |
| --- | --- |
| `list` | List running leader processes. |
| `info` | Show details for a leader process. |
| `kill` | Stop all running leader processes. |
| `profile status` | Show profiling status. |
| `profile start` | Start CPU profiling. |
| `profile stop` | Stop CPU profiling. |

### 2.5 `grok login`

Sign in to Grok.

| Flag | Notes |
| --- | --- |
| `--oauth` | Use Grok OAuth via `auth.x.ai`. |
| `--device-auth` | Device-code auth for headless environments. |

### 2.6 `grok mcp`

Manage MCP server configurations.

| Subcommand | Args | Purpose |
| --- | --- | --- |
| `list` | none | List configured MCP servers. |
| `add` | `<NAME>` plus flags | Add or update an MCP server configuration. |
| `remove` | `<NAME>` | Remove an MCP server configuration. |
| `doctor` | none | Diagnose MCP server configuration and connectivity. |

`mcp add` flags:

| Flag | Notes |
| --- | --- |
| `--command <COMMAND>` | Command to run for stdio transport. |
| `--args <ARGS>...` | Variadic args. |
| `--env <KEY=VALUE>...` | Variadic environment overrides. |
| `--url <URL>` | Server URL for HTTP or SSE transport. |
| `--type <TRANSPORT_TYPE>` | Transport type for HTTP servers. |

### 2.7 `grok memory`

Manage cross-session memory.

| Subcommand | Flag | Notes |
| --- | --- | --- |
| `clear` | `--workspace` | Clear workspace memory (`MEMORY.md`, `sessions/`, `index.sqlite`). |
| `clear` | `--global` | Clear global `MEMORY.md`. |
| `clear` | `--all` | Clear both. |
| `clear` | `-y, --yes` | Skip confirmation prompt. |

### 2.8 `grok models`

```
grok models
```

List available models and exit. No flags. Returns plain text such as
`Default model: grok-build` followed by an enumerated list.

### 2.9 `grok sessions`

| Subcommand | Args | Purpose |
| --- | --- | --- |
| `list` | none | List recent sessions. Same as `search` with no query. |
| `search` | `<QUERY>` plus `-n, --limit` | Search sessions by keyword. Default limit 20. |

### 2.10 `grok setup`

```
grok setup
```

Fetch and install managed deployment configuration. No flags.

### 2.11 `grok share`

```
grok share <SESSION_ID>
```

Share a session and print the share URL.

### 2.12 `grok ssh`

```
grok ssh <SSH_ARG>...
```

Run `ssh` with local clipboard support. Arguments are forwarded to ssh
unchanged. Out of scope for the SDK in the first release except as a passthrough
helper.

### 2.13 `grok trace`

```
grok trace [OPTIONS] <SESSION_ID>
```

Export or upload session trace data.

| Flag | Notes |
| --- | --- |
| `--local` | Save locally only, skip remote upload. |
| `-o, --output <OUTPUT>` | Output path. Default `~/.grok/trace-exports/<session-id>.tar.gz`. |
| `--json` | Machine-readable JSON. |

### 2.14 `grok update`

Check for updates or install a specific version.

| Flag | Notes |
| --- | --- |
| `--check` | Check without installing. |
| `--json` | Machine-readable output for `--check`. |
| `--force-reinstall` | Re-download even if up to date. |
| `--version <VERSION>` | Install a specific version. |
| `--alpha` | Switch to alpha channel. |
| `--stable` | Switch to stable channel (default). |

### 2.15 `grok worktree`

Manage git worktrees that grok itself tracks. Distinct from the top-level
`-w/--worktree` flag, which spawns a session in a new worktree.

| Subcommand | Args | Purpose |
| --- | --- | --- |
| `list` | none | List tracked worktrees. |
| `show` | (id) | Show details. |
| `rm` | `<IDS>...`, `-f`, `--dry-run` | Remove worktrees. |
| `gc` | `--dry-run`, `--max-age`, `-f` | Garbage-collect orphaned or stale worktrees. |
| `db rebuild` | none | Rebuild DB from filesystem scan. |
| `db stats` | none | Show DB statistics. |
| `db path` | none | Print DB file path. |

### 2.16 `grok version`, `grok help`

`grok version` / `grok v`: print version info. `grok help [SUBCOMMAND]`:
print help for a given subcommand.

---

## 3. Headless Output Formats

The SDK consumes three `--output-format` shapes.

### 3.1 `plain`

Stdout is the raw response text. Stderr may contain warnings or error
diagnostics (the captured sample contained a transient
`worker quit with fatal: Transport channel closed, when Auth(AuthorizationRequired)`
diagnostic that did not block exit). Exit code zero on success.

### 3.2 `json`

A single JSON object on stdout. Captured shape (2026-05-15, `grok-build`,
`grok -p "say hello" --output-format json`):

```json
{
  "text": "Hello! How can I help you today?",
  "stopReason": "EndTurn",
  "sessionId": "019e2a92-d8bb-7b12-acc9-b48719ec2354",
  "requestId": "b8be7af6-27b8-4ed5-a529-4932d90dff58",
  "thought": "The user said \"say hello\". This is a simple greeting request. I should respond in a friendly, direct way.\n"
}
```

Fields the SDK must decode:

- `text` final assistant text
- `stopReason` enum tag (`EndTurn`, etc.)
- `sessionId` UUID-like string
- `requestId` UUID-like string
- `thought` reasoning text when reasoning is enabled (may be absent)

Open question: cost, token, duration fields are not in the captured sample.
The SDK structs include optional fields for these and `omitempty` round-trips
gracefully. Confirm against future `--include-...` style flags before adding
new required behavior.

### 3.3 `streaming-json`

Newline-delimited JSON events. The SDK exposes a Go channel of typed events
plus a parallel error channel, matching the `claude-code-go` `StreamPrompt`
signature. Exact event taxonomy must be captured against a real session
during implementation (open question logged).

---

## 4. Environment Variables and Files

| Name | Role |
| --- | --- |
| `GROK_SANDBOX` | Default sandbox profile if `--sandbox` not supplied. |
| `GROK_AGENT_SECRET` | Default auth secret for `agent serve`. |
| `~/.grok/` | Per-user state. |
| `~/.grok/bin/grok` | Default binary location used by the host install. |
| `~/.grok/trace-exports/` | Default trace export directory. |
| `config.toml` | Per-user CLI config, includes `[cli] use_leader`. |
| `MEMORY.md` (workspace and global) | Cross-session memory store. |

The SDK exposes `LocateBinary()` that searches `PATH` then `~/.grok/bin/grok`,
matching how a user-installed grok lands on disk.

---

## 5. Refresh Procedure

When `grok` ships a new minor version:

1. Re-run `grok --help`, `grok help <each subcommand>`, and `grok help agent
   <each subsubcommand>`. Diff against this document.
2. Re-run `grok -p "say hello" --output-format json` and capture stdout to
   confirm field set. Re-run with `--output-format streaming-json` and save
   the first three events as a fixture in the implementation repo's
   `testdata/`.
3. Update this file. Bump the captured version header.
4. Update `pkg/grok/options.go`, wrappers, tests, and README examples only
   after this file is accurate.
