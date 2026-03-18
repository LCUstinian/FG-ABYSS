<template>
  <div class="plugins-panel">
    <!-- 页面标题 -->
    <PageHeader 
      :title="t('plugins.title')" 
      :subtitle="t('plugins.subtitle')" 
    />

    <div class="content-body">
      <!-- PluginsManagementView 内容 -->
      <PluginsManagementView />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/shared/PageHeader.vue'
import PluginsManagementView from './PluginsManagementView.vue'

const { t } = useI18n()
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
