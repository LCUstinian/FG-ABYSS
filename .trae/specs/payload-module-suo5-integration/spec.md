# FG-ABYSS 载荷模块重构 - Suo5 集成版 - 产品需求文档

## Overview
- **Summary**: 对 FG-ABYSS WebShell 生成器进行全面升级，深度集成开源项目 Suo5 的核心代理逻辑，构建集高隐蔽文件型 WebShell、主流内存马注入器、定制化 Suo5 代理载荷于一体的专业红队武器平台。
- **Purpose**: 解决传统 WebShell 易被检测、代理功能薄弱的问题，提供 2026 年最新攻防对抗趋势下的高隐蔽性解决方案。
- **Target Users**: 高级安全研究人员、专业红队人员、资深渗透测试工程师。

## Goals
- 实现三大实战场景：文件落地、内存马注入、纯 Suo5 代理
- 深度集成 Suo5 (GitHub: zema1/suo5) 核心载荷逻辑
- 支持 PHP, JSP, ASP, ASPX (.NET) 多种脚本类型
- 提供强大的变量混淆和垃圾代码注入功能
- 支持 Java (Tomcat Filter/Valve, Spring Interceptor) 和 .NET (IIS HttpModule) 内存马注入
- 生成与官方 Suo5 客户端兼容的服务端代码
- 提供智能预检功能（语法检查、协议完整性验证、参数一致性检查）
- 自动生成对应的 Suo5 客户端命令字符串

## Non-Goals (Out of Scope)
- 不实现远程连接测试（仅本地验证）
- 不重构载荷历史管理功能
- 不重构载荷模板管理功能
- 不实现 WebShell 连接管理功能
- 不实现实时代码高亮功能

## Background & Context
- 当前项目已存在基础的载荷生成功能
- Suo5 是一个开源的高性能 HTTP 隧道代理工具，支持 Chunked-Encoding
- 传统 WebShell 易被 WAF/IDS 检测，需要更隐蔽的代理方案
- 内存马技术是红队实战的重要手段
- 项目技术栈：Tauri v2 (Rust) + Vue 3 (TypeScript) + Naive UI + Pinia

## Functional Requirements
- **FR-1**: 支持三模态切换（文件落地、内存马注入、纯 Suo5 代理）
- **FR-2**: 集成 Suo5 官方的 JSP/ASPX/PHP 服务端代码模板
- **FR-3**: 支持自定义 Suo5 认证密码、监听路径、超时时间
- **FR-4**: 对生成的 Suo5 载荷进行变量混淆和垃圾代码注入
- **FR-5**: 支持生成"纯 Suo5 代理文件"或"WebShell + Suo5 双功能文件"
- **FR-6**: 支持 Java 内存马注入（Tomcat Filter/Valve, Spring Interceptor）
- **FR-7**: 支持 .NET 内存马注入（IIS HttpModule, IHttpHandler 动态注册）
- **FR-8**: 支持内存马注入后动态注册 Suo5 代理端点
- **FR-9**: 提供智能预检功能（语法检查、协议完整性、参数一致性、尺寸预警）
- **FR-10**: 自动生成对应的 Suo5 客户端命令字符串
- **FR-11**: 支持载荷加密（AES-256-GCM, XOR, None）
- **FR-12**: 支持 ASP 特化处理（Chr() 拼接 + 纯 ASCII）
- **FR-13**: 支持自毁逻辑（仅内存马/一次性模式）
- **FR-14**: 提供 Suo5 专属配置面板和客户端命令生成按钮

## Non-Functional Requirements
- **NFR-1**: 生成响应时间 < 3 秒（常规配置）
- **NFR-2**: 界面采用深色/浅色双主题，专业感强
- **NFR-3**: 所有用户界面文本使用中文，支持国际化
- **NFR-4**: 生成的代码必须能与官方 Suo5 客户端正常握手通信
- **NFR-5**: 混淆不能改变 Suo5 的业务逻辑，确保代理功能正常
- **NFR-6**: ASP 生成确保纯 ASCII 输出，无乱码
- **NFR-7**: 代码健壮，错误处理完善
- **NFR-8**: 底部显著标注"本工具仅供授权渗透测试与安全研究使用"

## Constraints
- **Technical**: 必须使用现有技术栈（Tauri v2 + Vue 3 + TypeScript + Naive UI + Pinia），不能引入新的 UI 组件库（如 Shadcn-Vue）
- **Business**: 需要向后兼容现有数据结构
- **Suo5 兼容性**: 生成的代码必须能与官方 Suo5 客户端正常握手通信，严禁破坏核心协议逻辑
- **混淆安全性**: 混淆不能改变 Suo5 的业务逻辑，确保代理功能正常
- **无幻觉测试**: 不做远程连接测试，但必须在本地模拟参数替换和语法检查过程

## Assumptions
- 用户了解基本的 WebShell 生成和 Suo5 代理使用方法
- 现有项目结构可以支持新的功能扩展
- Naive UI 组件库足够实现新的 UI 需求
- 可以获取 Suo5 官方的最新稳定版源码
- Suo5 的核心逻辑可以通过变量重命名和字符串分割进行混淆而不破坏功能

## Acceptance Criteria

### AC-1: 三模态切换实现
- **Given**: 用户打开载荷生成页面
- **When**: 用户查看页面顶部
- **Then**: 显示显著的 Segmented Control 切换开关，包含三个模式：文件落地、内存马注入、纯 Suo5 代理
- **Verification**: `human-judgment`
- **Notes**: 切换时表单字段平滑过渡

### AC-2: 模式 A - 文件落地
- **Given**: 用户选择模式 A（文件落地）
- **When**: 用户配置选项
- **Then**: 可以选择仅包含命令执行，或集成 Suo5 代理功能（二选一）
- **Verification**: `programmatic`

### AC-3: 模式 B - 内存马注入
- **Given**: 用户选择模式 B（内存马注入）
- **When**: 用户选择 Java 或 .NET 脚本类型
- **Then**: 显示对应的注入类型选项（Tomcat Filter/Spring Interceptor 或 IIS HttpModule），可以选择动态注册 Suo5 代理端点
- **Verification**: `programmatic`

### AC-4: 模式 C - 纯 Suo5 代理
- **Given**: 用户选择模式 C（纯 Suo5 代理）
- **When**: 查看表单
- **Then**: 自动显示 Suo5 特有配置项（Auth, Path, Timeout），并隐藏"混淆强度"（仅做变量重命名）
- **Verification**: `programmatic`

### AC-5: Suo5 参数配置
- **Given**: 用户在任何 Suo5 相关模式
- **When**: 用户配置 Suo5 参数
- **Then**: 可以自定义认证密码（支持随机生成）、监听路径（支持随机生成）、超时时间
- **Verification**: `programmatic`

### AC-6: Suo5 模板集成
- **Given**: 用户生成 Suo5 相关载荷
- **When**: 后端处理生成请求
- **Then**: 使用 Suo5 官方的 JSP/ASPX/PHP 源码作为基础模板
- **Verification**: `programmatic`
- **Notes**: 模板必须是最新稳定版，支持 Chunked-Encoding

### AC-7: 参数动态注入
- **Given**: 用户配置了 Suo5 参数
- **When**: 生成载荷
- **Then**: 将模板中的硬编码密码、路径等替换为用户配置的值
- **Verification**: `programmatic`

### AC-8: 变量混淆
- **Given**: 用户配置了混淆强度
- **When**: 生成 Suo5 载荷
- **Then**: 将模板中的关键变量名随机化，严禁改变逻辑
- **Verification**: `programmatic`

### AC-9: 字符串分割
- **Given**: 用户配置了 L2 混淆强度
- **When**: 生成 Suo5 载荷
- **Then**: 将密码字符串分割成多段拼接，防止明文匹配
- **Verification**: `programmatic`

### AC-10: 垃圾代码注入
- **Given**: 用户配置了 L2 混淆强度
- **When**: 生成 Suo5 载荷
- **Then**: 在 Suo5 核心循环前后插入无用的计算逻辑或死代码
- **Verification**: `programmatic`

### AC-11: 内存马适配 - Java
- **Given**: 用户选择模式 B + Java 脚本类型
- **When**: 生成载荷
- **Then**: 生成能在一次性执行后将 Suo5 核心逻辑注册到容器事件管道的代码
- **Verification**: `programmatic`

### AC-12: 内存马适配 - .NET
- **Given**: 用户选择模式 B + .NET 脚本类型
- **When**: 生成载荷
- **Then**: 利用反射操作动态注册包含 Suo5 逻辑的 IHttpModule
- **Verification**: `programmatic`

### AC-13: 智能预检 - 语法检查
- **Given**: 用户点击生成
- **When**: 后端处理
- **Then**: 确保生成的 JSP/ASPX/PHP 语法正确，括号闭合
- **Verification**: `programmatic`

### AC-14: 智能预检 - 协议完整性
- **Given**: 用户生成 Suo5 载荷
- **When**: 后端处理
- **Then**: 检查生成的代码是否包含 Suo5 特有的 Chunked-Encoding 处理逻辑
- **Verification**: `programmatic`

### AC-15: 智能预检 - 参数一致性
- **Given**: 用户配置了参数
- **When**: 后端处理
- **Then**: 验证用户设置的密码和路径是否正确嵌入到代码中
- **Verification**: `programmatic`

### AC-16: 智能预检 - 尺寸预警
- **Given**: 生成的代码过大
- **When**: 后端处理完成
- **Then**: 提示可能触发 WAF 长度限制
- **Verification**: `human-judgment`

### AC-17: 客户端命令生成
- **Given**: 用户生成了 Suo5 载荷
- **When**: 用户查看结果
- **Then**: 自动生成对应的 ./suo5 -t <url> -auth <password> 命令字符串
- **Verification**: `programmatic`

### AC-18: 载荷加密
- **Given**: 用户选择了载荷加密选项
- **When**: 生成载荷
- **Then**: 对 Suo5 核心逻辑之外的辅助代码进行加密
- **Verification**: `programmatic`

### AC-19: ASP 特化处理
- **Given**: 用户选择 ASP 脚本类型
- **When**: 生成载荷
- **Then**: 强制使用 Chr() 拼接 + 纯 ASCII
- **Verification**: `programmatic`

### AC-20: 自毁逻辑
- **Given**: 用户在模式 B 且启用自毁逻辑
- **When**: 生成载荷
- **Then**: 注入后删除文件或返回 404
- **Verification**: `programmatic`

### AC-21: Suo5 专属面板
- **Given**: 用户选择 Suo5 相关模式
- **When**: 查看界面
- **Then**: 展示清晰的参数配置区，并提供"生成客户端命令"按钮
- **Verification**: `human-judgment`

### AC-22: 安全警示
- **Given**: 用户在载荷生成页面
- **When**: 查看页面底部
- **Then**: 显著位置标注"本工具仅供授权渗透测试与安全研究使用"
- **Verification**: `human-judgment`

## Open Questions
- [ ] Suo5 官方源码的具体版本是哪个？
- [ ] ASP 是否有 Suo5 的实现？如果没有，是否需要禁用 Suo5 集成？
- [ ] 内存马注入的具体实现细节需要进一步明确
- [ ] 是否需要保留之前的载荷生成功能作为兼容选项？
