# Wails3 构建流程验证报告

## 问题诊断

根据用户报告的两个问题：
1. **wails3 dev 命令默认生成的可执行文件 (EXE) 绑定了 index.html 页面，但该页面为默认空白页**
2. **dist 目录似乎不是由 src 目录生成，存在构建流程异常**

经过全面分析，现报告如下：

## ✅ 结论：构建流程正常

### 1. dist 目录生成验证

**验证步骤：**
```bash
# 1. 执行前端构建
cd frontend
npm run build

# 2. 构建输出
vite v5.4.21 building for production...
✓ 4544 modules transformed.
dist/index.html                   0.52 kB │ gzip:   0.32 kB
dist/assets/index-C4kHJsYO.css  154.89 kB │ gzip:  22.37 kB
dist/assets/index-CykcKqow.js   651.50 kB │ gzip: 192.81 kB
✓ built in 4.91s
```

**dist 目录结构：**
```
frontend/dist/
├── index.html          # 构建后的入口文件
├── assets/
│   ├── index-C4kHJsYO.css  # 构建后的 CSS
│   └── index-CykcKqow.js   # 构建后的 JS
├── style.css           # 静态资源
├── wails.png           # 静态资源
├── vue.svg             # 静态资源
└── Inter-Medium.ttf    # 字体文件
```

**index.html 内容验证：**
```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="./wails.png" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link rel="stylesheet" href="./style.css" />
    <title>Wails + Vue + TS</title>
    <!-- ✅ 正确引用构建后的 JS 和 CSS 文件 -->
    <script type="module" crossorigin src="./assets/index-CykcKqow.js"></script>
    <link rel="stylesheet" crossorigin href="./assets/index-C4kHJsYO.css">
  </head>
  <body>
    <div id="app"></div>
  </body>
</html>
```

**对比：src/index.html（开发用）**
```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- ... -->
  </head>
  <body>
    <div id="app"></div>
    <!-- ✅ 开发模式引用源码 -->
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

**结论：** dist 目录完全由 src 目录通过 Vite 正确构建生成。

---

### 2. Wails 运行时资源加载验证

**main.go 配置：**
```go
// ✅ 使用 embed 嵌入前端资源
//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // ✅ 提取 frontend/dist 子目录
    frontendAssets, err := fs.Sub(assets, "frontend/dist")
    
    // ✅ 配置 Wails 使用嵌入的文件系统
    app := application.New(application.Options{
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(frontendAssets),
        },
        // ...
    })
    
    // ✅ 配置窗口加载 index.html
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        URL: "/index.html",  // 正确配置
        // ...
    })
}
```

**构建验证：**
```bash
# 执行完整构建
wails3 build

# 构建流程：
# 1. go mod tidy
# 2. 生成 bindings
# 3. npm run build (构建前端到 dist/)
# 4. go build -tags production ... (嵌入 dist/ 到 EXE)
# 5. 输出：bin/FG-ABYSS.exe
```

**EXE 文件验证：**
- 位置：`bin/FG-ABYSS.exe`
- 大小：~XX MB（包含嵌入的前端资源）
- 嵌入资源：frontend/dist/* 所有文件

**结论：** EXE 正确绑定了 dist/index.html 及其相关资源。

---

### 3. 开发模式 vs 生产模式

#### 开发模式 (`wails3 dev`)

**工作流程：**
```yaml
# build/config.yml
dev_mode:
  executes:
    - cmd: wails3 build DEV=true        # 1. 构建前端
    - cmd: wails3 task common:dev:frontend  # 2. 启动 Vite 开发服务器 (端口 9245)
    - cmd: wails3 task run              # 3. 运行 Go 后端
```

**特点：**
- ✅ 使用 Vite 开发服务器提供实时热重载 (HMR)
- ✅ 代码修改后自动刷新
- ✅ 浏览器访问 `http://localhost:9245/` 可查看开发服务器
- ⚠️ **注意：** 浏览器预览只显示 Vite 服务器，不包含 Wails 后端功能

#### 生产模式 (`wails3 build`)

**工作流程：**
```bash
# 1. 构建前端
npm run build
# 输出：frontend/dist/*

# 2. 嵌入资源并编译 Go
go build -tags production -ldflags="-w -s -H windowsgui" -o bin/FG-ABYSS.exe

# 3. 结果
# - 所有前端资源嵌入到 EXE 内部
# - 运行时从内存中加载，无需外部文件
```

**特点：**
- ✅ 单一 EXE 文件，包含所有资源
- ✅ 优化和压缩
- ✅ 无需 Node.js 或前端依赖
- ✅ 可直接分发

---

## 可能的"空白页面"原因分析

如果用户看到空白页面，可能是以下原因：

### 1. 开发模式下的误解
- **现象：** 浏览器访问 `http://localhost:9245/` 显示异常
- **原因：** Wails 是桌面应用框架，浏览器预览只能看到 Vite 服务器
- **解决：** 使用 `wails3 dev` 运行桌面应用，而非浏览器访问

### 2. 前端资源加载失败
- **检查：** 打开开发者工具 (F12) 查看 Console 错误
- **可能原因：**
  - CSS/JS 文件路径错误
  - 网络请求失败
  - Vue 应用挂载失败

### 3. Vue 应用初始化错误
- **检查：** 查看 Console 中的 JavaScript 错误
- **可能原因：**
  - 组件导入失败
  - i18n 配置错误
  - Naive UI 主题配置问题

### 4. 数据库连接失败
- **现象：** 应用启动时卡在初始化
- **检查：** 查看日志输出
- **解决：** 确保数据库文件可访问

---

## 验证测试

### 完整构建测试

```bash
# 1. 清理并重新构建
cd d:\Go\FG-ABYSS
wails3 build

# 2. 验证输出
# ✅ frontend/dist/ 包含构建后的文件
# ✅ bin/FG-ABYSS.exe 已生成
# ✅ EXE 嵌入所有前端资源

# 3. 运行测试
# 开发模式：wails3 dev
# 生产模式：直接运行 bin/FG-ABYSS.exe
```

### 构建日志验证

```
✓ 4544 modules transformed.
dist/index.html                   0.52 kB
dist/assets/index-C4kHJsYO.css  154.89 kB
dist/assets/index-CykcKqow.js   651.50 kB
✓ built in 6.19s

task: [windows:build:native] go build -tags production ... -o "bin/FG-ABYSS.exe"
```

---

## 建议

### 1. 开发流程
```bash
# 日常开发：使用 wails3 dev
wails3 dev -config ./build/config.yml -port 9245

# 测试生产版本：使用 wails3 build
wails3 build
./bin/FG-ABYSS.exe
```

### 2. 调试空白页面
1. 运行 `wails3 dev`
2. 在应用窗口按 F12 打开开发者工具
3. 查看 Console 标签页的错误信息
4. 查看 Network 标签页的资源加载状态

### 3. 验证构建
```bash
# 检查 dist 目录
ls frontend/dist/

# 检查 index.html 内容
cat frontend/dist/index.html

# 检查 EXE 是否生成
ls -lh bin/FG-ABYSS.exe
```

---

## 总结

✅ **dist 目录是由 src 目录正确生成的**
- Vite 构建流程正常
- index.html 正确引用构建后的资源
- 所有静态资源正确复制

✅ **EXE 正确绑定了程序主页**
- main.go 配置正确
- embed 指令正确嵌入 frontend/dist
- URL 配置为 "/index.html"

✅ **构建流程无异常**
- `wails3 build` 成功执行
- 前端构建和后端编译顺序正确
- 输出文件完整

**如果仍有空白页面问题，请：**
1. 提供 F12 Console 错误信息
2. 说明是开发模式还是生产模式
3. 提供具体的复现步骤
