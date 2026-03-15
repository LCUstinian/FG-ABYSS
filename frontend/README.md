# FG-ABYSS Frontend

Vue 3 + TypeScript + Vite 前端应用

## 📦 技术栈

- **框架**: Vue 3.4+ (Composition API + `<script setup>`)
- **语言**: TypeScript 5+
- **构建工具**: Vite 5+
- **UI 框架**: Naive UI 2.37+
- **测试框架**: Vitest + Vue Test Utils
- **国际化**: Vue I18n

## 🚀 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

启动开发服务器，支持热重载（HMR）

### 构建

```bash
# 开发构建
npm run build:dev

# 生产构建
npm run build
```

### 预览构建结果

```bash
npm run preview
```

## 📁 项目结构

```
frontend/
├── public/              # 静态资源
├── src/
│   ├── api/            # API 调用
│   │   └── system.ts
│   ├── components/     # Vue 组件
│   │   ├── CreateProjectModal.vue
│   │   ├── CreateWebShellModal.vue
│   │   ├── HomeContent.vue
│   │   ├── ProjectsContent.vue
│   │   └── ...
│   ├── i18n/          # 国际化配置
│   │   ├── en-US.ts
│   │   ├── zh-CN.ts
│   │   └── index.ts
│   ├── styles/        # 全局样式
│   │   └── global.css
│   ├── types/         # TypeScript 类型定义
│   │   └── wails.d.ts
│   ├── utils/         # 工具函数
│   │   ├── formatTime.ts
│   │   └── formatTime.test.ts
│   ├── App.vue        # 根组件
│   ├── main.ts        # 入口文件
│   └── shims-vue.d.ts # Vue 类型声明
├── index.html         # HTML 模板
├── package.json       # 依赖配置
├── tsconfig.json      # TypeScript 配置
├── vite.config.ts     # Vite 配置
└── vitest.config.ts   # Vitest 配置
```

## 🧪 测试

### 运行测试

```bash
npm run test
```

### 监视模式

```bash
npm run test:watch
```

### 生成覆盖率报告

```bash
npm run test:coverage
```

##  代码规范

### 组件命名

- 使用 PascalCase
- 多单词命名（至少两个单词）
- 语义化命名

```vue
<!-- ✅ 正确 -->
ProjectsContent.vue
CreateProjectModal.vue

<!-- ❌ 错误 -->
Content.vue
projects-content.vue
```

### Props 定义

```typescript
// ✅ 推荐：使用 interface + withDefaults
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
```

### 事件命名

```typescript
// ✅ 推荐：使用 kebab-case
emit('project-created', data)
emit('update:loading', false)
```

## 🌐 国际化

支持中文和英文：

```typescript
// 使用示例
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const text = t('common.submit') // 提交 / Submit
```

## 🔧 IDE 配置

### 推荐扩展

- [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (Vue 3 支持)
- [TypeScript Vue Plugin](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin)

### Take Over 模式

为了获得更好的性能，可以启用 Volar 的 Take Over 模式：

1. 禁用 VSCode 内置 TypeScript 扩展
2. 重新加载 VSCode 窗口

详见：[Volar Take Over Mode](https://github.com/johnsoncodehk/volar/discussions/471)

## 📚 开发指南

### 创建新组件

```vue
<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Default Title'
})

const emit = defineEmits<{
  (e: 'update', value: string): void
}>()

const count = ref(0)
</script>

<template>
  <div>{{ title }}: {{ count }}</div>
</template>
```

### API 调用

```typescript
// src/api/system.ts
import { invoke } from '@wailsapp/runtime'

export async function getProjects() {
  return await invoke('GetAllProjects')
}

export async function createProject(name: string) {
  return await invoke('CreateProject', { name })
}
```

### 使用 Wails API

```typescript
import { EventsOn, EventsOff } from '@wailsapp/runtime'

// 监听事件
EventsOn('project-created', (project) => {
  console.log('Project created:', project)
})

// 清理
EventsOff('project-created')
```

## 🛠️ 构建配置

### Vite 配置

详见 [`vite.config.ts`](./vite.config.ts)

### TypeScript 配置

详见 [`tsconfig.json`](./tsconfig.json)

## 📊 测试覆盖率

目标覆盖率：
- 工具函数：100%
- 组件逻辑：80%+
- API 调用：50%+

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

遵循项目主仓库的许可证。

---

**注意**: 此前端应用通过 Wails 与 Go 后端集成，单独运行时部分功能可能不可用。
