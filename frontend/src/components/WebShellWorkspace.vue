<template>
  <div class="webshell-workspace">
    <!-- 标签页 -->
    <n-tabs
      v-model:value="activeTab"
      type="line"
      animated
      @update:value="handleTabChange"
    >
      <n-tab-pane name="terminal" tab="终端" display="flex" style="flex: 1">
        <WebShellTerminal ref="terminalRef" />
      </n-tab-pane>
      
      <n-tab-pane name="file" tab="文件管理" display="flex" style="flex: 1">
        <FileManager ref="fileManagerRef" />
      </n-tab-pane>
      
      <n-tab-pane name="command" tab="命令执行" display="flex" style="flex: 1">
        <CommandPanel ref="commandPanelRef" />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import WebShellTerminal from './WebShellTerminal.vue'
import FileManager from './FileManager.vue'
import CommandPanel from './CommandPanel.vue'

const activeTab = ref('terminal')

const terminalRef = ref<InstanceType<typeof WebShellTerminal> | null>(null)
const fileManagerRef = ref<InstanceType<typeof FileManager> | null>(null)
const commandPanelRef = ref<InstanceType<typeof CommandPanel> | null>(null)

const handleTabChange = (name: string) => {
  activeTab.value = name
}

// 暴露方法给父组件
defineExpose({
  connect: () => terminalRef.value?.connect(),
  disconnect: () => terminalRef.value?.disconnect()
})
</script>

<style scoped>
.webshell-workspace {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #ffffff;
}

.dark .webshell-workspace {
  background: #1e1e1e;
}

:deep(.n-tabs) {
  height: 100%;
}

:deep(.n-tabs-content) {
  height: 100%;
}

:deep(.n-tab-pane) {
  height: 100%;
}
</style>
