package grok

import (
	"context"
	"encoding/json"
)

// InspectInstruction is one entry in InspectReport.ProjectInstructions.
type InspectInstruction struct {
	Path         string `json:"path"`
	Scope        string `json:"scope"`
	FileType     string `json:"fileType"`
	SizeBytes    int64  `json:"sizeBytes"`
	ApproxTokens int64  `json:"approxTokens"`
	Vendor       string `json:"vendor"`
}

// InspectHook is one configured hook in InspectReport.Hooks.
type InspectHook struct {
	Event    string          `json:"event"`
	HookType string          `json:"hookType"`
	Target   string          `json:"target"`
	Source   json.RawMessage `json:"source"`
	Matcher  string          `json:"matcher"`
}

// InspectSkill is one skill in InspectReport.Skills.
type InspectSkill struct {
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Source        json.RawMessage `json:"source"`
	UserInvocable bool            `json:"userInvocable"`
}

// InspectAgent is one agent in InspectReport.Agents.
type InspectAgent struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Source      json.RawMessage `json:"source"`
}

// InspectPlugin is one plugin in InspectReport.Plugins.
type InspectPlugin struct {
	Name     string          `json:"name"`
	Scope    string          `json:"scope"`
	Path     string          `json:"path"`
	Enabled  bool            `json:"enabled"`
	Provides json.RawMessage `json:"provides"`
}

// InspectMCPServer is one MCP server in InspectReport.MCPServers.
type InspectMCPServer struct {
	Name      string          `json:"name"`
	Transport string          `json:"transport"`
	Target    string          `json:"target"`
	Source    json.RawMessage `json:"source"`
	Vendor    string          `json:"vendor"`
}

// InspectReport is the decoded `grok inspect --json` camelCase payload.
// Free-form objects are kept as json.RawMessage; any unknown top-level key lands
// in Extra for forward compatibility.
type InspectReport struct {
	GrokVersion         string               `json:"grokVersion"`
	Channel             string               `json:"channel"`
	CWD                 string               `json:"cwd"`
	ProjectRoot         string               `json:"projectRoot"`
	ProjectTrusted      bool                 `json:"projectTrusted"`
	BridgeTrusted       bool                 `json:"bridgeTrusted"`
	ProjectInstructions []InspectInstruction `json:"projectInstructions"`
	Permissions         json.RawMessage      `json:"permissions"`
	LoginPolicy         json.RawMessage      `json:"loginPolicy"`
	Hooks               []InspectHook        `json:"hooks"`
	Skills              []InspectSkill       `json:"skills"`
	Agents              []InspectAgent       `json:"agents"`
	Plugins             []InspectPlugin      `json:"plugins"`
	Marketplaces        []json.RawMessage    `json:"marketplaces"`
	MCPServers          []InspectMCPServer   `json:"mcpServers"`
	LSPServers          []json.RawMessage    `json:"lspServers"`
	ConfigSources       json.RawMessage      `json:"configSources"`
	ExternalCompat      json.RawMessage      `json:"externalCompat"`
	Extra               map[string]any       `json:"-"`
}

func (c *GrokClient) Inspect(ctx context.Context) (*InspectReport, error) {
	out, err := c.runSubcommand(ctx, []string{"inspect", "--json"})
	if err != nil {
		return nil, err
	}
	return parseInspect(out)
}

// inspectKnownKeys is the set of top-level keys InspectReport decodes directly;
// anything else is preserved in InspectReport.Extra.
var inspectKnownKeys = map[string]bool{
	"grokVersion": true, "channel": true, "cwd": true, "projectRoot": true,
	"projectTrusted": true, "bridgeTrusted": true, "projectInstructions": true,
	"permissions": true, "loginPolicy": true, "hooks": true, "skills": true,
	"agents": true, "plugins": true, "marketplaces": true, "mcpServers": true,
	"lspServers": true, "configSources": true, "externalCompat": true,
}

func parseInspect(out []byte) (*InspectReport, error) {
	var ir InspectReport
	if err := json.Unmarshal(out, &ir); err != nil {
		return nil, &GrokError{Type: ErrorValidation, Message: "inspect decode: " + err.Error()}
	}
	var raw map[string]any
	_ = json.Unmarshal(out, &raw)
	ir.Extra = map[string]any{}
	for k, v := range raw {
		if !inspectKnownKeys[k] {
			ir.Extra[k] = v
		}
	}
	return &ir, nil
}
