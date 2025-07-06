package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	_ "os"
	"time"

	"go-llm-tools/internal/chain"
	"go-llm-tools/internal/llm"
	"go-llm-tools/internal/prompt"
	"go-llm-tools/internal/rag"
	"go-llm-tools/internal/utils"
)

func main() {
	// 解析命令行参数
	var (
		query     = flag.String("query", "", "查询内容")
		template  = flag.String("template", "qa", "使用的 Prompt 模板")
		model     = flag.String("model", "gpt-3.5-turbo", "使用的模型")
		apiKey    = flag.String("api-key", "", "OpenAI API Key")
		baseURL   = flag.String("base-url", "", "OpenAI Base URL")
		verbose   = flag.Bool("verbose", false, "详细输出")
		chainMode = flag.Bool("chain", false, "使用链式调用模式")
	)
	flag.Parse()

	// 加载配置
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 验证配置
	if err := utils.ValidateConfig(config); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// 使用命令行参数覆盖配置
	if *apiKey != "" {
		config.OpenAIAPIKey = *apiKey
	}
	if *baseURL != "" {
		config.OpenAIBaseURL = *baseURL
	}

	// 创建 LLM 提供者
	llmConfig := &llm.Config{
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Model:       *model,
		Temperature: config.OpenAITemperature,
		MaxTokens:   config.OpenAIMaxTokens,
		Timeout:     config.RequestTimeout,
	}

	provider := llm.NewOpenAIProvider(llmConfig)

	// 创建 Prompt 引擎
	promptEngine := prompt.NewPromptEngine()

	// 添加默认模板
	for name, tmpl := range prompt.DefaultTemplates {
		if err := promptEngine.AddTemplate(tmpl); err != nil {
			log.Printf("Warning: failed to add template %s: %v", name, err)
		}
	}

	// 创建 RAG 引擎
	retriever := rag.NewSimpleRetriever()
	ragEngine := rag.NewRAGEngine(retriever)

	// 添加一些示例文档
	addSampleDocuments(retriever)

	if *chainMode {
		// 链式调用模式
		runChainMode(provider, promptEngine, ragEngine, *query, *template, *verbose)
	} else {
		// 简单模式
		runSimpleMode(provider, promptEngine, *query, *template, *verbose)
	}
}

func runChainMode(provider llm.Provider, promptEngine *prompt.PromptEngine, ragEngine *rag.RAGEngine, query, templateName string, verbose bool) {
	if query == "" {
		fmt.Println("请输入查询内容 (使用 -query 参数)")
		return
	}

	// 创建链式调用
	c := chain.NewChain()

	// 添加步骤：检索 -> 构建 Prompt -> 调用 LLM
	c.AddStep(rag.Retrieve)
	c.AddStep(prompt.BuildPrompt)
	c.AddStep(func(ctx context.Context, input interface{}) (interface{}, error) {
		if str, ok := input.(string); ok {
			// 调用 LLM
			req := &llm.ChatRequest{
				Model:       provider.GetConfig().Model,
				Messages:    []llm.Message{{Role: "user", Content: str}},
				Temperature: provider.GetConfig().Temperature,
				MaxTokens:   provider.GetConfig().MaxTokens,
			}

			resp, err := provider.Chat(ctx, req)
			if err != nil {
				return nil, err
			}

			if len(resp.Choices) > 0 {
				return resp.Choices[0].Message.Content, nil
			}

			return "抱歉，没有获得有效回复。", nil
		}
		return input, nil
	})

	// 执行链式调用
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := c.RunString(ctx, query)
	if err != nil {
		log.Fatalf("Chain execution failed: %v", err)
	}

	fmt.Printf("查询: %s\n", query)
	fmt.Printf("结果: %s\n", result)
}

func runSimpleMode(provider llm.Provider, promptEngine *prompt.PromptEngine, query, templateName string, verbose bool) {
	if query == "" {
		fmt.Println("请输入查询内容 (使用 -query 参数)")
		return
	}

	// 渲染 Prompt 模板
	data := map[string]interface{}{
		"question": query,
	}

	prompt, err := promptEngine.Render(templateName, data)
	if err != nil {
		log.Fatalf("Failed to render prompt: %v", err)
	}

	if verbose {
		fmt.Printf("使用的模板: %s\n", templateName)
		fmt.Printf("生成的 Prompt: %s\n", prompt)
	}

	// 调用 LLM
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &llm.ChatRequest{
		Model:       provider.GetConfig().Model,
		Messages:    []llm.Message{{Role: "user", Content: prompt}},
		Temperature: provider.GetConfig().Temperature,
		MaxTokens:   provider.GetConfig().MaxTokens,
	}

	resp, err := provider.Chat(ctx, req)
	if err != nil {
		log.Fatalf("LLM call failed: %v", err)
	}

	if len(resp.Choices) > 0 {
		fmt.Printf("查询: %s\n", query)
		fmt.Printf("回答: %s\n", resp.Choices[0].Message.Content)

		if verbose {
			fmt.Printf("Token 使用: %d\n", resp.Usage.TotalTokens)
		}
	} else {
		fmt.Println("抱歉，没有获得有效回复。")
	}
}

func addSampleDocuments(retriever *rag.SimpleRetriever) {
	// 添加一些示例文档
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
		{
			ID:      "doc3",
			Content: "Prompt Engineering 是设计和优化提示词的艺术和科学，目的是从语言模型中获得更好的输出。它包括理解模型的局限性、设计有效的提示词模板等。",
			Metadata: map[string]string{
				"source": "prompt_engineering_guide",
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
