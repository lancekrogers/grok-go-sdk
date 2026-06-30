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
	servers, err := c.MCPList(ctx)
	if err != nil {
		log.Fatalf("mcp list: %v", err)
	}
	fmt.Printf("%d MCP server(s) configured\n", len(servers))
	for _, s := range servers {
		fmt.Printf("  - %s: %s\n", s.Name, s.CommandOrURL)
	}
	report, err := c.MCPDoctor(ctx)
	fmt.Println("\n--- mcp doctor ---")
	if report != "" {
		fmt.Println(report)
	}
	if err != nil {
		fmt.Printf("(doctor reported: %v)\n", err)
	}
}
