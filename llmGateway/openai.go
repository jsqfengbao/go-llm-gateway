package llmGateway

import (
	"context"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
)

// OpenAIClient 是 OpenAI 客户端的结构体
type OpenAIClient struct {
	APIKey string
	Model  string
	// 这里可以添加 OpenAI 客户端的其他配置
}

func (o *OpenAIClient) Init(ctx context.Context) error {
	// 初始化 OpenAI 客户端
	// 这里可以添加验证 API Key 等操作
	if o.APIKey == "" {
		return errors.New("OpenAI API key is required")
	}
	return nil
}

func (d *OpenAIClient) StreamChatCompletion(ctx context.Context, messages []ChatMessage, tools []Tool, chanStream chan<- request.EventSSEData) (toolCallsResp []ToolCall, fullMessage string) {

	return nil, ""
}
