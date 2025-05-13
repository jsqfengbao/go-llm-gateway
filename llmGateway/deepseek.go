package llmGateway

import (
	"context"
	"errors"
	"github.com/cohesion-org/deepseek-go"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
	"go.uber.org/zap"
	"io"
)

// DeepSeekClient 是deepseek大模型客户端的结构体
type DeepSeekClient struct {
	APIKey string
	// 这里可以添加豆包客户端的其他配置
}

func (d *DeepSeekClient) Init(ctx context.Context) error {
	// 如果你需要初始化配置（如加载 API 密钥等），可以放这里
	// 当前 DeepSeek SDK 应该已经初始化了 client，所以无需特别处理
	if d.APIKey == "" {
		return errors.New("Deepseek API key is required")
	}
	return nil
}
func (d *DeepSeekClient) StreamChatCompletion(
	ctx context.Context,
	messages []ChatMessage,
	tools []Tool,
	chanStream chan<- request.EventSSEData,
) (toolCallsResp []ToolCall, fullMessage string) {

	var toolCallCount int

	systemPrompt := `你是一个多智能体系统，拥有不同工具能力。请遵守以下原则：
	1. 主体职责是理解用户意图并尽可能使用自然语言回答。
	2. 如果用户问题能通过模型知识直接回答，请优先自然语言回答，而不是调用agent。
	3. 严禁推测性调用工具。
	4. 工具调用流程需符合常识和用户授权。`

	// 构建 deepseek-go 所需的消息结构
	var chatMessages []deepseek.ChatCompletionMessage
	for _, m := range messages {
		if m.Role == "user" {
			chatMessages = append(chatMessages, deepseek.ChatCompletionMessage{
				Role:    "system",
				Content: systemPrompt,
			})
			chatMessages = append(chatMessages, deepseek.ChatCompletionMessage{
				Role:    "user",
				Content: m.Content,
			})
		} else if m.Role == "assistant" {
			chatMessages = append(chatMessages, deepseek.ChatCompletionMessage{
				Role:      "assistant",
				Content:   m.Content,
				ToolCalls: convertToDeepseekToolCalls(m.ToolCalls),
			})
		} else if m.Role == "tool" {
			chatMessages = append(chatMessages, deepseek.ChatCompletionMessage{
				Role:       "tool",
				Content:    m.Content,
				ToolCallID: m.ToolCallId,
				ToolCalls:  convertToDeepseekToolCalls(m.ToolCalls),
			})
		}
	}

	client := deepseek.NewClient(d.APIKey)
	toolDefs := buildToolDefinitionsForDeepSeek(tools)
	streamReq := &deepseek.StreamChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: chatMessages,
		Stream:   true,
		Tools:    toolDefs,
	}

	stream, err := client.CreateChatCompletionStream(ctx, streamReq)
	if err != nil {
		global.GVA_LOG.Error("DeepSeek ChatCompletionStream 失败", zap.Error(err))
		return nil, ""
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			global.GVA_LOG.Error("读取 DeepSeek 流失败", zap.Error(err))
			break
		}

		if len(resp.Choices) == 0 {
			continue
		}

		delta := resp.Choices[0].Delta
		if delta.Content != "" {
			fullMessage += delta.Content
			chanStream <- request.EventSSEData{Delta: delta.Content}
		}

		if len(delta.ToolCalls) > 0 {
			tc := delta.ToolCalls[0]
			if tc.ID != "" {
				toolCallCount++
				newCall := ToolCall{
					ID:        tc.ID,
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				}
				toolCallsResp = append(toolCallsResp, newCall)
			}
			if tc.Function.Arguments != "" {
				toolCallsResp[toolCallCount-1].Arguments += tc.Function.Arguments
			}
		}
	}

	return toolCallsResp, fullMessage
}

func buildToolDefinitionsForDeepSeek(tools []Tool) []deepseek.Tool {
	var result []deepseek.Tool
	for _, t := range tools {
		result = append(result, deepseek.Tool{
			Type: "function",
			Function: deepseek.Function{
				Name:        t.Name,
				Description: t.Description,
				Parameters: &deepseek.FunctionParameters{
					Type:       t.Parameters["type"].(string),
					Properties: t.Parameters["properties"].(map[string]interface{}),
					Required:   t.Parameters["required"].([]string),
				},
			},
		})
	}
	return result
}

func convertToDeepseekToolCalls(toolCalls []ToolCall) []deepseek.ToolCall {
	var result []deepseek.ToolCall
	for _, tc := range toolCalls {
		result = append(result, deepseek.ToolCall{
			ID:   tc.ID,
			Type: "function",
			Function: deepseek.ToolCallFunction{
				Name:      tc.Name,
				Arguments: tc.Arguments,
			},
		})
	}
	return result
}
