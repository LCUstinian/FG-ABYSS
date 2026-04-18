<template>
  <div class="appearance-panel">
    <div class="settings-card">
      <h3 class="card-title">主题</h3>
      <div class="theme-cards">
        <button
          v-for="m in themeModes"
          :key="m.value"
          class="theme-card"
          :class="{ 'is-active': themeStore.mode === m.value }"
          @click="themeStore.setMode(m.value)"
        >
          <div class="theme-preview" :data-preview="m.value" />
          <span class="theme-label">{{ m.label }}</span>
        </button>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">强调色</h3>
      <div class="accent-circles">
        <button
          v-for="(colors, key) in ACCENT_COLORS"
          :key="key"
          class="accent-circle"
          :style="{ background: colors[themeStore.resolvedMode] }"
          :aria-label="`强调色: ${key}`"
          :class="{ 'is-selected': themeStore.accentKey === key }"
          @click="themeStore.setAccent(key as AccentKey)"
        >
          <Check v-if="themeStore.accentKey === key" :size="14" color="#fff" />
        </button>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">字体大小</h3>
      <n-select
        :value="themeStore.fontSize"
        :options="fontSizeOptions"
        size="small"
        style="width: 100px"
        @update:value="themeStore.setFontSize"
      />
    </div>

    <div class="settings-card">
      <h3 class="card-title">语言</h3>
      <n-select
        :value="locale"
        :options="langOptions"
        size="small"
        style="width: 140px"
        @update:value="(v: string) => (locale = v)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Check } from 'lucide-vue-next'
import { useThemeStore, ACCENT_COLORS } from '@/stores/theme'
import type { AccentKey } from '@/stores/theme'
import type { ThemeMode, FontSize } from '@/stores/theme'

const { locale } = useI18n()
const themeStore = useThemeStore()

const themeModes: Array<{ value: ThemeMode; label: string }> = [
  { value: 'system', label: '跟随系统' },
  { value: 'dark',   label: '深色'     },
  { value: 'light',  label: '浅色'     },
]

const fontSizeOptions: Array<{ label: string; value: FontSize }> = [
  { label: '12px', value: '12px' },
  { label: '13px', value: '13px' },
  { label: '14px', value: '14px' },
  { label: '15px', value: '15px' },
]

const langOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English',  value: 'en-US' },
]
</script>

<style scoped>
.appearance-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

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

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
}

.theme-cards {
  display: flex;
  gap: 12px;
}

.theme-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border: 2px solid var(--border);
  border-radius: 6px;
  background: transparent;
  cursor: pointer;
  transition: border-color 100ms;
}

.theme-card:hover { border-color: var(--text-2); }
.theme-card.is-active { border-color: var(--accent); }

.theme-preview {
  width: 64px;
  height: 40px;
  border-radius: 4px;
}

.theme-preview[data-preview="dark"]   { background: #0d0e13; }
.theme-preview[data-preview="light"]  { background: #f6f7fb; border: 1px solid #dde0ec; }
.theme-preview[data-preview="system"] {
  background: linear-gradient(135deg, #0d0e13 50%, #f6f7fb 50%);
}

.theme-label {
  font-size: 12px;
  color: var(--text-2);
}

.accent-circles {
  display: flex;
  gap: 8px;
}

.accent-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: box-shadow 100ms, transform 100ms;
}

.accent-circle.is-selected {
  box-shadow: 0 0 0 2px var(--bg-base), 0 0 0 4px var(--accent);
  transform: scale(1.05);
}
</style>
