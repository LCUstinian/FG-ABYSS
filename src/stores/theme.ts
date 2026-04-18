import { defineStore } from 'pinia'
import { computed } from 'vue'
import { darkTheme, useOsTheme } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'
import { buildOverrides, ACCENT_COLORS } from '@/theme'
import type { AccentKey } from '@/theme'

export type ThemeMode = 'dark' | 'light' | 'system'
export type FontSize = '12px' | '13px' | '14px' | '15px'

export { ACCENT_COLORS }
export type { AccentKey }

export const useThemeStore = defineStore('theme', () => {
  const osTheme = useOsTheme()
  const mode      = useLocalStorage<ThemeMode>('fg-theme-mode',   'dark')
  const accentKey = useLocalStorage<AccentKey>('fg-theme-accent', 'blue')
  const fontSize  = useLocalStorage<FontSize>( 'fg-theme-fontsize', '13px')

  const resolvedMode = computed<'dark' | 'light'>(() =>
    mode.value === 'system' ? (osTheme.value ?? 'dark') : mode.value
  )

  const naiveTheme = computed(() =>
    resolvedMode.value === 'dark' ? darkTheme : null
  )

  const themeOverrides = computed(() =>
    buildOverrides(
      ACCENT_COLORS[accentKey.value][resolvedMode.value],
      resolvedMode.value === 'dark',
      fontSize.value,
    )
  )

  function setMode(m: ThemeMode)    { mode.value = m }
  function setAccent(k: AccentKey)  { accentKey.value = k }
  function setFontSize(s: FontSize) { fontSize.value = s }

  return {
    mode, accentKey, fontSize,
    resolvedMode, naiveTheme, themeOverrides,
    setMode, setAccent, setFontSize,
  }
})
