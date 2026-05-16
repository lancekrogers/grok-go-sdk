package grok

import "context"

type MemoryScope int

const (
	MemoryWorkspace MemoryScope = 1 << iota
	MemoryGlobal
)
const MemoryAll = MemoryWorkspace | MemoryGlobal

func (c *GrokClient) MemoryClear(ctx context.Context, scope MemoryScope) error {
	args := []string{"memory", "clear", "-y"}
	switch {
	case scope == MemoryAll:
		args = append(args, "--all")
	case scope&MemoryWorkspace != 0:
		args = append(args, "--workspace")
	case scope&MemoryGlobal != 0:
		args = append(args, "--global")
	default:
		return &GrokError{Type: ErrorValidation, Message: "MemoryClear: scope required"}
	}
	_, err := c.runSubcommand(ctx, args)
	return err
}
