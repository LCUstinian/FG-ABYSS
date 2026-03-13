<template>
  <div class="sidebar">
    <div class="nav-item" 
         v-for="item in navItems" 
         :key="item.id"
         :class="{ active: currentNavItem === item.id }"
         @click="switchNavItem(item.id)"
    >
      <component :is="item.icon" :size="24" class="nav-icon" />
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
/* 侧边栏样式 - 现代化设计 2.0 */
.sidebar {
  width: 96px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--space-6) 0;
  gap: var(--space-4);
  transition: all var(--duration-normal) var(--ease-standard);
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
  transition: all var(--duration-normal) var(--ease-emphasis);
  gap: var(--space-2);
  padding: var(--space-3) var(--space-2);
  box-sizing: border-box;
  border-radius: var(--radius-lg);
  position: relative;
  overflow: hidden;
}

/* 渐变闪光动画 */
.nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transition: left var(--duration-slower) var(--ease-smooth);
  z-index: 0;
}

.nav-item:hover::before {
  left: 100%;
}

/* 悬停状态 */
.nav-item:hover {
  background: var(--bg-hover);
  transform: translateY(-3px) scale(1.02);
  box-shadow: var(--shadow-md);
}

.nav-item:hover .nav-icon {
  transform: scale(1.15);
}

/* 激活状态 - 使用渐变背景 */
.nav-item.active {
  background: linear-gradient(135deg, var(--primary-500), var(--primary-600));
  color: white;
  box-shadow: var(--shadow-lg), 0 0 0 2px rgba(59, 130, 246, 0.2);
  transform: translateY(-2px);
}

.nav-item.active::before {
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.25), transparent);
}

.nav-item.active .nav-icon {
  transform: scale(1.1);
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

.nav-item span {
  font-size: var(--text-xs);
  text-align: center;
  font-weight: var(--font-medium);
  letter-spacing: 0.3px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  position: relative;
  z-index: 1;
  transition: color var(--duration-normal) var(--ease-standard);
}

.nav-icon {
  flex-shrink: 0;
  position: relative;
  z-index: 1;
  transition: transform var(--duration-normal) var(--ease-emphasis);
  color: inherit;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 76px;
    padding: var(--space-5) 0;
    gap: var(--space-3);
  }
  
  .nav-item {
    width: 64px;
    height: 68px;
    padding: var(--space-2) var(--space-1);
    gap: var(--space-1);
  }
  
  .nav-item span {
    font-size: 10px;
  }
  
  .nav-icon {
    width: 20px;
    height: 20px;
  }
}
</style>
