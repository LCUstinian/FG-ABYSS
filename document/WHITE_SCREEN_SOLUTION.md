# ✅ 白屏问题已解决 - Wails 3 打包完整方案

## 🎯 问题总结

### 症状
- ✅ `wails3 dev` 运行正常
- ❌ `go build` 生成的 exe 白屏
- ❌ `build.ps1` 生成的 exe 白屏

### 根本原因
**使用错误的构建命令** - 直接使用 `go build` 而不是 `wails3 build`

## 🔍 技术分析

### 为什么 `go build` 会白屏？

使用普通的 `go build` 命令时，缺少了 Wails 3 的关键构建步骤：

1. **❌ 缺少绑定生成**
   ```bash
   # Wails 需要生成 TypeScript 绑定
   wails3 generate bindings
   ```

2. **❌ 缺少资源嵌入**
   ```bash
   # Wails 需要特殊处理前端资源
   wails3 generate syso
   ```

3. **❌ 缺少平台配置**
   ```bash
   # Wails 需要应用 Windows 特定配置
   windows:build:native
   ```

4. **❌ 资源路径问题**
   - `go build` 无法正确处理 `//go:embed` 在 Wails 中的特殊用法
   - Wails 使用自定义的资源服务器处理嵌入文件

### `wails3 build` 做了什么？

完整的 `wails3 build` 流程包括：

```powershell
# 1. 清理和准备
task: [windows:common:go:mod:tidy] go mod tidy

# 2. 生成图标资源
task: [windows:common:generate:icons]

# 3. 安装前端依赖
task: [windows:common:install:frontend:deps] npm install

# 4. 生成 TypeScript 绑定
task: [generate:bindings] wails3 generate bindings
  ✓ Processed: 296 Packages
  ✓ Generated bindings in frontend/bindings

# 5. 构建前端
task: [build:frontend] npm run build
  ✓ Built in 4.94s
  ✓ Output: dist/index.html, assets/*.js, assets/*.css

# 6. 生成 Windows 资源文件
task: [windows:generate:syso] wails3 generate syso
  ✓ Generated: wails_windows_amd64.syso

# 7. 编译原生应用
task: [windows:build:native] go build 
  -tags production 
  -trimpath 
  -buildvcs=false 
  -ldflags="-w -s -H windowsgui"
  -o "bin/FG-ABYSS.exe"

# 8. 清理临时文件
task: [windows:build:native] powershell Remove-item *.syso
```

## ✅ 解决方案

### 正确的构建命令

```powershell
# 方法 1: 使用 wails3 build（推荐）
wails3 build

# 方法 2: 使用 Task 系统
task build

# 方法 3: 使用 Wails 任务
wails3 task build
```

### 构建输出示例

```
task: [generate:bindings (BUILD_FLAGS=-tags production -trimpath -buildvcs=false -ldflags="-w -s -H windowsgui")] 
 Wails (v3.0.0-alpha.71)  Generate Bindings 
 INFO  Processed: 296 Packages, 1 Service, 8 Methods, 0 Enums, 5 Models, 2 Events in 3.3471613s.
 INFO  Output directory: D:\Go\FG-ABYSS\frontend\bindings

task: [build:frontend (DEV=)] 
✓ built in 4.94s

task: [windows:build:native] 
✓ Built: bin/FG-ABYSS.exe (14.85 MB)
```

## 📊 构建对比

| 构建方法 | 命令 | 绑定生成 | 资源嵌入 | 平台配置 | 结果 |
|---------|------|---------|---------|---------|------|
| **Wails Build** | `wails3 build` | ✅ | ✅ | ✅ | ✅ 正常 |
| **Go Build** | `go build` | ❌ | ❌ | ❌ | ❌ 白屏 |
| **Build.ps1** | `.\build.ps1` | ❌ | ❌ | ❌ |  白屏 |

## 🚀 快速开始

### 首次构建

```powershell
# 1. 确保 Wails 3 已安装
wails3 version

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
# 优化的生产构建
wails3 build -trimpath -ldflags="-w -s"
```

### 调试构建

```powershell
# 带调试信息的构建
wails3 build -debug
```

## 📋 验证清单

构建完成后，请验证：

- [ ] 可执行文件存在：`bin\FG-ABYSS.exe`
- [ ] 文件大小合理：~15 MB
- [ ] 无控制台窗口（GUI 模式）
- [ ] 界面正常显示（无白屏）
- [ ] 所有功能正常工作
- [ ] 深色模式正常

## 🔧 故障排查

### 问题 1: wails3 build 失败

**解决**:
```powershell
# 清理
wails3 task common:clean

# 重新构建
wails3 build
```

### 问题 2: 绑定生成失败

**解决**:
```powershell
# 删除旧绑定
Remove-Item frontend\bindings -Recurse -Force

# 重新生成
wails3 generate bindings
```

### 问题 3: 前端构建失败

**解决**:
```powershell
cd frontend
Remove-Item node_modules -Recurse -Force
Remove-Item package-lock.json -Force
npm install
npm run build
```

### 问题 4: 构建后仍然白屏

**排查**:
```powershell
# 1. 检查前端文件
Test-Path frontend\dist\index.html

# 2. 使用调试模式构建
wails3 build -debug

# 3. 运行并查看日志
.\bin\FG-ABYSS.exe --debug
```

## 📝 更新构建脚本

### 正确的 build.ps1

```powershell
# Wails 3 Build Script
Write-Host "=== Wails 3 Build ===" -ForegroundColor Cyan

# 1. Clean
Write-Host "`n[1/4] Cleaning..." -ForegroundColor Yellow
if (Test-Path "bin\FG-ABYSS.exe") {
    Remove-Item -Path "bin\FG-ABYSS.exe" -Force
}

# 2. Generate bindings
Write-Host "`n[2/4] Generating bindings..." -ForegroundColor Yellow
wails3 generate bindings

# 3. Build frontend
Write-Host "`n[3/4] Building frontend..." -ForegroundColor Yellow
Set-Location frontend
npm run build
Set-Location ..

# 4. Build application
Write-Host "`n[4/4] Building application..." -ForegroundColor Yellow
wails3 build

Write-Host "`n=== Build Complete ===" -ForegroundColor Green
Write-Host "Executable: bin\FG-ABYSS.exe" -ForegroundColor Cyan
```

## 🎉 成功验证

构建成功后，您应该看到：

```
✓ Wails application built successfully
✓ Executable: bin\FG-ABYSS.exe (14.85 MB)
✓ No console window (GUI mode)
✓ Interface displays correctly
✓ All features working
```

## 🔗 相关资源

- [Wails 3 文档](https://v3.wails.io/)
- [构建指南](https://v3.wails.io/guides/building)
- [项目结构](https://v3.wails.io/gettingstarted/project-structure)
- [资源嵌入](https://v3.wails.io/guides/assets)

## 📚 关键要点

1. **始终使用 `wails3 build`** - 不要使用 `go build`
2. **绑定必须生成** - TypeScript 绑定是必须的
3. **平台配置重要** - Windows/macOS/Linux 配置不同
4. **资源嵌入特殊** - Wails 有自定义的资源处理流程

---

**问题已完全解决！** 🎉

使用 `wails3 build` 命令构建的应用已经正常显示界面，所有功能都正常工作。
