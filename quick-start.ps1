# Go LLM Tools - 快速启动脚本

Write-Host "========================================" -ForegroundColor Green
Write-Host "Go LLM Tools - 快速启动" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# 检查 Go 是否安装
Write-Host "检查 Go 环境..." -ForegroundColor Yellow
$goInstalled = $false

# 尝试多个可能的 Go 安装位置
$goPaths = @(
    "C:\Program Files\Go\bin\go.exe",
    "C:\Go\bin\go.exe",
    "$env:USERPROFILE\go\bin\go.exe",
    "$env:LOCALAPPDATA\Programs\Go\bin\go.exe"
)

foreach ($path in $goPaths) {
    if (Test-Path $path) {
        Write-Host "找到 Go: $path" -ForegroundColor Green
        $env:PATH += ";$(Split-Path $path)"
        $goInstalled = $true
        break
    }
}

# 如果没找到，尝试使用 go 命令
if (-not $goInstalled) {
    try {
        $goVersion = & go version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Go 环境检查通过！" -ForegroundColor Green
            Write-Host $goVersion -ForegroundColor Yellow
            $goInstalled = $true
        }
    } catch {
        # Go 未安装
    }
}

if (-not $goInstalled) {
    Write-Host "Go 未安装或未配置到 PATH 中" -ForegroundColor Red
    Write-Host ""
    Write-Host "请选择安装方式:" -ForegroundColor Cyan
    Write-Host "1. 手动下载安装 (推荐)" -ForegroundColor White
    Write-Host "2. 使用 winget 安装" -ForegroundColor White
    Write-Host "3. 跳过安装，仅查看项目" -ForegroundColor White
    Write-Host ""
    
    $installChoice = Read-Host "请输入选择 (1-3)"
    
    switch ($installChoice) {
        "1" {
            Write-Host "请访问 https://go.dev/dl/ 下载并安装 Go" -ForegroundColor Yellow
            Write-Host "安装完成后重启此脚本" -ForegroundColor Yellow
            Start-Process "https://go.dev/dl/"
            Read-Host "按任意键退出"
            exit
        }
        "2" {
            Write-Host "正在使用 winget 安装 Go..." -ForegroundColor Yellow
            winget install GoLang.Go
            Write-Host "安装完成！请重启 PowerShell 后再次运行此脚本" -ForegroundColor Green
            Read-Host "按任意键退出"
            exit
        }
        "3" {
            Write-Host "跳过安装，显示项目信息..." -ForegroundColor Yellow
        }
        default {
            Write-Host "无效选择" -ForegroundColor Red
            exit
        }
    }
}

# 如果 Go 已安装，继续启动项目
if ($goInstalled) {
    Write-Host ""
    Write-Host "Go 环境检查通过！开始启动项目..." -ForegroundColor Green
    
    # 检查 .env 文件
    if (-not (Test-Path ".env")) {
        Write-Host "创建 .env 文件..." -ForegroundColor Yellow
        Copy-Item "env.example" ".env"
        Write-Host "请编辑 .env 文件，添加你的 OpenAI API Key" -ForegroundColor Yellow
    }
    
    # 安装依赖
    Write-Host "安装项目依赖..." -ForegroundColor Yellow
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Host "依赖安装失败，请检查网络连接" -ForegroundColor Red
        Read-Host "按任意键退出"
        exit
    }
    Write-Host "依赖安装完成！" -ForegroundColor Green
    
    # 显示菜单
    do {
        Write-Host ""
        Write-Host "请选择要启动的服务:" -ForegroundColor Cyan
        Write-Host "1. 启动 API 服务" -ForegroundColor White
        Write-Host "2. 运行 CLI 工具 (简单模式)" -ForegroundColor White
        Write-Host "3. 运行 CLI 工具 (链式调用模式)" -ForegroundColor White
        Write-Host "4. 运行示例程序" -ForegroundColor White
        Write-Host "5. 查看项目结构" -ForegroundColor White
        Write-Host "6. 退出" -ForegroundColor White
        Write-Host ""
        
        $choice = Read-Host "请输入选择 (1-6)"
        
        switch ($choice) {
            "1" {
                Write-Host "启动 API 服务..." -ForegroundColor Yellow
                Write-Host "API 将在 http://localhost:8080 启动" -ForegroundColor Green
                Write-Host "按 Ctrl+C 停止服务" -ForegroundColor Yellow
                go run api/main.go
            }
            "2" {
                Write-Host "运行 CLI 工具 (简单模式)..." -ForegroundColor Yellow
                go run cmd/cli/main.go -query "什么是 LangChain？" -verbose
                Write-Host ""
                Read-Host "按任意键继续"
            }
            "3" {
                Write-Host "运行 CLI 工具 (链式调用模式)..." -ForegroundColor Yellow
                go run cmd/cli/main.go -query "什么是 RAG？" -chain -verbose
                Write-Host ""
                Read-Host "按任意键继续"
            }
            "4" {
                Write-Host "运行示例程序..." -ForegroundColor Yellow
                go run examples/basic_usage.go
                Write-Host ""
                Read-Host "按任意键继续"
            }
            "5" {
                Write-Host "项目结构:" -ForegroundColor Yellow
                Get-ChildItem -Recurse -Directory | Select-Object FullName | ForEach-Object {
                    Write-Host $_.FullName -ForegroundColor Gray
                }
                Write-Host ""
                Read-Host "按任意键继续"
            }
            "6" {
                Write-Host "再见！" -ForegroundColor Green
                break
            }
            default {
                Write-Host "无效选择，请重新输入" -ForegroundColor Red
            }
        }
    } while ($true)
} else {
    Write-Host "项目文件已准备就绪，但需要安装 Go 环境才能运行" -ForegroundColor Yellow
    Write-Host "请参考 INSTALL.md 文件了解详细安装步骤" -ForegroundColor Yellow
} 