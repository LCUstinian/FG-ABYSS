<template>
  <div class="database-manager">
    <!-- 工具栏 -->
    <div class="db-toolbar">
      <n-space>
        <n-button size="small" @click="handleRefresh">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 4 23 10 17 10"/>
              <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
            </svg>
          </template>
          刷新
        </n-button>
        <n-button size="small" @click="handleNewQuery">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 19l7-7 3 3-7 7-3-3z"/>
              <path d="M18 13l-1.5-7.5L2 2l3.5 14.5L13 18l5-5z"/>
              <path d="M2 2l7.586 7.586"/>
              <circle cx="11" cy="11" r="2"/>
            </svg>
          </template>
          新建查询
        </n-button>
        <n-button size="small" @click="handleExport">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
          </template>
          导出
        </n-button>
      </n-space>

      <n-space style="margin-left: auto">
        <n-select
          v-model:value="selectedDatabase"
          :options="databaseOptions"
          size="small"
          style="width: 200px"
          placeholder="选择数据库"
        />
      </n-space>
    </div>

    <n-grid :cols="2" :x-gap="12" style="flex: 1">
      <!-- 左侧：表列表 -->
      <n-card title="数据表" :bordered="false" style="height: 100%">
        <n-input
          v-model:value="tableSearchQuery"
          placeholder="搜索表..."
          size="small"
          style="margin-bottom: 8px"
          clearable
        />
        <n-tree
          :data="filteredTables"
          :selectable="true"
          block-line
          @update:selected-keys="handleSelectTable"
        />
      </n-card>

      <!-- 右侧：表数据/查询编辑器 -->
      <n-card :title="activeTab === 'data' ? `表数据：${selectedTable}` : 'SQL 查询'" :bordered="false" style="height: 100%">
        <template #header-extra>
          <n-tabs v-model:value="activeTab" type="line" animated>
            <n-tab-pane name="data" tab="数据" />
            <n-tab-pane name="query" tab="查询" />
            <n-tab-pane name="structure" tab="结构" />
          </n-tabs>
        </template>

        <!-- 数据视图 -->
        <div v-if="activeTab === 'data'" class="data-view">
          <n-data-table
            :columns="tableColumns"
            :data="tableData"
            :row-key="rowKey"
            :single-line="false"
            :max-height="400"
          />
        </div>

        <!-- 查询编辑器 -->
        <div v-if="activeTab === 'query'" class="query-view">
          <n-input
            v-model:value="sqlQuery"
            type="textarea"
            placeholder="输入 SQL 查询..."
            :rows="6"
            style="font-family: monospace; margin-bottom: 12px"
          />
          <n-space>
            <n-button type="primary" @click="executeQuery">
              <template #icon>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polygon points="5 3 19 12 5 21 5 3"/>
                </svg>
              </template>
              执行
            </n-button>
            <n-button @click="clearQuery">
              <template #icon>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </template>
              清空
            </n-button>
          </n-space>

          <!-- 查询结果 -->
          <div v-if="queryResult" class="query-result">
            <n-alert type="info" title="查询结果" style="margin-top: 12px">
              <n-data-table
                :columns="queryResultColumns"
                :data="queryResult.data"
                :row-key="rowKey"
                :max-height="300"
              />
            </n-alert>
          </div>
        </div>

        <!-- 结构视图 -->
        <div v-if="activeTab === 'structure'" class="structure-view">
          <n-data-table
            :columns="structureColumns"
            :data="tableStructure"
            :row-key="rowKey"
            :single-line="false"
          />
        </div>
      </n-card>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { NAlert, NButton, NCard, NDataTable, NGrid, NInput, NSpace, NSelect, NTabPane, NTabs, NTree, useMessage } from 'naive-ui'
import type { DataTableColumns, TreeOption } from 'naive-ui'
import { invoke } from '@tauri-apps/api/core'

interface TableData {
  key: string
  [key: string]: any
}

interface QueryResult {
  data: TableData[]
  columns: string[]
}

const message = useMessage()

// 状态
const selectedDatabase = ref<string | null>(null)
const selectedTable = ref<string>('')
const tableSearchQuery = ref('')
const activeTab = ref<'data' | 'query' | 'structure'>('data')
const sqlQuery = ref('')
const queryResult = ref<QueryResult | null>(null)

// 数据库选项
const databaseOptions = [
  { label: 'MySQL', value: 'mysql' },
  { label: 'PostgreSQL', value: 'postgresql' },
  { label: 'SQLite', value: 'sqlite' },
  { label: 'SQL Server', value: 'mssql' },
]

// 模拟表数据
const tableData = ref<TableData[]>([
  { key: '1', id: 1, name: '张三', email: 'zhangsan@example.com', created_at: '2024-01-01' },
  { key: '2', id: 2, name: '李四', email: 'lisi@example.com', created_at: '2024-01-02' },
  { key: '3', id: 3, name: '王五', email: 'wangwu@example.com', created_at: '2024-01-03' },
])

// 模拟表结构
const tableStructure = ref([
  { key: '1', field: 'id', type: 'INT', null: 'NO', key: 'PRI', default: null, extra: 'auto_increment' },
  { key: '2', field: 'name', type: 'VARCHAR(100)', null: 'NO', key: '', default: null, extra: '' },
  { key: '3', field: 'email', type: 'VARCHAR(255)', null: 'NO', key: 'UNI', default: null, extra: '' },
  { key: '4', field: 'created_at', type: 'TIMESTAMP', null: 'YES', key: '', default: 'CURRENT_TIMESTAMP', extra: '' },
])

// 表树形数据
const tablesTree = ref<TreeOption[]>([
  {
    key: 'users',
    label: 'users',
    isLeaf: true,
  },
  {
    key: 'products',
    label: 'products',
    isLeaf: true,
  },
  {
    key: 'orders',
    label: 'orders',
    isLeaf: true,
  },
])

// 过滤后的表
const filteredTables = computed(() => {
  if (!tableSearchQuery.value) return tablesTree.value
  return filterTree(tablesTree.value, tableSearchQuery.value)
})

// 数据视图列
const tableColumns = computed<DataTableColumns<TableData>>(() => {
  if (tableData.value.length === 0) return []
  const keys = Object.keys(tableData.value[0]).filter(k => k !== 'key')
  return keys.map(key => ({
    title: key,
    key: key,
    sortable: 'default',
  }))
})

// 查询结果列
const queryResultColumns = computed<DataTableColumns<TableData>>(() => {
  if (!queryResult.value || queryResult.value.data.length === 0) return []
  const keys = Object.keys(queryResult.value.data[0]).filter(k => k !== 'key')
  return keys.map(key => ({
    title: key,
    key: key,
  }))
})

// 结构视图列
const structureColumns: DataTableColumns<any> = [
  { title: '字段', key: 'field' },
  { title: '类型', key: 'type' },
  { title: 'NULL', key: 'null' },
  { title: '键', key: 'key' },
  { title: '默认值', key: 'default' },
  { title: '额外', key: 'extra' },
]

// 行键
const rowKey = (row: any) => row.key

// 过滤树形数据
const filterTree = (tree: TreeOption[], query: string): TreeOption[] => {
  return tree
    .filter(node => node.label.toLowerCase().includes(query.toLowerCase()))
    .map(node => ({
      ...node,
      children: node.children ? filterTree(node.children, query) : undefined,
    }))
    .filter(node => node.children?.length || node.label.toLowerCase().includes(query.toLowerCase()))
}

// 选择表
const handleSelectTable = (keys: Array<string | number>) => {
  if (keys.length > 0) {
    selectedTable.value = String(keys[0])
    loadTableData(selectedTable.value)
  }
}

// 加载表数据
const loadTableData = async (tableName: string) => {
  try {
    // TODO: 调用后端 API 加载表数据
    message.success(`加载表 ${tableName} 数据成功`)
  } catch (error) {
    message.error(`加载表数据失败：${error}`)
  }
}

// 刷新
const handleRefresh = () => {
  if (selectedTable.value) {
    loadTableData(selectedTable.value)
  }
  message.success('刷新成功')
}

// 新建查询
const handleNewQuery = () => {
  activeTab.value = 'query'
  sqlQuery.value = 'SELECT * FROM ' + (selectedTable.value || 'table_name') + ' LIMIT 100;'
}

// 导出
const handleExport = () => {
  // TODO: 实现导出功能
  message.info('导出功能开发中...')
}

// 执行查询
const executeQuery = async () => {
  if (!sqlQuery.value.trim()) {
    message.warning('请输入 SQL 查询')
    return
  }

  try {
    // TODO: 调用后端 API 执行 SQL 查询
    queryResult.value = {
      data: tableData.value,
      columns: ['id', 'name', 'email', 'created_at'],
    }
    message.success('查询执行成功')
  } catch (error) {
    message.error(`查询执行失败：${error}`)
  }
}

// 清空查询
const clearQuery = () => {
  sqlQuery.value = ''
  queryResult.value = null
}
</script>

<style scoped>
.database-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
}

.db-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.data-view,
.query-view,
.structure-view {
  height: 100%;
  overflow: auto;
}

.query-result {
  margin-top: 12px;
}
</style>
