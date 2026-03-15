# FG-ABYSS Development Build Script
# Note: For production builds, prefer using Taskfile.yml: `task build`
# This script builds development version to bin/dev/ directory.

# Get the script's directory and project root
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

# Change to project root directory
Set-Location $ProjectRoot

Write-Host "=== Building FG-ABYSS (Development) ===" -ForegroundColor Cyan
Write-Host "Project Root: $ProjectRoot" -ForegroundColor Gray
Write-Host "Note: For production builds, use: task build:prod" -ForegroundColor Yellow
Write-Host ""

# 1. Clean old build files
Write-Host "[1/5] Cleaning old development build files..." -ForegroundColor Yellow
if (Test-Path "bin\dev\FG-ABYSS-dev.exe") {
    Remove-Item -Path "bin\dev\FG-ABYSS-dev.exe" -Force
    Write-Host "  - Removed old development executable" -ForegroundColor Green
}

# Create bin/dev directory if not exists
if (-not (Test-Path "bin\dev")) {
    New-Item -ItemType Directory -Path "bin\dev" -Force | Out-Null
    Write-Host "  - Created bin/dev directory" -ForegroundColor Green
}

# Clean data directory for fresh start
if (Test-Path "data\app.db*") {
    Remove-Item -Path "data\app.db*" -Force
    Write-Host "  - Removed old database files" -ForegroundColor Green
}

# 2. Build frontend
Write-Host "`n[2/5] Building frontend assets..." -ForegroundColor Yellow
Set-Location frontend
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "  - Frontend build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  - Frontend build successful" -ForegroundColor Green
Set-Location ..

# 3. Verify frontend files
Write-Host "`n[3/5] Verifying frontend files..." -ForegroundColor Yellow
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

# 4. Build Windows GUI application (development mode)
Write-Host "`n[4/5] Building Windows GUI application (Development)..." -ForegroundColor Yellow
Write-Host "  Output: bin/dev/Fg-ABYSS-dev.exe" -ForegroundColor Gray
Write-Host "  Note: Development build with debug symbols" -ForegroundColor Gray

# Build development version (with console window for debugging)
go build -o bin/dev/FG-ABYSS-dev.exe .

if ($LASTEXITCODE -ne 0) {
    Write-Host "  - Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "  - Development build successful" -ForegroundColor Green

# 5. Copy database schema (if exists)
Write-Host "`n[5/5] Setting up data directory..." -ForegroundColor Yellow
if (-not (Test-Path "bin\dev\data")) {
    New-Item -ItemType Directory -Path "bin\dev\data" -Force | Out-Null
    Write-Host "  - Created bin/dev/data directory" -ForegroundColor Green
}

# Show build results
Write-Host "`n=== Development Build Complete ===" -ForegroundColor Cyan
$exePath = Resolve-Path "bin\dev\FG-ABYSS-dev.exe"
Write-Host "Executable: $exePath" -ForegroundColor Green

$exeSize = (Get-Item $exePath).Length / 1MB
Write-Host ("File size: " + [math]::Round($exeSize, 2) + " MB") -ForegroundColor Green

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Development build completed successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "`nNext steps:" -ForegroundColor White
Write-Host "  Run: .\bin\dev\FG-ABYSS-dev.exe" -ForegroundColor Cyan
Write-Host "  Or use Taskfile: task dev" -ForegroundColor Cyan
Write-Host "`nFor production builds, use:" -ForegroundColor Yellow
Write-Host "  task build:prod  - Build production version" -ForegroundColor Gray
Write-Host "  task package     - Create installer" -ForegroundColor Gray
