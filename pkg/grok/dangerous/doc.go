// Package dangerous exposes explicit opt-in helpers for weakening grok safety
// controls, such as bypassing permission prompts or disabling the sandbox.
//
// Callers must set GROK_ENABLE_DANGEROUS=i-accept-all-risks before creating a
// client. The package refuses to run when GO_ENV or NODE_ENV is production.
package dangerous
