# FG-ABYSS 项目规范检查脚本

用于自动化检查项目规范的遵守情况。

## 使用方法

### Linux/macOS

```bash
# 添加执行权限
chmod +x scripts/check-standards.sh

# 运行检查
./scripts/check-standards.sh
```

### Windows (PowerShell)

```powershell
# 使用 Git Bash 运行
bash scripts/check-standards.sh

# 或者使用 WSL
wsl bash scripts/check-standards.sh
```

## 检查项目

脚本会检查以下项目：

1. **文档位置** - 根目录不应有.md 文件（除 README、LICENSE）
2. **文档目录结构** - docs/ 下的分类目录
3. **Go 文件命名** - 应使用小写 + 下划线
4. **测试文件命名** - 应以 _test.go 结尾
5. **Vue 组件命名** - 应使用 PascalCase
6. **调试代码** - 检查 console.log 等
7. **Git 提交信息** - 检查 Conventional Commits 规范
8. **测试覆盖率** - 检查是否达标

## 输出示例

```
=========================================
FG-ABYSS 项目规范检查
=========================================

📋 检查 1: 文档位置规范
-----------------------------------------
✅ 通过：根目录没有违规的.md 文件

📋 检查 2: 文档目录结构
-----------------------------------------
✅ 通过：docs 目录结构完整

📋 检查 3: Go 文件命名规范
-----------------------------------------
✅ 通过：Go 文件命名规范

📋 检查 4: 测试文件命名
-----------------------------------------
✅ 通过：发现 4 个测试文件

...

=========================================
检查总结
=========================================

✅ 完美！所有检查通过
```

## 集成到 CI/CD

可以将此脚本添加到 CI/CD 流程中：

```yaml
# .github/workflows/ci.yml
jobs:
  check-standards:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Check Project Standards
        run: ./scripts/check-standards.sh
```

## 修复建议

如果检查失败，脚本会提供修复建议：

```bash
# 示例输出
❌ 错误：根目录发现不应存在的.md 文件:
PROJECT_REPORT.md

建议移动到 docs/ 目录下：
  mv PROJECT_REPORT.md docs/reports/
```

## 自定义检查

可以根据项目需求修改脚本，添加更多检查项。

编辑：`scripts/check-standards.sh`

## 相关文档

- [完整规范](docs/development/project-optimization-specification.md)
- [快速检查清单](docs/development/QUICK_CHECKLIST.md)
- [培训指南](docs/development/TRAINING_GUIDE.md)
