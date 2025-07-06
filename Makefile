# Go LLM Tools Makefile

.PHONY: help build run test clean deps lint format docs

# 默认目标
help:
	@echo "Go LLM Tools - 大模型应用开发工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make build     - 构建项目"
	@echo "  make run       - 运行 CLI 工具"
	@echo "  make api       - 启动 API 服务"
	@echo "  make test      - 运行测试"
	@echo "  make clean     - 清理构建文件"
	@echo "  make deps      - 安装依赖"
	@echo "  make lint      - 代码检查"
	@echo "  make format    - 格式化代码"
	@echo "  make docs      - 生成文档"
	@echo "  make example   - 运行示例"

# 构建项目
build:
	@echo "构建项目..."
	go build -o bin/cli cmd/cli/main.go
	go build -o bin/api api/main.go
	@echo "构建完成!"

# 运行 CLI 工具
run:
	@echo "运行 CLI 工具..."
	go run cmd/cli/main.go -query "什么是 LangChain？" -verbose

# 启动 API 服务
api:
	@echo "启动 API 服务..."
	go run api/main.go

# 运行测试
test:
	@echo "运行测试..."
	go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	go clean

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod tidy
	go mod download

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 格式化代码
format:
	@echo "格式化代码..."
	go fmt ./...
	go vet ./...

# 生成文档
docs:
	@echo "生成文档..."
	godoc -http=:6060

# 运行示例
example:
	@echo "运行示例..."
	go run examples/basic_usage.go

# 安装工具
install-tools:
	@echo "安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest

# 开发模式
dev:
	@echo "开发模式..."
	@echo "1. 启动 API 服务..."
	@echo "2. 在另一个终端运行: make run"
	@echo "3. 访问 http://localhost:8080/api/v1/health"
	go run api/main.go

# 完整构建
all: deps build test
	@echo "完整构建完成!" 