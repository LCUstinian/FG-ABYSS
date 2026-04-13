# Tasks

## 阶段 1: 项目初始化 (P0)

- [x] Task 1.1: 初始化 Tauri V2 + Vue 3 项目结构
  - [ ] 创建 package.json 配置（依赖：vue, vue-router, pinia, naive-ui, vue-i18n）
  - [ ] 创建 src-tauri/Cargo.toml 配置（依赖：tauri, tokio, reqwest, rusqlite, ring）
  - [ ] 创建 TypeScript 配置（tsconfig.json）
  - [ ] 创建 Vite 配置（vite.config.ts）
  - [ ] 创建 Tauri 配置（tauri.conf.json）
  - [ ] 创建基础目录结构（src/, src-tauri/src/）

- [ ] Task 1.2: 配置开发环境和 Git
  - [ ] 初始化 Git 仓库（使用当前目录）
  - [ ] 创建 .gitignore 文件
  - [ ] 创建 .env.example 文件
  - [ ] 安装前端依赖（pnpm install）
  - [ ] 安装 Rust 依赖（cargo fetch）
  - [ ] 创建首个 Git 提交（Initial commit）

- [ ] Task 1.3: 创建应用入口和基础组件
  - [ ] 创建 src/main.ts（Vue 应用入口）
  - [ ] 创建 src/App.vue（根组件）
  - [ ] 创建 src/index.html（HTML 入口）
  - [ ] 创建 src-tauri/src/main.rs（Tauri 入口）
  - [ ] 创建 src-tauri/src/lib.rs（核心逻辑入口）
  - [ ] 验证开发服务器启动（pnpm tauri dev）

## 阶段 2: 核心架构实现 (P0)

- [ ] Task 2.1: 实现布局系统
  - [ ] 创建 CustomTitlebar.vue（自定义标题栏，含 Logo、应用名、主题切换、窗口控制）
  - [ ] 创建 NavigationMenu.vue（左侧导航菜单，5 个菜单项）
  - [ ] 创建 StatusBar.vue（底部状态栏，显示连接状态、项目数、载荷数）
  - [ ] 创建 MainLayout.vue（主布局组件，整合所有布局组件）
  - [ ] 实现主题系统（亮色/暗色/跟随系统）
  - [ ] 实现语言切换（中文/英文）

- [ ] Task 2.2: 实现路由系统
  - [ ] 创建 src/router/index.ts（路由配置）
  - [ ] 创建路由定义（/, /project, /payload, /plugin, /settings）
  - [ ] 创建路由守卫
  - [ ] 实现页面切换动画

- [ ] Task 2.3: 实现状态管理（Pinia）
  - [ ] 创建 src/stores/index.ts（Pinia 初始化）
  - [ ] 创建 app.ts（应用全局状态）
  - [ ] 创建 payload.ts（载荷配置状态）
  - [ ] 创建 project.ts（项目状态）
  - [ ] 创建 webshell.ts（WebShell 连接状态）

- [ ] Task 2.4: 实现国际化（i18n）
  - [ ] 创建 src/i18n/index.ts（i18n 配置）
  - [ ] 创建 zh-CN.ts（中文语言包）
  - [ ] 创建 en-US.ts（英文语言包）
  - [ ] 集成到 Vue 应用

## 阶段 3: 数据库基础设施 (P0)

- [ ] Task 3.1: 实现 SQLite 数据库层
  - [ ] 创建 src-tauri/infra/database/mod.rs（数据库模块）
  - [ ] 创建连接池管理（pool.rs）
  - [ ] 创建数据模型（models/mod.rs）
  - [ ] 创建数据库迁移脚本（migrations/001_create_tables.sql）
  - [ ] 实现数据库初始化逻辑

- [ ] Task 3.2: 实现数据模型
  - [ ] 创建 PayloadConfig 模型（载荷配置）
  - [ ] 创建 Project 模型（项目）
  - [ ] 创建 WebShell 模型（WebShell）
  - [ ] 创建 Plugin 模型（插件）
  - [ ] 实现 CRUD 操作 trait

## 阶段 4: 载荷管理模块 (P0) - 本地功能

- [ ] Task 4.1: 实现载荷生成核心逻辑
  - [ ] 创建 src-tauri/core/payload/mod.rs（载荷模块）
  - [ ] 创建 generator.rs（载荷生成器）
  - [ ] 创建 php.rs（PHP 载荷生成）
  - [ ] 创建 jsp.rs（JSP 载荷生成）
  - [ ] 创建 asp.rs（ASP 载荷生成）
  - [ ] 创建 aspx.rs（ASPX 载荷生成）
  - [ ] 创建 template.rs（模板管理）

- [ ] Task 4.2: 实现载荷管理 Commands
  - [ ] 创建 src-tauri/commands/payload.rs
  - [ ] 实现 create_payload_config 命令
  - [ ] 实现 update_payload_config 命令
  - [ ] 实现 delete_payload_config 命令
  - [ ] 实现 get_all_payload_configs 命令
  - [ ] 实现 generate_payload 命令

- [ ] Task 4.3: 实现载荷管理前端组件
  - [ ] 创建 PayloadView.vue（载荷管理页面）
  - [ ] 创建 PayloadConfig.vue（载荷配置表单）
  - [ ] 创建 PayloadPreview.vue（代码预览）
  - [ ] 创建 PayloadHistory.vue（历史记录）
  - [ ] 创建 PayloadTemplate.vue（模板管理）
  - [ ] 创建 LanguageSelector.vue（脚本语言选择）
  - [ ] 创建 EncryptionSelector.vue（加密算法选择）
  - [ ] 创建 ObfuscationSlider.vue（混淆强度滑块）

- [ ] Task 4.4: 实现载荷管理 API 层
  - [ ] 创建 src/api/payload.ts（载荷 API）
  - [ ] 集成 invoke 调用后端命令
  - [ ] 实现错误处理
  - [ ] 实现类型定义（src/types/payload.ts）

- [ ] Task 4.5: 实现搜索/过滤/排序功能
  - [ ] 实现搜索框组件（通用搜索算法）
  - [ ] 实现过滤面板（按类型、标签过滤）
  - [ ] 实现排序功能（多字段排序）
  - [ ] 集成到 Pinia store

## 阶段 5: 项目管理模块 (P0) - 本地数据管理

- [ ] Task 5.1: 实现项目管理核心逻辑
  - [ ] 创建 src-tauri/core/project/mod.rs（项目模块）
  - [ ] 创建 manager.rs（项目管理器）
  - [ ] 创建 entity.rs（项目实体）
  - [ ] 创建 recycle_bin.rs（回收站）

- [ ] Task 5.2: 实现 WebShell 连接模块
  - [ ] 创建 src-tauri/core/webshell/mod.rs（WebShell 模块）
  - [ ] 创建 connection.rs（连接管理）
  - [ ] 创建 session.rs（会话管理）
  - [ ] 创建 health_check.rs（健康检查）

- [ ] Task 5.3: 实现项目管理 Commands
  - [ ] 创建 src-tauri/commands/project.rs
  - [ ] 实现 create_project 命令
  - [ ] 实现 update_project 命令
  - [ ] 实现 delete_project 命令（软删除）
  - [ ] 实现 get_all_projects 命令
  - [ ] 实现 restore_project 命令（从回收站恢复）
  - [ ] 创建 src-tauri/commands/webshell.rs
  - [ ] 实现 create_webshell 命令
  - [ ] 实现 update_webshell 命令
  - [ ] 实现 delete_webshell 命令（软删除）
  - [ ] 实现 test_connection 命令
  - [ ] 实现 connect_webshell 命令

- [ ] Task 5.4: 实现项目管理前端组件
  - [ ] 创建 ProjectView.vue（项目管理页面）
  - [ ] 创建 ProjectTree.vue（项目树组件）
  - [ ] 创建 ProjectSidebar.vue（左侧子栏，含新建、回收站）
  - [ ] 创建 WebShellList.vue（WebShell 列表表格）
  - [ ] 创建 WebShellForm.vue（WebShell 创建/编辑表单）
  - [ ] 创建 WebShellContextMenu.vue（右键菜单）
  - [ ] 创建 RecycleBin.vue（回收站视图）
  - [ ] 创建 Pagination.vue（分页器）

- [ ] Task 5.5: 实现项目管理 API 层
  - [ ] 创建 src/api/project.ts（项目 API）
  - [ ] 创建 src/api/webshell.ts（WebShell API）
  - [ ] 集成 invoke 调用后端命令
  - [ ] 实现类型定义（src/types/project.ts, src/types/webshell.ts）

## 阶段 6: 加密通信模块 (P0) - 远程功能

- [ ] Task 6.1: 实现加密核心逻辑
  - [ ] 创建 src-tauri/core/crypto/mod.rs（加密模块）
  - [ ] 创建 aes.rs（AES-256-GCM 加密）
  - [ ] 创建 xor.rs（XOR 加密）
  - [ ] 创建 key_manager.rs（密钥管理）
  - [ ] 实现 Argon2 密钥派生

- [ ] Task 6.2: 实现 HTTP 客户端
  - [ ] 创建 src-tauri/core/http/mod.rs（HTTP 模块）
  - [ ] 创建 client.rs（HTTP 客户端，使用 reqwest）
  - [ ] 创建 proxy.rs（代理支持）
  - [ ] 创建 rate_limiter.rs（速率限制）
  - [ ] 创建 retry.rs（重试机制）

- [ ] Task 6.3: 实现加密通信 Commands
  - [ ] 创建 src-tauri/commands/console.rs
  - [ ] 实现 execute_command 命令（执行远程命令）
  - [ ] 实现 list_files 命令（列出远程文件）
  - [ ] 实现 upload_file 命令（上传文件）
  - [ ] 实现 download_file 命令（下载文件）
  - [ ] 实现 execute_query 命令（执行 SQL 查询）

## 阶段 7: 控制台窗口系统 (P0) - 远程功能

- [ ] Task 7.1: 实现控制台窗口组件
  - [ ] 创建 ConsoleWindow.vue（控制台窗口，使用 Tauri window API）
  - [ ] 创建 ConsoleTabs.vue（TAB 栏组件，支持拖拽排序）
  - [ ] 实现窗口管理逻辑（打开、关闭、最小化、最大化）
  - [ ] 实现独立会话管理

- [ ] Task 7.2: 实现文件管理插件
  - [ ] 创建 FileBrowser.vue（文件浏览器，双栏布局）
  - [ ] 创建 FileUploader.vue（文件上传，支持拖拽）
  - [ ] 创建 FileEditor.vue（文件编辑器）
  - [ ] 创建 FileContextMenu.vue（文件右键菜单）
  - [ ] 集成后端文件操作命令

- [ ] Task 7.3: 实现数据库管理插件
  - [ ] 创建 DatabaseManager.vue（数据库管理器）
  - [ ] 创建 QueryEditor.vue（SQL 查询编辑器）
  - [ ] 创建 ResultTable.vue（查询结果表格）
  - [ ] 创建 DatabaseSelector.vue（数据库选择器）
  - [ ] 集成后端数据库操作命令

- [ ] Task 7.4: 实现终端管理插件
  - [ ] 创建 TerminalView.vue（终端视图，使用 xterm.js）
  - [ ] 创建 TerminalInput.vue（命令输入框）
  - [ ] 创建 CommandHistory.vue（命令历史）
  - [ ] 集成后端命令执行命令

## 阶段 8: 设置页面 (P1)

- [ ] Task 8.1: 实现设置页面布局
  - [ ] 创建 SettingsView.vue（设置页面）
  - [ ] 创建 SettingsLayout.vue（左侧子菜单 + 右侧内容）
  - [ ] 实现设置路由（/settings/appearance, /settings/connection, /settings/security, /settings/about）

- [ ] Task 8.2: 实现外观设置
  - [ ] 创建 AppearanceSettings.vue（外观设置页面）
  - [ ] 创建 ThemeSelector.vue（主题选择器）
  - [ ] 创建 LanguageSelector.vue（语言选择器）
  - [ ] 创建 ColorPicker.vue（强调色选择器）

- [ ] Task 8.3: 实现连接设置
  - [ ] 创建 ConnectionSettings.vue（连接设置页面）
  - [ ] 实现请求配置（超时、重试、UA 池）
  - [ ] 实现代理配置（HTTP/SOCKS5）

- [ ] Task 8.4: 实现安全设置
  - [ ] 创建 SecuritySettings.vue（安全设置页面）
  - [ ] 实现加密算法选择
  - [ ] 实现密钥管理
  - [ ] 实现数据安全选项

- [ ] Task 8.5: 实现关于页面
  - [ ] 创建 AboutPage.vue（关于页面）
  - [ ] 实现版本信息卡片
  - [ ] 实现检查更新功能

## 阶段 9: 插件系统基础 (P2)

- [ ] Task 9.1: 实现插件加载器
  - [ ] 创建 src-tauri/plugins/mod.rs（插件模块）
  - [ ] 创建 loader.rs（插件加载器）
  - [ ] 创建 sandbox.rs（沙箱隔离）
  - [ ] 创建 runtime.rs（运行时）

- [ ] Task 9.2: 实现插件管理前端
  - [ ] 创建 PluginView.vue（插件管理页面）
  - [ ] 创建 PluginList.vue（插件列表）
  - [ ] 创建 PluginManager.vue（插件管理）
  - [ ] 实现插件安装/卸载/启用/禁用功能

## 阶段 10: 通用功能和优化 (P1)

- [ ] Task 10.1: 实现通用组件
  - [ ] 创建 Loading.vue（加载组件）
  - [ ] 创建 EmptyState.vue（空状态组件）
  - [ ] 创建 ConfirmDialog.vue（确认对话框）
  - [ ] 创建 ErrorDisplay.vue（错误展示）
  - [ ] 创建 CodeEditor.vue（通用代码编辑器）
  - [ ] 创建 SearchBar.vue（通用搜索栏）

- [ ] Task 10.2: 实现错误处理系统
  - [ ] 创建 src-tauri/error.rs（错误处理）
  - [ ] 实现统一错误码系统（60+ 错误码，10 大分类）
  - [ ] 实现错误严重性分级
  - [ ] 实现用户友好的错误提示

- [ ] Task 10.3: 实现日志系统
  - [ ] 创建 src-tauri/logging.rs（日志系统）
  - [ ] 实现结构化日志
  - [ ] 实现多日志级别
  - [ ] 实现敏感信息脱敏
  - [ ] 实现日志轮转

- [ ] Task 10.4: 实现性能优化
  - [ ] 实现组件懒加载
  - [ ] 实现虚拟列表（用于大数据列表）
  - [ ] 实现防抖和节流（搜索、滚动）
  - [ ] 优化内存占用

## 阶段 11: 测试和文档 (P1)

- [ ] Task 11.1: 编写单元测试
  - [ ] 编写 Rust 单元测试（core/ 模块）
  - [ ] 编写前端单元测试（composables/）
  - [ ] 确保测试覆盖率 > 90%

- [ ] Task 11.2: 编写集成测试
  - [ ] 编写载荷生成集成测试
  - [ ] 编写加密通信集成测试
  - [ ] 编写项目管理集成测试

- [ ] Task 11.3: 编写用户文档
  - [ ] 编写用户手册
  - [ ] 编写安装部署文档
  - [ ] 更新 README.md

## 阶段 12: 构建和发布 (P0)

- [ ] Task 12.1: 配置构建流程
  - [ ] 创建构建脚本（scripts/build.ps1, scripts/build.sh）
  - [ ] 配置生产环境构建
  - [ ] 配置代码签名（可选）

- [ ] Task 12.2: 执行最终测试
  - [ ] 执行全量测试
  - [ ] 修复所有 Critical/High Bug
  - [ ] 性能测试和优化

- [ ] Task 12.3: 发布 v0.1.0
  - [ ] 打 Git 标签（v0.1.0）
  - [ ] 创建 Release
  - [ ] 发布安装包

# Task Dependencies

- Task 1.x 是所有任务的基础，必须首先完成
- Task 2.x 依赖于 Task 1.x 完成
- Task 3.x 依赖于 Task 2.x 完成（需要路由和状态管理）
- Task 4.x 和 Task 5.x 可以并行开发（都依赖于 Task 2.x 和 Task 3.x）
- Task 6.x 依赖于 Task 5.x（需要 WebShell 连接）
- Task 7.x 依赖于 Task 6.x（需要加密通信）
- Task 8.x 可以独立开发（依赖于 Task 2.x 布局系统）
- Task 9.x 可以独立开发（P2 优先级）
- Task 10.x 可以并行开发（依赖于各功能模块）
- Task 11.x 在所有功能完成后进行
- Task 12.x 是最后阶段
