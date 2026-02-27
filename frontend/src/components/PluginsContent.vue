<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('plugins.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('plugins.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="plugins-content">
        <div class="plugins-tabs">
          <button class="tab-button active">{{ t('plugins.localPlugins') }}</button>
          <button class="tab-button">{{ t('plugins.pluginStore') }}</button>
        </div>
        <div class="empty-state">
          <Plug :size="64" />
          <h3>暂无插件</h3>
          <p>您还没有安装任何插件</p>
          <p>从插件商店中浏览并安装插件来扩展功能</p>
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

/* 插件内容样式 - 深色主题风格 */
.plugins-content {
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  min-height: 400px;
  box-sizing: border-box;
  padding: 20px;
}

.plugins-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 12px;
}

.tab-button {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  background: var(--card-bg);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
  font-size: 14px;
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 90px;
}

.tab-button:hover {
  background: var(--hover-color);
  transform: translateY(-2px);
}

.tab-button.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  box-shadow: var(--shadow-md);
}

.plugins-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  margin-top: 20px;
}

.plugin-item {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-md);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 8px;
}

.plugin-item:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.plugin-item h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
}

.plugin-item p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.8;
}

.plugin-button {
  margin-top: 12px;
  padding: 12px 24px;
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  width: fit-content;
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 90px;
}

.plugin-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  background: var(--active-color);
  opacity: 0.9;
}

/* 空状态样式 */
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
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.7;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-header {
    padding: 12px 16px;
  }
  
  .content-body {
    padding: 16px;
  }
  
  .plugins-tabs {
    gap: 8px;
  }
  
  .tab-button {
    padding: 6px 12px;
    font-size: 12px;
  }
  
  .plugins-list {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .plugin-item {
    padding: 20px;
  }
  
  .plugin-button {
    padding: 10px 20px;
    font-size: 12px;
  }
}
</style>