import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'

// 创建 i18n 实例
const i18n = createI18n({
  // 使用 composition API 模式
  legacy: false,
  // 全局注入
  globalInjection: true,
  // 默认语言
  locale: localStorage.getItem('locale') || 'zh-CN',
  // 回退语言
  fallbackLocale: 'zh-CN',
  // 翻译消息
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  },
  // 静默警告
  silentTranslationWarn: true,
  silentFallbackWarn: true,
  // 确保使用全局作用域
  allowComposition: true
})

export default i18n
export const i18nGlobal = i18n.global