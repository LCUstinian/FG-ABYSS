# FG-ABYSS 构建输出目录规范

> 本文档规定了 FG-ABYSS 项目的构建输出目录结构、命名格式和文件组织规范，确保构建产物的一致性和可维护性。

---

## 📋 目录

1. [输出目录结构](#1-输出目录结构)
2. [命名规范](#2-命名规范)
3. [文件组织](#3-文件组织)
4. [构建脚本规范](#4-构建脚本规范)
5. [示例](#5-示例)

---

## 1. 输出目录结构

### 1.1 标准目录结构

```
FG-ABYSS/
├── bin/                          # 构建输出目录（.gitignore）
│   ├── dev/                      # 开发版本
│   │   ├── FG-ABYSS-dev.exe      # 开发版可执行文件
│   │   └── data/                 # 开发数据目录
│   │       └── app.db            # SQLite 数据库
│   └── prod/                     # 生产版本
│       ├── FG-ABYSS.exe          # 生产版可执行文件
│       └── data/                 # 生产数据目录
│           └── app.db            # SQLite 数据库
├── frontend/dist/                # 前端构建产物
│   ├── index.html               # HTML 入口
│   ├── assets/                  # 静态资源
│   │   ├── index-xxxxx.css      # 样式文件
│   │   └── index-xxxxx.js       # JavaScript 文件
│   └── *.png, *.ttf             # 字体和图片
└── build/                        # 构建配置和临时文件
    ├── windows/                 # Windows 特定配置
    ├── darwin/                  # macOS 特定配置
    └── linux/                   # Linux 特定配置
```

### 1.2 目录说明

| 目录 | 用途 | 是否 Git 跟踪 | 说明 |
|------|------|--------------|------|
| `bin/` | 构建输出总目录 | ❌ 否 | 包含 dev/ 和 prod/ |
| `bin/dev/` | 开发版本输出 | ❌ 否 | 调试版本，带符号 |
| `bin/prod/` | 生产版本输出 | ❌ 否 | 发布版本，优化后 |
| `bin/*/data/` | 运行时数据 | ❌ 否 | 数据库和配置 |
| `frontend/dist/` | 前端构建产物 | ❌ 否 | Vite 构建结果 |
| `build/` | 构建配置 | ✅ 是 | 配置模板 |

---

## 2. 命名规范

### 2.1 可执行文件命名

**开发版本**:
```powershell
# ✅ 正确
FG-ABYSS-dev.exe          # Windows 开发版
FG-ABYSS-dev              # Linux/macOS 开发版

# ❌ 错误
Fg-ABYSS-dev.exe          # 混合大小写
fg-abyss-dev.exe          # 全小写
```

**生产版本**:
```powershell
# ✅ 正确
FG-ABYSS.exe              # Windows 生产版
FG-ABYSS                  # Linux/macOS 生产版

# ❌ 错误
Fg-ABYSS.exe              # 混合大小写
fg-abyss.exe              # 全小写
```

**原因**:
- 保持与项目名称一致
- Windows 不区分大小写，但统一命名更清晰
- 遵循项目文件命名规范
- 开发版使用 `-dev` 后缀区分

### 2.2 目录命名

**规则**: 小写字母 + 连字符

```powershell
# ✅ 正确
bin/
data/
assets/
resources/

# ❌ 错误
Bin/                  # 首字母大写
BIN/                  # 全大写
bin_data/             # 下划线分隔
```

### 2.3 前端资源文件命名

**规则**: 遵循 Vite 的哈希命名

```powershell
# ✅ 正确（Vite 自动生成）
index-CmfSNYlp.css
index-CN5LEtZV.js

# ❌ 错误（手动重命名）
index.css             # 无哈希，无法缓存控制
```

---

## 3. 文件组织

### 3.1 bin/ 目录组织

**结构**:
```
bin/
├── FG-ABYSS.exe              # 主程序
├── data/                     # 运行时数据
│   ├── app.db                # 主数据库
│   ├── app.db-shm            # SQLite 共享内存
│   └── app.db-wal            # SQLite 预写日志
└── resources/                # 额外资源（可选）
    └── config.yaml           # 配置文件
```

**说明**:
- 可执行文件放在 `bin/` 根目录
- 数据文件放在 `bin/data/` 子目录
- 资源文件放在 `bin/resources/` 子目录

### 3.2 frontend/dist/ 目录组织

**结构**:
```
frontend/dist/
├── index.html                # HTML 入口
├── assets/                   # 资源文件
│   ├── *.css                 # 样式文件
│   ├── *.js                  # JavaScript 文件
│   └── *.woff2               # 字体文件
├── *.png                     # 图片文件
└── *.ico                     # 图标文件
```

**说明**:
- 所有静态资源由 Vite 自动管理
- 使用哈希命名确保缓存控制
- 嵌入到 Go 二进制文件中

---

## 4. 构建脚本规范

### 4.1 输出目录配置

**PowerShell 脚本**:
```powershell
# ✅ 正确：使用相对路径和标准命名
go build -o bin/FG-ABYSS.exe .

# ❌ 错误：绝对路径
go build -o C:\Projects\FG-ABYSS\bin\FG-ABYSS.exe .

# ❌ 错误：非标准命名
go build -o bin/Fg-ABYSS.exe .
go build -o bin/fg_abyss.exe .
```

**Taskfile.yml**:
```yaml
tasks:
  build:
    cmds:
      - go build -o bin/FG-ABYSS.exe .
```

### 4.2 清理操作

**构建前清理**:
```powershell
# ✅ 正确：清理旧的构建文件
if (Test-Path "bin\FG-ABYSS.exe") {
    Remove-Item -Path "bin\FG-ABYSS.exe" -Force
}

# ✅ 正确：清理数据文件（可选）
if (Test-Path "data\app.db*") {
    Remove-Item -Path "data\app.db*" -Force
}
```

### 4.3 验证操作

**构建后验证**:
```powershell
# ✅ 正确：验证构建产物
if (-not (Test-Path "bin\FG-ABYSS.exe")) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# ✅ 正确：显示构建信息
$exeSize = (Get-Item "bin\FG-ABYSS.exe").Length / 1MB
Write-Host "File size: $([math]::Round($exeSize, 2)) MB"
```

---

## 5. 示例

### 5.1 完整构建脚本示例

```powershell
# FG-ABYSS Build Script
# 符合项目构建输出目录规范

Write-Host "=== Building FG-ABYSS ===" -ForegroundColor Cyan

# 1. Clean old build files
Write-Host "[1/4] Cleaning old build files..." -ForegroundColor Yellow
if (Test-Path "bin\FG-ABYSS.exe") {
    Remove-Item -Path "bin\FG-ABYSS.exe" -Force
}

# 2. Build frontend
Write-Host "[2/4] Building frontend..." -ForegroundColor Yellow
Set-Location frontend
npm run build
Set-Location ..

# 3. Verify frontend files
Write-Host "[3/4] Verifying frontend files..." -ForegroundColor Yellow
if (-not (Test-Path "frontend\dist\index.html")) {
    Write-Host "Frontend build failed!" -ForegroundColor Red
    exit 1
}

# 4. Build Windows application
Write-Host "[4/4] Building Windows application..." -ForegroundColor Yellow
go build -ldflags="-H=windowsgui" -o bin/FG-ABYSS.exe .

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Show build results
Write-Host "`n=== Build Complete ===" -ForegroundColor Cyan
$exePath = Resolve-Path "bin\FG-ABYSS.exe"
Write-Host "Executable: $exePath" -ForegroundColor Green

$exeSize = (Get-Item $exePath).Length / 1MB
Write-Host "File size: $([math]::Round($exeSize, 2)) MB" -ForegroundColor Green

Write-Host "`nRun: .\bin\FG-ABYSS.exe" -ForegroundColor Cyan
```

### 5.2 Taskfile.yml 示例

```yaml
version: '3'

tasks:
  build:
    summary: Build the application with standard output directory
    cmds:
      - task: common:build:frontend
      - task: common:build:backend
      - task: common:verify:build
    
  common:build:frontend:
    dir: frontend
    cmds:
      - npm run build
  
  common:build:backend:
    cmds:
      - go build -o bin/FG-ABYSS.exe .
  
  common:verify:build:
    cmds:
      - test -f bin/FG-ABYSS.exe || exit 1
```

---

## 6. 检查清单

### 构建前检查

- [ ] 确认输出目录为 `bin/`
- [ ] 确认文件名为 `FG-ABYSS.exe`（全大写）
- [ ] 清理旧的构建文件
- [ ] 验证前端资源已更新

### 构建后验证

- [ ] 检查 `bin/FG-ABYSS.exe` 存在
- [ ] 检查文件大小合理（约 25-30 MB）
- [ ] 运行应用测试基本功能
- [ ] 验证数据库目录结构

### Git 提交前

- [ ] 确认 `bin/` 已加入 `.gitignore`
- [ ] 确认 `frontend/dist/` 已加入 `.gitignore`
- [ ] 确认没有提交构建产物

---

## 7. 常见问题

### Q1: 为什么使用 `bin/` 而不是 `build/`？

**A**: 
- `bin/` 是 Go 项目的标准输出目录
- `build/` 用于存放构建配置和脚本
- 分离配置和产物更清晰

### Q2: 为什么文件名要大写？

**A**:
- 与项目名称保持一致
- Windows 不区分大小写，统一命名更清晰
- 避免跨平台问题

### Q3: 数据库文件应该放在哪里？

**A**:
- 开发环境：`data/app.db`（项目根目录）
- 生产环境：`bin/data/app.db`（可执行文件同级目录）
- 用户数据：`%APPDATA%/FG-ABYSS/app.db`（Windows）

### Q4: 如何处理多平台构建？

**A**:
- Windows: `bin/FG-ABYSS.exe`
- macOS: `bin/FG-ABYSS.app` 或 `bin/FG-ABYSS`
- Linux: `bin/FG-ABYSS`

---

## 8. 相关文档

- [项目优化规范](project-optimization-specification.md) - 完整项目规范
- [文件命名规范](QUICK_CHECKLIST.md) - 文件命名规则
- [scripts/README.md](../../scripts/README.md) - 构建脚本说明

---

**版本**: v1.0  
**更新日期**: 2026-03-15  
**维护者**: FG-ABYSS Team
