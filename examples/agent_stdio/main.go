package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func main() {
	c, err := grok.NewClientFromPath()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	s, err := c.StartStdioAgent(ctx, nil)
	if err != nil {
		log.Fatalf("start: %v", err)
	}
	defer s.Close()

	if _, err := s.Initialize(ctx, "grok-go-sdk-example", "0.1.0"); err != nil {
		log.Fatalf("initialize: %v", err)
	}
	if _, err := s.Authenticate(ctx, "cached_token"); err != nil {
		log.Fatalf("authenticate: %v", err)
	}
	cwd, _ := os.Getwd()
	sessionID, err := s.NewSession(ctx, cwd, nil)
	if err != nil {
		log.Fatalf("new session: %v", err)
	}
	fmt.Printf("session: %s\n\n", sessionID)

	go func() {
		for u := range s.Updates() {
			if u.Update.SessionUpdate == grok.UpdateAgentMessageChunk {
				fmt.Print(u.Update.ContentText())
			}
		}
	}()

	for i, prompt := range []string{
		"What is 2+2? Reply with just the number.",
		"And 3+3? Just the number.",
		"What was my first question?",
	} {
		fmt.Printf("\n--- turn %d: %s ---\n", i+1, prompt)
		res, err := s.PromptText(ctx, sessionID, prompt)
		if err != nil {
			log.Fatalf("turn %d: %v", i+1, err)
		}
		fmt.Printf("\n[stop=%s tokens=%d]\n", res.StopReason, res.Meta.TotalTokens)
	}
}
