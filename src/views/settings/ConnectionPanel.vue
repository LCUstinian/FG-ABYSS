<template>
  <div class="connection-panel">
    <!-- Proxy -->
    <div class="settings-card">
      <h3 class="card-title">代理设置</h3>
      <div class="field-row">
        <label class="field-label">代理类型</label>
        <n-radio-group v-model:value="proxyType" name="proxy-type" size="small">
          <n-radio-button value="none">不使用</n-radio-button>
          <n-radio-button value="http">HTTP</n-radio-button>
          <n-radio-button value="socks5">SOCKS5</n-radio-button>
        </n-radio-group>
      </div>
      <Transition name="proxy-fields">
        <div v-if="proxyType !== 'none'" class="proxy-fields">
          <div class="input-row">
            <div class="input-group" style="flex:1">
              <label class="input-label">主机</label>
              <n-input v-model:value="proxyHost" placeholder="127.0.0.1" size="small" style="font-family:var(--font-mono)" />
            </div>
            <div class="input-group" style="width:90px">
              <label class="input-label">端口</label>
              <n-input-number v-model:value="proxyPort" :min="1" :max="65535" size="small" placeholder="1080" style="font-family:var(--font-mono)" />
            </div>
          </div>
        </div>
      </Transition>
    </div>

    <!-- Timeout & Retry -->
    <div class="settings-card">
      <h3 class="card-title">连接参数</h3>
      <div class="field-row">
        <label class="field-label">超时时间</label>
        <div class="slider-row">
          <n-slider v-model:value="timeout" :min="1" :max="60" :step="1" style="flex:1" />
          <span class="value-badge">{{ timeout }}s</span>
        </div>
      </div>
      <div class="field-row">
        <label class="field-label">重试次数</label>
        <div style="display:flex;align-items:center;gap:12px">
          <span style="flex:1;font-size:13px;color:var(--text-2)">失败后自动重试</span>
          <n-input-number v-model:value="retryCount" :min="0" :max="5" size="small" style="width:80px" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const proxyType  = ref<'none' | 'http' | 'socks5'>('none')
const proxyHost  = ref('127.0.0.1')
const proxyPort  = ref(1080)
const timeout    = ref(30)
const retryCount = ref(3)
</script>

<style scoped>
.connection-panel { display: flex; flex-direction: column; gap: 16px; }

.settings-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

[data-theme="light"] .settings-card {
  border: none;
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}

.card-title { font-size: 14px; font-weight: 600; color: var(--text-1); }

.field-row { display: flex; flex-direction: column; gap: 8px; }
.field-label { font-size: 13px; font-weight: 500; color: var(--text-2); }

.proxy-fields {
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
}

.proxy-fields-enter-active { transition: opacity 200ms ease-out, max-height 200ms ease-out; max-height: 200px; }
.proxy-fields-leave-active { transition: opacity 140ms ease-in, max-height 140ms ease-in; }
.proxy-fields-enter-from, .proxy-fields-leave-to { opacity: 0; max-height: 0; }

.input-row { display: flex; gap: 12px; }
.input-group { display: flex; flex-direction: column; gap: 6px; }
.input-label { font-size: 12px; color: var(--text-3); }

.slider-row { display: flex; align-items: center; gap: 12px; }

.value-badge {
  min-width: 42px;
  text-align: center;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-2);
  background: var(--bg-hover);
  padding: 2px 6px;
  border-radius: 4px;
}
</style>
