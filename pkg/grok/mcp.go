package grok

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"strings"
)

type MCPTransport string

const (
	MCPTransportStdio MCPTransport = "stdio"
	MCPTransportHTTP  MCPTransport = "http"
	MCPTransportSSE   MCPTransport = "sse"
)

type MCPServerConfig struct {
	Name      string
	Command   string
	Args      []string
	Env       map[string]string
	URL       string
	Transport MCPTransport
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
	args := []string{"mcp", "add", cfg.Name}
	if cfg.Command != "" {
		args = append(args, "--command", cfg.Command)
	}
	if len(cfg.Args) > 0 {
		args = append(args, "--args")
		args = append(args, cfg.Args...)
	}
	for k, v := range cfg.Env {
		args = append(args, "--env", k+"="+v)
	}
	if cfg.URL != "" {
		args = append(args, "--url", cfg.URL)
	}
	if cfg.Transport != "" {
		args = append(args, "--type", string(cfg.Transport))
	}
	_, err := c.runSubcommand(ctx, args)
	return err
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

func parseMCPList(b []byte) []MCPServerConfig {
	if bytes.HasPrefix(bytes.TrimSpace(b), []byte("No MCP servers configured")) {
		return nil
	}
	var out []MCPServerConfig
	sc := bufio.NewScanner(bytes.NewReader(b))
	var cur *MCPServerConfig
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)
		if trim == "" {
			if cur != nil {
				out = append(out, *cur)
				cur = nil
			}
			continue
		}
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			if cur != nil {
				out = append(out, *cur)
			}
			cur = &MCPServerConfig{Name: trim}
			continue
		}
		if cur == nil {
			continue
		}
		switch {
		case strings.HasPrefix(trim, "command:"):
			cur.Command = strings.TrimSpace(strings.TrimPrefix(trim, "command:"))
		case strings.HasPrefix(trim, "url:"):
			cur.URL = strings.TrimSpace(strings.TrimPrefix(trim, "url:"))
		case strings.HasPrefix(trim, "transport:"):
			cur.Transport = MCPTransport(strings.TrimSpace(strings.TrimPrefix(trim, "transport:")))
		}
	}
	if cur != nil {
		out = append(out, *cur)
	}
	return out
}
