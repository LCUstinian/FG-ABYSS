<template>
  <div class="content-section">
    <PageHeader :title="t('settings.title')" :subtitle="t('settings.subtitle')" />
    <div class="content-body">
      <div class="settings-layout">
        <div class="settings-sidebar">
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'appearance' }" @click="currentSettingsTab = 'appearance'">
            <span class="nav-icon">🎨</span>
            <span>{{ t('settings.appearance') }}</span>
          </div>
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'proxy' }" @click="currentSettingsTab = 'proxy'">
            <span class="nav-icon">🌐</span>
            <span>{{ t('settings.proxy') }}</span>
          </div>
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'network' }" @click="currentSettingsTab = 'network'">
            <span class="nav-icon">📡</span>
            <span>{{ t('settings.network') }}</span>
          </div>
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'about' }" @click="currentSettingsTab = 'about'">
            <span class="nav-icon">ℹ️</span>
            <span>{{ t('settings.about') }}</span>
          </div>
        </div>

        <div class="settings-main">
          <div class="settings-panel">
            <template v-if="currentSettingsTab === 'appearance'">
              <div class="settings-card">
                <h4>{{ t('settings.theme') }}</h4>
                <div class="theme-options">
                  <button 
                    class="theme-option" 
                    :class="{ active: localThemeMode === 'light' }"
                    @click="localThemeMode = 'light'; handleThemeChange()"
                  >
                    <span class="theme-icon">☀️</span>
                    <span>{{ t('settings.lightMode') }}</span>
                    <span v-if="localThemeMode === 'light'" class="theme-check">✓</span>
                  </button>
                  <button 
                    class="theme-option" 
                    :class="{ active: localThemeMode === 'dark' }"
                    @click="localThemeMode = 'dark'; handleThemeChange()"
                  >
                    <span class="theme-icon">🌙</span>
                    <span>{{ t('settings.darkMode') }}</span>
                    <span v-if="localThemeMode === 'dark'" class="theme-check">✓</span>
                  </button>
                  <button 
                    class="theme-option" 
                    :class="{ active: localThemeMode === 'system' }"
                    @click="localThemeMode = 'system'; handleThemeChange()"
                  >
                    <span class="theme-icon">🖥️</span>
                    <span>{{ t('settings.systemMode') }}</span>
                    <span v-if="localThemeMode === 'system'" class="theme-check">✓</span>
                  </button>
                </div>
              </div>
              <div class="settings-card">
                <h4>{{ t('settings.language') }}</h4>
                <div class="language-options">
                  <button 
                    class="language-option" 
                    :class="{ active: currentLanguage === 'zh-CN' }"
                    @click="changeLanguage('zh-CN')"
                  >
                    <span class="language-icon">🇨🇳</span>
                    <span>{{ t('settings.chinese') }}</span>
                    <span v-if="currentLanguage === 'zh-CN'" class="language-check">✓</span>
                  </button>
                  <button 
                    class="language-option" 
                    :class="{ active: currentLanguage === 'en-US' }"
                    @click="changeLanguage('en-US')"
                  >
                    <span class="language-icon">🇺🇸</span>
                    <span>{{ t('settings.english') }}</span>
                    <span v-if="currentLanguage === 'en-US'" class="language-check">✓</span>
                  </button>
                </div>
              </div>
              <div class="settings-card">
                <h4>{{ t('settings.accentColor') }}</h4>
                <div class="accent-color-options">
                  <button 
                    v-for="color in accentColors" 
                    :key="color.value"
                    class="accent-color-option"
                    :class="{ active: currentAccentColor === color.value }"
                    :style="{ backgroundColor: color.value }"
                    @click="changeAccentColor(color.value)"
                  >
                    <span v-if="currentAccentColor === color.value" class="accent-color-check">✓</span>
                  </button>
                </div>
              </div>
            </template>

            <template v-else-if="currentSettingsTab === 'proxy'">
              <div class="settings-card">
                <h4>{{ t('settings.proxySettings') }}</h4>
                <div class="placeholder-content">
                  <span class="placeholder-icon">🌐</span>
                  <p>{{ t('settings.proxyUnderDevelopment') }}</p>
                </div>
              </div>
            </template>

            <template v-else-if="currentSettingsTab === 'network'">
              <div class="settings-card">
                <h4>{{ t('settings.networkSettings') }}</h4>
                <div class="placeholder-content">
                  <span class="placeholder-icon">📡</span>
                  <p>{{ t('settings.networkUnderDevelopment') }}</p>
                </div>
              </div>
            </template>

            <template v-else-if="currentSettingsTab === 'about'">
              <div class="about-page-container">
                <div class="about-header">
                  <div class="app-logo">
                    <span class="logo-icon" style="font-size: 48px;">ℹ️</span>
                  </div>
                  <h2 class="app-name">FG-ABYSS</h2>
                  <p class="app-version">{{ t('settings.currentVersion') }} v1.0.0</p>
                  <p class="app-description">{{ t('settings.appDescription') }}</p>
                </div>

                <div class="about-content">
                  <div class="about-card">
                    <div class="card-header">
                      <h3>{{ t('settings.author') }}</h3>
                    </div>
                    <div class="card-body">
                      <div class="author-info">
                        <span style="font-size: 48px;">👨‍💻</span>
                        <div class="author-details">
                          <h4>{{ t('settings.authorName') }}</h4>
                          <p>{{ t('settings.authorIntro') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="about-card">
                    <div class="card-header">
                      <h3>{{ t('settings.githubRepo') }}</h3>
                    </div>
                    <div class="card-body">
                      <a 
                        href="https://github.com/FG-ABYSS" 
                        target="_blank" 
                        class="github-link-btn"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                        </svg>
                        <span>{{ t('settings.visitGithub') }}</span>
                      </a>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/shared/PageHeader.vue'

const { t, locale } = useI18n()

const currentSettingsTab = ref('appearance')
const localThemeMode = ref(localStorage.getItem('themeMode') || 'light')
const currentLanguage = ref(locale.value)
const currentAccentColor = ref(localStorage.getItem('accentColor') || '#3b82f6')

const accentColors = [
  { value: '#3b82f6' },
  { value: '#8b5cf6' },
  { value: '#ec4899' },
  { value: '#f59e0b' },
  { value: '#10b981' },
  { value: '#06b6d4' },
]

const handleThemeChange = () => {
  localStorage.setItem('themeMode', localThemeMode.value)
  
  let isDark = false
  if (localThemeMode.value === 'dark') {
    isDark = true
  } else if (localThemeMode.value === 'system') {
    isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  
  localStorage.setItem('theme', isDark ? 'dark' : 'light')
  document.documentElement.classList.toggle('dark', isDark)
  
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'themeMode',
    newValue: localThemeMode.value
  }))
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'theme',
    newValue: isDark ? 'dark' : 'light'
  }))
}

const changeLanguage = (lang: string) => {
  currentLanguage.value = lang
  locale.value = lang
  localStorage.setItem('locale', lang)
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'locale',
    newValue: lang
  }))
}

const changeAccentColor = (color: string) => {
  currentAccentColor.value = color
  localStorage.setItem('accentColor', color)
  document.documentElement.style.setProperty('--active-color', color)
  
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'accentColor',
    newValue: color
  }))
}

watch(() => locale.value, (newVal) => {
  currentLanguage.value = newVal
})
</script>

<style scoped>
.content-section {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
}

.content-body {
  flex: 1;
  overflow: hidden;
  background: var(--content-bg);
}

.settings-layout {
  display: flex;
  height: 100%;
  background: var(--border-color);
}

.settings-sidebar {
  width: 240px;
  background: var(--sidebar-bg);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.settings-nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: all var(--transition-normal);
}

.settings-nav-item:hover {
  background: var(--hover-color);
}

.settings-nav-item.active {
  background: var(--active-color);
  color: white;
}

.nav-icon {
  font-size: 20px;
}

.settings-main {
  flex: 1;
  overflow: auto;
  background: var(--card-bg);
  padding: 24px;
}

.settings-panel {
  max-width: 900px;
  margin: 0 auto;
}

.settings-card {
  background: var(--content-bg);
  border-radius: var(--border-radius-lg);
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: var(--shadow-sm);
}

.settings-card h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
}

.theme-options,
.language-options,
.accent-color-options {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.theme-option,
.language-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border: 2px solid var(--border-color);
  border-radius: var(--border-radius-md);
  background: var(--card-bg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.theme-option:hover,
.language-option:hover {
  border-color: var(--active-color);
  transform: translateY(-2px);
}

.theme-option.active,
.language-option.active {
  border-color: var(--active-color);
  background: var(--active-color-suppl);
}

.theme-icon,
.language-icon {
  font-size: 20px;
}

.accent-color-option {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 3px solid var(--border-color);
  cursor: pointer;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: white;
}

.accent-color-option:hover {
  transform: scale(1.1);
}

.accent-color-option.active {
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px var(--active-color-suppl);
}

.placeholder-content {
  text-align: center;
  padding: 40px 20px;
}

.placeholder-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 16px;
}

.about-page-container {
  max-width: 800px;
  margin: 0 auto;
}

.about-header {
  text-align: center;
  padding: 40px 20px;
  background: var(--active-color-suppl);
  border-radius: var(--border-radius-lg);
  margin-bottom: 24px;
}

.app-logo {
  margin-bottom: 16px;
}

.app-name {
  font-size: 32px;
  font-weight: 700;
  color: var(--active-color);
  margin: 0 0 8px 0;
}

.app-version {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
}

.app-description {
  font-size: 14px;
  color: var(--text-color);
  opacity: 0.8;
  margin: 0;
}

.about-card {
  background: var(--content-bg);
  border-radius: var(--border-radius-lg);
  margin-bottom: 20px;
  overflow: hidden;
}

.card-header {
  padding: 16px 20px;
  background: var(--active-color-suppl);
  border-bottom: 1px solid var(--border-color);
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
}

.card-body {
  padding: 20px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.author-details h4 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: var(--text-color);
}

.author-details p {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
}

.github-link-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: var(--card-bg);
  border: 2px solid var(--border-color);
  border-radius: var(--border-radius-md);
  text-decoration: none;
  color: var(--text-color);
  font-weight: 600;
  transition: all var(--transition-fast);
}

.github-link-btn:hover {
  border-color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.github-link-btn svg {
  flex-shrink: 0;
}
</style>
