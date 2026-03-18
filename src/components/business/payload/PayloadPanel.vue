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

        <n-tab-pane name="templates" tab="模板管理">
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
  height: 100%;
  width: 100%;
  background: var(--content-bg);
}

.content-body {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: var(--content-bg);
}

:deep(.n-tabs) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.n-tabs-nav) {
  background: var(--content-bg);
  border-bottom: 1px solid var(--border-color-light);
  padding: 0;
  flex-shrink: 0;
}

:deep(.n-tabs-tab) {
  padding: 16px 24px !important;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  transition: all 0.3s ease;
  border-radius: 8px 8px 0 0;
  margin-right: 4px;
}

:deep(.n-tabs-tab:hover) {
  color: var(--text-primary);
  background: var(--card-bg-hover);
}

:deep(.n-tabs-tab--active) {
  color: var(--active-color);
  font-weight: 600;
}

:deep(.n-tabs-tab-pane) {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding: 24px 0;
  background: transparent;
}

:deep(.n-tabs-content) {
  flex: 1;
  overflow: auto;
  min-height: 0;
}
</style>
