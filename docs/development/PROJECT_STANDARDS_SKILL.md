# FG-ABYSS 项目规范技能包

> 完整的项目管理和开发规范体系，包含文档、技能和检查工具

## 📦 技能包组成

### 1. 核心规范文档

**位置**: `docs/development/project-optimization-specification.md`

**内容**:
- 文档管理规范
- 文件命名规范
- Git 提交规范
- 项目结构规范
- 代码开发规范
- 测试规范
- 版本管理规范

### 2. 快速检查清单

**位置**: `docs/development/QUICK_CHECKLIST.md`

**用途**: 日常开发快速参考，打印或保存为书签

**内容**:
- 创建新文件检查项
- 编写代码检查项
- 提交代码前检查项
- 常见错误及修正

### 3. 培训指南

**位置**: `docs/development/TRAINING_GUIDE.md`

**用途**: 新成员培训和团队学习

**内容**:
- 学习路径（分阶段）
- 实践练习
- 考核标准
- 常见问题解答

### 4. 执行技能 (Skill)

**位置**: `.trae/skills/project-standard-checker/SKILL.md`

**用途**: AI 助手自动检查规范遵守情况

**触发场景**:
- 创建新文件时
- 组织项目结构时
- 审查代码合规性时
- 准备提交代码时

### 5. 自动化检查脚本

**位置**: `scripts/check-standards.ps1`

**用途**: 自动化检查项目规范

**运行方式**:
```powershell
# Windows
powershell -ExecutionPolicy Bypass -File scripts\check-standards.ps1

# Linux/macOS
bash scripts/check-standards.sh
```

**检查项目**:
1. 文档位置规范
2. 文档目录结构
3. Go 文件命名
4. 测试文件
5. Vue 组件命名
6. Git 提交信息

---

## 🚀 快速开始

### 新成员入职

1. **阅读规范文档** (1-2 小时)
   ```bash
   # 打开完整规范
   code docs/development/project-optimization-specification.md
   
   # 打开快速清单
   code docs/development/QUICK_CHECKLIST.md
   ```

2. **完成培训练习** (2-3 小时)
   - 按照 `TRAINING_GUIDE.md` 进行实践
   - 完成所有练习项目

3. **通过考核**
   - 理论考核（选择题）
   - 实践考核（功能开发）

### 日常开发

**开发前**:
```bash
# 运行检查脚本
powershell -ExecutionPolicy Bypass -File scripts/check-standards.ps1
```

**开发中**:
- 遵循规范文档
- 参考快速清单
- 使用 AI 技能助手

**提交前**:
- [ ] 运行测试
- [ ] 检查覆盖率
- [ ] 验证命名规范
- [ ] 编写规范 commit message

### 团队管理

**每周检查**:
```bash
# 运行完整检查
./scripts/check-standards.ps1

# 查看检查结果
# - Errors: 必须修复
# - Warnings: 建议修复
```

**每月审查**:
- 审查规范执行情况
- 讨论规范改进建议
- 更新规范文档

---

## 📊 检查结果示例

```
=========================================
FG-ABYSS Project Standards Check
=========================================

1. Checking documentation location...
  [ERROR] Found .md files in root directory:
    - BUGFIX_REPORT.md
    - TEST_REPORT.md

2. Checking docs directory structure...
  [OK] docs directory exists

3. Checking Go file naming...
  [OK] Go file naming is correct

4. Checking test files...
  [OK] Found 4 test files

5. Checking Vue component naming...
  [OK] Vue component naming is correct

6. Checking Git commit messages...
  [WARNING] Found 8 non-standard commits

=========================================
Summary
=========================================
Errors: 1
Warnings: 1
```

---

## 🔧 立即修复

根据检查结果，需要修复：

### 1. 移动文档到正确位置

```bash
# 移动报告文档
mv BUGFIX_REPORT.md docs/testing/bugfixes/
mv TEST_REPORT.md docs/testing/

# 重命名为规范格式
cd docs/testing/bugfixes/
ren BUGFIX_REPORT.md webshell-restore-fix.md

cd ../
ren TEST_REPORT.md unit-test-report.md
```

### 2. 规范 Git 提交

下次提交时使用正确格式：
```bash
git commit -m "type(scope): description"

# 示例
git commit -m "fix(webshell): restore deleted webshell functionality"
git commit -m "test(service): add comprehensive unit tests"
```

---

## 📚 文档索引

| 文档 | 用途 | 位置 |
|------|------|------|
| 项目优化规范 | 完整规范文档 | `docs/development/project-optimization-specification.md` |
| 快速检查清单 | 日常参考 | `docs/development/QUICK_CHECKLIST.md` |
| 培训指南 | 团队培训 | `docs/development/TRAINING_GUIDE.md` |
| 技能说明 | AI 助手配置 | `.trae/skills/project-standard-checker/SKILL.md` |
| 检查脚本 | 自动化工具 | `scripts/check-standards.ps1` |

---

## 🎯 规范要点

### 文档管理
- ✅ 所有文档在 `docs/` 目录
- ✅ 按类型分类存放
- ❌ 禁止在根目录创建文档

### 文件命名
- ✅ Go: 小写 + 下划线 (`project_service.go`)
- ✅ Vue: PascalCase (`ProjectsContent.vue`)
- ✅ 文档：小写 + 连字符 (`unit-test-report.md`)

### Git 提交
- ✅ 格式：`<type>(<scope>): <subject>`
- ✅ 类型：feat, fix, docs, style, refactor, test, chore
- ✅ 时态：祈使句，现在时

### 项目结构
- ✅ `internal/app/handlers/` - 请求处理
- ✅ `internal/app/services/` - 业务逻辑
- ✅ `internal/domain/entity/` - 实体定义
- ✅ `internal/infrastructure/repositories/` - 数据访问

---

## 💡 最佳实践

### 1. 文档及时更新
- 新功能 → 更新 API 文档
- Bug 修复 → 编写修复报告
- 架构变更 → 更新架构文档

### 2. 小步提交
- 一个功能一个提交
- 提交信息清晰完整
- 遵循 Conventional Commits

### 3. 测试先行
- 先写测试再写实现
- 保证覆盖率达标
- 测试命名清晰

### 4. 代码审查
- 使用检查清单
- 自动化检查工具
- 团队互相监督

---

## 📞 获取帮助

- 📖 查看完整规范文档
- ✅ 运行检查脚本
- 🤖 使用 AI 技能助手
- 💬 团队讨论

---

**版本**: v1.0  
**创建日期**: 2026-03-15  
**维护者**: FG-ABYSS Team
