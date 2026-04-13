<template>
  <div class="terminal-manager">
    <!-- 工具栏 -->
    <div class="terminal-toolbar">
      <n-space>
        <n-button size="small" @click="handleNewTerminal">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="16 18 22 12 16 6"/>
              <polyline points="8 6 2 12 8 18"/>
            </svg>
          </template>
          新建终端
        </n-button>
        <n-button size="small" @click="handleClear">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
          </template>
          清空
        </n-button>
        <n-button size="small" @click="handleCopy">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
            </svg>
          </template>
          复制
        </n-button>
        <n-button size="small" @click="handlePaste">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
              <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
            </svg>
          </template>
          粘贴
        </n-button>
      </n-space>

      <n-space style="margin-left: auto">
        <n-select
          v-model:value="selectedShell"
          :options="shellOptions"
          size="small"
          style="width: 120px"
        />
        <n-tag :type="terminalStatus === 'connected' ? 'success' : 'error'" size="small">
          {{ terminalStatus === 'connected' ? '已连接' : '未连接' }}
        </n-tag>
      </n-space>
    </div>

    <!-- 终端标签页 -->
    <n-tabs
      v-model:value="activeTerminal"
      type="card"
      closable
      @close="handleCloseTerminal"
      @update:value="handleSwitchTerminal"
      style="flex: 1"
    >
      <n-tab-pane
        v-for="terminal in terminals"
        :key="terminal.id"
        :name="terminal.id"
        :tab="terminal.name"
      >
        <div class="terminal-container">
          <!-- 终端输出 -->
          <div class="terminal-output" ref="terminalOutputRef">
            <div
              v-for="(line, index) in terminal.output"
              :key="index"
              class="terminal-line"
              :class="line.type"
            >
              <span v-if="line.type === 'prompt'" class="prompt">{{ line.content }}</span>
              <span v-else class="output">{{ line.content }}</span>
            </div>
            <div v-if="terminalStatus === 'connected'" class="input-line">
              <span class="prompt">{{ currentPrompt }}</span>
              <input
                ref="terminalInputRef"
                v-model="currentInput"
                type="text"
                class="terminal-input"
                @keydown="handleInputKeydown"
                @focus="handleFocus"
                @blur="handleBlur"
                :disabled="terminalStatus !== 'connected'"
              />
            </div>
          </div>
        </div>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { NButton, NInput, NSpace, NSelect, NTabPane, NTabs, NTag, useMessage } from 'naive-ui'
import { invoke } from '@tauri-apps/api/core'

interface TerminalOutput {
  type: 'prompt' | 'output' | 'error' | 'info'
  content: string
}

interface Terminal {
  id: string
  name: string
  output: TerminalOutput[]
  history: string[]
  historyIndex: number
}

const message = useMessage()

// 状态
const terminalStatus = ref<'connected' | 'disconnected'>('disconnected')
const selectedShell = ref('bash')
const activeTerminal = ref<string>('')
const currentInput = ref('')
const currentPrompt = ref('$')

// Shell 选项
const shellOptions = [
  { label: 'Bash', value: 'bash' },
  { label: 'PowerShell', value: 'powershell' },
  { label: 'CMD', value: 'cmd' },
  { label: 'Zsh', value: 'zsh' },
]

// 终端列表
const terminals = ref<Terminal[]>([
  {
    id: '1',
    name: 'Terminal 1',
    output: [
      { type: 'info', content: '欢迎使用 FG-ABYSS 终端管理器' },
      { type: 'info', content: '连接 WebShell 后即可开始使用' },
    ],
    history: [],
    historyIndex: -1,
  },
])

// 终端输出引用
const terminalOutputRef = ref<HTMLElement | null>(null)
const terminalInputRef = ref<HTMLInputElement | null>(null)

// 当前终端
const currentTerminal = computed(() => {
  return terminals.value.find(t => t.id === activeTerminal.value)
})

// 自动滚动到底部
const scrollToBottom = async () => {
  await nextTick()
  if (terminalOutputRef.value) {
    terminalOutputRef.value.scrollTop = terminalOutputRef.value.scrollHeight
  }
}

// 添加输出
const addOutput = (type: TerminalOutput['type'], content: string) => {
  const terminal = currentTerminal.value
  if (terminal) {
    terminal.output.push({ type, content })
    scrollToBottom()
  }
}

// 新建终端
const handleNewTerminal = () => {
  const id = Date.now().toString()
  const terminal: Terminal = {
    id,
    name: `Terminal ${terminals.value.length + 1}`,
    output: [
      { type: 'info', content: '新终端已创建' },
    ],
    history: [],
    historyIndex: -1,
  }
  terminals.value.push(terminal)
  activeTerminal.value = id
  message.success('新终端已创建')
}

// 关闭终端
const handleCloseTerminal = (tabName: string) => {
  const index = terminals.value.findIndex(t => t.id === tabName)
  if (index !== -1) {
    terminals.value.splice(index, 1)
    if (terminals.value.length > 0) {
      activeTerminal.value = terminals.value[Math.max(0, index - 1)].id
    } else {
      activeTerminal.value = ''
    }
    message.success('终端已关闭')
  }
}

// 切换终端
const handleSwitchTerminal = (tabName: string) => {
  activeTerminal.value = tabName
  nextTick(() => {
    terminalInputRef.value?.focus()
  })
}

// 清空当前终端
const handleClear = () => {
  const terminal = currentTerminal.value
  if (terminal) {
    terminal.output = []
    message.success('终端已清空')
  }
}

// 复制
const handleCopy = async () => {
  try {
    const selection = window.getSelection()
    if (selection && selection.toString()) {
      await navigator.clipboard.writeText(selection.toString())
      message.success('已复制到剪贴板')
    } else {
      message.warning('请先选择要复制的文本')
    }
  } catch (error) {
    message.error(`复制失败：${error}`)
  }
}

// 粘贴
const handlePaste = async () => {
  try {
    const text = await navigator.clipboard.readText()
    currentInput.value += text
    terminalInputRef.value?.focus()
  } catch (error) {
    message.error(`粘贴失败：${error}`)
  }
}

// 输入键盘事件
const handleInputKeydown = async (e: KeyboardEvent) => {
  if (e.key === 'Enter') {
    e.preventDefault()
    await executeCommand()
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    navigateHistory('up')
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    navigateHistory('down')
  } else if (e.key === 'Tab') {
    e.preventDefault()
    // TODO: 实现自动补全
  }
}

// 执行命令
const executeCommand = async () => {
  const command = currentInput.value.trim()
  if (!command) return

  const terminal = currentTerminal.value
  if (!terminal) return

  // 添加到历史
  terminal.history.push(command)
  terminal.historyIndex = terminal.history.length

  // 显示命令
  addOutput('prompt', currentPrompt.value)
  addOutput('output', command)

  try {
    // TODO: 调用后端 API 执行命令
    const response = await invoke('execute_system_command', {
      webshellId: 'placeholder',
      command: command,
    })

    addOutput('output', String(response))
  } catch (error) {
    addOutput('error', `错误：${error}`)
  }

  currentInput.value = ''
  scrollToBottom()
}

// 导航历史
const navigateHistory = (direction: 'up' | 'down') => {
  const terminal = currentTerminal.value
  if (!terminal || terminal.history.length === 0) return

  if (direction === 'up') {
    if (terminal.historyIndex > 0) {
      terminal.historyIndex--
      currentInput.value = terminal.history[terminal.historyIndex]
    }
  } else {
    if (terminal.historyIndex < terminal.history.length) {
      terminal.historyIndex++
      if (terminal.historyIndex === terminal.history.length) {
        currentInput.value = ''
      } else {
        currentInput.value = terminal.history[terminal.historyIndex]
      }
    }
  }
}

// 焦点事件
const handleFocus = () => {
  // 终端获得焦点
}

const handleBlur = () => {
  // 终端失去焦点
}

// 初始化
activeTerminal.value = terminals.value[0]?.id || ''
</script>

<style scoped>
.terminal-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
}

.terminal-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.terminal-container {
  display: flex;
  flex-direction: column;
  height: calc(100% - 40px);
}

.terminal-output {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  background-color: #1e1e1e;
  color: #d4d4d4;
}

.terminal-line {
  margin-bottom: 4px;
  white-space: pre-wrap;
  word-break: break-all;
}

.terminal-line.prompt {
  color: #6a9955;
}

.terminal-line.output {
  color: #d4d4d4;
}

.terminal-line.error {
  color: #f48771;
}

.terminal-line.info {
  color: #4ec9b0;
}

.input-line {
  display: flex;
  align-items: center;
  margin-top: 8px;
}

.prompt {
  color: #6a9955;
  margin-right: 8px;
  white-space: nowrap;
}

.terminal-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #d4d4d4;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
}

.terminal-input:focus {
  outline: none;
}
</style>
