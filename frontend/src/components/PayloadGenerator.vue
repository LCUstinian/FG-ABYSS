<template>
  <div class="payload-generator">
    <n-card title="WebShell 生成器">
      <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="left" :label-width="120">
        <n-form-item label="脚本类型" path="type">
          <n-select v-model:value="formData.type" :options="typeOptions" />
        </n-form-item>

        <n-form-item label="功能类型" path="function">
          <n-select v-model:value="formData.function" :options="functionOptions" />
        </n-form-item>

        <n-form-item label="连接密码" path="password">
          <n-input v-model:value="formData.password" placeholder="请输入密码" />
        </n-form-item>

        <n-form-item label="编码器" path="encoder">
          <n-select v-model:value="formData.encoder" :options="encoderOptions" />
        </n-form-item>

        <n-form-item label="混淆级别" path="obfuscationLevel">
          <n-select v-model:value="formData.obfuscationLevel" :options="obfuscationOptions" />
        </n-form-item>

        <n-form-item label="输出文件名" path="outputFilename">
          <n-input v-model:value="formData.outputFilename" placeholder="可选，留空自动生成" />
        </n-form-item>

        <n-space style="margin-top: 24px">
          <n-button type="primary" :loading="generating" @click="handleGenerate">
            生成 Payload
          </n-button>
          <n-button @click="handlePreview">预览代码</n-button>
        </n-space>
      </n-form>
    </n-card>

    <!-- 代码预览 -->
    <n-card title="代码预览" style="margin-top: 20px">
      <div class="code-preview">
        <pre v-if="previewCode"><code>{{ previewCode }}</code></pre>
        <n-empty v-else description="点击生成或预览按钮查看代码" />
      </div>
    </n-card>

    <!-- 内置模板 -->
    <n-card title="内置模板" style="margin-top: 20px">
      <n-data-table :columns="templateColumns" :data="templateList" :pagination="false" size="small" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormRules, FormInst, DataTableColumns } from 'naive-ui'
import { Generate, GetTemplates } from '../../bindings/fg-abyss/internal/app/handlers/payloadhandler'

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

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const generating = ref(false)
const previewCode = ref('')
const templateList = ref<TemplateInfo[]>([])

const formData = reactive<FormData>({
  type: 'php',
  function: 'basic',
  password: '',
  encoder: 'none',
  obfuscationLevel: 'none',
  outputFilename: '',
})

const formRules: FormRules = {
  type: { required: true, message: '请选择类型', trigger: 'change' },
  function: { required: true, message: '请选择功能', trigger: 'change' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
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

const getTypeColor = (type: string): string => {
  const map: Record<string, string> = {
    php: 'success',
    asp: 'warning',
    aspx: 'info',
    jsp: 'error',
  }
  return map[type] || 'default'
}

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

const loadTemplates = async () => {
  try {
    const response = await GetTemplates()
    templateList.value = response.map((t) => ({
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
  const map: Record<string, string> = {
    'PHP Basic': 'PHP 基础命令执行',
    'PHP Full': 'PHP 完整功能',
    'ASP Basic': 'ASP 基础命令执行',
    'ASPX Basic': 'ASPX 基础命令执行',
    'JSP Basic': 'JSP 基础命令执行',
  }
  return map[name] || '未知模板'
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped lang="scss">
.payload-generator {
  padding: 20px;

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
</style>
