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
	sessions, dropped := parseSessionsReport(out)
	if dropped > 0 {
		return sessions, validationErrorf("sessions parse dropped %d unrecognized rows", dropped)
	}
	return sessions, nil
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
	sessions, dropped := parseSessionsReport(out)
	if dropped > 0 {
		return sessions, validationErrorf("sessions parse dropped %d unrecognized rows", dropped)
	}
	return sessions, nil
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
	sessions, _ := parseSessionsReport(data)
	return sessions
}

func parseSessionsReport(data []byte) ([]SessionSummary, int) {
	var out []SessionSummary
	dropped := 0
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if isSessionListMetadata(line) {
			continue
		}
		s := parseSessionLine(line)
		if s.ID != "" {
			out = append(out, s)
			continue
		}
		dropped++
	}
	return out, dropped
}

func isSessionListMetadata(line string) bool {
	return line == "" ||
		line == "(no label)" ||
		strings.HasPrefix(line, "#") ||
		strings.HasPrefix(line, "SESSION ID")
}

func parseSessionLine(line string) SessionSummary {
	if strings.Contains(line, "\t") {
		return parseTabSessionLine(line)
	}
	return parseTableSessionLine(line)
}

func parseTabSessionLine(line string) SessionSummary {
	parts := strings.Split(line, "\t")
	if len(parts) < 4 || strings.TrimSpace(parts[0]) == "" || !isDateLike(strings.TrimSpace(parts[1])) {
		return SessionSummary{}
	}
	s := SessionSummary{ID: strings.TrimSpace(parts[0])}
	if len(parts) > 1 {
		s.UpdatedAt = strings.TrimSpace(parts[1])
	}
	if len(parts) > 2 {
		s.CWD = strings.TrimSpace(parts[2])
	}
	if len(parts) > 3 {
		s.Summary = strings.TrimSpace(parts[3])
	}
	return s
}

func parseTableSessionLine(line string) SessionSummary {
	parts := strings.Fields(line)
	if len(parts) < 4 || parts[0] == "" || !isDateLike(parts[1]) || !isDateLike(parts[2]) {
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

func isDateLike(s string) bool {
	if len(s) < len("2006-01-02") {
		return false
	}
	for i, r := range s[:len("2006-01-02")] {
		switch i {
		case 4, 7:
			if r != '-' {
				return false
			}
		default:
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}
