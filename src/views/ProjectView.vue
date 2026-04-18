<template>
  <div class="project-view">
    <!-- Project tree -->
    <div class="project-tree">
      <div class="tree-header">
        <n-button size="small" type="primary" ghost>
          <template #icon><Plus :size="14" /></template>
          新建项目
        </n-button>
      </div>
      <div class="tree-list">
        <div
          v-for="proj in projects"
          :key="proj.id"
          class="tree-item"
          :class="{ 'is-active': selectedProject === proj.id }"
          @click="selectedProject = proj.id"
        >
          <span class="tree-name">{{ proj.name }}</span>
          <span class="tree-count">{{ proj.shellCount }}</span>
        </div>
        <div
          class="tree-item tree-item--recycle"
          :class="{ 'is-active': selectedProject === '__recycle__' }"
          @click="selectedProject = '__recycle__'"
        >
          <Trash2 :size="14" />
          <span class="tree-name">回收站</span>
          <span class="tree-count">0</span>
        </div>
      </div>
    </div>

    <!-- Shell table -->
    <div class="shell-table-area">
      <div class="table-toolbar">
        <n-input
          v-model:value="searchText"
          placeholder="搜索..."
          clearable
          size="small"
          style="width: 220px"
        >
          <template #prefix><Search :size="14" /></template>
        </n-input>
        <n-select
          v-model:value="typeFilter"
          :options="typeOptions"
          placeholder="类型"
          clearable
          size="small"
          style="width: 100px"
        />
        <div style="flex:1" />
        <n-button type="primary" size="small">
          <template #icon><Plus :size="14" /></template>
          添加 Shell
        </n-button>
      </div>

      <n-data-table
        :columns="columns"
        :data="filteredShells"
        :row-key="(row: Shell) => row.id"
        :loading="false"
        size="small"
        class="shell-table"
      />

      <div class="table-footer">
        <n-pagination :page="1" :page-count="1" size="small" />
        <span class="total-count">共 {{ filteredShells.length }} 条</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns } from 'naive-ui'
import { Plus, Search, Trash2 } from 'lucide-vue-next'

interface Project { id: string; name: string; shellCount: number }
interface Shell {
  id: string; name: string; url: string
  type: 'PHP' | 'JSP' | 'ASP' | 'ASPX'
  status: 'active' | 'offline' | 'error'
  lastConnected: string
  projectId: string
}

const projects = ref<Project[]>([
  { id: 'all', name: '全部', shellCount: 0 },
])

const selectedProject = ref('all')
const searchText = ref('')
const typeFilter = ref<string | null>(null)

const typeOptions = [
  { label: 'PHP',  value: 'PHP'  },
  { label: 'JSP',  value: 'JSP'  },
  { label: 'ASP',  value: 'ASP'  },
  { label: 'ASPX', value: 'ASPX' },
]

const shells = ref<Shell[]>([])

const filteredShells = computed(() => {
  let list = shells.value
  if (typeFilter.value) list = list.filter(s => s.type === typeFilter.value)
  if (searchText.value) {
    const q = searchText.value.toLowerCase()
    list = list.filter(s => s.name.toLowerCase().includes(q) || s.url.toLowerCase().includes(q))
  }
  return list
})

const STATUS_DOT: Record<string, string> = {
  active:  'var(--color-success)',
  offline: 'var(--text-3)',
  error:   'var(--color-error)',
}

const TYPE_COLORS: Record<string, { bg: string; text: string }> = {
  PHP:  { bg: 'rgba(79,156,255,0.15)',  text: '#4f9cff' },
  JSP:  { bg: 'rgba(251,146,60,0.15)',  text: '#fb923c' },
  ASP:  { bg: 'rgba(167,139,250,0.15)', text: '#a78bfa' },
  ASPX: { bg: 'rgba(34,211,238,0.15)',  text: '#22d3ee' },
}

const columns: DataTableColumns<Shell> = [
  {
    type: 'selection',
    width: 40,
  },
  {
    key: 'status',
    title: '状',
    width: 40,
    render: (row) => h('span', {
      style: {
        display: 'inline-block',
        width: '8px', height: '8px',
        borderRadius: '50%',
        background: STATUS_DOT[row.status],
      },
      'aria-label': `状态: ${row.status}`,
    }),
  },
  {
    key: 'name',
    title: '名称',
    width: 160,
    render: (row) => h('span', { style: { color: 'var(--text-1)' } }, row.name),
  },
  {
    key: 'url',
    title: 'URL',
    render: (row) => h('span', {
      style: { fontFamily: 'var(--font-mono)', fontSize: '12px', color: 'var(--text-2)' },
      title: row.url,
    }, row.url),
  },
  {
    key: 'type',
    title: '类型',
    width: 80,
    render: (row) => {
      const c = TYPE_COLORS[row.type] ?? TYPE_COLORS.PHP
      return h('span', {
        style: {
          padding: '2px 6px',
          borderRadius: '4px',
          background: c.bg,
          color: c.text,
          fontSize: '12px',
          fontFamily: 'var(--font-mono)',
        },
      }, row.type)
    },
  },
  {
    key: 'lastConnected',
    title: '最近连接',
    width: 120,
    render: (row) => h('span', { style: { color: 'var(--text-2)', fontSize: '12px' } }, row.lastConnected),
  },
]
</script>

<style scoped>
.project-view {
  height: 100%;
  display: flex;
  overflow: hidden;
}

.project-tree {
  width: 240px;
  flex-shrink: 0;
  background: var(--bg-deep);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.tree-header {
  padding: 12px;
  border-bottom: 1px solid var(--border);
}

.tree-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.tree-item {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  padding: 0 16px;
  cursor: pointer;
  color: var(--text-2);
  font-size: 13px;
  transition: background 80ms ease-out;
}

.tree-item:hover { background: var(--bg-hover); color: var(--text-1); }
.tree-item.is-active { background: var(--accent-bg); color: var(--accent); font-weight: 600; }

.tree-name { flex: 1; }
.tree-count { font-size: 12px; color: var(--text-3); font-family: var(--font-mono); }

.tree-item--recycle {
  border-top: 1px solid var(--border);
  color: var(--text-3);
  margin-top: 8px;
}

.shell-table-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.table-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
}

.shell-table {
  flex: 1;
  overflow: hidden;
}

.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-top: 1px solid var(--border);
}

.total-count {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}
</style>
