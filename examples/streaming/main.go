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
	events, errs := c.StreamPrompt(ctx, "Explain Go channels in three sentences", &grok.RunOptions{})
	for {
		select {
		case ev, ok := <-events:
			if !ok {
				fmt.Println()
				return
			}
			if ev.Type == grok.EventText {
				fmt.Print(ev.Content())
			}
		case err, ok := <-errs:
			if !ok {
				continue
			}
			if err != nil {
				log.Fatalf("stream error: %v", err)
			}
		}
	}
}
