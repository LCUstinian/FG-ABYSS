#!/bin/bash

# FG-ABYSS 项目规范检查脚本
# 用于自动化检查项目规范的遵守情况

set -e

echo "========================================="
echo "FG-ABYSS 项目规范检查"
echo "========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 计数器
ERRORS=0
WARNINGS=0

# 检查根目录的.md 文件（除了 README.md 和 LICENSE）
echo "📋 检查 1: 文档位置规范"
echo "-----------------------------------------"

ROOT_MD_FILES=$(ls -1 *.md 2>/dev/null | grep -v -E "^(README|LICENSE)" || true)

if [ -n "$ROOT_MD_FILES" ]; then
    echo -e "${RED}❌ 错误：根目录发现不应存在的.md 文件:${NC}"
    echo "$ROOT_MD_FILES"
    echo ""
    echo "建议移动到 docs/ 目录下："
    for file in $ROOT_MD_FILES; do
        echo "  mv $file docs/reports/"
    done
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ 通过：根目录没有违规的.md 文件${NC}"
fi

echo ""

# 检查 docs 目录结构
echo "📋 检查 2: 文档目录结构"
echo "-----------------------------------------"

REQUIRED_DIRS=("docs/architecture" "docs/development" "docs/testing" "docs/api" "docs/deployment" "docs/reports")

for dir in "${REQUIRED_DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        echo -e "${YELLOW}⚠️  警告：缺少目录 $dir${NC}"
        WARNINGS=$((WARNINGS + 1))
    fi
done

if [ -d "docs" ]; then
    echo -e "${GREEN}✅ 通过：docs 目录结构完整${NC}"
else
    echo -e "${RED}❌ 错误：docs 目录不存在${NC}"
    ERRORS=$((ERRORS + 1))
fi

echo ""

# 检查 Go 文件命名
echo "📋 检查 3: Go 文件命名规范"
echo "-----------------------------------------"

# 查找大写字母开头的 Go 文件（除了 vendor 和 bindings）
INVALID_GO=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./frontend/bindings/*" | xargs -I {} basename {} | grep -E "^[A-Z]" || true)

if [ -n "$INVALID_GO" ]; then
    echo -e "${RED}❌ 错误：发现命名不规范的 Go 文件:${NC}"
    echo "$INVALID_GO"
    echo ""
    echo "建议：Go 文件应该使用小写字母 + 下划线命名"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ 通过：Go 文件命名规范${NC}"
fi

echo ""

# 检查测试文件
echo "📋 检查 4: 测试文件命名"
echo "-----------------------------------------"

# 查找不以 _test.go 结尾的 Go 测试文件
INVALID_TEST=$(find . -name "*_test.go" -not -path "./vendor/*" || true)

if [ -n "$INVALID_TEST" ]; then
    TEST_COUNT=$(echo "$INVALID_TEST" | wc -l | tr -d ' ')
    echo -e "${GREEN}✅ 通过：发现 $TEST_COUNT 个测试文件${NC}"
else
    echo -e "${YELLOW}⚠️  警告：没有发现测试文件${NC}"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""

# 检查 Vue 组件命名
echo "📋 检查 5: Vue 组件命名规范"
echo "-----------------------------------------"

# 查找小写字母开头的 Vue 文件
INVALID_VUE=$(find ./frontend/src -name "*.vue" 2>/dev/null | xargs -I {} basename {} | grep -E "^[a-z]" || true)

if [ -n "$INVALID_VUE" ]; then
    echo -e "${RED}❌ 错误：发现命名不规范的 Vue 组件:${NC}"
    echo "$INVALID_VUE"
    echo ""
    echo "建议：Vue 组件应该使用 PascalCase 命名"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ 通过：Vue 组件命名规范${NC}"
fi

echo ""

# 检查是否有调试代码
echo "📋 检查 6: 调试代码检查"
echo "-----------------------------------------"

# 查找 console.log 和 fmt.Println 用于调试
CONSOLE_LOG=$(grep -r "console\.log" ./frontend/src --include="*.ts" --include="*.vue" 2>/dev/null | grep -v "node_modules" || true)
# FMT_PRINT=$(grep -r "fmt\.Print" . --include="*.go" --exclude-dir=vendor 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$CONSOLE_LOG" ]; then
    COUNT=$(echo "$CONSOLE_LOG" | wc -l | tr -d ' ')
    echo -e "${YELLOW}⚠️  警告：发现 $COUNT 处 console.log${NC}"
    echo "请检查是否为调试代码："
    echo "$CONSOLE_LOG" | head -5
    WARNINGS=$((WARNINGS + 1))
else
    echo -e "${GREEN}✅ 通过：没有发现调试代码${NC}"
fi

echo ""

# 检查 Git 提交历史（如果在 Git 仓库中）
if [ -d ".git" ]; then
    echo "📋 检查 7: Git 提交信息规范（最近 10 条）"
    echo "-----------------------------------------"
    
    # 检查最近 10 条提交
    INVALID_COMMITS=$(git log --oneline -10 2>/dev/null | grep -v -E "^(feat|fix|docs|style|refactor|test|chore|perf|ci)\(" || true)
    
    if [ -n "$INVALID_COMMITS" ]; then
        echo -e "${YELLOW}⚠️  警告：发现不规范的提交信息:${NC}"
        echo "$INVALID_COMMITS"
        echo ""
        echo "建议格式：<type>(<scope>): <subject>"
        echo "例如：feat(webshell): add batch delete support"
        WARNINGS=$((WARNINGS + 1))
    else
        echo -e "${GREEN}✅ 通过：最近提交信息规范${NC}"
    fi
    
    echo ""
fi

# 检查测试覆盖率（如果存在覆盖率文件）
echo "📋 检查 8: 测试覆盖率"
echo "-----------------------------------------"

if [ -f "coverage.out" ]; then
    COVERAGE=$(go tool cover -func=coverage.out 2>/dev/null | grep "total:" | awk '{print $3}')
    if [ -n "$COVERAGE" ]; then
        echo "总体覆盖率：$COVERAGE"
        
        # 提取数字部分
        COVERAGE_NUM=$(echo "$COVERAGE" | sed 's/%//')
        
        if (( $(echo "$COVERAGE_NUM >= 40" | bc -l 2>/dev/null || echo "0") )); then
            echo -e "${GREEN}✅ 通过：测试覆盖率达标${NC}"
        else
            echo -e "${YELLOW}⚠️  警告：测试覆盖率低于 40%${NC}"
            WARNINGS=$((WARNINGS + 1))
        fi
    else
        echo -e "${YELLOW}⚠️  信息：无法解析覆盖率文件${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  信息：没有找到覆盖率文件 (coverage.out)${NC}"
    echo "运行测试生成覆盖率：task test"
fi

echo ""

# 总结
echo "========================================="
echo "检查总结"
echo "========================================="
echo ""

if [ $ERRORS -gt 0 ]; then
    echo -e "${RED}❌ 发现 $ERRORS 个错误${NC}"
fi

if [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}⚠️  发现 $WARNINGS 个警告${NC}"
fi

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ 完美！所有检查通过${NC}"
fi

echo ""
echo "详细规范文档：docs/development/project-optimization-specification.md"
echo "快速参考：docs/development/QUICK_CHECKLIST.md"
echo ""

# 退出码
if [ $ERRORS -gt 0 ]; then
    exit 1
else
    exit 0
fi
