<template>
  <div class="webshell-control-window">
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

    <div class="window-content">
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
import { ref, onMounted, onUnmounted } from 'vue'
import { NTabs, NTabPane, NTag, NButton, NSpace, NDivider } from 'naive-ui'
import { getCurrentWindow } from '@tauri-apps/api/window'
import { invoke } from '@/utils/tauri-mock-adapter'
import WebShellTerminal from '@/components/business/webshell/WebShellTerminal.vue'
import FileManager from '@/components/business/webshell/FileManager.vue'
import DatabaseManager from '@/components/business/database/DatabaseManager.vue'
import CommandPanel from '@/components/business/webshell/CommandPanel.vue'

const getWebShellIdFromURL = () => {
  const hash = window.location.hash
  if (!hash) return ''
  
  const hashWithoutSharp = hash.startsWith('#/') ? hash.slice(2) : hash.slice(1)
  const [path, queryString] = hashWithoutSharp.split('?')
  
  if (path !== 'webshell-control') {
    return ''
  }
  
  if (queryString) {
    const params = new URLSearchParams(queryString)
    return params.get('id') || ''
  }
  
  return ''
}

const webshellId = getWebShellIdFromURL()

const activeTab = ref('terminal')
const connectionStatus = ref('connected')
const connecting = ref(false)
const webshellInfo = ref<any>(null)

const terminalRef = ref<InstanceType<typeof WebShellTerminal> | null>(null)
const fileManagerRef = ref<InstanceType<typeof FileManager> | null>(null)
const databaseManagerRef = ref<InstanceType<typeof DatabaseManager> | null>(null)
const commandPanelRef = ref<InstanceType<typeof CommandPanel> | null>(null)

const handleConnect = async () => {
  if (!webshellId) {
    console.error('WebShell ID is missing')
    return
  }
  
  connecting.value = true
  try {
    const result = await invoke('connect_webshell', { id: webshellId })
    if (result && result.success) {
      connectionStatus.value = 'connected'
      terminalRef.value?.connect()
      fileManagerRef.value?.load()
      databaseManagerRef.value?.load()
    } else {
      throw new Error(result?.message || '连接失败')
    }
  } catch (error) {
    console.error('Connection failed:', error)
    connectionStatus.value = 'disconnected'
  } finally {
    connecting.value = false
  }
}

const handleDisconnect = async () => {
  try {
    const result = await invoke('disconnect_webshell', { id: webshellId })
    if (result && result.success) {
      connectionStatus.value = 'disconnected'
      terminalRef.value?.disconnect()
      fileManagerRef.value?.cleanup()
      databaseManagerRef.value?.cleanup()
    }
  } catch (error) {
    console.error('Disconnect failed:', error)
  }
}

const handleTabChange = (name: string) => {
  activeTab.value = name
}

const handleClose = async () => {
  if (connectionStatus.value === 'connected') {
    handleDisconnect()
  }
  
  try {
    const appWindow = getCurrentWindow()
    await appWindow.close()
  } catch (error) {
    console.error('关闭窗口失败:', error)
    window.close()
  }
}

const loadWebShellInfo = async () => {
  if (!webshellId) return
  
  try {
    const result = await invoke('get_webshell', { id: webshellId })
    if (result) {
      webshellInfo.value = result
    }
  } catch (error) {
    console.error('Failed to load WebShell info:', error)
    
    webshellInfo.value = {
      id: webshellId,
      name: 'WebShell Control Window',
      url: 'http://example.com/shell.php',
    }
  }
}

onMounted(() => {
  loadWebShellInfo()
})

onUnmounted(() => {
  if (connectionStatus.value === 'connected') {
    handleDisconnect()
  }
})
</script>

<style scoped>
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
}

.webshell-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--text-color);
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
