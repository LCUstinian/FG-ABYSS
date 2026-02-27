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
            <span>代理</span>
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
                            placeholder="建议8-24px"
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
                <h4>代理设置</h4>
                <div class="placeholder-content">
                  <div class="placeholder-icon">🌐</div>
                  <p>代理设置功能正在开发中...</p>
                  <p>这里将提供HTTP代理、SOCKS代理等配置选项</p>
                  <div class="placeholder-tips">
                    <span>提示：</span>
                    <p>您可以在这里设置代理服务器，用于网络请求的转发</p>
                  </div>
                </div>
              </div>
            </template>
            
            <!-- 网络设置 -->
            <template v-else-if="currentSettingsTab === 'network'">
              <div class="settings-card">
                <h4>网络设置</h4>
                <div class="placeholder-content">
                  <div class="placeholder-icon">📡</div>
                  <p>网络设置功能正在开发中...</p>
                  <p>这里将提供网络连接、超时设置等配置选项</p>
                  <div class="placeholder-tips">
                    <span>提示：</span>
                    <p>您可以在这里配置网络连接参数，优化网络性能</p>
                  </div>
                </div>
              </div>
            </template>
            
            <!-- 关于 -->
            <template v-else-if="currentSettingsTab === 'about'">
              <div class="settings-card">
                <h4>关于</h4>
                <div class="placeholder-content">
                  <div class="placeholder-icon">ℹ️</div>
                  <p>FG-ABYSS</p>
                  <p>版本：1.0.0</p>
                  <p>一个功能强大的网络安全工具</p>
                  <div class="placeholder-tips">
                    <span>提示：</span>
                    <p>这里将显示软件版本信息、更新日志和相关链接</p>
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

// 强调色选项
const accentColors = [
  { value: '#667eea' }, // 默认蓝色
  { value: '#764ba2' }, // 紫色
  { value: '#f093fb' }, // 粉色
  { value: '#f5576c' }, // 红色
  { value: '#4facfe' }, // 浅蓝色
  { value: '#43e97b' }, // 绿色
  { value: '#fa709a' }, // 玫红色
  { value: '#fee140' }, // 黄色
  { value: '#30cfd0' }, // 青色
  { value: '#949494' }, // 灰色
  { value: '#1e40af' }, // 深蓝色
  { value: '#059669' }, // 深绿色
  { value: '#f97316' }, // 橙色
  { value: '#7e22ce' }, // 深紫色
  { value: '#92400e' }  // 棕褐色
]

// 当前强调色
const currentAccentColor = ref(localStorage.getItem('accentColor') || '#667eea')

// 字体选项
const fontFamilies = [
  { name: '系统默认', value: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, sans-serif' },
  { name: '微软雅黑', value: '"Microsoft YaHei", sans-serif' },
  { name: '宋体', value: 'SimSun, serif' },
  { name: '黑体', value: 'SimHei, sans-serif' },
  { name: 'Arial', value: 'Arial, sans-serif' },
  { name: 'Helvetica', value: 'Helvetica, sans-serif' },
  { name: 'Times New Roman', value: '"Times New Roman", serif' },
  { name: 'Courier New', value: '"Courier New", monospace' }
]

// 字体大小选项
const fontSizes = [
  { label: '极小', value: '10px' },
  { label: '小', value: '12px' },
  { label: '中', value: '14px' },
  { label: '大', value: '16px' },
  { label: '极大', value: '18px' }
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
  localStorage.setItem('language', lang)
}

// 处理强调色切换
const changeAccentColor = (color: string) => {
  currentAccentColor.value = color
  localStorage.setItem('accentColor', color)
  document.documentElement.style.setProperty('--active-color', color)
  // 触发storage事件，让其他组件知道颜色变化
  window.dispatchEvent(new StorageEvent('storage', {
    key: 'accentColor',
    newValue: color
  }))
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

// 初始化语言、强调色、字体和字体大小
onMounted(() => {
  const savedLanguage = localStorage.getItem('language')
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
  
  // 监听localStorage中变化
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
  padding: 16px 24px;
  margin-bottom: 0;
  background: var(--panel-bg);
  border-bottom: none;
  box-shadow: none;
}

.content-section h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-color);
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
  line-height: 1.2;
}

.content-body {
  flex: 1;
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;
  box-sizing: border-box;
  background: var(--content-bg);
  border-top: none;
  display: flex;
  align-items: stretch;
}

.content-section h1 .title {
  font-weight: 700;
  color: white;
}

.content-section h1 .separator {
  color: white;
  opacity: 0.7;
  font-weight: 400;
}

.content-section h1 .subtitle {
  font-size: 14px;
  font-weight: 400;
  color: white;
  opacity: 0.8;
  font-style: normal;
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
  height: 100%;
  background: var(--content-bg);
  padding: 0;
  margin: 0;
  overflow: hidden;
  border: none;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.05);
  z-index: 5;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
}

.dark .settings-main {
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.4);
}

.settings-panel {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px;
  border: none;
  overflow: hidden;
  box-sizing: border-box;
}

.settings-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.settings-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.settings-card h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  opacity: 0.8;
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

.theme-option:hover {
  border-color: var(--active-color);
}

.theme-option.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.theme-icon {
  font-size: 28px;
  display: block;
  margin: 0 auto;
}

.theme-option span:nth-child(2) {
  font-size: 14px;
  font-weight: 500;
  display: block;
  margin: 0 auto;
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
}

.dark .theme-check {
  background: white;
  color: var(--active-color);
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

.language-option:hover {
  border-color: var(--active-color);
}

.language-option.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.language-icon {
  font-size: 28px;
  display: block;
  margin: 0 auto;
}

.language-option span:nth-child(2) {
  font-size: 14px;
  font-weight: 500;
  display: block;
  margin: 0 auto;
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
}

.dark .language-check {
  background: white;
  color: var(--active-color);
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

.accent-color-option.active {
  border-color: var(--text-color);
  box-shadow: 0 0 12px rgba(0, 0, 0, 0.3);
  transform: scale(1.15);
}

.accent-color-check {
  color: white;
  font-size: 16px;
  font-weight: bold;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
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

.font-family-select:hover {
  border-color: var(--active-color);
}

.font-family-select:focus {
  border-color: var(--active-color);
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
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
}

.apply-button {
  padding: 10px 20px;
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
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
}

.font-size-input::placeholder {
  color: var(--text-color);
  opacity: 0.5;
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
}

.apply-button:hover {
  opacity: 0.9;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
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
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.placeholder-content p {
  margin: 0;
  color: var(--text-color);
  opacity: 0.8;
  line-height: 1.6;
}

.placeholder-tips {
  margin-top: 20px;
  padding: 16px;
  background: var(--hover-color);
  border-radius: 8px;
  width: 100%;
  max-width: 400px;
  text-align: left;
}

.placeholder-tips span {
  font-weight: 600;
  color: var(--text-color);
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
</style>