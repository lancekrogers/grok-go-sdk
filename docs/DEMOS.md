# Demos

Each demo is a runnable example under `examples/<name>/main.go`. All demos
require a working `grok` install and an authenticated grok login.

| Name | Demonstrates | Recipe |
|------|--------------|--------|
| basic | Single JSON-output prompt | `just demo basic` |
| streaming | Channel-driven streaming | `just demo streaming` |
| sessions | Continue then resume | `just demo sessions` |
| mcp | Attach MCP, run doctor | `just demo mcp` |
| agent_stdio | Long-running stdio session | `just demo agent-stdio` |
| worktree | -w worktree with cleanup | `just demo worktree` |
| best_of_n | 4-way headless race | `just demo best-of-n` |
| check_loop | Self-verifying response | `just demo check-loop` |
| subagents | Inline --agents JSON | `just demo subagents` |
| agent_runtime | Embedded long-running agent host | `just demo agent-runtime` |

The `pkg/grok/dangerous` subpackage is intentionally excluded from the
default demo set so users do not accidentally invoke bypass shortcuts.
