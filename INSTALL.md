# Go LLM Tools - 安装指南

## 🚀 快速安装

### 第一步：安装 Go 环境

#### 方法一：官方安装包（推荐）
1. 访问 https://go.dev/dl/
2. 下载 Windows 版本的 MSI 安装包
3. 运行安装程序，按默认设置安装
4. 重启 PowerShell 或命令提示符

#### 方法二：使用包管理器
```powershell
# 使用 winget
winget install GoLang.Go

# 或使用 Chocolatey
choco install golang

# 或使用 Scoop
scoop install go
```

### 第二步：验证安装
```powershell
go version
```
应该显示类似：`go version go1.21.5 windows/amd64`

### 第三步：配置项目

1. **复制环境变量文件**
```powershell
copy env.example .env
```

2. **编辑 .env 文件**
   - 打开 `.env` 文件
   - 将 `your_openai_api_key_here` 替换为你的 OpenAI API Key
   - 可以从 https://platform.openai.com/api-keys 获取

### 第四步：安装依赖
```powershell
go mod tidy
```

### 第五步：启动项目

#### 使用启动脚本（推荐）
```powershell
powershell -ExecutionPolicy Bypass -File simple-start.ps1
```

#### 手动启动
```powershell
# 启动 API 服务
go run api/main.go

# 运行 CLI 工具
go run cmd/cli/main.go -query "什么是 LangChain？" -verbose

# 运行示例
go run examples/basic_usage.go
```

## 🔧 故障排除

### 问题 1：Go 命令未找到
**解决方案：**
1. 检查 Go 是否正确安装
2. 重启终端或重启系统
3. 检查 PATH 环境变量是否包含 Go 路径

### 问题 2：依赖下载失败
**解决方案：**
```powershell
# 设置代理（中国大陆用户）
go env -w GOPROXY=https://goproxy.cn,direct

# 重新安装依赖
go mod tidy
```

### 问题 3：API 调用失败
**解决方案：**
1. 检查 API Key 是否正确
2. 检查网络连接
3. 如果在中国大陆，可能需要配置代理

### 问题 4：端口被占用
**解决方案：**
1. 修改 `.env` 文件中的 `SERVER_PORT`
2. 或停止占用端口的其他服务

## 📝 配置说明

### 环境变量
- `OPENAI_API_KEY`: OpenAI API Key（必需）
- `OPENAI_BASE_URL`: OpenAI API 地址
- `OPENAI_MODEL`: 使用的模型
- `SERVER_PORT`: API 服务端口
- `LOG_LEVEL`: 日志级别

### 支持的模型
- gpt-3.5-turbo
- gpt-4
- gpt-4-turbo
- 其他 OpenAI 兼容的模型

## 🎯 快速测试

安装完成后，可以运行以下命令测试：

```powershell
# 测试 API 服务
curl http://localhost:8080/api/v1/health

# 测试聊天功能
curl -X POST http://localhost:8080/api/v1/chat -H "Content-Type: application/json" -d '{"query": "Hello"}'
```

## 📚 更多信息

- 项目文档：查看 `README.md`
- API 文档：查看 `docs/API.md`
- 使用示例：查看 `examples/` 目录 