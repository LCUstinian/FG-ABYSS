<template>
  <div class="status-bar" :class="{ 'fullscreen-hidden': isFullscreen }">
    <div class="status-left">
      <div class="status-item" :title="t('status.memoryTip')">
        <svg class="status-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="2" y="6" width="20" height="12" rx="2"/>
          <path d="M6 10h.01M10 10h.01"/>
        </svg>
        <span>{{ systemStatus.memoryUsage }}</span>
      </div>
    </div>
    
    <div class="status-center">
      <div class="status-item" :title="t('status.processIdTip')">
        <svg class="status-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="4" y="4" width="16" height="16" rx="2"/>
          <path d="M9 9h6M9 13h6M9 17h2"/>
        </svg>
        <span>ID: {{ systemStatus.processId }}</span>
      </div>
      
      <div class="status-item" :title="t('status.cpuTip')">
        <svg class="status-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="4" y="4" width="16" height="16" rx="2"/>
          <path d="M9 9h6v6H9z"/>
          <path d="M9 1v3M15 1v3M9 20v3M15 20v3M20 9h3M20 14h3M1 9h3M1 14h3"/>
        </svg>
        <span>{{ systemStatus.cpuUsage }}</span>
      </div>
      
      <div class="status-item" :title="t('status.uptimeTip')">
        <svg class="status-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M12 6v6l4 2"/>
        </svg>
        <span>{{ formatUptime(systemStatus.uptime) }}</span>
      </div>
    </div>
    
    <div class="status-right">
      <div class="status-item status-time" :title="t('status.currentTime')">
        <svg class="status-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M12 3v1M12 20v1M3 12h1M20 12h1"/>
        </svg>
        <span>{{ currentTime }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  systemStatus: {
    type: Object,
    required: true,
    default: () => ({
      memoryUsage: '0 GB / 0 GB',
      processId: '0',
      cpuUsage: '0%',
      uptime: '0'
    })
  }
})

const currentTime = ref('')
const isFullscreen = ref(false)

// 格式化运行时间
const formatUptime = (uptime: string): string => {
  const seconds = parseFloat(uptime)
  if (isNaN(seconds)) return '0s'
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  
  if (hours > 0) {
    return `${hours}h ${minutes}m ${secs}s`
  } else if (minutes > 0) {
    return `${minutes}m ${secs}s`
  } else {
    return `${secs}s`
  }
}

// 更新时间
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { 
    hour: '2-digit', 
    minute: '2-digit',
    second: '2-digit'
  })
}

// 监听全屏状态
const checkFullscreen = () => {
  isFullscreen.value = document.fullscreenElement !== null
}

onMounted(() => {
  updateTime()
  const timer = setInterval(updateTime, 1000)
  
  document.addEventListener('fullscreenchange', checkFullscreen)
  
  onUnmounted(() => {
    clearInterval(timer)
    document.removeEventListener('fullscreenchange', checkFullscreen)
  })
})
</script>

<style scoped>
.status-bar {
  height: 32px;
  background: var(--status-bar-bg);
  border-top: 1px solid var(--border-color);
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  font-size: 12px;
  color: var(--text-color);
  transition: all 0.3s ease;
  width: 100%;
  box-sizing: border-box;
  z-index: 100;
}

.status-bar.fullscreen-hidden {
  opacity: 0;
  transform: translateY(100%);
  position: absolute;
  bottom: 0;
}

.status-bar:hover.fullscreen-hidden {
  opacity: 1;
  transform: translateY(0);
}

.status-left,
.status-center,
.status-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
  cursor: default;
  white-space: nowrap;
}

.status-item:hover {
  background: var(--hover-color);
  transform: translateY(-1px);
}

.status-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  opacity: 0.7;
  transition: opacity 0.2s ease;
}

.status-item:hover .status-icon {
  opacity: 1;
}

.status-time {
  font-weight: 600;
  font-family: 'SF Mono', 'Consolas', 'Monaco', monospace;
  letter-spacing: 0.5px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .status-bar {
    padding: 0 12px;
    gap: 8px;
  }
  
  .status-left,
  .status-center,
  .status-right {
    gap: 12px;
  }
  
  .status-item {
    padding: 2px 6px;
    font-size: 11px;
  }
  
  .status-icon {
    width: 12px;
    height: 12px;
  }
}

@media (max-width: 768px) {
  .status-bar {
    height: 28px;
    padding: 0 8px;
    font-size: 10px;
    justify-content: space-between;
  }
  
  .status-left {
    /* 小屏幕下仍然显示内存信息 */
    display: flex;
  }
  
  .status-center {
    gap: 8px;
  }
  
  .status-item {
    padding: 2px 4px;
  }
  
  /* 在小屏幕下隐藏部分文本，但保留图标 */
  .status-center .status-item span {
    display: inline;
  }
  
  .status-right {
    margin-left: auto;
  }
}

/* 超小屏幕优化 */
@media (max-width: 480px) {
  .status-bar {
    padding: 0 4px;
    font-size: 9px;
  }
  
  .status-left .status-item span {
    display: none;
  }
  
  .status-center .status-item span {
    display: none;
  }
  
  .status-icon {
    width: 11px;
    height: 11px;
  }
}

/* 深色主题优化 */
:deep(.dark) .status-bar {
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.3);
}

:deep(.dark) .status-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

/* 浅色主题优化 */
:deep(:root:not(.dark)) .status-item:hover {
  background: rgba(0, 0, 0, 0.03);
}
</style>