<template>
  <div class="audit-logs">
    <n-card title="安全审计日志" :bordered="false">
      <n-alert type="info" title="审计说明" style="margin-bottom: 20px">
        记录所有关键操作和安全事件，支持搜索、过滤和导出
      </n-alert>

      <!-- 统计信息 -->
      <n-grid :cols="4" :x-gap="16" :y-gap="16" style="margin-bottom: 20px">
        <n-grid-item>
          <n-statistic label="总日志数">
            {{ statistics.total_logs }}
          </n-statistic>
        </n-grid-item>
        <n-grid-item>
          <n-statistic label="信息">
            <n-text type="info">{{ statistics.info_count }}</n-text>
          </n-statistic>
        </n-grid-item>
        <n-grid-item>
          <n-statistic label="警告">
            <n-text type="warning">{{ statistics.warning_count }}</n-text>
          </n-statistic>
        </n-grid-item>
        <n-grid-item>
          <n-statistic label="错误">
            <n-text type="error">{{ statistics.error_count }}</n-text>
          </n-statistic>
        </n-grid-item>
      </n-grid>

      <!-- 工具栏 -->
      <n-space justify="space-between" style="margin-bottom: 16px">
        <n-space>
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索日志..."
            clearable
            style="width: 300px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <n-icon :component="SearchOutline" />
            </template>
            <template #suffix>
              <n-button text type="primary" @click="handleSearch">
                搜索
              </n-button>
            </template>
          </n-input>

          <n-select
            v-model:value="filterLevel"
            :options="levelOptions"
            placeholder="级别过滤"
            clearable
            style="width: 150px"
            @update:value="handleFilter"
          />
        </n-space>

        <n-space>
          <n-button @click="handleRefresh">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            刷新
          </n-button>
          <n-button @click="handleExport">
            <template #icon>
              <n-icon :component="DownloadOutline" />
            </template>
            导出
          </n-button>
          <n-button type="error" @click="handleClear" secondary>
            <template #icon>
              <n-icon :component="TrashOutline" />
            </template>
            清空
          </n-button>
        </n-space>
      </n-space>

      <!-- 日志表格 -->
      <n-data-table
        :columns="columns"
        :data="logData"
        :loading="loading"
        :pagination="pagination"
        :scroll-x="1200"
        size="small"
        striped
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import {
  SearchOutline,
  RefreshOutline,
  DownloadOutline,
  TrashOutline,
} from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'
import { GetLogs, GetStatistics, SearchLogs, ExportLogs, ClearLogs } from '@bindings/AuditHandler'

const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
const searchKeyword = ref('')
const filterLevel = ref('')

const statistics = reactive({
  total_logs: 0,
  info_count: 0,
  warning_count: 0,
  error_count: 0,
  critical_count: 0,
  success_count: 0,
  failure_count: 0,
})

const logData = ref<any[]>([])

const levelOptions = [
  { label: 'INFO', value: 'INFO' },
  { label: 'WARNING', value: 'WARNING' },
  { label: 'ERROR', value: 'ERROR' },
  { label: 'CRITICAL', value: 'CRITICAL' },
]

const pagination = reactive({
  page: 1,
  pageSize: 20,
  pageSizes: [20, 50, 100],
  showSizePicker: true,
  onChange: (page: number) => {
    pagination.page = page
    loadLogs()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.page = 1
    pagination.pageSize = pageSize
    loadLogs()
  },
})

const columns: DataTableColumns = [
  {
    title: '时间',
    key: 'timestamp',
    width: 180,
    fixed: 'left',
  },
  {
    title: '级别',
    key: 'level',
    width: 80,
    render: (row: any) => {
      const typeMap: Record<string, string> = {
        INFO: 'info',
        WARNING: 'warning',
        ERROR: 'error',
        CRITICAL: 'error',
      }
      return h('n-tag', { type: typeMap[row.level] || 'default', size: 'small' }, () => row.level)
    },
  },
  {
    title: '操作类型',
    key: 'operation',
    width: 120,
  },
  {
    title: '用户',
    key: 'user',
    width: 120,
  },
  {
    title: '资源',
    key: 'resource',
    width: 150,
    ellipsis: { tooltip: true },
  },
  {
    title: '动作',
    key: 'action',
    width: 150,
    ellipsis: { tooltip: true },
  },
  {
    title: '状态',
    key: 'success',
    width: 80,
    render: (row: any) => {
      return row.success
        ? h('n-tag', { type: 'success', size: 'small' }, () => '成功')
        : h('n-tag', { type: 'error', size: 'small' }, () => '失败')
    },
  },
  {
    title: '耗时 (ms)',
    key: 'duration',
    width: 90,
  },
  {
    title: '详情',
    key: 'details',
    ellipsis: { tooltip: true },
  },
]

const loadStatistics = async () => {
  try {
    const stats = await GetStatistics()
    statistics.total_logs = stats.total_logs
    statistics.info_count = stats.info_count
    statistics.warning_count = stats.warning_count
    statistics.error_count = stats.error_count
    statistics.critical_count = stats.critical_count
    statistics.success_count = stats.success_count
    statistics.failure_count = stats.failure_count
  } catch (error) {
    console.error('加载统计失败:', error)
  }
}

const loadLogs = async () => {
  try {
    loading.value = true
    const response = await GetLogs({
      limit: pagination.pageSize,
      offset: (pagination.page - 1) * pagination.pageSize,
    })
    logData.value = response.logs
  } catch (error) {
    message.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = async () => {
  if (!searchKeyword.value) {
    loadLogs()
    return
  }

  try {
    loading.value = true
    const response = await SearchLogs({
      keyword: searchKeyword.value,
      limit: 100,
    })
    logData.value = response.logs
    message.success('搜索完成')
  } catch (error) {
    message.error('搜索失败')
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  loadLogs()
}

const handleRefresh = () => {
  loadStatistics()
  loadLogs()
  message.success('刷新成功')
}

const handleExport = async () => {
  try {
    const filePath = `data/audit_logs_${Date.now()}.json`
    await ExportLogs({ file_path: filePath })
    message.success('导出成功')
  } catch (error) {
    message.error('导出失败')
  }
}

const handleClear = () => {
  dialog.warning({
    title: '确认清空',
    content: '确定要清空所有审计日志吗？此操作不可恢复！',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await ClearLogs()
        loadStatistics()
        loadLogs()
        message.success('清空成功')
      } catch (error) {
        message.error('清空失败')
      }
    },
  })
}

onMounted(() => {
  loadStatistics()
  loadLogs()
})
</script>

<style scoped lang="scss">
.audit-logs {
  padding: 20px;
}
</style>
