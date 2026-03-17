<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('projects.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('projects.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="projects-content">
        <div class="projects-sidebar">
          <!-- 顶部新建按钮区域 -->
          <div class="sidebar-top-section">
            <Tooltip :text="t('projects.newProject')" :always-show="true">
              <button class="new-project-button" @click="handleNewProject">
                +
              </button>
            </Tooltip>
          </div>
          <div style="width: 100%; height: 1px; background: var(--border-color); margin-bottom: 16px;"></div>
          <!-- 项目列表区域 -->
          <div class="directory-tree">
            <!-- 正常项目列表 -->
            <div 
              v-for="project in projects" 
              :key="project.id"
              class="tree-item" 
              :class="{ active: selectedProject === project.id }"
              @click="handleProjectSelect(project.id)"
              @mouseenter="hoveredProject = project.id"
              @mouseleave="hoveredProject = null"
            >
              <span class="tree-item-icon">📁</span>
              <Tooltip :text="project.name">
                <span class="tree-item-text">{{ project.name }}</span>
              </Tooltip>
              <span class="tree-item-actions" :class="{ 'show-actions': hoveredProject === project.id || selectedProject === project.id }">
                <button 
                  class="delete-project-btn" 
                  @click.stop="handleDeleteProject(project)"
                  :title="t('projects.deleteProject')"
                >
                  🗑️
                </button>
              </span>
            </div>
          </div>
          
          <!-- 底部回收站区域 - 固定于底部 -->
          <div class="sidebar-bottom-section">
            <div class="recycle-bin-section">
              <div class="recycle-bin-divider">
                <span class="divider-line"></span>
                <span class="divider-text">{{ t('projects.recycleBin') }}</span>
                <span class="divider-line"></span>
              </div>
              <!-- 统一的恢复按钮 - 始终显示 -->
              <button class="recover-all-btn" @click="showRecoverDialog = true">
                <span class="btn-icon">↩️</span>
                <span class="btn-text">{{ t('projects.recoverProjects') }}</span>
                <span class="btn-count" v-if="deletedProjects.length > 0">{{ deletedProjects.length }}</span>
              </button>
            </div>
          </div>
        </div>
        
        <!-- 恢复项目弹窗组件 -->
        <RecoverProjectModal
          v-model="showRecoverDialog"
          :deleted-projects="deletedProjects"
          @recover="handleRecoverProject"
          @close="handleRecoverDialogClose"
        />
        <div class="projects-main">
          <NCard class="webshell-table-card" style="height: 100%; display: flex; flex-direction: column; overflow: hidden;">
            <template #header>
              <div style="width: 100%; flex-shrink: 0;">
                <div style="display: flex; justify-content: space-between; align-items: center; padding: 16px 16px 0 16px;">
                  <NInput
                    v-model:value="searchQuery"
                    placeholder="搜索 WebShell..."
                    size="small"
                    class="search-input"
                    style="width: 300px;"
                    @input="handleSearch"
                  >
                    <template #prefix>
                      <span style="margin-right: 8px;">🔍</span>
                    </template>
                  </NInput>
                  <div class="toolbar-container">
                    <span style="font-size: 14px; color: var(--text-color);">{{ t('projects.total') }}: <span style="color: var(--active-color); font-weight: 500;">{{ total }}</span></span>
                    <span style="font-size: 14px; color: var(--text-color);">{{ t('projects.active') }}: <span style="color: #4CAF50; font-weight: 500;">{{ activeCount }}</span></span>
                    <span style="font-size: 14px; color: var(--text-color);">{{ t('projects.inactive') }}: <span style="color: #FF9800; font-weight: 500;">{{ inactiveCount }}</span></span>
                    <!-- 回收站切换按钮 -->
                    <NButton 
                      :type="showDeleted ? 'warning' : 'default'"
                      size="small"
                      @click="toggleView"
                      class="recycle-bin-btn"
                    >
                      <template #icon>
                        <span class="recycle-icon">{{ showDeleted ? '🗑️' : '♻️' }}</span>
                      </template>
                      <span class="recycle-text">{{ showDeleted ? t('projects.showNormal') : t('projects.showDeleted') }}</span>
                    </NButton>
                    <Tooltip :text="t('projects.newWebShell')" :always-show="true" placement="bottom">
                      <NButton 
                        type="primary" 
                        size="small"
                        @click="handleNewWebShell"
                        class="new-webshell-btn"
                      >
                        <template #icon>
                          <span class="btn-icon">+</span>
                        </template>
                      </NButton>
                    </Tooltip>
                    <!-- 每页条数选择器容器 -->
                    <div class="page-size-container">
                      <NSelect
                        v-model:value="pageSize"
                        :options="pageSizeOptions.map(size => ({
                          label: `${size} ${t('projects.itemsPerPage')}`,
                          value: size
                        }))"
                        @update:value="handlePageSizeChange"
                        size="small"
                        class="page-size-select"
                        :option-props="{
                          class: 'custom-option'
                        }"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </template>
            <!-- 完整表格 -->
            <div class="webshell-table-container">
              <!-- 空状态提示 -->
              <div v-if="tableData.length === 0" style="text-align: center; padding: 60px 20px; color: var(--text-color); opacity: 0.6;">
                <div style="font-size: 48px; margin-bottom: 16px;">
                  {{ showDeleted ? '🗑️' : '📦' }}
                </div>
                <div style="font-size: 16px; margin-bottom: 8px;">
                  {{ showDeleted ? t('projects.noDeletedData') : t('projects.noData') }}
                </div>
                <div style="font-size: 14px; opacity: 0.8;">
                  {{ showDeleted ? t('projects.noDeletedDataTip') : t('projects.noDataTip') }}
                </div>
              </div>
              
              <!-- 表格容器 -->
              <div v-else style="flex: 1; overflow: auto; min-height: 0;">
                <table id="webshellTable" class="webshell-table" style="width: 100%; border-collapse: collapse; table-layout: fixed;">
                  <thead>
                  <tr class="webshell-table-header-row">
                    <th class="webshell-table-header" style="text-align: left; min-width: 60px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('id')">
                      ID <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('id') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'id')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 200px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('url')">
                      {{ t('projects.url') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('url') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'url')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 100px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('payload')">
                      {{ t('projects.payloadType') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('payload') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'payload')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 100px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('cryption')">
                      {{ t('projects.cryption') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('cryption') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'cryption')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 80px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('encoding')">
                      {{ t('projects.encoding') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('encoding') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'encoding')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 100px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('proxyType')">
                      {{ t('projects.proxyType') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('proxyType') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'proxyType')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 150px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('remark')">
                      {{ t('projects.remark') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('remark') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'remark')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 150px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('createdAt')">
                      {{ t('projects.createTime') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('createdAt') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'createdAt')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 150px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('updatedAt')">
                      {{ t('projects.updateTime') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('updatedAt') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'updatedAt')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                    <th class="webshell-table-header" style="text-align: left; min-width: 80px; cursor: pointer; user-select: none; position: relative;" @click="handleSort('status')">
                      {{ t('projects.status') }} <span style="font-size: 10px; margin-left: 4px;">{{ getSortIcon('status') }}</span>
                      <div class="resize-handle" @mousedown.stop="handleResizeStart($event, 'status')" style="position: absolute; right: 0; top: 0; bottom: 0; width: 5px; cursor: col-resize; z-index: 10;"></div>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in tableData" :key="item.id" 
                    @click="handleTableRowClick(item)"
                    @contextmenu="(event) => handleContextMenu(item, event)"
                    :class="{ 'table-row-selected': selectedTableRow && selectedTableRow.id === item.id }"
                    class="webshell-table-row"
                  >
                    <td class="webshell-table-cell">
                      <Tooltip :text="item.id">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.id }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="item.url">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.url }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="item.payload">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.payload }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="item.cryption">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.cryption }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="item.encoding">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.encoding }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="item.proxyType">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.proxyType }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip v-if="item.remark" :text="item.remark">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.remark }}</span>
                      </Tooltip>
                      <span v-else style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">-</span>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="formatTime(item.createdAt)">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ formatTime(item.createdAt) }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="formatTime(item.updatedAt)">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ formatTime(item.updatedAt) }}</span>
                      </Tooltip>
                    </td>
                    <td class="webshell-table-cell">
                      <Tooltip :text="item.status">
                        <span style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ item.status }}</span>
                      </Tooltip>
                    </td>
                  </tr>
                </tbody>
                </table>
              </div>
            </div>
            <template #footer>
              <div class="pagination-container">
                <NPagination
                  v-model:page="page"
                  :page-size="pageSize"
                  :item-count="total"
                  @update:page="handlePageChange"
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
    <!-- 新建项目弹窗 -->
    <CreateProjectModal 
      v-model="newProjectDialogVisible" 
      @created="handleProjectCreated"
    />
    
    <!-- 新建 WebShell 弹窗 -->
    <CreateWebShellModal 
      v-model="newWebShellDialogVisible" 
      :project-id="selectedProject || undefined"
      @created="handleWebShellCreated"
    />

  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { 
  NButton, 
  NCard, 
  NPagination, 
  NSelect,
  NSpace, 
  NDropdown, 
  NIcon, 
  NMenu, 
  NText, 
  NInput,
  NDialogProvider,
  useDialog,
  useMessage
} from 'naive-ui'
import Tooltip from './Tooltip.vue'
import CreateProjectModal from './CreateProjectModal.vue'
import CreateWebShellModal from './CreateWebShellModal.vue'
import RecoverProjectModal from './RecoverProjectModal.vue'
import { useSmartPagination } from '@/composables/useSmartPagination'

// 导入 Wails 运行时和绑定
import { Events } from '@wailsio/runtime'
import * as ProjectHandler from '../../bindings/fg-abyss/internal/app/handlers/projecthandler'
import * as WebShellHandler from '../../bindings/fg-abyss/internal/app/handlers/webshellhandler'

// 导入时间格式化工具
import { formatTime } from '@/utils/formatTime'

// 导入工具函数
import { componentLogger } from '@/utils/logger'
import type { Project as ProjectType } from '@/types'
import type { Project } from '@bindings/fg-abyss/internal/domain/entity/models'

// 定义 WebShell 接口
interface WebShell {
  id: string
  projectId: string
  url: string
  payload: string
  cryption: string
  encoding: string
  proxyType: string
  remark: string
  createdAt: string
  updatedAt: string
  status: string
}

// 统一的事件发送函数
const emitEvent = (event: string, data?: any) => {
  try {
    componentLogger.debug('发送事件:', event)
    Events.Emit(event, data);
  } catch (error) {
    componentLogger.error('事件发送失败:', event, error)
    message.error('Wails 运行时不可用')
  }
};


const { t } = useI18n()

// 使用 Naive UI 的对话框和消息组件
const dialog = useDialog()
const message = useMessage()

// 智能分页状态
const total = ref(0)

// 使用智能分页组合函数
const {
  page,
  pageSize,
  isManualMode,
  pageSizeOptions,
  setPageSize,
  setPage,
  resetToAuto
} = useSmartPagination({
  total: () => total.value,
  enableAuto: false, // 禁用自动调整，使用默认分页大小
  onPageSizeChange: (size: number) => {
    console.log('分页大小调整为:', size)
    // 分页大小变化时，重置到第一页并重新获取数据
    page.value = 1
    fetchData()
  },
  onPageChange: () => {
    // 页码变化时只获取数据，不重置分页
    fetchData()
  }
})

// 搜索状态
const searchQuery = ref('')

// 排序状态
const sortField = ref<string>('id')
const sortDirection = ref<'asc' | 'desc'>('asc')

// 回收站视图状态
const showDeleted = ref(false)

// 切换视图（正常/回收站）
const toggleView = () => {
  showDeleted.value = !showDeleted.value
  page.value = 1
  fetchData()
}

// 右键菜单状态
const menuOptions = computed(() => {
  const baseOptions = [
    { 
      label: t('projects.control'),
      key: 'enter',
      icon: () => h('span', { class: 'menu-icon-enter' }, '🖥️')
    },
    { 
      label: t('projects.cache'),
      key: 'cache',
      icon: () => h('span', { class: 'menu-icon-cache' }, '💾')
    },
    { 
      label: t('projects.edit'),
      key: 'edit',
      icon: () => h('span', { class: 'menu-icon-edit' }, '✏️')
    }
  ]
  
  // 根据视图显示不同的操作
  if (showDeleted.value) {
    // 回收站视图：显示恢复按钮
    baseOptions.push({
      label: t('projects.recover'),
      key: 'recover',
      icon: () => h('span', { class: 'menu-icon-recover' }, '↩️')
    })
  } else {
    // 正常视图：显示删除按钮
    baseOptions.push({
      label: t('projects.delete'),
      key: 'delete',
      icon: () => h('span', { class: 'menu-icon-delete' }, '🗑️')
    })
  }
  
  return baseOptions
})

const selectedRow = ref<WebShell | null>(null)
const menuVisible = ref(false)
const menuPosition = ref({ x: 0, y: 0 })
const selectedProject = ref<string | null>(null)  // 存储项目 ID，初始为 null
const selectedTableRow = ref<WebShell | null>(null)

// 项目列表
const projects = ref<Project[]>([])
const deletedProjects = ref<Project[]>([])  // 前端临时存储已删除的项目

// UI 状态
const hoveredProject = ref<string | null>(null)
const showRecoverDialog = ref(false)

// 新建项目弹窗
const newProjectDialogVisible = ref(false)

// 新建 WebShell 弹窗
const newWebShellDialogVisible = ref(false)

// 表格数据
const tableData = ref<WebShell[]>([])

// 计算活跃和非活跃 webshell 数量
const activeCount = computed(() => {
  return tableData.value.filter(item => item.status === 'active').length
})

const inactiveCount = computed(() => {
  return tableData.value.filter(item => item.status === 'inactive').length
})

// 处理项目选择（左侧列表点击）
const handleProjectSelect = (projectId: string) => {
  console.log('选择项目:', projectId)
  selectedProject.value = projectId
  // 重置分页和搜索
  page.value = 1
  searchQuery.value = ''
  // 刷新数据
  fetchData()
}

// 处理项目切换
const handleProjectChange = (projectId: string) => {
  console.log('切换项目:', projectId)
  // 重置分页和搜索
  page.value = 1
  searchQuery.value = ''
  // 刷新数据
  fetchData()
}

// 处理新建项目
const handleNewProject = () => {
  newProjectDialogVisible.value = true
}

// 处理新建 WebShell
const handleNewWebShell = () => {
  newWebShellDialogVisible.value = true
}

// 处理 WebShell 创建成功
const handleWebShellCreated = async () => {
  await fetchData()
}

// 处理项目创建成功
const handleProjectCreated = async () => {
  await fetchProjects()
  newProjectDialogVisible.value = false
}

// 格式化时间（使用导入的工具函数）
// const formatTime = (time: string | number | Date) => {
//   if (!time) return ''
//   const date = new Date(time)
//   return date.toLocaleString('zh-CN', {
//     year: 'numeric',
//     month: '2-digit',
//     day: '2-digit',
//     hour: '2-digit',
//     minute: '2-digit',
//     second: '2-digit'
//   })
// }

// 获取项目列表
const fetchProjects = async () => {
  try {
    const projectList = await ProjectHandler.GetProjects()
    
    // 按创建时间排序（从旧到新）
    projects.value = projectList.sort((a, b) => {
      const timeA = new Date(a.createdAt).getTime()
      const timeB = new Date(b.createdAt).getTime()
      return timeA - timeB
    })
    
    console.log('获取项目列表:', projects.value)
    
    // 如果没有选中项目，且项目列表不为空，选择第一个项目
    if (!selectedProject.value && projects.value.length > 0) {
      selectedProject.value = projects.value[0].id
      console.log('自动选择第一个项目:', selectedProject.value)
    }
  } catch (error) {
    console.error('获取项目列表失败:', error)
  }
}

// 获取已删除项目列表
const fetchDeletedProjects = async () => {
  try {
    const deletedList = await ProjectHandler.GetDeletedProjects()
    
    // 按删除时间排序（从新到旧）
    deletedProjects.value = deletedList.sort((a, b) => {
      const timeA = new Date(a.deletedAt).getTime()
      const timeB = new Date(b.deletedAt).getTime()
      return timeB - timeA
    })
    
    console.log('获取已删除项目列表:', deletedProjects.value)
  } catch (error) {
    console.error('获取已删除项目列表失败:', error)
  }
}

// 处理删除项目
const handleDeleteProject = async (project: any) => {
  // 显示确认对话框
  dialog.warning({
    title: t('projects.deleteProjectConfirm'),
    content: t('projects.deleteProjectConfirmContent', { name: project.name }),
    positiveText: t('projects.confirm'),
    negativeText: t('projects.cancel'),
    onPositiveClick: async () => {
      // 显示加载状态
      const loading = message.loading(t('projects.deleting'), {
        duration: 0
      })
      
      try {
        // 调用后端删除接口
        await ProjectHandler.DeleteProject(project.id)
        
        // 关闭加载提示
        loading.destroy()
        
        // 显示成功消息
        message.success(t('projects.deleteSuccess'))
        
        // 刷新项目列表和已删除项目列表
        await fetchProjects()
        await fetchDeletedProjects()
        
        // 如果删除的是当前选中的项目，清空选中
        if (selectedProject.value === project.id) {
          selectedProject.value = null
          tableData.value = []
          total.value = 0
        }
      } catch (error: any) {
        // 关闭加载提示
        loading.destroy()
        
        // 显示错误消息
        console.error('删除失败:', error)
        let errorMessage = ''
        if (error.message) {
          if (typeof error.message === 'string') {
            errorMessage = error.message
          } else if (error.message.message) {
            errorMessage = error.message.message
          } else {
            errorMessage = String(error.message)
          }
        } else {
          errorMessage = String(error)
        }
        message.error(t('projects.deleteError') + ': ' + errorMessage)
      }
    },
    onNegativeClick: () => {
      console.log('取消删除')
    }
  })
}

// 处理恢复项目
const handleRecoverProject = async (project: any) => {
  // 显示加载状态
  const loading = message.loading(t('projects.recovering'), {
    duration: 0
  })
  
  try {
    // 调用后端恢复接口
    await ProjectHandler.RecoverProject(project.id)
    
    // 关闭加载提示
    loading.destroy()
    
    // 显示成功消息
    message.success(t('projects.recoverSuccess'))
    
    // 刷新项目列表和已删除项目列表
    await fetchProjects()
    await fetchDeletedProjects()
  } catch (error: any) {
    // 关闭加载提示
    loading.destroy()
    
    // 显示错误消息
    console.error('恢复失败:', error)
    let errorMessage = ''
    if (error.message) {
      if (typeof error.message === 'string') {
        errorMessage = error.message
      } else if (error.message.message) {
        errorMessage = error.message.message
      } else {
        errorMessage = String(error.message)
      }
    } else {
      errorMessage = String(error)
    }
    message.error(t('projects.recoverError') + ': ' + errorMessage)
  }
}

// 处理恢复全部（由弹窗内部处理，此处只需刷新列表）
const handleRecoverAll = async () => {
  await fetchProjects()
  await fetchDeletedProjects()
}

// 处理关闭弹窗
const handleRecoverDialogClose = () => {
  console.log('关闭恢复项目弹窗')
}

// 处理右键菜单
const handleContextMenu = (row: WebShell, event: MouseEvent) => {
  event.preventDefault()
  selectedRow.value = row
  menuPosition.value = { x: event.clientX, y: event.clientY }
  menuVisible.value = true
}

// 打开 WebShell 控制窗口
const openWebShellControlWindow = async (webshell: WebShell) => {
  try {
    // 使用 Wails API 打开新窗口
    // @ts-ignore
    if (window.runtime) {
      // 发送事件到后端，请求打开新窗口
      // @ts-ignore
      await window.runtime.EventsEmit('open-webshell-window', {
        id: webshell.id,
        name: webshell.name || webshell.url,
        url: webshell.url,
      })
      
      message.success(`已打开 WebShell 控制窗口：${webshell.name || webshell.url}`)
    } else {
      // 浏览器环境下，使用 window.open 作为降级方案
      // 使用完整路径
      const baseUrl = window.location.origin + window.location.pathname
      const url = `${baseUrl}#/webshell-control?id=${webshell.id}`
      
      const newWindow = window.open(url, '_blank', 'width=1200,height=800,resizable=yes,scrollbars=yes')
      
      if (newWindow) {
        message.success('新窗口已打开（浏览器模式）')
      } else {
        message.warning('浏览器阻止了弹窗，请允许弹窗后重试')
      }
    }
  } catch (error: any) {
    console.error('打开控制窗口失败:', error)
    message.error('打开控制窗口失败：' + (error.message || '未知错误'))
  }
}

// 处理菜单点击
const handleMenuClick = (key: string) => {
  if (!selectedRow.value) return
  
  switch (key) {
    case 'enter':
      console.log('Control webshell:', selectedRow.value)
      // 打开 WebShell 控制窗口
      openWebShellControlWindow(selectedRow.value)
      break
    case 'cache':
      console.log('Cache webshell:', selectedRow.value)
      // TODO: 实现缓存功能
      message.info('缓存功能尚未实现')
      break
    case 'edit':
      console.log('Edit webshell:', selectedRow.value)
      // TODO: 打开编辑窗口
      message.info('编辑功能尚未实现')
      break
    case 'delete':
      console.log('Delete webshell:', selectedRow.value)
      // 显示确认对话框
      dialog.warning({
        title: t('projects.deleteConfirm'),
        content: t('projects.deleteConfirmContent'),
        positiveText: t('projects.confirm'),
        negativeText: t('projects.cancel'),
        onPositiveClick: async () => {
          // 检查 selectedRow 是否为 null
          if (!selectedRow.value) {
            message.error('未选择要删除的 WebShell')
            return
          }
          
          // 显示加载状态
          const loading = message.loading('正在删除...', {
            duration: 0
          })
          
          try {
            // 调用后端删除方法
            await WebShellHandler.DeleteWebShell(selectedRow.value.id)
            
            // 关闭加载提示
            loading.destroy()
            
            // 显示成功消息
            message.success(t('projects.deleteSuccess'))
            
            // 刷新数据
            await fetchData()
            
            // 清空选中项
            selectedRow.value = null
          } catch (error: any) {
            // 关闭加载提示
            loading.destroy()
            
            // 显示错误消息
            console.error('删除失败:', error)
            let errorMessage = ''
            if (error.message) {
              if (typeof error.message === 'string') {
                errorMessage = error.message
              } else if (error.message.message) {
                errorMessage = error.message.message
              } else {
                errorMessage = String(error.message)
              }
            } else {
              errorMessage = String(error)
            }
            message.error(t('projects.deleteError') + ': ' + errorMessage)
          }
        },
        onNegativeClick: () => {
          console.log('取消删除')
        }
      })
      break
    case 'recover':
      console.log('Recover webshell:', selectedRow.value)
      // 显示确认对话框
      dialog.success({
        title: t('projects.recoverConfirm'),
        content: t('projects.recoverConfirmContent'),
        positiveText: t('projects.confirm'),
        negativeText: t('projects.cancel'),
        onPositiveClick: async () => {
          // 检查 selectedRow 是否为 null
          if (!selectedRow.value) {
            message.error('未选择要恢复的 WebShell')
            return
          }
          
          // 显示加载状态
          const loading = message.loading('正在恢复...', {
            duration: 0
          })
          
          try {
            // 调用后端恢复方法
            await WebShellHandler.RecoverWebShell(selectedRow.value.id)
            
            // 关闭加载提示
            loading.destroy()
            
            // 显示成功消息
            message.success(t('projects.recoverSuccess'))
            
            // 刷新数据
            await fetchData()
            
            // 清空选中项
            selectedRow.value = null
          } catch (error: any) {
            // 关闭加载提示
            loading.destroy()
            
            // 显示错误消息
            console.error('恢复失败:', error)
            let errorMessage = ''
            if (error.message) {
              if (typeof error.message === 'string') {
                errorMessage = error.message
              } else if (error.message.message) {
                errorMessage = error.message.message
              } else {
                errorMessage = String(error.message)
              }
            } else {
              errorMessage = String(error)
            }
            message.error(t('projects.recoverError') + ': ' + errorMessage)
          }
        },
        onNegativeClick: () => {
          console.log('取消恢复')
        }
      })
      break
  }
  menuVisible.value = false
}

// 处理表格行点击
const handleTableRowClick = (item: WebShell) => {
  selectedTableRow.value = item
}

// 处理搜索
const handleSearch = () => {
  page.value = 1 // 搜索时重置到第一页
  fetchData()
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
  fetchData()
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
  setPage(pageNum)
}

const handlePageSizeChange = (size: number) => {
  setPageSize(size)
}

// 从后端获取数据
const fetchData = async () => {
  try {
    let data: WebShell[]
    let count: number
    
    // 如果没有选中项目，返回空数据
    if (!selectedProject.value) {
      console.log('没有选中的项目，返回空数据')
      tableData.value = []
      total.value = 0
      return
    }
    
    console.log('获取项目数据，项目 ID:', selectedProject.value, '回收站视图:', showDeleted.value)
    
    if (showDeleted.value) {
      // 获取已删除的 WebShell（回收站）
      data = await WebShellHandler.GetDeletedWebShells(selectedProject.value)
      count = data.length
      console.log('获取回收站数据成功，数量:', count)
    } else {
      // 获取正常的 WebShell
      const [result, totalCount] = await WebShellHandler.GetWebShells(
        selectedProject.value,
        page.value,
        pageSize.value,
        searchQuery.value,
        sortField.value,
        sortDirection.value
      )
      data = result
      count = totalCount
      console.log('获取 WebShell 数据成功，数量:', count)
    }
    
    tableData.value = data
    total.value = count
  } catch (error: any) {
    console.error('获取数据失败:', error)
    // 显示错误消息
    let errorMessage = ''
    if (error.message) {
      if (typeof error.message === 'string') {
        errorMessage = error.message
      } else if (error.message.message) {
        errorMessage = error.message.message
      } else {
        errorMessage = String(error.message)
      }
    } else {
      errorMessage = String(error)
    }
    message.error('获取数据失败：' + errorMessage)
  }
}

// 监听项目变化
watch(selectedProject, () => {
  page.value = 1
  fetchData()
})

// 表格列宽调整功能
let resizing = false
let currentTh: HTMLElement | null = null
let startX = 0
let startWidth = 0

const handleResizeStart = (e: MouseEvent, field: string) => {
  resizing = true
  // 根据 field 名称找到对应的 th 元素
  const th = document.querySelector(`th[style*="min-width"]`)  // 这只是示例，实际需要更精确的匹配
  if (e.target instanceof HTMLElement) {
    currentTh = e.target.parentElement as HTMLElement
  }
  startX = e.clientX
  startWidth = currentTh?.offsetWidth || 100
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

const handleMouseMove = (e: MouseEvent) => {
  if (!resizing || !currentTh) return
  
  const width = Math.max(60, startWidth + (e.clientX - startX))
  currentTh.style.width = `${width}px`
  
  // 同步更新同列的所有单元格
  const columnIndex = Array.from(currentTh.parentNode?.children || []).indexOf(currentTh)
  const rows = document.querySelectorAll('.webshell-table-row')
  
  rows.forEach((row) => {
    const cell = row.children[columnIndex] as HTMLElement
    if (cell) {
      cell.style.width = `${width}px`
    }
  })
}

const handleMouseUp = () => {
  resizing = false
  currentTh = null
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}

onMounted(async () => {
  // 添加全局鼠标事件监听用于列宽拖动
  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
  
  // 添加点击空白处关闭菜单的事件
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('contextmenu', handleContextMenuOutside)
  
  // 初始加载数据
  await fetchProjects()
  await fetchDeletedProjects()
  // 等待项目列表加载完成后再获取数据
  if (selectedProject.value) {
    await fetchData()
  }
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
  padding: 24px 24px 20px 24px;
  margin-bottom: 0;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  box-sizing: border-box;
}

.dark .content-header {
  border-bottom-color: var(--border-strong);
}

.content-section h1 {
  margin: 0;
  font-size: 0;
  line-height: 1;
  color: var(--text-color);
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
}

.content-section h1 .title {
  font-size: 24px;
  font-weight: 600;
  color: var(--active-color);
  letter-spacing: 0;
}

.content-section h1 .separator {
  color: var(--text-tertiary);
  font-weight: 300;
  font-size: 20px;
}

.content-section h1 .subtitle {
  font-size: 16px;
  font-weight: 400;
  color: var(--text-secondary);
  letter-spacing: 0;
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
  overflow: hidden; /* 添加 overflow 防止溢出 */
}

/* 项目内容样式 - 深色主题风格 */
.projects-content {
  display: flex;
  gap: 0;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  background: var(--border-color);
  overflow: hidden;
}

.projects-sidebar {
  width: 18%;
  min-width: 220px;
  max-width: 320px;
  background: var(--sidebar-bg);
  padding: 20px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  overflow-y: auto;
  overflow-x: hidden;
  height: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

/* 顶部区域 - 新建项目按钮 */
.sidebar-top-section {
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  padding-bottom: 8px;
}

/* 项目列表区域 - 可滚动 */
.directory-tree {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  margin-bottom: 16px;
}

/* 底部区域 - 回收站固定于底部 */
.sidebar-bottom-section {
  flex-shrink: 0;
  margin-top: auto;
  padding-top: 16px;
}

.projects-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
  background: var(--card-bg);
  height: 100%;
}

.webshell-table-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100%;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.webshell-table-container {
  flex: 1;
  overflow: auto;
  padding: 16px;
  border-radius: 6px;
}

/* 表格样式优化 */
.webshell-table-container :deep(.n-data-table) {
  border-radius: 6px;
  overflow: hidden;
  font-size: 13px;
  table-layout: auto; /* 列宽自适应 */
}

/* 表头样式 */
.webshell-table-container :deep(.n-data-table .n-data-table-th) {
  background: var(--table-header-bg) !important;
  font-weight: 600;
  padding: 10px 12px !important;
  color: var(--text-color);
  border-bottom: 2px solid var(--border-color);
  white-space: nowrap !important;
  max-width: 200px;
}

/* 单元格样式 - 关键修复 */
.webshell-table-container :deep(.n-data-table .n-data-table-td) {
  padding: 8px 12px !important;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.2s ease;
  height: 40px !important;
  max-height: 40px !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  white-space: nowrap !important;
  line-height: 24px !important; /* 使用固定行高，不是 40px */
  vertical-align: middle !important;
}

/* 行样式 */
.webshell-table-container :deep(.n-data-table .n-data-table-tr) {
  height: 40px !important;
}

.webshell-table-container :deep(.n-data-table .n-data-table-tr:hover) {
  background: var(--hover-color) !important;
}

/* 单元格内部容器 - 关键修复 */
.webshell-table-container :deep(.n-data-table .n-data-table-td .n-data-table-cell) {
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  white-space: nowrap !important;
  max-width: 100%;
  display: block;
}

/* 单元格内的任何元素 */
.webshell-table-container :deep(.n-data-table .n-data-table-td *) {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}

.webshell-table-container :deep(.n-data-table .n-data-table-wrapper) {
  border-radius: 6px;
}

/* 特定列的宽度优化 */
.webshell-table-container :deep(.n-data-table .n-data-table-col-id) {
  min-width: 180px;
  max-width: 200px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-url) {
  min-width: 120px;
  max-width: 180px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-payload) {
  min-width: 60px;
  max-width: 80px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-encryption) {
  min-width: 80px;
  max-width: 100px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-encoding) {
  min-width: 60px;
  max-width: 80px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-proxyType) {
  min-width: 80px;
  max-width: 100px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-remark) {
  min-width: 80px;
  max-width: 120px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-createTime) {
  min-width: 140px;
  max-width: 160px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-updateTime) {
  min-width: 140px;
  max-width: 160px;
}

.webshell-table-container :deep(.n-data-table .n-data-table-col-status) {
  min-width: 60px;
  max-width: 80px;
}

/* 滚动条美化 */
.webshell-table-container::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.webshell-table-container::-webkit-scrollbar-track {
  background: var(--scrollbar-track);
  border-radius: 4px;
}

.webshell-table-container::-webkit-scrollbar-thumb {
  background: var(--scrollbar-thumb);
  border-radius: 4px;
}

.webshell-table-container::-webkit-scrollbar-thumb:hover {
  background: var(--scrollbar-thumb-hover);
}

.directory-tree {
  margin-top: 16px;
}

.tree-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  margin-bottom: 4px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  border-radius: 4px;
  position: relative;
}

.tree-item-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.tree-item-actions.show-actions {
  opacity: 1;
}

.tree-item-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 120px;
  display: block;
}

.tree-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
}

.tree-item.active {
  background: var(--active-color);
  color: white;
  box-shadow: var(--shadow-sm);
  font-weight: 500;
}

.tree-item-icon {
  font-size: 16px;
  transition: all 0.3s ease;
}

.tree-item:hover .tree-item-icon {
  transform: scale(1.1);
}

.tree-item.active .tree-item-icon {
  filter: brightness(0) invert(1);
}

.tree-item-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

/* 删除项目按钮样式 */
.delete-project-btn {
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid transparent;
  cursor: pointer;
  font-size: 14px;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  opacity: 0;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color);
}

.tree-item:hover .delete-project-btn,
.tree-item.active .delete-project-btn {
  opacity: 1;
}

.delete-project-btn:hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.3);
  transform: scale(1.08);
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.2);
}

.dark .delete-project-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
}

/* 回收站区域样式优化 */
.recycle-bin-section {
  margin-top: 20px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(0, 0, 0, 0.05) 0%, rgba(0, 0, 0, 0.02) 100%);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 8px;
  backdrop-filter: blur(8px);
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
}

.dark .recycle-bin-section {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.03) 0%, rgba(255, 255, 255, 0.01) 100%);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.recycle-bin-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 12px;
  position: relative;
  width: 100%;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(to right, transparent, var(--border-color), transparent);
  min-width: 30px;
}

.divider-text {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 600;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 6px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  text-align: center;
}

.divider-text::before {
  content: '🗑️';
  font-size: 14px;
  filter: grayscale(0.2);
}

/* 统一的恢复按钮样式 - 始终显示 */
.recover-all-btn {
  width: calc(100% - 32px);  /* 减去左右各 16px 的 padding，确保与容器边缘对齐 */
  margin: 0 16px;             /* 确保与容器边缘有适当间距 */
  padding: 12px 16px;
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.08) 0%, rgba(76, 175, 80, 0.05) 100%);
  border: 1px solid rgba(76, 175, 80, 0.25);
  border-radius: 10px;
  color: var(--text-color);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  text-align: center;
  line-height: 1;
  letter-spacing: 0.3px;
}

.recover-all-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transition: left 0.6s ease;
}

.recover-all-btn:hover::before {
  left: 100%;
}

.recover-all-btn:hover {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.18) 0%, rgba(76, 175, 80, 0.12) 100%);
  border-color: rgba(76, 175, 80, 0.45);
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(76, 175, 80, 0.3);
}

.dark .recover-all-btn {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.12) 0%, rgba(76, 175, 80, 0.08) 100%);
  border-color: rgba(76, 175, 80, 0.3);
}

.dark .recover-all-btn:hover {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.2) 0%, rgba(76, 175, 80, 0.15) 100%);
  box-shadow: 0 6px 16px rgba(76, 175, 80, 0.3);
}

.btn-icon {
  font-size: 18px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  flex-shrink: 0;
  width: 20px;
  height: 20px;
}

.btn-text {
  font-weight: 600;
  letter-spacing: 0.5px;
  text-align: center;
  display: inline-block;
  white-space: nowrap;
  line-height: 1;
}

.btn-count {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.25) 0%, rgba(76, 175, 80, 0.35) 100%);
  color: #4CAF50;
  font-size: 12px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
  min-width: 20px;
  height: 20px;
  text-align: center;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(76, 175, 80, 0.3), inset 0 1px 0 rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
  border: 1px solid rgba(76, 175, 80, 0.4);
  line-height: 1;
}

.dark .btn-count {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.3) 0%, rgba(76, 175, 80, 0.4) 100%);
  box-shadow: 0 2px 6px rgba(76, 175, 80, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.15);
  border-color: rgba(76, 175, 80, 0.5);
}

.recover-all-btn:hover .btn-count {
  transform: scale(1.1);
  box-shadow: 0 3px 10px rgba(76, 175, 80, 0.4);
}

/* 恢复项目弹窗样式 - 基于主界面居中 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(8px);
  animation: fadeIn 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideIn {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-content {
  background: var(--content-bg);
  border-radius: 8px;
  width: 90%;
  max-width: 550px;
  max-height: 75vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3), 0 0 0 1px rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  animation: slideIn 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.modal-header {
  padding: 24px 24px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.05) 0%, transparent 100%);
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--text-color);
  letter-spacing: -0.3px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.modal-header h3::before {
  content: '↩️';
  font-size: 20px;
}

.modal-close {
  background: rgba(0, 0, 0, 0.05);
  border: none;
  font-size: 22px;
  cursor: pointer;
  color: var(--text-color);
  opacity: 0.6;
  transition: all 0.2s ease;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 300;
}

.modal-close:hover {
  opacity: 1;
  background: rgba(0, 0, 0, 0.1);
  transform: rotate(90deg);
}

.dark .modal-close {
  background: rgba(255, 255, 255, 0.05);
}

.dark .modal-close:hover {
  background: rgba(255, 255, 255, 0.1);
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: var(--content-bg);
}

/* 自定义滚动条 */
.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-track {
  background: transparent;
}

.modal-body::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-color);
  opacity: 0.6;
}

.empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 16px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.recover-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recover-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: linear-gradient(135deg, var(--bg-secondary) 0%, rgba(0, 0, 0, 0.02) 100%);
  border-radius: 10px;
  border: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.recover-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(180deg, rgba(239, 68, 68, 0.5) 0%, rgba(239, 68, 68, 0.2) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.recover-item:hover::before {
  opacity: 1;
}

.recover-item:hover {
  transform: translateX(6px);
  border-color: rgba(76, 175, 80, 0.4);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(76, 175, 80, 0.1);
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.05) 0%, var(--bg-secondary) 100%);
}

.dark .recover-item:hover {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(76, 175, 80, 0.2);
}

.recover-item-icon {
  font-size: 22px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.recover-item-info {
  flex: 1;
  min-width: 0;
}

.recover-item-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 6px;
  letter-spacing: -0.2px;
}

.recover-item-time {
  font-size: 12px;
  color: var(--text-secondary);
  opacity: 0.75;
  font-weight: 500;
}

.recover-item-btn {
  flex-shrink: 0;
  padding: 10px 18px;
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.1) 0%, rgba(76, 175, 80, 0.08) 100%);
  border: 1px solid rgba(76, 175, 80, 0.3);
  border-radius: 8px;
  color: var(--text-color);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.recover-item-btn span:first-child {
  font-size: 16px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.recover-item-btn:hover {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.2) 0%, rgba(76, 175, 80, 0.15) 100%);
  border-color: rgba(76, 175, 80, 0.5);
  transform: scale(1.06) translateY(-1px);
  box-shadow: 0 6px 16px rgba(76, 175, 80, 0.25), 0 0 0 1px rgba(76, 175, 80, 0.2);
}

.dark .recover-item-btn {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.12) 100%);
}

.dark .recover-item-btn:hover {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.25) 0%, rgba(76, 175, 80, 0.2) 100%);
  box-shadow: 0 6px 16px rgba(76, 175, 80, 0.3), 0 0 0 1px rgba(76, 175, 80, 0.25);
}

.modal-footer {
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  background: linear-gradient(0deg, rgba(0, 0, 0, 0.02) 0%, transparent 100%);
}

.dark .modal-footer {
  background: linear-gradient(0deg, rgba(255, 255, 255, 0.02) 0%, transparent 100%);
}

.modal-btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  letter-spacing: 0.3px;
}

.modal-btn-secondary {
  background: var(--bg-secondary);
  color: var(--text-color);
  border: 1px solid var(--border-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.modal-btn-secondary:hover {
  background: var(--hover-color);
  border-color: var(--text-secondary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.deleted-projects-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.tree-item.deleted {
  opacity: 0.5;
  background: rgba(0, 0, 0, 0.1);
  border: 1px dashed var(--border-color);
  border-radius: 4px;
  position: relative;
  overflow: hidden;
}

.tree-item.deleted::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--border-color);
}

.tree-item.deleted:hover {
  opacity: 0.8;
  transform: translateX(2px);
  background: rgba(0, 0, 0, 0.15);
  border-color: var(--active-color);
}

.dark .tree-item.deleted {
  background: rgba(255, 255, 255, 0.05);
}

.dark .tree-item.deleted:hover {
  background: rgba(255, 255, 255, 0.1);
}

/* 恢复项目按钮样式 */
.recover-project-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 4px 6px;
  border-radius: 4px;
  transition: all 0.2s ease;
  opacity: 0.6;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.tree-item.deleted:hover .recover-project-btn {
  opacity: 1;
}

.recover-project-btn:hover {
  background: rgba(76, 175, 80, 0.15);
  transform: scale(1.15);
  box-shadow: 0 2px 4px rgba(76, 175, 80, 0.2);
}

.dark .recover-project-btn:hover {
  background: rgba(76, 175, 80, 0.25);
  box-shadow: 0 2px 4px rgba(76, 175, 80, 0.3);
}

/* 空回收站提示 */
.recycle-bin-section:empty::after {
  content: '回收站为空';
  display: block;
  text-align: center;
  padding: 20px;
  color: var(--text-color);
  opacity: 0.4;
  font-size: 13px;
}

/* 新建 WebShell 按钮样式 - 与主题和其他元素保持一致 */
.new-webshell-btn {
  width: 32px !important;
  height: 32px !important;
  min-width: 32px !important;
  padding: 0 !important;
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  border-radius: 8px !important;
  font-size: 20px !important;
  font-weight: 600 !important;
  line-height: 1 !important;
  background: var(--active-color) !important;
  color: #ffffff !important;
  border: none !important;
  box-shadow: var(--shadow-sm) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  cursor: pointer !important;
  vertical-align: middle !important;
  margin: 0 !important;
}

.new-webshell-btn .btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  margin: 0;
  padding: 0;
}

.new-webshell-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: var(--shadow-md) !important;
  opacity: 0.95 !important;
}

.new-webshell-btn:active {
  transform: translateY(0) !important;
  box-shadow: var(--shadow-sm) !important;
}

/* 回收站按钮样式 - 与新建按钮保持统一 */
.recycle-bin-btn {
  height: 32px !important;
  min-width: 100px !important;
  padding: 0 16px !important;
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  border-radius: 8px !important;
  font-size: 14px !important;
  font-weight: 500 !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  cursor: pointer !important;
  border: 1px solid transparent !important;
  outline: none !important;
  vertical-align: middle !important;
  gap: 6px !important;
  margin: 0 !important;
}

.recycle-bin-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: var(--shadow-md) !important;
  opacity: 0.9 !important;
}

.recycle-bin-btn:active {
  transform: translateY(0) !important;
  box-shadow: var(--shadow-sm) !important;
}

/* 回收站按钮在警告状态下的样式 */
.recycle-bin-btn[type="warning"] {
  background-color: var(--warning-color) !important;
  color: white !important;
  border-color: var(--warning-color) !important;
}

.recycle-bin-btn[type="warning"]:hover {
  background-color: var(--warning-color-hover) !important;
  border-color: var(--warning-color-hover) !important;
}

/* 回收站按钮在默认状态下的样式 - 使用强调色 */
.recycle-bin-btn[type="default"] {
  background-color: var(--active-color) !important;
  color: white !important;
  border-color: var(--active-color) !important;
}

.recycle-bin-btn[type="default"]:hover {
  background-color: var(--active-color) !important;
  opacity: 0.9 !important;
  border-color: var(--active-color) !important;
  color: white !important;
}

.recycle-icon {
  font-size: 16px !important;
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  line-height: 1 !important;
  margin-right: 4px !important;
  transition: all 0.3s ease !important;
}

.recycle-text {
  font-size: 14px !important;
  font-weight: 500 !important;
  transition: all 0.3s ease !important;
}

.recycle-bin-btn:hover .recycle-icon {
  transform: scale(1.1) rotate(5deg) !important;
}

.recycle-bin-btn:hover .recycle-text {
  letter-spacing: 0.5px !important;
}

/* 工具栏容器 - 使用 flexbox 确保所有元素完美对齐 */
.toolbar-container {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 32px;
}

/* 分页条数选择器容器 */
.page-size-container {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  margin-left: 4px;
  height: 32px;
  vertical-align: middle;
}

/* 分隔线样式 - 统一主题表现 */
.page-size-divider {
  width: 1px;
  height: 20px;
  background-color: rgba(128, 128, 128, 0.3);
  border-radius: 2px;
  flex-shrink: 0;
  /* 确保在深色和浅色模式下都有一致的视觉效果 */
  opacity: 0.5;
  transition: all 0.3s ease;
}

/* 分页条数选择框样式 */
.page-size-select {
  min-width: 95px !important;
  width: auto !important;
  height: 32px !important;
  vertical-align: middle !important;
  font-size: 13px !important;
  /* 确保与按钮高度一致 */
  line-height: 32px !important;
  margin: 0 !important;
}

.page-size-select :deep(.n-base-selection) {
  background-color: var(--card-bg) !important;
  border: 1px solid var(--border-color) !important;
  border-radius: 6px !important;
  color: var(--text-color) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  height: 32px !important;
  font-size: 13px !important;
}

.page-size-select :deep(.n-base-selection:hover) {
  background-color: var(--hover-color) !important;
  border-color: var(--active-color) !important;
}

.page-size-select :deep(.n-base-selection--active) {
  border-color: var(--active-color) !important;
  box-shadow: 0 0 0 2px rgba(var(--active-color-rgb), 0.1) !important;
}

.page-size-select :deep(.n-base-selection-label) {
  color: var(--text-color) !important;
  font-size: 13px !important;
  padding: 0 8px !important;
  display: flex;
  align-items: center;
}

.page-size-select :deep(.n-base-selection-arrow) {
  color: var(--text-color-3) !important;
  font-size: 14px !important;
  transition: all 0.3s ease !important;
  margin-right: 8px !important;
}

.page-size-select :deep(.n-base-selection:hover .n-base-selection-arrow) {
  color: var(--active-color) !important;
}

.page-size-select :deep(.n-base-selection) {
  min-width: 120px !important; /* 确保选择器本身也足够宽 */
}

.page-size-select :deep(.n-base-option) {
  background-color: var(--card-bg) !important;
  color: var(--text-color) !important;
  font-size: 13px !important;
  padding: 10px 12px !important;
  margin: 2px 0 !important;
  border-radius: 6px !important;
  min-width: 120px !important; /* 确保足够宽度显示完整文字 */
  white-space: nowrap !important; /* 防止文字换行 */
  cursor: pointer !important;
  transition: all 0.2s ease !important;
}

.page-size-select :deep(.n-base-option:hover) {
  background-color: var(--hover-color) !important;
  transform: translateX(2px) !important;
}

.page-size-select :deep(.n-base-option.n-base-option--selected) {
  background-color: var(--active-color-suppl) !important;
  color: var(--active-color) !important;
  font-weight: 600 !important;
  padding: 10px 12px !important;
  margin: 2px 0 !important;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1) !important;
}

/* 强制覆盖 Naive UI 默认样式 */
.page-size-select :deep(.n-base-option--selected .n-base-option__checked) {
  color: var(--active-color) !important;
}

.page-size-select :deep(.n-base-option--selected:hover) {
  background-color: var(--active-color-suppl) !important;
  opacity: 0.95 !important;
  transform: translateX(2px) !important;
}

/* ========== 深色模式完整样式覆盖 ========== */
/* 使用 CSS 变量而非硬编码颜色，确保主题一致性 */
.dark .page-size-select :deep(.n-base-selection) {
  background-color: var(--bg-tertiary) !important;
  border: 1px solid var(--border-strong) !important;
  border-radius: 6px !important;
}

.dark .page-size-select :deep(.n-base-selection-label) {
  color: var(--text-primary) !important;
}

.dark .page-size-select :deep(.n-base-selection:hover) {
  border-color: var(--border-strong) !important;
}

/* 深色模式下拉菜单弹出层 */
.dark .page-size-select :deep(.n-base-selection__menu) {
  background-color: var(--bg-tertiary) !important;
  border: 1px solid var(--border-subtle) !important;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5) !important;
  border-radius: 8px !important;
  padding: 6px !important;
}

/* 深色模式下的选项 */
.dark .page-size-select :deep(.n-base-option) {
  background-color: var(--bg-secondary) !important;
  color: var(--text-secondary) !important;
  border-radius: 6px !important;
  padding: 10px 12px !important;
  margin: 2px 0 !important;
  transition: all 0.2s ease !important;
  cursor: pointer !important;
}

.dark .page-size-select :deep(.n-base-option:hover) {
  background-color: var(--bg-hover) !important;
  color: var(--text-primary) !important;
  transform: translateX(2px) !important;
}

/* 深色模式下的选中项 */
.dark .page-size-select :deep(.n-base-option.n-base-option--selected) {
  background-color: var(--active-color) !important;
  color: #ffffff !important;
  font-weight: 600 !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3) !important;
  padding: 10px 12px !important;
  margin: 2px 0 !important;
}

.dark .page-size-select :deep(.n-base-option--selected .n-base-option__checked) {
  color: #ffffff !important;
}

.dark .page-size-select :deep(.n-base-option--selected:hover) {
  opacity: 0.95 !important;
  transform: translateX(2px) !important;
}

/* 分页组件样式 - 应用强调色 */
.pagination-container :deep(.n-pagination) {
  --n-item-text-color: var(--text-color) !important;
  --n-item-text-color-hover: var(--active-color) !important;
  --n-item-text-color-active: var(--active-color) !important;
  --n-item-border-color: var(--border-color) !important;
  --n-item-border-color-hover: var(--active-color) !important;
  --n-item-border-color-active: var(--active-color) !important;
  --n-item-bg-color-hover: rgba(var(--active-color-rgb), 0.1) !important;
  --n-item-bg-color-active: rgba(var(--active-color-rgb), 0.15) !important;
}

.pagination-container :deep(.n-pagination-item) {
  border-radius: 6px !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.pagination-container :deep(.n-pagination-item:hover) {
  transform: translateY(-2px) !important;
  box-shadow: var(--shadow-sm) !important;
}

.pagination-container :deep(.n-pagination-item.n-pagination-item--active) {
  background-color: var(--active-color) !important;
  color: white !important;
  border-color: var(--active-color) !important;
  font-weight: 600 !important;
}

.pagination-container :deep(.n-pagination-item.n-pagination-item--active:hover) {
  opacity: 0.9 !important;
  transform: translateY(-2px) !important;
  box-shadow: var(--shadow-md) !important;
}

.pagination-container :deep(.n-pagination-item--disabled) {
  opacity: 0.5 !important;
  cursor: not-allowed !important;
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

/* 深色主题下增强表格卡片边框 */
.dark .webshell-table-card {
  border: 1px solid var(--border-strong);
}

/* 覆盖 NCard 的默认内容区域样式 */
.webshell-table-card :deep(.n-card__content) {
  flex: 1;
  padding: 0 !important;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

/* 覆盖 NCard 的默认页脚区域样式 */
.webshell-table-card :deep(.n-card__footer) {
  padding: 16px !important;
  border-top: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 深色模式下的卡片页脚样式 */
.dark .webshell-table-card :deep(.n-card__footer) {
  border-top: 1px solid var(--border-strong) !important;
}

/* WebShell 表格样式 */
.webshell-table-container {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.webshell-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed; /* 固定表格布局，列宽自适应 */
}

.webshell-table-header-row {
  background: var(--hover-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 确保所有表头单元格的基线对齐 */
.webshell-table-header-row th {
  box-sizing: border-box !important;
}

/* 确保所有行单元格的基线对齐 */
.webshell-table-row td {
  box-sizing: border-box !important;
}

.webshell-table-header {
  padding: 0 8px !important; /* 紧凑布局：上下 0，左右 8px */
  border: 1px solid var(--border-color);
  text-align: left;
  cursor: pointer;
  user-select: none;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 53px !important; /* 固定表头高度 53px */
  line-height: 53px !important; /* 垂直居中 */
  overflow: hidden !important; /* 隐藏溢出 */
  text-overflow: ellipsis !important; /* 显示省略号 */
  white-space: nowrap !important; /* 不换行 */
  display: table-cell !important; /* 确保作为表格单元格显示 */
  vertical-align: middle !important; /* 垂直居中 */
  color: var(--text-color);
}

.webshell-table-header span {
  display: inline-block !important;
  vertical-align: middle !important;
  line-height: normal !important; /* 使用默认行高，确保文字清晰 */
  transition: color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.webshell-table-header:hover {
  background: var(--bg-hover);
}

.webshell-table-row {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 53px !important; /* 固定行高 53px（比 40px 增加约三分之一） */
}

.webshell-table-row:hover {
  background: var(--bg-hover);
}

.webshell-table-cell {
  padding: 0 8px !important; /* 紧凑布局：上下 0，左右 8px */
  border: 1px solid var(--border-color) !important;
  text-align: left;
  background-color: transparent !important;
  color: var(--text-color) !important;
  height: 53px !important; /* 固定单元格高度 53px */
  max-height: 53px !important; /* 最大高度 53px */
  vertical-align: middle !important; /* 垂直居中 */
  overflow: hidden !important; /* 隐藏溢出 */
  text-overflow: ellipsis !important; /* 显示省略号 */
  white-space: nowrap !important; /* 不换行 */
  display: table-cell !important; /* 确保作为表格单元格显示 */
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 所有单元格内的内容都使用 flex 布局确保垂直居中 */
.webshell-table-cell > * {
  display: flex !important;
  align-items: center !important;
  height: 100% !important;
  width: 100% !important;
}

/* 单元格内的 Tooltip wrapper 需要特殊处理 */
.webshell-table-cell .tooltip-wrapper {
  display: block !important;
  width: 100% !important;
  overflow: hidden !important;
}

.webshell-table-cell .tooltip-content {
  display: block !important;
  width: 100% !important;
  overflow: hidden !important;
}

.webshell-table-cell .tooltip-content span {
  display: block !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  white-space: nowrap !important;
  width: 100% !important;
}

.webshell-table-cell-truncate {
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  white-space: nowrap !important;
  display: flex !important;
  align-items: center !important;
  width: 100% !important;
  height: 53px !important; /* 与单元格高度一致 */
}

.webshell-table-cell-truncate span {
  display: inline-block !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  white-space: nowrap !important;
  max-width: 100% !important;
  vertical-align: middle !important;
  line-height: normal !important;
}

/* Tooltip 在表格中的样式 - URL 和备注列使用 flex 布局 */
.webshell-table-cell .tooltip-wrapper {
  display: flex !important;
  align-items: center !important;
  width: 100% !important;
  height: 100% !important;
  overflow: hidden !important;
}

.webshell-table-cell .tooltip-content {
  display: flex !important;
  align-items: center !important;
  width: 100% !important;
  height: 100% !important;
  overflow: hidden !important;
}

.webshell-table-cell .tooltip-content span {
  display: inline-block !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  white-space: nowrap !important;
  width: 100% !important;
  vertical-align: middle !important;
  line-height: normal !important;
}

/* 所有单元格内的文字都应用省略号和垂直居中 */
.webshell-table-cell span {
  display: inline-block !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  white-space: nowrap !important;
  width: 100% !important;
  vertical-align: middle !important;
  line-height: normal !important;
}

/* 列宽拖动手柄样式优化 */
.resize-handle {
  opacity: 0;
  transition: opacity 0.2s ease;
  background-color: var(--active-color);
}

.resize-handle:hover,
.resize-handle:active {
  opacity: 0.5;
}

.webshell-table-header:hover .resize-handle {
  opacity: 0.3;
}

.webshell-table-header:hover .resize-handle:hover {
  opacity: 0.5;
}

/* 深色主题下优化表格样式，确保边框清晰可见 */
.dark .webshell-table-container {
  border: 1px solid var(--border-strong) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .webshell-table-header {
  border: 1px solid var(--border-strong) !important;
  background-color: var(--bg-tertiary) !important; /* 表头背景色 */
  color: var(--text-primary) !important; /* 表头文字颜色 */
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .webshell-table-cell {
  border: 1px solid var(--border-strong) !important;
  background-color: transparent !important; /* 单元格背景透明 */
  color: var(--text-primary) !important; /* 单元格文字颜色 */
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 深色模式下增强表头与单元格的视觉区分 */
.dark .webshell-table-header .webshell-table-cell {
  background-color: var(--bg-tertiary) !important; /* 表头单元格稍深的背景 */
  font-weight: var(--font-semibold) !important; /* 表头文字加粗 */
  border-bottom: 2px solid var(--border-strong) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 深色模式下悬停效果优化 */
.dark .webshell-table-row:hover {
  background-color: var(--bg-hover) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 搜索输入框样式 */
.search-input {
  border-radius: 8px !important;
  border: 1px solid #d1d5db !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.search-input:hover {
  border-color: #9ca3af !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1) !important;
}

.search-input:focus,
.search-input:focus-within {
  border-color: var(--active-color) !important;
  box-shadow: 0 0 0 3px var(--active-color-suppl) !important;
}

/* 深色模式下的搜索框样式 */
.dark .search-input {
  border-color: var(--border-subtle) !important;
}

.dark .search-input:hover {
  border-color: var(--border-strong) !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3) !important;
}

.dark .search-input:focus,
.dark .search-input:focus-within {
  border-color: var(--active-color) !important;
  box-shadow: 0 0 0 3px var(--active-color-suppl) !important;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding-top: 16px;
}

/* 表格行选中状态 - 浅色模式 */
.table-row-selected {
  background-color: rgba(var(--active-color-rgb), 0.1) !important; /* 使用系统强调色的淡色背景 */
  border-left: 4px solid var(--active-color) !important;
  border-top: 2px solid var(--active-color) !important;
  border-bottom: 2px solid var(--active-color) !important;
  box-shadow: 0 0 0 1px var(--active-color) !important; /* 外发光效果 */
}

/* 深色主题下选中行样式 - 与浅色模式保持一致的样式参数 */
.dark .table-row-selected {
  background-color: rgba(var(--active-color-rgb), 0.15) !important; /* 深色模式下适度增加背景透明度 */
  border-left: 4px solid var(--active-color) !important; /* 使用系统强调色 */
  border-top: 2px solid var(--active-color) !important; /* 与浅色模式一致的边框粗细 */
  border-bottom: 2px solid var(--active-color) !important; /* 与浅色模式一致的边框粗细 */
  box-shadow: 0 0 0 1px var(--active-color) !important; /* 与浅色模式一致的外发光效果 */
}

/* 选中行单元格边框优化 - 移除单元格边框，让行的边框显示 */
.table-row-selected .webshell-table-cell {
  border: none !important; /* 移除所有单元格边框 */
}

/* 选中行第一个单元格的左边框优化 */
.dark .table-row-selected .webshell-table-cell:first-child,
.table-row-selected .webshell-table-cell:first-child {
  border-left: none !important; /* 移除单元格左边框，让行的左边框显示 */
}

/* 响应式设计 - 大屏幕优化 */
@media (min-width: 1920px) {
  .projects-sidebar {
    width: 15%;
    min-width: 280px;
    max-width: 380px;
  }
}

/* 中等屏幕优化 */
@media (min-width: 1440px) and (max-width: 1919px) {
  .projects-sidebar {
    width: 18%;
    min-width: 240px;
    max-width: 340px;
  }
}

/* 小屏幕优化 */
@media (min-width: 1024px) and (max-width: 1439px) {
  .projects-sidebar {
    width: 22%;
    min-width: 220px;
    max-width: 300px;
  }
}

/* 平板 landscape */
@media (min-width: 768px) and (max-width: 1023px) {
  .projects-sidebar {
    width: 28%;
    min-width: 200px;
    max-width: 280px;
  }
}

/* 平板 portrait 和移动端 */
@media (max-width: 767px) {
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
    max-height: 200px;
    min-width: 100%;
    max-width: 100%;
  }
  
  /* 移动端按钮优化 */
  .recycle-bin-btn {
    min-width: 80px !important;
    padding: 0 12px !important;
    font-size: 13px !important;
  }
  
  .new-webshell-btn {
    width: 32px !important;
    height: 32px !important;
    min-width: 32px !important;
    font-size: 18px !important;
  }
  
  /* 移动端分页选择框优化 */
  .page-size-container {
    gap: 10px !important;
    margin-left: 2px !important;
  }
  
  .page-size-divider {
    height: 18px !important;
  }
  
  .page-size-select {
    min-width: 85px !important;
    height: 30px !important;
    font-size: 12px !important;
    line-height: 30px !important;
  }
  
  .page-size-select :deep(.n-base-selection) {
    height: 30px !important;
    font-size: 12px !important;
  }
  
  /* 移动端表格行高优化 */
  .webshell-table-row {
    height: 44px !important;
  }
  
  .webshell-table-cell {
    padding: 6px 10px !important;
    font-size: 12px !important;
  }
}

@media (max-width: 480px) {
  /* 小屏幕进一步优化 */
  .recycle-bin-btn {
    min-width: 70px !important;
    padding: 0 10px !important;
    font-size: 12px !important;
  }
  
  .recycle-icon {
    font-size: 14px !important;
    margin-right: 3px !important;
  }
  
  .recycle-text {
    font-size: 12px !important;
  }
  
  /* 小屏幕分页选择框优化 */
  .page-size-container {
    gap: 8px !important;
    margin-left: 0 !important;
  }
  
  .page-size-divider {
    height: 16px !important;
  }
  
  .page-size-select {
    min-width: 75px !important;
    height: 28px !important;
    font-size: 11px !important;
    line-height: 28px !important;
  }
  
  .page-size-select :deep(.n-base-selection) {
    height: 28px !important;
    font-size: 11px !important;
  }
  
  /* 小屏幕表格行高进一步优化 */
  .webshell-table-row {
    height: 40px !important;
  }
  
  .webshell-table-cell {
    padding: 5px 8px !important;
    font-size: 11px !important;
  }
}

/* 右键菜单样式 - 现代化主题风格 */

/* 菜单图标样式 */
.menu-icon-enter {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.menu-icon-cache {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.menu-icon-edit {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.menu-icon-delete {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.menu-item-enter {
  color: var(--active-color);
}

.menu-item-cache {
  color: var(--text-secondary);
}

.menu-item-edit {
  color: var(--primary-500);
}

.menu-item-delete {
  color: var(--error-500);
}

/* 菜单容器样式 - 紧凑布局 */
.n-dropdown {
  animation: menuFadeIn 0.2s ease-out;
  background-color: transparent !important; /* 容器本身透明 */
  border: none !important;
}

/* 菜单内部容器 - 实际显示背景的区域 */
.n-dropdown-menu {
  background-color: var(--card-bg) !important;
  border: 1px solid var(--border-strong) !important;
  border-radius: var(--radius-lg) !important;
  box-shadow: var(--shadow-lg) !important;
  backdrop-filter: blur(8px);
  padding: 4px !important;
  min-width: fit-content !important;
}

/* 深色主题下优化菜单背景色，使用精确的深蓝色调确保与整体主题协调 */
.dark .n-dropdown-menu {
  border: 1px solid var(--border-strong) !important;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5) !important;
  background-color: var(--bg-secondary) !important;
  backdrop-filter: blur(16px); /* 增强毛玻璃效果 */
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

/* 菜单项样式 - 紧凑布局 */
.n-menu-item {
  background-color: transparent !important;
  color: var(--text-primary) !important;
  padding: 8px 12px !important; /* 减小内边距 */
  transition: all var(--duration-fast) var(--ease-standard);
  border-radius: var(--radius-md) !important;
  margin: 2px 4px !important; /* 减小外边距 */
  font-size: var(--text-sm) !important;
  font-weight: var(--font-medium) !important;
  white-space: nowrap !important; /* 防止文字换行 */
}

/* 菜单项悬停效果 */
.n-menu-item:hover {
  background-color: var(--hover-color) !important;
  color: var(--text-primary) !important;
  transform: translateX(2px);
}

/* 菜单项选中效果 */
.n-menu-item.n-menu-item--selected {
  background-color: var(--active-color-suppl) !important;
  color: var(--active-color) !important;
  font-weight: var(--font-semibold) !important;
}

/* 深色主题下优化菜单项对比度和背景色，使用精确颜色值 */
.dark .n-menu-item {
  color: var(--text-primary) !important;
}

.dark .n-menu-item:hover {
  background-color: var(--bg-hover) !important;
  color: var(--text-primary) !important;
}

.dark .n-menu-item.n-menu-item--selected {
  background-color: var(--active-color-suppl) !important;
  color: var(--active-color) !important;
}

/* 菜单项图标样式 */
.n-menu-item-icon {
  margin-right: 10px; /* 减小图标间距 */
  font-size: 16px;
  width: 18px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0; /* 防止图标被压缩 */
}

/* 菜单项文字样式 */
.n-menu-item-content {
  font-size: var(--text-sm) !important;
  font-weight: var(--font-medium) !important;
  letter-spacing: 0.025em;
  flex-shrink: 0; /* 防止文字被压缩 */
}

/* 菜单项内容区域优化 */
.n-menu-item-content__main {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 菜单分隔线样式 */
.n-dropdown-menu-item-divider {
  background-color: var(--border-subtle) !important;
  margin: 6px 8px !important;
}

/* 深色主题下增强分隔线对比度 */
.dark .n-dropdown-menu-item-divider {
  background-color: var(--border-strong) !important;
}
</style>

<!-- 全局样式 - 右键菜单背景色（必须在 scoped 外，因为菜单被 teleport 到 body） -->
<style>
/* 深色主题下强制应用正确的背景色 */
.dark .n-dropdown-menu {
  background-color: var(--bg-secondary) !important;
  border: 1px solid var(--border-strong) !important;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5) !important;
  backdrop-filter: blur(16px) !important;
}

/* 浅色主题下的背景色 */
.n-dropdown-menu {
  background-color: var(--bg-primary) !important;
  border: 1px solid var(--border-subtle) !important;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1) !important;
  backdrop-filter: blur(8px) !important;
}
</style>