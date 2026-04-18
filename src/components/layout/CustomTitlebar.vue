<template>
  <div class="titlebar" data-tauri-drag-region>
    <div class="titlebar-left" data-tauri-drag-region>
      <div class="logo" aria-hidden="true">◈</div>
      <span class="app-name">{{ t('app.name') }}</span>
      <span class="app-sub">{{ t('app.sub') }}</span>
    </div>

    <div class="titlebar-right" data-tauri-drag-region="false">
      <!-- Theme toggle -->
      <button
        class="tb-icon-btn"
        :aria-label="`切换主题: ${modeLabel}`"
        :title="`切换主题: ${modeLabel}`"
        @click="cycleMode"
      >
        <Moon v-if="themeStore.mode === 'dark'" :size="16" />
        <Sun v-else-if="themeStore.mode === 'light'" :size="16" />
        <Monitor v-else :size="16" />
      </button>

      <!-- Language toggle -->
      <button
        class="tb-icon-btn"
        aria-label="切换语言"
        title="切换语言"
        @click="toggleLang"
      >
        <Languages :size="16" />
      </button>

      <!-- Window controls -->
      <div class="win-controls">
        <button class="wc-btn wc-min" aria-label="最小化" @click="minimize">
          <Minus :size="12" />
        </button>
        <button class="wc-btn wc-max" aria-label="最大化" @click="toggleMaximize">
          <Square :size="12" />
        </button>
        <button class="wc-btn wc-close" aria-label="关闭" @click="close">
          <X :size="12" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Moon, Sun, Monitor, Languages, Minus, Square, X } from 'lucide-vue-next'
import { getCurrentWindow } from '@tauri-apps/api/window'
import { useThemeStore } from '@/stores/theme'
import type { ThemeMode } from '@/stores/theme'

const { t, locale } = useI18n()
const themeStore = useThemeStore()
const win = getCurrentWindow()

const modeLabel = computed(() => ({
  dark: '深色', light: '浅色', system: '跟随系统',
}[themeStore.mode]))

const modeOrder: ThemeMode[] = ['dark', 'light', 'system']
function cycleMode() {
  const idx = modeOrder.indexOf(themeStore.mode)
  themeStore.setMode(modeOrder[(idx + 1) % 3])
}

function toggleLang() {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
}

const minimize       = () => win.minimize()
const toggleMaximize = () => win.toggleMaximize()
const close          = () => win.close()
</script>

<style scoped>
.titlebar {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-deep);
  border-bottom: 1px solid var(--border);
  user-select: none;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-left: 16px;
  height: 100%;
}

.logo {
  font-size: 20px;
  color: var(--accent);
  line-height: 1;
}

.app-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-1);
}

.app-sub {
  font-size: 11px;
  color: var(--text-3);
  margin-left: 2px;
  letter-spacing: 0.08em;
}

.titlebar-right {
  display: flex;
  align-items: center;
  height: 100%;
}

.tb-icon-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-2);
  cursor: pointer;
  border-radius: 4px;
  margin: 0 2px;
  transition: background 80ms ease-out, color 80ms ease-out;
}

.tb-icon-btn:hover {
  background: var(--bg-hover);
  color: var(--text-1);
}

.win-controls {
  display: flex;
  align-items: center;
  height: 100%;
  margin-left: 4px;
}

.wc-btn {
  width: 40px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-2);
  cursor: default;
  transition: background 80ms ease-out, color 80ms ease-out;
}

.wc-btn:hover {
  background: var(--bg-hover);
  color: var(--text-1);
}

.wc-close {
  width: 48px;
}

.wc-close:hover {
  background: var(--wc-close-bg);
  color: #fff;
}
</style>
