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
	ctx := context.Background()
	name := "example-" + grok.GenerateSessionID()
	defer func() {
		if err := c.WorktreeRemove(ctx, []string{name}, true, false); err != nil {
			log.Printf("cleanup: %v", err)
		}
	}()
	res, err := c.RunPromptCtx(ctx, "List the files in the current directory", &grok.RunOptions{
		Format:       grok.JSONOutput,
		Worktree:     grok.WorktreeOption{Enabled: true, Name: name},
		MaxBudgetUSD: 0.05,
	})
	if err != nil {
		log.Fatalf("run: %v", err)
	}
	fmt.Println(res.Text)
}
