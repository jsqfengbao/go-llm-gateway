# LLM Gateway 抽象层

这个模块提供了一个抽象层，用于与不同的大型语言模型（LLM）进行交互。目前支持 DeepSeek 模型，并可以轻松扩展以支持其他 LLM。

## 架构

LLM Gateway 使用了以下设计模式：

1. **接口抽象**：通过 `LLMClient` 接口定义与 LLM 交互的标准方法
2. **工厂模式**：使用 `LLMFactory` 创建不同类型的 LLM 客户端
3. **适配器模式**：每个具体的 LLM 客户端实现都适配特定 LLM 的 API

## 主要组件

### LLMClient 接口

定义了与 LLM 交互的标准方法：

```go
type LLMClient interface {
    Init(ctx context.Context) error
    ChatCompletionStream(ctx context.Context, messages []ChatMessage, tools []Tool) (<-chan CompletionResult, error)
}
```

### LLM 工厂

`LLMFactory` 负责创建不同类型的 LLM 客户端：

```go
func (f *LLMFactory) CreateLLMClient(ctx context.Context, llmType LLMType) (LLMClient, error)
func (f *LLMFactory) GetDefaultLLMClient(ctx context.Context) (LLMClient, error)
```

### 具体实现

- **DeepSeekClient**：DeepSeek 模型的客户端实现
- **OpenAIClient**：OpenAI 模型的客户端实现（示例，需要完善）
- **DoubaoClient**：豆包大模型的客户端实现

## 如何使用

### 1. 配置系统参数

在系统参数中设置以下参数：

- `deepseek_apiKey`：DeepSeek API 密钥
- `openai_apiKey`：OpenAI API 密钥
- `doubao_apiKey`：豆包 API 密钥
- `default_llm_type`：默认使用的 LLM 类型（如 "deepseek"、"openai" 或 "doubao"）
- `use_llm_abstraction`：是否使用 LLM 抽象层（"true" 或 "false"）

### 2. 使用 LLM 工厂创建客户端

```go
// 创建 LLM 工厂
llmFactory := llmGateway.NewLLMFactory()

// 获取默认的 LLM 客户端
llmClient, err := llmFactory.GetDefaultLLMClient(context.Background())
if err != nil {
    // 处理错误
}

// 或者指定 LLM 类型
llmClient, err := llmFactory.CreateLLMClient(context.Background(), llmGateway.DeepSeekLLM)
if err != nil {
    // 处理错误
}
```

### 3. 使用 LLM 客户端

```go
// 准备消息
messages := []llmGateway.ChatMessage{
    {
        Role:    "user",
        Content: "你好，请帮我查询北京的天气",
    },
}

// 准备工具
tools := []llmGateway.Tool{
    {
        Name:        "get_weather",
        Description: "获取城市天气",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "city": map[string]interface{}{
                    "type":        "string",
                    "description": "城市名称",
                },
            },
            "required": []string{"city"},
        },
    },
}

// 调用流式聊天完成
stream, err := llmClient.ChatCompletionStream(ctx, messages, tools)
if err != nil {
    // 处理错误
}

// 处理流式响应
for result := range stream {
    if result.Message != "" {
        // 处理消息
        fmt.Println(result.Message)
    }
    
    // 处理工具调用
    for _, toolCall := range result.ToolCalls {
        // 处理工具调用
        fmt.Printf("工具调用: %s, 参数: %s\n", toolCall.Name, toolCall.Arguments)
    }
}
```

## 添加新的 LLM 支持

要添加新的 LLM 支持，需要以下步骤：

1. 创建新的 LLM 客户端结构体，实现 `LLMClient` 接口
2. 在 `LLMType` 中添加新的 LLM 类型
3. 在 `LLMFactory.CreateLLMClient` 方法中添加新的 case

例如，添加 Anthropic Claude 支持：

```go
// 1. 定义新的 LLM 类型
const (
    DeepSeekLLM  LLMType = "deepseek"
    OpenAILLM    LLMType = "openai"
    DoubaoLLM    LLMType = "doubao"
    AnthropicLLM LLMType = "anthropic"  // 新增
)

// 2. 创建 Anthropic 客户端结构体
type AnthropicClient struct {
    APIKey string
    Model  string
}

// 3. 实现 LLMClient 接口的方法
func (a *AnthropicClient) Init(ctx context.Context) error {
    // 初始化实现
}

func (a *AnthropicClient) GetTools(ctx context.Context) ([]Tool, error) {
    // 获取工具实现
}

func (a *AnthropicClient) ChatCompletionStream(ctx context.Context, messages []ChatMessage, tools []Tool) (<-chan CompletionResult, error) {
    // 流式聊天实现
}

// 4. 在工厂中添加新的 case
func (f *LLMFactory) CreateLLMClient(ctx context.Context, llmType LLMType) (LLMClient, error) {
    switch llmType {
    case DeepSeekLLM:
        return f.createDeepSeekClient(ctx)
    case OpenAILLM:
        return f.createOpenAIClient(ctx)
    case DoubaoLLM:
        return f.createDoubaoClient(ctx)
    case AnthropicLLM:
        return f.createAnthropicClient(ctx)  // 新增
    default:
        return nil, errors.New("unsupported LLM type")
    }
}

// 5. 添加创建方法
func (f *LLMFactory) createAnthropicClient(ctx context.Context) (LLMClient, error) {
    param, err := system.SysParamsServiceApp.GetSysParam("anthropic_apiKey")
    if err != nil {
        return nil, err
    }
    
    client := &AnthropicClient{
        APIKey: param.Value,
        Model:  "claude-3-opus-20240229",  // 默认模型
    }
    
    if err := client.Init(ctx); err != nil {
        return nil, err
    }
    
    return client, nil
}
