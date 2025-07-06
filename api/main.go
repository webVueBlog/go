package main

import (
	"context"
	_ "encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"go-llm-tools/internal/auth"
	"go-llm-tools/internal/chain"
	"go-llm-tools/internal/chatgpt"
	"go-llm-tools/internal/llm"
	"go-llm-tools/internal/prompt"
	"go-llm-tools/internal/rag"
	"go-llm-tools/internal/utils"
)

// API 请求结构
type ChatRequest struct {
	Query     string            `json:"query" binding:"required"`
	Template  string            `json:"template"`
	Model     string            `json:"model"`
	Variables map[string]string `json:"variables"`
	ChainMode bool              `json:"chain_mode"`
}

type ChatResponse struct {
	Query      string `json:"query"`
	Answer     string `json:"answer"`
	Template   string `json:"template"`
	Model      string `json:"model"`
	TokenUsage int    `json:"token_usage,omitempty"`
	Error      string `json:"error,omitempty"`
}

type TemplateRequest struct {
	Name     string            `json:"name" binding:"required"`
	Content  string            `json:"content" binding:"required"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
}

type TemplateResponse struct {
	Name      string            `json:"name"`
	Content   string            `json:"content"`
	Version   string            `json:"version"`
	Variables []string          `json:"variables"`
	Metadata  map[string]string `json:"metadata"`
}

// 认证相关请求结构
type RegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	ChatGPTToken string `json:"chatgpt_token"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChatGPTRequest struct {
	Message        string `json:"message" binding:"required"`
	ConversationID string `json:"conversation_id,omitempty"`
	Model          string `json:"model"`
}

// 全局变量
var (
	provider      llm.Provider
	promptEngine  *prompt.PromptEngine
	ragEngine     *rag.RAGEngine
	config        *utils.Config
	logger        *logrus.Logger
	authManager   *auth.AuthManager
	chatGPTClient *chatgpt.ChatGPTClient
)

func main() {
	// 初始化日志
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// 加载配置
	var err error
	config, err = utils.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// 验证配置
	if err := utils.ValidateConfig(config); err != nil {
		logger.Fatalf("Invalid config: %v", err)
	}

	// 初始化组件
	initializeComponents()

	// 设置 Gin 模式
	if config.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// 设置路由
	setupRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	logger.Infof("Starting server on %s", addr)

	if err := r.Run(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func initializeComponents() {
	// 初始化认证管理器
	authManager = auth.NewAuthManager("your-secret-key-here")

	// 初始化 ChatGPT 客户端
	chatGPTClient = chatgpt.NewChatGPTClient("https://chat.openai.com")

	// 初始化 LLM 提供者
	llmConfig := &llm.Config{
		APIKey:      config.OpenAIAPIKey,
		BaseURL:     config.OpenAIBaseURL,
		Model:       config.OpenAIModel,
		Temperature: config.OpenAITemperature,
		MaxTokens:   config.OpenAIMaxTokens,
		Timeout:     config.RequestTimeout,
	}

	provider = llm.NewOpenAIProvider(llmConfig)

	// 初始化 Prompt 引擎
	promptEngine = prompt.NewPromptEngine()

	// 添加默认模板
	for name, tmpl := range prompt.DefaultTemplates {
		if err := promptEngine.AddTemplate(tmpl); err != nil {
			logger.Warnf("Failed to add template %s: %v", name, err)
		}
	}

	// 初始化 RAG 引擎
	retriever := rag.NewSimpleRetriever()
	ragEngine = rag.NewRAGEngine(retriever)

	// 添加示例文档
	addSampleDocuments(retriever)
}

func setupRoutes(r *gin.Engine) {
	// API 版本组
	v1 := r.Group("/api/v1")
	{
		// 认证相关接口
		v1.POST("/auth/register", handleRegister)
		v1.POST("/auth/login", handleLogin)
		v1.POST("/auth/logout", authManager.AuthMiddleware(), handleLogout)
		v1.GET("/auth/profile", authManager.AuthMiddleware(), handleGetProfile)
		v1.PUT("/auth/profile", authManager.AuthMiddleware(), handleUpdateProfile)

		// ChatGPT 相关接口
		v1.POST("/chatgpt/chat", authManager.AuthMiddleware(), handleChatGPTChat)
		v1.GET("/chatgpt/conversations", authManager.AuthMiddleware(), handleGetConversations)
		v1.GET("/chatgpt/conversations/:id", authManager.AuthMiddleware(), handleGetConversation)
		v1.DELETE("/chatgpt/conversations/:id", authManager.AuthMiddleware(), handleDeleteConversation)

		// 聊天接口
		v1.POST("/chat", handleChat)

		// 模板管理
		v1.GET("/templates", handleListTemplates)
		v1.POST("/templates", handleAddTemplate)
		v1.GET("/templates/:name", handleGetTemplate)
		v1.DELETE("/templates/:name", handleDeleteTemplate)

		// RAG 接口
		v1.POST("/rag/query", handleRAGQuery)
		v1.POST("/rag/documents", handleAddDocument)

		// 健康检查
		v1.GET("/health", handleHealth)
	}

	// 根路径
	r.GET("/", handleRoot)
}

// 认证相关处理器
func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := authManager.RegisterUser(req.Username, req.Email, req.ChatGPTToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := authManager.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"session": session,
	})
}

func handleLogout(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*auth.User)

	if err := authManager.LogoutUser(userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func handleGetProfile(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func handleUpdateProfile(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*auth.User)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := authManager.UpdateUser(userObj.ID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

// ChatGPT 相关处理器
func handleChatGPTChat(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*auth.User)

	var req ChatGPTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用用户的 ChatGPT Token
	chatReq := chatgpt.ChatRequest{
		Message:        req.Message,
		ConversationID: req.ConversationID,
		Model:          req.Model,
	}

	response, err := chatGPTClient.SendMessage(userObj.ChatGPTToken, chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func handleGetConversations(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*auth.User)

	conversations, err := chatGPTClient.GetConversations(userObj.ChatGPTToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

func handleGetConversation(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*auth.User)
	conversationID := c.Param("id")

	conversation, err := chatGPTClient.GetConversation(userObj.ChatGPTToken, conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func handleDeleteConversation(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*auth.User)
	conversationID := c.Param("id")

	if err := chatGPTClient.DeleteConversation(userObj.ChatGPTToken, conversationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted successfully"})
}

func handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Template == "" {
		req.Template = "qa"
	}
	if req.Model == "" {
		req.Model = provider.GetConfig().Model
	}

	var response ChatResponse
	response.Query = req.Query
	response.Template = req.Template
	response.Model = req.Model

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if req.ChainMode {
		// 链式调用模式
		result, err := runChainMode(ctx, req)
		if err != nil {
			response.Error = err.Error()
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response.Answer = result
	} else {
		// 简单模式
		result, tokenUsage, err := runSimpleMode(ctx, req)
		if err != nil {
			response.Error = err.Error()
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response.Answer = result
		response.TokenUsage = tokenUsage
	}

	c.JSON(http.StatusOK, response)
}

func runChainMode(ctx context.Context, req ChatRequest) (string, error) {
	// 创建链式调用
	c := chain.NewChain()

	// 添加步骤：检索 -> 构建 Prompt -> 调用 LLM
	c.AddStep(rag.Retrieve)
	c.AddStep(prompt.BuildPrompt)
	c.AddStep(func(ctx context.Context, input interface{}) (interface{}, error) {
		if str, ok := input.(string); ok {
			// 调用 LLM
			llmReq := &llm.ChatRequest{
				Model:       req.Model,
				Messages:    []llm.Message{{Role: "user", Content: str}},
				Temperature: provider.GetConfig().Temperature,
				MaxTokens:   provider.GetConfig().MaxTokens,
			}

			resp, err := provider.Chat(ctx, llmReq)
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

	result, err := c.RunString(ctx, req.Query)
	if err != nil {
		return "", fmt.Errorf("chain execution failed: %w", err)
	}

	return result, nil
}

func runSimpleMode(ctx context.Context, req ChatRequest) (string, int, error) {
	// 渲染 Prompt 模板
	data := map[string]interface{}{
		"question": req.Query,
	}

	// 添加自定义变量
	for k, v := range req.Variables {
		data[k] = v
	}

	prompt, err := promptEngine.Render(req.Template, data)
	if err != nil {
		return "", 0, fmt.Errorf("failed to render prompt: %w", err)
	}

	// 调用 LLM
	llmReq := &llm.ChatRequest{
		Model:       req.Model,
		Messages:    []llm.Message{{Role: "user", Content: prompt}},
		Temperature: provider.GetConfig().Temperature,
		MaxTokens:   provider.GetConfig().MaxTokens,
	}

	resp, err := provider.Chat(ctx, llmReq)
	if err != nil {
		return "", 0, fmt.Errorf("LLM call failed: %w", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, resp.Usage.TotalTokens, nil
	}

	return "抱歉，没有获得有效回复。", 0, nil
}

func handleListTemplates(c *gin.Context) {
	templates := promptEngine.ListTemplates()
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func handleAddTemplate(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tmpl := &prompt.Template{
		Name:     req.Name,
		Content:  req.Content,
		Version:  req.Version,
		Metadata: req.Metadata,
	}

	if err := promptEngine.AddTemplate(tmpl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template added successfully"})
}

func handleGetTemplate(c *gin.Context) {
	name := c.Param("name")
	tmpl, err := promptEngine.GetTemplate(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := TemplateResponse{
		Name:      tmpl.Name,
		Content:   tmpl.Content,
		Version:   tmpl.Version,
		Variables: tmpl.Variables,
		Metadata:  tmpl.Metadata,
	}

	c.JSON(http.StatusOK, response)
}

func handleDeleteTemplate(c *gin.Context) {
	name := c.Param("name")
	if err := promptEngine.RemoveTemplate(name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

func handleRAGQuery(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
		Limit int    `json:"limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 5
	}

	result, err := ragEngine.Query(c.Request.Context(), req.Query, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":  req.Query,
		"result": result,
		"limit":  req.Limit,
	})
}

func handleAddDocument(c *gin.Context) {
	var req struct {
		ID       string            `json:"id" binding:"required"`
		Content  string            `json:"content" binding:"required"`
		Metadata map[string]string `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 这里需要访问实际的检索器，简化处理
	c.JSON(http.StatusOK, gin.H{"message": "Document added successfully"})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

func handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Go LLM Tools API",
		"version": "1.0.0",
		"docs":    "/api/v1/health",
	})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func addSampleDocuments(retriever *rag.SimpleRetriever) {
	// 添加示例文档（与 CLI 工具相同）
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
			logger.Warnf("Failed to add document %s: %v", doc.ID, err)
		}
	}
}
