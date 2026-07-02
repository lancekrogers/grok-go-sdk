package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lancekrogers/grok-go-sdk/pkg/grok"
)

// Demonstrates --json-schema structured output: the model is constrained to
// emit JSON matching the schema, which we decode into a typed Go struct.
func main() {
	c, err := grok.NewClientFromPath()
	if err != nil {
		log.Fatalf("locate grok: %v", err)
	}

	const schema = `{
	  "type": "object",
	  "properties": {
	    "language": {"type": "string"},
	    "year":     {"type": "integer"}
	  },
	  "required": ["language", "year"]
	}`

	// JSONSchema implies --output-format json, so Format need not be set.
	res, err := c.RunPrompt("In what year was the Go programming language first released, and what is its name?", &grok.RunOptions{
		JSONSchema: schema,
	})
	if err != nil {
		log.Fatalf("run: %v", err)
	}

	var out struct {
		Language string `json:"language"`
		Year     int    `json:"year"`
	}
	if err := json.Unmarshal([]byte(res.Text), &out); err != nil {
		log.Fatalf("decode structured output %q: %v", res.Text, err)
	}
	fmt.Printf("language=%s year=%d\n", out.Language, out.Year)
}
