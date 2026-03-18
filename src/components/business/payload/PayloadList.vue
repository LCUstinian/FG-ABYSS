<template>
  <div class="content-section payloads-section">
    <PageHeader title="载荷" subtitle="WebShell 生成器" />
    <div class="content-body">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <n-tab-pane name="generator" tab="Payload 生成">
          <div class="generator-container">
            <n-grid :cols="24" :x-gap="24" :y-gap="24">
              <n-grid-item :span="12">
                <n-card title="Payload 配置" :bordered="false" class="config-card">
                  <n-form ref="formRef" :model="formData" label-placement="left" :label-width="100">
                    <n-form-item label="脚本类型" path="type">
                      <n-select v-model:value="formData.type" :options="typeOptions" />
                    </n-form-item>

                    <n-form-item label="功能类型" path="function">
                      <n-select v-model:value="formData.function" :options="functionOptions" />
                    </n-form-item>

                    <n-form-item label="连接密码" path="password">
                      <n-input v-model:value="formData.password" placeholder="请输入密码">
                        <template #prefix>
                          <span>🔑</span>
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
                          <span>📄</span>
                        </template>
                      </n-input>
                    </n-form-item>

                    <n-space vertical style="margin-top: 24px">
                      <n-button type="primary" :loading="generating" @click="handleGenerate" block>
                        <template #icon>
                          <span>⚡</span>
                        </template>
                        生成 Payload
                      </n-button>

                      <n-button @click="handlePreview" :disabled="!formData.type" block>
                        <template #icon>
                          <span>👁️</span>
                        </template>
                        预览代码
                      </n-button>
                    </n-space>
                  </n-form>
                </n-card>
              </n-grid-item>

              <n-grid-item :span="12">
                <n-card title="Payload 预览" :bordered="false" class="preview-card">
                  <template #header-extra>
                    <n-space v-if="generatedResult">
                      <n-button quaternary size="small" @click="handleCopy">
                        <template #icon>
                          <span>📋</span>
                        </template>
                        复制
                      </n-button>
                      <n-button quaternary size="small" @click="handleDownload">
                        <template #icon>
                          <span>⬇️</span>
                        </template>
                        下载
                      </n-button>
                    </n-space>
                  </template>

                  <div v-if="generatedResult" class="generation-status">
                    <div class="status-content">
                      <n-space align="center" :size="16">
                        <n-tag :type="generatedResult.success ? 'success' : 'error'" size="medium" round>
                          <template #icon>
                            <span>{{ generatedResult.success ? '✓' : '✗' }}</span>
                          </template>
                          {{ generatedResult.success ? '生成成功' : '生成失败' }}
                        </n-tag>
                        
                        <div class="status-divider"></div>
                        
                        <div class="status-item">
                          <span style="font-size: 16px;">📄</span>
                          <n-text depth="2" style="margin-left: 6px;">
                            {{ generatedResult.filename }}
                          </n-text>
                        </div>
                        
                        <div class="status-divider"></div>
                        
                        <div class="status-item">
                          <span style="font-size: 16px;">⚡</span>
                          <n-text depth="2" style="margin-left: 6px;">
                            {{ formatFileSize(generatedResult.size) }}
                          </n-text>
                        </div>
                      </n-space>
                    </div>
                  </div>

                  <div class="code-preview">
                    <n-scrollbar style="max-height: 500px">
                      <pre v-if="previewCode"><code>{{ previewCode }}</code></pre>
                      <n-empty v-else description="点击 &quot;预览代码&quot; 或 &quot;生成 Payload&quot; 查看结果" size="small" />
                    </n-scrollbar>
                  </div>
                </n-card>
              </n-grid-item>
            </n-grid>
          </div>
        </n-tab-pane>

        <n-tab-pane name="list" tab="Payload 列表">
          <n-card :bordered="false" size="large">
            <template #header>
              <div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                <div>
                  <n-text style="font-size: 16px; font-weight: 600; color: var(--text-primary);">已生成的 Payload</n-text>
                  <n-text depth="3" style="font-size: 13px; margin-left: 12px;">共 {{ filteredPayloads.length }} 个</n-text>
                </div>
                <n-space>
                  <n-input v-model:value="searchQuery" placeholder="搜索 Payload..." clearable style="width: 240px">
                    <template #prefix>
                      <span>🔍</span>
                    </template>
                  </n-input>
                  <n-button type="primary" @click="loadPayloads">
                    <template #icon>
                      <span>🔄</span>
                    </template>
                    刷新
                  </n-button>
                </n-space>
              </div>
            </template>

            <n-data-table 
              :columns="payloadColumns" 
              :data="filteredPayloads" 
              :row-key="row => row.id" 
              :pagination="pagination" 
              :scroll-x="1200" 
              striped 
            />
          </n-card>
        </n-tab-pane>

        <n-tab-pane name="templates" tab="模板管理">
          <n-card :bordered="false" size="large">
            <template #header>
              <div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                <div>
                  <n-text style="font-size: 16px; font-weight: 600; color: var(--text-primary);">自定义模板</n-text>
                  <n-text depth="3" style="font-size: 13px; margin-left: 12px;">共 {{ customTemplates.length }} 个</n-text>
                </div>
                <n-space>
                  <n-popconfirm @positive-click="handleDeleteAllTemplates">
                    <template #trigger>
                      <n-button type="error" quaternary>
                        <template #icon>
                          <span>🗑️</span>
                        </template>
                        清空全部
                      </n-button>
                    </template>
                    确定要清空所有自定义模板吗？
                  </n-popconfirm>
                  <n-button type="primary" @click="showCreateTemplateModal = true">
                    <template #icon>
                      <span>➕</span>
                    </template>
                    创建模板
                  </n-button>
                </n-space>
              </div>
            </template>

            <n-data-table 
              :columns="templateColumns" 
              :data="customTemplates" 
              :row-key="row => row.name" 
              :pagination="false" 
              striped 
            />
          </n-card>
        </n-tab-pane>
      </n-tabs>

      <n-modal 
        v-model:show="showCreateTemplateModal" 
        preset="dialog" 
        title="创建自定义模板" 
        style="width: 800px;" 
        :close-on-esc="true" 
        :mask-closable="true"
      >
        <template #header>
          <div style="display: flex; align-items: center; gap: 12px;">
            <span style="font-size: 20px;">📝</span>
            <n-text style="font-size: 18px; font-weight: 600;">创建自定义模板</n-text>
          </div>
        </template>
        
        <n-form 
          ref="templateFormRef" 
          :model="templateForm" 
          :rules="templateFormRules" 
          label-placement="top" 
          label-width="100"
        >
          <n-grid :cols="2" :x-gap="20">
            <n-grid-item>
              <n-form-item label="模板名称" path="name">
                <n-input 
                  v-model:value="templateForm.name" 
                  placeholder="输入模板名称，例如：My Custom PHP Shell" 
                  clearable 
                />
              </n-form-item>
            </n-grid-item>
            <n-grid-item>
              <n-form-item label="Payload 类型" path="type">
                <n-select 
                  v-model:value="templateForm.type" 
                  :options="typeOptions" 
                  placeholder="选择类型" 
                />
              </n-form-item>
            </n-grid-item>
          </n-grid>

          <n-form-item label="功能模式" path="function">
            <n-select 
              v-model:value="templateForm.function" 
              :options="functionOptions" 
              placeholder="选择功能" 
            />
          </n-form-item>

          <n-form-item label="模板描述" path="description">
            <n-input 
              v-model:value="templateForm.description" 
              placeholder="可选，描述模板的用途" 
              clearable 
            />
          </n-form-item>

          <n-form-item label="模板内容" path="content">
            <n-input 
              v-model:value="templateForm.content" 
              type="textarea" 
              :rows="15" 
              placeholder="输入模板代码，使用 {{.Password}} 作为密码占位符" 
              style="font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 13px;"
              show-count
            />
          </n-form-item>
        </n-form>

        <template #action>
          <n-space justify="end">
            <n-button @click="showCreateTemplateModal = false">取消</n-button>
            <n-button type="primary" @click="handleCreateTemplate">
              <template #icon>
                <span>✓</span>
              </template>
              创建
            </n-button>
          </n-space>
        </template>
      </n-modal>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import PageHeader from '@/components/shared/PageHeader.vue'
import { useMessage } from 'naive-ui'
import { invoke, mockStore } from '@/utils/tauri-mock-adapter'
import {
  NTabs,
  NTabPane,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NButton,
  NSpace,
  NGrid,
  NGridItem,
  NDataTable,
  NModal,
  NScrollbar,
  NEmpty,
  NText,
  NTag,
  NPopconfirm
} from 'naive-ui'

const message = useMessage()
const activeTab = ref('generator')
const generating = ref(false)
const searchQuery = ref('')
const showCreateTemplateModal = ref(false)

const formData = reactive({
  type: 'php',
  function: 'basic',
  password: '',
  encoder: 'none',
  obfuscationLevel: 'low',
  outputFilename: ''
})

const typeOptions = [
  { label: 'PHP', value: 'php' },
  { label: 'ASP', value: 'asp' },
  { label: 'ASPX', value: 'aspx' },
  { label: 'JSP', value: 'jsp' }
]

const functionOptions = [
  { label: '基础连接', value: 'basic' },
  { label: '文件管理', value: 'file' },
  { label: '数据库管理', value: 'database' },
  { label: '命令执行', value: 'command' }
]

const encoderOptions = [
  { label: '无编码', value: 'none' },
  { label: 'Base64', value: 'base64' },
  { label: 'ROT13', value: 'rot13' },
  { label: 'XOR', value: 'xor' },
  { label: 'URL 编码', value: 'urlencode' },
  { label: 'Hex 编码', value: 'hex' }
]

const obfuscationOptions = [
  { label: '低', value: 'low' },
  { label: '中', value: 'medium' },
  { label: '高', value: 'high' }
]

const generatedResult = ref<any>(null)
const previewCode = ref('')
const payloads = ref<any[]>([])
const customTemplates = ref<any[]>([])

const filteredPayloads = computed(() => {
  if (!searchQuery.value) return payloads.value
  return payloads.value.filter(p => 
    p.name?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    p.type?.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const pagination = reactive({
  pageSize: 20,
  pageSizes: [20, 50, 100],
  showSizePicker: true
})

const payloadColumns = computed(() => [
  { title: '名称', key: 'name', width: 220 },
  { title: '类型', key: 'type', width: 80 },
  { title: '功能', key: 'function', width: 100 },
  { title: '编码器', key: 'encoder', width: 100 },
  { title: '混淆级别', key: 'obfuscationLevel', width: 100 },
  { title: '文件大小', key: 'size', width: 100 },
  { title: '创建时间', key: 'createdAt', width: 180 },
  { 
    title: '操作', 
    key: 'actions', 
    width: 180,
    render: (row: any) => [
      {
        type: 'button',
        text: '查看',
        props: {
          size: 'small',
          onClick: () => viewPayload(row)
        }
      },
      {
        type: 'button',
        text: '复制',
        props: {
          size: 'small',
          quaternary: true,
          onClick: () => copyPayload(row)
        }
      },
      {
        type: 'button',
        text: '删除',
        props: {
          size: 'small',
          quaternary: true,
          type: 'error',
          onClick: () => deletePayload(row)
        }
      }
    ]
  }
])

const templateColumns = computed(() => [
  { 
    title: '模板名称', 
    key: 'name', 
    width: 200,
    ellipsis: { tooltip: true }
  },
  { 
    title: '类型', 
    key: 'type', 
    width: 80,
    render: (row: any) => {
      const typeMap: Record<string, string> = {
        php: 'PHP',
        asp: 'ASP',
        aspx: 'ASPX',
        jsp: 'JSP'
      }
      return typeMap[row.type] || row.type
    }
  },
  { 
    title: '功能', 
    key: 'function', 
    width: 100,
    render: (row: any) => {
      const funcMap: Record<string, string> = {
        basic: '基础',
        file: '文件',
        database: '数据库',
        command: '命令'
      }
      return funcMap[row.function] || row.function
    }
  },
  { 
    title: '描述', 
    key: 'description', 
    ellipsis: { tooltip: true },
    width: 400
  },
  { 
    title: '创建时间', 
    key: 'createdAt', 
    width: 160,
    render: (row: any) => {
      if (!row.createdAt) return '-'
      try {
        return new Date(row.createdAt).toLocaleString('zh-CN')
      } catch {
        return row.createdAt
      }
    }
  },
  { 
    title: '操作', 
    key: 'actions', 
    width: 180,
    fixed: 'right',
    render: (row: any) => [
      {
        type: 'button',
        text: '删除',
        props: {
          size: 'small',
          quaternary: true,
          type: 'error',
          onClick: () => deleteTemplate(row)
        }
      }
    ]
  }
])

const deleteTemplate = (row: any) => {
  const index = mockStore.templates.findIndex((t: any) => t.name === row.name)
  if (index !== -1) {
    mockStore.templates.splice(index, 1)
    message.success('模板已删除')
    loadTemplates()
  }
}

const handleGenerate = async () => {
  generating.value = true
  try {
    const result = await invoke('generate_payload', {
      type: formData.type,
      function: formData.function,
      password: formData.password,
      encoder: formData.encoder,
      obfuscation_level: formData.obfuscationLevel,
      output_filename: formData.outputFilename
    })

    if (result && result.success) {
      generatedResult.value = {
        success: true,
        filename: result.filename,
        size: result.size,
        code: result.content
      }
      previewCode.value = result.content || ''
      message.success('Payload 生成成功')
      
      // 刷新 Payload 列表
      await loadPayloads()
    } else {
      throw new Error('生成失败')
    }
  } catch (error: any) {
    message.error('Payload 生成失败：' + (error.message || '未知错误'))
  } finally {
    generating.value = false
  }
}

const handlePreview = async () => {
  if (!formData.type) {
    message.warning('请先选择脚本类型')
    return
  }
  
  try {
    // 调用生成接口但不保存，仅用于预览
    const result = await invoke('generate_payload', {
      type: formData.type,
      function: formData.function,
      password: formData.password,
      encoder: formData.encoder,
      obfuscation_level: formData.obfuscationLevel,
      output_filename: 'preview.tmp'
    })
    
    if (result && result.success) {
      previewCode.value = result.content || ''
      message.success('代码预览已更新')
    }
  } catch (error: any) {
    message.error('预览失败：' + (error.message || '未知错误'))
  }
}

const handleCopy = () => {
  if (previewCode.value) {
    navigator.clipboard.writeText(previewCode.value)
    message.success('代码已复制到剪贴板')
  }
}

const copyPayload = (row: any) => {
  if (row.content) {
    navigator.clipboard.writeText(row.content)
    message.success('Payload 代码已复制')
  } else {
    message.warning('该 Payload 无可用代码')
  }
}

const viewPayload = (row: any) => {
  if (row.content) {
    previewCode.value = row.content
    generatedResult.value = {
      success: true,
      filename: row.name,
      size: row.size
    }
    message.success('已加载 Payload 代码')
  } else {
    message.warning('该 Payload 代码不可用')
  }
}

const deletePayload = (row: any) => {
  // 从 mockStore 中删除（实际应该调用后端 API）
  const index = mockStore.payloads.findIndex(p => p.id === row.id)
  if (index !== -1) {
    mockStore.payloads.splice(index, 1)
    message.success('Payload 已删除')
    // 刷新列表
    loadPayloads()
  }
}

const handleDownload = () => {
  if (generatedResult.value && previewCode.value) {
    const blob = new Blob([previewCode.value], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = generatedResult.value.filename || 'payload.txt'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    message.success('文件已开始下载')
  } else {
    message.warning('没有可下载的内容')
  }
}

const loadPayloads = async () => {
  try {
    const result = await invoke('get_payloads', {})
    if (result) {
      payloads.value = result.payloads || []
      message.success('Payload 列表已刷新')
    }
  } catch (error: any) {
    message.error('加载 Payload 列表失败')
  }
}

const loadTemplates = async () => {
  try {
    const result = await invoke('get_templates', {})
    if (result) {
      customTemplates.value = result.templates || []
    }
  } catch (error: any) {
    customTemplates.value = mockStore.templates || []
  }
}

const handleDeleteAllTemplates = async () => {
  try {
    await invoke('delete_all_templates', {})
    mockStore.templates = []
    customTemplates.value = []
    message.success('已清空所有自定义模板')
  } catch (error: any) {
    mockStore.templates = []
    customTemplates.value = []
    message.success('已清空所有自定义模板')
  }
}

const handleCreateTemplate = async () => {
  if (!templateForm.value.name || !templateForm.value.type) {
    message.warning('请填写模板名称和类型')
    return
  }
  
  try {
    await invoke('add_template', {
      name: templateForm.value.name,
      type: templateForm.value.type,
      function: templateForm.value.function || 'basic',
      content: templateForm.value.content || ''
    })
    
    mockStore.templates.push({
      name: templateForm.value.name,
      type: templateForm.value.type,
      function: templateForm.value.function || 'basic',
      description: templateForm.value.description || '',
      content: templateForm.value.content || '',
      createdAt: new Date().toISOString()
    })
    
    message.success('模板创建成功')
    showCreateTemplateModal.value = false
    loadTemplates()
  } catch (error: any) {
    mockStore.templates.push({
      name: templateForm.value.name,
      type: templateForm.value.type,
      function: templateForm.value.function || 'basic',
      description: templateForm.value.description || '',
      content: templateForm.value.content || '',
      createdAt: new Date().toISOString()
    })
    message.success('模板创建成功')
    showCreateTemplateModal.value = false
    loadTemplates()
  }
}

const templateForm = reactive({
  name: '',
  type: 'php',
  function: 'basic',
  content: '',
  description: ''
})

const templateFormRules = {
  name: {
    required: true,
    message: '请输入模板名称',
    trigger: 'blur'
  },
  type: {
    required: true,
    message: '请选择模板类型',
    trigger: 'change'
  },
  content: {
    required: true,
    message: '请输入模板内容',
    trigger: 'blur'
  }
}

const formatFileSize = (size: number) => {
  if (!size) return '0 B'
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  return `${(size / 1024 / 1024).toFixed(2)} MB`
}
</script>

<style scoped>
.payloads-section {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  background: var(--content-bg);
}

.content-body {
  flex: 1;
  padding: 24px;
  background: var(--content-bg);
  overflow: auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.generator-container {
  height: calc(100vh - 200px);
  min-height: 600px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.config-card,
.preview-card {
  height: 100%;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  border: 1px solid var(--border-color-light);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: var(--card-bg);
  overflow: hidden;
}

.config-card:hover,
.preview-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  border-color: var(--border-color);
}

.generation-status {
  margin-bottom: 16px;
  padding: 12px 20px;
  background: var(--card-bg-hover);
  border-radius: 12px 12px 0 0;
  border-bottom: 1px solid var(--border-color-light);
}

.status-content {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.status-divider {
  width: 1px;
  height: 20px;
  background: var(--border-color);
  margin: 0 12px;
  flex-shrink: 0;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
  flex-shrink: 0;
}

.status-item .n-text {
  overflow: hidden;
  text-overflow: ellipsis;
}

.code-preview {
  background: var(--code-background, var(--card-bg));
  border-radius: 10px;
  padding: 24px;
  border: 1px solid var(--border-color-light);
  flex: 1;
  min-height: 420px;
  max-height: none;
  overflow: auto;
  position: relative;
  transition: all 0.3s ease;
}

.code-preview:hover {
  border-color: var(--active-color);
  box-shadow: inset 0 0 0 1px var(--active-color);
}

.code-preview pre {
  margin: 0;
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-feature-settings: 'liga' 1;
  color: var(--code-text, var(--text-color));
}

.code-preview code {
  font-family: inherit;
}

.code-preview .n-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 320px;
  color: var(--text-tertiary);
}

.code-preview .n-empty :deep(.n-empty__icon) {
  font-size: 64px;
  opacity: 0.3;
  margin-bottom: 16px;
}

.code-preview .n-empty :deep(.n-empty__description) {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
}

:deep(.n-tabs) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.n-tabs-nav) {
  background: var(--content-bg);
  border-bottom: 1px solid var(--border-color-light);
  padding: 0;
  flex-shrink: 0;
}

:deep(.n-tabs-tab) {
  padding: 16px 24px !important;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  transition: all 0.3s ease;
  border-radius: 8px 8px 0 0;
  margin-right: 4px;
}

:deep(.n-tabs-tab:hover) {
  color: var(--text-primary);
  background: var(--card-bg-hover);
}

:deep(.n-tabs-tab--active) {
  color: var(--active-color);
  font-weight: 600;
}

:deep(.n-tabs-tab-pane) {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding: 24px;
  background: transparent;
}

:deep(.n-card__header) {
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color-light);
  background: var(--card-bg-hover);
  min-height: 64px;
  display: flex;
  align-items: center;
}

:deep(.n-card__header .n-card-header__title) {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.4;
  letter-spacing: 0.3px;
}

:deep(.n-card__content) {
  padding: 24px;
  flex: 1;
  overflow: auto;
  min-height: 0;
}

:deep(.n-data-table) {
  background: var(--card-bg);
  border-radius: 8px;
  overflow: hidden;
}

:deep(.n-data-table-th) {
  background: var(--card-bg-hover);
  font-weight: 600;
  font-size: 13px;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border-color);
}

:deep(.n-data-table-td) {
  font-size: 13px;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color-light);
  transition: all 0.2s ease;
}

:deep(.n-data-table-td:hover) {
  background: var(--table-hover-color);
}

:deep(.n-button) {
  border-radius: 8px;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.n-button--type-primary) {
  background: linear-gradient(135deg, var(--active-color) 0%, var(--active-color-hover) 100%);
  border: none;
  box-shadow: 0 2px 8px rgba(var(--active-color-rgb, 0, 128, 255), 0.3);
}

:deep(.n-button--type-primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(var(--active-color-rgb, 0, 128, 255), 0.4);
}

:deep(.n-button--type-primary:active) {
  transform: translateY(0);
}

:deep(.n-form-item) {
  margin-bottom: 22px;
}

:deep(.n-form-item-label__text) {
  font-weight: 500;
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.5;
}

:deep(.n-input),
:deep(.n-select) {
  border-radius: 8px;
  transition: all 0.3s ease;
  font-size: 14px;
}

:deep(.n-input:hover),
:deep(.n-select:hover) {
  border-color: var(--active-color);
}

:deep(.n-input--focus),
:deep(.n-select--focus) {
  border-color: var(--active-color);
  box-shadow: 0 0 0 2px var(--active-color-bg, rgba(0, 128, 255, 0.1));
}

:deep(.n-tag) {
  border-radius: 6px;
  font-weight: 500;
  font-size: 12px;
  padding: 4px 16px;
}

:deep(.n-divider) {
  margin: 16px 0;
  border-color: var(--border-color);
}

@media (max-width: 1200px) {
  .content-body {
    padding: 20px;
  }

  .generator-container {
    height: auto;
    min-height: auto;
  }

  .config-card,
  .preview-card {
    min-height: 450px;
  }

  .generation-status {
    padding: 10px 16px;
  }

  .status-item {
    max-width: 200px;
    font-size: 12px;
  }

  .status-divider {
    margin: 0 8px;
  }

  .code-preview {
    min-height: 350px;
    padding: 16px;
  }
}

@media (max-width: 768px) {
  .payloads-section {
    min-height: 100vh;
  }

  .content-body {
    padding: 16px;
  }

  :deep(.n-tabs-tab) {
    padding: 14px 16px !important;
    font-size: 13px;
    min-height: 48px;
  }

  .generator-container {
    gap: 16px;
  }

  .config-card,
  .preview-card {
    min-height: auto;
    height: auto;
  }

  .generation-status {
    padding: 10px 12px;
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .status-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .status-divider {
    width: 100%;
    height: 1px;
    margin: 8px 0;
  }

  .status-item {
    max-width: 100%;
    width: 100%;
  }

  .code-preview {
    min-height: 300px;
    padding: 12px;
  }

  :deep(.n-data-table) {
    font-size: 12px;
  }

  :deep(.n-button) {
    font-size: 13px;
    padding: 4px 12px;
  }
}

@media (max-width: 480px) {
  .content-body {
    padding: 12px;
  }

  :deep(.n-form-item) {
    margin-bottom: 16px;
  }

  :deep(.n-card__header) {
    padding: 16px;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  :deep(.n-card__content) {
    padding: 16px;
  }

  .code-preview {
    min-height: 250px;
    padding: 12px;
    font-size: 12px;
  }

  .code-preview pre {
    font-size: 12px;
    line-height: 1.5;
  }
}
</style>
