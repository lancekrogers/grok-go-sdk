package grok

import "context"

// supportedShells is the set of shells `grok completions` can generate scripts for.
var supportedShells = map[string]bool{
	"bash": true, "elvish": true, "fish": true, "powershell": true, "zsh": true,
}

// Completions returns a shell completion script for the given shell
// (bash, elvish, fish, powershell, zsh) as stdout bytes.
//
// Note: the interactive `grok wrap` (PTY/clipboard) and `grok dashboard`
// (TUI) commands are intentionally not wrapped by this headless SDK.
func (c *GrokClient) Completions(ctx context.Context, shell string) ([]byte, error) {
	if !supportedShells[shell] {
		return nil, &GrokError{Type: ErrorValidation, Message: "completions: unsupported shell " + shell}
	}
	return c.runSubcommand(ctx, []string{"completions", shell})
}
