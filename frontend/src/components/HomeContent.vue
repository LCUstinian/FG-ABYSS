<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('home.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('home.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="home-content">
        <div class="home-card">
          <Server :size="48" />
          <h3>项目介绍</h3>
          <p>FG-ABYSS是一个功能强大的安全工具，用于管理项目、生成载荷和使用插件。</p>
          <p>这里将展示项目的主要功能和使用方法。</p>
        </div>
        <div class="home-card">
          <Cpu :size="48" />
          <h3>系统状态</h3>
          <p>内存使用: {{ systemStatus.memoryUsage }}</p>
          <p>CPU使用率: {{ systemStatus.cpuUsage }}</p>
          <p>运行时间: {{ systemStatus.uptime }}</p>
        </div>
        <div class="home-card">
          <Database :size="48" />
          <h3>许可证</h3>
          <p>本软件采用 MIT 许可证开源。</p>
          <p>您可以自由使用、修改和分发本软件。</p>
        </div>
        <div class="home-card">
          <FileText :size="48" />
          <h3>快速开始</h3>
          <p>创建您的第一个项目，开始使用 FG-ABYSS 的强大功能。</p>
          <p>探索载荷生成和插件管理功能。</p>
        </div>
        <div class="home-card">
          <Shield :size="48" />
          <h3>安全提示</h3>
          <p>请在合法授权的环境中使用本工具。</p>
          <p>遵守相关法律法规，保护网络安全。</p>
        </div>
        <div class="home-card">
          <HelpCircle :size="48" />
          <h3>帮助与支持</h3>
          <p>查看文档了解更多功能细节。</p>
          <p>加入社区获取技术支持。</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Server, Cpu, Database, FileText, Shield, HelpCircle } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps({
  systemStatus: {
    type: Object,
    required: true,
    default: () => ({
      memoryUsage: '0 GB / 0 GB',
      processId: '0',
      cpuUsage: '0%',
      uptime: '0 hours'
    })
  }
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

/* 首页内容样式 - 深色主题风格 */
.home-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  padding: 20px;
  width: 100%;
  min-height: 400px;
  box-sizing: border-box;
}

.home-card {
  background: var(--card-bg);
  border: 2px solid var(--border-color);
  box-shadow: var(--shadow-md);
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  border-radius: 12px;
}

.home-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg);
}

.home-card svg {
  flex-shrink: 0;
  color: var(--active-color);
}

.home-card h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
}

.home-card p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-header {
    padding: 12px 16px;
  }
  
  .content-body {
    padding: 16px;
  }
  
  .home-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>