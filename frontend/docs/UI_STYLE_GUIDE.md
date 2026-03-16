# FG-ABYSS UI 样式规范文档

## 📋 目录

1. [设计原则](#设计原则)
2. [CSS 变量系统](#css-变量系统)
3. [样式层级架构](#样式层级架构)
4. [深色/浅色模式规范](#深色浅色模式规范)
5. [最佳实践](#最佳实践)
6. [代码审查清单](#代码审查清单)

---

## 🎯 设计原则

### 1. 单一数据源
- **原则**：所有样式值必须来自 CSS 变量
- **禁止**：硬编码颜色值、间距值、圆角值
- **例外**：无（零容忍）

### 2. 主题一致性
- **原则**：深色/浅色模式使用相同的变量名，仅值不同
- **禁止**：为不同模式定义不同的变量
- **示例**：
  ```css
  /* ✅ 正确 */
  :root { --bg-primary: #ffffff; }
  .dark { --bg-primary: #0f172a; }
  
  /* ❌ 错误 */
  :root { --bg-light: #ffffff; }
  .dark { --bg-dark: #0f172a; }
  ```

### 3. 层级清晰
- **原则**：样式必须遵循四层架构
- **禁止**：跨层级覆盖（除非必要）
- **优先级**：设计令牌 < 全局样式 < 主题覆盖 < 组件样式

---

## 🎨 CSS 变量系统

### 变量命名规范

#### 1. 背景色变量
```css
--bg-primary:      /* 主背景色 */
--bg-secondary:    /* 次背景色 */
--bg-tertiary:     /* 第三级背景色 */
--bg-hover:        /* 悬停背景色 */
```

#### 2. 文字色变量
```css
--text-primary:    /* 主文字色 */
--text-secondary:  /* 次文字色 */
--text-tertiary:   /* 第三级文字色 */
```

#### 3. 边框色变量
```css
--border-subtle:   /* 细边框 */
--border-strong:   /* 强调边框 */
```

#### 4. 强调色变量
```css
--active-color:         /* 主强调色 */
--active-color-rgb:     /* RGB 格式（用于透明度） */
--active-color-suppl:   /* 辅助强调色（10% 透明） */
--active-color-hover:   /* 强调色悬停（提亮 10%） */
```

### 变量使用规范

#### ✅ 正确使用
```css
/* 使用 CSS 变量 */
.component {
  background-color: var(--bg-primary);
  color: var(--text-primary);
  border: 1px solid var(--border-subtle);
}
```

#### ❌ 错误使用
```css
/* 硬编码颜色值 */
.component {
  background-color: #ffffff;  /* ❌ */
  color: #1f2937;             /* ❌ */
  border: 1px solid #e5e7eb;  /* ❌ */
}
```

---

## 🏗️ 样式层级架构

### 四层架构

```
┌─────────────────────────────────────┐
│ Layer 4: 组件特定样式              │
│ - 组件 scoped 样式                  │
│ - 最高优先级（不使用 !important）   │
└─────────────────────────────────────┘
           ↑
┌─────────────────────────────────────┐
│ Layer 3: Naive UI 主题覆盖          │
│ - App.vue themeOverrides            │
│ - 统一的组件主题配置                │
└─────────────────────────────────────┘
           ↑
┌─────────────────────────────────────┐
│ Layer 2: 全局功能样式              │
│ - global.css                        │
│ - 通用组件样式（dialog, message 等） │
└─────────────────────────────────────┘
           ↑
┌─────────────────────────────────────┐
│ Layer 1: 设计令牌（CSS 变量）       │
│ - 颜色、间距、圆角、阴影            │
│ - 深色/浅色模式变量定义             │
└─────────────────────────────────────┘
```

### 各层级职责

#### Layer 1: 设计令牌
- **位置**：`global.css:root` 和 `global.css.dark`
- **职责**：定义所有 CSS 变量
- **禁止**：直接应用于组件

#### Layer 2: 全局功能样式
- **位置**：`global.css`
- **职责**：通用组件样式（Naive UI 覆盖）
- **允许**：使用 `!important` 覆盖 Naive UI

#### Layer 3: Naive UI 主题覆盖
- **位置**：`App.vue themeOverrides`
- **职责**：统一配置 Naive UI 组件主题
- **推荐**：优先使用此方式

#### Layer 4: 组件样式
- **位置**：各组件 `.vue` 文件
- **职责**：组件特定样式
- **允许**：使用 `:deep()` 修改子组件

---

## 🌓 深色/浅色模式规范

### 变量定义一致性

#### ✅ 正确示例
```css
/* 浅色模式 */
:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #ffffff;
  --text-primary: #1f2937;
  --text-secondary: #4b5563;
  --border-subtle: #e5e7eb;
  --border-strong: #d1d5db;
}

/* 深色模式 - 变量名称和数量完全一致 */
.dark {
  --bg-primary: #0f172a;
  --bg-secondary: #1e293b;
  --bg-tertiary: #334155;
  --text-primary: #f1f5f9;
  --text-secondary: #94a3b8;
  --border-subtle: #334155;
  --border-strong: #475569;
}
```

#### ❌ 错误示例
```css
/* 浅色模式 */
:root {
  --bg-primary: #ffffff;
  --card-bg: #f9fafb;  /* ❌ 变量名不一致 */
}

/* 深色模式 */
.dark {
  --bg-primary: #0f172a;
  --bg-card: #1e293b;  /* ❌ 变量名改变 */
}
```

### 组件样式自动适配

#### ✅ 正确：使用 CSS 变量自动适配
```css
.component {
  background-color: var(--bg-primary);  /* 自动适配主题 */
  color: var(--text-primary);
  border: 1px solid var(--border-subtle);
}
/* 无需写 .dark .component */
```

#### ⚠️ 仅在特殊情况下使用 .dark
```css
/* 深色模式需要特殊处理时 */
.component {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.dark .component {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);  /* 深色模式阴影更深 */
}
```

---

## 📖 最佳实践

### 1. Naive UI 组件样式处理

#### 方案 A：使用 themeOverrides（推荐）
```typescript
// App.vue
const themeOverrides = computed(() => ({
  common: {
    primaryColor: getVar('--active-color', '#3b82f6'),
    primaryColorSuppl: getVar('--active-color-suppl', 'rgba(59, 130, 246, 0.1)'),
  },
  Button: {
    colorPrimary: 'var(--active-color)',
    colorPrimaryHover: 'var(--active-color-hover)',
  },
  Select: {
    colorActive: 'var(--active-color-suppl)',
    borderColorActive: 'var(--active-color)',
  }
}))
```

#### 方案 B：全局 CSS 覆盖
```css
/* global.css - 统一覆盖 */
.n-button--primary {
  --n-color: var(--active-color);
  --n-color-hover: var(--active-color-hover);
}

.n-select--active {
  --n-border: var(--active-color);
  --n-color-active: var(--active-color-suppl);
}
```

### 2. !important 使用规范

#### ✅ 允许使用 !important 的场景
1. 覆盖 Naive UI 内联样式
2. 全局样式覆盖组件 scoped 样式
3. 主题覆盖默认样式

#### ❌ 禁止使用 !important 的场景
1. 组件内部样式覆盖
2. 可以通过选择器优先级解决的场景
3. 为了方便而滥用

### 3. 深色模式特殊效果

#### 阴影处理
```css
/* 浅色模式 */
.component {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 深色模式 - 加深阴影 */
.dark .component {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}
```

#### 毛玻璃效果
```css
/* 浅色模式 */
.component {
  backdrop-filter: blur(8px);
  background: rgba(255, 255, 255, 0.9);
}

/* 深色模式 - 增强毛玻璃 */
.dark .component {
  backdrop-filter: blur(16px);
  background: rgba(30, 41, 59, 0.95);
}
```

---

## ✅ 代码审查清单

### CSS 变量使用
- [ ] 所有颜色值都使用 CSS 变量
- [ ] 无硬编码的 HEX/RGB 颜色值
- [ ] 变量命名符合规范
- [ ] 深色/浅色模式变量名一致

### 主题适配
- [ ] 组件自动适配深色/浅色模式
- [ ] 无模式特定的硬编码颜色
- [ ] 特殊效果（阴影、毛玻璃）已适配
- [ ] 主题切换流畅无闪烁

### 样式层级
- [ ] 遵循四层架构
- [ ] 无跨层级覆盖
- [ ] !important 使用合理
- [ ] 选择器优先级正确

### Naive UI 组件
- [ ] 优先使用 themeOverrides
- [ ] 全局覆盖使用 CSS 变量
- [ ] 组件特定样式在组件内定义
- [ ] 无重复的样式定义

### 性能优化
- [ ] 使用 transform 而非 margin/padding 实现位移
- [ ] 避免在 hover 时改变布局
- [ ] 使用 will-change 优化动画
- [ ] 过渡时间符合规范（150ms/200ms/300ms）

---

## 📊 已优化项目清单

### 已完成的优化
- ✅ ProjectsContent.vue - 移除所有硬编码颜色
- ✅ RecoverProjectModal.vue - 使用 CSS 变量
- ✅ Tooltip.vue - 统一文字颜色变量
- ✅ global.css - 建立完整的变量系统

### 待优化的组件
- [ ] CreateProjectModal.vue
- [ ] CreateWebShellModal.vue
- [ ] EditWebShellModal.vue
- [ ] 其他模态框组件

---

## 🔄 持续改进

### 定期审查
- 每月审查新增样式代码
- 检查是否有硬编码颜色
- 验证主题适配完整性
- 更新文档和最佳实践

### 工具辅助
- 使用 Stylelint 检查硬编码颜色
- 配置 CSS 变量使用提示
- 自动化主题对比度测试

---

**最后更新**: 2026-03-17  
**维护者**: FG-ABYSS Team
