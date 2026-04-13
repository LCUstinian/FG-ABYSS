# FG-ABYSS 完整开发规范

## Why
基于 PRD、DEV、GUI 三个核心文档，启动 FG-ABYSS 项目的完整开发流程，实现从 0 到 1 的 WebShell 管理工具。

## What Changes
- 初始化 Tauri V2 + Vue 3 + TypeScript 项目结构
- 实现核心功能模块（载荷管理、项目管理、加密通信）
- 实现 UI 界面（自定义标题栏、导航、状态栏、各功能页面）
- 建立数据库存储系统（SQLite）
- 实现插件系统基础架构
- 实现控制台窗口系统（文件、数据库、终端插件）

**BREAKING**: 这是从 0 到 1 的完整项目开发，无历史代码需要兼容

## Impact
- Affected specs: 无（新项目）
- Affected code: 整个项目代码库

## ADDED Requirements

### Requirement: 项目初始化
系统 SHALL 使用 Tauri V2 + Vue 3 + TypeScript + Naive UI 技术栈初始化项目

#### Scenario: 成功初始化
- WHEN 开发者运行项目设置脚本
- THEN 项目结构、依赖、配置全部正确设置，可以运行 `pnpm tauri dev`

### Requirement: 载荷管理模块（本地）
系统 SHALL 提供载荷配置管理功能（CRUD），支持生成 PHP/JSP/ASP/ASPX 载荷代码

#### Scenario: 创建载荷配置
- WHEN 用户在 /payload 页面填写配置表单
- THEN 系统保存配置到 SQLite 数据库，并显示在列表中

#### Scenario: 生成载荷代码
- WHEN 用户选择载荷配置并点击"生成载荷"
- THEN 系统生成对应语言的 WebShell 代码，支持预览、复制、下载

### Requirement: 项目管理模块（远程）
系统 SHALL 提供项目管理功能，支持创建项目、添加 WebShell、连接测试

#### Scenario: 创建项目
- WHEN 用户点击新建项目按钮
- THEN 系统创建新项目，显示在项目树中

#### Scenario: 添加 WebShell
- WHEN 用户在项目中添加 WebShell
- THEN 系统保存 WebShell 信息（加密存储密码），支持连接测试

### Requirement: 加密通信模块（远程）
系统 SHALL 提供 AES-256-GCM 和 XOR 加密方式，支持密钥管理

#### Scenario: 连接 WebShell
- WHEN 用户右键 WebShell 选择"连接控制台"
- THEN 系统使用配置的加密方式建立安全连接，打开控制台窗口

### Requirement: 控制台窗口系统
系统 SHALL 提供独立控制台窗口，支持多 TAB 展示插件功能

#### Scenario: 打开控制台
- WHEN 用户右键 WebShell 选择"连接控制台"
- THEN 系统打开新窗口，显示文件、数据库、终端等 TAB

### Requirement: 文件管理插件（远程）
系统 SHALL 提供远程文件浏览、上传、下载、编辑功能

#### Scenario: 浏览文件
- WHEN 用户在控制台选择"文件"TAB
- THEN 系统显示远程服务器文件列表，支持目录导航

### Requirement: 数据库管理插件（远程）
系统 SHALL 提供远程 SQL 查询执行、结果展示功能

#### Scenario: 执行查询
- WHEN 用户在控制台选择"数据库"TAB 并输入 SQL
- THEN 系统执行查询并展示结果表格

### Requirement: 终端管理插件（远程）
系统 SHALL 提供远程命令执行功能

#### Scenario: 执行命令
- WHEN 用户在控制台选择"终端"TAB 并输入命令
- THEN 系统执行命令并显示输出

### Requirement: 设置页面
系统 SHALL 提供外观、连接、安全、日志、备份、关于等设置页面

#### Scenario: 修改主题
- WHEN 用户在设置页面选择深色主题
- THEN 系统立即切换为深色主题，并保存到本地存储

### Requirement: 软删除机制
系统 SHALL 对项目和 WebShell 实现软删除，支持回收站功能

#### Scenario: 删除项目
- WHEN 用户删除项目
- THEN 项目标记为已删除，移入回收站，支持恢复

## MODIFIED Requirements
无（新项目）

## REMOVED Requirements
无（新项目）
