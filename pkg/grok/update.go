package grok

import (
	"context"
	"encoding/json"
)

type UpdateInfo struct {
	Current         string `json:"current"`
	Latest          string `json:"latest"`
	Channel         string `json:"channel"`
	UpdateAvailable bool   `json:"update_available"`
}

func (c *GrokClient) UpdateCheck(ctx context.Context) (*UpdateInfo, error) {
	out, err := c.runSubcommand(ctx, []string{"update", "--check", "--json"})
	if err != nil {
		return nil, err
	}
	var ui UpdateInfo
	if err := json.Unmarshal(out, &ui); err != nil {
		return nil, &GrokError{Type: ErrorValidation, Message: "update check decode: " + err.Error()}
	}
	return &ui, nil
}

func (c *GrokClient) UpdateInstall(ctx context.Context, version string, force, alpha, stable bool) error {
	args := []string{"update"}
	if force {
		args = append(args, "--force-reinstall")
	}
	if alpha {
		args = append(args, "--alpha")
	}
	if stable {
		args = append(args, "--stable")
	}
	if version != "" {
		args = append(args, "--version", version)
	}
	_, err := c.runSubcommand(ctx, args)
	return err
}
