<template>
  <div class="audit-log-view">
    <NCard title="审计日志" class="audit-log-card">
      <div class="audit-log-header">
        <NText depth="2">系统操作审计记录</NText>
        <NButton type="error" @click="handleClearLogs" :disabled="loading">
          <template #icon>
            <NIcon><Trash2 /></NIcon>
          </template>
          清空日志
        </NButton>
      </div>

      <NDataTable
        ref="dataTable"
        :data="logs"
        :loading="loading"
        :columns="columns"
        :pagination="pagination"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
        class="audit-log-table"
      >
        <template #empty>
          <div class="empty-state">
            <NIcon size="48" color="var(--text-sub)">
              <FileText />
            </NIcon>
            <div class="empty-text">暂无审计日志</div>
          </div>
        </template>
      </NDataTable>
    </NCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { invoke } from '@/utils/tauri-mock-adapter';
import { NIcon, NCard, NDataTable, NButton, NText, useMessage } from 'naive-ui';
import { Trash2, FileText } from 'lucide-vue-next';

const message = useMessage();

interface AuditLog {
  id: number;
  action_type: string;
  payload_hash?: string;
  content: string;
  created_at: string;
}

const loading = ref(false);
const logs = ref<AuditLog[]>([]);
const total = ref(0);

const pagination = reactive({
  page: 1,
  pageSize: 20,
  pageSizes: [10, 20, 50, 100],
  showSizePicker: true,
  showQuickJumper: true,
  showTotal: true,
  total: 0
});

const columns = [
  {
    title: '操作类型',
    key: 'action_type',
    width: 150,
    render(row: AuditLog) {
      const typeMap: Record<string, string> = {
        'PayloadGenerate': '生成载荷',
        'WebShellConnect': '连接WebShell',
        'WebShellCommand': '执行命令',
        'PluginLoad': '加载插件',
        'ProjectCreate': '创建项目',
        'ProjectDelete': '删除项目',
        'WebShellCreate': '创建WebShell',
        'WebShellDelete': '删除WebShell',
        'SettingsChange': '修改设置',
        'Other': '其他操作'
      };
      return typeMap[row.action_type] || row.action_type;
    }
  },
  {
    title: '操作内容',
    key: 'content',
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '载荷哈希',
    key: 'payload_hash',
    width: 200,
    ellipsis: {
      tooltip: true
    },
    render(row: AuditLog) {
      return row.payload_hash || '-';
    }
  },
  {
    title: '操作时间',
    key: 'created_at',
    width: 200,
    render(row: AuditLog) {
      return new Date(row.created_at).toLocaleString();
    }
  }
];

const fetchLogs = async () => {
  loading.value = true;
  try {
    const response = await invoke('get_audit_logs', {
      limit: pagination.pageSize,
      offset: (pagination.page - 1) * pagination.pageSize
    });
    
    if (response && typeof response === 'object') {
      logs.value = response.logs || [];
      total.value = response.total || 0;
      pagination.total = total.value;
    }
  } catch (error: any) {
    message.error('获取审计日志失败: ' + (error.message || '未知错误'));
  } finally {
    loading.value = false;
  }
};

const handlePageChange = (page: number) => {
  pagination.page = page;
  fetchLogs();
};

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  fetchLogs();
};

const handleClearLogs = async () => {
  try {
    await invoke('clear_audit_logs');
    message.success('审计日志已清空');
    fetchLogs();
  } catch (error: any) {
    message.error('清空审计日志失败: ' + (error.message || '未知错误'));
  }
};

onMounted(() => {
  fetchLogs();
});
</script>

<style scoped>
.audit-log-view {
  padding: 20px;
  height: 100%;
  overflow: auto;
}

.audit-log-card {
  height: 100%;
  border-radius: 10px;
}

.audit-log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
}

.audit-log-table {
  height: calc(100% - 100px);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: var(--text-sub);
}

.empty-text {
  margin-top: 16px;
  font-size: 14px;
}
</style>
