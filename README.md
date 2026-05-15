<h1 align="center">Grok Go SDK</h1>

A Go library for programmatically integrating the [grok](https://github.com/xai-org/grok)
CLI into Go applications. Wraps the headless and stdio-agent surfaces of the
`grok` binary so Go programs can drive Grok Build sessions, stream output,
manage sessions, attach MCP servers, and enforce permissions, budgets, and
retry policies.

> This SDK is in early development. Public API is unstable until `v1.0`.

## Status

- Repository: private during initial development.
- License: MIT.
- Reference design: see the
  [`workflow/design/go-grok-sdk/`](../../workflow/design/go-grok-sdk/)
  package in the obey-campaign planning repo for the full architecture,
  CLI reference, API surface, testing strategy, and execution plan.
- Captured `grok` CLI snapshot: `grok 0.1.210 (8b63e9068c)`.

## Prerequisites

- Go (latest stable; this module targets `go1.26.1`).
- The [grok CLI](https://github.com/xai-org/grok) installed and authenticated.
- [just](https://github.com/casey/just) for running project commands.

## Quick Start

```
just              # list top-level recipes
just build all    # build the library + examples
just test all     # run unit tests
just lint         # fmt + vet
```

Full recipe surface is documented dynamically by `just --list`.

## License

MIT. See [LICENSE](LICENSE).
