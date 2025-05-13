# LLM Gateway Abstraction Layer

[中文文档](README_zh.md)

This module provides an abstraction layer for interacting with different Large Language Models (LLMs). Currently supports DeepSeek models, with easy extensibility for other LLMs.

## Detail URL： https://plugin.gin-vue-admin.com/details/116


## Architecture

The LLM Gateway uses the following design patterns:
1. **Interface Abstraction**: Defines standard methods for LLM interaction through the `LLMClient` interface
2. **Factory Pattern**: Uses `LLMFactory` to create different types of LLM clients
3. **Adapter Pattern**: Each concrete LLM client implementation adapts to a specific LLM's API

## Main Components

### LLMClient Interface

Defines standard methods for LLM interaction:
```go
type LLMClient interface {
    Init(ctx context.Context) error
    ChatCompletionStream(ctx context.Context, messages []ChatMessage, tools []Tool) (<-chan CompletionResult, error)
}
// Create LLM factory
llmFactory := llmGateway.NewLLMFactory()

// Get default client
llmClient, err := llmFactory.GetDefaultLLMClient(context.Background())
if err != nil {
    // Handle error
}

// Or specify LLM type
llmClient, err := llmFactory.CreateLLMClient(context.Background(), llmGateway.DeepSeekLLM)
if err != nil {
    // Handle error
}

// Prepare messages
messages := []llmGateway.ChatMessage{
    {
        Role:    "user",
        Content: "Please help check the weather in Beijing",
    },
}

// Prepare tools
tools := []llmGateway.Tool{
    {
        Name:        "get_weather",
        Description: "Get city weather information",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "city": map[string]interface{}{
                    "type":        "string",
                    "description": "City name",
                },
            },
            "required": []string{"city"},
        },
    },
}

// Stream chat completion
stream, err := llmClient.ChatCompletionStream(ctx, messages, tools)
if err != nil {
    // Handle error
}

// Process stream
for result := range stream {
    if result.Message != "" {
        fmt.Println(result.Message)
    }
    
    for _, toolCall := range result.ToolCalls {
        fmt.Printf("Tool call: %s, Arguments: %s\n", toolCall.Name, toolCall.Arguments)
    }
}

