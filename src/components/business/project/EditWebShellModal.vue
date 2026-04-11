<template>
  <n-modal
    v-model:show="visible"
    preset="card"
    title="编辑 WebShell"
    size="large"
    :bordered="false"
    :close-on-esc="false"
  >
    <n-form
      v-model:model="formData"
      :rules="rules"
      :label-width="100"
    >
      <n-form-item label="URL" path="url">
        <n-input
          v-model:value="formData.url"
          placeholder="请输入 WebShell URL"
          size="medium"
        />
      </n-form-item>
      
      <n-form-item label="载荷类型" path="payload">
        <n-select
          v-model:value="formData.payload"
          :options="payloadOptions"
          placeholder="请选择载荷类型"
          size="medium"
        />
      </n-form-item>
      
      <n-form-item label="加密方式" path="cryption">
        <n-select
          v-model:value="formData.cryption"
          :options="cryptionOptions"
          placeholder="请选择加密方式"
          size="medium"
        />
      </n-form-item>
      
      <n-form-item label="编码方式" path="encoding">
        <n-select
          v-model:value="formData.encoding"
          :options="encodingOptions"
          placeholder="请选择编码方式"
          size="medium"
        />
      </n-form-item>
      
      <n-form-item label="代理类型" path="proxyType">
        <n-select
          v-model:value="formData.proxyType"
          :options="proxyOptions"
          placeholder="请选择代理类型"
          size="medium"
        />
      </n-form-item>
      
      <n-form-item label="备注" path="remark">
        <n-input
          v-model:value="formData.remark"
          type="textarea"
          placeholder="请输入备注信息"
          size="medium"
          :autosize="{ minRows: 3, maxRows: 5 }"
        />
      </n-form-item>
    </n-form>
    
    <template #footer>
      <div class="modal-footer">
        <n-button @click="handleCancel">取消</n-button>
        <n-button type="primary" @click="handleSubmit" :loading="loading">保存</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { NModal, NForm, NFormItem, NInput, NSelect, NButton } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { invoke } from '@/utils/tauri-mock-adapter'

const props = defineProps<{
  modelValue: boolean
  webshellId?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'updated': []
}>()

const message = useMessage()
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const loading = ref(false)
const formData = reactive({
  url: '',
  payload: '',
  cryption: '',
  encoding: '',
  proxyType: '',
  remark: ''
})

const rules = {
  url: [
    { required: true, message: '请输入 URL', trigger: 'blur' }
  ],
  payload: [
    { required: true, message: '请选择载荷类型', trigger: 'change' }
  ],
  cryption: [
    { required: true, message: '请选择加密方式', trigger: 'change' }
  ],
  encoding: [
    { required: true, message: '请选择编码方式', trigger: 'change' }
  ]
}

const payloadOptions = [
  { label: 'PHP', value: 'php' },
  { label: 'JSP', value: 'jsp' },
  { label: 'ASPX', value: 'aspx' },
  { label: 'ASP', value: 'asp' }
]

const cryptionOptions = [
  { label: 'AES-256-GCM', value: 'aes-256-gcm' },
  { label: 'AES-256-CBC', value: 'aes-256-cbc' },
  { label: 'RC4', value: 'rc4' }
]

const encodingOptions = [
  { label: 'Base64', value: 'base64' },
  { label: 'Hex', value: 'hex' },
  { label: 'URL', value: 'url' }
]

const proxyOptions = [
  { label: '无代理', value: 'none' },
  { label: 'HTTP', value: 'http' },
  { label: 'SOCKS5', value: 'socks5' }
]

const loadWebShellData = async () => {
  if (!props.webshellId) return
  
  try {
    const webshell = await invoke('get_webshell', { id: props.webshellId })
    if (webshell) {
      formData.url = webshell.url || ''
      formData.payload = webshell.payload || ''
      formData.cryption = webshell.cryption || ''
      formData.encoding = webshell.encoding || ''
      formData.proxyType = webshell.proxyType || ''
      formData.remark = webshell.remark || ''
    }
  } catch (error) {
    console.error('加载 WebShell 数据失败:', error)
    message.error('加载数据失败')
  }
}

watch(() => props.webshellId, () => {
  if (props.webshellId) {
    loadWebShellData()
  }
})

watch(() => props.modelValue, (newValue) => {
  if (newValue && props.webshellId) {
    loadWebShellData()
  }
})

const handleCancel = () => {
  visible.value = false
}

const handleSubmit = async () => {
  if (!props.webshellId) return
  
  loading.value = true
  try {
    const result = await invoke('update_webshell', {
      id: props.webshellId,
      ...formData
    })
    
    if (result.success) {
      message.success('编辑成功')
      visible.value = false
      emit('updated')
    } else {
      throw new Error(result.message || '编辑失败')
    }
  } catch (error: any) {
    console.error('编辑失败:', error)
    message.error(error.message || '编辑失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}
</style>