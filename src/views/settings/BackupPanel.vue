<template>
  <div class="backup-panel">
    <div class="settings-card">
      <h3 class="card-title">数据导出</h3>
      <p class="card-desc">将所有 WebShell、项目和配置导出为加密包</p>
      <div style="display:flex;gap:8px">
        <n-button size="small" type="primary" ghost>导出备份</n-button>
        <n-button size="small">导入备份</n-button>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">自动备份</h3>
      <div class="field-row">
        <span style="flex:1;font-size:13px;color:var(--text-2)">启用每日自动备份</span>
        <n-switch v-model:value="autoBackup" size="small" />
      </div>
      <div class="field-row" v-if="autoBackup">
        <div style="flex:1">
          <div class="field-label">最多保留</div>
        </div>
        <div style="display:flex;align-items:center;gap:8px">
          <n-input-number v-model:value="keepCount" :min="1" :max="30" size="small" style="width:80px" />
          <span style="font-size:13px;color:var(--text-3)">份</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const autoBackup = ref(false)
const keepCount  = ref(7)
</script>

<style scoped>
.backup-panel { display: flex; flex-direction: column; gap: 16px; }
.settings-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
[data-theme="light"] .settings-card {
  border: none;
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}
.card-title { font-size: 14px; font-weight: 600; color: var(--text-1); }
.card-desc  { font-size: 13px; color: var(--text-2); line-height: 1.6; }
.field-label { font-size: 13px; font-weight: 500; color: var(--text-2); }
.field-row { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
</style>
