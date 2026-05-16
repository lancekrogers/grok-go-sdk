package grok

import "encoding/json"

func MarshalAgents(agents map[string]*SubagentConfig) (string, error) {
	if len(agents) == 0 {
		return "", nil
	}
	b, err := json.Marshal(agents)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
