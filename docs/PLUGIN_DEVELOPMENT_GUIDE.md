# FG-ABYSS 插件开发指南

## 目录

1. [概述](#概述)
2. [插件架构](#插件架构)
3. [快速开始](#快速开始)
4. [插件元数据](#插件元数据)
5. [插件生命周期](#插件生命周期)
6. [插件 API](#插件-api)
7. [权限系统](#权限系统)
8. [依赖管理](#依赖管理)
9. [发布与分发](#发布与分发)
10. [最佳实践](#最佳实践)

## 概述

FG-ABYSS 插件系统允许开发者扩展软件的核心功能。系统支持两种类型的插件：

- **内置插件**：随软件一起发布的官方插件，提供核心功能
- **外置插件**：由第三方开发者开发的扩展插件

## 插件架构

```
┌─────────────────────────────────────┐
│         FG-ABYSS 应用程序            │
├─────────────────────────────────────┤
│          插件管理器 (Manager)         │
├─────────────────────────────────────┤
│          插件加载器 (Loader)          │
├──────────┬──────────┬───────────────┤
│ 内置插件 1 │ 内置插件 2 │ 外置插件 N     │
└──────────┴──────────┴───────────────┘
```

## 快速开始

### 1. 环境要求

- Go 1.16+（支持 plugin 构建模式）
- FG-ABYSS SDK

### 2. 创建插件项目

```bash
mkdir my-plugin
cd my-plugin
go mod init my-plugin
```

### 3. 添加依赖

```go
require (
    fg-abyss v1.0.0
)
```

### 4. 编写插件代码

```go
package main

import (
    "context"
    "fg-abyss/internal/plugin"
)

// 插件元数据
var PluginMetadata = &plugin.PluginMetadata{
    ID:          "example.my_plugin",
    Name:        "我的插件",
    Version:     "1.0.0",
    Description: "这是一个示例插件",
    Author:      "Your Name",
    Type:        plugin.PluginTypeExternal,
    Category:    "tools",
}

// 插件实现
type MyPlugin struct {
    *plugin.BasePlugin
    api plugin.PluginAPI
}

// 创建插件实例
func NewPlugin() plugin.Plugin {
    return &MyPlugin{
        BasePlugin: plugin.NewBasePlugin(PluginMetadata),
    }
}

// 初始化
func (p *MyPlugin) Initialize(ctx context.Context, api plugin.PluginAPI) error {
    p.api = api
    return p.BasePlugin.Initialize(ctx, api)
}
```

### 5. 编译插件

```bash
go build -buildmode=plugin -o my_plugin.so
```

### 6. 安装插件

将编译后的 `.so` 或 `.dll` 文件放入 `plugins` 目录，或在应用中使用"安装插件"功能。

## 插件元数据

```go
type PluginMetadata struct {
    ID              string              // 唯一标识符（必需）
    Name            string              // 插件名称（必需）
    Version         string              // 语义化版本号（必需）
    Description     string              // 描述（必需）
    Author          string              // 作者（必需）
    Type            PluginType          // 插件类型（自动设置）
    Category        string              // 分类
    Tags            []string            // 标签
    Homepage        string              // 主页 URL
    Repository      string              // 代码仓库
    License         string              // 许可证
    MinAppVersion   string              // 最低应用版本
    Dependencies    []PluginDependency  // 依赖列表
    Permissions     []Permission        // 权限声明
    CreatedAt       time.Time           // 创建时间
    UpdatedAt       time.Time           // 更新时间
}
```

### 版本号规范

使用语义化版本号（Semantic Versioning）：`MAJOR.MINOR.PATCH`

- `MAJOR`: 不兼容的 API 变更
- `MINOR`: 向后兼容的功能新增
- `PATCH`: 向后兼容的问题修正

### 依赖声明示例

```go
Dependencies: []PluginDependency{
    {
        ID:       "builtin.file_manager",
        Version:  "^1.0.0",
        Required: true,
    },
    {
        ID:       "example.other_plugin",
        Version:  ">=2.0.0",
        Required: false,
    },
}
```

## 插件生命周期

```
加载 (Load) → 初始化 (Initialize) → 启动 (Start) → 运行 (Running)
                                              ↓
关闭 (Shutdown) ← 停止 (Stop) ← 禁用 (Disabled)
```

### 生命周期方法

```go
// 初始化插件
func (p *MyPlugin) Initialize(ctx context.Context, api plugin.PluginAPI) error {
    // 注册命令、事件等
    return nil
}

// 启动插件
func (p *MyPlugin) Start(ctx context.Context) error {
    // 启动后台任务、连接等
    return nil
}

// 停止插件
func (p *MyPlugin) Stop(ctx context.Context) error {
    // 停止后台任务、保存状态等
    return nil
}

// 关闭插件
func (p *MyPlugin) Shutdown(ctx context.Context) error {
    // 清理资源、关闭连接等
    return nil
}
```

## 插件 API

### 配置管理

```go
// 获取配置
value, err := api.GetConfig("database_url")

// 设置配置
err := api.SetConfig("database_url", "mysql://localhost:3306/db")
```

### 命令注册

```go
// 注册命令
err := api.RegisterCommand("my_command", func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // 处理命令
    return map[string]interface{}{
        "result": "success",
    }, nil
})
```

### 事件系统

```go
// 触发事件
err := api.EmitEvent("user.login", map[string]interface{}{
    "user_id": "123",
    "time":    time.Now(),
})

// 订阅事件
err := api.SubscribeEvent("user.login", func(event string, data interface{}) error {
    // 处理事件
    return nil
})
```

### 日志记录

```go
api.Log(plugin.LogLevelInfo, "插件已启动", map[string]interface{}{
    "version": "1.0.0",
})
```

### 调用其他插件

```go
result, err := api.CallPlugin("builtin.file_manager", "ListFiles", map[string]interface{}{
    "path": "/home/user",
})
```

## 权限系统

### 可用权限

| 权限 | 说明 |
|------|------|
| `fs:readwrite` | 文件系统读写 |
| `network` | 网络访问 |
| `database` | 数据库访问 |
| `execute` | 命令执行 |
| `clipboard` | 剪贴板访问 |
| `settings` | 设置访问 |
| `plugin` | 插件管理 |

### 权限声明

```go
Permissions: []plugin.Permission{
    plugin.PermissionFileSystem,
    plugin.PermissionNetwork,
}
```

### 权限检查

```go
if !api.HasPermission(plugin.PermissionFileSystem) {
    return nil, errors.New("permission denied")
}
```

## 依赖管理

### 依赖解析顺序

1. 检查依赖插件是否已安装
2. 检查版本是否匹配
3. 加载依赖插件（如果未加载）
4. 初始化当前插件

### 循环依赖处理

系统会检测循环依赖并拒绝加载，确保插件依赖图是无环的。

## 发布与分发

### 1. 编译插件

```bash
# Linux/macOS
go build -buildmode=plugin -o my_plugin.so

# Windows
go build -buildmode=plugin -o my_plugin.dll
```

### 2. 打包发布

```bash
# 创建发布包
mkdir my-plugin
cp my_plugin.so my-plugin/
cp README.md my-plugin/
cp LICENSE my-plugin/

# 压缩
tar -czf my-plugin-1.0.0.tar.gz my-plugin/
```

### 3. 发布渠道

- GitHub Releases
- 插件商店（未来）
- 手动分发

## 最佳实践

### 1. 错误处理

```go
func (p *MyPlugin) Start(ctx context.Context) error {
    if err := p.initializeResource(); err != nil {
        p.api.Log(plugin.LogLevelError, "Failed to start", map[string]interface{}{
            "error": err.Error(),
        })
        return err
    }
    return nil
}
```

### 2. 资源清理

```go
func (p *MyPlugin) Shutdown(ctx context.Context) error {
    // 确保资源被清理
    defer p.cleanup()
    return p.BasePlugin.Shutdown(ctx)
}
```

### 3. 并发安全

```go
type MyPlugin struct {
    *plugin.BasePlugin
    api   plugin.PluginAPI
    mu    sync.RWMutex
    data  map[string]interface{}
}
```

### 4. 性能优化

- 避免在事件处理器中执行耗时操作
- 使用缓存减少重复计算
- 合理使用后台 goroutine

### 5. 安全考虑

- 验证所有用户输入
- 最小权限原则
- 定期更新依赖

## 示例代码

完整的示例插件请参考：`plugins/example/hello_plugin.go`

## 故障排查

### 常见问题

1. **插件无法加载**
   - 检查插件文件路径
   - 验证插件元数据
   - 查看应用日志

2. **权限错误**
   - 确认权限声明
   - 检查用户授权

3. **依赖冲突**
   - 检查依赖版本
   - 更新插件版本

### 日志查看

```bash
# 查看应用日志
tail -f data/logs/app.log
```

## 获取帮助

- 文档：查看本指南
- 示例：参考 `plugins/example/` 目录
- 问题反馈：提交 Issue

---

**最后更新**: 2026-03-17
**版本**: 1.0.0
