package llmGateway

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
)

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}

type ChatMessage struct {
	Role       string
	Content    string
	ToolCalls  []ToolCall
	Name       string
	ToolCallId string
}

type CompletionResult struct {
	ToolCalls []ToolCall
	Message   string
	Error     error
}

type LLMClient interface {
	Init(ctx context.Context) error
	StreamChatCompletion(ctx context.Context, messages []ChatMessage, tools []Tool, chanStream chan<- request.EventSSEData) ([]ToolCall, string)
}
