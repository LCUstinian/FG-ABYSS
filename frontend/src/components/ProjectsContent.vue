<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('projects.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('projects.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="projects-content">
        <div class="projects-sidebar">
          <h3>{{ t('projects.projectDir') }}</h3>
          <div class="directory-tree">
            <div class="tree-item">{{ t('projects.title') }} 1</div>
            <div class="tree-item">{{ t('projects.title') }} 2</div>
            <div class="tree-item">{{ t('projects.title') }} 3</div>
          </div>
        </div>
        <div class="projects-main">
          <h3>{{ t('projects.webshellList') }}</h3>
          <NCard class="webshell-table-card">
            <template #header>
              <NSpace justify="space-between" align="center">
                <NText>{{ t('projects.webshellList') }}</NText>
                <NSelect
                  v-model:value="pageSize"
                  :options="[
                    { label: `5 ${t('projects.itemsPerPage')}`, value: 5 },
                    { label: `10 ${t('projects.itemsPerPage')}`, value: 10 },
                    { label: `20 ${t('projects.itemsPerPage')}`, value: 20 }
                  ]"
                  @update:value="handlePageSizeChange"
                  size="small"
                />
              </NSpace>
            </template>
            <NTable
              :columns="columns"
              :data="currentPageData"
              @sort="handleSort"
              @contextmenu="handleContextMenu"
              row-key="id"
            />
            <template #footer>
              <div class="pagination-container">
                <NPagination
                  v-model:page="page"
                  v-model:page-size="pageSize"
                  :page-sizes="[5, 10, 20]"
                  :item-count="total"
                  @update:page="handlePageChange"
                  @update:page-size="handlePageSizeChange"
                />
              </div>
            </template>
          </NCard>
          <NDropdown
            v-model:show="menuVisible"
            :x="menuPosition.x"
            :y="menuPosition.y"
            trigger="manual"
          >
            <NMenu
              :options="menuOptions"
              @select="handleMenuClick"
            />
          </NDropdown>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { NTable, NButton, NCard, NPagination, NSpace, NSelect, NDropdown, NIcon, NMenu, NText } from 'naive-ui'

const { t } = useI18n()

// 分页状态
const page = ref(1)
const pageSize = ref(5)
const total = ref(12)

// 排序状态
const sortState = ref({
  columnKey: 'date',
  order: 'descend' as 'ascend' | 'descend' | null
})

// 右键菜单状态
const menuOptions = computed(() => [
  { label: t('projects.control'), key: 'control' },
  { label: t('projects.cache'), key: 'cache' },
  { label: t('projects.edit'), key: 'edit' },
  { label: t('projects.delete'), key: 'delete' }
])

const selectedRow = ref<any>(null)
const menuVisible = ref(false)
const menuPosition = ref({ x: 0, y: 0 })

// WebShell列表数据（模拟）
const webshellData = ref([
  { id: 1, name: 'webshell1.php', ip: '192.168.1.100', date: '2024-01-01', status: 'active' },
  { id: 2, name: 'webshell2.asp', ip: '192.168.1.101', date: '2024-01-02', status: 'inactive' },
  { id: 3, name: 'webshell3.aspx', ip: '192.168.1.102', date: '2024-01-03', status: 'active' },
  { id: 4, name: 'webshell4.jsp', ip: '192.168.1.103', date: '2024-01-04', status: 'active' },
  { id: 5, name: 'webshell5.php', ip: '192.168.1.104', date: '2024-01-05', status: 'inactive' },
  { id: 6, name: 'webshell6.asp', ip: '192.168.1.105', date: '2024-01-06', status: 'active' },
  { id: 7, name: 'webshell7.aspx', ip: '192.168.1.106', date: '2024-01-07', status: 'active' },
  { id: 8, name: 'webshell8.jsp', ip: '192.168.1.107', date: '2024-01-08', status: 'inactive' },
  { id: 9, name: 'webshell9.php', ip: '192.168.1.108', date: '2024-01-09', status: 'active' },
  { id: 10, name: 'webshell10.asp', ip: '192.168.1.109', date: '2024-01-10', status: 'active' },
  { id: 11, name: 'webshell11.aspx', ip: '192.168.1.110', date: '2024-01-11', status: 'inactive' },
  { id: 12, name: 'webshell12.jsp', ip: '192.168.1.111', date: '2024-01-12', status: 'active' }
])

// 处理右键菜单
const handleContextMenu = (row: any, event: MouseEvent) => {
  event.preventDefault()
  selectedRow.value = row
  menuPosition.value = { x: event.clientX, y: event.clientY }
  menuVisible.value = true
}

// 处理菜单点击
const handleMenuClick = (key: string) => {
  console.log('Menu clicked:', key, 'for row:', selectedRow.value)
  menuVisible.value = false
}

// 处理排序
const handleSort = (columnKey: string, order: 'ascend' | 'descend' | null) => {
  sortState.value = { columnKey, order }
  // 这里可以实现实际的排序逻辑
  if (order) {
    webshellData.value.sort((a, b) => {
      if (a[columnKey] < b[columnKey]) return order === 'ascend' ? -1 : 1
      if (a[columnKey] > b[columnKey]) return order === 'ascend' ? 1 : -1
      return 0
    })
  }
}

// 处理分页
const handlePageChange = (pageNum: number) => {
  page.value = pageNum
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  page.value = 1
}

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (page.value - 1) * pageSize.value
  const end = start + pageSize.value
  return webshellData.value.slice(start, end)
})

// 表格列配置
const columns = computed(() => [
  {
    title: t('projects.filename'),
    key: 'name',
    sortable: true,
    sorter: (a: any, b: any) => a.name.localeCompare(b.name)
  },
  {
    title: t('projects.ipAddress'),
    key: 'ip',
    sortable: true,
    sorter: (a: any, b: any) => a.ip.localeCompare(b.ip)
  },
  {
    title: t('projects.createDate'),
    key: 'date',
    sortable: true,
    sorter: (a: any, b: any) => new Date(a.date).getTime() - new Date(b.date).getTime()
  },
  {
    title: t('projects.status'),
    key: 'status',
    sortable: true,
    render: (row: any) => {
      return h('span', {}, row.status === 'active' ? t('projects.active') : t('projects.inactive'))
    }
  },
  {
    title: t('projects.action'),
    key: 'action',
    render: (row: any) => {
      return h(NButton, {
        type: 'primary',
        size: 'small',
        onClick: () => console.log('Enter webshell:', row)
      }, {
        default: () => t('projects.enter')
      })
    }
  }
])
</script>

<style scoped>
.content-section {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

.content-header {
  width: 100%;
  padding: 16px 24px;
  margin-bottom: 0;
  background: var(--panel-bg);
  border-bottom: none;
  box-shadow: none;
}

.content-section h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-color);
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
  line-height: 1.2;
}

.content-body {
  flex: 1;
  width: 100%;
  padding: 0;
  margin: 0;
  box-sizing: border-box;
  background: var(--content-bg);
  border-top: none;
  display: flex;
  align-items: stretch;
}

.content-section h1 .title {
  font-weight: 700;
  color: white;
}

.content-section h1 .separator {
  color: white;
  opacity: 0.7;
  font-weight: 400;
}

.content-section h1 .subtitle {
  font-size: 14px;
  font-weight: 400;
  color: white;
  opacity: 0.8;
  font-style: normal;
}

/* 项目内容样式 - 深色主题风格 */
.projects-content {
  display: flex;
  gap: 1px;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  background: var(--border-color);
}

.projects-sidebar {
  width: 200px;
  background: var(--sidebar-bg);
  padding: 20px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow-y: auto;
  height: 100%;
  box-sizing: border-box;
}

.directory-tree {
  margin-top: 16px;
}

.tree-item {
  padding: 10px 12px;
  cursor: pointer;
  margin-bottom: 4px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  border-radius: 4px;
}

.tree-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
}

.tree-item.active {
  background: var(--active-color);
  color: white;
  box-shadow: var(--shadow-sm);
}

.projects-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--content-bg);
  padding: 20px;
  overflow-y: auto;
  height: 100%;
  box-sizing: border-box;
}

.projects-main h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
}

.webshell-table-card {
  margin-bottom: 16px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding-top: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-header {
    padding: 12px 16px;
  }
  
  .content-body {
    padding: 16px;
  }
  
  .projects-content {
    flex-direction: column;
  }
  
  .projects-sidebar {
    width: 100%;
    height: auto;
    max-height: 150px;
  }
}
</style>