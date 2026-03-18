<template>
  <div class="settings-panel">
    <!-- 页面标题 -->
    <PageHeader 
      title="设置" 
      subtitle="自定义应用外观与行为" 
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
            {{ tab.label }}
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
  { name: 'appearance', label: '外观' },
  { name: 'connection', label: '连接' },
  { name: 'about', label: '关于' }
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
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
  cursor: pointer;
  text-align: left;
}

.tab-item:hover {
  color: var(--text-primary);
  background: rgba(128, 128, 128, 0.1);
}

.tab-item.active {
  color: var(--active-color);
  font-weight: 600;
  border-left-color: var(--active-color);
  background: rgba(var(--active-color-rgb), 0.1);
}

.settings-main {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  overflow-x: hidden;
  background: var(--content-bg);
  padding: 40px 0;
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
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .settings-sidebar {
    width: 260px;
  }
  
  .tab-item {
    padding: 18px 28px;
    font-size: 15px;
  }
}
</style>
