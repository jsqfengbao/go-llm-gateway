package llmGateway

import (
	"context"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

// LLMType 定义支持的LLM类型
type LLMType string

const (
	DeepSeekLLM LLMType = "deepseek"
	OpenAILLM   LLMType = "openai"
	DoubaoLLM   LLMType = "doubao"
	// 可以添加更多的LLM类型
)

// LLMFactory 创建LLM客户端的工厂
type LLMFactory struct{}

// NewLLMFactory 创建一个新的LLM工厂
func NewLLMFactory() *LLMFactory {
	return &LLMFactory{}
}

// CreateLLMClient 根据类型创建对应的LLM客户端
func (f *LLMFactory) CreateLLMClient(ctx context.Context, llmType LLMType) (LLMClient, error) {
	switch llmType {
	case DeepSeekLLM:
		return f.createDeepSeekClient(ctx)
	case OpenAILLM:
		return f.createOpenAIClient(ctx)
	case DoubaoLLM:
		return f.createDoubaoClient(ctx)
	default:
		return nil, errors.New("unsupported LLM type")
	}
}

// createDeepSeekClient 创建DeepSeek客户端
func (f *LLMFactory) createDeepSeekClient(ctx context.Context) (LLMClient, error) {
	param, err := system.SysParamsServiceApp.GetSysParam("deepseek_apiKey")
	if err != nil {
		global.GVA_LOG.Error("Failed to get deepseek_apiKey", zap.Error(err))
		return nil, err
	}

	client := &DeepSeekClient{
		APIKey: param.Value,
	}

	if err := client.Init(ctx); err != nil {
		global.GVA_LOG.Error("Failed to initialize DeepSeek client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

// createOpenAIClient 创建OpenAI客户端
func (f *LLMFactory) createOpenAIClient(ctx context.Context) (LLMClient, error) {
	param, err := system.SysParamsServiceApp.GetSysParam("openai_apiKey")
	if err != nil {
		global.GVA_LOG.Error("Failed to get openai_apiKey", zap.Error(err))
		return nil, err
	}

	modelParam, err := system.SysParamsServiceApp.GetSysParam("openai_model")
	if err != nil {
		global.GVA_LOG.Warn("Failed to get openai_model, using default", zap.Error(err))
		modelParam.Value = "gpt-4-turbo"
	}

	client := &OpenAIClient{
		APIKey: param.Value,
		Model:  modelParam.Value,
	}

	if err := client.Init(ctx); err != nil {
		global.GVA_LOG.Error("Failed to initialize OpenAI client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

// createDoubaoClient 创建豆包客户端
func (f *LLMFactory) createDoubaoClient(ctx context.Context) (LLMClient, error) {
	param, err := system.SysParamsServiceApp.GetSysParam("doubao_apiKey")
	if err != nil {
		global.GVA_LOG.Error("Failed to get doubao_apiKey", zap.Error(err))
		return nil, err
	}

	modelParam, err := system.SysParamsServiceApp.GetSysParam("doubao_model")
	if err != nil {
		global.GVA_LOG.Error("Failed to get doubao_model, using default", zap.Error(err))
		modelParam.Value = "doubao-lite"
	}

	client := &DoubaoClient{
		APIKey: param.Value,
		Model:  modelParam.Value,
	}

	if err := client.Init(ctx); err != nil {
		global.GVA_LOG.Error("Failed to initialize Doubao client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

// GetDefaultLLMClient 获取默认的LLM客户端
func (f *LLMFactory) GetDefaultLLMClient(ctx context.Context) (LLMClient, error) {
	// 从系统参数中获取默认LLM类型
	param, err := system.SysParamsServiceApp.GetSysParam("default_llm_type")
	if err != nil {
		// 如果没有设置，默认使用DeepSeek
		global.GVA_LOG.Warn("Failed to get default_llm_type, using DeepSeek as default", zap.Error(err))
		return f.CreateLLMClient(ctx, DeepSeekLLM)
	}
	return f.CreateLLMClient(ctx, LLMType(param.Value))
}
