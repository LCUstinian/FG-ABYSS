# FG-ABYSS (非攻 - 渊渟)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](package.json)

## 项目介绍

FG-ABYSS (非攻 - 渊渟) 是一款基于 Rust + Tauri 架构开发的下一代跨平台桌面端 WebShell 管理与载荷生成工具。旨在提供高隐蔽性、强对抗能力及模块化的红队作战支持，对标并超越现有开源工具（如哥斯拉），解决传统工具在流量特征、混淆强度及扩展性上的痛点。

## 技术栈

- **核心框架**: Rust (后端逻辑), Tauri (桌面容器)
- **前端界面**: Vue 3 + TypeScript + Naive UI
- **数据存储**: SQLite (本地加密存储)
- **关键依赖**:
  - 混淆引擎: Tree-sitter (多语言 AST 解析)
  - 加密算法: AES-256-GCM (载荷), ChaCha20-Poly1305 (日志), Ed25519 (插件签名)
  - 代码编辑器: Monaco Editor
  - 终端模拟: xterm.js

## 功能特性

### 核心功能

1. **项目管理模块**
   - 项目生命周期管理（创建、编辑、删除、回收站）
   - WebShell 实例管理（创建、编辑、测试连接、状态监控）
   - 列表交互（列宽拖拽、分页、全文搜索）
   - 右键上下文菜单（控制、清理缓存、编辑、删除、导出）

2. **载荷工厂模块**
   - 智能载荷生成（文件类：通用型、代理型、混合型）
   - 脚本类型适配（PHP, JSP, ASPX, ASP）
   - 高级混淆引擎（L1 轻量级、L2 标准级、L3 高对抗级）
   - 流量伪装（JSON、Form-Data、图片上传、API 请求）
   - 载荷资产库（历史记录、模板管理）

3. **插件生态模块**
   - 内置插件集（命令执行、文件管理、数据库管理）
   - 本地插件加载（动态加载/卸载第三方插件）
   - 插件签名验证（Ed25519 数字签名）

4. **系统设置模块**
   - 外观定制（深色/浅色主题、国际化、全局强调色）
   - 网络连接配置（超时阈值、重试次数、User-Agent 池）
   - 代理配置（全局 HTTP/SOCKS5 代理，单 Shell 独立代理）

5. **首页模块**
   - 数据概览（项目总数、WebShell 存量、插件数量、载荷类型分布）
   - 风险热力图（Shell 存活率、最后活跃时间、WAF 拦截频次）
   - 系统监控（进程资源占用情况）
   - 快捷操作（新建项目、快速生成载荷）
   - 活动审计（最近的关键操作历史）

### 安全特性

- **内存安全**: 利用 Rust 语言特性，确保无内存泄漏及缓冲区溢出风险
- **数据加密**: 敏感数据使用 AES-256-GCM 加密存储，审计日志使用 ChaCha20-Poly1305 加密
- **插件安全**: 强制 Ed25519 数字签名验证，杜绝恶意代码注入
- **输入验证**: 防止 XSS 和 SQL 注入攻击
- **审计日志**: 记录所有敏感操作，确保可追溯性

### 性能特性

- **响应速度**: 界面操作响应时间 < 200ms
- **大数据处理**: 支持大规模列表（>1000 条）无卡顿渲染
- **资源占用**: 内存和 CPU 占用合理，运行稳定
- **启动速度**: 应用启动时间短，无明显等待时间
- **网络请求**: 优化网络请求，减少延迟

## 安装说明

### 前置条件

- **Rust**: 1.70.0 或更高版本
- **Node.js**: 18.0.0 或更高版本
- **npm**: 9.0.0 或更高版本
- **Tauri CLI**: 2.0.0 或更高版本

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/yourusername/fg-abyss.git
   cd fg-abyss
   ```

2. **安装依赖**
   ```bash
   npm install
   ```

3. **构建项目**
   ```bash
   npm run build
   ```

4. **运行开发服务器**
   ```bash
   npm run dev
   ```

5. **构建可执行文件**
   ```bash
   npm run tauri build
   ```

## 使用方法

### 首次启动

1. 启动应用后，进入首页模块，可以看到系统概览
2. 点击左侧导航栏的「项目」模块，创建一个新的项目
3. 在项目中添加 WebShell 实例，配置 URL、密码等信息
4. 点击「测试连接」按钮，验证 WebShell 是否可用
5. 进入「载荷」模块，生成所需的载荷

### 日常使用

- **项目管理**: 在「项目」模块中管理项目和 WebShell 实例
- **载荷生成**: 在「载荷」模块中生成和管理载荷
- **插件管理**: 在「插件」模块中管理内置插件和第三方插件
- **系统设置**: 在「设置」模块中配置系统参数和外观
- **系统监控**: 在首页查看系统资源占用情况和操作历史

## 项目结构

```
├── src/                  # 前端源码
│   ├── components/       # 组件
│   ├── composables/      # 组合式函数
│   ├── i18n/             # 国际化
│   ├── stores/           # 状态管理
│   ├── styles/           # 样式
│   ├── types/            # 类型定义
│   ├── utils/            # 工具函数
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── src-tauri/            # Tauri 后端源码
│   ├── capabilities/     # 权限配置
│   ├── icons/            # 应用图标
│   ├── src/              # Rust 源码
│   │   ├── cmd/          # 命令处理
│   │   ├── core/         # 核心逻辑
│   │   ├── templates/    # 载荷模板
│   │   └── types/        # 类型定义
│   └── tauri.conf.json   # Tauri 配置
├── docs/                 # 文档
├── package.json          # 前端依赖
├── tsconfig.json         # TypeScript 配置
└── vite.config.ts        # Vite 配置
```

## 贡献指南

1. **Fork 项目**
2. **创建分支** (`git checkout -b feature/AmazingFeature`)
3. **提交更改** (`git commit -m 'Add some AmazingFeature'`)
4. **推送到分支** (`git push origin feature/AmazingFeature`)
5. **打开 Pull Request**

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 免责声明

本工具仅供授权渗透测试和安全研究使用，严禁用于非法用途。使用本工具造成的一切后果由使用者自行承担，作者不承担任何责任。

## 联系方式

- **项目地址**: [https://github.com/yourusername/fg-abyss](https://github.com/yourusername/fg-abyss)
- **问题反馈**: [https://github.com/yourusername/fg-abyss/issues](https://github.com/yourusername/fg-abyss/issues)

---

**FG-ABYSS (非攻 - 渊渟)** - 下一代 WebShell 管理与载荷生成工具，为红队作战提供强大支持！