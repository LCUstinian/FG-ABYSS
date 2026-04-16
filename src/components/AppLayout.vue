<template>
  <div class="app-layout">
    <CustomTitlebar v-if="showTitlebar" />
    <div class="app-container">
      <div class="app-sidebar" :style="{ width: sidebarWidth + 'px' }">
        <div class="sidebar-header">
          <n-icon size="24" color="#409eff">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path fill="currentColor" d="M12 2A10 10 0 0 0 2 12a10 10 0 0 0 10 10a10 10 0 0 0 10-10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 0-16a8 8 0 0 1 0 16z"/>
              <path fill="currentColor" d="M12 6a1 1 0 0 0-1 1v5a1 1 0 0 0 2 0V7a1 1 0 0 0-1-1zm0 8a1 1 0 1 0 0 2a1 1 0 0 0 0-2z"/>
            </svg>
          </n-icon>
          <span class="sidebar-brand">FG-ABYSS</span>
        </div>
        <div class="sidebar-nav">
          <n-menu
            v-model:value="activeKey"
            mode="vertical"
            :options="menuOptions"
            @update:value="handleMenuSelect"
          />
        </div>
      </div>
      <div class="app-content">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { RouterLink } from 'vue-router'
import { NIcon } from 'naive-ui'
import CustomTitlebar from './CustomTitlebar.vue'
import type { MenuOption } from 'naive-ui'

interface Props {
  showTitlebar?: boolean
  sidebarWidth?: number
}

withDefaults(defineProps<Props>(), {
  showTitlebar: true,
  sidebarWidth: 220,
})

const activeKey = ref<string>('home')

const renderIcon = (iconPath: string) => {
  return () => h(NIcon, null, {
    default: () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', viewBox: '0 0 24 24', innerHTML: iconPath })
  })
}

const menuOptions: MenuOption[] = [
  {
    key: 'home',
    label: () => h(RouterLink, { to: '/' }, { default: () => '首页' }),
    icon: renderIcon('<path fill="currentColor" d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/>'),
  },
  {
    key: 'project',
    label: () => h(RouterLink, { to: '/project' }, { default: () => '项目' }),
    icon: renderIcon('<path fill="currentColor" d="M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8z"/>'),
  },
  {
    key: 'payload',
    label: () => h(RouterLink, { to: '/payload' }, { default: () => '载荷' }),
    icon: renderIcon('<path fill="currentColor" d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" points="14 2 14 8 20 8"/>'),
  },
  {
    key: 'plugin',
    label: () => h(RouterLink, { to: '/plugin' }, { default: () => '插件' }),
    icon: renderIcon('<path fill="currentColor" d="M20.5 11H19V7c0-1.1-.9-2-2-2h-4V3.5C13 2.12 11.88 1 10.5 1S8 2.12 8 3.5V5H4c-1.1 0-2 .9-2 2v4H1.5C.12 11-1 12.12-1 13.5S.12 16 1.5 16H4v4c0 1.1.9 2 2 2h4v1.5c0 1.38 1.12 2.5 2.5 2.5s2.5-1.12 2.5-2.5V22h4c1.1 0 2-.9 2-2v-4h1.5c1.38 0 2.5-1.12 2.5-2.5s-1.12-2.5-2.5-2.5z"/>'),
  },
  {
    key: 'console',
    label: () => h(RouterLink, { to: '/console' }, { default: () => '控制台' }),
    icon: renderIcon('<path fill="currentColor" d="M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 14H4V6h16v12zM6 10h2v2H6zm0 4h2v2H6zm4-4h8v2h-8zm0 4h5v2h-5z"/>'),
  },
  {
    key: 'settings',
    label: () => h(RouterLink, { to: '/settings' }, { default: () => '设置' }),
    icon: renderIcon('<path fill="currentColor" d="M19.14 12.94c.04-.3.06-.61.06-.94c0-.32-.02-.64-.07-.94l2.03-1.58a.49.49 0 0 0 .12-.61l-1.92-3.32a.48.48 0 0 0-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.484.484 0 0 0-.48-.41h-3.84a.484.484 0 0 0-.48.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96a.48.48 0 0 0-.59.22L2.8 8.87a.49.49 0 0 0 .12.61l2.03 1.58c-.05.3-.07.62-.07.94s.02.64.07.94l-2.03 1.58a.49.49 0 0 0-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.27.41.48.41h3.84c.24 0 .44-.17.48-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32a.49.49 0 0 0-.12-.61zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6s3.6 1.62 3.6 3.6s-1.62 3.6-3.6 3.6z"/>'),
  },
]

const handleMenuSelect = (key: string) => {
  activeKey.value = key
}
</script>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  width: 100vw;
  height: 100vh;
  background: #0f0f23;
  overflow: hidden;
  padding-top: 40px; /* 为固定标题栏留出空间 */
  box-sizing: border-box;
}

.app-container {
  display: flex;
  flex: 1;
  overflow: hidden;
  min-height: 0;
}

.app-sidebar {
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

.sidebar-header {
  display: flex;
  align-items: center;
  padding: 16px;
  gap: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

.sidebar-brand {
  font-size: 18px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 1px;
  white-space: nowrap;
}

.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  overflow-y: auto;
  overflow-x: hidden;
}

.app-content {
  flex: 1;
  overflow: auto;
  background: #0f0f23;
  min-width: 0;
  min-height: 0;
}

:deep(.n-menu) {
  --n-item-text-color: #ffffff;
  --n-item-text-color-hover: #ffffff;
  --n-item-color-hover: rgba(64, 158, 255, 0.15);
  --n-item-color-active: rgba(64, 158, 255, 0.2);
  --n-item-height: 44px;
}
</style>
