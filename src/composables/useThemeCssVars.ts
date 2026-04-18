import { watchEffect } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { ACCENT_COLORS, DARK_VARS, LIGHT_VARS } from '@/theme'

export function useThemeCssVars() {
  const store = useThemeStore()

  watchEffect(() => {
    const isDark = store.resolvedMode === 'dark'
    const accent = ACCENT_COLORS[store.accentKey][store.resolvedMode]
    const vars   = isDark ? DARK_VARS : LIGHT_VARS
    const root   = document.documentElement

    root.setAttribute('data-theme', store.resolvedMode)
    Object.entries(vars).forEach(([k, v]) => root.style.setProperty(k, v))
    root.style.setProperty('--accent', accent)
    // accent-bg: 12% opacity dark (0x1f ≈ 12%), 9% opacity light (0x17 ≈ 9%)
    root.style.setProperty('--accent-bg', isDark ? `${accent}1f` : `${accent}17`)
    root.style.setProperty('--font-mono', "'JetBrains Mono', 'Fira Code', monospace")
  })
}
