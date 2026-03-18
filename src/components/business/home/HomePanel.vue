<template>
  <div class="home-panel">
    <!-- 页面标题 -->
    <PageHeader 
      :title="t('home.title')" 
      :subtitle="t('home.subtitle')" 
    />

    <div class="content-body">
      <!-- DashboardView 内容 -->
      <DashboardView />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/shared/PageHeader.vue'
import DashboardView from './DashboardView.vue'

const { t } = useI18n()
</script>

<style scoped>
.home-panel {
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
