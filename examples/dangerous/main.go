package main

import (
	"fmt"
	"log"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
	"github.com/lancekrogers/grok-go-sdk/pkg/grok/dangerous"
)

func main() {
	bin, err := grok.LocateBinary()
	if err != nil {
		log.Fatalf("locate grok: %v", err)
	}
	c, err := dangerous.NewDangerousClient(bin)
	if err != nil {
		log.Fatalf("enable dangerous client: %v", err)
	}
	res, err := c.BypassAllPermissions("Reply with exactly: dangerous mode opt-in confirmed.", &grok.RunOptions{
		Format: grok.JSONOutput,
	})
	if err != nil {
		log.Fatalf("run: %v", err)
	}
	fmt.Println("text:", res.Text)
}
