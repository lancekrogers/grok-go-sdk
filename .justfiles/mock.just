# Mock Binary Commands
# Run `just mock` to see all mock binary commands
#
# The grok-go-sdk mock is a compiled binary (not a long-running HTTP server).
# It impersonates the `grok` CLI for unit tests: argv inspection, stdout
# fixtures, configurable exit codes. Tests point the SDK at this binary
# instead of the real `grok`.

set dotenv-load := true

PROJECT_DIR := justfile_directory()
MOCK_DIR := PROJECT_DIR / "test/mockserver"
MOCK_BIN := MOCK_DIR / "bin/grok-mock"

[private]
@default:
    @echo "Mock Binary Commands"
    @echo "===================="
    @just --list --justfile {{source_file()}}

# Build the mock grok binary
build:
    @echo "Building mock grok binary..."
    @mkdir -p {{MOCK_DIR}}/bin
    cd {{PROJECT_DIR}} && go build -o {{MOCK_BIN}} ./test/mockserver
    @echo "Built: {{MOCK_BIN}}"

# Rebuild (clean + build)
rebuild: clean build

# Print the path to the mock binary
path:
    @echo "{{MOCK_BIN}}"

# Smoke test the mock binary against the captured say-hello fixture
smoke: build
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -x "{{MOCK_BIN}}" ]; then
        echo "Mock binary missing at {{MOCK_BIN}}"
        exit 1
    fi
    echo "Invoking mock with -p 'say hello' --output-format json..."
    GROK_MOCK_SCENARIO="say-hello" "{{MOCK_BIN}}" -p "say hello" --output-format json

# Clean the mock binary
clean:
    @echo "Cleaning mock binary..."
    rm -rf {{MOCK_DIR}}/bin
    @echo "Mock binary cleaned"

# Show captured testdata fixtures
fixtures:
    @echo "Captured fixtures under test/testdata/:"
    @find {{PROJECT_DIR}}/test/testdata -type f 2>/dev/null || echo "  (none yet)"
