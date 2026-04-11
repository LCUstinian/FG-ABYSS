<template>
  <div class="payload-generator-view">
    <div class="generator-layout">
      <!-- 左侧配置区 -->
      <div class="config-panel">
        <n-card :title="$t('payloads.configTitle')" :bordered="false" size="small">
          <!-- 多模态切换 -->
          <div class="mode-selector-wrapper">
            <div class="mode-selector-label">生成模式</div>
            <n-segmented
              v-model:value="currentMode"
              :options="[
                { label: '文件落地', value: 'file_based' },
                { label: '文件通用', value: 'file_common' },
                { label: '文件代理', value: 'file_proxy' },
                { label: '文件混合', value: 'file_hybrid' },
                { label: '内存马注入', value: 'memory_shell' },
                { label: '纯 Suo5 代理', value: 'suo5_only' },
              ]"
              block
              @update:value="onModeChange"
            />
          </div>

          <n-divider style="margin: 12px 0;" />

          <!-- 基于模式切换配置 -->
          <FileShellConfig v-if="isFileMode" />

          <!-- 传统配置表单（其他模式） -->
          <div v-else class="config-form">
            <!-- 第一行：脚本类型 + 注入类型（仅内存马模式） -->
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
                  placeholder="选择脚本类型"
                  size="small"
                />
              </n-form-item>

              <n-form-item v-if="isMemoryShellMode" label="注入类型" label-placement="top" size="small" class="form-item-half">
                <n-select
                  v-model:value="config.injection_type"
                  :options="INJECTION_TYPE_OPTIONS"
                  placeholder="选择注入类型"
                  size="small"
                />
              </n-form-item>
              <div v-else class="form-item-half"></div>
            </div>

            <!-- 第二行：混淆强度 + 自毁逻辑（仅内存马模式） -->
            <div class="form-row">
              <n-form-item label="混淆强度" label-placement="top" size="small" class="form-item-half">
                <n-select
                  v-model:value="obfuscationValue"
                  :options="[
                    { label: 'L1 (轻量)', value: 1 },
                    { label: 'L2 (中等)', value: 2 },
                    { label: 'L3 (高级)', value: 3 },
                  ]"
                  placeholder="选择混淆强度"
                  size="small"
                />
              </n-form-item>

              <n-form-item v-if="isMemoryShellMode" label="自毁逻辑" label-placement="top" size="small" class="form-item-half">
                <n-switch v-model:value="config.self_destruct" />
              </n-form-item>
              <div v-else class="form-item-half"></div>
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
                  message: '密码不能为空',
                  trigger: 'blur',
                }"
              >
                <n-input
                  v-model:value="config.password"
                  placeholder="输入连接密码"
                  type="password"
                  show-password-on="click"
                  clearable
                  size="small"
                />
              </n-form-item>

              <n-form-item label="输出文件名" label-placement="top" size="small" class="form-item-half">
                <n-input
                  v-model:value="config.output_filename"
                  placeholder="输入输出文件名"
                  clearable
                  size="small"
                />
              </n-form-item>
            </div>

            <!-- Suo5 专属配置面板 -->
            <n-collapse v-if="isSuo5OnlyMode || isFileProxyMode || isFileHybridMode">
              <n-collapse-item title="Suo5 配置">
                <!-- 第一行：认证密码 + 随机生成 -->
                <div class="form-row">
                  <n-form-item label="Suo5 认证密码" label-placement="top" size="small" class="form-item-half">
                    <n-input
                      v-model:value="suo5Auth"
                      placeholder="输入认证密码"
                      type="password"
                      show-password-on="click"
                      clearable
                      size="small"
                    />
                    <n-text v-if="suo5Auth" depth="3" style="font-size: 12px; margin-top: 4px;">
                      密码强度: {{ getPasswordStrength(suo5Auth) }}
                    </n-text>
                  </n-form-item>
                  <n-form-item label=" " label-placement="top" size="small" class="form-item-half">
                    <n-button
                      secondary
                      size="small"
                      block
                      @click="generateSuo5Password"
                    >
                      随机生成密码
                    </n-button>
                  </n-form-item>
                </div>

                <!-- 第二行：代理路径 + 随机生成 -->
                <div class="form-row">
                  <n-form-item label="代理路径" label-placement="top" size="small" class="form-item-half">
                    <n-input
                      v-model:value="suo5Path"
                      placeholder="输入代理路径（如 /api/proxy）"
                      clearable
                      size="small"
                    />
                    <n-text depth="3" style="font-size: 12px; margin-top: 4px;">
                      路径格式: /api/xxx 或 /path/to/proxy
                    </n-text>
                  </n-form-item>
                  <n-form-item label=" " label-placement="top" size="small" class="form-item-half">
                    <n-button
                      secondary
                      size="small"
                      block
                      @click="generateSuo5Path"
                    >
                      随机生成路径
                    </n-button>
                  </n-form-item>
                </div>

                <!-- 第三行：超时时间 -->
                <div class="form-row">
                  <n-form-item label="超时时间（秒）" label-placement="top" size="small" class="form-item-half">
                    <n-input-number
                      v-model:value="suo5Timeout"
                      :min="1"
                      :max="300"
                      size="small"
                    />
                    <n-text depth="3" style="font-size: 12px; margin-top: 4px;">
                      建议值: 30-60秒
                    </n-text>
                  </n-form-item>
                  <div class="form-item-half"></div>
                </div>
              </n-collapse-item>
            </n-collapse>
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
              {{ $t('payloads.generate') }}
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
              {{ $t('payloads.realTimePreview') }}
            </n-button>
          </div>
        </n-card>
      </div>

      <!-- 右侧预览区 -->
      <div class="preview-panel">
        <n-card :title="$t('payloads.codePreview')" :bordered="false" size="medium">
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
                {{ $t('payloads.copyCode') }}
              </n-tooltip>

              <n-tooltip>
                <template #trigger>
                  <n-button quaternary circle @click="downloadFile">
                    <template #icon>
                      <IconDownload />
                    </template>
                  </n-button>
                </template>
                {{ $t('payloads.downloadFile') }}
              </n-tooltip>

              <n-tooltip v-if="isAdvancedMode && generatedResult?.client_config">
                <template #trigger>
                  <n-button quaternary circle @click="exportConfig">
                    <template #icon>
                      <IconSettings />
                    </template>
                  </n-button>
                </template>
                {{ $t('payloads.exportConfig') }}
              </n-tooltip>
            </n-space>
          </template>

          <!-- 加载状态 -->
          <n-spin :show="isGenerating">
            <template #description>
              {{ $t('payloads.generating') }}
            </template>

            <!-- 代码预览区 -->
            <div class="code-preview-container">
              <MonacoEditor
                v-if="generatedResult?.code"
                :value="generatedResult.code"
                :language="getScriptLanguage(generatedResult.filename)"
                :read-only="true"
              />
              
              <n-empty
                v-else
                :description="$t('payloads.clickToPreview')"
                style="padding: 60px 20px;"
              />
            </div>
          </n-spin>

          <!-- 状态信息 -->
          <div v-if="generatedResult" class="status-bar">
            <n-space justify="space-between">
              <n-text depth="3">
                文件名: {{ generatedResult.filename }}
              </n-text>
              <n-text depth="3">
                大小: {{ formatSize(generatedResult.size) }}
              </n-text>
            </n-space>
          </div>

          <!-- Suo5 客户端命令 -->
          <div v-if="clientCommand" class="client-command-panel">
            <n-card title="Suo5 客户端命令" size="small">
              <div class="command-info">
                <n-text depth="3" style="font-size: 12px; margin-bottom: 8px;">
                  复制以下命令到终端执行，连接到生成的 Suo5 代理：
                </n-text>
              </div>
              <pre class="command-block">{{ clientCommand }}</pre>
              <div class="command-actions">
                <n-button
                  secondary
                  size="small"
                  @click="copyClientCommand"
                >
                  复制命令
                </n-button>
              </div>
            </n-card>
          </div>

          <!-- 安全警示 -->
          <div class="security-warning">
            <n-alert
              type="warning"
              title="安全警示"
              size="small"
            >
              <div>
                <p>本工具仅供授权渗透测试与安全研究使用，严禁用于非法用途。</p>
                <p style="font-size: 12px; margin-top: 4px;">
                  使用本工具生成的载荷应在授权范围内使用，遵守相关法律法规。
                </p>
              </div>
            </n-alert>
          </div>
        </n-card>


      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineAsyncComponent } from 'vue'
import { useMessage } from 'naive-ui'
import { 
  IconCode, 
  IconEye, 
  IconCopy, 
  IconDownload, 
  IconSettings 
} from '@tabler/icons-vue'
import { usePayloadStore } from '@/stores/payload'
import type { PayloadConfig, ObfuscationLevel, FunctionType, InjectionType } from '@/types/payload'
import { FUNCTION_TYPE_OPTIONS, INJECTION_TYPE_OPTIONS } from '@/types/payload'
import FileShellConfig from './modes/FileShellConfig.vue'
// 动态导入 Monaco Editor 以优化性能
const MonacoEditor = defineAsyncComponent(() => import('@/components/shared/MonacoEditor.vue'))
import { AuditLogger } from '@/utils/auditLogger'

const message = useMessage()
const payloadStore = usePayloadStore()

// 本地状态 - 直接使用 store 中的 config
const config = computed(() => payloadStore.config)

// 当前模式（用于 v-model）
const currentMode = computed({
  get: () => config.value.mode,
  set: (mode: any) => payloadStore.setMode(mode)
})

// 混淆级别滑块值 (1-3)
const obfuscationValue = ref(1)

// 计算属性
const isFileBasedMode = computed(() => config.value.mode === 'file_based')
const isFileCommonMode = computed(() => config.value.mode === 'file_common')
const isFileProxyMode = computed(() => config.value.mode === 'file_proxy')
const isFileHybridMode = computed(() => config.value.mode === 'file_hybrid')
const isMemoryShellMode = computed(() => config.value.mode === 'memory_shell')
const isSuo5OnlyMode = computed(() => config.value.mode === 'suo5_only')
const isFileMode = computed(() => isFileBasedMode.value || isFileCommonMode.value || isFileProxyMode.value || isFileHybridMode.value)
const isGenerating = computed(() => payloadStore.isGenerating)
const generatedResult = computed(() => payloadStore.generatedResult)
const clientCommand = computed(() => payloadStore.generateClientCommand())

// Suo5 配置双向绑定
const suo5Auth = computed({
  get: () => config.value.suo5_config?.auth || '',
  set: (value) => {
    if (config.value.suo5_config) {
      config.value.suo5_config.auth = value
    }
  }
})

const suo5Path = computed({
  get: () => config.value.suo5_config?.path || '',
  set: (value) => {
    if (config.value.suo5_config) {
      config.value.suo5_config.path = value
    }
  }
})

const suo5Timeout = computed({
  get: () => config.value.suo5_config?.timeout || 30,
  set: (value) => {
    if (config.value.suo5_config) {
      config.value.suo5_config.timeout = value
    }
  }
})

// 模式切换
const onModeChange = (mode: string) => {
  // 模式切换逻辑由 store 中的 watch 处理
  // 这里可以添加额外的 UI 逻辑
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
    message.success($t('payloads.generateSuccess'))
    // 记录审计日志
    if (generatedResult.value) {
      const payloadTypeMap: Record<string, string> = {
        'file_based': '文件落地',
        'file_common': '文件通用模式',
        'file_proxy': '文件代理模式',
        'file_hybrid': '文件混合模式',
        'memory_shell': '内存马注入',
        'suo5_only': '纯 Suo5 代理'
      };
      const payloadType = payloadTypeMap[currentMode.value] || currentMode.value;
      const scriptType = config.value.script_type || 'php';
      // 生成简单的哈希作为 payload_hash
      const payloadHash = btoa(generatedResult.value.code.substring(0, 100)).substring(0, 32);
      await AuditLogger.logPayloadGenerate(payloadType, scriptType, payloadHash);
    }
  } catch (error: any) {
    message.error(error.message || $t('payloads.generateFailed'))
  }
}

// 实时预览
const handlePreview = async () => {
  payloadStore.setObfuscationLevel(mapObfuscationLevel(obfuscationValue.value))
  
  try {
    await payloadStore.generatePreview(300)
  } catch (error: any) {
    message.error(error.message || $t('payloads.previewFailed'))
  }
}

// 复制代码
const copyCode = async () => {
  if (!generatedResult.value?.code) return
  
  try {
    await payloadStore.copyCode(generatedResult.value.code)
    message.success($t('payloads.copySuccess'))
  } catch (error: any) {
    message.error($t('payloads.copyFailed'))
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
  
  message.success($t('payloads.downloadStarted'))
}

// 导出客户端配置
const exportConfig = async () => {
  try {
    // 这里可以集成 Tauri 文件对话框
    message.info($t('payloads.exportInProgress'))
  } catch (error: any) {
    message.error($t('payloads.exportFailed'))
  }
}

// 格式化文件大小
const formatSize = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

// 生成 Suo5 随机密码
const generateSuo5Password = () => {
  const password = payloadStore.generateRandomSuo5Password()
  message.success('已生成随机密码')
}

// 生成 Suo5 随机路径
const generateSuo5Path = () => {
  const path = payloadStore.generateRandomSuo5Path()
  message.success('已生成随机路径')
}

// 复制客户端命令
const copyClientCommand = async () => {
  if (!clientCommand.value) return
  
  try {
    await navigator.clipboard.writeText(clientCommand.value)
    message.success('客户端命令已复制到剪贴板')
  } catch (error: any) {
    message.error('复制失败')
  }
}

// 密码强度检测
const getPasswordStrength = (password: string): string => {
  if (password.length < 8) return '弱'
  if (password.length < 12) return '中等'
  if (password.length < 16) return '强'
  return '极强'
}

// 根据文件名获取脚本语言
const getScriptLanguage = (filename: string): string => {
  const ext = filename.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'php':
      return 'php'
    case 'jsp':
      return 'java'
    case 'aspx':
      return 'csharp'
    case 'asp':
      return 'vb'
    default:
      return 'plaintext'
  }
}
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

  /* 客户端命令面板 */
  .client-command-panel {
    margin-top: 16px;
  }

  .command-block {
    margin: 0 0 12px 0;
    padding: 12px;
    background: var(--code-bg);
    border: 1px solid var(--border-color);
    border-radius: 6px;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 12px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-wrap: break-word;
    color: var(--code-text);
  }

  .command-info {
    margin-bottom: 8px;
  }

  .command-actions {
    display: flex;
    justify-content: flex-end;
  }

  /* 安全警示 */
  .security-warning {
    margin-top: 16px;
  }

  :deep(.n-alert) {
    border-radius: 8px;
  }

/* 分割线优化 */
:deep(.n-divider) {
  border-color: var(--border-color);
  margin: 10px 0;
}

/* 表单样式优化 - 紧凑版 */
:deep(.n-form-item) {
  margin-bottom: 12px;
}

:deep(.n-form-item-label) {
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  font-size: 13px;
  letter-spacing: 0.2px;
}

:deep(.n-form-item-blank) {
  padding: 0;
}

/* 模式选择器包装器 */
.mode-selector-wrapper {
  margin-bottom: 16px;
}

.mode-selector-label {
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  font-size: 13px;
  letter-spacing: 0.2px;
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
  padding: 6px;
  background: var(--card-bg-hover);
  border: 1px solid var(--border-color);
  margin-bottom: 12px;
}

:deep(.n-segmented-item) {
  padding: 10px 12px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 6px;
  transition: all 0.2s ease;
  min-width: 80px;
}

:deep(.n-segmented-item:hover) {
  background: var(--active-color-bg);
}

:deep(.n-segmented-item--selected) {
  background: var(--active-color) !important;
  color: white !important;
  box-shadow: 0 2px 8px rgba(var(--active-color-rgb), 0.4);
  font-weight: 600;
}

:deep(.n-segmented-item__label) {
  color: inherit;
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
