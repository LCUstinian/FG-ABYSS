<template>
  <div class="payload-generator-view">
    <div class="generator-layout">
      <!-- 左侧配置区 -->
      <div class="config-panel">
        <n-card title="生成配置" :bordered="false" size="small">
          <!-- 模式切换 -->
          <n-form-item label="生成模式" label-placement="top" size="small">
            <n-segmented
              v-model:value="config.mode"
              :options="[
                { label: '极简模式', value: 'simple' },
                { label: '高级加密', value: 'advanced' },
              ]"
              block
              @update:value="onModeChange"
            />
          </n-form-item>

          <n-divider style="margin: 12px 0;" />

          <!-- 配置表单 -->
          <div class="config-form">
            <!-- 第一行：脚本类型 + 功能类型 -->
            <div class="form-row">
              <n-form-item label="脚本类型" label-placement="top" size="small" class="form-item-half">
                <n-select
                  v-model:value="config.script_type"
                  :options="[
                    { label: 'PHP', value: 'php' },
                    { label: 'JSP', value: 'jsp' },
                    { label: 'ASPX', value: 'aspx' },
                    { label: 'ASP', value: 'asp' },
                  ]"
                  placeholder="请选择"
                  size="small"
                />
              </n-form-item>

              <n-form-item label="功能类型" label-placement="top" size="small" class="form-item-half">
                <n-select
                  v-model:value="config.function_type"
                  :options="FUNCTION_TYPE_OPTIONS"
                  placeholder="请选择"
                  size="small"
                />
              </n-form-item>
            </div>

            <!-- 第二行：编码器/加密算法 + 混淆强度 -->
            <div class="form-row">
              <!-- Simple 模式专属选项 -->
              <template v-if="isSimpleMode">
                <n-form-item label="编码器" label-placement="top" size="small" class="form-item-half">
                  <n-select
                    v-model:value="config.encode_type"
                    :options="[
                      { label: '无编码', value: 'none' },
                      { label: 'Base64', value: 'base64' },
                      { label: 'XOR', value: 'xor' },
                      { label: 'GZInflate', value: 'gzinflate' },
                      { label: 'Hex', value: 'hex' },
                      { label: 'URL', value: 'urlencode' },
                      { label: 'ROT13', value: 'rot13' },
                    ]"
                    placeholder="请选择"
                    size="small"
                  />
                </n-form-item>
              </template>

              <!-- Advanced 模式专属选项 -->
              <template v-if="isAdvancedMode">
                <n-form-item label="加密算法" label-placement="top" size="small" class="form-item-half">
                  <n-select
                    v-model:value="config.encrypt_algo"
                    :options="[
                      { label: 'AES-128', value: 'aes128_cbc' },
                      { label: 'AES-256', value: 'aes256_cbc' },
                      { label: 'XOR', value: 'xor' },
                    ]"
                    placeholder="请选择"
                    size="small"
                  />
                </n-form-item>
              </template>

              <n-form-item label="混淆强度" label-placement="top" size="small" class="form-item-half">
                <n-select
                  v-model:value="obfuscationValue"
                  :options="[
                    { label: '低', value: 1 },
                    { label: '中', value: 2 },
                    { label: '高', value: 3 },
                  ]"
                  placeholder="请选择"
                  size="small"
                />
              </n-form-item>
            </div>

            <!-- 第三行：连接密码 + 输出文件名 -->
            <div class="form-row">
              <n-form-item 
                label="连接密码" 
                label-placement="top"
                size="small"
                class="form-item-half"
                :rule="{
                  required: true,
                  message: '请输入密码',
                  trigger: 'blur',
                }"
              >
                <n-input
                  v-model:value="config.password"
                  placeholder="请输入密码"
                  type="password"
                  show-password-on="click"
                  clearable
                  size="small"
                />
              </n-form-item>

              <n-form-item label="输出文件名" label-placement="top" size="small" class="form-item-half">
                <n-input
                  v-model:value="config.output_filename"
                  placeholder="留空自动生成"
                  clearable
                  size="small"
                />
              </n-form-item>
            </div>
          </div>

          <!-- 生成按钮 -->
          <div class="button-group">
            <n-button
              type="primary"
              size="medium"
              block
              :loading="isGenerating"
              @click="handleGenerate"
              class="primary-btn"
            >
              <template #icon>
                <IconCode />
              </template>
              生成载荷
            </n-button>

            <n-button
              secondary
              size="medium"
              block
              @click="handlePreview"
              class="secondary-btn"
            >
              <template #icon>
                <IconEye />
              </template>
              实时预览
            </n-button>
          </div>
        </n-card>
      </div>

      <!-- 右侧预览区 -->
      <div class="preview-panel">
        <n-card title="代码预览" :bordered="false" size="medium">
          <template #header-extra>
            <n-space>
              <n-tooltip>
                <template #trigger>
                  <n-button quaternary circle @click="copyCode">
                    <template #icon>
                      <IconCopy />
                    </template>
                  </n-button>
                </template>
                复制代码
              </n-tooltip>

              <n-tooltip>
                <template #trigger>
                  <n-button quaternary circle @click="downloadFile">
                    <template #icon>
                      <IconDownload />
                    </template>
                  </n-button>
                </template>
                下载文件
              </n-tooltip>

              <n-tooltip v-if="isAdvancedMode && generatedResult?.client_config">
                <template #trigger>
                  <n-button quaternary circle @click="exportConfig">
                    <template #icon>
                      <IconSettings />
                    </template>
                  </n-button>
                </template>
                导出客户端配置
              </n-tooltip>
            </n-space>
          </template>

          <!-- 加载状态 -->
          <n-spin :show="isGenerating">
            <template #description>
              正在生成载荷代码...
            </template>

            <!-- 代码预览区 -->
            <div class="code-preview-container">
              <pre v-if="generatedResult?.code" class="code-block"><code>{{ generatedResult.code }}</code></pre>
              
              <n-empty
                v-else
                description="点击生成按钮预览代码"
                style="padding: 60px 20px;"
              />
            </div>
          </n-spin>

          <!-- 状态信息 -->
          <div v-if="generatedResult" class="status-bar">
            <n-space justify="space-between">
              <n-text depth="3">
                文件名：{{ generatedResult.filename }}
              </n-text>
              <n-text depth="3">
                大小：{{ formatSize(generatedResult.size) }}
              </n-text>
            </n-space>
          </div>
        </n-card>

        <!-- 历史记录 -->
        <n-card title="生成历史" :bordered="false" size="small" style="margin-top: 16px;">
          <template #header-extra>
            <n-button text type="primary" @click="clearHistory">
              清空
            </n-button>
          </template>

          <n-scrollbar style="max-height: 200px;">
            <div v-if="history.length > 0" class="history-list">
              <div
                v-for="(item, index) in history"
                :key="index"
                class="history-item"
                @click="loadHistoryItem(item)"
              >
                <n-tag :type="item.success ? 'success' : 'error'" size="small">
                  {{ item.filename }}
                </n-tag>
                <n-text depth="3" style="font-size: 12px;">
                  {{ formatSize(item.size) }}
                </n-text>
              </div>
            </div>
            <n-empty v-else description="暂无历史记录" size="small" />
          </n-scrollbar>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { 
  IconCode, 
  IconEye, 
  IconCopy, 
  IconDownload, 
  IconSettings 
} from '@tabler/icons-vue'
import { usePayloadStore } from '@/stores/payload'
import type { PayloadConfig, ObfuscationLevel, FunctionType } from '@/types/payload'
import { FUNCTION_TYPE_OPTIONS } from '@/types/payload'

const message = useMessage()
const payloadStore = usePayloadStore()

// 本地状态 - 直接使用 store 中的 config
const config = computed(() => payloadStore.config)

// 混淆级别滑块值 (1-3)
const obfuscationValue = ref(1)

// 计算属性
const isSimpleMode = computed(() => config.value.mode === 'simple')
const isAdvancedMode = computed(() => config.value.mode === 'advanced')
const isGenerating = computed(() => payloadStore.isGenerating)
const generatedResult = computed(() => payloadStore.generatedResult)
const history = computed(() => payloadStore.history)

// 模式切换
const onModeChange = (mode: string) => {
  if (mode === 'advanced') {
    config.value.encode_type = undefined
    config.value.encrypt_algo = 'aes128_cbc'
  } else {
    config.value.encrypt_algo = undefined
    config.value.encode_type = 'none'
  }
}

// 混淆级别映射
const mapObfuscationLevel = (value: number): ObfuscationLevel => {
  const map: Record<number, ObfuscationLevel> = {
    1: 'low',
    2: 'medium',
    3: 'high',
  }
  return map[value] || 'low'
}

// 生成载荷
const handleGenerate = async () => {
  // 更新混淆级别
  payloadStore.setObfuscationLevel(mapObfuscationLevel(obfuscationValue.value))
  
  try {
    await payloadStore.generate()
    message.success('载荷生成成功')
  } catch (error: any) {
    message.error(error.message || '生成失败')
  }
}

// 实时预览
const handlePreview = async () => {
  payloadStore.setObfuscationLevel(mapObfuscationLevel(obfuscationValue.value))
  
  try {
    await payloadStore.generatePreview(300)
  } catch (error: any) {
    message.error(error.message || '预览失败')
  }
}

// 复制代码
const copyCode = async () => {
  if (!generatedResult.value?.code) return
  
  try {
    await payloadStore.copyCode(generatedResult.value.code)
    message.success('代码已复制到剪贴板')
  } catch (error: any) {
    message.error('复制失败')
  }
}

// 下载文件
const downloadFile = () => {
  if (!generatedResult.value?.code) return
  
  const blob = new Blob([generatedResult.value.code], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = generatedResult.value.filename || 'payload.txt'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  message.success('文件已开始下载')
}

// 导出客户端配置
const exportConfig = async () => {
  try {
    // 这里可以集成 Tauri 文件对话框
    message.info('导出配置功能开发中')
  } catch (error: any) {
    message.error('导出失败')
  }
}

// 加载历史记录
const loadHistoryItem = (item: any) => {
  payloadStore.generatedResult = item
  message.success('已加载历史记录')
}

// 清空历史
const clearHistory = () => {
  payloadStore.clearHistory()
  message.success('历史记录已清空')
}

// 格式化文件大小
const formatSize = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

// 加载历史记录
onMounted(async () => {
  await payloadStore.loadHistory()
})
</script>

<style scoped>
.payload-generator-view {
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

.generator-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  max-width: 1600px;
  margin: 0 auto;
  padding: 0;
}

/* 自适应窗口布局 */
@media (max-width: 1200px) {
  .generator-layout {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .generator-layout {
    gap: 12px;
  }
}

.config-panel,
.preview-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

/* 统一卡片样式 - 所有卡片一致 */
:deep(.n-card) {
  background: linear-gradient(145deg, var(--card-bg) 0%, var(--card-bg-hover) 100%);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

:deep(.n-card::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--active-color), transparent, var(--active-color));
  opacity: 0.6;
}

:deep(.n-card:hover) {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

:deep(.n-card-header) {
  border-bottom: 1px solid var(--border-color);
  padding: 14px 18px;
  background: transparent;
}

:deep(.n-card-header__main) {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.3px;
}

:deep(.n-card__content) {
  padding: 16px 18px;
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

/* 代码预览容器优化 */
.code-preview-container {
  background: linear-gradient(135deg, var(--code-bg, var(--card-bg-hover)) 0%, var(--card-bg) 100%);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 24px;
  min-height: 400px;
  max-height: 600px;
  overflow: auto;
  position: relative;
  box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.05);
}

.code-preview-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--active-color) 0%, transparent 100%);
  border-radius: 12px 12px 0 0;
}

/* 自定义滚动条 */
.code-preview-container::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.code-preview-container::-webkit-scrollbar-track {
  background: var(--sidebar-bg);
  border-radius: 4px;
}

.code-preview-container::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.code-preview-container::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

.code-block {
  margin: 0;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--code-text, var(--text-primary));
  position: relative;
  z-index: 1;
}

.status-bar {
  margin-top: 16px;
  padding: 16px 20px;
  background: linear-gradient(135deg, var(--card-bg-hover) 0%, var(--card-bg) 100%);
  border-radius: 12px;
  border: 1px solid var(--border-color);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.status-bar:hover {
  border-color: var(--active-color);
  box-shadow: 0 4px 12px rgba(var(--active-color-rgb), 0.1);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--card-bg-hover);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid transparent;
  position: relative;
  overflow: hidden;
}

.history-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, var(--active-color-bg) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.history-item:hover {
  background: var(--active-color-bg);
  border-color: var(--active-color);
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(var(--active-color-rgb), 0.15);
}

.history-item:hover::before {
  opacity: 1;
}

/* 分割线优化 */
:deep(.n-divider) {
  border-color: var(--border-color);
  margin: 10px 0;
}

/* 表单样式优化 - 紧凑版 */
:deep(.n-form-item) {
  margin-bottom: 10px;
}

:deep(.n-form-item-label) {
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 5px;
  font-size: 12px;
  letter-spacing: 0.2px;
}

:deep(.n-form-item-blank) {
  padding: 0;
}

:deep(.n-input),
:deep(.n-select) {
  border-radius: 6px;
  transition: all 0.2s ease;
}

:deep(.n-input-wrapper),
:deep(.n-select-trigger) {
  border-radius: 6px;
}

:deep(.n-input:hover),
:deep(.n-select:hover) {
  border-color: var(--text-secondary);
}

:deep(.n-input:focus-within),
:deep(.n-select:focus-within) {
  box-shadow: 0 0 0 2px rgba(var(--active-color-rgb), 0.15);
  border-color: var(--active-color);
}

:deep(.n-input__wrapper),
:deep(.n-select__trigger) {
  background: var(--card-bg);
}

/* 表单布局 */
.config-form {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-item-half {
  margin-bottom: 12px;
}

.form-item-half :deep(.n-form-item-blank) {
  padding: 0;
}

/* 响应式设计 - 表单在小屏幕单列 */
@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
    gap: 0;
  }
  
  .form-item-half {
    margin-bottom: 10px;
  }
  
  .config-panel :deep(.n-card__content),
  .preview-panel :deep(.n-card__content) {
    padding: 12px 14px;
  }
  
  .config-panel :deep(.n-card-header),
  .preview-panel :deep(.n-card-header) {
    padding: 12px 14px;
  }
  
  .button-group {
    margin-top: 10px;
    padding-top: 10px;
  }
}

/* 中等屏幕优化 */
@media (min-width: 769px) and (max-width: 1024px) {
  .form-row {
    gap: 12px;
  }
  
  .form-item-half {
    margin-bottom: 10px;
  }
}

/* 按钮组样式 */
.button-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px dashed var(--border-color);
}

.button-group :deep(.n-button) {
  border-radius: 8px;
  font-weight: 500;
  letter-spacing: 0.2px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.button-group :deep(.primary-btn) {
  background: linear-gradient(135deg, var(--active-color) 0%, var(--active-color) 100%);
  box-shadow: 0 4px 12px rgba(var(--active-color-rgb), 0.3);
}

.button-group :deep(.primary-btn:hover) {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(var(--active-color-rgb), 0.4);
}

.button-group :deep(.primary-btn:active) {
  transform: translateY(0);
}

.button-group :deep(.secondary-btn:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 分段器优化 */
:deep(.n-segmented) {
  border-radius: 10px;
  padding: 4px;
  background: var(--card-bg-hover);
}

:deep(.n-segmented-item) {
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 6px;
  transition: all 0.2s ease;
}

:deep(.n-segmented-item--selected) {
  background: var(--active-color) !important;
  color: white !important;
  box-shadow: 0 2px 8px rgba(var(--active-color-rgb), 0.3);
}

/* 主题适配 */
:root {
  --card-bg: #ffffff;
  --card-bg-hover: #f8fafc;
  --border-color: #e2e8f0;
  --text-primary: #1e293b;
  --text-secondary: #64748b;
  --active-color: #3b82f6;
  --active-color-bg: rgba(59, 130, 246, 0.08);
  --active-color-rgb: 59, 130, 246;
  --sidebar-bg: #f1f5f9;
  --code-bg: #f8fafc;
  --code-text: #1e293b;
}

html.dark {
  --card-bg: #1e293b;
  --card-bg-hover: #334155;
  --border-color: #475569;
  --text-primary: #f1f5f9;
  --text-secondary: #94a3b8;
  --active-color: #60a5fa;
  --active-color-bg: rgba(96, 165, 250, 0.15);
  --active-color-rgb: 96, 165, 250;
  --sidebar-bg: #0f172a;
  --code-bg: #0f172a;
  --code-text: #e2e8f0;
}

/* 桌面端优化 */
@media (min-width: 1440px) {
  .generator-layout {
    gap: 32px;
    max-width: 1800px;
  }
  
  .code-preview-container {
    padding: 28px;
  }
  
  .code-block {
    font-size: 14px;
    line-height: 1.8;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .generator-layout {
    gap: 16px;
  }
  
  .code-preview-container {
    padding: 16px;
    min-height: 300px;
  }
  
  .code-block {
    font-size: 12px;
    line-height: 1.5;
  }
  
  .status-bar {
    padding: 12px 16px;
  }
}
</style>
