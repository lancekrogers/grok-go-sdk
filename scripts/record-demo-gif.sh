#!/usr/bin/env bash
# Record demo GIFs using asciinema + agg.
#
# Most grok-go-sdk demos are one-shot Go binaries that run to completion,
# so they are spawned directly under asciinema. The agent_runtime demo is
# a stdin-reading REPL; it uses an expect script at scripts/demo-expect/
# to drive the interaction.
#
# Each recording runs in a fresh temp directory sandbox so the demo cannot
# touch the host repo. Demos that make real grok API calls cost real money.
#
# Usage:
#   scripts/record-demo-gif.sh <demo-name>
#   scripts/record-demo-gif.sh all       (records the cheap default set)
#   scripts/record-demo-gif.sh list      (prints available demos)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="${PROJECT_DIR}/docs/gif"
EXPECT_DIR="${SCRIPT_DIR}/demo-expect"

RED=$'\033[0;31m'
GREEN=$'\033[0;32m'
BLUE=$'\033[0;34m'
YELLOW=$'\033[1;33m'
NC=$'\033[0m'

ALL_DEMOS=(basic streaming sessions mcp agent-stdio worktree best-of-n check-loop subagents agent-runtime)
DEFAULT_DEMOS=(basic streaming sessions mcp agent-stdio subagents)
REPL_DEMOS=(agent-runtime)

usage() {
    cat <<EOF
Usage: $0 <demo|all|list>

Records demo GIFs to docs/gif/<demo>.gif.

Demos available:
$(for d in "${ALL_DEMOS[@]}"; do echo "  - $d"; done)

'all' records the default set (cheap demos only):
$(for d in "${DEFAULT_DEMOS[@]}"; do echo "  - $d"; done)

Recording demos NOT in the default set is supported one-at-a-time but
they are excluded from 'all' because they are slow or expensive:
  - worktree     (creates+removes a real worktree)
  - best-of-n    (N parallel grok calls; expensive)
  - check-loop   (multiple verification passes; slow)
  - agent-runtime (interactive REPL; needs expect script)
EOF
}

check_deps() {
    local missing=()
    for bin in asciinema agg just; do
        command -v "$bin" >/dev/null 2>&1 || missing+=("$bin")
    done
    if [ ${#missing[@]} -gt 0 ]; then
        echo "${RED}Missing dependencies: ${missing[*]}${NC}" >&2
        cat <<EOF >&2

Install with:
  brew install asciinema just
  cargo install --git https://github.com/asciinema/agg
EOF
        return 1
    fi
    if ! command -v grok >/dev/null 2>&1; then
        echo "${RED}grok CLI not found in PATH${NC}" >&2
        echo "Install from https://github.com/xai-org/grok and authenticate with 'grok login'." >&2
        return 1
    fi
}

needs_expect() {
    local d="$1"
    for r in "${REPL_DEMOS[@]}"; do
        [ "$d" = "$r" ] && return 0
    done
    return 1
}

valid_demo() {
    local d="$1"
    for a in "${ALL_DEMOS[@]}"; do
        [ "$d" = "$a" ] && return 0
    done
    return 1
}

record_oneshot() {
    local demo="$1"
    local cast_file="${OUTPUT_DIR}/${demo}.cast"
    local gif_file="${OUTPUT_DIR}/${demo}.gif"

    echo "${BLUE}Recording demo: ${demo}${NC}"

    local sandbox
    sandbox=$(mktemp -d)
    echo "${BLUE}Sandbox: ${sandbox}${NC}"
    trap 'rm -rf "$sandbox"' EXIT

    local force_term="TERM=xterm-256color FORCE_COLOR=1 CLICOLOR_FORCE=1"
    local typed_cmd="just demo ${demo}"
    local inner="cd $sandbox && echo '\$ $typed_cmd' && sleep 0.6 && cd '$PROJECT_DIR' && $typed_cmd"

    asciinema rec "$cast_file" \
        --overwrite \
        --cols=100 \
        --rows=25 \
        --command="env $force_term bash -lc \"$inner\""

    trap - EXIT
    rm -rf "$sandbox"

    if [ ! -f "$cast_file" ]; then
        echo "${RED}Recording failed: no cast file produced${NC}" >&2
        return 1
    fi

    convert_gif "$cast_file" "$gif_file"
}

record_with_expect() {
    local demo="$1"
    local expect_script="${EXPECT_DIR}/${demo}.exp"
    local cast_file="${OUTPUT_DIR}/${demo}.cast"
    local gif_file="${OUTPUT_DIR}/${demo}.gif"

    if [ ! -f "$expect_script" ]; then
        echo "${RED}Expect script not found: ${expect_script}${NC}" >&2
        return 1
    fi
    if ! command -v expect >/dev/null 2>&1; then
        echo "${RED}expect not installed. brew install expect${NC}" >&2
        return 1
    fi

    echo "${BLUE}Recording REPL demo: ${demo}${NC}"

    local sandbox
    sandbox=$(mktemp -d)
    trap 'rm -rf "$sandbox"' EXIT

    cd "$sandbox"
    export PROJECT_DIR
    export TERM="xterm-256color"
    export FORCE_COLOR="1"
    export CLICOLOR_FORCE="1"

    asciinema rec "$cast_file" \
        --overwrite \
        --cols=100 \
        --rows=25 \
        --command="PROJECT_DIR=$PROJECT_DIR TERM=$TERM FORCE_COLOR=$FORCE_COLOR expect $expect_script"

    trap - EXIT
    rm -rf "$sandbox"

    convert_gif "$cast_file" "$gif_file"
}

convert_gif() {
    local cast_file="$1"
    local gif_file="$2"
    echo "${BLUE}Converting to GIF...${NC}"
    agg \
        --theme=monokai \
        --font-size=14 \
        --cols=100 \
        --rows=25 \
        --speed=1.0 \
        "$cast_file" \
        "$gif_file"
    local cast_size
    cast_size=$(du -h "$cast_file" | cut -f1)
    local gif_size
    gif_size=$(du -h "$gif_file" | cut -f1)
    echo "${GREEN}Generated:${NC}"
    echo "  Cast: ${cast_file} (${cast_size})"
    echo "  GIF:  ${gif_file} (${gif_size})"
}

record_one() {
    local demo="$1"
    if ! valid_demo "$demo"; then
        echo "${RED}Unknown demo: ${demo}${NC}" >&2
        usage >&2
        exit 2
    fi
    mkdir -p "$OUTPUT_DIR"

    echo "${YELLOW}Note: real grok API credits will be consumed.${NC}"

    if needs_expect "$demo"; then
        record_with_expect "$demo"
    else
        record_oneshot "$demo"
    fi
}

record_all() {
    mkdir -p "$OUTPUT_DIR"
    echo "${BLUE}Recording default demo set (cheap demos only):${NC}"
    for d in "${DEFAULT_DEMOS[@]}"; do
        echo "  - $d"
    done
    echo "${YELLOW}Each recording makes real grok API calls.${NC}"
    for d in "${DEFAULT_DEMOS[@]}"; do
        record_one "$d" || echo "${YELLOW}Failed to record ${d}, continuing...${NC}"
    done
    echo "${GREEN}Default demo set recorded.${NC}"
}

main() {
    local arg="${1:-}"
    if [ -z "$arg" ] || [ "$arg" = "-h" ] || [ "$arg" = "--help" ]; then
        usage
        exit 0
    fi
    if [ "$arg" = "list" ]; then
        for d in "${ALL_DEMOS[@]}"; do echo "$d"; done
        exit 0
    fi
    check_deps
    if [ "$arg" = "all" ]; then
        record_all
    else
        record_one "$arg"
    fi
}

main "$@"
