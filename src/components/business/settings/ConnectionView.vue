<template>
  <div class="connection-view">
    <div class="connection-container">
      <!-- 代理设置卡片 -->
      <div class="settings-card proxy-card">
        <div class="card-header-section">
          <div class="card-icon-wrapper">
            <span class="card-icon">🌐</span>
          </div>
          <div class="card-title-section">
            <h4 class="card-title">{{ t('settings.proxySettings') }}</h4>
            <p class="card-description">{{ t('settings.proxySettingsDescription') }}</p>
          </div>
        </div>
        <div class="card-content">
          <div class="proxy-settings-form">
            <div class="form-item">
              <label class="form-label">{{ t('settings.proxyEnabled') }}</label>
              <n-switch v-model:value="proxyEnabled" />
            </div>
            
            <div class="form-item" v-show="proxyEnabled">
              <label class="form-label">{{ t('settings.proxyType') }}</label>
              <n-select
                v-model:value="proxyType"
                :options="proxyTypeOptions"
                style="width: 200px;"
              />
            </div>
            
            <div class="form-item" v-show="proxyEnabled">
              <label class="form-label">{{ t('settings.proxyHost') }}</label>
              <n-input
                v-model:value="proxyHost"
                placeholder="127.0.0.1"
                style="width: 200px;"
              />
            </div>
            
            <div class="form-item" v-show="proxyEnabled">
              <label class="form-label">{{ t('settings.proxyPort') }}</label>
              <n-input-number
                v-model:value="proxyPort"
                :min="1"
                :max="65535"
                placeholder="7890"
                style="width: 120px;"
              />
            </div>
            
            <div class="form-actions" v-show="proxyEnabled">
              <n-button @click="testProxy" :loading="testingProxy">
                {{ t('settings.testProxy') }}
              </n-button>
              <n-button type="primary" @click="saveProxySettings">
                {{ t('settings.save') }}
              </n-button>
            </div>
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
            <p class="card-description">{{ t('settings.networkSettingsDescription') }}</p>
          </div>
        </div>
        <div class="card-content">
          <div class="network-settings-form">
            <div class="form-item">
              <label class="form-label">{{ t('settings.connectionTimeout') }}</label>
              <n-input-number
                v-model:value="connectionTimeout"
                :min="1"
                :max="300"
                :suffix="t('settings.seconds')"
                style="width: 150px;"
              />
            </div>
            
            <div class="form-item">
              <label class="form-label">{{ t('settings.maxRetries') }}</label>
              <n-input-number
                v-model:value="maxRetries"
                :min="0"
                :max="10"
                style="width: 150px;"
              />
            </div>
            
            <div class="form-item">
              <label class="form-label">{{ t('settings.verifySSL') }}</label>
              <n-switch v-model:value="verifySSL" />
            </div>
            
            <div class="form-actions">
              <n-button type="primary" @click="saveNetworkSettings">
                {{ t('settings.save') }}
              </n-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useMessage } from 'naive-ui'

const { t } = useI18n()
const message = useMessage()

// 代理设置
const proxyEnabled = ref(false)
const proxyType = ref('http')
const proxyHost = ref('127.0.0.1')
const proxyPort = ref<number | null>(7890)
const testingProxy = ref(false)

// 网络设置
const connectionTimeout = ref(30)
const maxRetries = ref(3)
const verifySSL = ref(true)

// 代理类型选项
const proxyTypeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' }
]

// 测试代理连接
const testProxy = async () => {
  testingProxy.value = true
  try {
    // 模拟测试
    await new Promise(resolve => setTimeout(resolve, 1500))
    message.success(t('settings.proxyTestSuccess'))
  } catch (error) {
    message.error(t('settings.proxyTestFailed'))
  } finally {
    testingProxy.value = false
  }
}

// 保存代理设置
const saveProxySettings = () => {
  const settings = {
    enabled: proxyEnabled.value,
    type: proxyType.value,
    host: proxyHost.value,
    port: proxyPort.value
  }
  localStorage.setItem('proxySettings', JSON.stringify(settings))
  message.success(t('settings.settingsSaved'))
}

// 保存网络设置
const saveNetworkSettings = () => {
  const settings = {
    timeout: connectionTimeout.value,
    maxRetries: maxRetries.value,
    verifySSL: verifySSL.value
  }
  localStorage.setItem('networkSettings', JSON.stringify(settings))
  message.success(t('settings.settingsSaved'))
}

// 加载保存的设置
onMounted(() => {
  const savedProxy = localStorage.getItem('proxySettings')
  if (savedProxy) {
    const settings = JSON.parse(savedProxy)
    proxyEnabled.value = settings.enabled || false
    proxyType.value = settings.type || 'http'
    proxyHost.value = settings.host || '127.0.0.1'
    proxyPort.value = settings.port || 7890
  }
  
  const savedNetwork = localStorage.getItem('networkSettings')
  if (savedNetwork) {
    const settings = JSON.parse(savedNetwork)
    connectionTimeout.value = settings.timeout || 30
    maxRetries.value = settings.maxRetries || 3
    verifySSL.value = settings.verifySSL !== false
  }
})
</script>

<style scoped>
.connection-view {
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

.connection-container {
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

.placeholder-content:hover .placeholder-icon {
  opacity: 0.8;
  transform: scale(1.1);
  transition: all 0.3s ease;
}

.proxy-settings-form,
.network-settings-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  animation: formFadeIn 0.5s ease;
}

@keyframes formFadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.form-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  transition: all 0.3s ease;
  border-bottom: 1px solid transparent;
}

.form-item:hover {
  border-bottom-color: var(--border-color);
  padding-left: 8px;
}

.form-label {
  min-width: 140px;
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
  transition: color 0.3s ease;
}

.form-item:hover .form-label {
  color: var(--text-primary);
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.form-actions:hover {
  border-top-color: var(--active-color);
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .connection-container {
    padding: 0 40px 60px;
  }
  
  .settings-card {
    max-width: 1100px;
    padding: 32px 36px;
  }
  
  .card-title {
    font-size: 20px;
  }
  
  .card-description {
    font-size: 15px;
  }
  
  .form-item {
    gap: 20px;
  }
  
  .form-label {
    min-width: 160px;
    font-size: 15px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .connection-container {
    padding: 0 24px 32px;
    gap: 16px;
  }
  
  .settings-card {
    padding: 24px 28px;
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
  
  .form-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    padding: 12px 0;
  }
  
  .form-label {
    min-width: auto;
    font-size: 13px;
  }
  
  .form-actions {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
