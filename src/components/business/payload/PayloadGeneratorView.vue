<template>
  <div class="payload-generator-view">
    <div class="generator-layout">
      <!-- 左侧配置区 -->
      <div class="config-panel">
        <n-card title="生成配置" :bordered="false" size="medium">
          <!-- 模式切换 -->
          <n-form-item label="生成模式" label-placement="top">
            <n-segmented
              v-model:value="config.mode"
              :options="[
                { label: '极简模式', value: 'simple', description: '编码/混淆，绕过 WAF' },
                { label: '高级加密', value: 'advanced', description: 'AES 加密，抗分析' },
              ]"
              block
              @update:value="onModeChange"
            />
          </n-form-item>

          <n-divider />

          <!-- 基础配置 -->
          <n-form-item label="脚本类型" label-placement="top">
            <n-select
              v-model:value="config.script_type"
              :options="[
                { label: 'PHP', value: 'php' },
                { label: 'JSP', value: 'jsp' },
                { label: 'ASPX', value: 'aspx' },
                { label: 'ASP', value: 'asp' },
              ]"
              placeholder="选择脚本类型"
            />
          </n-form-item>

          <n-form-item label="功能类型" label-placement="top">
            <n-select
              v-model:value="config.function_type"
              :options="FUNCTION_TYPE_OPTIONS"
              placeholder="选择功能类型"
            />
          </n-form-item>

          <n-form-item 
            label="连接密码" 
            label-placement="top"
            :rule="{
              required: true,
              message: '请输入连接密码',
              trigger: 'blur',
            }"
          >
            <n-input
              v-model:value="config.password"
              placeholder="输入连接密码"
              type="password"
              show-password-on="click"
              clearable
            />
          </n-form-item>

          <!-- Simple 模式专属选项 -->
          <template v-if="isSimpleMode">
            <n-form-item label="编码器" label-placement="top">
              <n-select
                v-model:value="config.encode_type"
                :options="[
                  { label: '无编码', value: 'none' },
                  { label: 'Base64', value: 'base64' },
                  { label: 'XOR', value: 'xor' },
                  { label: 'GZInflate', value: 'gzinflate' },
                  { label: 'Hex', value: 'hex' },
                  { label: 'URL Encode', value: 'urlencode' },
                  { label: 'ROT13', value: 'rot13' },
                ]"
                placeholder="选择编码器"
              />
              <n-text depth="3" style="font-size: 12px; margin-top: 8px;">
                编码器用于绕过 WAF 静态规则检测
              </n-text>
            </n-form-item>
          </template>

          <!-- Advanced 模式专属选项 -->
          <template v-if="isAdvancedMode">
            <n-form-item label="加密算法" label-placement="top">
              <n-select
                v-model:value="config.encrypt_algo"
                :options="[
                  { label: 'AES-128-CBC', value: 'aes128_cbc' },
                  { label: 'AES-256-CBC', value: 'aes256_cbc' },
                  { label: 'XOR', value: 'xor' },
                ]"
                placeholder="选择加密算法"
              />
              <n-alert 
                type="info" 
                title="加密说明"
                style="margin-top: 8px;"
              >
                高级加密模式包含完整的密钥协商和动态解密逻辑，可有效抗流量分析
              </n-alert>
            </n-form-item>
          </template>

          <!-- 混淆级别 -->
          <n-form-item label="混淆强度" label-placement="top">
            <n-slider
              v-model:value="obfuscationValue"
              :marks="{
                1: '低',
                2: '中',
                3: '高',
              }"
              :step="1"
              :min="1"
              :max="3"
              show-tooltip
            />
          </n-form-item>

          <!-- 输出文件名 -->
          <n-form-item label="输出文件名" label-placement="top">
            <n-input
              v-model:value="config.output_filename"
              placeholder="留空则自动生成"
              clearable
            />
          </n-form-item>

          <!-- 生成按钮 -->
          <n-space vertical style="margin-top: 24px;">
            <n-button
              type="primary"
              size="large"
              block
              :loading="isGenerating"
              @click="handleGenerate"
            >
              <template #icon>
                <IconCode />
              </template>
              生成载荷
            </n-button>

            <n-button
              secondary
              block
              @click="handlePreview"
            >
              <template #icon>
                <IconEye />
              </template>
              实时预览
            </n-button>
          </n-space>
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
  height: 100%;
}

.generator-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  max-width: 1600px;
  margin: 0 auto;
}

@media (max-width: 1200px) {
  .generator-layout {
    grid-template-columns: 1fr;
  }
}

.config-panel,
.preview-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.code-preview-container {
  background: var(--code-bg, var(--card-bg-hover));
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 20px;
  min-height: 400px;
  max-height: 600px;
  overflow: auto;
}

.code-block {
  margin: 0;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--code-text);
}

.status-bar {
  margin-top: 16px;
  padding: 12px;
  background: var(--card-bg-hover);
  border-radius: 8px;
  border-top: 1px solid var(--border-color);
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
  padding: 8px 12px;
  background: var(--card-bg-hover);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.history-item:hover {
  background: var(--active-color-bg);
  transform: translateX(4px);
}
</style>
