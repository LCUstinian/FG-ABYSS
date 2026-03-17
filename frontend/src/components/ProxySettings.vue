<template>
  <div class="proxy-settings">
    <n-card title="代理设置" :bordered="false">
      <n-alert type="info" title="代理配置" style="margin-bottom: 20px">
        配置 HTTP/HTTPS/SOCKS5 代理，用于 WebShell 连接的流量转发
      </n-alert>

      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-placement="top"
        label-width="100"
      >
        <n-form-item label="代理类型">
          <n-select
            v-model:value="formData.type"
            :options="typeOptions"
            placeholder="选择代理类型"
            style="width: 200px"
          />
        </n-form-item>

        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="代理主机" path="host">
              <n-input
                v-model:value="formData.host"
                placeholder="127.0.0.1"
                clearable
              >
                <template #prefix>
                  <n-icon :component="ServerOutline" />
                </template>
              </n-input>
            </n-form-item>
          </n-grid-item>

          <n-grid-item>
            <n-form-item label="代理端口" path="port">
              <n-input-number
                v-model:value="formData.port"
                placeholder="8080"
                :min="1"
                :max="65535"
                style="width: 100%"
              />
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="用户名">
              <n-input
                v-model:value="formData.username"
                placeholder="可选"
                clearable
              >
                <template #prefix>
                  <n-icon :component="PersonOutline" />
                </template>
              </n-input>
            </n-form-item>
          </n-grid-item>

          <n-grid-item>
            <n-form-item label="密码">
              <n-input
                v-model:value="formData.password"
                type="password"
                placeholder="可选"
                show-password-on="click"
              >
                <template #prefix>
                  <n-icon :component="LockClosedOutline" />
                </template>
              </n-input>
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-form-item label="超时时间">
          <n-input-number
            v-model:value="formData.timeout"
            :min="1"
            :max="300"
            placeholder="30"
            style="width: 200px"
          >
            <template #suffix>
              秒
            </template>
          </n-input-number>
        </n-form-item>

        <n-divider />

        <n-space>
          <n-button
            type="primary"
            @click="handleSetProxy"
            :loading="setting"
            :disabled="!formData.host || !formData.port"
          >
            <template #icon>
              <n-icon :component="CheckmarkCircleOutline" />
            </template>
            应用配置
          </n-button>

          <n-button
            @click="handleTestProxy"
            :loading="testing"
            :disabled="!formData.host || !formData.port"
          >
            <template #icon>
              <n-icon :component="PulseOutline" />
            </template>
            测试连接
          </n-button>

          <n-button
            @click="handleDisableProxy"
            :disabled="!proxyStatus.enabled"
          >
            <template #icon>
              <n-icon :component="CloseCircleOutline" />
            </template>
            禁用代理
          </n-button>

          <n-button @click="handleReset">
            <template #icon>
              <n-icon :component="RefreshOutline" />
            </template>
            重置
          </n-button>
        </n-space>
      </n-form>

      <n-divider />

      <!-- 代理状态 -->
      <n-space vertical>
        <n-text strong>当前代理状态</n-text>
        
        <n-alert
          v-if="proxyStatus.enabled"
          type="success"
          title="代理已启用"
          style="max-width: 400px"
        >
          <template #icon>
            <n-icon :component="CheckmarkCircleOutline" />
          </template>
          {{ proxyStatus.type }}://{{ proxyStatus.host }}:{{ proxyStatus.port }}
          <n-tag v-if="proxyStatus.hasAuth" type="warning" size="small" style="margin-left: 8px">
            需要认证
          </n-tag>
        </n-alert>

        <n-alert
          v-else
          type="default"
          title="代理未启用"
          style="max-width: 400px"
        >
          当前未使用代理，直接连接目标服务器
        </n-alert>
      </n-space>

      <!-- 预设代理 -->
      <n-divider />

      <n-space vertical>
        <n-text strong>预设代理配置</n-text>
        <n-space>
          <n-button
            size="small"
            @click="handleLoadPreset('tor')"
            quaternary
          >
            Tor (9050)
          </n-button>
          <n-button
            size="small"
            @click="handleLoadPreset('burp')"
            quaternary
          >
            Burp Suite (8080)
          </n-button>
          <n-button
            size="small"
            @click="handleLoadPreset('owasp')"
            quaternary
          >
            OWASP ZAP (8090)
          </n-button>
        </n-space>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  ServerOutline,
  PersonOutline,
  LockClosedOutline,
  CheckmarkCircleOutline,
  PulseOutline,
  CloseCircleOutline,
  RefreshOutline,
} from '@vicons/ionicons5'
import type { FormRules, FormInst } from 'naive-ui'
import {
  SetProxy,
  TestProxy,
  GetProxyStatus,
  DisableProxy,
  ValidateProxy,
} from '@bindings/ProxyHandler'

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const setting = ref(false)
const testing = ref(false)

const formData = reactive({
  type: 'http',
  host: '',
  port: 0,
  username: '',
  password: '',
  timeout: 30,
})

const formRules: FormRules = {
  type: {
    required: true,
    message: '请选择代理类型',
    trigger: 'change',
  },
  host: {
    required: true,
    message: '请输入代理主机',
    trigger: 'blur',
  },
  port: {
    required: true,
    message: '请输入代理端口',
    trigger: 'blur',
  },
}

const proxyStatus = reactive({
  enabled: false,
  type: '',
  host: '',
  port: 0,
  hasAuth: false,
})

const typeOptions = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'SOCKS5', value: 'socks5' },
]

const handleSetProxy = async () => {
  try {
    await formRef.value?.validate()
    setting.value = true

    const response = await SetProxy({
      type: formData.type,
      host: formData.host,
      port: formData.port,
      username: formData.username,
      password: formData.password,
      timeout: formData.timeout,
    })

    if (response.success) {
      message.success('代理配置成功')
      loadProxyStatus()
    } else {
      message.error('代理配置失败：' + response.message)
    }
  } catch (error: any) {
    if (error.errors) return
    message.error('配置失败：' + (error.message || '未知错误'))
  } finally {
    setting.value = false
  }
}

const handleTestProxy = async () => {
  try {
    await formRef.value?.validate()
    testing.value = true

    const response = await TestProxy({
      type: formData.type,
      host: formData.host,
      port: formData.port,
      username: formData.username,
      password: formData.password,
      timeout: formData.timeout || 10,
    })

    if (response.success) {
      message.success('代理连接测试成功')
    } else {
      message.error('代理连接测试失败：' + response.message)
    }
  } catch (error: any) {
    if (error.errors) return
    message.error('测试失败：' + (error.message || '未知错误'))
  } finally {
    testing.value = false
  }
}

const handleDisableProxy = async () => {
  try {
    const response = await DisableProxy()
    if (response.success) {
      message.success('代理已禁用')
      loadProxyStatus()
    } else {
      message.error('禁用代理失败')
    }
  } catch (error: any) {
    message.error('禁用代理失败')
  }
}

const handleReset = () => {
  formData.type = 'http'
  formData.host = ''
  formData.port = 0
  formData.username = ''
  formData.password = ''
  formData.timeout = 30
}

const handleLoadPreset = (preset: string) => {
  switch (preset) {
    case 'tor':
      formData.type = 'socks5'
      formData.host = '127.0.0.1'
      formData.port = 9050
      formData.username = ''
      formData.password = ''
      formData.timeout = 30
      break
    case 'burp':
      formData.type = 'http'
      formData.host = '127.0.0.1'
      formData.port = 8080
      formData.username = ''
      formData.password = ''
      formData.timeout = 10
      break
    case 'owasp':
      formData.type = 'http'
      formData.host = '127.0.0.1'
      formData.port = 8090
      formData.username = ''
      formData.password = ''
      formData.timeout = 10
      break
  }
  message.success('已加载预设配置')
}

const loadProxyStatus = async () => {
  try {
    const status = await GetProxyStatus()
    proxyStatus.enabled = status.enabled
    proxyStatus.type = status.type
    proxyStatus.host = status.host
    proxyStatus.port = status.port
    proxyStatus.hasAuth = status.has_auth
  } catch (error) {
    console.error('加载代理状态失败:', error)
  }
}

onMounted(() => {
  loadProxyStatus()
})
</script>

<style scoped lang="scss">
.proxy-settings {
  padding: 20px;
}
</style>
