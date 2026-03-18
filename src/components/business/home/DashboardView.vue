<template>
  <div class="dashboard-view">
    <div class="dashboard-container">
        <!-- 核心数据指标卡片 -->
        <div class="metrics-grid">
          <div class="metric-card" @click="navigateTo('projects')">
            <div class="metric-card-header">
              <div class="metric-icon projects">
                <Folder :size="24" />
              </div>
              <div class="metric-trend positive">
                <TrendingUp :size="16" />
                <span>+12%</span>
              </div>
            </div>
            <div class="metric-card-body">
              <div class="metric-value">24</div>
              <div class="metric-label">{{ t('home.totalProjects') }}</div>
              <div class="metric-progress">
                <div class="progress-bar">
                  <div class="progress-fill projects" style="width: 75%"></div>
                </div>
                <span class="progress-text">75% 使用率</span>
              </div>
            </div>
          </div>

          <div class="metric-card" @click="navigateTo('webshells')">
            <div class="metric-card-header">
              <div class="metric-icon webshells">
                <Terminal :size="24" />
              </div>
              <div class="metric-trend positive">
                <TrendingUp :size="16" />
                <span>+8%</span>
              </div>
            </div>
            <div class="metric-card-body">
              <div class="metric-value">156</div>
              <div class="metric-label">{{ t('home.totalWebShells') }}</div>
              <div class="metric-progress">
                <div class="progress-bar">
                  <div class="progress-fill webshells" style="width: 60%"></div>
                </div>
                <span class="progress-text">60% 活跃</span>
              </div>
            </div>
          </div>

          <div class="metric-card" @click="navigateTo('payloads')">
            <div class="metric-card-header">
              <div class="metric-icon payloads">
                <Box :size="24" />
              </div>
              <div class="metric-trend negative">
                <TrendingDown :size="16" />
                <span>-3%</span>
              </div>
            </div>
            <div class="metric-card-body">
              <div class="metric-value">89</div>
              <div class="metric-label">{{ t('home.totalPayloads') }}</div>
              <div class="metric-progress">
                <div class="progress-bar">
                  <div class="progress-fill payloads" style="width: 45%"></div>
                </div>
                <span class="progress-text">45% 使用率</span>
              </div>
            </div>
          </div>

          <div class="metric-card" @click="navigateTo('plugins')">
            <div class="metric-card-header">
              <div class="metric-icon plugins">
                <Puzzle :size="24" />
              </div>
              <div class="metric-trend positive">
                <TrendingUp :size="16" />
                <span>+15%</span>
              </div>
            </div>
            <div class="metric-card-body">
              <div class="metric-value">12</div>
              <div class="metric-label">{{ t('home.totalPlugins') }}</div>
              <div class="metric-progress">
                <div class="progress-bar">
                  <div class="progress-fill plugins" style="width: 30%"></div>
                </div>
                <span class="progress-text">30% 使用率</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 系统状态卡片 -->
        <div class="dashboard-grid">
          <!-- 系统状态 -->
          <div class="dashboard-card system-status">
            <div class="card-header">
              <Cpu :size="20" />
              <h3>{{ t('home.systemStatus') }}</h3>
            </div>
            <div class="card-body">
              <div class="status-grid">
                <div class="status-item">
                  <div class="status-label">{{ t('home.memoryUsage') }}</div>
                  <div class="status-value">
                    <n-progress
                      type="line"
                      :percentage="getMemoryPercentage(props.systemStatus?.memoryUsage || '0 GB / 0 GB')"
                      :indicator-placement="'inside'"
                      :height="8"
                      :border-radius="4"
                    />
                    <span class="status-text">{{ props.systemStatus?.memoryUsage || 'N/A' }}</span>
                  </div>
                </div>

                <div class="status-item">
                  <div class="status-label">{{ t('home.cpuUsage') }}</div>
                  <div class="status-value cpu-bars">
                    <div
                      v-for="i in 8"
                      :key="i"
                      class="cpu-bar"
                      :style="{ height: getCPUBarHeight(props.systemStatus?.cpuUsage || '0%', i - 1) + '%' }"
                    />
                    <span class="status-text">{{ props.systemStatus?.cpuUsage || 'N/A' }}</span>
                  </div>
                </div>

                <div class="status-item">
                  <div class="status-label">{{ t('home.uptime') }}</div>
                  <div class="status-value">
                    <Zap :size="16" />
                    <span class="status-text">{{ formatUptimeDisplay(props.systemStatus?.uptime || '0') }}</span>
                  </div>
                </div>

                <div class="status-item">
                  <div class="status-label">{{ t('home.processId') }}</div>
                  <div class="status-value">
                    <Clock :size="16" />
                    <span class="status-text">PID: {{ props.systemStatus?.processId || 'N/A' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 快捷操作 -->
          <div class="dashboard-card quick-actions">
            <div class="card-header">
              <Zap :size="20" />
              <h3>{{ t('home.quickActions') }}</h3>
            </div>
            <div class="card-body">
              <div class="actions-grid">
                <button class="action-button" @click="navigateTo('projects')">
                  <Plus :size="20" />
                  <span>{{ t('home.newProject') }}</span>
                </button>
                <button class="action-button" @click="navigateTo('webshells')">
                  <Terminal :size="20" />
                  <span>{{ t('home.newWebShell') }}</span>
                </button>
                <button class="action-button" @click="navigateTo('payloads')">
                  <Box :size="20" />
                  <span>{{ t('home.generatePayload') }}</span>
                </button>
                <button class="action-button" @click="navigateTo('plugins')">
                  <Puzzle :size="20" />
                  <span>{{ t('home.managePlugins') }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 活动日志 -->
        <div class="dashboard-card activity-log">
          <div class="card-header">
            <Clock :size="20" />
            <h3>{{ t('home.recentActivity') }}</h3>
          </div>
          <div class="card-body">
            <div class="activity-list">
              <div class="activity-item">
                <div class="activity-icon success">
                  <Check :size="14" />
                </div>
                <div class="activity-content">
                  <div class="activity-text">成功连接 WebShell</div>
                  <div class="activity-time">5 分钟前</div>
                </div>
              </div>
              <div class="activity-item">
                <div class="activity-icon info">
                  <Plus :size="14" />
                </div>
                <div class="activity-content">
                  <div class="activity-text">创建新项目 "测试项目"</div>
                  <div class="activity-time">10 分钟前</div>
                </div>
              </div>
              <div class="activity-item">
                <div class="activity-icon info">
                  <Edit :size="14" />
                </div>
                <div class="activity-content">
                  <div class="activity-text">更新 WebShell 配置</div>
                  <div class="activity-time">15 分钟前</div>
                </div>
              </div>
              <div class="activity-item">
                <div class="activity-icon warning">
                  <AlertTriangle :size="14" />
                </div>
                <div class="activity-content">
                  <div class="activity-text">插件版本需要更新</div>
                  <div class="activity-time">1 小时前</div>
                </div>
              </div>
              <div class="activity-item">
                <div class="activity-icon success">
                  <Download :size="14" />
                </div>
                <div class="activity-content">
                  <div class="activity-text">下载新的 Payload 模板</div>
                  <div class="activity-time">3 小时前</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 使用统计图表区域 -->
        <div class="dashboard-card usage-stats">
          <div class="card-header">
            <BarChart3 :size="20" />
            <h3>{{ t('home.usageStatistics') }}</h3>
          </div>
          <div class="card-body">
            <div class="stats-grid">
              <div class="stat-item">
                <div class="stat-label">{{ t('home.projectsByType') }}</div>
                <div class="stat-bars">
                  <div class="stat-bar">
                    <div class="stat-bar-label">Web</div>
                    <div class="stat-bar-fill" style="width: 80%"></div>
                    <span class="stat-bar-value">75%</span>
                  </div>
                  <div class="stat-bar">
                    <div class="stat-bar-label">Mobile</div>
                    <div class="stat-bar-fill" style="width: 60%"></div>
                    <span class="stat-bar-value">60%</span>
                  </div>
                  <div class="stat-bar">
                    <div class="stat-bar-label">Desktop</div>
                    <div class="stat-bar-fill" style="width: 45%"></div>
                    <span class="stat-bar-value">45%</span>
                  </div>
                  <div class="stat-bar">
                    <div class="stat-bar-label">API</div>
                    <div class="stat-bar-fill" style="width: 70%"></div>
                    <span class="stat-bar-value">70%</span>
                  </div>
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">{{ t('home.shellByLanguage') }}</div>
                <div class="stat-tags">
                  <span class="stat-tag">PHP: 45</span>
                  <span class="stat-tag">ASPX: 32</span>
                  <span class="stat-tag">JSP: 28</span>
                  <span class="stat-tag">Node.js: 24</span>
                  <span class="stat-tag">Python: 27</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import {
  Folder,
  Terminal,
  Puzzle,
  Box,
  Cpu,
  Zap,
  Clock,
  BarChart3,
  TrendingUp,
  TrendingDown,
  Minus,
  Plus,
  Check,
  Edit,
  AlertTriangle,
  Download
} from 'lucide-vue-next'
import { NProgress } from 'naive-ui'
import { useDashboard } from '@/composables/useDashboard'

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

// 使用 Dashboard Composable
const { stats, activities, refresh } = useDashboard()

const formatUptimeDisplay = (uptime: string): string => {
  const seconds = parseFloat(uptime)
  if (isNaN(seconds)) return 'N/A'
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  
  if (hours > 0) {
    return `${hours}h ${minutes}m ${secs}s`
  } else if (minutes > 0) {
    return `${minutes}m ${secs}s`
  } else {
    return `${secs}s`
  }
}

const getCPUBarHeight = (cpuUsage: string, index: number): number => {
  const cpuPercent = parseFloat(cpuUsage) || 0
  const baseHeight = 20
  const variance = Math.sin(index * 1.5 + cpuPercent / 10) * 20
  const height = baseHeight + (cpuPercent * 0.8) + variance
  return Math.min(Math.max(height, 10), 100)
}

const getMemoryPercentage = (memoryUsage: string): number => {
  try {
    const parts = memoryUsage.split('/')
    if (parts.length !== 2) return 0
    
    const used = parseFloat(parts[0].trim())
    const total = parseFloat(parts[1].trim())
    
    if (isNaN(used) || isNaN(total) || total === 0) return 0
    
    return Math.min((used / total) * 100, 100)
  } catch {
    return 0
  }
}

const navigateTo = (page: string) => {
  console.log('Navigate to:', page)
}
</script>

<style scoped>
/* ========================================
   全局容器样式
   ======================================== */
.dashboard-view {
  width: 100%;
  height: 100%;
}

.dashboard-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  background: transparent;
}

.dashboard-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

/* ========================================
   核心数据指标卡片样式
   ======================================== */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

/* 响应式断点 */
@media (max-width: 1400px) {
  .metrics-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .metrics-grid {
    grid-template-columns: 1fr;
  }
  
  .content-body {
    padding: 16px;
  }
}

.metric-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color-light);
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.metric-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, 
    var(--metric-gradient-start, var(--active-color)) 0%, 
    var(--metric-gradient-end, var(--active-color-hover)) 100%);
  opacity: 0;
  transition: opacity 0.4s ease;
}

.metric-card.projects::before {
  --metric-gradient-start: #667eea;
  --metric-gradient-end: #764ba2;
}

.metric-card.webshells::before {
  --metric-gradient-start: #f093fb;
  --metric-gradient-end: #f5576c;
}

.metric-card.payloads::before {
  --metric-gradient-start: #4facfe;
  --metric-gradient-end: #00f2fe;
}

.metric-card.plugins::before {
  --metric-gradient-start: #43e97b;
  --metric-gradient-end: #38f9d7;
}

.metric-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
  border-color: var(--active-color);
}

.metric-card:hover::before {
  opacity: 1;
}

.metric-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.metric-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: transform 0.3s ease;
}

.metric-card:hover .metric-icon {
  transform: scale(1.1) rotate(5deg);
}

.metric-icon.projects { 
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.metric-icon.webshells { 
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  box-shadow: 0 4px 12px rgba(240, 147, 251, 0.3);
}

.metric-icon.payloads { 
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.metric-icon.plugins { 
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  box-shadow: 0 4px 12px rgba(67, 233, 123, 0.3);
}

.metric-trend {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.metric-trend.positive { 
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.metric-trend.negative { 
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.metric-card-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.metric-value {
  font-size: 36px;
  font-weight: 800;
  color: var(--text-primary);
  line-height: 1.1;
  letter-spacing: -0.5px;
  transition: color 0.3s ease;
}

.metric-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  line-height: 1.5;
  letter-spacing: 0.2px;
  transition: color 0.3s ease;
}

.metric-progress {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: var(--border-color);
  border-radius: 4px;
  overflow: hidden;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.05);
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.progress-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.2) 50%,
    transparent 100%
  );
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.progress-fill.projects { background: linear-gradient(90deg, #667eea, #764ba2); }
.progress-fill.webshells { background: linear-gradient(90deg, #f093fb, #f5576c); }
.progress-fill.payloads { background: linear-gradient(90deg, #4facfe, #00f2fe); }
.progress-fill.plugins { background: linear-gradient(90deg, #43e97b, #38f9d7); }

.progress-text {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-tertiary);
  white-space: nowrap;
  min-width: 70px;
  text-align: right;
  transition: color 0.3s ease;
}

/* ========================================
   系统状态和快捷操作网格
   ======================================== */
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

@media (max-width: 1024px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

.dashboard-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color-light);
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.dashboard-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border-color: var(--border-color);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 18px 24px;
  background: var(--card-bg-hover);
  border-bottom: 1px solid var(--border-color-light);
  transition: all 0.3s ease;
}

.card-header svg {
  color: var(--active-color);
  transition: transform 0.3s ease;
}

.dashboard-card:hover .card-header svg {
  transform: rotate(15deg) scale(1.1);
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 0.3px;
  line-height: 1.4;
  transition: color 0.3s ease;
}

.card-body {
  padding: 24px;
  transition: background 0.3s ease;
}

/* ========================================
   系统状态样式
   ======================================== */
.status-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

@media (max-width: 640px) {
  .status-grid {
    grid-template-columns: 1fr;
  }
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.status-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  letter-spacing: 0.3px;
  text-transform: uppercase;
  line-height: 1.5;
  transition: color 0.3s ease;
}

.status-value {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 600;
  line-height: 1.4;
  transition: color 0.3s ease;
}

.status-text {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
  transition: color 0.3s ease;
}

.cpu-bars {
  display: flex;
  align-items: flex-end;
  gap: 5px;
  height: 36px;
}

.cpu-bar {
  width: 10px;
  background: linear-gradient(180deg, var(--active-color) 0%, var(--active-color-hover) 100%);
  border-radius: 3px;
  min-height: 4px;
  transition: height 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* ========================================
   快捷操作样式
   ======================================== */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

@media (max-width: 640px) {
  .actions-grid {
    grid-template-columns: 1fr;
  }
}

.action-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 16px 24px;
  background: var(--card-bg-hover);
  border: 1px solid var(--border-color-light);
  border-radius: 12px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.action-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, 
    var(--active-color) 0%, 
    var(--active-color-hover) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 0;
}

.action-button:hover::before {
  opacity: 1;
}

.action-button > * {
  position: relative;
  z-index: 1;
}

.action-button:hover {
  color: white;
  border-color: transparent;
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(var(--active-color-rgb, 59, 130, 246), 0.3);
}

.action-button:active {
  transform: translateY(-1px);
}

/* ========================================
   活动日志样式
   ======================================== */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 12px;
  border-radius: 10px;
  transition: all 0.3s ease;
}

.activity-item:hover {
  background: var(--card-bg-hover);
  transform: translateX(4px);
}

.activity-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.activity-item:hover .activity-icon {
  transform: scale(1.15) rotate(10deg);
}

.activity-icon.success { 
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.3);
}

.activity-icon.info { 
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

.activity-icon.warning { 
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.3);
}

.activity-icon.error { 
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
}

.activity-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.activity-text {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.5;
  letter-spacing: 0.2px;
  transition: color 0.3s ease;
}

.activity-time {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-tertiary);
  line-height: 1.4;
  transition: color 0.3s ease;
}

/* ========================================
   使用统计样式
   ======================================== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.stat-label {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: 0.3px;
  line-height: 1.4;
  transition: color 0.3s ease;
}

.stat-bars {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.stat-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-bar-label {
  width: 60px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  letter-spacing: 0.2px;
  transition: color 0.3s ease;
}

.stat-bar-fill {
  flex: 1;
  height: 10px;
  background: linear-gradient(90deg, var(--active-color) 0%, var(--active-color-hover) 100%);
  border-radius: 5px;
  max-width: 140px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.stat-bar-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.3) 50%,
    transparent 100%
  );
  animation: shimmer 2s infinite;
}

.stat-bar-value {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-secondary);
  width: 45px;
  text-align: right;
  transition: color 0.3s ease;
}

.stat-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.stat-tag {
  padding: 8px 16px;
  background: var(--card-bg-hover);
  border: 1px solid var(--border-color-light);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  transition: all 0.3s ease;
  cursor: default;
}

.stat-tag:hover {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(var(--active-color-rgb, 59, 130, 246), 0.3);
}

/* ========================================
   深色主题适配
   ======================================== */
@media (prefers-color-scheme: dark) {
  .metric-card {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
  
  .metric-card:hover {
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.3);
  }
  
  .dashboard-card {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
  
  .dashboard-card:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
  }
  
  .cpu-bar {
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
  
  .activity-icon {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }
  
  .stat-bar-fill {
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
}

/* ========================================
   打印样式
   ======================================== */
@media print {
  .metric-card,
  .dashboard-card {
    break-inside: avoid;
  }
  
  .action-button,
  .metric-trend {
    display: none;
  }
}
</style>
