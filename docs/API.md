# Go LLM Tools API 文档

## 概述

Go LLM Tools 提供了完整的 RESTful API 接口，支持聊天、Prompt 模板管理、RAG 查询等功能。

## 基础信息

- **基础 URL**: `http://localhost:8080`
- **API 版本**: `v1`
- **内容类型**: `application/json`

## 认证

目前 API 不需要认证，但建议在生产环境中添加适当的认证机制。

## 端点

### 1. 健康检查

**GET** `/api/v1/health`

检查服务状态。

**响应示例:**
```json
{
  "status": "healthy",
  "timestamp": 1640995200,
  "version": "1.0.0"
}
```

### 2. 聊天接口

**POST** `/api/v1/chat`

与 LLM 进行对话。

**请求体:**
```json
{
  "query": "什么是 LangChain？",
  "template": "qa",
  "model": "gpt-3.5-turbo",
  "variables": {
    "custom_var": "value"
  },
  "chain_mode": false
}
```

**参数说明:**
- `query` (必需): 查询内容
- `template` (可选): 使用的 Prompt 模板，默认为 "qa"
- `model` (可选): 使用的模型，默认为配置中的模型
- `variables` (可选): 自定义变量
- `chain_mode` (可选): 是否使用链式调用模式

**响应示例:**
```json
{
  "query": "什么是 LangChain？",
  "answer": "LangChain 是一个用于开发由语言模型驱动的应用程序的框架...",
  "template": "qa",
  "model": "gpt-3.5-turbo",
  "token_usage": 150
}
```

### 3. 模板管理

#### 3.1 列出所有模板

**GET** `/api/v1/templates`

获取所有可用的 Prompt 模板。

**响应示例:**
```json
{
  "templates": ["qa", "translation", "summary", "code_review"]
}
```

#### 3.2 添加模板

**POST** `/api/v1/templates`

添加新的 Prompt 模板。

**请求体:**
```json
{
  "name": "custom_template",
  "content": "请回答以下问题：\n\n{{.question}}\n\n请提供详细的答案。",
  "version": "1.0",
  "metadata": {
    "type": "custom",
    "author": "user"
  }
}
```

**响应示例:**
```json
{
  "message": "Template added successfully"
}
```

#### 3.3 获取模板详情

**GET** `/api/v1/templates/{name}`

获取指定模板的详细信息。

**响应示例:**
```json
{
  "name": "qa",
  "content": "请回答以下问题：\n\n{{.question}}\n\n请提供详细、准确的答案。",
  "version": "1.0",
  "variables": ["question"],
  "metadata": {
    "type": "question-answer"
  }
}
```

#### 3.4 删除模板

**DELETE** `/api/v1/templates/{name}`

删除指定的模板。

**响应示例:**
```json
{
  "message": "Template deleted successfully"
}
```

### 4. RAG 功能

#### 4.1 RAG 查询

**POST** `/api/v1/rag/query`

执行 RAG 查询。

**请求体:**
```json
{
  "query": "Go 语言的并发特性",
  "limit": 5
}
```

**参数说明:**
- `query` (必需): 查询内容
- `limit` (可选): 返回结果数量限制，默认为 5

**响应示例:**
```json
{
  "query": "Go 语言的并发特性",
  "result": "基于以下上下文信息回答问题：\n\n上下文：\n文档 1 (相关度: 2.00):\nGo 语言支持并发编程，通过 goroutine 和 channel 实现轻量级线程和通信。\n\n问题：Go 语言的并发特性",
  "limit": 5
}
```

#### 4.2 添加文档

**POST** `/api/v1/rag/documents`

添加文档到 RAG 系统。

**请求体:**
```json
{
  "id": "doc_001",
  "content": "Go 语言是由 Google 开发的开源编程语言...",
  "metadata": {
    "source": "go_docs",
    "type": "language"
  }
}
```

**响应示例:**
```json
{
  "message": "Document added successfully"
}
```

## 错误处理

所有 API 端点都返回标准的 HTTP 状态码：

- `200 OK`: 请求成功
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误

错误响应格式：
```json
{
  "error": "错误描述信息"
}
```

## 使用示例

### cURL 示例

1. **健康检查:**
```bash
curl -X GET http://localhost:8080/api/v1/health
```

2. **聊天请求:**
```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "什么是 RAG？",
    "template": "qa",
    "chain_mode": false
  }'
```

3. **添加模板:**
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "custom_qa",
    "content": "请回答：{{.question}}",
    "version": "1.0"
  }'
```

### JavaScript 示例

```javascript
// 聊天请求
const response = await fetch('http://localhost:8080/api/v1/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    query: '什么是 LangChain？',
    template: 'qa',//qa是 默认模板
    chain_mode: false
  })
});

const result = await response.json();
console.log(result.answer);
```

## 配置

API 服务通过环境变量进行配置，详见 `env.example` 文件。

## 限制

- 请求超时时间：30 秒
- 最大 Token 数：1000（可配置）
- 并发请求数：无限制（建议根据服务器性能调整） 