<template>
  <div class="database-manager">
    <n-grid :cols="24" :x-gap="20" :y-gap="20">
      <n-grid-item :span="8">
        <n-card :bordered="false" class="connection-list-card">
          <template #header>
            <n-space justify="space-between">
              <n-text>数据库连接</n-text>
              <n-button type="primary" size="small" @click="showAddModal = true">
                <template #icon>
                  <span>➕</span>
                </template>
                新建
              </n-button>
            </n-space>
          </template>

          <n-scrollbar style="max-height: 600px">
            <n-list>
              <n-list-item
                v-for="conn in connections"
                :key="conn.id"
                @click="selectConnection(conn)"
                :class="{ active: selectedConn?.id === conn.id }"
              >
                <template #prefix>
                  <span style="font-size: 24px;">{{ getDatabaseIcon(conn.type) }}</span>
                </template>
                <n-space vertical>
                  <n-text strong>{{ conn.name }}</n-text>
                  <n-text depth="3" style="font-size: 12px">
                    {{ conn.host }}:{{ conn.port }}/{{ conn.database }}
                  </n-text>
                </n-space>
                <template #suffix>
                  <n-space vertical>
                    <n-tag :type="getConnectionStatusType(conn.id)" size="small">
                      {{ getConnectionStatus(conn.id) }}
                    </n-tag>
                    <n-button
                      quaternary
                      size="tiny"
                      @click.stop="handleTestConnection(conn)"
                    >
                      测试
                    </n-button>
                  </n-space>
                </template>
              </n-list-item>
            </n-list>

            <n-empty
              v-if="connections.length === 0"
              description="暂无数据库连接"
              size="small"
            />
          </n-scrollbar>
        </n-card>
      </n-grid-item>

      <n-grid-item :span="16">
        <n-tabs v-model:value="activeTab" type="line" animated>
          <n-tab-pane name="query" tab="SQL 查询" display="flex">
            <n-card :bordered="false" class="query-card">
              <n-space vertical>
                <n-input
                  v-model:value="sqlQuery"
                  type="textarea"
                  :rows="6"
                  placeholder="输入 SQL 查询语句..."
                  style="font-family: monospace"
                />
                <n-space justify="space-between">
                  <n-space>
                    <n-button
                      type="primary"
                      @click="handleExecuteQuery"
                      :loading="executing"
                    >
                      <template #icon>
                        <span>▶️</span>
                      </template>
                      执行
                    </n-button>
                    <n-button
                      @click="handleClearQuery"
                      :disabled="!sqlQuery"
                    >
                      清空
                    </n-button>
                  </n-space>
                  <n-space>
                    <n-input
                      v-model:value="queryLimit"
                      type="number"
                      :min="1"
                      :max="10000"
                      placeholder="限制行数"
                      style="width: 120px"
                    />
                  </n-space>
                </n-space>
              </n-space>

              <n-divider />

              <n-data-table
                :columns="resultColumns"
                :data="queryResults"
                :pagination="pagination"
                :scroll-x="1000"
                size="small"
                striped
              />

              <n-space v-if="queryResult" justify="space-between" style="margin-top: 12px">
                <n-text depth="3">
                  执行时间：{{ queryResult.duration }} | 影响行数：{{ queryResult.affected }}
                </n-text>
                <n-space>
                  <n-button
                    size="small"
                    @click="handleExportResults"
                    :disabled="!queryResults.length"
                  >
                    <template #icon>
                      <span>⬇️</span>
                    </template>
                    导出
                  </n-button>
                </n-space>
              </n-space>
            </n-card>
          </n-tab-pane>

          <n-tab-pane name="tables" tab="表结构" display="flex">
            <n-card :bordered="false" class="tables-card">
              <template #header>
                <n-space justify="space-between">
                  <n-text>数据表</n-text>
                  <n-button
                    size="small"
                    @click="handleRefreshTables"
                    :disabled="!selectedConn"
                  >
                    <template #icon>
                      <span>🔄</span>
                    </template>
                    刷新
                  </n-button>
                </n-space>
              </template>

              <n-grid :cols="2" :x-gap="12">
                <n-grid-item
                  v-for="table in tables"
                  :key="table.name"
                  @click="handleSelectTable(table)"
                >
                  <n-card
                    :class="{ 'table-card': true, selected: selectedTable?.name === table.name }"
                    :bordered="false"
                    size="small"
                  >
                    <n-space vertical>
                      <n-text strong>{{ table.name }}</n-text>
                      <n-text depth="3" style="font-size: 12px">
                        {{ table.rows }} 行 | {{ formatSize(table.size) }}
                      </n-text>
                      <n-text depth="3" style="font-size: 12px" v-if="table.comment">
                        {{ table.comment }}
                      </n-text>
                    </n-space>
                  </n-card>
                </n-grid-item>
              </n-grid>

              <n-empty
                v-if="tables.length === 0"
                description="暂无数据表"
                size="small"
              />
            </n-card>
          </n-tab-pane>

          <n-tab-pane name="columns" tab="表结构详情" display="flex">
            <n-card :bordered="false" class="columns-card">
              <template #header>
                <n-space justify="space-between">
                  <n-text>
                    {{ selectedTable ? `${selectedTable.name} - 列信息` : '请选择表' }}
                  </n-text>
                </n-space>
              </template>

              <n-data-table
                :columns="columnColumns"
                :data="tableColumns"
                :pagination="false"
                size="small"
                striped
              />

              <n-empty
                v-if="tableColumns.length === 0"
                description="请先选择数据表"
                size="small"
              />
            </n-card>
          </n-tab-pane>
        </n-tabs>
      </n-grid-item>
    </n-grid>

    <n-modal
      v-model:show="showAddModal"
      preset="dialog"
      title="添加数据库连接"
      style="width: 700px"
    >
      <n-form
        ref="connFormRef"
        :model="connForm"
        :rules="connFormRules"
        label-placement="top"
      >
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item :span="2">
            <n-form-item label="连接名称" path="name">
              <n-input v-model:value="connForm.name" placeholder="自定义连接名称" />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="1">
            <n-form-item label="数据库类型" path="type">
              <n-select
                v-model:value="connForm.type"
                :options="dbTypeOptions"
                placeholder="选择类型"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="1">
            <n-form-item label="字符集" path="charset">
              <n-input v-model:value="connForm.charset" placeholder="utf8" />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="1">
            <n-form-item label="主机地址" path="host">
              <n-input v-model:value="connForm.host" placeholder="localhost" />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="1">
            <n-form-item label="端口" path="port">
              <n-input-number
                v-model:value="connForm.port"
                placeholder="默认端口"
                style="width: 100%"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="1">
            <n-form-item label="用户名" path="username">
              <n-input v-model:value="connForm.username" placeholder="root" />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="1">
            <n-form-item label="密码" path="password">
              <n-input
                v-model:value="connForm.password"
                type="password"
                show-password-on="click"
                placeholder="密码"
              />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="2">
            <n-form-item label="数据库名" path="database">
              <n-input v-model:value="connForm.database" placeholder="数据库名称" />
            </n-form-item>
          </n-grid-item>

          <n-grid-item :span="2">
            <n-form-item label="SSL 模式">
              <n-switch v-model:value="connForm.ssl_mode" />
            </n-form-item>
          </n-grid-item>
        </n-grid>
      </n-form>

      <template #action>
        <n-space justify="end">
          <n-button @click="showAddModal = false">取消</n-button>
          <n-button @click="handleTestNewConnection">测试连接</n-button>
          <n-button type="primary" @click="handleAddConnection">添加</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { invoke } from '@/utils/tauri-mock-adapter'

interface DatabaseConnection {
  id: number
  name: string
  type: string
  host: string
  port: number
  username: string
  password: string
  database: string
  charset: string
  ssl_mode: boolean
}

interface DatabaseTable {
  name: string
  schema: string
  rows: number
  size: number
  engine: string
  collation: string
  comment: string
}

interface TableColumn {
  name: string
  type: string
  length: string
  nullable: boolean
  default: string
  auto_inc: boolean
  primary_key: boolean
  comment: string
}

const message = useMessage()
const activeTab = ref('query')
const showAddModal = ref(false)
const selectedConn = ref<DatabaseConnection | null>(null)
const selectedTable = ref<DatabaseTable | null>(null)
const executing = ref(false)

const connections = ref<DatabaseConnection[]>([])
const tables = ref<DatabaseTable[]>([])
const tableColumns = ref<TableColumn[]>([])
const sqlQuery = ref('')
const queryLimit = ref(100)
const queryResults = ref<any[]>([])
const queryResult = ref<any>(null)

const connFormRef = ref<any | null>(null)
const connForm = reactive({
  name: '',
  type: 'mysql',
  host: 'localhost',
  port: 3306,
  username: 'root',
  password: '',
  database: '',
  charset: 'utf8',
  ssl_mode: false,
})

const connFormRules = {
  name: { required: true, message: '请输入连接名称', trigger: 'blur' },
  type: { required: true, message: '请选择数据库类型', trigger: 'change' },
  host: { required: true, message: '请输入主机地址', trigger: 'blur' },
  port: { required: true, message: '请输入端口', trigger: 'blur' },
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
  database: { required: true, message: '请输入数据库名', trigger: 'blur' },
}

const dbTypeOptions = [
  { label: 'MySQL', value: 'mysql' },
  { label: 'PostgreSQL', value: 'postgresql' },
  { label: 'SQLite', value: 'sqlite' },
  { label: 'SQL Server', value: 'mssql' },
]

const pagination = reactive({
  pageSize: 20,
  pageSizes: [20, 50, 100],
  showSizePicker: true,
})

const resultColumns = computed(() => {
  if (!queryResults.value.length) return []
  const keys = Object.keys(queryResults.value[0])
  return keys.map((key) => ({
    title: key,
    key: key,
    ellipsis: { tooltip: true },
    width: 150,
  }))
})

const columnColumns = [
  { title: '列名', key: 'name', width: 150 },
  { title: '类型', key: 'type', width: 120 },
  { title: '长度', key: 'length', width: 80 },
  {
    title: '可空',
    key: 'nullable',
    width: 60,
    render: (row: TableColumn) => row.nullable ? '是' : '否',
  },
  { title: '默认值', key: 'default', width: 120 },
  {
    title: '主键',
    key: 'primary_key',
    width: 60,
    render: (row: TableColumn) => row.primary_key ? '✓' : '',
  },
  {
    title: '自增',
    key: 'auto_inc',
    width: 60,
    render: (row: TableColumn) => row.auto_inc ? '✓' : '',
  },
  { title: '注释', key: 'comment', ellipsis: { tooltip: true } },
]

const getDatabaseIcon = (type: string) => {
  const icons: Record<string, string> = {
    mysql: '🐬',
    postgresql: '🐘',
    sqlite: '🪨',
    mssql: '🔷',
  }
  return icons[type] || '📊'
}

const getConnectionStatus = (connId: number) => '已连接'
const getConnectionStatusType = (connId: number) => 'success'

const selectConnection = (conn: DatabaseConnection) => {
  selectedConn.value = conn
  handleRefreshTables()
}

const handleAddConnection = async () => {
  try {
    await connFormRef.value?.validate()

    await invoke('database_connect', {
      id: Date.now(),
      name: connForm.name,
      type: connForm.type,
      host: connForm.host,
      port: connForm.port,
      username: connForm.username,
      password: connForm.password,
      database: connForm.database,
      charset: connForm.charset,
      ssl_mode: connForm.ssl_mode,
    })

    connections.value.push({ ...connForm, id: Date.now() })
    message.success('连接添加成功')
    showAddModal.value = false

    Object.assign(connForm, {
      name: '',
      type: 'mysql',
      host: 'localhost',
      port: 3306,
      username: 'root',
      password: '',
      database: '',
      charset: 'utf8',
      ssl_mode: false,
    })
  } catch (error: any) {
    if (error.errors) return
    message.error('添加连接失败：' + (error.message || '未知错误'))
  }
}

const handleTestNewConnection = async () => {
  try {
    await connFormRef.value?.validate()

    const response = await invoke('database_test_connection', {
      type: connForm.type,
      host: connForm.host,
      port: connForm.port,
      username: connForm.username,
      password: connForm.password,
      database: connForm.database,
      charset: connForm.charset,
      ssl_mode: connForm.ssl_mode,
    })

    if (response && response.success) {
      message.success('连接测试成功')
    } else {
      message.error('连接测试失败')
    }
  } catch (error: any) {
    if (error.errors) return
    message.error('测试失败：' + (error.message || '未知错误'))
  }
}

const handleTestConnection = async (conn: DatabaseConnection) => {
  try {
    const response = await invoke('database_test_connection', {
      type: conn.type,
      host: conn.host,
      port: conn.port,
      username: conn.username,
      password: conn.password,
      database: conn.database,
      charset: conn.charset,
      ssl_mode: conn.ssl_mode,
    })

    if (response && response.success) {
      message.success('连接测试成功')
    } else {
      message.error('连接测试失败')
    }
  } catch (error: any) {
    message.error('测试失败')
  }
}

const handleExecuteQuery = async () => {
  if (!selectedConn.value) {
    message.warning('请先选择数据库连接')
    return
  }

  if (!sqlQuery.value.trim()) {
    message.warning('请输入 SQL 查询语句')
    return
  }

  executing.value = true
  try {
    const response = await invoke('database_execute_query', {
      sql: sqlQuery.value,
      connectionId: selectedConn.value.id.toString(),
      limit: queryLimit.value,
    })

    if (response) {
      queryResults.value = response.rows || []
      queryResult.value = {
        duration: response.duration || 0,
        affected: response.affected || 0,
      }
      message.success('查询执行成功')
    }
  } catch (error: any) {
    message.error('查询执行失败：' + (error.message || '未知错误'))
  } finally {
    executing.value = false
  }
}

const handleClearQuery = () => {
  sqlQuery.value = ''
  queryResults.value = []
  queryResult.value = null
}

const handleRefreshTables = async () => {
  if (!selectedConn.value) return

  try {
    const response = await invoke('database_get_tables', {
      connectionId: selectedConn.value.id.toString(),
    })

    if (response) {
      tables.value = response.tables || []
      message.success('表列表刷新成功')
    }
  } catch (error: any) {
    message.error('刷新表列表失败')
  }
}

const handleSelectTable = async (table: DatabaseTable) => {
  selectedTable.value = table

  if (!selectedConn.value) return

  try {
    const response = await invoke('database_get_table_columns', {
      tableName: table.name,
      connectionId: selectedConn.value.id.toString(),
    })

    if (response) {
      tableColumns.value = response.columns || []
    }
  } catch (error: any) {
    message.error('获取列信息失败')
  }
}

const handleExportResults = () => {
  message.info('导出功能开发中')
}

const formatSize = (size: number) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / 1024 / 1024).toFixed(2)} MB`
  return `${(size / 1024 / 1024 / 1024).toFixed(2)} GB`
}
</script>

<style scoped>
.database-manager {
  padding: 20px;
  height: 100%;
  overflow: auto;
}

.connection-list-card,
.query-card,
.tables-card,
.columns-card {
  height: 100%;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

:deep(.n-list-item) {
  cursor: pointer;
  transition: all 0.2s ease;
}

:deep(.n-list-item.active) {
  background: var(--active-color-suppl);
}

:deep(.n-list-item:hover) {
  background: var(--hover-color);
}

.table-card {
  cursor: pointer;
  transition: all 0.2s ease;
}

.table-card.selected {
  background: var(--active-color-suppl);
  border-color: var(--active-color);
}

.table-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}
</style>
