# FG-ABYSS 现代化用户界面设计规范

## 设计愿景

打造一套**专业、高效、安全**的渗透测试工具界面，在保留所有现有功能的基础上，提升用户体验和视觉美感。

## 核心设计原则

### 1. 专业性（Professional）
- 符合安全工具的专业形象
- 清晰的信息层次和视觉重点
- 严谨的交互逻辑和反馈机制

### 2. 高效性（Efficient）
- 优化工作流程，减少操作步骤
- 快捷键支持，提升专家用户效率
- 智能默认值和上下文感知

### 3. 一致性（Consistent）
- 统一的视觉语言和设计规范
- 跨平台一致的交互体验
- 符合各操作系统的人机界面指南

### 4. 可用性（Usable）
- 直观的信息架构
- 清晰的导航和定位
- 友好的错误提示和恢复机制

---

## 视觉设计系统

### 色彩系统 2.0

#### 主色调（Primary Palette）

保留现有的蓝色系作为主色，但进行优化：

```css
:root {
  /* 主色系 - 科技蓝 */
  --primary-50: #eff6ff;
  --primary-100: #dbeafe;
  --primary-200: #bfdbfe;
  --primary-300: #93c5fd;
  --primary-400: #60a5fa;
  --primary-500: #3b82f6;  /* 当前 --active-color */
  --primary-600: #2563eb;
  --primary-700: #1d4ed8;
  --primary-800: #1e40af;
  --primary-900: #1e3a8a;
  
  /* 语义色系 */
  --success-500: #10b981;  /* 成功/安全 */
  --warning-500: #f59e0b;  /* 警告/注意 */
  --error-500: #ef4444;    /* 错误/危险 */
  --info-500: #3b82f6;     /* 信息/提示 */
  
  /* 中性色系 */
  --gray-50: #f9fafb;
  --gray-100: #f3f4f6;
  --gray-200: #e5e7eb;
  --gray-300: #d1d5db;
  --gray-400: #9ca3af;
  --gray-500: #6b7280;
  --gray-600: #4b5563;
  --gray-700: #374151;
  --gray-800: #1f2937;
  --gray-900: #111827;
}
```

#### 暗色主题优化

```css
.dark {
  --bg-primary: #0f172a;      /* 主背景 - 深空蓝 */
  --bg-secondary: #1e293b;    /* 次级背景 - 侧边栏 */
  --bg-tertiary: #334155;     /* 第三背景 - 卡片 */
  
  --text-primary: #f1f5f9;    /* 主文字 */
  --text-secondary: #94a3b8;  /* 次级文字 */
  --text-tertiary: #64748b;   /* 提示文字 */
  
  --border-subtle: #334155;   /* 细边框 */
  --border-strong: #475569;   /* 强边框 */
  
  /* 强调色保持与亮色主题一致 */
  --primary-500: #3b82f6;
}
```

#### 功能色板

为不同功能模块定义专属色板：

```css
/* 项目管理 - 蓝色系 */
--project-color: #3b82f6;

/* WebShell - 绿色系 */
--webshell-color: #10b981;

/* 载荷生成 - 紫色系 */
--payload-color: #8b5cf6;

/* 插件管理 - 橙色系 */
--plugin-color: #f59e0b;

/* 系统设置 - 青色系 */
--settings-color: #06b6d4;
```

### 字体系统 2.0

#### 字体栈优化

```css
:root {
  /* 英文优先使用 Inter，中文优先使用系统字体 */
  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 
               'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 
               'Helvetica Neue', Helvetica, Arial, sans-serif;
  
  /* 代码字体 - 适合渗透测试场景 */
  --font-mono: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 
               'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', 
               Consolas, 'Courier New', monospace;
  
  /* 字号系统 - 基于 1.25 比例 */
  --text-xs: 0.75rem;    /* 12px - 辅助文字 */
  --text-sm: 0.875rem;   /* 14px - 正文 */
  --text-base: 1rem;     /* 16px - 标准 */
  --text-lg: 1.25rem;    /* 20px - 标题 */
  --text-xl: 1.5rem;     /* 24px - 大标题 */
  --text-2xl: 1.875rem;  /* 30px - 超大标题 */
}
```

#### 字重系统

```css
--font-normal: 400;
--font-medium: 500;
--font-semibold: 600;
--font-bold: 700;
```

### 间距系统

采用 **8px 基准网格系统**：

```css
:root {
  --space-1: 0.25rem;  /* 4px */
  --space-2: 0.5rem;   /* 8px */
  --space-3: 0.75rem;  /* 12px */
  --space-4: 1rem;     /* 16px */
  --space-5: 1.25rem;  /* 20px */
  --space-6: 1.5rem;   /* 24px */
  --space-8: 2rem;     /* 32px */
  --space-10: 2.5rem;  /* 40px */
  --space-12: 3rem;    /* 48px */
  --space-16: 4rem;    /* 64px */
}
```

### 圆角系统

```css
:root {
  --radius-sm: 0.25rem;  /* 4px - 小按钮、标签 */
  --radius-md: 0.375rem; /* 6px - 输入框、卡片 */
  --radius-lg: 0.5rem;   /* 8px - 大卡片 */
  --radius-xl: 0.75rem;  /* 12px - 模态框 */
  --radius-2xl: 1rem;    /* 16px - 特殊组件 */
  --radius-full: 9999px; /* 圆形 */
}
```

### 阴影系统 2.0

增加层次感，优化暗色主题阴影：

```css
:root {
  /* 亮色主题阴影 */
  --shadow-xs: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-sm: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  --shadow-2xl: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  
  /* 暗色主题阴影 - 使用更深的阴影 */
  --shadow-dark-sm: 0 1px 3px 0 rgba(0, 0, 0, 0.3), 0 1px 2px 0 rgba(0, 0, 0, 0.2);
  --shadow-dark-md: 0 4px 6px -1px rgba(0, 0, 0, 0.4), 0 2px 4px -1px rgba(0, 0, 0, 0.3);
  --shadow-dark-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.5), 0 4px 6px -2px rgba(0, 0, 0, 0.4);
  --shadow-dark-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.6), 0 10px 10px -5px rgba(0, 0, 0, 0.5);
  
  /* 强调色阴影 - 用于激活状态 */
  --shadow-primary: 0 0 0 3px rgba(59, 130, 246, 0.2);
  --shadow-success: 0 0 0 3px rgba(16, 185, 129, 0.2);
  --shadow-warning: 0 0 0 3px rgba(245, 158, 11, 0.2);
  --shadow-error: 0 0 0 3px rgba(239, 68, 68, 0.2);
}
```

---

## 布局系统升级

### 响应式断点 2.0

采用 **移动优先** 的断点策略：

```css
/* 小屏手机 */
--breakpoint-sm: 640px;

/* 大屏手机 */
--breakpoint-md: 768px;

/* 平板 */
--breakpoint-lg: 1024px;

/* 小屏桌面 */
--breakpoint-xl: 1280px;

/* 大屏桌面 */
--breakpoint-2xl: 1536px;
```

### 桌面应用布局优化

保持经典的三栏布局，但进行优化：

```
┌──────────────────────────────────────────────┐
│            TitleBar (48px)                   │
│  [Logo] FG-ABYSS              [Controls]     │
├─────────┬────────────────────────────────────┤
│         │  Toolbar (56px)                    │
│Sidebar  │  [Breadcrumbs] [Actions] [Search]  │
│(240px)  ├────────────────────────────────────┤
│         │                                    │
│[导航]   │         Main Content               │
│         │         (flex: 1)                  │
│         │                                    │
│         ├────────────────────────────────────┤
│         │  Status Bar (32px)                 │
│         │  [Info] [Progress] [Tips]          │
└─────────┴────────────────────────────────────┘
```

### 侧边栏升级

**宽度调整**：
- 收起状态：`64px`（仅图标）
- 展开状态：`240px`（图标 + 文字）
- 可拖拽调整宽度

**导航结构**：
```
FG-ABYSS Logo
─────────────
📊 首页
📁 项目管理
  ├─ 项目 A
  ├─ 项目 B
  └─ 项目 C
🎯 载荷生成
🔌 插件管理
⚙️ 设置
─────────────
[收起/展开]
```

### 内容区域布局模式

定义 5 种标准布局模式：

#### 模式 1：单栏布局（首页、设置）
```css
.content-single {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--space-6);
}
```

#### 模式 2：双栏布局（项目管理）
```css
.content-two-column {
  display: flex;
  gap: var(--space-4);
  
  .sidebar {
    width: 280px;
    flex-shrink: 0;
  }
  
  .main {
    flex: 1;
    min-width: 0;
  }
}
```

#### 模式 3：表单居中布局（载荷生成）
```css
.content-form {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 140px);
  
  .form-container {
    width: 100%;
    max-width: 600px;
    padding: var(--space-8);
  }
}
```

#### 模式 4：卡片网格布局（插件管理）
```css
.content-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--space-6);
  padding: var(--space-6);
}
```

#### 模式 5：全宽表格布局（数据列表）
```css
.content-full-width {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: var(--space-4);
  
  .table-container {
    flex: 1;
    overflow: auto;
  }
}
```

---

## 组件设计规范

### 1. 导航组件

#### 侧边导航项

```css
.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-normal);
  
  /* 默认状态 */
  background: transparent;
  color: var(--text-secondary);
  
  /* 悬停状态 */
  &:hover {
    background: var(--hover-color);
    color: var(--text-primary);
    transform: translateX(4px);
  }
  
  /* 激活状态 */
  &.active {
    background: linear-gradient(135deg, var(--primary-500), var(--primary-600));
    color: white;
    box-shadow: var(--shadow-md), 0 0 0 2px rgba(59, 130, 246, 0.2);
  }
  
  /* 图标 */
  .icon {
    width: 20px;
    height: 20px;
    flex-shrink: 0;
    transition: transform var(--transition-normal);
  }
  
  &:hover .icon {
    transform: scale(1.1);
  }
  
  /* 文字 */
  .label {
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    white-space: nowrap;
  }
}
```

#### 面包屑导航（新增）

```vue
<Breadcrumb>
  <BreadcrumbItem>首页</BreadcrumbItem>
  <BreadcrumbSeparator>/</BreadcrumbSeparator>
  <BreadcrumbItem>项目管理</BreadcrumbItem>
  <BreadcrumbSeparator>/</BreadcrumbSeparator>
  <BreadcrumbItem active>项目 A</BreadcrumbItem>
</Breadcrumb>
```

```css
.breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  
  .item {
    color: var(--text-secondary);
    cursor: pointer;
    transition: color var(--transition-fast);
    
    &:hover {
      color: var(--primary-500);
    }
    
    &.active {
      color: var(--text-primary);
      font-weight: var(--font-medium);
      cursor: default;
    }
  }
  
  .separator {
    color: var(--text-tertiary);
  }
}
```

### 2. 按钮组件

#### 按钮变体

```css
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  cursor: pointer;
  transition: all var(--transition-fast);
  border: none;
  outline: none;
  
  /* 主要按钮 */
  &.btn-primary {
    background: var(--primary-500);
    color: white;
    
    &:hover {
      background: var(--primary-600);
      transform: translateY(-1px);
      box-shadow: var(--shadow-md);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
  
  /* 次要按钮 */
  &.btn-secondary {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-strong);
    
    &:hover {
      background: var(--hover-color);
      border-color: var(--primary-500);
    }
  }
  
  /* 幽灵按钮 */
  &.btn-ghost {
    background: transparent;
    color: var(--text-secondary);
    
    &:hover {
      background: var(--hover-color);
      color: var(--primary-500);
    }
  }
  
  /* 危险按钮 */
  &.btn-danger {
    background: var(--error-500);
    color: white;
    
    &:hover {
      background: #dc2626;
      box-shadow: var(--shadow-md);
    }
  }
  
  /* 尺寸变体 */
  &.btn-sm {
    padding: var(--space-1) var(--space-2);
    font-size: var(--text-xs);
  }
  
  &.btn-lg {
    padding: var(--space-3) var(--space-6);
    font-size: var(--text-base);
  }
  
  /* 图标按钮 */
  &.btn-icon {
    padding: var(--space-2);
    
    .icon {
      width: 20px;
      height: 20px;
    }
  }
  
  /* 禁用状态 */
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none !important;
  }
}
```

### 3. 输入框组件

```css
.input {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: var(--text-sm);
  font-family: var(--font-sans);
  transition: all var(--transition-fast);
  
  /* 占位符 */
  &::placeholder {
    color: var(--text-tertiary);
  }
  
  /* 焦点状态 */
  &:focus {
    outline: none;
    border-color: var(--primary-500);
    box-shadow: var(--shadow-primary);
  }
  
  /* 错误状态 */
  &.error {
    border-color: var(--error-500);
    
    &:focus {
      box-shadow: var(--shadow-error);
    }
  }
  
  /* 禁用状态 */
  &:disabled {
    background: var(--bg-secondary);
    cursor: not-allowed;
    opacity: 0.6;
  }
  
  /* 带图标的输入框 */
  &.input-with-icon {
    padding-left: var(--space-8);
    
    & + .input-icon {
      position: absolute;
      left: var(--space-3);
      top: 50%;
      transform: translateY(-50%);
      color: var(--text-tertiary);
      pointer-events: none;
    }
  }
}
```

### 4. 卡片组件

```css
.card {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  transition: all var(--transition-normal);
  
  /* 悬停效果 */
  &:hover {
    border-color: var(--primary-500);
    box-shadow: var(--shadow-lg);
    transform: translateY(-2px);
  }
  
  /* 卡片头部 */
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--space-4);
    padding-bottom: var(--space-3);
    border-bottom: 1px solid var(--border-subtle);
    
    .card-title {
      font-size: var(--text-lg);
      font-weight: var(--font-semibold);
      color: var(--text-primary);
    }
    
    .card-subtitle {
      font-size: var(--text-sm);
      color: var(--text-secondary);
      margin-top: var(--space-1);
    }
  }
  
  /* 卡片内容 */
  .card-content {
    color: var(--text-primary);
    line-height: 1.6;
  }
  
  /* 卡片底部 */
  .card-footer {
    margin-top: var(--space-4);
    padding-top: var(--space-3);
    border-top: 1px solid var(--border-subtle);
    display: flex;
    gap: var(--space-2);
  }
  
  /* 可选变体：无边框 */
  &.card-borderless {
    border: none;
    box-shadow: var(--shadow-md);
  }
  
  /* 可选变体：可点击 */
  &.card-clickable {
    cursor: pointer;
    
    &:active {
      transform: translateY(0);
    }
  }
}
```

### 5. 表格组件

```css
.table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  
  thead {
    th {
      padding: var(--space-3) var(--space-4);
      background: var(--bg-secondary);
      border-bottom: 2px solid var(--border-strong);
      text-align: left;
      font-size: var(--text-sm);
      font-weight: var(--font-semibold);
      color: var(--text-primary);
      cursor: pointer;
      transition: all var(--transition-fast);
      
      &:hover {
        background: var(--hover-color);
      }
      
      /* 可排序列 */
      &.sortable {
        user-select: none;
        
        .sort-icon {
          margin-left: var(--space-2);
          opacity: 0.5;
        }
        
        &:hover .sort-icon {
          opacity: 1;
        }
      }
      
      /* 可调整列宽 */
      &.resizable {
        position: relative;
        
        .resize-handle {
          position: absolute;
          right: 0;
          top: 0;
          bottom: 0;
          width: 5px;
          cursor: col-resize;
          background: transparent;
          
          &:hover {
            background: var(--primary-500);
          }
        }
      }
    }
  }
  
  tbody {
    tr {
      transition: all var(--transition-fast);
      
      /* 默认行 */
      td {
        padding: var(--space-3) var(--space-4);
        border-bottom: 1px solid var(--border-subtle);
        font-size: var(--text-sm);
        color: var(--text-primary);
      }
      
      /* 悬停状态 */
      &:hover {
        background: var(--hover-color);
      }
      
      /* 选中状态 */
      &.selected {
        background: rgba(59, 130, 246, 0.1);
        border-left: 3px solid var(--primary-500);
        
        td {
          color: var(--primary-500);
          font-weight: var(--font-medium);
        }
      }
      
      /* 禁用状态 */
      &.disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }
    }
  }
  
  /* 紧凑模式 */
  &.table-compact {
    th, td {
      padding: var(--space-2) var(--space-3);
    }
  }
  
  /* 条纹模式 */
  &.table-striped {
    tbody tr:nth-child(even) {
      background: rgba(0, 0, 0, 0.02);
      
      &:hover {
        background: var(--hover-color);
      }
    }
  }
}
```

### 6. 模态框/弹窗

```css
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn var(--transition-fast) ease-out;
}

.modal {
  background: var(--bg-tertiary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2xl);
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
  animation: slideUp var(--transition-normal) ease-out;
  
  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-5);
    border-bottom: 1px solid var(--border-subtle);
    
    .modal-title {
      font-size: var(--text-xl);
      font-weight: var(--font-semibold);
      color: var(--text-primary);
    }
    
    .close-button {
      width: 32px;
      height: 32px;
      border-radius: var(--radius-md);
      border: none;
      background: transparent;
      color: var(--text-secondary);
      cursor: pointer;
      transition: all var(--transition-fast);
      
      &:hover {
        background: var(--hover-color);
        color: var(--text-primary);
      }
    }
  }
  
  .modal-content {
    padding: var(--space-5);
    overflow-y: auto;
    max-height: calc(90vh - 140px);
  }
  
  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: var(--space-3);
    padding: var(--space-5);
    border-top: 1px solid var(--border-subtle);
    background: var(--bg-secondary);
  }
  
  /* 尺寸变体 */
  &.modal-sm { max-width: 400px; }
  &.modal-lg { max-width: 800px; }
  &.modal-xl { max-width: 1000px; }
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

### 7. 工具提示（Tooltip）

```css
.tooltip {
  position: relative;
  display: inline-block;
  
  &:hover .tooltip-content {
    visibility: visible;
    opacity: 1;
  }
}

.tooltip-content {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(-8px);
  padding: var(--space-2) var(--space-3);
  background: var(--bg-tertiary);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  font-size: var(--text-xs);
  color: var(--text-primary);
  white-space: nowrap;
  z-index: 1000;
  visibility: hidden;
  opacity: 0;
  transition: all var(--transition-fast);
  backdrop-filter: blur(8px);
  
  /* 小三角 */
  &::after {
    content: '';
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -4px;
    border-left: 4px solid transparent;
    border-right: 4px solid transparent;
    border-top: 4px solid var(--border-strong);
  }
  
  /* 位置变体 */
  &.tooltip-right {
    left: 100%;
    top: 50%;
    transform: translateY(-50%) translateX(8px);
    
    &::after {
      top: 50%;
      left: 0;
      margin-left: -4px;
      margin-top: -4px;
      border-left: none;
      border-right: 4px solid var(--border-strong);
      border-top: 4px solid transparent;
      border-bottom: 4px solid transparent;
    }
  }
}
```

---

## 交互设计规范

### 1. 动画原则

#### 缓动函数

```css
/* 标准缓动 - 用于大多数过渡 */
--ease-standard: cubic-bezier(0.4, 0.0, 0.2, 1);

/* 强调缓动 - 用于重要元素 */
--ease-emphasis: cubic-bezier(0.4, 0.0, 0.6, 1);

/* 弹跳缓动 - 用于 playful 场景 */
--ease-bounce: cubic-bezier(0.68, -0.55, 0.265, 1.55);

/* 平滑缓动 - 用于持续动画 */
--ease-smooth: cubic-bezier(0.25, 0.1, 0.25, 1);
```

#### 过渡时长

```css
/* 快速 - 用于小元素 */
--duration-fast: 150ms;

/* 标准 - 用于大多数元素 */
--duration-normal: 200ms;

/* 慢速 - 用于大元素或重要动画 */
--duration-slow: 300ms;

/* 极慢 - 用于背景或特殊效果 */
--duration-slower: 500ms;
```

### 2. 微交互设计

#### 按钮点击反馈

```css
.btn {
  position: relative;
  overflow: hidden;
  
  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle, rgba(255,255,255,0.3) 0%, transparent 70%);
    transform: scale(0);
    opacity: 0;
    transition: transform 0.5s, opacity 0.3s;
  }
  
  &:active::after {
    transform: scale(2);
    opacity: 1;
    transition: transform 0s, opacity 0s;
  }
}
```

#### 加载状态

```css
.loading-spinner {
  width: 24px;
  height: 24px;
  border: 3px solid var(--border-subtle);
  border-top-color: var(--primary-500);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 骨架屏 */
.skeleton {
  background: linear-gradient(
    90deg,
    var(--bg-secondary) 0%,
    var(--bg-tertiary) 50%,
    var(--bg-secondary) 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
```

#### 成功/错误反馈

```css
/* 成功动画 */
.success-checkmark {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--success-500);
  display: flex;
  align-items: center;
  justify-content: center;
  animation: successPop 0.5s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

@keyframes successPop {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

/* 错误抖动 */
.shake {
  animation: shake 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97);
}

@keyframes shake {
  10%, 90% { transform: translateX(-1px); }
  20%, 80% { transform: translateX(2px); }
  30%, 50%, 70% { transform: translateX(-4px); }
  40%, 60% { transform: translateX(4px); }
}
```

### 3. 键盘快捷键

定义全局快捷键系统：

```javascript
const shortcuts = {
  // 导航
  'CmdOrCtrl+1': '切换到首页',
  'CmdOrCtrl+2': '切换到项目管理',
  'CmdOrCtrl+3': '切换到载荷生成',
  'CmdOrCtrl+4': '切换到插件管理',
  'CmdOrCtrl+,': '打开设置',
  
  // 操作
  'CmdOrCtrl+N': '新建项目',
  'CmdOrCtrl+S': '保存',
  'CmdOrCtrl+F': '搜索',
  'CmdOrCtrl+R': '刷新',
  
  // 编辑
  'CmdOrCtrl+Z': '撤销',
  'CmdOrCtrl+Y': '重做',
  'Delete': '删除选中项',
  'F2': '重命名',
  
  // 视图
  'CmdOrCtrl++': '放大',
  'CmdOrCtrl+-': '缩小',
  'CmdOrCtrl+0': '重置缩放',
  'F11': '全屏',
  'Escape': '关闭弹窗/取消操作'
}
```

---

## 跨平台适配规范

### Windows 特定优化

```css
/* 使用系统字体 */
.platform-windows {
  --font-sans: 'Segoe UI', 'Microsoft YaHei', sans-serif;
  --font-mono: 'Cascadia Code', Consolas, monospace;
  
  /* 窗口控制按钮位置 */
  .window-controls {
    order: 1; /* 右侧 */
  }
}
```

### macOS 特定优化

```css
/* 使用系统字体 */
.platform-macos {
  --font-sans: -apple-system, BlinkMacSystemFont, 'SF Pro Text', sans-serif;
  --font-mono: 'SF Mono', Monaco, monospace;
  
  /* 交通灯窗口控制 */
  .window-controls {
    order: 0; /* 左侧 */
    
    .control {
      border-radius: 50%; /* 圆形按钮 */
    }
  }
  
  /* 毛玻璃效果 */
  .sidebar {
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
  }
}
```

### Linux 特定优化

```css
/* 使用开源字体 */
.platform-linux {
  --font-sans: 'Inter', 'Roboto', 'Ubuntu', sans-serif;
  --font-mono: 'Fira Code', 'JetBrains Mono', monospace;
}
```

---

## 性能优化指南

### 1. CSS 性能

- 使用 `transform` 和 `opacity` 实现动画（GPU 加速）
- 避免使用 `*` 通配符选择器
- 使用 CSS 变量而非预处理器变量
- 合理使用 `will-change` 属性

### 2. 渲染优化

- 大列表使用虚拟滚动
- 图片使用懒加载
- 复杂计算使用 Web Workers
- 防抖/节流频繁触发的操作

### 3. 代码分割

```javascript
// 路由懒加载
const ProjectsContent = () => import('./components/ProjectsContent.vue')
const PayloadsContent = () => import('./components/PayloadsContent.vue')
```

---

## 实施路线图

### 阶段 1：基础系统（1-2 周）
- [ ] 更新 CSS 变量系统
- [ ] 优化字体和排版
- [ ] 实现新的间距和圆角系统
- [ ] 更新阴影系统

### 阶段 2：组件升级（2-3 周）
- [ ] 重构导航组件
- [ ] 升级按钮和输入框
- [ ] 优化卡片和表格
- [ ] 改进模态框和弹窗

### 阶段 3：交互优化（1-2 周）
- [ ] 实现动画系统
- [ ] 添加微交互反馈
- [ ] 实现键盘快捷键
- [ ] 优化加载状态

### 阶段 4：跨平台适配（1 周）
- [ ] Windows 特定优化
- [ ] macOS 特定优化
- [ ] Linux 特定优化
- [ ] 统一测试和验证

### 阶段 5：测试和迭代（持续）
- [ ] 用户测试
- [ ] 性能测试
- [ ] 可访问性测试
- [ ] 持续优化

---

## 设计资源

### 参考设计系统
- [Material Design 3](https://m3.material.io/)
- [Apple Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/)
- [Fluent Design System](https://www.microsoft.com/design/fluent/)
- [Carbon Design System](https://carbondesignsystem.com/)

### 工具推荐
- **设计工具**: Figma, Sketch, Adobe XD
- **原型工具**: Principle, Framer
- **色彩工具**: Coolors, Adobe Color
- **字体工具**: Font Pair, Typewolf

---

## 总结

本设计规范在**严格保留所有现有功能**的基础上，对 FG-ABYSS 的 UI 进行了全面升级：

### 核心改进
1. ✅ **视觉系统**：更现代的色彩、字体、间距系统
2. ✅ **组件库**：统一、可复用的组件设计
3. ✅ **交互体验**：流畅的动画和微交互
4. ✅ **跨平台**：适配 Windows、macOS、Linux
5. ✅ **性能优化**：GPU 加速、代码分割、懒加载

### 设计承诺
- 🔒 **功能完整性**：不修改或移除任何现有功能
- 🎨 **一致性**：统一的视觉语言和交互模式
- 🚀 **效率提升**：优化工作流程，减少操作步骤
- ♿ **可访问性**：符合 WCAG 2.1 AA 标准
- 📱 **响应式**：适配不同屏幕尺寸

这套设计规范将为 FG-ABYSS 打造专业、高效、美观的用户体验！
