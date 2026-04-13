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
import { ref, computed, h } from 'vue'
import { RouterLink } from 'vue-router'
import { NIcon } from 'naive-ui'
import CustomTitlebar from './CustomTitlebar.vue'
import type { MenuOption } from 'naive-ui'

interface Props {
  showTitlebar?: boolean
  sidebarWidth?: number
}

const props = withDefaults(defineProps<Props>(), {
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
    key: 'terminal',
    label: () => h(RouterLink, { to: '/terminal' }, { default: () => '终端' }),
    icon: renderIcon('<path fill="currentColor" d="M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm-8 14H4v-8h8v8zm0-10H4V6h8v2zm10 10h-8v-8h8v8zm0-10h-8V6h8v2z"/>'),
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
  width: 100%;
  height: 100%;
  background: #0f0f23;
}

.app-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.app-sidebar {
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar-header {
  display: flex;
  align-items: center;
  padding: 16px;
  gap: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.sidebar-brand {
  font-size: 18px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 1px;
}

.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  overflow-y: auto;
}

.app-content {
  flex: 1;
  overflow: auto;
  background: #0f0f23;
}

:deep(.n-menu) {
  --n-item-text-color: #ffffff;
  --n-item-text-color-hover: #ffffff;
  --n-item-color-hover: rgba(64, 158, 255, 0.15);
  --n-item-color-active: rgba(64, 158, 255, 0.2);
}
</style>
