# FG-ABYSS 项目使用总览

**版本**: v1.0.0  
**更新时间**: 2024-01-XX

---

## 📖 目录

1. [项目简介](#项目简介)
2. [快速开始](#快速开始)
3. [核心功能](#核心功能)
4. [使用指南](#使用指南)
5. [最佳实践](#最佳实践)
6. [常见问题](#常见问题)
7. [技术架构](#技术架构)

---

## 项目简介

FG-ABYSS 是一款**下一代 WebShell 管理工具**，完整复刻了"哥斯拉"的核心功能，并进行了多项增强：

### 核心特性

✅ **完整的 WebShell 管理** - 支持 PHP/ASP/ASPX/JSP  
✅ **命令执行** - 实时终端模拟器  
✅ **文件管理** - 上传/下载/编辑/权限  
✅ **数据库管理** - MySQL/PostgreSQL/SQLite/MSSQL  
✅ **Payload 生成** - 多语言模板 + 代码混淆  
✅ **批量操作** - 导入/导出/测试  
✅ **代理支持** - HTTP/HTTPS/SOCKS5  
✅ **流量加密** - AES-256-GCM + HMAC-SHA256  
✅ **安全审计** - 完整的操作日志系统  
✅ **性能优化** - LRU 缓存 + 连接池  

### 技术栈

- **后端**: Go 1.21+ + Wails v3
- **前端**: Vue 3 + TypeScript + Naive UI v3
- **数据库**: SQLite + GORM
- **终端**: xterm.js

---

## 快速开始

### 1. 安装

```bash
# 克隆项目
git clone https://github.com/your-org/fg-abyss.git
cd fg-abyss

# 安装前端依赖
cd frontend
npm install
npm run build

# 安装后端依赖
cd ..
go mod download

# 构建应用
wails build

# 运行应用
wails dev
```

### 2. 首次使用

1. **创建第一个 WebShell 连接**
   - 点击左侧导航栏的 **WebShell 管理**
   - 点击 **新建 WebShell**
   - 填写连接信息（名称、URL、类型、密码）
   - 点击 **测试连接** 验证
   - 点击 **保存**

2. **连接 WebShell**
   - 在列表中选择刚创建的连接
   - 点击 **连接** 按钮
   - 进入终端界面执行命令

3. **管理文件**
   - 切换到 **文件管理** 标签
   - 浏览远程文件系统
   - 支持上传/下载/删除/编辑

---

## 核心功能

### 1. WebShell 管理

**功能列表**:
- 创建/编辑/删除 WebShell
- 连接测试
- 状态管理（正常/异常/离线）
- 快速重连
- 历史记录

**支持的 WebShell 类型**:
- PHP (Basic/Full)
- ASP
- ASPX
- JSP

**编码器**:
- Base64
- URL Encode
- Raw (不编码)

### 2. 命令执行

**终端功能**:
- 实时命令执行
- 命令输出显示
- 命令历史记录
- 快捷键支持（Ctrl+C, Ctrl+V）
- 系统信息获取

**支持的命令**:
- 系统命令（whoami, pwd, ls 等）
- 网络命令（ipconfig, ifconfig, netstat 等）
- 文件命令（cat, dir, cd 等）
- 自定义命令

### 3. 文件管理

**文件操作**:
- 📤 上传文件
- 📥 下载文件
- 🗑️ 删除文件
- ✏️ 重命名文件
- 🔐 修改权限
- 📝 编辑内容

**目录操作**:
- 目录浏览
- 创建目录
- 删除目录
- 路径导航

### 4. Payload 生成

**内置模板**:
- PHP Basic - 基础命令执行
- PHP Full - 完整功能
- ASP - ASP 基础
- ASPX - ASPX 完整
- JSP - JSP 完整

**混淆级别**:
- None - 无混淆
- Low - 变量名混淆
- Medium - + 字符串编码
- High - + 垃圾代码 + 控制流

**使用步骤**:
1. 点击 **Payload 生成**
2. 选择类型和功能
3. 设置密码和混淆级别
4. 预览代码
5. 生成并下载

### 5. 数据库管理

**支持数据库**:
- MySQL
- PostgreSQL
- SQLite
- SQL Server

**功能**:
- 连接管理
- SQL 查询执行
- 表结构查看
- 列信息展示
- 查询结果导出（CSV）

### 6. 批量操作

**导入导出**:
- JSON 格式
- CSV 格式
- XML 格式

**批量功能**:
- 批量导入 WebShell
- 批量导出 WebShell
- 批量删除
- 批量测试连接
- 批量更新状态

### 7. 代理设置

**支持代理**:
- HTTP
- HTTPS
- SOCKS5

**预设配置**:
- Tor (9050)
- Burp Suite (8080)
- OWASP ZAP (8090)

### 8. 流量加密

**加密算法**:
- AES-256-GCM
- HMAC-SHA256

**安全特性**:
- 时间戳防重放攻击
- 请求签名验证
- 可配置时间窗口

### 9. 安全审计

**审计级别**:
- INFO - 信息
- WARNING - 警告
- ERROR - 错误
- CRITICAL - 严重

**操作类型**:
- LOGIN/LOGOUT
- CREATE/UPDATE/DELETE
- QUERY/EXECUTE
- UPLOAD/DOWNLOAD
- CONNECT/DISCONNECT
- 等等...

---

## 使用指南

### 场景 1: 管理多个 WebShell

1. **批量导入**
   - 准备 CSV/JSON 格式的 WebShell 数据
   - 点击 **批量操作** > **批量导入**
   - 粘贴数据或上传文件
   - 验证数据
   - 导入

2. **批量测试**
   - 选择要测试的 WebShell
   - 点击 **批量操作** > **批量测试**
   - 查看测试结果

3. **分类管理**
   - 使用搜索功能快速定位
   - 使用状态标签筛选

### 场景 2: 生成安全的 Payload

1. **选择高级混淆**
   - 类型：PHP Full
   - 混淆级别：High
   - 编码器：Base64

2. **启用流量加密**
   - 点击 **流量加密**
   - 生成加密配置
   - 应用配置
   - 测试加密

3. **下载 Payload**
   - 预览代码
   - 生成
   - 下载文件

### 场景 3: 数据库管理

1. **添加连接**
   - 点击 **数据库管理**
   - 点击 **新建**
   - 填写连接信息
   - 测试连接
   - 保存

2. **执行查询**
   - 选择数据库连接
   - 切换到 **SQL 查询** 标签
   - 输入 SQL 语句
   - 执行
   - 查看结果

3. **查看表结构**
   - 切换到 **表结构** 标签
   - 选择表
   - 查看列信息

### 场景 4: 安全审计

1. **查看操作日志**
   - 点击 **安全审计**
   - 查看日志列表
   - 使用搜索功能定位

2. **导出日志**
   - 点击 **导出**
   - 选择保存路径
   - 导出为 JSON

3. **统计分析**
   - 查看统计面板
   - 分析操作分布
   - 识别异常行为

---

## 最佳实践

### 安全建议

1. **启用流量加密**
   - 所有生产环境连接都应启用加密
   - 定期更换加密密钥
   - 使用强密码策略

2. **配置审计日志**
   - 启用自动保存
   - 定期导出备份
   - 监控异常操作

3. **使用代理**
   - 敏感操作使用 Tor 代理
   - 配置代理链增加安全性
   - 定期测试代理连接

4. **Payload 混淆**
   - 生产环境使用 High 级别混淆
   - 定期更换 Payload
   - 避免使用默认密码

### 性能优化

1. **使用缓存**
   - 频繁访问的数据使用缓存
   - 配置合理的过期时间
   - 定期清理过期缓存

2. **连接池管理**
   - 配置合理的最大连接数
   - 及时关闭不用的连接
   - 监控连接池状态

3. **批量操作**
   - 使用批量导入代替单个添加
   - 批量测试代替逐个测试
   - 定期清理无效连接

### 组织管理

1. **命名规范**
   - 使用有意义的名称
   - 添加详细描述
   - 分类管理

2. **权限管理**
   - 限制敏感操作权限
   - 定期审查操作日志
   - 实施最小权限原则

3. **备份策略**
   - 定期导出配置
   - 备份审计日志
   - 保存重要 Payload

---

## 常见问题

### Q1: 无法连接 WebShell？

**检查步骤**:
1. URL 是否正确
2. 密码是否匹配
3. 编码器选择是否正确
4. 网络连接是否正常
5. 是否需要配置代理

**解决方案**:
- 使用测试连接功能验证
- 检查 WebShell 是否在线
- 查看错误日志

### Q2: Payload 无法使用？

**可能原因**:
1. 密码不匹配
2. 编码器选择错误
3. 服务器环境不支持
4. 被杀毒软件拦截

**解决方案**:
- 提高混淆级别
- 更换 Payload 类型
- 检查服务器环境
- 使用免杀技术

### Q3: 数据库连接失败？

**检查步骤**:
1. 数据库服务是否运行
2. 主机地址和端口是否正确
3. 用户名密码是否正确
4. 防火墙是否阻止
5. 数据库是否允许远程访问

**解决方案**:
- 使用测试连接功能
- 检查防火墙配置
- 配置数据库权限

### Q4: 审计日志不显示？

**可能原因**:
1. 服务未启动
2. 日志文件路径错误
3. 权限问题

**解决方案**:
- 检查服务状态
- 验证文件路径
- 检查文件权限

---

## 技术架构

### 整体架构

```
┌─────────────────────────────────────┐
│         Frontend (Vue 3)            │
│  ┌─────────────────────────────┐    │
│  │   Components (20+)          │    │
│  │   - WebShell Management     │    │
│  │   - Terminal                │    │
│  │   - File Manager            │    │
│  │   - Payload Generator       │    │
│  │   - Database Manager        │    │
│  │   - Batch Operations        │    │
│  │   - Proxy Settings          │    │
│  │   - Traffic Encryption      │    │
│  │   - Audit Logs              │    │
│  └─────────────────────────────┘    │
└─────────────────────────────────────┘
              ↕ (Wails Runtime)
┌─────────────────────────────────────┐
│         Backend (Go)                │
│  ┌─────────────────────────────┐    │
│  │   Handlers (10+)            │    │
│  │   - Connection Handler      │    │
│  │   - Command Handler         │    │
│  │   - File Handler            │    │
│  │   - Payload Handler         │    │
│  │   - Database Handler        │    │
│  │   - Batch Handler           │    │
│  │   - Proxy Handler           │    │
│  │   - Encryption Handler      │    │
│  │   - Audit Handler           │    │
│  └─────────────────────────────┘    │
│  ┌─────────────────────────────┐    │
│  │   Services (10+)            │    │
│  │   - HTTP Engine             │    │
│  │   - Connection Manager      │    │
│  │   - Payload Generator       │    │
│  │   - Database Service        │    │
│  │   - Batch Service           │    │
│  │   - Proxy Service           │    │
│  │   - Traffic Encryption      │    │
│  │   - Cache Service           │    │
│  │   - Audit Service           │    │
│  └─────────────────────────────┘    │
│  ┌─────────────────────────────┐    │
│  │   Entities                  │    │
│  │   - WebShell                │    │
│  │   - Connection              │    │
│  │   - Payload                 │    │
│  │   - Database                │    │
│  └─────────────────────────────┘    │
└─────────────────────────────────────┘
              ↕
┌─────────────────────────────────────┐
│         Database (SQLite)           │
│  - Projects                         │
│  - WebShells                        │
│  - Connections                      │
│  - Settings                         │
│  - Audit Logs                       │
└─────────────────────────────────────┘
```

### DDD 分层架构

```
Entity Layer (实体层)
    ↓
Service Layer (服务层)
    ↓
Handler Layer (处理器层)
    ↓
Wails Runtime
    ↓
Frontend (Vue 3)
```

### 数据流

```
用户操作 → Frontend Component → Wails Runtime → Handler → Service → Entity → Database
          ↑                                                                       ↓
          └───────────────────────────────────────────────────────────────────────┘
```

---

## 附录

### 快捷键

- `Ctrl + N` - 新建
- `Ctrl + S` - 保存
- `Ctrl + F` - 搜索
- `Ctrl + R` - 刷新
- `F5` - 刷新当前页面
- `Esc` - 关闭弹窗

### 文件格式

**WebShell 导入格式**:
- JSON: `webshells.json`
- CSV: `webshells.csv`
- XML: `webshells.xml`

**审计日志格式**:
- JSON: `audit_logs_YYYYMMDD.json`

**Payload 格式**:
- PHP: `shell.php`
- ASP: `shell.asp`
- ASPX: `shell.aspx`
- JSP: `shell.jsp`

### 默认配置

**数据库路径**: `data/app.db`  
**审计日志路径**: `data/audit_logs.json`  
**缓存大小**: 1000 条  
**连接池大小**: 100 个  
**加密时间窗口**: 300 秒（5 分钟）

---

**文档结束**

*FG-ABYSS Development Team*  
*2024-01-XX*
