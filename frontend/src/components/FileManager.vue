<template>
  <div class="file-manager-container">
    <!-- 工具栏 -->
    <div class="file-toolbar">
      <div class="toolbar-left">
        <n-button size="small" @click="handleRefresh">
          <template #icon>
            <n-icon><Refresh /></n-icon>
          </template>
          刷新
        </n-button>
        
        <n-button size="small" @click="handleGoUp" :disabled="!canGoUp">
          <template #icon>
            <n-icon><ArrowUp /></n-icon>
          </template>
          上级
        </n-button>
        
        <n-button size="small" @click="showCreateFolder = true">
          <template #icon>
            <n-icon><FolderAdd /></n-icon>
          </template>
          新建文件夹
        </n-button>
        
        <n-button size="small" @click="handleUpload">
          <template #icon>
            <n-icon><Upload /></n-icon>
          </template>
          上传
        </n-button>
        
        <n-button size="small" @click="handleDownload" :disabled="!selectedFile">
          <template #icon>
            <n-icon><Download /></n-icon>
          </template>
          下载
        </n-button>
        
        <n-button size="small" @click="handleDelete" :disabled="!selectedFile" type="error">
          <template #icon>
            <n-icon><Delete /></n-icon>
          </template>
          删除
        </n-button>
      </div>
      
      <div class="toolbar-right">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索文件..."
          size="small"
          style="width: 200px"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <n-icon><Search /></n-icon>
          </template>
        </n-input>
      </div>
    </div>
    
    <!-- 路径导航 -->
    <div class="breadcrumb-bar">
      <n-breadcrumb>
        <n-breadcrumb-item
          v-for="(item, index) in breadcrumbPath"
          :key="index"
          @click="handleNavigateTo(index)"
        >
          {{ item }}
        </n-breadcrumb-item>
      </n-breadcrumb>
    </div>
    
    <!-- 文件列表 -->
    <div class="file-list">
      <n-data-table
        :columns="columns"
        :data="fileList"
        :row-key="rowKey"
        :loading="loading"
        :pagination="false"
        :bordered="false"
        @update:checked-row-keys="handleCheck"
      />
    </div>
    
    <!-- 上传文件弹窗 -->
    <n-modal
      v-model:show="showUploadModal"
      preset="dialog"
      title="上传文件"
      style="width: 500px"
    >
      <n-upload
        :action="uploadUrl"
        :headers="uploadHeaders"
        :on-finish="handleUploadFinish"
        :on-error="handleUploadError"
        multiple
      >
        <n-upload-dragger>
          <div style="margin-bottom: 12px">
            <n-icon size="48" :depth="3">
              <Archive />
            </n-icon>
          </div>
          <n-text style="font-size: 16px">
            点击或者拖动文件到此处上传
          </n-text>
        </n-upload-dragger>
      </n-upload>
    </n-modal>
    
    <!-- 新建文件夹弹窗 -->
    <n-modal
      v-model:show="showCreateFolder"
      preset="dialog"
      title="新建文件夹"
      style="width: 400px"
    >
      <n-input
        v-model:value="newFolderName"
        placeholder="文件夹名称"
        @keyup.enter="handleCreateFolder"
      />
      <template #action>
        <n-button @click="showCreateFolder = false">取消</n-button>
        <n-button type="primary" @click="handleCreateFolder">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import {
  Refresh,
  ArrowUp,
  FolderAdd,
  Upload,
  Download,
  Delete,
  Search,
  Archive
} from '@vicons/carbon'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'

const message = useMessage()
const dialog = useDialog()

// 状态
const loading = ref(false)
const currentPath = ref('/')
const searchKeyword = ref('')
const selectedFile = ref<string | null>(null)
const showUploadModal = ref(false)
const showCreateFolder = ref(false)
const newFolderName = ref('')

// 模拟数据
const fileList = ref([
  {
    key: '1',
    name: '..',
    size: 0,
    isDir: true,
    modTime: '',
    permissions: 'drwxr-xr-x'
  },
  {
    key: '2',
    name: 'www',
    size: 0,
    isDir: true,
    modTime: '2024-01-15 10:30:00',
    permissions: 'drwxr-xr-x'
  },
  {
    key: '3',
    name: 'index.php',
    size: 15360,
    isDir: false,
    modTime: '2024-01-15 11:20:00',
    permissions: '-rw-r--r--'
  },
  {
    key: '4',
    name: 'config.php',
    size: 2048,
    isDir: false,
    modTime: '2024-01-14 09:15:00',
    permissions: '-rw-r--r--'
  },
  {
    key: '5',
    name: 'uploads',
    size: 0,
    isDir: true,
    modTime: '2024-01-13 16:45:00',
    permissions: 'drwxr-xr-x'
  }
])

// 面包屑路径
const breadcrumbPath = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  return ['/', ...parts]
})

const canGoUp = computed(() => {
  return currentPath.value !== '/'
})

// 表格列配置
const columns: DataTableColumns = [
  {
    type: 'selection'
  },
  {
    title: '名称',
    key: 'name',
    width: 300,
    render: (row) => {
      const icon = row.isDir ? '📁' : '📄'
      return h('span', {
        style: { cursor: 'pointer', color: row.isDir ? '#1890ff' : 'inherit' },
        onClick: () => handleFileClick(row)
      }, `${icon} ${row.name}`)
    }
  },
  {
    title: '大小',
    key: 'size',
    width: 100,
    render: (row) => {
      if (row.isDir) return '-'
      return formatFileSize(row.size)
    }
  },
  {
    title: '修改时间',
    key: 'modTime',
    width: 180
  },
  {
    title: '权限',
    key: 'permissions',
    width: 120
  }
]

// 上传相关
const uploadUrl = ref('')
const uploadHeaders = ref({})

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// 行键
const rowKey = (row: any) => row.key

// 加载文件列表
const loadFileList = async () => {
  loading.value = true
  try {
    // TODO: 调用后端 API
    // const result = await FileHandler.ListFiles(webshellId, currentPath.value)
    // fileList.value = result.files
    
    // 模拟加载
    await new Promise(resolve => setTimeout(resolve, 500))
  } catch (error: any) {
    message.error('加载文件列表失败：' + error.message)
  } finally {
    loading.value = false
  }
}

// 处理文件点击
const handleFileClick = (row: any) => {
  if (row.name === '..') {
    handleGoUp()
    return
  }
  
  if (row.isDir) {
    currentPath.value = row.path || currentPath.value + '/' + row.name
    loadFileList()
  } else {
    selectedFile.value = row.name
  }
}

// 刷新
const handleRefresh = () => {
  loadFileList()
}

// 上级目录
const handleGoUp = () => {
  const parts = currentPath.value.split('/').filter(Boolean)
  if (parts.length > 0) {
    parts.pop()
    currentPath.value = '/' + parts.join('/')
    loadFileList()
  }
}

// 导航到指定层级
const handleNavigateTo = (index: number) => {
  if (index === 0) {
    currentPath.value = '/'
  } else {
    currentPath.value = '/' + breadcrumbPath.value.slice(1, index + 1).join('/')
  }
  loadFileList()
}

// 搜索
const handleSearch = () => {
  if (!searchKeyword.value) {
    loadFileList()
    return
  }
  
  // TODO: 调用搜索 API
  message.info('搜索：' + searchKeyword.value)
}

// 上传
const handleUpload = () => {
  showUploadModal.value = true
}

const handleUploadFinish = () => {
  message.success('上传成功')
  showUploadModal.value = false
  loadFileList()
}

const handleUploadError = (error: any) => {
  message.error('上传失败：' + error.message)
}

// 下载
const handleDownload = async () => {
  if (!selectedFile.value) {
    message.warning('请先选择文件')
    return
  }
  
  try {
    // TODO: 调用下载 API
    // const content = await FileHandler.DownloadFile(webshellId, currentPath.value + '/' + selectedFile.value)
    
    message.success('下载成功')
  } catch (error: any) {
    message.error('下载失败：' + error.message)
  }
}

// 删除
const handleDelete = () => {
  if (!selectedFile.value) {
    message.warning('请先选择文件')
    return
  }
  
  dialog.warning({
    title: '确认删除',
    content: '确定要删除 "' + selectedFile.value + '" 吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        // TODO: 调用删除 API
        // await FileHandler.DeleteFile(webshellId, currentPath.value + '/' + selectedFile.value)
        
        message.success('删除成功')
        loadFileList()
      } catch (error: any) {
        message.error('删除失败：' + error.message)
      }
    }
  })
}

// 新建文件夹
const handleCreateFolder = async () => {
  if (!newFolderName.value) {
    message.warning('请输入文件夹名称')
    return
  }
  
  try {
    // TODO: 调用创建目录 API
    // await FileHandler.CreateDirectory(webshellId, currentPath.value + '/' + newFolderName.value)
    
    message.success('创建成功')
    showCreateFolder.value = false
    newFolderName.value = ''
    loadFileList()
  } catch (error: any) {
    message.error('创建失败：' + error.message)
  }
}

// 选择文件
const handleCheck = (keys: DataTableRowKey[]) => {
  selectedFile.value = keys.length > 0 ? String(keys[0]) : null
}

// 生命周期
onMounted(() => {
  loadFileList()
})
</script>

<style scoped>
.file-manager-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #ffffff;
}

.dark .file-manager-container {
  background: #1e1e1e;
}

.file-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f5f5f5;
  border-bottom: 1px solid #e8e8e8;
}

.dark .file-toolbar {
  background: #2d2d2d;
  border-bottom-color: #3d3d3d;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.breadcrumb-bar {
  padding: 8px 12px;
  background: #fafafa;
  border-bottom: 1px solid #e8e8e8;
}

.dark .breadcrumb-bar {
  background: #252525;
  border-bottom-color: #3d3d3d;
}

.file-list {
  flex: 1;
  overflow: auto;
  padding: 12px;
}

:deep(.n-data-table) {
  font-size: 13px;
}

:deep(.n-data-table-th) {
  background: #fafafa;
}

.dark :deep(.n-data-table-th) {
  background: #2d2d2d;
}
</style>
