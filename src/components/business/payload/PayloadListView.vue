<template>
  <div class="payload-list-view">
    <n-card :title="t('payload.historyTitle')" :bordered="false">
      <template #header-extra>
        <n-space>
          <n-button size="small" @click="refreshList">
            <template #icon>
              <span>🔄</span>
            </template>
            {{ t('common.refresh') }}
          </n-button>
          <n-button size="small" type="error" @click="clearAll">
            <template #icon>
              <span>🗑️</span>
            </template>
            {{ t('common.clear') }}
          </n-button>
        </n-space>
      </template>

      <!-- 搜索和过滤 -->
      <div class="filter-bar" style="margin-bottom: 16px;">
        <n-space>
          <n-input
            v-model:value="searchQuery"
            :placeholder="t('common.searchFilenamePlaceholder')"
            style="width: 200px;"
            clearable
          />
          <n-select
            v-model:value="filterType"
            :options="[
              { label: t('payload.allTypes'), value: 'all' },
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
import { useI18n } from 'vue-i18n'
import { useMessage, NTag, NButton } from 'naive-ui'
import { usePayloadStore } from '@/stores/payload'
import type { DataTableColumns } from 'naive-ui'

const { t } = useI18n()
const message = useMessage()
const payloadStore = usePayloadStore()

const loading = ref(false)
const searchQuery = ref('')
const filterType = ref('all')

// 表格列定义
const columns: DataTableColumns = [
  {
    title: t('payload.filename'),
    key: 'filename',
    width: 250,
    sorter: 'default',
  },
  {
    title: t('payload.type'),
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
    title: t('payload.function'),
    key: 'function_type',
    width: 120,
    render(row) {
      const funcMap: Record<string, string> = {
        basic: t('payload.functionBasic'),
        file_manager: t('payload.functionFileManager'),
        process_manager: t('payload.functionProcessManager'),
        registry: t('payload.functionRegistry'),
        network: t('payload.functionNetwork'),
      }
      return funcMap[row.function_type] || t('payload.functionUnknown')
    },
  },
  {
    title: t('payload.size'),
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
    title: t('payload.status'),
    key: 'success',
    width: 100,
    render(row) {
      return h(
        NTag,
        { type: row.success ? 'success' : 'error', size: 'small' },
        { default: () => (row.success ? t('payload.statusSuccess') : t('payload.statusFailed')) }
      )
    },
  },
  {
    title: t('payload.action'),
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
        { default: () => t('payload.actionView') }
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
    message.success(t('payload.refreshSuccess'))
  } catch (error: any) {
    message.error(t('payload.refreshError'))
  } finally {
    loading.value = false
  }
}

const clearAll = () => {
  payloadStore.clearHistory()
  message.success(t('payload.clearSuccess'))
}

const viewPayload = (row: any) => {
  payloadStore.generatedResult = row
  message.success(t('payload.loadSuccess'))
}

// 初始化加载
refreshList()
</script>

<style scoped>
.payload-list-view {
  width: 100%;
  min-height: 100%;
  animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 卡片样式优化 */
:deep(.n-card) {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.n-card-header) {
  border-bottom: 1px solid var(--border-color);
  padding: 20px 24px;
  background: linear-gradient(135deg, var(--card-bg) 0%, var(--card-bg-hover) 100%);
  flex-shrink: 0;
}

:deep(.n-card-header__main) {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

:deep(.n-card__content) {
  padding: 24px;
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--card-bg-hover);
  border-radius: 12px;
  margin-bottom: 20px;
  border: 1px solid var(--border-color);
  flex-wrap: wrap;
}

/* 按钮样式优化 */
:deep(.n-button) {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.n-button:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 输入框和选择器样式优化 */
:deep(.n-input),
:deep(.n-select) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.n-input:hover),
:deep(.n-select:hover) {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

:deep(.n-input:focus-within),
:deep(.n-select:focus-within) {
  box-shadow: 0 0 0 3px rgba(var(--active-color-rgb), 0.1);
}

/* 数据表格样式优化 */
:deep(.n-data-table) {
  border-radius: 12px;
  border: 1px solid var(--border-color);
  overflow: hidden;
  flex: 1;
  display: flex;
  flex-direction: column;
}

:deep(.n-data-table-wrapper) {
  border-radius: 12px;
}

:deep(.n-data-table-th) {
  background: linear-gradient(135deg, var(--card-bg-hover) 0%, var(--card-bg) 100%);
  font-weight: 600;
  color: var(--text-primary);
  padding: 16px 20px;
  border-bottom: 2px solid var(--border-color);
}

:deep(.n-data-table-td) {
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

:deep(.n-data-table-tr:hover .n-data-table-td) {
  background: rgba(var(--active-color-rgb), 0.05);
}

:deep(.n-data-table-tr--selected .n-data-table-td) {
  background: rgba(var(--active-color-rgb), 0.1);
}

/* 标签样式优化 */
:deep(.n-tag) {
  border-radius: 6px;
  font-weight: 500;
  padding: 4px 12px;
}

/* 分页样式优化 */
:deep(.n-pagination) {
  margin-top: 20px;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

:deep(.n-pagination-item) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.n-pagination-item:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

:deep(.n-pagination-item--active) {
  background: var(--active-color) !important;
  color: white !important;
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  :deep(.n-card-header) {
    padding: 24px 28px;
  }
  
  :deep(.n-card-header__main) {
    font-size: 20px;
  }
  
  :deep(.n-card__content) {
    padding: 28px;
  }
  
  .filter-bar {
    gap: 20px;
    padding: 20px;
  }
  
  :deep(.n-data-table-th),
  :deep(.n-data-table-td) {
    padding: 18px 24px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  :deep(.n-card-header) {
    padding: 16px 20px;
  }
  
  :deep(.n-card__content) {
    padding: 16px;
  }
  
  .filter-bar {
    gap: 12px;
    padding: 12px;
  }
  
  :deep(.n-data-table-th),
  :deep(.n-data-table-td) {
    padding: 12px 16px;
    font-size: 13px;
  }
  
  :deep(.n-pagination) {
    justify-content: center;
  }
}
</style>
