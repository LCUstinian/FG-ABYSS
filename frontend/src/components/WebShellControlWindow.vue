<template>
  <div class="webshell-control-window">
    <!-- 顶部工具栏 -->
    <div class="window-toolbar">
      <div class="toolbar-left">
        <div class="webshell-info">
          <span class="webshell-name">{{ webshellInfo?.name || 'WebShell Control' }}</span>
          <n-tag 
            :type="connectionStatus === 'connected' ? 'success' : 'error'" 
            size="small"
          >
            {{ connectionStatus === 'connected' ? '已连接' : '未连接' }}
          </n-tag>
        </div>
      </div>
      
      <div class="toolbar-right">
        <n-space>
          <n-button 
            v-if="connectionStatus !== 'connected'"
            type="primary" 
            size="small"
            @click="handleConnect"
            :loading="connecting"
          >
            <template #icon>
              <span>🔌</span>
            </template>
            连接
          </n-button>
          
          <n-button 
            v-else
            type="error" 
            size="small"
            @click="handleDisconnect"
          >
            <template #icon>
              <span>⏹️</span>
            </template>
            断开
          </n-button>
          
          <n-divider vertical />
          
          <n-button 
            size="small"
            @click="handleClose"
          >
            <template #icon>
              <span>❌</span>
            </template>
            关闭窗口
          </n-button>
        </n-space>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="window-content">
      <!-- 标签页 -->
      <n-tabs
        v-model:value="activeTab"
        type="line"
        animated
        @update:value="handleTabChange"
      >
        <n-tab-pane name="terminal" tab="终端" display="flex" style="flex: 1">
          <WebShellTerminal 
            ref="terminalRef" 
            :webshell-id="webshellId"
          />
        </n-tab-pane>
        
        <n-tab-pane name="file" tab="文件管理" display="flex" style="flex: 1">
          <FileManager 
            ref="fileManagerRef"
            :webshell-id="webshellId"
          />
        </n-tab-pane>
        
        <n-tab-pane name="database" tab="数据库管理" display="flex" style="flex: 1">
          <DatabaseManager 
            ref="databaseManagerRef"
            :webshell-id="webshellId"
          />
        </n-tab-pane>
        
        <n-tab-pane name="command" tab="命令执行" display="flex" style="flex: 1">
          <CommandPanel 
            ref="commandPanelRef"
            :webshell-id="webshellId"
          />
        </n-tab-pane>
      </n-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import WebShellTerminal from './WebShellTerminal.vue'
import FileManager from './FileManager.vue'
import DatabaseManager from './DatabaseManager.vue'
import CommandPanel from './CommandPanel.vue'

// 从 URL hash 参数获取 WebShell ID
// 支持两种格式：
// 1. #/webshell-control?id=xxx
// 2. #webshell-control?id=xxx
const getWebShellIdFromURL = () => {
  const hash = window.location.hash
  if (!hash) return ''
  
  // 移除 # 符号
  const hashWithoutSharp = hash.startsWith('#/') ? hash.slice(2) : hash.slice(1)
  
  // 分割路径和参数
  const [path, queryString] = hashWithoutSharp.split('?')
  
  // 检查是否是 webshell-control 路由
  if (path !== 'webshell-control') {
    return ''
  }
  
  // 解析参数
  if (queryString) {
    const params = new URLSearchParams(queryString)
    return params.get('id') || ''
  }
  
  return ''
}

const webshellId = getWebShellIdFromURL()

// 状态
const activeTab = ref('terminal')
const connectionStatus = ref('connected') // 默认已连接
const connecting = ref(false)
const webshellInfo = ref<any>(null)

// Refs
const terminalRef = ref<InstanceType<typeof WebShellTerminal> | null>(null)
const fileManagerRef = ref<InstanceType<typeof FileManager> | null>(null)
const databaseManagerRef = ref<InstanceType<typeof DatabaseManager> | null>(null)
const commandPanelRef = ref<InstanceType<typeof CommandPanel> | null>(null)

// 连接
const handleConnect = async () => {
  if (!webshellId) {
    console.error('WebShell ID is missing')
    return
  }
  
  connecting.value = true
  try {
    // TODO: 调用连接 API
    // await WebShellConnection.Connect(webshellId)
    connectionStatus.value = 'connected'
    
    // 连接到终端
    terminalRef.value?.connect()
    
    // 加载其他组件
    fileManagerRef.value?.load()
    databaseManagerRef.value?.load()
  } catch (error) {
    console.error('Connection failed:', error)
    connectionStatus.value = 'disconnected'
  } finally {
    connecting.value = false
  }
}

// 断开连接
const handleDisconnect = async () => {
  try {
    // TODO: 调用断开 API
    // await WebShellConnection.Disconnect(webshellId)
    connectionStatus.value = 'disconnected'
    
    // 断开终端
    terminalRef.value?.disconnect()
    
    // 清理其他组件
    fileManagerRef.value?.cleanup()
    databaseManagerRef.value?.cleanup()
  } catch (error) {
    console.error('Disconnect failed:', error)
  }
}

// 标签页切换
const handleTabChange = (name: string) => {
  activeTab.value = name
}

// 关闭窗口
const handleClose = () => {
  // 先断开连接
  if (connectionStatus.value === 'connected') {
    handleDisconnect()
  }
  
  // 关闭窗口（使用 Wails API）
  // @ts-ignore
  if (window.__WAILS__) {
    // @ts-ignore
    window.runtime.WindowClose()
  } else {
    window.close()
  }
}

// 生命周期
onMounted(() => {
  // 加载 WebShell 信息
  loadWebShellInfo()
  
  // 自动连接（可选）
  // handleConnect()
})

onUnmounted(() => {
  // 清理资源
  if (connectionStatus.value === 'connected') {
    handleDisconnect()
  }
})

// 加载 WebShell 信息
const loadWebShellInfo = async () => {
  if (!webshellId) return
  
  try {
    // TODO: 获取 WebShell 详情
    // const info = await WebShell.GetWebShell(webshellId)
    // webshellInfo.value = info
    
    // 模拟数据
    webshellInfo.value = {
      id: webshellId,
      name: 'WebShell Control Window',
      url: 'http://example.com/shell.php',
    }
  } catch (error) {
    console.error('Failed to load WebShell info:', error)
  }
}
</script>

<style scoped lang="scss">
.webshell-control-window {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  background: var(--body-bg-color);
  overflow: hidden;
}

.window-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.webshell-info {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .webshell-name {
    font-weight: 600;
    font-size: 14px;
    color: var(--text-color);
  }
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.window-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

:deep(.n-tabs) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.n-tabs-content) {
  flex: 1;
  overflow: hidden;
}

:deep(.n-tab-pane) {
  height: 100%;
  overflow: auto;
}
</style>
