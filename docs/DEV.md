# FG-ABYSS 开发文档 (DEV)

> **项目名称**: 非攻 - 渊渟 (FG-ABYSS)  
> **版本**: v1.0.0  
> **最后更新**: 2026-04-13  
> **状态**: 开发准备就绪

---

## 🤖 TRAE CN 快速理解指南

**如果你是用 TRAE CN 阅读本文档，请先理解以下关键点**：

### 技术栈

```yaml
前端：
  - Vue 3.4+ (Composition API + <script setup>)
  - TypeScript 5.3+
  - Naive UI 2.38+ (组件库)
  - Pinia 2.1+ (状态管理)
  - Vue Router 4.3+
  - vue-i18n (国际化)

后端：
  - Rust 1.70+
  - Tauri V2 (桌面框架)
  - tokio (异步运行时)
  - reqwest (HTTP 客户端)
  - rusqlite (SQLite 数据库)
  - ring (加密库)
  - mlua (Lua 插件)
  - wasmtime (WebAssembly 插件)

数据库:
  - SQLite (本地存储)
```

### 项目结构关键点

```
FG-ABYSS/
├── src/                    # 前端
│   ├── components/
│   │   ├── payload/       # 载荷组件（本地）
│   │   ├── project/       # 项目组件（本地数据）
│   │   └── console/       # 控制台组件（远程）
│   └── views/
│       ├── PayloadView/   # /payload (本地)
│       └── ProjectView/   # /project (远程)
│
└── src-tauri/              # 后端
    ├── commands/
    │   ├── payload.rs     # 载荷生成（本地）
    │   ├── project.rs     # 项目管理（本地）
    │   └── webshell.rs    # WebShell 连接（远程）
    └── core/
        ├── payload/       # 载荷生成逻辑（本地）
        ├── project/       # 项目管理逻辑（本地）
        └── webshell/      # WebShell 连接逻辑（远程）
```

### 开发注意事项

1. **载荷管理模块** (`/payload`):
   - 仅生成代码，不涉及远程操作
   - 前端：`components/payload/`, `views/PayloadView/`
   - 后端：`commands/payload.rs`, `core/payload/`

2. **项目管理模块** (`/project`):
   - 管理远程 WebShell，有远程操作
   - 前端：`components/project/`, `views/ProjectView/`
   - 后端：`commands/project.rs`, `core/project/`, `core/webshell/`

3. **控制台窗口**:
   - 右键 WebShell 打开的独立窗口
   - 以 TAB 形式展示插件（文件、数据库、终端）
   - 前端：`components/console/`
   - 后端：`commands/console.rs`, `core/console/`

### 开发命令

```bash
# 开发模式（热重载）
pnpm tauri dev

# 构建生产版本
pnpm tauri build

# 运行测试
pnpm test
```

---

## 1. 开发环境设置

### 1.1 系统要求

#### 最低要求

- **操作系统**: Windows 10 / macOS 11 / Linux (Ubuntu 20.04+)
- **CPU**: 双核 2.0 GHz
- **内存**: 4 GB RAM
- **磁盘**: 10 GB 可用空间

#### 推荐配置

- **操作系统**: Windows 11 / macOS 12+ / Linux (Ubuntu 22.04+)
- **CPU**: 四核 2.5 GHz+
- **内存**: 8 GB RAM
- **磁盘**: 20 GB SSD

### 1.2 依赖安装

#### Rust 环境

```bash
# 安装 Rust (使用 rustup)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 验证安装
rustc --version
cargo --version

# 安装所需 Rust 版本 (1.70+)
rustup install 1.70.0
rustup default 1.70.0
```

#### Node.js 环境

```bash
# 安装 Node.js (LTS 版本)
# Windows/macOS: 下载安装包
# Linux:
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 验证安装
node --version
npm --version

# 安装 pnpm (推荐)
npm install -g pnpm
```

#### 系统依赖

**Windows**:

```powershell
# 安装 Visual C++ Redistributable
# 安装 WebView2 (Windows 10/11 自带)
```

**macOS**:

```bash
# 安装 Xcode Command Line Tools
xcode-select --install
```

**Linux**:

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y \
  build-essential \
  curl \
  wget \
  libssl-dev \
  libgtk-3-dev \
  libayatana-appindicator3-dev \
  librsvg2-dev \
  libsoup2.4-dev \
  libjavascriptcoregtk-4.0-dev \
  libwebkit2gtk-4.0-dev

# Fedora/RHEL
sudo dnf install -y \
  gcc \
  gcc-c++ \
  make \
  openssl-devel \
  gtk3-devel \
  libappindicator-gtk3-devel \
  librsvg2-devel \
  webkit2gtk3-devel
```

### 1.3 项目克隆

```bash
# 克隆项目
git clone https://github.com/your-org/fg-abyss.git
cd fg-abyss

# 安装前端依赖
pnpm install

# 安装 Rust 依赖
cd src-tauri
cargo fetch
```

### 1.4 IDE 配置

#### VS Code (推荐)

**安装扩展**:

- Rust Analyzer
- Vue Language Features (Volar)
- TypeScript Vue Plugin (Volar)
- Tauri
- ESLint
- Prettier

**settings.json**:

```json
{
  "rust-analyzer.checkOnSave.command": "clippy",
  "rust-analyzer.cargo.features": "all",
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "[rust]": {
    "editor.defaultFormatter": "rust-lang.rust-analyzer"
  },
  "[vue]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

#### RustRover (可选)

- 自动识别 Rust 项目
- 需要安装 Vue/TypeScript 插件

### 1.5 开发服务器

```bash
# 启动开发服务器 (热重载)
pnpm tauri dev

# 或分别启动
# 前端
cd src
pnpm dev

# 后端 (新终端)
cd src-tauri
cargo watch -x run
```

***

## 2. 项目结构

### 2.1 目录结构

```
FG-ABYSS/
├── package.json                      # Node.js 依赖配置
├── pnpm-lock.yaml                    # 包锁定文件
├── README.md                         # 项目说明
├── LICENSE                           # 开源许可
├── .gitignore                        # Git 忽略规则
├── .env.example                      # 环境变量示例
├── tsconfig.json                     # TypeScript 配置
├── vite.config.ts                    # Vite 构建配置
├── tailwind.config.js                # Tailwind CSS 配置
│
├── src/                              # 前端源代码
│   ├── main.ts                       # Vue 应用入口
│   ├── App.vue                       # 根组件
│   ├── index.html                    # HTML 入口
│   │
│   ├── components/                   # UI 组件
│   │   ├── layout/                   # 布局组件
│   │   │   ├── CustomTitlebar.vue    # 自定义标题栏
│   │   │   ├── NavigationMenu.vue    # 导航菜单
│   │   │   ├── StatusBar.vue         # 底部状态栏
│   │   │   └── index.ts
│   │   │
│   │   ├── payload/                  # 载荷管理组件（本地代码生成）
│   │   │   ├── PayloadConfig.vue     # 载荷配置表单
│   │   │   ├── PayloadPreview.vue    # 代码预览
│   │   │   ├── PayloadHistory.vue    # 历史记录
│   │   │   ├── PayloadTemplate.vue   # 模板管理
│   │   │   └── index.ts
│   │   │
│   │   ├── project/                  # 项目管理组件（远程 WebShell）
│   │   │   ├── ProjectTree.vue       # 项目树
│   │   │   ├── WebShellList.vue      # WebShell 列表
│   │   │   ├── WebShellForm.vue      # WebShell 表单
│   │   │   ├── WebShellContextMenu.vue # 右键菜单
│   │   │   └── index.ts
│   │   │
│   │   ├── console/                  # 控制台窗口组件
│   │   │   ├── ConsoleWindow.vue     # 控制台窗口
│   │   │   ├── ConsoleTabs.vue       # 插件 TAB 栏
│   │   │   └── index.ts
│   │   │
│   │   ├── file/                     # 文件管理组件（控制台插件）
│   │   │   ├── FileBrowser.vue       # 文件浏览器（双栏布局）
│   │   │   ├── FileUploader.vue      # 文件上传组件（拖拽上传）
│   │   │   ├── FileEditor.vue        # 文件编辑器（代码编辑）
│   │   │   ├── FileContextMenu.vue   # 文件右键菜单
│   │   │   └── index.ts
│   │   │
│   │   ├── database/                 # 数据库管理组件（控制台插件）
│   │   │   ├── DatabaseManager.vue   # 数据库管理器
│   │   │   ├── QueryEditor.vue       # SQL 查询编辑器
│   │   │   ├── ResultTable.vue       # 查询结果表格
│   │   │   ├── DatabaseSelector.vue  # 数据库选择器
│   │   │   └── index.ts
│   │   │
│   │   ├── terminal/                 # 终端组件（控制台插件）
│   │   │   ├── TerminalView.vue      # 终端视图（xterm.js）
│   │   │   ├── TerminalInput.vue     # 命令输入框
│   │   │   ├── CommandHistory.vue    # 命令历史
│   │   │   └── index.ts
│   │   │
│   │   ├── plugin/                   # 插件组件
│   │   │   ├── PluginList.vue        # 插件列表
│   │   │   ├── PluginManager.vue     # 插件管理
│   │   │   └── index.ts
│   │   │
│   │   └── common/                   # 通用组件
│   │       ├── Loading.vue
│   │       ├── ErrorDisplay.vue
│   │       ├── ConfirmDialog.vue
│   │       ├── CodeEditor.vue        # 代码编辑器（通用）
│   │       ├── SearchBar.vue         # 搜索栏（通用）
│   │       ├── FilterPanel.vue       # 过滤面板（通用）
│   │       └── index.ts
│   │
│   ├── views/                        # 页面视图
│   │   ├── Dashboard.vue             # 首页
│   │   ├── PayloadView.vue           # 载荷管理页（/payload）
│   │   ├── ProjectView.vue           # 项目管理页（/project）
│   │   ├── PluginView.vue            # 插件管理页（/plugin）
│   │   └── SettingsView.vue          # 设置页（/settings）
│   │
│   ├── stores/                       # Pinia 状态管理
│   │   ├── index.ts
│   │   ├── payload.ts                # 载荷配置状态
│   │   ├── project.ts                # 项目状态
│   │   ├── webshell.ts               # WebShell 连接状态
│   │   ├── plugin.ts                 # 插件状态
│   │   └── app.ts                    # 应用全局状态
│   │
│   ├── router/                       # 路由配置
│   │   ├── index.ts
│   │   ├── routes.ts                 # 路由定义
│   │   └── guards.ts                 # 路由守卫
│   │
│   ├── api/                          # API 调用层
│   │   ├── index.ts
│   │   ├── payload.ts                # 载荷生成 API
│   │   ├── project.ts                # 项目管理 API
│   │   ├── webshell.ts               # WebShell 连接 API
│   │   ├── plugin.ts                 # 插件 API
│   │   └── console.ts                # 控制台 API
│   │
│   ├── types/                        # TypeScript 类型
│   │   ├── index.ts
│   │   ├── payload.ts                # 载荷类型
│   │   ├── project.ts                # 项目类型
│   │   ├── webshell.ts               # WebShell 类型
│   │   ├── plugin.ts                 # 插件类型
│   │   └── common.ts                 # 通用类型
│   │
│   ├── composables/                  # Composables
│   │   ├── index.ts
│   │   ├── usePayload.ts             # 载荷生成逻辑
│   │   ├── useProject.ts             # 项目管理逻辑
│   │   ├── useWebShell.ts            # WebShell 连接逻辑
│   │   ├── usePlugin.ts              # 插件加载逻辑
│   │   └── useI18n.ts                # 国际化逻辑
│   │
│   ├── utils/                        # 工具函数
│   │   ├── index.ts
│   │   ├── format.ts                 # 格式化工具
│   │   ├── validation.ts             # 验证工具
│   │   ├── storage.ts                # 存储工具（localStorage）
│   │   ├── export.ts                 # 导出工具（Payload/项目/查询结果）
│   │   ├── search.ts                 # 搜索工具（通用搜索算法）
│   │   ├── filter.ts                 # 过滤工具（数组过滤）
│   │   ├── sort.ts                   # 排序工具（多字段排序）
│   │   └── download.ts               # 下载工具（文件下载）
│   │
│   ├── i18n/                         # 国际化
│   │   ├── index.ts
│   │   ├── zh-CN.ts
│   │   └── en-US.ts
│   │
│   ├── assets/                       # 静态资源
│   │   ├── images/
│   │   ├── icons/
│   │   └── fonts/
│   │
│   └── styles/                       # 样式
│       ├── main.css
│       ├── variables.css
│       └── components/
│
├── src-tauri/                        # Rust 后端
│   ├── Cargo.toml                    # Rust 依赖配置
│   ├── Cargo.lock                    # Rust 依赖锁定
│   ├── tauri.conf.json               # Tauri 配置
│   ├── build.rs                      # 构建脚本
│   ├── capabilities/                 # 权限配置
│   │   └── default.json
│   ├── icons/                        # 应用图标
│   │   ├── icon.ico
│   │   ├── icon.icns
│   │   └── icon.png
│   └── src/                          # Rust 源代码
│       ├── main.rs                   # 桌面端入口
│       ├── lib.rs                    # 核心逻辑入口
│       │
│       ├── commands/                 # Tauri Commands（IPC 接口）
│       │   ├── mod.rs
│       │   ├── payload.rs            # 载荷生成命令（本地）
│       │   ├── project.rs            # 项目管理命令（远程）
│       │   ├── webshell.rs           # WebShell 连接命令（远程）
│       │   ├── console.rs            # 控制台命令（远程）
│       │   └── plugin.rs             # 插件管理命令（本地）
│       │
│       ├── core/                     # 核心业务逻辑
│       │   ├── mod.rs
│       │   │
│       │   ├── payload/              # 载荷生成模块（本地）
│       │   │   ├── mod.rs
│       │   │   ├── generator.rs      # 载荷生成器
│       │   │   ├── php.rs            # PHP 载荷
│       │   │   ├── jsp.rs            # JSP 载荷
│       │   │   ├── asp.rs            # ASP 载荷
│       │   │   ├── aspx.rs           # ASPX 载荷
│       │   │   ├── template.rs       # 模板管理
│       │   │   └── history.rs        # 历史记录
│       │   │
│       │   ├── project/              # 项目管理模块（本地数据）
│       │   │   ├── mod.rs
│       │   │   ├── manager.rs        # 项目管理器
│       │   │   ├── entity.rs         # 项目实体
│       │   │   └── recycle_bin.rs    # 回收站
│       │   │
│       │   ├── webshell/             # WebShell 连接模块（远程）
│       │   │   ├── mod.rs
│       │   │   ├── connection.rs     # 连接管理
│       │   │   ├── session.rs        # 会话管理
│       │   │   └── health_check.rs   # 健康检查
│       │   │
│       │   ├── crypto/               # 加密通信模块（远程）
│       │   │   ├── mod.rs
│       │   │   ├── aes.rs            # AES 加密
│       │   │   ├── xor.rs            # XOR 加密
│       │   │   └── key_manager.rs    # 密钥管理
│       │   │
│       │   ├── http/                 # HTTP 客户端模块（远程）
│       │   │   ├── mod.rs
│       │   │   ├── client.rs         # HTTP 客户端
│       │   │   ├── proxy.rs          # 代理支持
│       │   │   ├── rate_limiter.rs   # 速率限制
│       │   │   └── retry.rs          # 重试机制
│       │   │
│       │   ├── console/              # 控制台模块（远程）
│       │   │   ├── mod.rs
│       │   │   ├── file_manager.rs   # 文件管理插件
│       │   │   ├── database_manager.rs # 数据库管理插件
│       │   │   ├── terminal.rs       # 终端管理插件
│       │   │   └── plugin_host.rs    # 插件宿主
│       │   │
│       │   └── obfuscation/          # 代码混淆模块（本地）
│       │       ├── mod.rs
│       │       ├── engine.rs         # 混淆引擎
│       │       └── tree_sitter.rs    # 语法分析
│       │
│       ├── infra/                    # 基础设施
│       │   ├── mod.rs
│       │   ├── database/             # SQLite 数据库
│       │   │   ├── mod.rs
│       │   │   ├── pool.rs           # 连接池
│       │   │   └── models/           # 数据模型
│       │   │       ├── mod.rs
│       │   │       ├── payload_config.rs  # 载荷配置
│       │   │       ├── project.rs    # 项目
│       │   │       ├── webshell.rs   # WebShell
│       │   │       └── plugin.rs     # 插件
│       │   │
│       │   ├── metrics/              # 指标收集
│       │   │   ├── mod.rs
│       │   │   └── collector.rs
│       │   │
│       │   └── health/               # 健康检查
│       │       ├── mod.rs
│       │       └── checker.rs
│       │
│       ├── plugins/                  # 插件系统（本地）
│       │   ├── mod.rs
│       │   ├── loader.rs             # 插件加载器
│       │   ├── sandbox.rs            # 沙箱隔离
│       │   ├── runtime.rs            # 运行时
│       │   └── types.rs              # 插件类型定义
│       │
│       ├── error.rs                  # 错误处理
│       ├── config.rs                 # 配置管理
│       ├── logging.rs                # 日志系统
│       └── utils.rs                  # 工具函数
│
├── migrations/                       # 数据库迁移
│   ├── 001_create_shells_table.sql
│   ├── 002_create_projects_table.sql
│   └── 003_create_plugins_table.sql
│
├── scripts/                          # 脚本
│   ├── build.ps1                     # Windows 构建
│   ├── build.sh                      # Linux/macOS 构建
│   └── dev-setup.ps1                 # 开发环境设置
│
├── tests/                            # 测试
│   ├── integration/
│   │   ├── shell_tests.rs
│   │   ├── crypto_tests.rs
│   │   └── payload_tests.rs
│   └── e2e/
│       └── main_flow_tests.rs
│
└── benches/                          # 基准测试
    ├── crypto_bench.rs
    └── payload_bench.rs
```

### 2.2 模块职责

#### 前端模块

| 模块 | 职责 | 技术 |
|------|------|------|
| `components/payload/` | 载荷配置、预览、历史记录（本地） | Vue 3 + TypeScript |
| `components/project/` | 项目树、WebShell 列表（本地数据） | Vue 3 + TypeScript |
| `components/console/` | 控制台窗口、插件 TAB（远程操作） | Vue 3 + TypeScript |
| `views/PayloadView/` | 载荷管理页面（/payload，本地） | Vue Router |
| `views/ProjectView/` | 项目管理页面（/project，远程） | Vue Router |
| `stores/payload/` | 载荷配置状态（本地） | Pinia |
| `stores/project/` | 项目和 WebShell 状态（本地数据） | Pinia |
| `stores/webshell/` | WebShell 连接状态（远程会话） | Pinia |

#### 后端模块

| 模块 | 职责 | 是否远程 | 技术 |
|------|------|---------|------|
| `commands/payload.rs` | 载荷生成命令 | ❌ 否 | Tauri IPC |
| `commands/project.rs` | 项目管理命令 | ❌ 否 | Tauri IPC |
| `commands/webshell.rs` | WebShell 连接命令 | ✅ 是 | Tauri IPC |
| `commands/console.rs` | 控制台插件命令 | ✅ 是 | Tauri IPC |
| `core/payload/` | 载荷生成、模板、历史 | ❌ 否 | Rust |
| `core/project/` | 项目管理、回收站 | ❌ 否 | Rust |
| `core/webshell/` | WebShell 连接、会话 | ✅ 是 | Rust + reqwest |
| `core/console/` | 文件/数据库/终端插件 | ✅ 是 | Rust |
| `core/crypto/` | AES/XOR加密通信 | ✅ 是 | Rust + ring |
| `infra/database/` | SQLite 数据存储 | ❌ 否 | rusqlite |
| `plugins/` | 插件系统（加载/沙箱） | ❌ 否 | mlua + wasmtime |

---

## 2.3 状态管理设计 (Pinia Stores)

### 2.3.1 payload.ts - 载荷配置状态

```typescript
import { defineStore } from 'pinia'

interface PayloadConfig {
  id: string
  name: string
  url: string
  password: string
  payloadType: 'php' | 'jsp' | 'asp' | 'aspx'
  encryption: 'aes-256-gcm' | 'xor' | 'base64'
  obfuscationLevel: 1 | 2 | 3
  tags?: string[]              // 标签数组
  group?: string               // 分组名称
  createdAt: number
  updatedAt: number
}

interface PayloadTemplate {
  id: string
  name: string
  type: 'php' | 'jsp' | 'asp' | 'aspx'
  code: string
  isBuiltIn: boolean
}

interface PayloadHistory {
  id: string
  configId: string
  generatedCode: string
  generatedAt: number
}

interface PayloadState {
  configs: PayloadConfig[]
  templates: PayloadTemplate[]
  history: PayloadHistory[]
  selectedConfig: PayloadConfig | null
  isGenerating: boolean
  // 搜索/过滤/排序
  searchQuery: string
  filterType: 'all' | 'php' | 'jsp' | 'asp' | 'aspx'
  filterTag: string | null
  sortBy: 'name' | 'createdAt' | 'updatedAt' | 'type'
  sortOrder: 'asc' | 'desc'
  tags: string[]
}

export const usePayloadStore = defineStore('payload', {
  state: (): PayloadState => ({
    configs: [],
    templates: [],
    history: [],
    selectedConfig: null,
    isGenerating: false,
    // 搜索/过滤/排序
    searchQuery: '',
    filterType: 'all',
    filterTag: null,
    sortBy: 'name',
    sortOrder: 'asc',
    tags: [],
  }),
  
  getters: {
    configCount: (state) => state.configs.length,
    templateCount: (state) => state.templates.length,
    historyCount: (state) => state.history.length,
    
    // 搜索/过滤/排序后的配置列表
    filteredConfigs: (state) => {
      let result = [...state.configs]
      
      // 搜索
      if (state.searchQuery) {
        const query = state.searchQuery.toLowerCase()
        result = result.filter(config =>
          config.name.toLowerCase().includes(query) ||
          config.url.toLowerCase().includes(query)
        )
      }
      
      // 类型过滤
      if (state.filterType !== 'all') {
        result = result.filter(config => config.payloadType === state.filterType)
      }
      
      // 标签过滤
      if (state.filterTag) {
        result = result.filter(config => config.tags?.includes(state.filterTag))
      }
      
      // 排序
      result.sort((a, b) => {
        const multiplier = state.sortOrder === 'asc' ? 1 : -1
        if (state.sortBy === 'name') {
          return a.name.localeCompare(b.name) * multiplier
        } else if (state.sortBy === 'type') {
          return a.payloadType.localeCompare(b.payloadType) * multiplier
        } else {
          return (a[state.sortBy] - b[state.sortBy]) * multiplier
        }
      })
      
      return result
    },
    
    // 所有唯一标签
    uniqueTags: (state) => {
      const tagSet = new Set<string>()
      state.configs.forEach(config => {
        config.tags?.forEach(tag => tagSet.add(tag))
      })
      return Array.from(tagSet)
    },
  },
  
  actions: {
    // 搜索相关
    setSearchQuery(query: string) {
      this.searchQuery = query
    },
    
    // 过滤相关
    setFilterType(type: 'all' | 'php' | 'jsp' | 'asp' | 'aspx') {
      this.filterType = type
    },
    
    setFilterTag(tag: string | null) {
      this.filterTag = tag
    },
    
    clearFilters() {
      this.searchQuery = ''
      this.filterType = 'all'
      this.filterTag = null
    },
    
    // 排序相关
    setSortBy(sortBy: 'name' | 'createdAt' | 'updatedAt' | 'type') {
      this.sortBy = sortBy
    },
    
    toggleSortOrder() {
      this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc'
    },
    
    // 标签管理
    addTag(tag: string) {
      if (!this.tags.includes(tag)) {
        this.tags.push(tag)
      }
    },
    
    removeTag(tag: string) {
      this.tags = this.tags.filter(t => t !== tag)
    },
    
    // 配置管理
    async loadConfigs() {
      // 从后端加载载荷配置列表
    },
    
    async createConfig(config: PayloadConfig) {
      // 创建新的载荷配置
    },
    
    async updateConfig(id: string, config: Partial<PayloadConfig>) {
      // 更新载荷配置
    },
    
    async deleteConfig(id: string) {
      // 软删除载荷配置
    },
    
    // 标签操作
    async updateConfigTags(id: string, tags: string[]) {
      // 更新配置的标签
    },
    
    // 导入/导出
    async exportConfigs(ids: string[]) {
      // 导出选中的载荷配置
    },
    
    async importConfigs(file: File) {
      // 导入载荷配置
    },
    
    async generatePayload(configId: string) {
      // 生成载荷代码
      this.isGenerating = true
      try {
        // 调用后端 API
      } finally {
        this.isGenerating = false
      }
    },
    
    async loadTemplates() {
      // 加载载荷模板
    },
    
    async loadHistory(configId?: string) {
      // 加载载荷历史记录
    },
  },
})
```

### 2.3.2 project.ts - 项目状态

```typescript
import { defineStore } from 'pinia'

interface Project {
  id: string
  name: string
  description?: string
  webshellCount: number
  tags?: string[]              // 项目标签
  group?: string               // 项目分组
  color?: string               // 分组颜色
  createdAt: number
  updatedAt: number
}

interface ProjectState {
  projects: Project[]
  selectedProject: string | null
  isLoading: boolean
  // 搜索/过滤/排序
  searchQuery: string
  filterGroup: string | null
  sortBy: 'name' | 'createdAt' | 'updatedAt' | 'webshellCount'
  sortOrder: 'asc' | 'desc'
  // 分组管理
  groups: Array<{
    id: string
    name: string
    color: string
    order: number
  }>
}

export const useProjectStore = defineStore('project', {
  state: (): ProjectState => ({
    projects: [],
    selectedProject: null,
    isLoading: false,
    // 搜索/过滤/排序
    searchQuery: '',
    filterGroup: null,
    sortBy: 'name',
    sortOrder: 'asc',
    // 分组管理
    groups: [],
  }),
  
  getters: {
    projectCount: (state) => state.projects.length,
    activeProject: (state) => 
      state.projects.find(p => p.id === state.selectedProject),
    
    // 搜索/过滤/排序后的项目列表
    filteredProjects: (state) => {
      let result = [...state.projects]
      
      // 搜索
      if (state.searchQuery) {
        const query = state.searchQuery.toLowerCase()
        result = result.filter(project =>
          project.name.toLowerCase().includes(query) ||
          project.description?.toLowerCase().includes(query)
        )
      }
      
      // 分组过滤
      if (state.filterGroup) {
        result = result.filter(project => project.group === state.filterGroup)
      }
      
      // 排序
      result.sort((a, b) => {
        const multiplier = state.sortOrder === 'asc' ? 1 : -1
        if (state.sortBy === 'name') {
          return a.name.localeCompare(b.name) * multiplier
        } else if (state.sortBy === 'webshellCount') {
          return (a.webshellCount - b.webshellCount) * multiplier
        } else {
          return (a[state.sortBy] - b[state.sortBy]) * multiplier
        }
      })
      
      return result
    },
    
    // 按分组统计项目数
    projectsByGroup: (state) => {
      const map = new Map<string, number>()
      state.projects.forEach(project => {
        const group = project.group || 'ungrouped'
        map.set(group, (map.get(group) || 0) + 1)
      })
      return map
    },
  },
  
  actions: {
    // 搜索/过滤/排序
    setSearchQuery(query: string) {
      this.searchQuery = query
    },
    
    setFilterGroup(group: string | null) {
      this.filterGroup = group
    },
    
    setSortBy(sortBy: 'name' | 'createdAt' | 'updatedAt' | 'webshellCount') {
      this.sortBy = sortBy
    },
    
    toggleSortOrder() {
      this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc'
    },
    
    clearFilters() {
      this.searchQuery = ''
      this.filterGroup = null
    },
    
    // 分组管理
    async createGroup(group: { name: string; color: string; order: number }) {
      // 创建新分组
    },
    
    async updateGroup(groupId: string, group: Partial<{ name: string; color: string; order: number }>) {
      // 更新分组
    },
    
    async deleteGroup(groupId: string) {
      // 删除分组（将项目设为无分组）
    },
    
    async assignProjectToGroup(projectId: string, groupId: string | null) {
      // 将项目分配到分组
    },
    
    // 项目管理
    async loadProjects() {
      this.isLoading = true
      try {
        // 从后端加载项目列表
      } finally {
        this.isLoading = false
      }
    },
    
    async createProject(project: Omit<Project, 'id' | 'createdAt' | 'updatedAt'>) {
      // 创建新项目
    },
    
    async updateProject(id: string, project: Partial<Project>) {
      // 更新项目
    },
    
    async deleteProject(id: string) {
      // 软删除项目（移入回收站）
    },
    
    async restoreProject(id: string) {
      // 从回收站恢复项目
    },
    
    selectProject(projectId: string | null) {
      this.selectedProject = projectId
    },
    
    // 导出/导入
    async exportProject(id: string) {
      // 导出项目数据（包含 WebShell 列表）
    },
    
    async importProject(file: File) {
      // 导入项目数据
    },
  },
})
```

### 2.3.3 webshell.ts - WebShell 连接状态

```typescript
import { defineStore } from 'pinia'

interface WebShell {
  id: string
  projectId: string
  name: string
  url: string
  password: string
  payloadType: 'php' | 'jsp' | 'asp' | 'aspx'
  status: 'online' | 'offline' | 'unknown'
  lastConnectedAt?: number
  createdAt: number
  updatedAt: number
}

interface WebShellSession {
  webshellId: string
  isConnected: boolean
  lastActivity: number
  metadata?: Record<string, any>
}

interface WebShellState {
  webshells: WebShell[]
  sessions: Record<string, WebShellSession>
  activeSession: string | null
  isConnecting: boolean
}

export const useWebShellStore = defineStore('webshell', {
  state: (): WebShellState => ({
    webshells: [],
    sessions: {},
    activeSession: null,
    isConnecting: false,
  }),
  
  getters: {
    webshellCount: (state) => state.webshells.length,
    onlineCount: (state) => 
      Object.values(state.sessions).filter(s => s.isConnected).length,
    offlineCount: (state) => 
      Object.values(state.sessions).filter(s => !s.isConnected).length,
    
    getWebShellById: (state) => (id: string) => 
      state.webshells.find(w => w.id === id),
    
    getSessionById: (state) => (id: string) => 
      state.sessions[id],
  },
  
  actions: {
    async loadWebShells(projectId?: string) {
      // 从后端加载 WebShell 列表
    },
    
    async createWebShell(webshell: Omit<WebShell, 'id' | 'createdAt' | 'updatedAt'>) {
      // 创建新的 WebShell
    },
    
    async updateWebShell(id: string, webshell: Partial<WebShell>) {
      // 更新 WebShell
    },
    
    async deleteWebShell(id: string) {
      // 软删除 WebShell
    },
    
    async connectWebShell(id: string) {
      this.isConnecting = true
      try {
        // 连接到 WebShell
        // 建立会话
      } finally {
        this.isConnecting = false
      }
    },
    
    async disconnectWebShell(id: string) {
      // 断开 WebShell 连接
    },
    
    async testConnection(id: string) {
      // 测试 WebShell 连接
    },
    
    openConsole(webshellId: string) {
      // 打开控制台窗口
      this.activeSession = webshellId
    },
    
    closeConsole() {
      this.activeSession = null
    },
  },
})
```

### 2.3.4 plugin.ts - 插件状态

```typescript
import { defineStore } from 'pinia'

interface Plugin {
  id: string
  name: string
  version: string
  description: string
  author: string
  enabled: boolean
  installed: boolean
  type: 'lua' | 'wasm' | 'native'
  icon?: string
}

interface PluginState {
  plugins: Plugin[]
  isLoading: boolean
  isInstalling: boolean
}

export const usePluginStore = defineStore('plugin', {
  state: (): PluginState => ({
    plugins: [],
    isLoading: false,
    isInstalling: false,
  }),
  
  getters: {
    pluginCount: (state) => state.plugins.length,
    enabledCount: (state) => state.plugins.filter(p => p.enabled).length,
    installedPlugins: (state) => state.plugins.filter(p => p.installed),
    availablePlugins: (state) => state.plugins.filter(p => !p.installed),
  },
  
  actions: {
    async loadPlugins() {
      this.isLoading = true
      try {
        // 从后端加载插件列表
      } finally {
        this.isLoading = false
      }
    },
    
    async installPlugin(pluginId: string) {
      this.isInstalling = true
      try {
        // 安装插件
      } finally {
        this.isInstalling = false
      }
    },
    
    async uninstallPlugin(pluginId: string) {
      // 卸载插件
    },
    
    async enablePlugin(pluginId: string) {
      // 启用插件
    },
    
    async disablePlugin(pluginId: string) {
      // 禁用插件
    },
    
    async updatePlugin(pluginId: string) {
      // 更新插件
    },
  },
})
```

### 2.3.5 app.ts - 应用全局状态

```typescript
import { defineStore } from 'pinia'

interface AppSettings {
  theme: 'light' | 'dark' | 'system'
  language: 'zh-CN' | 'en-US'
  accentColor: string
  connectionTimeout: number
  retryCount: number
  proxyEnabled: boolean
  proxyType: 'http' | 'https' | 'socks5'
  proxyHost: string
  proxyPort: number
  encryptionAlgorithm: 'aes-256-gcm' | 'xor'
  dataEncryption: boolean
  logDesensitization: boolean
}

interface AppState {
  settings: AppSettings
  isConnected: boolean
  isLoading: boolean
  notifications: any[]
}

export const useAppStore = defineStore('app', {
  state: (): AppState => ({
    settings: {
      theme: 'system',
      language: 'zh-CN',
      accentColor: '#18a058',
      connectionTimeout: 30,
      retryCount: 3,
      proxyEnabled: false,
      proxyType: 'http',
      proxyHost: '127.0.0.1',
      proxyPort: 7890,
      encryptionAlgorithm: 'aes-256-gcm',
      dataEncryption: true,
      logDesensitization: true,
    },
    isConnected: false,
    isLoading: false,
    notifications: [],
  }),
  
  getters: {
    isDarkTheme: (state) => {
      if (state.settings.theme === 'system') {
        return window.matchMedia('(prefers-color-scheme: dark)').matches
      }
      return state.settings.theme === 'dark'
    },
  },
  
  actions: {
    async loadSettings() {
      // 从本地存储或后端加载设置
    },
    
    async saveSettings(settings: Partial<AppSettings>) {
      // 保存设置到本地存储或后端
      this.settings = { ...this.settings, ...settings }
    },
    
    async updateTheme(theme: 'light' | 'dark' | 'system') {
      this.settings.theme = theme
    },
    
    async updateLanguage(language: 'zh-CN' | 'en-US') {
      this.settings.language = language
    },
    
    async updateAccentColor(color: string) {
      this.settings.accentColor = color
    },
    
    setConnected(connected: boolean) {
      this.isConnected = connected
    },
    
    addNotification(notification: any) {
      this.notifications.push(notification)
    },
    
    removeNotification(notificationId: string) {
      this.notifications = this.notifications.filter(
        n => n.id !== notificationId
      )
    },
  },
})
```

---

## 2.4 工具函数设计

### 2.4.1 export.ts - 导出工具

```typescript
import { save } from '@tauri-apps/api/dialog'
import { writeTextFile } from '@tauri-apps/api/fs'

/**
 * 导出 Payload 配置
 */
export async function exportPayloadConfigs(
  configs: PayloadConfig[],
  format: 'json' | 'csv' = 'json'
) {
  const defaultPath = `payload_configs_${Date.now()}.${format}`
  
  const filePath = await save({
    defaultPath,
    filters: [{
      name: format === 'json' ? 'JSON' : 'CSV',
      extensions: [format]
    }]
  })
  
  if (!filePath) return false
  
  try {
    if (format === 'json') {
      const content = JSON.stringify(configs, null, 2)
      await writeTextFile(filePath, content)
    } else {
      // CSV 格式
      const headers = ['id', 'name', 'url', 'payloadType', 'encryption', 'tags']
      const rows = configs.map(c => headers.map(h => c[h as keyof PayloadConfig] || ''))
      const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
      await writeTextFile(filePath, csv)
    }
    
    return true
  } catch (error) {
    console.error('导出失败:', error)
    return false
  }
}

/**
 * 导入 Payload 配置
 */
export async function importPayloadConfigs(file: File): Promise<PayloadConfig[]> {
  const text = await file.text()
  
  try {
    if (file.name.endsWith('.json')) {
      const configs = JSON.parse(text)
      if (!Array.isArray(configs)) {
        throw new Error('无效的 JSON 格式')
      }
      return configs
    } else if (file.name.endsWith('.csv')) {
      const lines = text.split('\n')
      const headers = lines[0].split(',')
      const configs = lines.slice(1).map(line => {
        const values = line.split(',')
        const config: any = {}
        headers.forEach((h, i) => {
          config[h.trim()] = values[i]?.trim() || ''
        })
        return config as PayloadConfig
      })
      return configs
    } else {
      throw new Error('不支持的文件格式')
    }
  } catch (error) {
    console.error('导入失败:', error)
    throw error
  }
}

/**
 * 导出项目数据
 */
export async function exportProject(
  project: Project,
  webshells: WebShell[]
) {
  const defaultPath = `project_${project.name}_${Date.now()}.json`
  
  const filePath = await save({
    defaultPath,
    filters: [{
      name: 'FG-ABYSS Project',
      extensions: ['json']
    }]
  })
  
  if (!filePath) return false
  
  try {
    const data = {
      version: '1.0',
      exportedAt: new Date().toISOString(),
      project,
      webshells: webshells.map(w => ({
        ...w,
        password: encrypt(w.password)  // 密码加密
      }))
    }
    
    const content = JSON.stringify(data, null, 2)
    await writeTextFile(filePath, content)
    
    return true
  } catch (error) {
    console.error('导出失败:', error)
    return false
  }
}

/**
 * 导出查询结果
 */
export async function exportQueryResult(
  data: any[],
  format: 'csv' | 'json' | 'sql' = 'csv'
) {
  const defaultPath = `query_result_${Date.now()}.${format}`
  
  const filePath = await save({
    defaultPath,
    filters: [{
      name: format.toUpperCase(),
      extensions: [format]
    }]
  })
  
  if (!filePath) return false
  
  try {
    if (format === 'csv') {
      const headers = Object.keys(data[0] || {})
      const rows = data.map(row => 
        headers.map(h => {
          const value = row[h]
          if (typeof value === 'string' && (value.includes(',') || value.includes('"'))) {
            return `"${value.replace(/"/g, '""')}"`
          }
          return value
        }).join(',')
      )
      const csv = [headers.join(','), ...rows].join('\n')
      await writeTextFile(filePath, csv)
    } else if (format === 'json') {
      const content = JSON.stringify(data, null, 2)
      await writeTextFile(filePath, content)
    } else if (format === 'sql') {
      const tableName = 'query_result'
      const inserts = data.map(row => {
        const keys = Object.keys(row).join(', ')
        const values = Object.values(row).map(v => 
          typeof v === 'string' ? `'${v.replace(/'/g, "''")}'` : v
        ).join(', ')
        return `INSERT INTO ${tableName} (${keys}) VALUES (${values});`
      })
      const sql = inserts.join('\n')
      await writeTextFile(filePath, sql)
    }
    
    return true
  } catch (error) {
    console.error('导出失败:', error)
    return false
  }
}
```

### 2.4.2 search.ts - 搜索工具

```typescript
/**
 * 搜索选项
 */
interface SearchOptions {
  caseSensitive?: boolean
  fuzzy?: boolean
  fields?: string[]
}

/**
 * 在数组中搜索
 */
export function searchInArray<T>(
  items: T[],
  query: string,
  options: SearchOptions = {}
): T[] {
  const {
    caseSensitive = false,
    fuzzy = false,
    fields
  } = options
  
  if (!query.trim()) return items
  
  const searchQuery = caseSensitive ? query : query.toLowerCase()
  
  return items.filter(item => {
    const searchFields = fields || Object.keys(item)
    
    return searchFields.some(field => {
      const value = (item as any)[field]
      if (value === null || value === undefined) return false
      
      const searchText = caseSensitive ? String(value) : String(value).toLowerCase()
      
      if (fuzzy) {
        return fuzzyMatch(searchQuery, searchText)
      }
      
      return searchText.includes(searchQuery)
    })
  })
}

/**
 * 模糊匹配算法
 */
function fuzzyMatch(query: string, text: string): boolean {
  let queryIndex = 0
  let textIndex = 0
  
  while (queryIndex < query.length && textIndex < text.length) {
    if (query[queryIndex] === text[textIndex]) {
      queryIndex++
    }
    textIndex++
  }
  
  return queryIndex === query.length
}

/**
 * 高亮搜索匹配
 */
export function highlightMatch(text: string, query: string): string {
  if (!query) return text
  
  const regex = new RegExp(`(${query})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}
```

### 2.4.3 sort.ts - 排序工具

```typescript
/**
 * 多字段排序
 */
export function sortArray<T>(
  items: T[],
  sortConfigs: Array<{
    field: keyof T
    order: 'asc' | 'desc'
  }>
): T[] {
  return [...items].sort((a, b) => {
    for (const { field, order } of sortConfigs) {
      const aVal = a[field]
      const bVal = b[field]
      
      let comparison = 0
      
      if (typeof aVal === 'string' && typeof bVal === 'string') {
        comparison = aVal.localeCompare(bVal)
      } else if (typeof aVal === 'number' && typeof bVal === 'number') {
        comparison = aVal - bVal
      } else if (aVal instanceof Date && bVal instanceof Date) {
        comparison = aVal.getTime() - bVal.getTime()
      }
      
      if (comparison !== 0) {
        return order === 'asc' ? comparison : -comparison
      }
    }
    
    return 0
  })
}

/**
 * 自然排序（适用于文件名等）
 */
export function naturalSort(a: string, b: string): number {
  return a.localeCompare(b, undefined, { numeric: true, sensitivity: 'base' })
```

### 2.3.5 batch.ts - 批量操作状态

```typescript
import { defineStore } from 'pinia'

/**
 * 任务状态
 */
type TaskStatus = 'pending' | 'running' | 'success' | 'failed' | 'cancelled'

/**
 * 任务类型
 */
type TaskType = 
  | 'test_connection'
  | 'execute_command'
  | 'file_operation'
  | 'payload_generate'
  | 'data_export'
  | 'custom'

/**
 * 任务项
 */
interface BatchTask {
  id: string
  type: TaskType
  name: string
  description?: string
  status: TaskStatus
  progress: number           // 0-100
  total: number              // 总任务数
  current: number            // 当前完成数
  result?: any               // 执行结果
  error?: string             // 错误信息
  createdAt: number
  startedAt?: number
  completedAt?: number
  // 任务数据
  data?: {
    webshellIds?: string[]   // 目标 WebShell 列表
    command?: string         // 批量执行的命令
    files?: string[]         // 批量操作的文件
    [key: string]: any
  }
}

/**
 * 任务队列配置
 */
interface QueueConfig {
  maxConcurrent: number      // 最大并发数
  retryCount: number         // 失败重试次数
  retryDelay: number         // 重试延迟（毫秒）
  autoStart: boolean         // 是否自动开始
}

interface BatchState {
  tasks: BatchTask[]
  activeTaskId: string | null
  isRunning: boolean
  queue: BatchTask[]
  config: QueueConfig
  // 统计信息
  stats: {
    total: number
    pending: number
    running: number
    success: number
    failed: number
    cancelled: number
  }
}

export const useBatchStore = defineStore('batch', {
  state: (): BatchState => ({
    tasks: [],
    activeTaskId: null,
    isRunning: false,
    queue: [],
    config: {
      maxConcurrent: 5,
      retryCount: 3,
      retryDelay: 1000,
      autoStart: false,
    },
    stats: {
      total: 0,
      pending: 0,
      running: 0,
      success: 0,
      failed: 0,
      cancelled: 0,
    },
  }),
  
  getters: {
    // 获取活动任务
    activeTask: (state) => 
      state.tasks.find(t => t.id === state.activeTaskId),
    
    // 获取 pending 任务
    pendingTasks: (state) => 
      state.tasks.filter(t => t.status === 'pending'),
    
    // 获取运行中任务
    runningTasks: (state) => 
      state.tasks.filter(t => t.status === 'running'),
    
    // 获取已完成任务
    completedTasks: (state) => 
      state.tasks.filter(t => t.status === 'success' || t.status === 'failed'),
    
    // 获取失败任务
    failedTasks: (state) => 
      state.tasks.filter(t => t.status === 'failed'),
    
    // 总体进度
    overallProgress: (state) => {
      if (state.tasks.length === 0) return 0
      const total = state.tasks.reduce((sum, t) => sum + t.total, 0)
      const current = state.tasks.reduce((sum, t) => sum + t.current, 0)
      return total > 0 ? Math.round((current / total) * 100) : 0
    },
  },
  
  actions: {
    /**
     * 创建批量任务
     */
    createTask(task: Omit<BatchTask, 'id' | 'status' | 'progress' | 'createdAt'>) {
      const newTask: BatchTask = {
        ...task,
        id: `task_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
        status: 'pending',
        progress: 0,
        createdAt: Date.now(),
      }
      
      this.tasks.push(newTask)
      this.updateStats()
      
      if (this.config.autoStart) {
        this.startTask(newTask.id)
      }
      
      return newTask.id
    },
    
    /**
     * 批量连接测试
     */
    async batchTestConnection(webshellIds: string[]) {
      const taskId = this.createTask({
        type: 'test_connection',
        name: '批量连接测试',
        description: `测试 ${webshellIds.length} 个 WebShell 的连接`,
        total: webshellIds.length,
        current: 0,
        data: { webshellIds },
      })
      
      await this.executeTask(taskId, async (task) => {
        const results = []
        for (const id of task.data!.webshellIds!) {
          try {
            // 调用后端 API 测试连接
            const result = await this.testConnectionAPI(id)
            results.push({ id, success: true, result })
            task.current++
            task.progress = Math.round((task.current / task.total) * 100)
          } catch (error) {
            results.push({ 
              id, 
              success: false, 
              error: error instanceof Error ? error.message : '未知错误' 
            })
            task.current++
            task.progress = Math.round((task.current / task.total) * 100)
          }
        }
        return results
      })
    },
    
    /**
     * 批量执行命令
     */
    async batchExecuteCommand(webshellIds: string[], command: string) {
      const taskId = this.createTask({
        type: 'execute_command',
        name: '批量命令执行',
        description: `在 ${webshellIds.length} 个 WebShell 上执行命令`,
        total: webshellIds.length,
        current: 0,
        data: { webshellIds, command },
      })
      
      await this.executeTask(taskId, async (task) => {
        const results = []
        for (const id of task.data!.webshellIds!) {
          try {
            const result = await this.executeCommandAPI(id, task.data!.command!)
            results.push({ id, success: true, result })
            task.current++
            task.progress = Math.round((task.current / task.total) * 100)
          } catch (error) {
            results.push({ id, success: false, error: String(error) })
            task.current++
            task.progress = Math.round((task.current / task.total) * 100)
          }
        }
        return results
      })
    },
    
    /**
     * 批量文件操作
     */
    async batchFileOperation(
      webshellIds: string[], 
      operation: 'upload' | 'download' | 'delete',
      files: string[]
    ) {
      const taskId = this.createTask({
        type: 'file_operation',
        name: `批量文件${operation === 'upload' ? '上传' : operation === 'download' ? '下载' : '删除'}`,
        total: webshellIds.length * files.length,
        current: 0,
        data: { webshellIds, files, operation },
      })
      
      await this.executeTask(taskId, async (task) => {
        const results = []
        for (const webshellId of task.data!.webshellIds!) {
          for (const file of task.data!.files!) {
            try {
              await this.fileOperationAPI(webshellId, operation, file)
              results.push({ webshellId, file, success: true })
              task.current++
              task.progress = Math.round((task.current / task.total) * 100)
            } catch (error) {
              results.push({ webshellId, file, success: false, error: String(error) })
              task.current++
              task.progress = Math.round((task.current / task.total) * 100)
            }
          }
        }
        return results
      })
    },
    
    /**
     * 批量生成 Payload
     */
    async batchGeneratePayload(configIds: string[]) {
      const taskId = this.createTask({
        type: 'payload_generate',
        name: '批量生成 Payload',
        total: configIds.length,
        current: 0,
        data: { configIds },
      })
      
      await this.executeTask(taskId, async (task) => {
        const results = []
        for (const configId of task.data!.configIds!) {
          try {
            const payload = await this.generatePayloadAPI(configId)
            results.push({ configId, success: true, payload })
            task.current++
            task.progress = Math.round((task.current / task.total) * 100)
          } catch (error) {
            results.push({ configId, success: false, error: String(error) })
            task.current++
            task.progress = Math.round((task.current / task.total) * 100)
          }
        }
        return results
      })
    },
    
    /**
     * 执行任务（核心逻辑）
     */
    async executeTask(
      taskId: string, 
      executor: (task: BatchTask) => Promise<any>
    ) {
      const task = this.tasks.find(t => t.id === taskId)
      if (!task) return
      
      try {
        task.status = 'running'
        task.startedAt = Date.now()
        this.activeTaskId = taskId
        this.isRunning = true
        
        // 执行任务
        const result = await executor(task)
        
        // 执行成功
        task.status = 'success'
        task.result = result
        task.completedAt = Date.now()
        
      } catch (error) {
        // 执行失败
        task.status = 'failed'
        task.error = error instanceof Error ? error.message : '未知错误'
        task.completedAt = Date.now()
        
      } finally {
        this.activeTaskId = null
        this.isRunning = false
        this.updateStats()
      }
    },
    
    /**
     * 开始任务
     */
    async startTask(taskId: string) {
      const task = this.tasks.find(t => t.id === taskId)
      if (!task || task.status !== 'pending') return
      
      await this.executeTask(taskId, async (t) => {
        // 根据任务类型调用不同的执行函数
        switch (t.type) {
          case 'test_connection':
            return this.batchTestConnection(t.data!.webshellIds!)
          case 'execute_command':
            return this.batchExecuteCommand(t.data!.webshellIds!, t.data!.command!)
          case 'file_operation':
            return this.batchFileOperation(
              t.data!.webshellIds!, 
              t.data!.operation!, 
              t.data!.files!
            )
          case 'payload_generate':
            return this.batchGeneratePayload(t.data!.configIds!)
          default:
            throw new Error(`未知任务类型：${t.type}`)
        }
      })
    },
    
    /**
     * 取消任务
     */
    cancelTask(taskId: string) {
      const task = this.tasks.find(t => t.id === taskId)
      if (!task || task.status !== 'running') return
      
      task.status = 'cancelled'
      task.completedAt = Date.now()
      this.activeTaskId = null
      this.isRunning = false
      this.updateStats()
    },
    
    /**
     * 重试失败任务
     */
    async retryTask(taskId: string) {
      const task = this.tasks.find(t => t.id === taskId)
      if (!task || task.status !== 'failed') return
      
      // 重置任务状态
      task.status = 'pending'
      task.progress = 0
      task.current = 0
      task.error = undefined
      task.completedAt = undefined
      
      // 重新开始
      await this.startTask(taskId)
    },
    
    /**
     * 清空任务列表
     */
    clearTasks(status?: TaskStatus) {
      if (status) {
        this.tasks = this.tasks.filter(t => t.status !== status)
      } else {
        this.tasks = []
      }
      this.updateStats()
    },
    
    /**
     * 删除任务
     */
    deleteTask(taskId: string) {
      this.tasks = this.tasks.filter(t => t.id !== taskId)
      this.updateStats()
    },
    
    /**
     * 更新统计信息
     */
    updateStats() {
      this.stats.total = this.tasks.length
      this.stats.pending = this.tasks.filter(t => t.status === 'pending').length
      this.stats.running = this.tasks.filter(t => t.status === 'running').length
      this.stats.success = this.tasks.filter(t => t.status === 'success').length
      this.stats.failed = this.tasks.filter(t => t.status === 'failed').length
      this.stats.cancelled = this.tasks.filter(t => t.status === 'cancelled').length
    },
    
    /**
     * 更新配置
     */
    updateConfig(config: Partial<QueueConfig>) {
      this.config = { ...this.config, ...config }
    },
    
    // API 调用占位函数（实际由后端实现）
    async testConnectionAPI(id: string) {
      // 调用后端 API
    },
    
    async executeCommandAPI(id: string, command: string) {
      // 调用后端 API
    },
    
    async fileOperationAPI(
      id: string, 
      operation: 'upload' | 'download' | 'delete', 
      file: string
    ) {
      // 调用后端 API
    },
    
    async generatePayloadAPI(configId: string) {
      // 调用后端 API
    },
  },
})
```

---

## 2.5 组件设计补充

### 2.5.1 批量操作组件

```vue
<!-- components/batch/BatchOperationPanel.vue -->
<template>
  <div class="batch-operation-panel">
    <!-- 任务列表 -->
    <div class="task-list">
      <div v-for="task in tasks" :key="task.id" class="task-item">
        <div class="task-header">
          <span class="task-name">{{ task.name }}</span>
          <span class="task-status" :class="task.status">
            {{ getStatusText(task.status) }}
          </span>
        </div>
        
        <!-- 进度条 -->
        <n-progress 
          type="line"
          :percentage="task.progress"
          :status="getProgressStatus(task.status)"
        />
        
        <!-- 详细信息 -->
        <div class="task-details">
          <span>总数：{{ task.total }}</span>
          <span>完成：{{ task.current }}</span>
          <span>失败：{{ task.total - task.current }}</span>
        </div>
        
        <!-- 操作按钮 -->
        <div class="task-actions">
          <n-button 
            v-if="task.status === 'pending'"
            @click="startTask(task.id)"
          >
            开始
          </n-button>
          <n-button 
            v-if="task.status === 'running'"
            @click="cancelTask(task.id)"
          >
            取消
          </n-button>
          <n-button 
            v-if="task.status === 'failed'"
            @click="retryTask(task.id)"
          >
            重试
          </n-button>
          <n-button 
            @click="deleteTask(task.id)"
          >
            删除
          </n-button>
        </div>
      </div>
    </div>
    
    <!-- 总体统计 -->
    <div class="overall-stats">
      <n-statistic label="总任务数" :value="stats.total" />
      <n-statistic label="进行中" :value="stats.running" />
      <n-statistic label="成功" :value="stats.success" />
      <n-statistic label="失败" :value="stats.failed" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useBatchStore } from '@/stores/batch'

const batchStore = useBatchStore()

const tasks = computed(() => batchStore.tasks)
const stats = computed(() => batchStore.stats)

const startTask = (id: string) => batchStore.startTask(id)
const cancelTask = (id: string) => batchStore.cancelTask(id)
const retryTask = (id: string) => batchStore.retryTask(id)
const deleteTask = (id: string) => batchStore.deleteTask(id)

const getStatusText = (status: string) => {
  const map = {
    pending: '等待中',
    running: '进行中',
    success: '成功',
    failed: '失败',
    cancelled: '已取消',
  }
  return map[status as keyof typeof map]
}

const getProgressStatus = (status: string) => {
  const map = {
    pending: 'default',
    running: 'default',
    success: 'success',
    failed: 'error',
    cancelled: 'warning',
  }
  return map[status as keyof typeof map]
}
</script>

<style scoped>
.batch-operation-panel {
  padding: 16px;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-item {
  padding: 16px;
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
}

.task-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}

.task-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.task-status.pending { background: #f0f0f0; }
.task-status.running { background: #1890ff; color: white; }
.task-status.success { background: #52c41a; color: white; }
.task-status.failed { background: #ff4d4f; color: white; }
.task-status.cancelled { background: #faad14; color: white; }

.task-details {
  display: flex;
  gap: 16px;
  margin: 12px 0;
  font-size: 12px;
  color: var(--n-text-color-3);
}

.task-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.overall-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--n-border-color);
}
</style>
```

### 2.5.2 进度显示组件

```vue
<!-- components/common/ProgressPanel.vue -->
<template>
  <div class="progress-panel">
    <div class="progress-header">
      <h3>{{ title }}</h3>
      <n-button text @click="onClose">
        <template #icon>
          <i class="i-close"></i>
        </template>
      </n-button>
    </div>
    
    <div class="progress-content">
      <!-- 总体进度 -->
      <div class="overall-progress">
        <n-progress 
          type="line" 
          :percentage="overallProgress"
          :status="progressStatus"
        />
        <div class="progress-text">
          {{ currentCount }} / {{ totalCount }} ({{ overallProgress }}%)
        </div>
      </div>
      
      <!-- 详细列表 -->
      <div class="progress-list">
        <div 
          v-for="item in items" 
          :key="item.id"
          class="progress-item"
          :class="item.status"
        >
          <div class="item-icon">
            <i v-if="item.status === 'success'" class="i-check-circle" />
            <i v-else-if="item.status === 'failed'" class="i-x-circle" />
            <i v-else class="i-loader" />
          </div>
          <div class="item-content">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-status">{{ item.statusText }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface ProgressItem {
  id: string
  name: string
  status: 'pending' | 'running' | 'success' | 'failed'
  statusText: string
}

const props = defineProps<{
  title: string
  items: ProgressItem[]
  overallProgress: number
  onClose: () => void
}>()

const currentCount = computed(() => 
  props.items.filter(i => i.status === 'success' || i.status === 'failed').length
)

const totalCount = computed(() => props.items.length)

const progressStatus = computed(() => {
  const failed = props.items.filter(i => i.status === 'failed').length
  return failed > 0 ? 'error' : 'success'
})
</script>

<style scoped>
.progress-panel {
  position: fixed;
  top: 60px;
  right: 20px;
  width: 400px;
  max-height: 600px;
  background: var(--n-color);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--n-border-color);
}

.progress-content {
  padding: 16px;
  max-height: 500px;
  overflow-y: auto;
}

.overall-progress {
  margin-bottom: 16px;
}

.progress-text {
  text-align: center;
  margin-top: 8px;
  font-size: 12px;
  color: var(--n-text-color-3);
}

.progress-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.progress-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 4px;
  background: var(--n-color-2);
}

.item-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.item-icon i {
  font-size: 16px;
}

.item-content {
  flex: 1;
}

.item-name {
  font-size: 14px;
  font-weight: 500;
}

.item-status {
  font-size: 12px;
  color: var(--n-text-color-3);
}
</style>
```

***

## 2.6 后端模块详细设计

### 2.6.1 混淆引擎 UI 设计

```vue
<!-- components/payload/ObfuscationPanel.vue -->
<template>
  <div class="obfuscation-panel">
    <div class="obfuscation-header">
      <h3>代码混淆引擎</h3>
      <n-tag type="info">Tree-sitter 驱动</n-tag>
    </div>
    
    <!-- 混淆级别选择 -->
    <div class="obfuscation-level">
      <h4>混淆级别</h4>
      <n-slider
        v-model:value="level"
        :min="0"
        :max="5"
        :step="1"
        :marks="levelMarks"
        :tooltip="false"
        @update:value="onLevelChange"
      />
      <div class="level-description">
        {{ levelDescriptions[level] }}
      </div>
    </div>
    
    <!-- 混淆选项 -->
    <div class="obfuscation-options">
      <h4>混淆选项</h4>
      
      <n-checkbox
        v-model:checked="options.variableRenaming"
        :disabled="level < 1"
      >
        变量重命名
      </n-checkbox>
      
      <n-checkbox
        v-model:checked="options.controlFlowFlattening"
        :disabled="level < 2"
      >
        控制流平坦化
      </n-checkbox>
      
      <n-checkbox
        v-model:checked="options.stringEncryption"
        :disabled="level < 3"
      >
        字符串加密
      </n-checkbox>
      
      <n-checkbox
        v-model:checked="options.deadCodeInjection"
        :disabled="level < 4"
      >
        垃圾代码插入
      </n-checkbox>
      
      <n-checkbox
        v-model:checked="options.combinedObfuscation"
        :disabled="level < 5"
      >
        组合混淆（所有技术）
      </n-checkbox>
    </div>
    
    <!-- 预览区域 -->
    <div class="preview-section">
      <h4>混淆预览</h4>
      <div class="preview-tabs">
        <n-tabs type="line">
          <n-tab-pane name="original" tab="原始代码">
            <CodeEditor
              v-model="originalCode"
              :readonly="true"
              language="php"
              height="300px"
            />
          </n-tab-pane>
          <n-tab-pane name="obfuscated" tab="混淆后代码">
            <CodeEditor
              v-model="obfuscatedCode"
              :readonly="true"
              language="php"
              height="300px"
            />
          </n-tab-pane>
          <n-tab-pane name="diff" tab="差异对比">
            <DiffViewer
              :original="originalCode"
              :modified="obfuscatedCode"
            />
          </n-tab-pane>
        </n-tabs>
      </div>
    </div>
    
    <!-- 统计信息 -->
    <div class="obfuscation-stats">
      <n-statistic label="原始行数" :value="originalLines" />
      <n-statistic label="混淆后行数" :value="obfuscatedLines" />
      <n-statistic label="代码膨胀率" :value="expansionRate" suffix="%" />
      <n-statistic label="混淆强度" :value="strengthScore" suffix="/100" />
    </div>
    
    <!-- 操作按钮 -->
    <div class="obfuscation-actions">
      <n-button @click="previewObfuscation">
        <template #icon>
          <i class="i-eye"></i>
        </template>
        预览混淆
      </n-button>
      
      <n-button @click="exportObfuscated" type="primary">
        <template #icon>
          <i class="i-download"></i>
        </template>
        导出混淆代码
      </n-button>
      
      <n-button @click="applyAndGenerate" type="success">
        <template #icon>
          <i class="i-play"></i>
        </template>
        应用并生成 Payload
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import CodeEditor from '@/components/common/CodeEditor.vue'
import DiffViewer from '@/components/common/DiffViewer.vue'

const level = ref(3)
const originalCode = ref(`<?php
password = $_POST['pass'];
exec = $_POST['exec'];
if (password == 'admin') {
    system(exec);
}
?>`)

const obfuscatedCode = ref('')

const levelMarks = {
  0: '无混淆',
  1: '基础',
  2: '中等',
  3: '高级',
  4: '专家',
  5: '地狱',
}

const levelDescriptions = {
  0: '无混淆（原始代码）',
  1: '变量重命名 - 将变量名替换为无意义的字符',
  2: '控制流平坦化 - 打乱代码执行顺序',
  3: '字符串加密 - 加密所有字符串常量',
  4: '垃圾代码插入 - 插入无用代码混淆视听',
  5: '组合混淆 - 使用所有混淆技术',
}

const options = ref({
  variableRenaming: true,
  controlFlowFlattening: false,
  stringEncryption: false,
  deadCodeInjection: false,
  combinedObfuscation: false,
})

const originalLines = computed(() => 
  originalCode.value.split('\n').length
)

const obfuscatedLines = computed(() => 
  obfuscatedCode.value ? obfuscatedCode.value.split('\n').length : 0
)

const expansionRate = computed(() => {
  if (!obfuscatedCode.value) return 0
  return Math.round(
    ((obfuscatedLines.value - originalLines.value) / originalLines.value) * 100
  )
})

const strengthScore = computed(() => {
  const baseScore = level.value * 15
  const optionBonus = Object.values(options.value).filter(Boolean).length * 5
  return Math.min(baseScore + optionBonus, 100)
})

const onLevelChange = (newLevel: number) => {
  // 根据级别自动启用对应选项
  if (newLevel >= 1) options.value.variableRenaming = true
  if (newLevel >= 2) options.value.controlFlowFlattening = true
  if (newLevel >= 3) options.value.stringEncryption = true
  if (newLevel >= 4) options.value.deadCodeInjection = true
  if (newLevel >= 5) options.value.combinedObfuscation = true
}

const previewObfuscation = async () => {
  // 调用后端 API 预览混淆
  obfuscatedCode.value = await obfuscateAPI(
    originalCode.value,
    level.value,
    options.value
  )
}

const exportObfuscated = () => {
  // 导出混淆后的代码
  downloadText(obfuscatedCode.value, 'obfuscated_payload.php')
}

const applyAndGenerate = () => {
  // 应用混淆设置并生成 Payload
  emit('generate', {
    code: obfuscatedCode.value,
    level: level.value,
    options: options.value,
  })
}

async function obfuscateAPI(
  code: string, 
  level: number, 
  options: any
): Promise<string> {
  // 调用后端混淆 API
  return ''
}

function downloadText(content: string, filename: string) {
  // 下载文本
}

const emit = defineEmits(['generate'])
</script>

<style scoped>
.obfuscation-panel {
  padding: 20px;
}

.obfuscation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.obfuscation-level {
  margin-bottom: 24px;
}

.level-description {
  margin-top: 12px;
  padding: 12px;
  background: var(--n-color-2);
  border-radius: 4px;
  font-size: 13px;
  line-height: 1.6;
}

.obfuscation-options {
  margin-bottom: 24px;
}

.obfuscation-options h4 {
  margin-bottom: 12px;
}

.obfuscation-options .n-checkbox {
  display: block;
  margin-bottom: 8px;
}

.preview-section {
  margin-bottom: 24px;
}

.preview-tabs {
  margin-top: 12px;
}

.obfuscation-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: var(--n-color-2);
  border-radius: 8px;
}

.obfuscation-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}
</style>
```

### 2.6.2 插件 API 设计

#### A. Lua 插件 API

```rust
use mlua::{Lua, LuaSerdeExt, Table, Value};
use serde::{Serialize, Deserialize};

/// Lua 插件上下文
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LuaPluginContext {
    pub plugin_id: String,
    pub plugin_name: String,
    pub plugin_version: String,
    pub webshell_id: Option<String>,
}

/// Lua 插件接口
pub struct LuaPluginAPI {
    lua: Lua,
    context: LuaPluginContext,
}

impl LuaPluginAPI {
    pub fn new(context: LuaPluginContext) -> Result<Self, Box<dyn std::error::Error>> {
        let lua = Lua::new();
        
        // 注册全局函数
        let globals = lua.globals();
        
        // http_get(url, headers) -> {status, body, error}
        globals.set("http_get", lua.create_function(|_, (url, headers): (String, Option<Table>)| {
            // 调用后端 HTTP 客户端
            let client = reqwest::blocking::Client::new();
            let mut request = client.get(&url);
            
            if let Some(h) = headers {
                for pair in h.pairs::<String, String>() {
                    let (key, value) = pair?;
                    request = request.header(&key, &value);
                }
            }
            
            let response = request.send()?;
            let status = response.status().as_u16();
            let body = response.text()?;
            
            Ok(lua.create_table_from([
                ("status", status),
                ("body", body),
                ("error", Value::Nil),
            ])?)
        })?)?;
        
        // http_post(url, data, headers) -> {status, body, error}
        globals.set("http_post", lua.create_function(|_, (url, data, headers): (String, Table, Option<Table>)| {
            // 实现 POST 请求
            Ok(lua.create_table()?)
        })?)?;
        
        // shell_exec(webshell_id, command) -> {output, error}
        globals.set("shell_exec", lua.create_function(|lua, (webshell_id, command): (String, String)| {
            // 调用后端执行命令
            let output = format!("Output of {} on {}", command, webshell_id);
            
            Ok(lua.create_table_from([
                ("output", output),
                ("error", Value::Nil),
            ])?)
        })?)?;
        
        // log_info(message)
        globals.set("log_info", lua.create_function(|_, message: String| {
            log::info!("[Plugin] {}", message);
            Ok(())
        })?)?;
        
        // log_error(message)
        globals.set("log_error", lua.create_function(|_, message: String| {
            log::error!("[Plugin] {}", message);
            Ok(())
        })?)?;
        
        // fs_list(path) -> {files: [], dirs: [], error: nil}
        globals.set("fs_list", lua.create_function(|lua, path: String| {
            // 列出目录内容
            let files = lua.create_table();
            let dirs = lua.create_table();
            
            Ok(lua.create_table_from([
                ("files", files),
                ("dirs", dirs),
                ("error", Value::Nil),
            ])?)
        })?)?;
        
        // fs_read(path) -> {content, error}
        globals.set("fs_read", lua.create_function(|lua, path: String| {
            // 读取文件
            Ok(lua.create_table_from([
                ("content", ""),
                ("error", Value::Nil),
            ])?)
        })?)?;
        
        // fs_write(path, content) -> {success, error}
        globals.set("fs_write", lua.create_function(|_, (path, content): (String, String)| {
            // 写入文件
            Ok(())
        })?)?;
        
        // db_connect(connection_string) -> {connection_id, error}
        globals.set("db_connect", lua.create_function(|lua, conn_str: String| {
            // 连接数据库
            Ok(lua.create_table_from([
                ("connection_id", "conn_123"),
                ("error", Value::Nil),
            ])?)
        })?)?;
        
        // db_query(connection_id, sql) -> {rows, columns, error}
        globals.set("db_query", lua.create_function(|lua, (conn_id, sql): (String, String)| {
            // 执行 SQL 查询
            let columns = lua.create_table();
            let rows = lua.create_table();
            
            Ok(lua.create_table_from([
                ("columns", columns),
                ("rows", rows),
                ("error", Value::Nil),
            ])?)
        })?)?;
        
        Ok(Self { lua, context })
    }
    
    /// 加载插件脚本
    pub fn load_script(&self, script: &str) -> Result<(), Box<dyn std::error::Error>> {
        self.lua.load(script).exec()?;
        Ok(())
    }
    
    /// 调用插件函数
    pub fn call_function<T: serde::de::DeserializeOwned>(
        &self,
        func_name: &str,
        args: Vec<Value>
    ) -> Result<T, Box<dyn std::error::Error>> {
        let func: mlua::Function = self.lua.globals().get(func_name)?;
        let result = func.call(args)?;
        Ok(self.lua.from_value(result)?)
    }
    
    /// 执行插件
    pub fn execute(&self, function_name: &str, params: serde_json::Value) -> Result<serde_json::Value, Box<dyn std::error::Error>> {
        let lua_params = self.lua.to_value(&params)?;
        let result: mlua::Value = self.call_function(function_name, vec![lua_params])?;
        Ok(self.lua.from_value(result)?)
    }
}

/// Lua 插件元数据
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LuaPluginMetadata {
    pub name: String,
    pub version: String,
    pub author: String,
    pub description: String,
    pub entry_point: String,  // 入口函数名
    pub required_api_version: String,
}
```

#### B. WebAssembly 插件 API

```rust
use wasmtime::{Engine, Store, Module, Instance, Func, Val, ValType};
use serde::{Serialize, Deserialize};

/// WASM 插件上下文
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WasmPluginContext {
    pub plugin_id: String,
    pub memory_limit: u32,  // 内存限制（MB）
    pub timeout: u32,       // 超时时间（ms）
}

/// WASM 插件运行时
pub struct WasmPluginRuntime {
    engine: Engine,
    store: Store<WasmPluginContext>,
    instance: Option<Instance>,
}

impl WasmPluginRuntime {
    pub fn new(context: WasmPluginContext) -> Result<Self, Box<dyn std::error::Error>> {
        let engine = Engine::default();
        let mut store = Store::new(&engine, context);
        
        // 定义主机函数（Host Functions）
        let http_get = Func::wrap(&mut store, |ptr: i32, len: i32| -> i64 {
            // ptr 和 len 指向 WASM 内存中的 URL
            // 返回 (result_ptr, result_len)
            0
        });
        
        let shell_exec = Func::wrap(&mut store, |webshell_id: i64, cmd_ptr: i32, cmd_len: i32| -> i64 {
            // 执行命令
            0
        });
        
        let log_info = Func::wrap(&mut store, |msg_ptr: i32, msg_len: i32| {
            // 记录日志
        });
        
        let malloc = Func::wrap(&mut store, |size: i32| -> i32 {
            // 内存分配
            0
        });
        
        let free = Func::wrap(&mut store, |ptr: i32| {
            // 内存释放
        });
        
        Ok(Self {
            engine,
            store,
            instance: None,
        })
    }
    
    /// 加载 WASM 模块
    pub fn load_module(&mut self, wasm_bytes: &[u8]) -> Result<(), Box<dyn std::error::Error>> {
        let module = Module::new(self.engine(), wasm_bytes)?;
        let instance = Instance::new(&mut self.store, &module, &[
            http_get.into(),
            shell_exec.into(),
            log_info.into(),
            malloc.into(),
            free.into(),
        ])?;
        
        self.instance = Some(instance);
        Ok(())
    }
    
    /// 调用插件函数
    pub fn call(&mut self, func_name: &str, args: &[Val]) -> Result<Vec<Val>, Box<dyn std::error::Error>> {
        let instance = self.instance.as_ref().ok_or("实例未加载")?;
        let func = instance.get_func(&mut self.store, func_name)
            .ok_or_else(|| format!("函数不存在：{}", func_name))?;
        
        let results = func.call(&mut self.store, args)?;
        Ok(results)
    }
    
    /// 执行插件（JSON 接口）
    pub fn execute_json(&mut self, input: &str) -> Result<String, Box<dyn std::error::Error>> {
        // 1. 分配内存存储输入
        // 2. 调用 plugin_execute(input_ptr, input_len)
        // 3. 读取返回结果
        // 4. 释放内存
        
        Ok("{\"success\": true, \"data\": {}}".to_string())
    }
    
    pub fn engine(&self) -> &Engine {
        &self.engine
    }
}

/// WASM 插件元数据
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WasmPluginMetadata {
    pub name: String,
    pub version: String,
    pub author: String,
    pub description: String,
    pub exports: Vec<String>,  // 导出的函数列表
}
```

#### C. Rust 原生插件 API

```rust
use std::sync::Arc;
use serde::{Serialize, Deserialize};

/// 插件 trait（所有插件必须实现）
pub trait NativePlugin: Send + Sync {
    /// 获取插件元数据
    fn metadata(&self) -> PluginMetadata;
    
    /// 初始化插件
    fn initialize(&mut self, context: PluginContext) -> Result<(), PluginError>;
    
    /// 执行插件功能
    fn execute(&self, params: serde_json::Value) -> Result<serde_json::Value, PluginError>;
    
    /// 关闭插件
    fn shutdown(&self) -> Result<(), PluginError>;
}

/// 插件元数据
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PluginMetadata {
    pub name: String,
    pub version: String,
    pub author: String,
    pub description: String,
    pub plugin_type: PluginType,
}

/// 插件类型
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum PluginType {
    FileManager,      // 文件管理
    DatabaseManager,  // 数据库管理
    Terminal,         // 终端
    Custom,           // 自定义
}

/// 插件上下文
#[derive(Debug, Clone)]
pub struct PluginContext {
    pub plugin_id: String,
    pub webshell_id: Option<String>,
    pub http_client: Arc<dyn HttpClient>,
    pub logger: Arc<dyn Logger>,
}

/// HTTP 客户端 trait
pub trait HttpClient: Send + Sync {
    fn get(&self, url: &str, headers: std::collections::HashMap<String, String>) 
        -> Result<HttpResponse, HttpError>;
    
    fn post(&self, url: &str, data: &[u8], headers: std::collections::HashMap<String, String>) 
        -> Result<HttpResponse, HttpError>;
}

/// HTTP 响应
#[derive(Debug, Clone)]
pub struct HttpResponse {
    pub status: u16,
    pub headers: std::collections::HashMap<String, String>,
    pub body: Vec<u8>,
}

/// 日志 trait
pub trait Logger: Send + Sync {
    fn info(&self, message: &str);
    fn error(&self, message: &str);
    fn debug(&self, message: &str);
}

/// 插件错误
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PluginError {
    pub code: String,
    pub message: String,
}

impl std::fmt::Display for PluginError {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "[{}] {}", self.code, self.message)
    }
}

impl std::error::Error for PluginError {}

/// HTTP 错误
#[derive(Debug, Clone)]
pub struct HttpError {
    pub code: String,
    pub message: String,
}

// ============ 示例：文件管理插件实现 ============

pub struct FileManagerPlugin {
    context: Option<PluginContext>,
}

impl NativePlugin for FileManagerPlugin {
    fn metadata(&self) -> PluginMetadata {
        PluginMetadata {
            name: "File Manager".to_string(),
            version: "1.0.0".to_string(),
            author: "FG-ABYSS Team".to_string(),
            description: "文件管理插件".to_string(),
            plugin_type: PluginType::FileManager,
        }
    }
    
    fn initialize(&mut self, context: PluginContext) -> Result<(), PluginError> {
        self.context = Some(context);
        Ok(())
    }
    
    fn execute(&self, params: serde_json::Value) -> Result<serde_json::Value, PluginError> {
        let action = params["action"].as_str().ok_or(PluginError {
            code: "INVALID_PARAMS".to_string(),
            message: "缺少 action 参数".to_string(),
        })?;
        
        match action {
            "list" => self.list_files(&params),
            "upload" => self.upload_file(&params),
            "download" => self.download_file(&params),
            "delete" => self.delete_file(&params),
            "edit" => self.edit_file(&params),
            _ => Err(PluginError {
                code: "UNKNOWN_ACTION".to_string(),
                message: format!("未知操作：{}", action),
            }),
        }
    }
    
    fn shutdown(&self) -> Result<(), PluginError> {
        Ok(())
    }
}

impl FileManagerPlugin {
    fn list_files(&self, params: &serde_json::Value) -> Result<serde_json::Value, PluginError> {
        let path = params["path"].as_str().unwrap_or("/");
        // 实现文件列表逻辑
        Ok(serde_json::json!({
            "success": true,
            "data": {
                "files": [],
                "dirs": []
            }
        }))
    }
    
    fn upload_file(&self, params: &serde_json::Value) -> Result<serde_json::Value, PluginError> {
        // 实现文件上传逻辑
        Ok(serde_json::json!({
            "success": true,
            "data": {}
        }))
    }
    
    fn download_file(&self, params: &serde_json::Value) -> Result<serde_json::Value, PluginError> {
        // 实现文件下载逻辑
        Ok(serde_json::json!({
            "success": true,
            "data": {}
        }))
    }
    
    fn delete_file(&self, params: &serde_json::Value) -> Result<serde_json::Value, PluginError> {
        // 实现文件删除逻辑
        Ok(serde_json::json!({
            "success": true,
            "data": {}
        }))
    }
    
    fn edit_file(&self, params: &serde_json::Value) -> Result<serde_json::Value, PluginError> {
        // 实现文件编辑逻辑
        Ok(serde_json::json!({
            "success": true,
            "data": {}
        }))
    }
}

// ============ 插件加载器 ============

pub struct PluginLoader {
    plugins: std::collections::HashMap<String, Box<dyn NativePlugin>>,
}

impl PluginLoader {
    pub fn new() -> Self {
        Self {
            plugins: std::collections::HashMap::new(),
        }
    }
    
    /// 注册插件
    pub fn register(&mut self, plugin: Box<dyn NativePlugin>) {
        let metadata = plugin.metadata();
        self.plugins.insert(metadata.name.clone(), plugin);
    }
    
    /// 获取插件
    pub fn get_plugin(&self, name: &str) -> Option<&dyn NativePlugin> {
        self.plugins.get(name).map(|p| p.as_ref())
    }
    
    /// 执行插件
    pub fn execute_plugin(
        &self,
        name: &str,
        params: serde_json::Value
    ) -> Result<serde_json::Value, PluginError> {
        let plugin = self.get_plugin(name).ok_or_else(|| PluginError {
            code: "PLUGIN_NOT_FOUND".to_string(),
            message: format!("插件不存在：{}", name),
        })?;
        
        plugin.execute(params)
    }
    
    /// 列出所有插件
    pub fn list_plugins(&self) -> Vec<PluginMetadata> {
        self.plugins.values()
            .map(|p| p.metadata())
            .collect()
    }
}
```

### 2.6.4 备份恢复系统设计

```rust
use std::fs;
use std::path::{Path, PathBuf};
use std::io::{Read, Write};
use flate2::read::GzDecoder;
use flate2::write::GzEncoder;
use flate2::Compression;
use tar::{Archive, Builder};
use chrono::{DateTime, Utc, Duration};
use serde::{Serialize, Deserialize};
use rusqlite::Connection;

/// 备份类型
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum BackupType {
    Full,           // 完整备份
    Incremental,    // 增量备份
    Differential,   // 差异备份
}

/// 备份状态
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum BackupStatus {
    Pending,
    Running,
    Success,
    Failed,
    Restoring,
}

/// 备份元数据
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BackupMetadata {
    pub backup_id: String,
    pub backup_type: BackupType,
    pub status: BackupStatus,
    pub created_at: DateTime<Utc>,
    pub size: u64,
    pub file_count: u32,
    pub description: String,
    pub backup_path: PathBuf,
    pub checksum: String,  // SHA256 校验和
    pub is_encrypted: bool,
    pub fg_abyss_version: String,
}

/// 备份配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BackupConfig {
    pub backup_dir: PathBuf,
    pub auto_backup_enabled: bool,
    pub auto_backup_interval_hours: u32,
    pub max_backups: u32,
    pub backup_on_exit: bool,
    pub encryption_enabled: bool,
    pub encryption_key: Option<String>,
    pub compression_level: u32,  // 0-9
}

/// 备份项
#[derive(Debug, Clone)]
pub struct BackupItem {
    pub path: PathBuf,
    pub relative_path: String,
    pub size: u64,
    pub modified_at: DateTime<Utc>,
    pub item_type: BackupItemType,
}

/// 备份项类型
#[derive(Debug, Clone)]
pub enum BackupItemType {
    Database,
    Config,
    Logs,
    Payloads,
    Projects,
    Custom,
}

/// 备份管理器
pub struct BackupManager {
    config: BackupConfig,
    last_backup_time: Option<DateTime<Utc>>,
    last_backup_hash: Option<String>,
}

impl BackupManager {
    /// 创建新的备份管理器
    pub fn new(config: BackupConfig) -> Result<Self, Box<dyn std::error::Error>> {
        // 确保备份目录存在
        fs::create_dir_all(&config.backup_dir)?;
        
        Ok(Self {
            config,
            last_backup_time: None,
            last_backup_hash: None,
        })
    }
    
    /// 创建完整备份
    pub fn create_full_backup(&mut self, description: String) -> Result<BackupMetadata, Box<dyn std::error::Error>> {
        log::info!("开始创建完整备份...");
        
        let backup_id = format!("backup_{}", Utc::now().format("%Y%m%d_%H%M%S"));
        let backup_path = self.config.backup_dir.join(&backup_id).with_extension("tar.gz");
        
        // 收集所有备份项
        let items = self.collect_backup_items()?;
        
        // 创建临时目录
        let temp_dir = self.config.backup_dir.join(&backup_id);
        fs::create_dir_all(&temp_dir)?;
        
        // 复制文件到临时目录
        for item in &items {
            let dest = temp_dir.join(&item.relative_path);
            if let Some(parent) = dest.parent() {
                fs::create_dir_all(parent)?;
            }
            fs::copy(&item.path, &dest)?;
        }
        
        // 创建元数据文件
        let metadata = BackupMetadata {
            backup_id: backup_id.clone(),
            backup_type: BackupType::Full,
            status: BackupStatus::Running,
            created_at: Utc::now(),
            size: 0,
            file_count: items.len() as u32,
            description,
            backup_path: backup_path.clone(),
            checksum: String::new(),
            is_encrypted: self.config.encryption_enabled,
            fg_abyss_version: env!("CARGO_PKG_VERSION").to_string(),
        };
        
        let metadata_path = temp_dir.join("metadata.json");
        fs::write(&metadata_path, serde_json::to_string_pretty(&metadata)?)?;
        
        // 压缩备份
        let tar_gz = fs::File::create(&backup_path)?;
        let enc = GzEncoder::new(tar_gz, Compression::new(self.config.compression_level));
        let mut tar = Builder::new(enc);
        
        tar.append_dir_all(&backup_id, &temp_dir)?;
        
        // 完成压缩
        let mut enc = tar.into_inner()?.into_inner()?;
        enc.finish()?;
        
        // 计算校验和
        let mut file = fs::File::open(&backup_path)?;
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer)?;
        let checksum = format!("{:x}", sha2::Sha256::digest(&buffer));
        
        // 更新元数据
        let size = fs::metadata(&backup_path)?.len();
        let mut metadata = metadata;
        metadata.size = size;
        metadata.checksum = checksum;
        metadata.status = BackupStatus::Success;
        
        // 保存元数据到备份文件
        if self.config.encryption_enabled {
            // TODO: 加密备份
            log::info!("备份已加密");
        }
        
        // 清理临时目录
        fs::remove_dir_all(&temp_dir)?;
        
        self.last_backup_time = Some(metadata.created_at);
        self.last_backup_hash = Some(metadata.checksum.clone());
        
        log::info!("完整备份创建完成：{:?}, 大小：{} bytes", backup_path, size);
        
        // 清理旧备份
        self.cleanup_old_backups()?;
        
        Ok(metadata)
    }
    
    /// 创建增量备份
    pub fn create_incremental_backup(&mut self, description: String) -> Result<BackupMetadata, Box<dyn std::error::Error>> {
        if self.last_backup_hash.is_none() {
            return self.create_full_backup(description);
        }
        
        log::info!("开始创建增量备份...");
        
        let backup_id = format!("backup_inc_{}", Utc::now().format("%Y%m%d_%H%M%S"));
        let backup_path = self.config.backup_dir.join(&backup_id).with_extension("tar.gz");
        
        // 只备份自上次备份以来更改的文件
        let changed_items = self.collect_changed_items()?;
        
        if changed_items.is_empty() {
            log::info!("没有更改的文件，跳过增量备份");
            return Err("没有需要备份的更改".into());
        }
        
        // 创建临时目录
        let temp_dir = self.config.backup_dir.join(&backup_id);
        fs::create_dir_all(&temp_dir)?;
        
        // 复制更改的文件
        for item in &changed_items {
            let dest = temp_dir.join(&item.relative_path);
            if let Some(parent) = dest.parent() {
                fs::create_dir_all(parent)?;
            }
            fs::copy(&item.path, &dest)?;
        }
        
        // 保存上次备份的哈希用于对比
        let last_hash_path = temp_dir.join("last_backup.hash");
        fs::write(&last_hash_path, self.last_backup_hash.as_ref().unwrap())?;
        
        // 创建元数据
        let metadata = BackupMetadata {
            backup_id: backup_id.clone(),
            backup_type: BackupType::Incremental,
            status: BackupStatus::Running,
            created_at: Utc::now(),
            size: 0,
            file_count: changed_items.len() as u32,
            description,
            backup_path: backup_path.clone(),
            checksum: String::new(),
            is_encrypted: self.config.encryption_enabled,
            fg_abyss_version: env!("CARGO_PKG_VERSION").to_string(),
        };
        
        // 压缩和保存（与完整备份类似）
        // ...（省略压缩代码，与完整备份相同）
        
        Ok(metadata)
    }
    
    /// 恢复备份
    pub fn restore_backup(&self, backup_path: &Path) -> Result<(), Box<dyn std::error::Error>> {
        log::info!("开始恢复备份：{:?}", backup_path);
        
        // 验证备份文件
        if !backup_path.exists() {
            return Err("备份文件不存在".into());
        }
        
        // 解压备份
        let temp_dir = self.config.backup_dir.join("restoring");
        fs::create_dir_all(&temp_dir)?;
        
        let tar_gz = fs::File::open(backup_path)?;
        let dec = GzDecoder::new(tar_gz);
        let mut archive = Archive::new(dec);
        
        archive.unpack(&temp_dir)?;
        
        // 读取元数据
        let metadata_path = temp_dir.join("metadata.json");
        let metadata_content = fs::read_to_string(&metadata_path)?;
        let metadata: BackupMetadata = serde_json::from_str(&metadata_content)?;
        
        // 验证校验和
        if !self.verify_checksum(backup_path, &metadata.checksum)? {
            return Err("备份文件校验和验证失败".into());
        }
        
        // 恢复文件
        for entry in fs::read_dir(&temp_dir)? {
            let entry = entry?;
            let path = entry.path();
            
            if path.is_file() && path.file_name().unwrap() != "metadata.json" {
                let relative = path.strip_prefix(&temp_dir)?;
                let dest = self.get_restore_destination(relative)?;
                
                if let Some(parent) = dest.parent() {
                    fs::create_dir_all(parent)?;
                }
                
                fs::copy(&path, &dest)?;
            }
        }
        
        // 清理临时目录
        fs::remove_dir_all(&temp_dir)?;
        
        log::info!("备份恢复完成");
        
        Ok(())
    }
    
    /// 列出所有备份
    pub fn list_backups(&self) -> Result<Vec<BackupMetadata>, Box<dyn std::error::Error>> {
        let mut backups = Vec::new();
        
        for entry in fs::read_dir(&self.config.backup_dir)? {
            let entry = entry?;
            let path = entry.path();
            
            if path.extension().and_then(|s| s.to_str()) == Some("tar.gz") {
                // 从备份文件中提取元数据
                if let Ok(metadata) = self.extract_metadata(&path) {
                    backups.push(metadata);
                }
            }
        }
        
        // 按时间倒序排序
        backups.sort_by(|a, b| b.created_at.cmp(&a.created_at));
        
        Ok(backups)
    }
    
    /// 删除备份
    pub fn delete_backup(&self, backup_id: &str) -> Result<(), Box<dyn std::error::Error>> {
        let backups = self.list_backups()?;
        
        if let Some(backup) = backups.iter().find(|b| b.backup_id == backup_id) {
            fs::remove_file(&backup.backup_path)?;
            log::info!("备份已删除：{}", backup_id);
        } else {
            return Err("备份不存在".into());
        }
        
        Ok(())
    }
    
    /// 验证备份
    pub fn verify_backup(&self, backup_path: &Path) -> Result<bool, Box<dyn std::error::Error>> {
        let metadata = self.extract_metadata(backup_path)?;
        
        // 验证校验和
        if !self.verify_checksum(backup_path, &metadata.checksum)? {
            return Ok(false);
        }
        
        // 验证版本兼容性
        let current_version = env!("CARGO_PKG_VERSION");
        if !self.is_version_compatible(&metadata.fg_abyss_version, current_version)? {
            log::warn!("备份版本可能不兼容：{} vs {}", metadata.fg_abyss_version, current_version);
        }
        
        Ok(true)
    }
    
    /// 收集所有备份项
    fn collect_backup_items(&self) -> Result<Vec<BackupItem>, Box<dyn std::error::Error>> {
        let mut items = Vec::new();
        
        // 数据库文件
        if let Some(db_path) = self.get_database_path() {
            items.push(BackupItem {
                path: db_path,
                relative_path: "data/database.sqlite".to_string(),
                size: 0,
                modified_at: Utc::now(),
                item_type: BackupItemType::Database,
            });
        }
        
        // 配置文件
        if let Some(config_path) = self.get_config_path() {
            items.push(BackupItem {
                path: config_path,
                relative_path: "config/config.toml".to_string(),
                size: 0,
                modified_at: Utc::now(),
                item_type: BackupItemType::Config,
            });
        }
        
        // 日志文件（可选）
        // items.push(...);
        
        Ok(items)
    }
    
    /// 收集更改的项
    fn collect_changed_items(&self) -> Result<Vec<BackupItem>, Box<dyn std::error::Error>> {
        let all_items = self.collect_backup_items()?;
        
        // 过滤出上次备份后更改的项
        let changed = all_items.into_iter().filter(|item| {
            if let Some(last_backup) = self.last_backup_time {
                item.modified_at > last_backup
            } else {
                true
            }
        }).collect();
        
        Ok(changed)
    }
    
    /// 清理旧备份
    fn cleanup_old_backups(&self) -> Result<(), Box<dyn std::error::Error>> {
        let backups = self.list_backups()?;
        
        if backups.len() as u32 > self.config.max_backups {
            let to_delete = backups.iter().skip(self.config.max_backups as usize);
            
            for backup in to_delete {
                self.delete_backup(&backup.backup_id)?;
            }
            
            log::info!("已清理 {} 个旧备份", backups.len() - self.config.max_backups as usize);
        }
        
        Ok(())
    }
    
    /// 提取备份元数据
    fn extract_metadata(&self, backup_path: &Path) -> Result<BackupMetadata, Box<dyn std::error::Error>> {
        // 从 tar.gz 文件中读取 metadata.json
        let tar_gz = fs::File::open(backup_path)?;
        let dec = GzDecoder::new(tar_gz);
        let mut archive = Archive::new(dec);
        
        for entry in archive.entries()? {
            let mut entry = entry?;
            let path = entry.path()?;
            
            if path.file_name().and_then(|s| s.to_str()) == Some("metadata.json") {
                let mut content = String::new();
                entry.read_to_string(&mut content)?;
                return Ok(serde_json::from_str(&content)?);
            }
        }
        
        Err("备份中未找到元数据文件".into())
    }
    
    /// 验证校验和
    fn verify_checksum(&self, backup_path: &Path, expected: &str) -> Result<bool, Box<dyn std::error::Error>> {
        let mut file = fs::File::open(backup_path)?;
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer)?;
        
        let checksum = format!("{:x}", sha2::Sha256::digest(&buffer));
        
        Ok(checksum == expected)
    }
    
    /// 检查版本兼容性
    fn is_version_compatible(&self, backup_version: &str, current_version: &str) -> Result<bool, Box<dyn std::error::Error>> {
        // 简单的版本兼容性检查
        // 实际实现中应该更复杂
        Ok(true)
    }
    
    /// 获取恢复目标路径
    fn get_restore_destination(&self, relative_path: &Path) -> Result<PathBuf, Box<dyn std::error::Error>> {
        // 根据相对路径确定恢复位置
        if relative_path.starts_with("data") {
            Ok(self.get_database_path().unwrap_or_else(|| PathBuf::from("data/database.sqlite")))
        } else if relative_path.starts_with("config") {
            Ok(self.get_config_path().unwrap_or_else(|| PathBuf::from("config/config.toml")))
        } else {
            Err("未知的恢复路径".into())
        }
    }
    
    fn get_database_path(&self) -> Option<PathBuf> {
        // 获取数据库路径
        None
    }
    
    fn get_config_path(&self) -> Option<PathBuf> {
        // 获取配置路径
        None
    }
}

// 自动备份调度器
pub struct AutoBackupScheduler {
    backup_manager: BackupManager,
    interval: Duration,
    running: bool,
}

impl AutoBackupScheduler {
    pub fn new(backup_manager: BackupManager, interval_hours: u32) -> Self {
        Self {
            backup_manager,
            interval: Duration::hours(interval_hours as i64),
            running: false,
        }
    }
    
    pub fn start(&mut self) {
        self.running = true;
        
        std::thread::spawn(move || {
            while self.running {
                std::thread::sleep(self.interval);
                
                match self.backup_manager.create_incremental_backup("自动备份".to_string()) {
                    Ok(metadata) => {
                        log::info!("自动备份完成：{}", metadata.backup_id);
                    }
                    Err(e) => {
                        log::error!("自动备份失败：{}", e);
                    }
                }
            }
        });
    }
    
    pub fn stop(&mut self) {
        self.running = false;
    }
}
```

### 2.6.5 事件总线系统设计

```rust
use std::collections::HashMap;
use std::sync::{Arc, RwLock};
use serde::{Serialize, Deserialize};
use chrono::{DateTime, Utc};
use tokio::sync::broadcast;
use uuid::Uuid;

/// 事件 trait（所有事件必须实现）
pub trait Event: Send + Sync + Serialize + Clone {
    /// 事件类型名称
    fn event_type(&self) -> &'static str;
    
    /// 事件发生时间
    fn timestamp(&self) -> DateTime<Utc>;
    
    /// 事件来源模块
    fn source(&self) -> Option<&str>;
}

/// 事件处理器 trait
pub trait EventHandler: Send + Sync {
    /// 处理事件
    fn handle(&self, event: Arc<dyn Event>) -> Result<(), EventError>;
    
    /// 过滤器（可选）
    fn filter(&self, _event_type: &str) -> bool {
        true
    }
}

/// 事件错误
#[derive(Debug, Clone)]
pub struct EventError {
    pub code: String,
    pub message: String,
}

impl std::fmt::Display for EventError {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "[{}] {}", self.code, self.message)
    }
}

impl std::error::Error for EventError {}

/// 事件总线配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EventBusConfig {
    pub channel_size: usize,           // 广播频道大小
    pub enable_persistence: bool,      // 是否启用持久化
    pub persistence_interval: u64,     // 持久化间隔（秒）
    pub max_stored_events: u32,        // 最大存储事件数
    pub event_types: Vec<String>,      // 启用的事件类型
}

/// 事件订阅者
#[derive(Debug, Clone)]
pub struct EventSubscriber {
    pub subscriber_id: String,
    pub event_types: Vec<String>,
    pub handler: Arc<dyn EventHandler>,
    pub priority: u32,  // 优先级（数字越小优先级越高）
}

/// 事件总线核心结构
pub struct EventBus {
    config: EventBusConfig,
    broadcaster: broadcast::Sender<Arc<dyn Event>>,
    subscribers: Arc<RwLock<HashMap<String, EventSubscriber>>>,
    event_store: Arc<RwLock<Vec<StoredEvent>>>,
    running: Arc<RwLock<bool>>,
}

impl EventBus {
    /// 创建新的事件总线
    pub fn new(config: EventBusConfig) -> Self {
        let (broadcaster, _) = broadcast::channel(config.channel_size);
        
        Self {
            config,
            broadcaster,
            subscribers: Arc::new(RwLock::new(HashMap::new())),
            event_store: Arc::new(RwLock::new(Vec::new())),
            running: Arc::new(RwLock::new(false)),
        }
    }
    
    /// 注册事件处理器
    pub fn subscribe(&self, subscriber: EventSubscriber) -> Result<String, EventError> {
        let mut subscribers = self.subscribers.write().unwrap();
        subscribers.insert(subscriber.subscriber_id.clone(), subscriber);
        Ok(subscriber.subscriber_id.clone())
    }
    
    /// 取消订阅
    pub fn unsubscribe(&self, subscriber_id: &str) -> Result<(), EventError> {
        let mut subscribers = self.subscribers.write().unwrap();
        subscribers.remove(subscriber_id)
            .ok_or_else(|| EventError {
                code: "SUBSCRIBER_NOT_FOUND".to_string(),
                message: format!("订阅者不存在：{}", subscriber_id),
            })?;
        Ok(())
    }
    
    /// 发布事件
    pub fn publish(&self, event: Arc<dyn Event>) -> Result<(), EventError> {
        // 1. 存储事件（如果启用持久化）
        if self.config.enable_persistence {
            self.store_event(Arc::clone(&event));
        }
        
        // 2. 广播给所有订阅者
        let subscribers = self.subscribers.read().unwrap();
        let mut sorted_subscribers: Vec<_> = subscribers.values().collect();
        
        // 按优先级排序
        sorted_subscribers.sort_by_key(|s| s.priority);
        
        // 3. 异步分发事件
        for subscriber in sorted_subscribers {
            if subscriber.filter(event.event_type()) {
                let handler = Arc::clone(&subscriber.handler);
                let event_clone = Arc::clone(&event);
                
                tokio::spawn(async move {
                    if let Err(e) = handler.handle(event_clone) {
                        log::error!("事件处理失败：{}", e);
                    }
                });
            }
        }
        
        // 4. 广播到频道
        let _ = self.broadcaster.send(event);
        
        Ok(())
    }
    
    /// 存储事件
    fn store_event(&self, event: Arc<dyn Event>) {
        let mut store = self.event_store.write().unwrap();
        
        let stored_event = StoredEvent {
            event_id: Uuid::new_v4().to_string(),
            event_type: event.event_type().to_string(),
            timestamp: event.timestamp(),
            source: event.source().map(|s| s.to_string()),
            data: serde_json::to_value(Arc::clone(&event)).unwrap_or_default(),
        };
        
        store.push(stored_event);
        
        // 限制存储数量
        if store.len() > self.config.max_stored_events as usize {
            store.remove(0);
        }
    }
    
    /// 查询历史事件
    pub fn query_events(
        &self,
        query: EventQuery
    ) -> Result<Vec<StoredEvent>, EventError> {
        let store = self.event_store.read().unwrap();
        
        let mut results: Vec<_> = store.iter()
            .filter(|event| {
                // 类型过滤
                if let Some(ref event_type) = query.event_type {
                    if &event.event_type != event_type {
                        return false;
                    }
                }
                
                // 时间范围过滤
                if let Some(start) = query.start_time {
                    if event.timestamp < start {
                        return false;
                    }
                }
                if let Some(end) = query.end_time {
                    if event.timestamp > end {
                        return false;
                    }
                }
                
                // 来源过滤
                if let Some(ref source) = query.source {
                    if let Some(ref event_source) = event.source {
                        if event_source != source {
                            return false;
                        }
                    } else {
                        return false;
                    }
                }
                
                true
            })
            .cloned()
            .collect();
        
        // 排序（默认时间倒序）
        results.sort_by(|a, b| {
            if query.ascending {
                a.timestamp.cmp(&b.timestamp)
            } else {
                b.timestamp.cmp(&a.timestamp)
            }
        });
        
        // 分页
        let results = results
            .into_iter()
            .skip(query.offset as usize)
            .take(query.limit as usize)
            .collect();
        
        Ok(results)
    }
    
    /// 清空事件存储
    pub fn clear_events(&self) -> Result<(), EventError> {
        let mut store = self.event_store.write().unwrap();
        store.clear();
        Ok(())
    }
    
    /// 导出事件
    pub fn export_events(
        &self,
        query: EventQuery,
        format: &str
    ) -> Result<String, EventError> {
        let events = self.query_events(query)?;
        
        match format {
            "json" => Ok(serde_json::to_string_pretty(&events).unwrap_or_default()),
            "csv" => {
                let mut csv = String::from("event_id,event_type,timestamp,source\n");
                for event in &events {
                    csv.push_str(&format!(
                        "{},{},{},{}\n",
                        event.event_id,
                        event.event_type,
                        event.timestamp.to_rfc3339(),
                        event.source.as_deref().unwrap_or("")
                    ));
                }
                Ok(csv)
            },
            _ => Err(EventError {
                code: "UNSUPPORTED_FORMAT".to_string(),
                message: format!("不支持的导出格式：{}", format),
            }),
        }
    }
    
    /// 启动事件持久化任务
    pub fn start_persistence_task(&self) -> Result<(), EventError> {
        if !self.config.enable_persistence {
            return Err(EventError {
                code: "PERSISTENCE_DISABLED".to_string(),
                message: "持久化未启用".to_string(),
            });
        }
        
        let event_store = Arc::clone(&self.event_store);
        let interval = self.config.persistence_interval;
        let running = Arc::clone(&self.running);
        
        *running.write().unwrap() = true;
        
        tokio::spawn(async move {
            while *running.read().unwrap() {
                tokio::time::sleep(tokio::time::Duration::from_secs(interval)).await;
                
                // 持久化事件到磁盘
                let store = event_store.read().unwrap();
                if !store.is_empty() {
                    // TODO: 实际持久化逻辑
                    log::info!("持久化 {} 个事件", store.len());
                }
            }
        });
        
        Ok(())
    }
    
    /// 停止事件总线
    pub fn stop(&self) {
        *self.running.write().unwrap() = false;
    }
    
    /// 获取订阅者数量
    pub fn subscriber_count(&self) -> usize {
        self.subscribers.read().unwrap().len()
    }
    
    /// 获取存储的事件数量
    pub fn stored_event_count(&self) -> usize {
        self.event_store.read().unwrap().len()
    }
}

/// 存储的事件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StoredEvent {
    pub event_id: String,
    pub event_type: String,
    pub timestamp: DateTime<Utc>,
    pub source: Option<String>,
    pub data: serde_json::Value,
}

/// 事件查询条件
#[derive(Debug, Clone)]
pub struct EventQuery {
    pub event_type: Option<String>,
    pub start_time: Option<DateTime<Utc>>,
    pub end_time: Option<DateTime<Utc>>,
    pub source: Option<String>,
    pub limit: u32,
    pub offset: u32,
    pub ascending: bool,
}

// ============ 具体事件类型实现 ============

/// Shell 连接事件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ShellConnectedEvent {
    pub shell_id: String,
    pub shell_name: String,
    pub connected_at: DateTime<Utc>,
    pub response_time_ms: u64,
}

impl Event for ShellConnectedEvent {
    fn event_type(&self) -> &'static str {
        "SHELL_CONNECTED"
    }
    
    fn timestamp(&self) -> DateTime<Utc> {
        self.connected_at
    }
    
    fn source(&self) -> Option<&str> {
        Some("shell_manager")
    }
}

/// Shell 断开连接事件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ShellDisconnectedEvent {
    pub shell_id: String,
    pub shell_name: String,
    pub disconnected_at: DateTime<Utc>,
    pub reason: String,
}

impl Event for ShellDisconnectedEvent {
    fn event_type(&self) -> &'static str {
        "SHELL_DISCONNECTED"
    }
    
    fn timestamp(&self) -> DateTime<Utc> {
        self.disconnected_at
    }
    
    fn source(&self) -> Option<&str> {
        Some("shell_manager")
    }
}

/// 文件操作事件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FileOperationEvent {
    pub operation: String,  // upload, download, delete, edit
    pub path: String,
    pub shell_id: String,
    pub timestamp: DateTime<Utc>,
    pub success: bool,
    pub error: Option<String>,
}

impl Event for FileOperationEvent {
    fn event_type(&self) -> &'static str {
        "FILE_OPERATION"
    }
    
    fn timestamp(&self) -> DateTime<Utc> {
        self.timestamp
    }
    
    fn source(&self) -> Option<&str> {
        Some("file_manager")
    }
}

/// 数据库操作事件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DatabaseOperationEvent {
    pub operation: String,  // query, export, backup
    pub connection_id: String,
    pub sql: Option<String>,
    pub timestamp: DateTime<Utc>,
    pub success: bool,
    pub row_count: Option<u32>,
    pub error: Option<String>,
}

impl Event for DatabaseOperationEvent {
    fn event_type(&self) -> &'static str {
        "DATABASE_OPERATION"
    }
    
    fn timestamp(&self) -> DateTime<Utc> {
        self.timestamp
    }
    
    fn source(&self) -> Option<&str> {
        Some("database_manager")
    }
}

/// 批量操作事件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BatchOperationEvent {
    pub batch_id: String,
    pub operation_type: String,  // delete, export, test_connection
    pub total_count: u32,
    pub success_count: u32,
    pub failed_count: u32,
    pub started_at: DateTime<Utc>,
    pub completed_at: Option<DateTime<Utc>>,
    pub error: Option<String>,
}

impl Event for BatchOperationEvent {
    fn event_type(&self) -> &'static str {
        "BATCH_OPERATION"
    }
    
    fn timestamp(&self) -> DateTime<Utc> {
        self.started_at
    }
    
    fn source(&self) -> Option<&str> {
        Some("batch_manager")
    }
}

/// 系统事件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SystemEvent {
    pub event_kind: String,  // startup, shutdown, error, config_change
    pub message: String,
    pub timestamp: DateTime<Utc>,
    pub severity: String,  // info, warning, error, critical
}

impl Event for SystemEvent {
    fn event_type(&self) -> &'static str {
        "SYSTEM_EVENT"
    }
    
    fn timestamp(&self) -> DateTime<Utc> {
        self.timestamp
    }
    
    fn source(&self) -> Option<&str> {
        Some("system")
    }
}

// ============ 事件处理器示例 ============

/// 日志事件处理器
pub struct LoggingEventHandler {
    logger: Arc<dyn crate::logging::Logger>,
}

impl LoggingEventHandler {
    pub fn new(logger: Arc<dyn crate::logging::Logger>) -> Self {
        Self { logger }
    }
}

impl EventHandler for LoggingEventHandler {
    fn handle(&self, event: Arc<dyn Event>) -> Result<(), EventError> {
        let event_type = event.event_type();
        let serialized = serde_json::to_string(&event).unwrap_or_default();
        
        self.logger.info(&format!("事件：{} - {}", event_type, serialized));
        
        Ok(())
    }
}

/// 统计事件处理器
pub struct StatisticsEventHandler {
    stats: Arc<RwLock<EventStatistics>>,
}

impl StatisticsEventHandler {
    pub fn new() -> Self {
        Self {
            stats: Arc::new(RwLock::new(EventStatistics::default())),
        }
    }
    
    pub fn get_statistics(&self) -> EventStatistics {
        self.stats.read().unwrap().clone()
    }
}

impl EventHandler for StatisticsEventHandler {
    fn handle(&self, event: Arc<dyn Event>) -> Result<(), EventError> {
        let mut stats = self.stats.write().unwrap();
        
        stats.total_events += 1;
        stats
            .events_by_type
            .entry(event.event_type().to_string())
            .or_insert(0) += 1;
        
        if let Some(source) = event.source() {
            stats
                .events_by_source
                .entry(source.to_string())
                .or_insert(0) += 1;
        }
        
        stats.last_event_time = event.timestamp();
        
        Ok(())
    }
}

/// 事件统计
#[derive(Debug, Clone, Default)]
pub struct EventStatistics {
    pub total_events: u64,
    pub events_by_type: HashMap<String, u64>,
    pub events_by_source: HashMap<String, u64>,
    pub last_event_time: DateTime<Utc>,
}

/// 告警事件处理器
pub struct AlertEventHandler {
    alert_threshold: u32,
    alert_callback: Box<dyn Fn(&str) + Send + Sync>,
}

impl AlertEventHandler {
    pub fn new<F>(alert_threshold: u32, callback: F) -> Self
    where
        F: Fn(&str) + Send + Sync + 'static,
    {
        Self {
            alert_threshold,
            alert_callback: Box::new(callback),
        }
    }
}

impl EventHandler for AlertEventHandler {
    fn handle(&self, event: Arc<dyn Event>) -> Result<(), EventError> {
        // 对错误事件进行告警
        if event.event_type() == "SYSTEM_EVENT" {
            if let Ok(system_event) = Arc::downcast::<SystemEvent>(event) {
                if system_event.severity == "error" || system_event.severity == "critical" {
                    (self.alert_callback)(&system_event.message);
                }
            }
        }
        
        Ok(())
    }
}

// ============ 使用示例 ============

pub fn init_event_bus() -> Result<Arc<EventBus>, EventError> {
    let config = EventBusConfig {
        channel_size: 1000,
        enable_persistence: true,
        persistence_interval: 60,  // 每分钟持久化一次
        max_stored_events: 10000,
        event_types: vec![
            "SHELL_CONNECTED".to_string(),
            "SHELL_DISCONNECTED".to_string(),
            "FILE_OPERATION".to_string(),
            "DATABASE_OPERATION".to_string(),
            "BATCH_OPERATION".to_string(),
            "SYSTEM_EVENT".to_string(),
        ],
    };
    
    let event_bus = Arc::new(EventBus::new(config));
    
    // 注册日志处理器
    let logger = Arc::new(crate::logging::FgAbyssLogger::default());
    let logging_handler = LoggingEventHandler::new(logger);
    
    event_bus.subscribe(EventSubscriber {
        subscriber_id: "logging_handler".to_string(),
        event_types: vec!["*".to_string()],  // 订阅所有事件
        handler: Arc::new(logging_handler),
        priority: 1,  // 高优先级
    })?;
    
    // 注册统计处理器
    let stats_handler = StatisticsEventHandler::new();
    event_bus.subscribe(EventSubscriber {
        subscriber_id: "statistics_handler".to_string(),
        event_types: vec!["*".to_string()],
        handler: Arc::new(stats_handler),
        priority: 2,
    })?;
    
    // 注册告警处理器
    let alert_handler = AlertEventHandler::new(5, |message| {
        log::warn!("告警：{}", message);
        // 这里可以发送通知到外部系统
    });
    
    event_bus.subscribe(EventSubscriber {
        subscriber_id: "alert_handler".to_string(),
        event_types: vec!["SYSTEM_EVENT".to_string()],
        handler: Arc::new(alert_handler),
        priority: 0,  // 最高优先级
    })?;
    
    // 启动持久化任务
    event_bus.start_persistence_task()?;
    
    log::info!("事件总线初始化完成");
    
    Ok(event_bus)
}
```

### 2.6.6 日志系统设计
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum LogLevel {
    Trace,
    Debug,
    Info,
    Warn,
    Error,
}

/// 日志配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogConfig {
    pub level: LogLevel,
    pub max_file_size: u64,           // 单个日志文件最大大小（字节）
    pub max_files: u32,               // 保留的日志文件数量
    pub log_dir: PathBuf,             // 日志目录
    pub enable_console: bool,         // 是否输出到控制台
    pub enable_file: bool,            // 是否输出到文件
    pub enable_json: bool,            // 是否使用 JSON 格式
}

/// 结构化日志条目
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogEntry {
    pub timestamp: DateTime<Utc>,
    pub level: String,
    pub target: String,
    pub message: String,
    pub module: Option<String>,
    pub file: Option<String>,
    pub line: Option<u32>,
    pub context: Option<serde_json::Value>,
}

/// 日志记录器
pub struct FgAbyssLogger {
    config: LogConfig,
    current_file: PathBuf,
    current_size: u64,
    file_index: u32,
}

impl FgAbyssLogger {
    /// 创建新的日志记录器
    pub fn new(config: LogConfig) -> Result<Self, Box<dyn std::error::Error>> {
        // 创建日志目录
        fs::create_dir_all(&config.log_dir)?;
        
        let current_file = config.log_dir.join("fg-abyss.log");
        let current_size = if current_file.exists() {
            fs::metadata(&current_file)?.len()
        } else {
            0
        };
        
        Ok(Self {
            config,
            current_file,
            current_size,
            file_index: 0,
        })
    }
    
    /// 初始化全局日志记录器
    pub fn init(config: LogConfig) -> Result<(), Box<dyn std::error::Error>> {
        let logger = Self::new(config)?;
        
        log::set_boxed_logger(Box::new(logger.clone()))?;
        log::set_max_level(LevelFilter::Info);
        
        Ok(())
    }
    
    /// 检查是否需要轮转日志
    fn should_rotate(&self) -> bool {
        self.current_size >= self.config.max_file_size
    }
    
    /// 轮转日志文件
    fn rotate_logs(&mut self) -> Result<(), Box<dyn std::error::Error>> {
        // 删除最旧的日志文件
        let oldest_file = self.config.log_dir.join(format!(
            "fg-abyss.log.{}",
            self.config.max_files - 1
        ));
        if oldest_file.exists() {
            fs::remove_file(oldest_file)?;
        }
        
        // 重命名现有文件
        for i in (1..self.config.max_files).rev() {
            let old_file = self.config.log_dir.join(format!("fg-abyss.log.{}", i - 1));
            let new_file = self.config.log_dir.join(format!("fg-abyss.log.{}", i));
            if old_file.exists() {
                fs::rename(old_file, new_file)?;
            }
        }
        
        // 重命名当前文件
        if self.current_file.exists() {
            let new_file = self.config.log_dir.join("fg-abyss.log.0");
            fs::rename(&self.current_file, new_file)?;
        }
        
        // 重置当前文件
        self.current_file = self.config.log_dir.join("fg-abyss.log");
        self.current_size = 0;
        
        Ok(())
    }
    
    /// 格式化日志消息
    fn format_log(&self, record: &Record, context: Option<serde_json::Value>) -> String {
        let timestamp = Utc::now().format("%Y-%m-%d %H:%M:%S%.3f");
        
        if self.config.enable_json {
            // JSON 格式
            let entry = LogEntry {
                timestamp: Utc::now(),
                level: record.level().to_string(),
                target: record.target().to_string(),
                message: record.args().to_string(),
                module: record.module_path().map(|s| s.to_string()),
                file: record.file().map(|s| s.to_string()),
                line: record.line(),
                context,
            };
            serde_json::to_string(&entry).unwrap_or_default()
        } else {
            // 文本格式
            let module = record.module_path().unwrap_or("");
            format!(
                "{} [{}] {} - {}",
                timestamp,
                record.level(),
                module,
                record.args()
            )
        }
    }
    
    /// 敏感信息脱敏
    fn desensitize(message: &str) -> String {
        let mut result = message.to_string();
        
        // URL 中的密码脱敏
        let url_pattern = regex::Regex::new(r"(https?://[^:]+:)([^@]+)(@)").unwrap();
        result = url_pattern.replace_all(&result, "$1******$3").to_string();
        
        // 密码字段脱敏
        let password_pattern = regex::Regex::new(r#"(?"password"\s*:\s*")([^"]+)("")"#).unwrap();
        result = password_pattern.replace_all(&result, r#"$1******$2"#).to_string();
        
        // API 密钥脱敏
        let api_key_pattern = regex::Regex::new(r#"(?"api_key"\s*:\s*")([^"]+)("")"#).unwrap();
        result = api_key_pattern.replace_all(&result, r#"$1******$2"#).to_string();
        
        // IP 地址脱敏（保留前两段）
        let ip_pattern = regex::Regex::new(r"(\d{1,3}\.\d{1,3}\.)\d{1,3}\.\d{1,3}").unwrap();
        result = ip_pattern.replace_all(&result, "${1}***.***").to_string();
        
        result
    }
}

impl log::Log for FgAbyssLogger {
    fn enabled(&self, metadata: &Metadata) -> bool {
        metadata.level() <= self.config.level.into()
    }
    
    fn log(&self, record: &Record) {
        if !self.enabled(record.metadata()) {
            return;
        }
        
        // 敏感信息脱敏
        let message = Self::desensitize(&record.args().to_string());
        
        // 控制台输出
        if self.config.enable_console {
            let log_msg = self.format_log(record, None);
            println!("{}", log_msg);
        }
        
        // 文件输出
        if self.config.enable_file {
            let mut file = OpenOptions::new()
                .create(true)
                .append(true)
                .open(&self.current_file)
                .unwrap();
            
            let mut writer = BufWriter::new(file);
            let log_msg = self.format_log(record, None);
            writeln!(writer, "{}", log_msg).unwrap();
            writer.flush().unwrap();
            
            // 检查是否需要轮转
            if self.should_rotate() {
                // 这里需要 Mutex 保护，实际实现中需要使用 Arc<Mutex<>>
                // 简化示例
            }
        }
    }
    
    fn flush(&self) {
        // 刷新缓冲区
    }
}

/// 日志查询条件
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogQuery {
    pub level: Option<LogLevel>,
    pub start_time: Option<DateTime<Utc>>,
    pub end_time: Option<DateTime<Utc>>,
    pub keyword: Option<String>,
    pub module: Option<String>,
    pub limit: u32,
    pub offset: u32,
}

/// 日志查询结果
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogQueryResult {
    pub total: u32,
    pub logs: Vec<LogEntry>,
}

/// 日志分析器
pub struct LogAnalyzer {
    log_dir: PathBuf,
}

impl LogAnalyzer {
    pub fn new(log_dir: PathBuf) -> Self {
        Self { log_dir }
    }
    
    /// 查询日志
    pub fn query_logs(&self, query: LogQuery) -> Result<LogQueryResult, Box<dyn std::error::Error>> {
        let mut all_logs = Vec::new();
        
        // 读取所有日志文件
        for entry in fs::read_dir(&self.log_dir)? {
            let entry = entry?;
            let path = entry.path();
            
            if path.extension().and_then(|s| s.to_str()) == Some("log") {
                let content = fs::read_to_string(path)?;
                for line in content.lines() {
                    if let Ok(log_entry) = serde_json::from_str::<LogEntry>(line) {
                        // 过滤条件
                        if self.matches_query(&log_entry, &query) {
                            all_logs.push(log_entry);
                        }
                    }
                }
            }
        }
        
        // 排序（时间倒序）
        all_logs.sort_by(|a, b| b.timestamp.cmp(&a.timestamp));
        
        let total = all_logs.len() as u32;
        
        // 分页
        let logs = all_logs
            .into_iter()
            .skip(query.offset as usize)
            .take(query.limit as usize)
            .collect();
        
        Ok(LogQueryResult { total, logs })
    }
    
    /// 检查日志条目是否匹配查询条件
    fn matches_query(&self, entry: &LogEntry, query: &LogQuery) -> bool {
        // 级别过滤
        if let Some(level) = &query.level {
            if entry.level != level.to_string() {
                return false;
            }
        }
        
        // 时间范围过滤
        if let Some(start) = &query.start_time {
            if entry.timestamp < *start {
                return false;
            }
        }
        if let Some(end) = &query.end_time {
            if entry.timestamp > *end {
                return false;
            }
        }
        
        // 关键词过滤
        if let Some(keyword) = &query.keyword {
            if !entry.message.contains(keyword) && 
               !entry.target.contains(keyword) {
                return false;
            }
        }
        
        // 模块过滤
        if let Some(module) = &query.module {
            if let Some(entry_module) = &entry.module {
                if !entry_module.contains(module) {
                    return false;
                }
            } else {
                return false;
            }
        }
        
        true
    }
    
    /// 导出日志
    pub fn export_logs(
        &self, 
        query: LogQuery, 
        format: &str, 
        output_path: &PathBuf
    ) -> Result<(), Box<dyn std::error::Error>> {
        let result = self.query_logs(query)?;
        
        let content = match format {
            "json" => serde_json::to_string_pretty(&result.logs)?,
            "csv" => {
                let mut csv = String::from("timestamp,level,target,message\n");
                for log in &result.logs {
                    csv.push_str(&format!(
                        "{},{},{},\"{}\"\n",
                        log.timestamp.to_rfc3339(),
                        log.level,
                        log.target,
                        log.message.replace('"', "\"\"")
                    ));
                }
                csv
            },
            _ => return Err("不支持的导出格式".into()),
        };
        
        fs::write(output_path, content)?;
        Ok(())
    }
    
    /// 生成日志分析报告
    pub fn generate_report(&self, days: u32) -> Result<String, Box<dyn std::error::Error>> {
        let end_time = Utc::now();
        let start_time = end_time - chrono::Duration::days(days as i64);
        
        let query = LogQuery {
            level: None,
            start_time: Some(start_time),
            end_time: Some(end_time),
            keyword: None,
            module: None,
            limit: 10000,
            offset: 0,
        };
        
        let result = self.query_logs(query)?;
        
        // 统计信息
        let mut level_counts = std::collections::HashMap::new();
        let mut module_counts = std::collections::HashMap::new();
        
        for log in &result.logs {
            *level_counts.entry(log.level.clone()).or_insert(0) += 1;
            
            if let Some(module) = &log.module {
                *module_counts.entry(module.clone()).or_insert(0) += 1;
            }
        }
        
        // 生成报告
        let mut report = String::new();
        report.push_str(&format!("日志分析报告 ({} 天)\n\n", days));
        report.push_str(&format!("时间范围：{} - {}\n\n", start_time, end_time));
        
        report.push_str("日志级别统计:\n");
        for (level, count) in &level_counts {
            report.push_str(&format!("  {}: {}\n", level, count));
        }
        
        report.push_str("\n模块统计 (Top 10):\n");
        let mut modules: Vec<_> = module_counts.iter().collect();
        modules.sort_by(|a, b| b.1.cmp(a.1));
        for (module, count) in modules.iter().take(10) {
            report.push_str(&format!("  {}: {}\n", module, count));
        }
        
        report.push_str(&format!("\n总日志数：{}\n", result.total));
        
        Ok(report)
    }
}

// 使用示例
pub fn init_logging() -> Result<(), Box<dyn std::error::Error>> {
    let config = LogConfig {
        level: LogLevel::Debug,
        max_file_size: 10 * 1024 * 1024,  // 10MB
        max_files: 5,                      // 保留 5 个文件
        log_dir: PathBuf::from("logs"),
        enable_console: true,
        enable_file: true,
        enable_json: false,
    };
    
    FgAbyssLogger::init(config)?;
    
    log::info!("日志系统初始化完成");
    
    Ok(())
}
```

### 2.6.7 性能监控系统设计

```rust
use std::sync::{Arc, RwLock};
use std::time::{Duration, Instant};
use std::collections::HashMap;
use chrono::{DateTime, Utc};
use serde::{Serialize, Deserialize};
use rusqlite::Connection;

/// 性能指标类型
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum MetricType {
    // 系统资源
    CpuUsage,           // CPU 使用率
    MemoryUsage,        // 内存使用量
    DiskIO,            // 磁盘 I/O
    
    // 应用性能
    ResponseTime,       // 响应时间
    RequestCount,       // 请求数量
    ActiveConnections,  // 活跃连接数
    
    // WebShell 相关
    ShellConnectTime,   // Shell 连接时间
    ShellResponseTime,  // Shell 响应时间
    ShellSuccessRate,   // Shell 成功率
}

/// 性能指标数据
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MetricData {
    pub metric_type: MetricType,
    pub value: f64,
    pub unit: String,
    pub timestamp: DateTime<Utc>,
    pub tags: HashMap<String, String>,  // 标签（用于过滤）
}

/// 性能告警级别
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum AlertLevel {
    Info,       // 信息
    Warning,    // 警告
    Error,      // 错误
    Critical,   // 严重
}

/// 性能告警
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PerformanceAlert {
    pub alert_id: String,
    pub metric_type: MetricType,
    pub alert_level: AlertLevel,
    pub threshold: f64,
    pub current_value: f64,
    pub message: String,
    pub triggered_at: DateTime<Utc>,
    pub acknowledged: bool,
}

/// 性能监控配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PerformanceConfig {
    pub enabled: bool,
    pub collection_interval_secs: u64,    // 采集间隔（秒）
    pub retention_days: u32,               // 数据保留天数
    pub alert_enabled: bool,               // 是否启用告警
    pub thresholds: HashMap<MetricType, f64>,  // 告警阈值
}

/// 性能监控器
pub struct PerformanceMonitor {
    config: PerformanceConfig,
    metrics: Arc<RwLock<Vec<MetricData>>>,
    alerts: Arc<RwLock<Vec<PerformanceAlert>>>,
    running: Arc<RwLock<bool>>,
    db_connection: Arc<RwLock<Connection>>,
}

impl PerformanceMonitor {
    /// 创建新的性能监控器
    pub fn new(config: PerformanceConfig, db_path: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let db_connection = Connection::open(db_path)?;
        
        // 初始化数据库表
        Self::init_database(&db_connection)?;
        
        Ok(Self {
            config,
            metrics: Arc::new(RwLock::new(Vec::new())),
            alerts: Arc::new(RwLock::new(Vec::new())),
            running: Arc::new(RwLock::new(false)),
            db_connection: Arc::new(RwLock::new(db_connection)),
        })
    }
    
    /// 初始化数据库表
    fn init_database(conn: &Connection) -> Result<(), Box<dyn std::error::Error>> {
        conn.execute(
            "CREATE TABLE IF NOT EXISTS metrics (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                metric_type TEXT NOT NULL,
                value REAL NOT NULL,
                unit TEXT NOT NULL,
                timestamp TEXT NOT NULL,
                tags TEXT
            )",
            [],
        )?;
        
        conn.execute(
            "CREATE TABLE IF NOT EXISTS alerts (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                alert_id TEXT NOT NULL,
                metric_type TEXT NOT NULL,
                alert_level TEXT NOT NULL,
                threshold REAL NOT NULL,
                current_value REAL NOT NULL,
                message TEXT NOT NULL,
                triggered_at TEXT NOT NULL,
                acknowledged INTEGER DEFAULT 0
            )",
            [],
        )?;
        
        Ok(())
    }
    
    /// 启动监控任务
    pub fn start(&self) -> Result<(), Box<dyn std::error::Error>> {
        if !self.config.enabled {
            return Err("性能监控未启用".into());
        }
        
        *self.running.write().unwrap() = true;
        
        let metrics_clone = Arc::clone(&self.metrics);
        let alerts_clone = Arc::clone(&self.alerts);
        let running_clone = Arc::clone(&self.running);
        let config_clone = self.config.clone();
        let db_clone = Arc::clone(&self.db_connection);
        
        tokio::spawn(async move {
            let mut interval = tokio::time::interval(
                Duration::from_secs(config_clone.collection_interval_secs)
            );
            
            while *running_clone.read().unwrap() {
                interval.tick().await;
                
                // 收集系统指标
                let system_metrics = Self::collect_system_metrics();
                
                // 收集应用指标
                let app_metrics = Self::collect_app_metrics();
                
                // 合并指标
                let mut all_metrics = system_metrics.into_iter()
                    .chain(app_metrics.into_iter())
                    .collect::<Vec<_>>();
                
                // 检查告警
                let alerts = Self::check_alerts(&all_metrics, &config_clone);
                
                // 存储指标
                {
                    let mut metrics = metrics_clone.write().unwrap();
                    metrics.extend(all_metrics);
                    
                    // 限制内存中的指标数量
                    if metrics.len() > 10000 {
                        let drain = metrics.drain(0..metrics.len() - 5000);
                        // 持久化到数据库
                        Self::persist_metrics(&db_clone, drain.collect());
                    }
                }
                
                // 存储告警
                if !alerts.is_empty() {
                    let mut alerts_store = alerts_clone.write().unwrap();
                    alerts_store.extend(alerts);
                }
            }
        });
        
        Ok(())
    }
    
    /// 收集系统指标
    fn collect_system_metrics() -> Vec<MetricData> {
        let mut metrics = Vec::new();
        let timestamp = Utc::now();
        
        // CPU 使用率
        if let Some(cpu_usage) = Self::get_cpu_usage() {
            metrics.push(MetricData {
                metric_type: MetricType::CpuUsage,
                value: cpu_usage,
                unit: "%".to_string(),
                timestamp,
                tags: HashMap::new(),
            });
        }
        
        // 内存使用量
        if let Some(memory_usage) = Self::get_memory_usage() {
            metrics.push(MetricData {
                metric_type: MetricType::MemoryUsage,
                value: memory_usage,
                unit: "MB".to_string(),
                timestamp,
                tags: HashMap::new(),
            });
        }
        
        metrics
    }
    
    /// 收集应用指标
    fn collect_app_metrics() -> Vec<MetricData> {
        let mut metrics = Vec::new();
        let timestamp = Utc::now();
        
        // 活跃连接数
        let active_connections = Self::get_active_connections();
        metrics.push(MetricData {
            metric_type: MetricType::ActiveConnections,
            value: active_connections as f64,
            unit: "connections".to_string(),
            timestamp,
            tags: HashMap::new(),
        });
        
        // 请求数量
        let request_count = Self::get_request_count();
        metrics.push(MetricData {
            metric_type: MetricType::RequestCount,
            value: request_count as f64,
            unit: "requests".to_string(),
            timestamp,
            tags: HashMap::new(),
        });
        
        metrics
    }
    
    /// 检查告警
    fn check_alerts(metrics: &[MetricData], config: &PerformanceConfig) -> Vec<PerformanceAlert> {
        let mut alerts = Vec::new();
        
        for metric in metrics {
            if let Some(threshold) = config.thresholds.get(&metric.metric_type) {
                if metric.value > *threshold {
                    let alert = PerformanceAlert {
                        alert_id: uuid::Uuid::new_v4().to_string(),
                        metric_type: metric.metric_type.clone(),
                        alert_level: Self::calculate_alert_level(metric.value, *threshold),
                        threshold: *threshold,
                        current_value: metric.value,
                        message: format!(
                            "{} 超过阈值：{:.2} > {:.2}",
                            format!("{:?}", metric.metric_type),
                            metric.value,
                            *threshold
                        ),
                        triggered_at: metric.timestamp,
                        acknowledged: false,
                    };
                    
                    alerts.push(alert);
                }
            }
        }
        
        alerts
    }
    
    /// 计算告警级别
    fn calculate_alert_level(value: f64, threshold: f64) -> AlertLevel {
        let ratio = value / threshold;
        
        if ratio > 2.0 {
            AlertLevel::Critical
        } else if ratio > 1.5 {
            AlertLevel::Error
        } else if ratio > 1.2 {
            AlertLevel::Warning
        } else {
            AlertLevel::Info
        }
    }
    
    /// 获取实时指标
    pub fn get_current_metrics(&self) -> Vec<MetricData> {
        let metrics = self.metrics.read().unwrap();
        
        // 返回最近 1 分钟的指标
        let one_minute_ago = Utc::now() - Duration::from_secs(60);
        
        metrics.iter()
            .filter(|m| m.timestamp > one_minute_ago)
            .cloned()
            .collect()
    }
    
    /// 获取历史指标
    pub fn get_historical_metrics(
        &self,
        metric_type: MetricType,
        start_time: DateTime<Utc>,
        end_time: DateTime<Utc>,
    ) -> Result<Vec<MetricData>, Box<dyn std::error::Error>> {
        let conn = self.db_connection.read().unwrap();
        
        let mut stmt = conn.prepare(
            "SELECT metric_type, value, unit, timestamp, tags 
             FROM metrics 
             WHERE metric_type = ? 
             AND timestamp BETWEEN ? AND ?
             ORDER BY timestamp ASC"
        )?;
        
        let metrics = stmt.query_map(
            [format!("{:?}", metric_type), start_time.to_rfc3339(), end_time.to_rfc3339()],
            |row| {
                Ok(MetricData {
                    metric_type: serde_json::from_str(&row.get::<_, String>(0)?).unwrap_or(MetricType::CpuUsage),
                    value: row.get(1)?,
                    unit: row.get(2)?,
                    timestamp: DateTime::parse_from_rfc3339(&row.get::<_, String>(3)?).unwrap().into(),
                    tags: serde_json::from_str(&row.get::<_, String>(4)?).unwrap_or_default(),
                })
            },
        )?;
        
        let mut result = Vec::new();
        for metric in metrics {
            result.push(metric?);
        }
        
        Ok(result)
    }
    
    /// 获取未确认的告警
    pub fn get_unacknowledged_alerts(&self) -> Vec<PerformanceAlert> {
        let alerts = self.alerts.read().unwrap();
        
        alerts.iter()
            .filter(|a| !a.acknowledged)
            .cloned()
            .collect()
    }
    
    /// 确认告警
    pub fn acknowledge_alert(&self, alert_id: &str) -> Result<(), Box<dyn std::error::Error>> {
        let mut alerts = self.alerts.write().unwrap();
        
        for alert in alerts.iter_mut() {
            if alert.alert_id == alert_id {
                alert.acknowledged = true;
                return Ok(());
            }
        }
        
        Err("告警不存在".into())
    }
    
    /// 生成性能报告
    pub fn generate_report(&self, days: u32) -> Result<String, Box<dyn std::error::Error>> {
        let end_time = Utc::now();
        let start_time = end_time - Duration::from_secs(days as u64 * 86400);
        
        let mut report = String::new();
        report.push_str(&format!("性能报告 ({} 天)\n\n", days));
        report.push_str(&format!("时间范围：{} - {}\n\n", start_time, end_time));
        
        // 统计各指标的平均值
        let metrics = self.get_historical_metrics(MetricType::CpuUsage, start_time, end_time)?;
        if !metrics.is_empty() {
            let avg: f64 = metrics.iter().map(|m| m.value).sum::<f64>() / metrics.len() as f64;
            report.push_str(&format!("CPU 平均使用率：{:.2}%\n", avg));
        }
        
        // 告警统计
        let alerts = self.alerts.read().unwrap();
        let alert_count = alerts.iter()
            .filter(|a| a.triggered_at >= start_time && a.triggered_at <= end_time)
            .count();
        
        report.push_str(&format!("\n总告警数：{}\n", alert_count));
        
        Ok(report)
    }
    
    /// 持久化指标到数据库
    fn persist_metrics(db: &Arc<RwLock<Connection>>, metrics: impl Iterator<Item = MetricData>) {
        if let Ok(conn) = db.write() {
            for metric in metrics {
                let _ = conn.execute(
                    "INSERT INTO metrics (metric_type, value, unit, timestamp, tags) 
                     VALUES (?1, ?2, ?3, ?4, ?5)",
                    [
                        format!("{:?}", metric.metric_type),
                        metric.value.to_string(),
                        metric.unit,
                        metric.timestamp.to_rfc3339(),
                        serde_json::to_string(&metric.tags).unwrap_or_default(),
                    ],
                );
            }
        }
    }
    
    // ============ 系统指标采集辅助函数 ============
    
    fn get_cpu_usage() -> Option<f64> {
        // 使用 sysinfo crate 获取 CPU 使用率
        // 这里简化实现
        Some(25.5)  // 示例值
    }
    
    fn get_memory_usage() -> Option<f64> {
        // 使用 sysinfo crate 获取内存使用量
        // 这里简化实现
        Some(512.0)  // 示例值：512MB
    }
    
    fn get_active_connections() -> u32 {
        // 从连接管理器获取活跃连接数
        12  // 示例值
    }
    
    fn get_request_count() -> u32 {
        // 从 HTTP 客户端获取请求计数
        1500  // 示例值
    }
}

/// 性能监控 API（供前端调用）
pub mod performance_api {
    use super::*;
    
    /// 获取实时性能数据
    #[tauri::command]
    pub fn get_current_metrics(monitor: tauri::State<PerformanceMonitor>) -> Vec<MetricData> {
        monitor.get_current_metrics()
    }
    
    /// 获取历史性能数据
    #[tauri::command]
    pub fn get_historical_metrics(
        monitor: tauri::State<PerformanceMonitor>,
        metric_type: String,
        start_time: String,
        end_time: String,
    ) -> Result<Vec<MetricData>, String> {
        let start = DateTime::parse_from_rfc3339(&start_time).map_err(|e| e.to_string())?;
        let end = DateTime::parse_from_rfc3339(&end_time).map_err(|e| e.to_string())?;
        
        monitor.get_historical_metrics(
            serde_json::from_str(&format!("\"{}\"", metric_type)).map_err(|e| e.to_string())?,
            start.into(),
            end.into(),
        ).map_err(|e| e.to_string())
    }
    
    /// 获取未确认告警
    #[tauri::command]
    pub fn get_unacknowledged_alerts(monitor: tauri::State<PerformanceMonitor>) -> Vec<PerformanceAlert> {
        monitor.get_unacknowledged_alerts()
    }
    
    /// 确认告警
    #[tauri::command]
    pub fn acknowledge_alert(monitor: tauri::State<PerformanceMonitor>, alert_id: String) -> Result<(), String> {
        monitor.acknowledge_alert(&alert_id).map_err(|e| e.to_string())
    }
    
    /// 生成性能报告
    #[tauri::command]
    pub fn generate_report(monitor: tauri::State<PerformanceMonitor>, days: u32) -> Result<String, String> {
        monitor.generate_report(days).map_err(|e| e.to_string())
    }
}

// 使用示例
pub fn init_performance_monitoring() -> Result<Arc<PerformanceMonitor>, Box<dyn std::error::Error>> {
    let config = PerformanceConfig {
        enabled: true,
        collection_interval_secs: 10,  // 每 10 秒采集一次
        retention_days: 30,             // 保留 30 天数据
        alert_enabled: true,
        thresholds: {
            let mut thresholds = HashMap::new();
            thresholds.insert(MetricType::CpuUsage, 80.0);        // CPU > 80% 告警
            thresholds.insert(MetricType::MemoryUsage, 800.0);    // 内存 > 800MB 告警
            thresholds.insert(MetricType::ActiveConnections, 200); // 连接 > 200 告警
            thresholds
        },
    };
    
    let monitor = Arc::new(PerformanceMonitor::new(config, "data/metrics.db")?);
    
    // 启动监控任务
    monitor.start()?;
    
    log::info!("性能监控系统已启动");
    
    Ok(monitor)
}
```

### 2.6.2 配置管理系统设计 (config.rs)

```rust
use serde::{Serialize, Deserialize};
use std::fs;
use std::path::{Path, PathBuf};
use toml;
use notify::{Config, RecommendedWatcher, RecursiveMode, Watcher};
use std::sync::{Arc, Mutex};
use std::time::Duration;

/// 应用配置
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppConfig {
    pub app: AppSection,
    pub ui: UiSection,
    pub connection: ConnectionSection,
    pub security: SecuritySection,
    pub log: LogSection,
}

/// 应用配置段
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppSection {
    pub language: String,
    pub check_update: bool,
    pub start_on_boot: bool,
}

/// UI 配置段
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UiSection {
    pub theme: String,
    pub accent_color: String,
    pub font_size: u32,
}

/// 连接配置段
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConnectionSection {
    pub timeout: u32,
    pub retry_count: u32,
    pub proxy_enabled: bool,
    pub proxy_type: String,
    pub proxy_host: String,
    pub proxy_port: u32,
}

/// 安全配置段
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecuritySection {
    pub encryption_algorithm: String,
    pub key_rotation_days: u32,
    pub data_encryption: bool,
    pub log_desensitization: bool,
}

/// 日志配置段
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogSection {
    pub level: String,
    pub max_file_size: u64,
    pub max_files: u32,
}

/// 配置管理器
pub struct ConfigManager {
    config_path: PathBuf,
    config: Arc<Mutex<AppConfig>>,
    watchers: Arc<Mutex<Vec<RecommendedWatcher>>>,
    callbacks: Arc<Mutex<Vec<Box<dyn Fn(&AppConfig) + Send + Sync>>>>,
}

impl ConfigManager {
    /// 创建新的配置管理器
    pub fn new(config_path: PathBuf) -> Result<Self, Box<dyn std::error::Error>> {
        let config = if config_path.exists() {
            // 加载现有配置
            Self::load_from_file(&config_path)?
        } else {
            // 创建默认配置
            Self::default_config()
        };
        
        Ok(Self {
            config_path,
            config: Arc::new(Mutex::new(config)),
            watchers: Arc::new(Mutex::new(Vec::new())),
            callbacks: Arc::new(Mutex::new(Vec::new())),
        })
    }
    
    /// 创建默认配置
    fn default_config() -> AppConfig {
        AppConfig {
            app: AppSection {
                language: "zh-CN".to_string(),
                check_update: true,
                start_on_boot: false,
            },
            ui: UiSection {
                theme: "system".to_string(),
                accent_color: "#18a058".to_string(),
                font_size: 14,
            },
            connection: ConnectionSection {
                timeout: 30,
                retry_count: 3,
                proxy_enabled: false,
                proxy_type: "http".to_string(),
                proxy_host: "127.0.0.1".to_string(),
                proxy_port: 7890,
            },
            security: SecuritySection {
                encryption_algorithm: "aes-256-gcm".to_string(),
                key_rotation_days: 30,
                data_encryption: true,
                log_desensitization: true,
            },
            log: LogSection {
                level: "info".to_string(),
                max_file_size: 10 * 1024 * 1024,
                max_files: 5,
            },
        }
    }
    
    /// 从文件加载配置
    fn load_from_file(path: &Path) -> Result<AppConfig, Box<dyn std::error::Error>> {
        let content = fs::read_to_string(path)?;
        let config: AppConfig = toml::from_str(&content)?;
        
        // 验证配置
        Self::validate_config(&config)?;
        
        Ok(config)
    }
    
    /// 验证配置
    fn validate_config(config: &AppConfig) -> Result<(), Box<dyn std::error::Error>> {
        // 验证语言
        if !["zh-CN", "en-US"].contains(&config.app.language.as_str()) {
            return Err("无效的语言设置".into());
        }
        
        // 验证主题
        if !["light", "dark", "system"].contains(&config.ui.theme.as_str()) {
            return Err("无效的主题设置".into());
        }
        
        // 验证加密算法
        if !["aes-256-gcm", "xor"].contains(&config.security.encryption_algorithm.as_str()) {
            return Err("无效的加密算法".into());
        }
        
        // 验证超时时间
        if config.connection.timeout == 0 || config.connection.timeout > 300 {
            return Err("超时时间必须在 1-300 秒之间".into());
        }
        
        // 验证重试次数
        if config.connection.retry_count > 10 {
            return Err("重试次数不能超过 10 次".into());
        }
        
        Ok(())
    }
    
    /// 保存配置到文件
    pub fn save(&self) -> Result<(), Box<dyn std::error::Error>> {
        let config = self.config.lock().unwrap();
        let content = toml::to_string_pretty(&*config)?;
        
        // 确保目录存在
        if let Some(parent) = self.config_path.parent() {
            fs::create_dir_all(parent)?;
        }
        
        // 写入临时文件
        let temp_path = self.config_path.with_extension("toml.tmp");
        fs::write(&temp_path, &content)?;
        
        // 原子替换
        fs::rename(&temp_path, &self.config_path)?;
        
        Ok(())
    }
    
    /// 获取配置
    pub fn get_config(&self) -> AppConfig {
        self.config.lock().unwrap().clone()
    }
    
    /// 更新配置
    pub fn update<F>(&self, f: F) -> Result<(), Box<dyn std::error::Error>>
    where
        F: FnOnce(&mut AppConfig),
    {
        let mut config = self.config.lock().unwrap();
        f(&mut config);
        
        // 验证配置
        Self::validate_config(&config)?;
        
        // 保存配置
        drop(config);
        self.save()?;
        
        // 通知回调
        self.notify_callbacks();
        
        Ok(())
    }
    
    /// 导入配置
    pub fn import(&self, path: &Path) -> Result<(), Box<dyn std::error::Error>> {
        let config = Self::load_from_file(path)?;
        
        let mut current = self.config.lock().unwrap();
        *current = config;
        
        self.notify_callbacks();
        
        Ok(())
    }
    
    /// 导出配置
    pub fn export(&self, path: &Path) -> Result<(), Box<dyn std::error::Error>> {
        let config = self.config.lock().unwrap();
        let content = toml::to_string_pretty(&*config)?;
        fs::write(path, content)?;
        Ok(())
    }
    
    /// 恢复默认配置
    pub fn reset(&self) -> Result<(), Box<dyn std::error::Error>> {
        let mut config = self.config.lock().unwrap();
        *config = Self::default_config();
        
        self.save()?;
        self.notify_callbacks();
        
        Ok(())
    }
    
    /// 回滚配置（从备份恢复）
    pub fn rollback(&self, backup_path: &Path) -> Result<(), Box<dyn std::error::Error>> {
        self.import(backup_path)
    }
    
    /// 创建配置备份
    pub fn create_backup(&self) -> Result<PathBuf, Box<dyn std::error::Error>> {
        let timestamp = chrono::Utc::now().format("%Y%m%d_%H%M%S");
        let backup_path = self.config_path.with_extension(format!("toml.backup.{}", timestamp));
        
        self.export(&backup_path)?;
        
        Ok(backup_path)
    }
    
    /// 启用配置热重载
    pub fn enable_hot_reload(&self) -> Result<(), Box<dyn std::error::Error>> {
        let config_path = self.config_path.clone();
        let config_clone = Arc::clone(&self.config);
        let callbacks_clone = Arc::clone(&self.callbacks);
        
        let mut watcher = RecommendedWatcher::new(
            move |event: notify::Result<notify::Event>| {
                if let Ok(event) = event {
                    if event.kind.is_modify() {
                        // 重新加载配置
                        if let Ok(new_config) = Self::load_from_file(&config_path) {
                            let mut config = config_clone.lock().unwrap();
                            *config = new_config;
                            
                            // 通知回调
                            for callback in callbacks_clone.lock().unwrap().iter() {
                                callback(&config);
                            }
                        }
                    }
                }
            },
            Config::default(),
        )?;
        
        watcher.watch(&self.config_path, RecursiveMode::NonRecursive)?;
        
        self.watchers.lock().unwrap().push(watcher);
        
        Ok(())
    }
    
    /// 注册配置变更回调
    pub fn on_change<F>(&self, callback: F)
    where
        F: Fn(&AppConfig) + Send + Sync + 'static,
    {
        self.callbacks.lock().unwrap().push(Box::new(callback));
    }
    
    /// 通知所有回调
    fn notify_callbacks(&self) {
        let config = self.config.lock().unwrap();
        for callback in self.callbacks.lock().unwrap().iter() {
            callback(&config);
        }
    }
}

// 使用示例
pub fn init_config() -> Result<ConfigManager, Box<dyn std::error::Error>> {
    let config_path = PathBuf::from("config/config.toml");
    
    let manager = ConfigManager::new(config_path)?;
    
    // 启用热重载
    manager.enable_hot_reload()?;
    
    // 注册配置变更回调
    manager.on_change(|config| {
        log::info!("配置已更新：{:?}", config);
    });
    
    Ok(manager)
}
```

***

## 3. 开发规范

### 3.1 Rust 代码规范

#### 命名规范

```rust
// 类型使用 PascalCase
pub struct ShellManager { }
pub enum ErrorCode { }
pub trait Plugin { }

// 函数和变量使用 snake_case
pub fn create_shell() { }
let shell_id = "123";

// 常量使用 SCREAMING_SNAKE_CASE
pub const MAX_CONNECTIONS: usize = 100;

// 模块文件使用 snake_case
// shell_manager.rs
// error_handler.rs
```

#### 错误处理

```rust
// 使用 Result<T, E>
pub fn create_shell(url: String) -> Result<Shell, FgAbyssError> {
    // 使用 ? 传播错误
    let shell = validate_url(&url)?;
    Ok(shell)
}

// 使用 thiserror 定义错误
#[derive(Debug, thiserror::Error)]
pub enum FgAbyssError {
    #[error("Shell 不存在：{0}")]
    ShellNotFound(String),
    
    #[error("HTTP 错误：{0}")]
    Http(#[from] reqwest::Error),
}
```

#### 异步编程

```rust
// 使用 tokio 异步运行时
use tokio::sync::RwLock;

pub async fn get_shell(&self, id: &str) -> Result<Shell> {
    // 异步操作
    let shell = self.db.get_shell(id).await?;
    Ok(shell)
}

// 使用 Arc 共享状态
use std::sync::Arc;

pub struct App {
    shell_manager: Arc<ShellManager>,
}
```

#### 文档注释

````rust
/// Shell 管理器
/// 
/// 负责 WebShell 的 CRUD 操作和连接管理
/// 
/// # Examples
/// 
/// ```
/// let manager = ShellManager::new(db);
/// let shell = manager.create_shell(url, password).await?;
/// ```
pub struct ShellManager {
    // ...
}
````

### 3.2 Vue 代码规范

#### 组件结构

```vue
<template>
  <!-- 模板 -->
</template>

<script setup lang="ts">
// 导入
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

// Props
interface Props {
  shellId: string
  showActions?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showActions: true,
})

// Emits
const emit = defineEmits<{
  (e: 'update', id: string): void
  (e: 'delete', id: string): void
}>()

// 状态
const loading = ref(false)

// 计算属性
const shellName = computed(() => {
  return props.shellId
})

// 函数
const handleUpdate = async () => {
  loading.value = true
  try {
    // 业务逻辑
    emit('update', props.shellId)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 样式 */
</style>
```

#### TypeScript 类型

```typescript
// 定义类型
export interface Shell {
  id: string
  name: string
  url: string
  payloadType: string
  status: 'online' | 'offline'
  createdAt: Date
}

// 定义 API 响应
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}

// 使用泛型
export type StoreState<T> = {
  items: T[]
  loading: boolean
  error: string | null
}
```

### 3.3 Git 规范

#### Commit Message

```bash
# 格式：<type>(<scope>): <subject>

# type 类型
feat:     新功能
fix:      修复 bug
docs:     文档更新
style:    代码格式（不影响代码运行）
refactor: 重构
test:     测试相关
chore:    构建/工具链相关

# 示例
feat(shell): 添加 Shell 批量删除功能
fix(crypto): 修复 AES 解密边界条件
docs(readme): 更新安装说明
refactor(http): 重构 HTTP 客户端错误处理
```

#### 分支管理

```bash
# 主分支
main          # 生产环境代码
develop       # 开发分支

# 功能分支
feature/shell-crud
feature/file-manager

# 修复分支
fix/connection-issue
fix/memory-leak

# 发布分支
release/v0.2.0
```

### 3.4 代码审查清单

#### Rust

- [ ] 使用 `cargo clippy` 检查代码
- [ ] 使用 `cargo fmt` 格式化代码
- [ ] 错误处理完整
- [ ] 文档注释完整
- [ ] 单元测试覆盖
- [ ] 无内存泄漏（使用 `cargo leak-check`）

#### Vue

- [ ] 使用 ESLint 检查代码
- [ ] 使用 Prettier 格式化代码
- [ ] Props 类型定义完整
- [ ] Emits 类型定义完整
- [ ] 组件文档注释
- [ ] 无 TypeScript 错误

***

## 4. 核心模块开发指南

### 4.1 Shell 管理模块

#### 后端实现

```rust
// src-tauri/src/core/shell/manager.rs

use crate::error::Result;
use crate::database::DatabasePool;
use std::collections::HashMap;
use tokio::sync::RwLock;
use std::sync::Arc;

pub struct ShellManager {
    db: Arc<DatabasePool>,
    cache: Arc<RwLock<HashMap<String, ShellCacheEntry>>>,
}

impl ShellManager {
    pub fn new(db: Arc<DatabasePool>) -> Self {
        Self {
            db,
            cache: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    /// 创建 Shell
    pub async fn create_shell(
        &self,
        url: String,
        password: String,
        payload_type: String,
        project_id: Option<String>,
    ) -> Result<ShellEntity> {
        // 1. 验证 URL
        validate_url(&url)?;

        // 2. 生成 ID
        let id = uuid::Uuid::new_v4().to_string();

        // 3. 创建实体
        let shell = ShellEntity {
            id: id.clone(),
            name: format!("Shell {}", id[..8].to_string()),
            url,
            password, // TODO: 加密存储
            payload_type,
            project_id,
            created_at: chrono::Utc::now(),
            updated_at: chrono::Utc::now(),
        };

        // 4. 保存数据库
        self.db.save_shell(&shell).await?;

        Ok(shell)
    }

    /// 删除 Shell
    pub async fn delete_shell(&self, id: &str) -> Result<()> {
        // 1. 软删除
        self.db.delete_shell(id).await?;

        // 2. 清除缓存
        let mut cache = self.cache.write().await;
        cache.remove(id);

        Ok(())
    }

    /// 测试连接
    pub async fn test_connection(&self, shell_id: &str) -> Result<ConnectionTestResult> {
        let shell = self.db.get_shell(shell_id).await?
            .ok_or_else(|| FgAbyssError::ShellNotFound(shell_id.to_string()))?;

        // 使用 HTTP 客户端测试
        let client = reqwest::Client::new();
        let response = client.get(&shell.url).send().await?;

        Ok(ConnectionTestResult {
            success: response.status().is_success(),
            status_code: Some(response.status().as_u16()),
            response_time: 0, // TODO: 计算响应时间
        })
    }
}
```

#### Tauri Command

```rust
// src-tauri/src/commands/shell.rs

use crate::core::shell::ShellManager;
use tauri::State;

#[tauri::command]
pub async fn create_shell(
    manager: State<'_, ShellManager>,
    url: String,
    password: String,
    payload_type: String,
    project_id: Option<String>,
) -> Result<ShellEntity, String> {
    manager.create_shell(url, password, payload_type, project_id)
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn delete_shell(
    manager: State<'_, ShellManager>,
    shell_id: String,
) -> Result<(), String> {
    manager.delete_shell(&shell_id)
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn test_connection(
    manager: State<'_, ShellManager>,
    shell_id: String,
) -> Result<ConnectionTestResult, String> {
    manager.test_connection(&shell_id)
        .await
        .map_err(|e| e.to_string())
}
```

#### 前端实现

```typescript
// src/api/shell.ts
import { invoke } from '@tauri-apps/api/core'
import type { Shell, ConnectionTestResult } from '@/types'

export async function createShell(
  url: string,
  password: string,
  payloadType: string,
  projectId?: string
): Promise<Shell> {
  return invoke('create_shell', {
    url,
    password,
    payloadType,
    projectId,
  })
}

export async function deleteShell(shellId: string): Promise<void> {
  return invoke('delete_shell', { shellId })
}

export async function testConnection(shellId: string): Promise<ConnectionTestResult> {
  return invoke('test_connection', { shellId })
}
```

```vue
<!-- src/components/shell/ShellForm.vue -->
<template>
  <n-form ref="formRef" :model="formValue" :rules="rules">
    <n-form-item label="URL" path="url">
      <n-input v-model:value="formValue.url" placeholder="http://example.com/shell.php" />
    </n-form-item>

    <n-form-item label="密码" path="password">
      <n-input v-model:value="formValue.password" type="password" />
    </n-form-item>

    <n-form-item label="类型" path="payloadType">
      <n-select
        v-model:value="formValue.payloadType"
        :options="[
          { label: 'PHP', value: 'php' },
          { label: 'JSP', value: 'jsp' },
          { label: 'ASP', value: 'asp' },
          { label: 'ASPX', value: 'aspx' },
        ]"
      />
    </n-form-item>

    <n-form-item>
      <n-button type="primary" @click="handleSubmit" :loading="loading">
        创建
      </n-button>
    </n-form-item>
  </n-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { createShell } from '@/api/shell'

const message = useMessage()

const formRef = ref()
const loading = ref(false)

const formValue = ref({
  url: '',
  password: '',
  payloadType: 'php',
})

const rules = {
  url: {
    required: true,
    message: '请输入 URL',
    trigger: 'blur',
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: 'blur',
  },
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true

    await createShell(
      formValue.value.url,
      formValue.value.password,
      formValue.value.payloadType
    )

    message.success('创建成功')
    emit('success')
  } catch (error) {
    message.error(`创建失败：${error}`)
  } finally {
    loading.value = false
  }
}
</script>
```

### 4.2 加密模块

```rust
// src-tauri/src/core/crypto/aes.rs

use aes_gcm::{
    aead::{Aead, KeyInit},
    Aes256Gcm, Nonce,
};
use rand::Rng;

/// AES-256-GCM 加密
pub fn aes_encrypt(data: &[u8], key: &[u8]) -> Result<Vec<u8>, CryptoError> {
    if key.len() != 32 {
        return Err(CryptoError::InvalidKeyLength);
    }

    let cipher = Aes256Gcm::new_from_slice(key)?;
    let nonce = generate_nonce();

    let ciphertext = cipher.encrypt(&nonce, data)?;

    // 格式：nonce(12) + ciphertext
    let mut result = Vec::with_capacity(12 + ciphertext.len());
    result.extend_from_slice(nonce.as_slice());
    result.extend_from_slice(&ciphertext);

    Ok(result)
}

/// AES-256-GCM 解密
pub fn aes_decrypt(encrypted: &[u8], key: &[u8]) -> Result<Vec<u8>, CryptoError> {
    if key.len() != 32 {
        return Err(CryptoError::InvalidKeyLength);
    }

    if encrypted.len() < 12 {
        return Err(CryptoError::InvalidData);
    }

    let nonce = Nonce::from_slice(&encrypted[0..12]);
    let ciphertext = &encrypted[12..];

    let cipher = Aes256Gcm::new_from_slice(key)?;
    let plaintext = cipher.decrypt(nonce, ciphertext)?;

    Ok(plaintext)
}

/// 生成随机 Nonce
fn generate_nonce() -> Nonce<Aes256Gcm> {
    let mut nonce_bytes = [0u8; 12];
    rand::thread_rng().fill_bytes(&mut nonce_bytes);
    Nonce::from(nonce_bytes)
}
```

### 4.3 HTTP 客户端

```rust
// src-tauri/src/core/http/client.rs

use reqwest::{Client, Proxy};
use std::time::Duration;

pub struct HttpClient {
    client: Client,
    proxy: Option<ProxyConfig>,
    timeout: Duration,
}

impl HttpClient {
    pub fn new() -> Self {
        let client = Client::builder()
            .timeout(Duration::from_secs(30))
            .danger_accept_invalid_certs(false) // 生产环境应为 true
            .build()
            .expect("Failed to create HTTP client");

        Self {
            client,
            proxy: None,
            timeout: Duration::from_secs(30),
        }
    }

    pub fn with_proxy(mut self, proxy: ProxyConfig) -> Self {
        let proxy = Proxy::all(&proxy.url).expect("Failed to create proxy");
        self.client = Client::builder()
            .proxy(proxy)
            .timeout(self.timeout)
            .build()
            .expect("Failed to create HTTP client with proxy");
        self.proxy = Some(proxy);
        self
    }

    pub async fn get(&self, url: &str) -> Result<String, reqwest::Error> {
        let response = self.client.get(url).send().await?;
        let text = response.text().await?;
        Ok(text)
    }

    pub async fn post(&self, url: &str, body: &[u8]) -> Result<String, reqwest::Error> {
        let response = self.client.post(url).body(body.to_vec()).send().await?;
        let text = response.text().await?;
        Ok(text)
    }
}
```

***

## 5. 测试

### 5.1 单元测试

#### Rust 测试

```rust
// src-tauri/src/core/crypto/tests.rs

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_aes_encrypt_decrypt_round_trip() {
        let key = b"this_is_a_32_byte_key_for_aes256";
        let data = b"Hello, World!";

        let encrypted = aes_encrypt(data, key).unwrap();
        let decrypted = aes_decrypt(&encrypted, key).unwrap();

        assert_eq!(data, &decrypted[..]);
    }

    #[test]
    fn test_aes_decrypt_invalid_data() {
        let key = b"this_is_a_32_byte_key_for_aes256";
        let invalid_data = b"too_short";

        let result = aes_decrypt(invalid_data, key);
        assert!(result.is_err());
    }
}
```

运行测试：

```bash
cd src-tauri
cargo test
```

#### Vue 组件测试

```typescript
// src/components/shell/__tests__/ShellForm.test.ts
import { mount } from '@vue/test-utils'
import ShellForm from '../ShellForm.vue'
import { createTestingPinia } from '@pinia/testing'

describe('ShellForm', () => {
  it('should create shell successfully', async () => {
    const wrapper = mount(ShellForm, {
      global: {
        plugins: [createTestingPinia()],
      },
    })

    // 填写表单
    await wrapper.find('input[type="text"]').setValue('http://example.com/shell.php')
    await wrapper.find('input[type="password"]').setValue('password123')

    // 提交
    await wrapper.find('button[type="submit"]').trigger('click')

    // 验证
    expect(wrapper.emitted('success')).toBeTruthy()
  })
})
```

运行测试：

```bash
pnpm test
```

### 5.2 集成测试

```rust
// tests/integration/shell_manager.rs

#[tokio::test]
async fn test_shell_crud_operations() {
    // 1. 准备测试数据库
    let db = DatabasePool::new(":memory:").await.unwrap();
    let manager = ShellManager::new(Arc::new(db));

    // 2. 创建 Shell
    let shell = manager
        .create_shell(
            "http://example.com/shell.php".to_string(),
            "password123".to_string(),
            "php".to_string(),
            None,
        )
        .await
        .unwrap();

    // 3. 验证创建成功
    assert_eq!(shell.payload_type, "php");

    // 4. 获取 Shell
    let retrieved = manager.get_shell(&shell.id).await.unwrap();
    assert_eq!(retrieved.id, shell.id);

    // 5. 更新 Shell
    let updates = ShellUpdateDto {
        name: Some("Updated Shell".to_string()),
        ..Default::default()
    };
    manager.update_shell(&retrieved.id, updates).await.unwrap();

    // 6. 验证更新成功
    let updated = manager.get_shell(&retrieved.id).await.unwrap();
    assert_eq!(updated.name, "Updated Shell");

    // 7. 删除 Shell
    manager.delete_shell(&retrieved.id).await.unwrap();

    // 8. 验证删除成功
    let result = manager.get_shell(&retrieved.id).await;
    assert!(result.is_err());
}
```

### 5.3 基准测试

```rust
// benches/crypto_benchmark.rs

use criterion::{black_box, criterion_group, criterion_main, Criterion};

fn benchmark_aes_encrypt(c: &mut Criterion) {
    let key = b"this_is_a_32_byte_key_for_aes256";
    let data = b"Hello, World! This is a test message.";

    c.bench_function("aes_encrypt_64bytes", |b| {
        b.iter(|| {
            aes_encrypt(black_box(data), black_box(key)).unwrap()
        })
    });
}

criterion_group!(benches, benchmark_aes_encrypt);
criterion_main!(benches);
```

运行基准测试：

```bash
cd src-tauri
cargo bench
```

***

## 6. 构建和发布

### 6.1 开发构建

```bash
# 开发模式（热重载）
pnpm tauri dev

# 调试构建
pnpm tauri build --debug
```

### 6.2 生产构建

```bash
# 发布构建
pnpm tauri build

# 构建产物位置
# Windows: src-tauri/target/release/FG-ABYSS.exe
# macOS: src-tauri/target/release/bundle/macos/FG-ABYSS.app
# Linux: src-tauri/target/release/bundle/deb/fg-abyss.deb
```

### 6.3 跨平台构建

```bash
# 使用 GitHub Actions 自动构建
# .github/workflows/build.yml
name: Build
on: push
jobs:
  build:
    strategy:
      matrix:
        platform: [windows-latest, macos-latest, ubuntu-latest]
    runs-on: ${{ matrix.platform }}
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
      - name: Setup Rust
        uses: actions-rs/toolchain@v1
      - name: Install dependencies
        run: pnpm install
      - name: Build
        run: pnpm tauri build
```

### 6.4 代码签名

#### Windows

```bash
# 使用 SignTool 签名
signtool sign /f certificate.pfx /p password /t http://timestamp.digicert.com FG-ABYSS.exe
```

#### macOS

```bash
# 使用 codesign 签名
codesign --sign "Developer ID Application: Your Name" --timestamp --options runtime FG-ABYSS.app
```

***

## 7. 调试

### 7.1 Rust 调试

```bash
# 使用 lldb 调试
lldb target/debug/fg-abyss
(lldb) break set -name create_shell
(lldb) run
```

### 7.2 前端调试

- 使用浏览器 DevTools
- 启用 Vue Devtools 扩展
- 使用 console.log 调试

### 7.3 日志调试

```rust
// 启用调试日志
RUST_LOG=debug cargo run

// 代码中添加日志
log::debug!("Creating shell: {}", url);
log::info!("Shell created successfully");
log::error!("Failed to create shell: {}", error);
```

***

## 8. 性能优化

### 8.1 Rust 性能优化

```rust
// 使用性能分析工具
cargo install cargo-flamegraph
cargo flamegraph --root -- cargo test --release

// 优化建议：
// 1. 使用迭代器而非循环
// 2. 避免不必要的克隆
// 3. 使用合适的集合类型
// 4. 利用并行计算
```

### 8.2 前端性能优化

```typescript
// 1. 组件懒加载
const ShellManager = defineAsyncComponent(
  () => import('@/views/ShellManager.vue')
)

// 2. 虚拟列表
import { useVirtualList } from '@vueuse/core'

// 3. 防抖节流
import { useDebounceFn, useThrottleFn } from '@vueuse/core'
```

***

## 9. 常见问题

### Q1: 构建失败 "WebView2 未找到"

**解决**: 安装 WebView2 Runtime

```bash
# Windows 10/11 自带
# 旧版本 Windows: 下载安装
```

### Q2: Rust 编译慢

**解决**:

```bash
# 使用 sccache 加速
cargo install sccache
# 配置 .cargo/config.toml
[build]
rustc-wrapper = "sccache"
```

### Q3: 前端热重载不工作

**解决**:

```bash
# 清理缓存
rm -rf node_modules
pnpm install

# 重启开发服务器
```

***

## 10. 贡献指南

### 10.1 提交 PR

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 10.2 代码风格

- 遵循 Rust 官方代码风格
- 遵循 Vue 3 最佳实践
- 使用 Clippy 和 ESLint 检查代码

### 10.3 测试要求

- 新功能必须包含单元测试
- 核心功能需要集成测试
- 所有测试必须通过

***

**文档结束**
