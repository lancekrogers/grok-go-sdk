package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func main() {
	c, err := grok.NewClientFromPath()
	if err != nil {
		log.Fatal(err)
	}
	agents := map[string]*grok.SubagentConfig{
		"security": {
			Description: "Security analysis and vulnerability detection",
			Prompt:      "You are a security expert. Analyze code for vulnerabilities.",
			Tools:       []string{grok.ToolRead, grok.ToolGrep},
			Model:       "grok-build",
		},
		"testing": {
			Description: "Test generation and coverage analysis",
			Prompt:      "You are a testing expert. Generate comprehensive tests.",
			Tools:       []string{grok.ToolRead, grok.ToolWrite},
		},
	}
	res, err := c.RunPromptCtx(context.Background(),
		"Review this snippet for SQL injection: db.Exec(\"SELECT * FROM users WHERE id=\" + req.ID)",
		&grok.RunOptions{
			Format:       grok.JSONOutput,
			Agent:        "security",
			Agents:       agents,
			MaxBudgetUSD: 0.05,
		})
	if err != nil {
		log.Fatalf("run: %v", err)
	}
	fmt.Println(res.Text)
}
