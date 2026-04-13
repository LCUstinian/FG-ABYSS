import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    welcome: 'Welcome to FG-ABYSS',
    terminal: 'Terminal',
    settings: 'Settings',
  },
  zh: {
    welcome: '欢迎来到 FG-ABYSS',
    terminal: '终端',
    settings: '设置',
  },
}

const i18n = createI18n({
  legacy: false,
  locale: 'zh',
  fallbackLocale: 'en',
  messages,
})

export default i18n
