<template>
  <div class="payload-workspace">
    <n-tabs v-model:value="activeTab" type="line" animated>
      <!-- Payload 生成标签页 -->
      <n-tab-pane name="generator" tab="Payload 生成" display="flex">
        <PayloadGenerator />
      </n-tab-pane>

      <!-- Payload 列表标签页 -->
      <n-tab-pane name="list" tab="Payload 列表" display="flex">
        <n-card :bordered="false">
          <template #header>
            <n-space justify="space-between">
              <n-text>已生成的 Payload</n-text>
              <n-space>
                <n-input
                  v-model:value="searchQuery"
                  placeholder="搜索 Payload..."
                  clearable
                  style="width: 240px"
                >
                  <template #prefix>
                    <n-icon :component="SearchOutline" />
                  </template>
                </n-input>
                <n-button type="primary" @click="handleRefresh">
                  <template #icon>
                    <n-icon :component="RefreshOutline" />
                  </template>
                  刷新
                </n-button>
              </n-space>
            </n-space>
          </template>

          <n-data-table
            :columns="payloadColumns"
            :data="filteredPayloads"
            :row-key="row => row.id"
            :pagination="pagination"
            :scroll-x="1200"
          />
        </n-card>
      </n-tab-pane>

      <!-- 模板管理标签页 -->
      <n-tab-pane name="templates" tab="模板管理" display="flex">
        <n-card :bordered="false">
          <template #header>
            <n-space justify="space-between">
              <n-text>自定义模板</n-text>
              <n-button type="primary" @click="showCreateTemplateModal = true">
                <template #icon>
                  <n-icon :component="AddOutline" />
                </template>
                创建模板
              </n-button>
            </n-space>
          </template>

          <n-data-table
            :columns="templateColumns"
            :data="customTemplates"
            :row-key="row => row.name"
            :pagination="false"
          />
        </n-card>
      </n-tab-pane>
    </n-tabs>

    <!-- 创建模板弹窗 -->
    <n-modal
      v-model:show="showCreateTemplateModal"
      preset="dialog"
      title="创建自定义模板"
      style="width: 800px"
    >
      <n-form
        ref="templateFormRef"
        :model="templateForm"
        :rules="templateFormRules"
        label-placement="top"
      >
        <n-form-item label="模板名称" path="name">
          <n-input v-model:value="templateForm.name" placeholder="输入模板名称" />
        </n-form-item>

        <n-form-item label="Payload 类型" path="type">
          <n-select
            v-model:value="templateForm.type"
            :options="typeOptions"
            placeholder="选择类型"
          />
        </n-form-item>

        <n-form-item label="功能模式" path="function">
          <n-select
            v-model:value="templateForm.function"
            :options="functionOptions"
            placeholder="选择功能"
          />
        </n-form-item>

        <n-form-item label="模板内容" path="content">
          <n-input
            v-model:value="templateForm.content"
            type="textarea"
            :rows="15"
            placeholder="输入模板代码，使用 {{.password}} 作为密码占位符"
            style="font-family: monospace"
          />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space justify="end">
          <n-button @click="showCreateTemplateModal = false">取消</n-button>
          <n-button type="primary" @click="handleCreateTemplate">创建</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 查看 Payload 详情弹窗 -->
    <n-modal
      v-model:show="showViewModal"
      preset="dialog"
      title="Payload 详情"
      style="width: 800px"
      positive-text="复制"
      negative-text="关闭"
      @positive-click="handleCopyPayload"
    >
      <n-scrollbar style="max-height: 400px">
        <pre class="payload-content"><code>{{ selectedPayload?.content }}</code></pre>
      </n-scrollbar>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import {
  SearchOutline,
  RefreshOutline,
  AddOutline,
  EyeOutline,
  DownloadOutline,
  TrashOutline,
  CopyOutline,
} from '@vicons/ionicons5'
import type { FormRules, FormInst, DataTableColumns } from 'naive-ui'
import PayloadGenerator from './PayloadGenerator.vue'
import { GetTemplates, AddTemplate, DeleteTemplate } from '../../bindings/fg-abyss/internal/app/handlers/payloadhandler.js'

interface PayloadItem {
  id: number
  name: string
  type: string
  function: string
  createdAt: string
  size: number
  content: string
}

interface TemplateItem {
  name: string
  type: string
  function: string
  content: string
  isCustom: boolean
}

const message = useMessage()
const dialog = useDialog()
const activeTab = ref('generator')
const searchQuery = ref('')
const showCreateTemplateModal = ref(false)
const showViewModal = ref(false)
const selectedPayload = ref<PayloadItem | null>(null)

const payloads = ref<PayloadItem[]>([])
const customTemplates = ref<TemplateItem[]>([])

const templateFormRef = ref<FormInst | null>(null)
const templateForm = reactive({
  name: '',
  type: 'php',
  function: 'basic',
  content: '',
})

const templateFormRules: FormRules = {
  name: {
    required: true,
    message: '请输入模板名称',
    trigger: 'blur',
  },
  type: {
    required: true,
    message: '请选择类型',
    trigger: 'change',
  },
  function: {
    required: true,
    message: '请选择功能',
    trigger: 'change',
  },
  content: {
    required: true,
    message: '请输入模板内容',
    trigger: 'blur',
  },
}

const typeOptions = [
  { label: 'PHP', value: 'php' },
  { label: 'ASP', value: 'asp' },
  { label: 'ASPX', value: 'aspx' },
  { label: 'JSP', value: 'jsp' },
]

const functionOptions = [
  { label: '基础命令 (Basic)', value: 'basic' },
  { label: '完整功能 (Full)', value: 'full' },
]

const pagination = reactive({
  pageSize: 10,
  pageSizes: [10, 20, 50, 100],
  showSizePicker: true,
})

const filteredPayloads = computed(() => {
  if (!searchQuery.value) return payloads.value
  return payloads.value.filter((p) =>
    p.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const payloadColumns: DataTableColumns = [
  {
    title: '名称',
    key: 'name',
    width: 200,
    ellipsis: { tooltip: true },
  },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row: PayloadItem) => {
      const typeMap: Record<string, string> = {
        php: 'success',
        asp: 'warning',
        aspx: 'info',
        jsp: 'error',
      }
      return h('n-tag', { type: typeMap[row.type] || 'default' }, () => row.type.toUpperCase())
    },
  },
  {
    title: '功能',
    key: 'function',
    width: 100,
    render: (row: PayloadItem) => {
      const funcMap: Record<string, string> = {
        basic: '基础',
        full: '完整',
      }
      return funcMap[row.function] || row.function
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 100,
    render: (row: PayloadItem) => formatFileSize(row.size),
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 180,
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    fixed: 'right',
    render: (row: PayloadItem) => {
      return h('n-space', {}, [
        h(
          'n-button',
          {
            size: 'small',
            quaternary: true,
            onClick: () => handleViewPayload(row),
          },
          {
            icon: () => h('n-icon', { component: EyeOutline }),
            default: () => '查看',
          }
        ),
        h(
          'n-button',
          {
            size: 'small',
            quaternary: true,
            onClick: () => handleDownloadPayload(row),
          },
          {
            icon: () => h('n-icon', { component: DownloadOutline }),
            default: () => '下载',
          }
        ),
        h(
          'n-popconfirm',
          {
            onPositiveClick: () => handleDeletePayload(row),
          },
          {
            trigger: () =>
              h(
                'n-button',
                {
                  size: 'small',
                  quaternary: true,
                  type: 'error',
                },
                {
                  icon: () => h('n-icon', { component: TrashOutline }),
                  default: () => '删除',
                }
              ),
            default: () => '确定要删除这个 Payload 吗？',
          }
        ),
      ])
    },
  },
]

const templateColumns: DataTableColumns = [
  {
    title: '名称',
    key: 'name',
    width: 150,
  },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row: TemplateItem) => {
      const typeMap: Record<string, string> = {
        php: 'success',
        asp: 'warning',
        aspx: 'info',
        jsp: 'error',
      }
      return h('n-tag', { type: typeMap[row.type] || 'default' }, () => row.type.toUpperCase())
    },
  },
  {
    title: '功能',
    key: 'function',
    width: 100,
    render: (row: TemplateItem) => {
      const funcMap: Record<string, string> = {
        basic: '基础',
        full: '完整',
      }
      return funcMap[row.function] || row.function
    },
  },
  {
    title: '内容预览',
    key: 'content',
    ellipsis: { tooltip: true },
    render: (row: TemplateItem) => {
      const preview = row.content.substring(0, 100)
      return h('n-text', { depth: 3 }, () => (preview + '...').replace(/\n/g, ' '))
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row: TemplateItem) => {
      if (!row.isCustom) return null
      return h(
        'n-popconfirm',
        {
          onPositiveClick: () => handleDeleteTemplate(row),
        },
        {
          trigger: () =>
            h(
              'n-button',
              {
                size: 'small',
                quaternary: true,
                type: 'error',
              },
              {
                icon: () => h('n-icon', { component: TrashOutline }),
                default: () => '删除',
              }
            ),
          default: () => '确定要删除这个模板吗？',
        }
      )
    },
  },
]

const handleRefresh = () => {
  loadPayloads()
  message.success('刷新成功')
}

const handleViewPayload = (row: PayloadItem) => {
  selectedPayload.value = row
  showViewModal.value = true
}

const handleCopyPayload = () => {
  if (!selectedPayload.value) return
  navigator.clipboard.writeText(selectedPayload.value.content)
  message.success('代码已复制')
}

const handleDownloadPayload = (row: PayloadItem) => {
  const blob = new Blob([row.content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = row.name
  a.click()
  URL.revokeObjectURL(url)
  message.success('下载已开始')
}

const handleDeletePayload = (row: PayloadItem) => {
  const index = payloads.value.findIndex((p) => p.id === row.id)
  if (index !== -1) {
    payloads.value.splice(index, 1)
    message.success('删除成功')
  }
}

const handleCreateTemplate = async () => {
  try {
    await templateFormRef.value?.validate()

    await AddTemplate({
      name: templateForm.name,
      type: templateForm.type,
      function: templateForm.function,
      content: templateForm.content,
    })

    message.success('模板创建成功')
    showCreateTemplateModal.value = false
    loadTemplates()

    // 重置表单
    templateForm.name = ''
    templateForm.type = 'php'
    templateForm.function = 'basic'
    templateForm.content = ''
  } catch (error: any) {
    if (error.errors) return // 表单验证失败
    message.error('创建模板失败：' + (error.message || '未知错误'))
  }
}

const handleDeleteTemplate = async (row: TemplateItem) => {
  try {
    await DeleteTemplate(row.name)
    message.success('模板删除成功')
    loadTemplates()
  } catch (error: any) {
    message.error('删除模板失败：' + (error.message || '未知错误'))
  }
}

const loadPayloads = () => {
  // 模拟加载已生成的 Payload
  payloads.value = [
    {
      id: 1,
      name: 'shell.php',
      type: 'php',
      function: 'basic',
      createdAt: '2024-01-15 10:30:00',
      size: 512,
      content: '<?php ... ?>',
    },
  ]
}

const loadTemplates = async () => {
  try {
    const response = await GetTemplates()
    customTemplates.value = response.templates.map((t) => ({
      ...t,
      isCustom: true,
    }))
  } catch (error) {
    console.error('加载模板失败:', error)
  }
}

const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

onMounted(() => {
  loadPayloads()
  loadTemplates()
})
</script>

<style scoped lang="scss">
.payload-workspace {
  padding: 20px;
  height: calc(100vh - 120px);

  :deep(.n-tabs) {
    height: 100%;
  }

  :deep(.n-tab-pane) {
    padding: 0;
    height: 100%;
  }

  .payload-content {
    background-color: var(--code-background);
    border-radius: 8px;
    padding: 16px;
    font-family: 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.5;
    color: var(--code-text);
    white-space: pre-wrap;
    word-wrap: break-word;
    margin: 0;
  }
}
</style>
