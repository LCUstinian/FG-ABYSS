<template>
  <div class="payload-panel">
    <!-- 页面标题 -->
    <PageHeader 
      title="载荷" 
      subtitle="WebShell 生成器" 
    />

    <div class="content-body">
      <!-- Tab 切换 -->
      <n-tabs v-model:value="activeTab" type="line" animated>
        <n-tab-pane name="generator" tab="Payload 生成">
          <PayloadGeneratorView />
        </n-tab-pane>

        <n-tab-pane name="list" tab="Payload 列表">
          <PayloadListView />
        </n-tab-pane>

        <n-tab-pane name="templates" tab="Payload 模板">
          <PayloadTemplateManagerView />
        </n-tab-pane>
      </n-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import PageHeader from '@/components/shared/PageHeader.vue'
import PayloadGeneratorView from './PayloadGeneratorView.vue'
import PayloadListView from './PayloadListView.vue'
import PayloadTemplateManagerView from './PayloadTemplateManagerView.vue'

const activeTab = ref('generator')
</script>

<style scoped>
.payload-panel {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  background: var(--content-bg);
  transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.content-body {
  flex: 1;
  padding: 24px;
  background: var(--content-bg);
  animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

:deep(.n-tabs) {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

:deep(.n-tabs-nav) {
  background: var(--content-bg);
  border-bottom: 1px solid var(--border-color);
  padding: 0 24px;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

:deep(.n-tabs-tab) {
  padding: 18px 28px !important;
  font-size: 15px;
  font-weight: 500;
  color: var(--text-secondary);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  margin-right: 8px;
  position: relative;
  overflow: hidden;
}

:deep(.n-tabs-tab:hover) {
  color: var(--text-primary);
  background: rgba(128, 128, 128, 0.08);
}

:deep(.n-tabs-tab--active) {
  color: var(--active-color);
  font-weight: 600;
  background: rgba(var(--active-color-rgb), 0.08);
}

:deep(.n-tabs-tab-pane) {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  padding: 24px 0;
  background: transparent;
  animation: slideIn 0.3s ease-in-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

:deep(.n-tabs-content) {
  flex: 1;
  overflow: hidden;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

/* 自定义滚动条 */
:deep(.n-tabs-tab-pane)::-webkit-scrollbar {
  width: 8px;
}

:deep(.n-tabs-tab-pane)::-webkit-scrollbar-track {
  background: var(--sidebar-bg);
  border-radius: 4px;
}

:deep(.n-tabs-tab-pane)::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
  transition: background 0.3s ease;
}

:deep(.n-tabs-tab-pane)::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .content-body {
    padding: 32px 40px;
  }
  
  :deep(.n-tabs-nav) {
    padding: 0 40px;
  }
  
  :deep(.n-tabs-tab) {
    padding: 20px 32px !important;
    font-size: 16px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-body {
    padding: 16px;
  }
  
  :deep(.n-tabs-nav) {
    padding: 0 12px;
  }
  
  :deep(.n-tabs-tab) {
    padding: 14px 18px !important;
    font-size: 14px;
    margin-right: 4px;
  }
}
</style>
