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
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'connection' }" @click="currentSettingsTab = 'connection'">
            <span class="nav-icon">🔗</span>
            <span>{{ t('settings.connection') }}</span>
          </div>
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'about' }" @click="currentSettingsTab = 'about'">
            <span class="nav-icon">ℹ️</span>
            <span>{{ t('settings.about') }}</span>
          </div>
        </div>

        <div class="settings-main">
          <div class="settings-panel">
            <template v-if="currentSettingsTab === 'appearance'">
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
                        @click="localThemeMode = 'light'; handleThemeChange()"
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
                        @click="localThemeMode = 'dark'; handleThemeChange()"
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
                        @click="localThemeMode = 'system'; handleThemeChange()"
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
                      v-model="currentAccentColor"
                      @apply="handleAccentColorChange"
                    />
                  </div>
                </div>
              </div>
            </template>

            <template v-else-if="currentSettingsTab === 'connection'">
              <div class="connection-container">
                <!-- 代理设置卡片 -->
                <div class="settings-card proxy-card">
                  <div class="card-header-section">
                    <div class="card-icon-wrapper">
                      <span class="card-icon">🌐</span>
                    </div>
                    <div class="card-title-section">
                      <h4 class="card-title">{{ t('settings.proxySettings') }}</h4>
                      <p class="card-description">{{ t('settings.proxyDescription') }}</p>
                    </div>
                  </div>
                  <div class="card-content">
                    <div class="placeholder-content">
                      <span class="placeholder-icon">🌐</span>
                      <p>{{ t('settings.proxyUnderDevelopment') }}</p>
                    </div>
                  </div>
                </div>

                <!-- 网络设置卡片 -->
                <div class="settings-card network-card">
                  <div class="card-header-section">
                    <div class="card-icon-wrapper">
                      <span class="card-icon">📡</span>
                    </div>
                    <div class="card-title-section">
                      <h4 class="card-title">{{ t('settings.networkSettings') }}</h4>
                      <p class="card-description">{{ t('settings.networkDescription') }}</p>
                    </div>
                  </div>
                  <div class="card-content">
                    <div class="placeholder-content">
                      <span class="placeholder-icon">📡</span>
                      <p>{{ t('settings.networkUnderDevelopment') }}</p>
                    </div>
                  </div>
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
import AccentColorPicker from '@/components/shared/AccentColorPicker.vue'

const { t, locale } = useI18n()

const currentSettingsTab = ref('appearance')
const localThemeMode = ref(localStorage.getItem('themeMode') || 'light')
const currentLanguage = ref(locale.value)
const currentAccentColor = ref(localStorage.getItem('accentColor') || '#3b82f6')

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

const handleAccentColorChange = (color: string) => {
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
  background: var(--content-bg);
  gap: 1px;
}

.settings-sidebar {
  width: 260px;
  flex-shrink: 0;
  background: var(--sidebar-bg);
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  border-right: 1px solid var(--border-color);
}

.dark .settings-sidebar {
  border-right-color: var(--border-strong);
}

.settings-nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: all var(--transition-normal);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  position: relative;
  overflow: hidden;
}

.settings-nav-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: linear-gradient(180deg, var(--active-color), var(--active-color-suppl));
  transform: scaleY(0);
  transition: transform var(--transition-fast);
  opacity: 0;
}

.settings-nav-item:hover {
  background: var(--hover-color);
}

.settings-nav-item:hover::before {
  transform: scaleY(0.5);
  opacity: 0.7;
}

.settings-nav-item.active {
  background: linear-gradient(90deg, var(--active-color-suppl), transparent);
  color: var(--active-color);
  font-weight: 600;
  box-shadow: inset 0 0 0 1px rgba(var(--active-color-rgb), 0.1);
}

.settings-nav-item.active::before {
  transform: scaleY(1);
  opacity: 1;
}

.dark .settings-nav-item.active {
  background: linear-gradient(90deg, rgba(var(--active-color-rgb), 0.15), transparent);
  box-shadow: inset 0 0 0 1px rgba(var(--active-color-rgb), 0.2);
}

.nav-icon {
  font-size: 18px;
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.settings-main {
  flex: 1;
  overflow: auto;
  background: var(--content-bg);
  padding: 32px;
}

.settings-panel {
  max-width: 900px;
  margin: 0 auto;
}

/* ===== 外观设置容器 ===== */
.appearance-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ===== 连接设置容器 ===== */
.connection-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ===== 设置卡片基础样式 ===== */
.settings-card {
  background: var(--card-bg);
  border-radius: var(--border-radius-lg);
  padding: 0;
  margin-bottom: 0;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
  transition: all var(--transition-normal);
  overflow: hidden;
}

.dark .settings-card {
  border-color: var(--border-strong);
}

.settings-card:hover {
  box-shadow: var(--shadow-md);
  border-color: var(--active-color-suppl);
  transform: translateY(-2px);
}

/* ===== 卡片头部区域 ===== */
.card-header-section {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 24px 28px;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, var(--card-bg), var(--content-bg));
}

.dark .card-header-section {
  border-bottom-color: var(--border-strong);
}

.card-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: var(--border-radius-md);
  background: linear-gradient(135deg, var(--active-color-suppl), var(--card-bg));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: var(--shadow-sm);
}

.card-icon {
  font-size: 24px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.card-title-section {
  flex: 1;
  min-width: 0;
}

.card-title {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
  letter-spacing: 0.3px;
  line-height: 1.4;
}

.card-description {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

/* ===== 卡片内容区域 ===== */
.card-content {
  padding: 24px 28px;
}

.theme-options,
.language-options,
.accent-color-options {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

/* ===== 主题选项按钮 ===== */
.theme-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px 20px;
  border: 2px solid var(--border-color);
  border-radius: var(--border-radius-md);
  background: var(--content-bg);
  cursor: pointer;
  transition: all var(--transition-fast);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  min-width: 160px;
  flex: 1;
  position: relative;
  overflow: hidden;
}

.dark .theme-option {
  border-color: var(--border-strong);
}

.theme-option::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, transparent, rgba(255,255,255,0.08), transparent);
  transform: translateX(-100%);
  transition: transform 0.6s ease;
}

.theme-option:hover::before {
  transform: translateX(100%);
}

.theme-option:hover {
  border-color: var(--active-color);
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.dark .theme-option:hover {
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
}

.theme-option.active {
  border-color: var(--active-color);
  background: linear-gradient(135deg, var(--active-color-suppl), var(--card-bg));
  color: var(--active-color);
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.dark .theme-option.active {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.option-content {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.option-text {
  font-weight: 500;
}

.theme-icon {
  font-size: 20px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.option-check {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--active-color);
  flex-shrink: 0;
  animation: checkmark 0.3s ease-out;
}

@keyframes checkmark {
  from {
    transform: scale(0);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

/* ===== 语言选项按钮 ===== */
.language-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px 20px;
  border: 2px solid var(--border-color);
  border-radius: var(--border-radius-md);
  background: var(--content-bg);
  cursor: pointer;
  transition: all var(--transition-fast);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  flex: 1;
  min-width: 200px;
  position: relative;
  overflow: hidden;
}

.dark .language-option {
  border-color: var(--border-strong);
}

.language-option::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, transparent, rgba(255,255,255,0.08), transparent);
  transform: translateX(-100%);
  transition: transform 0.6s ease;
}

.language-option:hover::before {
  transform: translateX(100%);
}

.language-option:hover {
  border-color: var(--active-color);
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
}

.dark .language-option:hover {
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
}

.language-option.active {
  border-color: var(--active-color);
  background: linear-gradient(135deg, var(--active-color-suppl), var(--card-bg));
  color: var(--active-color);
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.dark .language-option.active {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.language-flag {
  font-size: 24px;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.option-text-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.option-label {
  font-size: 15px;
  font-weight: 600;
  line-height: 1.2;
}

.option-sublabel {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 400;
  opacity: 0.8;
}

/* ===== 强调色选择器 ===== */
.accent-color-options {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.accent-color-option {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: 3px solid var(--border-color);
  cursor: pointer;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  flex-shrink: 0;
}

.dark .accent-color-option {
  border-color: var(--border-strong);
}

.accent-color-option::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.3), transparent);
  border-radius: 50%;
}

.accent-color-option::after {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at 30% 30%, rgba(255,255,255,0.2), transparent);
  border-radius: 50%;
}

.accent-color-option:hover {
  transform: scale(1.15);
  box-shadow: var(--shadow-md);
  border-color: var(--active-color);
}

.accent-color-option.active {
  border-color: white;
  box-shadow: 0 0 0 4px var(--active-color-suppl), var(--shadow-md);
  animation: pulse-ring 2s infinite;
}

.accent-check {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 1;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
  animation: checkmark-bounce 0.4s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

@keyframes pulse-ring {
  0%, 100% {
    box-shadow: 0 0 0 4px var(--active-color-suppl), var(--shadow-md);
  }
  50% {
    box-shadow: 0 0 0 8px var(--active-color-suppl), var(--shadow-md);
  }
}

@keyframes checkmark-bounce {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.about-page-container {
  max-width: 800px;
  margin: 0 auto;
}

.about-header {
  text-align: center;
  padding: 48px 32px;
  background: linear-gradient(135deg, var(--active-color-suppl), var(--card-bg));
  border-radius: var(--border-radius-lg);
  margin-bottom: 28px;
  border: 1px solid var(--border-color);
  position: relative;
  overflow: hidden;
}

.dark .about-header {
  border-color: var(--border-strong);
}

.about-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, var(--active-color-suppl) 0%, transparent 70%);
  opacity: 0.1;
  animation: rotate 20s linear infinite;
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.app-logo {
  margin-bottom: 20px;
  position: relative;
  z-index: 1;
}

.app-name {
  font-size: 36px;
  font-weight: 800;
  background: linear-gradient(135deg, var(--active-color), var(--active-color-suppl));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 12px 0;
  position: relative;
  z-index: 1;
  letter-spacing: 1px;
}

.app-version {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 12px 0;
  font-weight: 500;
  position: relative;
  z-index: 1;
}

.app-description {
  font-size: 15px;
  color: var(--text-color);
  opacity: 0.85;
  margin: 0;
  line-height: 1.6;
  position: relative;
  z-index: 1;
}

.about-card {
  background: var(--card-bg);
  border-radius: var(--border-radius-lg);
  margin-bottom: 24px;
  overflow: hidden;
  border: 1px solid var(--border-color);
  transition: all var(--transition-normal);
}

.dark .about-card {
  border-color: var(--border-strong);
}

.about-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.card-header {
  padding: 18px 24px;
  background: linear-gradient(135deg, var(--active-color-suppl), var(--card-bg));
  border-bottom: 1px solid var(--border-color);
}

.dark .card-header {
  border-bottom-color: var(--border-strong);
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color);
  letter-spacing: 0.3px;
}

.card-body {
  padding: 24px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 20px;
}

.author-details h4 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: var(--text-color);
  font-weight: 600;
}

.author-details p {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.github-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 14px 24px;
  background: linear-gradient(135deg, var(--card-bg), var(--content-bg));
  border: 2px solid var(--border-color);
  border-radius: var(--border-radius-md);
  text-decoration: none;
  color: var(--text-color);
  font-weight: 600;
  font-size: 14px;
  transition: all var(--transition-fast);
  box-shadow: var(--shadow-sm);
}

.dark .github-link-btn {
  border-color: var(--border-strong);
}

.github-link-btn:hover {
  border-color: var(--active-color);
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
  background: linear-gradient(135deg, var(--active-color-suppl), var(--card-bg));
}

.github-link-btn:active {
  transform: translateY(-1px);
}

.github-link-btn svg {
  flex-shrink: 0;
  transition: transform var(--transition-fast);
}

.github-link-btn:hover svg {
  transform: scale(1.1);
}

@media (max-width: 768px) {
  .settings-layout {
    flex-direction: column;
  }
  
  .settings-sidebar {
    width: 100%;
    flex-direction: row;
    overflow-x: auto;
    padding: 16px;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
  }
  
  .dark .settings-sidebar {
    border-bottom-color: var(--border-strong);
  }
  
  .settings-nav-item {
    flex-shrink: 0;
    padding: 12px 14px;
  }
  
  .settings-main {
    padding: 20px 16px;
  }
  
  .settings-card {
    border-radius: var(--border-radius-md);
  }
  
  .card-header-section {
    padding: 20px;
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
  
  .card-content {
    padding: 20px;
  }
  
  .theme-option {
    min-width: 140px;
    padding: 14px 16px;
  }
  
  .language-option {
    min-width: 180px;
    padding: 14px 16px;
  }
  
  .accent-color-option {
    width: 48px;
    height: 48px;
  }
  
  .about-header {
    padding: 32px 20px;
  }
  
  .app-name {
    font-size: 28px;
  }
}

@media (max-width: 480px) {
  .settings-sidebar {
    gap: 8px;
  }
  
  .settings-nav-item {
    font-size: 13px;
    gap: 8px;
  }
  
  .nav-icon {
    font-size: 16px;
  }
  
  .card-header-section {
    flex-direction: column;
    text-align: center;
    align-items: center;
  }
  
  .card-icon-wrapper {
    margin-bottom: 8px;
  }
  
  .card-title-section {
    text-align: center;
  }
  
  .card-title {
    font-size: 17px;
  }
  
  .card-description {
    font-size: 13px;
  }
  
  .card-content {
    padding: 16px;
  }
  
  .theme-options,
  .language-options {
    flex-direction: column;
  }
  
  .theme-option,
  .language-option {
    width: 100%;
    min-width: unset;
    justify-content: center;
  }
  
  .accent-color-options {
    justify-content: center;
    gap: 12px;
  }
  
  .accent-color-option {
    width: 44px;
    height: 44px;
  }
  
  .appearance-container {
    gap: 16px;
  }
}
</style>
