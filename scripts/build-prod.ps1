# FG-ABYSS Production Build Script
# Note: For production builds, prefer using Taskfile.yml: `task build:prod`
# This script builds optimized production version to bin/prod/ directory.

# Get the script's directory and project root
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

# Change to project root directory
Set-Location $ProjectRoot

Write-Host "=== Building FG-ABYSS (Production) ===" -ForegroundColor Cyan
Write-Host "Project Root: $ProjectRoot" -ForegroundColor Gray
Write-Host "Note: This creates an optimized release build" -ForegroundColor Yellow
Write-Host ""

# 1. Clean old build files
Write-Host "[1/5] Cleaning old production build files..." -ForegroundColor Yellow
if (Test-Path "bin\prod\FG-ABYSS.exe") {
    Remove-Item -Path "bin\prod\FG-ABYSS.exe" -Force
    Write-Host "  - Removed old production executable" -ForegroundColor Green
}

# Create bin/prod directory if not exists
if (-not (Test-Path "bin\prod")) {
    New-Item -ItemType Directory -Path "bin\prod" -Force | Out-Null
    Write-Host "  - Created bin/prod directory" -ForegroundColor Green
}

# Clean data directory for fresh start
if (Test-Path "bin\prod\data\app.db*") {
    Remove-Item -Path "bin\prod\data\app.db*" -Force
    Write-Host "  - Removed old production database files" -ForegroundColor Green
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

# 4. Build Windows GUI application (production mode - optimized)
Write-Host "`n[4/5] Building Windows GUI application (Production)..." -ForegroundColor Yellow
Write-Host "  Output: bin/prod/Fg-ABYSS.exe" -ForegroundColor Gray
Write-Host "  Note: Production build with optimizations (-s -w flags)" -ForegroundColor Gray
Write-Host "  Console: Hidden (GUI mode)" -ForegroundColor Gray

# Build production version (optimized, no console window)
# -s: Omit symbol table
# -w: Omit DWARF debugging information
# -H=windowsgui: Hide console window
go build -ldflags="-s -w -H=windowsgui" -o bin/prod/FG-ABYSS.exe .

if ($LASTEXITCODE -ne 0) {
    Write-Host "  - Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "  - Production build successful" -ForegroundColor Green

# 5. Setup data directory
Write-Host "`n[5/5] Setting up data directory..." -ForegroundColor Yellow
if (-not (Test-Path "bin\prod\data")) {
    New-Item -ItemType Directory -Path "bin\prod\data" -Force | Out-Null
    Write-Host "  - Created bin/prod/data directory" -ForegroundColor Green
}

# Show build results
Write-Host "`n=== Production Build Complete ===" -ForegroundColor Cyan
$exePath = Resolve-Path "bin\prod\FG-ABYSS.exe"
Write-Host "Executable: $exePath" -ForegroundColor Green

$exeSize = (Get-Item $exePath).Length / 1MB
Write-Host ("File size: " + [math]::Round($exeSize, 2) + " MB") -ForegroundColor Green
Write-Host "Note: Production builds are smaller and faster" -ForegroundColor Yellow

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Production build completed successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "`nNext steps:" -ForegroundColor White
Write-Host "  Run: .\bin\prod\FG-ABYSS.exe" -ForegroundColor Cyan
Write-Host "  Test: Verify all features work correctly" -ForegroundColor Cyan
Write-Host "  Deploy: Distribute to users" -ForegroundColor Cyan
Write-Host "`nRelated commands:" -ForegroundColor Yellow
Write-Host "  task package   - Create Windows installer (.msi)" -ForegroundColor Gray
Write-Host "  task dev       - Run development version" -ForegroundColor Gray
