package grok

type GrokResult struct {
	Text       string `json:"text"`
	StopReason string `json:"stopReason"`
	SessionID  string `json:"sessionId"`
	RequestID  string `json:"requestId"`
	Thought    string `json:"thought,omitempty"`

	CostUSD       float64 `json:"costUsd,omitempty"`
	DurationMS    int64   `json:"durationMs,omitempty"`
	DurationAPIMS int64   `json:"durationApiMs,omitempty"`
	NumTurns      int     `json:"numTurns,omitempty"`
	IsError       bool    `json:"isError,omitempty"`
	Subtype       string  `json:"subtype,omitempty"`
}
