package grok

import "context"

type ImportOptions struct {
	List bool
	JSON bool
}

func (c *GrokClient) Import(ctx context.Context, targets []string, opts *ImportOptions) ([]byte, error) {
	args := []string{"import"}
	if opts != nil {
		if opts.List {
			args = append(args, "--list")
		}
		if opts.JSON {
			args = append(args, "--json")
		}
	}
	args = append(args, targets...)
	return c.runSubcommand(ctx, args)
}
