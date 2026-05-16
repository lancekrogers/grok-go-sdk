package grok

import (
	"context"
	"encoding/json"
)

type TraceOptions struct {
	LocalOnly bool
	Output    string
}

func (c *GrokClient) Trace(ctx context.Context, sessionID string, opts *TraceOptions) (string, bool, error) {
	args := []string{"trace"}
	if opts != nil {
		if opts.LocalOnly {
			args = append(args, "--local")
		}
		if opts.Output != "" {
			args = append(args, "-o", opts.Output)
		}
	}
	args = append(args, "--json", sessionID)
	out, err := c.runSubcommand(ctx, args)
	if err != nil {
		return "", false, err
	}
	var r struct {
		Path     string `json:"path"`
		Uploaded bool   `json:"uploaded"`
	}
	if err := json.Unmarshal(out, &r); err != nil {
		return "", false, &GrokError{Type: ErrorValidation, Message: "trace decode: " + err.Error()}
	}
	return r.Path, r.Uploaded, nil
}
