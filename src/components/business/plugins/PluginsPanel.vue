<template>
  <div class="plugins-panel">
    <!-- 页面标题 -->
    <PageHeader 
      title="插件" 
      subtitle="扩展功能与插件管理" 
    />

    <div class="content-body">
      <!-- PluginsManagementView 内容 -->
      <PluginsManagementView />
    </div>
  </div>
</template>

<script setup lang="ts">
import PageHeader from '@/components/shared/PageHeader.vue'
import PluginsManagementView from './PluginsManagementView.vue'
</script>

<style scoped>
.plugins-panel {
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
</style>
