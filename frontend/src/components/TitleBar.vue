<template>
  <div class="title-bar" :class="{ 'dark': isDarkTheme }">
    <!-- 左侧：应用图标和名称 -->
    <div class="title-bar-left">
      <span class="app-name">{{ appName }}</span>
    </div>

    <!-- 间距元素 -->
    <div class="spacer"></div>
    
    <!-- 语言按钮和窗口控制按钮区域 -->
    <div class="button-group">
      <!-- 主题切换按钮 -->
      <Tooltip :text="isDarkTheme ? '切换到浅色主题' : '切换到深色主题'">
        <button @click="toggleTheme" class="control-button theme-button">
          <Sun v-if="isDarkTheme" :size="20" />
          <Moon v-else :size="20" />
        </button>
      </Tooltip>
      
      <!-- 语言切换按钮 -->
      <Tooltip :text="locale === 'zh-CN' ? '切换到英文' : '切换到中文'">
        <button @click="toggleLanguage" class="control-button language-button">
          <span class="language-icon" :class="locale === 'zh-CN' ? 'china-flag' : 'us-flag'"></span>
        </button>
      </Tooltip>
      
      <div class="divider"></div>
      
      <!-- 窗口控制按钮 -->
      <Tooltip :text="t('status.minimize')">
        <button @click="minimizeWindow" class="window-control minimize">
          <Minus :size="20" />
        </button>
      </Tooltip>
      <Tooltip :text="isMaximized ? t('status.restore') : t('status.maximize')">
        <button @click="toggleMaximize" class="window-control maximize">
          <Maximize2 v-if="!isMaximized" :size="20" />
          <Minimize2 v-else :size="20" />
        </button>
      </Tooltip>
      <Tooltip :text="t('status.close')">
        <button @click="closeWindow" class="window-control close">
          <X :size="20" />
        </button>
      </Tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Minus, Maximize2, Minimize2, X, Sun, Moon } from 'lucide-vue-next'
import * as wails from '@wailsio/runtime'
import Tooltip from './Tooltip.vue'

const { t, locale } = useI18n()

// 接收主题状态作为prop
const props = defineProps({
  isDarkTheme: {
    type: Boolean,
    default: false
  }
})

// 窗口状态
const isMaximized = ref(false)

// 应用名称
const appName = ref('FG-ABYSS 非攻-渊渟')

// 主题切换
const toggleTheme = () => {
  const newTheme = !props.isDarkTheme
  localStorage.setItem('theme', newTheme ? 'dark' : 'light')
  localStorage.setItem('themeMode', newTheme ? 'dark' : 'light')
  document.documentElement.classList.toggle('dark', newTheme)
  // 触发storage事件，让其他组件知道主题变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'theme',
    newValue: newTheme ? 'dark' : 'light'
  }))
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'themeMode',
    newValue: newTheme ? 'dark' : 'light'
  }))
}

// 监听localStorage中主题变化
window.addEventListener('storage', (event) => {
  if (event.key === 'theme') {
    // 由于我们使用的是prop，这里不需要直接修改isDarkTheme
    // 而是通过App.vue中的watch来更新prop值
  }
})

// 语言切换
const toggleLanguage = () => {
  const currentLang = locale.value
  const newLang = currentLang === 'zh-CN' ? 'en-US' : 'zh-CN'
  locale.value = newLang
  localStorage.setItem('locale', newLang)
  // 触发storage事件，让其他组件知道语言变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'locale',
    newValue: newLang
  }))
}

// 窗口控制
const minimizeWindow = async () => {
  try {
    await wails.Window.Minimise()
  } catch (error) {
    console.error('Error minimizing window:', error)
  }
}

const toggleMaximize = async () => {
  try {
    if (isMaximized.value) {
      await wails.Window.UnMaximise()
      isMaximized.value = false
    } else {
      await wails.Window.Maximise()
      isMaximized.value = true
    }
  } catch (error) {
    console.error('Error toggling maximize:', error)
  }
}

const closeWindow = async () => {
  try {
    await wails.Window.Close()
  } catch (error) {
    console.error('Error closing window:', error)
  }
}

// 窗口拖拽功能通过 CSS 属性 --wails-draggable: drag 实现
// 无需显式调用方法，Wails 运行时会自动处理

// 检查窗口最大化状态
const checkMaximizeState = async () => {
  try {
    isMaximized.value = await wails.Window.IsMaximised()
  } catch (error) {
    console.error('Error checking maximize state:', error)
  }
}

onMounted(() => {
  // 检查窗口状态
  checkMaximizeState()
  window.addEventListener('resize', checkMaximizeState)
})
</script>

<style scoped>
/* 标题栏样式 - 现代风格 */
.title-bar {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
  height: 52px;
  padding: 0 24px;
  background: var(--status-bar-bg);
  border-bottom: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  color: var(--text-color);
  user-select: none;
  cursor: default;
  transition: all var(--transition-normal);
  min-width: 0;
  position: relative;
}

.title-bar-left {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
  cursor: grab;
  --wails-draggable: drag;
  position: relative;
  z-index: 1;
}

.title-bar-left:active {
  cursor: grabbing;
}

.app-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  margin-right: 12px;
}

.app-name {
  font-size: 16px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.5px;
  font-family: var(--font-family);
}

.spacer {
  flex: 1;
  min-width: 0;
  position: relative;
  z-index: 1;
}

.button-group {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
  gap: 6px;
}

.button-group .divider {
  width: 1px;
  height: 28px;
  background: var(--border-color);
  flex-shrink: 0;
  margin: 0 8px;
}

/* 控制按钮通用样式 */
.control-button,
.window-control {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid var(--border-color);
  background: var(--card-bg);
  border-radius: var(--border-radius-md);
  cursor: pointer;
  font-size: 16px;
  transition: all var(--transition-normal);
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
  color: var(--text-color);
}

.control-button::before,
.window-control::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left var(--transition-slow);
}

.control-button:hover::before,
.window-control:hover::before {
  left: 100%;
}

.control-button:hover,
.window-control:hover {
  background: var(--hover-color);
  border-color: var(--active-color);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.control-button svg,
.window-control svg {
  position: relative;
  z-index: 1;
  transition: transform var(--transition-normal);
}

.control-button:hover svg,
.window-control:hover svg {
  transform: scale(1.1);
}

/* 主题切换按钮 */
.theme-button {
  color: var(--warning-color);
}

/* 语言切换按钮 */
.language-button {
  color: var(--info-color);
}

.language-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 15px;
  position: relative;
  z-index: 1;
  flex-shrink: 0;
  margin: 0 auto;
}

/* 中国国旗 - 使用 emoji */
.china-flag::before {
  content: '🇨🇳';
  font-size: 16px;
}

/* 美国国旗 - 使用 emoji */
.us-flag::before {
  content: '🇺🇸';
  font-size: 16px;
}

/* 窗口控制按钮 */
.window-control.minimize {
  color: var(--text-color);
  opacity: 0.7;
}

.window-control.maximize {
  color: var(--text-color);
  opacity: 0.7;
}

.window-control.close {
  color: var(--error-color);
}

.window-control.close:hover {
  background: var(--error-color);
  color: white;
  border-color: var(--error-color);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.window-control.close:hover::before {
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .title-bar {
    padding: 0 16px;
    height: 48px;
  }
  
  .app-name {
    font-size: 14px;
  }
  
  .control-button,
  .window-control {
    width: 32px;
    height: 32px;
  }
  
  .button-group {
    gap: 6px;
  }
  
  .button-group .divider {
    height: 24px;
    margin: 0 6px;
  }
}

/* 小窗口适配 */
@media (max-width: 480px) {
  .app-name {
    font-size: 12px;
  }
  
  .control-button,
  .window-control {
    width: 28px;
    height: 28px;
  }
  
  .title-bar {
    padding: 0 12px;
  }
  
  .button-group {
    gap: 4px;
  }
}

/* Tooltip 组件样式调整 */
.title-bar :deep(.tooltip-container) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.title-bar :deep(.tooltip-container) > * {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
