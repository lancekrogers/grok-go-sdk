package grok

const (
	ToolRead      = "Read"
	ToolGrep      = "Grep"
	ToolGlob      = "Glob"
	ToolBash      = "Bash"
	ToolWrite     = "Write"
	ToolEdit      = "Edit"
	ToolWebFetch  = "WebFetch"
	ToolWebSearch = "WebSearch"
)

func BuildToolGlob(name, pattern string) string {
	if pattern == "" {
		return name
	}
	return name + "(" + pattern + ")"
}
