package main

import (
	"fmt"
	"log"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

func main() {
	c, err := grok.NewClientFromPath()
	if err != nil {
		log.Fatalf("locate grok: %v", err)
	}
	res, err := c.RunPrompt("Write a one-line Go function that returns 42", &grok.RunOptions{
		Format: grok.JSONOutput,
	})
	if err != nil {
		log.Fatalf("run: %v", err)
	}
	fmt.Println("session:", res.SessionID)
	fmt.Println("text:", res.Text)
}
