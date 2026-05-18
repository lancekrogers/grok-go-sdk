# Demo Commands
# Run `just demo` to see all demo commands
#
# Each demo builds the corresponding example under examples/<name>/ and
# runs the resulting binary. Demos require a working grok install and a
# valid grok login.

set dotenv-load := true

PROJECT_DIR := justfile_directory()
BIN_DIR := PROJECT_DIR / "bin"

[private]
@default:
    @echo "Demo Commands"
    @echo "============="
    @just --list --justfile {{source_file()}}

# Build and run basic demo
basic: build-basic
    @echo "Starting basic demo..."
    {{BIN_DIR}}/basic

# Build and run streaming demo
streaming: build-streaming
    @echo "Starting streaming demo..."
    {{BIN_DIR}}/streaming

# Build and run sessions demo
sessions: build-sessions
    @echo "Starting sessions demo..."
    {{BIN_DIR}}/sessions

# Build and run MCP demo
mcp: build-mcp
    @echo "Starting MCP demo..."
    {{BIN_DIR}}/mcp

# Build and run stdio agent demo
agent-stdio: build-agent-stdio
    @echo "Starting agent_stdio demo..."
    {{BIN_DIR}}/agent_stdio

# Build and run worktree demo
worktree: build-worktree
    @echo "Starting worktree demo..."
    {{BIN_DIR}}/worktree

# Build and run best-of-n demo
best-of-n: build-best-of-n
    @echo "Starting best_of_n demo..."
    {{BIN_DIR}}/best_of_n

# Build and run check-loop demo
check-loop: build-check-loop
    @echo "Starting check_loop demo..."
    {{BIN_DIR}}/check_loop

# Build and run subagents demo
subagents: build-subagents
    @echo "Starting subagents demo..."
    {{BIN_DIR}}/subagents

# Build and run agent_runtime demo
agent-runtime: build-agent-runtime
    @echo "Starting agent_runtime demo..."
    {{BIN_DIR}}/agent_runtime

# ===== Build recipes =====

build-basic:
    @just build basic

build-streaming:
    @just build streaming

build-sessions:
    @just build sessions

build-mcp:
    @just build mcp

build-agent-stdio:
    @just build agent-stdio

build-worktree:
    @just build worktree

build-best-of-n:
    @just build best-of-n

build-check-loop:
    @just build check-loop

build-subagents:
    @just build subagents

build-agent-runtime:
    @just build agent-runtime

# Build all demos
build-all: build-basic build-streaming build-sessions build-mcp build-agent-stdio build-worktree build-best-of-n build-check-loop build-subagents build-agent-runtime
    @echo "All demos built"

# ===== GIF generation (placeholders; script lands in scripts/) =====

gif demo-name:
    @{{PROJECT_DIR}}/scripts/record-demo-gif.sh {{demo-name}}

gif-all:
    @{{PROJECT_DIR}}/scripts/record-demo-gif.sh all

gif-list:
    @echo "Available demos for GIF generation:"
    @echo "  - basic"
    @echo "  - streaming"
    @echo "  - sessions"
    @echo "  - mcp"
    @echo "  - agent-stdio"
    @echo "  - worktree"
    @echo "  - best-of-n"
    @echo "  - check-loop"
    @echo "  - subagents"
    @echo "  - agent-runtime"
    @echo ""
    @echo "Usage: just demo gif <demo-name>"
    @echo "       just demo gif-all"
