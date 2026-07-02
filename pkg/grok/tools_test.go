package grok

import "testing"

func TestBuiltinToolConstants_VerifiedGrokBuildNames(t *testing.T) {
	want := map[string]string{
		"read":            "GrokBuild:read_file",
		"grep":            "GrokBuild:grep",
		"list_dir":        "GrokBuild:list_dir",
		"bash":            "GrokBuild:run_terminal_cmd",
		"edit":            "GrokBuild:search_replace",
		"web_fetch":       "GrokBuild:web_fetch",
		"web_search":      "GrokBuild:web_search",
		"task":            "GrokBuild:task",
		"get_task_output": "GrokBuild:get_task_output",
		"kill_task":       "GrokBuild:kill_task",
		"enter_plan":      "GrokBuild:enter_plan_mode",
		"exit_plan":       "GrokBuild:exit_plan_mode",
		"ask_user":        "GrokBuild:ask_user_question",
	}
	got := map[string]string{
		"read":            ToolRead,
		"grep":            ToolGrep,
		"list_dir":        ToolListDir,
		"bash":            ToolBash,
		"edit":            ToolEdit,
		"web_fetch":       ToolWebFetch,
		"web_search":      ToolWebSearch,
		"task":            ToolTask,
		"get_task_output": ToolGetTaskOutput,
		"kill_task":       ToolKillTask,
		"enter_plan":      ToolPlanEnter,
		"exit_plan":       ToolPlanExit,
		"ask_user":        ToolAskUser,
	}
	for name, wantValue := range want {
		if got[name] != wantValue {
			t.Fatalf("%s got %q want %q", name, got[name], wantValue)
		}
	}
	if ToolGlob != ToolListDir {
		t.Fatalf("ToolGlob got %q want ToolListDir %q", ToolGlob, ToolListDir)
	}
	if ToolWrite != ToolEdit {
		t.Fatalf("ToolWrite got %q want ToolEdit %q", ToolWrite, ToolEdit)
	}
}
