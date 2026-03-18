<template>
  <div class="appearance-view">
    <div class="appearance-container">
      <!-- 主题设置卡片 -->
      <div class="settings-card theme-card">
        <div class="card-header-section">
          <div class="card-icon-wrapper">
            <span class="card-icon">🎨</span>
          </div>
          <div class="card-title-section">
            <h4 class="card-title">{{ t('settings.theme') }}</h4>
            <p class="card-description">{{ t('settings.themeDescription') }}</p>
          </div>
        </div>
        <div class="card-content">
          <div class="theme-options">
            <button 
              class="theme-option" 
              :class="{ active: localThemeMode === 'light' }"
              @click="handleThemeChange('light')"
            >
              <div class="option-content">
                <span class="theme-icon">☀️</span>
                <span class="option-text">{{ t('settings.lightMode') }}</span>
              </div>
              <span v-if="localThemeMode === 'light'" class="option-check">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
              </span>
            </button>
            <button 
              class="theme-option" 
              :class="{ active: localThemeMode === 'dark' }"
              @click="handleThemeChange('dark')"
            >
              <div class="option-content">
                <span class="theme-icon">🌙</span>
                <span class="option-text">{{ t('settings.darkMode') }}</span>
              </div>
              <span v-if="localThemeMode === 'dark'" class="option-check">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
              </span>
            </button>
            <button 
              class="theme-option" 
              :class="{ active: localThemeMode === 'system' }"
              @click="handleThemeChange('system')"
            >
              <div class="option-content">
                <span class="theme-icon">🖥️</span>
                <span class="option-text">{{ t('settings.systemMode') }}</span>
              </div>
              <span v-if="localThemeMode === 'system'" class="option-check">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- 语言设置卡片 -->
      <div class="settings-card language-card">
        <div class="card-header-section">
          <div class="card-icon-wrapper">
            <span class="card-icon">🌐</span>
          </div>
          <div class="card-title-section">
            <h4 class="card-title">{{ t('settings.language') }}</h4>
            <p class="card-description">{{ t('settings.languageDescription') }}</p>
          </div>
        </div>
        <div class="card-content">
          <div class="language-options">
            <button 
              class="language-option" 
              :class="{ active: currentLanguage === 'zh-CN' }"
              @click="changeLanguage('zh-CN')"
            >
              <div class="option-content">
                <span class="language-flag">🇨🇳</span>
                <div class="option-text-group">
                  <span class="option-label">中文</span>
                  <span class="option-sublabel">Chinese</span>
                </div>
              </div>
              <span v-if="currentLanguage === 'zh-CN'" class="option-check">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
              </span>
            </button>
            <button 
              class="language-option" 
              :class="{ active: currentLanguage === 'en-US' }"
              @click="changeLanguage('en-US')"
            >
              <div class="option-content">
                <span class="language-flag">🇺🇸</span>
                <div class="option-text-group">
                  <span class="option-label">English</span>
                  <span class="option-sublabel">English</span>
                </div>
              </div>
              <span v-if="currentLanguage === 'en-US'" class="option-check">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- 强调色设置卡片 -->
      <div class="settings-card accent-card">
        <div class="card-header-section">
          <div class="card-icon-wrapper">
            <span class="card-icon">🎨</span>
          </div>
          <div class="card-title-section">
            <h4 class="card-title">{{ t('settings.accentColor') }}</h4>
            <p class="card-description">{{ t('settings.accentColorDescription') }}</p>
          </div>
        </div>
        <div class="card-content">
          <AccentColorPicker 
            v-model="accentColor"
            @apply="handleApplyAccentColor"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AccentColorPicker from '@/components/shared/AccentColorPicker.vue'

const props = defineProps({
  isDarkTheme: {
    type: Boolean,
    default: false
  },
  themeMode: {
    type: String,
    default: 'system'
  },
  currentLanguage: {
    type: String,
    default: 'zh-CN'
  }
})

const emit = defineEmits(['update:theme-mode', 'change-language'])

const { t } = useI18n()

const localThemeMode = ref(props.themeMode)
const currentLanguage = ref(props.currentLanguage)
const accentColor = ref(localStorage.getItem('accentColor') || '#3b82f6')

// 监听外部主题变化
watch(() => props.themeMode, (newVal) => {
  localThemeMode.value = newVal
})

// 监听语言变化
watch(() => props.currentLanguage, (newVal) => {
  currentLanguage.value = newVal
})

const handleThemeChange = (mode: string) => {
  localThemeMode.value = mode
  emit('update:theme-mode', mode)
}

const changeLanguage = (lang: string) => {
  currentLanguage.value = lang
  emit('change-language', lang)
}

// 应用强调色
const handleApplyAccentColor = (color: string) => {
  localStorage.setItem('accentColor', color)
  document.documentElement.style.setProperty('--active-color', color)
  
  // 发送事件通知全局主题更新
  window.dispatchEvent(new CustomEvent('accent-color-change', {
    detail: { color }
  }))
}

// 初始化强调色
onMounted(() => {
  const savedColor = localStorage.getItem('accentColor')
  if (savedColor) {
    accentColor.value = savedColor
    document.documentElement.style.setProperty('--active-color', savedColor)
  }
})
</script>

<style scoped>
.appearance-view {
  width: 100%;
  height: auto;
  display: flex;
  flex-direction: column;
  animation: slideIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.appearance-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 32px 40px;
  box-sizing: border-box;
}

.settings-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 28px 32px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  width: 100%;
  max-width: 1000px;
  margin: 0 auto;
  position: relative;
  overflow: hidden;
}

.settings-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--active-color) 0%, transparent 100%);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.3s ease;
}

.settings-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.settings-card:hover::before {
  transform: scaleX(1);
}

.card-header-section {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 24px;
  position: relative;
  z-index: 1;
}

.card-icon-wrapper {
  width: 48px;
  height: 48px;
  background: var(--active-color-bg);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.settings-card:hover .card-icon-wrapper {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(var(--active-color-rgb), 0.2);
}

.card-icon {
  font-size: 24px;
  transition: transform 0.3s ease;
}

.settings-card:hover .card-icon {
  transform: rotate(5deg);
}

.card-title-section {
  flex: 1;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  transition: color 0.3s ease;
}

.card-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.5;
  transition: color 0.3s ease;
}

.card-content {
  margin-top: 20px;
  position: relative;
  z-index: 1;
}

.theme-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 18px;
}

.theme-option {
  padding: 22px 18px;
  background: var(--card-bg-hover);
  border: 2px solid var(--border-color);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}

.theme-option::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, var(--active-color-bg) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.theme-option:hover {
  border-color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.theme-option:hover::before {
  opacity: 0.5;
}

.theme-option.active {
  border-color: var(--active-color);
  background: var(--active-color-bg);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  }
  50% {
    box-shadow: 0 8px 24px rgba(var(--active-color-rgb), 0.2);
  }
}

.theme-option.active::before {
  opacity: 1;
}

.option-content {
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
  z-index: 1;
}

.theme-icon {
  font-size: 24px;
  transition: transform 0.3s ease;
}

.theme-option:hover .theme-icon {
  transform: scale(1.1);
}

.option-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  transition: color 0.3s ease;
}

.option-check {
  color: var(--active-color);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.theme-option.active .option-check {
  transform: scale(1.1);
  animation: checkBounce 0.5s ease;
}

@keyframes checkBounce {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
  }
}

.language-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}

.language-option {
  padding: 22px 18px;
  background: var(--card-bg-hover);
  border: 2px solid var(--border-color);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}

.language-option::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, var(--active-color-bg) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.language-option:hover {
  border-color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.language-option:hover::before {
  opacity: 0.5;
}

.language-option.active {
  border-color: var(--active-color);
  background: var(--active-color-bg);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  animation: pulse 2s infinite;
}

.language-option.active::before {
  opacity: 1;
}

.language-flag {
  font-size: 32px;
  transition: transform 0.3s ease;
}

.language-option:hover .language-flag {
  transform: rotate(5deg);
}

.option-text-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
  position: relative;
  z-index: 1;
}

.option-label {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  transition: color 0.3s ease;
}

.option-sublabel {
  font-size: 12px;
  color: var(--text-secondary);
  transition: color 0.3s ease;
}

.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px 20px;
  color: var(--text-secondary);
  animation: fadeIn 0.6s ease;
}

.placeholder-icon {
  font-size: 48px;
  opacity: 0.5;
  transition: opacity 0.3s ease;
}

.placeholder-content:hover .placeholder-icon {
  opacity: 0.8;
  transform: scale(1.1);
  transition: all 0.3s ease;
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .appearance-container {
    padding: 0 40px 60px;
  }
  
  .theme-options {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  }
  
  .language-options {
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  }
  
  .settings-card {
    padding: 32px 36px;
  }
  
  .card-title {
    font-size: 20px;
  }
  
  .card-description {
    font-size: 15px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .appearance-container {
    padding: 0 24px 32px;
    gap: 16px;
  }
  
  .settings-card {
    padding: 24px 28px;
  }
  
  .theme-options {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .language-options {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .card-header-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .card-icon-wrapper {
    width: 40px;
    height: 40px;
  }
  
  .card-icon {
    font-size: 20px;
  }
  
  .card-title {
    font-size: 16px;
  }
  
  .card-description {
    font-size: 13px;
  }
}
</style>
