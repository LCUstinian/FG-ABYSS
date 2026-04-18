<template>
  <div class="statusbar" role="status" aria-live="polite">
    <div class="status-left">
      <span
        class="status-dot"
        :class="{ 'status-dot-active': activeCount > 0 }"
        :aria-label="`状态: ${activeCount > 0 ? '活跃' : '离线'}`"
      />
      <span class="status-text">{{ activeCount > 0 ? '已连接' : '未连接' }}</span>
    </div>

    <div class="divider" aria-hidden="true" />
    <span class="stat-item">
      项目 <span class="stat-num">{{ projectCount }}</span>
    </span>

    <div class="divider" aria-hidden="true" />
    <span class="stat-item">
      Shell <span class="stat-num">{{ shellCount }}</span>
    </span>

    <div class="divider" aria-hidden="true" />
    <span class="stat-item">
      载荷 <span class="stat-num">{{ payloadCount }}</span>
    </span>
  </div>
</template>

<script setup lang="ts">
// Mock data — will be replaced by store values in feature plans
const activeCount  = 0
const projectCount = 0
const shellCount   = 0
const payloadCount = 0
</script>

<style scoped>
.statusbar {
  height: 28px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 16px;
  background: var(--bg-deep);
  border-top: 1px solid var(--border);
  font-size: 12px;
  color: var(--text-2);
}

.status-left {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-3);
  flex-shrink: 0;
}

.status-dot.status-dot-active {
  background: var(--color-success);
}

.status-text {
  font-size: 12px;
}

.divider {
  width: 1px;
  height: 12px;
  background: var(--border);
}

.stat-item {
  font-size: 12px;
  color: var(--text-2);
}

.stat-num {
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  color: var(--text-1);
  margin-left: 4px;
}
</style>
