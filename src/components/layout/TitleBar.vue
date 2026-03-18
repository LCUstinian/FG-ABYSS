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
        <Tooltip :text="isDarkTheme ? '切换到浅色主题' : '切换到深色主题'">
          <button @click="toggleTheme" class="control-button theme-button" type="button" aria-label="切换主题">
            <div class="button-icon">
              <Sun v-if="isDarkTheme" :size="18" :stroke-width="2" />
              <Moon v-else :size="18" :stroke-width="2" />
            </div>
          </button>
        </Tooltip>
        
        <!-- 语言切换按钮 -->
        <Tooltip :text="locale === 'zh-CN' ? 'Switch to English' : '切换到中文'">
          <button @click="toggleLanguage" class="control-button language-button" type="button" aria-label="切换语言">
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
          <button @click="minimizeWindow" class="window-control minimize" type="button" aria-label="最小化">
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
          <button @click="handleClose" class="window-control close" type="button" aria-label="关闭">
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
const appName = ref('FG-ABYSS 非攻 - 渊渟')

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
    message.error('窗口操作失败')
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
    message.error('关闭窗口失败')
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
   右侧按钮区域
   ============================================ */

.title-bar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

/* 按钮组容器 */
.button-group {
  display: flex;
  align-items: center;
  gap: 2px;
}

/* 窗口控制按钮组特殊处理 */
.button-group.window-controls {
  gap: 2px;
}

/* ============================================
   分隔线
   ============================================ */

.divider {
  width: 1px;
  height: 24px;
  background: var(--border-color);
  margin: 0 4px;
  flex-shrink: 0;
  transition: background var(--transition-normal);
}

/* ============================================
   按钮通用样式 - 统一视觉规范
   ============================================ */

.control-button,
.window-control {
  /* 尺寸规格 */
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  
  /* 视觉样式 */
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  color: var(--text-color);
  padding: 0;
  outline: none;
  -webkit-tap-highlight-color: transparent;
  
  /* 过渡动画 */
  transition: all var(--transition-normal);
  
  /* 图标统一规格 */
  font-size: 18px;
  font-weight: 500;
}

/* 按钮图标容器 - 确保图标完美居中 */
.button-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  /* 确保图标视觉居中 */
  line-height: 1;
}

/* 图标统一视觉权重 */
.control-button svg,
.window-control svg {
  /* 统一图标线条粗细 */
  stroke-width: 2px;
  /* 统一图标尺寸 */
  width: 18px;
  height: 18px;
  /* 平滑过渡 */
  transition: all var(--transition-fast);
}

/* 按钮悬停效果 - 统一反馈 */
.control-button:hover,
.window-control:hover {
  background: var(--hover-color);
}

/* 按钮按下效果 - 统一反馈 */
.control-button:active,
.window-control:active {
  transform: scale(0.9);
  background: var(--active-color-suppl);
}

/* 按钮禁用状态 */
.control-button:disabled,
.window-control:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: transparent !important;
  transform: none !important;
}

/* 图标悬停效果 - 统一缩放 */
.control-button:hover svg,
.window-control:hover svg {
  transform: scale(1.1);
}

/* 图标按下效果 */
.control-button:active svg,
.window-control:active svg {
  transform: scale(0.95);
}

/* ============================================
   主题切换按钮 - 统一视觉风格
   ============================================ */

.theme-button {
  /* 统一颜色 */
  color: var(--warning-color);
}

.theme-button svg {
  /* 统一图标视觉权重 */
  stroke-width: 2px;
}

.theme-button:hover {
  color: var(--warning-color);
  /* 主题色背景 */
  background: rgba(245, 158, 11, 0.12);
  /* 添加阴影增强视觉 */
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.2);
}

.theme-button:active {
  background: rgba(245, 158, 11, 0.2);
  box-shadow: 0 1px 4px rgba(245, 158, 11, 0.3);
}

/* ============================================
   语言切换按钮 - 精确居中对齐
   ============================================ */

.language-button {
  /* 统一颜色 */
  color: var(--info-color);
}

/* 确保 .button-icon 容器完全居中 */
.language-button .button-icon {
  /* 继承父容器的 flex 居中 */
  display: flex;
  align-items: center;
  justify-content: center;
  /* 确保完全填充 */
  width: 100%;
  height: 100%;
}

/* 语言图标精确居中 */
.language-icon {
  /* 字体设置 */
  font-size: 20px;
  line-height: 1;
  
  /* Emoji 字体 */
  font-family: 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji', sans-serif;
  font-style: normal;
  font-weight: normal;
  
  /* 精确居中对齐 - 三重保障 */
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  vertical-align: middle;
  
  /* 消除所有可能的间距 */
  margin: 0;
  padding: 0;
  letter-spacing: 0;
  word-spacing: 0;
  white-space: nowrap;
  
  /* 确保绝对居中 - 使用 transform 微调 */
  transform: translateY(0);
  position: relative;
  
  /* 过渡动画 */
  transition: transform var(--transition-fast);
}

.language-button:hover {
  /* 主题色背景 */
  background: rgba(59, 130, 246, 0.12);
  /* 添加阴影增强视觉 */
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.2);
}

.language-button:hover .language-icon {
  transform: scale(1.1);
}

.language-button:active {
  background: rgba(59, 130, 246, 0.2);
  box-shadow: 0 1px 4px rgba(59, 130, 246, 0.3);
}

.language-button:active .language-icon {
  transform: scale(0.95);
}

/* ============================================
   窗口控制按钮 - 统一视觉风格
   ============================================ */

/* 最小化按钮 */
.window-control.minimize {
  /* 统一颜色 */
  color: var(--text-color);
  /* 统一视觉权重 */
  opacity: 0.9;
}

.window-control.minimize svg {
  /* 统一图标线条 */
  stroke-width: 2px;
}

.window-control.minimize:hover {
  opacity: 1;
  /* 统一悬停背景 */
  background: var(--hover-color);
}

.window-control.minimize:active {
  background: var(--active-color-suppl);
}

/* 最大化按钮 */
.window-control.maximize {
  /* 统一颜色 */
  color: var(--text-color);
  /* 统一视觉权重 */
  opacity: 0.9;
}

.window-control.maximize svg {
  /* 统一图标线条 */
  stroke-width: 2px;
  /* 优化视觉大小 - 微调使其与其他图标一致 */
  width: 17px;
  height: 17px;
}

.window-control.maximize:hover {
  opacity: 1;
  /* 统一悬停背景 */
  background: var(--hover-color);
}

.window-control.maximize:active {
  background: var(--active-color-suppl);
}

/* 关闭按钮 - 特殊处理但保持视觉统一 */
.window-control.close {
  /* 统一颜色 */
  color: var(--error-color);
  /* 统一视觉权重 */
  opacity: 0.95;
}

.window-control.close svg {
  /* 统一图标线条 */
  stroke-width: 2px;
}

.window-control.close:hover {
  /* 红色背景强调 */
  background: var(--error-color);
  color: white;
  opacity: 1;
  /* 增强阴影效果 */
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}

.window-control.close:active {
  /* 按下效果 */
  background: var(--error-color);
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.5);
}

/* ============================================
   深色模式适配 - 统一视觉表现
   ============================================ */

.title-bar.dark {
  background: var(--title-bar-bg, var(--bg-color));
  border-bottom-color: var(--border-color);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.title-bar.dark .app-name {
  color: var(--text-color);
  opacity: 0.95;
}

.title-bar.dark .divider {
  background: var(--border-color);
  opacity: 0.5;
}

/* 深色模式下的按钮效果 - 统一视觉 */
.title-bar.dark .control-button:hover,
.title-bar.dark .window-control:hover {
  /* 深色模式专用悬停背景 */
  background: rgba(255, 255, 255, 0.1);
}

.title-bar.dark .control-button:active,
.title-bar.dark .window-control:active {
  /* 深色模式专用按下背景 */
  background: rgba(255, 255, 255, 0.15);
}

/* 深色模式下主题按钮 */
.title-bar.dark .theme-button:hover {
  background: rgba(245, 158, 11, 0.15);
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.3);
}

/* 深色模式下语言按钮 */
.title-bar.dark .language-button:hover {
  background: rgba(59, 130, 246, 0.15);
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

/* 深色模式下关闭按钮 */
.title-bar.dark .window-control.close:hover {
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.5);
}

/* ============================================
   响应式设计 - 统一视觉表现
   ============================================ */

/* 平板尺寸 */
@media (max-width: 768px) {
  .title-bar {
    height: 44px;
    padding: 0 10px;
  }
  
  .app-name {
    font-size: 13px;
  }
  
  /* 按钮尺寸适度缩小 */
  .control-button,
  .window-control {
    width: 38px;
    height: 38px;
  }
  
  /* 图标尺寸保持不变，确保可见性 */
  .control-button svg,
  .window-control svg {
    width: 17px;
    height: 17px;
  }
  
  .divider {
    height: 20px;
    margin: 0 3px;
  }
  
  .title-bar-right {
    gap: 6px;
  }
  
  .button-group {
    gap: 1px;
  }
}

/* 手机尺寸 */
@media (max-width: 480px) {
  .title-bar {
    height: 40px;
    padding: 0 8px;
  }
  
  .app-name {
    font-size: 12px;
  }
  
  /* 按钮尺寸进一步缩小 */
  .control-button,
  .window-control {
    width: 36px;
    height: 36px;
  }
  
  /* 图标尺寸适度缩小 */
  .control-button svg,
  .window-control svg {
    width: 16px;
    height: 16px;
    stroke-width: 1.8px;
  }
  
  /* emoji 图标适度缩小 */
  .language-icon {
    font-size: 18px;
  }
  
  .divider {
    height: 18px;
    margin: 0 2px;
  }
  
  .title-bar-right {
    gap: 4px;
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
