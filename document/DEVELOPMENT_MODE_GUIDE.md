# Wails3 开发模式使用指南

## 问题说明

在 Wails3 开发过程中，您可能会遇到以下情况：
- 通过 `wails3 dev` 命令运行应用时，显示正常 ✅
- 通过浏览器访问 `http://localhost:9245/` 时，出现"画中画"或内容缺失 ❌

## 原因分析

### Wails3 开发模式的工作原理

```
wails3 dev 命令执行流程:
┌─────────────────────────────────────────┐
│ 1. 启动 Wails 后端 (Go)                  │
│    - 编译 Go 代码                         │
│    - 启动后端服务                        │
├─────────────────────────────────────────┤
│ 2. 启动前端开发服务器 (Vite)             │
│    - 运行 `npm run dev`                  │
│    - 监听端口：9245                      │
├─────────────────────────────────────────┤
│ 3. 启动 Wails WebView 窗口               │
│    - 创建桌面应用窗口                    │
│    - 加载前端界面                        │
└─────────────────────────────────────────┘
```

### 为什么浏览器预览会出问题？

当您通过浏览器访问 `http://localhost:9245/` 时：

1. **访问的是 Vite 开发服务器**，不是 Wails 应用
2. **Wails 运行时检测到浏览器环境**，尝试适配
3. **缺少 Wails 后端 API**，导致功能异常
4. **可能出现自我渲染**，导致"画中画"效果

## 正确的开发方式

### ✅ 推荐：使用 `wails3 dev` 命令

**这是 Wails3 的标准开发模式**：

```bash
# 在项目根目录执行
wails3 dev
```

**特点**：
- ✅ 启动完整的 Wails 应用（后端 + 前端）
- ✅ 支持热重载（前端代码修改自动刷新）
- ✅ 支持 Go 代码热重载（部分修改）
- ✅ 完整的 Wails API 支持
- ✅ 正确的桌面应用环境

**工作流程**：
```
1. 执行 wails3 dev
   ↓
2. Wails 编译后端代码
   ↓
3. 启动 Vite 开发服务器 (端口 9245)
   ↓
4. Wails 创建 WebView 窗口
   ↓
5. 应用正常运行 ✅
```

### ⚠️ 不推荐：通过浏览器访问开发服务器

**不要这样做**：
```bash
# ❌ 错误的方式
打开浏览器访问 http://localhost:9245/
```

**问题**：
- ❌ 缺少 Wails 后端 API
- ❌ 可能出现"画中画"效果
- ❌ 功能不完整
- ❌ 与桌面应用环境不一致

## 如果需要在浏览器中测试前端

如果您确实需要在浏览器中单独测试前端（不带 Wails 后端），可以：

### 方案 A：使用 Mock 数据

1. **创建 Mock API**：
   ```typescript
   // src/api/mock.ts
   export const mockWebShells = [
     { id: '1', url: 'http://test.com', ... }
   ];
   ```

2. **修改组件使用 Mock 数据**：
   ```typescript
   // 开发模式下使用 Mock
   if (import.meta.env.DEV) {
     data.value = mockWebShells;
   } else {
     data.value = await fetchWebShells();
   }
   ```

3. **单独启动前端**：
   ```bash
   cd frontend
   npm run dev
   ```

### 方案 B：配置代理服务器

1. **修改 `vite.config.ts`**：
   ```typescript
   export default defineConfig({
     server: {
       proxy: {
         '/api': {
           target: 'http://localhost:8080', // Wails 后端
           changeOrigin: true
         }
       }
     }
   });
   ```

2. **同时运行 Wails 后端和 Vite 前端**：
   ```bash
   # 终端 1: 运行 Wails 后端
   go run main.go
   
   # 终端 2: 运行 Vite 前端
   cd frontend
   npm run dev
   ```

## 常见问题解答

### Q1: 为什么 `wails3 dev` 正常，浏览器访问就不正常？

**A**: Wails3 是一个桌面应用框架，不是纯 Web 应用。`wails3 dev` 会：
- 启动 Go 后端
- 创建 WebView 窗口
- 提供完整的 Wails API

而浏览器访问只能访问 Vite 开发服务器，缺少这些关键组件。

### Q2: 如何在浏览器中调试前端代码？

**A**: 有几种方法：

**方法 1**: 在 `wails3 dev` 运行的应用中调试
- 按 `Ctrl + Shift + I`（Windows）或 `Cmd + Option + I`（Mac）
- 打开 WebView 开发者工具

**方法 2**: 使用 Mock 数据单独运行前端
```bash
cd frontend
npm run dev
# 访问 http://localhost:5173/
```

### Q3: 可以修改端口号吗？

**A**: 可以，修改 `vite.config.ts`：
```typescript
export default defineConfig({
  server: {
    port: 3000, // 修改端口
  }
});
```

### Q4: 如何让浏览器预览也正常工作？

**A**: 这需要：
1. 创建 Mock 数据层
2. 条件化 Wails API 调用
3. 单独运行前端服务器

但**不推荐**这样做，因为：
- 增加维护成本
- 与桌面应用环境不一致
- 可能掩盖真实问题

## 最佳实践

### ✅ 推荐做法

1. **使用 `wails3 dev` 进行开发**
   ```bash
   wails3 dev
   ```

2. **在 WebView 中调试**
   - 按 `Ctrl + Shift + I` 打开开发者工具
   - 使用 Console、Elements、Network 等面板

3. **使用热重载**
   - 修改 Vue 组件自动刷新
   - 修改 Go 代码部分情况自动重载

4. **测试生产构建**
   ```bash
   # 构建生产版本
   wails3 build
   
   # 运行测试
   ./bin/prod/FG-ABYSS.exe
   ```

### ❌ 避免的做法

1. **不要通过浏览器访问 `wails3 dev` 的开发服务器**
2. **不要依赖浏览器特定的 API**（如 localStorage 在某些情况下可能不同）
3. **不要假设浏览器环境等同于 Wails 环境**

## 开发工作流建议

### 标准开发流程

```bash
# 1. 启动开发模式
wails3 dev

# 2. 编辑代码
# - 前端：src/components/*.vue
# - 后端：app.go, backend/*.go

# 3. 查看效果
# - 应用窗口自动刷新
# - 或使用 Ctrl+Shift+I 调试

# 4. 测试功能
# - 在应用窗口中操作
# - 检查开发者工具

# 5. 提交前测试生产构建
wails3 build
./bin/prod/FG-ABYSS.exe
```

### 前端单独开发流程（仅限 UI 开发）

```bash
# 1. 进入前端目录
cd frontend

# 2. 安装依赖（首次）
npm install

# 3. 启动开发服务器
npm run dev

# 4. 访问 http://localhost:5173/
# 注意：只有 UI，没有后端功能

# 5. 完成 UI 开发后，使用 wails3 dev 测试完整功能
cd ..
wails3 dev
```

## 总结

**核心要点**：

1. ✅ **`wails3 dev` 是正确的开发方式**
2. ❌ **不要通过浏览器访问开发服务器**
3. 🔧 **如需浏览器调试，使用 Mock 数据**
4. 🎯 **在 WebView 中使用开发者工具调试**

Wails3 是桌面应用框架，应该在其目标环境（桌面应用窗口）中开发和测试，而不是浏览器。

---

**相关文档**：
- [Wails3 官方文档](https://wails.io/)
- [开发模式配置](https://v3.wails.io/docs/guides/dev-mode)
- [前端开发指南](./frontend/README.md)
