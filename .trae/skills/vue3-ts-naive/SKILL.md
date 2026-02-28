---
name: "vue3-ts-naive"
description: "Provides assistance for Vue 3 + TypeScript + Naive UI projects, including component creation, styling, theming, and best practices. Invoke when working on Vue 3 projects using TypeScript and Naive UI."
---

# Vue 3 + TypeScript + Naive UI Assistant

This skill provides specialized assistance for projects using Vue 3, TypeScript, and Naive UI component library. It helps with:

## Features

### Component Creation
- Generates Vue 3 Composition API components with TypeScript
- Creates Naive UI component templates with proper props and events
- Provides best practices for component structure and organization

### Theming and Styling
- Assists with Naive UI theme configuration
- Helps with CSS variables and dark/light mode implementation
- Provides guidance on responsive design patterns

### Internationalization
- Supports vue-i18n integration and usage
- Helps with translation file management
- Provides examples for multi-language support

### Best Practices
- TypeScript type definitions and interfaces
- Vue 3 Composition API usage patterns
- Naive UI component best practices
- Code organization and project structure

## Usage Examples

### Creating a new component
```vue
<template>
  <NCard>
    <template #header>
      <h3>{{ title }}</h3>
    </template>
    <div class="content">
      {{ content }}
    </div>
  </NCard>
</template>

<script setup lang="ts">
import { NCard } from 'naive-ui'

const props = defineProps<{
  title: string
  content: string
}>()
</script>
```

### Theme configuration
```typescript
import { createTheme } from 'naive-ui'

const theme = createTheme({
  common: {
    primaryColor: '#1890ff',
    primaryColorHover: '#40a9ff',
    primaryColorPressed: '#096dd9'
  }
})
```

### Internationalization setup
```typescript
import { createI18n } from 'vue-i18n'

const i18n = createI18n({
  locale: 'zh-CN',
  messages: {
    'zh-CN': require('./zh-CN.ts'),
    'en-US': require('./en-US.ts')
  }
})
```

## Project Structure

- `src/components/` - Vue components
- `src/i18n/` - Internationalization files
- `src/styles/` - Global styles and themes
- `src/utils/` - Utility functions
- `src/types/` - TypeScript type definitions

## Common Issues and Solutions

### TypeScript errors
- Ensure proper type definitions for props and refs
- Use `defineProps` and `defineEmits` with TypeScript generics
- Check for null/undefined values with optional chaining

### Naive UI component issues
- Verify proper import of components
- Check theme configuration for consistent styling
- Ensure responsive design for different screen sizes

### Theming problems
- Use CSS variables for consistent theming
- Test both light and dark modes
- Ensure proper contrast and accessibility

This skill is designed to help developers working on Vue 3 + TypeScript + Naive UI projects by providing best practices, code examples, and troubleshooting assistance.