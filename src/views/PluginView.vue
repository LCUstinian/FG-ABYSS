<template>
  <div class="plugin-view">
    <div class="plugin-toolbar">
      <span class="plugin-count">已安装 ({{ filteredPlugins.length }})</span>
      <div class="filter-tabs">
        <button
          v-for="tab in filterTabs"
          :key="tab.value"
          class="filter-tab"
          :class="{ 'is-active': activeFilter === tab.value }"
          @click="activeFilter = tab.value"
        >{{ tab.label }}</button>
      </div>
      <div style="flex:1" />
      <n-button type="primary" size="small">
        <template #icon><Plus :size="14" /></template>
        安装
      </n-button>
    </div>

    <div class="plugin-grid" v-if="filteredPlugins.length > 0">
      <div
        v-for="plugin in filteredPlugins"
        :key="plugin.id"
        class="plugin-card"
        :class="{ 'is-disabled': !plugin.enabled }"
      >
        <div class="card-header">
          <span class="plugin-name">{{ plugin.name }}</span>
          <div class="card-badges">
            <span class="version-badge">{{ plugin.version }}</span>
            <n-tag v-if="plugin.builtin" size="tiny" :color="{ color: 'var(--accent-bg)', textColor: 'var(--accent)', borderColor: 'transparent' }">内置</n-tag>
            <n-tag v-else size="tiny">第三方</n-tag>
          </div>
        </div>
        <p class="plugin-desc">{{ plugin.description }}</p>
        <div class="card-footer">
          <div style="display:flex;align-items:center;gap:8px">
            <n-switch
              :value="plugin.enabled"
              size="small"
              @update:value="(v: boolean) => togglePlugin(plugin.id, v)"
            />
            <span class="switch-label">{{ plugin.enabled ? '已启用' : '已禁用' }}</span>
          </div>
          <n-button v-if="!plugin.builtin" text size="tiny" type="error">卸载</n-button>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <Puzzle :size="64" class="empty-icon" />
      <div class="empty-title">还没有安装插件</div>
      <div class="empty-sub">内置插件随应用分发</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Plus, Puzzle } from 'lucide-vue-next'

interface Plugin {
  id: string; name: string; version: string
  builtin: boolean; description: string; enabled: boolean
}

const plugins = ref<Plugin[]>([
  { id: 'file-manager', name: '文件管理', version: 'v1.2.0', builtin: true,  description: '远程文件系统浏览、上传、下载、编辑', enabled: true  },
  { id: 'db-manager',   name: '数据库管理', version: 'v1.0.0', builtin: true,  description: 'SQL 查询执行与结果导出', enabled: true  },
  { id: 'terminal',     name: '终端',   version: 'v1.1.0', builtin: true,  description: '远程命令执行与进程管理', enabled: true  },
])

const activeFilter = ref('all')

const filterTabs = [
  { label: '全部',   value: 'all'      },
  { label: '已启用', value: 'enabled'  },
  { label: '已禁用', value: 'disabled' },
  { label: '内置',   value: 'builtin'  },
  { label: '第三方', value: 'third'    },
]

const filteredPlugins = computed(() => {
  switch (activeFilter.value) {
    case 'enabled':  return plugins.value.filter(p => p.enabled)
    case 'disabled': return plugins.value.filter(p => !p.enabled)
    case 'builtin':  return plugins.value.filter(p => p.builtin)
    case 'third':    return plugins.value.filter(p => !p.builtin)
    default:         return plugins.value
  }
})

function togglePlugin(id: string, enabled: boolean) {
  const p = plugins.value.find(x => x.id === id)
  if (p) p.enabled = enabled
}
</script>

<style scoped>
.plugin-view {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 16px;
  overflow-y: auto;
}

.plugin-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.plugin-count {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
}

.filter-tabs {
  display: flex;
  gap: 4px;
}

.filter-tab {
  padding: 4px 12px;
  border-radius: 4px;
  border: none;
  background: transparent;
  color: var(--text-2);
  font-size: 13px;
  cursor: pointer;
  transition: background 80ms;
}

.filter-tab:hover { background: var(--bg-hover); color: var(--text-1); }
.filter-tab.is-active { background: var(--accent-bg); color: var(--accent); }

.plugin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.plugin-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: opacity 200ms;
}

[data-theme="light"] .plugin-card {
  border: none;
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}

.plugin-card.is-disabled { opacity: 0.6; }

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.plugin-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
}

.card-badges {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.version-badge {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}

.plugin-desc {
  font-size: 13px;
  color: var(--text-2);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
}

.switch-label {
  font-size: 12px;
  color: var(--text-2);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 80px 0;
  color: var(--text-3);
}

.empty-icon { color: var(--text-3); }
.empty-title { font-size: 14px; color: var(--text-1); }
.empty-sub   { font-size: 13px; color: var(--text-2); }
</style>
