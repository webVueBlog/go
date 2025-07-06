# Go LLM Tools - 大模型应用开发工具

一个用 Go 语言开发的大模型应用开发工具，提供类似 LangChain、RAG、Prompt Engineering 等功能，提升大模型应用开发效率。

## 功能特性

- 🔗 **链式调用 (Chain)**: 支持多步骤串联处理
- 🔍 **RAG (检索增强生成)**: 支持向量检索和文档增强
- 📝 **Prompt 工程**: 支持模板、变量替换、版本管理
- 🤖 **多模型支持**: 统一接口支持 OpenAI、Azure、百度文心等
- 🚀 **Web API**: 提供 RESTful API 接口
- 🛠️ **CLI 工具**: 命令行工具支持

## 项目结构

```
├── cmd/                    # 命令行工具
│   └── cli/
├── internal/               # 内部包
│   ├── chain/             # 链式调用
│   ├── rag/               # RAG 功能
│   ├── prompt/            # Prompt 工程
│   ├── llm/               # 大模型接口
│   └── utils/             # 工具函数
├── pkg/                   # 可导出的包
├── api/                   # API 服务
├── examples/              # 示例代码
└── docs/                  # 文档
```

## 快速开始

### 安装依赖
```bash
go mod tidy
```

### 运行示例
```bash
go run cmd/cli/main.go
```

### 启动 API 服务
```bash
go run api/main.go
```

## 使用示例

```go
package main

import (
    "fmt"
    "go-llm-tools/internal/chain"
    "go-llm-tools/internal/rag"
    "go-llm-tools/internal/prompt"
    "go-llm-tools/internal/llm"
)

func main() {
    // 创建链式调用
    c := chain.NewChain()
    c.AddStep(rag.Retrieve)
    c.AddStep(prompt.BuildPrompt)
    c.AddStep(llm.CallLLM)
    
    result := c.Run("请介绍一下LangChain")
    fmt.Println(result)
}
```

## 配置

创建 `.env` 文件：
```
OPENAI_API_KEY=your_openai_api_key
OPENAI_BASE_URL=https://api.openai.com/v1
```

## 许可证

MIT License 