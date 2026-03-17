# 修复 TypeScript 编译错误

## 已修复的问题

1. ✅ PayloadGenerator.vue - 中文引号问题
2. ✅ DatabaseManager.vue - sslMode → ssl_mode
3. ✅ PayloadWorkspace.vue - 重复的 h 导入
4. ✅ ProxySettings.vue - has_auth → hasAuth
5. ✅ package.json - 添加缺失的依赖 (@vicons/carbon, @xterm/*)
6. ✅ 安装依赖

## 待修复的问题

### 1. Bindings 导入路径问题
所有使用 `@bindings/*` 的地方需要改为相对路径：
- `@bindings/PayloadHandler` → `../../bindings/fg-abyss/internal/app/handlers/payloadhandler.js`
- `@bindings/DatabaseHandler` → `../../bindings/fg-abyss/internal/app/handlers/databasehandler.js`
- `@bindings/ProxyHandler` → `../../bindings/fg-abyss/internal/app/handlers/proxyhandler.js`
- `@bindings/EncryptionHandler` → `../../bindings/fg-abyss/internal/app/handlers/encryptionhandler.js`
- `@bindings/BatchHandler` → `../../bindings/fg-abyss/internal/app/handlers/batchhandler.js`
- `@bindings/AuditHandler` → `../../bindings/fg-abyss/internal/app/handlers/audithandler.js`

### 2. 缺失的图标
- DatabaseManager.vue: AddOutline, PlayOutline, DatabaseOutline
- FileManager.vue: Refresh
- WebShellTerminal.vue: Disconnect

### 3. 类型错误
- DataTable render 函数类型不匹配
- 缺失 h 函数导入
- row.size 类型未知

### 4. xterm 主题配置
- selection → selectionBackground

## 快速修复方案

由于错误较多，建议使用以下命令跳过 TypeScript 检查直接构建：

```bash
# 临时方案：跳过类型检查
npm run build:skip-typecheck

# 或者修改 vite.config.ts 禁用类型检查
```

或者修复所有类型错误后再构建。
