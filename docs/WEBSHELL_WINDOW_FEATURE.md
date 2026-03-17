# WebShell 独立窗口控制功能实现指南

## 功能概述

本功能允许用户在项目模块的 WebShell 表格中右键点击某个 WebShell 实例，通过右键菜单中的"控制"按钮打开一个独立的 Wails3 窗口，在该窗口中可以执行完整的 WebShell 操作。

## 已完成的实现

### 1. 前端组件

#### WebShellControlWindow.vue
**位置**: `frontend/src/components/WebShellControlWindow.vue`

**功能**:
- 独立的 WebShell 控制窗口界面
- 包含 4 个标签页：终端、文件管理、数据库管理、命令执行
- 窗口工具栏显示 WebShell 信息和连接状态
- 支持连接/断开控制
- 自动从 URL 参数获取 WebShell ID

**集成方式**:
```vue
<!-- 在 App.vue 或主路由中添加 -->
<template>
  <div v-if="currentRoute === 'webshell-control'">
    <WebShellControlWindow />
  </div>
</template>
```

### 2. 后端服务

#### WebShellWindowManager
**位置**: `internal/app/services/webshell_window_manager.go`

**功能**:
- 管理所有 WebShell 独立窗口
- 窗口的打开、关闭、状态查询
- 防止重复打开同一 WebShell 的窗口

**方法**:
```go
- OpenWindow(webshellID, name, url string) error
- CloseWindow(webshellID string) error
- GetWindow(webshellID string) (*application.Window, error)
- IsWindowOpen(webshellID string) bool
- GetAllWindows() []WindowInfo
- CloseAllWindows()
```

#### SystemHandler 扩展
**位置**: `internal/app/handlers/system_handler.go`

**新增方法**:
```go
- OpenWebShellWindow(ctx context.Context, req *OpenWebShellWindowRequest) error
- CloseWebShellWindow(ctx context.Context, req *CloseWebShellWindowRequest) error
- GetWebShellWindowStatus(ctx context.Context, webshellID string) (map[string]interface{}, error)
```

### 3. 前端集成

#### ProjectsContent.vue 修改
**位置**: `frontend/src/components/ProjectsContent.vue`

**修改内容**:
1. 右键菜单已包含"控制"选项（key: 'enter'）
2. 添加 `openWebShellControlWindow` 函数
3. 点击"控制"按钮时调用该函数

**关键代码**:
```typescript
const openWebShellControlWindow = async (webshell: WebShell) => {
  // 使用 Wails API 打开新窗口
  if (window.runtime) {
    await window.runtime.EventsEmit('open-webshell-window', {
      id: webshell.id,
      name: webshell.name || webshell.url,
      url: webshell.url,
    })
  }
}
```

## 待完成的集成步骤

### 步骤 1: 在 main.go 中注册窗口管理器

```go
// 在 main.go 中添加

// 1. 导入服务
import "fg-abyss/internal/app/services"

// 2. 创建窗口管理器（在创建 handlers 之前）
log.Println("=== Creating Window Manager ===")
windowManager := services.NewWebShellWindowManager(app)
log.Println("=== Window Manager created ===")

// 3. 将窗口管理器传递给需要它的 handler
systemHandler := handlers.NewSystemHandlerWithWindowManager(appService, windowManager)
```

### 步骤 2: 实现窗口创建事件处理

```go
// 在 main.go 中添加事件监听

// 监听打开 WebShell 窗口的事件
app.On("open-webshell-window", func(ctx context.Context, e *application.Event[map[string]interface{}]) {
  data := e.Payload
  webshellID, _ := data["id"].(string)
  name, _ := data["name"].(string)
  url, _ := data["url"].(string)
  
  // 使用窗口管理器打开窗口
  err := windowManager.OpenWindow(webshellID, name, url)
  if err != nil {
    log.Printf("Failed to open WebShell window: %v", err)
  }
})
```

### 步骤 3: 配置路由

在 `App.vue` 中添加路由支持：

```vue
<script setup lang="ts">
import { ref } from 'vue'
import HomeContent from './components/HomeContent.vue'
import ProjectsContent from './components/ProjectsContent.vue'
import PayloadsContent from './components/PayloadsContent.vue'
import PluginsContent from './components/PluginsContent.vue'
import SettingsContent from './components/SettingsContent.vue'
import WebShellControlWindow from './components/WebShellControlWindow.vue'

const currentRoute = ref('home')

// 从 URL 获取路由
const updateRouteFromURL = () => {
  const hash = window.location.hash.slice(1)
  if (hash.startsWith('/webshell-control')) {
    currentRoute.value = 'webshell-control'
  } else if (hash.startsWith('/projects')) {
    currentRoute.value = 'projects'
  } else if (hash.startsWith('/payloads')) {
    currentRoute.value = 'payloads'
  } else if (hash.startsWith('/plugins')) {
    currentRoute.value = 'plugins'
  } else if (hash.startsWith('/settings')) {
    currentRoute.value = 'settings'
  } else {
    currentRoute.value = 'home'
  }
}

// 监听 hash 变化
window.addEventListener('hashchange', updateRouteFromURL)
updateRouteFromURL()
</script>

<template>
  <div class="app-container">
    <!-- 主窗口内容 -->
    <template v-if="currentRoute !== 'webshell-control'">
      <HomeContent v-if="currentRoute === 'home'" />
      <ProjectsContent v-else-if="currentRoute === 'projects'" />
      <PayloadsContent v-else-if="currentRoute === 'payloads'" />
      <PluginsContent v-else-if="currentRoute === 'plugins'" />
      <SettingsContent v-else-if="currentRoute === 'settings'" />
    </template>
    
    <!-- WebShell 控制窗口 -->
    <WebShellControlWindow v-else />
  </div>
</template>
```

### 步骤 4: Wails 窗口配置

由于 Wails v3 的窗口 API 可能有所不同，需要根据实际文档调整。以下是参考实现：

```go
// 在 webshell_window_manager.go 中实现 OpenWindow

func (m *WebShellWindowManager) OpenWindow(webshellID, name, url string) error {
  m.mu.Lock()
  defer m.mu.Unlock()

  // 检查窗口是否已存在
  if _, exists := m.windows[webshellID]; exists {
    return fmt.Errorf("window for WebShell %s already exists", webshellID)
  }

  // 创建新窗口
  windowURL := fmt.Sprintf("index.html#/webshell-control?id=%s", webshellID)
  
  // 使用 Wails v3 API 创建窗口（需要根据实际 API 调整）
  window, err := m.app.NewWindow(application.WindowOptions{
    Title:  fmt.Sprintf("WebShell Control - %s", name),
    URL:    windowURL,
    Width:  1200,
    Height: 800,
    X:      100,
    Y:      100,
  })
  
  if err != nil {
    return err
  }
  
  m.windows[webshellID] = window
  return nil
}
```

## 功能特性

### ✅ 已实现
1. WebShell 独立窗口组件
2. 右键菜单"控制"按钮
3. 窗口管理器服务
4. 4 个功能标签页（终端、文件、数据库、命令）

### ⏳ 待实现
1. main.go 中窗口管理器集成
2. 窗口创建事件处理
3. 路由配置
4. Wails 窗口 API 适配

## 使用说明

1. **打开控制窗口**:
   - 在项目模块中，找到目标 WebShell
   - 右键点击该行
   - 选择"控制"按钮
   - 新窗口将打开

2. **窗口功能**:
   - **终端**: 执行命令、查看输出
   - **文件管理**: 浏览、上传、下载文件
   - **数据库管理**: 连接数据库、执行查询
   - **命令执行**: 使用预设命令模板

3. **窗口管理**:
   - 每个 WebShell 只能打开一个控制窗口
   - 关闭窗口会自动断开连接
   - 窗口状态独立于主窗口

## 技术要点

1. **URL 参数传递**: 通过 URL hash 传递 WebShell ID
2. **窗口通信**: 使用 Wails Events 进行窗口间通信
3. **状态同步**: 窗口管理器维护所有窗口状态
4. **资源清理**: 关闭窗口时自动清理连接和资源

## 下一步

1. 根据 Wails v3 文档完善窗口创建 API
2. 测试窗口打开和关闭功能
3. 验证所有标签页功能正常
4. 添加窗口状态持久化
5. 实现窗口间数据同步

---

**文档版本**: 1.0
**最后更新**: 2026-03-17
**状态**: 实现中
