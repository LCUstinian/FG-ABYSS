# Wails 3 项目构建脚本
# 使用 Wails 3 官方构建系统，确保正确的资源嵌入和绑定生成

Write-Host "=== Wails 3 Build ===" -ForegroundColor Cyan

# 1. Clean old build files
Write-Host "`n[1/4] Cleaning old build files..." -ForegroundColor Yellow
if (Test-Path "bin\FG-ABYSS.exe") {
    Remove-Item -Path "bin\FG-ABYSS.exe" -Force
    Write-Host "  - Removed old executable" -ForegroundColor Green
}

# 2. Generate bindings
Write-Host "`n[2/4] Generating TypeScript bindings..." -ForegroundColor Yellow
wails3 generate bindings
if ($LASTEXITCODE -ne 0) {
    Write-Host "  - Bindings generation failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  - Bindings generated successfully" -ForegroundColor Green

# 3. Build frontend
Write-Host "`n[3/4] Building frontend assets..." -ForegroundColor Yellow
Set-Location frontend
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "  - Frontend build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  - Frontend build successful" -ForegroundColor Green
Set-Location ..

# 4. Build application using Wails
Write-Host "`n[4/4] Building application with Wails..." -ForegroundColor Yellow
Write-Host "  This will:" -ForegroundColor Gray
Write-Host "    - Embed frontend assets" -ForegroundColor Gray
Write-Host "    - Apply platform-specific configurations" -ForegroundColor Gray
Write-Host "    - Generate Windows resources" -ForegroundColor Gray
Write-Host "    - Build optimized executable" -ForegroundColor Gray

wails3 build
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

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "To run the application:" -ForegroundColor White
Write-Host "  .\bin\FG-ABYSS.exe" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
