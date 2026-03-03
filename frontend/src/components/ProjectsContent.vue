<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('projects.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('projects.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="projects-content">
        <div class="projects-sidebar">
          <div style="display: flex; justify-content: center; align-items: center; margin-bottom: 16px;">
            <Tooltip :text="t('projects.newProject')">
              <button class="new-project-button" @click="handleNewProject">
                +
              </button>
            </Tooltip>
          </div>
          <div style="width: 100%; height: 1px; background: var(--border-color); margin-bottom: 16px;"></div>
          <div class="directory-tree">
            <div 
              class="tree-item" 
              :class="{ active: selectedProject === t('projects.title') + ' 1' }"
              @click="selectedProject = t('projects.title') + ' 1'"
            >
              {{ t('projects.title') }} 1
            </div>
            <div 
              class="tree-item" 
              :class="{ active: selectedProject === t('projects.title') + ' 2' }"
              @click="selectedProject = t('projects.title') + ' 2'"
            >
              {{ t('projects.title') }} 2
            </div>
            <div 
              class="tree-item" 
              :class="{ active: selectedProject === t('projects.title') + ' 3' }"
              @click="selectedProject = t('projects.title') + ' 3'"
            >
              {{ t('projects.title') }} 3
            </div>
          </div>
        </div>
        <div class="projects-main">
          <NCard class="webshell-table-card">
            <template #header>
              <div style="width: 100%;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
                  <NInput
                    v-model:value="searchQuery"
                    placeholder="搜索 WebShell..."
                    size="small"
                    style="width: 300px; border-radius: 8px;"
                    @input="handleSearch"
                  >
                    <template #prefix>
                      <span style="margin-right: 8px;">🔍</span>
                    </template>
                  </NInput>
                  <NSpace align="center" style="gap: 16px;">
                    <NSpace align="center" style="gap: 8px;">
                      <span style="font-size: 14px; color: var(--text-color);">{{ t('projects.total') }}: <span style="color: var(--active-color); font-weight: 500;">{{ filteredData.length }}</span></span>
                      <span style="font-size: 14px; color: var(--text-color);">{{ t('projects.active') }}: <span style="color: #4CAF50; font-weight: 500;">{{ activeCount }}</span></span>
                      <span style="font-size: 14px; color: var(--text-color);">{{ t('projects.inactive') }}: <span style="color: #FF9800; font-weight: 500;">{{ inactiveCount }}</span></span>
                    </NSpace>
                    <NSelect
                      v-model:value="pageSize"
                      :options="[
                        { label: `5 ${t('projects.itemsPerPage')}`, value: 5 },
                        { label: `10 ${t('projects.itemsPerPage')}`, value: 10 },
                        { label: `20 ${t('projects.itemsPerPage')}`, value: 20 }
                      ]"
                      @update:value="handlePageSizeChange"
                      size="small"
                      style="min-width: 100px"
                    />
                  </NSpace>
                </div>
              </div>
            </template>
            <!-- 完整表格 -->
            <div style="overflow-x: auto; margin-bottom: 16px; background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 8px; padding: 16px;">
              <!-- 表格 -->
              <table id="webshellTable" style="width: 100%; border-collapse: collapse; table-layout: fixed;">
                <thead>
                  <tr style="background: var(--hover-color);">
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 60px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('id')">
                      ID <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('id') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 120px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('name')">
                      {{ t('projects.filename') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('name') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 200px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('url')">
                      {{ t('projects.url') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('url') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 100px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('payload')">
                      {{ t('projects.payloadType') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('payload') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 100px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('cryption')">
                      {{ t('projects.cryption') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('cryption') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 80px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('encoding')">
                      {{ t('projects.encoding') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('encoding') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 100px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('proxyType')">
                      {{ t('projects.proxyType') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('proxyType') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 150px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('remark')">
                      {{ t('projects.remark') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('remark') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 150px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('createTime')">
                      {{ t('projects.createTime') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('createTime') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 150px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('updateTime')">
                      {{ t('projects.updateTime') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('updateTime') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                    <th style="padding: 10px; border: 1px solid var(--border-color); text-align: left; min-width: 80px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('status')">
                      {{ t('projects.status') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('status') }}</span>
                      <div class="resize-handle" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize;"></div>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in sortedData" :key="item.id" 
                    @click="handleTableRowClick(item)"
                    @contextmenu="(event) => handleContextMenu(item, event)"
                    :class="{ 'table-row-selected': selectedTableRow && selectedTableRow.id === item.id }"
                    style="cursor: pointer; transition: background-color 0.2s;"
                  >
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.id }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.name }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.url }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.payload }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.cryption }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.encoding }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.proxyType }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.remark }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.createTime }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.updateTime }}</td>
                    <td style="padding: 10px; border: 1px solid var(--border-color); text-align: left;">{{ item.status }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
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
            :options="menuOptions"
            @select="handleMenuClick"
            placement="bottom"
            :width="150"
            :max-width="200"
            :z-index="1000"
            :bordered="false"
            :show-icon="true"
            :animation="{
              name: 'fade',
              duration: 200
            }"

          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NCard, NPagination, NSpace, NSelect, NDropdown, NIcon, NMenu, NText, NInput } from 'naive-ui'
import Tooltip from './Tooltip.vue'

// 导入 Wails 运行时
import { Events } from '@wailsio/runtime'

// 统一的事件发送函数
const emitEvent = (event: string, data?: any) => {
  console.log('Attempting to emit event:', event, 'with data:', data);
  try {
    console.log('Using Wails 3 Events from @wailsio/runtime');
    Events.Emit(event, data);
  } catch (error) {
    console.error('Error emitting event:', error);
    // 显示详细的错误信息，帮助调试
    alert('Wails 运行时不可用，无法创建原生窗口。请检查 Wails 配置和版本。');
  }
};


const { t } = useI18n()

// 分页状态
const page = ref(1)
const pageSize = ref(5)
const total = ref(12)

// 搜索状态
const searchQuery = ref('')

// 排序状态
const sortField = ref<string>('id')
const sortDirection = ref<'asc' | 'desc'>('asc')

// 右键菜单状态
const menuOptions = computed(() => [
  { 
    label: t('projects.control'), 
    key: 'enter'
  },
  { 
    label: t('projects.cache'), 
    key: 'cache'
  },
  { 
    label: t('projects.edit'), 
    key: 'edit'
  },
  { 
    label: t('projects.delete'), 
    key: 'delete'
  }
])

const selectedRow = ref<any>(null)
const menuVisible = ref(false)
const menuPosition = ref({ x: 0, y: 0 })
const selectedProject = ref('项目 1')
const selectedTableRow = ref<any>(null)

// 处理新建项目
const handleNewProject = () => {
  console.log('新建项目')
  // 这里可以添加新建项目的逻辑
}

// 处理右键菜单
const handleContextMenu = (row: any, event: MouseEvent) => {
  event.preventDefault()
  selectedRow.value = row
  menuPosition.value = { x: event.clientX, y: event.clientY }
  menuVisible.value = true
}

// 处理菜单点击
const handleMenuClick = (key: string) => {
  switch (key) {
    case 'enter':
      console.log('Enter webshell:', selectedRow.value)
      // 发送事件到后端创建Wails 3原生窗口
      emitEvent('createWindow', {
        title: '控制 WebShell',
        width: 800,
        height: 600,
        x: 100,
        y: 100
      })
      break
    case 'cache':
      console.log('Cache webshell:', selectedRow.value)
      // 发送事件到后端创建Wails 3原生窗口
      emitEvent('createWindow', {
        title: '缓存 WebShell',
        width: 800,
        height: 600,
        x: 150,
        y: 150
      })
      break
    case 'edit':
      console.log('Edit webshell:', selectedRow.value)
      // 发送事件到后端创建Wails 3原生窗口
      emitEvent('createWindow', {
        title: '编辑 WebShell',
        width: 800,
        height: 600,
        x: 200,
        y: 200
      })
      break
    case 'delete':
      console.log('Delete webshell:', selectedRow.value)
      break
  }
  menuVisible.value = false
}

// 处理表格行点击
const handleTableRowClick = (item: any) => {
  selectedTableRow.value = item
}

// 处理搜索
const handleSearch = () => {
  page.value = 1 // 搜索时重置到第一页
}

// 处理排序
const handleSort = (field: string) => {
  if (sortField.value === field) {
    // 如果点击的是当前排序字段，切换排序方向
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    // 如果点击的是新字段，设置为升序
    sortField.value = field
    sortDirection.value = 'asc'
  }
}

// 获取排序图标
const getSortIcon = (field: string) => {
  if (sortField.value !== field) {
    return '↑↓'
  }
  return sortDirection.value === 'asc' ? '↑' : '↓'
}

// 处理分页
const handlePageChange = (pageNum: number) => {
  page.value = pageNum
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  page.value = 1
}

// 项目webshell数据
const projectWebshells = ref({
  '项目 1': [
    {
      id: 1,
      name: 'webshell1.php',
      url: 'http://example.com/shell.php',
      payload: 'PHP',
      cryption: 'base64',
      encoding: 'utf-8',
      proxyType: 'http',
      remark: '测试webshell',
      createTime: '2024-01-01 10:00:00',
      updateTime: '2024-01-01 10:00:00',
      status: 'active'
    },
    {
      id: 2,
      name: 'webshell2.asp',
      url: 'http://example.com/shell.asp',
      payload: 'ASP',
      cryption: 'none',
      encoding: 'gb2312',
      proxyType: 'socks5',
      remark: '生产环境webshell',
      createTime: '2024-01-02 11:30:00',
      updateTime: '2024-01-02 11:30:00',
      status: 'inactive'
    }
  ],
  '项目 2': [
    {
      id: 3,
      name: 'webshell3.aspx',
      url: 'http://example.com/shell.aspx',
      payload: 'ASPX',
      cryption: 'aes',
      encoding: 'utf-8',
      proxyType: 'none',
      remark: 'ASP.NET webshell',
      createTime: '2024-01-03 14:20:00',
      updateTime: '2024-01-03 14:20:00',
      status: 'active'
    }
  ],
  '项目 3': [
    {
      id: 4,
      name: 'webshell4.jsp',
      url: 'http://example.com/shell.jsp',
      payload: 'JSP',
      cryption: 'base64',
      encoding: 'utf-8',
      proxyType: 'http',
      remark: 'Java webshell',
      createTime: '2024-01-04 09:15:00',
      updateTime: '2024-01-04 09:15:00',
      status: 'active'
    },
    {
      id: 5,
      name: 'webshell5.php',
      url: 'http://example.com/shell5.php',
      payload: 'PHP',
      cryption: 'none',
      encoding: 'utf-8',
      proxyType: 'none',
      remark: '简单webshell',
      createTime: '2024-01-05 16:45:00',
      updateTime: '2024-01-05 16:45:00',
      status: 'inactive'
    },
    {
      id: 6,
      name: 'webshell6.php',
      url: 'http://example.com/shell6.php',
      payload: 'PHP',
      cryption: 'base64',
      encoding: 'utf-8',
      proxyType: 'http',
      remark: '高级webshell',
      createTime: '2024-01-06 10:30:00',
      updateTime: '2024-01-06 10:30:00',
      status: 'active'
    }
  ]
})

// 获取当前项目的webshell数据
const webshellData = computed(() => {
  return projectWebshells.value[selectedProject.value] || []
})

// 计算过滤后的数据
const filteredData = computed(() => {
  if (!searchQuery.value) {
    return webshellData.value
  }
  const query = searchQuery.value.toLowerCase()
  return webshellData.value.filter(item => {
    return (
      item.id.toString().includes(query) ||
      item.name.toLowerCase().includes(query) ||
      item.url.toLowerCase().includes(query) ||
      item.payload.toLowerCase().includes(query) ||
      item.cryption.toLowerCase().includes(query) ||
      item.encoding.toLowerCase().includes(query) ||
      item.proxyType.toLowerCase().includes(query) ||
      item.remark.toLowerCase().includes(query) ||
      item.status.toLowerCase().includes(query)
    )
  })
})

// 计算活跃和非活跃webshell数量
const activeCount = computed(() => {
  return filteredData.value.filter(item => item.status === 'active').length
})

const inactiveCount = computed(() => {
  return filteredData.value.filter(item => item.status === 'inactive').length
})

// 计算排序后的数据
const sortedData = computed(() => {
  return [...filteredData.value].sort((a, b) => {
    let aValue = a[sortField.value as keyof typeof a]
    let bValue = b[sortField.value as keyof typeof b]
    
    // 处理不同类型的比较
    if (typeof aValue === 'string' && typeof bValue === 'string') {
      return sortDirection.value === 'asc' 
        ? aValue.localeCompare(bValue) 
        : bValue.localeCompare(aValue)
    } else if (typeof aValue === 'number' && typeof bValue === 'number') {
      return sortDirection.value === 'asc' 
        ? aValue - bValue 
        : bValue - aValue
    } else {
      // 默认比较
      return sortDirection.value === 'asc' 
        ? String(aValue).localeCompare(String(bValue)) 
        : String(bValue).localeCompare(String(aValue))
    }
  })
})

// 表格列宽调整功能

let resizing = false
let currentTh: HTMLElement | null = null
let startX = 0
let startWidth = 0

const handleMouseDown = (e: MouseEvent, th: HTMLElement) => {
  if (e.target instanceof HTMLElement && e.target.classList.contains('resize-handle')) {
    resizing = true
    currentTh = th
    startX = e.clientX
    startWidth = th.offsetWidth
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
  }
}

const handleMouseMove = (e: MouseEvent) => {
  if (!resizing || !currentTh) return
  
  const width = startWidth + (e.clientX - startX)
  if (width > 50) { // 最小宽度
    currentTh.style.width = `${width}px`
  }
}

const handleMouseUp = () => {
  resizing = false
  currentTh = null
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}

onMounted(() => {
  const table = document.getElementById('webshellTable')
  if (table) {
    const thElements = table.querySelectorAll('th')
    thElements.forEach(th => {
      th.addEventListener('mousedown', (e) => handleMouseDown(e, th as HTMLElement))
    })
    
    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
  }
  
  // 添加点击空白处关闭菜单的事件
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('contextmenu', handleContextMenuOutside)
})

onUnmounted(() => {
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', handleMouseUp)
  
  // 移除事件监听器
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('contextmenu', handleContextMenuOutside)
})

// 处理点击空白处关闭菜单
const handleClickOutside = (event: MouseEvent) => {
  // 只有在左键点击时才关闭菜单，避免右键点击时关闭
  if (event.button === 0) {
    menuVisible.value = false
  }
}

// 处理右键点击空白处关闭菜单
const handleContextMenuOutside = (event: MouseEvent) => {
  // 只有在菜单可见时才处理，并且点击不是来自表格行
  if (menuVisible.value) {
    // 检查点击目标是否是表格行或其子元素
    const target = event.target as HTMLElement
    const isTableRow = target.closest('tr')
    if (!isTableRow) {
      menuVisible.value = false
    }
  }
}

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

/* 新建项目按钮样式 */
.new-project-button {
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 8px;
  width: 60px;
  height: 36px;
  font-size: 20px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin: 0 auto;
}

.new-project-button:hover {
  opacity: 0.9;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  scale: 1.05;
}

.projects-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--content-bg);
  padding: 20px;
  overflow: auto;
  height: 100%;
  box-sizing: border-box;
  min-height: 400px;
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
  flex: 1;
  min-height: 350px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding-top: 16px;
}

.table-row-selected {
  background-color: rgba(99, 102, 241, 0.1) !important;
  border-left: 4px solid var(--active-color) !important;
  box-shadow: 0 0 0 1px var(--active-color) !important;
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

/* 右键菜单样式 */
.menu-item-enter {
  color: var(--active-color);
}

.menu-item-cache {
  color: #64748b;
}

.menu-item-edit {
  color: #3b82f6;
}

.menu-item-delete {
  color: #ef4444;
}

/* 菜单容器样式 */
.n-dropdown {
  animation: menuFadeIn 0.2s ease-out;
  background-color: var(--card-bg) !important;
  border: 1px solid var(--border-color) !important;
  border-radius: 8px !important;
  box-shadow: var(--shadow-md) !important;
}

@keyframes menuFadeIn {
  from {
    opacity: 0;
    transform: translateY(-5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 菜单项样式 */
.n-menu-item {
  background-color: transparent !important;
  color: var(--text-color) !important;
  padding: 8px 12px !important;
  transition: all 0.2s ease;
  border-radius: 4px !important;
  margin: 2px 4px !important;
}

/* 菜单项悬停效果 */
.n-menu-item:hover {
  background-color: var(--hover-color) !important;
  color: var(--text-color) !important;
}

/* 菜单项选中效果 */
.n-menu-item.n-menu-item--selected {
  background-color: var(--active-color) !important;
  color: white !important;
}

/* 菜单项图标样式 */
.n-menu-item-icon {
  margin-right: 10px;
  font-size: 14px;
  width: 16px;
  text-align: center;
}

/* 菜单项文字样式 */
.n-menu-item-content {
  font-size: 14px !important;
  font-weight: 500 !important;
}
</style>