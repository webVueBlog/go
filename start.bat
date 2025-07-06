@echo off
echo ========================================
echo Go LLM Tools - 启动脚本
echo ========================================

REM 检查 Go 是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: Go 未安装或未配置到 PATH 中
    echo 请先安装 Go: https://go.dev/dl/
    pause
    exit /b 1
)

echo Go 环境检查通过！

REM 安装依赖
echo 正在安装依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo 错误: 依赖安装失败
    pause
    exit /b 1
)

echo 依赖安装完成！

REM 显示菜单
:menu
echo.
echo 请选择要启动的服务:
echo 1. 启动 API 服务
echo 2. 运行 CLI 工具 (简单模式)
echo 3. 运行 CLI 工具 (链式调用模式)
echo 4. 运行示例程序
echo 5. 退出
echo.
set /p choice=请输入选择 (1-5): 

if "%choice%"=="1" goto start_api
if "%choice%"=="2" goto start_cli_simple
if "%choice%"=="3" goto start_cli_chain
if "%choice%"=="4" goto start_example
if "%choice%"=="5" goto exit
goto menu

:start_api
echo 启动 API 服务...
echo API 将在 http://localhost:8080 启动
echo 按 Ctrl+C 停止服务
go run api/main.go
goto menu

:start_cli_simple
echo 运行 CLI 工具 (简单模式)...
go run cmd/cli/main.go -query "什么是 LangChain？" -verbose
echo.
pause
goto menu

:start_cli_chain
echo 运行 CLI 工具 (链式调用模式)...
go run cmd/cli/main.go -query "什么是 RAG？" -chain -verbose
echo.
pause
goto menu

:start_example
echo 运行示例程序...
go run examples/basic_usage.go
echo.
pause
goto menu

:exit
echo 再见！
pause 