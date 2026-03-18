<template>
  <NModal
    :show="modelValue"
    @update:show="handleUpdateShow"
    preset="dialog"
    :title="t('projects.recoverProjects')"
    :z-index="10100"
    :closable="true"
    class="recover-project-modal"
    :mask-closable="false"
  >
    <div class="modal-content">
      <div v-if="deletedProjects.length === 0" class="empty-state">
        <span class="empty-icon">📦</span>
        <p class="empty-text">{{ t('projects.noDeletedData') }}</p>
      </div>
      <div v-else class="recover-list">
        <div 
          v-for="project in deletedProjects" 
          :key="project.id"
          class="recover-item"
        >
          <span class="recover-item-icon">🗑️</span>
          <div class="recover-item-info">
            <div class="recover-item-name">{{ project.name }}</div>
            <div class="recover-item-time">
              {{ t('projects.deletedAt') }}: {{ formatTime(project.deletedAt) }}
            </div>
          </div>
          <NButton 
            size="small" 
            type="primary" 
            @click="showRecoverConfirm(project)"
            class="recover-item-btn"
          >
            <template #icon>
              <span class="btn-icon">↩️</span>
            </template>
            {{ t('projects.recover') }}
          </NButton>
        </div>
      </div>
    </div>
    <template #action>
      <div class="modal-actions">
        <NButton 
          v-if="deletedProjects.length > 0"
          type="success" 
          @click="handleRecoverAll" 
          class="recover-all-action-btn"
          :loading="recoveringAll"
        >
          <template #icon>
            <span class="btn-icon">↩️</span>
          </template>
          {{ t('projects.recoverAll') }}
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<script setup lang="ts">
import { NModal, NButton, useDialog, useMessage } from 'naive-ui'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { componentLogger } from '@/utils/logger'
import type { Project } from '@/types'
import { formatTime } from '@/utils/formatTime'

const { t } = useI18n()
const dialog = useDialog()
const message = useMessage()
const recoveringAll = ref(false)

const props = defineProps<{
  modelValue: boolean
  deletedProjects: Project[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'recover', project: any): void
  (e: 'recoverAll', projects: any[]): void
  (e: 'close'): void
}>()

const handleUpdateShow = (value: boolean) => {
  emit('update:modelValue', value)
}

const handleClose = () => {
  emit('update:modelValue', false)
  emit('close')
}

const showRecoverConfirm = (project: any) => {
  dialog.warning({
    title: t('projects.recoverProjectConfirm'),
    content: t('projects.recoverProjectConfirmContent', { name: project.name }),
    positiveText: t('projects.recover'),
    negativeText: t('projects.cancel'),
    onPositiveClick: () => {
      handleRecover(project)
    },
    zIndex: 10200
  })
}

const handleRecover = (project: any) => {
  emit('recover', project)
  emit('update:modelValue', false)
  emit('close')
}

const handleRecoverAll = async () => {
  if (props.deletedProjects.length === 0) return
  
  dialog.warning({
    title: t('projects.recoverAllConfirm'),
    content: t('projects.recoverAllConfirmContent', { count: props.deletedProjects.length }),
    positiveText: t('projects.recoverAll'),
    negativeText: t('projects.cancel'),
    onPositiveClick: async () => {
      recoveringAll.value = true
      
      try {
        const recoverPromises = props.deletedProjects.map(project => 
          emit('recover', project)
        )
        
        await Promise.all(recoverPromises)
        
        componentLogger.log('批量恢复成功:', props.deletedProjects.length)
        message.success(t('projects.recoverAllSuccess', { count: props.deletedProjects.length }))
        emit('update:modelValue', false)
        emit('close')
      } catch (error) {
        componentLogger.error('批量恢复失败:', error)
        message.error(t('projects.recoverAllError'))
      } finally {
        recoveringAll.value = false
      }
    },
    zIndex: 10200
  })
}
</script>

<style scoped>
.recover-project-modal :deep(.n-dialog) {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: var(--shadow-xl);
  border: 1px solid var(--border-color);
  overflow: hidden;
  max-width: 600px;
  width: 90%;
}

.recover-project-modal :deep(.n-dialog__header) {
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--warning-color);
}

.recover-project-modal :deep(.n-dialog__title) {
  color: #ffffff;
  font-size: 18px;
  font-weight: 600;
}

.recover-project-modal :deep(.n-dialog__content) {
  padding: 24px;
  color: var(--text-color);
  background-color: var(--card-bg);
  max-height: 400px;
  overflow-y: auto;
}

.recover-project-modal :deep(.n-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  background-color: var(--content-bg);
}

.recover-project-modal :deep(.n-dialog__close) {
  color: rgba(255, 255, 255, 0.9);
  transition: all 0.2s ease;
}

.recover-project-modal :deep(.n-dialog__close:hover) {
  color: #ffffff;
  transform: scale(1.1);
}

.dark .recover-project-modal :deep(.n-dialog__header) {
  background-color: var(--warning-color);
}

.modal-content {
  padding: 0;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 16px;
}

.empty-text {
  color: var(--text-secondary);
  font-size: 14px;
}

.recover-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recover-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: var(--border-radius-md);
  background: var(--content-bg);
  transition: all var(--transition-fast);
}

.recover-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
}

.recover-item-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.recover-item-info {
  flex: 1;
  min-width: 0;
}

.recover-item-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 4px;
}

.recover-item-time {
  font-size: 12px;
  color: var(--text-secondary);
}

.recover-item-btn {
  flex-shrink: 0;
}

.modal-actions {
  display: flex;
  justify-content: center;
  width: 100%;
  padding: 8px 0 4px 0;
}

.recover-all-action-btn {
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

.recover-all-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.recover-all-action-btn:active {
  transform: translateY(0);
}
</style>
