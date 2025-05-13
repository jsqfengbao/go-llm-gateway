package llmGateway

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
	"io"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"go.uber.org/zap"
)

// DoubaoClient 是豆包大模型客户端的结构体
type DoubaoClient struct {
	APIKey string
	Model  string
	// 这里可以添加豆包客户端的其他配置
}

func (d *DoubaoClient) Init(ctx context.Context) error {
	// 初始化豆包客户端
	// 这里可以添加验证 API Key 等操作
	if d.APIKey == "" {
		return errors.New("Doubao API key is required")
	}
	return nil
}

// buildToolDefinitions 工具转化为标准定义
func buildToolDefinitionsForDoubao(tools []Tool) []*model.Tool {
	var result []*model.Tool
	for _, t := range tools {
		result = append(result, &model.Tool{
			Type: "function",
			Function: &model.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		})
	}
	return result
}

func (d *DoubaoClient) StreamChatCompletion(ctx context.Context, messages []ChatMessage, tools []Tool, chanStream chan<- request.EventSSEData) (toolCallsResp []ToolCall, fullMessage string) {
	toolDefs := buildToolDefinitionsForDoubao(tools)
	var (
		toolCallCount int
	)
	systemPrompt := `你是一个多智能体系统，拥有不同工具能力。请遵守以下原则：
	1. 主体职责是理解用户意图并尽可能使用自然语言回答。
	2. 如果用户问题能通过模型知识直接回答，请优先自然语言回答，而不是调用agent。
	3. 严禁推测性调用工具。
	4. 工具调用流程需符合常识和用户授权。`

	var chatMessages []*model.ChatCompletionMessage
	for _, m := range messages {
		if m.Role == "user" {
			chatMessages = append(chatMessages, &model.ChatCompletionMessage{
				Role: "system",
				Content: &model.ChatCompletionMessageContent{
					StringValue: &systemPrompt,
				},
			})
			chatMessages = append(chatMessages, &model.ChatCompletionMessage{
				Role: m.Role,
				Content: &model.ChatCompletionMessageContent{
					StringValue: &m.Content,
				},
			})
		} else if m.Role == "assistant" {
			chatMessages = append(chatMessages, &model.ChatCompletionMessage{
				Role: m.Role,
				Content: &model.ChatCompletionMessageContent{
					StringValue: &m.Content,
				},
				ToolCalls: convertFromLLMToolcallToDoubao(m.ToolCalls),
			})
		} else if m.Role == "tool" {
			chatMessages = append(chatMessages, &model.ChatCompletionMessage{
				Role: m.Role,
				Content: &model.ChatCompletionMessageContent{
					StringValue: &m.Content,
				},
				Name:       &m.Name,
				ToolCallID: m.ToolCallId,
			})
		}

	}

	client := arkruntime.NewClientWithApiKey(d.APIKey)

	streamFlag := true
	input := &model.CreateChatCompletionRequest{
		Model:    d.Model,
		Messages: chatMessages,
		Stream:   &streamFlag,
		Tools:    toolDefs,
	}
	global.GVA_LOG.Info("豆包请求参数", zap.Any("input", input))

	resp, err := client.CreateChatCompletionStream(ctx, input)
	if err != nil {
		global.GVA_LOG.Error("Ark ChatWithStream failed", zap.Error(err))
		return nil, ""
	}
	defer resp.Response.Body.Close()

	reader := bufio.NewReader(resp.Response.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			global.GVA_LOG.Error("Error reading ark stream", zap.Error(err))
			break
		}
		global.GVA_LOG.Info("豆包流式响应", zap.String("line", line))

		// 去除 "data: " 前缀
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}
		line = strings.TrimSpace(line)
		if line == "" || line == "\n" || line == "\r" || line == "[DONE]" {
			continue
		}

		var chunk model.ChatCompletionStreamResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			global.GVA_LOG.Error("Unmarshal chunk failed", zap.String("line", line), zap.Error(err))
			continue
		}

		// 处理内容
		if chunk.Choices != nil && len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				fullMessage += content
				chanStream <- request.EventSSEData{Delta: content}
			}

			if len(chunk.Choices[0].Delta.ToolCalls) > 0 {
				toolCall := chunk.Choices[0].Delta.ToolCalls[0]
				if toolCall.ID != "" {
					toolCallCount++
					toolCall := ToolCall{
						ID:        toolCall.ID,
						Name:      toolCall.Function.Name,
						Arguments: toolCall.Function.Arguments,
					}
					toolCallsResp = append(toolCallsResp, toolCall)
				}
				if toolCall.Function.Arguments != "" {
					toolCallsResp[toolCallCount-1].Arguments += toolCall.Function.Arguments
				}
			}
		}
	}

	return toolCallsResp, fullMessage
}

func convertFromLLMToolcallToDoubao(toolCalls []ToolCall) []*model.ToolCall {
	var result []*model.ToolCall
	for _, tc := range toolCalls {
		result = append(result, &model.ToolCall{
			ID:   tc.ID,
			Type: "function",
			Function: model.FunctionCall{
				Name:      tc.Name,
				Arguments: tc.Arguments,
			},
		})
	}
	return result
}
