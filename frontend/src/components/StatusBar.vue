<template>
  <div class="status-bar">
    <div class="status-item">
      <span>{{ t('status.memory') }}: {{ systemStatus.memoryUsage }}</span>
    </div>
    <div class="status-item">
      <span>{{ t('status.processId') }}: {{ systemStatus.processId }}</span>
    </div>
    <div class="status-item">
      <span>{{ t('status.cpu') }}: {{ systemStatus.cpuUsage }}</span>
    </div>
    <div class="status-item">
      <span>{{ t('status.uptime') }}: {{ systemStatus.uptime }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  systemStatus: {
    type: Object,
    required: true,
    default: () => ({
      memoryUsage: '0 GB / 0 GB',
      processId: '0',
      cpuUsage: '0%',
      uptime: '0 hours'
    })
  }
})
</script>

<style scoped>
.status-bar {
  height: 28px;
  background: var(--status-bar-bg);
  border-top: 1px solid var(--border-color);
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 20px;
  gap: 24px;
  font-size: 12px;
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  flex-wrap: wrap;
}

.status-item {
  display: flex;
  align-items: center;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .status-bar {
    gap: 12px;
    padding: 0 12px;
    font-size: 11px;
  }
  
  .status-item {
    white-space: nowrap;
  }
}
</style>