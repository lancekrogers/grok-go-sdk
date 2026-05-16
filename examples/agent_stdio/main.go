package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func main() {
	c, err := grok.NewClientFromPath()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	s, err := c.StartStdioAgent(ctx, &grok.StdioConfig{Model: "grok-build"})
	if err != nil {
		log.Fatalf("start: %v", err)
	}
	defer s.Close()
	for i, prompt := range []string{
		"What is 2+2?",
		"And 3+3?",
		"What was my first question?",
	} {
		resp, err := s.RequestResponse(ctx, grok.AgentMessage{Type: "user", Text: prompt})
		if err != nil {
			log.Fatalf("turn %d: %v", i+1, err)
		}
		fmt.Printf("turn %d: %s\n", i+1, resp.Text)
	}
}
