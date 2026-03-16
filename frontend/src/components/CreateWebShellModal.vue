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
        <!-- 第一行：URL（全宽） -->
        <div class="form-row single-column">
          <NFormItem :label="t('projects.webShellUrl')" class="form-item">
            <NInput
              v-model:value="webshell.url"
              :placeholder="t('projects.webShellUrl')"
              class="webshell-input"
            />
          </NFormItem>
        </div>
        
        <!-- 第二行：Payload 类型和加密方式 -->
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
        
        <!-- 第三行：编码方式和代理类型 -->
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
        
        <!-- 第四行：备注（全宽） -->
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
import { NModal, NForm, NFormItem, NInput, NSelect, NButton } from 'naive-ui'
import * as WebShellHandler from '../../bindings/fg-abyss/internal/app/handlers/webshellhandler'
import type { WebShell } from '../../bindings/fg-abyss/internal/domain/entity/models'

const { t } = useI18n()

const props = defineProps<{
  modelValue: boolean
  projectId?: string  // 当前选中的项目 ID
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'created'): void
}>()

const loading = ref(false)

const webshell = reactive<Partial<WebShell>>({
  url: '',
  payload: 'php',
  cryption: 'none',
  encoding: 'UTF-8',
  proxyType: 'none',
  remark: '',
  status: 'active'
})

// 自动生成名称
const generateName = (url: string) => {
  if (!url) return ''
  try {
    // 从 URL 提取域名作为名称
    const urlObj = new URL(url.startsWith('http') ? url : `http://${url}`)
    const hostname = urlObj.hostname || urlObj.host
    return hostname || `WebShell_${Date.now()}`
  } catch {
    // 如果 URL 解析失败，使用时间戳生成名称
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

// 处理弹窗显示状态更新
const handleUpdateShow = (value: boolean) => {
  emit('update:modelValue', value)
}

// 处理取消按钮点击
const handleCancel = () => {
  resetForm()
  emit('update:modelValue', false)
}

// 重置表单
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
    alert('请输入目标 URL')
    return
  }

  loading.value = true

  try {
    // 使用传入的项目 ID，如果没有则报错
    if (!props.projectId) {
      throw new Error('未选择项目')
    }
    
    // 自动生成名称
    const autoName = generateName(webshell.url)
    
    // 调用后端创建方法（传递 8 个参数）
    await WebShellHandler.CreateWebShell(
      props.projectId,
      webshell.url || '',
      webshell.payload || 'php',
      webshell.cryption || 'none',
      webshell.encoding || 'UTF-8',
      webshell.proxyType || 'none',
      webshell.remark || '',
      webshell.status || 'active'
    )
    
    resetForm()
    emit('update:modelValue', false)
    emit('created')
  } catch (error) {
    console.error('创建 WebShell 失败:', error)
    alert('创建 WebShell 失败：' + error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 弹窗容器样式 - 使用主题变量 */
.create-webshell-modal :deep(.n-dialog) {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: var(--shadow-xl);
  border: 1px solid var(--border-color);
  overflow: hidden;
  max-width: 700px;
  width: 90%;
}

/* 弹窗头部样式 */
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

/* 弹窗内容样式 */
.create-webshell-modal :deep(.n-dialog__content) {
  padding: 20px 24px;
  color: var(--text-color);
  background-color: var(--card-bg);
  max-height: 75vh;
  overflow-y: auto;
}

/* 弹窗底部样式 */
.create-webshell-modal :deep(.n-dialog__footer) {
  padding: 12px 24px;
  border-top: 1px solid var(--border-color);
  background-color: var(--content-bg);
}

/* 关闭按钮样式 */
.create-webshell-modal :deep(.n-dialog__close) {
  color: rgba(255, 255, 255, 0.9);
  transition: all 0.2s ease;
}

.create-webshell-modal :deep(.n-dialog__close:hover) {
  color: #ffffff;
  transform: scale(1.1);
}

/* 表单样式 */
.modal-content {
  padding: 0;
}

.webshell-form {
  width: 100%;
}

/* 双列布局行 */
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

/* 单列布局行（全宽） */
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

/* 表单标签样式优化 */
.form-item :deep(.n-form-item-label) {
  margin-bottom: 8px !important;
  padding-bottom: 2px;
}

.form-item :deep(.n-form-item-label__text) {
  color: var(--text-color) !important;
  font-weight: 600 !important;
  font-size: 13px !important;
  line-height: 1.5 !important;
  letter-spacing: 0.02em !important;
  text-rendering: optimizeLegibility !important;
  -webkit-font-smoothing: antialiased !important;
  -moz-osx-font-smoothing: grayscale !important;
}

/* 表单反馈信息 */
.form-item :deep(.n-form-item-feedback) {
  margin-top: 6px !important;
  font-size: 12px !important;
  color: var(--text-color) !important;
  opacity: 0.8 !important;
}
/* 输入框样式 - 移除所有边框效果 */
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

/* 移除悬停时的边框颜色变化 */
.webshell-input:hover,
.webshell-select:hover,
.webshell-textarea:hover {
  background-color: var(--card-bg);
}

/* 移除聚焦时的边框颜色和光晕效果 */
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
  min-height: 60px;
  max-height: 100px;
  resize: vertical;
  font-family: inherit;
  height: auto;
  line-height: 1.5;
}

/* 下拉选择框样式优化 - 移除冗余输入框视觉 */
.webshell-select {
  position: relative;
  display: block;
  width: 100%;
}

/* 移除 Naive UI Select 的默认输入框样式 - 使用更高优先级 */
.webshell-select .n-select {
  background-color: var(--content-bg) !important;
  border: 1px solid var(--border-color) !important;
  border-radius: var(--border-radius-md) !important;
  font-size: 13px !important;
  width: 100% !important;
  height: 36px !important;
  transition: all 0.2s ease !important;
  box-sizing: border-box !important;
}

/* 移除内部输入框的背景和边框 */
.webshell-select .n-select .n-select__input {
  background-color: transparent !important;
  color: var(--text-color) !important;
  height: 36px !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 8px !important;
  outline: none !important;
}

/* 选中项显示 */
.webshell-select .n-select .n-select__selected-item {
  color: var(--text-color) !important;
  font-size: 13px !important;
  padding: 0 8px !important;
  display: flex;
  align-items: center;
  height: 100% !important;
  line-height: 36px !important;
}

/* 箭头图标样式 */
.webshell-select .n-select .n-select__arrow {
  color: var(--text-color) !important;
  opacity: 0.7 !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  font-size: 14px !important;
  display: flex;
  align-items: center;
  justify-content: center;
}

.webshell-select:hover .n-select .n-select__arrow {
  opacity: 1 !important;
  transform: translateY(1px) !important;
}

.webshell-select:active .n-select .n-select__arrow {
  transform: translateY(0) !important;
}

/* 下拉菜单容器 */
.webshell-select .n-select .n-select-menu {
  background-color: var(--card-bg) !important;
  border: 1px solid var(--border-color) !important;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12) !important;
  border-radius: var(--border-radius-md) !important;
  padding: 4px 0 !important;
  margin-top: 4px !important;
  backdrop-filter: blur(8px) !important;
  -webkit-backdrop-filter: blur(8px) !important;
  animation: selectMenuFadeIn 0.2s ease-out !important;
}

@keyframes selectMenuFadeIn {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 选项样式 */
.webshell-select .n-select .n-select-option {
  color: var(--text-color) !important;
  font-size: 13px !important;
  padding: 10px 16px !important;
  margin: 2px 8px !important;
  border-radius: 6px !important;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
  cursor: pointer !important;
}

.webshell-select .n-select .n-select-option:hover {
  background-color: var(--bg-hover) !important;
  transform: translateX(2px) !important;
}

.webshell-select .n-select .n-select-option--selected {
  background-color: rgba(59, 130, 246, 0.12) !important;
  color: var(--active-color) !important;
  font-weight: 600 !important;
}

.webshell-select .n-select .n-select-option--selected:hover {
  background-color: rgba(59, 130, 246, 0.2) !important;
}

/* 占位符文本 */
.webshell-select .n-select .n-select-placeholder {
  color: var(--text-color) !important;
  opacity: 0.5 !important;
  transition: opacity 0.2s ease !important;
}

.webshell-select:focus-within .n-select .n-select-placeholder {
  opacity: 0.3 !important;
}

/* 移除多余的边框和背景 */
.webshell-select .n-select .n-base-selection {
  background-color: transparent !important;
  border: none !important;
  box-shadow: none !important;
}

/* 移除聚焦状态的边框和光晕 */
.webshell-select:focus-within .n-select {
  box-shadow: none !important;
}

/* 移除悬停时的边框颜色变化 */
.webshell-select:hover .n-select {
  background-color: var(--card-bg) !important;
}

/* 深色模式优化 */
.dark .webshell-select :deep(.n-select-menu) {
  background-color: var(--bg-tertiary) !important;
  border-color: var(--border-strong) !important;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3) !important;
}

.dark .webshell-select :deep(.n-select-option:hover) {
  background-color: var(--bg-hover) !important;
}

.dark .webshell-input,
.dark .webshell-select,
.dark .webshell-textarea {
  background-color: rgba(15, 23, 42, 0.5) !important;
}

.dark .webshell-input:hover,
.dark .webshell-select:hover,
.dark .webshell-textarea:hover,
.dark .webshell-input:focus,
.dark .webshell-select:focus,
.dark .webshell-textarea:focus {
  background-color: rgba(30, 41, 59, 0.8) !important;
}

/* 浏览器兼容性优化 */
@supports (-webkit-backdrop-filter: blur(8px)) {
  .webshell-select :deep(.n-select-menu) {
    -webkit-backdrop-filter: blur(8px);
  }
}

/* 响应式布局优化 */
@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .form-row.single-column {
    grid-template-columns: 1fr;
  }
  
  /* 移动端优化下拉菜单 */
  .webshell-select :deep(.n-select-menu) {
    max-height: 300px !important;
  }
}

/* 确保表单项目垂直对齐 */
.form-item {
  align-items: stretch;
}

.form-item :deep(.n-form-item-feedback) {
  margin-top: 4px;
  font-size: 12px;
}

/* 输入框和选择框高度一致性保证 */
.webshell-input :deep(.n-input__input),
.webshell-select :deep(.n-select__input) {
  height: 100% !important;
  display: flex;
  align-items: center;
}

/* 标签与输入框的视觉连接 */
.form-item :deep(.n-form-item-label) {
  padding-bottom: 2px;
}

/* Firefox 特定优化 */
@-moz-document url-prefix() {
  .webshell-select :deep(.n-select__arrow) {
    font-size: 16px !important;
  }
}

/* Safari 特定优化 */
@media not all and (min-resolution:.001dpc) {
  .webshell-select :deep(.n-select-menu) {
    -webkit-backdrop-filter: blur(8px);
    backdrop-filter: blur(8px);
  }
}

/* Edge/Chromium 优化 */
@supports (-ms-ime-align:auto) {
  .webshell-select :deep(.n-select__input) {
    line-height: 36px;
  }
}

/* 按钮样式 */
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
  box-shadow: 0 2px 8px rgba(var(--active-color-rgb, 59, 130, 246), 0.3);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(var(--active-color-rgb, 59, 130, 246), 0.4);
  opacity: 0.95;
}

.action-btn.primary:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(var(--active-color-rgb, 59, 130, 246), 0.3);
}

.action-btn.primary:focus {
  box-shadow: 0 0 0 3px rgba(var(--active-color-rgb, 59, 130, 246), 0.2);
}

.action-btn.primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* 滚动条美化 */
.create-webshell-modal :deep(.n-dialog__content)::-webkit-scrollbar {
  width: 6px;
}

.create-webshell-modal :deep(.n-dialog__content)::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
}

.create-webshell-modal :deep(.n-dialog__content)::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.dark .create-webshell-modal :deep(.n-dialog__content)::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

.dark .create-webshell-modal :deep(.n-dialog__content)::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}

/* 深色主题优化 */
.dark .create-webshell-modal :deep(.n-dialog__header) {
  background-color: var(--active-color);
}

.dark .webshell-input,
.dark .webshell-select,
.dark .webshell-textarea {
  background-color: rgba(15, 23, 42, 0.5);
}

.dark .webshell-input:hover,
.dark .webshell-select:hover,
.dark .webshell-textarea:hover,
.dark .webshell-input:focus,
.dark .webshell-select:focus,
.dark .webshell-textarea:focus {
  background-color: rgba(30, 41, 59, 0.8);
}

/* 浅色主题优化 */
@media (prefers-color-scheme: light) {
  .create-webshell-modal :deep(.n-dialog__header) {
    background-color: #3b82f6;
  }
  
  .webshell-input,
  .webshell-select,
  .webshell-textarea {
    background-color: #f9fafb;
  }
  
  .webshell-input:hover,
  .webshell-select:hover,
  .webshell-textarea:hover,
  .webshell-input:focus,
  .webshell-select:focus,
  .webshell-textarea:focus {
    background-color: #ffffff;
  }
}
</style>
