package grok

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"sort"
	"strings"
)

type MCPTransport string

const (
	MCPTransportStdio MCPTransport = "stdio"
	MCPTransportHTTP  MCPTransport = "http"
	MCPTransportSSE   MCPTransport = "sse"
)

type MCPScope string

const (
	MCPScopeUser    MCPScope = "user"
	MCPScopeProject MCPScope = "project"
)

// MCPServerConfig describes an MCP server for `grok mcp add`.
//
// grok 0.2.77 uses a positional interface: a server name plus a single
// command-or-URL positional, with the server's own arguments after `--`.
type MCPServerConfig struct {
	Name         string
	CommandOrURL string            // command to launch (stdio) or URL to connect to (http/sse)
	Args         []string          // server command arguments, emitted after `--`
	Transport    MCPTransport      // -t/--transport (defaults to stdio when empty)
	Scope        MCPScope          // -s/--scope (defaults to user when empty)
	Env          map[string]string // -e KEY=value (repeatable)
	Headers      map[string]string // -H "Name: Value" (repeatable, remote transports)
}

func (c *GrokClient) MCPList(ctx context.Context) ([]MCPServerConfig, error) {
	out, err := c.runSubcommand(ctx, []string{"mcp", "list"})
	if err != nil {
		return nil, err
	}
	return parseMCPList(out), nil
}

func (c *GrokClient) MCPAdd(ctx context.Context, cfg MCPServerConfig) error {
	if cfg.Name == "" {
		return errors.New("mcp add: name required")
	}
	_, err := c.runSubcommand(ctx, mcpAddArgs(cfg))
	return err
}

// mcpAddArgs builds the argv for `grok mcp add` (0.2.77 positional interface):
//
//	mcp add [-t T] [-s S] [-e K=V...] [-H "N: V"...] <name> <commandOrURL> -- <args...>
func mcpAddArgs(cfg MCPServerConfig) []string {
	args := []string{"mcp", "add"}
	if cfg.Transport != "" {
		args = append(args, "-t", string(cfg.Transport))
	}
	if cfg.Scope != "" {
		args = append(args, "-s", string(cfg.Scope))
	}
	for _, k := range sortedKeys(cfg.Env) {
		args = append(args, "-e", k+"="+cfg.Env[k])
	}
	for _, k := range sortedKeys(cfg.Headers) {
		args = append(args, "-H", k+": "+cfg.Headers[k])
	}
	args = append(args, cfg.Name)
	if cfg.CommandOrURL != "" {
		args = append(args, cfg.CommandOrURL)
	}
	if len(cfg.Args) > 0 {
		args = append(args, "--")
		args = append(args, cfg.Args...)
	}
	return args
}

func (c *GrokClient) MCPRemove(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("mcp remove: name required")
	}
	_, err := c.runSubcommand(ctx, []string{"mcp", "remove", name})
	return err
}

func (c *GrokClient) MCPDoctor(ctx context.Context) (string, error) {
	out, err := c.runSubcommandTolerant(ctx, []string{"mcp", "doctor"})
	return string(out), err
}

// sortedKeys returns the map keys in deterministic order so emitted argv is stable.
func sortedKeys(m map[string]string) []string {
	if len(m) == 0 {
		return nil
	}
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// parseMCPList parses `grok mcp list` output. The 0.2.77 format is one line per
// server, "<name>: <command-or-url>"; the empty state is a single
// "No MCP servers configured" line.
func parseMCPList(b []byte) []MCPServerConfig {
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) == 0 || bytes.HasPrefix(trimmed, []byte("No MCP servers configured")) {
		return nil
	}
	var out []MCPServerConfig
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			continue
		}
		name := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		if name == "" {
			continue
		}
		out = append(out, MCPServerConfig{Name: name, CommandOrURL: val})
	}
	return out
}
