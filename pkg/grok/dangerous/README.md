# pkg/grok/dangerous

This subpackage exposes intentionally-risky operations on top of `pkg/grok`.

Operations here weaken safety controls (bypass permission prompts, disable
the sandbox, auto-approve all tool calls). They require explicit opt-in:

    export GROK_ENABLE_DANGEROUS="i-accept-all-risks"

`NewDangerousClient` returns an error if the env var is not set to the
sentinel value, or if `GO_ENV` or `NODE_ENV` is `production`.

Use cases:
- Trusted automation pipelines on isolated machines
- Local development with full grok permissions
- CI on disposable containers where speed matters more than sandbox

Do not import this package in code that runs against shared infrastructure,
in containers with mounted host paths, or in any context where untrusted
prompts can reach the grok process.
