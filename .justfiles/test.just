# Test Commands
# Run `just test` to see all test commands

set dotenv-load := true

PROJECT_DIR := justfile_directory()

[private]
@default:
    @echo "Test Commands"
    @echo "============="
    @just --list --justfile {{source_file()}}

# Run all unit tests (library + dangerous subpackage)
all: lib dangerous

# Run core library tests
lib:
    @echo "Testing core library..."
    go test -v {{PROJECT_DIR}}/pkg/grok
    @echo "Core library tests complete"

# Run dangerous package tests
dangerous:
    @echo "Testing dangerous package..."
    go test -v {{PROJECT_DIR}}/pkg/grok/dangerous
    @echo "Dangerous package tests complete"

# Run integration tests against the mock grok binary
integration:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Running integration tests against mock binary..."
    just mock build
    GROK_MOCK_BIN="{{PROJECT_DIR}}/test/mockserver/bin/grok-mock" \
    INTEGRATION_TESTS=1 USE_MOCK_BINARY=1 \
        go test -tags=integration {{PROJECT_DIR}}/test/integration -v
    echo "Mock-binary integration tests complete"

# Run integration tests against the real grok CLI
integration-real:
    @echo "Running real integration tests..."
    @echo "Requires grok CLI installed and authenticated"
    INTEGRATION_TESTS=1 go test -tags=integration {{PROJECT_DIR}}/test/integration -v
    @echo "Real integration tests complete"

# Run tests with race detection
race:
    @echo "Running tests with race detection..."
    go test -race {{PROJECT_DIR}}/pkg/grok/...
    @echo "Race detection tests complete"

# Run benchmark tests
bench:
    @echo "Running benchmarks..."
    go test -bench=. {{PROJECT_DIR}}/pkg/grok/...
    @echo "Benchmarks complete"

# Run tests with memory profiling
memory:
    @echo "Running memory profiling..."
    go test -memprofile=mem.out {{PROJECT_DIR}}/pkg/grok/...
    @echo "Memory profile saved to mem.out"

# Run quick tests (no verbose)
quick:
    @echo "Running quick tests..."
    go test {{PROJECT_DIR}}/pkg/grok/...
    @echo "Quick tests complete"

# Run tests for a specific package
[no-exit-message]
pkg package:
    @echo "Testing package: {{package}}"
    go test -v {{PROJECT_DIR}}/{{package}}

# Verify examples compile
examples:
    @echo "Compiling examples (build-only)..."
    cd {{PROJECT_DIR}} && go build ./examples/...
    @echo "All examples compile"
