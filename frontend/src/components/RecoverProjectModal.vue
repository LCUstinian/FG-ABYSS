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

const { t } = useI18n()
const dialog = useDialog()
const message = useMessage()
const recoveringAll = ref(false)

const props = defineProps<{
  modelValue: boolean
  deletedProjects: any[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'recover', project: any): void
  (e: 'recoverAll', projects: any[]): void
  (e: 'close'): void
}>()

// 处理弹窗显示状态更新
const handleUpdateShow = (value: boolean) => {
  emit('update:modelValue', value)
}

// 处理关闭按钮
const handleClose = () => {
  emit('update:modelValue', false)
  emit('close')
}

// 显示恢复确认框
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

// 处理恢复
const handleRecover = (project: any) => {
  emit('recover', project)
  message.success(t('projects.recoverSuccess', { name: project.name }))
  emit('update:modelValue', false)
  emit('close')
}

// 处理恢复全部
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
        // 批量恢复所有项目
        const recoverPromises = props.deletedProjects.map(project => 
          emit('recover', project)
        )
        
        await Promise.all(recoverPromises)
        
        message.success(t('projects.recoverAllSuccess', { count: props.deletedProjects.length }))
        emit('update:modelValue', false)
        emit('close')
      } catch (error) {
        console.error('批量恢复失败:', error)
        message.error(t('projects.recoverAllError'))
      } finally {
        recoveringAll.value = false
      }
    },
    zIndex: 10200
  })
}

// 格式化时间
const formatTime = (time: string | number | Date) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}
</script>

<style scoped>
/* 弹窗容器样式 */
.recover-project-modal :deep(.n-dialog) {
  background-color: var(--card-bg);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12), 0 0 0 1px rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-color);
  overflow: hidden;
  max-width: 650px;
  width: 90%;
}

/* 弹窗头部样式 */
.recover-project-modal :deep(.n-dialog__header) {
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.12) 0%, rgba(76, 175, 80, 0.06) 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.recover-project-modal :deep(.n-dialog__title) {
  color: var(--text-color);
  font-size: 18px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 10px;
}

.recover-project-modal :deep(.n-dialog__title)::before {
  content: '↩️';
  font-size: 20px;
}

/* 弹窗内容样式 */
.recover-project-modal :deep(.n-dialog__content) {
  padding: 0;
  max-height: 500px;
  overflow-y: auto;
  background-color: var(--card-bg);
}

/* 自定义滚动条 */
.recover-project-modal :deep(.n-dialog__content::-webkit-scrollbar) {
  width: 6px;
}

.recover-project-modal :deep(.n-dialog__content::-webkit-scrollbar-track) {
  background: transparent;
}

.recover-project-modal :deep(.n-dialog__content::-webkit-scrollbar-thumb) {
  background: var(--border-color);
  border-radius: 3px;
}

.recover-project-modal :deep(.n-dialog__content::-webkit-scrollbar-thumb:hover) {
  background: var(--text-secondary);
}

/* 弹窗底部样式 */
.recover-project-modal :deep(.n-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  background-color: var(--content-bg);
}

/* 关闭按钮样式 */
.recover-project-modal :deep(.n-dialog__close) {
  color: var(--text-color);
  opacity: 0.6;
  transition: all 0.2s ease;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.recover-project-modal :deep(.n-dialog__close:hover) {
  opacity: 1;
  background-color: var(--hover-color);
  transform: rotate(90deg);
}

/* 空状态样式 */
.modal-content {
  padding: 0;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-color);
}

.empty-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
  opacity: 0.6;
  filter: grayscale(0.3);
}

.empty-text {
  font-size: 14px;
  color: var(--text-secondary);
  opacity: 0.8;
  margin: 0;
}

/* 恢复列表样式 */
.recover-list {
  display: flex;
  flex-direction: column;
  padding: 20px;
  gap: 12px;
}

.recover-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  background: linear-gradient(135deg, var(--bg-secondary) 0%, rgba(0, 0, 0, 0.02) 100%);
  border-radius: 10px;
  border: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.recover-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(180deg, rgba(239, 68, 68, 0.5) 0%, rgba(239, 68, 68, 0.2) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.recover-item:hover::before {
  opacity: 1;
}

.recover-item:hover {
  transform: translateX(6px);
  border-color: rgba(76, 175, 80, 0.4);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(76, 175, 80, 0.1);
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.05) 0%, var(--bg-secondary) 100%);
}

.dark .recover-item:hover {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(76, 175, 80, 0.2);
}

.recover-item-icon {
  font-size: 22px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.recover-item-info {
  flex: 1;
  min-width: 0;
}

.recover-item-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 6px;
  letter-spacing: -0.2px;
}

.recover-item-time {
  font-size: 12px;
  color: var(--text-secondary);
  opacity: 0.75;
  font-weight: 500;
}

/* 恢复按钮样式 */
.recover-item-btn {
  flex-shrink: 0;
  font-weight: 600;
  border-radius: 8px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.recover-item-btn :deep(.btn-icon) {
  font-size: 16px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.recover-item-btn:hover {
  transform: scale(1.06) translateY(-1px);
  box-shadow: 0 6px 16px rgba(76, 175, 80, 0.25);
}

/* 按钮动作区域 */
.modal-actions {
  display: flex;
  justify-content: center;
  width: 100%;
  padding: 8px 0;
}

/* 恢复全部按钮样式 */
.recover-all-action-btn {
  min-width: 160px;
  height: 40px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.95) 0%, rgba(76, 175, 80, 0.85) 100%);
  color: #ffffff;
  border: 1px solid rgba(76, 175, 80, 0.3);
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.25), inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.recover-all-action-btn:hover {
  background: linear-gradient(135deg, rgba(76, 175, 80, 1) 0%, rgba(76, 175, 80, 0.9) 100%);
  border-color: rgba(76, 175, 80, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(76, 175, 80, 0.35), inset 0 1px 0 rgba(255, 255, 255, 0.15);
}

.recover-all-action-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(76, 175, 80, 0.25), inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.recover-all-action-btn :deep(.btn-icon) {
  font-size: 18px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

/* 深色主题优化 */
.dark .recover-all-action-btn {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.85) 0%, rgba(76, 175, 80, 0.75) 100%);
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.3), inset 0 1px 0 rgba(255, 255, 255, 0.08);
}

.dark .recover-all-action-btn:hover {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.95) 0%, rgba(76, 175, 80, 0.85) 100%);
  box-shadow: 0 6px 20px rgba(76, 175, 80, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.12);
}

/* 浅色主题优化 */
@media (prefers-color-scheme: light) {
  .recover-all-action-btn {
    background: linear-gradient(135deg, rgba(76, 175, 80, 0.9) 0%, rgba(76, 175, 80, 0.8) 100%);
    box-shadow: 0 4px 12px rgba(76, 175, 80, 0.2), inset 0 1px 0 rgba(255, 255, 255, 0.15);
  }
  
  .recover-all-action-btn:hover {
    background: linear-gradient(135deg, rgba(76, 175, 80, 1) 0%, rgba(76, 175, 80, 0.9) 100%);
    box-shadow: 0 6px 20px rgba(76, 175, 80, 0.3), inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
}

/* 深色主题优化 */
.dark .recover-project-modal :deep(.n-dialog__header) {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.15) 0%, rgba(76, 175, 80, 0.08) 100%);
}

.dark .empty-icon {
  filter: grayscale(0.5) brightness(0.8);
}

/* 浅色主题优化 */
@media (prefers-color-scheme: light) {
  .recover-project-modal :deep(.n-dialog__header) {
    background: linear-gradient(135deg, rgba(76, 175, 80, 0.12) 0%, rgba(76, 175, 80, 0.06) 100%);
  }
  
  .recover-item {
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  }
  
  .recover-item:hover {
    background: linear-gradient(135deg, rgba(76, 175, 80, 0.08) 0%, #f8fafc 100%);
  }
}
</style>
