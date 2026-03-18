<template>
  <div class="project-list-view">
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
                    :placeholder="t('common.searchPlaceholder')"
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
import Tooltip from '@/components/shared/Tooltip.vue'
import CreateProjectModal from '@/components/business/project/CreateProjectModal.vue'
import CreateWebShellModal from './CreateWebShellModal.vue'
import RecoverProjectModal from '@/components/business/project/RecoverProjectModal.vue'
import { useSmartPagination } from '@/composables/useSmartPagination'
import { useProject } from '@/composables'

// Tauri Mock 适配层 - 替换 Wails API
import { emitEvent, invoke } from '@/utils/tauri-mock-adapter'
import type { Project as ProjectType } from '@/types'

// 导入时间格式化工具
import { formatTime } from '@/utils/formatTime'

// 导入工具函数
import { componentLogger } from '@/utils/logger'

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

const { t } = useI18n()

// 使用 Naive UI 的对话框和消息组件
const dialog = useDialog()
const message = useMessage()

// 使用项目 composable
const { 
  projects: projectList,
  fetchProjects: fetchProjectList 
} = useProject()

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

// 项目列表 - 使用 composable
const projects = projectList
const deletedProjects = ref<ProjectType[]>([])  // 前端临时存储已删除的项目

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
  await fetchProjectList()
  newProjectDialogVisible.value = false
}

// 获取项目列表 - 使用 composable
const fetchProjects = async () => {
  await fetchProjectList()
  
  // 如果没有选中项目，且项目列表不为空，选择第一个项目
  if (!selectedProject.value && projects.value.length > 0) {
    selectedProject.value = projects.value[0].id
    console.log('自动选择第一个项目:', selectedProject.value)
  }
}

// 获取已删除项目列表
const fetchDeletedProjects = async () => {
  try {
    const deletedList = await invoke<ProjectType[]>('get_deleted_projects')
    
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
        // 调用 Tauri Mock 删除接口
        await invoke('delete_project', { projectId: project.id })
        
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
      console.log(t('project.cancelDelete'))
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
    // 调用 Tauri Mock 恢复接口
    await invoke('recover_project', { projectId: project.id })
    
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

// 处理关闭弹窗
const handleRecoverDialogClose = () => {
  console.log(t('project.closeRecoverDialog'))
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
    // Tauri Mock 事件 - 替换 Wails EventsEmit
    await emitEvent('open-webshell-window', {
      id: webshell.id,
      name: webshell.remark || webshell.url,
      url: webshell.url,
    })
    
    message.success(`已打开 WebShell 控制窗口：${webshell.remark || webshell.url}`)
  } catch (error: any) {
    componentLogger.error('发送打开 WebShell 窗口事件失败:', error)
    message.error(t('project.openWebShellFailed') + (error.message || t('common.unknownError')))
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
      message.info(t('project.cacheNotImplemented'))
      break
    case 'edit':
      console.log('Edit webshell:', selectedRow.value)
      // TODO: 打开编辑窗口
      message.info(t('project.editNotImplemented'))
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
          const loading = message.loading(t('projects.deleting'), {
            duration: 0
          })
          
          try {
            // 调用 Tauri Mock 删除方法
            await invoke('delete_webshell', { webshellId: selectedRow.value.id })
            
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
          console.log(t('project.cancelDelete'))
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
          const loading = message.loading(t('projects.recovering'), {
            duration: 0
          })
          
          try {
            // 调用 Tauri Mock 恢复方法
            await invoke('recover_webshell', { webshellId: selectedRow.value.id })
            
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
      // 获取已删除的 WebShell（回收站）- 使用 Tauri Mock API
      data = await invoke<WebShell[]>('get_deleted_webshells', { projectId: selectedProject.value })
      count = data.length
      console.log('获取回收站数据成功，数量:', count)
    } else {
      // 获取正常的 WebShell - 使用 Tauri Mock API
      data = await invoke<WebShell[]>('get_webshells', { 
        projectId: selectedProject.value,
        page: page.value,
        pageSize: pageSize.value,
        searchQuery: searchQuery.value,
        sortField: sortField.value,
        sortDirection: sortDirection.value
      })
      count = data.length
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
  currentTh = e.currentTarget as HTMLElement
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
</script>

<style scoped>
.project-list-view {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
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
  overflow: hidden;
}

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

.sidebar-top-section {
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  padding-bottom: 8px;
}

.directory-tree {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  margin-bottom: 16px;
}

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

.pagination-container {
  display: flex;
  justify-content: center;
  padding: 12px;
  border-top: 1px solid var(--border-color);
}

.toolbar-container {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.page-size-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-size-select {
  width: 120px;
}

.search-input {
  background: var(--input-bg);
  border-color: var(--border-color);
}

.new-project-button {
  width: 100%;
  padding: 12px;
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: var(--border-radius-md);
  font-size: 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-normal);
  display: flex;
  align-items: center;
  justify-content: center;
}

.new-project-button:hover {
  background: var(--active-color-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.tree-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  margin-bottom: 4px;
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: all var(--transition-normal);
  gap: 8px;
}

.tree-item:hover {
  background: var(--hover-color);
}

.tree-item.active {
  background: var(--active-color);
  color: white;
}

.tree-item-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.tree-item-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tree-item-actions {
  opacity: 0;
  transition: opacity var(--transition-normal);
  display: flex;
  gap: 4px;
}

.tree-item-actions.show-actions {
  opacity: 1;
}

.delete-project-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 4px;
  border-radius: 4px;
  transition: all var(--transition-fast);
}

.delete-project-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  transform: scale(1.1);
}

.recycle-bin-section {
  padding: 16px 0;
}

.recycle-bin-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: var(--border-color);
}

.divider-text {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}

.recover-all-btn {
  width: 100%;
  padding: 10px;
  background: var(--warning-color);
  color: white;
  border: none;
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: all var(--transition-normal);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
}

.recover-all-btn:hover {
  background: var(--warning-color-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.btn-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.recycle-bin-btn {
  display: flex;
  align-items: center;
  gap: 6px;
}

.recycle-icon {
  font-size: 14px;
}

.recycle-text {
  font-size: 12px;
}

.new-webshell-btn {
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon {
  font-size: 16px;
  font-weight: 600;
}

.webshell-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.webshell-table-header-row {
  background: var(--table-header-bg);
}

.webshell-table-header {
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  color: var(--text-color);
  border-bottom: 2px solid var(--border-color);
  user-select: none;
  position: relative;
}

.webshell-table-header:hover {
  background: var(--hover-color);
}

.webshell-table-row {
  transition: all var(--transition-fast);
}

.webshell-table-row:hover {
  background: var(--hover-color);
}

.table-row-selected {
  background: var(--active-color) !important;
  color: white;
}

.webshell-table-cell {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.resize-handle {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 5px;
  cursor: col-resize;
  z-index: 10;
  background: transparent;
}

.resize-handle:hover {
  background: var(--active-color);
  opacity: 0.3;
}

:deep(.n-dropdown) {
  z-index: 1000;
}

:deep(.n-menu-item) {
  transition: all var(--transition-fast);
}

:deep(.n-menu-item:hover) {
  background: var(--hover-color);
}
</style>
