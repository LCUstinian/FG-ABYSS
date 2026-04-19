<template>
  <div class="about-panel">
    <div class="settings-card app-info-card">
      <div class="app-logo">◈</div>
      <div class="app-title-group">
        <span class="app-title">渊渟 ABYSS</span>
        <span class="app-version">{{ version }}</span>
      </div>
      <p class="app-desc">现代化 WebShell 管理工具，基于 Tauri V2 + Rust + Vue 3 构建</p>
    </div>

    <div class="settings-card">
      <h3 class="card-title">系统信息</h3>
      <div class="info-rows">
        <div class="info-row">
          <span class="info-key">数据目录</span>
          <span class="info-val">{{ dataDir || '加载中…' }}</span>
        </div>
        <div class="info-row">
          <span class="info-key">日志目录</span>
          <span class="info-val">{{ logsDir || '加载中…' }}</span>
        </div>
        <div class="info-row">
          <span class="info-key">数据库</span>
          <span class="info-val">{{ dbPath || '加载中…' }}</span>
        </div>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">开源协议</h3>
      <p style="font-size:13px;color:var(--text-2);line-height:1.6">
        本软件基于 MIT 许可证开源发布。仅供授权测试、学习研究使用。
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { commands } from '@/bindings'

const version = ref('--')
const dataDir = ref('')
const logsDir = ref('')
const dbPath  = ref('')

onMounted(async () => {
  const res = await commands.getAppInfo()
  if (res.status === 'ok') {
    version.value = `v${res.data.version}`
    dataDir.value = res.data.dataDir
    logsDir.value = res.data.logsDir
    dbPath.value  = res.data.dbPath
  }
})
</script>

<style scoped>
.about-panel { display: flex; flex-direction: column; gap: 16px; }

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

.app-info-card { align-items: center; text-align: center; padding: 24px 16px; }

.app-logo {
  font-size: 36px;
  color: var(--accent);
  line-height: 1;
}

.app-title-group { display: flex; align-items: baseline; gap: 10px; }
.app-title { font-size: 18px; font-weight: 700; color: var(--text-1); }
.app-version {
  font-size: 12px;
  font-family: var(--font-mono);
  color: var(--text-3);
  background: var(--bg-hover);
  padding: 2px 8px;
  border-radius: 4px;
}

.app-desc { font-size: 13px; color: var(--text-2); line-height: 1.6; }

.info-rows { display: flex; flex-direction: column; gap: 8px; }
.info-row { display: flex; gap: 16px; font-size: 12px; }
.info-key { color: var(--text-3); min-width: 64px; flex-shrink: 0; }
.info-val {
  font-family: var(--font-mono);
  color: var(--text-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}
</style>
