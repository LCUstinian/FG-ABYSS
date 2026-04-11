<template>
  <NModal
    :show="modelValue"
    @update:show="handleUpdateShow"
    preset="dialog"
    title="新建项目"
    :z-index="10000"
    :closable="true"
    class="create-project-modal"
  >
    <div class="modal-content">
      <NForm labelPlacement="top" class="project-form">
        <NFormItem label="项目名称" class="form-item">
          <NInput
            v-model:value="projectName"
            placeholder="请输入项目名称"
            class="project-input"
          />
        </NFormItem>
        <NFormItem label="项目描述" class="form-item">
          <NInput
            v-model:value="projectDescription"
            type="textarea"
            placeholder="请输入项目描述"
            class="project-textarea"
          />
        </NFormItem>
      </NForm>
    </div>
    <template #action>
      <div class="modal-actions">
        <NButton type="primary" @click="handleCreate" class="action-btn primary">
          创建项目
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NModal, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { invoke } from '@/utils/tauri-mock-adapter'
import { AuditLogger } from '@/utils/auditLogger'

defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'created'): void
}>()

const message = useMessage()

const projectName = ref('')
const projectDescription = ref('')

const handleUpdateShow = (value: boolean) => {
  emit('update:modelValue', value)
}

const handleCancel = () => {
  projectName.value = ''
  projectDescription.value = ''
  emit('update:modelValue', false)
}

const handleCreate = async () => {
  if (!projectName.value.trim()) {
    message.warning('请输入项目名称')
    return
  }

  try {
    await invoke('create_project', {
      name: projectName.value.trim(),
      description: projectDescription.value.trim()
    })
    // 记录审计日志
    await AuditLogger.logProjectCreate(projectName.value.trim())
    projectName.value = ''
    projectDescription.value = ''
    emit('update:modelValue', false)
    emit('created')
  } catch (error: any) {
    console.error('创建项目失败:', error)
    let errorMessage = ''
    
    if (error && typeof error === 'object') {
      if (error.message) {
        if (typeof error.message === 'string') {
          errorMessage = error.message
        } else if (typeof error.message === 'object' && error.message.message) {
          errorMessage = error.message.message
        } else {
          errorMessage = String(error.message)
        }
      } else if (error.error) {
        if (typeof error.error === 'string') {
          errorMessage = error.error
        } else if (typeof error.error === 'object' && error.error.message) {
          errorMessage = error.error.message
        } else {
          errorMessage = String(error.error)
        }
      } else {
        errorMessage = String(error)
      }
    } else if (typeof error === 'string') {
      errorMessage = error
    } else {
      errorMessage = String(error)
    }
    
    if (errorMessage.includes('UNIQUE constraint failed') && errorMessage.includes('projects.name')) {
      message.error('创建项目失败：项目名称已存在，请使用其他名称')
    } else if (errorMessage.includes('项目名称已存在')) {
      message.error('创建项目失败：项目名称已存在，请使用其他名称')
    } else if (errorMessage) {
      message.error('创建项目失败：' + errorMessage)
    } else {
      message.error('创建项目失败：' + String(error))
    }
  }
}
</script>

<style scoped>
.create-project-modal :deep(.n-dialog) {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: var(--shadow-xl);
  border: 1px solid var(--border-color);
  overflow: hidden;
}

.create-project-modal :deep(.n-dialog__header) {
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--active-color);
}

.create-project-modal :deep(.n-dialog__title) {
  color: #ffffff;
  font-size: 18px;
  font-weight: 600;
}

.create-project-modal :deep(.n-dialog__content) {
  padding: 24px;
  color: var(--text-color);
  background-color: var(--card-bg);
}

.create-project-modal :deep(.n-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  background-color: var(--content-bg);
}

.create-project-modal :deep(.n-dialog__close) {
  color: rgba(255, 255, 255, 0.9);
  transition: all 0.2s ease;
}

.create-project-modal :deep(.n-dialog__close:hover) {
  color: #ffffff;
  transform: scale(1.1);
}

.dark .create-project-modal :deep(.n-dialog__close) {
  color: rgba(255, 255, 255, 0.9);
}

.dark .create-project-modal :deep(.n-dialog__close:hover) {
  color: #ffffff;
}

.modal-content {
  padding: 0;
}

.project-form {
  width: 100%;
}

.form-item {
  margin-bottom: 20px;
}

.form-item:last-child {
  margin-bottom: 0;
}

.form-item :deep(.n-form-item-label__text) {
  color: var(--text-color);
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 8px;
}

.project-input,
.project-textarea {
  width: 100%;
  background-color: var(--content-bg);
  border: 1px solid var(--border-color);
  color: var(--text-color);
  border-radius: var(--border-radius-md);
  padding: 10px 12px;
  font-size: 14px;
  transition: all 0.2s ease;
}

.project-input:hover,
.project-textarea:hover {
  border-color: var(--active-color);
  background-color: var(--card-bg);
}

.project-input:focus,
.project-textarea:focus {
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px var(--active-color-suppl);
  background-color: var(--card-bg);
}

.project-input::placeholder,
.project-textarea::placeholder {
  color: var(--text-color);
  opacity: 0.5;
}

.project-textarea {
  min-height: 100px;
  resize: vertical;
  font-family: inherit;
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
  box-shadow: 0 2px 8px var(--active-color-suppl);
}

.action-btn.primary:focus {
  box-shadow: 0 0 0 3px var(--active-color-suppl);
}

.dark .create-project-modal :deep(.n-dialog__header) {
  background-color: var(--active-color);
}

.dark .create-project-modal :deep(.n-dialog) {
  background-color: var(--card-bg);
}

.dark .create-project-modal :deep(.n-dialog__content) {
  background-color: var(--card-bg);
}

.dark .create-project-modal :deep(.n-dialog__footer) {
  background-color: var(--content-bg);
}

.dark .project-input,
.dark .project-textarea {
  background-color: var(--content-bg);
}

.dark .project-input:hover,
.dark .project-textarea:hover,
.dark .project-input:focus,
.dark .project-textarea:focus {
  background-color: var(--card-bg);
}
</style>
