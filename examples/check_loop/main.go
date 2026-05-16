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
	res, err := c.RunPromptCtx(ctx, "Find the largest prime below 100.", &grok.RunOptions{
		Format:       grok.JSONOutput,
		Check:        true,
		MaxBudgetUSD: 0.10,
	})
	if err != nil {
		log.Fatalf("run: %v", err)
	}
	fmt.Println(res.Text)
}
