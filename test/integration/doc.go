//go:build integration

// Package integration contains end-to-end tests for the grok-go-sdk.
//
// Two lanes:
//   - Mock lane (default `go test -tags integration`): runs against the compiled
//     fake grok binary at test/mockserver/bin/grok-mock.
//   - Real lane (CI nightly): runs against the real grok CLI on PATH.
//
// Lane selection is controlled by env vars:
//   - GROK_MOCK_BIN: absolute path to the mock binary (forces mock lane)
//   - INTEGRATION_REAL=1: forces real lane (and the mock var is ignored)
package integration
