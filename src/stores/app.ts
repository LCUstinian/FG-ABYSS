import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  // State
  const theme = ref<'light' | 'dark'>('dark')
  const language = ref<'zh-CN' | 'en-US'>('zh-CN')
  const accentColor = ref('#409eff')
  const connectionStatus = ref<'connected' | 'disconnected'>('disconnected')
  const projectCount = ref(0)
  const payloadCount = ref(0)

  // Getters
  const isDarkTheme = computed(() => theme.value === 'dark')
  const isConnected = computed(() => connectionStatus.value === 'connected')

  // Actions
  function setTheme(newTheme: 'light' | 'dark') {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
  }

  function setLanguage(lang: 'zh-CN' | 'en-US') {
    language.value = lang
    localStorage.setItem('language', lang)
  }

  function setAccentColor(color: string) {
    accentColor.value = color
    localStorage.setItem('accentColor', color)
  }

  function setConnectionStatus(status: 'connected' | 'disconnected') {
    connectionStatus.value = status
  }

  function updateProjectCount(count: number) {
    projectCount.value = count
  }

  function updatePayloadCount(count: number) {
    payloadCount.value = count
  }

  // 初始化时从 localStorage 加载
  function loadFromStorage() {
    const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
    const savedLanguage = localStorage.getItem('language') as 'zh-CN' | 'en-US' | null
    const savedAccentColor = localStorage.getItem('accentColor')

    if (savedTheme) theme.value = savedTheme
    if (savedLanguage) language.value = savedLanguage
    if (savedAccentColor) accentColor.value = savedAccentColor
  }

  return {
    // State
    theme,
    language,
    accentColor,
    connectionStatus,
    projectCount,
    payloadCount,
    // Getters
    isDarkTheme,
    isConnected,
    // Actions
    setTheme,
    setLanguage,
    setAccentColor,
    setConnectionStatus,
    updateProjectCount,
    updatePayloadCount,
    loadFromStorage,
  }
})
