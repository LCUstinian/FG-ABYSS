# FG-ABYSS 项目规范培训指南

> 本指南用于帮助团队成员快速理解和掌握项目规范，建议新成员入职时首先学习

## 📚 目录

1. [培训目标](#培训目标)
2. [学习路径](#学习路径)
3. [实践练习](#实践练习)
4. [考核标准](#考核标准)
5. [常见问题](#常见问题)

---

## 培训目标

完成本培训后，团队成员应能够：

- ✅ 理解并遵循文档管理规范
- ✅ 正确使用文件命名规则
- ✅ 编写规范的 Git 提交信息
- ✅ 按照项目结构组织代码
- ✅ 应用代码开发最佳实践
- ✅ 编写符合规范的单元测试

---

## 学习路径

### 第一阶段：基础规范（1-2 小时）

#### 1.1 文档管理规范（15 分钟）

**学习内容**：
- 文档存储位置规则
- 文档分类方法
- 文档模板结构

**关键要点**：
```bash
# ✅ 所有文档必须在 docs/ 目录下
docs/architecture/     # 架构设计
docs/development/      # 开发指南
docs/testing/          # 测试报告
docs/api/             # API 文档
docs/deployment/      # 部署文档
docs/reports/         # 临时报告

# ❌ 禁止在根目录创建文档
/REPORT.md  # 错误！
```

**练习**：
1. 查看现有文档位置
2. 将不合规的文档移动到正确位置
3. 创建一个新的测试报告文档

#### 1.2 文件命名规范（15 分钟）

**学习内容**：
- 通用命名规则
- 各类型文件命名
- 目录命名规则

**关键要点**：
```bash
# 通用规则：小写 + 连字符
✅ project-service.go
✅ webshell-handler.ts
✅ unit-test-report.md

# Go 文件：小写 + 下划线
✅ project_service.go
✅ project_service_test.go

# Vue 组件：PascalCase
✅ ProjectsContent.vue
✅ CreateProjectModal.vue
```

**练习**：
1. 找出项目中命名不规范的文件
2. 重命名这些文件
3. 更新相关导入语句

#### 1.3 Git 提交规范（30 分钟）

**学习内容**：
- Commit Message 格式
- Type 类型选择
- Scope 范围使用
- Subject 编写规则

**关键要点**：
```bash
# 格式：<type>(<scope>): <subject>

# Type 类型
feat     - 新功能
fix      - Bug 修复
docs     - 文档更新
style    - 代码格式
refactor - 重构
test     - 测试
chore    - 构建/配置

# 示例
✅ feat(webshell): add batch delete support
✅ fix(repository): correct GORM model specification
✅ test(service): add unit tests for CRUD

# ❌ 错误示例
❌ fixed bug              # 过去时，无 scope
❌ Update code            # 太 vague
❌ Fix: Fixed issue       # 过去时，有句号
```

**练习**：
1. 查看最近的提交历史
2. 找出不规范的提交信息
3. 重写这些提交信息（使用 `git commit --amend`）

#### 1.4 项目结构规范（30 分钟）

**学习内容**：
- 标准项目结构
- 模块职责划分
- 分层架构理解

**关键要点**：
```
# 分层架构
Handler     → 接收请求，参数验证，响应格式化
   ↓
Service     → 业务逻辑，事务管理，业务规则
   ↓
Repository  → 数据访问，ORM 操作
```

**练习**：
1. 绘制当前项目的架构图
2. 说明每个模块的职责
3. 找出职责不清晰的模块

---

### 第二阶段：代码规范（2-3 小时）

#### 2.1 Go 代码规范（1 小时）

**学习内容**：
- 包组织
- 错误处理
- 注释规范
- 测试编写

**关键要点**：
```go
// ✅ 错误处理
func (s *Service) Create(name string) (*entity.Project, error) {
    if err := validate(name); err != nil {
        return nil, fmt.Errorf("validate name: %w", err)
    }
    // ...
}

// ✅ 导出函数注释
// CreateProject creates a new project with validation
func (s *Service) CreateProject(name string) (*entity.Project, error) {
    // ...
}

// ✅ 测试命名
func TestProjectService_Create_Success(t *testing.T)
func TestProjectService_Create_NameExists(t *testing.T)
```

**练习**：
1. 审查现有 Go 代码
2. 改进错误处理
3. 补充缺失的注释
4. 编写单元测试

#### 2.2 Vue/TypeScript 规范（1 小时）

**学习内容**：
- 组件命名
- Props 定义
- 事件命名
- 类型定义

**关键要点**：
```typescript
// ✅ 组件命名
<script setup lang="ts">
// ProjectsContent.vue (PascalCase)
</script>

// ✅ Props 定义
interface Props {
  projectId?: string
  showDeleted?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  projectId: '',
  showDeleted: false
})

// ✅ 事件命名
emit('project-created', data)
emit('update:loading', false)

// ✅ 类型定义（不使用 any）
interface WebShell {
  id: string
  url: string
  projectId: string
}
```

**练习**：
1. 检查现有组件
2. 改进 Props 类型定义
3. 规范化事件命名

#### 2.3 测试规范（1 小时）

**学习内容**：
- 测试文件位置
- 测试命名
- 覆盖率要求
- 测试编写最佳实践

**关键要点**：
```bash
# 测试文件位置
✅ internal/app/services/project_service_test.go  # 同目录
✅ frontend/src/utils/formatTime.test.ts          # 同目录

# 测试命名
✅ TestProjectService_Create_Success
✅ TestWebShell_Validate_URLEmpty

# 覆盖率要求
Service 层：≥ 40%
Entity 验证：≥ 80%
Handler 层：≥ 30%
Repository: ≥ 50%
```

**练习**：
1. 运行测试覆盖率报告
2. 找出覆盖率低的模块
3. 补充测试用例

---

### 第三阶段：实战演练（2-3 小时）

#### 3.1 完整功能开发（2 小时）

**任务**：实现一个简单的 CRUD 功能

**要求**：
1. 创建 Entity 实体
2. 定义 Repository 接口
3. 实现 Repository
4. 编写 Service 层业务逻辑
5. 创建 Handler 处理请求
6. 编写完整的单元测试
7. 编写 API 文档

**检查点**：
- [ ] 文件命名规范
- [ ] 目录结构正确
- [ ] 代码分层清晰
- [ ] 错误处理完整
- [ ] 测试覆盖达标
- [ ] 文档齐全

#### 3.2 Bug 修复演练（1 小时）

**任务**：修复一个模拟的 Bug

**要求**：
1. 创建 Bug 报告文档
2. 编写复现测试
3. 修复 Bug
4. 验证修复
5. 更新相关文档

**检查点**：
- [ ] Bug 报告规范
- [ ] 提交信息规范
- [ ] 测试验证完整
- [ ] 文档更新及时

---

## 实践练习

### 练习 1：文档整理

**目标**：将混乱的文档整理规范

**步骤**：
```bash
# 1. 找出根目录的所有.md 文件
ls *.md

# 2. 创建正确的目录结构
mkdir -p docs/{architecture,development,testing,api,deployment,reports}

# 3. 移动文档到正确位置
mv PROJECT_REPORT.md docs/testing/project-report.md
mv BUGFIX.md docs/testing/bugfixes/bugfix-description.md
mv ARCHITECTURE.md docs/architecture/system-architecture.md

# 4. 重命名不规范的文件
rename 'y/A-Z/a-z/' docs/*/*.MD
rename 's/_/-/g' docs/*/*.md
```

**验收标准**：
- 根目录只有 README.md 和 LICENSE
- 所有文档在 docs/ 目录下
- 文件命名符合规范

### 练习 2：代码重构

**目标**：重构一段不符合规范的代码

**原始代码**：
```go
// ❌ 问题代码
package app

func CreateProject(n string)(*Project,error){
    if n==""{
        return nil,errors.New("name empty")
    }
    p:=new(Project)
    p.Name=n
    db.Create(p)  // 错误未处理
    return p,nil
}
```

**要求**：
1. 移动到正确的包
2. 重命名文件和函数
3. 改进错误处理
4. 添加注释
5. 编写测试

**参考答案**：
```go
// ✅ 重构后
// File: internal/app/services/project_service.go

package services

import (
    "fg-abyss/internal/domain/entity"
    "fg-abyss/internal/domain/repository"
)

// ProjectService handles project business logic
type ProjectService struct {
    repo repository.ProjectRepository
}

// CreateProject creates a new project with validation
func (s *ProjectService) CreateProject(name string) (*entity.Project, error) {
    if name == "" {
        return nil, fmt.Errorf("validate name: name cannot be empty")
    }
    
    project := &entity.Project{
        Name: name,
    }
    
    if err := s.repo.Save(project); err != nil {
        return nil, fmt.Errorf("save project: %w", err)
    }
    
    return project, nil
}
```

### 练习 3：提交信息重写

**目标**：将不规范的提交信息改写为规范格式

**原始提交**：
```bash
"fixed bug"
"update code"
"add new feature for webshell"
"Fixed the restore function error"
"merged branch"
```

**重写为**：
```bash
fix(webshell): resolve null pointer in create method
refactor(service): simplify validation logic
feat(webshell): add batch delete support
fix(repository): correct GORM model in Recover method
chore(merge): merge feature/webshell-batch into develop
```

---

## 考核标准

### 理论考核（选择题）

1. 文档应该存放在哪个目录？
   - A. 项目根目录
   - B. docs/ 目录 ✅
   - C. src/ 目录
   - D. 任意位置

2. 以下哪个文件命名是正确的？
   - A. ProjectService.go
   - B. project-service.go
   - C. project_service.go ✅
   - D. PROJECT_SERVICE.go

3. 正确的 Commit 格式是？
   - A. "Fixed bug"
   - B. "fix: bug fix"
   - C. "fix(webshell): resolve null pointer" ✅
   - D. "bug fix"

### 实践考核

**任务**：完成一个完整的功能开发

**评分标准**：

| 项目 | 分值 | 评分标准 |
|------|------|---------|
| 文件命名 | 10 分 | 完全规范 10 分，部分规范 5 分，不规范 0 分 |
| 目录结构 | 10 分 | 完全正确 10 分，基本正确 5 分，混乱 0 分 |
| 代码规范 | 20 分 | 完全符合 20 分，少量问题 10 分，不规范 0 分 |
| 测试覆盖 | 20 分 | ≥80% 20 分，≥60% 15 分，≥40% 10 分，<40% 0 分 |
| 提交信息 | 20 分 | 完全规范 20 分，部分规范 10 分，不规范 0 分 |
| 文档质量 | 20 分 | 完整规范 20 分，基本完整 10 分，缺失 0 分 |

**通过标准**：总分 ≥ 70 分

---

## 常见问题

### Q1: 为什么文档不能放在根目录？

**A**: 
- 保持根目录整洁，只显示最重要的文件
- 便于文档分类和管理
- 符合 Go 项目标准结构
- 方便自动化处理和工具集成

### Q2: Go 文件为什么用下划线而不是连字符？

**A**:
- Go 社区约定俗成的规范
- 与 Go 标准库保持一致
- 工具链支持更好
- 避免与包名混淆

### Q3: 为什么 Commit Message 要用现在时？

**A**:
- Git 提交信息描述的是"这个提交做了什么"
- 使用祈使句更符合 Git 的设计哲学
- 便于生成 CHANGELOG
- 与 Git 自动生成的信息风格一致

### Q4: 测试覆盖率要求是否必须达到？

**A**:
- 是硬性要求，但可以根据实际情况调整
- 核心业务逻辑（Service、Entity）必须达到
- 简单代码（Handler、配置）可以适当降低
- 关键是测试质量，不是单纯追求数字

### Q5: 如何处理遗留的不规范代码？

**A**:
1. 新建功能严格按照规范
2. 修改旧代码时顺便重构（童子军规则）
3. 制定迁移计划，逐步改进
4. 不要一次性全部重构，风险太大

### Q6: 规范是否一成不变？

**A**:
- 规范可以改进，但需要团队讨论
- 每季度审查一次规范
- 改进需要记录在文档中
- 确保团队成员都了解变更

---

## 培训资源

### 必读文档
- [项目优化规范](file://d:\Go\FG-ABYSS\docs\development\project-optimization-specification.md) - 完整规范
- [快速检查清单](file://d:\Go\FG-ABYSS\docs\development\QUICK_CHECKLIST.md) - 日常参考
- 本培训指南

### 参考链接
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Vue.js Style Guide](https://vuejs.org/style-guide/)
- [TypeScript Deep Dive](https://basarat.gitbook.io/typescript/)

### 工具推荐
- **Commitlint** - Commit message 检查工具
- **gofmt** - Go 代码格式化
- **ESLint** - TypeScript/JavaScript 检查
- **Vitest** - 前端测试框架

---

## 持续改进

### 反馈渠道
- 在团队会议中提出规范问题
- 在 PR 中讨论规范改进
- 更新规范文档时通知团队

### 定期审查
- 每月：检查规范执行情况
- 每季：审查规范合理性
- 每年：全面更新规范

### 激励机制
- 表彰规范执行优秀的成员
- 分享最佳实践案例
- 建立规范大使制度

---

**培训版本**: v1.0  
**更新日期**: 2026-03-15  
**负责人**: FG-ABYSS Team Lead
