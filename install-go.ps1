# Go Installation Script

Write-Host "========================================" -ForegroundColor Green
Write-Host "Go Installation Helper" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

Write-Host ""
Write-Host "Installing Go..." -ForegroundColor Yellow

# Try to download Go installer
try {
    $url = "https://go.dev/dl/go1.21.5.windows-amd64.msi"
    $output = "go-installer.msi"
    
    Write-Host "Downloading Go installer..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $url -OutFile $output -UseBasicParsing
    
    if (Test-Path $output) {
        Write-Host "Download completed!" -ForegroundColor Green
        Write-Host "Installing Go..." -ForegroundColor Yellow
        
        # Install Go
        Start-Process msiexec.exe -Wait -ArgumentList "/i $output /quiet"
        
        # Clean up
        Remove-Item $output
        
        Write-Host "Go installation completed!" -ForegroundColor Green
        Write-Host "Please restart PowerShell and try again." -ForegroundColor Yellow
    }
} catch {
    Write-Host "Failed to download Go installer automatically." -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Go manually:" -ForegroundColor Yellow
    Write-Host "1. Visit https://go.dev/dl/" -ForegroundColor White
    Write-Host "2. Download 'go1.21.5.windows-amd64.msi'" -ForegroundColor White
    Write-Host "3. Run the installer" -ForegroundColor White
    Write-Host "4. Restart PowerShell" -ForegroundColor White
    Write-Host ""
    
    # Open the download page
    Start-Process "https://go.dev/dl/"
}

Read-Host "Press any key to exit" 