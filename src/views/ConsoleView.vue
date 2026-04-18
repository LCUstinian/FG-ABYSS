<template>
  <div class="console-view">
    <!-- Console header -->
    <div class="console-header">
      <div class="console-shell-info">
        <span class="console-icon" aria-hidden="true">◈</span>
        <span class="console-name">{{ shellName }}</span>
        <span class="console-url">{{ shellUrl }}</span>
        <span
          class="console-dot"
          :class="{ 'status-dot-active': isActive }"
          aria-label="状态: 活跃"
        />
      </div>
      <button class="console-close" aria-label="关闭" @click="closeWindow">
        <X :size="14" />
      </button>
    </div>

    <!-- Tab bar -->
    <n-tabs
      v-model:value="activeTab"
      type="card"
      size="small"
      class="console-tabs"
    >
      <n-tab-pane name="file" tab="文件管理">
        <div class="tab-content-placeholder">
          文件管理器将在 Console 功能计划中实现
        </div>
      </n-tab-pane>
      <n-tab-pane name="database" tab="数据库">
        <div class="tab-content-placeholder">
          数据库管理器将在 Console 功能计划中实现
        </div>
      </n-tab-pane>
      <n-tab-pane name="terminal" tab="终端">
        <div class="tab-content-placeholder">
          xterm.js 终端将在 Console 功能计划中实现
        </div>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { X } from 'lucide-vue-next'
import { getCurrentWindow } from '@tauri-apps/api/window'

// Console window reads webshell_id from URL query param (hash router)
// URL format: tauri://localhost/#/console?id={webshell_id}
const route = useRoute()
const webshellId = route.query.id as string | undefined

const shellName = webshellId ? `shell-${webshellId.slice(0, 8)}` : 'Unknown'
const shellUrl  = 'https://example.com/shell.php'
const isActive  = ref(true)
const activeTab = ref('terminal')

async function closeWindow() {
  await getCurrentWindow().close()
}
</script>

<style scoped>
.console-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-base);
}

.console-header {
  height: 32px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background: var(--bg-deep);
  border-bottom: 1px solid var(--border);
  user-select: none;
}

.console-shell-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.console-icon { color: var(--accent); }
.console-name { color: var(--text-1); font-weight: 600; }
.console-url  { color: var(--text-2); font-family: var(--font-mono); font-size: 12px; }

.console-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-success);
}

.console-close {
  width: 32px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-2);
  cursor: pointer;
  border-radius: 4px;
}

.console-close:hover { background: var(--wc-close-bg); color: #fff; }

.console-tabs {
  flex: 1;
  overflow: hidden;
}

.tab-content-placeholder {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-3);
  font-size: 13px;
  padding: 40px;
}
</style>
