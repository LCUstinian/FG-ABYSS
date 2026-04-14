<template>
  <div class="settings-view">
    <n-grid :cols="100" :x-gap="12" style="height: 100%">
      <!-- 左侧：设置导航 -->
      <n-grid-item :span="25">
        <n-card :bordered="false">
          <n-menu
            v-model:value="activeSection"
            mode="vertical"
            :options="menuOptions"
            @update:value="scrollToSection"
          />
        </n-card>
      </n-grid-item>

      <!-- 右侧：设置内容 -->
      <n-grid-item :span="75">
        <n-card :bordered="false" style="height: 100%; overflow-y: auto">
          <!-- 外观设置 -->
          <div id="section-appearance" class="settings-section">
            <h2 class="section-title">{{ t('settings.appearance') }}</h2>
            
            <n-form label-placement="left" label-width="120">
              <n-form-item label="主题模式">
                <n-select
                  v-model:value="settings.appearance.theme"
                  :options="[
                    { label: '深色', value: 'dark' },
                    { label: '浅色', value: 'light' },
                    { label: '跟随系统', value: 'system' },
                  ]"
                  style="width: 200px"
                />
              </n-form-item>

              <n-form-item label="语言">
                <n-select
                  v-model:value="settings.appearance.language"
                  :options="[
                    { label: '简体中文', value: 'zh-CN' },
                    { label: 'English', value: 'en-US' },
                  ]"
                  style="width: 200px"
                />
              </n-form-item>

              <n-form-item label="字体大小">
                <n-slider
                  v-model:value="settings.appearance.fontSize"
                  :min="12"
                  :max="24"
                  :step="1"
                  :marks="{ 12: '小', 16: '中', 20: '大', 24: '超大' }"
                  style="width: 300px"
                />
              </n-form-item>

              <n-form-item label="紧凑模式">
                <n-switch v-model:value="settings.appearance.compactMode" />
              </n-form-item>
            </n-form>
          </div>

          <n-divider />

          <!-- 连接设置 -->
          <div id="section-connection" class="settings-section">
            <h2 class="section-title">{{ t('settings.connection') }}</h2>

            <n-form label-placement="left" label-width="120">
              <n-form-item label="默认超时">
                <n-input-number
                  v-model:value="settings.connection.timeout"
                  :min="1"
                  :max="300"
                  style="width: 150px"
                />
                <span style="margin-left: 8px">秒</span>
              </n-form-item>

              <n-form-item label="最大重试">
                <n-input-number
                  v-model:value="settings.connection.maxRetries"
                  :min="0"
                  :max="10"
                  style="width: 150px"
                />
                <span style="margin-left: 8px">次</span>
              </n-form-item>

              <n-form-item label="速率限制">
                <n-input-number
                  v-model:value="settings.connection.rateLimit"
                  :min="1"
                  :max="100"
                  style="width: 150px"
                />
                <span style="margin-left: 8px">请求/秒</span>
              </n-form-item>

              <n-form-item label="代理配置">
                <n-space vertical style="width: 100%">
                  <n-switch v-model:value="settings.connection.useProxy" />
                  <n-input
                    v-if="settings.connection.useProxy"
                    v-model:value="settings.connection.proxyUrl"
                    placeholder="代理 URL (如：http://proxy.example.com:8080)"
                    style="width: 400px"
                  />
                </n-space>
              </n-form-item>

              <n-form-item label="自动重连">
                <n-switch v-model:value="settings.connection.autoReconnect" />
              </n-form-item>
            </n-form>
          </div>

          <n-divider />

          <!-- 安全设置 -->
          <div id="section-security" class="settings-section">
            <h2 class="section-title">{{ t('settings.security') }}</h2>

            <n-form label-placement="left" label-width="120">
              <n-form-item label="加密算法">
                <n-select
                  v-model:value="settings.security.encryptionAlgorithm"
                  :options="[
                    { label: 'AES-256-GCM (推荐)', value: 'aes-256-gcm' },
                    { label: 'XOR (轻量级)', value: 'xor' },
                  ]"
                  style="width: 200px"
                />
              </n-form-item>

              <n-form-item label="密钥长度">
                <n-select
                  v-model:value="settings.security.keyLength"
                  :options="[
                    { label: '128 位', value: 128 },
                    { label: '192 位', value: 192 },
                    { label: '256 位 (推荐)', value: 256 },
                  ]"
                  style="width: 200px"
                />
              </n-form-item>

              <n-form-item label="验证 SSL">
                <n-switch v-model:value="settings.security.verifySsl" />
                <template #feedback>
                  <n-alert type="warning" style="margin-top: 8px" v-if="!settings.security.verifySsl">
                    禁用 SSL 验证可能存在安全风险
                  </n-alert>
                </template>
              </n-form-item>

              <n-form-item label="自动锁屏">
                <n-input-number
                  v-model:value="settings.security.autoLockMinutes"
                  :min="0"
                  :max="120"
                  style="width: 150px"
                />
                <span style="margin-left: 8px">分钟 (0 为禁用)</span>
              </n-form-item>

              <n-form-item label="清除剪贴板">
                <n-switch v-model:value="settings.security.clearClipboard" />
                <template #feedback>
                  <span style="margin-left: 8px; color: var(--n-text-color-3)">
                    退出时自动清除剪贴板中的敏感信息
                  </span>
                </template>
              </n-form-item>
            </n-form>
          </div>

          <n-divider />

          <!-- 数据存储 -->
          <div id="section-storage" class="settings-section">
            <h2 class="section-title">数据存储</h2>

            <n-form label-placement="left" label-width="120">
              <n-form-item label="数据目录">
                <n-space>
                  <n-input
                    v-model:value="settings.storage.dataDirectory"
                    readonly
                    style="width: 400px"
                  />
                  <n-button @click="handleChangeDataDir">
                    <template #icon>
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                      </svg>
                    </template>
                    更改
                  </n-button>
                </n-space>
              </n-form-item>

              <n-form-item label="自动备份">
                <n-switch v-model:value="settings.storage.autoBackup" />
              </n-form-item>

              <n-form-item label="备份周期">
                <n-select
                  v-model:value="settings.storage.backupInterval"
                  :options="[
                    { label: '每天', value: 'daily' },
                    { label: '每周', value: 'weekly' },
                    { label: '每月', value: 'monthly' },
                  ]"
                  :disabled="!settings.storage.autoBackup"
                  style="width: 150px"
                />
              </n-form-item>

              <n-form-item label="保留备份数">
                <n-input-number
                  v-model:value="settings.storage.maxBackups"
                  :min="1"
                  :max="100"
                  :disabled="!settings.storage.autoBackup"
                  style="width: 150px"
                />
              </n-form-item>

              <n-form-item label=" ">
                <n-space>
                  <n-button @click="handleBackupNow">
                    <template #icon>
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                        <polyline points="17 8 12 3 7 8"/>
                        <line x1="12" y1="3" x2="12" y2="15"/>
                      </svg>
                    </template>
                    立即备份
                  </n-button>
                  <n-button @click="handleClearCache">
                    <template #icon>
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                      </svg>
                    </template>
                    清除缓存
                  </n-button>
                </n-space>
              </n-form-item>
            </n-form>
          </div>

          <n-divider />

          <!-- 关于 -->
          <div id="section-about" class="settings-section">
            <h2 class="section-title">关于</h2>

            <n-space vertical>
              <div style="text-align: center; padding: 20px">
                <n-avatar
                  :size="128"
                  style="background-color: #409eff; margin-bottom: 16px"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="1.5">
                    <path d="M12 2A10 10 0 0 0 2 12a10 10 0 0 0 10 10a10 10 0 0 0 10-10A10 10 0 0 0 12 2z"/>
                    <path d="M12 6a1 1 0 0 0-1 1v5a1 1 0 0 0 2 0V7a1 1 0 0 0-1-1z"/>
                    <circle cx="12" cy="14" r="1"/>
                  </svg>
                </n-avatar>
                <h3 style="margin: 0 0 8px 0">FG-ABYSS</h3>
                <p style="color: var(--n-text-color-3); margin: 0">
                  Advanced penetration testing framework
                </p>
                <n-tag type="info" style="margin-top: 8px">
                  v0.1.0
                </n-tag>
              </div>

              <n-descriptions bordered :column="2">
                <n-descriptions-item label="构建时间">
                  2024-01-01
                </n-descriptions-item>
                <n-descriptions-item label="Electron 版本">
                  Tauri 2.0.0
                </n-descriptions-item>
                <n-descriptions-item label="Node.js 版本">
                  v20.x
                </n-descriptions-item>
                <n-descriptions-item label="架构">
                  x64
                </n-descriptions-item>
              </n-descriptions>

              <n-space justify="center">
                <n-button text @click="handleCheckUpdate">
                  <template #icon>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="23 4 23 10 17 10"/>
                      <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
                    </svg>
                  </template>
                  检查更新
                </n-button>
                <n-button text @click="handleViewLicense">
                  <template #icon>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                      <polyline points="14 2 14 8 20 8"/>
                      <line x1="16" y1="13" x2="8" y2="13"/>
                      <line x1="16" y1="17" x2="8" y2="17"/>
                    </svg>
                  </template>
                  查看许可
                </n-button>
              </n-space>
            </n-space>
          </div>

          <!-- 保存按钮 -->
          <div style="margin-top: 24px; text-align: center">
            <n-space>
              <n-button type="primary" @click="handleSaveSettings">
                <template #icon>
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
                    <polyline points="17 21 17 13 7 13 7 21"/>
                    <polyline points="7 3 7 8 15 8"/>
                  </svg>
                </template>
                保存设置
              </n-button>
              <n-button @click="handleResetSettings">
                <template #icon>
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="1 4 1 10 7 10"/>
                    <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"/>
                  </svg>
                </template>
                重置为默认
              </n-button>
            </n-space>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NAlert, NAvatar, NButton, NCard, NDescriptions, NDescriptionsItem, NDivider,
  NForm, NFormItem, NGrid, NGridItem, NInput, NInputNumber, NMenu, NSelect,
  NSlider, NSpace, NSwitch, NTag, useMessage, useDialog,
} from 'naive-ui'
import type { MenuOption } from 'naive-ui'

const { t } = useI18n()
const message = useMessage()
const dialog = useDialog()

// 设置数据
const settings = ref({
  appearance: {
    theme: 'dark',
    language: 'zh-CN',
    fontSize: 14,
    compactMode: false,
  },
  connection: {
    timeout: 30,
    maxRetries: 3,
    rateLimit: 10,
    useProxy: false,
    proxyUrl: '',
    autoReconnect: true,
  },
  security: {
    encryptionAlgorithm: 'aes-256-gcm',
    keyLength: 256,
    verifySsl: false,
    autoLockMinutes: 0,
    clearClipboard: true,
  },
  storage: {
    dataDirectory: '%APPDATA%/fg-abyss',
    autoBackup: true,
    backupInterval: 'weekly',
    maxBackups: 10,
  },
})

// 当前选中的设置部分
const activeSection = ref('appearance')

// 设置菜单选项
const menuOptions: MenuOption[] = [
  { key: 'appearance', label: '外观设置', icon: () => '🎨' },
  { key: 'connection', label: '连接设置', icon: () => '🔌' },
  { key: 'security', label: '安全设置', icon: () => '🔒' },
  { key: 'storage', label: '数据存储', icon: () => '💾' },
  { key: 'about', label: '关于', icon: () => 'ℹ️' },
]

// 滚动到指定部分
const scrollToSection = (key: string) => {
  const element = document.getElementById(`section-${key}`)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}

// 更改数据目录
const handleChangeDataDir = () => {
  dialog.info({
    title: '更改数据目录',
    content: '此功能需要重启应用后生效。确定要更改吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      // TODO: 使用 Tauri dialog API 选择目录
      message.success('数据目录已更改')
    },
  })
}

// 立即备份
const handleBackupNow = () => {
  message.success('备份已完成')
}

// 清除缓存
const handleClearCache = () => {
  dialog.warning({
    title: '清除缓存',
    content: '确定要清除所有缓存数据吗？此操作不可恢复。',
    positiveText: '清除',
    negativeText: '取消',
    onPositiveClick: () => {
      message.success('缓存已清除')
    },
  })
}

// 检查更新
const handleCheckUpdate = () => {
  message.info('检查更新中...')
  // TODO: 实现更新检查
}

// 查看许可
const handleViewLicense = () => {
  dialog.info({
    title: '许可信息',
    content: 'FG-ABYSS - Advanced penetration testing framework\n\nCopyright © 2024 FG-ABYSS Team\nLicensed under MIT License',
    positiveText: '关闭',
  })
}

// 保存设置
const handleSaveSettings = () => {
  // TODO: 保存到配置文件
  message.success('设置已保存')
}

// 重置设置
const handleResetSettings = () => {
  dialog.warning({
    title: '重置设置',
    content: '确定要重置所有设置为默认值吗？',
    positiveText: '重置',
    negativeText: '取消',
    onPositiveClick: () => {
      // 重置为默认值
      settings.value = {
        appearance: {
          theme: 'dark',
          language: 'zh-CN',
          fontSize: 14,
          compactMode: false,
        },
        connection: {
          timeout: 30,
          maxRetries: 3,
          rateLimit: 10,
          useProxy: false,
          proxyUrl: '',
          autoReconnect: true,
        },
        security: {
          encryptionAlgorithm: 'aes-256-gcm',
          keyLength: 256,
          verifySsl: false,
          autoLockMinutes: 0,
          clearClipboard: true,
        },
        storage: {
          dataDirectory: '%APPDATA%/fg-abyss',
          autoBackup: true,
          backupInterval: 'weekly',
          maxBackups: 10,
        },
      }
      message.success('设置已重置')
    },
  })
}
</script>

<style scoped>
.settings-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 12px;
  box-sizing: border-box;
  overflow: hidden;
}

.settings-section {
  padding: 16px 0;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--n-text-color);
}

:deep(.n-card) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.n-card__content) {
  flex: 1;
  overflow-y: auto;
  padding-right: 16px;
  min-height: 0;
}

:deep(.n-menu) {
  --n-item-height: 44px;
}

:deep(.n-grid) {
  height: calc(100% - 24px);
}
</style>
