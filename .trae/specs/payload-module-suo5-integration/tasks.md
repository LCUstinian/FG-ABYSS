# FG-ABYSS 载荷模块重构 - Suo5 集成版 - 实施计划

## [/] Task 1: 扩展 TypeScript 类型定义 - Suo5 集成
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 扩展 `src/types/payload.ts` 中的类型定义
  - 添加三种生成模式类型（文件落地、内存马注入、纯 Suo5 代理）
  - 添加注入类型（Tomcat Filter, Spring Interceptor, IIS HttpModule）
  - 添加 Suo5 配置类型（auth, path, timeout）
  - 更新相关的选项常量数组
- **Acceptance Criteria Addressed**: [FR-1, FR-3, FR-6, FR-7]
- **Test Requirements**:
  - `programmatic` TR-1.1: 类型定义编译无错误
  - `programmatic` TR-1.2: 新增的选项常量包含所有要求的值
- **Notes**: 确保向后兼容现有类型

## [ ] Task 2: 更新国际化配置 - Suo5 集成
- **Priority**: P0
- **Depends On**: Task 1
- **Description**: 
  - 更新 `src/i18n/zh-CN.ts`，添加所有新 UI 文本的中文翻译
  - 更新 `src/i18n/en-US.ts`，添加对应的英文翻译
  - 确保所有文本支持国际化，无硬编码
- **Acceptance Criteria Addressed**: [FR-1, FR-3, FR-6, FR-7, FR-13, FR-14, NFR-3]
- **Test Requirements**:
  - `human-judgement` TR-2.1: 所有新增 UI 文本都有对应的翻译
  - `human-judgement` TR-2.2: 中文翻译准确专业
- **Notes**: 保持与现有翻译风格一致，添加安全警示文本

## [ ] Task 3: 重构 Pinia Store - Suo5 集成
- **Priority**: P0
- **Depends On**: Task 1
- **Description**: 
  - 更新 `src/stores/payload.ts` 中的 config 状态结构
  - 添加三模式状态管理
  - 添加 Suo5 配置项（auth, path, timeout）
  - 添加注入类型配置
  - 添加自毁逻辑开关
  - 实现配置选项之间的联动逻辑
  - 添加 Suo5 密码随机生成功能
  - 添加 Suo5 路径随机生成功能
  - 添加客户端命令生成功能
- **Acceptance Criteria Addressed**: [FR-3, FR-10, FR-13, AC-5, AC-17]
- **Test Requirements**:
  - `programmatic` TR-3.1: store 更新编译无错误
  - `programmatic` TR-3.2: 配置联动逻辑正确响应
  - `programmatic` TR-3.3: 随机密码和路径生成功能正常
  - `programmatic` TR-3.4: 客户端命令生成正确
- **Notes**: 确保现有功能不被破坏

## [ ] Task 4: 获取并集成 Suo5 模板
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 从 GitHub zema1/suo5 获取最新稳定版源码
  - 收集 suo5.jsp, suo5.aspx, suo5.php 模板
  - 创建 `src-tauri/resources/` 目录
  - 将模板存入 Rust 资源文件
  - 确保模板支持 Chunked-Encoding
- **Acceptance Criteria Addressed**: [FR-2, AC-6]
- **Test Requirements**:
  - `programmatic` TR-4.1: 模板文件成功集成
  - `human-judgement` TR-4.2: 模板是最新稳定版
- **Notes**: 确保遵守 Suo5 的开源协议

## [x] Task 5: 扩展 Rust 类型定义 - Suo5 集成
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 更新 `src-tauri/src/types/payload.rs` 中的类型定义
  - 与前端 TypeScript 类型保持一致
  - 添加新的枚举类型（生成模式、注入类型、Suo5 配置）
  - 更新 PayloadConfig 结构体
- **Acceptance Criteria Addressed**: [FR-1, FR-3, FR-6, FR-7]
- **Test Requirements**:
  - `programmatic` TR-5.1: Rust 类型定义编译无错误
  - `programmatic` TR-5.2: 与前端类型序列化/反序列化兼容
- **Notes**: 使用 serde 确保与 Tauri 命令兼容

## [ ] Task 6: 实现 Suo5 载荷定制器 - 核心模块
- **Priority**: P0
- **Depends On**: Task 4, Task 5
- **Description**: 
  - 创建 `src-tauri/src/suo5_generator.rs`
  - 实现模板加载功能（从资源文件加载 suo5.jsp/aspx/php）
  - 实现参数动态注入（替换硬编码的密码、路径等）
  - 实现变量重命名混淆（L1）
  - 实现字符串分割混淆（L2）
  - 实现垃圾代码注入（L2）
  - 确保不改变 Suo5 核心协议逻辑
- **Acceptance Criteria Addressed**: [FR-2, FR-4, AC-6, AC-7, AC-8, AC-9, AC-10, NFR-4, NFR-5]
- **Test Requirements**:
  - `programmatic` TR-6.1: suo5_generator.rs 编译无错误
  - `programmatic` TR-6.2: 参数注入正确工作
  - `programmatic` TR-6.3: 变量重命名不破坏逻辑
  - `programmatic` TR-6.4: 字符串分割正确实现
  - `programmatic` TR-6.5: 垃圾代码注入正确实现
- **Notes**: 这是核心难点，重点实现

## [ ] Task 7: 实现内存马注入器 - Java
- **Priority**: P1
- **Depends On**: Task 6
- **Description**: 
  - 在 `src-tauri/src/` 创建或更新内存马模块
  - 实现 Tomcat Filter 注入
  - 实现 Spring Interceptor 注入
  - 实现将 Suo5 逻辑集成到内存马
  - 实现自毁逻辑
- **Acceptance Criteria Addressed**: [FR-6, FR-8, FR-13, AC-3, AC-11, AC-20]
- **Test Requirements**:
  - `programmatic` TR-7.1: 内存马模块编译无错误
  - `programmatic` TR-7.2: Tomcat Filter 注入代码生成正确
  - `programmatic` TR-7.3: Spring Interceptor 注入代码生成正确
  - `programmatic` TR-7.4: Suo5 逻辑集成正确
- **Notes**: 确保不破坏 Suo5 协议逻辑

## [ ] Task 8: 实现内存马注入器 - .NET
- **Priority**: P1
- **Depends On**: Task 6
- **Description**: 
  - 实现 IIS HttpModule 动态注册
  - 实现 IHttpHandler 动态注册
  - 利用反射操作 System.Web.Hosting.HostingEnvironment
  - 实现将 Suo5 逻辑集成到 .NET 内存马
  - 实现自毁逻辑
- **Acceptance Criteria Addressed**: [FR-7, FR-8, FR-13, AC-3, AC-12, AC-20]
- **Test Requirements**:
  - `programmatic` TR-8.1: .NET 内存马代码编译无错误
  - `programmatic` TR-8.2: HttpModule 注册代码生成正确
  - `programmatic` TR-8.3: Suo5 逻辑集成正确
- **Notes**: 确保不破坏 Suo5 协议逻辑

## [ ] Task 9: 实现智能预检模块
- **Priority**: P0
- **Depends On**: Task 6
- **Description**: 
  - 实现语法检查（确保 JSP/ASPX/PHP 语法正确，括号闭合）
  - 实现 Suo5 协议完整性检查（正则匹配关键函数名）
  - 实现参数一致性检查（验证密码和路径是否正确嵌入）
  - 实现尺寸预警（代码过大时提示）
- **Acceptance Criteria Addressed**: [FR-9, AC-13, AC-14, AC-15, AC-16]
- **Test Requirements**:
  - `programmatic` TR-9.1: 预检模块编译无错误
  - `programmatic` TR-9.2: 语法检查正确工作
  - `programmatic` TR-9.3: 协议完整性检查正确工作
  - `programmatic` TR-9.4: 参数一致性检查正确工作
- **Notes**: 严禁真实网络测试，仅本地验证

## [ ] Task 10: 实现客户端命令生成器
- **Priority**: P0
- **Depends On**: Task 5
- **Description**: 
  - 实现 generate_client_command 函数
  - 根据用户配置生成对应的 ./suo5 -t <url> -auth <password> 命令字符串
  - 支持其他 Suo5 参数（如超时）
- **Acceptance Criteria Addressed**: [FR-10, AC-17]
- **Test Requirements**:
  - `programmatic` TR-10.1: 命令生成器编译无错误
  - `programmatic` TR-10.2: 生成的命令格式正确
- **Notes**: 确保与官方 Suo5 客户端命令格式一致

## [ ] Task 11: 更新载荷加密模块 - Suo5 兼容
- **Priority**: P1
- **Depends On**: Task 5
- **Description**: 
  - 更新或创建载荷加密模块
  - 实现 AES-256-GCM 加密
  - 实现 XOR 加密
  - 确保仅对 Suo5 核心逻辑之外的辅助代码加密
- **Acceptance Criteria Addressed**: [FR-11, AC-18]
- **Test Requirements**:
  - `programmatic` TR-11.1: 加密模块编译无错误
  - `programmatic` TR-11.2: 所有加密算法正确实现
- **Notes**: Suo5 内部协议已加密，此处为额外保护层

## [ ] Task 12: 更新 ASP ASCII 转换器
- **Priority**: P1
- **Depends On**: Task 5
- **Description**: 
  - 完善 asp_ascii_converter
  - 确保选择 ASP 时强制使用 Chr() 拼接
  - 确保 100% 纯 ASCII 输出
- **Acceptance Criteria Addressed**: [FR-12, AC-19, NFR-6]
- **Test Requirements**:
  - `programmatic` TR-12.1: 转换器编译无错误
  - `programmatic` TR-12.2: 输出为纯 ASCII
- **Notes**: 仅文件模式适用

## [ ] Task 13: 重构洋葱模型生成流水线 - Suo5 集成
- **Priority**: P0
- **Depends On**: Task 6, Task 7, Task 8, Task 9, Task 10, Task 11, Task 12
- **Description**: 
  - 重构 `src-tauri/src/core/generator.rs`
  - 集成 Suo5 定制器
  - 集成内存马注入器
  - 集成智能预检模块
  - 集成客户端命令生成器
  - 根据三种模式分别处理
  - 更新 Tauri 命令
- **Acceptance Criteria Addressed**: [FR-1, FR-2, FR-5, FR-6, FR-7, FR-8, FR-9, FR-10, FR-11, FR-12, FR-13]
- **Test Requirements**:
  - `programmatic` TR-13.1: 生成流水线编译无错误
  - `programmatic` TR-13.2: 三种模式都正确处理
  - `programmatic` TR-13.3: 所有模块正确集成
- **Notes**: 保持向后兼容

## [/] Task 14: 重构 UI 组件 - 三模态切换
- **Priority**: P0
- **Depends On**: Task 1, Task 2, Task 3
- **Description**: 
  - 更新 `src/components/business/payload/PayloadGeneratorView.vue`
  - 实现顶部显著的 Segmented Control 三模态切换
  - 实现模式切换时表单字段的平滑过渡
  - 实现模式 A（文件落地）的表单
  - 实现模式 B（内存马注入）的表单
  - 实现模式 C（纯 Suo5 代理）的表单
- **Acceptance Criteria Addressed**: [FR-1, AC-1, AC-2, AC-3, AC-4]
- **Test Requirements**:
  - `human-judgement` TR-14.1: 三模态切换器美观显著
  - `programmatic` TR-14.2: 模式切换正确联动表单
  - `human-judgement` TR-14.3: 表单字段平滑过渡
- **Notes**: 保持 Naive UI 组件库的使用

## [/] Task 15: 重构 UI 组件 - Suo5 专属面板
- **Priority**: P0
- **Depends On**: Task 14
- **Description**: 
  - 实现 Suo5 专属配置面板（密码、路径、超时）
  - 实现密码随机生成按钮
  - 实现路径随机生成按钮
  - 实现"生成客户端命令"按钮
  - 实现一键复制客户端命令功能
  - 更新代码预览显示
  - 在底部添加安全警示
- **Acceptance Criteria Addressed**: [FR-3, FR-10, FR-14, AC-5, AC-17, AC-21, AC-22, NFR-8]
- **Test Requirements**:
  - `human-judgement` TR-15.1: Suo5 专属面板清晰易用
  - `programmatic` TR-15.2: 客户端命令生成正确
  - `programmatic` TR-15.3: 复制功能正常工作
  - `human-judgement` TR-15.4: 安全警示显著
- **Notes**: 确保用户体验良好

## [ ] Task 16: 集成测试和联调 - Suo5 版
- **Priority**: P0
- **Depends On**: Task 13, Task 15
- **Description**: 
  - 测试完整的端到端生成流程
  - 测试三种模式都正常工作
  - 测试所有脚本类型生成
  - 测试 Suo5 参数配置
  - 测试混淆功能
  - 测试智能预检
  - 测试客户端命令生成
  - 测试下载和复制功能
  - 测试异常情况处理
- **Acceptance Criteria Addressed**: [NFR-1, NFR-4, NFR-5, NFR-7]
- **Test Requirements**:
  - `programmatic` TR-16.1: 端到端流程无崩溃
  - `programmatic` TR-16.2: 三种模式都正常工作
  - `programmatic` TR-16.3: 所有配置组合正常工作
  - `human-judgement` TR-16.4: 生成响应时间 < 3 秒
- **Notes**: 确保与现有载荷历史模块的集成
