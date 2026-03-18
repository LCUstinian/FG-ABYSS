<template>
  <div class="home-panel">
    <!-- 页面标题 -->
    <PageHeader 
      title="首页" 
      subtitle="欢迎使用 FG-ABYSS WebShell 管理工具" 
    />

    <div class="content-body">
      <!-- DashboardView 内容 -->
      <DashboardView />
    </div>
  </div>
</template>

<script setup lang="ts">
import PageHeader from '@/components/shared/PageHeader.vue'
import DashboardView from './DashboardView.vue'
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
