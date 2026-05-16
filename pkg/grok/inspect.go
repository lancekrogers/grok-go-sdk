package grok

import (
	"context"
	"encoding/json"
)

type InspectReport struct {
	Version             string         `json:"version"`
	CWD                 string         `json:"cwd"`
	GitRoot             string         `json:"git_root"`
	ProjectTrusted      bool           `json:"project_trusted"`
	ProjectInstructions []string       `json:"project_instructions"`
	Permissions         map[string]any `json:"permissions"`
	Skills              []string       `json:"skills"`
	Extra               map[string]any `json:"-"`
}

func (c *GrokClient) Inspect(ctx context.Context) (*InspectReport, error) {
	out, err := c.runSubcommand(ctx, []string{"inspect", "--json"})
	if err != nil {
		return nil, err
	}
	var ir InspectReport
	if err := json.Unmarshal(out, &ir); err != nil {
		return nil, &GrokError{Type: ErrorValidation, Message: "inspect decode: " + err.Error()}
	}
	var raw map[string]any
	_ = json.Unmarshal(out, &raw)
	known := map[string]bool{
		"version": true, "cwd": true, "git_root": true, "project_trusted": true,
		"project_instructions": true, "permissions": true, "skills": true,
	}
	ir.Extra = map[string]any{}
	for k, v := range raw {
		if !known[k] {
			ir.Extra[k] = v
		}
	}
	return &ir, nil
}
