# FG-ABYSS 应用构建指南

## 问题解决

### ✅ 已解决的问题

1. **白屏问题** - 通过 `fs.Sub()` 正确提取前端资源子目录
2. **控制台黑窗口** - 使用 `-ldflags="-H=windowsgui"` 隐藏控制台窗口

### 🔧 技术细节

#### 1. 白屏问题修复

**问题原因**：
- Wails 3 的 `AssetFileServerFS` 需要的是资源目录的 `fs.FS` 实例
- 直接使用嵌入的 `assets` 会导致路径错误

**解决方案**：
```go
// 获取前端资源子目录
frontendAssets, err := fs.Sub(assets, "frontend/dist")
if err != nil {
    log.Fatalf("Failed to load frontend assets: %v", err)
}

// 使用子目录作为资源处理器
Assets: application.AssetOptions{
    Handler: application.AssetFileServerFS(frontendAssets),
},
```

#### 2. 控制台窗口隐藏

**问题原因**：
- Go 默认构建的是控制台应用
- Windows 会显示黑色的控制台窗口

**解决方案**：
使用 `-ldflags="-H=windowsgui"` 告诉 Go 编译器这是一个 Windows GUI 应用：

```bash
go build -ldflags="-H=windowsgui" -o bin/Fg-ABYSS.exe .
```

## 📋 构建方法

### 方法 1：使用 PowerShell 脚本（推荐）

```powershell
.\build.ps1
```

脚本会自动：
1. 清理旧文件
2. 构建前端
3. 验证文件
4. 构建 GUI 应用

### 方法 2：手动构建

```powershell
# 1. 构建前端
cd frontend
npm run build
cd ..

# 2. 验证前端文件
Test-Path frontend\dist\index.html  # 应该返回 True

# 3. 构建 Windows GUI 应用
go build -ldflags="-H=windowsgui" -o bin\FG-ABYSS.exe .

# 4. 运行应用
.\bin\FG-ABYSS.exe
```

### 方法 3：开发模式（带控制台窗口）

开发时可以显示控制台窗口以便调试：

```powershell
go build -o bin\FG-ABYSS.exe .
.\bin\FG-ABYSS.exe
```

## 🔍 验证清单

构建完成后，请验证：

- [ ] 可执行文件存在：`bin\FG-ABYSS.exe`
- [ ] 文件大小合理（约 20-30 MB）
- [ ] 双击运行时不显示控制台窗口
- [ ] 应用界面正常显示（无白屏）
- [ ] 所有功能正常工作

## 🐛 故障排查

### 问题 1：仍然白屏

**检查点**：
1. 前端是否已构建：`Test-Path frontend\dist\index.html`
2. 检查 `main.go` 中的 `fs.Sub()` 调用
3. 查看应用日志输出

**调试方法**：
```go
// 在 main 函数中添加调试代码
log.Printf("Loading assets from: frontend/dist")
entries, _ := fs.ReadDir(frontendAssets, ".")
for _, entry := range entries {
    log.Printf("  - %s", entry.Name())
}
```

### 问题 2：仍然显示控制台窗口

**检查点**：
1. 构建命令是否包含 `-ldflags="-H=windowsgui"`
2. 是否使用了正确的输出路径

**解决方案**：
```powershell
# 删除旧文件
Remove-Item bin\FG-ABYSS.exe -Force

# 重新构建
go build -ldflags="-H=windowsgui" -o bin\FG-ABYSS.exe .
```

### 问题 3：应用无法启动

**可能原因**：
- 数据库文件损坏
- 端口被占用
- WebView2 未安装

**解决方案**：
```powershell
# 1. 删除数据库文件
Remove-Item data\app.db* -Force

# 2. 检查 WebView2
# Windows 10/11 应该自带 WebView2
# 如果没有，从微软官网下载安装

# 3. 重新运行
.\bin\FG-ABYSS.exe
```

## 📊 构建输出示例

```
=== 开始构建 FG-ABYSS ===

[1/4] 清理旧的构建文件...
  ✓ 已删除旧的可执行文件

[2/4] 构建前端资源...
✓ built in 5.61s
  ✓ 前端构建成功

[3/4] 验证前端文件...
  ✓ 前端文件验证通过

[4/4] 构建 Windows GUI 应用...
  使用 -ldflags=-H=windowsgui 隐藏控制台窗口
  ✓ 构建成功

=== 构建完成 ===
可执行文件位置：D:\Go\FG-ABYSS\bin\FG-ABYSS.exe
文件大小：25.43 MB

运行应用：.\bin\FG-ABYSS.exe
```

## 🎯 开发 vs 生产构建

| 特性 | 开发模式 | 生产模式 |
|------|---------|---------|
| **控制台窗口** | 显示 | 隐藏 |
| **构建命令** | `go build` | `go build -ldflags="-H=windowsgui"` |
| **日志输出** | 完整 | 精简 |
| **调试模式** | 启用 | 禁用 |
| **用途** | 开发调试 | 用户分发 |

## 📝 注意事项

1. **开发时保留控制台**：开发时建议使用普通构建，方便查看日志
2. **分发时使用 GUI 模式**：给用户使用时使用 `-ldflags="-H=windowsgui"`
3. **日志记录**：生产环境建议使用文件日志而不是控制台日志
4. **错误处理**：GUI 模式下，严重错误应该使用对话框而不是控制台输出

## 🔗 相关资源

- [Wails 3 文档](https://v3.wails.io/)
- [Go embed 文档](https://pkg.go.dev/embed)
- [Windows GUI 应用](https://docs.microsoft.com/en-us/windows/win32/learnwin32/learn-to-program-for-windows)
