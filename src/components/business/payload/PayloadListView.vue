<template>
  <div class="payload-list-view">
    <n-card title="历史载荷列表" :bordered="false">
      <template #header-extra>
        <n-space>
          <n-button size="small" @click="refreshList">
            <template #icon>
              <span>🔄</span>
            </template>
            刷新
          </n-button>
          <n-button size="small" type="error" @click="clearAll">
            <template #icon>
              <span>🗑️</span>
            </template>
            清空
          </n-button>
        </n-space>
      </template>

      <!-- 搜索和过滤 -->
      <div class="filter-bar" style="margin-bottom: 16px;">
        <n-space>
          <n-input
            v-model:value="searchQuery"
            placeholder="搜索文件名..."
            style="width: 200px;"
            clearable
          />
          <n-select
            v-model:value="filterType"
            :options="[
              { label: '全部类型', value: 'all' },
              { label: 'PHP', value: 'php' },
              { label: 'JSP', value: 'jsp' },
              { label: 'ASPX', value: 'aspx' },
              { label: 'ASP', value: 'asp' },
            ]"
            style="width: 150px;"
          />
        </n-space>
      </div>

      <!-- 载荷列表 -->
      <n-data-table
        :columns="columns"
        :data="filteredPayloads"
        :loading="loading"
        :pagination="pagination"
        :single-line="false"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useMessage, NTag, NButton } from 'naive-ui'
import { usePayloadStore } from '@/stores/payload'
import type { DataTableColumns } from 'naive-ui'

const message = useMessage()
const payloadStore = usePayloadStore()

const loading = ref(false)
const searchQuery = ref('')
const filterType = ref('all')

// 表格列定义
const columns: DataTableColumns = [
  {
    title: '文件名',
    key: 'filename',
    width: 250,
    sorter: 'default',
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render(row) {
      const typeMap: Record<string, string> = {
        php: 'php',
        jsp: 'jsp',
        aspx: 'aspx',
        asp: 'asp',
      }
      const type = typeMap[row.script_type] || 'unknown'
      return h(
        NTag,
        { type: 'info', size: 'small' },
        { default: () => type.toUpperCase() }
      )
    },
  },
  {
    title: '功能',
    key: 'function_type',
    width: 120,
    render(row) {
      const funcMap: Record<string, string> = {
        basic: '基础连接',
        file_manager: '文件管理',
        process_manager: '进程管理',
        registry: '注册表',
        network: '网络',
      }
      return funcMap[row.function_type] || '未知'
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 100,
    render(row) {
      const bytes = row.size
      if (bytes < 1024) return `${bytes} B`
      if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
      return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
    },
  },
  {
    title: '状态',
    key: 'success',
    width: 100,
    render(row) {
      return h(
        NTag,
        { type: row.success ? 'success' : 'error', size: 'small' },
        { default: () => (row.success ? '成功' : '失败') }
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render(row) {
      return h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          onClick: () => viewPayload(row),
        },
        { default: () => '查看' }
      )
    },
  },
]

// 计算属性
const allPayloads = computed(() => payloadStore.history)

const filteredPayloads = computed(() => {
  return allPayloads.value.filter((item: any) => {
    const matchSearch = !searchQuery.value || 
      item.filename.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchType = filterType.value === 'all' || 
      item.script_type === filterType.value
    return matchSearch && matchType
  })
})

// 分页
const pagination = {
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
}

// 方法
const refreshList = async () => {
  loading.value = true
  try {
    await payloadStore.loadHistory()
    message.success('列表已刷新')
  } catch (error: any) {
    message.error('刷新失败')
  } finally {
    loading.value = false
  }
}

const clearAll = () => {
  payloadStore.clearHistory()
  message.success('已清空所有历史记录')
}

const viewPayload = (row: any) => {
  payloadStore.generatedResult = row
  message.success('已加载载荷详情')
}

// 初始化加载
refreshList()
</script>

<style scoped>
.payload-list-view {
  width: 100%;
  height: 100%;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
