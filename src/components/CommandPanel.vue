<template>
  <div class="command-panel-container">
    <!-- 命令输入区 -->
    <div class="command-input-section">
      <n-input
        v-model:value="commandInput"
        type="textarea"
        placeholder="输入命令，按 Ctrl+Enter 执行"
        :rows="3"
        @keydown.ctrl.enter="handleExecute"
      >
        <template #suffix>
          <n-button type="primary" @click="handleExecute" :loading="executing">
            <template #icon>
              <n-icon><Play /></n-icon>
            </template>
            执行
          </n-button>
        </template>
      </n-input>
    </div>
    
    <!-- 快捷命令 -->
    <div class="quick-commands">
      <n-space>
        <n-tag
          v-for="cmd in quickCommands"
          :key="cmd.label"
          :type="cmd.type || 'default'"
          size="small"
          style="cursor: pointer"
          @click="handleQuickCommand(cmd.command)"
        >
          {{ cmd.label }}
        </n-tag>
      </n-space>
    </div>
    
    <!-- 输出区域 -->
    <div class="command-output-section">
      <div class="output-header">
        <span class="output-title">执行结果</span>
        <n-space>
          <n-button size="small" @click="handleClearOutput">
            <template #icon>
              <n-icon><Delete /></n-icon>
            </template>
            清空
          </n-button>
          <n-button size="small" @click="handleCopyOutput">
            <template #icon>
              <n-icon><Copy /></n-icon>
            </template>
            复制
          </n-button>
        </n-space>
      </div>
      
      <div ref="outputContainer" class="output-content">
        <pre v-if="commandOutput">{{ commandOutput }}</pre>
        <n-empty v-else description="暂无输出" />
      </div>
    </div>
    
    <!-- 命令历史 -->
    <div class="command-history">
      <n-collapse>
        <n-collapse-item title="命令历史" name="history">
          <n-timeline>
            <n-timeline-item
              v-for="(item, index) in commandHistory"
              :key="index"
              :type="item.success ? 'success' : 'error'"
              :title="item.command"
              :content="item.time"
              :detail="item.output"
            />
          </n-timeline>
        </n-collapse-item>
      </n-collapse>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { Play, Delete, Copy } from '@vicons/carbon'

const message = useMessage()

// 状态
const commandInput = ref('')
const commandOutput = ref('')
const executing = ref(false)
const outputContainer = ref<HTMLElement | null>(null)

// 命令历史
const commandHistory = ref<any[]>([])

// 快捷命令
const quickCommands = reactive([
  { label: '系统信息', command: 'uname -a', type: 'info' },
  { label: '当前用户', command: 'whoami', type: 'info' },
  { label: '工作目录', command: 'pwd', type: 'info' },
  { label: '进程列表', command: 'ps aux', type: 'warning' },
  { label: '网络信息', command: 'netstat -an', type: 'warning' },
  { label: '磁盘使用', command: 'df -h', type: 'success' },
  { label: '内存使用', command: 'free -m', type: 'success' },
  { label: 'PHP 版本', command: 'php -v', type: 'info' },
  { label: 'Python 版本', command: 'python --version', type: 'info' }
])

// 执行命令
const handleExecute = async () => {
  const command = commandInput.value.trim()
  if (!command) {
    message.warning('请输入命令')
    return
  }

  executing.value = true
  const startTime = Date.now()

  try {
    // TODO: 调用后端命令执行 API
    // const result = await CommandHandler.ExecuteCommand(webshellId, command, 30)
    
    // 模拟执行
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const mockOutput = `命令执行成功\n\n输出结果示例...\n`
    commandOutput.value = mockOutput
    
    // 添加到历史
    commandHistory.value.unshift({
      command: command,
      output: mockOutput,
      success: true,
      time: new Date().toLocaleString(),
      duration: Date.now() - startTime
    })
    
    // 限制历史记录数量
    if (commandHistory.value.length > 50) {
      commandHistory.value.pop()
    }
    
    message.success('命令执行成功')
  } catch (error: any) {
    commandOutput.value = `执行失败：${error.message}`
    
    // 添加到历史
    commandHistory.value.unshift({
      command: command,
      output: error.message,
      success: false,
      time: new Date().toLocaleString(),
      duration: Date.now() - startTime
    })
    
    message.error('命令执行失败')
  } finally {
    executing.value = false
  }
}

// 快捷命令
const handleQuickCommand = (command: string) => {
  commandInput.value = command
  handleExecute()
}

// 清空输出
const handleClearOutput = () => {
  commandOutput.value = ''
}

// 复制输出
const handleCopyOutput = () => {
  if (!commandOutput.value) {
    message.warning('暂无输出可复制')
    return
  }

  navigator.clipboard.writeText(commandOutput.value)
    .then(() => {
      message.success('复制成功')
    })
    .catch(() => {
      message.error('复制失败')
    })
}

// 生命周期
onMounted(() => {
  // 可以加载保存的命令历史
})

// 暴露方法
defineExpose({
  execute: handleExecute,
  clear: handleClearOutput
})
</script>

<style scoped>
.command-panel-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
  gap: 12px;
  background: #ffffff;
}

.dark .command-panel-container {
  background: #1e1e1e;
}

.command-input-section {
  flex-shrink: 0;
}

.quick-commands {
  flex-shrink: 0;
  padding: 8px 0;
  border-bottom: 1px solid #e8e8e8;
}

.dark .quick-commands {
  border-bottom-color: #3d3d3d;
}

.command-output-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

.dark .command-output-section {
  border-color: #3d3d3d;
}

.output-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f5f5f5;
  border-bottom: 1px solid #e8e8e8;
}

.dark .output-header {
  background: #2d2d2d;
  border-bottom-color: #3d3d3d;
}

.output-title {
  font-weight: 600;
  color: #333;
}

.dark .output-title {
  color: #fff;
}

.output-content {
  flex: 1;
  overflow: auto;
  padding: 12px;
  background: #1e1e1e;
  color: #ffffff;
}

.output-content pre {
  margin: 0;
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.command-history {
  flex-shrink: 0;
  max-height: 200px;
  overflow: auto;
}
</style>
