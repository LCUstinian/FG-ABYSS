# FG-ABYSS 载荷模块重构 - 产品需求文档

## Overview
- **Summary**: 对现有 FG-ABYSS WebShell 生成器的载荷模块进行全面重构，采用全新的洋葱模型生成流水线，提供更丰富的配置选项、更强的免杀能力和更专业的用户界面。
- **Purpose**: 解决当前载荷生成器功能有限、免杀能力不足、配置选项不够灵活的问题，构建一个真正具有高隐蔽性的专业红队工具。
- **Target Users**: 安全研究人员、红队人员、渗透测试工程师。

## Goals
- 重构 UI 配置表单，支持更丰富的生成选项（生成模式、脚本类型、基础连接、载荷加密、传输编码、混淆强度等）
- 实现基于洋葱模型的后端生成流水线（模板组装 → 代码混淆 → 载荷加密 → 传输编码 → 最终组装）
- 提供实时安全等级预估面板，增强用户体验
- 支持多种脚本类型（PHP, JSP, ASP, ASPX, Python, Node.js）
- 实现强大的代码混淆功能（L1-L4 级别，支持 ASP 特殊 Chr() 混淆）
- 支持多种加密算法（XOR, AES-128, AES-256-GCM, ChaCha20）
- 分离载荷加密（安全核心）和传输编码（格式兼容）的概念

## Non-Goals (Out of Scope)
- 不重构载荷历史管理功能（保持现有实现）
- 不重构载荷模板管理功能（保持现有实现）
- 不实现 WebShell 连接管理功能（由其他模块负责）
- 不实现实时代码高亮功能（超出当前范围）

## Background & Context
- 当前项目已存在基础的载荷生成功能，但配置选项简单，免杀能力有限
- 现有架构使用 Simple/Advanced 两种模式，需要重构为更灵活的配置模型
- 项目技术栈：Tauri v2 (Rust) + Vue 3 (TypeScript) + Naive UI + Pinia
- 现有代码结构清晰，便于进行模块化重构

## Functional Requirements
- **FR-1**: 支持两种生成模式（经典回显、流量代理）
- **FR-2**: 支持 6 种脚本类型（PHP, JSP, ASP, ASPX, Python, Node.js）
- **FR-3**: 支持配置基础连接（信道：POST Body/Header(Cookie/UA)/GET；格式：Raw/JSON）
- **FR-4**: 支持连接密码输入和一键随机生成
- **FR-5**: 支持多种载荷加密选项（无、XOR、AES-128、AES-256-GCM、ChaCha20）
- **FR-6**: 支持多种传输编码选项（Base64、Hex、URL Encode、None）
- **FR-7**: 支持 4 级混淆强度（L1-L4）
- **FR-8**: 支持自定义输出文件名并提供智能推荐
- **FR-9**: 实现洋葱模型后端生成流水线
- **FR-10**: 提供实时安全等级预估面板
- **FR-11**: 生成过程中显示 Rust 处理日志
- **FR-12**: 生成成功后提供下载文件和复制代码按钮

## Non-Functional Requirements
- **NFR-1**: 生成响应时间 < 2 秒（常规配置）
- **NFR-2**: 界面采用深色/浅色双主题，专业感强
- **NFR-3**: 所有用户界面文本使用中文，支持国际化
- **NFR-4**: 生成的 WebShell 无硬编码密钥
- **NFR-5**: ASP 生成确保纯 ASCII 输出，无乱码
- **NFR-6**: 代码健壮，处理异常情况（如密码为空的警告）

## Constraints
- **Technical**: 必须使用现有技术栈（Tauri v2 + Vue 3 + TypeScript + Naive UI + Pinia），不能引入新的 UI 组件库（如 Shadcn-Vue）
- **Business**: 需要向后兼容现有数据结构（如历史记录）
- **Dependencies**: 需要确保与现有载荷历史、载荷模板模块的集成

## Assumptions
- 用户了解基本的 WebShell 生成和免杀概念
- 现有项目结构可以支持新的功能扩展
- Naive UI 组件库足够实现新的 UI 需求
- 可以在现有类型系统基础上扩展新的配置选项

## Acceptance Criteria

### AC-1: 新配置表单实现
- **Given**: 用户打开载荷生成页面
- **When**: 用户查看配置区域
- **Then**: 显示完整的新配置表单，包括所有要求的字段（生成模式、脚本类型、基础连接、连接密码、载荷加密、传输编码、混淆强度、输出文件名）
- **Verification**: `human-judgment`
- **Notes**: 确保使用折叠面板优化布局

### AC-2: 生成模式切换
- **Given**: 用户在载荷生成页面
- **When**: 用户切换生成模式（经典回显 ↔ 流量代理）
- **Then**: 相应的配置选项根据模式显示/隐藏
- **Verification**: `programmatic`

### AC-3: 脚本类型选择
- **Given**: 用户在载荷生成页面
- **When**: 用户选择不同的脚本类型
- **Then**: 配置选项正确响应（如选择 ASP 时启用特殊 Chr() 混淆）
- **Verification**: `programmatic`

### AC-4: 连接密码随机生成
- **Given**: 用户在载荷生成页面
- **When**: 用户点击随机生成密码按钮
- **Then**: 连接密码字段填充随机生成的强密码
- **Verification**: `programmatic`

### AC-5: 洋葱模型生成流水线
- **Given**: 用户配置完所有选项并点击生成
- **When**: 后端处理生成请求
- **Then**: 严格按照洋葱模型处理：模板组装 → 代码混淆 → 载荷加密 → 传输编码 → 最终组装
- **Verification**: `programmatic`

### AC-6: 代码混淆级别
- **Given**: 用户选择不同的混淆级别（L1-L4）
- **When**: 生成载荷
- **Then**: 应用对应的混淆技术：L1(变量重命名) → L2(垃圾代码) → L3(控制流平坦化) → L4(沙箱检测)
- **Verification**: `programmatic`

### AC-7: ASP 特殊混淆
- **Given**: 用户选择 ASP 脚本类型
- **When**: 生成载荷
- **Then**: 所有关键字（如 Execute）转换为 Chr() 拼接，确保纯 ASCII 输出
- **Verification**: `programmatic`

### AC-8: 载荷加密和传输编码分离
- **Given**: 用户配置载荷加密和传输编码
- **When**: 查看生成过程
- **Then**: 载荷加密（安全层）和传输编码（兼容层）作为独立步骤清晰分离
- **Verification**: `human-judgment`

### AC-9: 实时安全等级预估
- **Given**: 用户正在配置选项
- **When**: 用户修改任何配置选项
- **Then**: 右侧面板实时更新安全等级预估（如 ⭐⭐⭐⭐☆）
- **Verification**: `human-judgment`

### AC-10: 生成日志显示
- **Given**: 用户点击生成按钮
- **When**: 生成过程进行中
- **Then**: 显示 Rust 处理日志
- **Verification**: `human-judgment`

### AC-11: 下载和复制功能
- **Given**: 载荷生成成功
- **When**: 用户点击下载或复制按钮
- **Then**: 成功下载文件或复制代码到剪贴板
- **Verification**: `programmatic`

### AC-12: 异常处理
- **Given**: 用户输入无效配置（如空密码）
- **When**: 用户点击生成
- **Then**: 显示清晰的错误提示，不崩溃
- **Verification**: `programmatic`

## Open Questions
- [ ] 是否需要保留现有的 Simple/Advanced 模式作为兼容选项？
- [ ] 流量代理模式的具体功能需求是什么？（SOCKS/HTTP 隧道实现细节）
- [ ] 是否需要支持撤销/重做配置操作？
- [ ] 智能推荐文件名的具体规则是什么？
