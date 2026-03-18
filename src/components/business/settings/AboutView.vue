<template>
  <div class="about-view">
    <div class="about-container">
      <div class="about-header">
        <div class="app-logo">
          <span class="logo-icon">ℹ️</span>
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
              <span class="author-avatar">👨‍💻</span>
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
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<style scoped>
.about-view {
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

.about-container {
  display: flex;
  flex-direction: column;
  gap: 28px;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 32px 40px;
  box-sizing: border-box;
}

.about-header {
  text-align: center;
  padding: 48px 32px;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  width: 100%;
  max-width: 1000px;
  margin: 0 auto;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.about-header::before {
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

.about-header:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.about-header:hover::before {
  transform: scaleX(1);
}

.app-logo {
  margin-bottom: 20px;
  transition: all 0.3s ease;
}

.about-header:hover .app-logo {
  transform: scale(1.05);
}

.logo-icon {
  font-size: 72px;
  transition: transform 0.3s ease;
}

.about-header:hover .logo-icon {
  transform: rotate(5deg);
}

.app-name {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 12px 0;
  transition: all 0.3s ease;
  position: relative;
  display: inline-block;
}

.app-name::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 2px;
  background: var(--active-color);
  transition: width 0.3s ease;
}

.about-header:hover .app-name::after {
  width: 80%;
}

.app-version {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 12px 0;
  transition: color 0.3s ease;
}

.app-description {
  font-size: 15px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.7;
  max-width: 700px;
  margin: 0 auto;
  transition: color 0.3s ease;
}

.about-content {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.about-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
  width: 100%;
  max-width: 1000px;
  margin: 0 auto;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  position: relative;
}

.about-card::before {
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

.about-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.about-card:hover::before {
  transform: scaleX(1);
}

.card-header {
  padding: 24px 28px;
  background: var(--card-bg-hover);
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.about-card:hover .card-header {
  background: var(--card-bg);
  border-bottom-color: var(--active-color);
}

.card-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  transition: color 0.3s ease;
}

.about-card:hover .card-header h3 {
  color: var(--active-color);
}

.card-body {
  padding: 28px 32px;
  transition: all 0.3s ease;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 20px;
  transition: all 0.3s ease;
}

.about-card:hover .author-info {
  transform: translateX(8px);
}

.author-avatar {
  font-size: 56px;
  transition: transform 0.3s ease;
}

.about-card:hover .author-avatar {
  transform: scale(1.1) rotate(5deg);
}

.author-details h4 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  transition: color 0.3s ease;
}

.about-card:hover .author-details h4 {
  color: var(--active-color);
}

.author-details p {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
  transition: color 0.3s ease;
}

.github-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 12px 24px;
  background: var(--active-color-bg);
  color: var(--active-color);
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.github-link-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.2) 50%, transparent 100%);
  transition: left 0.6s ease;
}

.github-link-btn:hover::before {
  left: 100%;
}

.github-link-btn:hover {
  background: var(--active-color-bg-hover);
  transform: translateX(8px);
  box-shadow: 0 4px 12px rgba(var(--active-color-rgb), 0.2);
}

.github-link-btn svg {
  flex-shrink: 0;
  transition: transform 0.3s ease;
}

.github-link-btn:hover svg {
  transform: rotate(10deg);
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .about-container {
    padding: 0 40px 60px;
  }
  
  .about-header {
    max-width: 1100px;
    padding: 56px 40px;
  }
  
  .logo-icon {
    font-size: 80px;
  }
  
  .app-name {
    font-size: 40px;
  }
  
  .app-description {
    font-size: 16px;
  }
  
  .about-card {
    max-width: 1100px;
  }
  
  .card-body {
    padding: 32px 40px;
  }
  
  .card-header h3 {
    font-size: 20px;
  }
  
  .author-details h4 {
    font-size: 22px;
  }
  
  .author-details p {
    font-size: 15px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .about-container {
    padding: 0 24px 32px;
    gap: 16px;
  }
  
  .about-header {
    padding: 32px 24px;
  }
  
  .logo-icon {
    font-size: 56px;
  }
  
  .app-name {
    font-size: 28px;
  }
  
  .app-description {
    font-size: 14px;
  }
  
  .about-card {
    padding: 0;
  }
  
  .card-header {
    padding: 20px 24px;
  }
  
  .card-body {
    padding: 24px 28px;
  }
  
  .author-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 16px;
  }
  
  .author-avatar {
    font-size: 48px;
  }
  
  .github-link-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
