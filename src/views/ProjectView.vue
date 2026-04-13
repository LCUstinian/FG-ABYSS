<template>
  <div class="project-view">
    <n-space vertical>
      <n-h1>{{ $t('project.manager') }}</n-h1>
      <n-empty 
        v-if="projectStore.projects.length === 0"
        description="暂无项目，点击右上角 '+' 创建新项目"
      >
        <template #extra>
          <n-button type="primary">
            <template #icon>
              <n-icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/></svg></n-icon>
            </template>
            新建项目
          </n-button>
        </template>
      </n-empty>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { useProjectStore } from '@/stores/project'

const projectStore = useProjectStore()
</script>

<style scoped>
.project-view {
  padding: 16px;
}
</style>
