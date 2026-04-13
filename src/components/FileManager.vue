<template>
  <div class="file-manager">
    <!-- 工具栏 -->
    <div class="file-toolbar">
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
        <n-button size="small" @click="handleGoUp" :disabled="currentPath === '/'">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="19" x2="12" y2="5"/>
              <polyline points="5 12 12 5 19 12"/>
            </svg>
          </template>
          上级
        </n-button>
        <n-button size="small" @click="handleUpload">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="17 8 12 3 7 8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
          </template>
          上传
        </n-button>
        <n-button size="small" @click="handleNewFolder">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
              <line x1="12" y1="11" x2="12" y2="17"/>
              <line x1="9" y1="14" x2="15" y2="14"/>
            </svg>
          </template>
          新建文件夹
        </n-button>
      </n-space>

      <n-space style="margin-left: auto">
        <n-input
          v-model:value="searchQuery"
          size="small"
          placeholder="搜索文件..."
          style="width: 200px"
          clearable
        />
      </n-space>
    </div>

    <!-- 路径导航 -->
    <div class="path-bar">
      <n-breadcrumb>
        <n-breadcrumb-item
          v-for="(segment, index) in pathSegments"
          :key="index"
          @click="handleNavigateTo(index)"
          style="cursor: pointer"
        >
          {{ segment || '根目录' }}
        </n-breadcrumb-item>
      </n-breadcrumb>
    </div>

    <!-- 文件列表 -->
    <div class="file-list">
      <n-data-table
        :columns="columns"
        :data="filteredFiles"
        :row-key="rowKey"
        :single-line="false"
        @update:checked-row-keys="handleCheck"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { NButton, NBreadcrumb, NBreadcrumbItem, NDataTable, NInput, NSpace, NTag, useMessage, useDialog } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { invoke } from '@tauri-apps/api/core'

interface FileItem {
  key: string
  name: string
  type: 'file' | 'folder'
  size: number
  modified: string
  permissions: string
}

const message = useMessage()
const dialog = useDialog()

// 状态
const currentPath = ref('/')
const searchQuery = ref('')
const fileList = ref<FileItem[]>([])
const checkedFiles = ref<string[]>([])

// 路径分段
const pathSegments = computed(() => {
  return currentPath.value.split('/')
})

// 过滤后的文件列表
const filteredFiles = computed(() => {
  if (!searchQuery.value) return fileList.value
  return fileList.value.filter(file =>
    file.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 表格列定义
const columns = computed<DataTableColumns<FileItem>>(() => [
  {
    type: 'selection',
  },
  {
    title: '名称',
    key: 'name',
    render: (row) => {
      return h('div', {
        style: {
          display: 'flex',
          alignItems: 'center',
          gap: '8px',
          cursor: row.type === 'folder' ? 'pointer' : 'default',
        },
        onClick: () => row.type === 'folder' && handleOpenFolder(row.name),
      }, [
        h('span', {
          style: { fontSize: '16px' },
        }, row.type === 'folder' ? '📁' : '📄'),
        h('span', {}, row.name),
      ])
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 100,
    render: (row) => {
      if (row.type === 'folder') return '-'
      return formatSize(row.size)
    },
  },
  {
    title: '修改时间',
    key: 'modified',
    width: 180,
  },
  {
    title: '权限',
    key: 'permissions',
    width: 100,
    render: (row) => h(NTag, {
      type: row.permissions.includes('w') ? 'warning' : 'info',
      size: 'small',
    }, { default: () => row.permissions }),
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row) => {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, {
            size: 'tiny',
            onClick: () => handleDownload(row),
          }, { default: () => '下载' }),
          h(NButton, {
            size: 'tiny',
            type: 'warning',
            onClick: () => handleRename(row),
          }, { default: () => '重命名' }),
          h(NButton, {
            size: 'tiny',
            type: 'error',
            onClick: () => handleDelete(row),
          }, { default: () => '删除' }),
        ],
      })
    },
  },
])

// 行键
const rowKey = (row: FileItem) => row.key

// 格式化文件大小
const formatSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`
}

// 加载文件列表
const loadFileList = async (path: string = currentPath.value) => {
  try {
    const response = await invoke('list_directory', {
      webshellId: 'placeholder', // TODO: 使用真实的 webshell ID
      path: path,
    })
    
    // TODO: 解析真实的响应数据
    fileList.value = [
      { key: '1', name: 'Documents', type: 'folder', size: 0, modified: '2024-01-01 12:00', permissions: 'rwx' },
      { key: '2', name: 'Downloads', type: 'folder', size: 0, modified: '2024-01-02 14:30', permissions: 'rwx' },
      { key: '3', name: 'test.php', type: 'file', size: 2048, modified: '2024-01-03 09:15', permissions: 'rw-' },
    ]
    
    currentPath.value = path
    message.success('文件列表加载成功')
  } catch (error) {
    message.error(`加载文件列表失败：${error}`)
  }
}

// 刷新
const handleRefresh = () => {
  loadFileList()
}

// 上级目录
const handleGoUp = () => {
  const segments = currentPath.value.split('/').filter(Boolean)
  segments.pop()
  const newPath = '/' + segments.join('/')
  loadFileList(newPath || '/')
}

// 导航到指定层级
const handleNavigateTo = (index: number) => {
  const segments = pathSegments.value.slice(0, index + 1)
  const newPath = segments.join('/')
  loadFileList(newPath || '/')
}

// 打开文件夹
const handleOpenFolder = (folderName: string) => {
  const newPath = currentPath.value === '/' 
    ? `/${folderName}`
    : `${currentPath.value}/${folderName}`
  loadFileList(newPath)
}

// 选中文件
const handleCheck = (keys: string[]) => {
  checkedFiles.value = keys
}

// 上传文件
const handleUpload = async () => {
  // TODO: 实现文件上传对话框
  message.info('上传功能开发中...')
}

// 新建文件夹
const handleNewFolder = () => {
  dialog.input({
    title: '新建文件夹',
    placeholder: '文件夹名称',
    onConfirm: async (name) => {
      if (!name) {
        message.warning('请输入文件夹名称')
        return
      }
      // TODO: 调用后端 API 创建文件夹
      message.success(`文件夹 "${name}" 创建成功`)
      handleRefresh()
    },
  })
}

// 下载文件
const handleDownload = (file: FileItem) => {
  if (file.type === 'folder') {
    message.warning('不能下载文件夹')
    return
  }
  // TODO: 实现文件下载
  message.info(`下载文件：${file.name}`)
}

// 重命名文件
const handleRename = (file: FileItem) => {
  dialog.input({
    title: '重命名',
    placeholder: '新名称',
    defaultValue: file.name,
    onConfirm: async (newName) => {
      if (!newName) {
        message.warning('请输入新名称')
        return
      }
      // TODO: 调用后端 API 重命名
      message.success(`文件已重命名为 "${newName}"`)
      handleRefresh()
    },
  })
}

// 删除文件
const handleDelete = (file: FileItem) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除 ${file.type === 'folder' ? '文件夹' : '文件'} "${file.name}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      // TODO: 调用后端 API 删除
      message.success(`"${file.name}" 已删除`)
      handleRefresh()
    },
  })
}

// 初始化
loadFileList()
</script>

<style scoped>
.file-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
}

.file-toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.path-bar {
  margin-bottom: 12px;
  padding: 8px 12px;
  background-color: var(--n-color);
  border-radius: 4px;
}

.file-list {
  flex: 1;
  overflow: auto;
}
</style>
