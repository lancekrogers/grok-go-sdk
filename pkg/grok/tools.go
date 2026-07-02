package grok

const (
	ToolRead          = "GrokBuild:read_file"
	ToolGrep          = "GrokBuild:grep"
	ToolListDir       = "GrokBuild:list_dir"
	ToolGlob          = ToolListDir
	ToolBash          = "GrokBuild:run_terminal_cmd"
	ToolEdit          = "GrokBuild:search_replace"
	ToolWrite         = ToolEdit
	ToolWebFetch      = "GrokBuild:web_fetch"
	ToolWebSearch     = "GrokBuild:web_search"
	ToolTask          = "GrokBuild:task"
	ToolGetTaskOutput = "GrokBuild:get_task_output"
	ToolKillTask      = "GrokBuild:kill_task"
	ToolPlanEnter     = "GrokBuild:enter_plan_mode"
	ToolPlanExit      = "GrokBuild:exit_plan_mode"
	ToolAskUser       = "GrokBuild:ask_user_question"
)

func BuildToolGlob(name, pattern string) string {
	if pattern == "" {
		return name
	}
	return name + "(" + pattern + ")"
}
