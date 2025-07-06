package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 配置结构
type Config struct {
	// OpenAI 配置
	OpenAIAPIKey      string `json:"openai_api_key"`
	OpenAIBaseURL     string `json:"openai_base_url"`
	OpenAIModel       string `json:"openai_model"`
	OpenAITemperature float64 `json:"openai_temperature"`
	OpenAIMaxTokens   int    `json:"openai_max_tokens"`
	
	// 服务器配置
	ServerPort int    `json:"server_port"`
	ServerHost string `json:"server_host"`
	
	// 日志配置
	LogLevel string `json:"log_level"`
	LogFile  string `json:"log_file"`
	
	// RAG 配置
	RAGMaxResults int `json:"rag_max_results"`
	
	// 超时配置
	RequestTimeout time.Duration `json:"request_timeout"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 尝试加载 .env 文件
	if err := godotenv.Load(); err != nil {
		// 如果 .env 文件不存在，使用环境变量
		fmt.Println("Warning: .env file not found, using environment variables")
	}
	
	config := &Config{}
	
	// 加载 OpenAI 配置
	config.OpenAIAPIKey = getEnv("OPENAI_API_KEY", "")
	config.OpenAIBaseURL = getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1")
	config.OpenAIModel = getEnv("OPENAI_MODEL", "gpt-3.5-turbo")
	config.OpenAITemperature = getEnvFloat("OPENAI_TEMPERATURE", 0.7)
	config.OpenAIMaxTokens = getEnvInt("OPENAI_MAX_TOKENS", 1000)
	
	// 加载服务器配置
	config.ServerPort = getEnvInt("SERVER_PORT", 8080)
	config.ServerHost = getEnv("SERVER_HOST", "localhost")
	
	// 加载日志配置
	config.LogLevel = getEnv("LOG_LEVEL", "info")
	config.LogFile = getEnv("LOG_FILE", "")
	
	// 加载 RAG 配置
	config.RAGMaxResults = getEnvInt("RAG_MAX_RESULTS", 5)
	
	// 加载超时配置
	timeout := getEnvInt("REQUEST_TIMEOUT_SECONDS", 30)
	config.RequestTimeout = time.Duration(timeout) * time.Second
	
	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvFloat 获取浮点数环境变量
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// ValidateConfig 验证配置
func ValidateConfig(config *Config) error {
	if config.OpenAIAPIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required")
	}
	
	if config.ServerPort <= 0 || config.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", config.ServerPort)
	}
	
	if config.OpenAITemperature < 0 || config.OpenAITemperature > 2 {
		return fmt.Errorf("invalid temperature: %f", config.OpenAITemperature)
	}
	
	if config.OpenAIMaxTokens <= 0 {
		return fmt.Errorf("invalid max tokens: %d", config.OpenAIMaxTokens)
	}
	
	return nil
} 