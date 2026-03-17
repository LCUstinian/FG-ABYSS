<template>
  <div class="plugin-manager">
    <n-card title="插件管理" size="small">
      <template #header-extra>
        <n-space>
          <n-button type="primary" @click="handleInstall">
            <template #icon>
              <i class="i-carbon-add"></i>
            </template>
            安装插件
          </n-button>
          <n-button @click="refreshList">
            <template #icon>
              <i class="i-carbon-refresh"></i>
            </template>
            刷新
          </n-button>
        </n-space>
      </template>

      <!-- 统计信息 -->
      <n-alert type="info" title="插件统计" style="margin-bottom: 16px;">
        <n-space>
          <n-tag type="info">总插件数：{{ stats.total }}</n-tag>
          <n-tag type="success">已启用：{{ stats.enabled }}</n-tag>
          <n-tag type="warning">内置插件：{{ stats.builtin }}</n-tag>
          <n-tag type="error">外置插件：{{ stats.external }}</n-tag>
        </n-space>
      </n-alert>

      <!-- 插件列表 -->
      <n-data-table
        :columns="columns"
        :data="pluginList"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row) => row.id"
      />
    </n-card>

    <!-- 安装插件对话框 -->
    <n-modal v-model:show="showInstallModal" preset="dialog" title="安装插件" :mask-closable="false">
      <n-form ref="installFormRef" :model="installForm">
        <n-form-item label="插件文件路径" required>
          <n-input
            v-model:value="installForm.pluginPath"
            placeholder="选择 .so 或 .dll 插件文件"
            readonly
          >
            <template #suffix>
              <n-button size="small" @click="selectPluginFile">浏览</n-button>
            </template>
          </n-input>
        </n-form-item>
      </n-form>
      <template #action>
        <n-space justify="end">
          <n-button @click="showInstallModal = false">取消</n-button>
          <n-button type="primary" @click="confirmInstall" :loading="installing">安装</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 插件详情对话框 -->
    <n-modal v-model:show="showDetailModal" preset="dialog" title="插件详情" :mask-closable="false" style="width: 600px;">
      <n-descriptions bordered :column="2">
        <n-descriptions-item label="插件 ID">{{ selectedPlugin?.metadata?.id }}</n-descriptions-item>
        <n-descriptions-item label="版本号">{{ selectedPlugin?.metadata?.version }}</n-descriptions-item>
        <n-descriptions-item label="插件名称" :span="2">{{ selectedPlugin?.metadata?.name }}</n-descriptions-item>
        <n-descriptions-item label="插件描述" :span="2">{{ selectedPlugin?.metadata?.description }}</n-descriptions-item>
        <n-descriptions-item label="作者">{{ selectedPlugin?.metadata?.author }}</n-descriptions-item>
        <n-descriptions-item label="类型">
          <n-tag :type="selectedPlugin?.metadata?.type === 'builtin' ? 'warning' : 'info'">
            {{ selectedPlugin?.metadata?.type === 'builtin' ? '内置插件' : '外置插件' }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="分类">{{ selectedPlugin?.metadata?.category }}</n-descriptions-item>
        <n-descriptions-item label="状态">
          <n-tag :type="getStatusType(selectedPlugin)">
            {{ getStatusText(selectedPlugin) }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="标签" :span="2" v-if="selectedPlugin?.metadata?.tags">
          <n-space>
            <n-tag v-for="tag in selectedPlugin?.metadata?.tags" :key="tag">{{ tag }}</n-tag>
          </n-space>
        </n-descriptions-item>
        <n-descriptions-item label="许可证" :span="2" v-if="selectedPlugin?.metadata?.license">
          {{ selectedPlugin?.metadata?.license }}
        </n-descriptions-item>
      </n-descriptions>
      <template #action>
        <n-space justify="end">
          <n-button @click="showDetailModal = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h, onMounted } from 'vue';
import { NButton, NTag, NSpace, NIcon, useMessage, useDialog } from 'naive-ui';
import { AddOutline, RefreshOutline, InformationCircleOutline, PlayOutline, StopOutline, TrashOutline } from '@vicons/ionicons5';
import { Plugins } from '../../bindings/bindings';

interface PluginInfo {
  id: string;
  metadata: {
    id: string;
    name: string;
    version: string;
    description: string;
    author: string;
    type: 'builtin' | 'external';
    category: string;
    tags?: string[];
    license?: string;
  };
  status: string;
  is_enabled: boolean;
  can_disable: boolean;
  error?: string;
}

interface PluginStats {
  total: number;
  enabled: number;
  builtin: number;
  external: number;
}

const message = useMessage();
const dialog = useDialog();

const loading = ref(false);
const pluginList = ref<PluginInfo[]>([]);
const showInstallModal = ref(false);
const showDetailModal = ref(false);
const installing = ref(false);
const selectedPlugin = ref<PluginInfo | null>(null);

const installForm = reactive({
  pluginPath: '',
});

const stats = computed<PluginStats>(() => {
  return {
    total: pluginList.value.length,
    enabled: pluginList.value.filter(p => p.is_enabled).length,
    builtin: pluginList.value.filter(p => p.metadata.type === 'builtin').length,
    external: pluginList.value.filter(p => p.metadata.type === 'external').length,
  };
});

const columns = [
  {
    title: '插件名称',
    key: 'name',
    width: 200,
    render: (row: PluginInfo) => {
      return h('div', { style: 'display: flex; align-items: center; gap: 8px;' }, [
        h('strong', {}, row.metadata.name),
        h(NTag, { type: row.metadata.type === 'builtin' ? 'warning' : 'info', size: 'small' }, 
          { default: () => row.metadata.type === 'builtin' ? '内置' : '外置' }
        ),
      ]);
    },
  },
  {
    title: '版本',
    key: 'version',
    width: 80,
    render: (row: PluginInfo) => row.metadata.version,
  },
  {
    title: '描述',
    key: 'description',
    ellipsis: { tooltip: true },
    render: (row: PluginInfo) => row.metadata.description,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: PluginInfo) => {
      return h(NTag, { type: getStatusType(row), size: 'small' }, 
        { default: () => getStatusText(row) }
      );
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 280,
    fixed: 'right',
    render: (row: PluginInfo) => {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, {
            size: 'small',
            type: 'primary',
            onClick: () => showDetail(row),
          }, {
            default: () => '详情',
            icon: () => h(NIcon, {}, { default: () => h(InformationCircleOutline) }),
          }),
          row.metadata.type === 'external' ? (
            row.is_enabled ?
              h(NButton, {
                size: 'small',
                type: 'warning',
                onClick: () => handleDisable(row),
              }, {
                default: () => '禁用',
                icon: () => h(NIcon, {}, { default: () => h(StopOutline) }),
              }) :
              h(NButton, {
                size: 'small',
                type: 'success',
                onClick: () => handleEnable(row),
              }, {
                default: () => '启用',
                icon: () => h(NIcon, {}, { default: () => h(PlayOutline) }),
              })
          ) : null,
          row.metadata.type === 'external' ?
            h(NButton, {
              size: 'small',
              type: 'error',
              onClick: () => handleUninstall(row),
            }, {
              default: () => '卸载',
              icon: () => h(NIcon, {}, { default: () => h(TrashOutline) }),
            }) : null,
        ],
      });
    },
  },
];

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
});

const getStatusType = (row: PluginInfo) => {
  if (row.error) return 'error';
  if (row.is_enabled) return 'success';
  return 'default';
};

const getStatusText = (row: PluginInfo) => {
  if (row.error) return '错误';
  if (row.is_enabled) return '已启用';
  if (row.status === 'loaded') return '已加载';
  return '已禁用';
};

const loadPluginList = async () => {
  loading.value = true;
  try {
    const response = await Plugins.GetPluginList();
    pluginList.value = response.plugins || [];
  } catch (error: any) {
    message.error('加载插件列表失败：' + error.message);
  } finally {
    loading.value = false;
  }
};

const refreshList = () => {
  loadPluginList();
};

const handleInstall = () => {
  installForm.pluginPath = '';
  showInstallModal.value = true;
};

const selectPluginFile = async () => {
  // TODO: 实现文件选择器
  message.info('文件选择功能尚未实现，请手动输入路径');
};

const confirmInstall = async () => {
  if (!installForm.pluginPath) {
    message.warning('请选择插件文件');
    return;
  }

  installing.value = true;
  try {
    const response = await Plugins.InstallPlugin({ plugin_path: installForm.pluginPath });
    if (response.success) {
      message.success('插件安装成功');
      showInstallModal.value = false;
      await loadPluginList();
    } else {
      message.error('安装失败：' + response.message);
    }
  } catch (error: any) {
    message.error('安装失败：' + error.message);
  } finally {
    installing.value = false;
  }
};

const showDetail = (row: PluginInfo) => {
  selectedPlugin.value = row;
  showDetailModal.value = true;
};

const handleEnable = async (row: PluginInfo) => {
  try {
    const response = await Plugins.EnablePlugin({ plugin_id: row.id });
    if (response.success) {
      message.success('插件已启用');
      await loadPluginList();
    } else {
      message.error('启用失败：' + response.message);
    }
  } catch (error: any) {
    message.error('启用失败：' + error.message);
  }
};

const handleDisable = async (row: PluginInfo) => {
  dialog.warning({
    title: '禁用插件',
    content: `确定要禁用插件 "${row.metadata.name}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const response = await Plugins.DisablePlugin({ plugin_id: row.id });
        if (response.success) {
          message.success('插件已禁用');
          await loadPluginList();
        } else {
          message.error('禁用失败：' + response.message);
        }
      } catch (error: any) {
        message.error('禁用失败：' + error.message);
      }
    },
  });
};

const handleUninstall = async (row: PluginInfo) => {
  dialog.warning({
    title: '卸载插件',
    content: `确定要卸载插件 "${row.metadata.name}" 吗？此操作不可恢复。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const response = await Plugins.UninstallPlugin({ plugin_id: row.id });
        if (response.success) {
          message.success('插件已卸载');
          await loadPluginList();
        } else {
          message.error('卸载失败：' + response.message);
        }
      } catch (error: any) {
        message.error('卸载失败：' + error.message);
      }
    },
  });
};

onMounted(() => {
  loadPluginList();
});
</script>

<style scoped lang="scss">
.plugin-manager {
  padding: 16px;
  height: 100%;
  overflow: auto;
}
</style>
