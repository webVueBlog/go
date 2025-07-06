package llm

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider OpenAI 提供者
type OpenAIProvider struct {
	client *openai.Client
	config *Config
}

// NewOpenAIProvider 创建 OpenAI 提供者
func NewOpenAIProvider(config *Config) *OpenAIProvider {
	if config == nil {
		config = &Config{
			BaseURL:     "https://api.openai.com/v1",
			Model:       "gpt-3.5-turbo",
			Timeout:     30 * time.Second,
			MaxRetries:  3,
			Temperature: 0.7,
			MaxTokens:   1000,
		}
	}

	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	return &OpenAIProvider{
		client: openai.NewClientWithConfig(clientConfig),
		config: config,
	}
}

// Chat 实现聊天接口
func (p *OpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// 转换消息格式
	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 构建请求
	completionReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
		TopP:        float32(req.TopP),
		Stream:      req.Stream,
	}

	// 调用 API
	resp, err := p.client.CreateChatCompletion(ctx, completionReq)
	if err != nil {
		return nil, fmt.Errorf("openai chat completion failed: %w", err)
	}

	// 转换响应格式
	chatResp := &ChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: make([]struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		}, len(resp.Choices)),
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}

	for i, choice := range resp.Choices {
		chatResp.Choices[i].Index = choice.Index
		chatResp.Choices[i].Message.Role = choice.Message.Role
		chatResp.Choices[i].Message.Content = choice.Message.Content
		chatResp.Choices[i].FinishReason = string(choice.FinishReason)
	}

	return chatResp, nil
}

// Complete 实现补全接口
func (p *OpenAIProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// 构建请求
	completionReq := openai.CompletionRequest{
		Model:       req.Model,
		Prompt:      req.Prompt,
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
		TopP:        float32(req.TopP),
		Stop:        req.Stop,
	}

	// 调用 API
	resp, err := p.client.CreateCompletion(ctx, completionReq)
	if err != nil {
		return nil, fmt.Errorf("openai completion failed: %w", err)
	}

	// 转换响应格式
	completionResp := &CompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: make([]struct {
			Text         string      `json:"text"`
			Index        int         `json:"index"`
			Logprobs     interface{} `json:"logprobs"`
			FinishReason string      `json:"finish_reason"`
		}, len(resp.Choices)),
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}

	for i, choice := range resp.Choices {
		completionResp.Choices[i].Text = choice.Text
		completionResp.Choices[i].Index = choice.Index
		completionResp.Choices[i].Logprobs = choice.LogProbs
		completionResp.Choices[i].FinishReason = string(choice.FinishReason)
	}

	return completionResp, nil
}

// GetConfig 获取配置
func (p *OpenAIProvider) GetConfig() *Config {
	return p.config
}

// SetConfig 设置配置
func (p *OpenAIProvider) SetConfig(config *Config) {
	p.config = config
	if config != nil {
		clientConfig := openai.DefaultConfig(config.APIKey)
		if config.BaseURL != "" {
			clientConfig.BaseURL = config.BaseURL
		}
		p.client = openai.NewClientWithConfig(clientConfig)
	}
}

// GetModelType 获取模型类型
func (p *OpenAIProvider) GetModelType() ModelType {
	return ModelTypeOpenAI
} 