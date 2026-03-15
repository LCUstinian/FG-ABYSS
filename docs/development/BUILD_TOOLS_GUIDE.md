# FG-ABYSS 构建工具使用指南

> 📋 本文档详细说明项目中各种构建工具的区别和使用场景，帮助开发者选择正确的构建方式。

---

## 📦 构建工具总览

| 工具/脚本 | 用途 | 输出位置 | 推荐场景 | 推荐度 |
|-----------|------|----------|----------|--------|
| `build-dev.ps1` | 快速构建开发版 | `bin/dev/` | 本地开发测试 | ⭐⭐ |
| `build-prod.ps1` | 快速构建生产版 | `bin/prod/` | 生产发布准备 | ⭐⭐ |
| `wails3 dev` | Wails 开发模式 | 内存中 | 日常开发调试 | ⭐⭐⭐ |
| `wails3 build` | Wails 标准构建 | `build/bin/` | 标准构建 | ⭐⭐⭐ |
| `task dev` | Taskfile 开发 | `bin/dev/` | 推荐开发方式 | ⭐⭐⭐⭐ |
| `task build:prod` | Taskfile 生产构建 | `bin/prod/` | 推荐发布方式 | ⭐⭐⭐⭐⭐ |

---

## 🔧 PowerShell 脚本

### 1. build-dev.ps1（开发版构建）

**文件名**: `scripts/build-dev.ps1`（原 build-app.ps1）

**用途**: 快速构建开发版本到 `bin/dev/` 目录

**特点**:
- ✅ 包含调试符号
- ✅ 保留控制台窗口（便于调试）
- ✅ 构建速度快
- ❌ 文件较大（约 27.7 MB）
- ❌ 未优化性能

**输出**:
```
bin/dev/
└── FG-ABYSS-dev.exe    # 开发版可执行文件
```

**使用场景**:
- ✅ 开发完成后测试完整应用
- ✅ 验证前后端集成
- ✅ 打包前功能验证
- ❌ 不适合日常开发（使用 `wails3 dev`）
- ❌ 不适合生产发布

**使用方法**:
```powershell
# 从任何目录运行
powershell -File D:\Go\FG-ABYSS\scripts\build-dev.ps1

# 或在 scripts 目录运行
cd scripts
.\build-dev.ps1
```

**构建步骤**:
```
[1/5] 清理旧开发版文件
[2/5] 构建前端资源
[3/5] 验证前端文件
[4/5] 构建 Go 后端（带调试符号）
[5/5] 设置数据目录
```

---

### 2. build-prod.ps1（生产版构建）

**文件名**: `scripts/build-prod.ps1`

**用途**: 构建优化的生产版本到 `bin/prod/` 目录

**特点**:
- ✅ 移除调试符号（-s -w 标志）
- ✅ 隐藏控制台窗口（纯 GUI）
- ✅ 优化性能和文件大小
- ✅ 接近最终发布版本
- ❌ 构建时间稍长

**输出**:
```
bin/prod/
└── FG-ABYSS.exe    # 生产版可执行文件（约 20-25 MB）
```

**使用场景**:
- ✅ 生产发布前测试
- ✅ 性能测试
- ✅ 用户 beta 测试
- ✅ 创建安装程序前的验证
- ❌ 不适合日常开发

**使用方法**:
```powershell
# 从任何目录运行
powershell -File D:\Go\FG-ABYSS\scripts\build-prod.ps1
```

**构建步骤**:
```
[1/5] 清理旧生产版文件
[2/5] 构建前端资源
[3/5] 验证前端文件
[4/5] 构建 Go 后端（优化，-s -w -H=windowsgui）
[5/5] 设置数据目录
```

---

## 🚀 Wails 命令

### 3. wails3 dev（Wails 开发模式）

**命令**: `wails3 dev`

**用途**: Wails 官方开发模式，支持热重载

**特点**:
- ✅ 实时热重载（HMR）
- ✅ 自动重新编译
- ✅ 详细的错误输出
- ✅ 最快的开发迭代速度
- ❌ 不生成可执行文件
- ❌ 需要保持运行状态

**输出**:
```
内存中运行，不生成文件
```

**使用场景**:
- ✅ **日常开发（强烈推荐）**
- ✅ 前端 UI 开发
- ✅ 快速测试功能变更
- ✅ 调试和排错
- ❌ 不适合性能测试
- ❌ 不适合打包分发

**使用方法**:
```powershell
# 在项目根目录运行
wails3 dev

# 支持参数
wails3 dev -debug          # 调试模式
wails3 dev -noreload       # 禁用热重载
```

**工作流程**:
```
1. 启动 Wails 开发服务器
2. 监控文件变化
3. 自动重新编译前端和后端
4. 浏览器/窗口自动刷新
5. 按 Ctrl+C 停止
```

---

### 4. wails3 build（Wails 标准构建）

**命令**: `wails3 build`

**用途**: Wails 官方标准构建命令

**特点**:
- ✅ 生成标准可执行文件
- ✅ 遵循 Wails 最佳实践
- ✅ 跨平台支持
- ✅ 包含所有资源
- ❌ 默认输出到 build/bin/

**输出**:
```
build/bin/
└── fg-abyss.exe    # 标准构建的可执行文件
```

**使用场景**:
- ✅ 标准构建需求
- ✅ 跨平台测试
- ✅ Wails 项目验证
- ❌ 不符合项目目录规范

**使用方法**:
```powershell
# 标准构建
wails3 build

# 生产构建（优化）
wails3 build -ldflags="-s -w"

# 指定平台
wails3 build -platform windows/amd64
```

---

## ⭐ Taskfile 任务（推荐）

### 5. task dev（开发任务）

**任务名**: `task dev`

**用途**: 使用 Taskfile.yml 进行开发构建

**特点**:
- ✅ 遵循项目规范
- ✅ 输出到标准目录（bin/dev/）
- ✅ 包含所有必要步骤
- ✅ 可自定义和扩展
- ✅ 跨平台一致

**输出**:
```
bin/dev/
└── FG-ABYSS-dev.exe
```

**使用场景**:
- ✅ **推荐开发构建方式**
- ✅ 团队统一构建流程
- ✅ CI/CD 集成

**使用方法**:
```powershell
task dev
```

---

### 6. task build:prod（生产构建任务）

**任务名**: `task build:prod`

**用途**: 使用 Taskfile.yml 进行生产构建

**特点**:
- ✅ 遵循项目规范
- ✅ 输出到标准目录（bin/prod/）
- ✅ 完整的优化和验证
- ✅ 包含所有必要步骤
- ✅ 跨平台一致

**输出**:
```
bin/prod/
└── FG-ABYSS.exe
```

**使用场景**:
- ✅ **推荐生产构建方式**
- ✅ 正式发布
- ✅ CI/CD 流水线

**使用方法**:
```powershell
task build:prod
```

---

## 📊 对比总结

### 开发场景推荐

| 场景 | 推荐工具 | 理由 |
|------|----------|------|
| 日常编码 | `wails3 dev` ⭐⭐⭐ | 热重载，最快迭代 |
| 功能测试 | `task dev` ⭐⭐⭐⭐ | 完整构建，符合规范 |
| 集成验证 | `build-dev.ps1` ⭐⭐ | 独立可执行文件 |
| 性能测试 | `task build:prod` ⭐⭐⭐⭐⭐ | 优化后的版本 |

### 发布场景推荐

| 场景 | 推荐工具 | 理由 |
|------|----------|------|
| Beta 测试 | `build-prod.ps1` ⭐⭐⭐ | 快速构建测试版 |
| 正式发布 | `task build:prod` ⭐⭐⭐⭐⭐ | 标准流程，可重复 |
| 创建安装包 | `task package` ⭐⭐⭐⭐⭐ | 生成 MSI 安装程序 |
| 跨平台发布 | `wails3 build` ⭐⭐⭐ | 官方工具，跨平台 |

---

## 🎯 决策树

```
需要构建应用？
├─ 日常开发？
│  ├─ 是 → 使用 wails3 dev（热重载）
│  └─ 否 → 继续
│
├─ 需要可执行文件？
│  ├─ 开发测试 → 使用 task dev 或 build-dev.ps1
│  ├─ 生产发布 → 使用 task build:prod 或 build-prod.ps1
│  └─ 继续
│
├─ 团队项目？
│  ├─ 是 → 使用 Taskfile（task dev / task build:prod）
│  └─ 否 → 继续
│
└─ 快速测试？
   ├─ 是 → 使用 PowerShell 脚本
   └─ 否 → 使用 Wails 命令
```

---

## 📝 最佳实践

### 日常开发流程

```powershell
# 1. 启动开发模式（热重载）
wails3 dev

# 2. 编码和测试
# ... 编写代码 ...
# 自动重新编译和刷新

# 3. 完成功能后，构建开发版验证
task dev

# 4. 运行测试
.\bin\dev\FG-ABYSS-dev.exe

# 5. 运行测试套件
task test
```

### 发布流程

```powershell
# 1. 运行所有测试
task test

# 2. 构建生产版本
task build:prod

# 3. 验证生产版本
.\bin\prod\FG-ABYSS.exe

# 4. 创建安装程序（可选）
task package

# 5. 打标签并发布
git tag v1.0.0
git push origin v1.0.0
```

---

## 🔍 常见问题

### Q1: 为什么有多个构建工具？

**A**: 
- **wails3 dev**: 官方开发工具，热重载最快
- **PowerShell 脚本**: 快速构建，适合特定场景
- **Taskfile**: 项目标准，团队统一
- 每种工具都有其最佳使用场景

### Q2: 应该使用哪个？

**A**: 
- **日常开发**: `wails3 dev`
- **构建测试**: `task dev`
- **生产发布**: `task build:prod`
- **快速验证**: `build-dev.ps1` 或 `build-prod.ps1`

### Q3: build-dev.ps1 和 build-prod.ps1 的区别？

**A**:
| 特性 | build-dev.ps1 | build-prod.ps1 |
|------|---------------|----------------|
| 输出 | bin/dev/FG-ABYSS-dev.exe | bin/prod/FG-ABYSS.exe |
| 调试符号 | ✅ 包含 | ❌ 移除 |
| 文件大小 | ~27.7 MB | ~20-25 MB |
| 控制台 | ✅ 显示 | ❌ 隐藏 |
| 优化 | ❌ 无 | ✅ 完整优化 |
| 构建速度 | 快 | 稍慢 |

### Q4: wails3 dev 和 task dev 的区别？

**A**:
- **wails3 dev**: 内存中运行，热重载，最快开发
- **task dev**: 生成可执行文件到 bin/dev/，完整构建

### Q5: 为什么不直接用 wails3 build？

**A**:
- `wails3 build` 输出到 `build/bin/`，不符合项目规范
- 项目规范要求输出到 `bin/dev/` 和 `bin/prod/`
- Taskfile 和 PowerShell 脚本确保符合规范

---

## 📚 相关文档

- [构建输出目录规范](build-output-specification.md) - 输出目录结构
- [scripts/README.md](../../scripts/README.md) - 脚本说明
- [Taskfile.yml](../../Taskfile.yml) - 任务配置
- [项目结构规范](project-optimization-specification.md) - 项目规范

---

**更新日期**: 2026-03-15  
**维护者**: FG-ABYSS Team
