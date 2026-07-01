// Mock grok CLI impersonator used by unit tests.
//
// Phase 001/02: skeleton with synthetic responses.
// Phase 003/01: enhanced with fixture routing from test/testdata/.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]

	if os.Getenv("GROK_MOCK_ARGV_DUMP") == "1" {
		for _, a := range args {
			fmt.Fprintln(os.Stderr, a)
		}
	}
	if s := os.Getenv("GROK_MOCK_STDERR"); s != "" {
		fmt.Fprint(os.Stderr, s)
	}

	code := route(args)

	if v := os.Getenv("GROK_MOCK_EXIT_CODE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			os.Exit(n)
		}
	}
	os.Exit(code)
}

func route(args []string) int {
	if len(args) == 0 {
		return 0
	}
	switch args[0] {
	case "sessions":
		return doSessions(args[1:])
	case "share":
		fmt.Println("https://example.invalid/share/mock-skeleton")
		return 0
	case "trace":
		fmt.Println(`{"path":"/tmp/mock-skeleton-trace.tar.gz","uploaded":false}`)
		return 0
	case "import":
		return 0
	case "export":
		fmt.Println("# Mock transcript\n\nhello")
		return 0
	case "inspect":
		return doInspect(args[1:])
	case "models":
		fmt.Println("Default model: grok-build")
		fmt.Println("* grok-build (default)")
		return 0
	case "mcp":
		return doMCP(args[1:])
	case "memory":
		return 0
	case "plugin":
		return doPlugin(args[1:])
	case "worktree":
		return doWorktree(args[1:])
	case "update":
		return doUpdate(args[1:])
	case "setup":
		return 0
	case "login":
		return 0
	case "logout":
		return 0
	case "agent":
		return doAgent(args[1:])
	case "leader":
		return doLeader(args[1:])
	}
	return doPrompt(args)
}

func doPrompt(args []string) int {
	format := flagValue(args, "--output-format")
	scenario := os.Getenv("GROK_MOCK_SCENARIO")
	if scenario == "" {
		scenario = "say-hello"
	}
	switch format {
	case "json":
		return emitFixture("json/" + scenario + ".json")
	case "streaming-json":
		return emitFixture("streaming-json/" + scenario + ".jsonl")
	case "plain", "":
		return emitPlain(scenario)
	default:
		fmt.Fprintf(os.Stderr, "mock: unknown output-format %q\n", format)
		return 2
	}
}

func emitFixture(rel string) int {
	base := os.Getenv("GROK_MOCK_TESTDATA")
	if base == "" {
		base = locateTestdata()
	}
	if base == "" {
		fmt.Fprintln(os.Stderr, "mock: cannot locate test/testdata directory")
		return 99
	}
	path := filepath.Join(base, rel)
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mock: fixture missing: %s: %v\n", path, err)
		return 99
	}
	os.Stdout.Write(b)
	return 0
}

func emitPlain(scenario string) int {
	switch scenario {
	case "say-hello":
		fmt.Println("Hello! How can I help you today?")
	default:
		fmt.Println("mock plain response: " + scenario)
	}
	return 0
}

func locateTestdata() string {
	if exe, err := os.Executable(); err == nil {
		for d := filepath.Dir(exe); ; {
			candidate := filepath.Join(d, "test", "testdata")
			if st, err := os.Stat(candidate); err == nil && st.IsDir() {
				return candidate
			}
			parent := filepath.Dir(d)
			if parent == d {
				break
			}
			d = parent
		}
	}
	cwd, _ := os.Getwd()
	for d := cwd; ; {
		candidate := filepath.Join(d, "test", "testdata")
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			return candidate
		}
		parent := filepath.Dir(d)
		if parent == d {
			return ""
		}
		d = parent
	}
}

func doPlugin(args []string) int {
	if len(args) == 0 {
		return 0
	}
	switch args[0] {
	case "list":
		fmt.Println("mock-plugin v1.0 (enabled)")
	case "details":
		fmt.Println("Components: 1 command, 0 agents")
	case "validate":
		fmt.Println("Manifest OK")
	case "marketplace":
		if len(args) > 1 && args[1] == "list" {
			fmt.Println("mock-marketplace: https://example.invalid/mp")
		}
	}
	return 0
}

func doSessions(args []string) int {
	if len(args) == 0 {
		return 0
	}
	switch args[0] {
	case "list", "search":
		fmt.Println("mock-session-01HMOCK0000000000000000000\t2026-05-16T02:00:00Z\t/tmp\tmock summary")
	case "delete":
		// success, no output
	}
	return 0
}

func doInspect(args []string) int {
	if flagSet(args, "--json") {
		fmt.Println(`{"grokVersion":"mock","channel":"stable","cwd":"/tmp","projectRoot":"/tmp/","projectTrusted":true,"bridgeTrusted":true,"projectInstructions":[],"permissions":{},"loginPolicy":{},"hooks":[],"skills":[],"agents":[],"plugins":[],"marketplaces":[],"mcpServers":[],"lspServers":[],"configSources":{},"externalCompat":{}}`)
		return 0
	}
	fmt.Println("Mock inspect output (use --json for machine-readable)")
	return 0
}

func doMCP(args []string) int {
	if len(args) == 0 {
		return 0
	}
	switch args[0] {
	case "list":
		fmt.Println("# mock mcp servers (empty)")
	case "add", "remove", "doctor":
	}
	return 0
}

func doWorktree(args []string) int {
	if len(args) == 0 {
		return 0
	}
	switch args[0] {
	case "list", "show":
		fmt.Println("mock-wt-01\t/tmp/mock-wt-01\tmock-branch\t2026-05-16")
	case "db":
		if len(args) > 1 && args[1] == "path" {
			fmt.Println("/tmp/mock-worktree.db")
		}
	}
	return 0
}

func doUpdate(args []string) int {
	if flagSet(args, "--check") && flagSet(args, "--json") {
		fmt.Println(`{"current":"mock","latest":"mock","channel":"stable","update_available":false}`)
	}
	return 0
}

func doAgent(args []string) int {
	if len(args) == 0 {
		return 0
	}
	switch args[0] {
	case "stdio":
		return doAgentStdio()
	case "serve":
		return doAgentServe(args[1:])
	case "headless", "leader":
		return 0
	}
	return 0
}

func doAgentStdio() int {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	var sessionCounter int
	for sc.Scan() {
		line := sc.Bytes()
		var m map[string]any
		if err := json.Unmarshal(line, &m); err != nil {
			continue
		}
		method, _ := m["method"].(string)
		id, hasID := m["id"]
		if method == "" || !hasID {
			continue
		}
		switch method {
		case "initialize":
			writeJSON(map[string]any{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]any{
					"protocolVersion": 1,
					"agentCapabilities": map[string]any{
						"loadSession":        true,
						"promptCapabilities": map[string]bool{"image": false, "audio": false, "embeddedContext": true},
						"mcpCapabilities":    map[string]bool{"http": true, "sse": true},
					},
					"authMethods": []map[string]any{
						{"id": "cached_token", "name": "cached_token", "description": "Cached token"},
					},
				},
			})
		case "authenticate":
			writeJSON(map[string]any{
				"jsonrpc": "2.0",
				"id":      id,
				"result":  map[string]any{"_meta": map[string]any{"email": "mock@example.invalid", "auth_mode": "mock"}},
			})
		case "session/new":
			sessionCounter++
			sessID := fmt.Sprintf("mock-session-%08d", sessionCounter)
			writeJSON(map[string]any{
				"jsonrpc": "2.0",
				"id":      id,
				"result":  map[string]any{"sessionId": sessID},
			})
		case "session/prompt":
			params, _ := m["params"].(map[string]any)
			sessID, _ := params["sessionId"].(string)
			promptBlocks, _ := params["prompt"].([]any)
			var promptText string
			for _, b := range promptBlocks {
				if bb, ok := b.(map[string]any); ok {
					if t, _ := bb["text"].(string); t != "" {
						promptText += t
					}
				}
			}
			emitStreamingChunks(sessID, "Mock response to: "+promptText)
			writeJSON(map[string]any{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]any{
					"stopReason": "end_turn",
					"_meta": map[string]any{
						"sessionId":    sessID,
						"requestId":    fmt.Sprintf("mock-req-%08d", sessionCounter),
						"promptId":     fmt.Sprintf("mock-prompt-%08d", sessionCounter),
						"modelId":      "grok-mock",
						"totalTokens":  100,
						"inputTokens":  60,
						"outputTokens": 40,
					},
				},
			})
		default:
			writeJSON(map[string]any{
				"jsonrpc": "2.0",
				"id":      id,
				"error":   map[string]any{"code": -32601, "message": "Method not found"},
			})
		}
	}
	return 0
}

func emitStreamingChunks(sessID, text string) {
	for _, word := range strings.Fields(text) {
		writeJSON(map[string]any{
			"jsonrpc": "2.0",
			"method":  "session/update",
			"params": map[string]any{
				"sessionId": sessID,
				"update": map[string]any{
					"sessionUpdate": "agent_message_chunk",
					"content":       map[string]string{"type": "text", "text": word + " "},
				},
			},
		})
	}
}

func writeJSON(v any) {
	b, _ := json.Marshal(v)
	os.Stdout.Write(append(b, '\n'))
}

func doAgentServe(args []string) int {
	bind := flagValue(args, "--bind")
	if bind == "" {
		bind = "127.0.0.1:0"
	}
	fmt.Fprintf(os.Stderr, "listening on %s\n", bind)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	return 0
}

func doLeader(args []string) int {
	if len(args) > 0 && args[0] == "list" {
		if flagSet(args, "--json") {
			fmt.Println(`[]`)
		}
	}
	return 0
}

func flagValue(args []string, name string) string {
	for i, a := range args {
		if a == name && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(a, name+"=") {
			return strings.TrimPrefix(a, name+"=")
		}
	}
	return ""
}

func flagSet(args []string, name string) bool {
	for _, a := range args {
		if a == name || strings.HasPrefix(a, name+"=") {
			return true
		}
	}
	return false
}
