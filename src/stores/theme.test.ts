import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore } from './theme'

describe('useThemeStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('default mode is dark', () => {
    const store = useThemeStore()
    expect(store.mode).toBe('dark')
    expect(store.resolvedMode).toBe('dark')
  })

  it('setMode changes mode', () => {
    const store = useThemeStore()
    store.setMode('light')
    expect(store.mode).toBe('light')
    expect(store.resolvedMode).toBe('light')
  })

  it('setAccent changes accentKey', () => {
    const store = useThemeStore()
    store.setAccent('cyan')
    expect(store.accentKey).toBe('cyan')
  })

  it('naiveTheme is null for light mode', () => {
    const store = useThemeStore()
    store.setMode('light')
    expect(store.naiveTheme).toBeNull()
  })
})
