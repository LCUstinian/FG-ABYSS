<template>
  <div class="plugins-management-view">
    <div class="content-body">
      <div class="plugins-content">
        <div class="plugins-tabs">
          <button class="tab-button active">{{ t('plugins.localPlugins') }}</button>
          <button class="tab-button">{{ t('plugins.pluginStore') }}</button>
        </div>
        <div class="empty-state">
          <Plug :size="64" />
          <h3>{{ t('plugins.noPlugins') }}</h3>
          <p>{{ t('plugins.noPluginsDescription1') }}</p>
          <p>{{ t('plugins.noPluginsDescription2') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Plug } from 'lucide-vue-next'

const { t } = useI18n()
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

.plugins-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
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
</style>
