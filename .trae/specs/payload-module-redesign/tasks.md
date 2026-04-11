# FG-ABYSS 载荷模块重构 - 实施计划

## [ ] Task 1: 扩展 TypeScript 类型定义
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 扩展 `src/types/payload.ts` 中的类型定义
  - 添加新的生成模式类型（经典回显、流量代理）
  - 添加新的脚本类型（Python、Node.js）
  - 添加基础连接配置类型（信道、格式）
  - 添加新的载荷加密算法类型（AES-256-GCM、ChaCha20）
  - 添加新的传输编码类型
  - 添加 4 级混淆强度类型（L1-L4）
  - 更新相关的选项常量数组
- **Acceptance Criteria Addressed**: [FR-1, FR-2, FR-3, FR-5, FR-6, FR-7]
- **Test Requirements**:
  - `programmatic` TR-1.1: 类型定义编译无错误
  - `programmatic` TR-1.2: 新增的选项常量包含所有要求的值
- **Notes**: 确保向后兼容现有类型

## [ ] Task 2: 更新国际化配置
- **Priority**: P0
- **Depends On**: Task 1
- **Description**: 
  - 更新 `src/i18n/zh-CN.ts`，添加所有新 UI 文本的中文翻译
  - 更新 `src/i18n/en-US.ts`，添加对应的英文翻译
  - 确保所有文本支持国际化，无硬编码
- **Acceptance Criteria Addressed**: [FR-1, FR-2, FR-3, FR-4, FR-5, FR-6, FR-7, FR-8, NFR-3]
- **Test Requirements**:
  - `human-judgement` TR-2.1: 所有新增 UI 文本都有对应的翻译
  - `human-judgement` TR-2.2: 中文翻译准确专业
- **Notes**: 保持与现有翻译风格一致

## [ ] Task 3: 重构 Pinia Store
- **Priority**: P0
- **Depends On**: Task 1
- **Description**: 
  - 更新 `src/stores/payload.ts` 中的 config 状态结构
  - 添加新配置项的 setter 方法
  - 实现配置选项之间的联动逻辑
  - 添加安全等级计算逻辑
  - 添加密码随机生成功能
  - 添加文件名智能推荐功能
  - 保留与历史记录的兼容性
- **Acceptance Criteria Addressed**: [FR-4, FR-8, FR-10, AC-4]
- **Test Requirements**:
  - `programmatic` TR-3.1: store 更新编译无错误
  - `programmatic` TR-3.2: 配置联动逻辑正确响应
  - `programmatic` TR-3.3: 随机密码生成功能正常
- **Notes**: 确保现有功能不被破坏

## [ ] Task 4: 重构 UI 组件 - 配置表单
- **Priority**: P0
- **Depends On**: Task 1, Task 2, Task 3
- **Description**: 
  - 完全重写 `src/components/business/payload/PayloadGeneratorView.vue` 的配置表单部分
  - 使用折叠面板优化布局
  - 实现生成模式切换（经典回显/流量代理）
  - 实现脚本类型选择（6种）
  - 实现基础连接配置（信道+格式）
  - 实现连接密码输入（含随机生成按钮）
  - 实现载荷加密选择
  - 实现传输编码选择
  - 实现混淆强度滑块（L1-L4）
  - 实现输出文件名输入（含智能推荐）
- **Acceptance Criteria Addressed**: [FR-1, FR-2, FR-3, FR-4, FR-5, FR-6, FR-7, FR-8, AC-1, AC-2, AC-3, AC-4]
- **Test Requirements**:
  - `human-judgement` TR-4.1: UI 布局合理，使用折叠面板
  - `programmatic` TR-4.2: 所有表单控件正确绑定到 store
  - `programmatic` TR-4.3: 生成模式切换正确显示/隐藏选项
- **Notes**: 保持 Naive UI 组件库的使用，不引入新库

## [ ] Task 5: 重构 UI 组件 - 预览和安全等级面板
- **Priority**: P1
- **Depends On**: Task 4
- **Description**: 
  - 更新 `PayloadGeneratorView.vue` 的预览区
  - 实现实时安全等级预估面板（⭐⭐⭐⭐☆）
  - 更新代码预览显示
  - 实现生成日志显示区域
  - 更新状态信息显示
- **Acceptance Criteria Addressed**: [FR-10, FR-11, AC-9, AC-10]
- **Test Requirements**:
  - `human-judgement` TR-5.1: 安全等级面板美观专业
  - `human-judgement` TR-5.2: 安全等级随配置变化实时更新
  - `human-judgement` TR-5.3: 生成日志清晰可读
- **Notes**: 安全等级计算需要考虑：加密强度、混淆级别、编码方式等

## [ ] Task 6: 扩展 Rust 类型定义
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 更新 `src-tauri/src/types/payload.rs` 中的类型定义
  - 与前端 TypeScript 类型保持一致
  - 添加新的枚举类型（生成模式、脚本类型、信道、格式、加密算法、传输编码、混淆级别）
  - 更新 PayloadConfig 结构体
- **Acceptance Criteria Addressed**: [FR-1, FR-2, FR-3, FR-5, FR-6, FR-7]
- **Test Requirements**:
  - `programmatic` TR-6.1: Rust 类型定义编译无错误
  - `programmatic` TR-6.2: 与前端类型序列化/反序列化兼容
- **Notes**: 使用 serde 确保与 Tauri 命令兼容

## [ ] Task 7: 实现 Rust 代码混淆模块
- **Priority**: P0
- **Depends On**: Task 6
- **Description**: 
  - 在 `src-tauri/src/core/` 创建或更新混淆模块
  - 实现 L1 混淆：变量重命名
  - 实现 L2 混淆：垃圾代码注入
  - 实现 L3 混淆：控制流平坦化
  - 实现 L4 混淆：沙箱检测
  - 实现 ASP 特殊 Chr() 混淆
- **Acceptance Criteria Addressed**: [FR-7, AC-6, AC-7]
- **Test Requirements**:
  - `programmatic` TR-7.1: 混淆模块编译无错误
  - `programmatic` TR-7.2: L1-L4 各级别混淆正确应用
  - `programmatic` TR-7.3: ASP Chr() 混淆生成纯 ASCII 输出
- **Notes**: 混淆逻辑需要针对不同脚本类型分别实现

## [ ] Task 8: 实现 Rust 载荷加密模块
- **Priority**: P0
- **Depends On**: Task 6
- **Description**: 
  - 在 `src-tauri/src/core/` 创建或更新加密模块
  - 实现 XOR 加密
  - 实现 AES-128 加密
  - 实现 AES-256-GCM 加密（推荐）
  - 实现 ChaCha20 加密
  - 实现密钥派生（从密码派生）
- **Acceptance Criteria Addressed**: [FR-5, NFR-4]
- **Test Requirements**:
  - `programmatic` TR-8.1: 加密模块编译无错误
  - `programmatic` TR-8.2: 所有加密算法正确实现
  - `programmatic` TR-8.3: 密钥派生正确工作
- **Notes**: 使用成熟的 Rust 加密库（如 aes、chacha20poly1305 等）

## [ ] Task 9: 实现 Rust 传输编码模块
- **Priority**: P0
- **Depends On**: Task 6
- **Description**: 
  - 在 `src-tauri/src/core/` 创建或更新编码模块
  - 实现 Base64 编码
  - 实现 Hex 编码
  - 实现 URL Encode 编码
  - 实现无编码选项
- **Acceptance Criteria Addressed**: [FR-6, AC-8]
- **Test Requirements**:
  - `programmatic` TR-9.1: 编码模块编译无错误
  - `programmatic` TR-9.2: 所有编码算法正确实现
- **Notes**: 确保编码与加密作为独立步骤

## [ ] Task 10: 实现洋葱模型生成流水线
- **Priority**: P0
- **Depends On**: Task 6, Task 7, Task 8, Task 9
- **Description**: 
  - 重构 `src-tauri/src/core/generator.rs`
  - 实现模板组装步骤（根据生成模式+脚本类型）
  - 实现代码混淆步骤
  - 实现载荷加密步骤
  - 实现传输编码步骤
  - 实现最终组装步骤（解密 Stub + 密文 + 通信接收逻辑）
  - 更新 Tauri 命令 `generate_payload_cmd`
  - 添加生成日志输出
- **Acceptance Criteria Addressed**: [FR-9, FR-11, AC-5, AC-8, AC-10]
- **Test Requirements**:
  - `programmatic` TR-10.1: 生成流水线编译无错误
  - `programmatic` TR-10.2: 洋葱模型步骤按正确顺序执行
  - `programmatic` TR-10.3: 生成日志正确输出
- **Notes**: 保持向后兼容现有命令接口

## [ ] Task 11: 更新模板系统
- **Priority**: P1
- **Depends On**: Task 10
- **Description**: 
  - 更新或创建新的载荷模板
  - 为经典回显模式创建模板（所有6种脚本类型）
  - 为流量代理模式创建模板（可选，视需求而定）
  - 实现解密 Stub 模板
  - 实现通信接收逻辑模板
- **Acceptance Criteria Addressed**: [FR-1, FR-2, FR-9]
- **Test Requirements**:
  - `programmatic` TR-11.1: 模板系统正确加载
  - `human-judgement` TR-11.2: 生成的代码功能完整
- **Notes**: 流量代理模式可后续实现

## [ ] Task 12: 集成测试和联调
- **Priority**: P0
- **Depends On**: Task 5, Task 11
- **Description**: 
  - 测试完整的端到端生成流程
  - 测试所有脚本类型生成
  - 测试所有加密算法
  - 测试所有编码方式
  - 测试所有混淆级别
  - 测试 ASP 特殊混淆
  - 测试下载和复制功能
  - 测试异常情况处理
- **Acceptance Criteria Addressed**: [AC-11, AC-12, NFR-1, NFR-5, NFR-6]
- **Test Requirements**:
  - `programmatic` TR-12.1: 端到端流程无崩溃
  - `programmatic` TR-12.2: 所有配置组合正常工作
  - `programmatic` TR-12.3: 异常情况正确处理
  - `human-judgement` TR-12.4: 生成响应时间 < 2 秒
- **Notes**: 确保与现有载荷历史模块的集成
