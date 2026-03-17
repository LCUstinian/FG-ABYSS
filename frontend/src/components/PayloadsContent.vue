<template>
  <div class="content-section payloads-section">
    <!-- 标准标题栏 -->
    <div class="content-header">
      <h1>
        <span class="title">载荷</span>
        <span class="separator">|</span>
        <span class="subtitle">WebShell 生成器</span>
      </h1>
    </div>

    <!-- 内容主体 -->
    <div class="content-body">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <!-- Payload 生成标签页 -->
        <n-tab-pane name="generator" tab="Payload 生成">
          <div class="generator-container">
            <n-grid :cols="24" :x-gap="20" :y-gap="20">
              <!-- 配置区域 -->
              <n-grid-item :span="12">
                <n-card title="Payload 配置" :bordered="false" class="config-card">
                  <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="left" :label-width="100">
                    <n-form-item label="脚本类型" path="type">
                      <n-select v-model:value="formData.type" :options="typeOptions" />
                    </n-form-item>

                    <n-form-item label="功能类型" path="function">
                      <n-select v-model:value="formData.function" :options="functionOptions" />
                    </n-form-item>

                    <n-form-item label="连接密码" path="password">
                      <n-input v-model:value="formData.password" placeholder="请输入密码">
                        <template #prefix>
                          <n-icon :component="KeyOutline" />
                        </template>
                      </n-input>
                    </n-form-item>

                    <n-form-item label="编码器" path="encoder">
                      <n-select v-model:value="formData.encoder" :options="encoderOptions" />
                    </n-form-item>

                    <n-form-item label="混淆级别" path="obfuscationLevel">
                      <n-select v-model:value="formData.obfuscationLevel" :options="obfuscationOptions" />
                    </n-form-item>

                    <n-form-item label="输出文件名" path="outputFilename">
                      <n-input v-model:value="formData.outputFilename" placeholder="可选，留空自动生成">
                        <template #prefix>
                          <n-icon :component="DocumentOutline" />
                        </template>
                      </n-input>
                    </n-form-item>

                    <n-space vertical style="margin-top: 24px">
                      <n-button type="primary" :loading="generating" @click="handleGenerate" block>
                        <template #icon>
                          <n-icon :component="FlashOutline" />
                        </template>
                        生成 Payload
                      </n-button>

                      <n-button @click="handlePreview" :disabled="!formData.type" block>
                        <template #icon>
                          <n-icon :component="EyeOutline" />
                        </template>
                        预览代码
                      </n-button>
                    </n-space>
                  </n-form>
                </n-card>
              </n-grid-item>

              <!-- 预览区域 -->
              <n-grid-item :span="12">
                <n-card title="代码预览" :bordered="false" class="preview-card">
                  <template #header-extra>
                    <n-space v-if="generatedResult">
                      <n-button quaternary size="small" @click="handleCopy">
                        <template #icon>
                          <n-icon :component="CopyOutline" />
                        </template>
                        复制
                      </n-button>
                      <n-button quaternary size="small" @click="handleDownload">
                        <template #icon>
                          <n-icon :component="DownloadOutline" />
                        </template>
                        下载
                      </n-button>
                    </n-space>
                  </template>

                  <div class="code-preview">
                    <n-scrollbar style="max-height: 500px">
                      <pre v-if="previewCode"><code>{{ previewCode }}</code></pre>
                      <n-empty v-else description="点击 &quot;预览代码&quot; 或 &quot;生成 Payload&quot; 查看结果" size="small" />
                    </n-scrollbar>
                  </div>

                  <n-divider v-if="generatedResult" />

                  <n-space v-if="generatedResult" vertical>
                    <n-statistic label="生成状态">
                      <n-tag :type="generatedResult.success ? 'success' : 'error'">
                        {{ generatedResult.success ? '成功' : '失败' }}
                      </n-tag>
                    </n-statistic>
                    <n-grid :cols="2" :x-gap="12">
                      <n-grid-item>
                        <n-statistic label="文件名">
                          <n-text depth="3">{{ generatedResult.filename }}</n-text>
                        </n-statistic>
                      </n-grid-item>
                      <n-grid-item>
                        <n-statistic label="文件大小">
                          <n-text depth="3">{{ formatFileSize(generatedResult.size) }}</n-text>
                        </n-statistic>
                      </n-grid-item>
                    </n-grid>
                  </n-space>
                </n-card>
              </n-grid-item>
            </n-grid>

            <!-- 内置模板列表 -->
            <n-card title="内置模板" :bordered="false" style="margin-top: 20px">
              <n-data-table :columns="templateColumns" :data="templateList" :pagination="false" size="small" />
            </n-card>
          </div>
        </n-tab-pane>

        <!-- Payload 列表标签页 -->
        <n-tab-pane name="list" tab="Payload 列表">
          <n-card :bordered="false">
            <template #header>
              <n-space justify="space-between">
                <n-text>已生成的 Payload</n-text>
                <n-space>
                  <n-input v-model:value="searchQuery" placeholder="搜索 Payload..." clearable style="width: 240px">
                    <template #prefix>
                      <n-icon :component="SearchOutline" />
                    </template>
                  </n-input>
                  <n-button type="primary" @click="loadPayloads">
                    <template #icon>
                      <n-icon :component="RefreshOutline" />
                    </template>
                    刷新
                  </n-button>
                </n-space>
              </n-space>
            </template>

            <n-data-table :columns="payloadColumns" :data="filteredPayloads" :row-key="row => row.id" :pagination="pagination" :scroll-x="1200" />
          </n-card>
        </n-tab-pane>

        <!-- 模板管理标签页 -->
        <n-tab-pane name="templates" tab="模板管理">
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

            <n-data-table :columns="templateColumns2" :data="customTemplates" :row-key="row => row.name" :pagination="false" />
          </n-card>
        </n-tab-pane>
      </n-tabs>
    </div>

    <!-- 创建模板弹窗 -->
    <n-modal v-model:show="showCreateTemplateModal" preset="dialog" title="创建自定义模板" style="width: 800px">
      <n-form ref="templateFormRef" :model="templateForm" :rules="templateFormRules" label-placement="top">
        <n-form-item label="模板名称" path="name">
          <n-input v-model:value="templateForm.name" placeholder="输入模板名称" />
        </n-form-item>

        <n-form-item label="Payload 类型" path="type">
          <n-select v-model:value="templateForm.type" :options="typeOptions" placeholder="选择类型" />
        </n-form-item>

        <n-form-item label="功能模式" path="function">
          <n-select v-model:value="templateForm.function" :options="functionOptions" placeholder="选择功能" />
        </n-form-item>

        <n-form-item label="模板内容" path="content">
          <n-input v-model:value="templateForm.content" type="textarea" :rows="15" placeholder="输入模板代码，使用 {{.password}} 作为密码占位符" style="font-family: monospace" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space justify="end">
          <n-button @click="showCreateTemplateModal = false">取消</n-button>
          <n-button type="primary" @click="handleCreateTemplate">创建</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useMessage } from 'naive-ui'
import {
  FlashOutline,
  KeyOutline,
  DocumentOutline,
  EyeOutline,
  CopyOutline,
  DownloadOutline,
  SearchOutline,
  RefreshOutline,
  AddOutline,
  TrashOutline,
} from '@vicons/ionicons5'
import type { FormRules, FormInst, DataTableColumns } from 'naive-ui'
import { Generate, GetTemplates, AddTemplate, DeleteTemplate } from '../../bindings/fg-abyss/internal/app/handlers/payloadhandler'

interface FormData {
  type: string
  function: string
  password: string
  encoder: string
  obfuscationLevel: string
  outputFilename: string
}

interface TemplateInfo {
  name: string
  type: string
  function: string
  description: string
}

interface PayloadItem {
  id: number
  name: string
  type: string
  function: string
  createdAt: string
  size: number
  content: string
}

interface GenerateResult {
  success: boolean
  content: string
  filename: string
  size: number
}

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const activeTab = ref('generator')
const generating = ref(false)
const previewCode = ref('')
const generatedResult = ref<GenerateResult | null>(null)
const templateList = ref<TemplateInfo[]>([])
const customTemplates = ref<TemplateInfo[]>([])
const searchQuery = ref('')
const showCreateTemplateModal = ref(false)

const formData = reactive<FormData>({
  type: 'php',
  function: 'basic',
  password: '',
  encoder: 'none',
  obfuscationLevel: 'none',
  outputFilename: '',
})

const templateFormRef = ref<FormInst | null>(null)
const templateForm = reactive({
  name: '',
  type: 'php',
  function: 'basic',
  content: '',
})

const formRules: FormRules = {
  type: { required: true, message: '请选择类型', trigger: 'change' },
  function: { required: true, message: '请选择功能', trigger: 'change' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
}

const templateFormRules: FormRules = {
  name: { required: true, message: '请输入模板名称', trigger: 'blur' },
  type: { required: true, message: '请选择类型', trigger: 'change' },
  function: { required: true, message: '请选择功能', trigger: 'change' },
  content: { required: true, message: '请输入模板内容', trigger: 'blur' },
}

const typeOptions = [
  { label: 'PHP', value: 'php' },
  { label: 'ASP', value: 'asp' },
  { label: 'ASPX', value: 'aspx' },
  { label: 'JSP', value: 'jsp' },
]

const functionOptions = [
  { label: '基础命令执行', value: 'basic' },
  { label: '完整功能（文件/命令/数据库）', value: 'full' },
]

const encoderOptions = [
  { label: '无编码', value: 'none' },
  { label: 'Base64', value: 'base64' },
  { label: 'ROT13', value: 'rot13' },
  { label: 'URL 编码', value: 'urlencode' },
  { label: '十六进制', value: 'hex' },
]

const obfuscationOptions = [
  { label: '无混淆', value: 'none' },
  { label: '轻度混淆', value: 'low' },
  { label: '中度混淆', value: 'medium' },
  { label: '高度混淆', value: 'high' },
]

const payloads = ref<PayloadItem[]>([])

const pagination = reactive({
  pageSize: 10,
  pageSizes: [10, 20, 50, 100],
  showSizePicker: true,
})

const filteredPayloads = computed(() => {
  if (!searchQuery.value) return payloads.value
  return payloads.value.filter((p) => p.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

const getTypeColor = (type: string): string => {
  const map: Record<string, string> = {
    php: 'success',
    asp: 'warning',
    aspx: 'info',
    jsp: 'error',
  }
  return map[type] || 'default'
}

const templateColumns: DataTableColumns = [
  { title: '名称', key: 'name', width: 150 },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row: TemplateInfo) => h('n-tag', { type: getTypeColor(row.type) }, () => row.type.toUpperCase()),
  },
  { title: '功能', key: 'function', width: 100 },
  { title: '描述', key: 'description' },
]

const templateColumns2: DataTableColumns = [
  { title: '名称', key: 'name', width: 150 },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row: TemplateInfo) => h('n-tag', { type: getTypeColor(row.type) }, () => row.type.toUpperCase()),
  },
  { title: '功能', key: 'function', width: 100 },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row: TemplateInfo) => h(
      'n-popconfirm',
      {
        onPositiveClick: () => handleDeleteTemplate(row),
      },
      {
        trigger: () => h(
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
    ),
  },
]

const payloadColumns: DataTableColumns = [
  { title: '名称', key: 'name', width: 200, ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row: PayloadItem) => h('n-tag', { type: getTypeColor(row.type) }, () => row.type.toUpperCase()),
  },
  { title: '功能', key: 'function', width: 100 },
  { title: '大小', key: 'size', width: 100, render: (row: PayloadItem) => formatFileSize(row.size) },
  { title: '创建时间', key: 'createdAt', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row: PayloadItem) => h('n-space', {}, [
      h('n-button', { size: 'small', quaternary: true }, {
        icon: () => h('n-icon', { component: EyeOutline }),
        default: () => '查看',
      }),
      h('n-button', { size: 'small', quaternary: true }, {
        icon: () => h('n-icon', { component: DownloadOutline }),
        default: () => '下载',
      }),
      h('n-popconfirm', {
        onPositiveClick: () => handleDeletePayload(row),
      }, {
        trigger: () => h('n-button', { size: 'small', quaternary: true, type: 'error' }, {
          icon: () => h('n-icon', { component: TrashOutline }),
          default: () => '删除',
        }),
        default: () => '确定要删除这个 Payload 吗？',
      }),
    ]),
  },
]

const handleGenerate = async () => {
  try {
    await formRef.value?.validate()
    generating.value = true

    const response = await Generate({
      type: formData.type,
      function: formData.function,
      password: formData.password,
      encoder: formData.encoder,
      obfuscation_level: formData.obfuscationLevel,
      output_filename: formData.outputFilename,
    })

    if (response.success) {
      previewCode.value = response.content || ''
      generatedResult.value = {
        success: response.success,
        content: response.content || '',
        filename: response.filename || '',
        size: response.size || 0,
      }
      message.success('Payload 生成成功')
    } else {
      message.error('生成失败：' + (response.message || '未知错误'))
    }
  } catch (error: any) {
    if (error.errors) return
    message.error(error.message || '生成失败')
  } finally {
    generating.value = false
  }
}

const handlePreview = async () => {
  try {
    await formRef.value?.validate()
    const response = await Generate({
      type: formData.type,
      function: formData.function,
      password: formData.password,
      encoder: formData.encoder,
      obfuscation_level: formData.obfuscationLevel,
      output_filename: formData.outputFilename,
    })

    previewCode.value = response.content || ''
    message.success('预览已更新')
  } catch (error: any) {
    if (error.errors) return
    message.error(error.message || '预览失败')
  }
}

const handleCopy = () => {
  if (!previewCode.value) return
  navigator.clipboard.writeText(previewCode.value)
  message.success('代码已复制')
}

const handleDownload = () => {
  if (!generatedResult.value) return
  const blob = new Blob([generatedResult.value.content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = generatedResult.value.filename
  a.click()
  URL.revokeObjectURL(url)
  message.success('下载已开始')
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
    templateForm.name = ''
    templateForm.type = 'php'
    templateForm.function = 'basic'
    templateForm.content = ''
  } catch (error: any) {
    if (error.errors) return
    message.error('创建模板失败：' + (error.message || '未知错误'))
  }
}

const handleDeleteTemplate = async (row: TemplateInfo) => {
  try {
    await DeleteTemplate(row.name)
    message.success('模板删除成功')
    loadTemplates()
  } catch (error: any) {
    message.error('删除模板失败：' + (error.message || '未知错误'))
  }
}

const handleDeletePayload = (row: PayloadItem) => {
  const index = payloads.value.findIndex((p) => p.id === row.id)
  if (index !== -1) {
    payloads.value.splice(index, 1)
    message.success('删除成功')
  }
}

const loadPayloads = () => {
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
    templateList.value = response.map((t) => ({
      name: t.name,
      type: t.type,
      function: t.function,
      description: getTemplateDescription(t.name),
    }))
    customTemplates.value = response.map((t) => ({
      ...t,
      isCustom: true,
    }))
  } catch (error) {
    console.error('加载模板失败:', error)
  }
}

const getTemplateDescription = (name: string): string => {
  const map: Record<string, string> = {
    'PHP Basic': 'PHP 基础命令执行',
    'PHP Full': 'PHP 完整功能',
    'ASP Basic': 'ASP 基础命令执行',
    'ASPX Basic': 'ASPX 基础命令执行',
    'JSP Basic': 'JSP 基础命令执行',
  }
  return map[name] || '未知模板'
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
.payloads-section {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;

  .content-header {
    flex-shrink: 0;
    padding: 24px 24px 20px 24px;
    border-bottom: 1px solid var(--border-color);
    background: var(--card-bg);
    box-sizing: border-box;
    display: flex;
    align-items: center;

    h1 {
      display: flex;
      align-items: center;
      gap: 12px;
      margin: 0;
      font-size: 0;
      line-height: 1;
    }

    .title {
      font-size: 24px;
      font-weight: 600;
      color: var(--active-color);
      letter-spacing: 0;
    }

    .separator {
      color: var(--text-tertiary);
      font-weight: 300;
      font-size: 20px;
    }

    .subtitle {
      font-size: 16px;
      font-weight: 400;
      color: var(--text-secondary);
      letter-spacing: 0;
    }
  }

  .content-body {
    flex: 1;
    overflow-y: auto;
    padding: 24px;

    :deep(.n-tabs) {
      height: 100%;
      display: flex;
      flex-direction: column;
    }

    :deep(.n-tabs-tab-wrapper) {
      flex-shrink: 0;
    }

    :deep(.n-tabs-content-wrapper) {
      flex: 1;
      overflow-y: auto;
    }

    :deep(.n-tab-pane) {
      padding: 0;
      height: 100%;
    }

    .generator-container {
      padding: 20px 0;

      .config-card,
      .preview-card {
        height: 100%;
        min-height: 400px;
      }

      .code-preview {
        background-color: var(--code-background);
        border-radius: 8px;
        padding: 16px;
        min-height: 200px;
        max-height: 500px;
        overflow-y: auto;

        pre {
          margin: 0;
          font-family: 'Courier New', monospace;
          font-size: 13px;
          line-height: 1.5;
          color: var(--code-text);
          white-space: pre-wrap;
          word-wrap: break-word;

          code {
            font-family: inherit;
          }
        }
      }
    }
  }
}
</style>
