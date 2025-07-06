# Go LLM Tools - Simple Start Script

Write-Host "========================================" -ForegroundColor Green
Write-Host "Go LLM Tools - Quick Start" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# Check if Go is installed
Write-Host "Checking Go environment..." -ForegroundColor Yellow

try {
    $goVersion = & go version 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Go environment check passed!" -ForegroundColor Green
        Write-Host $goVersion -ForegroundColor Yellow
        
        # Check .env file
        if (-not (Test-Path ".env")) {
            Write-Host "Creating .env file..." -ForegroundColor Yellow
            Copy-Item "env.example" ".env"
            Write-Host "Please edit .env file and add your OpenAI API Key" -ForegroundColor Yellow
        }
        
        # Install dependencies
        Write-Host "Installing dependencies..." -ForegroundColor Yellow
        go mod tidy
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Failed to install dependencies" -ForegroundColor Red
            Read-Host "Press any key to exit"
            exit
        }
        Write-Host "Dependencies installed successfully!" -ForegroundColor Green
        
        # Show menu
        do {
            Write-Host ""
            Write-Host "Choose service to start:" -ForegroundColor Cyan
            Write-Host "1. Start API Server" -ForegroundColor White
            Write-Host "2. Run CLI Tool (Simple Mode)" -ForegroundColor White
            Write-Host "3. Run CLI Tool (Chain Mode)" -ForegroundColor White
            Write-Host "4. Run Examples" -ForegroundColor White
            Write-Host "5. Exit" -ForegroundColor White
            Write-Host ""
            
            $choice = Read-Host "Enter choice (1-5)"
            
            switch ($choice) {
                "1" {
                    Write-Host "Starting API server..." -ForegroundColor Yellow
                    Write-Host "API will be available at http://localhost:8080" -ForegroundColor Green
                    Write-Host "Press Ctrl+C to stop" -ForegroundColor Yellow
                    go run api/main.go
                }
                "2" {
                    Write-Host "Running CLI tool (simple mode)..." -ForegroundColor Yellow
                    go run cmd/cli/main.go -query "What is LangChain?" -verbose
                    Write-Host ""
                    Read-Host "Press any key to continue"
                }
                "3" {
                    Write-Host "Running CLI tool (chain mode)..." -ForegroundColor Yellow
                    go run cmd/cli/main.go -query "What is RAG?" -chain -verbose
                    Write-Host ""
                    Read-Host "Press any key to continue"
                }
                "4" {
                    Write-Host "Running examples..." -ForegroundColor Yellow
                    go run examples/basic_usage.go
                    Write-Host ""
                    Read-Host "Press any key to continue"
                }
                "5" {
                    Write-Host "Goodbye!" -ForegroundColor Green
                    break
                }
                default {
                    Write-Host "Invalid choice, please try again" -ForegroundColor Red
                }
            }
        } while ($true)
        
    } else {
        throw "Go not found"
    }
} catch {
    Write-Host "Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Go first:" -ForegroundColor Yellow
    Write-Host "1. Visit https://go.dev/dl/" -ForegroundColor White
    Write-Host "2. Download and install Go for Windows" -ForegroundColor White
    Write-Host "3. Restart PowerShell after installation" -ForegroundColor White
    Write-Host "4. Run this script again" -ForegroundColor White
    Write-Host ""
    Read-Host "Press any key to exit"
} 