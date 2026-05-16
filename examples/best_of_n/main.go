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
	res, err := c.RunPromptCtx(ctx, "Write a haiku about goroutines.", &grok.RunOptions{
		Format:       grok.JSONOutput,
		BestOfN:      4,
		MaxBudgetUSD: 0.20,
	})
	if err != nil {
		log.Fatalf("run: %v", err)
	}
	fmt.Println(res.Text)
}
