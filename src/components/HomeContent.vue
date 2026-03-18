<template>
  <div class="content-section">
    <PageHeader :title="t('home.title')" :subtitle="t('home.subtitle')" />
    <div class="content-body">
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

          <div class="metric-card" @click="navigateTo('plugins')">
            <div class="metric-card-header">
              <div class="metric-icon plugins">
                <Puzzle :size="24" />
              </div>
              <div class="metric-trend neutral">
                <Minus :size="16" />
                <span>0%</span>
              </div>
            </div>
            <div class="metric-card-body">
              <div class="metric-value">12</div>
              <div class="metric-label">{{ t('home.totalPlugins') }}</div>
              <div class="metric-progress">
                <div class="progress-bar">
                  <div class="progress-fill plugins" style="width: 40%"></div>
                </div>
                <span class="progress-text">40% 启用</span>
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
                  <div class="progress-fill payloads" style="width: 85%"></div>
                </div>
                <span class="progress-text">85% 容量</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 系统状态和快速操作 -->
        <div class="dashboard-grid">
          <!-- 系统状态卡片 -->
          <div class="dashboard-card system-status">
            <div class="card-header">
              <Cpu :size="20" />
              <h3>{{ t('home.systemStatus') }}</h3>
            </div>
            <div class="card-body">
              <div class="status-item">
                <div class="status-label">
                  <span class="status-dot"></span>
                  {{ t('home.cpuUsage') }}
                </div>
                <div class="status-value">{{ systemStatus.cpuUsage }}</div>
                <div class="mini-chart">
                  <div 
                    v-for="i in 5" 
                    :key="i" 
                    class="chart-bar" 
                    :style="{ height: getCPUBarHeight(systemStatus.cpuUsage, i) + '%' }"
                  ></div>
                </div>
              </div>
              <div class="status-item">
                <div class="status-label">
                  <span class="status-dot"></span>
                  {{ t('home.memoryUsage') }}
                </div>
                <div class="status-value">{{ systemStatus.memoryUsage }}</div>
                <div class="progress-bar">
                  <div class="progress-fill memory" :style="{ width: getMemoryPercentage(systemStatus.memoryUsage) + '%' }"></div>
                </div>
              </div>
              <div class="status-item">
                <div class="status-label">
                  <span class="status-dot active"></span>
                  {{ t('home.uptime') }}
                </div>
                <div class="status-value">{{ formatUptimeDisplay(systemStatus.uptime) }}</div>
              </div>
              <div class="status-item">
                <div class="status-label">
                  <span class="status-dot"></span>
                  {{ t('home.processId') }}
                </div>
                <div class="status-value">{{ systemStatus.processId || 'N/A' }}</div>
              </div>
            </div>
          </div>

          <!-- 快速操作卡片 -->
          <div class="dashboard-card quick-actions">
            <div class="card-header">
              <Zap :size="20" />
              <h3>{{ t('home.quickActions') }}</h3>
            </div>
            <div class="card-body">
              <button class="action-btn" @click="navigateTo('projects')">
                <div class="action-btn-icon">
                  <Plus :size="16" />
                </div>
                <span class="action-btn-text">{{ t('projects.newProject') }}</span>
              </button>
              <button class="action-btn" @click="navigateTo('webshells')">
                <div class="action-btn-icon">
                  <Plus :size="16" />
                </div>
                <span class="action-btn-text">{{ t('projects.newWebShell') }}</span>
              </button>
              <button class="action-btn" @click="navigateTo('plugins')">
                <div class="action-btn-icon">
                  <Plus :size="16" />
                </div>
                <span class="action-btn-text">{{ t('plugins.newPlugin') }}</span>
              </button>
              <button class="action-btn" @click="navigateTo('payloads')">
                <div class="action-btn-icon">
                  <Plus :size="16" />
                </div>
                <span class="action-btn-text">{{ t('payloads.generate') }}</span>
              </button>
            </div>
          </div>

          <!-- 最近活动卡片 -->
          <div class="dashboard-card recent-activity">
            <div class="card-header">
              <Clock :size="20" />
              <h3>{{ t('home.recentActivity') }}</h3>
            </div>
            <div class="card-body">
              <div class="activity-item">
                <div class="activity-icon success">
                  <Check :size="14" />
                </div>
                <div class="activity-content">
                  <div class="activity-text">创建新项目 "Project Alpha"</div>
                  <div class="activity-time">2 分钟前</div>
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
                    <span class="stat-bar-value">80%</span>
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
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/shared/PageHeader.vue'
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

// 格式化运行时间显示
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

// 计算 CPU 柱状图高度（根据 CPU 使用率和柱子位置）
const getCPUBarHeight = (cpuUsage: string, index: number): number => {
  const cpuPercent = parseFloat(cpuUsage) || 0
  // 基础高度 20%，根据 CPU 使用率波动
  const baseHeight = 20
  // 波动范围：每个柱子有不同的波动系数
  const variance = Math.sin(index * 1.5 + cpuPercent / 10) * 20
  // 最终高度 = 基础高度 + (CPU 使用率 * 系数) + 波动
  const height = baseHeight + (cpuPercent * 0.8) + variance
  return Math.min(Math.max(height, 10), 100) // 限制在 10-100% 之间
}

// 计算内存使用百分比
const getMemoryPercentage = (memoryUsage: string): number => {
  try {
    // 解析 "11.75 GB / 32 GB" 格式
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
  // TODO: 实现页面导航逻辑
  console.log('Navigate to:', page)
}
</script>

<style scoped>
/* 内容区域布局 */
.content-section {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}



.content-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.dashboard-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

/* 核心数据指标网格 */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 20px;
  width: 100%;
  flex-shrink: 0;
}

.metric-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: visible;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 180px;
}

.dark .metric-card {
  border-color: var(--border-strong);
  background: linear-gradient(135deg, var(--card-bg), rgba(255, 255, 255, 0.02));
}

.metric-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.2);
  border-color: var(--active-color);
}

.dark .metric-card:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05);
}

.metric-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.metric-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
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
}

.metric-icon.webshells {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.metric-icon.plugins {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.metric-icon.payloads {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.metric-trend {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 700;
  padding: 6px 10px;
  border-radius: 8px;
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.metric-trend.positive {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  box-shadow: 0 2px 8px rgba(34, 197, 94, 0.2);
}

.metric-trend.neutral {
  background: rgba(107, 114, 128, 0.15);
  color: #6b7280;
  box-shadow: 0 2px 8px rgba(107, 114, 128, 0.2);
}

.metric-trend.negative {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.2);
}

.metric-card-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}

.metric-value {
  font-size: 42px;
  font-weight: 800;
  color: var(--text-color);
  line-height: 1;
  letter-spacing: -1px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.dark .metric-value {
  text-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);
}

.metric-label {
  font-size: 15px;
  color: var(--text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.metric-progress {
  margin-top: 12px;
}

.progress-bar {
  height: 8px;
  background: var(--bg-secondary);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 6px;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
}

.dark .progress-bar {
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.3);
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
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
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  animation: shimmer 3s infinite;
}

.progress-fill.projects {
  background: linear-gradient(90deg, #667eea, #764ba2);
}

.progress-fill.webshells {
  background: linear-gradient(90deg, #f093fb, #f5576c);
}

.progress-fill.plugins {
  background: linear-gradient(90deg, #4facfe, #00f2fe);
}

.progress-fill.payloads {
  background: linear-gradient(90deg, #43e97b, #38f9d7);
}

.progress-fill.memory {
  background: linear-gradient(90deg, #667eea, #764ba2);
}

.progress-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.progress-text::before {
  content: '';
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--active-color);
  box-shadow: 0 0 8px var(--active-color);
}

/* 仪表盘网格 */
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
  flex-shrink: 0;
}

.dashboard-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: visible;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 1;
}

.dark .dashboard-card {
  border-color: var(--border-strong);
  background: linear-gradient(135deg, var(--card-bg), rgba(255, 255, 255, 0.02));
}

.dashboard-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: var(--active-color);
  z-index: 2;
}

.dark .dashboard-card:hover {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.05);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-color);
  flex-shrink: 0;
}

.dark .card-header {
  border-bottom-color: var(--border-strong);
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.3px;
}

.card-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  overflow: visible;
}

/* 系统状态 */
.status-item {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 24px;
  position: relative;
}

.status-item:last-child {
  margin-bottom: 0;
}

.status-label {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--text-tertiary);
  box-shadow: 0 0 8px rgba(0, 0, 0, 0.2);
  flex-shrink: 0;
}

.status-dot.active {
  background: #22c55e;
  box-shadow: 0 0 12px rgba(34, 197, 94, 0.6);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 0 12px rgba(34, 197, 94, 0.6);
  }
  50% {
    box-shadow: 0 0 20px rgba(34, 197, 94, 0.8);
  }
}

.status-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-color);
  letter-spacing: -0.5px;
}

.mini-chart {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  height: 50px;
  margin-top: 12px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}

.dark .mini-chart {
  background: rgba(255, 255, 255, 0.05);
}

.chart-bar {
  flex: 1;
  background: linear-gradient(180deg, var(--active-color), rgba(102, 126, 234, 0.4));
  border-radius: 4px 4px 0 0;
  min-width: 10px;
  transition: height 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.dark .chart-bar {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);
}

/* 快速操作卡片特殊优化 */
.quick-actions {
  display: flex;
  flex-direction: column;
}

.quick-actions .card-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  overflow: visible;
}

/* 快速操作 */
.action-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 11px 14px;
  margin: 0;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-color);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  overflow: visible;
  flex-shrink: 0;
}

.dark .action-btn {
  border-color: var(--border-strong);
  background: rgba(255, 255, 255, 0.03);
}

.action-btn:hover {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.dark .action-btn:hover {
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
}

.action-btn:active {
  transform: translateX(2px) scale(0.98);
}

.action-btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 6px;
  background: rgba(59, 130, 246, 0.1);
  color: var(--active-color);
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.action-btn:hover .action-btn-icon {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.dark .action-btn-icon {
  background: rgba(59, 130, 246, 0.15);
}

.action-btn-text {
  flex: 1;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

/* 最近活动 */
.activity-item {
  display: flex;
  gap: 14px;
  padding: 14px 0;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s ease;
  position: relative;
}

.dark .activity-item {
  border-bottom-color: var(--border-strong);
}

.activity-item:hover {
  background: rgba(0, 0, 0, 0.02);
  padding-left: 8px;
  padding-right: 8px;
  margin: 0 -8px;
  border-radius: 8px;
}

.dark .activity-item:hover {
  background: rgba(255, 255, 255, 0.02);
}

.activity-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.activity-item:first-child {
  padding-top: 0;
}

.activity-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.activity-item:hover .activity-icon {
  transform: scale(1.1) rotate(5deg);
}

.dark .activity-icon {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
}

.activity-icon.success {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 0.1));
  color: #22c55e;
  box-shadow: 0 2px 8px rgba(34, 197, 94, 0.3);
}

.activity-icon.info {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(59, 130, 246, 0.1));
  color: #3b82f6;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

.activity-icon.warning {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(245, 158, 11, 0.1));
  color: #f59e0b;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.3);
}

.activity-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.activity-text {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-color);
  line-height: 1.4;
  word-wrap: break-word;
}

.activity-time {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-tertiary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.activity-time::before {
  content: '•';
  color: var(--active-color);
  font-size: 16px;
  line-height: 1;
}

/* 使用统计 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 28px;
  width: 100%;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.stat-label {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-color);
  letter-spacing: -0.3px;
}

.stat-bars {
  display: flex;
  flex-direction: column;
  gap: 14px;
  width: 100%;
}

.stat-bar {
  display: grid;
  grid-template-columns: 70px 1fr auto;
  align-items: center;
  gap: 14px;
  position: relative;
}

.stat-bar-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  white-space: nowrap;
}

.stat-bar-fill {
  height: 10px;
  background: linear-gradient(90deg, var(--active-color), var(--active-color-suppl));
  border-radius: 6px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .stat-bar-fill {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);
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
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.stat-bar-value {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-color);
  min-width: 50px;
  text-align: right;
}

.stat-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  width: 100%;
}

.stat-tag {
  padding: 8px 14px;
  background: linear-gradient(135deg, var(--bg-secondary), rgba(255, 255, 255, 0.02));
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  transition: all 0.3s ease;
  cursor: default;
}

.dark .stat-tag {
  border-color: var(--border-strong);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.03), rgba(255, 255, 255, 0.01));
}

.stat-tag:hover {
  transform: translateY(-2px);
  border-color: var(--active-color);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.dark .stat-tag:hover {
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.05);
}

/* 响应式设计 */
@media (max-width: 1280px) {
  .metrics-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  
  .quick-actions .card-body {
    gap: 8px;
  }
  
  .action-btn {
    padding: 10px 12px;
    font-size: 13px;
  }
  
  .action-btn-icon {
    width: 20px;
    height: 20px;
  }
  
  .action-btn-icon svg {
    width: 14px;
    height: 14px;
  }
}

@media (max-width: 1024px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .content-header h1 {
    font-size: 24px;
  }
  
  .content-header .title {
    font-size: 28px;
  }
  
  .content-header .subtitle {
    font-size: 16px;
  }
  
  .quick-actions .card-body {
    gap: 8px;
    padding: 14px;
  }
  
  .action-btn {
    padding: 10px 12px;
  }
}

@media (max-width: 768px) {
  .metrics-grid {
    grid-template-columns: 1fr;
  }
  
  .content-body {
    padding: 16px;
  }
  
  .dashboard-container {
    gap: 16px;
  }
  
  .metric-card {
    padding: 20px;
  }
  
  .metric-icon {
    width: 48px;
    height: 48px;
  }
  
  .metric-value {
    font-size: 32px;
  }
  
  .metric-label {
    font-size: 13px;
  }
  
  .dashboard-card {
    border-radius: 12px;
  }
  
  .card-header {
    padding: 14px 18px;
  }
  
  .card-header h3 {
    font-size: 15px;
  }
  
  .card-body {
    padding: 16px;
  }
  
  .quick-actions .card-body {
    gap: 8px;
    padding: 14px;
  }
  
  .action-btn {
    padding: 10px 12px;
    font-size: 13px;
    gap: 10px;
  }
  
  .action-btn-icon {
    width: 20px;
    height: 20px;
  }
  
  .action-btn-icon svg {
    width: 14px;
    height: 14px;
  }
  
  .stats-grid {
    gap: 20px;
  }
  
  .stat-bar {
    grid-template-columns: 60px 1fr auto;
  }
  
  .stat-tag {
    padding: 6px 12px;
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  
  .metric-card {
    padding: 18px;
  }
  
  .metric-value {
    font-size: 28px;
  }
  
  .status-value {
    font-size: 20px;
  }
  
  .mini-chart {
    height: 40px;
  }
  
  .action-btn {
    padding: 9px 11px;
    font-size: 12px;
  }
  
  .action-btn-icon {
    width: 18px;
    height: 18px;
  }
  
  .action-btn-icon svg {
    width: 12px;
    height: 12px;
  }
  
  .quick-actions .card-body {
    gap: 6px;
    padding: 12px;
  }
}

/* 滚动条美化 */
.content-body::-webkit-scrollbar {
  width: 8px;
}

.content-body::-webkit-scrollbar-track {
  background: var(--bg-secondary);
  border-radius: 4px;
}

.content-body::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.content-body::-webkit-scrollbar-thumb:hover {
  background: var(--active-color);
}

.dark .content-body::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

.dark .content-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}

.dark .content-body::-webkit-scrollbar-thumb:hover {
  background: var(--active-color);
}

/* 卡片内部滚动条美化 */
.card-body::-webkit-scrollbar {
  width: 6px;
}

.card-body::-webkit-scrollbar-track {
  background: transparent;
}

.card-body::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
}

.dark .card-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}
</style>
