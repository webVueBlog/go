# Go LLM Tools - Project Showcase

Write-Host "========================================" -ForegroundColor Green
Write-Host "Go LLM Tools - Project Overview" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

Write-Host ""
Write-Host "Project Structure:" -ForegroundColor Cyan
Write-Host "├── api/                    # API Server" -ForegroundColor White
Write-Host "├── cmd/cli/               # CLI Tool" -ForegroundColor White
Write-Host "├── internal/              # Core Modules" -ForegroundColor White
Write-Host "│   ├── chain/            # Chain Processing" -ForegroundColor Gray
Write-Host "│   ├── llm/              # LLM Interface" -ForegroundColor Gray
Write-Host "│   ├── prompt/           # Prompt Engineering" -ForegroundColor Gray
Write-Host "│   ├── rag/              # RAG Functionality" -ForegroundColor Gray
Write-Host "│   └── utils/            # Utilities" -ForegroundColor Gray
Write-Host "├── examples/              # Usage Examples" -ForegroundColor White
Write-Host "├── docs/                  # Documentation" -ForegroundColor White
Write-Host "└── scripts/               # Startup Scripts" -ForegroundColor White

Write-Host ""
Write-Host "Features:" -ForegroundColor Cyan
Write-Host "✅ Chain Processing (类似 LangChain)" -ForegroundColor Green
Write-Host "✅ Multi-Model Support (OpenAI, Azure, etc.)" -ForegroundColor Green
Write-Host "✅ Prompt Engineering & Templates" -ForegroundColor Green
Write-Host "✅ RAG (Retrieval-Augmented Generation)" -ForegroundColor Green
Write-Host "✅ RESTful API Server" -ForegroundColor Green
Write-Host "✅ CLI Tool" -ForegroundColor Green
Write-Host "✅ Configuration Management" -ForegroundColor Green

Write-Host ""
Write-Host "To start the project, you need to:" -ForegroundColor Yellow
Write-Host "1. Install Go: https://go.dev/dl/" -ForegroundColor White
Write-Host "2. Configure OpenAI API Key in .env file" -ForegroundColor White
Write-Host "3. Run: go mod tidy" -ForegroundColor White
Write-Host "4. Run: go run api/main.go (for API server)" -ForegroundColor White
Write-Host "5. Or run: go run cmd/cli/main.go (for CLI tool)" -ForegroundColor White

Write-Host ""
Write-Host "Available Scripts:" -ForegroundColor Cyan
Write-Host "• simple-start.ps1    - Quick start script" -ForegroundColor White
Write-Host "• start.ps1           - Full featured start script" -ForegroundColor White
Write-Host "• start.bat           - Batch file version" -ForegroundColor White
Write-Host "• INSTALL.md          - Detailed installation guide" -ForegroundColor White

Write-Host ""
Write-Host "Example Usage:" -ForegroundColor Cyan
Write-Host "• API Server: http://localhost:8080/api/v1/chat" -ForegroundColor White
Write-Host "• CLI Tool: go run cmd/cli/main.go -query 'What is LangChain?'" -ForegroundColor White
Write-Host "• Examples: go run examples/basic_usage.go" -ForegroundColor White

Write-Host ""
Write-Host "Documentation:" -ForegroundColor Cyan
Write-Host "• README.md           - Project overview" -ForegroundColor White
Write-Host "• docs/API.md         - API documentation" -ForegroundColor White
Write-Host "• INSTALL.md          - Installation guide" -ForegroundColor White

Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Yellow
Write-Host "1. Install Go from https://go.dev/dl/" -ForegroundColor White
Write-Host "2. Restart PowerShell after installation" -ForegroundColor White
Write-Host "3. Run: powershell -ExecutionPolicy Bypass -File simple-start.ps1" -ForegroundColor White

Read-Host "Press any key to exit" 