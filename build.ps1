# Wails 3 Application Build Script
# Build Windows GUI Application (without console window)

Write-Host "=== Building FG-ABYSS ===" -ForegroundColor Cyan

# 1. Clean old build files
Write-Host "`n[1/4] Cleaning old build files..." -ForegroundColor Yellow
if (Test-Path "bin\FG-ABYSS.exe") {
    Remove-Item -Path "bin\FG-ABYSS.exe" -Force
    Write-Host "  - Removed old executable" -ForegroundColor Green
}

# Clean data directory for fresh start
if (Test-Path "data\app.db*") {
    Remove-Item -Path "data\app.db*" -Force
    Write-Host "  - Removed old database files" -ForegroundColor Green
}

# 2. Build frontend
Write-Host "`n[2/4] Building frontend assets..." -ForegroundColor Yellow
Set-Location frontend
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "  - Frontend build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  - Frontend build successful" -ForegroundColor Green
Set-Location ..

# 3. Verify frontend files
Write-Host "`n[3/4] Verifying frontend files..." -ForegroundColor Yellow
if (-not (Test-Path "frontend\dist\index.html")) {
    Write-Host "  - Frontend files not found!" -ForegroundColor Red
    exit 1
}
Write-Host "  - Frontend files verified" -ForegroundColor Green

# List frontend dist files
Write-Host "`nFrontend dist files:" -ForegroundColor Gray
Get-ChildItem "frontend\dist" | ForEach-Object {
    Write-Host "  - $($_.Name)" -ForegroundColor Gray
}

# 4. Build Windows GUI application (hide console window)
Write-Host "`n[4/4] Building Windows GUI application..." -ForegroundColor Yellow
Write-Host "  Using -ldflags=-H=windowsgui to hide console window" -ForegroundColor Gray

# Build Windows GUI mode (no console window)
go build -ldflags="-H=windowsgui" -o bin/Fg-ABYSS.exe .

if ($LASTEXITCODE -ne 0) {
    Write-Host "  - Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "  - Build successful" -ForegroundColor Green

# Show build results
Write-Host "`n=== Build Complete ===" -ForegroundColor Cyan
$exePath = Resolve-Path "bin\FG-ABYSS.exe"
Write-Host "Executable: $exePath" -ForegroundColor Green

$exeSize = (Get-Item $exePath).Length / 1MB
Write-Host ("File size: " + [math]::Round($exeSize, 2) + " MB") -ForegroundColor Green

# Verify embedded files
Write-Host "`nVerifying embedded resources..." -ForegroundColor Yellow
go run -ldflags="-H=windowsgui" . --help 2>&1 | Select-String -Pattern "error|panic|failed" -Context 0,2

Write-Host "`nTo run the application:" -ForegroundColor Cyan
Write-Host "  .\bin\FG-ABYSS.exe" -ForegroundColor White

Write-Host "`nTo run with debug output (development mode):" -ForegroundColor Cyan
Write-Host "  go run ." -ForegroundColor White
