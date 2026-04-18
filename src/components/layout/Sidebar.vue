<template>
  <nav
    class="sidebar"
    :class="{ 'is-collapsed': collapsed }"
    role="navigation"
    aria-label="主导航"
  >
    <div class="nav-items">
      <n-tooltip
        v-for="item in navItems"
        :key="item.path"
        placement="right"
        :delay="400"
        :disabled="!collapsed"
      >
        <template #trigger>
          <RouterLink
            :to="item.path"
            class="nav-item"
            :class="{ 'is-active': route.path === item.path }"
            :aria-label="collapsed ? item.label : undefined"
          >
            <component :is="item.icon" :size="18" class="nav-icon" />
            <span class="nav-label">{{ item.label }}</span>
          </RouterLink>
        </template>
        {{ item.label }}
      </n-tooltip>
    </div>

    <button
      class="collapse-btn"
      :aria-label="collapsed ? '展开侧边栏' : '折叠侧边栏'"
      @click="collapsed = !collapsed"
    >
      <ChevronLeft
        :size="16"
        class="collapse-icon"
        :class="{ 'is-rotated': collapsed }"
      />
      <span class="nav-label">折叠</span>
    </button>
  </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useLocalStorage } from '@vueuse/core'
import { Home, FolderOpen, Package, Puzzle, Settings } from 'lucide-vue-next'
import { ChevronLeft } from 'lucide-vue-next'

const route = useRoute()
const collapsed = useLocalStorage('fg-sidebar-collapsed', false)

const navItems = [
  { path: '/',         label: '首页', icon: Home       },
  { path: '/project',  label: '项目', icon: FolderOpen  },
  { path: '/payload',  label: '载荷', icon: Package     },
  { path: '/plugin',   label: '插件', icon: Puzzle      },
  { path: '/settings', label: '设置', icon: Settings    },
]
</script>

<style scoped>
.sidebar {
  height: 100%;
  background: var(--bg-deep);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.nav-items {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 12px;
  text-decoration: none;
  color: var(--text-2);
  border-left: 2px solid transparent; /* permanent — prevents 2px jitter on activation */
  cursor: pointer;
  transition: background 80ms ease-out, color 80ms ease-out;
  white-space: nowrap;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-1);
}

.nav-item.is-active {
  border-left-color: var(--accent);
  background: var(--accent-bg);
  color: var(--accent);
  font-weight: 600;
}

.nav-item.is-active .nav-icon {
  color: var(--accent);
}

.nav-icon {
  flex-shrink: 0;
}

.nav-label {
  font-size: 13px;
  overflow: hidden;
  white-space: nowrap;
  transition: opacity 200ms ease-out, width 200ms ease-out;
  width: auto;
  opacity: 1;
}

/* Collapsed: hide text */
.sidebar.is-collapsed .nav-label {
  opacity: 0;
  width: 0;
}

.sidebar.is-collapsed .nav-item {
  padding: 0;
  justify-content: center;
}

.collapse-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 12px;
  border: none;
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  border-top: 1px solid var(--border);
  width: 100%;
  transition: background 80ms ease-out, color 80ms ease-out;
  white-space: nowrap;
}

.collapse-btn:hover {
  background: var(--bg-hover);
  color: var(--text-2);
}

.collapse-icon {
  flex-shrink: 0;
  transition: transform 220ms cubic-bezier(0.4, 0, 0.2, 1);
}

.collapse-icon.is-rotated {
  transform: rotate(180deg);
}

.sidebar.is-collapsed .collapse-btn {
  padding: 0;
  justify-content: center;
}
</style>
