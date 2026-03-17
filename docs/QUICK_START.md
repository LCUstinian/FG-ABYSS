# FG-ABYSS 快速开始指南

欢迎使用 FG-ABYSS！本指南将帮助您快速上手这款下一代 WebShell 管理工具。

---

## 📦 安装与构建

### 前置要求

**后端**:
- Go 1.21+
- GCC (MinGW for Windows)

**前端**:
- Node.js 18+
- npm 或 pnpm

### 安装步骤

1. **克隆项目**
```bash
git clone https://github.com/your-org/fg-abyss.git
cd fg-abyss
```

2. **安装前端依赖**
```bash
cd frontend
npm install
npm run build
```

3. **安装后端依赖**
```bash
cd ..
go mod download
```

4. **构建应用**
```bash
wails build
```

5. **运行应用**
```bash
wails dev
```

---

## 🚀 快速上手

### 1. 创建第一个 WebShell 连接

1. 点击左侧导航栏的 **WebShell 管理**
2. 点击 **新建 WebShell** 按钮
3. 填写连接信息：
   - **名称**: 自定义连接名称
   - **URL**: WebShell 地址（如：http://example.com/shell.php）
   - **类型**: 选择脚本类型（PHP/ASP/ASPX/JSP）
   - **密码**: 连接密码
   - **编码器**: 选择编码方式（Base64/URL/Raw）
4. 点击 **测试连接** 验证配置
5. 点击 **保存** 创建连接

### 2. 使用命令执行功能

1. 在 WebShell 列表中选择一个连接
2. 点击 **连接** 按钮进入终端界面
3. 在命令行中输入命令（如：`whoami`, `pwd`, `ls`）
4. 查看命令执行结果

### 3. 文件管理

1. 在终端界面切换到 **文件管理** 标签
2. 浏览远程服务器文件系统
3. 支持的操作：
   - 📤 上传文件
   - 📥 下载文件
   - 🗑️ 删除文件
   - ✏️ 重命名文件
   - 🔐 修改权限

### 4. 生成 Payload

1. 点击左侧导航栏的 **Payload 生成**
2. 配置 Payload 参数：
   - **类型**: PHP/ASP/ASPX/JSP
   - **功能**: Basic（基础）/ Full（完整）
   - **密码**: 连接密码
   - **混淆级别**: None/Low/Medium/High
3. 点击 **预览代码** 查看生成的代码
4. 点击 **生成 Payload** 下载文件

### 5. 数据库管理

1. 点击左侧导航栏的 **数据库管理**
2. 添加数据库连接：
   - **类型**: MySQL/PostgreSQL/SQLite/MSSQL
   - **主机**: 数据库地址
   - **端口**: 数据库端口
   - **用户名/密码**: 认证信息
3. 执行 SQL 查询
4. 查看表结构和数据

### 6. 批量操作

1. 点击左侧导航栏的 **批量操作**
2. 选择操作类型：
   - **批量导入**: 导入 JSON/CSV/XML 格式的 WebShell 数据
   - **批量导出**: 导出选中的 WebShell
   - **批量删除**: 删除多个 WebShell
   - **批量测试**: 测试多个 WebShell 连接

### 7. 配置代理

1. 点击左侧导航栏的 **代理设置**
2. 配置代理参数：
   - **类型**: HTTP/HTTPS/SOCKS5
   - **主机**: 代理服务器地址
   - **端口**: 代理端口
   - **认证**: 用户名/密码（可选）
3. 点击 **测试连接** 验证代理
4. 点击 **应用配置** 启用代理

---

## 💡 使用技巧

### Payload 混淆建议

- **无混淆**: 用于测试环境
- **低级混淆**: 变量名混淆，适合基础场景
- **中级混淆**: + 字符串编码，推荐用于生产环境
- **高级混淆**: + 垃圾代码 + 控制流平坦化，最高安全性

### 数据库查询技巧

1. 使用 `LIMIT` 限制返回行数，提高性能
2. 导出查询结果为 CSV 格式
3. 查看表结构了解字段信息

### 批量导入格式

**JSON 格式示例**:
```json
[
  {
    "name": "示例 Shell",
    "url": "http://example.com/shell.php",
    "type": "php",
    "password": "password",
    "encoder": "base64",
    "description": "描述信息"
  }
]
```

**CSV 格式示例**:
```csv
name,url,type,password,encoder,description
示例 Shell,http://example.com/shell.php,php,password,base64,描述信息
```

### 代理预设

快速加载常用代理配置：
- **Tor**: SOCKS5 127.0.0.1:9050
- **Burp Suite**: HTTP 127.0.0.1:8080
- **OWASP ZAP**: HTTP 127.0.0.1:8090

---

## 🔧 常见问题

### Q1: 无法连接 WebShell？

**检查项**:
1. URL 是否正确
2. 密码是否匹配
3. 编码器选择是否正确
4. 网络连接是否正常
5. 是否需要配置代理

### Q2: Payload 生成后无法连接？

**可能原因**:
1. 密码不匹配
2. 编码器选择错误
3. Payload 被杀毒软件拦截（尝试提高混淆级别）
4. 服务器环境不支持

### Q3: 数据库连接失败？

**检查项**:
1. 数据库服务是否运行
2. 主机地址和端口是否正确
3. 用户名密码是否正确
4. 防火墙是否阻止连接
5. 数据库是否允许远程访问

### Q4: 批量导入失败？

**检查项**:
1. 数据格式是否正确
2. 必填字段是否完整
3. 字段值是否合法
4. 先使用 **验证数据** 功能检查

---

## 📚 更多资源

### 文档
- [开发文档](./DEVELOPMENT_PLAN.md)
- [进度总结](./PROGRESS_SUMMARY.md)
- [完成报告](./COMPLETION_REPORT.md)

### 技术支持
- GitHub Issues: 提交问题和功能请求
- 邮件支持：support@fg-abyss.com

---

## ⚠️ 免责声明

本工具仅供安全研究和授权测试使用。请勿用于非法用途。使用本工具进行未授权的系统访问是违法行为。

---

**祝您使用愉快！**

*FG-ABYSS Team*
