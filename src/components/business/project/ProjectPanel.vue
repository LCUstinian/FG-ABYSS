<template>
  <div class="project-panel">
    <!-- 页面标题 -->
    <PageHeader 
      title="项目" 
      subtitle="管理您的 WebShell 项目" 
    />

    <div class="content-body">
      <!-- ProjectListView 内容 -->
      <ProjectListView />
    </div>
  </div>
</template>

<script setup lang="ts">
import PageHeader from '@/components/shared/PageHeader.vue'
import ProjectListView from './ProjectListView.vue'
</script>

<style scoped>
.project-panel {
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
