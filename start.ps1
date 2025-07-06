# Go LLM Tools - PowerShell 启动脚本

Write-Host "========================================" -ForegroundColor Green
Write-Host "Go LLM Tools - 启动脚本" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# 检查 Go 是否安装
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Go not found"
    }
    Write-Host "Go 环境检查通过！" -ForegroundColor Green
    Write-Host $goVersion -ForegroundColor Yellow
} catch {
    Write-Host "错误: Go 未安装或未配置到 PATH 中" -ForegroundColor Red
    Write-Host "请先安装 Go: https://go.dev/dl/" -ForegroundColor Yellow
    Read-Host "按任意键退出"
    exit 1
}

# 安装依赖
Write-Host "正在安装依赖..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: 依赖安装失败" -ForegroundColor Red
    Read-Host "按任意键退出"
    exit 1
}
Write-Host "依赖安装完成！" -ForegroundColor Green

# 显示菜单
function Show-Menu {
    Write-Host ""
    Write-Host "请选择要启动的服务:" -ForegroundColor Cyan
    Write-Host "1. 启动 API 服务" -ForegroundColor White
    Write-Host "2. 运行 CLI 工具 (简单模式)" -ForegroundColor White
    Write-Host "3. 运行 CLI 工具 (链式调用模式)" -ForegroundColor White
    Write-Host "4. 运行示例程序" -ForegroundColor White
    Write-Host "5. 退出" -ForegroundColor White
    Write-Host ""
}

function Start-API {
    Write-Host "启动 API 服务..." -ForegroundColor Yellow
    Write-Host "API 将在 http://localhost:8080 启动" -ForegroundColor Green
    Write-Host "按 Ctrl+C 停止服务" -ForegroundColor Yellow
    go run api/main.go
}

function Start-CLI-Simple {
    Write-Host "运行 CLI 工具 (简单模式)..." -ForegroundColor Yellow
    go run cmd/cli/main.go -query "什么是 LangChain？" -verbose
    Write-Host ""
    Read-Host "按任意键继续"
}

function Start-CLI-Chain {
    Write-Host "运行 CLI 工具 (链式调用模式)..." -ForegroundColor Yellow
    go run cmd/cli/main.go -query "什么是 RAG？" -chain -verbose
    Write-Host ""
    Read-Host "按任意键继续"
}

function Start-Example {
    Write-Host "运行示例程序..." -ForegroundColor Yellow
    go run examples/basic_usage.go
    Write-Host ""
    Read-Host "按任意键继续"
}

# 主循环
do {
    Show-Menu
    $choice = Read-Host "请输入选择 (1-5)"
    
    switch ($choice) {
        "1" { Start-API }
        "2" { Start-CLI-Simple }
        "3" { Start-CLI-Chain }
        "4" { Start-Example }
        "5" { 
            Write-Host "再见！" -ForegroundColor Green
            break 
        }
        default { 
            Write-Host "无效选择，请重新输入" -ForegroundColor Red
        }
    }
} while ($true) 