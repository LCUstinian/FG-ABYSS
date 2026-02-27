<template>
  <div class="title-bar" :class="{ 'dark': isDarkTheme }">
    <!-- 左侧：应用图标和名称 -->
    <div class="title-bar-left">
      <span class="app-name">{{ appName }}</span>
    </div>

    <!-- 中间：主题和语言切换 -->
    <Tooltip :text="isDarkTheme ? '切换到浅色主题' : '切换到深色主题'">
      <button @click="toggleTheme" class="control-button theme-button">
        <Sun v-if="isDarkTheme" :size="20" />
        <Moon v-else :size="20" />
      </button>
    </Tooltip>
    
    <!-- 间距元素 -->
    <div class="spacer"></div>
    
    <!-- 语言按钮和窗口控制按钮区域 -->
    <div class="button-group">
      <Tooltip :text="locale === 'zh-CN' ? '切换到英文' : '切换到中文'">
        <button @click="toggleLanguage" class="control-button language-button">
          <span class="language-icon" :class="locale === 'zh-CN' ? 'china-flag' : 'us-flag'"></span>
        </button>
      </Tooltip>
      <div class="divider"></div>
      <Tooltip :text="t('status.minimize')">
        <button @click="minimizeWindow" class="window-control minimize">
          <Minus :size="20" />
        </button>
      </Tooltip>
    </div>
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
.title-bar {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 4px;
  height: 48px;
  padding: 0 20px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 
    0 2px 10px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  color: #1a1a1a;
  user-select: none;
  cursor: default;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 0;
  border-radius: 12px 12px 0 0;
  position: relative;
}

.title-bar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.4) 0%, rgba(255, 255, 255, 0.1) 100%);
  pointer-events: none;
  z-index: 0;
}

.title-bar.dark {
  background: rgba(20, 20, 20, 0.8);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow: 
    0 2px 10px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  color: #e0e0e0;
}

.title-bar.dark::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
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
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.app-name {
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.5px;
}

.title-bar-center {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.title-bar-right {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.spacer {
  width: 12px;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.button-group {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.button-group .control-button,
.button-group .window-control {
  margin: 0;
}

.button-group .divider {
  width: 1px;
  height: 24px;
  background: rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
  margin: 0 12px;
}

.title-bar.dark .button-group .divider {
  background: rgba(255, 255, 255, 0.15);
}

.control-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.25);
  border-radius: 10px;
  cursor: pointer;
  font-size: 18px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  position: relative;
  overflow: hidden;
}

.control-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.3) 0%, rgba(255, 255, 255, 0.1) 100%);
  pointer-events: none;
  z-index: 0;
}

.control-button:hover {
  background: var(--active-color);
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: 
    0 6px 16px rgba(0, 0, 0, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  border-color: var(--active-color);
  color: white;
}

.title-bar.dark .control-button {
  border: 1px solid rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.title-bar.dark .control-button::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.15) 0%, rgba(255, 255, 255, 0.05) 100%);
}

.title-bar.dark .control-button:hover {
  background: var(--active-color);
  opacity: 0.9;
  box-shadow: 
    0 6px 16px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  border-color: var(--active-color);
  color: white;
}

/* 主题切换按钮 */
.theme-button {
  color: #ffb300; /* 黄色，代表太阳 */
}

.theme-button svg {
  position: relative;
  z-index: 1;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.title-bar.dark .theme-button {
  color: #ffd700; /* 金色，在深色模式下更明显 */
}

.title-bar.dark .theme-button svg {
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

/* 语言切换按钮 */
.language-button {
  color: #1e88e5; /* 蓝色，代表国际化 */
}

.language-button svg {
  position: relative;
  z-index: 1;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.title-bar.dark .language-button {
  color: #64b5f6; /* 浅蓝色，在深色模式下更明显 */
}

.title-bar.dark .language-button svg {
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

.language-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 21px;
  position: relative;
  z-index: 1;
  flex-shrink: 0;
  margin: 0 auto;
}

/* 中国国旗 - 使用 emoji */
.china-flag::before {
  content: '🇨🇳';
  font-size: 18px;
}

/* 美国国旗 - 使用 emoji */
.us-flag::before {
  content: '🇺🇸';
  font-size: 18px;
}

.window-control {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.4);
  position: relative;
  overflow: hidden;
}

.window-control::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.3) 0%, rgba(255, 255, 255, 0.1) 100%);
  pointer-events: none;
  z-index: 0;
}

.window-control:hover {
  background: var(--active-color);
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: 
    0 6px 16px rgba(0, 0, 0, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  border-color: var(--active-color);
  color: white;
}

.title-bar.dark .window-control {
  border: 1px solid rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.title-bar.dark .window-control::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
}

.title-bar.dark .window-control:hover {
  background: var(--active-color);
  opacity: 0.9;
  box-shadow: 
    0 6px 16px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  border-color: var(--active-color);
  color: white;
}

/* 最小化按钮 */
.window-control.minimize {
  color: #757575; /* 灰色 */
}

.title-bar.dark .window-control.minimize {
  color: #bdbdbd; /* 浅灰色 */
}

/* 最大化按钮 */
.window-control.maximize {
  color: #4caf50; /* 绿色 */
}

.title-bar.dark .window-control.maximize {
  color: #81c784; /* 浅绿色 */
}

/* 关闭按钮 */
.window-control.close {
  color: #f44336; /* 红色 */
}

.title-bar.dark .window-control.close {
  color: #e57373; /* 浅红色 */
}

.window-control.close:hover {
  background-color: #ff5f57;
  color: white;
}

.title-bar.dark .window-control.close:hover {
  background-color: #ff5f57;
  color: white;
}

/* 窗口控制按钮图标 */
.window-control svg {
  position: relative;
  z-index: 1;
}

/* 图标切换动画 */
.window-control svg {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 关闭按钮特殊悬停效果 */
.window-control.close:hover {
  background: #ff5f57;
  color: white;
  border-color: #ff5f57;
  box-shadow: 
    0 4px 12px rgba(255, 95, 87, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.window-control.close:hover::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.3) 0%, rgba(255, 255, 255, 0.1) 100%);
}

.title-bar.dark .window-control.close:hover {
  background: #ff5f57;
  color: white;
  border-color: #ff5f57;
  box-shadow: 
    0 4px 12px rgba(255, 95, 87, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.title-bar.dark .window-control.close:hover::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.3) 0%, rgba(255, 255, 255, 0.1) 100%);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .title-bar {
    padding: 0 8px;
  }
  
  .app-name {
    font-size: 11px;
  }
  
  .control-button,
  .window-control {
    width: 20px;
    height: 20px;
  }
  
  .title-bar-center {
    gap: 4px;
  }
}

/* 小窗口适配 */
@media (max-width: 480px) {
  .app-name {
    font-size: 10px;
  }
  
  .control-button,
  .window-control {
    width: 18px;
    height: 18px;
  }
  
  .title-bar {
    padding: 0 6px;
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
