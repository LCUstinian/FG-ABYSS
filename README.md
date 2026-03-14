# FG-ABYSS

<div align="center">

![FG-ABYSS Banner](https://img.shields.io/badge/FG--ABYSS-网络安全工具-blue)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)
![TypeScript](https://img.shields.io/badge/TypeScript-5.3+-3178C6?logo=typescript)

**一个功能强大的网络安全工具**

[特性](#-核心特性) • [快速开始](#-快速开始) • [文档](#-文档) • [贡献](#-贡献) • [许可证](#-许可证)

</div>

---

## 📖 项目概述

FG-ABYSS 是一款基于现代技术栈构建的桌面端网络安全工具，采用 Wails v3 框架，结合了 Go 语言的后端性能和 Vue 3 前端的美观界面。它提供了项目管理、WebShell 管理、Payload 生成和插件扩展等核心功能，旨在为网络安全专业人员提供一个高效、易用且可扩展的工作平台。

### ✨ 主要亮点

- 🚀 **现代化架构** - Wails v3 + Vue 3 + TypeScript，性能与美观并重
- 🎨 **精美 UI** - 支持浅色/深色主题，自定义强调色，流畅动画效果
- 🔧 **功能丰富** - 项目管理、WebShell、Payload 生成、插件系统
- 🌍 **国际化** - 完整的中英文双语支持
- 📱 **响应式设计** - 适配不同屏幕尺寸，优秀的跨平台体验
- 🔌 **可扩展** - 插件化架构，支持自定义扩展

## 🎯 核心特性

### 1. 项目管理
- 创建和管理多个安全测试项目
- 项目分类和标签系统
- 项目数据持久化存储
- 快速切换和搜索功能

### 2. WebShell 管理
- 支持多种 WebShell 类型（PHP、ASPX、JSP 等）
- WebShell 配置管理
- 连接测试和状态监控
- 批量操作支持

### 3. Payload 生成
- 多种 Payload 类型支持
- 自定义配置选项
- 编码和加密方式
- 一键生成和导出

### 4. 插件系统
- 本地插件管理
- 插件商店（规划中）
- 插件启用/禁用
- 版本管理和更新

### 5. 界面特性
- **主题系统**：浅色/深色/跟随系统
- **自定义颜色**：支持自定义全局强调色
- **字体设置**：可调节字体类型和大小
- **响应式布局**：自适应不同屏幕尺寸
- **流畅动画**：微妙的交互动画效果

## 🛠️ 技术栈

### 后端
- **Go 1.21+** - 高性能后端语言
- **Wails v3** - 桌面应用框架
- **SQLite** - 轻量级数据库
- **GORM** - Go ORM 库

### 前端
- **Vue 3.4+** - 渐进式 JavaScript 框架
- **TypeScript 5.3+** - 类型安全的 JavaScript
- **Vite 5+** - 下一代前端构建工具
- **vue-i18n** - 国际化解决方案
- **Lucide Icons** - 现代化图标库

### UI 组件
- **Naive UI** - Vue 3 组件库
- **自定义主题** - 完整的主题系统
- **CSS 变量** - 动态主题切换

## 📋 环境要求

### 开发环境
- **Go**: 1.21 或更高版本
- **Node.js**: 18.0 或更高版本
- **npm**: 9.0 或更高版本
- **Wails CLI**: v3.0 或更高版本

### 构建环境
- **Windows**: Windows 10/11
- **macOS**: macOS 10.15+
- **Linux**: Ubuntu 18.04+ 或同等发行版

### 内存和存储
- **内存**: 最低 4GB，推荐 8GB+
- **存储**: 至少 500MB 可用空间

## 🚀 快速开始

### 1. 安装依赖

#### 安装 Go
```bash
# Windows (使用 Chocolatey)
choco install go

# macOS (使用 Homebrew)
brew install go

# Linux
sudo apt-get install golang-go
```

#### 安装 Node.js
```bash
# 访问 https://nodejs.org 下载安装
# 或使用版本管理工具 nvm

# Windows
choco install nodejs-lts

# macOS
brew install node

# Linux
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs
```

#### 安装 Wails CLI
```bash
go install github.com/wailsapp/wails/v3/cmd/wails@latest
```

### 2. 克隆项目
```bash
git clone https://github.com/FG-ABYSS/FG-ABYSS.git
cd FG-ABYSS
```

### 3. 安装前端依赖
```bash
cd frontend
npm install
cd ..
```

### 4. 启动开发环境
```bash
# 方式一：使用 Wails 命令（推荐）
wails dev

# 方式二：使用 PowerShell 脚本
.\wails-build.ps1
```

### 5. 构建生产版本
```bash
# 方式一：使用 Wails 命令
wails build

# 方式二：使用构建脚本
.\build.ps1
```

构建完成后，可执行文件将位于 `build/bin` 目录中。

## 📁 项目结构

```
FG-ABYSS/
├── backend/                 # Go 后端代码
│   ├── db/                 # 数据库相关
│   │   └── init.go        # 数据库初始化
│   └── models/            # 数据模型
│       ├── project.go     # 项目模型
│       └── webshell.go    # WebShell 模型
├── frontend/              # Vue 前端代码
│   ├── public/           # 静态资源
│   ├── src/
│   │   ├── api/         # API 调用
│   │   ├── components/  # Vue 组件
│   │   ├── i18n/       # 国际化文件
│   │   ├── styles/     # 全局样式
│   │   ├── types/      # TypeScript 类型
│   │   ├── App.vue     # 根组件
│   │   └── main.ts     # 入口文件
│   ├── index.html       # HTML 模板
│   ├── package.json     # 前端依赖
│   └── vite.config.ts   # Vite 配置
├── build/               # 构建输出目录
├── documentation/       # 项目文档
├── app.go              # Wails 应用逻辑
├── main.go             # 程序入口
└── README.md           # 项目说明
```

## 💻 使用指南

### 界面导航

#### 侧边栏
- **首页** - 数据看板和快速操作
- **项目** - 项目管理功能
- **载荷** - Payload 生成和管理
- **插件** - 插件管理
- **设置** - 系统配置

### 功能使用

#### 创建项目
1. 点击侧边栏「项目」
2. 点击「+」按钮或「新建项目」
3. 填写项目信息
4. 点击「创建」完成

#### 管理 WebShell
1. 进入项目详情页
2. 点击「新建 WebShell」
3. 配置 WebShell 参数
4. 保存并测试连接

#### 生成 Payload
1. 点击侧边栏「载荷」
2. 选择 Payload 类型
3. 配置相关参数
4. 点击「生成」获取 Payload

#### 自定义主题
1. 点击侧边栏「设置」
2. 选择外观设置
3. 调整主题模式（浅色/深色）
4. 自定义强调色
5. 调整字体设置

### 快捷键

| 快捷键 | 功能 |
|--------|------|
| `Ctrl + ,` | 打开设置 |
| `Ctrl + N` | 新建项目 |
| `F5` | 刷新当前页面 |
| `Esc` | 关闭弹窗 |

## 📚 API 文档

### 数据库模型

#### Project（项目）
```go
type Project struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### WebShell
```go
type WebShell struct {
    ID        string    `json:"id"`
    ProjectID string    `json:"project_id"`
    URL       string    `json:"url"`
    Password  string    `json:"password"`
    Type      string    `json:"type"`
    Encoder   string    `json:"encoder"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 系统事件

#### 时间更新事件
```typescript
// 每秒触发一次，携带当前时间信息
app.Event.On("time", (time: string) => {
    console.log("Current time:", time)
})
```

#### 创建窗口事件
```typescript
// 触发创建新窗口
app.Event.Emit("createWindow", {
    title: "New Window",
    width: 800,
    height: 600
})
```

## 🤝 贡献规范

我们欢迎各种形式的贡献！

### 贡献方式

1. **报告 Bug**
   - 使用 GitHub Issues
   - 提供详细的复现步骤
   - 附上环境信息和截图

2. **功能建议**
   - 在 GitHub Discussions 发起讨论
   - 描述功能需求和使用场景
   - 说明期望的实现方式

3. **代码贡献**
   - Fork 项目仓库
   - 创建功能分支 (`git checkout -b feature/AmazingFeature`)
   - 提交更改 (`git commit -m 'Add some AmazingFeature'`)
   - 推送到分支 (`git push origin feature/AmazingFeature`)
   - 创建 Pull Request

### 代码规范

#### Git 提交信息
```
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 重构代码
test: 测试相关
chore: 构建/工具相关
```

#### 代码风格
- Go: 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- TypeScript: 使用 ESLint + Prettier
- Vue: 遵循 [Vue Style Guide](https://vuejs.org/style-guide/)

### 开发流程

1. **设置开发环境**
   ```bash
   git clone https://github.com/YOUR-USERNAME/FG-ABYSS.git
   cd FG-ABYSS
   npm install -C frontend
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **开发并测试**
   ```bash
   wails dev
   ```

4. **提交代码**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **推送并创建 PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)

```
MIT License

Copyright (c) 2024 FG-ABYSS Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 📬 联系方式

### 项目信息
- **作者**: FG-ABYSS Team
- **专业领域**: 网络安全、渗透测试、工具开发
- **GitHub**: [FG-ABYSS](https://github.com/FG-ABYSS)

### 获取支持

- 📖 **文档**: [项目文档](documentation/)
- 🐛 **问题反馈**: [GitHub Issues](https://github.com/FG-ABYSS/FG-ABYSS/issues)
- 💬 **讨论交流**: [GitHub Discussions](https://github.com/FG-ABYSS/FG-ABYSS/discussions)
- 📧 **邮件联系**: [发送邮件](mailto:support@fg-abyss.com)

### 相关链接

- [Wails 官方文档](https://wails.io/)
- [Vue 3 官方文档](https://vuejs.org/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [Naive UI 组件库](https://www.naiveui.com/)

## 🙏 特别感谢

感谢所有为 FG-ABYSS 项目做出贡献的开发者和用户！

### 开源项目致谢
- [Wails](https://github.com/wailsapp/wails) - 优秀的桌面应用框架
- [Vue.js](https://github.com/vuejs/core) - 渐进式 JavaScript 框架
- [TypeScript](https://github.com/microsoft/TypeScript) - 类型化的 JavaScript
- [Naive UI](https://github.com/TuSimple/naive-ui) - Vue 3 组件库
- [Lucide Icons](https://github.com/lucide-icons/lucide) - 美丽的开源图标库

---

<div align="center">

**Made with ❤️ by FG-ABYSS Team**

[返回顶部](#fg-abyss)

</div>
