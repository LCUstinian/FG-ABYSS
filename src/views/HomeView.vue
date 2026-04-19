<template>
  <div class="home-view">
    <div class="home-header">
      <h1 class="home-title">欢迎回来</h1>
      <span class="home-date">{{ todayStr }}</span>
    </div>

    <!-- Stat cards -->
    <div class="stat-cards">
      <div class="stat-card" v-for="card in statCards" :key="card.label">
        <component :is="card.icon" :size="20" class="stat-icon" />
        <div class="stat-body">
          <div class="stat-num">{{ card.value }}</div>
          <div class="stat-label">{{ card.label }}</div>
        </div>
      </div>
    </div>

    <!-- Recent connections -->
    <div class="recent-section">
      <h2 class="section-title">最近连接</h2>
      <div v-if="recentShells.length === 0" class="empty-state">
        <Globe :size="64" class="empty-icon" />
        <div class="empty-title">还没有 WebShell</div>
        <div class="empty-sub">前往「项目」页面添加第一个 WebShell</div>
        <RouterLink to="/project">
          <n-button type="primary" size="medium">前往项目</n-button>
        </RouterLink>
      </div>
      <div v-else class="recent-list">
        <div
          v-for="shell in recentShells"
          :key="shell.id"
          class="recent-item"
        >
          <span
            class="recent-dot"
            :class="shell.active ? 'dot-active' : 'dot-offline'"
            :aria-label="`状态: ${shell.active ? '活跃' : '离线'}`"
          />
          <span class="recent-name">{{ shell.name }}</span>
          <span class="recent-url">{{ shell.url }}</span>
          <n-tag size="small" :color="typeColor(shell.type)">{{ shell.type }}</n-tag>
          <span class="recent-time">{{ shell.relativeTime }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Terminal, FolderOpen, Package, Zap, Globe } from 'lucide-vue-next'

const todayStr = new Date().toLocaleDateString('zh-CN', {
  year: 'numeric', month: '2-digit', day: '2-digit',
}).replace(/\//g, '/')

const statCards = [
  { icon: Terminal,   value: 0, label: '活跃 Shell' },
  { icon: FolderOpen, value: 0, label: '项目'       },
  { icon: Package,    value: 0, label: '载荷'       },
  { icon: Zap,        value: 0, label: '今日请求'   },
]

const recentShells: Array<{
  id: string; name: string; url: string
  type: string; active: boolean; relativeTime: string
}> = []

const TYPE_COLORS: Record<string, { color: string; textColor: string; borderColor: string }> = {
  PHP:  { color: 'rgba(79,156,255,0.15)',  textColor: '#4f9cff', borderColor: 'transparent' },
  JSP:  { color: 'rgba(251,146,60,0.15)',  textColor: '#fb923c', borderColor: 'transparent' },
  ASP:  { color: 'rgba(167,139,250,0.15)', textColor: '#a78bfa', borderColor: 'transparent' },
  ASPX: { color: 'rgba(34,211,238,0.15)',  textColor: '#22d3ee', borderColor: 'transparent' },
}

function typeColor(type: string) {
  return TYPE_COLORS[type] ?? TYPE_COLORS.PHP
}
</script>

<style scoped>
.home-view {
  height: 100%;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.home-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.home-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-1);
}

.home-date {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}

.stat-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 16px;
}

.stat-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-left: 3px solid var(--accent);
}

[data-theme="light"] .stat-card {
  border: none;
  border-left: 3px solid var(--accent);
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}

.stat-icon {
  color: var(--accent);
  flex-shrink: 0;
}

.stat-num {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-1);
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
}

.stat-label {
  font-size: 12px;
  color: var(--text-2);
  margin-top: 2px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
  margin-bottom: 12px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 48px 0;
  color: var(--text-3);
}

.empty-icon { color: var(--text-3); }

.empty-title {
  font-size: 14px;
  color: var(--text-1);
}

.empty-sub {
  font-size: 13px;
  color: var(--text-2);
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 40px;
  padding: 0 12px;
  border-radius: 6px;
  background: var(--bg-elevated);
}

.recent-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.dot-active  { background: var(--color-success); }
.dot-offline { background: var(--text-3); }

.recent-name {
  font-size: 13px;
  color: var(--text-1);
  min-width: 120px;
}

.recent-url {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-2);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.recent-time {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}
</style>
