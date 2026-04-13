<template>
  <div class="console-window">
    <!-- 工具栏 -->
    <div class="console-toolbar">
      <n-button-group>
        <n-button size="small" @click="handleConnect" :disabled="!selectedWebshell">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12.55a11 11 0 0 1 14.08 0"/>
              <path d="M1.42 9a16 16 0 0 1 21.16 0"/>
              <path d="M8.53 16.11a6 6 0 0 1 6.95 0"/>
              <line x1="12" y1="20" x2="12.01" y2="20"/>
            </svg>
          </template>
          连接
        </n-button>
        <n-button size="small" @click="handleDisconnect" :disabled="!isConnected">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18"/>
              <path d="M6 6l12 12"/>
            </svg>
          </template>
          断开
        </n-button>
      </n-button-group>

      <n-divider vertical />

      <n-button-group>
        <n-button size="small" @click="handleClear">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
          </template>
          清空
        </n-button>
        <n-button size="small" @click="handleExport">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
          </template>
          导出
        </n-button>
      </n-button-group>

      <n-divider vertical />

      <n-select
        v-model:value="selectedPlugin"
        :options="pluginOptions"
        size="small"
        style="width: 150px"
        placeholder="选择插件"
      />

      <n-space style="margin-left: auto">
        <n-tag :type="connectionStatus === 'connected' ? 'success' : 'error'" size="small">
          {{ connectionStatus === 'connected' ? '已连接' : '未连接' }}
        </n-tag>
      </n-space>
    </div>

    <!-- 终端输出区域 -->
    <div class="console-output" ref="outputRef">
      <div v-for="(line, index) in outputLines" :key="index" class="console-line" :class="line.type">
        <span class="line-prefix">{{ line.prefix }}</span>
        <span class="line-content">{{ line.content }}</span>
      </div>
    </div>

    <!-- 命令输入区域 -->
    <div class="console-input">
      <n-input
        v-model:value="currentCommand"
        type="textarea"
        placeholder="输入命令..."
        :disabled="!isConnected"
        @keydown="handleKeydown"
        :rows="3"
      />
      <n-button type="primary" @click="handleSendCommand" :disabled="!isConnected || !currentCommand">
        <template #icon>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="22" y1="2" x2="11" y2="13"/>
            <polygon points="22 2 15 22 11 13 2 9 22 2"/>
          </svg>
        </template>
        发送
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { NButton, NButtonGroup, NDivider, NInput, NSelect, NSpace, NTag, useMessage } from 'naive-ui'
import { invoke } from '@tauri-apps/api/core'

interface OutputLine {
  type: 'info' | 'success' | 'error' | 'command'
  prefix: string
  content: string
}

const message = useMessage()

// 状态
const isConnected = ref(false)
const connectionStatus = ref<'connected' | 'disconnected'>('disconnected')
const selectedWebshell = ref<string | null>(null)
const selectedPlugin = ref<string | null>(null)
const currentCommand = ref('')
const outputLines = ref<OutputLine[]>([])

// 插件选项
const pluginOptions = [
  { label: '终端管理', value: 'terminal' },
  { label: '文件管理', value: 'file' },
  { label: '数据库管理', value: 'database' },
  { label: '系统信息', value: 'system' },
]

// 输出区域引用
const outputRef = ref<HTMLElement | null>(null)

// 自动滚动到底部
const scrollToBottom = async () => {
  await nextTick()
  if (outputRef.value) {
    outputRef.value.scrollTop = outputRef.value.scrollHeight
  }
}

// 添加输出行
const addOutputLine = (type: OutputLine['type'], prefix: string, content: string) => {
  outputLines.value.push({ type, prefix, content })
  scrollToBottom()
}

// 连接处理
const handleConnect = async () => {
  if (!selectedWebshell.value) {
    message.warning('请先选择一个 WebShell')
    return
  }

  try {
    addOutputLine('info', 'SYSTEM', '正在连接 WebShell...')
    
    // TODO: 实现真实的连接逻辑
    isConnected.value = true
    connectionStatus.value = 'connected'
    
    addOutputLine('success', 'SYSTEM', '连接成功！')
    message.success('连接成功')
  } catch (error) {
    addOutputLine('error', 'SYSTEM', `连接失败：${error}`)
    message.error(`连接失败：${error}`)
  }
}

// 断开连接处理
const handleDisconnect = () => {
  isConnected.value = false
  connectionStatus.value = 'disconnected'
  addOutputLine('info', 'SYSTEM', '已断开连接')
  message.info('已断开连接')
}

// 清空输出
const handleClear = () => {
  outputLines.value = []
  message.success('已清空输出')
}

// 导出会话
const handleExport = async () => {
  try {
    const content = outputLines.value
      .map(line => `${line.prefix}: ${line.content}`)
      .join('\n')
    
    // TODO: 使用 Tauri 文件系统 API 导出
    message.success('会话已导出')
  } catch (error) {
    message.error(`导出失败：${error}`)
  }
}

// 发送命令
const handleSendCommand = async () => {
  if (!currentCommand.value.trim()) return

  const command = currentCommand.value.trim()
  addOutputLine('command', 'USER', command)
  
  try {
    // TODO: 调用后端命令执行 API
    const response = await invoke('execute_system_command', {
      webshellId: selectedWebshell.value,
      command: command,
    })
    
    addOutputLine('success', 'RESPONSE', String(response))
    message.success('命令执行成功')
  } catch (error) {
    addOutputLine('error', 'ERROR', String(error))
    message.error(`命令执行失败：${error}`)
  }
  
  currentCommand.value = ''
}

// 键盘事件处理（支持 Enter 发送）
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
    e.preventDefault()
    handleSendCommand()
  }
}

// 暴露方法供外部调用
const selectWebshell = (id: string) => {
  selectedWebshell.value = id
  message.info(`已选择 WebShell: ${id}`)
}
</script>

<style scoped>
.console-window {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--n-color);
  border-radius: 8px;
  overflow: hidden;
}

.console-toolbar {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background-color: var(--n-color);
  border-bottom: 1px solid var(--n-border-color);
  gap: 8px;
}

.console-output {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  background-color: #1e1e1e;
  color: #d4d4d4;
}

.console-line {
  margin-bottom: 4px;
  white-space: pre-wrap;
  word-break: break-all;
}

.console-line.info {
  color: #4ec9b0;
}

.console-line.success {
  color: #6a9955;
}

.console-line.error {
  color: #f48771;
}

.console-line.command {
  color: #569cd6;
}

.line-prefix {
  color: #808080;
  margin-right: 8px;
  user-select: none;
}

.line-content {
  color: inherit;
}

.console-input {
  display: flex;
  gap: 8px;
  padding: 12px;
  border-top: 1px solid var(--n-border-color);
  background-color: var(--n-color);
}

.console-input :deep(.n-input) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
}
</style>
