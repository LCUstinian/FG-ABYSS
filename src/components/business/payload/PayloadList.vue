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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import PageHeader from '@/components/shared/PageHeader.vue'
import { useMessage } from 'naive-ui'
import { invoke } from '@/utils/tauri-mock-adapter'

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
  { label: 'XOR', value: 'xor' }
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
  { title: '名称', key: 'name', width: 200 },
  { title: '类型', key: 'type', width: 80 },
  { title: '功能', key: 'function', width: 100 },
  { title: '编码器', key: 'encoder', width: 100 },
  { title: '混淆级别', key: 'obfuscationLevel', width: 100 },
  { title: '文件大小', key: 'size', width: 100 },
  { title: '创建时间', key: 'createdAt', width: 180 },
  { title: '操作', key: 'actions', width: 150 }
])

const templateColumns = computed(() => [
  { title: '模板名称', key: 'name', width: 200 },
  { title: '类型', key: 'type', width: 80 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  { title: '操作', key: 'actions', width: 150 }
])

const handleGenerate = async () => {
  generating.value = true
  try {
    const result = await invoke('generate_payload', {
      type: formData.type,
      function: formData.function,
      password: formData.password,
      encoder: formData.encoder,
      obfuscationLevel: formData.obfuscationLevel,
      outputFilename: formData.outputFilename
    })

    if (result && result.success) {
      generatedResult.value = result
      previewCode.value = result.code || ''
      message.success('Payload 生成成功')
    } else {
      throw new Error('生成失败')
    }
  } catch (error: any) {
    message.error('Payload 生成失败：' + (error.message || '未知错误'))
  } finally {
    generating.value = false
  }
}

const handlePreview = () => {
  if (!formData.type) {
    message.warning('请先选择脚本类型')
    return
  }
  
  previewCode.value = `<?php
// Payload Preview
// Type: ${formData.type}
// Function: ${formData.function}
// Password: ${formData.password || 'none'}
// Encoder: ${formData.encoder}
// Obfuscation: ${formData.obfuscationLevel}

echo "Hello World";
`
  message.success('代码预览已更新')
}

const handleCopy = () => {
  if (previewCode.value) {
    navigator.clipboard.writeText(previewCode.value)
    message.success('代码已复制到剪贴板')
  }
}

const handleDownload = () => {
  message.info('下载功能开发中')
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
}

.content-body {
  flex: 1;
  padding: 24px;
  background: var(--content-bg);
  overflow: auto;
}

.generator-container {
  height: calc(100vh - 200px);
  min-height: 600px;
}

.config-card,
.preview-card {
  height: 100%;
  border-radius: 8px;
  box-shadow: var(--shadow-sm);
}

.generation-status {
  margin-bottom: 16px;
  padding: 16px;
  background: var(--content-bg);
  border-radius: var(--border-radius-md);
}

.status-content {
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-divider {
  width: 1px;
  height: 24px;
  background: var(--border-color);
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.code-preview {
  background: var(--card-bg);
  border-radius: var(--border-radius-md);
  padding: 16px;
  border: 1px solid var(--border-color);
}

.code-preview pre {
  margin: 0;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.code-preview code {
  color: var(--text-color);
}

:deep(.n-tabs) {
  height: 100%;
}

:deep(.n-tab-pane) {
  height: 100%;
}
</style>
