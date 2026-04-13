# Checklist

## 阶段 1: 项目初始化

- [ ] 1.1.1: package.json 包含所有必需依赖（vue, vue-router, pinia, naive-ui, vue-i18n）
- [ ] 1.1.2: src-tauri/Cargo.toml 包含所有必需依赖（tauri, tokio, reqwest, rusqlite, ring）
- [ ] 1.1.3: TypeScript 配置正确（tsconfig.json）
- [ ] 1.1.4: Vite 配置正确（vite.config.ts）
- [ ] 1.1.5: Tauri 配置正确（tauri.conf.json）
- [ ] 1.1.6: 目录结构完整（src/, src-tauri/src/）
- [ ] 1.2.1: Git 仓库初始化成功
- [ ] 1.2.2: .gitignore 文件正确配置
- [ ] 1.2.3: .env.example 文件创建
- [ ] 1.2.4: pnpm install 成功执行
- [ ] 1.2.5: cargo fetch 成功执行
- [ ] 1.2.6: 首个 Git 提交创建（Initial commit）
- [ ] 1.3.1: src/main.ts 创建并配置正确
- [ ] 1.3.2: src/App.vue 创建并配置正确
- [ ] 1.3.3: src/index.html 创建并配置正确
- [ ] 1.3.4: src-tauri/src/main.rs 创建并配置正确
- [ ] 1.3.5: src-tauri/src/lib.rs 创建并配置正确
- [ ] 1.3.6: pnpm tauri dev 成功启动，显示空白窗口

## 阶段 2: 核心架构实现

- [ ] 2.1.1: CustomTitlebar.vue 实现（Logo、应用名、主题切换、窗口控制）
- [ ] 2.1.2: NavigationMenu.vue 实现（5 个菜单项：首页、项目、载荷、插件、设置）
- [ ] 2.1.3: StatusBar.vue 实现（连接状态、项目数、载荷数）
- [ ] 2.1.4: MainLayout.vue 整合所有布局组件
- [ ] 2.1.5: 主题系统工作（亮色/暗色/跟随系统）
- [ ] 2.1.6: 语言切换工作（中文/英文）
- [ ] 2.2.1: 路由配置完成（/, /project, /payload, /plugin, /settings）
- [ ] 2.2.2: 路由守卫实现
- [ ] 2.2.3: 页面切换动画实现
- [ ] 2.3.1: Pinia 初始化完成
- [ ] 2.3.2: app.ts（应用全局状态）实现
- [ ] 2.3.3: payload.ts（载荷配置状态）实现
- [ ] 2.3.4: project.ts（项目状态）实现
- [ ] 2.3.5: webshell.ts（WebShell 连接状态）实现
- [ ] 2.4.1: i18n 配置完成
- [ ] 2.4.2: 中文语言包完整
- [ ] 2.4.3: 英文语言包完整
- [ ] 2.4.4: i18n 集成到 Vue 应用

## 阶段 3: 数据库基础设施

- [ ] 3.1.1: SQLite 数据库模块创建
- [ ] 3.1.2: 连接池管理实现
- [ ] 3.1.3: 数据模型 trait 定义
- [ ] 3.1.4: 数据库迁移脚本创建（001_create_tables.sql）
- [ ] 3.1.5: 数据库初始化逻辑实现
- [ ] 3.2.1: PayloadConfig 模型实现（含 CRUD）
- [ ] 3.2.2: Project 模型实现（含 CRUD）
- [ ] 3.2.3: WebShell 模型实现（含 CRUD）
- [ ] 3.2.4: Plugin 模型实现（含 CRUD）

## 阶段 4: 载荷管理模块

- [ ] 4.1.1: 载荷生成模块创建（src-tauri/core/payload/）
- [ ] 4.1.2: generator.rs（载荷生成器）实现
- [ ] 4.1.3: php.rs（PHP 载荷生成）实现
- [ ] 4.1.4: jsp.rs（JSP 载荷生成）实现
- [ ] 4.1.5: asp.rs（ASP 载荷生成）实现
- [ ] 4.1.6: aspx.rs（ASPX 载荷生成）实现
- [ ] 4.1.7: template.rs（模板管理）实现
- [ ] 4.2.1: create_payload_config 命令实现
- [ ] 4.2.2: update_payload_config 命令实现
- [ ] 4.2.3: delete_payload_config 命令实现（软删除）
- [ ] 4.2.4: get_all_payload_configs 命令实现
- [ ] 4.2.5: generate_payload 命令实现
- [ ] 4.3.1: PayloadView.vue（载荷管理页面）实现
- [ ] 4.3.2: PayloadConfig.vue（载荷配置表单）实现
- [ ] 4.3.3: PayloadPreview.vue（代码预览）实现
- [ ] 4.3.4: PayloadHistory.vue（历史记录）实现
- [ ] 4.3.5: PayloadTemplate.vue（模板管理）实现
- [ ] 4.3.6: LanguageSelector.vue（脚本语言选择）实现
- [ ] 4.3.7: EncryptionSelector.vue（加密算法选择）实现
- [ ] 4.3.8: ObfuscationSlider.vue（混淆强度滑块）实现
- [ ] 4.4.1: src/api/payload.ts 创建
- [ ] 4.4.2: invoke 调用后端命令集成
- [ ] 4.4.3: 错误处理实现
- [ ] 4.4.4: src/types/payload.ts 类型定义
- [ ] 4.5.1: 搜索功能实现
- [ ] 4.5.2: 过滤功能实现（按类型、标签）
- [ ] 4.5.3: 排序功能实现（多字段）
- [ ] 4.5.4: 集成到 Pinia store

## 阶段 5: 项目管理模块

- [ ] 5.1.1: 项目管理模块创建（src-tauri/core/project/）
- [ ] 5.1.2: manager.rs（项目管理器）实现
- [ ] 5.1.3: entity.rs（项目实体）实现
- [ ] 5.1.4: recycle_bin.rs（回收站）实现
- [ ] 5.2.1: WebShell 模块创建（src-tauri/core/webshell/）
- [ ] 5.2.2: connection.rs（连接管理）实现
- [ ] 5.2.3: session.rs（会话管理）实现
- [ ] 5.2.4: health_check.rs（健康检查）实现
- [ ] 5.3.1: create_project 命令实现
- [ ] 5.3.2: update_project 命令实现
- [ ] 5.3.3: delete_project 命令实现（软删除）
- [ ] 5.3.4: get_all_projects 命令实现
- [ ] 5.3.5: restore_project 命令实现
- [ ] 5.3.6: create_webshell 命令实现
- [ ] 5.3.7: update_webshell 命令实现
- [ ] 5.3.8: delete_webshell 命令实现（软删除）
- [ ] 5.3.9: test_connection 命令实现
- [ ] 5.3.10: connect_webshell 命令实现
- [ ] 5.4.1: ProjectView.vue（项目管理页面）实现
- [ ] 5.4.2: ProjectTree.vue（项目树组件）实现
- [ ] 5.4.3: ProjectSidebar.vue（左侧子栏）实现
- [ ] 5.4.4: WebShellList.vue（WebShell 列表表格）实现
- [ ] 5.4.5: WebShellForm.vue（WebShell 表单）实现
- [ ] 5.4.6: WebShellContextMenu.vue（右键菜单）实现
- [ ] 5.4.7: RecycleBin.vue（回收站视图）实现
- [ ] 5.4.8: Pagination.vue（分页器）实现
- [ ] 5.5.1: src/api/project.ts 创建
- [ ] 5.5.2: src/api/webshell.ts 创建
- [ ] 5.5.3: invoke 调用后端命令集成
- [ ] 5.5.4: src/types/project.ts 类型定义
- [ ] 5.5.5: src/types/webshell.ts 类型定义

## 阶段 6: 加密通信模块

- [ ] 6.1.1: 加密模块创建（src-tauri/core/crypto/）
- [ ] 6.1.2: aes.rs（AES-256-GCM 加密）实现
- [ ] 6.1.3: xor.rs（XOR 加密）实现
- [ ] 6.1.4: key_manager.rs（密钥管理）实现
- [ ] 6.1.5: Argon2 密钥派生实现
- [ ] 6.2.1: HTTP 模块创建（src-tauri/core/http/）
- [ ] 6.2.2: client.rs（HTTP 客户端）实现
- [ ] 6.2.3: proxy.rs（代理支持）实现
- [ ] 6.2.4: rate_limiter.rs（速率限制）实现
- [ ] 6.2.5: retry.rs（重试机制）实现
- [ ] 6.3.1: execute_command 命令实现
- [ ] 6.3.2: list_files 命令实现
- [ ] 6.3.3: upload_file 命令实现
- [ ] 6.3.4: download_file 命令实现
- [ ] 6.3.5: execute_query 命令实现

## 阶段 7: 控制台窗口系统

- [ ] 7.1.1: ConsoleWindow.vue（控制台窗口）实现
- [ ] 7.1.2: ConsoleTabs.vue（TAB 栏组件）实现
- [ ] 7.1.3: 窗口管理逻辑实现（打开、关闭、最小化、最大化）
- [ ] 7.1.4: 独立会话管理实现
- [ ] 7.2.1: FileBrowser.vue（文件浏览器）实现
- [ ] 7.2.2: FileUploader.vue（文件上传）实现
- [ ] 7.2.3: FileEditor.vue（文件编辑器）实现
- [ ] 7.2.4: FileContextMenu.vue（文件右键菜单）实现
- [ ] 7.2.5: 文件操作命令集成
- [ ] 7.3.1: DatabaseManager.vue（数据库管理器）实现
- [ ] 7.3.2: QueryEditor.vue（SQL 查询编辑器）实现
- [ ] 7.3.3: ResultTable.vue（查询结果表格）实现
- [ ] 7.3.4: DatabaseSelector.vue（数据库选择器）实现
- [ ] 7.3.5: 数据库操作命令集成
- [ ] 7.4.1: TerminalView.vue（终端视图）实现
- [ ] 7.4.2: TerminalInput.vue（命令输入框）实现
- [ ] 7.4.3: CommandHistory.vue（命令历史）实现
- [ ] 7.4.4: 命令执行命令集成

## 阶段 8: 设置页面

- [ ] 8.1.1: SettingsView.vue（设置页面）实现
- [ ] 8.1.2: SettingsLayout.vue（设置布局）实现
- [ ] 8.1.3: 设置路由实现（/settings/appearance, /settings/connection, /settings/security, /settings/about）
- [ ] 8.2.1: AppearanceSettings.vue（外观设置）实现
- [ ] 8.2.2: ThemeSelector.vue（主题选择器）实现
- [ ] 8.2.3: LanguageSelector.vue（语言选择器）实现
- [ ] 8.2.4: ColorPicker.vue（强调色选择器）实现
- [ ] 8.3.1: ConnectionSettings.vue（连接设置）实现
- [ ] 8.3.2: 请求配置实现（超时、重试、UA 池）
- [ ] 8.3.3: 代理配置实现（HTTP/SOCKS5）
- [ ] 8.4.1: SecuritySettings.vue（安全设置）实现
- [ ] 8.4.2: 加密算法选择实现
- [ ] 8.4.3: 密钥管理实现
- [ ] 8.4.4: 数据安全选项实现
- [ ] 8.5.1: AboutPage.vue（关于页面）实现
- [ ] 8.5.2: 版本信息卡片实现
- [ ] 8.5.3: 检查更新功能实现

## 阶段 9: 插件系统基础

- [ ] 9.1.1: 插件模块创建（src-tauri/plugins/）
- [ ] 9.1.2: loader.rs（插件加载器）实现
- [ ] 9.1.3: sandbox.rs（沙箱隔离）实现
- [ ] 9.1.4: runtime.rs（运行时）实现
- [ ] 9.2.1: PluginView.vue（插件管理页面）实现
- [ ] 9.2.2: PluginList.vue（插件列表）实现
- [ ] 9.2.3: PluginManager.vue（插件管理）实现
- [ ] 9.2.4: 插件安装/卸载/启用/禁用功能实现

## 阶段 10: 通用功能和优化

- [ ] 10.1.1: Loading.vue（加载组件）实现
- [ ] 10.1.2: EmptyState.vue（空状态组件）实现
- [ ] 10.1.3: ConfirmDialog.vue（确认对话框）实现
- [ ] 10.1.4: ErrorDisplay.vue（错误展示）实现
- [ ] 10.1.5: CodeEditor.vue（通用代码编辑器）实现
- [ ] 10.1.6: SearchBar.vue（通用搜索栏）实现
- [ ] 10.2.1: src-tauri/error.rs（错误处理）创建
- [ ] 10.2.2: 统一错误码系统实现（60+ 错误码）
- [ ] 10.2.3: 错误严重性分级实现
- [ ] 10.2.4: 用户友好的错误提示实现
- [ ] 10.3.1: src-tauri/logging.rs（日志系统）创建
- [ ] 10.3.2: 结构化日志实现
- [ ] 10.3.3: 多日志级别实现
- [ ] 10.3.4: 敏感信息脱敏实现
- [ ] 10.3.5: 日志轮转实现
- [ ] 10.4.1: 组件懒加载实现
- [ ] 10.4.2: 虚拟列表实现
- [ ] 10.4.3: 防抖和节流实现
- [ ] 10.4.4: 内存占用优化

## 阶段 11: 测试和文档

- [ ] 11.1.1: Rust 单元测试编写（core/ 模块）
- [ ] 11.1.2: 前端单元测试编写（composables/）
- [ ] 11.1.3: 测试覆盖率 > 90%
- [ ] 11.2.1: 载荷生成集成测试编写
- [ ] 11.2.2: 加密通信集成测试编写
- [ ] 11.2.3: 项目管理集成测试编写
- [ ] 11.3.1: 用户手册编写
- [ ] 11.3.2: 安装部署文档编写
- [ ] 11.3.3: README.md 更新

## 阶段 12: 构建和发布

- [ ] 12.1.1: scripts/build.ps1 创建
- [ ] 12.1.2: scripts/build.sh 创建
- [ ] 12.1.3: 生产环境构建配置
- [ ] 12.1.4: 代码签名配置（可选）
- [ ] 12.2.1: 全量测试执行
- [ ] 12.2.2: 所有 Critical/High Bug 修复
- [ ] 12.2.3: 性能测试和优化
- [ ] 12.3.1: Git 标签 v0.1.0 创建
- [ ] 12.3.2: Release 创建
- [ ] 12.3.3: 安装包发布
