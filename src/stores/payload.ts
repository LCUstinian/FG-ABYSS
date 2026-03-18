/**
 * 载荷生成模块状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { 
  PayloadConfig, 
  PayloadResult, 
  PayloadMode,
  ScriptType,
  FunctionType,
  EncodeType,
  EncryptAlgo,
  ObfuscationLevel 
} from '@/types/payload'
import { invoke } from '@tauri-apps/api/core'

export const usePayloadStore = defineStore('payload', () => {
  // 状态
  const config = ref<PayloadConfig>({
    mode: 'simple',
    script_type: 'asp',
    function_type: 'basic',
    password: '',
    encode_type: 'none',
    encrypt_algo: undefined,
    obfuscation_level: 'low',
    output_filename: undefined,
    template_name: undefined,
  })

  const generatedResult = ref<PayloadResult | null>(null)
  const isGenerating = ref(false)
  const error = ref<string | null>(null)
  const history = ref<PayloadResult[]>([])

  // 计算属性
  const isSimpleMode = computed(() => config.value.mode === 'simple')
  const isAdvancedMode = computed(() => config.value.mode === 'advanced')

  // 模式切换时的联动逻辑
  watch(
    () => config.value.mode,
    (newMode: PayloadMode) => {
      if (newMode === 'advanced') {
        // 切换到高级模式：清空并禁用 encode_type，启用 encrypt_algo
        config.value.encode_type = undefined
        config.value.encrypt_algo = 'aes128_cbc'
      } else {
        // 切换到简单模式：清空并禁用 encrypt_algo，启用 encode_type
        config.value.encrypt_algo = undefined
        config.value.encode_type = 'none'
      }
    }
  )

  // 方法
  const setMode = (mode: PayloadMode) => {
    config.value.mode = mode
  }

  const setScriptType = (type: ScriptType) => {
    config.value.script_type = type
  }

  const setFunctionType = (funcType: FunctionType) => {
    config.value.function_type = funcType
  }

  const setPassword = (password: string) => {
    config.value.password = password
  }

  const setEncodeType = (encodeType: EncodeType) => {
    if (isSimpleMode.value) {
      config.value.encode_type = encodeType === 'none' ? undefined : encodeType
    }
  }

  const setEncryptAlgo = (algo: EncryptAlgo) => {
    if (isAdvancedMode.value) {
      config.value.encrypt_algo = algo
    }
  }

  const setObfuscationLevel = (level: ObfuscationLevel) => {
    config.value.obfuscation_level = level
  }

  const setOutputFilename = (filename: string) => {
    config.value.output_filename = filename || undefined
  }

  // 生成预览 (带防抖)
  let debounceTimer: NodeJS.Timeout | null = null

  const generatePreview = async (delay = 500) => {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }

    return new Promise<void>((resolve, reject) => {
      debounceTimer = setTimeout(async () => {
        await generate()
        resolve()
      }, delay)
    })
  }

  // 生成载荷
  const generate = async () => {
    isGenerating.value = true
    error.value = null

    try {
      // 验证输入
      if (!config.value.password || config.value.password.trim() === '') {
        throw new Error('密码不能为空')
      }

      // 调用 Tauri 命令
      const result = await invoke<PayloadResult>('generate_payload_cmd', {
        config: config.value,
      })

      generatedResult.value = result
      
      // 添加到历史记录
      history.value.unshift(result)
      
      // 限制历史记录数量
      if (history.value.length > 50) {
        history.value.pop()
      }

      return result
    } catch (err: any) {
      error.value = err.message || '生成失败'
      throw err
    } finally {
      isGenerating.value = false
    }
  }

  // 获取历史记录
  const loadHistory = async () => {
    try {
      history.value = await invoke<PayloadResult[]>('get_generated_payloads')
    } catch (err: any) {
      console.error('Failed to load history:', err)
    }
  }

  // 保存文件
  const saveToFile = async (path: string, content: string) => {
    try {
      await invoke('save_file_cmd', {
        path,
        content,
      })
    } catch (err: any) {
      throw new Error(`保存文件失败：${err.message}`)
    }
  }

  // 导出客户端配置
  const exportClientConfig = async (path: string) => {
    if (!generatedResult.value?.client_config) {
      throw new Error('没有可导出的客户端配置')
    }

    try {
      await invoke('export_client_config_cmd', {
        config: generatedResult.value.client_config,
        path,
      })
    } catch (err: any) {
      throw new Error(`导出配置失败：${err.message}`)
    }
  }

  // 清空历史
  const clearHistory = () => {
    history.value = []
  }

  // 复制代码到剪贴板
  const copyCode = async (code: string) => {
    try {
      await navigator.clipboard.writeText(code)
    } catch (err: any) {
      throw new Error(`复制失败：${err.message}`)
    }
  }

  return {
    // 状态
    config,
    generatedResult,
    isGenerating,
    error,
    history,
    
    // 计算属性
    isSimpleMode,
    isAdvancedMode,
    
    // 方法
    setMode,
    setScriptType,
    setFunctionType,
    setPassword,
    setEncodeType,
    setEncryptAlgo,
    setObfuscationLevel,
    setOutputFilename,
    generate,
    generatePreview,
    loadHistory,
    saveToFile,
    exportClientConfig,
    clearHistory,
    copyCode,
  }
})
