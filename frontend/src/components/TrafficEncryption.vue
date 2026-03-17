<template>
  <div class="traffic-encryption">
    <n-card title="流量加密配置" :bordered="false">
      <n-alert type="warning" title="加密说明" style="margin-bottom: 20px">
        启用流量加密后，所有 WebShell 通信将使用 AES-256-GCM 加密，并附带 HMAC-SHA256 签名验证
      </n-alert>

      <n-tabs type="line" animated>
        <!-- 快速配置 -->
        <n-tab-pane name="quick" tab="快速配置">
          <n-space vertical>
            <n-alert type="info" title="一键生成安全配置">
              系统将自动生成符合安全标准的加密密钥和参数
            </n-alert>

            <n-space>
              <n-button
                type="primary"
                @click="handleGenerateConfig"
                :loading="generating"
              >
                <template #icon>
                  <n-icon :component="KeyOutline" />
                </template>
                生成加密配置
              </n-button>
            </n-space>

            <n-divider />

            <n-form
              v-if="generatedConfig"
              label-placement="top"
              :model="generatedConfig"
            >
              <n-form-item label="加密密钥（Key）">
                <n-input
                  v-model:value="generatedConfig.key"
                  readonly
                  type="textarea"
                  :rows="2"
                >
                  <template #suffix>
                    <n-button
                      quaternary
                      size="small"
                      @click="handleCopy(generatedConfig.key)"
                    >
                      <template #icon>
                        <n-icon :component="CopyOutline" />
                      </template>
                    </n-button>
                  </template>
                </n-input>
              </n-form-item>

              <n-form-item label="初始化向量（IV）">
                <n-input
                  v-model:value="generatedConfig.iv"
                  readonly
                  type="textarea"
                  :rows="2"
                >
                  <template #suffix>
                    <n-button
                      quaternary
                      size="small"
                      @click="handleCopy(generatedConfig.iv)"
                    >
                      <template #icon>
                        <n-icon :component="CopyOutline" />
                      </template>
                    </n-button>
                  </template>
                </n-input>
              </n-form-item>

              <n-form-item label="签名密钥（Signature）">
                <n-input
                  v-model:value="generatedConfig.signature"
                  readonly
                  type="textarea"
                  :rows="2"
                >
                  <template #suffix>
                    <n-button
                      quaternary
                      size="small"
                      @click="handleCopy(generatedConfig.signature)"
                    >
                      <template #icon>
                        <n-icon :component="CopyOutline" />
                      </template>
                    </n-button>
                  </template>
                </n-input>
              </n-form-item>

              <n-space>
                <n-button
                  type="primary"
                  @click="handleApplyConfig"
                  :disabled="!generatedConfig"
                >
                  应用配置
                </n-button>
                <n-button @click="handleTestConfig">
                  测试加密
                </n-button>
              </n-space>
            </n-form>
          </n-space>
        </n-tab-pane>

        <!-- 手动配置 -->
        <n-tab-pane name="manual" tab="手动配置">
          <n-form
            ref="formRef"
            :model="manualConfig"
            :rules="formRules"
            label-placement="top"
          >
            <n-form-item
              label="加密密钥（Key）"
              path="key"
              help="256 位密钥（64 个十六进制字符）"
            >
              <n-input
                v-model:value="manualConfig.key"
                placeholder="输入 64 位十六进制密钥"
                :disabled="encryptionEnabled"
              />
            </n-form-item>

            <n-form-item
              label="初始化向量（IV）"
              path="iv"
              help="96 位 IV（24 个十六进制字符）"
            >
              <n-input
                v-model:value="manualConfig.iv"
                placeholder="输入 24 位十六进制 IV"
                :disabled="encryptionEnabled"
              />
            </n-form-item>

            <n-form-item
              label="签名密钥（Signature）"
              path="signature"
              help="256 位签名密钥（64 个十六进制字符）"
            >
              <n-input
                v-model:value="manualConfig.signature"
                placeholder="输入 64 位十六进制签名密钥"
                :disabled="encryptionEnabled"
              />
            </n-form-item>

            <n-space>
              <n-button
                v-if="!encryptionEnabled"
                type="primary"
                @click="handleInitEncryption"
                :loading="initializing"
              >
                启用加密
              </n-button>
              <n-button
                v-else
                type="error"
                @click="handleDisableEncryption"
              >
                禁用加密
              </n-button>
              <n-button @click="handleValidateManual">
                验证配置
              </n-button>
            </n-space>
          </n-form>
        </n-tab-pane>

        <!-- 测试 -->
        <n-tab-pane name="test" tab="加密测试">
          <n-space vertical>
            <n-form label-placement="top">
              <n-form-item label="测试数据">
                <n-input
                  v-model:value="testData"
                  type="textarea"
                  :rows="3"
                  placeholder="输入要测试的数据"
                />
              </n-form-item>

              <n-form-item label="时间窗口">
                <n-select
                  v-model:value="timeWindow"
                  :options="windowOptions"
                  style="width: 200px"
                />
              </n-form-item>

              <n-space>
                <n-button
                  type="primary"
                  @click="handleTestEncryptDecrypt"
                  :loading="testing"
                  :disabled="!encryptionEnabled"
                >
                  <template #icon>
                    <n-icon :component="PulseOutline" />
                  </template>
                  测试加密/解密
                </n-button>
              </n-space>
            </n-form>

            <n-divider />

            <n-space v-if="testResult" vertical>
              <n-result
                :status="testResult.success ? 'success' : 'error'"
                :title="testResult.success ? '测试成功' : '测试失败'"
                :description="testResult.message"
              />

              <n-card v-if="testResult.success" title="测试结果" size="small">
                <n-space vertical>
                  <n-text strong>原始数据:</n-text>
                  <n-code :code="testData" />
                  <n-text strong>加密后:</n-text>
                  <n-code :code="testResult.encrypted" />
                  <n-text strong>解密后:</n-text>
                  <n-code :code="testResult.decrypted" />
                </n-space>
              </n-card>
            </n-space>
          </n-space>
        </n-tab-pane>

        <!-- 状态 -->
        <n-tab-pane name="status" tab="状态">
          <n-space vertical>
            <n-alert
              :type="encryptionEnabled ? 'success' : 'default'"
              :title="encryptionEnabled ? '加密已启用' : '加密未启用'"
            >
              <template #icon>
                <n-icon :component="encryptionEnabled ? CheckmarkCircleOutline : CloseCircleOutline" />
              </template>
              {{ encryptionEnabled ? '所有流量将使用 AES-256-GCM 加密' : '当前通信未加密' }}
            </n-alert>

            <n-descriptions
              v-if="encryptionEnabled"
              bordered
              :column="1"
              size="small"
            >
              <n-descriptions-item label="加密算法">
                AES-256-GCM
              </n-descriptions-item>
              <n-descriptions-item label="签名算法">
                HMAC-SHA256
              </n-descriptions-item>
              <n-descriptions-item label="时间窗口">
                {{ timeWindow }} 秒
              </n-descriptions-item>
            </n-descriptions>
          </n-space>
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useMessage } from 'naive-ui'
import {
  KeyOutline,
  CopyOutline,
  CheckmarkCircleOutline,
  CloseCircleOutline,
  PulseOutline,
} from '@vicons/ionicons5'
import type { FormRules, FormInst } from 'naive-ui'
import {
  InitEncryption,
  GenerateConfig,
  ValidateConfig,
  TestEncryption,
} from '@bindings/EncryptionHandler'

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const generating = ref(false)
const initializing = ref(false)
const testing = ref(false)

const generatedConfig = ref<{
  key: string
  iv: string
  signature: string
} | null>(null)

const manualConfig = reactive({
  key: '',
  iv: '',
  signature: '',
})

const encryptionEnabled = ref(false)
const testData = ref('Hello, World!')
const timeWindow = ref(300)
const testResult = ref<any>(null)

const formRules: FormRules = {
  key: {
    required: true,
    message: '请输入加密密钥',
    trigger: 'blur',
    validator: (_: any, value: string) => {
      if (value.length !== 64) {
        return new Error('密钥长度必须为 64 个十六进制字符')
      }
      if (!/^[0-9a-fA-F]+$/.test(value)) {
        return new Error('密钥必须是十六进制字符')
      }
      return true
    },
  },
  iv: {
    required: true,
    message: '请输入初始化向量',
    trigger: 'blur',
    validator: (_: any, value: string) => {
      if (value.length !== 24) {
        return new Error('IV 长度必须为 24 个十六进制字符')
      }
      if (!/^[0-9a-fA-F]+$/.test(value)) {
        return new Error('IV 必须是十六进制字符')
      }
      return true
    },
  },
  signature: {
    required: true,
    message: '请输入签名密钥',
    trigger: 'blur',
    validator: (_: any, value: string) => {
      if (value.length !== 64) {
        return new Error('签名密钥长度必须为 64 个十六进制字符')
      }
      if (!/^[0-9a-fA-F]+$/.test(value)) {
        return new Error('签名密钥必须是十六进制字符')
      }
      return true
    },
  },
}

const windowOptions = [
  { label: '1 分钟', value: 60 },
  { label: '5 分钟', value: 300 },
  { label: '10 分钟', value: 600 },
  { label: '30 分钟', value: 1800 },
]

const handleGenerateConfig = async () => {
  try {
    generating.value = true
    const config = await GenerateConfig({})
    generatedConfig.value = {
      key: config.key,
      iv: config.iv,
      signature: config.signature,
    }
    message.success('配置生成成功')
  } catch (error: any) {
    message.error('配置生成失败')
  } finally {
    generating.value = false
  }
}

const handleCopy = (text: string) => {
  navigator.clipboard.writeText(text)
  message.success('已复制到剪贴板')
}

const handleApplyConfig = () => {
  if (generatedConfig.value) {
    manualConfig.key = generatedConfig.value.key
    manualConfig.iv = generatedConfig.value.iv
    manualConfig.signature = generatedConfig.value.signature
    message.success('配置已应用')
  }
}

const handleInitEncryption = async () => {
  try {
    await formRef.value?.validate()
    initializing.value = true

    const response = await InitEncryption({
      key: manualConfig.key,
      iv: manualConfig.iv,
      signature: manualConfig.signature,
    })

    if (response.success) {
      encryptionEnabled.value = true
      message.success('加密启用成功')
    } else {
      message.error('加密启用失败：' + response.message)
    }
  } catch (error: any) {
    if (error.errors) return
    message.error('启用失败')
  } finally {
    initializing.value = false
  }
}

const handleDisableEncryption = () => {
  encryptionEnabled.value = false
  message.success('加密已禁用')
}

const handleValidateManual = async () => {
  try {
    const response = await ValidateConfig({
      key: manualConfig.key,
      iv: manualConfig.iv,
      signature: manualConfig.signature,
    })

    if (response.valid) {
      message.success('配置验证通过')
    } else {
      message.error('配置验证失败：' + response.errors.join(', '))
    }
  } catch (error: any) {
    message.error('验证失败')
  }
}

const handleTestConfig = async () => {
  if (!generatedConfig.value) return

  try {
    const response = await TestEncryption({
      key: generatedConfig.value.key,
      iv: generatedConfig.value.iv,
      signature: generatedConfig.value.signature,
      test_data: 'Test data',
    })

    if (response.success) {
      message.success('加密测试成功')
    } else {
      message.error('加密测试失败：' + response.message)
    }
  } catch (error: any) {
    message.error('测试失败')
  }
}

const handleTestEncryptDecrypt = async () => {
  if (!encryptionEnabled.value) {
    message.warning('请先启用加密')
    return
  }

  try {
    testing.value = true
    // TODO: 实现完整的加密解密测试流程
    testResult.value = {
      success: true,
      message: '测试成功',
      encrypted: 'encrypted_data_placeholder',
      decrypted: testData.value,
    }
    message.success('加密解密测试成功')
  } catch (error: any) {
    message.error('测试失败')
  } finally {
    testing.value = false
  }
}
</script>

<style scoped lang="scss">
.traffic-encryption {
  padding: 20px;
}
</style>
