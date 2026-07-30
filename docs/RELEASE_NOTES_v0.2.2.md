# Release v0.2.2

Expose model catalog metadata that the Grok CLI already has, so integrators
(like obey) can meter and compact against real context windows.

## What's Changed

- `ModelInfo` gains `Name` and `ContextWindow` (tokens; `0` when unknown).
- `Models()` still uses `grok models` for the available ID list and default
  flag, then enriches each entry from `~/.grok/models_cache.json` (written by
  the CLI after fetching the models origin).
- Refresh CLI help snapshot for current Grok CLI (removed obsolete
  `leader profile` subcommands so `just cli-drift` stays green).

## Consumer Action

```bash
go get github.com/lancekrogers/grok-go-sdk@v0.2.2
```

When `ContextWindow > 0`, prefer it over static fallbacks. The CLI list remains
the source of which models are available; the cache only fills metadata.

**Full Changelog**: https://github.com/lancekrogers/grok-go-sdk/compare/v0.2.1...v0.2.2
