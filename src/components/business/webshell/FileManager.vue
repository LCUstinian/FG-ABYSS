<template>
  <div class="file-manager">
    <div class="file-manager-container">
      <!-- 工具栏 -->
      <div class="file-toolbar">
        <div class="toolbar-left">
          <n-space>
            <n-button 
              size="small" 
              quaternary
              @click="handleRefresh"
              :disabled="!connected"
            >
              <template #icon>
                <span>🔄</span>
              </template>
            </n-button>
            <n-button 
              size="small" 
              quaternary
              @click="handleGoUp"
              :disabled="!connected || currentPath === '/'"
            >
              <template #icon>
                <span>⬆️</span>
              </template>
            </n-button>
          </n-space>
        </div>
        <div class="toolbar-center">
          <n-input 
            v-model:value="currentPath" 
            size="small"
            placeholder="当前路径"
            @keydown.enter="handleNavigate"
          >
            <template #prefix>
              <span>📂</span>
            </template>
          </n-input>
        </div>
        <div class="toolbar-right">
          <n-space>
            <n-button 
              size="small" 
              type="primary"
              @click="handleUpload"
              :disabled="!connected"
            >
              <template #icon>
                <span>⬆️</span>
              </template>
              上传
            </n-button>
            <n-button 
              size="small"
              @click="handleNewFolder"
              :disabled="!connected"
            >
              <template #icon>
                <span>📁</span>
              </template>
              新建文件夹
            </n-button>
          </n-space>
        </div>
      </div>

      <!-- 文件列表 -->
      <div class="file-list-container">
        <n-scrollbar ref="scrollbarRef">
          <div v-if="!connected" class="file-manager-placeholder">
            <span class="placeholder-icon">📁</span>
            <h3 class="placeholder-title">文件管理器</h3>
            <p class="placeholder-description">
              连接 WebShell 后可以使用文件管理功能
            </p>
            <p class="placeholder-text">
              WebShell ID: {{ webshellId || '未指定' }}
            </p>
          </div>
          
          <div v-else class="file-list">
            <div class="file-list-header">
              <div class="file-item">
                <div class="file-icon">图标</div>
                <div class="file-name">名称</div>
                <div class="file-size">大小</div>
                <div class="file-modified">修改时间</div>
                <div class="file-actions">操作</div>
              </div>
            </div>
            
            <div class="file-list-body">
              <!-- 父目录 -->
              <div 
                v-if="currentPath !== '/'" 
                class="file-item directory"
                @click="handleGoUp"
              >
                <div class="file-icon">📁</div>
                <div class="file-name">..</div>
                <div class="file-size">-</div>
                <div class="file-modified">-</div>
                <div class="file-actions">-</div>
              </div>
              
              <!-- 文件和文件夹列表 -->
              <div 
                v-for="file in mockFiles" 
                :key="file.name"
                :class="['file-item', file.type]"
                @dblclick="handleFileClick(file)"
                @contextmenu="(e) => handleContextMenu(e, file)"
              >
                <div class="file-icon">{{ getFileIcon(file.type) }}</div>
                <div class="file-name">
                  <span>{{ file.name }}</span>
                </div>
                <div class="file-size">{{ formatFileSize(file.size) }}</div>
                <div class="file-modified">{{ file.modified }}</div>
                <div class="file-actions">
                  <n-space>
                    <n-button 
                      size="tiny" 
                      quaternary
                      @click.stop="handleDownload(file)"
                    >
                      <template #icon>
                        <span>⬇️</span>
                      </template>
                    </n-button>
                    <n-button 
                      size="tiny" 
                      quaternary
                      @click.stop="handleRename(file)"
                    >
                      <template #icon>
                        <span>✏️</span>
                      </template>
                    </n-button>
                    <n-button 
                      size="tiny" 
                      quaternary
                      @click.stop="handleDelete(file)"
                    >
                      <template #icon>
                        <span>🗑️</span>
                      </template>
                    </n-button>
                  </n-space>
                </div>
              </div>
            </div>
          </div>
        </n-scrollbar>
      </div>

      <!-- 状态栏 -->
      <div class="file-statusbar">
        <div class="status-info">
          <span v-if="connected">
            共 {{ mockFiles.length }} 个项目
          </span>
          <span v-else>未连接</span>
        </div>
        <div class="status-path">
          当前路径：{{ currentPath }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { NSpace, NButton, NInput, NScrollbar, useMessage } from 'naive-ui'

interface FileItem {
  name: string
  type: 'file' | 'directory'
  size: number
  modified: string
  path: string
}

defineProps<{
  webshellId?: string
}>()

const message = useMessage()
const connected = ref(false)
const currentPath = ref('/')
const scrollbarRef = ref<any | null>(null)

// Mock 文件数据
const mockFiles = ref<FileItem[]>([
  { name: 'index.php', type: 'file', size: 15360, modified: '2024-01-15 10:30', path: '/index.php' },
  { name: 'config.php', type: 'file', size: 2048, modified: '2024-01-14 15:20', path: '/config.php' },
  { name: 'uploads', type: 'directory', size: 0, modified: '2024-01-13 09:15', path: '/uploads' },
  { name: 'includes', type: 'directory', size: 0, modified: '2024-01-12 14:45', path: '/includes' },
  { name: 'assets', type: 'directory', size: 0, modified: '2024-01-11 11:30', path: '/assets' },
  { name: '.htaccess', type: 'file', size: 512, modified: '2024-01-10 16:00', path: '/.htaccess' },
  { name: 'robots.txt', type: 'file', size: 256, modified: '2024-01-09 08:45', path: '/robots.txt' },
])

const load = async () => {
  // 模拟连接
  connected.value = true
  currentPath.value = '/'
  message.success('文件管理器已加载')
}

const cleanup = () => {
  connected.value = false
  console.log('File manager cleaned up')
}

const handleRefresh = () => {
  message.success('刷新文件列表')
}

const handleGoUp = () => {
  if (currentPath.value !== '/') {
    const parts = currentPath.value.split('/').filter(Boolean)
    parts.pop()
    currentPath.value = '/' + parts.join('/') || '/'
  }
}

const handleNavigate = () => {
  message.info('导航到：' + currentPath.value)
}

const handleUpload = () => {
  message.info('上传文件功能开发中')
}

const handleNewFolder = () => {
  message.info('新建文件夹功能开发中')
}

const handleFileClick = (file: FileItem) => {
  if (file.type === 'directory') {
    currentPath.value = file.path
    message.success('进入目录：' + file.name)
  }
}

const handleContextMenu = (e: MouseEvent, file: FileItem) => {
  e.preventDefault()
  message.info('右键菜单功能开发中')
}

const handleDownload = (file: FileItem) => {
  message.info('下载文件：' + file.name)
}

const handleRename = (file: FileItem) => {
  message.info('重命名文件：' + file.name)
}

const handleDelete = (file: FileItem) => {
  message.warning('删除文件：' + file.name)
}

const getFileIcon = (type: string) => {
  const icons: Record<string, string> = {
    'file': '📄',
    'directory': '📁',
  }
  return icons[type] || '📄'
}

const formatFileSize = (size: number) => {
  if (size === 0) return '-'
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  return `${(size / 1024 / 1024).toFixed(2)} MB`
}

defineExpose({
  load,
  cleanup
})
</script>

<style scoped>
.file-manager {
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: var(--card-bg);
}

.file-manager-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 工具栏 */
.file-toolbar {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  gap: 12px;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
}

.toolbar-center {
  flex: 1;
  max-width: 600px;
}

/* 文件列表 */
.file-list-container {
  flex: 1;
  overflow: hidden;
}

.file-list {
  padding: 0;
}

.file-list-header {
  background: var(--content-bg);
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 1;
}

.file-list-body {
  padding: 8px 0;
}

.file-item {
  display: grid;
  grid-template-columns: 40px 1fr 100px 150px 120px;
  align-items: center;
  padding: 8px 16px;
  gap: 12px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-item:hover {
  background: var(--hover-color);
}

.file-item.directory {
  font-weight: 500;
}

.file-item.directory:hover {
  background: var(--hover-color);
}

.file-icon {
  font-size: 18px;
  text-align: center;
}

.file-name {
  font-size: 13px;
  color: var(--text-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size,
.file-modified {
  font-size: 12px;
  color: var(--text-secondary);
}

.file-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
}

/* 占位状态 */
.file-manager-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.placeholder-icon {
  font-size: 64px;
  margin-bottom: 20px;
  opacity: 0.6;
}

.placeholder-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-color);
  margin: 0 0 12px 0;
}

.placeholder-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 16px 0;
}

.placeholder-text {
  font-size: 13px;
  color: var(--text-tertiary);
}

/* 状态栏 */
.file-statusbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 16px;
  background: var(--card-bg);
  border-top: 1px solid var(--border-color);
  font-size: 12px;
  color: var(--text-secondary);
}

.status-info,
.status-path {
  flex: 1;
}

.status-path {
  text-align: right;
}
</style>
