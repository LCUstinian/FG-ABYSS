<template>
  <div class="payload-template-manager-view">
    <n-card title="模板管理" :bordered="false">
      <template #header-extra>
        <n-button type="primary" size="small" @click="showCreateModal = true">
          <template #icon>
            <span>➕</span>
          </template>
          新建模板
        </n-button>
      </template>

      <!-- 模板列表 -->
      <div class="template-grid">
        <n-grid :cols="24" :x-gap="16" :y-gap="16">
          <n-grid-item
            v-for="template in templates"
            :key="template.id"
            :span="8"
          >
            <n-card
              :title="template.name"
              :bordered="false"
              hoverable
              class="template-card"
            >
              <template #header-extra>
                <n-space>
                  <n-button quaternary circle size="small" @click="editTemplate(template)">
                    <template #icon>
                      <span>✏️</span>
                    </template>
                  </n-button>
                  <n-button
                    quaternary
                    circle
                    size="small"
                    type="error"
                    @click="deleteTemplate(template.id)"
                  >
                    <template #icon>
                      <span>🗑️</span>
                    </template>
                  </n-button>
                </n-space>
              </template>

              <div class="template-info">
                <n-space vertical>
                  <n-tag :type="getScriptTypeColor(template.script_type)" size="small">
                    {{ template.script_type.toUpperCase() }}
                  </n-tag>
                  <n-text depth="3" style="font-size: 13px;">
                    {{ template.description || '暂无描述' }}
                  </n-text>
                  <n-divider style="margin: 12px 0;" />
                  <n-space justify="space-between">
                    <n-text depth="3" style="font-size: 12px;">
                      创建：{{ formatDate(template.created_at) }}
                    </n-text>
                    <n-text depth="3" style="font-size: 12px;">
                      更新：{{ formatDate(template.updated_at) }}
                    </n-text>
                  </n-space>
                </n-space>
              </div>

              <template #action>
                <n-button block @click="useTemplate(template)">
                  使用模板
                </n-button>
              </template>
            </n-card>
          </n-grid-item>

          <!-- 空状态 -->
          <n-grid-item v-if="templates.length === 0" :span="24">
            <n-empty
              description="暂无模板，点击右上角创建第一个模板"
              style="padding: 60px 20px;"
            />
          </n-grid-item>
        </n-grid>
      </div>
    </n-card>

    <!-- 创建/编辑模板对话框 -->
    <n-modal
      v-model:show="showCreateModal"
      :title="editingTemplate ? '编辑模板' : '创建模板'"
      preset="dialog"
      :style="{ width: '600px' }"
    >
      <n-form
        ref="formRef"
        :model="formData"
        label-placement="left"
        :label-width="100"
      >
        <n-form-item
          label="模板名称"
          :rule="{ required: true, message: '请输入模板名称', trigger: 'blur' }"
        >
          <n-input v-model:value="formData.name" placeholder="例如：PHP 基础连接" />
        </n-form-item>

        <n-form-item
          label="脚本类型"
          :rule="{ required: true, message: '请选择脚本类型', trigger: 'change' }"
        >
          <n-select
            v-model:value="formData.script_type"
            :options="[
              { label: 'PHP', value: 'php' },
              { label: 'JSP', value: 'jsp' },
              { label: 'ASPX', value: 'aspx' },
              { label: 'ASP', value: 'asp' },
            ]"
          />
        </n-form-item>

        <n-form-item label="功能类型">
          <n-select
            v-model:value="formData.function_type"
            :options="[
              { label: '基础连接', value: 'basic' },
              { label: '文件管理', value: 'file_manager' },
              { label: '进程管理', value: 'process_manager' },
              { label: '注册表', value: 'registry' },
              { label: '网络', value: 'network' },
            ]"
          />
        </n-form-item>

        <n-form-item label="模板描述">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            placeholder="简要描述模板的用途"
            :rows="3"
          />
        </n-form-item>

        <n-form-item label="模板代码">
          <n-input
            v-model:value="formData.code"
            type="textarea"
            placeholder="输入模板代码"
            :rows="10"
            style="font-family: monospace;"
          />
        </n-form-item>
      </n-form>

      <template #action>
        <n-button @click="showCreateModal = false">取消</n-button>
        <n-button type="primary" @click="saveTemplate">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useMessage } from 'naive-ui'
import type { ScriptType, FunctionType } from '@/types/payload'

const message = useMessage()

const showCreateModal = ref(false)
const editingTemplate = ref<any>(null)

const formData = reactive({
  name: '',
  script_type: 'php' as ScriptType,
  function_type: 'basic' as FunctionType,
  description: '',
  code: '',
})

// 模拟模板数据
const templates = ref([
  {
    id: 1,
    name: 'PHP 基础连接',
    script_type: 'php',
    function_type: 'basic',
    description: '最基础的 PHP 一句话木马',
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
  {
    id: 2,
    name: 'ASP 文件管理',
    script_type: 'asp',
    function_type: 'file_manager',
    description: 'ASP 文件管理模板',
    created_at: '2024-01-02',
    updated_at: '2024-01-02',
  },
])

// 方法
const getScriptTypeColor = (type: string) => {
  const colorMap: Record<string, any> = {
    php: 'success',
    jsp: 'warning',
    aspx: 'info',
    asp: 'error',
  }
  return colorMap[type] || 'default'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN')
}

const editTemplate = (template: any) => {
  editingTemplate.value = template
  formData.name = template.name
  formData.script_type = template.script_type
  formData.function_type = template.function_type
  formData.description = template.description
  formData.code = template.code || ''
  showCreateModal.value = true
}

const deleteTemplate = (id: number) => {
  templates.value = templates.value.filter(t => t.id !== id)
  message.success('模板已删除')
}

const useTemplate = (template: any) => {
  message.success(`已使用模板：${template.name}`)
  // 这里可以集成到 PayloadGenerator
}

const saveTemplate = () => {
  if (!formData.name) {
    message.error('请输入模板名称')
    return
  }

  if (editingTemplate.value) {
    // 编辑模式
    const index = templates.value.findIndex(t => t.id === editingTemplate.value.id)
    if (index !== -1) {
      templates.value[index] = {
        ...templates.value[index],
        ...formData,
        updated_at: new Date().toISOString(),
      }
      message.success('模板已更新')
    }
  } else {
    // 创建模式
    templates.value.push({
      id: Date.now(),
      ...formData,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    })
    message.success('模板已创建')
  }

  showCreateModal.value = false
  editingTemplate.value = null
  Object.assign(formData, {
    name: '',
    script_type: 'php',
    function_type: 'basic',
    description: '',
    code: '',
  })
}
</script>

<style scoped>
.payload-template-manager-view {
  width: 100%;
  min-height: 100%;
  animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 主卡片样式优化 */
:deep(.n-card) {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.n-card-header) {
  border-bottom: 1px solid var(--border-color);
  padding: 20px 24px;
  background: linear-gradient(135deg, var(--card-bg) 0%, var(--card-bg-hover) 100%);
  flex-shrink: 0;
}

:deep(.n-card-header__main) {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

:deep(.n-card__content) {
  padding: 24px;
  flex: 1;
  overflow: auto;
}

/* 按钮样式优化 */
:deep(.n-button) {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.n-button:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.template-grid {
  margin-top: 16px;
}

.template-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  height: 100%;
  position: relative;
  overflow: hidden;
  border-radius: 12px;
}

.template-card::before {
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
  z-index: 1;
}

.template-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

.template-card:hover::before {
  transform: scaleX(1);
}

:deep(.template-card .n-card-header) {
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
  padding: 16px 20px;
}

:deep(.template-card .n-card-header__main) {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  transition: color 0.3s ease;
}

.template-card:hover :deep(.n-card-header__main) {
  color: var(--active-color);
}

.template-info {
  min-height: 120px;
}

/* 标签样式优化 */
:deep(.n-tag) {
  border-radius: 6px;
  font-weight: 500;
  padding: 4px 12px;
}

/* 分割线样式优化 */
:deep(.n-divider) {
  border-color: var(--border-color);
  margin: 12px 0;
}

/* 对话框样式优化 */
:deep(.n-modal-content) {
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

:deep(.n-dialog__header) {
  padding: 24px 28px 20px;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, var(--card-bg) 0%, var(--card-bg-hover) 100%);
}

:deep(.n-dialog__title) {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

:deep(.n-dialog__content) {
  padding: 24px 28px;
}

:deep(.n-dialog__footer) {
  padding: 20px 28px 24px;
  border-top: 1px solid var(--border-color);
}

/* 表单样式优化 */
:deep(.n-form-item-label) {
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

:deep(.n-input),
:deep(.n-select),
:deep(.n-input-number) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.n-input:hover),
:deep(.n-select:hover),
:deep(.n-input-number:hover) {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

:deep(.n-input:focus-within),
:deep(.n-select:focus-within),
:deep(.n-input-number:focus-within) {
  box-shadow: 0 0 0 3px rgba(var(--active-color-rgb), 0.1);
}

/* 空状态样式优化 */
:deep(.n-empty) {
  padding: 60px 20px;
}

:deep(.n-empty-description) {
  color: var(--text-secondary);
  margin-top: 12px;
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  :deep(.n-card-header) {
    padding: 24px 28px;
  }
  
  :deep(.n-card-header__main) {
    font-size: 20px;
  }
  
  :deep(.n-card__content) {
    padding: 28px;
  }
  
  :deep(.template-card .n-card-header) {
    padding: 20px 24px;
  }
  
  :deep(.template-card .n-card-header__main) {
    font-size: 17px;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  :deep(.n-card-header) {
    padding: 16px 20px;
  }
  
  :deep(.n-card__content) {
    padding: 16px;
  }
  
  :deep(.template-card .n-card-header) {
    padding: 14px 16px;
  }
  
  :deep(.template-card .n-card-header__main) {
    font-size: 15px;
  }
  
  .template-card:hover {
    transform: translateY(-4px);
  }
  
  :deep(.n-dialog__header) {
    padding: 20px 20px 16px;
  }
  
  :deep(.n-dialog__title) {
    font-size: 18px;
  }
  
  :deep(.n-dialog__content) {
    padding: 20px;
  }
  
  :deep(.n-dialog__footer) {
    padding: 16px 20px 20px;
  }
}
</style>
