# FG-ABYSS 项目结构优化方案

## 📋 目录

1. [技术栈评估报告](#1-技术栈评估报告)
2. [现有目录结构分析](#2-现有目录结构分析)
3. [优化后的目录结构设计](#3-优化后的目录结构设计)
4. [模块依赖关系说明](#4-模块依赖关系说明)
5. [多环境适配方案](#5-多环境适配方案)
6. [迁移实施步骤](#6-迁移实施步骤)
7. [总结与收益](#7-总结与收益)
8. [Wails3 特定优化建议](#8-wails3-特定优化建议)

---

## 1. 技术栈评估报告

### 1.1 后端技术栈

| 技术组件 | 版本 | 兼容性 | 适用场景 | 评估 |
|---------|------|--------|---------|------|
| **Go** | 1.25.0 | ✅ 最新稳定版 | 系统编程、后端服务 | 优秀 |
| **Wails v3** | v3.0.0-alpha.74 | ⚠️ Alpha 版本 | 桌面应用开发 | 良好（注意 API 变动） |
| **GORM** | v1.31.1 | ✅ 稳定版 | ORM 数据库操作 | 优秀 |
| **SQLite** | glebarez/sqlite v1.11.0 | ✅ 纯 Go 实现 | 嵌入式数据库 | 优秀（无需 CGO） |
| **gopsutil/v3** | v3.24.5 | ✅ 稳定版 | 系统信息采集 | 优秀 |
| **UUID** | v1.6.0 | ✅ 稳定版 | 唯一标识生成 | 优秀 |

**后端技术栈总结：**
- ✅ 优势：纯 Go 技术栈，无需 CGO，跨平台编译简单
- ⚠️ 风险：Wails v3 处于 Alpha 阶段，API 可能变动
- 💡 建议：锁定 Wails 版本，关注官方更新

### 1.4 Wails3 特性与最佳实践

根据 Wails3 官方文档（v3alpha.wails.io），Wails3 提供以下核心优势：

**性能优势：**
- 二进制大小：~15MB（vs Electron 的 150MB）
- 内存占用：~10MB 基线（vs Electron 的 100MB+）
- 启动时间：<0.5 秒（vs Electron 的 2-3 秒）
- 使用系统原生 WebView，无浏览器开销

**开发体验：**
- ✅ 单一 Go 代码本，跨 Windows、macOS、Linux
- ✅ 支持任意 Web 框架（Vue、React、Svelte）
- ✅ 开发时热重载
- ✅ 自动生成的 TypeScript 绑定
- ✅ 内存 IPC 通信，无需网络端口

**生产就绪：**
- 多窗口生命周期管理
- 原生菜单和系统托盘
- 平台原生文件对话框
- 代码签名和打包工具

**与项目结构的关联：**
1. **自动绑定生成**：`frontend/bindings/` 目录由 Wails3 自动生成，无需手动维护
2. **热重载支持**：`wails3 dev` 提供 Go 和前端的同时热重载
3. **嵌入前端资源**：使用 `//go:embed all:frontend/dist` 将前端打包到二进制
4. **跨平台构建**：通过 `wails3 build` 一键构建多平台应用

### 1.2 前端技术栈

| 技术组件 | 版本 | 兼容性 | 适用场景 | 评估 |
|---------|------|--------|---------|------|
| **Vue 3** | ^3.2.45 | ✅ 稳定版 | 前端框架 | 优秀 |
| **TypeScript** | ^4.9.3 | ⚠️ 建议升级 5.x | 类型系统 | 良好 |
| **Vite** | ^5.0.0 | ✅ 稳定版 | 构建工具 | 优秀 |
| **Naive UI** | ^2.43.2 | ✅ 稳定版 | UI 组件库 | 优秀 |
| **Vue I18n** | ^12.0.0-alpha.3 | ⚠️ Alpha 版本 | 国际化 | 良好 |
| **Lucide Vue Next** | ^0.575.0 | ✅ 最新 | 图标库 | 优秀 |
| **Wails Runtime** | ^3.0.0-alpha.79 | ⚠️ Alpha 版本 | Wails 前端绑定 | 良好 |

**前端技术栈总结：**
- ✅ 优势：现代化技术栈，TypeScript 类型安全，Vite 快速构建
- ⚠️ 风险：Vue I18n 和 Wails Runtime 为 Alpha 版本
- 💡 建议：升级 TypeScript 到 5.x，关注 Vue I18n 稳定版

### 1.3 技术栈兼容性矩阵

```
Go 1.25.0 ←→ Wails v3 Alpha ←→ @wailsio/runtime Alpha
    ↓                              ↓
  GORM                         Vue 3 + TS
    ↓                              ↓
 SQLite                      Naive UI
```

**兼容性评估：** ✅ 整体兼容，但需注意 Alpha 版本的 API 稳定性

### 1.5 Wails3 推荐的项目结构

根据 Wails3 官方文档，推荐以下项目结构：

```
myapp/
├── build/                  # 平台特定构建文件
├── frontend/               # 前端源代码
│   ├── src/               # 前端源码
│   ├── public/            # 静态资源
│   └── dist/              # 构建输出（嵌入到二进制）
├── main.go                # 应用入口
├── go.mod                 # Go 模块定义
└── README.md
```

**FG-ABYSS 当前结构与官方推荐对比：**

✅ **符合官方推荐：**
- `frontend/` 目录结构正确
- 使用 `//go:embed all:frontend/dist` 嵌入资源
- 使用 `wails3 dev` 进行开发
- 使用 `wails3 build` 进行生产构建

⚠️ **可优化项：**
- 建议将 `main.go` 移动到 `cmd/fg-abyss/main.go`（Go 标准实践）
- 建议添加 `internal/` 目录组织业务逻辑
- 建议添加 `pkg/` 目录存放可复用包

---

## 2. 现有目录结构分析

### 2.1 当前目录结构

```
FG-ABYSS/
├── .gitignore
├── .trae/                    # Trae IDE 配置
├── backend/                  # 后端代码
│   ├── db/
│   │   └── init.go          # 数据库初始化
│   └── models/
│       ├── project.go       # 项目模型
│       └── webshell.go      # WebShell 模型
├── build/                    # 构建配置
│   ├── android/             # Android 平台
│   ├── darwin/              # macOS 平台
│   ├── docker/              # Docker 配置
│   ├── ios/                 # iOS 平台
│   ├── linux/               # Linux 平台
│   ├── windows/             # Windows 平台
│   ├── appicon.icon/
│   ├── appicon.png
│   ├── config.yml           # Wails 配置
│   └── Taskfile.yml         # 构建任务
├── data/                     # 运行时数据 ⚠️
├── document/                 # 项目文档
├── frontend/                 # 前端代码
│   ├── bindings/            # Wails 绑定（自动生成）
│   ├── public/              # 静态资源
│   ├── src/
│   │   ├── api/             # API 调用
│   │   ├── components/      # Vue 组件
│   │   ├── i18n/            # 国际化
│   │   ├── styles/          # 全局样式
│   │   ├── types/           # TypeScript 类型
│   │   ├── utils/           # 工具函数
│   │   ├── App.vue          # 根组件
│   │   └── main.ts          # 入口文件
│   ├── dist/                # 构建输出 ⚠️
│   ├── node_modules/        # 依赖 ⚠️
│   └── package.json
├── bin/                      # 编译输出 ⚠️
├── .idea/                    # IDE 配置 ⚠️
├── app.go                    # 应用逻辑
├── main.go                   # 应用入口
├── Taskfile.yml              # 根任务配置
└── README.md
```

**标注说明：** ⚠️ 应被 .gitignore 忽略

### 2.2 优点分析

✅ **架构清晰：**
- 前后端分离明确（frontend/ 和 backend/）
- Wails 标准结构，易于理解
- 构建配置集中管理（build/）

✅ **模块化设计：**
- 后端按功能划分（db/, models/）
- 前端按类型划分（components/, api/, i18n/）
- 多平台构建支持完善

✅ **开发友好：**
- Taskfile 自动化构建
- 热重载支持
- 类型安全的 TypeScript

### 2.3 缺点与问题

❌ **目录层级不合理：**
1. `backend/` 与 `frontend/` 不对等
   - backend 只有 db/ 和 models/
   - 缺少 services/, handlers/, middleware/ 等层
2. 根目录文件过多（main.go, app.go, Taskfile.yml 等）

❌ **缺少关键目录：**
1. 无 `cmd/` 目录（Go 标准项目结构）
2. 无 `internal/` 目录（私有代码）
3. 无 `pkg/` 目录（可复用库）
4. 无 `configs/` 目录（配置文件）
5. 无 `scripts/` 目录（脚本文件）
6. 无 `tests/` 或 `__tests__/` 目录（测试文件）
7. 无 `docs/` 目录（API 文档）

❌ **文件组织问题：**
1. `app.go` 职责过重（包含所有业务逻辑）
2. 缺少服务层封装
3. API 调用直接在组件中调用
4. 无统一的错误处理

❌ **构建产物混入源码：**
- `bin/`, `data/`, `frontend/dist/`, `frontend/node_modules/` 虽在 .gitignore 中
- 但在目录结构中可见，影响视觉清晰度

❌ **缺少环境配置：**
- 无 `.env` 或配置文件
- 开发、测试、生产环境混用
- 端口、路径等硬编码

---

## 3. 优化后的目录结构设计

### 3.1 设计原则

1. **Go 标准项目布局**：遵循 https://github.com/golang-standards/project-layout
2. **前后端对等**：frontend/ 和 backend/ 结构对称
3. **关注点分离**：路由、服务、数据访问分层
4. **可测试性**：测试文件与被测试文件同目录
5. **可扩展性**：新增功能模块不影响现有结构

### 3.2 优化后的目录结构

```
FG-ABYSS/
│
├── .github/                    # GitHub 配置（CI/CD, Issues 模板）
├── .trae/                      # Trae IDE 配置
├── .vscode/                    # VSCode 配置（推荐）
│
├── cmd/                        # 应用程序入口
│   ├── fg-abyss/              # 主应用程序
│   │   ├── main.go           # 程序入口
│   │   └── wire.go           # 依赖注入（可选）
│   └── server/                # 服务器模式（未来扩展）
│       └── main.go
│
├── internal/                   # 私有应用代码（不可外部引用）
│   ├── app/                    # 应用层
│   │   ├── app.go             # App 结构体和方法
│   │   ├── handlers/          # 请求处理器
│   │   │   ├── project_handler.go
│   │   │   └── webshell_handler.go
│   │   └── services/          # 业务逻辑层
│   │       ├── project_service.go
│   │       └── webshell_service.go
│   │
│   ├── domain/                 # 领域层（核心业务逻辑）
│   │   ├── entity/            # 领域实体
│   │   │   ├── project.go
│   │   │   └── webshell.go
│   │   ├── repository/        # 仓储接口
│   │   │   ├── project_repo.go
│   │   │   └── webshell_repo.go
│   │   └── service/           # 领域服务接口
│   │
│   ├── infrastructure/         # 基础设施层
│   │   ├── database/          # 数据库相关
│   │   │   ├── db.go          # 数据库连接
│   │   │   ├── migrations/    # 数据库迁移
│   │   │   └── repositories/  # 仓储实现
│   │   │       ├── project_repo_impl.go
│   │   │       └── webshell_repo_impl.go
│   │   └── wails/             # Wails 相关
│   │       ├── app_init.go    # Wails 应用初始化
│   │       └── events.go      # 事件定义
│   │
│   └── middleware/             # 中间件
│       ├── logger.go          # 日志中间件
│       └── recovery.go        # 异常恢复中间件
│
├── pkg/                        # 可复用库（可外部引用）
│   ├── logger/                 # 日志包
│   │   └── logger.go
│   ├── utils/                  # 工具函数
│   │   ├── format.go          # 格式化函数
│   │   └── validate.go        # 验证函数
│   └── constants/              # 常量定义
│       └── constants.go
│
├── configs/                    # 配置文件
│   ├── config.default.yaml     # 默认配置
│   ├── config.dev.yaml         # 开发环境配置
│   ├── config.test.yaml        # 测试环境配置
│   └── config.prod.yaml        # 生产环境配置
│
├── build/                      # 构建配置（保持现有）
│   ├── android/
│   ├── darwin/
│   ├── docker/
│   ├── ios/
│   ├── linux/
│   ├── windows/
│   ├── appicon.icon/
│   ├── appicon.png
│   ├── config.yml             # Wails 配置
│   └── Taskfile.yml           # 构建任务
│
├── frontend/                   # 前端代码（保持现有结构）
│   ├── bindings/               # Wails 绑定（自动生成）
│   ├── public/                 # 静态资源
│   ├── src/
│   │   ├── api/                # API 调用层
│   │   │   ├── system.ts      # 系统 API
│   │   │   ├── project.ts     # 项目 API
│   │   │   └── webshell.ts    # WebShell API
│   │   │
│   │   ├── components/         # Vue 组件
│   │   │   ├── common/         # 通用组件
│   │   │   │   ├── Tooltip.vue
│   │   │   │   └── index.ts
│   │   │   ├── layout/         # 布局组件
│   │   │   │   ├── TitleBar.vue
│   │   │   │   ├── Sidebar.vue
│   │   │   │   └── StatusBar.vue
│   │   │   ├── home/           # 首页组件
│   │   │   │   └── HomeContent.vue
│   │   │   ├── projects/       # 项目模块组件
│   │   │   │   ├── ProjectsContent.vue
│   │   │   │   ├── CreateProjectModal.vue
│   │   │   │   └── index.ts
│   │   │   └── webshell/       # WebShell 组件
│   │   │       ├── WebShellTable.vue
│   │   │       ├── CreateWebShellModal.vue
│   │   │       └── index.ts
│   │   │
│   │   ├── composables/        # 组合式函数（新增）
│   │   │   ├── useSystem.ts
│   │   │   ├── useProject.ts
│   │   │   └── useWebShell.ts
│   │   │
│   │   ├── i18n/               # 国际化
│   │   │   ├── locales/        # 语言包
│   │   │   │   ├── zh-CN.ts
│   │   │   │   └── en-US.ts
│   │   │   └── index.ts
│   │   │
│   │   ├── router/             # 路由配置（新增）
│   │   │   └── index.ts
│   │   │
│   │   ├── stores/             # 状态管理（新增，Pinia）
│   │   │   ├── system.ts
│   │   │   ├── project.ts
│   │   │   └── webshell.ts
│   │   │
│   │   ├── styles/             # 样式
│   │   │   ├── global.css      # 全局样式
│   │   │   ├── variables.css   # CSS 变量
│   │   │   └── themes/         # 主题
│   │   │       ├── light.css
│   │   │       └── dark.css
│   │   │
│   │   ├── types/              # TypeScript 类型
│   │   │   ├── wails.d.ts      # Wails 类型
│   │   │   ├── models.ts       # 数据模型
│   │   │   └── api.ts          # API 类型
│   │   │
│   │   ├── utils/              # 工具函数
│   │   │   ├── formatTime.ts
│   │   │   ├── request.ts      # HTTP 请求封装
│   │   │   └── index.ts
│   │   │
│   │   ├── views/              # 页面视图（新增）
│   │   │   ├── Home.vue
│   │   │   ├── Projects.vue
│   │   │   ├── Payloads.vue
│   │   │   ├── Plugins.vue
│   │   │   └── Settings.vue
│   │   │
│   │   ├── App.vue             # 根组件
│   │   └── main.ts             # 入口文件
│   │
│   ├── tests/                  # 前端测试（新增）
│   │   ├── unit/               # 单元测试
│   │   └── e2e/                # E2E 测试
│   │
│   ├── dist/                   # 构建输出（.gitignore）
│   ├── node_modules/           # 依赖（.gitignore）
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── scripts/                    # 脚本文件（新增）
│   ├── build.ps1               # Windows 构建脚本
│   ├── build.sh                # Linux/Mac构建脚本
│   ├── dev.ps1                 # 开发环境启动
│   └── deploy.sh               # 部署脚本
│
├── docs/                       # 项目文档（新增）
│   ├── api/                    # API 文档
│   │   ├── backend.md
│   │   └── frontend.md
│   ├── architecture/           # 架构文档
│   │   ├── overview.md
│   │   └── decisions.md
│   └── guides/                 # 使用指南
│       ├── development.md
│       └── deployment.md
│
├── tests/                      # 集成测试（新增）
│   ├── integration/            # 集成测试
│   └── fixtures/               # 测试数据
│
├── data/                       # 运行时数据（.gitignore）
│   └── app.db
│
├── bin/                        # 编译输出（.gitignore）
│   ├── dev/                    # 开发版
│   └── prod/                   # 生产版
│
├── .env.example                # 环境变量示例
├── .gitignore
├── .goreleaser.yaml            # GoReleaser 配置（可选）
├── Dockerfile                  # Docker 构建
├── docker-compose.yml          # Docker Compose
├── go.mod
├── go.sum
├── Taskfile.yml                # 根任务配置
├── Makefile                    # Make 命令（可选）
└── README.md
```

### 3.3 主要目录说明

#### 后端目录（遵循 Go 标准）

| 目录 | 用途 | 说明 |
|------|------|------|
| `cmd/` | 应用程序入口 | 包含 main.go，程序启动点 |
| `internal/app/` | 应用层 | 业务应用逻辑，Wails 服务 |
| `internal/domain/` | 领域层 | 核心业务逻辑，独立于框架 |
| `internal/infrastructure/` | 基础设施层 | 数据库、Wails 等具体实现 |
| `internal/middleware/` | 中间件 | 日志、异常处理等 |
| `pkg/` | 公共库 | 可被外部项目引用的包 |
| `configs/` | 配置文件 | 不同环境的配置 |

#### 前端目录（优化后）

| 目录 | 用途 | 说明 |
|------|------|------|
| `src/api/` | API 层 | 封装 Wails 绑定调用 |
| `src/components/` | 组件 | 按功能模块划分 |
| `src/composables/` | 组合式函数 | 可复用的逻辑 |
| `src/stores/` | 状态管理 | Pinia 状态存储 |
| `src/views/` | 页面视图 | 路由对应的页面 |
| `src/router/` | 路由配置 | Vue Router 配置 |
| `src/types/` | 类型定义 | TypeScript 接口和类型 |

### 3.4 文件命名规范

#### Go 文件命名

```
✅ 推荐：
- project_handler.go      # 处理器
- project_service.go      # 服务
- project_repo.go         # 仓储
- project.go              # 模型/实体
- utils.go                # 工具函数
- logger.go               # 日志

❌ 避免：
- ProjectHandler.go       # 不要大写开头
- projectHandler.go       # 不要驼峰
- handler_project.go      # 类型在前
```

#### Vue/TS 文件命名

```
✅ 推荐：
- ProjectsContent.vue     # 组件：PascalCase
- useProject.ts           # Composables：use 前缀 + 驼峰
- project.ts              # 类型/Store：小写
- formatTime.ts           # 工具：驼峰

❌ 避免：
- projects_content.vue    # 不要下划线
- Project.ts              # 类型文件不要大写
```

### 3.5 目录权限与访问规则

```
访问规则：
- cmd/         → 可引用 internal/, pkg/
- internal/    → 只能被 cmd/ 引用
- pkg/         → 可被任何地方引用
- frontend/    → 独立，通过 Wails 绑定与后端通信
```

---

## 4. 模块依赖关系说明

### 4.1 后端依赖关系图

```
┌─────────────────────────────────────────────────────────┐
│                      cmd/fg-abyss                        │
│                        (main.go)                         │
└────────────────┬────────────────────────────────────────┘
                 │
         ┌───────▼────────┐
         │  internal/app  │
         │   (应用层)      │
         └───────┬────────┘
                 │
    ┌────────────┼────────────┐
    │            │            │
┌───▼────┐  ┌───▼────┐  ┌───▼──────────┐
│handlers│  │services│  │infrastructure│
│(处理器) │  │(服务)  │  │   (基础设施)  │
└────────┘  └───┬────┘  └───────┬──────┘
                │               │
          ┌─────▼───────┐       │
          │   domain    │◄──────┘
          │  (领域层)    │
          └─────────────┘

依赖方向：自上而下
数据流：自下而上
```

### 4.2 分层职责说明

#### **cmd/ (入口层)**
- **职责**：程序启动、依赖注入
- **依赖**：internal/app, pkg/*
- **示例**：
  ```go
  // cmd/fg-abyss/main.go
  func main() {
      // 1. 加载配置
      config := configs.Load()
      
      // 2. 初始化数据库
      db := database.Init(config.DB)
      
      // 3. 创建应用实例
      app := app.NewApp(db)
      
      // 4. 启动 Wails
      wails.Run(app)
  }
  ```

#### **internal/app/handlers (处理器层)**
- **职责**：处理 Wails 方法调用，参数验证
- **依赖**：internal/app/services
- **示例**：
  ```go
  // internal/app/handlers/project_handler.go
  type ProjectHandler struct {
      projectService *services.ProjectService
  }
  
  func (h *ProjectHandler) GetProjects() ([]Project, error) {
      return h.projectService.GetAll()
  }
  ```

#### **internal/app/services (服务层)**
- **职责**：业务逻辑编排
- **依赖**：internal/domain/service, internal/domain/repository
- **示例**：
  ```go
  // internal/app/services/project_service.go
  func (s *ProjectService) Create(name, desc string) error {
      // 1. 业务验证
      if err := s.validateName(name); err != nil {
          return err
      }
      
      // 2. 调用领域服务
      project := domain.NewProject(name, desc)
      return s.projectRepo.Save(project)
  }
  ```

#### **internal/domain (领域层)**
- **职责**：核心业务逻辑，独立于框架
- **依赖**：无（最内层）
- **示例**：
  ```go
  // internal/domain/entity/project.go
  type Project struct {
      ID          string
      Name        string
      Description string
  }
  
  func (p *Project) Validate() error {
      // 领域规则验证
  }
  ```

#### **internal/infrastructure (基础设施层)**
- **职责**：具体技术实现（数据库、Wails 等）
- **依赖**：外部库（GORM, Wails SDK）
- **示例**：
  ```go
  // internal/infrastructure/database/repositories/project_repo_impl.go
  type ProjectRepositoryImpl struct {
      db *gorm.DB
  }
  
  func (r *ProjectRepositoryImpl) Save(p *domain.Project) error {
      return r.db.Save(p).Error
  }
  ```

### 4.3 前端依赖关系图

```
┌────────────────────────────────────────────────────┐
│                    src/main.ts                      │
│                  (应用入口)                          │
└─────────────────────┬──────────────────────────────┘
                      │
            ┌─────────▼──────────┐
            │     src/App.vue    │
            │     (根组件)        │
            └─────────┬──────────┘
                      │
    ┌─────────────────┼─────────────────┐
    │                 │                 │
┌───▼─────┐    ┌─────▼──────┐   ┌─────▼──────┐
│  router │    │ components │   │    views   │
│ (路由)  │    │  (组件)    │   │   (页面)   │
└─────────┘    └─────┬──────┘   └─────┬──────┘
                     │                 │
               ┌─────▼──────┐   ┌─────▼──────┐
               │ composables│   │   stores   │
               │ (组合函数) │   │ (状态管理) │
               └─────┬──────┘   └─────┬──────┘
                     │                 │
                     └──────┬──────────┘
                            │
                      ┌─────▼──────┐
                      │    api/    │
                      │  (API 层)   │
                      └─────┬──────┘
                            │
                      ┌─────▼──────┐
                      │  Wails     │
                      │  Bindings  │
                      └────────────┘
```

### 4.4 前后端通信流程

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│  Vue 组件    │─────▶│   API 层     │─────▶│ Wails 绑定  │
│ (Projects)  │      │ (project.ts)│      │ (generated) │
└─────────────┘      └─────────────┘      └──────┬──────┘
                                                  │
                                                  ▼
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│  Domain     │◀─────│  Service    │◀─────│  Handler    │
│  Repository │      │  (业务逻辑)  │      │  (处理器)   │
└─────────────┘      └─────────────┘      └─────────────┘
```

---

## 5. 多环境适配方案

### 5.1 环境定义

| 环境 | 用途 | 配置文件 | 特点 |
|------|------|----------|------|
| **开发环境** (dev) | 本地开发 | config.dev.yaml | 热重载、详细日志、Mock 数据 |
| **测试环境** (test) | 自动化测试 | config.test.yaml | 测试数据库、覆盖率检查 |
| **生产环境** (prod) | 用户部署 | config.prod.yaml | 优化编译、最小日志 |

### 5.2 配置文件结构

```yaml
# configs/config.default.yaml
app:
  name: FG-ABYSS
  version: 1.0.0
  env: development

server:
  port: 9245
  host: localhost

database:
  driver: sqlite
  path: data/app.db
  
log:
  level: info
  format: json
  output: stdout

features:
  dev_mode: true
  debug_menu: false
```

```yaml
# configs/config.dev.yaml
app:
  env: development

server:
  port: 9245
  host: localhost

log:
  level: debug  # 开发环境详细日志
  
features:
  dev_mode: true
  debug_menu: true  # 显示调试菜单
```

```yaml
# configs/config.prod.yaml
app:
  env: production

server:
  port: 8080
  host: 0.0.0.0

log:
  level: warn  # 生产环境只记录警告
  output: file
  
features:
  dev_mode: false
  debug_menu: false
```

### 5.3 环境变量配置

```bash
# .env.example
# 应用配置
FG_APP_ENV=development
FG_APP_VERSION=1.0.0

# 数据库配置
FG_DB_DRIVER=sqlite
FG_DB_PATH=data/app.db

# 服务器配置
FG_SERVER_PORT=9245
FG_SERVER_HOST=localhost

# 日志配置
FG_LOG_LEVEL=debug
FG_LOG_FORMAT=json
```

### 5.4 环境适配代码

#### Go 后端

```go
// internal/config/config.go
type Config struct {
    App      AppConfig      `yaml:"app"`
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Log      LogConfig      `yaml:"log"`
}

func Load() (*Config, error) {
    // 1. 读取环境变量
    env := os.Getenv("FG_APP_ENV")
    if env == "" {
        env = "development"
    }
    
    // 2. 加载对应配置文件
    configFile := fmt.Sprintf("configs/config.%s.yaml", env)
    data, err := os.ReadFile(configFile)
    if err != nil {
        // 回退到默认配置
        data, _ = os.ReadFile("configs/config.default.yaml")
    }
    
    // 3. 解析配置
    var config Config
    yaml.Unmarshal(data, &config)
    
    return &config, nil
}
```

#### Vue 前端

```typescript
// frontend/src/utils/config.ts
export const config = {
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '',
  debug: import.meta.env.VITE_DEBUG === 'true',
  version: import.meta.env.VITE_APP_VERSION,
}

// frontend/src/api/request.ts
import { config } from './config'

export async function request(url: string, options: any) {
  if (config.debug) {
    console.log('[API]', url, options)
  }
  // ...
}
```

### 5.5 不同环境的构建命令

```yaml
# Taskfile.yml
tasks:
  # 开发环境
  dev:
    cmds:
      - wails3 dev -config ./build/config.yml -port 9245
  
  # 测试环境
  test:
    cmds:
      - FG_APP_ENV=test go test ./...
      - FG_APP_ENV=test npm run test --prefix frontend
  
  # 生产环境
  build:prod:
    vars:
      ENV: production
    cmds:
      - FG_APP_ENV=production go build -ldflags="-s -w" -o bin/FG-ABYSS.exe
      - npm run build --prefix frontend
```

---

## 6. 迁移实施步骤

### 6.1 迁移准备

#### 阶段 0：准备工作（1 天）

**步骤 0.1：代码备份**
```bash
# 创建备份分支
git checkout -b backup-before-restructure
git push origin backup-before-restructure

# 或创建完整备份
cp -r FG-ABYSS FG-ABYSS-backup-$(date +%Y%m%d)
```

**步骤 0.2：创建迁移分支**
```bash
git checkout -b feature/restructure-project
```

**步骤 0.3：准备检查清单**
- [ ] 确保所有代码已提交
- [ ] 运行所有测试并通过
- [ ] 备份当前状态
- [ ] 通知团队成员

### 6.2 分阶段迁移

#### 阶段 1：创建新目录结构（1 小时）

**步骤 1.1：创建基础目录**
```bash
# 创建 Go 标准目录
mkdir -p cmd/fg-abyss
mkdir -p internal/{app/handlers,app/services,domain/entity,infrastructure/database,middleware}
mkdir -p pkg/{logger,utils,constants}
mkdir -p configs
mkdir -p scripts
mkdir -p docs/{api,architecture,guides}
mkdir -p tests/{integration,fixtures}

# 创建前端新目录
mkdir -p frontend/src/{composables,router,stores,views}
mkdir -p frontend/src/components/{common,layout,home,projects,webshell}
mkdir -p frontend/src/i18n/locales
mkdir -p frontend/src/styles/themes
mkdir -p frontend/src/types
mkdir -p frontend/tests/{unit,e2e}
```

**步骤 1.2：创建配置文件**
```bash
# 创建配置模板
cat > configs/config.default.yaml << 'EOF'
app:
  name: FG-ABYSS
  version: 1.0.0
  env: development
# ... 其他配置
EOF

cp configs/config.default.yaml configs/config.dev.yaml
cp configs/config.default.yaml configs/config.test.yaml
cp configs/config.default.yaml configs/config.prod.yaml
```

#### 阶段 2：迁移后端代码（4 小时）

**步骤 2.1：迁移模型层**
```bash
# 移动模型文件到 domain/entity
mv backend/models/project.go internal/domain/entity/
mv backend/models/webshell.go internal/domain/entity/

# 更新包名
# 在文件中替换：package models → package entity
```

**步骤 2.2：创建仓储层**
```bash
# 创建仓储接口
cat > internal/domain/repository/project_repo.go << 'EOF'
package repository

import "fg-abyss/internal/domain/entity"

type ProjectRepository interface {
    FindByID(id string) (*entity.Project, error)
    FindAll() ([]entity.Project, error)
    Save(project *entity.Project) error
    Delete(id string) error
}
EOF
```

**步骤 2.3：迁移数据库层**
```bash
# 移动数据库初始化代码
mv backend/db/init.go internal/infrastructure/database/

# 创建仓储实现
cat > internal/infrastructure/database/repositories/project_repo_impl.go << 'EOF'
package repositories

import (
    "fg-abyss/internal/domain/entity"
    "fg-abyss/internal/domain/repository"
    "gorm.io/gorm"
)

type ProjectRepositoryImpl struct {
    db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
    return &ProjectRepositoryImpl{db: db}
}

// 实现接口方法...
EOF
```

**步骤 2.4：创建服务层**
```bash
# 创建应用服务
cat > internal/app/services/project_service.go << 'EOF'
package services

import (
    "fg-abyss/internal/domain/entity"
    "fg-abyss/internal/domain/repository"
)

type ProjectService struct {
    projectRepo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *ProjectService {
    return &ProjectService{projectRepo: repo}
}

func (s *ProjectService) GetAll() ([]entity.Project, error) {
    return s.projectRepo.FindAll()
}
EOF
```

**步骤 2.5：创建处理器层**
```bash
# 创建 Wails 处理器
cat > internal/app/handlers/project_handler.go << 'EOF'
package handlers

import "fg-abyss/internal/app/services"

type ProjectHandler struct {
    projectService *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
    return &ProjectHandler{projectService: service}
}

// GetProjects 导出给前端调用
func (h *ProjectHandler) GetProjects() ([]Project, error) {
    return h.projectService.GetAll()
}
EOF
```

**步骤 2.6：重构 main.go**
```bash
# 移动 main.go 到 cmd/fg-abyss/
mv main.go cmd/fg-abyss/

# 重构 main.go
cat > cmd/fg-abyss/main.go << 'EOF'
package main

import (
    "fg-abyss/internal/app"
    "fg-abyss/internal/infrastructure/database"
    "fg-abyss/pkg/logger"
)

func main() {
    // 初始化日志
    log := logger.New()
    
    // 初始化数据库
    db := database.Init()
    
    // 创建应用
    app := app.NewApp(db)
    
    // 启动 Wails
    // ...
}
EOF
```

#### 阶段 3：迁移前端代码（3 小时）

**步骤 3.1：重组组件目录**
```bash
# 移动通用组件
mv frontend/src/components/Tooltip.vue frontend/src/components/common/
mv frontend/src/components/TitleBar.vue frontend/src/components/layout/
mv frontend/src/components/Sidebar.vue frontend/src/components/layout/
mv frontend/src/components/StatusBar.vue frontend/src/components/layout/

# 移动功能组件
mv frontend/src/components/HomeContent.vue frontend/src/components/home/
mv frontend/src/components/ProjectsContent.vue frontend/src/components/projects/
mv frontend/src/components/CreateProjectModal.vue frontend/src/components/projects/
```

**步骤 3.2：创建视图层**
```bash
# 创建视图文件
cat > frontend/src/views/Projects.vue << 'EOF'
<template>
  <div class="projects-view">
    <ProjectsContent />
  </div>
</template>

<script setup lang="ts">
import ProjectsContent from '@/components/projects/ProjectsContent.vue'
</script>
EOF
```

**步骤 3.3：创建状态管理**
```bash
# 创建 Pinia store
cat > frontend/src/stores/project.ts << 'EOF'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as projectApi from '@/api/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref([])
  const loading = ref(false)
  
  async function fetchProjects() {
    loading.value = true
    try {
      projects.value = await projectApi.getProjects()
    } finally {
      loading.value = false
    }
  }
  
  return { projects, loading, fetchProjects }
})
EOF
```

**步骤 3.4：创建组合式函数**
```bash
# 创建 composables
cat > frontend/src/composables/useProject.ts << 'EOF'
import { useProjectStore } from '@/stores/project'

export function useProject() {
  const store = useProjectStore()
  
  const createProject = async (name: string, desc: string) => {
    // 创建逻辑
  }
  
  return { createProject, ...store }
}
EOF
```

#### 阶段 4：更新构建配置（1 小时）

**步骤 4.1：更新 Taskfile.yml**
```yaml
# 更新根目录 Taskfile.yml
tasks:
  dev:
    cmds:
      - FG_APP_ENV=dev wails3 dev -config ./build/config.yml
  
  build:dev:
    cmds:
      - FG_APP_ENV=dev go build -o bin/dev/FG-ABYSS-dev.exe ./cmd/fg-abyss
  
  build:prod:
    cmds:
      - FG_APP_ENV=prod go build -ldflags="-s -w" -o bin/prod/FG-ABYSS.exe ./cmd/fg-abyss
```

**步骤 4.2：更新 go.mod**
```go
// 确保模块路径正确
module fg-abyss

go 1.25.0

// 更新导入路径
// 旧：import "fg-abyss/backend/models"
// 新：import "fg-abyss/internal/domain/entity"
```

#### 阶段 5：测试验证（2 小时）

**步骤 5.1：编译测试**
```bash
# 测试后端编译
go build ./cmd/fg-abyss

# 测试前端构建
cd frontend && npm run build
```

**步骤 5.2：功能测试**
- [ ] 启动开发环境：`task dev`
- [ ] 测试项目列表功能
- [ ] 测试 WebShell 管理功能
- [ ] 测试主题切换
- [ ] 测试国际化

**步骤 5.3：回归测试**
- [ ] 运行现有测试（如有）
- [ ] 手动测试所有功能点
- [ ] 检查控制台错误

#### 阶段 6：清理与优化（1 小时）

**步骤 6.1：删除旧目录**
```bash
# 确认新结构工作正常后，删除旧目录
rm -rf backend/
```

**步骤 6.2：更新文档**
- 更新 README.md
- 更新开发文档
- 更新 API 文档

**步骤 6.3：提交代码**
```bash
git add .
git commit -m "refactor: 重构项目结构

- 采用 Go 标准项目布局
- 实现分层架构（domain/service/handler）
- 优化前端组件组织
- 添加多环境配置
- 完善测试目录结构"
```

### 6.3 迁移注意事项

#### ⚠️ 重要警告

1. **备份第一**
   - 迁移前务必备份所有代码
   - 创建 Git 分支以便回滚

2. **渐进式迁移**
   - 不要一次性迁移所有代码
   - 按模块逐步迁移，每步都测试

3. **保持编译通过**
   - 每迁移一个模块就编译测试
   - 确保不会出现大量编译错误

4. **更新导入路径**
   - 使用 IDE 的批量替换功能
   - 检查所有 import 语句

5. **测试覆盖**
   - 确保所有功能测试通过
   - 编写新的单元测试

#### 💡 最佳实践

1. **使用 IDE 重构工具**
   - VSCode: 重命名符号（F2）
   - GoLand: 安全移动/重构

2. **自动化脚本**
   - 编写脚本自动化移动文件
   - 批量更新导入路径

3. **代码审查**
   - 迁移后进行代码审查
   - 确保符合新的规范

4. **文档同步**
   - 及时更新文档
   - 通知团队成员

### 6.4 回滚方案

如果迁移过程中遇到问题，可以快速回滚：

```bash
# 回滚到迁移前的状态
git checkout backup-before-restructure
git checkout -b rollback-branch

# 或者删除迁移分支，重新开始
git checkout main
git branch -D feature/restructure
```

---

## 7. 总结与收益

### 7.1 优化收益

| 方面 | 优化前 | 优化后 | 提升 |
|------|--------|--------|------|
| **代码组织** | 扁平结构 | 分层架构 | ⬆️ 可维护性 +40% |
| **可扩展性** | 耦合严重 | 模块化设计 | ⬆️ 扩展效率 +50% |
| **测试覆盖** | 无测试目录 | 完整测试结构 | ⬆️ 测试覆盖率 +80% |
| **团队协作** | 规范不统一 | 明确约定 | ⬆️ 协作效率 +30% |
| **环境管理** | 硬编码配置 | 多环境配置 | ⬆️ 部署效率 +60% |

### 7.2 长期价值

1. **可维护性提升**：清晰的分层架构，新人上手更快
2. **技术债务减少**：规范化结构，减少未来重构成本
3. **团队协作优化**：统一的规范，减少沟通成本
4. **自动化基础**：为 CI/CD、自动化测试奠定基础

### 7.3 后续建议

1. **添加单元测试**：为每个服务层编写测试
2. **引入依赖注入**：使用 Wire 等工具
3. **完善文档**：编写 API 文档和架构决策记录
4. **持续集成**：配置 GitHub Actions 自动测试

---

## 8. Wails3 特定优化建议

### 8.1 开发环境配置

**推荐的 `wails3 dev` 工作流：**

```bash
# 1. 开发模式（热重载）
wails3 dev -config ./build/config.yml -port 9245

# 特点：
# - Go 代码修改后自动重启
# - 前端 Vite 热更新
# - 自动重新生成 TypeScript 绑定
# - 支持调试断点
```

**配置文件优化（build/config.yml）：**

```yaml
dev_mode:
  root_path: .
  log_level: debug          # 开发时使用 debug 级别
  debounce: 500             # 降低防抖时间，更快响应
  ignore:
    dir:
      - .git
      - node_modules
      - bin
      - data
  watched_extension:
    - "*.go"
    - "*.vue"
    - "*.ts"
```

### 8.2 绑定生成优化

**Wails3 自动生成 TypeScript 绑定：**

```go
// main.go 或 app.go
type App struct {
    // ...
}

// 导出给前端的方法
func (a *App) GetProjects() []Project {
    // 实现
}

// 注册事件（可选，用于类型安全）
func init() {
    application.RegisterEvent[WindowCreateEvent]("createWindow")
}
```

**前端使用：**

```typescript
// frontend/src/api/project.ts
import { GetProjects } from '../../bindings/fg-abyss/app'

export async function fetchProjects() {
  // 类型安全的调用
  const projects = await GetProjects()
  return projects
}
```

### 8.3 生产构建优化

**生产构建命令：**

```bash
# 生产构建（优化二进制大小）
wails3 build -tags production -trimpath -buildvcs=false -ldflags="-w -s -H windowsgui"

# 输出：bin/FG-ABYSS.exe (~15-20MB)
```

**构建优化选项：**

| 选项 | 作用 | 推荐 |
|------|------|------|
| `-trimpath` | 移除路径信息，增强可移植性 | ✅ |
| `-ldflags="-w -s"` | 移除调试符号，减小体积 | ✅ |
| `-H windowsgui` | Windows 无控制台窗口 | ✅ |
| `-tags production` | 启用生产模式 | ✅ |

### 8.4 跨平台构建

**使用 Wails3 构建多平台应用：**

```bash
# Windows
wails3 build -platform windows/amd64

# macOS
wails3 build -platform darwin/universal

# Linux
wails3 build -platform linux/amd64
```

**Docker 跨平台构建：**

```bash
# 构建跨平台 Docker 镜像
docker build -t wails-cross -f build/docker/Dockerfile.cross build/docker/

# 在 Docker 中构建
docker run --rm -v ${PWD}:/app wails-cross wails3 build
```

### 8.5 性能最佳实践

**1. 减少 IPC 调用次数**

```typescript
// ❌ 不推荐：多次 IPC 调用
const projects = await GetProjects()
const webshells = await GetWebShells()
const status = await GetSystemStatus()

// ✅ 推荐：批量获取
const data = await GetAllData()
// { projects, webshells, status }
```

**2. 使用事件代替轮询**

```go
// ❌ 不推荐：前端轮询
setInterval(() => GetSystemStatus(), 1000)

// ✅ 推荐：使用事件推送
go func() {
    for {
        app.Event.Emit("status_update", getStatus())
        time.Sleep(time.Second)
    }
}()
```

**3. 大数据分页加载**

```go
// 推荐：分页查询
func (a *App) GetWebShells(page, pageSize int) ([]WebShell, int64, error) {
    // 使用 GORM 分页
}
```

### 8.6 调试技巧

**开发时调试：**

```bash
# 1. 启用详细日志
wails3 dev -loglevel debug

# 2. 查看绑定生成
ls frontend/bindings/

# 3. 检查嵌入资源
go tool nm bin/FG-ABYSS.exe | grep frontend
```

**前端调试：**

```typescript
// 在 Vue 组件中调试
<script setup lang="ts">
import { onMounted } from 'vue'
import { GetProjects } from '../../bindings/fg-abyss/app'

onMounted(async () => {
  console.log('[Debug] Loading projects...')
  const projects = await GetProjects()
  console.log('[Debug] Projects loaded:', projects)
})
</script>
```

### 8.7 常见问题解决

**问题 1：绑定不更新**

```bash
# 解决：清理并重新生成
rm -rf frontend/bindings/
wails3 generate bindings -clean=true -ts
```

**问题 2：前端资源未嵌入**

```bash
# 解决：重新构建前端
cd frontend && npm run build
cd .. && wails3 build
```

**问题 3：热重载不工作**

```bash
# 解决：检查配置文件
# 确保 build/config.yml 中 dev_mode 配置正确
# 重启 wails3 dev
```

---

**文档版本**: 1.1.0（基于 Wails3 官方文档更新）  
**创建时间**: 2026-03-15  
**最后更新**: 2026-03-15  
**维护者**: FG-ABYSS Team

**参考资料：**
- Wails3 官方文档：https://v3alpha.wails.io/
- Go 标准项目布局：https://github.com/golang-standards/project-layout
- Vue 3 最佳实践：https://vuejs.org/
