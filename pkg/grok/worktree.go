package grok

import (
	"context"
	"strings"
	"time"
)

type Worktree struct {
	ID       string
	Path     string
	Branch   string
	LastUsed string
}

func (c *GrokClient) WorktreeList(ctx context.Context) ([]Worktree, error) {
	out, err := c.runSubcommand(ctx, []string{"worktree", "list"})
	if err != nil {
		return nil, err
	}
	return parseWorktreeList(out), nil
}

func (c *GrokClient) WorktreeShow(ctx context.Context, id string) (*Worktree, error) {
	out, err := c.runSubcommand(ctx, []string{"worktree", "show", id})
	if err != nil {
		return nil, err
	}
	list := parseWorktreeList(out)
	if len(list) == 0 {
		return nil, &GrokError{Type: ErrorValidation, Message: "no worktree " + id}
	}
	return &list[0], nil
}

func (c *GrokClient) WorktreeRemove(ctx context.Context, ids []string, force, dryRun bool) error {
	if len(ids) == 0 {
		return &GrokError{Type: ErrorValidation, Message: "worktree remove: no ids"}
	}
	args := []string{"worktree", "rm"}
	if force {
		args = append(args, "-f")
	}
	if dryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, ids...)
	_, err := c.runSubcommand(ctx, args)
	return err
}

func (c *GrokClient) WorktreeGC(ctx context.Context, maxAge time.Duration, force, dryRun bool) error {
	args := []string{"worktree", "gc"}
	if maxAge > 0 {
		args = append(args, "--max-age", maxAge.String())
	}
	if force {
		args = append(args, "-f")
	}
	if dryRun {
		args = append(args, "--dry-run")
	}
	_, err := c.runSubcommand(ctx, args)
	return err
}

func (c *GrokClient) WorktreeDBRebuild(ctx context.Context) error {
	_, err := c.runSubcommand(ctx, []string{"worktree", "db", "rebuild"})
	return err
}

func (c *GrokClient) WorktreeDBStats(ctx context.Context) (string, error) {
	out, err := c.runSubcommand(ctx, []string{"worktree", "db", "stats"})
	return strings.TrimSpace(string(out)), err
}

func (c *GrokClient) WorktreeDBPath(ctx context.Context) (string, error) {
	out, err := c.runSubcommand(ctx, []string{"worktree", "db", "path"})
	return strings.TrimSpace(string(out)), err
}

func parseWorktreeList(b []byte) []Worktree {
	var out []Worktree
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 1 {
			w := Worktree{ID: parts[0]}
			if len(parts) > 1 {
				w.Path = parts[1]
			}
			if len(parts) > 2 {
				w.Branch = parts[2]
			}
			if len(parts) > 3 {
				w.LastUsed = parts[3]
			}
			out = append(out, w)
		}
	}
	return out
}
