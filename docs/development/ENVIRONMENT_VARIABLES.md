# 环境变量配置指南

本文档说明 FG-ABYSS 项目的环境变量配置。

## 📋 配置步骤

1. 复制示例文件：
   ```bash
   cp .env.example .env
   ```

2. 根据实际需求修改 `.env` 文件中的值

3. 重启应用使配置生效

## 🔧 配置项说明

### 应用配置

| 变量名 | 说明 | 类型 | 默认值 | 示例 |
|--------|------|------|--------|------|
| `FG_APP_ENV` | 应用运行环境 | string | `development` | `development`, `production`, `test` |
| `FG_APP_VERSION` | 应用版本号 | string | `1.0.0` | `1.0.0`, `1.2.3` |
| `FG_APP_DEBUG` | 调试模式 | boolean | `true` | `true`, `false` |

**说明**：
- `development`: 开发环境，启用调试日志和热重载
- `production`: 生产环境，优化性能，减少日志
- `test`: 测试环境，用于运行自动化测试

### 服务器配置

| 变量名 | 说明 | 类型 | 默认值 | 示例 |
|--------|------|------|--------|------|
| `FG_SERVER_PORT` | 后端服务端口 | number | `9245` | `9245`, `8080` |
| `FG_SERVER_HOST` | 后端服务主机 | string | `localhost` | `localhost`, `0.0.0.0` |

**说明**：
- 端口号范围：1024-65535
- 主机设置为 `0.0.0.0` 可允许外部访问

### 数据库配置

| 变量名 | 说明 | 类型 | 默认值 | 示例 |
|--------|------|------|--------|------|
| `FG_DB_DRIVER` | 数据库驱动 | string | `sqlite` | `sqlite`, `mysql`, `postgres` |
| `FG_DB_PATH` | 数据库文件路径 | string | `data/app.db` | `data/app.db`, `/var/db/app.db` |

**说明**：
- 当前版本仅支持 SQLite
- 路径相对于项目根目录或绝对路径

### 日志配置

| 变量名 | 说明 | 类型 | 默认值 | 示例 |
|--------|------|------|--------|------|
| `FG_LOG_LEVEL` | 日志级别 | string | `debug` | `debug`, `info`, `warn`, `error` |
| `FG_LOG_FORMAT` | 日志格式 | string | `console` | `console`, `json` |
| `FG_LOG_OUTPUT` | 日志输出 | string | `stdout` | `stdout`, `file`, `all` |

**说明**：
- `debug`: 输出所有日志（开发推荐）
- `info`: 输出信息、警告、错误
- `warn`: 仅输出警告和错误（生产推荐）
- `error`: 仅输出错误

### Wails 配置

| 变量名 | 说明 | 类型 | 默认值 | 示例 |
|--------|------|------|--------|------|
| `FG_WAILS_PORT` | Wails 应用端口 | number | `9245` | `9245`, `8080` |
| `FG_WAILS_HOT_RELOAD` | 热重载 | boolean | `true` | `true`, `false` |

**说明**：
- 热重载功能仅在开发环境有效
- 端口号应与后端服务端口一致

## 📝 配置示例

### 开发环境配置

```bash
# .env
FG_APP_ENV=development
FG_APP_VERSION=1.0.0
FG_APP_DEBUG=true

FG_SERVER_PORT=9245
FG_SERVER_HOST=localhost

FG_DB_DRIVER=sqlite
FG_DB_PATH=data/app.db

FG_LOG_LEVEL=debug
FG_LOG_FORMAT=console
FG_LOG_OUTPUT=stdout

FG_WAILS_PORT=9245
FG_WAILS_HOT_RELOAD=true
```

### 生产环境配置

```bash
# .env.production
FG_APP_ENV=production
FG_APP_VERSION=1.0.0
FG_APP_DEBUG=false

FG_SERVER_PORT=9245
FG_SERVER_HOST=0.0.0.0

FG_DB_DRIVER=sqlite
FG_DB_PATH=/var/lib/fg-abyss/app.db

FG_LOG_LEVEL=warn
FG_LOG_FORMAT=json
FG_LOG_OUTPUT=file

FG_WAILS_PORT=9245
FG_WAILS_HOT_RELOAD=false
```

### 测试环境配置

```bash
# .env.test
FG_APP_ENV=test
FG_APP_VERSION=1.0.0
FG_APP_DEBUG=true

FG_SERVER_PORT=9246
FG_SERVER_HOST=localhost

FG_DB_DRIVER=sqlite
FG_DB_PATH=:memory:

FG_LOG_LEVEL=error
FG_LOG_FORMAT=console
FG_LOG_OUTPUT=stdout

FG_WAILS_PORT=9246
FG_WAILS_HOT_RELOAD=false
```

## 🔒 安全注意事项

### 不要提交敏感信息

- ✅ `.env` 文件已在 `.gitignore` 中忽略
- ✅ 只提交 `.env.example` 作为示例
- ❌ 不要将真实密码、密钥提交到版本控制

### 生产环境建议

1. 使用环境变量管理工具（如 dotenv-vault）
2. 定期更换敏感配置
3. 限制数据库访问权限
4. 使用强密码

## 🐛 故障排查

### 配置不生效

**问题**: 修改 `.env` 后配置不生效

**解决**:
1. 重启应用
2. 清理缓存：`rm -rf frontend/node_modules/.vite`
3. 重新构建：`npm run build`

### 端口冲突

**问题**: `Error: listen EADDRINUSE: address already in use`

**解决**:
```bash
# 查看端口占用
netstat -ano | findstr :9245

# 修改端口
FG_SERVER_PORT=9246
```

### 数据库路径错误

**问题**: `unable to open database file`

**解决**:
1. 确保路径存在：`mkdir -p data`
2. 检查权限：`chmod 755 data`
3. 使用绝对路径

## 📚 相关文档

- [项目启动指南](docs/development/getting-started.md)
- [部署指南](docs/deployment/production-deploy.md)
- [配置管理最佳实践](docs/development/config-best-practices.md)

## 🤝 贡献

添加新的环境变量配置时：

1. 更新 `.env.example`
2. 更新本文档
3. 确保向后兼容
4. 添加合理的默认值

---

**版本**: v1.0  
**更新日期**: 2026-03-15
