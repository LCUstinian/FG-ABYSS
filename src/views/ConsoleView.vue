<template>
  <div class="console-view">
    <n-grid :cols="100" :x-gap="12" style="height: calc(100% - 60px)">
      <!-- 左侧：WebShell 列表 -->
      <n-grid-item :span="25">
        <n-card :title="t('webshells')" :bordered="false" style="height: 100%">
          <template #header-extra>
            <n-button size="small" @click="handleRefresh">
              <template #icon>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"/>
                  <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
                </svg>
              </template>
            </n-button>
          </template>

          <n-input
            v-model:value="searchQuery"
            placeholder="搜索 WebShell..."
            size="small"
            style="margin-bottom: 12px"
            clearable
          />

          <n-list hoverable>
            <n-list-item
              v-for="webshell in filteredWebshells"
              :key="webshell.id"
              @click="handleSelectWebshell(webshell)"
            >
              <template #prefix>
                <n-avatar :style="{ backgroundColor: getStatusColor(webshell.status) }">
                  {{ webshell.name.charAt(0).toUpperCase() }}
                </n-avatar>
              </template>
              <n-thing :title="webshell.name" :description="webshell.url">
                <template #header-extra>
                  <n-tag :type="getStatusType(webshell.status)" size="small">
                    {{ webshell.status }}
                  </n-tag>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-grid-item>

      <!-- 右侧：控制台和插件 -->
      <n-grid-item :span="75">
        <n-card :bordered="false" style="height: 100%">
          <template #header>
            <n-space>
              <span>{{ selectedWebshell ? selectedWebshell.name : '选择 WebShell' }}</span>
              <n-divider vertical />
              <n-tabs v-model:value="activePlugin" type="line" animated>
                <n-tab-pane name="console" tab="控制台" />
                <n-tab-pane name="file" tab="文件管理" />
                <n-tab-pane name="database" tab="数据库" />
                <n-tab-pane name="terminal" tab="终端" />
              </n-tabs>
            </n-space>
          </template>

          <!-- 控制台窗口 -->
          <div v-show="activePlugin === 'console'" style="height: calc(100% - 100px)">
            <ConsoleWindow ref="consoleWindowRef" />
          </div>

          <!-- 文件管理器 -->
          <div v-show="activePlugin === 'file'" style="height: calc(100% - 100px)">
            <FileManager />
          </div>

          <!-- 数据库管理器 -->
          <div v-show="activePlugin === 'database'" style="height: calc(100% - 100px)">
            <DatabaseManager />
          </div>

          <!-- 终端管理器 -->
          <div v-show="activePlugin === 'terminal'" style="height: calc(100% - 100px)">
            <TerminalManager />
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NAvatar, NButton, NCard, NDivider, NGrid, NGridItem, NInput, NList, NListItem, NSpace, NTabPane, NTabs, NTag, NThing, useMessage } from 'naive-ui'
import ConsoleWindow from '../components/ConsoleWindow.vue'
import FileManager from '../components/FileManager.vue'
import DatabaseManager from '../components/DatabaseManager.vue'
import TerminalManager from '../components/TerminalManager.vue'

const { t } = useI18n()
const message = useMessage()

// 状态
const searchQuery = ref('')
const activePlugin = ref<'console' | 'file' | 'database' | 'terminal'>('console')
const selectedWebshell = ref<any>(null)
const consoleWindowRef = ref<any>(null)

// 模拟 WebShell 数据
const webshells = ref([
  { id: '1', name: 'Test Shell', url: 'http://example.com/shell.php', status: 'online', payload_type: 'PHP' },
  { id: '2', name: 'Production', url: 'https://prod.example.com/cmd.jsp', status: 'offline', payload_type: 'JSP' },
  { id: '3', name: 'Dev Server', url: 'http://dev.example.com/shell.aspx', status: 'online', payload_type: 'ASPX' },
])

// 过滤后的 WebShell 列表
const filteredWebshells = computed(() => {
  if (!searchQuery.value) return webshells.value
  return webshells.value.filter(ws =>
    ws.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    ws.url.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 获取状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'online':
      return 'success'
    case 'offline':
      return 'error'
    default:
      return 'default'
  }
}

// 获取状态颜色
const getStatusColor = (status: string) => {
  switch (status) {
    case 'online':
      return '#52c41a'
    case 'offline':
      return '#ff4d4f'
    default:
      return '#d9d9d9'
  }
}

// 刷新列表
const handleRefresh = async () => {
  // TODO: 调用后端 API 刷新 WebShell 列表
  message.success('刷新成功')
}

// 选择 WebShell
const handleSelectWebshell = (webshell: any) => {
  selectedWebshell.value = webshell
  message.info(`已选择：${webshell.name}`)
  
  // 通知控制台窗口
  if (consoleWindowRef.value) {
    consoleWindowRef.value.selectWebshell(webshell.id)
  }
}
</script>

<style scoped>
.console-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
  box-sizing: border-box;
}

.console-view :deep(.n-card) {
  height: 100%;
}

.console-view :deep(.n-card__content) {
  height: calc(100% - 60px);
  overflow: auto;
}
</style>
