package grok

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

type LeaderInfo struct {
	PID     int    `json:"pid"`
	Started string `json:"started"`
	Socket  string `json:"socket,omitempty"`
	Version string `json:"version,omitempty"`
}

func (c *GrokClient) LeaderList(ctx context.Context) ([]LeaderInfo, error) {
	out, err := c.runSubcommand(ctx, []string{"leader", "list", "--json"})
	if err != nil {
		return nil, err
	}
	var list []LeaderInfo
	if err := json.Unmarshal(out, &list); err != nil {
		return parseLeaderListText(out), nil
	}
	return list, nil
}

func (c *GrokClient) LeaderInfoCmd(ctx context.Context, pid int) (*LeaderInfo, error) {
	out, err := c.runSubcommand(ctx, []string{"leader", "info", strconv.Itoa(pid), "--json"})
	if err != nil {
		return nil, err
	}
	var info LeaderInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *GrokClient) LeaderKill(ctx context.Context) error {
	_, err := c.runSubcommand(ctx, []string{"leader", "kill"})
	return err
}

func (c *GrokClient) LeaderProfileStart(ctx context.Context, pid int) error {
	_, err := c.runSubcommand(ctx, []string{"leader", "profile", "start", strconv.Itoa(pid)})
	return err
}

func (c *GrokClient) LeaderProfileStop(ctx context.Context, pid int) error {
	_, err := c.runSubcommand(ctx, []string{"leader", "profile", "stop", strconv.Itoa(pid)})
	return err
}

func (c *GrokClient) LeaderProfileStatus(ctx context.Context, pid int) (string, error) {
	out, err := c.runSubcommand(ctx, []string{"leader", "profile", "status", strconv.Itoa(pid)})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func parseLeaderListText(b []byte) []LeaderInfo {
	var out []LeaderInfo
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "PID") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			pid, _ := strconv.Atoi(parts[0])
			out = append(out, LeaderInfo{PID: pid, Started: parts[1]})
		}
	}
	return out
}
