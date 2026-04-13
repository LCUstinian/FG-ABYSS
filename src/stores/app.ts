import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const theme = ref<'light' | 'dark'>('dark')
  const language = ref<'zh' | 'en'>('zh')

  function setTheme(newTheme: 'light' | 'dark') {
    theme.value = newTheme
  }

  function setLanguage(newLanguage: 'zh' | 'en') {
    language.value = newLanguage
  }

  return {
    theme,
    language,
    setTheme,
    setLanguage,
  }
})
