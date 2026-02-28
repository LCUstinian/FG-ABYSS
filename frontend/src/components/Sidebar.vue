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
/* 侧边栏样式 - 现代风格 */
.sidebar {
  width: 96px;
  background: var(--sidebar-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 0;
  gap: 16px;
  transition: all var(--transition-normal);
  flex-shrink: 0;
  box-shadow: var(--shadow-sm);
}

.nav-item {
  width: 80px;
  height: 80px;
  min-width: 80px;
  min-height: 80px;
  max-width: 80px;
  max-height: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-normal);
  gap: 8px;
  padding: 12px 8px;
  box-sizing: border-box;
  border-radius: var(--border-radius-lg);
  position: relative;
  overflow: hidden;
}

.nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transition: left var(--transition-slow);
}

.nav-item:hover::before {
  left: 100%;
}

.nav-item:hover {
  background: var(--hover-color);
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
}

.nav-item.active {
  background: var(--active-color);
  color: white;
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.nav-item.active::before {
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
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
  position: relative;
  z-index: 1;
}

.nav-item svg {
  flex-shrink: 0;
  position: relative;
  z-index: 1;
  transition: transform var(--transition-normal);
}

.nav-item:hover svg {
  transform: scale(1.1);
}

.nav-item.active svg {
  transform: scale(1.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 76px;
    padding: 20px 0;
    gap: 12px;
  }
  
  .nav-item {
    width: 64px;
    height: 68px;
    padding: 10px 6px;
    gap: 6px;
  }
  
  .nav-item span {
    font-size: 10px;
  }
  
  .nav-item svg {
    width: 20px;
    height: 20px;
  }
}
</style>