<template>
  <div class="plugins-management-view">
    <div class="content-body">
      <div class="plugins-content">
        <div class="plugins-header">
          <div class="plugins-tabs">
            <button 
              v-for="tab in tabs" 
              :key="tab.id"
              class="tab-button"
              :class="{ active: activeTab === tab.id }"
              @click="activeTab = tab.id"
            >
              {{ tab.label }}
            </button>
          </div>
          <button class="install-button" @click="openInstallDialog">
            <Plus :size="16" />
            {{ t('plugins.installPlugin') }}
          </button>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <p>{{ t('plugins.loadingPlugins') }}</p>
        </div>

        <div v-else-if="error" class="error-state">
          <AlertCircle :size="48" />
          <p>{{ error }}</p>
          <button class="retry-button" @click="loadPlugins">
            {{ t('plugins.retry') }}
          </button>
        </div>

        <div v-else-if="filteredPlugins.length === 0" class="empty-state">
          <Plug :size="64" />
          <h3>{{ t('plugins.noPlugins') }}</h3>
          <p>{{ t('plugins.noPluginsDescription1') }}</p>
          <p>{{ t('plugins.noPluginsDescription2') }}</p>
        </div>

        <div v-else class="plugins-list">
          <div v-for="plugin in filteredPlugins" :key="plugin.id" class="plugin-card">
            <div class="plugin-info">
              <div class="plugin-header">
                <h3>{{ plugin.name }}</h3>
                <span v-if="plugin.isBuiltIn" class="builtin-badge">{{ t('plugins.builtIn') }}</span>
              </div>
              <p class="plugin-version">{{ t('plugins.version') }}: {{ plugin.version }}</p>
              <p class="plugin-description">{{ plugin.description }}</p>
              <p class="plugin-author">{{ t('plugins.author') }}: {{ plugin.author }}</p>
              <div class="plugin-signature">
                <span class="signature-label">{{ t('plugins.signature') }}:</span>
                <span class="signature-value" :class="{ verified: plugin.signature === 'BUILT-IN' || plugin.signature === 'VERIFIED' }">
                  {{ plugin.signature === 'BUILT-IN' ? t('plugins.builtIn') : 
                     plugin.signature === 'VERIFIED' ? t('plugins.verified') : 
                     t('plugins.unverified') }}
                </span>
              </div>
            </div>
            <div class="plugin-actions">
              <button 
                class="action-button" 
                :class="{ primary: !plugin.enabled }"
                @click="togglePlugin(plugin)"
              >
                {{ plugin.enabled ? t('plugins.disable') : t('plugins.enable') }}
              </button>
              <button 
                v-if="!plugin.isBuiltIn" 
                class="action-button danger"
                @click="uninstallPlugin(plugin)"
              >
                {{ t('plugins.uninstall') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plug, Plus, AlertCircle } from 'lucide-vue-next'
import { usePlugin } from '@/composables/usePlugin'

const { t } = useI18n()
const { plugins, loading, error, loadPlugins, performPluginAction } = usePlugin()

const activeTab = ref('local')
const tabs = [
  { id: 'local', label: t('plugins.localPlugins') },
  { id: 'store', label: t('plugins.pluginStore') }
]

const filteredPlugins = computed(() => {
  if (activeTab.value === 'local') {
    return plugins.value
  }
  // 插件商店暂时返回空数组
  return []
})

const togglePlugin = async (plugin: any) => {
  const action = plugin.enabled ? 'disable' : 'enable'
  const result = await performPluginAction({
    pluginId: plugin.id,
    action: action as any
  })
  if (result.success) {
    // 可以添加通知
    console.log(result.message)
  }
}

const uninstallPlugin = async (plugin: any) => {
  if (confirm(t('plugins.confirmUninstall'))) {
    const result = await performPluginAction({
      pluginId: plugin.id,
      action: 'uninstall'
    })
    if (result.success) {
      // 可以添加通知
      console.log(result.message)
    }
  }
}

const openInstallDialog = () => {
  // 这里可以实现打开文件选择对话框的逻辑
  // 实际实现中，应该调用 Tauri 的文件选择 API
  alert(t('plugins.selectPluginFile'))
}

onMounted(() => {
  loadPlugins()
})
</script>

<style scoped>
.plugins-management-view {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

.content-body {
  flex: 1;
  width: 100%;
  padding: 0;
  margin: 0;
  box-sizing: border-box;
  background: var(--content-bg);
  border-top: none;
  display: flex;
  align-items: stretch;
}

.plugins-content {
  width: 100%;
  height: 100%;
  padding: 24px;
  box-sizing: border-box;
}

.plugins-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.plugins-tabs {
  display: flex;
  gap: 8px;
}

.tab-button {
  padding: 10px 20px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  transition: all 0.3s ease;
}

.tab-button:hover {
  background: var(--card-bg-hover);
  color: var(--text-primary);
}

.tab-button.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
}

.install-button {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: var(--active-color);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: white;
  transition: all 0.3s ease;
}

.install-button:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  gap: 16px;
  min-height: 400px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-color);
  border-top: 3px solid var(--active-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  gap: 16px;
  min-height: 400px;
}

.error-state svg {
  color: var(--error-color);
  opacity: 0.8;
}

.error-state p {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
}

.retry-button {
  padding: 8px 16px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.retry-button:hover {
  background: var(--card-bg-hover);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  gap: 16px;
  min-height: 400px;
}

.empty-state svg {
  color: var(--active-color);
  opacity: 0.6;
  margin-bottom: 16px;
}

.empty-state h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
}

.empty-state p {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.plugins-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.plugin-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.plugin-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.plugin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.plugin-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.builtin-badge {
  padding: 4px 12px;
  background: var(--active-color);
  color: white;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.plugin-version {
  margin: 0;
  font-size: 12px;
  color: var(--text-secondary);
}

.plugin-description {
  margin: 0;
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.5;
}

.plugin-author {
  margin: 0;
  font-size: 12px;
  color: var(--text-secondary);
}

.plugin-signature {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.signature-label {
  color: var(--text-secondary);
}

.signature-value {
  color: var(--error-color);
  font-weight: 500;
}

.signature-value.verified {
  color: var(--success-color);
}

.plugin-actions {
  display: flex;
  gap: 10px;
  margin-top: auto;
}

.action-button {
  flex: 1;
  padding: 8px 16px;
  background: var(--card-bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.action-button:hover {
  background: var(--card-bg);
}

.action-button.primary {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
}

.action-button.primary:hover {
  opacity: 0.9;
}

.action-button.danger {
  background: var(--error-color);
  color: white;
  border-color: var(--error-color);
}

.action-button.danger:hover {
  opacity: 0.9;
}

@media (max-width: 768px) {
  .plugins-list {
    grid-template-columns: 1fr;
  }
  
  .plugins-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
