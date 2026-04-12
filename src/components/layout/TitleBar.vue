<template>
  <div class="title-bar" :class="{ 'dark': isDarkTheme }" data-tauri-drag-region>
    <!-- 左侧：应用名称（拖拽区域） -->
    <div class="title-bar-left" data-tauri-drag-region>
      <span class="app-name">{{ appName }}</span>
    </div>

    <!-- 右侧：控制按钮区域 -->
    <div class="title-bar-right">
      <!-- 功能按钮组 -->
      <div class="button-group">
        <!-- 主题切换按钮 -->
        <Tooltip :text="isDarkTheme ? 'Switch to light theme' : 'Switch to dark theme'">
          <button @click="toggleTheme" class="control-button theme-button" type="button" aria-label="Toggle theme">
            <div class="button-icon">
              <Sun v-if="isDarkTheme" :size="18" :stroke-width="2" />
              <Moon v-else :size="18" :stroke-width="2" />
            </div>
          </button>
        </Tooltip>
        
        <!-- 语言切换按钮 -->
        <Tooltip :text="locale === 'zh-CN' ? 'Switch to English' : 'Switch to Chinese'">
          <button @click="toggleLanguage" class="control-button language-button" type="button" aria-label="Toggle language">
            <div class="button-icon">
              <span class="language-icon">{{ locale === 'zh-CN' ? '🇨🇳' : '🇺🇸' }}</span>
            </div>
          </button>
        </Tooltip>
      </div>
      
      <!-- 分隔线 -->
      <div class="divider"></div>
      
      <!-- 窗口控制按钮组 -->
      <div class="button-group window-controls">
        <!-- 最小化按钮 -->
        <Tooltip :text="t('status.minimize')">
          <button @click="minimizeWindow" class="window-control minimize" type="button" :aria-label="t('window.minimize')">
            <div class="button-icon">
              <Minus :size="18" :stroke-width="2" />
            </div>
          </button>
        </Tooltip>
        
        <!-- 最大化/还原按钮 -->
        <Tooltip :text="isMaximized ? t('status.restore') : t('status.maximize')">
          <button @click="handleToggleMaximize" class="window-control maximize" type="button" :aria-label="isMaximized ? t('status.restore') : t('status.maximize')">
            <div class="button-icon">
              <Maximize2 v-if="!isMaximized" :size="18" :stroke-width="2" />
              <Minimize2 v-else :size="18" :stroke-width="2" />
            </div>
          </button>
        </Tooltip>
        
        <!-- 关闭按钮 -->
        <Tooltip :text="t('status.close')">
          <button @click="handleClose" class="window-control close" type="button" :aria-label="t('window.close')">
            <div class="button-icon">
              <X :size="18" :stroke-width="2" />
            </div>
          </button>
        </Tooltip>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Minus, Maximize2, Minimize2, X, Sun, Moon } from 'lucide-vue-next'
import Tooltip from '@/components/shared/Tooltip.vue'
import { useWindowControl } from '@/composables'
import { useMessage } from 'naive-ui'

const { t, locale } = useI18n()
const message = useMessage()

// 接收主题状态作为 prop
const props = defineProps({
  isDarkTheme: {
    type: Boolean,
    default: false
  }
})

// 使用窗口控制 composable
const { isMaximized, minimizeWindow, toggleMaximize, closeWindow, checkMaximizeState } = useWindowControl()

// 应用名称
const appName = ref('FG-ABYSS')

// 窗口状态更新定时器
let resizeCheckTimer: number | null = null

/**
 * 处理最大化/还原切换
 * 添加错误处理和用户反馈
 */
const handleToggleMaximize = async () => {
  try {
    await toggleMaximize()
    // 延迟检查状态，确保状态已更新
    setTimeout(() => {
      checkMaximizeState()
    }, 100)
  } catch (error) {
    console.error('切换最大化状态失败:', error)
    message.error(t('window.windowOperationFailed'))
  }
}

/**
 * 处理关闭窗口
 * 添加确认提示（可选）
 */
const handleClose = async () => {
  try {
    await closeWindow()
  } catch (error) {
    console.error('关闭窗口失败:', error)
    message.error(t('window.closeWindowFailed'))
  }
}

/**
 * 主题切换
 */
const toggleTheme = () => {
  const newTheme = !props.isDarkTheme
  localStorage.setItem('theme', newTheme ? 'dark' : 'light')
  localStorage.setItem('themeMode', newTheme ? 'dark' : 'light')
  document.documentElement.classList.toggle('dark', newTheme)
  
  // 触发 storage 事件，让其他组件知道主题变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'theme',
    newValue: newTheme ? 'dark' : 'light'
  }))
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'themeMode',
    newValue: newTheme ? 'dark' : 'light'
  }))
}

/**
 * 语言切换
 */
const toggleLanguage = () => {
  const currentLang = locale.value
  const newLang = currentLang === 'zh-CN' ? 'en-US' : 'zh-CN'
  locale.value = newLang
  localStorage.setItem('locale', newLang)
  
  // 触发 storage 事件，让其他组件知道语言变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'locale',
    newValue: newLang
  }))
}

/**
 * 监听窗口 resize 事件，更新最大化状态
 */
const handleResize = () => {
  if (resizeCheckTimer) {
    clearTimeout(resizeCheckTimer)
  }
  
  resizeCheckTimer = window.setTimeout(() => {
    checkMaximizeState()
  }, 100)
}

// 组件挂载时初始化
onMounted(() => {
  // 检查初始最大化状态
  checkMaximizeState()
  
  // 监听窗口 resize 事件
  window.addEventListener('resize', handleResize)
})

// 组件卸载时清理
onUnmounted(() => {
  if (resizeCheckTimer) {
    clearTimeout(resizeCheckTimer)
  }
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
/* ============================================
   标题栏样式 - 现代简约风格
   ============================================ */

.title-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 48px;
  padding: 0 12px;
  background: var(--title-bar-bg, var(--bg-color));
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  color: var(--text-color);
  user-select: none;
  cursor: default;
  transition: all var(--transition-normal);
  position: relative;
  z-index: 1000;
}

/* ============================================
   左侧拖拽区域
   ============================================ */

.title-bar-left {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
  cursor: grab;
  -webkit-app-region: drag;
  -webkit-user-select: none;
  padding-left: 4px;
}

.title-bar-left:active {
  cursor: grabbing;
}

.app-name {
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.3px;
  color: var(--text-color);
  opacity: 0.95;
  transition: opacity var(--transition-normal);
}

/* ============================================
   右侧按钮区域 - 精确对齐
   ============================================ */

.title-bar-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

/* 按钮组容器 - 基线对齐 */
.button-group {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 窗口控制按钮组 */
.button-group.window-controls {
  gap: 4px;
  margin-left: 4px;
}

/* ============================================
   分隔线 - 视觉优化
   ============================================ */

.divider {
  width: 1px;
  height: 28px;
  background: var(--border-color);
  margin: 0 8px;
  flex-shrink: 0;
  transition: all var(--transition-normal);
  opacity: 0.7;
}

/* ============================================
   按钮通用样式 - 精确居中对齐
   ============================================ */

.control-button,
.window-control {
  /* 尺寸规格 - 统一 40x40px */
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  
  /* 视觉样式 - 移除所有边框和阴影 */
  border: none !important;
  background: transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  color: inherit;
  padding: 0;
  outline: none;
  -webkit-tap-highlight-color: transparent;
  box-shadow: none !important;
  
  /* 过渡动画 */
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  
  /* 图标统一规格 */
  font-size: 18px;
  font-weight: 500;
}

/* 按钮图标容器 - 三重居中保障 */
.button-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  /* 确保图标视觉居中 */
  line-height: 1;
  /* 消除所有可能的偏移 */
  transform: translateZ(0);
}

/* 图标统一视觉权重 - 确保所有图标尺寸一致 */
.control-button svg,
.window-control svg {
  stroke-width: 2px;
  width: 18px;
  height: 18px;
  display: block;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  /* 确保图标在容器中完美居中 */
  margin: auto;
}

/* 按钮悬停效果 */
.control-button:hover,
.window-control:hover {
  background: var(--hover-color);
}

/* 按钮按下效果 - 流畅缩放 */
.control-button:active,
.window-control:active {
  transform: scale(0.92);
  background: var(--active-color-suppl);
}

/* 按钮禁用状态 */
.control-button:disabled,
.window-control:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: transparent !important;
  transform: none !important;
}

/* 图标悬停效果 - 微妙缩放 */
.control-button:hover svg,
.window-control:hover svg {
  transform: scale(1.08);
}

/* 图标按下效果 */
.control-button:active svg,
.window-control:active svg {
  transform: scale(0.96);
}

/* ============================================
   主题切换按钮 - 颜色优化
   ============================================ */

.theme-button {
  color: var(--warning-600);
}

.theme-button svg {
  stroke-width: 2px;
}

.theme-button:hover {
  color: var(--warning-700);
  background: rgba(245, 158, 11, 0.1);
}

.theme-button:active {
  background: rgba(245, 158, 11, 0.15);
}

/* ============================================
   语言切换按钮 - 精确居中对齐
   ============================================ */

.language-button {
  color: var(--info-600);
}

.language-button .button-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

/* 语言图标 - 使用 flex 居中，避免 emoji 基线问题 */
.language-icon {
  font-size: 18px;
  line-height: 1;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  text-align: center;
  
  /* 消除所有间距 */
  margin: 0 !important;
  padding: 0 !important;
  
  /* 确保垂直居中 */
  height: 100%;
  width: 100%;
  
  /* 平滑过渡 */
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.language-button:hover {
  color: var(--info-700);
  background: rgba(59, 130, 246, 0.1);
}

.language-button:hover .language-icon {
  transform: scale(1.08);
}

.language-button:active {
  background: rgba(59, 130, 246, 0.15);
}

.language-button:active .language-icon {
  transform: scale(0.96);
}

/* ============================================
   窗口控制按钮 - 统一对齐
   ============================================ */

.window-control.minimize {
  color: var(--text-color);
  opacity: 0.85;
}

.window-control.minimize svg {
  stroke-width: 2px;
}

.window-control.minimize:hover {
  opacity: 1;
  background: rgba(125, 125, 125, 0.08);
}

.window-control.minimize:active {
  background: rgba(125, 125, 125, 0.12);
}

.window-control.maximize {
  color: var(--text-color);
  opacity: 0.85;
}

.window-control.maximize svg {
  stroke-width: 2px;
  /* 微调尺寸确保视觉一致 */
  width: 17px;
  height: 17px;
}

.window-control.maximize:hover {
  opacity: 1;
  background: rgba(125, 125, 125, 0.08);
}

.window-control.maximize:active {
  background: rgba(125, 125, 125, 0.12);
}

.window-control.close {
  color: var(--error-600);
  opacity: 0.9;
}

.window-control.close svg {
  stroke-width: 2px;
}

.window-control.close:hover {
  background: rgba(239, 68, 68, 0.1);
  color: var(--error-700);
  opacity: 1;
}

.window-control.close:active {
  background: rgba(239, 68, 68, 0.15);
}

/* ============================================
   深色模式适配 - 完全一致的对齐和样式
   ============================================ */

.title-bar.dark {
  background: var(--title-bar-bg, var(--bg-color));
  border-bottom-color: var(--border-color);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.4);
}

.title-bar.dark .app-name {
  color: var(--text-color);
  opacity: 0.95;
}

.title-bar.dark .divider {
  background: var(--border-color);
  opacity: 0.4;
}

/* 深色模式按钮效果 - 保持相同的对齐和间距 */
.title-bar.dark .control-button:hover,
.title-bar.dark .window-control:hover {
  background: rgba(255, 255, 255, 0.08);
}

.title-bar.dark .control-button:active,
.title-bar.dark .window-control:active {
  background: rgba(255, 255, 255, 0.12);
}

/* 深色模式主题按钮 */
.title-bar.dark .theme-button {
  color: var(--warning-400);
}

.title-bar.dark .theme-button:hover {
  color: var(--warning-300);
  background: rgba(245, 158, 11, 0.12);
}

.title-bar.dark .theme-button:active {
  background: rgba(245, 158, 11, 0.18);
}

/* 深色模式语言按钮 */
.title-bar.dark .language-button {
  color: var(--info-400);
}

.title-bar.dark .language-button:hover {
  color: var(--info-300);
  background: rgba(59, 130, 246, 0.12);
}

.title-bar.dark .language-button:active {
  background: rgba(59, 130, 246, 0.18);
}

/* 深色模式窗口按钮 */
.title-bar.dark .window-control.minimize,
.title-bar.dark .window-control.maximize {
  color: var(--text-color);
}

.title-bar.dark .window-control.close {
  color: var(--error-400);
}

.title-bar.dark .window-control.close:hover {
  color: var(--error-300);
  background: rgba(239, 68, 68, 0.12);
}

.title-bar.dark .window-control.close:active {
  background: rgba(239, 68, 68, 0.18);
}

/* ============================================
   响应式设计 - 保持对齐和样式一致性
   ============================================ */

/* 平板尺寸 (≤768px) */
@media (max-width: 768px) {
  .title-bar {
    height: 44px;
    padding: 0 12px;
  }
  
  .app-name {
    font-size: 13px;
  }
  
  .title-bar-right {
    gap: 10px;
  }
  
  .button-group {
    gap: 2px;
  }
  
  .button-group.window-controls {
    margin-left: 2px;
  }
  
  /* 按钮尺寸适度缩小 */
  .control-button,
  .window-control {
    width: 38px;
    height: 38px;
  }
  
  /* 图标尺寸保持一致 */
  .control-button svg,
  .window-control svg {
    width: 17px;
    height: 17px;
  }
  
  .language-icon {
    font-size: 19px;
  }
  
  .divider {
    height: 24px;
    margin: 0 6px;
  }
}

/* 手机尺寸 (≤480px) */
@media (max-width: 480px) {
  .title-bar {
    height: 40px;
    padding: 0 10px;
  }
  
  .app-name {
    font-size: 12px;
  }
  
  .title-bar-right {
    gap: 8px;
  }
  
  .button-group {
    gap: 1px;
  }
  
  /* 按钮尺寸进一步缩小 */
  .control-button,
  .window-control {
    width: 34px;
    height: 34px;
  }
  
  /* 图标尺寸适度缩小 */
  .control-button svg,
  .window-control svg {
    width: 16px;
    height: 16px;
    stroke-width: 1.8px;
  }
  
  .language-icon {
    font-size: 18px;
  }
  
  .divider {
    height: 20px;
    margin: 0 4px;
  }
}

/* ============================================
   Tooltip 容器样式
   ============================================ */

.title-bar :deep(.tooltip-container) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* ============================================
   无障碍优化
   ============================================ */

/* 键盘焦点样式 */
.control-button:focus-visible,
.window-control:focus-visible {
  outline: 2px solid var(--active-color);
  outline-offset: 2px;
}

/* 减少动画 */
@media (prefers-reduced-motion: reduce) {
  .control-button,
  .window-control,
  .control-button svg,
  .window-control svg,
  .app-name {
    transition: none;
  }
  
  .control-button:hover svg,
  .window-control:hover svg {
    transform: none;
  }
}
</style>
