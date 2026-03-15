# Wails3 升级快速参考

## ✅ 升级完成

**Wails3 v3.0.0-alpha.71 → v3.0.0-alpha.74**

## 🚀 快速命令

### 构建应用
```powershell
wails3 build
```

### 开发模式
```powershell
wails3 dev
```

### 运行应用
```powershell
.\bin\FG-ABYSS.exe
```

## 📦 版本信息

| 组件 | 版本 | 状态 |
|------|------|------|
| **Wails3** | v3.0.0-alpha.74 | ✅ 最新 |
| **Go** | go1.26.0 | ✅ 兼容 |
| **前端** | 已更新 | ✅ 最新 |

## 🔧 常用命令

```powershell
# 查看版本
wails3 version

# 构建生产版本
wails3 build

# 开发模式（热重载）
wails3 dev

# 生成绑定
wails3 generate bindings

# 清理构建
wails3 task common:clean

# 更新依赖
go mod tidy
cd frontend && npm update
```

## 📝 重要变更

### 1. 窗口 URL 配置
```go
// main.go
URL: "/index.html"  // 显式指定 index.html
```

### 2. 数据库初始化
- 自动迁移表结构
- 自动删除 name 字段
- 创建默认项目和示例数据

### 3. 资源路径
- 使用 `fs.Sub()` 提取资源
- vite.config.ts 配置 `base: './'`

## ⚠️ 注意事项

### 已知警告（不影响使用）
1. **前端安全警告**: 6 个中等漏洞
2. **CSS 嵌套警告**: 某些浏览器不支持
3. **代码大小警告**: 部分 chunk > 500KB

### 建议
- 定期关注 Wails3 更新
- 生产环境前充分测试
- 考虑处理安全警告

## 📋 验证清单

- [x] 版本正确 (v3.0.0-alpha.74)
- [x] 构建正常
- [x] 应用启动
- [x] 界面显示
- [x] 功能正常

## 🔗 相关文档

- [UPGRADE_LOG.md](./UPGRADE_LOG.md) - 详细升级日志
- [WHITE_SCREEN_SOLUTION.md](./WHITE_SCREEN_SOLUTION.md) - 白屏问题解决方案
- [DATABASE_INIT.md](./DATABASE_INIT.md) - 数据库初始化说明
- [BUILD_GUIDE.md](./BUILD_GUIDE.md) - 构建指南

## 💡 故障排查

### 如果构建失败
```powershell
# 清理
wails3 task common:clean

# 重新构建
wails3 build
```

### 如果白屏
```powershell
# 检查 main.go 中的 URL 配置
# 确保是 URL: "/index.html"

# 重新构建
wails3 build
```

### 如果依赖冲突
```powershell
# Go 依赖
go mod tidy

# 前端依赖
cd frontend
npm install
```

## 📞 获取帮助

- 查看文档：`wails3 docs`
- Discord 社区：https://discord.gg/J8vvK4c
- GitHub: https://github.com/wailsapp/wails

---

**最后更新**: 2026-03-14
**状态**: ✅ 运行正常
