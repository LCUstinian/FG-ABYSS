# Wails 3 项目完整构建指南

## 问题分析

### 症状
- ✅ `wails3 dev` 运行正常
- ❌ 直接 `go build` 生成的 exe 白屏
- ❌ `build.ps1` 生成的 exe 白屏

### 根本原因

使用普通的 `go build` 命令时，缺少了 Wails 3 的关键构建步骤：

1. **前端资源未正确嵌入** - Wails 需要特殊的资源处理
2. **缺少 Wails 运行时绑定** - TypeScript/Go 绑定未生成
3. **资源路径配置错误** - Wails 使用特殊的资源服务器
4. **缺少平台特定配置** - Windows/macOS/Linux 配置未应用

## ✅ 正确的构建方法

### 方法 1: 使用 Wails 3 构建命令（推荐）

```powershell
# 构建生产版本
wails3 build

# 构建开发版本（带调试信息）
wails3 build -debug

# 构建并打包（Windows）
wails3 build -platform windows/amd64
```

### 方法 2: 使用 Task 构建系统

```powershell
# 安装 Task（如果未安装）
go install github.com/go-task/task/v3/cmd/task@latest

# 构建
task build
```

### 方法 3: 完整的构建流程

```powershell
# 1. 清理
Remove-Item -Path "bin\*" -Recurse -Force
Remove-Item -Path "frontend\dist" -Recurse -Force

# 2. 生成绑定
wails3 task common:gen:bindings

# 3. 构建前端
cd frontend
npm run build
cd ..

# 4. 构建应用
wails3 build

# 5. 运行
.\bin\FG-ABYSS.exe
```

## 🔧 修复后的构建脚本

使用 Wails 3 的正确构建脚本应该是：

```powershell
# wails-build.ps1
Write-Host "=== Wails 3 Build ===" -ForegroundColor Cyan

# 1. 清理
Write-Host "`n[1/5] Cleaning..." -ForegroundColor Yellow
if (Test-Path "bin\FG-ABYSS.exe") {
    Remove-Item -Path "bin\FG-ABYSS.exe" -Force
}

# 2. 生成绑定
Write-Host "`n[2/5] Generating bindings..." -ForegroundColor Yellow
wails3 task common:gen:bindings

# 3. 构建前端
Write-Host "`n[3/5] Building frontend..." -ForegroundColor Yellow
Set-Location frontend
npm run build
Set-Location ..

# 4. 构建应用
Write-Host "`n[4/5] Building application..." -ForegroundColor Yellow
wails3 build

# 5. 完成
Write-Host "`n[5/5] Complete!" -ForegroundColor Green
Write-Host "Executable: bin\FG-ABYSS.exe" -ForegroundColor Cyan
```

## 📋 构建步骤详解

### 步骤 1: 生成绑定

```bash
wails3 task common:gen:bindings
```

这会：
- ✅ 扫描 Go 代码中的 `@wails` 注解
- ✅ 生成 TypeScript 类型定义
- ✅ 创建前端可用的 API 绑定

### 步骤 2: 构建前端

```bash
cd frontend
npm run build
```

这会：
- ✅ 编译 Vue 组件
- ✅ 打包 CSS/JS 资源
- ✅ 生成 `dist/` 目录

### 步骤 3: 构建应用

```bash
wails3 build
```

这会：
- ✅ 嵌入前端资源到二进制文件
- ✅ 应用平台特定配置
- ✅ 生成可执行文件

## 🎯 验证构建结果

### 检查清单

构建完成后，验证：

```powershell
# 1. 检查可执行文件
Test-Path bin\FG-ABYSS.exe  # 应该返回 True

# 2. 检查文件大小（应该在 25-30 MB）
(Get-Item bin\FG-ABYSS.exe).Length / 1MB

# 3. 检查嵌入的资源
go version -m bin\FG-ABYSS.exe

# 4. 运行应用
.\bin\FG-ABYSS.exe
```

### 预期结果

- ✅ 应用正常启动
- ✅ 界面正确显示
- ✅ 无控制台窗口（GUI 模式）
- ✅ 所有功能正常

## 🐛 常见问题排查

### 问题 1: wails3 命令不存在

**解决**:
```powershell
# 安装 Wails 3
go install github.com/wailsapp/wails/v3/cmd/wails@latest

# 验证安装
wails3 version
```

### 问题 2: 前端构建失败

**检查**:
```powershell
# 进入前端目录
cd frontend

# 检查依赖
npm install

# 手动构建
npm run build
```

### 问题 3: 绑定生成失败

**解决**:
```powershell
# 清理绑定
Remove-Item -Path "frontend\bindings" -Recurse -Force

# 重新生成
wails3 task common:gen:bindings
```

### 问题 4: 构建后仍然白屏

**排查步骤**:

1. **检查前端文件**:
   ```powershell
   Test-Path frontend\dist\index.html
   ```

2. **检查资源路径**:
   ```powershell
   Get-ChildItem frontend\dist
   ```

3. **使用调试模式构建**:
   ```powershell
   wails3 build -debug
   ```

4. **查看日志**:
   运行应用时添加 `--debug` 参数

## 📊 构建对比

| 方法 | 命令 | 结果 | 推荐度 |
|------|------|------|--------|
| **Wails 构建** | `wails3 build` | ✅ 正常 | ⭐⭐⭐⭐⭐ |
| **Task 构建** | `task build` | ✅ 正常 | ⭐⭐⭐⭐ |
| **Go Build** | `go build` | ❌ 白屏 | ❌ 不推荐 |
| **Build.ps1** | `.\build.ps1` | ❌ 白屏 | ❌ 不推荐 |

## 🚀 快速开始

### 首次构建

```powershell
# 1. 安装 Wails 3
go install github.com/wailsapp/wails/v3/cmd/wails@latest

# 2. 安装依赖
cd frontend
npm install
cd ..

# 3. 构建
wails3 build

# 4. 运行
.\bin\FG-ABYSS.exe
```

### 开发模式

```powershell
# 热重载开发
wails3 dev
```

### 生产构建

```powershell
# 优化后的生产构建
wails3 build -trimpath -ldflags="-s -w"
```

## 📝 注意事项

1. **不要使用 `go build`** - 会缺少 Wails 的资源处理
2. **始终使用 `wails3 build`** - 确保正确的资源嵌入
3. **定期清理构建** - 避免旧文件干扰
4. **检查绑定更新** - Go 代码变更后重新生成绑定

## 🔗 相关资源

- [Wails 3 文档](https://v3.wails.io/)
- [构建配置](https://v3.wails.io/guides/building)
- [资源嵌入](https://v3.wails.io/guides/assets)
- [项目结构](https://v3.wails.io/gettingstarted/project-structure)
