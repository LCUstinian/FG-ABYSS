<template>
  <div class="page-header">
    <h1>
      <span class="title">{{ title }}</span>
      <span class="separator">|</span>
      <span class="subtitle">{{ subtitle }}</span>
    </h1>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title: string
  subtitle?: string
}

withDefaults(defineProps<Props>(), {
  subtitle: ''
})
</script>

<style scoped>
.page-header {
  width: 100%;
  padding: 24px 24px 20px 24px;
  margin-bottom: 0;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  box-sizing: border-box;
}

.dark .page-header {
  border-bottom-color: var(--border-strong);
}

.page-header h1 {
  margin: 0;
  font-size: 0;
  line-height: 1;
  color: var(--text-color);
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-header h1 .title {
  font-size: 24px;
  font-weight: 600;
  color: var(--active-color);
  letter-spacing: 0;
}

.page-header h1 .separator {
  color: var(--text-tertiary);
  font-weight: 300;
  font-size: 20px;
}

.page-header h1 .subtitle {
  font-size: 16px;
  font-weight: 400;
  color: var(--text-secondary);
  letter-spacing: 0;
}
</style>
