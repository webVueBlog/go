package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-llm-tools/internal/chain"
	"go-llm-tools/internal/llm"
	"go-llm-tools/internal/prompt"
	"go-llm-tools/internal/rag"
	"go-llm-tools/internal/utils"
)

func main() {
	// 示例 1: 基础 LLM 调用
	exampleBasicLLM()

	// 示例 2: Prompt 模板使用
	examplePromptTemplates()

	// 示例 3: RAG 功能
	exampleRAG()

	// 示例 4: 链式调用
	exampleChain()

	// 示例 5: 完整工作流
	exampleCompleteWorkflow()
}

func exampleBasicLLM() {
	fmt.Println("\n=== 示例 1: 基础 LLM 调用 ===")

	// 加载配置
	config, err := utils.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	// 创建 LLM 提供者
	llmConfig := &llm.Config{
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		MaxTokens:   1000,
		Timeout:     30 * time.Second,
	}

	provider := llm.NewOpenAIProvider(llmConfig)

	// 创建聊天请求
	req := &llm.ChatRequest{
		Model: llmConfig.Model,
		Messages: []llm.Message{
			{Role: "user", Content: "请简单介绍一下 Go 语言的特点。"},
		},
		Temperature: llmConfig.Temperature,
		MaxTokens:   llmConfig.MaxTokens,
	}

	// 调用 LLM
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := provider.Chat(ctx, req)
	if err != nil {
		log.Printf("LLM call failed: %v", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Printf("回答: %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("Token 使用: %d\n", resp.Usage.TotalTokens)
	}
}

func examplePromptTemplates() {
	fmt.Println("\n=== 示例 2: Prompt 模板使用 ===")

	// 创建 Prompt 引擎
	engine := prompt.NewPromptEngine()

	// 添加自定义模板
	tmpl := &prompt.Template{
		Name:    "code_review",
		Content: "请对以下 {{.language}} 代码进行审查：\n\n```{{.language}}\n{{.code}}\n```\n\n请从代码质量、安全性、性能等方面进行评估。",
		Version: "1.0",
	}

	if err := engine.AddTemplate(tmpl); err != nil {
		log.Printf("Failed to add template: %v", err)
		return
	}

	// 渲染模板
	data := map[string]interface{}{
		"language": "go",
		"code": `func add(a, b int) int {
    return a + b
}`,
	}

	prompt, err := engine.Render("code_review", data)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		return
	}

	fmt.Printf("生成的 Prompt:\n%s\n", prompt)
}

func exampleRAG() {
	fmt.Println("\n=== 示例 3: RAG 功能 ===")

	// 创建检索器
	retriever := rag.NewSimpleRetriever()

	// 添加文档
	docs := []*rag.Document{
		{
			ID:      "doc1",
			Content: "Go 语言是由 Google 开发的开源编程语言，具有简洁、高效、并发安全等特点。",
			Metadata: map[string]string{
				"source": "go_docs",
				"type":   "language",
			},
		},
		{
			ID:      "doc2",
			Content: "Go 语言支持并发编程，通过 goroutine 和 channel 实现轻量级线程和通信。",
			Metadata: map[string]string{
				"source": "go_concurrency",
				"type":   "feature",
			},
		},
	}

	for _, doc := range docs {
		if err := retriever.AddDocument(context.Background(), doc); err != nil {
			log.Printf("Failed to add document: %v", err)
		}
	}

	// 创建 RAG 引擎
	ragEngine := rag.NewRAGEngine(retriever)

	// 执行查询
	query := "Go 语言的并发特性"
	result, err := ragEngine.Query(context.Background(), query, 3)
	if err != nil {
		log.Printf("RAG query failed: %v", err)
		return
	}

	fmt.Printf("查询: %s\n", query)
	fmt.Printf("结果: %s\n", result)
}

func exampleChain() {
	fmt.Println("\n=== 示例 4: 链式调用 ===")

	// 创建链式调用
	c := chain.NewChain()

	// 添加自定义步骤
	c.AddStep(func(ctx context.Context, input interface{}) (interface{}, error) {
		if str, ok := input.(string); ok {
			// 步骤 1: 预处理
			return "预处理: " + str, nil
		}
		return input, nil
	})

	c.AddStep(func(ctx context.Context, input interface{}) (interface{}, error) {
		if str, ok := input.(string); ok {
			// 步骤 2: 增强
			return str + " [已增强]", nil
		}
		return input, nil
	})

	c.AddStep(func(ctx context.Context, input interface{}) (interface{}, error) {
		if str, ok := input.(string); ok {
			// 步骤 3: 后处理
			return "最终结果: " + str, nil
		}
		return input, nil
	})

	// 执行链式调用
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := c.RunString(ctx, "测试输入")
	if err != nil {
		log.Printf("Chain execution failed: %v", err)
		return
	}

	fmt.Printf("链式调用结果: %s\n", result)
	fmt.Printf("步骤数量: %d\n", c.GetStepCount())
}

func exampleCompleteWorkflow() {
	fmt.Println("\n=== 示例 5: 完整工作流 ===")

	// 加载配置
	config, err := utils.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	// 创建组件
	llmConfig := &llm.Config{
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		MaxTokens:   1000,
		Timeout:     30 * time.Second,
	}

	provider := llm.NewOpenAIProvider(llmConfig)
	promptEngine := prompt.NewPromptEngine()
	retriever := rag.NewSimpleRetriever()
	ragEngine := rag.NewRAGEngine(retriever)

	// 添加默认模板
	for name, tmpl := range prompt.DefaultTemplates {
		if err := promptEngine.AddTemplate(tmpl); err != nil {
			log.Printf("Warning: failed to add template %s: %v", name, err)
		}
	}

	// 添加示例文档
	addSampleDocuments(retriever)

	// 创建完整工作流
	query := "什么是 LangChain？"
	fmt.Printf("查询: %s\n", query)

	// 步骤 1: RAG 检索
	enhancedQuery, err := ragEngine.Query(context.Background(), query, 3)
	if err != nil {
		log.Printf("RAG query failed: %v", err)
		return
	}

	// 步骤 2: 渲染 Prompt
	data := map[string]interface{}{
		"question": enhancedQuery,
	}

	prompt, err := promptEngine.Render("qa", data)
	if err != nil {
		log.Printf("Failed to render prompt: %v", err)
		return
	}

	// 步骤 3: 调用 LLM
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &llm.ChatRequest{
		Model:       llmConfig.Model,
		Messages:    []llm.Message{{Role: "user", Content: prompt}},
		Temperature: llmConfig.Temperature,
		MaxTokens:   llmConfig.MaxTokens,
	}

	resp, err := provider.Chat(ctx, req)
	if err != nil {
		log.Printf("LLM call failed: %v", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Printf("回答: %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("Token 使用: %d\n", resp.Usage.TotalTokens)
	}
}

func addSampleDocuments(retriever *rag.SimpleRetriever) {
	docs := []*rag.Document{
		{
			ID:      "doc1",
			Content: "LangChain 是一个用于开发由语言模型驱动的应用程序的框架。它提供了模块化的组件和预构建的链，使开发人员能够快速构建复杂的应用程序。",
			Metadata: map[string]string{
				"source": "langchain_docs",
				"type":   "framework",
			},
		},
		{
			ID:      "doc2",
			Content: "RAG (Retrieval-Augmented Generation) 是一种结合了信息检索和文本生成的技术。它首先从知识库中检索相关信息，然后使用这些信息来生成更准确、更相关的回答。",
			Metadata: map[string]string{
				"source": "rag_paper",
				"type":   "technique",
			},
		},
	}

	for _, doc := range docs {
		if err := retriever.AddDocument(context.Background(), doc); err != nil {
			log.Printf("Warning: failed to add document %s: %v", doc.ID, err)
		}
	}
} 