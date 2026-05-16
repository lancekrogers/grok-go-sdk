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
		log.Fatalf("locate grok: %v", err)
	}
	ctx := context.Background()
	first, err := c.RunPromptCtx(ctx, "Remember the number 7. Just say OK.", &grok.RunOptions{
		Format: grok.JSONOutput,
	})
	if err != nil {
		log.Fatalf("first: %v", err)
	}
	fmt.Println("first session id:", first.SessionID)
	second, err := c.ResumeConversationCtx(ctx, "What number did I tell you?", first.SessionID)
	if err != nil {
		log.Fatalf("resume: %v", err)
	}
	fmt.Println("resumed answer:", second.Text)
}
