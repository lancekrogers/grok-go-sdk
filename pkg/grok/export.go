package grok

import "context"

// ExportOptions configures Export (grok export <SESSION_ID> [OUTPUT] [-c]).
type ExportOptions struct {
	Output    string // optional output file path; empty => return stdout bytes
	Clipboard bool   // -c/--clipboard: copy to the clipboard instead of stdout
}

// Export writes a session transcript as Markdown. With no Output path the
// transcript is returned as stdout bytes; this is the surviving "shareable
// transcript" capability that replaced the removed `grok share`.
func (c *GrokClient) Export(ctx context.Context, sessionID string, opts ExportOptions) ([]byte, error) {
	if sessionID == "" {
		return nil, &GrokError{Type: ErrorValidation, Message: "export: sessionID required"}
	}
	return c.runSubcommand(ctx, exportArgs(sessionID, opts))
}

func exportArgs(sessionID string, opts ExportOptions) []string {
	args := []string{"export", sessionID}
	if opts.Output != "" {
		args = append(args, opts.Output)
	}
	if opts.Clipboard {
		args = append(args, "-c")
	}
	return args
}
