<template>
  <NModal
    :show="modelValue"
    @update:show="handleUpdateShow"
    preset="dialog"
    :title="t('projects.newWebShell')"
    :z-index="10000"
    :closable="true"
    class="create-webshell-modal"
  >
    <div class="modal-content">
      <NForm labelPlacement="top" class="webshell-form">
        <div class="form-row single-column">
          <NFormItem :label="t('projects.webShellUrl')" class="form-item">
            <NInput
              v-model:value="webshell.url"
              :placeholder="t('projects.webShellUrl')"
              class="webshell-input"
            />
          </NFormItem>
        </div>
        
        <div class="form-row">
          <NFormItem :label="t('projects.webShellPayload')" class="form-item">
            <NSelect
              v-model:value="webshell.payload"
              :options="payloadOptions"
              :placeholder="t('projects.webShellPayload')"
              class="webshell-select"
            />
          </NFormItem>
          <NFormItem :label="t('projects.webShellCryption')" class="form-item">
            <NSelect
              v-model:value="webshell.cryption"
              :options="cryptionOptions"
              :placeholder="t('projects.webShellCryption')"
              class="webshell-select"
            />
          </NFormItem>
        </div>
        
        <div class="form-row">
          <NFormItem :label="t('projects.webShellEncoding')" class="form-item">
            <NSelect
              v-model:value="webshell.encoding"
              :options="encodingOptions"
              :placeholder="t('projects.webShellEncoding')"
              class="webshell-select"
            />
          </NFormItem>
          <NFormItem :label="t('projects.webShellProxy')" class="form-item">
            <NSelect
              v-model:value="webshell.proxyType"
              :options="proxyOptions"
              :placeholder="t('projects.webShellProxy')"
              class="webshell-select"
            />
          </NFormItem>
        </div>
        
        <NFormItem :label="t('projects.webShellRemark')" class="form-item full-width">
          <NInput
            v-model:value="webshell.remark"
            type="textarea"
            :placeholder="t('projects.webShellRemark')"
            class="webshell-textarea"
          />
        </NFormItem>
      </NForm>
    </div>
    <template #action>
      <div class="modal-actions">
        <NButton type="primary" @click="handleCreate" class="action-btn primary" :loading="loading">
          {{ loading ? (t('projects.creating') || '创建中...') : (t('projects.createWebShell') || '创建 WebShell') }}
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { NModal, NForm, NFormItem, NInput, NSelect, NButton, useMessage } from 'naive-ui'
import { invoke } from '@/utils/tauri-mock-adapter'
import { validateUrl } from '@/utils/urlValidator'
import { componentLogger } from '@/utils/logger'
import { AuditLogger } from '@/utils/auditLogger'

const { t } = useI18n()
const message = useMessage()

const props = defineProps<{
  modelValue: boolean
  projectId?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'created'): void
}>()

const loading = ref(false)

const webshell = reactive<Partial<any>>({
  url: '',
  payload: 'php',
  cryption: 'none',
  encoding: 'UTF-8',
  proxyType: 'none',
  remark: '',
  status: 'active'
})

const generateName = (url: string) => {
  if (!url) return ''
  try {
    const urlObj = new URL(url.startsWith('http') ? url : `http://${url}`)
    const hostname = urlObj.hostname || urlObj.host
    return hostname || `WebShell_${Date.now()}`
  } catch {
    return `WebShell_${Date.now()}`
  }
}

const payloadOptions = [
  { label: 'PHP', value: 'php' },
  { label: 'ASP', value: 'asp' },
  { label: 'ASPX', value: 'aspx' },
  { label: 'JSP', value: 'jsp' }
]

const cryptionOptions = [
  { label: '无加密', value: 'none' },
  { label: 'Base64', value: 'base64' },
  { label: 'ROT13', value: 'rot13' },
  { label: 'XOR', value: 'xor' }
]

const encodingOptions = [
  { label: 'UTF-8', value: 'UTF-8' },
  { label: 'GBK', value: 'GBK' },
  { label: 'BIG5', value: 'BIG5' }
]

const proxyOptions = [
  { label: '无代理', value: 'none' },
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' }
]

const handleUpdateShow = (value: boolean) => {
  emit('update:modelValue', value)
}

const resetForm = () => {
  webshell.url = ''
  webshell.payload = 'php'
  webshell.cryption = 'none'
  webshell.encoding = 'UTF-8'
  webshell.proxyType = 'none'
  webshell.remark = ''
  webshell.status = 'active'
}

const handleCreate = async () => {
  if (!webshell.url?.trim()) {
    message.error('请输入目标 URL')
    return
  }

  const urlValidation = validateUrl(webshell.url.trim(), {
    allowedProtocols: ['http:', 'https:'],
    allowLocalhost: true,
    allowInternalIPs: false,
    allowIPAddresses: true
  })

  if (!urlValidation.valid) {
    message.error(urlValidation.error || 'URL 格式不正确')
    return
  }

  loading.value = true

  try {
    if (!props.projectId) {
      throw new Error('未选择项目')
    }
    
    const sanitizedUrl = urlValidation.sanitized || webshell.url.trim()
    const autoName = generateName(sanitizedUrl)
    
    componentLogger.log('创建 WebShell:', { url: sanitizedUrl, payload: webshell.payload, projectId: props.projectId })
    
    const result = await invoke('create_webshell', {
      projectId: props.projectId,
      url: sanitizedUrl,
      payload: webshell.payload || 'php',
      cryption: webshell.cryption || 'none',
      encoding: webshell.encoding || 'UTF-8',
      proxyType: webshell.proxyType || 'none',
      remark: webshell.remark || '',
      status: webshell.status || 'active'
    })
    
    if (result && result.success) {
      message.success('WebShell 创建成功')
      // 记录审计日志
      await AuditLogger.logWebShellCreate(`项目 ${props.projectId}`, sanitizedUrl)
      resetForm()
      emit('update:modelValue', false)
      emit('created')
    } else {
      throw new Error(result?.message || '创建失败')
    }
  } catch (error) {
    componentLogger.error('创建 WebShell 失败:', error)
    const errorMessage = error instanceof Error ? error.message : '创建 WebShell 失败'
    message.error(errorMessage)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.create-webshell-modal :deep(.n-dialog) {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: var(--shadow-xl);
  border: 1px solid var(--border-color);
  overflow: hidden;
  max-width: 700px;
  width: 90%;
}

.create-webshell-modal :deep(.n-dialog__header) {
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--active-color);
}

.create-webshell-modal :deep(.n-dialog__title) {
  color: #ffffff;
  font-size: 16px;
  font-weight: 600;
}

.create-webshell-modal :deep(.n-dialog__content) {
  padding: 20px 24px;
  color: var(--text-color);
  background-color: var(--card-bg);
  max-height: 75vh;
  overflow-y: auto;
}

.create-webshell-modal :deep(.n-dialog__footer) {
  padding: 12px 24px;
  border-top: 1px solid var(--border-color);
  background-color: var(--content-bg);
}

.create-webshell-modal :deep(.n-dialog__close) {
  color: rgba(255, 255, 255, 0.9);
  transition: all 0.2s ease;
}

.create-webshell-modal :deep(.n-dialog__close:hover) {
  color: #ffffff;
  transform: scale(1.1);
}

.modal-content {
  padding: 0;
}

.webshell-form {
  width: 100%;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.form-row.single-column {
  grid-template-columns: 1fr;
  margin-bottom: 20px;
}

.form-row:last-child {
  margin-bottom: 0;
}

.form-item {
  margin-bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.form-item.full-width {
  margin-top: 0;
}

.form-item :deep(.n-form-item-label) {
  margin-bottom: 8px !important;
  padding-bottom: 2px;
}

.form-item :deep(.n-form-item-label__text) {
  color: var(--text-color) !important;
  font-weight: 600 !important;
  font-size: 13px !important;
}

.webshell-input,
.webshell-select,
.webshell-textarea {
  width: 100%;
  background-color: var(--content-bg);
  border: 1px solid var(--border-color);
  color: var(--text-color);
  border-radius: var(--border-radius-md);
  padding: 8px 12px;
  font-size: 13px;
  transition: all 0.2s ease;
  height: 36px;
  box-sizing: border-box;
}

.webshell-input:hover,
.webshell-select:hover,
.webshell-textarea:hover {
  background-color: var(--card-bg);
}

.webshell-input:focus,
.webshell-select:focus,
.webshell-textarea:focus {
  background-color: var(--card-bg);
  outline: none;
  box-shadow: none;
}

.webshell-input::placeholder,
.webshell-textarea::placeholder {
  color: var(--text-color);
  opacity: 0.5;
}

.webshell-textarea {
  min-height: 80px;
  resize: vertical;
}

.modal-actions {
  display: flex;
  justify-content: center;
  width: 100%;
  padding: 8px 0 4px 0;
}

.action-btn {
  min-width: 140px;
  height: 42px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  letter-spacing: 0.5px;
  box-shadow: var(--shadow-sm);
}

.action-btn.primary {
  background-color: var(--active-color);
  border: none;
  color: #ffffff;
  box-shadow: 0 2px 8px var(--active-color-suppl);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px var(--active-color-suppl);
  opacity: 0.95;
}

.action-btn.primary:active {
  transform: translateY(0);
}

.dark .create-webshell-modal :deep(.n-dialog__header) {
  background-color: var(--active-color);
}
</style>
