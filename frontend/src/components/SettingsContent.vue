<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('settings.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('settings.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="settings-layout">
        <!-- 设置侧边栏 -->
        <div class="settings-sidebar">
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'appearance' }" @click="currentSettingsTab = 'appearance'">
            <Palette :size="18" />
            <span>{{ t('settings.appearance') }}</span>
          </div>
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'proxy' }" @click="currentSettingsTab = 'proxy'">
            <Globe :size="18" />
            <span>{{ t('settings.proxy') }}</span>
          </div>
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'network' }" @click="currentSettingsTab = 'network'">
            <Wifi :size="18" />
            <span>{{ t('settings.network') }}</span>
          </div>
          <div class="settings-nav-item" :class="{ active: currentSettingsTab === 'about' }" @click="currentSettingsTab = 'about'">
            <Info :size="18" />
            <span>{{ t('settings.about') }}</span>
          </div>
        </div>
        
        <!-- 设置内容 -->
        <div class="settings-main">
          <div class="settings-panel">
            <!-- 外观设置 -->
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
                
                <!-- 自定义颜色选择器 - 紧凑单行版本 -->
                <div class="custom-color-picker-compact">
                  <div class="color-picker-row">
                    <!-- 颜色预览 -->
                    <div class="color-preview-compact">
                      <div 
                        class="color-dot" 
                        :style="{ backgroundColor: customColorValue }"
                      ></div>
                      <span class="color-hex">{{ customColorValue.toUpperCase() }}</span>
                    </div>
                    
                    <!-- 颜色选择器和输入 -->
                    <div class="color-controls-compact">
                      <input
                        type="color"
                        v-model="customColorValue"
                        class="color-picker-btn"
                        @input="handleCustomColorPick"
                      />
                      <input
                        type="text"
                        v-model="customColorValue"
                        class="color-text-input"
                        placeholder="#3182CE"
                        @input="validateHexColor"
                        @blur="applyCustomColor"
                      />
                    </div>
                    
                    <!-- 操作按钮 -->
                    <div class="color-actions-compact">
                      <NButton 
                        type="primary" 
                        size="small"
                        @click="applyCustomColor"
                        class="btn-apply"
                      >
                        <template #icon>
                          <span class="btn-icon">✓</span>
                        </template>
                        {{ t('settings.applyColor') }}
                      </NButton>
                      <NButton 
                        size="small"
                        @click="resetToDefaultColor"
                        class="btn-reset"
                      >
                        <template #icon>
                          <span class="btn-icon">↺</span>
                        </template>
                        {{ t('settings.resetToDefault') }}
                      </NButton>
                    </div>
                  </div>
                </div>
              </div>
              <div class="settings-card">
                <h4>{{ t('settings.font') }}</h4>
                <div class="font-settings">
                  <div class="font-setting-item">
                    <label>{{ t('settings.fontFamily') }}</label>
                    <div class="font-family-selector">
                      <select 
                        v-model="currentFontFamily" 
                        @change="changeFontFamily(currentFontFamily)"
                        class="font-family-select"
                      >
                        <option 
                          v-for="font in fontFamilies" 
                          :key="font.value"
                          :value="font.value"
                          :style="{ fontFamily: font.value }"
                        >
                          {{ font.name }}
                        </option>
                      </select>
                    </div>
                  </div>
                  <div class="font-setting-item">
                    <label>{{ t('settings.fontSize') }}</label>
                    <div class="font-size-selector">
                      <div class="font-size-input-with-button">
                        <div class="font-size-input-group">
                          <input 
                            type="number" 
                            v-model.number="fontSizeValue" 
                            min="8" 
                            max="24" 
                            class="font-size-input"
                            :placeholder="t('settings.fontSizePlaceholder')"
                          >
                          <span class="font-size-unit">px</span>
                        </div>
                        <button @click="applyFontSize" class="apply-button">{{ t('settings.apply') }}</button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>
            
            <!-- 代理设置 -->
            <template v-else-if="currentSettingsTab === 'proxy'">
              <div class="settings-card">
                <h4>{{ t('settings.proxySettings') }}</h4>
                <div class="placeholder-content">
                  <div class="placeholder-icon">🌐</div>
                  <p>{{ t('settings.proxyUnderDevelopment') }}</p>
                  <p>{{ t('settings.proxyOptions') }}</p>
                  <div class="placeholder-tips">
                    <span>{{ t('settings.tip') }}</span>
                    <p>{{ t('settings.proxyTip') }}</p>
                  </div>
                </div>
              </div>
            </template>
            
            <!-- 网络设置 -->
            <template v-else-if="currentSettingsTab === 'network'">
              <div class="settings-card">
                <h4>{{ t('settings.networkSettings') }}</h4>
                <div class="placeholder-content">
                  <div class="placeholder-icon">📡</div>
                  <p>{{ t('settings.networkUnderDevelopment') }}</p>
                  <p>{{ t('settings.networkOptions') }}</p>
                  <div class="placeholder-tips">
                    <span>{{ t('settings.tip') }}</span>
                    <p>{{ t('settings.networkTip') }}</p>
                  </div>
                </div>
              </div>
            </template>
            
            <!-- 关于 -->
            <template v-else-if="currentSettingsTab === 'about'">
              <div class="about-page-container">
                <!-- 应用标志和版本信息 -->
                <div class="about-header">
                  <div class="app-logo">
                    <div class="logo-icon">
                      <Info :size="48" />
                    </div>
                  </div>
                  <h2 class="app-name">FG-ABYSS</h2>
                  <p class="app-version">
                    {{ t('settings.currentVersion') }} v1.0.0
                  </p>
                  <p class="app-description">{{ t('settings.appDescription') }}</p>
                </div>

                <!-- 主要内容区域 -->
                <div class="about-content">
                  <!-- 作者信息卡片 -->
                  <div class="about-card author-card">
                    <div class="card-header">
                      <h3>{{ t('settings.author') }}</h3>
                    </div>
                    <div class="card-body">
                      <div class="author-info">
                        <div class="author-avatar">
                          <span>👨‍💻</span>
                        </div>
                        <div class="author-details">
                          <h4>{{ t('settings.authorName') }}</h4>
                          <p class="author-intro">{{ t('settings.authorIntro') }}</p>
                          <p class="author-field">{{ t('settings.authorField') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- GitHub 链接卡片 -->
                  <div class="about-card github-card">
                    <div class="card-header">
                      <h3>{{ t('settings.githubRepo') }}</h3>
                    </div>
                    <div class="card-body">
                      <div class="github-links">
                        <a 
                          href="https://github.com/FG-ABYSS" 
                          target="_blank" 
                          class="github-link-btn"
                          @click.prevent="openGithubRepo"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor" class="github-icon">
                            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                          </svg>
                          <span>{{ t('settings.visitGithub') }}</span>
                        </a>
                      </div>
                    </div>
                  </div>

                  <!-- 项目信息卡片 -->
                  <div class="about-card info-card">
                    <div class="card-header">
                      <h3>{{ t('settings.projectVersion') }}</h3>
                    </div>
                    <div class="card-body">
                      <div class="info-grid">
                        <div class="info-item">
                          <span class="info-label">{{ t('settings.currentVersion') }}</span>
                          <span class="info-value">v1.0.0</span>
                        </div>
                        <div class="info-item">
                          <span class="info-label">{{ t('settings.buildDate') }}</span>
                          <span class="info-value">2024-01-15</span>
                        </div>
                        <div class="info-item">
                          <span class="info-label">{{ t('settings.updateDate') }}</span>
                          <span class="info-value">2024-01-15</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- 技术栈卡片 -->
                  <div class="about-card tech-card">
                    <div class="card-header">
                      <h3>{{ t('settings.technology') }}</h3>
                    </div>
                    <div class="card-body">
                      <div class="tech-tags">
                        <span class="tech-tag">Wails v3</span>
                        <span class="tech-tag">Vue 3</span>
                        <span class="tech-tag">TypeScript</span>
                        <span class="tech-tag">Naive UI</span>
                        <span class="tech-tag">Go</span>
                      </div>
                    </div>
                  </div>

                  <!-- 许可证卡片 -->
                  <div class="about-card license-card">
                    <div class="card-header">
                      <h3>{{ t('settings.license') }}</h3>
                    </div>
                    <div class="card-body">
                      <div class="license-info">
                        <div class="license-icon">📄</div>
                        <div class="license-details">
                          <h4>{{ t('settings.licenseType') }}</h4>
                          <p>{{ t('settings.copyright') }} © 2024 FG-ABYSS Team. {{ t('settings.allRightsReserved') }}.</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- 特别感谢卡片 -->
                  <div class="about-card thanks-card">
                    <div class="card-header">
                      <h3>{{ t('settings.thanks') }}</h3>
                    </div>
                    <div class="card-body">
                      <p class="thanks-text">{{ t('settings.thanksSupport') }}</p>
                    </div>
                  </div>
                </div>

                <!-- 页脚 -->
                <div class="about-footer">
                  <div class="footer-links">
                    <a href="#" class="footer-link">{{ t('settings.documentation') }}</a>
                    <a href="#" class="footer-link">{{ t('settings.changelog') }}</a>
                    <a href="#" class="footer-link">{{ t('settings.privacyPolicy') }}</a>
                    <a href="#" class="footer-link">{{ t('settings.termsOfService') }}</a>
                  </div>
                  <p class="copyright">
                    {{ t('settings.copyright') }} © 2024 FG-ABYSS Team
                  </p>
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
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Palette, Globe, Wifi, Info } from 'lucide-vue-next'

const { t, locale } = useI18n()

const props = defineProps({
  isDarkTheme: {
    type: Boolean,
    default: false
  },
  themeMode: {
    type: String,
    default: 'system'
  }
})

const emit = defineEmits(['update:theme-mode'])

const localThemeMode = ref(props.themeMode)
const currentLanguage = ref(locale.value)
const currentSettingsTab = ref('appearance')

// 强调色选项 - 极客风格
const accentColors = [
  { value: '#3182ce' }, // 科技蓝
  { value: '#48bb78' }, // 代码绿
  { value: '#805ad5' }, // 极客紫
  { value: '#f56565' }, // 错误红
  { value: '#4299e1' }, // 信息蓝
  { value: '#38a169' }, // 成功绿
  { value: '#dd6b20' }, // 警告橙
  { value: '#d53f8c' }, // 粉色
  { value: '#319795' }, // 青色
  { value: '#718096' }, // 灰色
  { value: '#2c5282' }, // 深蓝
  { value: '#0f766e' }, // 深绿
  { value: '#c05621' }, // 深橙
  { value: '#553c9a' }, // 深紫
  { value: '#b83280' }  // 深粉
]

// 当前强调色
const currentAccentColor = ref(localStorage.getItem('accentColor') || '#3182ce')

// 自定义颜色选择器
const customColorValue = ref(currentAccentColor.value)
const defaultAccentColor = '#3182ce'
const colorError = ref(false)

// 字体选项 - 极客风格
const fontFamilies = [
  { name: t('settings.fontDefault'), value: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, sans-serif' },
  { name: t('settings.fontMonospace'), value: 'SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace' },
  { name: 'Courier New', value: '"Courier New", monospace' },
  { name: 'Monaco', value: 'Monaco, monospace' },
  { name: 'Fira Code', value: '"Fira Code", monospace' },
  { name: 'Source Code Pro', value: '"Source Code Pro", monospace' },
  { name: 'Roboto Mono', value: '"Roboto Mono", monospace' },
  { name: 'Ubuntu Mono', value: '"Ubuntu Mono", monospace' }
]

// 字体大小选项
const fontSizes = [
  { label: t('settings.fontSizeTiny'), value: '10px' },
  { label: t('settings.fontSizeSmall'), value: '12px' },
  { label: t('settings.fontSizeMedium'), value: '14px' },
  { label: t('settings.fontSizeLarge'), value: '16px' },
  { label: t('settings.fontSizeHuge'), value: '18px' }
]

// 当前字体和字体大小
const currentFontFamily = ref(localStorage.getItem('fontFamily') || '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, sans-serif')
const currentFontSize = ref(localStorage.getItem('fontSize') || '14px')
const fontSizeValue = ref(parseInt(currentFontSize.value))

// 监听 props 变化
watch(() => props.themeMode, (newMode) => {
  localThemeMode.value = newMode
})

// 处理主题切换
const handleThemeChange = () => {
  emit('update:theme-mode', localThemeMode.value)
}

// 处理语言切换
const changeLanguage = (lang: string) => {
  locale.value = lang
  currentLanguage.value = lang
  localStorage.setItem('locale', lang)
  // 触发storage事件，让其他组件知道语言变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'locale',
    newValue: lang
  }))
}

// 处理强调色切换
const changeAccentColor = (color: string) => {
  currentAccentColor.value = color
  customColorValue.value = color
  localStorage.setItem('accentColor', color)
  document.documentElement.style.setProperty('--active-color', color)
  // 触发 storage 事件，让其他组件知道颜色变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'accentColor',
    newValue: color
  }))
}

// 验证 HEX 颜色格式
const validateHexColor = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value.trim()
  
  // 简单的 HEX 颜色验证
  const hexColorRegex = /^#?([0-9A-Fa-f]{6}|[0-9A-Fa-f]{3})$/
  
  if (hexColorRegex.test(value)) {
    colorError.value = false
    // 格式化颜色值（添加#号）
    if (!value.startsWith('#')) {
      customColorValue.value = '#' + value
    }
  } else {
    colorError.value = true
  }
}

// 处理自定义颜色选择器输入
const handleCustomColorPick = () => {
  colorError.value = false
  // 实时更新预览，但不应用
}

// 应用自定义颜色
const applyCustomColor = () => {
  if (colorError.value) {
    // 如果颜色值无效，恢复为当前强调色
    customColorValue.value = currentAccentColor.value
    colorError.value = false
    return
  }
  
  let color = customColorValue.value.trim()
  
  // 确保颜色值以#开头
  if (!color.startsWith('#')) {
    color = '#' + color
    customColorValue.value = color
  }
  
  // 验证颜色格式
  const hexColorRegex = /^#?([0-9A-Fa-f]{6}|[0-9A-Fa-f]{3})$/
  if (!hexColorRegex.test(color)) {
    // 颜色格式无效，恢复为当前强调色
    customColorValue.value = currentAccentColor.value
    return
  }
  
  // 应用颜色
  changeAccentColor(color)
}

// 重置为默认颜色
const resetToDefaultColor = () => {
  customColorValue.value = defaultAccentColor
  changeAccentColor(defaultAccentColor)
}

// 处理字体切换
const changeFontFamily = (font: string) => {
  currentFontFamily.value = font
  localStorage.setItem('fontFamily', font)
  document.documentElement.style.setProperty('--font-family', font)
  // 触发storage事件，让其他组件知道字体变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'fontFamily',
    newValue: font
  }))
}

// 处理字体大小切换
const changeFontSize = (size: string) => {
  currentFontSize.value = size
  localStorage.setItem('fontSize', size)
  document.documentElement.style.setProperty('--font-size', size)
  // 触发storage事件，让其他组件知道字体大小变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'fontSize',
    newValue: size
  }))
}

// 应用字体大小
const applyFontSize = () => {
  // 确保字体大小在有效范围内
  if (fontSizeValue.value < 8) fontSizeValue.value = 8
  if (fontSizeValue.value > 24) fontSizeValue.value = 24
  
  const size = `${fontSizeValue.value}px`
  changeFontSize(size)
}

// 打开 GitHub 仓库
const openGithubRepo = () => {
  window.open('https://github.com/FG-ABYSS', '_blank')
}

// 初始化语言、强调色、字体和字体大小
onMounted(() => {
  const savedLanguage = localStorage.getItem('locale')
  if (savedLanguage) {
    locale.value = savedLanguage
    currentLanguage.value = savedLanguage
  }
  
  // 初始化强调色
  const savedAccentColor = localStorage.getItem('accentColor')
  if (savedAccentColor) {
    currentAccentColor.value = savedAccentColor
    document.documentElement.style.setProperty('--active-color', savedAccentColor)
  } else {
    document.documentElement.style.setProperty('--active-color', currentAccentColor.value)
  }
  
  // 初始化字体
  const savedFontFamily = localStorage.getItem('fontFamily')
  if (savedFontFamily) {
    currentFontFamily.value = savedFontFamily
    document.documentElement.style.setProperty('--font-family', savedFontFamily)
  } else {
    document.documentElement.style.setProperty('--font-family', currentFontFamily.value)
  }
  
  // 初始化字体大小
  const savedFontSize = localStorage.getItem('fontSize')
  if (savedFontSize) {
    currentFontSize.value = savedFontSize
    fontSizeValue.value = parseInt(savedFontSize)
    document.documentElement.style.setProperty('--font-size', savedFontSize)
  } else {
    fontSizeValue.value = parseInt(currentFontSize.value)
    document.documentElement.style.setProperty('--font-size', currentFontSize.value)
  }
  
  // 监听 localStorage 中变化
  window.addEventListener('storage', (event) => {
    if (event.key === 'accentColor' && event.newValue) {
      currentAccentColor.value = event.newValue
      document.documentElement.style.setProperty('--active-color', event.newValue)
    } else if (event.key === 'fontFamily' && event.newValue) {
      currentFontFamily.value = event.newValue
      document.documentElement.style.setProperty('--font-family', event.newValue)
    } else if (event.key === 'fontSize' && event.newValue) {
      currentFontSize.value = event.newValue
      fontSizeValue.value = parseInt(event.newValue)
      document.documentElement.style.setProperty('--font-size', event.newValue)
    } else if (event.key === 'locale' && event.newValue) {
      locale.value = event.newValue
      currentLanguage.value = event.newValue
    }
  })
})
</script>

<style scoped>
.content-section {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

.content-header {
  width: 100%;
  padding: 24px 24px 20px 24px;
  margin-bottom: 0;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  box-sizing: border-box;
}

.dark .content-header {
  border-bottom-color: var(--border-strong);
}

.content-section h1 {
  margin: 0;
  font-size: 0;
  line-height: 1;
  color: var(--text-color);
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
}

.content-section h1 .title {
  font-size: 24px;
  font-weight: 600;
  color: var(--active-color);
  letter-spacing: 0;
}

.content-section h1 .separator {
  color: var(--text-tertiary);
  font-weight: 300;
  font-size: 20px;
}

.content-section h1 .subtitle {
  font-size: 16px;
  font-weight: 400;
  color: var(--text-secondary);
  letter-spacing: 0;
}

.content-body {
  flex: 1;
  width: 100%;
  min-height: 0;  /* 关键修复：允许 flex 子项缩小 */
  padding: 0;
  margin: 0;
  box-sizing: border-box;
  background: var(--content-bg);
  border-top: none;
  display: flex;
  align-items: stretch;
  overflow-y: auto;  /* 添加垂直滚动 */
}

/* 设置内容样式 - 深色主题风格 */
.settings-layout {
  display: flex;
  width: 100%;
  height: 100%;
  gap: 1px;
  background: var(--border-color);
  margin: 0;
  padding: 0;
  border-top: none;
  border-right: 1px solid var(--border-color);
  overflow: hidden;  /* 防止内容溢出 */
  flex-direction: row;
}

.settings-sidebar {
  width: 200px;
  height: 100%;
  background: var(--sidebar-bg);
  padding: 20px 0;
  overflow-y: auto;
  margin: 0;
  border: none;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  z-index: 10;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
}

.dark .settings-sidebar {
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.4);
}

.settings-nav-item {
  padding: 12px 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  border-left: 3px solid transparent;
  display: flex;
  align-items: center;
  gap: 12px;
}

.settings-nav-item:hover {
  background: var(--hover-color);
  border-left-color: var(--active-color);
}

.settings-nav-item.active {
  background: var(--active-color);
  color: white;
  border-left-color: var(--active-color);
}

.settings-main {
  flex: 1;
  width: 100%;
  min-height: 0;  /* 关键修复：允许 flex 子项缩小 */
  background: var(--content-bg);
  padding: 0;
  margin: 0;
  overflow: hidden;  /* 防止内容溢出 */
  border: none;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.05);
  z-index: 5;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.dark .settings-main {
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.4);
}

.settings-panel {
  width: 100%;
  min-height: 0;  /* 关键修复：允许 flex 子项缩小 */
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px;
  border: none;
  overflow-y: auto;  /* 添加垂直滚动 */
  box-sizing: border-box;
  flex: 1;  /* 占据剩余空间 */
}

.settings-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;  /* 确保占满宽度 */
  box-sizing: border-box;  /* 包含 padding 在宽度内 */
}

.dark .settings-card {
  border-color: var(--border-strong);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.settings-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.dark .settings-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.settings-card h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  opacity: 0.8;
  transition: color 0.3s ease;
}

.dark .settings-card h4 {
  opacity: 0.7;
}



/* 主题选项 */
.theme-options {
  display: flex;
  gap: 16px;
  width: 100%;
  margin: 0;
  padding: 0;
}

.theme-option {
  flex: 1;
  min-width: 0;
  padding: 24px 20px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: var(--text-color);
  text-align: center;
  position: relative;
  min-height: 120px;
}

.dark .theme-option {
  border-color: var(--border-strong);
}

.theme-option:hover {
  border-color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.dark .theme-option:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

.theme-option.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.3);
  transform: translateY(-2px);
}

.dark .theme-option.active {
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.4);
  filter: brightness(1.1);
}

.theme-icon {
  font-size: 28px;
  display: block;
  margin: 0 auto;
  transition: transform 0.3s ease;
}

.theme-option:hover .theme-icon {
  transform: scale(1.1);
}

.theme-option.active .theme-icon {
  transform: scale(1.15);
}

.theme-option span:nth-child(2) {
  font-size: 14px;
  font-weight: 500;
  display: block;
  margin: 0 auto;
  transition: color 0.3s ease;
}

.theme-check {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 16px;
  height: 16px;
  background: white;
  color: var(--active-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.dark .theme-check {
  background: white;
  color: var(--active-color);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

.theme-option.active:hover .theme-check {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.dark .theme-option.active:hover .theme-check {
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.4);
}

/* 语言选项 */
.language-options {
  display: flex;
  gap: 16px;
  width: 100%;
  margin: 0;
  padding: 0;
}

.language-option {
  flex: 1;
  min-width: 0;
  padding: 24px 20px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: var(--text-color);
  text-align: center;
  position: relative;
  min-height: 120px;
}

.dark .language-option {
  border-color: var(--border-strong);
}

.language-option:hover {
  border-color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.dark .language-option:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

.language-option.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.3);
  transform: translateY(-2px);
}

.dark .language-option.active {
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.4);
  filter: brightness(1.1);
}

.language-icon {
  font-size: 28px;
  display: block;
  margin: 0 auto;
  transition: transform 0.3s ease;
}

.language-option:hover .language-icon {
  transform: scale(1.1);
}

.language-option.active .language-icon {
  transform: scale(1.15);
}

.language-option span:nth-child(2) {
  font-size: 14px;
  font-weight: 500;
  display: block;
  margin: 0 auto;
  transition: color 0.3s ease;
}

.dark .language-option span:nth-child(2) {
  opacity: 0.7;
}

.language-check {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 16px;
  height: 16px;
  background: white;
  color: var(--active-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.dark .language-check {
  background: white;
  color: var(--active-color);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

.language-option.active:hover .language-check {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
}

.dark .language-option.active:hover .language-check {
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.4);
}

/* 强调色选项 */
.accent-color-options {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  width: 100%;
  margin: 0;
  padding: 0;
  justify-content: center;
}

.accent-color-option {
  width: 40px;
  height: 40px;
  border: 2px solid transparent;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.accent-color-option:hover {
  transform: scale(1.1);
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
}

.dark .accent-color-option:hover {
  box-shadow: 0 0 12px rgba(0, 0, 0, 0.4);
}

.accent-color-option.active {
  border-color: var(--text-color);
  box-shadow: 0 0 12px rgba(0, 0, 0, 0.3);
  transform: scale(1.15);
}

.dark .accent-color-option.active {
  box-shadow: 0 0 16px rgba(0, 0, 0, 0.5);
}

.accent-color-check {
  color: white;
  font-size: 16px;
  font-weight: bold;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

/* 自定义颜色选择器 - 紧凑单行版本 */
.custom-color-picker-compact {
  margin-top: 16px;
  padding: 12px 16px;
  background: var(--card-bg);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .custom-color-picker-compact {
  border-color: var(--border-strong);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.color-picker-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: nowrap;
  transition: opacity 0.3s ease;
}

/* 颜色预览 - 紧凑版 */
.color-preview-compact {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
}

.color-dot {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
  transition: background-color 0.4s cubic-bezier(0.4, 0, 0.2, 1), transform 0.2s ease, border-color 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .color-dot {
  border-color: var(--border-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.color-dot:hover {
  transform: scale(1.1);
}

.color-hex {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  min-width: 70px;
  transition: color 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .color-hex {
  color: var(--text-secondary);
  font-weight: 500;
}

/* 颜色控制 - 紧凑版 */
.color-controls-compact {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.color-picker-btn {
  width: 36px;
  height: 32px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  padding: 0;
  background: none;
  transition: transform 0.2s ease, box-shadow 0.3s ease;
  flex-shrink: 0;
}

.color-picker-btn::-webkit-color-swatch-wrapper {
  padding: 0;
  border-radius: 6px;
}

.color-picker-btn::-webkit-color-swatch {
  border: 2px solid var(--border-color);
  border-radius: 6px;
  transition: border-color 0.3s ease;
}

.color-picker-btn::-moz-color-swatch {
  border: 2px solid var(--border-color);
  border-radius: 6px;
  transition: border-color 0.3s ease;
}

.color-picker-btn:hover {
  transform: scale(1.08);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.dark .color-picker-btn:hover {
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.dark .color-picker-btn::-webkit-color-swatch {
  border-color: var(--border-strong);
}

.dark .color-picker-btn::-moz-color-swatch {
  border-color: var(--border-strong);
}

.color-text-input {
  flex: 1;
  min-width: 100px;
  padding: 6px 10px;
  background: var(--bg-primary);
  border: 2px solid var(--border-color);
  border-radius: 6px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  color: var(--text-primary);
  transition: all 0.3s ease;
  outline: none;
}

.color-text-input:focus {
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dark .color-text-input {
  background: var(--bg-tertiary);
  border-color: var(--border-strong);
}

.dark .color-text-input:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.color-text-input.error {
  border-color: var(--error-color);
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.dark .color-text-input.error {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2);
}

.color-text-input::placeholder {
  color: var(--text-tertiary);
  opacity: 0.6;
}

.dark .color-text-input::placeholder {
  opacity: 0.5;
}

/* 操作按钮 - 紧凑版 */
.color-actions-compact {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.btn-apply,
.btn-reset {
  padding: 6px 14px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 6px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 6px;
  border: none;
  position: relative;
  overflow: hidden;
}

/* 应用按钮 - 主操作 */
.btn-apply {
  background: var(--active-color);
  color: white;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

.btn-apply::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.btn-apply:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
  filter: brightness(1.1);
}

.btn-apply:hover::before {
  left: 100%;
}

.btn-apply:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
  filter: brightness(0.95);
}

/* 深色模式下的应用按钮 - 增强对比度 */
.dark .btn-apply {
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.4);
  filter: brightness(1.05);
}

.dark .btn-apply:hover {
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.5);
  filter: brightness(1.15);
}

.dark .btn-apply:active {
  filter: brightness(1);
}

/* 重置按钮 - 次操作 */
.btn-reset {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.btn-reset::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(0, 0, 0, 0.05), transparent);
  transition: left 0.5s ease;
}

.btn-reset:hover {
  transform: translateY(-2px);
  border-color: var(--active-color);
  color: var(--active-color);
  background: var(--bg-hover);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.btn-reset:hover::before {
  left: 100%;
}

.btn-reset:active {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

/* 深色模式下的重置按钮 - 增强对比度 */
.dark .btn-reset {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border-color: var(--border-strong);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

.dark .btn-reset:hover {
  border-color: var(--active-color);
  color: var(--active-color);
  background: var(--bg-hover);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
  filter: brightness(1.1);
}

.dark .btn-reset:active {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

/* 按钮图标 */
.btn-icon {
  font-size: 14px;
  font-weight: bold;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
}

.btn-apply .btn-icon {
  color: white;
}

.btn-reset .btn-icon {
  color: currentColor;
  transition: transform 0.3s ease;
}

.btn-reset:hover .btn-icon {
  transform: rotate(-180deg);
}

/* 响应式适配 */
@media (max-width: 768px) {
  .color-picker-row {
    flex-wrap: wrap;
    gap: 12px;
  }
  
  .color-preview-compact {
    width: 100%;
    justify-content: center;
    min-width: auto;
  }
  
  .color-controls-compact {
    width: 100%;
  }
  
  .color-actions-compact {
    width: 100%;
  }
  
  .btn-apply,
  .btn-reset {
    flex: 1;
  }
}

/* 字体设置 */
.font-settings {
  display: flex;
  gap: 20px;
  margin-top: 12px;
}

.font-setting-item {
  flex: 1;
  min-width: 0;
}

.font-setting-item label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  margin-bottom: 8px;
  opacity: 0.8;
  transition: color 0.3s ease, opacity 0.3s ease;
}

.dark .font-setting-item label {
  opacity: 0.7;
}

.font-family-selector {
  width: 100%;
}

.font-family-select {
  width: 100%;
  padding: 10px 16px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-color);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  outline: none;
}

.dark .font-family-select {
  border-color: var(--border-strong);
  background: var(--bg-tertiary);
}

.font-family-select:hover {
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dark .font-family-select:hover {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.font-family-select:focus {
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dark .font-family-select:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.font-family-select option {
  padding: 8px;
  background: var(--card-bg);
  color: var(--text-color);
}

.dark .font-family-select option {
  background: var(--bg-tertiary);
}

.font-size-selector {
  width: 100%;
}

.font-size-input-with-button {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.font-size-input-group {
  flex: 1;
  display: flex;
  align-items: stretch;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .font-size-input-group {
  border-color: var(--border-strong);
  background: var(--bg-tertiary);
}

.font-size-input-group:hover {
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dark .font-size-input-group:hover {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.font-size-input-group:focus-within {
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.dark .font-size-input-group:focus-within {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.apply-button {
  padding: 10px 20px;
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

.dark .apply-button {
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.4);
  filter: brightness(1.05);
}

.apply-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
  filter: brightness(1.1);
}

.dark .apply-button:hover {
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.5);
  filter: brightness(1.15);
}

.apply-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
  filter: brightness(0.95);
}

.dark .apply-button:active {
  filter: brightness(1);
}

.apply-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.font-size-input {
  flex: 1;
  padding: 10px 16px;
  border: none;
  background: transparent;
  color: var(--text-color);
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
  height: 100%;
  transition: color 0.3s ease;
}

.font-size-input::placeholder {
  color: var(--text-color);
  opacity: 0.5;
}

.dark .font-size-input::placeholder {
  opacity: 0.4;
}

.font-size-input:focus {
  box-shadow: none;
}

.font-size-unit {
  padding: 0 16px;
  background: rgba(0, 0, 0, 0.05);
  color: var(--text-color);
  font-size: 14px;
  border-left: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.dark .font-size-unit {
  background: rgba(255, 255, 255, 0.05);
  border-left-color: var(--border-strong);
}

.font-size-input-group:hover .font-size-unit {
  background: rgba(0, 0, 0, 0.08);
}

.dark .font-size-input-group:hover .font-size-unit {
  background: rgba(255, 255, 255, 0.08);
}

.font-size-hint {
  font-size: 12px;
  color: var(--text-color);
  opacity: 0.7;
  margin: 0;
  padding: 0;
  text-align: center;
  margin-top: 8px;
  transition: color 0.3s ease, opacity 0.3s ease;
}

.dark .font-size-hint {
  opacity: 0.6;
}

/* 空白页样式 */
.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  gap: 16px;
  width: 100%;  /* 确保占满宽度 */
  box-sizing: border-box;  /* 包含 padding 在宽度内 */
  min-height: 300px;  /* 最小高度，避免内容过少时显得空洞 */
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 16px;
  flex-shrink: 0;  /* 防止图标被压缩 */
}

.placeholder-content p {
  margin: 0;
  color: var(--text-color);
  opacity: 0.8;
  line-height: 1.6;
  max-width: 600px;  /* 限制文本最大宽度，提高可读性 */
  width: 100%;  /* 确保在小屏幕上也能正常显示 */
}

.placeholder-tips {
  margin-top: 20px;
  padding: 16px;
  background: var(--hover-color);
  border-radius: 8px;
  width: 100%;  /* 占满父容器宽度 */
  max-width: 600px;  /* 限制最大宽度，提高可读性 */
  text-align: left;
  box-sizing: border-box;  /* 包含 padding 在宽度内 */
}

.placeholder-tips span {
  font-weight: 600;
  color: var(--text-color);
  display: block;
  margin-bottom: 8px;
}

.placeholder-tips p {
  margin: 0;
  opacity: 0.9;
  font-size: 13px;
}

.placeholder-tips p {
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.5;
}

.font-size-hint {
  font-size: 12px;
  color: var(--text-color);
  opacity: 0.7;
  margin: 0;
  padding: 0;
  text-align: center;
  margin-top: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-header {
    padding: 12px 16px;
  }
  
  .content-body {
    padding: 16px;
  }
  
  .settings-layout {
    flex-direction: column;
  }
  
  .settings-sidebar {
    width: 100%;
    height: auto;
    max-height: 200px;
  }
  
  .theme-options {
    flex-direction: column;
  }
  
  .theme-option {
    flex-direction: row;
    min-height: 80px;
  }
  
  .language-options {
    flex-direction: column;
  }
  
  .language-option {
    flex-direction: row;
    min-height: 80px;
  }
  
  .accent-color-options {
    gap: 10px;
  }
  
  .accent-color-option {
    width: 36px;
    height: 36px;
  }
  
  .font-settings {
    flex-direction: column;
    gap: 16px;
  }
  
  .font-family-select,
  .font-size-input,
  .font-size-unit {
    font-size: 12px;
    padding: 8px 12px;
  }
  
  .apply-button {
    font-size: 12px;
    padding: 8px 16px;
  }
  
  .font-size-input-with-button {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }
  
  .apply-button {
    width: 100%;
  }
}

/* ==================== 关于页面样式 ==================== */
.about-page-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
  padding: 24px;
  box-sizing: border-box;
  animation: fadeIn 0.4s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 关于页面头部 */
.about-header {
  text-align: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, var(--card-bg), rgba(255, 255, 255, 0.02));
  border-radius: 16px;
  border: 1px solid var(--border-color);
}

.dark .about-header {
  border-color: var(--border-strong);
}

.app-logo {
  margin-bottom: 20px;
}

.logo-icon {
  width: 96px;
  height: 96px;
  margin: 0 auto;
  background: linear-gradient(135deg, var(--active-color), var(--active-color-suppl));
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 8px 24px rgba(59, 130, 246, 0.3);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.logo-icon:hover {
  transform: scale(1.05) rotate(5deg);
  box-shadow: 0 12px 32px rgba(59, 130, 246, 0.4);
}

.app-name {
  font-size: 32px;
  font-weight: 800;
  color: var(--text-color);
  margin: 0 0 12px 0;
  letter-spacing: -0.5px;
}

.app-version {
  font-size: 16px;
  font-weight: 600;
  color: var(--active-color);
  margin: 0 0 12px 0;
}

.app-description {
  font-size: 15px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
}

/* 关于内容区域 */
.about-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
}

/* 关于卡片 */
.about-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .about-card {
  border-color: var(--border-strong);
}

.about-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: var(--active-color);
}

.dark .about-card:hover {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  background: rgba(0, 0, 0, 0.02);
}

.dark .card-header {
  border-bottom-color: var(--border-strong);
  background: rgba(255, 255, 255, 0.02);
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-color);
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-body {
  padding: 20px;
}

/* 作者信息 */
.author-info {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.author-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--active-color), var(--active-color-suppl));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.author-details h4 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--text-color);
}

.author-intro,
.author-field {
  margin: 0 0 6px 0;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.author-field {
  margin-bottom: 0;
}

/* GitHub 链接按钮 */
.github-links {
  display: flex;
  gap: 12px;
}

.github-link-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: linear-gradient(135deg, var(--bg-secondary), rgba(255, 255, 255, 0.02));
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-color);
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  justify-content: center;
}

.dark .github-link-btn {
  border-color: var(--border-strong);
  background: rgba(255, 255, 255, 0.03);
}

.github-link-btn:hover {
  background: #24292e;
  border-color: #24292e;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(36, 41, 46, 0.4);
}

.dark .github-link-btn:hover {
  box-shadow: 0 6px 20px rgba(36, 41, 46, 0.5);
}

.github-icon {
  transition: transform 0.3s ease;
}

.github-link-btn:hover .github-icon {
  transform: scale(1.1) rotate(-5deg);
}

/* 信息网格 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-color);
}

/* 技术标签 */
.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tech-tag {
  padding: 8px 14px;
  background: linear-gradient(135deg, var(--bg-secondary), rgba(255, 255, 255, 0.02));
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color);
  transition: all 0.3s ease;
}

.dark .tech-tag {
  border-color: var(--border-strong);
  background: rgba(255, 255, 255, 0.03);
}

.tech-tag:hover {
  border-color: var(--active-color);
  color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

/* 许可证信息 */
.license-info {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.license-icon {
  font-size: 40px;
  flex-shrink: 0;
}

.license-details h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-color);
}

.license-details p {
  margin: 0;
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
}

/* 特别感谢 */
.thanks-text {
  margin: 0;
  font-size: 15px;
  color: var(--text-secondary);
  line-height: 1.6;
  text-align: center;
}

/* 页脚 */
.about-footer {
  text-align: center;
  padding: 24px;
  border-top: 1px solid var(--border-color);
  margin-top: 20px;
}

.dark .about-footer {
  border-top-color: var(--border-strong);
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.footer-link {
  font-size: 14px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.3s ease;
  position: relative;
}

.footer-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 1px;
  background: var(--active-color);
  transition: width 0.3s ease;
}

.footer-link:hover {
  color: var(--active-color);
}

.footer-link:hover::after {
  width: 100%;
}

.copyright {
  margin: 0;
  font-size: 13px;
  color: var(--text-tertiary);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .about-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .about-page-container {
    padding: 16px;
    gap: 16px;
  }
  
  .about-header {
    padding: 30px 16px;
  }
  
  .logo-icon {
    width: 80px;
    height: 80px;
  }
  
  .app-name {
    font-size: 26px;
  }
  
  .app-version {
    font-size: 14px;
  }
  
  .card-header {
    padding: 14px 16px;
  }
  
  .card-body {
    padding: 16px;
  }
  
  .author-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .license-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .footer-links {
    gap: 16px;
    flex-direction: column;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .logo-icon {
    width: 64px;
    height: 64px;
  }
  
  .app-name {
    font-size: 22px;
  }
  
  .github-links {
    flex-direction: column;
  }
  
  .tech-tags {
    justify-content: center;
  }
}
</style>