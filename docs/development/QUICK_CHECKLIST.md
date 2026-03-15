# 项目规范快速检查清单

> 打印此清单并贴在显眼位置，或保存为书签以便快速参考

## 📋 日常开发检查清单

### 创建新文件时

- [ ] **文件命名**：使用小写字母 + 连字符（Go 用下划线）
- [ ] **文件位置**：是否放在正确的目录
- [ ] **文档文件**：必须放在 `docs/<category>/` 目录
- [ ] **测试文件**：与被测文件同目录，命名为 `*_test.go` 或 `*.test.ts`

```bash
# ✅ 正确示例
docs/testing/unit-test-report.md
internal/app/services/project_service.go
internal/app/services/project_service_test.go
frontend/src/utils/formatTime.test.ts

# ❌ 错误示例
/REPORT.md                          # 文档在根目录
/internal/app/ProjectService.go     # 命名错误
/tests/service_test.go              # 测试位置错误
```

---

### 编写代码时

#### Go 代码
- [ ] 导出函数有注释说明
- [ ] 错误被正确处理（不忽略）
- [ ] 使用 `fmt.Errorf("context: %w", err)` 包装错误
- [ ] 遵循分层架构（Handler → Service → Repository）

#### Vue/TypeScript 代码
- [ ] 组件使用 PascalCase 命名
- [ ] Props 有完整的类型定义
- [ ] 事件使用 kebab-case 命名
- [ ] 不使用 `any` 类型

---

### 提交代码前

- [ ] **测试运行通过**
  ```bash
  task test              # Go
  cd frontend && npm run test  # Frontend
  ```
  
- [ ] **测试覆盖率达标**
  - Service 层 ≥ 40%
  - Entity 验证 ≥ 80%
  
- [ ] **Commit Message 格式正确**
  ```bash
  # ✅ 正确
  git commit -m "feat(webshell): add batch delete support"
  git commit -m "fix(repository): correct GORM model"
  
  # ❌ 错误
  git commit -m "fixed bug"
  git commit -m "update code"
  ```
  
- [ ] **没有调试代码**（删除 console.log、fmt.Println 调试语句）
- [ ] **没有敏感信息**（密码、密钥、Token）
- [ ] **文档已更新**（如有必要）

---

### 创建文档时

- [ ] 文件在 `docs/<category>/` 目录下
- [ ] 使用正确的分类：
  - `docs/architecture/` - 架构设计
  - `docs/development/` - 开发指南
  - `docs/testing/` - 测试报告
  - `docs/api/` - API 文档
  - `docs/deployment/` - 部署文档
  - `docs/reports/` - 临时报告
  
- [ ] 文件名使用小写 + 连字符
- [ ] 文档包含：标题、概述、正文、总结、参考资料

```bash
# ✅ 创建文档的正确方式
mkdir -p docs/testing/bugfixes
touch docs/testing/bugfixes/webshell-restore-fix.md

# ❌ 错误方式
touch BUGFIX_REPORT.md  # 在根目录创建
touch docs/BUGFIX.md    # 没有分类
```

---

## 🎯 Commit Message 快速参考

### 格式
```
<type>(<scope>): <subject>
```

### Type 类型选择

| 场景 | Type | 示例 |
|------|------|------|
| 新功能 | `feat` | `feat(webshell): add search feature` |
| Bug 修复 | `fix` | `fix(repository): resolve null pointer` |
| 文档更新 | `docs` | `docs(readme): update installation` |
| 代码格式 | `style` | `style(format): fix indentation` |
| 重构 | `refactor` | `refactor(service): simplify logic` |
| 测试 | `test` | `test(service): add CRUD tests` |
| 配置 | `chore` | `chore(deps): update dependencies` |

### Scope 选择

```bash
# 后端
feat(service): ...
fix(repository): ...
refactor(handler): ...

# 前端
feat(component): ...
fix(utils): ...

# 全局
docs(readme): ...
chore(deps): ...
```

### Subject 规则

- ✅ 使用动词原形（create, not created）
- ✅ 首字母小写
- ✅ 不加句号
- ✅ 50 字符以内

---

## 📁 项目结构速查

```
FG-ABYSS/
├── docs/                    # 📚 所有文档（除 README）
│   ├── architecture/       # 架构
│   ├── development/        # 开发
│   ├── testing/           # 测试
│   └── api/               # API
├── internal/               # 业务代码
│   ├── app/
│   │   ├── handlers/      # 请求处理
│   │   └── services/      # 业务逻辑
│   ├── domain/
│   │   ├── entity/        # 实体
│   │   └── repository/    # 仓库接口
│   └── infrastructure/
│       └── repositories/  # 仓库实现
├── frontend/src/
│   ├── components/        # 组件
│   ├── views/            # 页面
│   └── utils/            # 工具
└── tests/                # 测试
    ├── fixtures/         # 测试数据
    └── integration/      # 集成测试
```

**分层职责**：
- **Handler**: 接收请求 → 验证参数 → 调用 Service → 返回结果
- **Service**: 业务逻辑 → 事务管理 → 调用 Repository
- **Repository**: 数据库操作 → ORM 映射

---

## 🧪 测试命令速查

```bash
# Go 测试
task test                    # 运行所有测试
task test:coverage           # 显示覆盖率
go test ./... -v             # 详细输出
go test ./... -cover         # 显示覆盖率

# 前端测试
cd frontend
npm run test                 # 运行测试
npm run test:watch           # 监视模式
npm run test:coverage        # 生成报告

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## ⚠️ 常见错误及修正

### 1. 文档位置错误
```bash
# ❌ 错误
/PROJECT_REPORT.md

# ✅ 修正
mv PROJECT_REPORT.md docs/testing/project-report.md
```

### 2. 文件命名错误
```bash
# ❌ 错误
/internal/app/ProjectService.go

# ✅ 修正
mv /internal/app/ProjectService.go /internal/app/services/project_service.go
```

### 3. Commit 格式错误
```bash
# ❌ 错误
git commit -m "fixed the bug"

# ✅ 修正
git commit -m "fix(webshell): resolve restore function error"
```

### 4. 测试位置错误
```bash
# ❌ 错误
/tests/project_test.go

# ✅ 修正
mv /tests/project_test.go /internal/app/services/project_service_test.go
```

---

## 🔧 实用命令

### 创建标准化文档
```bash
# 创建 Bug 修复报告
mkdir -p docs/testing/bugfixes
touch docs/testing/bugfixes/<issue-name>.md

# 创建架构文档
mkdir -p docs/architecture
touch docs/architecture/<design-name>.md
```

### 重命名文件
```bash
# 批量重命名为小写 + 连字符
rename 'y/A-Z/a-z/' *.MD
rename 's/_/-/g' *.md
```

### 检查项目结构
```bash
# 查看文档目录结构
tree docs/ -L 2

# 查找根目录的.md 文件（应该只有 README）
ls *.md
```

---

## 📊 代码审查清单

提交 PR 前逐项检查：

### 代码质量
- [ ] 遵循命名规范
- [ ] 无重复代码（DRY 原则）
- [ ] 函数职责单一（SRP 原则）
- [ ] 错误处理完整

### 测试
- [ ] 添加了单元测试
- [ ] 测试覆盖率达到要求
- [ ] 所有测试通过

### 文档
- [ ] 代码注释完整
- [ ] 更新了相关文档
- [ ] Commit message 规范

### 安全
- [ ] 无敏感信息
- [ ] 无硬编码密码/密钥
- [ ] 输入验证完整

---

## 💡 最佳实践提示

### 1. 文档编写
- 及时更新，不要堆积
- 使用统一的模板结构
- 包含代码示例和截图

### 2. 代码提交
- 小步提交，频繁提交
- 一个提交只做一件事
- 写好 Commit message

### 3. 测试编写
- 先写测试（TDD）
- 测试命名清晰
- 覆盖边界条件

### 4. 代码组织
- 遵循分层架构
- 模块职责清晰
- 避免循环依赖

---

## 📞 需要帮助？

- 📖 完整规范：[docs/development/project-optimization-specification.md](file://d:\Go\FG-ABYSS\docs\development\project-optimization-specification.md)
- ️ 技能助手：使用 `project-standard-checker` 技能
- 💬 团队讨论：在团队会议中提出

---

**版本**: v1.0  
**更新日期**: 2026-03-15  
**维护者**: FG-ABYSS Team
