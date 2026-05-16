# Contributing

Thanks for your interest in `grok-go-sdk`. This document covers how to build
the SDK, run the tests, capture fixtures against a real `grok` CLI, and open
a pull request.

## Project Layout

- `pkg/grok/`: the public library
- `pkg/grok/dangerous/`: opt-in escape-hatch subpackage
- `examples/`: self-contained demo programs
- `test/mockserver/`: compiled fake grok binary for unit tests
- `test/testdata/`: captured fixtures from real `grok` runs
- `test/integration/`: integration tests behind the `integration` build tag

## Building

```
just build all
```

## Testing

```
just test all                            # unit tests
just test integration                    # integration tests against mock binary
INTEGRATION_REAL=1 just test integration # against real grok
just coverage report                     # HTML coverage
```

## Fixture Capture

Captured stdout from real grok sessions lives under `test/testdata/`. To
refresh fixtures after a grok release:

```
grok --version  # bump CLI_REFERENCE.md if changed
grok -p "say hello" --output-format json > test/testdata/json/say-hello.json
grok -p "say hello" --output-format streaming-json > test/testdata/streaming-json/basic.jsonl
```

See `test/testdata/<scenario>.about` files for capture metadata.

## PR Conventions

- One logical change per PR.
- Commits follow `<type>: <subject>` (chore, feat, fix, test, docs, refactor).
- Tests must pass on the mock lane; the real lane runs in CI.
- No emdashes in markdown.
- No comments in Go code unless the WHY is non-obvious.

## License

MIT. By contributing, you agree your contributions are licensed under MIT.
