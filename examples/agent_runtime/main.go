package main

import (
	"bufio"
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

	pm := grok.NewPluginManager()
	pm.Register(&grok.LoggingPlugin{SanitizeSecrets: true, TruncateLength: 200}, nil)
	pm.Register(&grok.MetricsPlugin{}, nil)
	pm.Register(grok.NewAuditPlugin(2048), nil)
	pm.Register(grok.NewToolFilterPlugin(map[string]string{
		"Bash(rm -rf*)":    "blanket rm -rf prohibited",
		"Bash(curl http*)": "outbound curl prohibited in this runtime",
	}), nil)

	budget := grok.NewBudgetTracker(&grok.BudgetConfig{
		MaxBudgetUSD:     5.00,
		WarningThreshold: 0.8,
		OnBudgetWarning:  func(c, m float64) { log.Printf("budget at %.0f%%", (c/m)*100) },
		OnBudgetExceeded: func(c, m float64) { log.Fatalf("budget exceeded: $%.2f > $%.2f", c, m) },
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = pm.Initialize(ctx)
	defer pm.Shutdown(ctx)

	s, err := c.StartStdioAgent(ctx, &grok.StdioConfig{Model: "grok-build"})
	if err != nil {
		log.Fatalf("stdio: %v", err)
	}
	defer s.Close()

	fmt.Println("agent_runtime: type a prompt and press Enter; Ctrl+C to exit")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		prompt := sc.Text()
		if prompt == "" {
			continue
		}
		rctx, rcancel := context.WithTimeout(ctx, 90*time.Second)
		resp, err := s.RequestResponse(rctx, grok.AgentMessage{Type: "user", Text: prompt})
		rcancel()
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		budget.Add(0.001)
		fmt.Println("agent:", resp.Text)
		fmt.Printf("[spent so far: $%.4f]\n", budget.TotalSpent())
	}
}
