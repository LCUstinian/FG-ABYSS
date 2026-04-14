<template>
  <div class="plugin-view">
    <n-grid :cols="100" :x-gap="12" style="height: 100%">
      <!-- 左侧：插件列表 -->
      <n-grid-item :span="30">
        <n-card :bordered="false" style="height: 100%">
          <template #header>
            <n-space justify="space-between">
              <span>插件列表</span>
              <n-button size="small" @click="handleRefresh">
                <template #icon>
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="23 4 23 10 17 10"/>
                    <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
                  </svg>
                </template>
              </n-button>
            </n-space>
          </template>

          <n-input
            v-model:value="searchQuery"
            placeholder="搜索插件..."
            size="small"
            style="margin-bottom: 12px"
            clearable
          />

          <n-tabs v-model:value="filterType" type="line" animated style="margin-bottom: 12px">
            <n-tab-pane name="all" tab="全部" />
            <n-tab-pane name="active" tab="已启用" />
            <n-tab-pane name="inactive" tab="已禁用" />
          </n-tabs>

          <n-list hoverable>
            <n-list-item
              v-for="plugin in filteredPlugins"
              :key="plugin.id"
              @click="handleSelectPlugin(plugin)"
            >
              <template #prefix>
                <n-avatar :style="{ backgroundColor: getPluginColor(plugin.type) }">
                  {{ plugin.name.charAt(0).toUpperCase() }}
                </n-avatar>
              </template>
              <n-thing :title="plugin.name" :description="plugin.description">
                <template #header-extra>
                  <n-tag :type="plugin.status === 'active' ? 'success' : 'default'" size="small">
                    {{ plugin.status === 'active' ? '已启用' : '已禁用' }}
                  </n-tag>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-grid-item>

      <!-- 右侧：插件详情和管理 -->
      <n-grid-item :span="70">
        <n-card :bordered="false" style="height: 100%" v-if="selectedPlugin">
          <template #header>
            <n-space justify="space-between">
              <n-space>
                <n-avatar :size="48" :style="{ backgroundColor: getPluginColor(selectedPlugin.type) }">
                  {{ selectedPlugin.name.charAt(0).toUpperCase() }}
                </n-avatar>
                <div>
                  <h2 style="margin: 0">{{ selectedPlugin.name }}</h2>
                  <span style="color: var(--n-text-color-3); font-size: 12px">
                    v{{ selectedPlugin.version }} - {{ selectedPlugin.author }}
                  </span>
                </div>
              </n-space>
              <n-space>
                <n-button
                  :type="selectedPlugin.status === 'active' ? 'warning' : 'success'"
                  @click="handleTogglePlugin"
                >
                  {{ selectedPlugin.status === 'active' ? '禁用' : '启用' }}
                </n-button>
                <n-button type="error" @click="handleUninstall">
                  <template #icon>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3 6 5 6 21 6"/>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                  </template>
                  卸载
                </n-button>
              </n-space>
            </n-space>
          </template>

          <n-tabs type="line" animated>
            <n-tab-pane name="info" tab="信息">
              <n-descriptions bordered :column="2">
                <n-descriptions-item label="插件 ID">
                  {{ selectedPlugin.id }}
                </n-descriptions-item>
                <n-descriptions-item label="版本">
                  {{ selectedPlugin.version }}
                </n-descriptions-item>
                <n-descriptions-item label="描述">
                  {{ selectedPlugin.description }}
                </n-descriptions-item>
                <n-descriptions-item label="类型">
                  <n-tag size="small">{{ selectedPlugin.type }}</n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="作者">
                  {{ selectedPlugin.author }}
                </n-descriptions-item>
                <n-descriptions-item label="最小版本要求">
                  {{ selectedPlugin.min_version }}
                </n-descriptions-item>
                <n-descriptions-item label="状态">
                  <n-tag :type="selectedPlugin.status === 'active' ? 'success' : 'default'">
                    {{ selectedPlugin.status }}
                  </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="路径">
                  {{ selectedPlugin.path }}
                </n-descriptions-item>
              </n-descriptions>

              <n-divider>依赖</n-divider>
              <n-empty v-if="!selectedPlugin.dependencies || selectedPlugin.dependencies.length === 0" description="无依赖" />
              <n-space v-else>
                <n-tag v-for="dep in selectedPlugin.dependencies" :key="dep">
                  {{ dep }}
                </n-tag>
              </n-space>

              <n-divider>权限</n-divider>
              <n-empty v-if="!selectedPlugin.permissions || selectedPlugin.permissions.length === 0" description="无权限" />
              <n-space v-else>
                <n-tag v-for="perm in selectedPlugin.permissions" :key="perm" type="warning">
                  {{ perm }}
                </n-tag>
              </n-space>
            </n-tab-pane>

            <n-tab-pane name="config" tab="配置">
              <n-form label-placement="left" label-width="100">
                <n-form-item label="自动加载">
                  <n-switch v-model:value="pluginConfig.autoLoad" />
                </n-form-item>
                <n-form-item label="优先级">
                  <n-input-number v-model:value="pluginConfig.priority" :min="0" :max="100" />
                </n-form-item>
                <n-form-item label=" ">
                  <n-button type="primary" @click="handleSaveConfig">
                    <template #icon>
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
                        <polyline points="17 21 17 13 7 13 7 21"/>
                        <polyline points="7 3 7 8 15 8"/>
                      </svg>
                    </template>
                    保存配置
                  </n-button>
                </n-form-item>
              </n-form>
            </n-tab-pane>
          </n-tabs>
        </n-card>

        <n-empty v-else description="请选择一个插件" style="height: 100%" />
      </n-grid-item>
    </n-grid>

    <!-- 安装插件按钮 -->
    <n-float-button
      position="bottom-right"
      :right="40"
      :bottom="40"
      type="primary"
      @click="handleInstallPlugin"
    >
      <template #icon>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="17 8 12 3 7 8"/>
          <line x1="12" y1="3" x2="12" y2="15"/>
        </svg>
      </template>
    </n-float-button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  NAvatar, NButton, NCard, NDescriptions, NDescriptionsItem, NDivider, NEmpty,
  NFloatButton, NForm, NFormItem, NGrid, NGridItem, NInput, NInputNumber, NList,
  NListItem, NSpace, NSwitch, NTabPane, NTabs, NTag, NThing, useMessage, useDialog,
} from 'naive-ui'
import { invoke } from '@tauri-apps/api/core'

interface Plugin {
  id: string
  name: string
  description: string
  version: string
  author: string
  min_version: string
  type: string
  status: 'active' | 'inactive'
  path: string
  dependencies?: string[]
  permissions?: string[]
}

interface PluginConfig {
  enabled: boolean
  autoLoad: boolean
  priority: number
  settings: any
}

const message = useMessage()
const dialog = useDialog()

// 状态
const searchQuery = ref('')
const filterType = ref<'all' | 'active' | 'inactive'>('all')
const selectedPlugin = ref<Plugin | null>(null)
const pluginConfig = ref<PluginConfig>({
  enabled: false,
  autoLoad: false,
  priority: 0,
  settings: null,
})

// 模拟插件数据
const plugins = ref<Plugin[]>([
  {
    id: 'terminal-plugin',
    name: '终端管理插件',
    description: '提供终端管理功能',
    version: '1.0.0',
    author: 'FG-ABYSS',
    min_version: '0.1.0',
    type: 'terminal',
    status: 'active',
    path: '/plugins/terminal-plugin',
    dependencies: [],
    permissions: ['execute_command'],
  },
  {
    id: 'file-manager-plugin',
    name: '文件管理插件',
    description: '提供文件管理功能',
    version: '1.0.0',
    author: 'FG-ABYSS',
    min_version: '0.1.0',
    type: 'file_manager',
    status: 'active',
    path: '/plugins/file-manager-plugin',
    dependencies: [],
    permissions: ['read_file', 'write_file'],
  },
])

// 过滤后的插件列表
const filteredPlugins = computed(() => {
  let result = plugins.value

  // 搜索
  if (searchQuery.value) {
    result = result.filter(p =>
      p.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      p.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  // 类型过滤
  if (filterType.value === 'active') {
    result = result.filter(p => p.status === 'active')
  } else if (filterType.value === 'inactive') {
    result = result.filter(p => p.status === 'inactive')
  }

  return result
})

// 获取插件类型颜色
const getPluginColor = (type: string) => {
  const colors: Record<string, string> = {
    terminal: '#409eff',
    file_manager: '#67c23a',
    database: '#e6a23c',
    network_tool: '#909399',
    crypto_tool: '#f56c6c',
  }
  return colors[type] || '#909399'
}

// 选择插件
const handleSelectPlugin = (plugin: Plugin) => {
  selectedPlugin.value = plugin
  pluginConfig.value = {
    enabled: plugin.status === 'active',
    autoLoad: false,
    priority: 0,
    settings: null,
  }
}

// 刷新列表
const handleRefresh = async () => {
  // TODO: 调用后端 API 刷新插件列表
  message.success('刷新成功')
}

// 切换插件状态
const handleTogglePlugin = async () => {
  if (!selectedPlugin.value) return

  const action = selectedPlugin.value.status === 'active' ? '禁用' : '启用'
  
  try {
    // TODO: 调用后端 API
    message.success(`${action}插件成功`)
    if (selectedPlugin.value) {
      selectedPlugin.value.status = selectedPlugin.value.status === 'active' ? 'inactive' : 'active'
    }
  } catch (error) {
    message.error(`${action}插件失败：${error}`)
  }
}

// 卸载插件
const handleUninstall = () => {
  if (!selectedPlugin.value) return

  dialog.warning({
    title: '确认卸载',
    content: `确定要卸载插件 "${selectedPlugin.value.name}" 吗？此操作不可恢复。`,
    positiveText: '卸载',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        // TODO: 调用后端 API
        message.success('卸载成功')
        const index = plugins.value.findIndex(p => p.id === selectedPlugin.value!.id)
        if (index !== -1) {
          plugins.value.splice(index, 1)
        }
        selectedPlugin.value = null
      } catch (error) {
        message.error(`卸载失败：${error}`)
      }
    },
  })
}

// 保存配置
const handleSaveConfig = async () => {
  if (!selectedPlugin.value) return

  try {
    // TODO: 调用后端 API
    message.success('配置保存成功')
  } catch (error) {
    message.error(`保存配置失败：${error}`)
  }
}

// 安装插件
const handleInstallPlugin = async () => {
  // TODO: 使用 Tauri dialog API 选择 ZIP 文件
  message.info('安装插件功能开发中...')
}
</script>

<style scoped>
.plugin-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
  box-sizing: border-box;
}

.plugin-view :deep(.n-card) {
  height: 100%;
}

.plugin-view :deep(.n-card__content) {
  height: calc(100% - 60px);
  overflow: auto;
}
</style>
