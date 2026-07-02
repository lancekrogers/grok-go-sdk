package grok

import (
	"context"
)

// This file wraps the `grok plugin` CLI subcommand tree (installable plugins and
// marketplace sources). It is unrelated to the in-process PluginManager / Plugin
// hook system in plugin.go — these are GrokClient methods that shell out to grok.

// PluginInstallOptions configures PluginInstall.
type PluginInstallOptions struct {
	Trust bool // --trust: trust immediately, skipping the confirmation prompt
}

// PluginUninstallOptions configures PluginUninstall.
type PluginUninstallOptions struct {
	Confirm  bool // --confirm: skip confirmation for multi-plugin repos
	KeepData bool // --keep-data: preserve the plugin's persistent data directory
}

// PluginTagOptions configures PluginTag.
type PluginTagOptions struct {
	Push   bool // --push: push the tag to the remote after creating it
	Force  bool // -f/--force: tag even if the tree is dirty or the tag exists
	DryRun bool // --dry-run: print what would be tagged without creating it
}

// PluginList lists installed plugins (grok plugin list).
func (c *GrokClient) PluginList(ctx context.Context) (string, error) {
	out, err := c.runSubcommand(ctx, []string{"plugin", "list"})
	return string(out), err
}

// PluginInstall installs a plugin from a git URL, GitHub shorthand, or local path.
func (c *GrokClient) PluginInstall(ctx context.Context, source string, opts PluginInstallOptions) error {
	if source == "" {
		return validationError("plugin install: source required")
	}
	_, err := c.runSubcommand(ctx, pluginInstallArgs(source, opts))
	return err
}

func pluginInstallArgs(source string, opts PluginInstallOptions) []string {
	args := []string{"plugin", "install", source}
	if opts.Trust {
		args = append(args, "--trust")
	}
	return args
}

// PluginUninstall uninstalls a plugin by name.
func (c *GrokClient) PluginUninstall(ctx context.Context, name string, opts PluginUninstallOptions) error {
	if name == "" {
		return validationError("plugin uninstall: name required")
	}
	_, err := c.runSubcommand(ctx, pluginUninstallArgs(name, opts))
	return err
}

func pluginUninstallArgs(name string, opts PluginUninstallOptions) []string {
	args := []string{"plugin", "uninstall", name}
	if opts.Confirm {
		args = append(args, "--confirm")
	}
	if opts.KeepData {
		args = append(args, "--keep-data")
	}
	return args
}

// PluginUpdate updates a named plugin, or all installed plugins when name is empty.
func (c *GrokClient) PluginUpdate(ctx context.Context, name string) error {
	args := []string{"plugin", "update"}
	if name != "" {
		args = append(args, name)
	}
	_, err := c.runSubcommand(ctx, args)
	return err
}

// PluginEnable enables a disabled plugin.
func (c *GrokClient) PluginEnable(ctx context.Context, name string) error {
	if name == "" {
		return validationError("plugin enable: name required")
	}
	_, err := c.runSubcommand(ctx, []string{"plugin", "enable", name})
	return err
}

// PluginDisable disables a plugin without uninstalling it.
func (c *GrokClient) PluginDisable(ctx context.Context, name string) error {
	if name == "" {
		return validationError("plugin disable: name required")
	}
	_, err := c.runSubcommand(ctx, []string{"plugin", "disable", name})
	return err
}

// PluginDetails shows a plugin's component inventory.
func (c *GrokClient) PluginDetails(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", validationError("plugin details: name required")
	}
	out, err := c.runSubcommand(ctx, []string{"plugin", "details", name})
	return string(out), err
}

// PluginValidate validates a plugin manifest at path (default "." when empty).
// It is tolerant of a non-zero exit so the caller still receives the report.
func (c *GrokClient) PluginValidate(ctx context.Context, path string) (string, error) {
	args := []string{"plugin", "validate"}
	if path != "" {
		args = append(args, path)
	}
	out, err := c.runSubcommandTolerant(ctx, args)
	return string(out), err
}

// PluginTag creates a release git tag from the plugin manifest version at path
// (default "." when empty).
func (c *GrokClient) PluginTag(ctx context.Context, path string, opts PluginTagOptions) (string, error) {
	out, err := c.runSubcommand(ctx, pluginTagArgs(path, opts))
	return string(out), err
}

func pluginTagArgs(path string, opts PluginTagOptions) []string {
	args := []string{"plugin", "tag"}
	if path != "" {
		args = append(args, path)
	}
	if opts.Push {
		args = append(args, "--push")
	}
	if opts.Force {
		args = append(args, "-f")
	}
	if opts.DryRun {
		args = append(args, "--dry-run")
	}
	return args
}

// MarketplaceList lists configured marketplace sources and their plugins.
func (c *GrokClient) MarketplaceList(ctx context.Context) (string, error) {
	out, err := c.runSubcommand(ctx, []string{"plugin", "marketplace", "list"})
	return string(out), err
}

// MarketplaceAdd adds a marketplace source (git URL or GitHub shorthand).
func (c *GrokClient) MarketplaceAdd(ctx context.Context, url string) error {
	if url == "" {
		return validationError("marketplace add: url required")
	}
	_, err := c.runSubcommand(ctx, []string{"plugin", "marketplace", "add", url})
	return err
}

// MarketplaceRemove removes a marketplace source (and uninstalls its plugins).
func (c *GrokClient) MarketplaceRemove(ctx context.Context, url string) error {
	if url == "" {
		return validationError("marketplace remove: url required")
	}
	_, err := c.runSubcommand(ctx, []string{"plugin", "marketplace", "remove", url})
	return err
}

// MarketplaceUpdate refreshes a marketplace source, or all when name is empty.
func (c *GrokClient) MarketplaceUpdate(ctx context.Context, name string) error {
	args := []string{"plugin", "marketplace", "update"}
	if name != "" {
		args = append(args, name)
	}
	_, err := c.runSubcommand(ctx, args)
	return err
}
