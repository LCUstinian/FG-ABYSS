<template>
  <div class="sidebar">
    <div class="nav-item" 
         v-for="item in navItems" 
         :key="item.id"
         :class="{ active: currentNavItem === item.id }"
         @click="switchNavItem(item.id)"
    >
      <component :is="item.icon" :size="24" />
      <span>{{ item.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Home, Folder, Package, Plug, Settings } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps({
  currentNavItem: {
    type: String,
    default: 'home'
  }
})

const emit = defineEmits(['switch-nav'])

const switchNavItem = (itemId: string) => {
  emit('switch-nav', itemId)
}

// 导航项配置
const navItems = computed(() => [
  { id: 'home', label: t('nav.home'), icon: Home },
  { id: 'projects', label: t('nav.projects'), icon: Folder },
  { id: 'payloads', label: t('nav.payloads'), icon: Package },
  { id: 'plugins', label: t('nav.plugins'), icon: Plug },
  { id: 'settings', label: t('nav.settings'), icon: Settings }
])
</script>

<style scoped>
/* 侧边栏样式 - 深色主题风格 */
.sidebar {
  width: 90px;
  background: var(--sidebar-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
  gap: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
}

.nav-item {
  width: 75px;
  height: 75px;
  min-width: 75px;
  min-height: 75px;
  max-width: 75px;
  max-height: 75px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  gap: 8px;
  padding: 10px 6px;
  box-sizing: border-box;
  border-radius: 8px;
}

.nav-item:hover {
  background: var(--hover-color);
  transform: translateY(-2px);
}

.nav-item.active {
  background: var(--active-color);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.nav-item span {
  font-size: 12px;
  text-align: center;
  font-weight: 500;
  letter-spacing: 0.3px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
}

.nav-item svg {
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 70px;
    padding: 16px 0;
    gap: 8px;
  }
  
  .nav-item {
    width: 58px;
    height: 62px;
    padding: 8px 4px;
    gap: 6px;
  }
  
  .nav-item span {
    font-size: 10px;
  }
}
</style>