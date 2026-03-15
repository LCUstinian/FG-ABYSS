<script setup lang="ts">
import TitleBar from './components/TitleBar.vue'
import StatusBar from './components/StatusBar.vue'
import Sidebar from './components/Sidebar.vue'
import HomeContent from './components/HomeContent.vue'
import ProjectsContent from './components/ProjectsContent.vue'
import PayloadsContent from './components/PayloadsContent.vue'
import PluginsContent from './components/PluginsContent.vue'
import SettingsContent from './components/SettingsContent.vue'
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { 
  NConfigProvider, 
  NMessageProvider, 
  NDialogProvider,
  darkTheme, 
  lightTheme 
} from 'naive-ui'
import { getSystemStatus, type SystemStatus } from './api/system'

const { t, locale } = useI18n()

const isDarkTheme = ref(localStorage.getItem('theme') === 'dark' || window.matchMedia('(prefers-color-scheme: dark)').matches)
const themeMode = ref(localStorage.getItem('themeMode') || 'system')
const currentNavItem = ref('home')

// 辅助函数：安全地获取 CSS 变量值
const getVar = (name: string, defaultVal: string): string => {
  if (typeof window === 'undefined') return defaultVal
  const val = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  return val || defaultVal
}

// 辅助函数：颜色变亮（用于 hover 状态）
const lightenColor = (color: string, percent: number): string => {
  const num = parseInt(color.replace('#', ''), 16)
  const amt = Math.round(2.55 * percent)
  const R = (num >> 16) + amt
  const G = ((num >> 8) & 0x00FF) + amt
  const B = (num & 0x0000FF) + amt
  return '#' + (0x1000000 + 
    (R < 255 ? (R < 1 ? 0 : R) : 255) * 0x10000 + 
    (G < 255 ? (G < 1 ? 0 : G) : 255) * 0x100 + 
    (B < 255 ? (B < 1 ? 0 : B) : 255)
  ).toString(16).slice(1)
}

// 辅助函数：颜色变暗（用于 pressed 状态）
const darkenColor = (color: string, percent: number): string => {
  const num = parseInt(color.replace('#', ''), 16)
  const amt = Math.round(2.55 * percent)
  const R = (num >> 16) - amt
  const G = ((num >> 8) & 0x00FF) - amt
  const B = (num & 0x0000FF) - amt
  return '#' + (0x1000000 + 
    (R < 255 ? (R < 1 ? 0 : R) : 255) * 0x10000 + 
    (G < 255 ? (G < 1 ? 0 : G) : 255) * 0x100 + 
    (B < 255 ? (B < 1 ? 0 : B) : 255)
  ).toString(16).slice(1)
}

// 辅助函数：添加透明度到颜色
const addOpacity = (color: string, opacity: number): string => {
  const alpha = Math.round(opacity * 255)
  return color + alpha.toString(16).padStart(2, '0')
}

// 定义主题覆盖配置 - 使用全局强调色 (在设置中可自定义)
const themeOverrides = computed(() => {
  // 触发重新计算（即使 CSS 变量没变，trigger 变化也会触发）
  themeTrigger.value
  
  // 从 CSS 变量中读取当前强调色
  const primaryColor = getVar('--active-color', '#3b82f6')
  
  // 计算衍生颜色
  const primaryColorHover = lightenColor(primaryColor, 10)
  const primaryColorPressed = darkenColor(primaryColor, 10)
  const primaryColorSuppl = addOpacity(primaryColor, 0.1)
  const primaryColorActive = primaryColor
  
  // 成功状态也使用强调色
  const successColor = primaryColor
  const successColorHover = primaryColorHover
  const successColorPressed = primaryColorPressed
  const successColorSuppl = primaryColorSuppl
  const successColorActive = primaryColorActive
  
  // 信息颜色使用强调色
  const infoColor = primaryColor
  const infoColorHover = primaryColorHover
  const infoColorPressed = primaryColorPressed
  const infoColorSuppl = primaryColorSuppl
  const infoColorActive = primaryColorActive
  
  // 警告颜色 (保持黄色系)
  const warningColor = '#d97706'
  const warningColorHover = lightenColor(warningColor, 10)
  const warningColorPressed = darkenColor(warningColor, 10)
  const warningColorSuppl = addOpacity(warningColor, 0.1)
  const warningColorActive = warningColor
  
  // 错误颜色 (保持红色系)
  const errorColor = '#dc2626'
  const errorColorHover = lightenColor(errorColor, 10)
  const errorColorPressed = darkenColor(errorColor, 10)
  const errorColorSuppl = addOpacity(errorColor, 0.1)
  const errorColorActive = errorColor
  
  return {
    common: {
      primaryColor,
      primaryColorHover,
      primaryColorPressed,
      primaryColorSuppl,
      primaryColorActive,
      successColor,
      successColorHover,
      successColorPressed,
      successColorSuppl,
      successColorActive,
      infoColor,
      infoColorHover,
      infoColorPressed,
      infoColorSuppl,
      infoColorActive,
      warningColor,
      warningColorHover,
      warningColorPressed,
      warningColorSuppl,
      warningColorActive,
      errorColor,
      errorColorHover,
      errorColorPressed,
      errorColorSuppl,
      errorColorActive,
    }
  }
})

// 深色模式的主题覆盖
const darkThemeOverrides = computed(() => {
  // 触发重新计算
  themeTrigger.value
  
  // 从 CSS 变量中读取当前强调色
  const primaryColor = getVar('--active-color', '#3b82f6')
  
  // 深色模式下使用不同的计算方式
  const primaryColorHover = lightenColor(primaryColor, 15)
  const primaryColorPressed = darkenColor(primaryColor, 5)
  const primaryColorSuppl = addOpacity(primaryColor, 0.15)
  const primaryColorActive = primaryColor
  
  // 成功状态也使用强调色
  const successColor = primaryColor
  const successColorHover = primaryColorHover
  const successColorPressed = primaryColorPressed
  const successColorSuppl = primaryColorSuppl
  const successColorActive = primaryColorActive
  
  // 信息颜色使用强调色
  const infoColor = primaryColor
  const infoColorHover = primaryColorHover
  const infoColorPressed = primaryColorPressed
  const infoColorSuppl = primaryColorSuppl
  const infoColorActive = primaryColorActive
  
  // 警告颜色 (保持黄色系)
  const warningColor = '#fbbf24'
  const warningColorHover = lightenColor(warningColor, 15)
  const warningColorPressed = darkenColor(warningColor, 10)
  const warningColorSuppl = addOpacity(warningColor, 0.2)
  const warningColorActive = warningColor
  
  // 错误颜色 (保持红色系)
  const errorColor = '#ef4444'
  const errorColorHover = lightenColor(errorColor, 15)
  const errorColorPressed = darkenColor(errorColor, 10)
  const errorColorSuppl = addOpacity(errorColor, 0.2)
  const errorColorActive = errorColor
  
  return {
    common: {
      primaryColor,
      primaryColorHover,
      primaryColorPressed,
      primaryColorSuppl,
      primaryColorActive,
      successColor,
      successColorHover,
      successColorPressed,
      successColorSuppl,
      successColorActive,
      infoColor,
      infoColorHover,
      infoColorPressed,
      infoColorSuppl,
      infoColorActive,
      warningColor,
      warningColorHover,
      warningColorPressed,
      warningColorSuppl,
      warningColorActive,
      errorColor,
      errorColorHover,
      errorColorPressed,
      errorColorSuppl,
      errorColorActive,
    }
  }
})

// 系统状态数据
const systemStatus = ref<SystemStatus>({
  memoryUsage: '0 GB / 0 GB',
  processId: '0',
  cpuUsage: '0%',
  uptime: '0 秒'
})

// 用于触发主题重新计算的信号
const themeTrigger = ref(0)

// 监听 localStorage 中 accentColor 的变化
const handleStorageChange = (event: StorageEvent) => {
  if (event.key === 'accentColor' && event.newValue) {
    // CSS 变量已经在 SettingsContent 中更新
    // 触发 themeOverrides 重新计算
    themeTrigger.value++
  }
}

// 定时器引用
let updateTimer: number | null = null

// 计算当前内容组件
const currentContent = computed(() => {
  return currentNavItem.value
})

// 切换导航项
const switchNavItem = (itemId: string) => {
  currentNavItem.value = itemId
}

// 处理主题切换
const handleThemeChange = (mode: string) => {
  themeMode.value = mode
  localStorage.setItem('themeMode', themeMode.value)
  updateTheme()
}

// 更新主题
const updateTheme = () => {
  if (themeMode.value === 'light') {
    isDarkTheme.value = false
  } else if (themeMode.value === 'dark') {
    isDarkTheme.value = true
  } else {
    // 跟随系统
    isDarkTheme.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  localStorage.setItem('theme', isDarkTheme.value ? 'dark' : 'light')
  document.documentElement.classList.toggle('dark', isDarkTheme.value)
}

// 监听系统主题变化
const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
mediaQuery.addEventListener('change', () => {
  if (themeMode.value === 'system') {
    updateTheme()
  }
})

// 监听 localStorage 中主题变化
window.addEventListener('storage', (event) => {
  if (event.key === 'theme') {
    isDarkTheme.value = event.newValue === 'dark'
    // 当手动切换主题时，更新 themeMode 为对应模式
    themeMode.value = isDarkTheme.value ? 'dark' : 'light'
    localStorage.setItem('themeMode', themeMode.value)
  } else if (event.key === 'accentColor' && event.newValue) {
    document.documentElement.style.setProperty('--active-color', event.newValue)
    // 触发 themeOverrides 重新计算
    themeTrigger.value++
  } else if (event.key === 'fontFamily' && event.newValue) {
    document.documentElement.style.setProperty('--font-family', event.newValue)
  } else if (event.key === 'fontSize' && event.newValue) {
    document.documentElement.style.setProperty('--font-size', event.newValue)
  } else if (event.key === 'locale' && event.newValue) {
    locale.value = event.newValue
  }
})

// 获取系统状态信息
const fetchSystemStatus = async () => {
  try {
    const status = await getSystemStatus()
    systemStatus.value = status
  } catch (error) {
    console.error('Failed to fetch system status:', error)
  }
}

// 启动定时更新
const startStatusUpdates = () => {
  // 立即获取一次数据
  fetchSystemStatus()
  
  // 每秒更新一次
  updateTimer = window.setInterval(() => {
    fetchSystemStatus()
  }, 1000)
}

// 停止定时更新
const stopStatusUpdates = () => {
  if (updateTimer) {
    window.clearInterval(updateTimer)
    updateTimer = null
  }
}

onMounted(() => {
  updateTheme()
  
  // 初始化语言
  const savedLanguage = localStorage.getItem('locale')
  if (savedLanguage) {
    locale.value = savedLanguage
  }
  
  // 初始化强调色
  const savedAccentColor = localStorage.getItem('accentColor')
  if (savedAccentColor) {
    document.documentElement.style.setProperty('--active-color', savedAccentColor)
  }
  
  // 初始化字体
  const savedFontFamily = localStorage.getItem('fontFamily')
  if (savedFontFamily) {
    document.documentElement.style.setProperty('--font-family', savedFontFamily)
  }
  
  // 初始化字体大小
  const savedFontSize = localStorage.getItem('fontSize')
  if (savedFontSize) {
    document.documentElement.style.setProperty('--font-size', savedFontSize)
  }
  
  // 启动系统状态定时更新
  startStatusUpdates()
})

// 组件卸载时清理定时器
onUnmounted(() => {
  stopStatusUpdates()
})
</script>

<template>
  <NConfigProvider 
    :theme="isDarkTheme ? darkTheme : null"
    :theme-overrides="isDarkTheme ? darkThemeOverrides : themeOverrides"
  >
    <NMessageProvider>
      <NDialogProvider>
        <div class="app-container" :class="{ 'dark': isDarkTheme }">
          <TitleBar :is-dark-theme="isDarkTheme" />
          <div class="main-content">
            <!-- 左边导航区 -->
            <Sidebar 
              :current-nav-item="currentNavItem"
              @switch-nav="switchNavItem"
            />
            
            <!-- 右边内容区 -->
            <div class="content-area">
              <!-- 首页内容 -->
              <HomeContent 
                v-if="currentContent === 'home'"
                :system-status="systemStatus"
              />
              
              <!-- 项目内容 -->
              <ProjectsContent 
                v-else-if="currentContent === 'projects'"
              />
              
              <!-- 载荷内容 -->
              <PayloadsContent 
                v-else-if="currentContent === 'payloads'"
              />
              
              <!-- 插件内容 -->
              <PluginsContent 
                v-else-if="currentContent === 'plugins'"
              />
              
              <!-- 设置内容 -->
              <SettingsContent 
                v-else-if="currentContent === 'settings'"
                :is-dark-theme="isDarkTheme"
                :theme-mode="themeMode"
                @update:theme-mode="handleThemeChange"
              />
              
              <!-- 默认内容（当没有匹配任何条件时显示） -->
              <div v-else class="default-content">
                <div class="default-content-inner">
                  <h2>{{ t('common.loading') }}</h2>
                  <p>正在加载内容...</p>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 底部状态栏 -->
          <StatusBar :system-status="systemStatus" />
        </div>
      </NDialogProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style>
/* 全局样式 - 现代极客风格 */
:root {
  --bg-color: #ffffff;
  --text-color: #1f2937;
  --border-color: #e5e7eb;
  --hover-color: #f3f4f6;
  --active-color: #3b82f6;
  --sidebar-bg: #ffffff;
  --content-bg: #f9fafb;
  --status-bar-bg: #ffffff;
  --card-bg: #ffffff;
  --panel-bg: #3b82f6;
  --success-color: #10b981;
  --warning-color: #f59e0b;
  --error-color: #ef4444;
  --info-color: #3b82f6;
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  --font-size: 14px;
  --border-radius-sm: 4px;
  --border-radius-md: 6px;
  --border-radius-lg: 8px;
  --border-radius-xl: 12px;
  --transition-fast: 150ms ease-in-out;
  --transition-normal: 200ms ease-in-out;
  --transition-slow: 300ms ease-in-out;
}

.dark {
  --bg-color: #0f172a;
  --text-color: #f1f5f9;
  --border-color: #334155;
  --hover-color: #1e293b;
  --sidebar-bg: #1e293b;
  --content-bg: #0f172a;
  --status-bar-bg: #1e293b;
  --card-bg: #1e293b;
  --panel-bg: #3b82f6;
  --success-color: #10b981;
  --warning-color: #f59e0b;
  --error-color: #ef4444;
  --info-color: #3b82f6;
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.4), 0 2px 4px -1px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.5), 0 4px 6px -2px rgba(0, 0, 0, 0.4);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.6), 0 10px 10px -5px rgba(0, 0, 0, 0.5);
}

body {
  margin: 0;
  padding: 0;
  font-family: var(--font-family);
  font-size: var(--font-size);
  line-height: 1.5;
  background: var(--bg-color);
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

* {
  font-size: inherit;
}

.dark body {
  background: var(--bg-color);
}

/* 应用容器样式 */
.app-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--bg-color);
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
}

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
  font-size: 0.857em; /* 12px / 14px */
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

/* 内容区域样式 - 深色主题风格 */
.content-area {
  flex: 1;
  background: var(--content-bg);
  padding: 0;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  min-width: 0;
  border-right: 1px solid var(--border-color);
}

.content-section {
  position: relative;
  z-index: 1;
}

.content-section h1 {
  font-size: 1.714em; /* 24px / 14px */
  font-weight: 600;
  margin-bottom: 24px;
  color: var(--text-color);
  transition: color 0.4s ease;
}

/* 首页卡片样式 */
.home-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 32px;
}

.home-card {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  padding: 24px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  z-index: 1;
}

.home-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.home-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
}

.home-card svg {
  margin-bottom: 16px;
  color: var(--active-color);
}

.home-card h3 {
  font-size: 1.286em; /* 18px / 14px */
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.home-card p {
  font-size: 1em; /* 14px / 14px */
  line-height: 1.5;
  color: var(--text-color);
  opacity: 0.8;
  position: relative;
  z-index: 1;
}

/* 项目管理样式 */
.projects-content {
  display: flex;
  gap: 20px;
  height: calc(100% - 40px);
}

.projects-sidebar {
  width: 250px;
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border-right: 1px solid var(--border-color);
  padding: 20px;
  overflow-y: auto;
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
}

.projects-sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.projects-sidebar h3 {
  font-size: 1.143em; /* 16px / 14px */
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.directory-tree {
  position: relative;
  z-index: 1;
}

.tree-item {
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 4px;
  color: var(--text-color);
}

.tree-item:hover {
  background: var(--hover-color);
}

.projects-main {
  flex: 1;
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  padding: 20px;
  overflow-y: auto;
  position: relative;
  overflow: hidden;
}

.projects-main::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.projects-main h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.webshell-list {
  position: relative;
  z-index: 1;
}

.webshell-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.2s ease;
  border: 1px solid var(--border-color);
}

.dark .webshell-item {
  background: rgba(255, 255, 255, 0.05);
}

.webshell-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
}

.webshell-item span {
  flex: 1;
  margin-right: 16px;
  font-size: 14px;
  color: var(--text-color);
}

.action-button {
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 0.857em; /* 12px / 14px */
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
  box-shadow: var(--shadow-sm);
}

.action-button:hover {
  background: var(--active-color);
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.action-button:active {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 32px;
  gap: 8px;
  position: relative;
  z-index: 1;
}

.page-button {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--text-color);
}

.page-button:hover {
  background: var(--hover-color);
  border-color: var(--active-color);
  color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.page-button.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  font-weight: 600;
  box-shadow: var(--shadow-md);
}

.page-button.active:hover {
  opacity: 0.9;
  transform: translateY(-2px);
}

/* 载荷生成样式 */
.payloads-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  position: relative;
  z-index: 1;
}

.payload-form {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  padding: 32px;
  width: 100%;
  max-width: 500px;
  position: relative;
  overflow: hidden;
}

.payload-form::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.form-group {
  margin-bottom: 24px;
  position: relative;
  z-index: 1;
}

.form-group label {
  display: block;
  font-size: 1em; /* 14px / 14px */
  font-weight: 500;
  margin-bottom: 8px;
  color: var(--text-color);
}

.form-group select,
.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.8);
  color: var(--text-color);
  font-size: 14px;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.dark .form-group select,
.dark .form-group input {
  background: rgba(255, 255, 255, 0.1);
}

.form-group select:hover,
.form-group input:hover {
  border-color: var(--active-color);
}

.form-group select:focus,
.form-group input:focus {
  outline: none;
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.generate-button {
  width: 100%;
  background: var(--active-color);
  color: white;
  border: none;
  padding: 14px 20px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  z-index: 1;
  border-radius: 8px;
  box-shadow: var(--shadow-sm);
  letter-spacing: 0.5px;
}

.generate-button:hover {
  background: var(--active-color);
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.generate-button:active {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

/* 插件管理样式 */
.plugins-content {
  position: relative;
  z-index: 1;
}

.plugins-tabs {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 12px;
}

.tab-button {
  background: none;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-color);
  font-size: 14px;
  font-weight: 500;
}

.tab-button:hover {
  background: var(--hover-color);
}

.tab-button.active {
  background: var(--active-color);
  color: white;
}

.plugins-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.plugin-item {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  padding: 20px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.plugin-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.plugin-item:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
}

.plugin-item h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.plugin-item p {
  font-size: 14px;
  line-height: 1.4;
  margin-bottom: 16px;
  color: var(--text-color);
  opacity: 0.8;
  position: relative;
  z-index: 1;
}

.plugin-button {
  background: var(--active-color);
  color: white;
  border: none;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  z-index: 1;
}

.plugin-button:hover {
  background: var(--active-color);
  opacity: 0.9;
  transform: translateY(-1px);
}

/* 设置样式 */
.settings-content {
  max-width: 600px;
  position: relative;
  z-index: 1;
}

.settings-group {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  padding: 24px;
  margin-bottom: 20px;
  position: relative;
  overflow: hidden;
}

.settings-group::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.settings-group h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
  position: relative;
  z-index: 1;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item label {
  font-size: 14px;
  color: var(--text-color);
}

.setting-item select {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 6px 12px;
  color: var(--text-color);
  font-size: 14px;
}

.setting-item input[type="checkbox"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.setting-item button {
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
}

.setting-item button:hover {
  background: var(--active-color);
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.setting-item button:active {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

.dark ::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.dark ::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 60px;
  }
  
  .nav-item {
    width: 50px;
    height: 50px;
    min-width: 50px;
    min-height: 50px;
    max-width: 50px;
    max-height: 50px;
    padding: 8px 4px;
  }
  
  .nav-item span {
    font-size: 10px;
  }
  
  .nav-item svg {
    width: 16px;
    height: 16px;
  }
  
  .main-content {
    padding: 8px;
    gap: 8px;
  }
  
  .content-area {
    padding: 16px;
  }
  
  .home-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .projects-content {
    flex-direction: column;
  }
  
  .projects-sidebar {
    width: 100%;
    max-height: 200px;
  }
  
  .status-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
    padding: 12px 20px;
  }
}
</style>

<style scoped>
.app-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--bg-color);
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
}

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
  font-size: 0.857em; /* 12px / 14px */
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

/* 内容区域样式 - 深色主题风格 */
.content-area {
  flex: 1;
  background: var(--content-bg);
  padding: 0;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  min-width: 0;
  border-right: 1px solid var(--border-color);
}

.content-section {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  width: 100%;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

.content-header {
  width: 100%;
  padding: 24px 24px 20px 24px;
  margin-bottom: 0;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  box-sizing: border-box;
}

.dark .content-header {
  border-bottom-color: var(--border-strong);
}

.content-section h1 {
  margin: 0;
  font-size: 0;
  line-height: 1;
  color: var(--text-color);
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
}

.content-body {
  flex: 1;
  width: 100%;
  padding: 0;
  margin: 0;
  box-sizing: border-box;
  background: var(--content-bg);
  border-top: none;
}

.content-section h1 .title {
  font-size: 24px;
  font-weight: 600;
  color: var(--active-color);
  letter-spacing: 0;
}

.content-section h1 .separator {
  color: var(--text-tertiary);
  font-weight: 300;
  font-size: 20px;
}

.content-section h1 .subtitle {
  font-size: 16px;
  font-weight: 400;
  color: var(--text-secondary);
  letter-spacing: 0;
}

/* 首页内容样式 - 深色主题风格 */
.home-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.home-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-md);
  padding: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  border-radius: 8px;
}

.home-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg);
}

.home-card svg {
  flex-shrink: 0;
  color: var(--active-color);
}

.home-card h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
}

.home-card p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.8;
}

/* 项目内容样式 - 深色主题风格 */
.projects-content {
  display: flex;
  gap: 1px;
  width: 100%;
  box-sizing: border-box;
  background: var(--border-color);
}

.projects-sidebar {
  width: 200px;
  background: var(--sidebar-bg);
  padding: 20px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow-y: auto;
}

.directory-tree {
  margin-top: 16px;
}

.tree-item {
  padding: 10px 12px;
  cursor: pointer;
  margin-bottom: 4px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  border-radius: 4px;
}

.tree-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
}

.tree-item.active {
  background: var(--panel-bg);
  color: var(--active-color);
  box-shadow: var(--shadow-sm);
}

.projects-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--content-bg);
  padding: 20px;
  overflow-y: auto;
}

.projects-main h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
}

.webshell-table-card {
  margin-bottom: 16px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding-top: 16px;
}

/* 载荷内容样式 - 深色主题风格 */
.payloads-content {
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-x: hidden;
}

.payload-form {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-md);
  padding: 28px;
  max-width: 560px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  border-radius: 8px;
}

.form-group {
  margin-bottom: 20px;
  width: 100%;
  box-sizing: border-box;
}

.form-group label {
  display: block;
  margin-bottom: 10px;
  font-weight: 600;
  font-size: 14px;
  color: var(--text-color);
  width: 100%;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  gap: 8px;
}

.form-group label::before {
  content: '';
  width: 4px;
  height: 16px;
  background: var(--active-color);
  border-radius: 2px;
  flex-shrink: 0;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--content-bg);
  color: var(--text-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input::placeholder {
  color: var(--text-color);
  opacity: 0.5;
}

.form-group select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23666' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  background-size: 16px;
  padding-right: 40px;
}

.form-group select:focus {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23667eea' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  border-color: var(--active-color);
}

.generate-button {
  padding: 16px 36px;
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 140px;
  max-width: 200px;
  margin: 0 auto;
  box-sizing: border-box;
  letter-spacing: 0.5px;
}

.generate-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  background: var(--active-color);
  opacity: 0.9;
}

.generate-button:active {
  transform: translateY(0);
}

/* 插件内容样式 - 深色主题风格 */
.plugins-content {
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.plugins-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 12px;
}

.tab-button {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  background: var(--card-bg);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
  font-size: 14px;
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 90px;
}

.tab-button:hover {
  background: var(--hover-color);
  transform: translateY(-2px);
}

.tab-button.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  box-shadow: var(--shadow-md);
}

.plugins-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
  box-sizing: border-box;
}

.plugin-item {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-md);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 8px;
}

.plugin-item:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.plugin-item h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
}

.plugin-item p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.8;
}

.plugin-button {
  margin-top: 12px;
  padding: 12px 24px;
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  width: fit-content;
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 90px;
}

.plugin-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  background: var(--active-color);
  opacity: 0.9;
}

/* 设置内容样式 - 深色主题风格 */
.settings-layout {
  display: flex;
  width: 100%;
  gap: 1px;
  background: var(--border-color);
  margin: 0;
  padding: 0;
  border-top: none;
}

.settings-sidebar {
  width: 200px;
  background: var(--sidebar-bg);
  padding: 20px 0;
  overflow-y: auto;
  margin: 0;
  border: none;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  z-index: 10;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .settings-sidebar {
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.4);
}

.settings-nav-item {
  padding: 12px 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  border-left: 3px solid transparent;
}

.settings-nav-item:hover {
  background: var(--hover-color);
  border-left-color: var(--active-color);
}

.settings-nav-item.active {
  background: var(--panel-bg);
  color: white;
  border-left-color: var(--active-color);
}

.settings-main {
  flex: 1;
  background: var(--content-bg);
  padding: 0;
  margin: 0;
  overflow-y: auto;
  border: none;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.05);
  z-index: 5;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .settings-main {
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.4);
}

.settings-panel {
  width: 100%;
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;
  border: none;
}



.settings-group {
  padding: 20px;
  margin: 0;
}

.settings-group h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  opacity: 0.8;
}

/* 主题选项 */
.theme-options {
  display: flex;
  gap: 16px;
  width: 100%;
  margin: 0;
  padding: 0;
}

.theme-option {
  flex: 1;
  min-width: 0;
  padding: 24px 20px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: var(--text-color);
  text-align: center;
  position: relative;
  min-height: 120px;
}

.theme-option:hover {
  border-color: var(--active-color);
}

.theme-option.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.theme-icon {
  font-size: 28px;
  display: block;
  margin: 0 auto;
}

.theme-option span:nth-child(2) {
  font-size: 14px;
  font-weight: 500;
  display: block;
  margin: 0 auto;
}

.theme-check {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 16px;
  height: 16px;
  background: white;
  color: var(--active-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
}

.dark .theme-check {
  background: var(--active-color);
  color: white;
}

@media (max-width: 1200px) {
  .home-content {
    max-width: 100%;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  }
  
  .projects-content {
    max-width: 100%;
  }
  
  .plugins-content {
    max-width: 100%;
  }
  
  .payloads-content {
    max-width: 100%;
  }
  
  .settings-content {
    max-width: 100%;
  }
}

@media (max-width: 768px) {
  .main-content {
    gap: 0;
  }
  
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
  
  .content-area {
    padding: 0;
  }
  
  .content-header {
    padding: 12px 16px;
  }
  
  .content-body {
    padding: 16px;
  }
  
  .home-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .projects-content {
    flex-direction: column;
  }
  
  .projects-sidebar {
    width: 100%;
    height: auto;
    max-height: 150px;
  }
}
</style>
