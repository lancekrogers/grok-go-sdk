package dangerous

import (
	"os"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

const enableEnv = "GROK_ENABLE_DANGEROUS"
const enableValue = "i-accept-all-risks"

var (
	ErrNotEnabled = &grok.GrokError{Type: grok.ErrorValidation, Message: "dangerous: GROK_ENABLE_DANGEROUS must be set to \"i-accept-all-risks\""}
	ErrProduction = &grok.GrokError{Type: grok.ErrorValidation, Message: "dangerous: refusing to run in production environment"}
)

type Client struct {
	inner *grok.GrokClient
}

func NewDangerousClient(binPath string) (*Client, error) {
	if os.Getenv(enableEnv) != enableValue {
		return nil, ErrNotEnabled
	}
	if isProduction() {
		return nil, ErrProduction
	}
	return &Client{inner: grok.NewClient(binPath)}, nil
}

func isProduction() bool {
	for _, k := range []string{"GO_ENV", "NODE_ENV"} {
		if os.Getenv(k) == "production" {
			return true
		}
	}
	return false
}

func (c *Client) BypassAllPermissions(prompt string, opts *grok.RunOptions) (*grok.GrokResult, error) {
	o := cloneOpts(opts)
	o.AlwaysApprove = true
	o.AllowDangerousMode = true
	return c.inner.RunPrompt(prompt, o)
}

func (c *Client) WithBypassPermissionsMode(prompt string, opts *grok.RunOptions) (*grok.GrokResult, error) {
	o := cloneOpts(opts)
	o.PermissionMode = grok.PermissionBypassPermissions
	o.AllowDangerousMode = true
	return c.inner.RunPrompt(prompt, o)
}

func (c *Client) DisableSandbox(prompt string, opts *grok.RunOptions) (*grok.GrokResult, error) {
	o := cloneOpts(opts)
	o.SandboxProfile = "off"
	o.AllowDangerousMode = true
	return c.inner.RunPrompt(prompt, o)
}

func cloneOpts(opts *grok.RunOptions) *grok.RunOptions {
	if opts == nil {
		return &grok.RunOptions{}
	}
	cp := *opts
	return &cp
}
