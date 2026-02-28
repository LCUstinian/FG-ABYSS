<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('home.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('home.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="home-content">
        <div class="home-card">
          <Server :size="48" />
          <h3>{{ t('home.title') }}</h3>
          <p>{{ t('home.projectIntro') }}</p>
          <p>{{ t('home.projectDescription') }}</p>
        </div>
        <div class="home-card">
          <Cpu :size="48" />
          <h3>{{ t('home.systemStatus') }}</h3>
          <p>{{ t('home.memoryUsage') }} {{ systemStatus.memoryUsage }}</p>
          <p>{{ t('home.cpuUsage') }} {{ systemStatus.cpuUsage }}</p>
          <p>{{ t('home.uptime') }} {{ systemStatus.uptime }}</p>
        </div>
        <div class="home-card">
          <Database :size="48" />
          <h3>{{ t('home.license') }}</h3>
          <p>{{ t('home.licenseText') }}</p>
          <p>{{ t('home.licenseDescription') }}</p>
        </div>
        <div class="home-card">
          <FileText :size="48" />
          <h3>{{ t('home.quickStart') }}</h3>
          <p>{{ t('home.quickStartDescription1') }}</p>
          <p>{{ t('home.quickStartDescription2') }}</p>
        </div>
        <div class="home-card">
          <Shield :size="48" />
          <h3>{{ t('home.securityTips') }}</h3>
          <p>{{ t('home.securityTip1') }}</p>
          <p>{{ t('home.securityTip2') }}</p>
        </div>
        <div class="home-card">
          <HelpCircle :size="48" />
          <h3>{{ t('home.helpSupport') }}</h3>
          <p>{{ t('home.helpSupport1') }}</p>
          <p>{{ t('home.helpSupport2') }}</p>
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
  padding: 20px 24px;
  margin-bottom: 0;
  background: linear-gradient(135deg, var(--panel-bg), #2563eb);
  border-bottom: none;
  box-shadow: var(--shadow-md);
  position: relative;
  overflow: hidden;
}

.content-header::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  animation: pulse 6s ease-in-out infinite;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.1); opacity: 0.8; }
  100% { transform: scale(1); opacity: 0.5; }
}

.content-section h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: white;
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
  line-height: 1.2;
  position: relative;
  z-index: 1;
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
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.content-section h1 .separator {
  color: white;
  opacity: 0.8;
  font-weight: 400;
}

.content-section h1 .subtitle {
  font-size: 16px;
  font-weight: 400;
  color: white;
  opacity: 0.9;
  font-style: normal;
}

/* 首页内容样式 - 现代风格 */
.home-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  padding: 32px;
  width: 100%;
  min-height: 400px;
  box-sizing: border-box;
  max-width: 1400px;
  margin: 0 auto;
}

.home-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-md);
  padding: 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
  transition: all var(--transition-normal);
  width: 100%;
  box-sizing: border-box;
  border-radius: var(--border-radius-xl);
  position: relative;
  overflow: hidden;
}

.home-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--active-color), #60a5fa);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform var(--transition-slow);
}

.home-card:hover::before {
  transform: scaleX(1);
}

.home-card:hover {
  transform: translateY(-8px);
  box-shadow: var(--shadow-xl);
  border-color: var(--active-color);
}

.home-card svg {
  flex-shrink: 0;
  color: var(--active-color);
  transition: transform var(--transition-normal);
  position: relative;
  z-index: 1;
}

.home-card:hover svg {
  transform: scale(1.1) rotate(5deg);
}

.home-card h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-color);
  position: relative;
  z-index: 1;
  transition: color var(--transition-normal);
}

.home-card:hover h3 {
  color: var(--active-color);
}

.home-card p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.8;
  position: relative;
  z-index: 1;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .home-content {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 20px;
    padding: 24px;
  }
}

@media (max-width: 768px) {
  .content-header {
    padding: 16px 20px;
  }
  
  .content-section h1 {
    font-size: 20px;
  }
  
  .content-section h1 .subtitle {
    font-size: 14px;
  }
  
  .home-content {
    grid-template-columns: 1fr;
    gap: 16px;
    padding: 20px;
  }
  
  .home-card {
    padding: 24px;
  }
  
  .home-card svg {
    width: 40px;
    height: 40px;
  }
  
  .home-card h3 {
    font-size: 18px;
  }
}

@media (max-width: 480px) {
  .content-header {
    padding: 12px 16px;
  }
  
  .content-section h1 {
    font-size: 18px;
    gap: 8px;
  }
  
  .home-content {
    padding: 16px;
  }
  
  .home-card {
    padding: 20px;
  }
}
</style>