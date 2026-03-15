# FG-ABYSS Project Standards Check Script

Write-Host "========================================="
Write-Host "FG-ABYSS Project Standards Check"
Write-Host "========================================="
Write-Host ""

$Errors = 0
$Warnings = 0

# Check 1: Root .md files
Write-Host "1. Checking documentation location..."
$RootMd = Get-ChildItem -Path . -Filter *.md | Where-Object { $_.Name -ne "README.md" -and $_.Name -ne "LICENSE" }
if ($RootMd) {
    Write-Host "  [ERROR] Found .md files in root directory:" -ForegroundColor Red
    $RootMd | ForEach-Object { Write-Host "    - $($_.Name)" }
    $Errors++
} else {
    Write-Host "  [OK] No违规 .md files in root" -ForegroundColor Green
}

# Check 2: docs directory
Write-Host "2. Checking docs directory structure..."
if (Test-Path "docs") {
    Write-Host "  [OK] docs directory exists" -ForegroundColor Green
} else {
    Write-Host "  [ERROR] docs directory not found" -ForegroundColor Red
    $Errors++
}

# Check 3: Go file naming
Write-Host "3. Checking Go file naming..."
$GoFiles = Get-ChildItem -Recurse -Filter *.go -Exclude *_test.go | Where-Object { $_.FullName -notmatch "vendor|bindings" }
$BadGo = $GoFiles | Where-Object { $_.Name -cmatch "^[A-Z]" }
if ($BadGo) {
    Write-Host "  [WARNING] Found $($BadGo.Count) Go files with uppercase naming" -ForegroundColor Yellow
    $Warnings++
} else {
    Write-Host "  [OK] Go file naming is correct" -ForegroundColor Green
}

# Check 4: Test files
Write-Host "4. Checking test files..."
$TestFiles = Get-ChildItem -Recurse -Filter *_test.go | Where-Object { $_.FullName -notmatch "vendor" }
if ($TestFiles) {
    Write-Host "  [OK] Found $($TestFiles.Count) test files" -ForegroundColor Green
} else {
    Write-Host "  [WARNING] No test files found" -ForegroundColor Yellow
    $Warnings++
}

# Check 5: Vue components
Write-Host "5. Checking Vue component naming..."
if (Test-Path "frontend/src") {
    $VueFiles = Get-ChildItem -Path "frontend/src" -Recurse -Filter *.vue
    $BadVue = $VueFiles | Where-Object { $_.Name -cmatch "^[a-z]" }
    if ($BadVue) {
        Write-Host "  [WARNING] Found $($BadVue.Count) Vue components with lowercase naming" -ForegroundColor Yellow
        $Warnings++
    } else {
        Write-Host "  [OK] Vue component naming is correct" -ForegroundColor Green
    }
}

# Check 6: Git commits
Write-Host "6. Checking Git commit messages..."
if (Test-Path ".git") {
    $Commits = git log --oneline -10 2>$null
    $BadCommits = $Commits | Where-Object { $_ -notmatch "^(feat|fix|docs|style|refactor|test|chore|perf|ci)\(" }
    if ($BadCommits) {
        Write-Host "  [WARNING] Found $($BadCommits.Count) non-standard commits" -ForegroundColor Yellow
        $Warnings++
    } else {
        Write-Host "  [OK] Commit messages follow convention" -ForegroundColor Green
    }
}

# Summary
Write-Host ""
Write-Host "========================================="
Write-Host "Summary"
Write-Host "========================================="
Write-Host "Errors: $Errors"
Write-Host "Warnings: $Warnings"
Write-Host ""

if ($Errors -eq 0 -and $Warnings -eq 0) {
    Write-Host "Perfect! All checks passed." -ForegroundColor Green
}

Write-Host ""
Write-Host "Documentation: docs/development/project-optimization-specification.md"
