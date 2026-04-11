<template>
  <div class="settings-panel">
    <!-- 页面标题 -->
    <PageHeader 
      :title="t('settings.title')" 
      :subtitle="t('settings.subtitle')" 
    />

    <!-- 左右布局容器 -->
    <div class="settings-content">
      <!-- 左侧 TAB 标签列 -->
      <div class="settings-sidebar">
        <div class="tab-list">
          <div 
            v-for="tab in tabs" 
            :key="tab.name"
            :class="['tab-item', { active: currentSettingsTab === tab.name }]"
            @click="currentSettingsTab = tab.name"
          >
            <span class="tab-icon">{{ tab.icon }}</span>
            <span class="tab-label">{{ tab.label }}</span>
          </div>
        </div>
      </div>
      
      <!-- 右侧内容区域 -->
      <div class="settings-main">
        <div v-if="currentSettingsTab === 'appearance'" class="tab-content">
          <AppearanceView 
            :theme-mode="themeMode"
            :current-language="currentLanguage"
            @update:theme-mode="handleThemeChange"
            @change-language="changeLanguage"
          />
        </div>
        <div v-if="currentSettingsTab === 'connection'" class="tab-content">
          <ConnectionView />
        </div>
        <div v-if="currentSettingsTab === 'audit'" class="tab-content">
          <AuditLogView />
        </div>
        <div v-if="currentSettingsTab === 'about'" class="tab-content">
          <AboutView />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/shared/PageHeader.vue'
import AppearanceView from './AppearanceView.vue'
import ConnectionView from './ConnectionView.vue'
import AboutView from './AboutView.vue'
import AuditLogView from './AuditLogView.vue'

const props = defineProps({
  isDarkTheme: {
    type: Boolean,
    default: false
  },
  themeMode: {
    type: String,
    default: 'system'
  }
})

const emit = defineEmits(['update:theme-mode'])

const { t, locale } = useI18n()

const currentSettingsTab = ref('appearance')
const currentLanguage = ref(locale.value)

const tabs = [
  { name: 'appearance', label: t('settings.appearance'), icon: '🎨' },
  { name: 'connection', label: t('settings.connection'), icon: '🔗' },
  { name: 'audit', label: '审计日志', icon: '📋' },
  { name: 'about', label: t('settings.about'), icon: 'ℹ️' }
]

watch(() => props.themeMode, (newVal) => {
  emit('update:theme-mode', newVal)
})

const handleThemeChange = (mode: string) => {
  emit('update:theme-mode', mode)
}

const changeLanguage = (lang: string) => {
  locale.value = lang
  currentLanguage.value = lang
}
</script>

<style scoped>
.settings-panel {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  background: var(--content-bg);
  overflow: hidden;
  transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.settings-content {
  flex: 1;
  display: flex;
  min-width: 0;
  overflow: hidden;
}

.settings-sidebar {
  width: 240px;
  flex-shrink: 0;
  border-right: 1px solid var(--border-color);
  background: var(--sidebar-bg);
  overflow-y: auto;
  overflow-x: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
}

.tab-list {
  display: flex;
  flex-direction: column;
  padding: 12px 0;
}

.tab-item {
  padding: 16px 24px;
  font-size: 15px;
  font-weight: 500;
  color: var(--text-secondary);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-left: 3px solid transparent;
  cursor: pointer;
  text-align: left;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 12px;
}

.tab-icon {
  font-size: 20px;
  line-height: 1;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.tab-label {
  font-size: 15px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.tab-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, rgba(var(--active-color-rgb), 0.05) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.tab-item:hover {
  color: var(--text-primary);
  background: rgba(128, 128, 128, 0.1);
  transform: translateX(4px);
}

.tab-item:hover .tab-icon {
  transform: scale(1.15) rotate(5deg);
}

.tab-item:hover .tab-label {
  font-weight: 600;
}

.tab-item:hover::before {
  opacity: 1;
}

.tab-item.active {
  color: var(--active-color);
  font-weight: 600;
  border-left-color: var(--active-color);
  background: rgba(var(--active-color-rgb), 0.1);
  box-shadow: 0 2px 8px rgba(var(--active-color-rgb), 0.15);
}

.tab-item.active .tab-icon {
  transform: scale(1.1);
  animation: iconPulse 2s infinite;
}

.tab-item.active .tab-label {
  font-weight: 600;
}

.tab-item.active::before {
  opacity: 1;
}

@keyframes iconPulse {
  0%, 100% {
    transform: scale(1.1);
  }
  50% {
    transform: scale(1.2);
  }
}

.settings-main {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  overflow-x: hidden;
  background: var(--content-bg);
  padding: 40px 0;
  transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 自定义滚动条 */
.settings-main::-webkit-scrollbar,
.settings-sidebar::-webkit-scrollbar {
  width: 8px;
}

.settings-main::-webkit-scrollbar-track,
.settings-sidebar::-webkit-scrollbar-track {
  background: var(--sidebar-bg);
  border-radius: 4px;
}

.settings-main::-webkit-scrollbar-thumb,
.settings-sidebar::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.settings-main::-webkit-scrollbar-thumb:hover,
.settings-sidebar::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .settings-main {
    padding: 60px 0;
  }
}

.tab-content {
  width: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .settings-sidebar {
    width: 260px;
  }
  
  .tab-item {
    padding: 18px 28px;
    gap: 14px;
  }
  
  .tab-icon {
    font-size: 22px;
  }
  
  .tab-label {
    font-size: 15px;
  }
  
  .settings-main {
    padding: 60px 0;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .settings-content {
    flex-direction: column;
  }
  
  .settings-sidebar {
    width: 100%;
    height: auto;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
  }
  
  .tab-list {
    flex-direction: row;
    padding: 0;
    overflow-x: auto;
  }
  
  .tab-item {
    white-space: nowrap;
    border-left: none;
    border-bottom: 3px solid transparent;
    padding: 16px 20px;
    flex-direction: column;
    gap: 8px;
  }
  
  .tab-icon {
    font-size: 24px;
  }
  
  .tab-label {
    font-size: 13px;
  }
  
  .tab-item.active {
    border-left: none;
    border-bottom-color: var(--active-color);
  }
  
  .settings-main {
    padding: 24px 0;
  }
}
</style>
