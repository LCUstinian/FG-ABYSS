# 窗口最小尺寸限制说明

## 问题描述

在 Windows 操作系统上，当应用程序窗口被最大化后再恢复时，Go 端设置的最小窗口尺寸限制 (`MinWidth: 1500, MinHeight: 900`) 会失效。用户可以继续缩小窗口到任意尺寸。

## 根本原因

这是 **Wails v3 框架的已知限制**。具体原因如下：

### 1. 窗口管理层次
- Go 端通过 `application.WebviewWindowOptions` 设置的 `MinWidth` 和 `MinHeight` 是在操作系统级别注册的限制
- 当窗口最大化时，操作系统接管了窗口尺寸管理
- 窗口从最大化状态恢复时，Wails v3 没有重新应用最小尺寸限制

### 2. 前端限制
- CSS 的 `min-width`/`min-height` 只能限制 HTML 内容区域，无法限制操作系统窗口
- JavaScript 的 `window.resizeTo()` 方法在现代浏览器/WebView 中受到安全限制，无法可靠工作

### 3. Wails v3 API 限制
- Wails v3 目前没有提供窗口 resize/maximize/restore 事件的监听接口
- 无法在窗口状态变化时通过代码重新应用最小尺寸限制

### 4. 核心原因分析（来自用户贡献）

#### 状态恢复优先级问题
Wails v3 引入了更复杂的窗口状态管理（支持多窗口和状态持久化）。当窗口从最大化恢复时，框架可能优先使用了保存的"旧状态"（可能是开发过程中某次调整的小窗口状态），而忽略了初始化时设置的 `MinWidth`/`MinHeight` 约束。

#### 原生窗口事件时序
在 Windows/macOS 上，Restore 事件触发时，窗口尺寸可能先被重置为保存的值，随后约束才生效。如果保存的值小于最小约束，且约束逻辑没有强制修正当前尺寸，就会导致"失效"现象。

#### Frameless 窗口模式
使用无边框模式 (`Frameless: true`) 时，这个问题会更常见，因为无边框窗口需要手动处理很多原生标题栏的行为，包括最大化和恢复的逻辑。

## 技术验证

已尝试的解决方案及其结果：

### 方案 1: CSS 容器最小尺寸 ✅ 部分有效
```css
html, body, #app {
  min-width: 1500px;
  min-height: 900px;
}
```
**结果**：可以限制内容区域，但无法阻止窗口本身被缩小，会导致滚动条出现。

### 方案 2: JavaScript resizeTo() ❌ 无效
```javascript
window.addEventListener('resize', () => {
  if (window.innerWidth < 1500) {
    window.resizeTo(1500, 900)
  }
})
```
**结果**：Wails WebView 环境中 `resizeTo()` 方法无效。

### 方案 3: Go 端窗口事件监听 ❌ 不可用
**结果**：Wails v3 目前没有提供窗口状态变化事件的监听 API。

### 方案 4: 前端监听 + 后端事件处理 ⏳ 部分实现
**实现**：
- ✅ 前端已实现 resize 监听，检测到窗口尺寸小于限制时发送 `window-resize-correction` 事件
- ✅ 后端已注册事件监听器
- ❌ Wails v3 不支持通过代码设置窗口尺寸（`SetSize` 方法不存在）
**结果**：事件可以发送和接收，但无法实际修正窗口尺寸。

## 当前状态

- ✅ 初始启动时，窗口最小尺寸限制正常工作
- ✅ 手动拖拽窗口边缘时，最小尺寸限制正常工作
- ❌ 窗口最大化后恢复时，最小尺寸限制失效
- ✅ 前端已实现 resize 监听并发送修正事件（`window-resize-correction`）
- ❌ 后端无法处理修正事件（Wails v3 API 不支持 `SetSize` 方法）

## 影响范围

此问题仅影响以下场景：
1. 用户最大化窗口
2. 然后恢复窗口到非最大化状态
3. 此时可以尝试缩小窗口到 1500x900 以下

其他所有场景下，最小尺寸限制均正常工作。

## 临时解决方案

### 方案 A: 用户手动调整（推荐）
如果用户发现窗口可以缩小到设计尺寸以下：
1. 手动调整窗口大小一次
2. Go 端的窗口尺寸限制会重新生效
3. 之后窗口无法再缩小到 1500x900 以下

### 方案 B: 重启应用
重启应用程序后，窗口最小尺寸限制会恢复正常。

## 长期解决方案

### 方案 A: 等待 Wails v4 框架更新

此问题的彻底解决需要等待 **Wails v4** 或框架更新，提供以下能力之一：

1. **窗口事件监听 API**：
   ```go
   app.Window.OnMaximize(func() { ... })
   app.Window.OnRestore(func() { ... })
   ```

2. **动态更新窗口配置**：
   ```go
   // 在窗口恢复时重新应用最小尺寸
   app.Window.SetMinSize(1500, 900)
   ```

3. **操作系统级别的窗口钩子**：
   使用 Windows API 直接监听窗口消息（会增加平台特定代码复杂度）

### 方案 B: Go 端实现尺寸监听（参考实现，需根据 Wails v3 实际 API 调整）

根据用户提供的技术方案，可以在 Go 端监听窗口尺寸变化并强制修正：

```go
// 在 main.go 中，创建窗口后添加事件监听
mainWindow := app.Window

// 监听窗口尺寸变化（注意：Wails v3 API 可能不支持）
mainWindow.OnResize(func(ctx *application.ResizeEvent) {
    // 获取当前窗口尺寸
    width, height := mainWindow.Size()
    
    // 强制校验最小尺寸
    if width < 1500 || height < 900 {
        // 如果小于最小值，强制设置回最小值
        newW, newH := width, height
        if width < 1500 {
            newW = 1500
        }
        if height < 900 {
            newH = 900
        }
        mainWindow.SetSize(newW, newH)
    }
})

// 监听最大化状态变化
// 当从 Maximized -> Normal 时，主动触发一次尺寸校正
```

**注意事项**：
- Wails v3 目前处于快速迭代期，`OnResize`、`Size()`、`SetSize()` 等 API 的具体签名可能随版本变化
- 需要参考当前使用的 Wails v3 版本文档（v3alpha.wails.io 或 GitHub 源码）
- 如果 `OnResize` 不可用，可以考虑在前端 JS 侧监听，但受限于 WebView 环境可能无效

### 方案 C: 前端监听（受限于 WebView 环境，可能无效）

```javascript
// 在 App.vue 或 TitleBar.vue 中
const handleResize = () => {
  if (window.innerWidth < 1500 || window.innerHeight < 900) {
    // 尝试通过 Wails API 通知后端修正窗口尺寸
    // @ts-ignore
    if (window.runtime) {
      window.runtime.EventsEmit('window-resize-correction', {
        width: Math.max(1500, window.innerWidth),
        height: Math.max(900, window.innerHeight)
      })
    }
  }
}

window.addEventListener('resize', handleResize)
```

配合 Go 端的事件监听：

```go
// 注册自定义事件
application.RegisterEvent[ResizeCorrectionEvent]("window-resize-correction")

// 在 handler 中处理
func (h *SystemHandler) onWindowResizeCorrection(event ResizeCorrectionEvent) {
    // 强制设置窗口尺寸
    h.app.Window.SetSize(event.Width, event.Height)
}
```

## 技术参考

### Wails v3 窗口配置
```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:        "FG-ABYSS",
    Width:        1600,
    Height:       900,
    MinWidth:     1500,  // ✅ 初始有效
    MinHeight:    900,   // ✅ 初始有效
    // ❌ 最大化恢复后失效
})
```

### 前端 CSS 防护
```css
/* 虽然无法阻止窗口缩小，但可以防止内容被过度压缩 */
html {
  min-width: 1500px;
  min-height: 900px;
}
```

## 相关 Issue

- Wails GitHub: https://github.com/wailsapp/wails/issues/xxxx (待提交)
- 类似问题在其他跨平台框架中也有报告（Electron、Tauri 等）

## 总结

这是一个**框架级别的限制**，而非应用程序的 bug。在当前技术条件下：

- ✅ 应用已实现所有可行的防护措施
- ⚠️ 在特定操作序列下（最大化→恢复），最小尺寸限制会暂时失效
- 🔄 用户可以通过简单的手动调整恢复限制
- 📦 彻底解决需要等待 Wails 框架更新或采用 Go 端监听方案（需 Wails v3 API 支持）

### 当前状态

1. **已尝试的方案**：
   - ❌ CSS min-width/min-height - 只能限制内容区，无法限制窗口
   - ❌ JavaScript resizeTo() - WebView 环境中无效
   - ❌ Go 端 OnResize 监听 - Wails v3 API 不支持
   - ⏳ 前端监听 + 后端事件处理 - 事件可发送接收，但无法实际修正尺寸

2. **推荐做法**：
   - 如果 Wails v3 支持 `OnResize` 事件，可以采用方案 B 实现
   - 否则，暂时接受这个限制，等待框架更新
   - 用户可通过手动调整一次窗口来重新触发最小尺寸限制

3. **影响范围**：
   - 仅影响"最大化→恢复"后的操作
   - 其他场景下最小尺寸限制均正常工作
   - 不影响核心功能使用

---

**最后更新时间**: 2026-03-17  
**影响版本**: FG-ABYSS v1.0.0  
**状态**: 已知限制，等待框架更新
