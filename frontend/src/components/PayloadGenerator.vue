<template>
  <div class="payload-generator">
    <n-grid :cols="24" :x-gap="20" :y-gap="20">
      <!-- 配置区域 -->
      <n-grid-item :span="12">
        <n-card title="Payload 配置" :bordered="false" class="config-card">
          <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top">
            <n-form-item label="Payload 类型" path="type">
              <n-select
                v-model:value="formData.type"
                :options="typeOptions"
                placeholder="选择 Payload 类型"
              />
            </n-form-item>

            <n-form-item label="功能模式" path="function">
              <n-select
                v-model:value="formData.function"
                :options="functionOptions"
                placeholder="选择功能模式"
              />
            </n-form-item>

            <n-form-item label="连接密码" path="password">
              <n-input
                v-model:value="formData.password"
                placeholder="设置连接密码"
                show-password-on="click"
              >
                <template #prefix>
                  <n-icon :component="KeyOutline" />
                </template>
              </n-input>
            </n-form-item>

            <n-form-item label="编码器" path="encoder">
              <n-select
                v-model:value="formData.encoder"
                :options="encoderOptions"
                placeholder="选择编码器"
              />
            </n-form-item>

            <n-form-item label="混淆级别" path="obfuscationLevel">
              <n-select
                v-model:value="formData.obfuscationLevel"
                :options="obfuscationOptions"
                placeholder="选择混淆级别"
              />
            </n-form-item>

            <n-form-item label="输出文件名" path="outputFilename">
              <n-input
                v-model:value="formData.outputFilename"
                placeholder="可选，留空自动生成"
              >
                <template #prefix>
                  <n-icon :component="DocumentOutline" />
                </template>
              </n-input>
            </n-form-item>

            <n-space vertical style="margin-top: 24px">
              <n-button
                type="primary"
                size="large"
                :loading="generating"
                @click="handleGenerate"
                block
              >
                <template #icon>
                  <n-icon :component="FlashOutline" />
                </template>
                生成 Payload
              </n-button>

              <n-button
                size="large"
                @click="handlePreview"
                :disabled="!formData.type"
                block
              >
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
            <n-space>
              <n-button
                v-if="generatedResult"
                quaternary
                size="small"
                @click="handleCopy"
              >
                <template #icon>
                  <n-icon :component="CopyOutline" />
                </template>
                复制
              </n-button>
              <n-button
                v-if="generatedResult"
                quaternary
                size="small"
                @click="handleDownload"
              >
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
              <n-empty
                v-else
                description="点击 &quot;预览代码&quot; 或 &quot;生成 Payload&quot; 查看结果"
                size="small"
              />
            </n-scrollbar>
          </div>

          <n-divider />

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
      <n-data-table
        :columns="templateColumns"
        :data="templateList"
        :row-key="row => row.name"
        :pagination="false"
        size="small"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useMessage } from 'naive-ui'
import {
  KeyOutline,
  DocumentOutline,
  FlashOutline,
  EyeOutline,
  CopyOutline,
  DownloadOutline,
} from '@vicons/ionicons5'
import type { FormRules, FormInst } from 'naive-ui'
import { 
  GeneratePayload, 
  GetTemplates,
  AddTemplate,
  DeleteTemplate,
} from '../../bindings/fg-abyss/internal/app/handlers/payloadhandler.js'

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

interface GenerateResult {
  success: boolean
  content: string
  filename: string
  size: number
}

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const generating = ref(false)
const previewCode = ref('')
const generatedResult = ref<GenerateResult | null>(null)
const templateList = ref<TemplateInfo[]>([])

const formData = reactive<FormData>({
  type: 'php',
  function: 'basic',
  password: '',
  encoder: 'base64',
  obfuscationLevel: 'none',
  outputFilename: '',
})

const formRules: FormRules = {
  type: {
    required: true,
    message: '请选择 Payload 类型',
    trigger: 'change',
  },
  function: {
    required: true,
    message: '请选择功能模式',
    trigger: 'change',
  },
  password: {
    required: true,
    message: '请输入连接密码',
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

const encoderOptions = [
  { label: 'Base64', value: 'base64' },
  { label: 'URL Encode', value: 'url' },
  { label: 'Raw (不编码)', value: 'raw' },
]

const obfuscationOptions = [
  { label: '无混淆', value: 'none' },
  { label: '低级混淆', value: 'low' },
  { label: '中级混淆', value: 'medium' },
  { label: '高级混淆', value: 'high' },
]

const templateColumns = [
  {
    title: '模板名称',
    key: 'name',
    width: 150,
  },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: (row: TemplateInfo) => {
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
    render: (row: TemplateInfo) => {
      const funcMap: Record<string, string> = {
        basic: '基础',
        full: '完整',
      }
      return funcMap[row.function] || row.function
    },
  },
  {
    title: '描述',
    key: 'description',
  },
]

const handleGenerate = async () => {
  try {
    await formRef.value?.validate()
    generating.value = true

    const response = await GeneratePayload({
      type: formData.type,
      function: formData.function,
      password: formData.password,
      encoder: formData.encoder,
      obfuscation_level: formData.obfuscationLevel,
      output_filename: formData.outputFilename,
    })

    if (response.success) {
      generatedResult.value = {
        success: response.success,
        content: response.content || '',
        filename: response.filename || '',
        size: response.size || 0,
      }
      previewCode.value = response.content || ''
      message.success('Payload 生成成功')
    } else {
      message.error('Payload 生成失败')
    }
  } catch (error: any) {
    message.error(error.message || '生成失败')
  } finally {
    generating.value = false
  }
}

const handlePreview = async () => {
  try {
    await formRef.value?.validate()
    generating.value = true

    const response = await GeneratePayload({
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
    message.error(error.message || '预览失败')
  } finally {
    generating.value = false
  }
}

const handleCopy = () => {
  if (!previewCode.value) return

  navigator.clipboard.writeText(previewCode.value)
  message.success('代码已复制到剪贴板')
}

const handleDownload = () => {
  if (!generatedResult.value) {
    message.warning('请先生成 Payload')
    return
  }

  const blob = new Blob([generatedResult.value.content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = generatedResult.value.filename
  a.click()
  URL.revokeObjectURL(url)

  message.success('文件下载已开始')
}

const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

const loadTemplates = async () => {
  try {
    const response = await GetTemplates()
    templateList.value = response.templates.map((t) => ({
      name: t.name,
      type: t.type,
      function: t.function,
      description: getTemplateDescription(t.name),
    }))
  } catch (error) {
    console.error('加载模板失败:', error)
  }
}

const getTemplateDescription = (name: string): string => {
  const descriptions: Record<string, string> = {
    'PHP Basic': 'PHP 基础命令执行',
    'PHP Full': 'PHP 完整功能（文件/命令/数据库）',
    ASP: 'ASP 基础命令执行',
    ASPX: 'ASPX 完整功能',
    JSP: 'JSP 完整功能',
  }
  return descriptions[name] || '未知模板'
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped lang="scss">
.payload-generator {
  padding: 20px;

  .config-card,
  .preview-card {
    height: 100%;
  }

  .code-preview {
    background-color: var(--code-background);
    border-radius: 8px;
    padding: 16px;
    margin-top: 12px;

    pre {
      margin: 0;
      font-family: 'Courier New', monospace;
      font-size: 13px;
      line-height: 1.5;
      color: var(--code-text);
      white-space: pre-wrap;
      word-wrap: break-word;
    }

    code {
      font-family: inherit;
    }
  }

  :deep(.n-statistic) {
    .n-statistic__label {
      font-size: 13px;
      color: var(--n-text-color-3);
    }

    .n-statistic__value {
      font-size: 15px;
      font-weight: 600;
    }
  }
}
</style>
