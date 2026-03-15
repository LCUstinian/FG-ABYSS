# FG-ABYSS 项目优化规范

> 本文档规定了 FG-ABYSS 项目的开发规范和管理标准，旨在提升项目的整洁度、可维护性和团队协作效率。

## 📋 目录

1. [文档管理规范](#1-文档管理规范)
2. [文件命名规范](#2-文件命名规范)
3. [代码提交规范](#3-代码提交规范)
4. [项目结构规范](#4-项目结构规范)
5. [代码开发规范](#5-代码开发规范)
6. [测试规范](#6-测试规范)
7. [版本管理规范](#7-版本管理规范)

---

## 1. 文档管理规范

### 1.1 文档存储位置

**原则**：所有文档必须存放在 `docs/` 目录下，禁止在项目根目录创建文档文件。

```
FG-ABYSS/
├── docs/                    # ✅ 所有文档存放目录
│   ├── architecture/       # 架构设计文档
│   ├── development/        # 开发指南
│   ├── api/               # API 文档
│   ├── testing/           # 测试文档
│   ├── deployment/        # 部署文档
│   └── reports/           # 临时报告
├── README.md               # ✅ 允许：项目说明
├── LICENSE                 # ✅ 允许：许可证
├── .gitignore             # ✅ 允许：Git 配置
└── ...
```

### 1.2 文档分类规则

| 文档类型 | 存放路径 | 示例 |
|---------|---------|------|
| 架构设计文档 | `docs/architecture/` | `system-architecture.md` |
| 优化方案文档 | `docs/architecture/optimization/` | `performance-optimization.md` |
| 测试报告 | `docs/testing/` | `unit-test-report.md` |
| Bug 修复报告 | `docs/testing/bugfixes/` | `webshell-restore-fix.md` |
| 开发指南 | `docs/development/` | `getting-started.md` |
| API 文档 | `docs/api/` | `webshell-api.md` |
| 部署文档 | `docs/deployment/` | `production-deploy.md` |
| 临时报告 | `docs/reports/` | `weekly-report-2026-w11.md` |

### 1.3 文档模板

所有技术文档应包含以下基本结构：

```markdown
# 文档标题

> 简短描述（可选）

## 概述

背景、目的、范围

## 正文内容

根据文档类型组织内容

## 总结

关键要点、下一步行动

## 参考资料

相关链接、文档
```

---

## 2. 文件命名规范

### 2.1 通用命名规则

**基本原则**：
- ✅ 使用小写字母
- ✅ 使用连字符 `-` 分隔单词
- ✅ 语义化命名，清晰表达文件用途
- ❌ 禁止使用驼峰命名
- ❌ 禁止使用下划线分隔
- ❌ 禁止使用空格和特殊字符

### 2.2 各类型文件命名

#### Go 代码文件
```bash
# ✅ 正确
project_service.go
webshell_handler.go
project_repository_test.go

# ❌ 错误
ProjectService.go
project-service.go
Project_Service.go
```

#### Vue 组件文件
```bash
# ✅ 正确
ProjectsContent.vue
CreateProjectModal.vue
WebShellTable.vue

# ❌ 错误
projects-content.vue
projectsContent.vue
Projects-Content.vue
```

#### 测试文件
```bash
# ✅ 正确
project_service_test.go
formatTime.test.ts
webshell_handler.spec.ts

# ❌ 错误
test_project.go
ProjectTest.go
```

#### 文档文件
```bash
# ✅ 正确
unit-test-report.md
bugfix-webshell-restore.md
optimization-plan-2026.md

# ❌ 错误
UnitTestReport.md
unit_test_report.md
Unit-Test-Report.md
```

#### 配置文件
```bash
# ✅ 正确（保持现有）
package.json
Taskfile.yml
tsconfig.json
vite.config.ts

# ❌ 错误
Package.json
package-json.json
```

### 2.3 目录命名

```bash
# ✅ 正确
docs/architecture/
internal/app/services/
frontend/src/components/

# ❌ 错误
docs/Architecture/
internal/app/Services/
frontend/src/Components/
```

---

## 3. 代码提交规范

### 3.1 Commit Message 格式

遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 3.2 Type 类型

| 类型 | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat(webshell): add batch delete feature` |
| `fix` | Bug 修复 | `fix(project): fix restore function error` |
| `docs` | 文档更新 | `docs(readme): update installation guide` |
| `style` | 代码格式（不影响功能） | `style(format): fix indentation` |
| `refactor` | 重构 | `refactor(service): simplify validation logic` |
| `test` | 测试相关 | `test(service): add unit tests for CRUD` |
| `chore` | 构建/工具/配置 | `chore(deps): update testify to v1.8.0` |
| `perf` | 性能优化 | `perf(query): improve search performance` |
| `ci` | CI/CD 配置 | `ci(github): add automated test workflow` |

### 3.3 Scope 范围

使用受影响的模块或组件名称：

```bash
# 后端模块
feat(service): ...
fix(repository): ...
refactor(handler): ...

# 前端模块
feat(component): ...
fix(utils): ...
style(view): ...

# 全局
docs(readme): ...
chore(deps): ...
```

### 3.4 Subject 描述

- 使用祈使句、现在时态
- 首字母不要大写
- 结尾不加句号
- 简洁明了（50 字符以内）

```bash
# ✅ 正确
fix(webshell): restore deleted webshell functionality
feat(project): add soft delete support
test(service): add unit tests for CRUD operations

# ❌ 错误
fix(webshell): Fixed the restore function  # 过去时
fix(webshell): Fixes restore function      # 第三人称
fix(webshell): restore function.           # 有句号
fix: Fixed webshell restore function error # 太长且无 scope
```

### 3.5 Body 详细说明（可选）

对于复杂的提交，添加详细说明：

```bash
feat(webshell): add batch delete support

- Add BatchDelete method to WebShellService
- Update WebShellHandler with BatchDelete endpoint
- Add frontend batch selection UI
- Add unit tests for batch operations

Closes #123
```

### 3.6 Footer 关联问题（可选）

```bash
# 关联 Issue
Closes #123
Fixes #456
Refs #789

# 破坏性变更
BREAKING CHANGE: API endpoint changed from /v1 to /v2
```

### 3.7 完整示例

```bash
# 功能开发
feat(webshell): add search and pagination support

Implement search functionality for WebShell list with:
- URL search
- Remark search
- Payload search
- Server-side pagination
- Sort by multiple fields

Closes #45

# Bug 修复
fix(repository): correct GORM model specification in Recover method

- Add .Model(&entity.WebShell{}) to Recover()
- Add .Model(&entity.WebShell{}) to DeleteSoft()
- Add .Unscoped() to FindByID() for soft delete support

The missing Model specification caused "Table not set" error
when restoring deleted webshells.

Fixes #78

# 测试添加
test(service): add comprehensive unit tests for ProjectService

- Test all CRUD operations
- Test edge cases (not found, duplicate name)
- Test default project protection
- Achieve 40%+ coverage for service layer
```

---

## 4. 项目结构规范

### 4.1 标准项目结构

```
FG-ABYSS/
├── .git/                          # Git 仓库
├── .trae/                         # Trae IDE 配置
│   ├── rules/                     # 项目规则
│   └── skills/                    # 自定义技能
├── build/                         # 构建配置
│   ├── windows/                   # Windows 构建
│   ├── darwin/                    # macOS 构建
│   └── linux/                     # Linux 构建
├── cmd/                           # 应用程序入口
│   └── fg-abyss/                  # 主程序
├── configs/                       # 配置文件
│   ├── config.yaml                # 主配置
│   └── config.test.yaml           # 测试配置
├── docs/                          # 📚 所有文档（除 README 外）
│   ├── architecture/              # 架构文档
│   ├── development/               # 开发指南
│   ├── api/                       # API 文档
│   ├── testing/                   # 测试文档
│   ├── deployment/                # 部署文档
│   └── reports/                   # 临时报告
├── frontend/                      # 前端代码
│   ├── bindings/                  # Wails 绑定（自动生成）
│   ├── public/                    # 静态资源
│   ├── src/
│   │   ├── api/                   # API 调用
│   │   ├── components/            # Vue 组件
│   │   ├── views/                 # 页面视图
│   │   ├── utils/                 # 工具函数
│   │   ├── assets/                # 资源文件
│   │   └── locales/               # 国际化
│   ├── index.html                 # HTML 模板
│   ├── package.json               # 依赖配置
│   ├── tsconfig.json              # TypeScript 配置
│   ├── vite.config.ts             # Vite 配置
│   └── vitest.config.ts           # Vitest 配置
├── internal/                      # 私有业务代码
│   ├── app/                       # 应用层
│   │   ├── handlers/              # 请求处理器
│   │   ├── services/              # 业务逻辑层
│   │   └── events/                # 事件处理
│   ├── domain/                    # 领域层
│   │   ├── entity/                # 实体定义
│   │   └── repository/            # 仓库接口
│   └── infrastructure/            # 基础设施层
│       ├── database/              # 数据库
│       └── repositories/          # 仓库实现
├── pkg/                           # 公共库代码
├── scripts/                       # 脚本文件
├── tests/                         # 测试文件
│   ├── fixtures/                  # 测试数据
│   └── integration/               # 集成测试
├── .gitignore                     # Git 忽略配置
├── .traeignore                    # Trae 忽略配置
├── go.mod                         # Go 模块定义
├── go.sum                         # Go 依赖锁定
├── main.go                        # 程序入口
├── README.md                      # 项目说明（允许在根目录）
├── LICENSE                        # 许可证（允许在根目录）
├── Taskfile.yml                   # Task 任务配置
└── wails.json                     # Wails 配置
```

### 4.2 模块职责划分

#### internal/app/handlers/
- 处理外部请求（Wails 调用、HTTP 请求）
- 参数验证和转换
- 调用 Service 层
- 返回结果格式化

#### internal/app/services/
- 核心业务逻辑
- 事务管理
- 业务规则验证
- 调用 Repository 层

#### internal/domain/entity/
- 业务实体定义
- 实体验证逻辑
- 领域模型

#### internal/domain/repository/
- 仓库接口定义
- 数据访问契约

#### internal/infrastructure/repositories/
- 仓库接口实现
- 数据库操作
- ORM 映射

#### internal/infrastructure/database/
- 数据库连接管理
- 数据库初始化
- 迁移脚本

---

## 5. 代码开发规范

### 5.1 Go 代码规范

#### 项目结构
```go
// ✅ 正确：清晰的包结构
package services

import (
    "fg-abyss/internal/domain/entity"
    "fg-abyss/internal/domain/repository"
)

// ❌ 错误：混乱的导入
import (
    "../domain/entity"
    "../../repository"
)
```

#### 错误处理
```go
// ✅ 正确：完整的错误处理
func (s *ProjectService) Create(name string) (*entity.Project, error) {
    if err := validate(name); err != nil {
        return nil, fmt.Errorf("validate name: %w", err)
    }
    // ...
}

// ❌ 错误：忽略错误
func (s *ProjectService) Create(name string) (*entity.Project, error) {
    validate(name) // 错误未处理
    // ...
}
```

#### 注释规范
```go
// ✅ 正确：导出函数有注释
// CreateProject creates a new project with the given name
func (s *ProjectService) CreateProject(name string) (*entity.Project, error) {
    // ...
}

// ❌ 错误：无注释或注释不清晰
func (s *ProjectService) CreateProject(name string) (*entity.Project, error) {
    // 创建项目
    // ...
}
```

### 5.2 Vue 组件规范

#### 组件命名
```vue
<!-- ✅ 正确：多单词命名，PascalCase -->
<script setup lang="ts">
// ProjectsContent.vue
</script>

<!-- ❌ 错误：单单词或 kebab-case -->
<!-- Content.vue -->
<!-- projects-content.vue -->
```

#### Props 定义
```typescript
// ✅ 正确：完整的类型定义和默认值
interface Props {
  projectId?: string
  showDeleted?: boolean
  pageSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  projectId: '',
  showDeleted: false,
  pageSize: 10
})

// ❌ 错误：缺少类型定义
const props = defineProps(['projectId', 'showDeleted'])
```

#### 事件命名
```typescript
// ✅ 正确：使用 kebab-case 发送事件
emit('project-created', projectData)
emit('update:loading', false)

// ❌ 错误：使用驼峰
emit('projectCreated', projectData)
```

### 5.3 TypeScript 规范

#### 类型定义
```typescript
// ✅ 正确：使用 interface 定义对象类型
interface WebShell {
  id: string
  url: string
  projectId: string
  status?: string
}

// 使用 type 定义联合类型
type SortDirection = 'asc' | 'desc'

// ❌ 错误：使用 any
const data: any = getData()

// ✅ 正确：使用 unknown 或具体类型
const data: unknown = getData()
```

---

## 6. 测试规范

### 6.1 测试文件位置

```bash
# Go 测试：与被测文件同目录
internal/app/services/
├── project_service.go
└── project_service_test.go

# 前端测试：与被测文件同目录或使用 tests 目录
frontend/src/utils/
├── formatTime.ts
└── formatTime.test.ts
```

### 6.2 测试命名规范

```go
// ✅ 正确：Test<功能>_<场景>
func TestProjectService_Create_Success(t *testing.T)
func TestProjectService_Create_NameExists(t *testing.T)
func TestWebShell_Validate_URLEmpty(t *testing.T)

// ❌ 错误：命名不清晰
func TestProject(t *testing.T)
func TestCreate1(t *testing.T)
```

### 6.3 测试覆盖要求

| 层级 | 覆盖率目标 | 优先级 |
|------|-----------|--------|
| Service 层 | ≥ 40% | 🔴 高 |
| Entity 验证 | ≥ 80% | 🔴 高 |
| Handler 层 | ≥ 30% | 🟡 中 |
| Repository 层 | ≥ 50% | 🟡 中 |
| 前端组件 | 关键逻辑覆盖 | 🟢 低 |

### 6.4 测试运行命令

```bash
# Go 测试
task test                    # 运行所有测试
task test:coverage           # 显示覆盖率
go test ./... -v             # 详细输出

# 前端测试
cd frontend
npm run test                 # 运行测试
npm run test:watch           # 监视模式
npm run test:coverage        # 生成覆盖率报告
```

---

## 7. 版本管理规范

### 7.1 Git Branch 策略

```
main              # 主分支，生产环境
├── develop       # 开发分支
│   ├── feature/login          # 功能分支
│   ├── feature/webshell-crud  # 功能分支
│   └── bugfix/restore-error   # Bug 修复分支
└── release/v1.2.0 # 发布分支
```

### 7.2 分支命名

```bash
# 功能开发
feature/<feature-name>
feature/user-authentication

# Bug 修复
bugfix/<issue-description>
bugfix/webshell-restore-error

# 发布版本
release/<version>
release/v1.2.0

# 紧急修复
hotfix/<critical-issue>
hotfix/security-patch
```

### 7.3 标签管理

```bash
# 语义化版本
git tag v1.0.0       # 主版本.次版本.修订版本
git tag v1.2.0
git tag v1.2.3

# 预发布版本
git tag v1.0.0-alpha.1
git tag v1.0.0-beta.2
git tag v1.0.0-rc.1
```

---

## 8. 最佳实践

### 8.1 代码审查清单

在提交代码前检查：

- [ ] 代码符合命名规范
- [ ] 添加了必要的单元测试
- [ ] 测试通过且覆盖率达标
- [ ] Commit message 符合规范
- [ ] 文档已更新（如需要）
- [ ] 没有调试代码和注释掉的代码
- [ ] 没有敏感信息（密码、密钥等）
- [ ] 代码已通过本地测试

### 8.2 文档更新时机

以下情况必须更新文档：

1. 新增功能或修改 API
2. 架构变更
3. 配置项变更
4. 部署流程变更
5. Bug 修复涉及重要逻辑

### 8.3 定期维护任务

- [ ] 每周清理 `docs/reports/` 中的临时报告
- [ ] 每月更新依赖版本
- [ ] 每季度审查和优化测试覆盖率
- [ ] 每次发布前更新 CHANGELOG（如需要）

---

## 附录

### A. 快速参考卡片

```bash
# 创建文档
mkdir -p docs/<category>
touch docs/<category>/<descriptive-name>.md

# 提交代码
git add .
git commit -m "type(scope): description"
git push

# 运行测试
task test              # Go
npm run test           # Frontend

# 检查规范
# - 文件命名是否符合规范？
# - 文档是否在 docs/ 目录？
# - Commit message 是否规范？
# - 测试是否通过？
```

### B. 违规示例及修正

```bash
# ❌ 错误：文档在根目录
/project-root/TEST_REPORT.md

# ✅ 正确：文档在 docs 目录
/docs/testing/TEST_REPORT.md

# ❌ 错误：命名不规范
/internal/app/ProjectService.go
/internal/app/services/project_service_test.go

# ✅ 正确：命名规范
/internal/app/services/project_service.go
/internal/app/services/project_service_test.go

# ❌ 错误：Commit 不规范
git commit -m "fixed bug"
git commit -m "update code"

# ✅ 正确：Commit 规范
git commit -m "fix(webshell): restore deleted webshell functionality"
git commit -m "feat(service): add batch delete support"
```

---

## 版本历史

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|---------|
| v1.0 | 2026-03-15 | FG-ABYSS Team | 初始版本 |

---

## 参考资料

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Vue.js Style Guide](https://vuejs.org/style-guide/)
- [TypeScript Deep Dive](https://basarat.gitbook.io/typescript/)
