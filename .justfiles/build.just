# Build Commands
# Run `just build` to see all build commands

PROJECT_DIR := justfile_directory()
BIN_DIR := PROJECT_DIR / "bin"

[private]
@default:
    @echo "Build Commands"
    @echo "=============="
    @just --list --justfile {{source_file()}}

# Build everything (library + dangerous + examples + mock)
all: lib dangerous examples mock

# Build core library
lib:
    @echo "Building core library..."
    cd {{PROJECT_DIR}} && go build ./pkg/grok
    @echo "Core library built"

# Build dangerous subpackage
dangerous:
    @echo "Building dangerous subpackage..."
    cd {{PROJECT_DIR}} && go build ./pkg/grok/dangerous
    @echo "Dangerous subpackage built"

# Build the mock grok binary used by unit tests
mock:
    @echo "Building mock grok binary..."
    @mkdir -p {{PROJECT_DIR}}/test/mockserver/bin
    cd {{PROJECT_DIR}} && go build -o test/mockserver/bin/grok-mock ./test/mockserver
    @echo "Built: test/mockserver/bin/grok-mock"

# Build all examples
examples: basic streaming sessions mcp agent-stdio worktree best-of-n check-loop subagents agent-runtime dangerous-example
    @echo "All examples built"

# Build basic example
basic:
    @echo "Building basic example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/basic ./examples/basic
    @echo "Built: {{BIN_DIR}}/basic"

# Build streaming example
streaming:
    @echo "Building streaming example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/streaming ./examples/streaming
    @echo "Built: {{BIN_DIR}}/streaming"

# Build sessions example
sessions:
    @echo "Building sessions example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/sessions ./examples/sessions
    @echo "Built: {{BIN_DIR}}/sessions"

# Build MCP example
mcp:
    @echo "Building MCP example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/mcp ./examples/mcp
    @echo "Built: {{BIN_DIR}}/mcp"

# Build stdio agent example
agent-stdio:
    @echo "Building agent_stdio example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/agent_stdio ./examples/agent_stdio
    @echo "Built: {{BIN_DIR}}/agent_stdio"

# Build worktree example
worktree:
    @echo "Building worktree example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/worktree ./examples/worktree
    @echo "Built: {{BIN_DIR}}/worktree"

# Build best-of-n example
best-of-n:
    @echo "Building best_of_n example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/best_of_n ./examples/best_of_n
    @echo "Built: {{BIN_DIR}}/best_of_n"

# Build check-loop example
check-loop:
    @echo "Building check_loop example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/check_loop ./examples/check_loop
    @echo "Built: {{BIN_DIR}}/check_loop"

# Build subagents example
subagents:
    @echo "Building subagents example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/subagents ./examples/subagents
    @echo "Built: {{BIN_DIR}}/subagents"

# Build agent_runtime example (embedded long-running agent host:
# budget tracker, audit plugin, tool filter, stdio agent session)
agent-runtime:
    @echo "Building agent_runtime example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/agent_runtime ./examples/agent_runtime
    @echo "Built: {{BIN_DIR}}/agent_runtime"

# Build dangerous example
dangerous-example:
    @echo "Building dangerous example..."
    @mkdir -p {{BIN_DIR}}
    cd {{PROJECT_DIR}} && go build -o {{BIN_DIR}}/dangerous_example ./examples/dangerous
    @echo "Built: {{BIN_DIR}}/dangerous_example"
