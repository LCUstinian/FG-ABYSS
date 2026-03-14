<template>
  <div class="content-section">
    <div class="content-header">
      <h1><span class="title">{{ t('payloads.title') }}</span> <span class="separator">|</span> <span class="subtitle">{{ t('payloads.subtitle') }}</span></h1>
    </div>
    <div class="content-body">
      <div class="payloads-content">
        <div class="payload-form">
          <div class="form-group">
            <label>{{ t('payloads.scriptType') }}</label>
            <select>
              <option>{{ t('payloads.php') }}</option>
              <option>{{ t('payloads.asp') }}</option>
              <option>{{ t('payloads.aspx') }}</option>
              <option>{{ t('payloads.jsp') }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>{{ t('payloads.funcType') }}</label>
            <select>
              <option>{{ t('payloads.basic') }}</option>
              <option>{{ t('payloads.fileManager') }}</option>
              <option>{{ t('payloads.cmdExec') }}</option>
              <option>{{ t('payloads.dbManager') }}</option>
              <option>{{ t('payloads.full') }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>{{ t('payloads.password') }}</label>
            <input type="password" :placeholder="t('payloads.passwordPlaceholder')">
          </div>
          <div class="form-group">
            <label>{{ t('payloads.obfuscation') }}</label>
            <select>
              <option>{{ t('payloads.noObfuscation') }}</option>
              <option>{{ t('payloads.lightObfuscation') }}</option>
              <option>{{ t('payloads.mediumObfuscation') }}</option>
              <option>{{ t('payloads.highObfuscation') }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>{{ t('payloads.filename') }}</label>
            <input type="text" :placeholder="t('payloads.filenamePlaceholder')">
          </div>
          <button class="generate-button">{{ t('payloads.generate') }}</button>
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
  padding: 0;
  margin: 0;
  box-sizing: border-box;
  background: var(--content-bg);
  border-top: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 载荷内容样式 - 深色主题风格 */
.payloads-content {
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  overflow-x: hidden;
  /* 确保背景纯净，无渐变或花纹 */
  background: transparent;
}

/* 深色主题下确保背景纯净 */
.dark .payloads-content {
  background: transparent;
}

.payload-form {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-md);
  padding: 28px;
  max-width: 560px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  border-radius: 8px;
  /* 确保表单背景为纯色 */
  background-image: none !important;
}

/* 深色主题下优化表单背景 */
.dark .payload-form {
  background: var(--bg-tertiary);
  border-color: var(--border-strong);
  background-image: none !important;
}

.form-group {
  margin-bottom: 20px;
  width: 100%;
  box-sizing: border-box;
}

.form-group label {
  display: block;
  margin-bottom: 10px;
  font-weight: 600;
  font-size: 14px;
  color: var(--text-color);
  width: 100%;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  gap: 8px;
}

.form-group label::before {
  content: '';
  width: 4px;
  height: 16px;
  background: var(--active-color);
  border-radius: 2px;
  flex-shrink: 0;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--content-bg);
  color: var(--text-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  box-sizing: border-box;
  /* 确保输入框背景纯净，无渐变或花纹 */
  background-image: none !important;
}

/* 深色主题下优化输入框和选择框 */
.dark .form-group input,
.dark .form-group select {
  background: var(--bg-tertiary);
  border-color: var(--border-strong);
  color: var(--text-primary);
  /* 深色主题下确保背景纯净 */
  background-image: none !important;
}

.form-group input::placeholder {
  color: var(--text-color);
  opacity: 0.5;
}

/* 深色主题下优化占位符 */
.dark .form-group input::placeholder {
  color: var(--text-secondary);
  opacity: 0.7;
}

.form-group select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23666' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  background-size: 16px;
  padding-right: 40px;
}

/* 深色主题下优化下拉箭头颜色 */
.dark .form-group select {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
}

.form-group select:focus {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23667eea' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  border-color: var(--active-color);
}

/* 深色主题下优化选中状态 */
.dark .form-group select:focus {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%233b82f6' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  border-color: var(--primary-500);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

/* 优化下拉选项在深色主题下的显示 */
.dark .form-group select option {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.generate-button {
  padding: 16px 36px;
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 140px;
  max-width: 200px;
  margin: 0 auto;
  box-sizing: border-box;
  letter-spacing: 0.5px;
}

.generate-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  background: var(--active-color);
  opacity: 0.9;
}

.generate-button:active {
  transform: translateY(0);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-header {
    padding: 12px 16px;
  }
  
  .content-body {
    padding: 16px;
  }
  
  .payload-form {
    padding: 20px;
  }
  
  .form-group {
    margin-bottom: 16px;
  }
  
  .generate-button {
    padding: 14px 32px;
    font-size: 14px;
  }
}
</style>