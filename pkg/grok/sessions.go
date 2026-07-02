package grok

import (
	"bufio"
	"bytes"
	"context"
	"strconv"
	"strings"
)

type SessionSummary struct {
	ID          string
	Summary     string
	UpdatedAt   string
	CWD         string
	FirstPrompt string
}

func (c *GrokClient) SessionsList(ctx context.Context, limit int) ([]SessionSummary, error) {
	args := []string{"sessions", "list"}
	if limit > 0 {
		args = append(args, "-n", strconv.Itoa(limit))
	}
	out, err := c.runSubcommand(ctx, args)
	if err != nil {
		return nil, err
	}
	return parseSessions(out), nil
}

func (c *GrokClient) SessionsSearch(ctx context.Context, query string, limit int) ([]SessionSummary, error) {
	args := []string{"sessions", "search", query}
	if limit > 0 {
		args = append(args, "-n", strconv.Itoa(limit))
	}
	out, err := c.runSubcommand(ctx, args)
	if err != nil {
		return nil, err
	}
	return parseSessions(out), nil
}

// SessionsDelete permanently deletes a session from history (sessions delete <id>).
func (c *GrokClient) SessionsDelete(ctx context.Context, id string) error {
	if id == "" {
		return &GrokError{Type: ErrorValidation, Message: "sessions delete: id required"}
	}
	_, err := c.runSubcommand(ctx, []string{"sessions", "delete", id})
	return err
}

func parseSessions(data []byte) []SessionSummary {
	var out []SessionSummary
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || line == "(no label)" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "SESSION ID") {
			continue
		}
		s := parseSessionLine(line)
		if s.ID != "" {
			out = append(out, s)
		}
	}
	return out
}

func parseSessionLine(line string) SessionSummary {
	if strings.Contains(line, "\t") {
		return parseTabSessionLine(line)
	}
	return parseTableSessionLine(line)
}

func parseTabSessionLine(line string) SessionSummary {
	parts := strings.Split(line, "\t")
	s := SessionSummary{}
	if len(parts) > 0 {
		s.ID = parts[0]
	}
	if len(parts) > 1 {
		s.UpdatedAt = parts[1]
	}
	if len(parts) > 2 {
		s.CWD = parts[2]
	}
	if len(parts) > 3 {
		s.Summary = parts[3]
	}
	return s
}

func parseTableSessionLine(line string) SessionSummary {
	parts := strings.Fields(line)
	if len(parts) < 4 || !isUUIDLike(parts[0]) {
		return SessionSummary{}
	}
	s := SessionSummary{
		ID:        parts[0],
		UpdatedAt: parts[2],
	}
	if len(parts) > 4 {
		s.Summary = strings.Join(parts[4:], " ")
	}
	return s
}

func isUUIDLike(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, r := range s {
		switch i {
		case 8, 13, 18, 23:
			if r != '-' {
				return false
			}
		default:
			if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
				return false
			}
		}
	}
	return true
}
