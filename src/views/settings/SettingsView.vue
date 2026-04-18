<template>
  <div class="settings-view">
    <!-- Left submenu -->
    <div class="settings-nav">
      <button
        v-for="item in navItems"
        :key="item.key"
        class="settings-nav-item"
        :class="{ 'is-active': activeSection === item.key }"
        @click="activeSection = item.key"
      >
        <component :is="item.icon" :size="16" />
        <span>{{ item.label }}</span>
      </button>
    </div>

    <!-- Right panels -->
    <div class="settings-content">
      <h2 class="section-heading">{{ currentLabel }}</h2>
      <AppearancePanel v-if="activeSection === 'appearance'" />
      <div v-else class="placeholder-panel">
        <span style="color: var(--text-3); font-size: 13px">{{ currentLabel }} — 配置项待实现</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Palette, Wifi, ShieldCheck, ScrollText, Archive, Info } from 'lucide-vue-next'
import AppearancePanel from './AppearancePanel.vue'

const activeSection = ref('appearance')

const navItems = [
  { key: 'appearance', label: '外观',   icon: Palette      },
  { key: 'connection', label: '连接',   icon: Wifi         },
  { key: 'security',   label: '安全',   icon: ShieldCheck  },
  { key: 'logs',       label: '日志',   icon: ScrollText   },
  { key: 'backup',     label: '备份',   icon: Archive      },
  { key: 'about',      label: '关于',   icon: Info         },
]

const currentLabel = computed(() =>
  navItems.find(n => n.key === activeSection.value)?.label ?? ''
)
</script>

<style scoped>
.settings-view {
  height: 100%;
  display: flex;
  overflow: hidden;
}

.settings-nav {
  width: 180px;
  flex-shrink: 0;
  background: var(--bg-deep);
  border-right: 1px solid var(--border);
  padding: 8px 0;
  display: flex;
  flex-direction: column;
}

.settings-nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 16px;
  border: none;
  border-left: 2px solid transparent;
  background: transparent;
  color: var(--text-2);
  font-size: 13px;
  cursor: pointer;
  text-align: left;
  transition: background 80ms;
}

.settings-nav-item:hover { background: var(--bg-hover); color: var(--text-1); }
.settings-nav-item.is-active {
  border-left-color: var(--accent);
  background: var(--accent-bg);
  color: var(--accent);
  font-weight: 600;
}

.settings-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.section-heading {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-1);
  margin-bottom: 20px;
}

.placeholder-panel {
  padding: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
