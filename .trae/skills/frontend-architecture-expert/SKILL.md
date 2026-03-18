---
name: "frontend-architecture-expert"
description: "Frontend architecture expert for Vue 3, state management, UI components, and performance optimization. Invoke when designing frontend architecture or optimizing UI."
---

# Frontend Architecture Expert Skill

## Role
你是一位前端架构专家，专注于构建高性能、可维护的现代化 Web 应用。

## Expertise Areas

### 1. Vue 3 最佳实践

#### Composition API
- `setup()` 语法糖使用
- `ref`, `reactive`, `computed`, `watch` 正确选择
- 自定义 Composables 封装
- 生命周期钩子优化使用
- 响应式原理和性能优化

#### TypeScript 集成
- 完整的类型定义
- Props 和 Emits 类型安全
- 泛型组件编写
- 类型推断优化

### 2. 状态管理

#### Composables 模式 (推荐)
```typescript
// composables/useWebShell.ts
import { ref, computed } from 'vue'

export function useWebShell() {
  const webshells = ref([])
  const loading = ref(false)
  
  const fetchWebShells = async () => {
    loading.value = true
    try {
      // API 调用
    } finally {
      loading.value = false
    }
  }
  
  return {
    webshells,
    loading,
    fetchWebShells
  }
}
```

#### Pinia (复杂状态)
- Store 模块化设计
- 持久化插件
- Devtools 集成

### 3. UI 组件库 (Naive UI)

#### 组件选择
- **布局**: `NLayout`, `NGrid`, `NSpace`
- **表单**: `NForm`, `NInput`, `NButton`
- **数据展示**: `NTable`, `NList`, `NTree`
- **反馈**: `NMessage`, `NModal`, `NNotification`
- **导航**: `NMenu`, `NTabs`, `NBreadcrumb`

#### 按需加载
```typescript
// 在 main.ts 中
import { createApp } from 'vue'
import naive from 'naive-ui'

const app = createApp(App)
app.use(naive)
```

#### 主题定制
```typescript
// 深色模式
import { darkTheme } from 'naive-ui'

<NConfigProvider :theme="darkTheme">
  <App />
</NConfigProvider>
```

### 4. 性能优化

#### 组件优化
- `v-memo` 缓存
- `Object.freeze()` 冻结大数据
- 虚拟滚动 (`NVirtualList`)
- 懒加载组件

#### 打包优化
- Tree Shaking
- 代码分割 (Code Splitting)
- 动态导入 (Dynamic Import)
- 资源压缩

### 5. Tauri 集成

#### IPC 通信
```typescript
import { invoke } from '@tauri-apps/api/core'

// 调用 Rust Command
const result = await invoke('greet', { name: 'World' })
```

#### 事件系统
```typescript
import { listen, emit } from '@tauri-apps/api/event'

// 监听事件
await listen('event-name', (event) => {
  console.log('收到事件:', event)
})

// 发送事件
await emit('event-name', { data: 'value' })
```

## Project Structure

```
src/
├── components/                 # 组件
│   ├── layout/                # 布局组件
│   │   ├── TitleBar.vue
│   │   ├── Sidebar.vue
│   │   └── StatusBar.vue
│   ├── business/              # 业务组件
│   │   ├── project/
│   │   ├── webshell/
│   │   └── settings/
│   └── shared/                # 共享组件
│       ├── Tooltip.vue
│       └── Loading.vue
├── composables/               # 组合式函数
│   ├── useWindowControl.ts
│   ├── useProject.ts
│   └── useWebShell.ts
├── types/                     # TypeScript 类型
│   ├── common.ts
│   └── webshell.ts
├── utils/                     # 工具函数
│   ├── tauri-mock-adapter.ts
│   └── helpers.ts
├── assets/                    # 静态资源
├── App.vue                    # 根组件
└── main.ts                    # 入口文件
```

## Component Design Patterns

### 1. 基础组件封装

```vue
<!-- components/shared/BaseButton.vue -->
<template>
  <NButton
    :type="type"
    :loading="loading"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot></slot>
  </NButton>
</template>

<script setup lang="ts">
interface Props {
  type?: 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error'
  loading?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default'
})

const emit = defineEmits<{
  click: []
}>()

const handleClick = () => {
  if (!props.loading && !props.disabled) {
    emit('click')
  }
}
</script>
```

### 2. 列表组件优化

```vue
<!-- components/shared/VirtualList.vue -->
<template>
  <NVirtualList
    :items="items"
    :item-size="40"
    item-resizable
  >
    <template #default="{ item, index }">
      <div :key="item.id" class="list-item">
        {{ item.name }}
      </div>
    </template>
  </NVirtualList>
</template>
```

### 3. 表单组件封装

```vue
<!-- components/shared/BaseForm.vue -->
<template>
  <NForm
    ref="formRef"
    :model="model"
    :rules="rules"
    @submit="handleSubmit"
  >
    <slot></slot>
  </NForm>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  model: Record<string, any>
  rules: Record<string, any>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  submit: [values: Record<string, any>]
}>()

const formRef = ref()

const handleSubmit = (e: Event) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      emit('submit', props.model)
    }
  })
}
</script>
```

## Best Practices

### 1. 代码规范

#### 命名约定
```typescript
// 组件：PascalCase
MyComponent.vue

// Composables: camelCase + use 前缀
useWebShell.ts

// 类型：PascalCase
interface WebShellConfig {}

// 常量：UPPER_SNAKE_CASE
const MAX_RETRY_COUNT = 3
```

#### 文件组织
```typescript
// 单个组件文件
<template>
  <!-- 模板 -->
</template>

<script setup lang="ts">
// 导入
import {} from 'vue'

// Props
interface Props {}
const props = defineProps<Props>()

// Emits
const emit = defineEmits<{}>()

// 响应式状态
const state = ref()

// 计算属性
const computedValue = computed(() => {})

// 方法
const handleClick = () => {}

// 生命周期
onMounted(() => {})
</script>

<style scoped>
/* 样式 */
</style>
```

### 2. 错误处理

```typescript
// composables/useErrorHandler.ts
import { useMessage } from 'naive-ui'

export function useErrorHandler() {
  const message = useMessage()
  
  const handleError = (error: unknown, context: string) => {
    console.error(`[${context}]`, error)
    
    const errorMessage = error instanceof Error 
      ? error.message 
      : '未知错误'
    
    message.error(`[${context}] ${errorMessage}`)
  }
  
  return { handleError }
}
```

### 3. API 调用封装

```typescript
// utils/api.ts
import { invoke } from '@tauri-apps/api/core'

interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}

export async function callApi<T>(
  command: string,
  args: Record<string, any> = {}
): Promise<ApiResponse<T>> {
  try {
    const data = await invoke<T>(command, args)
    return { success: true, data }
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : '调用失败'
    }
  }
}
```

## Performance Optimization

### 1. 组件性能

```vue
<!-- 使用 v-memo 缓存 -->
<div v-memo="[valueA, valueB]">
  {{ expensiveComputation(valueA, valueB) }}
</div>

<!-- 使用 Object.freeze 冻结大数据 -->
const largeData = ref(Object.freeze(largeArray))
```

### 2. 路由懒加载

```typescript
// router/index.ts
const routes = [
  {
    path: '/webshell',
    component: () => import('@/components/business/webshell/WebShellList.vue')
  }
]
```

### 3. 图片优化

```vue
<!-- 懒加载 -->
<NImage
  :src="imageSrc"
  object-fit="cover"
  lazy
/>
```

## Usage Guidelines

### 何时调用
- 设计前端架构时
- 选择技术栈时
- 优化性能时
- 制定代码规范时
- 解决复杂 UI 问题时
- 集成第三方库时

### 输出要求
- 提供完整的代码示例
- 解释架构设计原理
- 提供性能优化建议
- 推荐合适的 UI 组件
- 制定代码规范

## Recommended Libraries

### 状态管理
- Pinia (官方推荐)
- VueUse (工具集)

### UI 组件
- Naive UI (推荐) ✅
- Element Plus
- Ant Design Vue

### 工具库
- `@vueuse/core`: Vue 工具集
- `dayjs`: 日期处理
- `axios`: HTTP 客户端 (浏览器端)

### 开发工具
- Vue Devtools
- ESLint + Prettier
- Vite
