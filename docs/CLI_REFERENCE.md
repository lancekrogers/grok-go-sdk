# grok CLI Reference

This document is the captured `grok --help` surface that the SDK targets. When
the CLI ships new flags, refresh this file first, then update the Go option
structs and wrappers.

**Captured against:** `grok 0.2.77 (44e77bec3af6) [stable]`
Default model reported by `grok models`: `grok-build` (also `grok-composer-2.5-fast`).

> Changes since the previous capture: `grok share` and `grok ssh` were removed
> (`ssh` → `wrap`); `mcp add` moved to a positional interface; `inspect --json`
> switched to camelCase keys; and `plugin`, `wrap`, `dashboard`, `logout`,
> `completions`, `export` were added, along with top-level `--json-schema`,
> `--session-id`, `--fork-session`, `--worktree-ref`, and the global
> `--leader-socket`.

---

## 0. Global flags

Every subcommand accepts:

| Flag | Notes |
| --- | --- |
| `--debug` | Enable debug logging. |
| `--debug-file <FILE>` | Write debug logs to FILE. |
| `--leader-socket <PATH>` | Talk to a specific leader (`~/.grok/leader.sock` default). SDK: `GrokClient.LeaderSocket`. |

---

## 1. Top-level invocation

```
grok [OPTIONS] [PROMPT] [COMMAND]
```

No subcommand + no `-p` → interactive TUI. With `-p` / `--prompt-file` /
`--prompt-json` → single-turn headless. The SDK builds this argv in
`BuildArgs` (`RunPrompt`, `RunFromStdin`, `StreamPrompt`).

| Flag | SDK `RunOptions` field |
| --- | --- |
| `--agent <NAME>` | `Agent` / `AgentDefinitionFile` |
| `--agents <JSON>` | `Agents` / `AgentsJSON` |
| `--allow <RULE>` (repeatable) | `AllowRules` |
| `--always-approve` | `AlwaysApprove` (gated by `AllowDangerousMode`) |
| `--best-of-n <N>` | `BestOfN` (headless only) |
| `-c, --continue` | `Continue` |
| `--check` | `Check` (headless only) |
| `--cwd <CWD>` | `WorkingDirectory` |
| `--deny <RULE>` (repeatable) | `DenyRules` |
| `--disable-web-search` | `DisableWebSearch` |
| `--disallowed-tools <CSV>` | `DisallowedTools` |
| `--effort <LEVEL>` | `Effort` (`low`/`medium`/`high`/`xhigh`/`max`) |
| `--experimental-memory` | `ExperimentalMemory` |
| `--fork-session` | `ForkSession` (only with resume/continue) |
| `--json-schema <SCHEMA>` | `JSONSchema` (implies `--output-format json`) |
| `-m, --model <MODEL>` | `Model` |
| `--max-turns <N>` | `MaxTurns` |
| `--no-alt-screen` | `NoAltScreen` |
| `--no-memory` | `NoMemory` |
| `--no-plan` | `NoPlan` |
| `--no-subagents` | `NoSubagents` |
| `--oauth` | `OAuth` |
| `--output-format <FORMAT>` | `Format` (`plain`/`json`/`streaming-json`) |
| `-p, --single <PROMPT>` | `Prompt` |
| `--permission-mode <MODE>` | `PermissionMode` (`default`/`acceptEdits`/`auto`/`dontAsk`/`bypassPermissions`/`plan`) |
| `--prompt-file <PATH>` | `PromptFile` |
| `--prompt-json <JSON>` | `PromptJSON` |
| `-r, --resume [<SESSION_ID>]` | `ResumeID` |
| `--reasoning-effort <EFFORT>` | `ReasoningEffort` |
| `--restore-code` | `RestoreCode` (requires `ResumeID`) |
| `--rules <RULES>` | `Rules` / `RulesFile` (`@file`) |
| `--sandbox <PROFILE>` | `SandboxProfile` (env `GROK_SANDBOX`) |
| `-s, --session-id <UUID>` | `SessionID` |
| `--system-prompt-override <PROMPT>` | `SystemPromptOverride` |
| `--tools <CSV>` | `AllowedTools` |
| `--verbatim` | `Verbatim` |
| `-w, --worktree [<WORKTREE>]` | `Worktree{Enabled,Name}` |
| `--worktree-ref <REF>` (alias `--ref`) | `Worktree.Ref` |

> There is **no** `--input-format`, `--mcp-config`, or `--strict-mcp-config`
> flag (grok rejects them). MCP servers are configured via `grok mcp add` /
> `config.toml`.

---

## 2. `grok agent` — run without the interactive UI

Shared `agent` options (before the subcommand): `--reauth`, `-m/--model`,
`--reasoning-effort`, `--always-approve`, `--agent-profile`,
`--leader`/`--no-leader`, `--cli-chat-proxy-base-url`, `--xai-api-base-url`.
(SDK: `HeadlessAgentConfig` fields.)

| Subcommand | Args | SDK |
| --- | --- | --- |
| `stdio` | — | `StartStdioAgent` |
| `headless` | `--grok-ws-origin`, `--grok-ws-url` | `RunHeadlessAgent` |
| `serve` | `--bind` (default `127.0.0.1:2419`), `--secret` (env `GROK_AGENT_SECRET`), `--remote`, `--grok-ws-*` | `StartServeAgent` |
| `leader` | `--no-exit-on-disconnect`, `--relay-on-demand`, `--no-auto-update`, `--grok-ws-*` | `StartLeaderAgent` |

---

## 3. Subcommands

### `grok leader` — manage leader processes
`list` (`--json`), `info` (`--pid`, `--json`), `kill`, `profile {status,start,stop}`
(`--pid`; `start` also `--output`, `--frequency-hz`, `--json`).
SDK: `LeaderList`, `LeaderInfoCmd`, `LeaderKill`, `LeaderProfile{Start,Stop,Status}`.

### `grok mcp` — MCP servers
```
grok mcp add [OPTIONS] <NAME> [COMMAND_OR_URL] [ARGS]...
  -t/--transport <stdio|http|sse>   -s/--scope <user|project>
  -e/--env <KEY=value> (repeat)     -H/--header <NAME: VALUE> (repeat)
```
`list`, `remove <NAME>` (`-s/--scope`), `doctor [NAME]` (`--json`).
SDK: `MCPList`, `MCPAdd`, `MCPRemove`, `MCPDoctor` (`MCPServerConfig`, `MCPScope`).

### `grok memory`
`clear` (`--workspace`, `--global`, `--all`, `-y`). SDK: `MemoryClear`.

### `grok worktree`
`list`, `show <ID_OR_PATH>`, `rm <IDS>...` (`-f`, `--dry-run`), `gc`
(`--dry-run`, `--max-age`, `-f`), `db {rebuild,stats,path}`.
SDK: `WorktreeList/Show/Remove/GC/DB*`.

### `grok sessions`
`list` (`-n`), `search <QUERY>` (`-n`), `delete <ID>`.
SDK: `SessionsList`, `SessionsSearch`, `SessionsDelete`.

### `grok plugin`
`list`, `install <SOURCE>` (`--trust`), `uninstall <NAME>` (`--confirm`,
`--keep-data`), `update [NAME]`, `enable <NAME>`, `disable <NAME>`,
`details <NAME>`, `validate [PATH]`, `tag [PATH]` (`--push`, `-f`, `--dry-run`),
`marketplace {list, add <URL>, remove <URL>, update [NAME]}`.
SDK: `Plugin*` and `Marketplace*` methods (in `plugin_cli.go`; distinct from the
in-process `PluginManager`).

### Simple subcommands
| Command | Args | SDK |
| --- | --- | --- |
| `import [TARGETS]...` | `--list`, `--json` | `Import` |
| `inspect` | `--json` (camelCase payload) | `Inspect` → `InspectReport` |
| `export <SESSION_ID> [OUTPUT]` | `-c/--clipboard` | `Export` |
| `login` | `--oauth`, `--device-auth` | `Login` |
| `logout` | — | `Logout` |
| `models` | — | `Models` |
| `setup` | — | `Setup` |
| `trace <SESSION_ID>` | `--local`, `-o`, `--json` | `Trace` |
| `update` | `--check`, `--json`, `--force-reinstall`, `--version`, `--alpha`, `--stable` | `UpdateCheck`, `UpdateInstall` |
| `version` (alias `v`) | `--json` | — |
| `completions <SHELL>` | `bash`/`elvish`/`fish`/`powershell`/`zsh` | `Completions` |

### Interactive (not wrapped — out of scope for the headless SDK)
- `grok dashboard` — Agent Dashboard TUI.
- `grok wrap <CMD>...` — local PTY with OSC 52 clipboard forwarding (replaced `grok ssh`).

---

## 4. Headless output formats

`--output-format`: `plain` (raw text), `json` (single `GrokResult` object:
`text`, `stopReason`, `sessionId`, `requestId`, `thought`, …),
`streaming-json` (newline-delimited events → `StreamPrompt`).
`--json-schema` constrains `json` output to a caller-supplied JSON Schema.

---

## 5. Environment and files

| Name | Role |
| --- | --- |
| `GROK_SANDBOX` | Default sandbox profile. |
| `GROK_AGENT_SECRET` | Default `agent serve` auth secret. |
| `GROK_AGENT_DASHBOARD=0` | Disable the dashboard. |
| `GROK_HOME` | Root for `trace-exports/`. |
| `~/.grok/config.toml`, `./.grok/config.toml` | User/project config (MCP `--scope`). |
| `~/.grok/leader.sock` | Default leader socket (`--leader-socket`). |
| `MEMORY.md` | Cross-session memory store. |

---

## 6. Refresh procedure

1. Re-run `grok --help` and `grok help <subcommand>` for every command; diff here.
2. Re-capture fixtures under `test/testdata/` (json, streaming-json, models,
   mcp list, sessions list, leader info, version, inspect). See the
   `just cli-drift` recipe.
3. Bump the captured version header, then update `pkg/grok` option structs,
   wrappers, and tests.
