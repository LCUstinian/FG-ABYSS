# FG-ABYSS Scripts Directory

> 📋 **重要**: 本项目统一使用 Taskfile.yml 管理所有构建和开发任务。本目录的脚本仅供快速开发使用。

## 📁 目录内容

### 可用脚本

| 脚本 | 用途 | 输出位置 | 推荐度 | 替代方案 |
|------|------|----------|--------|----------|
| `build-dev.ps1` | 快速构建开发版 | `bin/dev/FG-ABYSS-dev.exe` | ⚠️ 可选 | `task dev` ✅ |
| `build-prod.ps1` | 快速构建生产版 | `bin/prod/FG-ABYSS.exe` | ⚠️ 可选 | `task build:prod` ✅ |
| `check-standards.ps1` | 项目规范检查 | - | ✅ 推荐 | 无替代 |

**说明**:
- ⭐ 推荐使用 Taskfile.yml 进行构建（符合项目规范）
- ⚠️ PowerShell 脚本仅用于快速开发测试

---

## 🚀 推荐用法

### 使用 Taskfile.yml（推荐）

**构建应用**:
```powershell
# 生产构建（推荐）
task build

# 开发构建（带调试信息）
task dev

# 打包安装程序
task package
```

**运行应用**:
```powershell
# 运行已构建的应用
task run

# 开发模式运行
task dev
```

**测试**:
```powershell
# 运行所有测试
task test

# 查看覆盖率
task test:coverage
```

### 使用 PowerShell 脚本（可选）

**快速构建**（仅开发使用）:
```powershell
.\scripts\build-app.ps1
```

**检查项目规范**:
```powershell
.\scripts\check-standards.ps1
```

---

## 📋 为什么推荐使用 Taskfile.yml？

### 优势

1. **跨平台**: Windows、macOS、Linux 统一使用
2. **标准化**: Wails 官方推荐的构建方式
3. **功能完整**: 包含所有必要的步骤（绑定生成、图标、前端构建等）
4. **可维护**: 单一配置文件，易于维护
5. **可扩展**: 轻松添加新任务

### Taskfile.yml vs PowerShell 脚本

| 特性 | Taskfile.yml | PowerShell 脚本 |
|------|--------------|-----------------|
| 跨平台 | ✅ 支持所有平台 | ❌ 仅 Windows |
| 官方推荐 | ✅ Wails 推荐 | ❌ 非官方 |
| 功能完整性 | ✅ 完整 | ⚠️ 部分 |
| 可维护性 | ✅ 高 | ⚠️ 中 |
| 扩展性 | ✅ 容易 | ⚠️ 一般 |
| 依赖管理 | ✅ 自动 | ❌ 手动 |

---

## 🔧 脚本说明

### build-app.ps1

**用途**: 快速构建 Windows 应用（开发使用）

**功能**:
- 清理旧文件
- 构建前端
- 编译 Go 后端
- 生成 Windows GUI 应用（无控制台）

**适用场景**:
- ✅ 快速本地开发构建
- ✅ 测试构建流程
- ❌ 不推荐用于生产发布

**输出**:
```
bin/FG-ABYSS.exe  (约 27.7 MB)
```

### check-standards.ps1

**用途**: 自动化检查项目规范遵守情况

**功能**:
- 检查文档位置
- 检查文件命名
- 检查测试文件
- 检查 Vue 组件命名
- 检查 Git 提交信息

**适用场景**:
- ✅ 提交前检查
- ✅ 代码审查
- ✅ 项目维护

**运行示例**:
```powershell
.\scripts\check-standards.ps1
```

---

## 📖 Taskfile.yml 完整任务列表

### 主要任务

```powershell
# 构建
task build          # 构建应用
task dev            # 开发模式
task package        # 打包安装程序

# 运行
task run            # 运行应用
task dev            # 开发模式运行

# 测试
task test           # 运行所有测试
task test:coverage  # 显示覆盖率
task test:watch     # 监视模式测试

# 服务器模式
task build:server   # 构建服务器版本
task run:server     # 运行服务器版本
task build:docker   # 构建 Docker 镜像
task run:docker     # 运行 Docker 容器
```

### 平台特定任务

```powershell
# Windows
task windows:build
task windows:package
task windows:run

# macOS
task darwin:build
task darwin:package
task darwin:run

# Linux
task linux:build
task linux:package
task linux:run
```

---

## 🎯 最佳实践

### 日常开发

```powershell
# 1. 开发模式运行（热重载）
task dev

# 2. 运行测试
task test

# 3. 检查规范
.\scripts\check-standards.ps1

# 4. 提交代码
git commit -m "feat: ..."
```

### 生产发布

```powershell
# 1. 运行所有测试
task test

# 2. 构建应用
task build

# 3. 打包安装程序
task package

# 4. 验证构建
.\bin\FG-ABYSS.exe
```

### 代码审查

```powershell
# 运行规范检查
.\scripts\check-standards.ps1

# 修复发现的问题
# ...

# 重新检查
.\scripts\check-standards.ps1
```

---

## 📚 相关文档

- [Taskfile.yml](../Taskfile.yml) - 主任务配置文件
- [build/Taskfile.yml](../build/Taskfile.yml) - 通用构建任务
- [build/windows/Taskfile.yml](../build/windows/Taskfile.yml) - Windows 特定任务
- [项目规范](../docs/development/project-optimization-specification.md) - 完整项目规范

---

## 🤝 贡献指南

### 添加新脚本

如需添加新脚本，请遵循以下规范：

1. **命名**: 使用小写 + 连字符
   ```powershell
   ✅ my-script.ps1
   ❌ MyScript.ps1
   ❌ my_script.ps1
   ```

2. **注释**: 添加清晰的头部注释
   ```powershell
   # Script Name
   # Description: What this script does
   # Usage: powershell -File script.ps1
   ```

3. **错误处理**: 使用 `$LASTEXITCODE` 检查命令执行结果
   ```powershell
   command
   if ($LASTEXITCODE -ne 0) {
       Write-Host "Command failed!" -ForegroundColor Red
       exit 1
   }
   ```

4. **文档**: 在 README.md 中说明脚本用途和用法

### 修改现有脚本

- 保持向后兼容
- 更新文档
- 测试修改
- 提交时说明原因

---

## ⚠️ 注意事项

1. **脚本不是必需的**: 所有功能都可以通过 Taskfile.yml 完成
2. **仅 Windows 可用**: PowerShell 脚本只能在 Windows 上运行
3. **开发使用**: 脚本主要用于快速开发，生产环境请使用 Taskfile
4. **维护责任**: 使用脚本需自行确保与 Taskfile 同步

---

## 📞 需要帮助？

- 查看 [Taskfile.yml](../Taskfile.yml) 了解所有可用任务
- 运行 `task --list` 查看任务列表
- 参考 [项目规范文档](../docs/development/)

---

**更新日期**: 2026-03-15  
**维护者**: FG-ABYSS Team
