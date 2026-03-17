<template>
  <div class="webshell-terminal-container">
    <!-- 工具栏 -->
    <div class="terminal-toolbar">
      <div class="toolbar-left">
        <n-button size="small" @click="handleConnect" :loading="connecting">
          <template #icon>
            <n-icon><Power /></n-icon>
          </template>
          {{ connecting ? '连接中...' : '连接' }}
        </n-button>
        
        <n-button size="small" @click="handleDisconnect" :disabled="!connected">
          <template #icon>
            <n-icon><Disconnect /></n-icon>
          </template>
          断开
        </n-button>
        
        <n-divider vertical />
        
        <n-select
          v-model:value="selectedEncoder"
          :options="encoderOptions"
          size="small"
          style="width: 120px"
          placeholder="编码器"
        />
        
        <n-input
          v-model:value="password"
          type="password"
          placeholder="密码"
          size="small"
          style="width: 150px"
          @keyup.enter="handleConnect"
        />
      </div>
      
      <div class="toolbar-right">
        <n-tag :type="connectionStatusTagType" size="small">
          {{ connectionStatusText }}
        </n-tag>
      </div>
    </div>
    
    <!-- 终端容器 -->
    <div ref="terminalContainer" class="terminal-container"></div>
    
    <!-- 连接配置弹窗 -->
    <n-modal
      v-model:show="showConfigModal"
      preset="dialog"
      title="WebShell 连接配置"
      style="width: 600px"
    >
      <n-form :model="configForm" label-placement="left" label-width="100">
        <n-form-item label="目标 URL">
          <n-input v-model:value="configForm.url" placeholder="http://example.com/shell.php" />
        </n-form-item>
        
        <n-form-item label="编码器">
          <n-select
            v-model:value="configForm.encoder"
            :options="encoderOptions"
            placeholder="选择编码器"
          />
        </n-form-item>
        
        <n-form-item label="密码">
          <n-input
            v-model:value="configForm.password"
            type="password"
            placeholder="连接密码"
          />
        </n-form-item>
        
        <n-form-item label="加密密钥">
          <n-input
            v-model:value="configForm.encryptionKey"
            placeholder="XOR 加密密钥（可选）"
          />
        </n-form-item>
        
        <n-form-item label="超时时间">
          <n-input-number
            v-model:value="configForm.timeout"
            :min="5"
            :max="300"
            placeholder="秒"
          />
        </n-form-item>
        
        <n-form-item label="代理类型">
          <n-select
            v-model:value="configForm.proxyType"
            :options="proxyOptions"
            placeholder="选择代理"
          />
        </n-form-item>
        
        <n-form-item label="代理地址" v-if="configForm.proxyType !== 'none'">
          <n-input
            v-model:value="configForm.proxyAddress"
            placeholder="http://127.0.0.1:8080"
          />
        </n-form-item>
        
        <n-form-item label="SSL 验证">
          <n-switch v-model:value="configForm.sslVerify" />
        </n-form-item>
      </n-form>
      
      <template #action>
        <n-button @click="showConfigModal = false">取消</n-button>
        <n-button type="primary" @click="handleConfigConfirm">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import { Power } from '@vicons/carbon'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'

const message = useMessage()

// 连接状态
const connecting = ref(false)
const connected = ref(false)
const showConfigModal = ref(false)

// 终端相关
const terminalContainer = ref<HTMLElement | null>(null)
let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null

// 配置表单
const configForm = reactive({
  url: '',
  encoder: 'none',
  password: '',
  encryptionKey: '',
  timeout: 30,
  proxyType: 'none',
  proxyAddress: '',
  sslVerify: false
})

// 快捷操作
const selectedEncoder = ref('none')
const password = ref('')

// 选项
const encoderOptions = [
  { label: '无编码', value: 'none' },
  { label: 'Base64', value: 'base64' },
  { label: 'ROT13', value: 'rot13' },
  { label: 'XOR', value: 'xor' }
]

const proxyOptions = [
  { label: '无代理', value: 'none' },
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' }
]

// 连接状态计算
const connectionStatusText = computed(() => {
  if (connecting.value) return '连接中...'
  if (connected.value) return '已连接'
  return '未连接'
})

const connectionStatusTagType = computed(() => {
  if (connecting.value) return 'warning'
  if (connected.value) return 'success'
  return 'default'
})

// 初始化终端
const initTerminal = async () => {
  if (!terminalContainer.value) return

  // 创建终端实例
  terminal = new Terminal({
    cursorBlink: true,
    cursorStyle: 'block',
    fontSize: 14,
    fontFamily: 'Consolas, "Courier New", monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#ffffff',
      cursor: '#ffffff',
      cursorAccent: '#000000',
      selection: 'rgba(255, 255, 255, 0.3)',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#e5e5e5'
    },
    scrollback: 10000,
    tabStopWidth: 4,
    rightClickSelectsWord: true
  })

  // 添加插件
  fitAddon = new FitAddon()
  const webLinksAddon = new WebLinksAddon()
  terminal.loadAddon(fitAddon)
  terminal.loadAddon(webLinksAddon)

  // 打开终端
  terminal.open(terminalContainer.value)
  fitAddon.fit()

  // 显示欢迎信息
  terminal.writeln('\x1b[1;32m=== FG-ABYSS WebShell Terminal ===\x1b[0m')
  terminal.writeln('\x1b[33m请点击"连接"按钮开始新的会话\x1b[0m')
  terminal.writeln('')

  // 监听输入
  terminal.onData((data) => {
    handleTerminalInput(data)
  })

  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
}

// 处理终端输入
const handleTerminalInput = (data: string) => {
  if (!connected.value) {
    message.warning('请先连接 WebShell')
    return
  }

  // TODO: 发送命令到后端
  terminal?.write(data)
}

// 处理窗口大小变化
const handleResize = () => {
  if (fitAddon && terminal) {
    fitAddon.fit()
  }
}

// 连接 WebShell
const handleConnect = async () => {
  if (!configForm.url) {
    showConfigModal.value = true
    return
  }

  connecting.value = true

  try {
    // TODO: 调用后端连接 API
    // const result = await ConnectionHandler.Connect(...)
    
    // 模拟连接
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    connected.value = true
    message.success('连接成功')
    terminal?.writeln('\x1b[1;32m[SUCCESS] WebShell 连接成功\x1b[0m')
    terminal?.writeln('')
  } catch (error: any) {
    message.error('连接失败：' + error.message)
    terminal?.writeln(`\x1b[1;31m[ERROR] 连接失败：${error.message}\x1b[0m`)
  } finally {
    connecting.value = false
  }
}

// 断开连接
const handleDisconnect = async () => {
  try {
    // TODO: 调用后端断开 API
    // await ConnectionHandler.Disconnect(...)
    
    connected.value = false
    message.success('已断开连接')
    terminal?.writeln('\x1b[1;33m[INFO] 连接已断开\x1b[0m')
    terminal?.writeln('')
  } catch (error: any) {
    message.error('断开失败：' + error.message)
  }
}

// 配置确认
const handleConfigConfirm = async () => {
  if (!configForm.url) {
    message.error('请输入目标 URL')
    return
  }

  showConfigModal.value = false
  await handleConnect()
}

// 生命周期
onMounted(() => {
  nextTick(() => {
    initTerminal()
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  terminal?.dispose()
})

// 暴露方法
defineExpose({
  connect: handleConnect,
  disconnect: handleDisconnect
})
</script>

<style scoped>
.webshell-terminal-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1e1e1e;
}

.terminal-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #3d3d3d;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.terminal-container {
  flex: 1;
  overflow: hidden;
  padding: 4px;
}

:deep(.xterm) {
  height: 100%;
}

:deep(.xterm-viewport) {
  overflow-y: auto;
}
</style>
