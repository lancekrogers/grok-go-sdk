package grok

import (
	"encoding/json"
	"fmt"
)

const acpProtocolVersion = "2024-11-05"

type RPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`

	Raw json.RawMessage `json:"-"`
}

type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

func (m *RPCMessage) IsRequest() bool      { return m.Method != "" && len(m.ID) > 0 }
func (m *RPCMessage) IsNotification() bool { return m.Method != "" && len(m.ID) == 0 }
func (m *RPCMessage) IsResponse() bool     { return m.Method == "" && len(m.ID) > 0 }

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeParams struct {
	ProtocolVersion string     `json:"protocolVersion"`
	ClientInfo      ClientInfo `json:"clientInfo"`
}

type InitializeResult struct {
	ProtocolVersion   int                  `json:"protocolVersion"`
	AgentCapabilities AgentCapabilities    `json:"agentCapabilities"`
	AuthMethods       []AuthMethod         `json:"authMethods"`
	Meta              *InitializeMeta      `json:"_meta,omitempty"`
}

type AgentCapabilities struct {
	LoadSession        bool                  `json:"loadSession"`
	PromptCapabilities PromptCapabilities    `json:"promptCapabilities"`
	MCPCapabilities    MCPCapabilities       `json:"mcpCapabilities"`
	Meta               json.RawMessage       `json:"_meta,omitempty"`
}

type PromptCapabilities struct {
	Image           bool `json:"image"`
	Audio           bool `json:"audio"`
	EmbeddedContext bool `json:"embeddedContext"`
}

type MCPCapabilities struct {
	HTTP bool `json:"http"`
	SSE  bool `json:"sse"`
}

type AuthMethod struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type InitializeMeta struct {
	GrokShell               bool             `json:"grokShell"`
	CurrentWorkingDirectory string           `json:"currentWorkingDirectory"`
	AgentVersion            string           `json:"agentVersion"`
	AgentID                 string           `json:"agentId"`
	AgentInstanceID         string           `json:"agentInstanceId"`
	Hostname                string           `json:"hostname"`
	ModelState              *ModelState      `json:"modelState,omitempty"`
	MCPServers              []MCPServerInfo  `json:"mcpServers,omitempty"`
	AvailableCommands       []AgentCommand   `json:"availableCommands,omitempty"`
	CancelRewind            bool             `json:"cancelRewind"`
}

type ModelState struct {
	CurrentModelID  string       `json:"currentModelId"`
	AvailableModels []ModelInfo2 `json:"availableModels"`
}

type ModelInfo2 struct {
	ModelID     string          `json:"modelId"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Meta        json.RawMessage `json:"_meta,omitempty"`
}

type MCPServerInfo struct {
	Name    string   `json:"name"`
	Source  string   `json:"source"`
	Type    string   `json:"type"`
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

type AgentCommand struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Input       *AgentCommandInput `json:"input"`
}

type AgentCommandInput struct {
	Hint string `json:"hint,omitempty"`
}

type AuthenticateParams struct {
	MethodID string `json:"methodId"`
}

type AuthenticateResult struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

type MCPServerSpec struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
}

type NewSessionParams struct {
	CWD        string          `json:"cwd"`
	MCPServers []MCPServerSpec `json:"mcpServers"`
}

type NewSessionResult struct {
	SessionID string `json:"sessionId"`
}

type PromptBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

func TextBlock(text string) PromptBlock {
	return PromptBlock{Type: "text", Text: text}
}

type PromptParams struct {
	SessionID string        `json:"sessionId"`
	Prompt    []PromptBlock `json:"prompt"`
}

type PromptResult struct {
	StopReason string       `json:"stopReason"`
	Meta       *PromptMeta  `json:"_meta,omitempty"`
}

type PromptMeta struct {
	SessionID         string `json:"sessionId"`
	RequestID         string `json:"requestId"`
	PromptID          string `json:"promptId"`
	ModelID           string `json:"modelId"`
	TotalTokens       int    `json:"totalTokens"`
	InputTokens       int    `json:"inputTokens"`
	OutputTokens      int    `json:"outputTokens"`
	CachedReadTokens  int    `json:"cachedReadTokens"`
	ReasoningTokens   int    `json:"reasoningTokens"`
}

type SessionUpdateKind string

const (
	UpdateAgentMessageChunk      SessionUpdateKind = "agent_message_chunk"
	UpdateAgentThoughtChunk      SessionUpdateKind = "agent_thought_chunk"
	UpdateAvailableCommands      SessionUpdateKind = "available_commands_update"
	UpdateToolCall               SessionUpdateKind = "tool_call"
	UpdateToolCallUpdate         SessionUpdateKind = "tool_call_update"
	UpdatePlan                   SessionUpdateKind = "plan"
)

type SessionUpdate struct {
	SessionID string             `json:"sessionId"`
	Update    SessionUpdateInner `json:"update"`
	Raw       json.RawMessage    `json:"-"`
}

type SessionUpdateInner struct {
	SessionUpdate     SessionUpdateKind `json:"sessionUpdate"`
	Content           *ContentBlock     `json:"content,omitempty"`
	AvailableCommands []AgentCommand    `json:"availableCommands,omitempty"`
	ToolCallID        string            `json:"toolCallId,omitempty"`
	Status            string            `json:"status,omitempty"`
	Title             string            `json:"title,omitempty"`
	Kind              string            `json:"kind,omitempty"`
	Raw               json.RawMessage   `json:"-"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

func (u SessionUpdateInner) ContentText() string {
	if u.Content == nil {
		return ""
	}
	return u.Content.Text
}
