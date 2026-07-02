package grok

import (
	"context"
	"reflect"
	"testing"
)

func TestPluginArgBuilders(t *testing.T) {
	if got := pluginInstallArgs("user/repo@v1", PluginInstallOptions{Trust: true}); !reflect.DeepEqual(got, []string{"plugin", "install", "user/repo@v1", "--trust"}) {
		t.Errorf("install: %v", got)
	}
	if got := pluginUninstallArgs("p", PluginUninstallOptions{Confirm: true, KeepData: true}); !reflect.DeepEqual(got, []string{"plugin", "uninstall", "p", "--confirm", "--keep-data"}) {
		t.Errorf("uninstall: %v", got)
	}
	if got := pluginTagArgs("", PluginTagOptions{Push: true, Force: true, DryRun: true}); !reflect.DeepEqual(got, []string{"plugin", "tag", "--push", "-f", "--dry-run"}) {
		t.Errorf("tag flags: %v", got)
	}
	if got := pluginTagArgs("./plug", PluginTagOptions{}); !reflect.DeepEqual(got, []string{"plugin", "tag", "./plug"}) {
		t.Errorf("tag path: %v", got)
	}
}

func TestPluginEmptyGuards(t *testing.T) {
	c := NewClient("/nonexistent")
	ctx := context.Background()
	if err := c.PluginInstall(ctx, "", PluginInstallOptions{}); err == nil {
		t.Error("install empty source should error")
	} else {
		requireValidationError(t, err)
	}
	if err := c.PluginUninstall(ctx, "", PluginUninstallOptions{}); err == nil {
		t.Error("uninstall empty name should error")
	}
	if err := c.PluginEnable(ctx, ""); err == nil {
		t.Error("enable empty name should error")
	}
	if err := c.PluginDisable(ctx, ""); err == nil {
		t.Error("disable empty name should error")
	}
	if _, err := c.PluginDetails(ctx, ""); err == nil {
		t.Error("details empty name should error")
	}
	if err := c.MarketplaceAdd(ctx, ""); err == nil {
		t.Error("marketplace add empty url should error")
	}
	if err := c.MarketplaceRemove(ctx, ""); err == nil {
		t.Error("marketplace remove empty url should error")
	}
}

func TestPluginCLI_AgainstMock(t *testing.T) {
	mock := buildOrLocateMock(t)
	c := NewClient(mock)
	ctx := context.Background()
	if out, err := c.PluginList(ctx); err != nil || out == "" {
		t.Fatalf("plugin list: out=%q err=%v", out, err)
	}
	if out, err := c.MarketplaceList(ctx); err != nil || out == "" {
		t.Fatalf("marketplace list: out=%q err=%v", out, err)
	}
	if out, err := c.PluginValidate(ctx, "."); err != nil || out == "" {
		t.Fatalf("plugin validate: out=%q err=%v", out, err)
	}
	if err := c.PluginUpdate(ctx, ""); err != nil {
		t.Fatalf("plugin update all: %v", err)
	}
	if err := c.MarketplaceUpdate(ctx, ""); err != nil {
		t.Fatalf("marketplace update all: %v", err)
	}
}
