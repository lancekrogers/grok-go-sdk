package grok

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestInspect_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ir, err := c.Inspect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ir.GrokVersion == "" {
		t.Fatalf("missing grokVersion: %#v", ir)
	}
}

func TestParseInspect_CamelCase(t *testing.T) {
	const sample = `{
	  "grokVersion":"0.2.77","channel":"stable",
	  "cwd":"/x","projectRoot":"/x/",
	  "projectTrusted":true,"bridgeTrusted":true,
	  "projectInstructions":[{"path":"/x/CLAUDE.md","scope":"project","fileType":"md","sizeBytes":12,"approxTokens":3,"vendor":"claude"}],
	  "permissions":{"loaded":[]},"loginPolicy":{"apiKeyAuthDisabled":false},
	  "hooks":[{"event":"PreToolUse","hookType":"command","target":"x","source":"user","matcher":"*"}],
	  "skills":[{"name":"deep-research","description":"d","source":"user","userInvocable":true}],
	  "agents":[{"name":"Explore","description":"d","source":"builtin"}],
	  "plugins":[{"name":"p","scope":"user","path":"/x","enabled":true,"provides":{}}],
	  "marketplaces":[],
	  "mcpServers":[{"name":"alpha","transport":"stdio","target":"/bin/x","source":"user","vendor":""}],
	  "lspServers":[],
	  "configSources":{"layers":[]},"externalCompat":{"cells":[]}
	}`
	ir, err := parseInspect([]byte(sample))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if ir.GrokVersion != "0.2.77" || ir.Channel != "stable" || ir.ProjectRoot != "/x/" {
		t.Fatalf("scalars mismatch: %#v", ir)
	}
	if !ir.ProjectTrusted || !ir.BridgeTrusted {
		t.Fatalf("trust flags mismatch: %#v", ir)
	}
	if len(ir.Skills) != 1 || ir.Skills[0].Name != "deep-research" || !ir.Skills[0].UserInvocable {
		t.Fatalf("skills mismatch: %#v", ir.Skills)
	}
	if len(ir.MCPServers) != 1 || ir.MCPServers[0].Transport != "stdio" {
		t.Fatalf("mcpServers mismatch: %#v", ir.MCPServers)
	}
	if len(ir.ProjectInstructions) != 1 || ir.ProjectInstructions[0].SizeBytes != 12 {
		t.Fatalf("instructions mismatch: %#v", ir.ProjectInstructions)
	}
	if len(ir.Permissions) == 0 || len(ir.ConfigSources) == 0 {
		t.Fatalf("raw objects should be populated: %#v", ir)
	}
	if len(ir.Extra) != 0 {
		t.Fatalf("Extra should be empty for a complete payload, got %#v", ir.Extra)
	}
}

func TestParseInspect_CapturedFixture(t *testing.T) {
	path := filepath.Join("..", "..", "test", "testdata", "inspect.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("fixture not yet captured: %v", err)
	}
	ir, err := parseInspect(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if ir.GrokVersion == "" {
		t.Fatalf("missing grokVersion: %#v", ir)
	}
	if ir.CWD == "" {
		t.Fatalf("missing cwd: %#v", ir)
	}
	if len(ir.Permissions) == 0 || len(ir.ConfigSources) == 0 {
		t.Fatalf("raw objects should be populated: %#v", ir)
	}
}

func TestParseInspect_UnknownKeyGoesToExtra(t *testing.T) {
	ir, err := parseInspect([]byte(`{"grokVersion":"x","futureField":42}`))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := ir.Extra["futureField"]; !ok {
		t.Fatalf("unknown key should land in Extra, got %#v", ir.Extra)
	}
}
