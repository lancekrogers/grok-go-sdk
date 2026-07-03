# Release v0.2.1

Recovery release for `v0.2.0` Go module checksum instability.

## Why This Exists

- `v0.2.0` was recreated after publication, which made Go module checksum resolution unstable across proxy edges.
- `v0.2.1` is the stable replacement tag for consumers.
- `go.mod` now retracts `v0.2.0` with guidance to use `v0.2.1`.

## Changes Since v0.2.0

- Added a `retract v0.2.0` directive.
- Kept the v0.2.0 CLI alignment code intact.
- Kept the release workflow fix that makes tag publishing idempotent.

## Consumer Action

Update dependencies to:

```bash
go get github.com/lancekrogers/grok-go-sdk@v0.2.1
```

Avoid `v0.2.0`.
