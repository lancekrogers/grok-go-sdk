# Grok Go SDK - Modular Justfile
# Run `just` to see all available commands

set dotenv-load := true

# Module imports
mod build '.justfiles/build.just'
mod demo '.justfiles/demos.just'
mod test '.justfiles/test.just'
mod coverage '.justfiles/coverage.just'
mod mock '.justfiles/mock.just'
mod util '.justfiles/util.just'
mod release '.justfiles/releases.just'

# Variables
PROJECT := "Grok Go SDK"
BIN_DIR := "./bin"
COVERAGE_DIR := "./coverage"

[private]
@default:
    @echo "Grok Go SDK"
    @echo "==========="
    @echo ""
    @just --list --unsorted

# Full pipeline: clean -> deps -> build -> test -> coverage
all: clean deps build-all test-all coverage-report

# Clean build artifacts
clean:
    just util clean

# Download dependencies
deps:
    just util deps

# Tidy all modules
tidy-all:
    just util tidy-all

# Run linters
lint:
    just util lint

# Check live grok CLI help against the committed snapshot
cli-drift:
    just util cli-drift

# Refresh the committed grok CLI help snapshot
cli-drift-update:
    just util cli-drift-update

# Generate coverage report
coverage-report:
    just coverage report

[private]
build-all:
  just build all

[private]
test-all:
  just test all
