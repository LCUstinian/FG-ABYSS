<template>
  <div class="payload-view">
    <!-- Left config panel -->
    <div class="config-panel">
      <div class="config-section">
        <label class="field-label">语言</label>
        <n-radio-group v-model:value="lang" name="lang" size="small">
          <n-radio-button value="PHP">PHP</n-radio-button>
          <n-radio-button value="JSP">JSP</n-radio-button>
          <n-radio-button value="ASP">ASP</n-radio-button>
          <n-radio-button value="ASPX">ASPX</n-radio-button>
        </n-radio-group>
      </div>

      <div class="config-section">
        <label class="field-label">加密方式</label>
        <n-radio-group v-model:value="encMode" name="enc">
          <div class="radio-list">
            <n-radio value="aes">AES-256-GCM <n-tag size="tiny" type="info">推荐</n-tag></n-radio>
            <n-radio value="xor">XOR <n-tag size="tiny">兼容</n-tag></n-radio>
            <n-radio value="none">无加密</n-radio>
          </div>
        </n-radio-group>
      </div>

      <div class="config-section">
        <label class="field-label">密钥</label>
        <div class="key-row">
          <n-input
            v-model:value="keyValue"
            :type="showKey ? 'text' : 'password'"
            placeholder="输入密钥"
            size="small"
            style="font-family: var(--font-mono); flex:1"
          />
          <n-button text size="small" :aria-label="showKey ? '隐藏密码' : '显示密码'" @click="showKey = !showKey">
            <EyeOff v-if="showKey" :size="14" />
            <Eye v-else :size="14" />
          </n-button>
          <n-button text size="small" aria-label="复制密钥" @click="copyKey">
            <Check v-if="keyCopied" :size="14" style="color: var(--color-success)" />
            <Copy v-else :size="14" />
          </n-button>
        </div>
        <!-- Password strength bar -->
        <div class="strength-bar">
          <div
            v-for="i in 4"
            :key="i"
            class="strength-seg"
            :class="{ 'seg-active': keyStrength >= i }"
            :style="{ background: strengthColor(keyStrength, i) }"
          />
        </div>
      </div>

      <div class="config-section">
        <label class="field-label">混淆等级</label>
        <div class="slider-row">
          <n-slider v-model:value="obfLevel" :min="0" :max="5" :step="1" style="flex:1" />
          <span class="level-badge">L{{ obfLevel }}</span>
        </div>
        <div class="obf-checks">
          <n-checkbox v-model:checked="obfVarRename"  :disabled="obfLevel === 0">变量重命名</n-checkbox>
          <n-checkbox v-model:checked="obfStrEncrypt" :disabled="obfLevel === 0">字符串加密</n-checkbox>
          <n-checkbox v-model:checked="obfJunk"       :disabled="obfLevel === 0">垃圾代码</n-checkbox>
        </div>
      </div>

      <n-button
        type="primary"
        block
        size="medium"
        style="margin-top: auto; height: 40px"
        :loading="generating"
        :style="generateBtnStyle"
        @click="generate"
      >
        生成载荷
      </n-button>
    </div>

    <!-- Right code panel -->
    <div class="code-panel">
      <div class="code-toolbar">
        <n-button text size="small" aria-label="复制代码" @click="copyCode">
          <Check v-if="codeCopied" :size="14" style="color: var(--color-success)" />
          <Copy v-else :size="14" />
        </n-button>
        <n-button text size="small" aria-label="下载">
          <Download :size="14" />
        </n-button>
      </div>
      <div class="code-area">
        <pre class="code-placeholder">{{ codeOutput || '// 点击「生成载荷」生成代码' }}</pre>
      </div>
      <div class="code-status">
        <span>{{ lang }} · {{ encLabel[encMode] }} · {{ codeOutput ? (codeOutput.length + ' B') : '--' }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Eye, EyeOff, Copy, Check, Download } from 'lucide-vue-next'

const lang     = ref<'PHP' | 'JSP' | 'ASP' | 'ASPX'>('PHP')
const encMode  = ref<'aes' | 'xor' | 'none'>('aes')
const keyValue = ref('')
const showKey  = ref(false)
const keyCopied = ref(false)
const obfLevel = ref(0)
const obfVarRename  = ref(true)
const obfStrEncrypt = ref(true)
const obfJunk       = ref(false)
const generating = ref(false)
const codeOutput = ref('')
const codeCopied = ref(false)
const generateError = ref(false)

const encLabel: Record<string, string> = {
  aes: 'AES-256', xor: 'XOR', none: '无加密',
}

const keyStrength = computed(() => {
  const k = keyValue.value
  if (k.length === 0) return 0
  if (k.length < 8)   return 1
  if (k.length < 16)  return 2
  const hasNum = /\d/.test(k)
  const hasSym = /[^a-zA-Z0-9]/.test(k)
  return hasNum && hasSym ? 4 : 3
})

function strengthColor(strength: number, seg: number): string {
  if (strength < seg) return 'var(--bg-hover)'
  if (strength === 1) return '#f87171'
  if (strength === 2) return '#fb923c'
  if (strength === 3) return '#fbbf24'
  return '#4ade80'
}

const generateBtnStyle = computed(() => ({
  ...(generateError.value ? { border: '1px solid #f87171' } : {}),
}))

async function copyKey() {
  if (!keyValue.value) return
  await navigator.clipboard.writeText(keyValue.value)
  keyCopied.value = true
  setTimeout(() => (keyCopied.value = false), 1500)
}

async function copyCode() {
  if (!codeOutput.value) return
  await navigator.clipboard.writeText(codeOutput.value)
  codeCopied.value = true
  setTimeout(() => (codeCopied.value = false), 1500)
}

async function generate() {
  generating.value = true
  generateError.value = false
  try {
    // Placeholder — actual Tauri invoke wired in payload feature plan
    await new Promise(r => setTimeout(r, 600))
    codeOutput.value = `<?php\n// Generated ${lang.value} payload\n// Key: ${keyValue.value || '(empty)'}\n// Enc: ${encMode.value}\n?>`
  } catch {
    generateError.value = true
  } finally {
    generating.value = false
  }
}
</script>

<style scoped>
.payload-view {
  height: 100%;
  display: flex;
  overflow: hidden;
}

.config-panel {
  width: 340px;
  flex-shrink: 0;
  border-right: 1px solid var(--border);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow-y: auto;
}

.config-section { display: flex; flex-direction: column; gap: 8px; }

.field-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-2);
}

.radio-list { display: flex; flex-direction: column; gap: 8px; }

.key-row { display: flex; align-items: center; gap: 4px; }

.strength-bar {
  display: flex;
  gap: 4px;
  height: 3px;
  margin-top: 6px;
}

.strength-seg {
  flex: 1;
  border-radius: 2px;
  background: var(--bg-hover);
  transition: background 200ms;
}

.slider-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.level-badge {
  min-width: 28px;
  text-align: center;
  background: var(--accent-bg);
  color: var(--accent);
  font-family: var(--font-mono);
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
}

.obf-checks { display: flex; flex-direction: column; gap: 6px; padding-top: 8px; }

.code-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--code-bg);
  border: 1px solid var(--code-border);
  margin: 16px;
  border-radius: 6px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.12), 0 0 0 1px rgba(0,0,0,0.08);
}

.code-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-bottom: 1px solid var(--code-border);
  background: var(--code-bg);
}

.code-area {
  flex: 1;
  overflow: auto;
  padding: 16px;
}

.code-placeholder {
  font-family: var(--font-mono);
  font-size: 13px;
  color: #abb2bf;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.code-status {
  padding: 6px 12px;
  font-size: 12px;
  color: var(--text-2);
  border-top: 1px solid var(--code-border);
  text-align: right;
  font-family: var(--font-mono);
  background: var(--code-bg);
}
</style>
