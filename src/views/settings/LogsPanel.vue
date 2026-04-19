<template>
  <div class="logs-panel">
    <div class="settings-card">
      <h3 class="card-title">日志级别</h3>
      <n-radio-group v-model:value="logLevel" name="log-level" size="small">
        <n-radio-button value="error">Error</n-radio-button>
        <n-radio-button value="warn">Warn</n-radio-button>
        <n-radio-button value="info">Info</n-radio-button>
        <n-radio-button value="debug">Debug</n-radio-button>
      </n-radio-group>
    </div>

    <div class="settings-card">
      <h3 class="card-title">日志文件</h3>
      <div class="field-row">
        <div style="flex:1">
          <div class="field-label">单文件大小上限</div>
        </div>
        <div style="display:flex;align-items:center;gap:8px">
          <n-input-number v-model:value="maxSizeMb" :min="1" :max="100" size="small" style="width:80px" />
          <span style="font-size:13px;color:var(--text-3)">MB</span>
        </div>
      </div>
      <div class="field-row">
        <div style="flex:1">
          <div class="field-label">保留文件数</div>
        </div>
        <n-input-number v-model:value="maxFiles" :min="1" :max="30" size="small" style="width:80px" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const logLevel = ref('info')
const maxSizeMb = ref(10)
const maxFiles  = ref(7)
</script>

<style scoped>
.logs-panel { display: flex; flex-direction: column; gap: 16px; }
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
.field-label { font-size: 13px; font-weight: 500; color: var(--text-2); }
.field-row { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
</style>
