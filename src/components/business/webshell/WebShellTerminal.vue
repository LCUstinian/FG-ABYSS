<template>
  <div class="webshell-terminal">
    <div class="terminal-container">
      <!-- 未连接状态 -->
      <div v-if="connectionStatus === 'disconnected'" class="terminal-placeholder">
        <div class="placeholder-content">
          <span class="placeholder-icon">📡</span>
          <h3 class="placeholder-title">WebShell 终端</h3>
          <p class="placeholder-description">
            连接到 WebShell 以使用终端功能
          </p>
          <div class="placeholder-info">
            <div class="info-item">
              <span class="info-label">WebShell ID:</span>
              <span class="info-value">{{ webshellId || '未指定' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">状态:</span>
              <n-tag size="small" type="default">未连接</n-tag>
            </div>
          </div>
          <n-button 
            type="primary" 
            @click="handleConnect"
            :loading="connecting"
            style="margin-top: 24px; min-width: 120px;"
          >
            <template #icon>
              <span>🔌</span>
            </template>
            连接终端
          </n-button>
        </div>
      </div>

      <!-- 连接中状态 -->
      <div v-else-if="connectionStatus === 'connecting'" class="terminal-placeholder">
        <div class="placeholder-content">
          <span class="placeholder-icon">⏳</span>
          <h3 class="placeholder-title">正在连接...</h3>
          <p class="placeholder-description">
            正在建立与 WebShell 的安全连接
          </p>
          <n-progress 
            type="line" 
            :percentage="connectionProgress"
            :show-indicator="false"
            :height="4"
            style="width: 300px; margin-top: 24px;"
          />
          <p class="placeholder-text" style="margin-top: 16px;">
            请稍候...
          </p>
        </div>
      </div>

      <!-- 已连接状态 -->
      <div v-else class="terminal-connected">
        <div class="terminal-header">
          <div class="terminal-status">
            <span class="status-dot connected"></span>
            <span class="status-text">已连接</span>
          </div>
          <div class="terminal-info">
            <span class="info-label">终端:</span>
            <span class="info-value">WebShell Terminal</span>
          </div>
          <n-button 
            size="small" 
            type="error" 
            @click="handleDisconnect"
            style="margin-left: auto;"
          >
            <template #icon>
              <span>⏹️</span>
            </template>
            断开连接
          </n-button>
        </div>
        <div class="terminal-body">
          <div class="terminal-mock-output">
            <div class="mock-line">
              <span class="mock-prompt">root@webshell:~$</span>
              <span class="mock-command"> whoami</span>
            </div>
            <div class="mock-output">root</div>
            <div class="mock-line">
              <span class="mock-prompt">root@webshell:~$</span>
              <span class="mock-command"> pwd</span>
            </div>
            <div class="mock-output">/var/www/html</div>
            <div class="mock-line">
              <span class="mock-prompt">root@webshell:~$</span>
              <span class="mock-cursor"></span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NTag, NButton, NProgress } from 'naive-ui'
import { invoke } from '@/utils/tauri-mock-adapter'

defineProps<{
  webshellId?: string
}>()

const connectionStatus = ref('disconnected')
const connecting = ref(false)
const connectionProgress = ref(0)

const connect = async () => {
  if (!webshellId) {
    console.error('WebShell ID is missing')
    return
  }
  
  connecting.value = true
  connectionStatus.value = 'connecting'
  connectionProgress.value = 0
  
  // 模拟连接进度
  const progressInterval = setInterval(() => {
    connectionProgress.value = Math.min(connectionProgress.value + 20, 90)
  }, 200)
  
  try {
    const result = await invoke('connect_webshell', { id: webshellId })
    
    clearInterval(progressInterval)
    connectionProgress.value = 100
    
    if (result && result.success) {
      setTimeout(() => {
        connectionStatus.value = 'connected'
        connecting.value = false
        console.log('Terminal connected successfully')
      }, 300)
    } else {
      throw new Error(result?.message || '连接失败')
    }
  } catch (error) {
    clearInterval(progressInterval)
    connectionStatus.value = 'disconnected'
    connecting.value = false
    console.error('Connection failed:', error)
    alert('连接失败：' + (error instanceof Error ? error.message : '未知错误'))
  }
}

const disconnect = async () => {
  try {
    const result = await invoke('disconnect_webshell', { id: webshellId })
    if (result && result.success) {
      connectionStatus.value = 'disconnected'
      connectionProgress.value = 0
      console.log('Terminal disconnected')
    }
  } catch (error) {
    console.error('Disconnect failed:', error)
  }
}

const handleConnect = () => {
  connect()
}

const handleDisconnect = () => {
  disconnect()
}

defineExpose({
  connect,
  disconnect,
  connectionStatus
})
</script>

<style scoped>
.webshell-terminal {
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: var(--card-bg);
}

.terminal-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 未连接状态 */
.terminal-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.placeholder-content {
  text-align: center;
  max-width: 500px;
}

.placeholder-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 20px;
  opacity: 0.6;
}

.placeholder-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-color);
  margin: 0 0 12px 0;
}

.placeholder-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 20px 0;
}

.placeholder-info {
  background: var(--content-bg);
  border-radius: var(--border-radius-md);
  padding: 16px;
  text-align: left;
}

.info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.info-value {
  font-size: 13px;
  color: var(--text-color);
  font-weight: 500;
}

.placeholder-text {
  font-size: 13px;
  color: var(--text-secondary);
}

/* 连接中状态 */
.terminal-connecting {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 已连接状态 */
.terminal-connected {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  gap: 12px;
}

.terminal-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-tertiary);
}

.status-dot.connected {
  background: #4CAF50;
  box-shadow: 0 0 8px rgba(76, 175, 80, 0.4);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.status-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color);
}

.terminal-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.terminal-body {
  flex: 1;
  background: #1e1e1e;
  padding: 12px 16px;
  overflow: auto;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
}

.terminal-mock-output {
  color: #d4d4d4;
}

.mock-line {
  display: flex;
  align-items: center;
  margin-bottom: 2px;
}

.mock-prompt {
  color: #4CAF50;
  font-weight: 600;
  margin-right: 8px;
}

.mock-command {
  color: #d4d4d4;
}

.mock-output {
  color: #d4d4d4;
  margin-bottom: 4px;
  padding-left: 8px;
}

.mock-cursor {
  width: 8px;
  height: 16px;
  background: #d4d4d4;
  animation: blink 1s infinite;
}

@keyframes blink {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}

/* 滚动条样式 */
.terminal-body::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

.terminal-body::-webkit-scrollbar-track {
  background: #1e1e1e;
}

.terminal-body::-webkit-scrollbar-thumb {
  background: #444;
  border-radius: 5px;
}

.terminal-body::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
