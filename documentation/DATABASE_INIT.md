# 数据库初始化流程优化说明

## ✅ 已解决的问题

### 问题描述
系统在首次运行或未检测到现有数据库时，未能自动执行完整的数据库初始化操作。

### 解决方案

#### 1. **优化数据库初始化流程**

在 `backend/db/init.go` 中添加了以下关键步骤：

```go
// 1. 先执行删除 name 字段的迁移（如果存在）
fmt.Println("开始检查数据库结构...")
if err := MigrateRemoveWebShellNameField(db); err != nil {
    fmt.Printf("警告：name 字段迁移失败：%v\n", err)
    // 不返回错误，继续执行，因为可能是字段已经不存在
}

// 2. 自动迁移数据库表结构
fmt.Println("正在执行数据库表结构迁移...")
if err := db.AutoMigrate(&models.Project{}, &models.WebShell{}); err != nil {
    return nil, fmt.Errorf("failed to migrate database: %w", err)
}
fmt.Println("数据库表结构迁移完成")

// 3. 创建默认项目和示例数据
// ... (原有代码)
```

### 📋 完整的初始化流程

现在系统启动时会按以下顺序执行：

#### **步骤 1: 创建数据目录**
```go
os.MkdirAll("data", 0755)
```
- ✅ 如果 `data` 目录不存在，自动创建
- ✅ 设置正确的权限 (0755)

#### **步骤 2: 连接数据库**
```go
gorm.Open(sqlite.Open("data/app.db"), config)
```
- ✅ 如果数据库文件不存在，自动创建
- ✅ 配置 GORM 参数
- ✅ 启用日志记录

#### **步骤 3: 数据库优化**
```go
PRAGMA foreign_keys = ON
PRAGMA journal_mode = WAL
```
- ✅ 启用外键约束
- ✅ 启用 WAL 模式（提高并发性能）

#### **步骤 4: 执行结构迁移**
```go
MigrateRemoveWebShellNameField(db)
db.AutoMigrate(&models.Project{}, &models.WebShell{})
```
- ✅ 检查并删除 WebShell 表的 name 字段（如果存在）
- ✅ 自动创建或更新表结构
- ✅ 确保数据库结构与模型定义一致

#### **步骤 5: 创建默认数据**
```go
// 创建默认项目
db.Where("name = ?", "默认项目").First(&defaultProject)
if result.Error == gorm.ErrRecordNotFound {
    db.Create(&defaultProject)
}

// 创建示例 WebShell 数据
db.Model(&models.WebShell{}).Count(&webshellCount)
if webshellCount == 0 {
    db.Create(&exampleWebShells)
}
```
- ✅ 检查并创建默认项目（如果不存在）
- ✅ 检查并创建示例 WebShell 数据（如果不存在）
- ✅ 避免重复插入数据

### 🎯 不同场景的处理

#### **场景 1: 首次运行（全新安装）**
```
1. 创建 data 目录
2. 创建 app.db 数据库文件
3. 执行表结构迁移（创建所有表）
4. 删除 name 字段（如果存在）
5. 创建默认项目
6. 创建示例 WebShell 数据
```

#### **场景 2: 数据库已存在（旧版本升级）**
```
1. 打开现有数据库
2. 执行 name 字段删除迁移
3. 执行表结构更新（如果有变更）
4. 跳过默认项目和示例数据（已存在）
```

#### **场景 3: 数据库文件损坏**
```
1. 手动删除 data/app.db*
2. 重新启动应用
3. 自动执行场景 1 的完整流程
```

### 📊 日志输出示例

#### **首次运行**
```
开始检查数据库结构...
开始迁移：删除 WebShell 表的 name 字段...
name 字段不存在，无需迁移
正在执行数据库表结构迁移...
数据库表结构迁移完成
默认项目创建成功
示例 WebShell 数据创建成功
Database initialized successfully (Pure Go SQLite)
```

#### **后续运行**
```
开始检查数据库结构...
开始迁移：删除 WebShell 表的 name 字段...
name 字段不存在，无需迁移
正在执行数据库表结构迁移...
数据库表结构迁移完成
Database initialized successfully (Pure Go SQLite)
```

### 🔧 手动重置数据库

如果需要完全重置数据库：

```powershell
# 删除数据库文件
Remove-Item -Path "data\app.db*" -Force

# 重新启动应用
.\bin\FG-ABYSS.exe
```

### ✅ 验证清单

系统启动后，请验证：

- [ ] 应用正常启动
- [ ] 数据目录存在：`data/`
- [ ] 数据库文件存在：`data/app.db`
- [ ] WAL 文件存在：`data/app.db-wal`
- [ ] 共享内存文件存在：`data/app.db-shm`
- [ ] 能看到默认项目
- [ ] 能看到示例 WebShell 数据

### 📝 修改文件清单

| 文件 | 修改内容 | 目的 |
|------|---------|------|
| `backend/db/init.go` | 添加迁移和 AutoMigrate 调用 | 确保表结构正确创建 |
| `backend/db/migrate_remove_name.go` | 已存在 | 删除 name 字段 |

### 🚀 使用说明

#### **开发环境**
```powershell
# 带控制台窗口（方便调试）
go run .

# 或构建后运行
go build -o bin\FG-ABYSS.exe .
.\bin\FG-ABYSS.exe
```

#### **生产环境**
```powershell
# 无控制台窗口
go build -ldflags="-H=windowsgui" -o bin\FG-ABYSS.exe .
.\bin\FG-ABYSS.exe
```

#### **使用构建脚本**
```powershell
.\build.ps1
```

### 🎉 优化效果

1. **自动化**：无需手动执行数据库初始化
2. **智能化**：自动检测并处理不同场景
3. **容错性**：迁移失败不影响启动（仅警告）
4. **幂等性**：多次运行不会产生重复数据
5. **兼容性**：支持旧版本平滑升级

### 🔍 故障排查

#### **问题：应用启动失败**
**解决**：
```powershell
# 查看错误信息
go run .

# 删除数据库重新初始化
Remove-Item data\app.db* -Force
go run .
```

#### **问题：表结构不正确**
**解决**：
```powershell
# 强制重新迁移
Remove-Item data\app.db* -Force
.\bin\FG-ABYSS.exe
```

#### **问题：示例数据未创建**
**检查**：
```powershell
# 查看日志输出
# 应该看到"示例 WebShell 数据创建成功"

# 如果没有，检查 WebShell 表是否为空
# 如果表中有数据，不会创建示例数据
```

## 📚 相关文档

- [GORM AutoMigrate](https://gorm.io/docs/migration.html)
- [SQLite WAL Mode](https://sqlite.org/wal.html)
- [Wails Database Guide](https://v3.wails.io/guides/database)
