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
  ObfuscationLevel,
  InjectionType,
  Suo5Config
} from '@/types/payload'
import { invoke } from '@tauri-apps/api/core'

export const usePayloadStore = defineStore('payload', () => {
  // 状态
  const config = ref<PayloadConfig>({
    mode: 'file_based',
    script_type: 'php',
    function_type: 'basic',
    password: '',
    encode_type: 'none',
    encrypt_algo: 'aes256_gcm',
    obfuscation_level: 'medium',
    output_filename: undefined,
    template_name: undefined,
    injection_type: undefined,
    suo5_config: {
      auth: '',
      path: '',
      timeout: 30
    },
    self_destruct: false
  })

  const generatedResult = ref<PayloadResult | null>(null)
  const isGenerating = ref(false)
  const error = ref<string | null>(null)
  const history = ref<PayloadResult[]>([])

  // 计算属性 - 多模式状态管理
  const isSimpleMode = computed(() => config.value.mode === 'simple')
  const isAdvancedMode = computed(() => config.value.mode === 'advanced')
  const isFileBasedMode = computed(() => config.value.mode === 'file_based')
  const isMemoryShellMode = computed(() => config.value.mode === 'memory_shell')
  const isSuo5OnlyMode = computed(() => config.value.mode === 'suo5_only')
  const isFileCommonMode = computed(() => config.value.mode === 'file_common')
  const isFileProxyMode = computed(() => config.value.mode === 'file_proxy')
  const isFileHybridMode = computed(() => config.value.mode === 'file_hybrid')
  const isShellMode = computed(() => isSimpleMode.value || isAdvancedMode.value || isFileBasedMode.value || isMemoryShellMode.value)
  const isFileMode = computed(() => isFileBasedMode.value || isFileCommonMode.value || isFileProxyMode.value || isFileHybridMode.value)

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
    config.value.encode_type = encodeType === 'none' ? undefined : encodeType
  }

  const setEncryptAlgo = (algo: EncryptAlgo) => {
    config.value.encrypt_algo = algo
  }

  const setObfuscationLevel = (level: ObfuscationLevel) => {
    config.value.obfuscation_level = level
  }

  const setOutputFilename = (filename: string) => {
    config.value.output_filename = filename || undefined
  }

  const setInjectionType = (injectionType: InjectionType | undefined) => {
    config.value.injection_type = injectionType
  }

  // Suo5 配置项管理
  const setSuo5Auth = (auth: string) => {
    if (config.value.suo5_config) {
      config.value.suo5_config.auth = auth
    }
  }

  const setSuo5Path = (path: string) => {
    if (config.value.suo5_config) {
      config.value.suo5_config.path = path
    }
  }

  const setSuo5Timeout = (timeout: number) => {
    if (config.value.suo5_config) {
      config.value.suo5_config.timeout = timeout
    }
  }

  // 自毁逻辑开关
  const setSelfDestruct = (enabled: boolean) => {
    config.value.self_destruct = enabled
  }

  // 随机生成字符串辅助函数
  const generateRandomString = (length: number, chars: string): string => {
    let result = ''
    for (let i = 0; i < length; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return result
  }

  // Suo5 密码随机生成
  const generateRandomSuo5Password = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
    const password = generateRandomString(16, chars)
    if (config.value.suo5_config) {
      config.value.suo5_config.auth = password
    }
    return password
  }

  // Suo5 路径随机生成
  const generateRandomSuo5Path = () => {
    const chars = 'abcdefghijklmnopqrstuvwxyz'
    const prefixes = ['api', 'admin', 'system', 'assets', 'resources', 'static', 'data', 'files', 'uploads', 'media']
    const suffixes = ['v1', 'v2', 'v3', 'internal', 'private', 'secure', 'backend', 'core']
    
    const prefix = prefixes[Math.floor(Math.random() * prefixes.length)]
    const suffix = suffixes[Math.floor(Math.random() * suffixes.length)]
    const random1 = generateRandomString(4 + Math.floor(Math.random() * 4), chars)
    const random2 = generateRandomString(3 + Math.floor(Math.random() * 3), chars)
    
    const path = `/${prefix}/${random1}/${suffix}/${random2}`
    if (config.value.suo5_config) {
      config.value.suo5_config.path = path
    }
    return path
  }

  // 文件模式密码随机生成
  const generateRandomFilePassword = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+'
    const password = generateRandomString(16, chars)
    config.value.password = password
    return password
  }

  // 随机生成输出文件名
  const generateRandomOutputFilename = () => {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789'
    const prefixes = ['admin', 'system', 'backup', 'config', 'update', 'api', 'upload']
    const extensions: Record<ScriptType, string> = {
      php: '.php',
      jsp: '.jsp',
      aspx: '.aspx',
      asp: '.asp'
    }
    const prefix = prefixes[Math.floor(Math.random() * prefixes.length)]
    const random = generateRandomString(6, chars)
    const ext = extensions[config.value.script_type]
    const filename = `${prefix}_${random}${ext}`
    config.value.output_filename = filename
    return filename
  }

  // 模式切换时的联动逻辑
  watch(
    () => config.value.mode,
    (newMode: PayloadMode, oldMode: PayloadMode) => {
      // 简单/高级模式切换逻辑（向后兼容）
      if ((newMode === 'advanced' && isShellMode.value) || 
          (oldMode === 'advanced' && isShellMode.value)) {
        if (newMode === 'advanced') {
          config.value.encode_type = undefined
          config.value.encrypt_algo = 'aes128_cbc'
        } else {
          config.value.encrypt_algo = undefined
          config.value.encode_type = 'none'
        }
      }

      // Suo5 专用模式联动
      if (isSuo5OnlyMode.value) {
        if (!config.value.suo5_config?.auth) {
          generateRandomSuo5Password()
        }
        if (!config.value.suo5_config?.path) {
          generateRandomSuo5Path()
        }
      }

      // 内存 Shell 模式联动 - 启用注入类型
      if (isMemoryShellMode.value) {
        if (!config.value.injection_type) {
          config.value.injection_type = 'tomcat_filter'
        }
      }

      // 文件通用模式联动 - 启用基础配置
      if (isFileCommonMode.value) {
        if (!config.value.password) {
          generateRandomFilePassword()
        }
        if (!config.value.output_filename) {
          generateRandomOutputFilename()
        }
      }

      // 文件代理模式联动 - 启用 Suo5 配置
      if (isFileProxyMode.value) {
        if (!config.value.suo5_config?.auth) {
          generateRandomSuo5Password()
        }
        if (!config.value.suo5_config?.path) {
          generateRandomSuo5Path()
        }
        if (!config.value.output_filename) {
          generateRandomOutputFilename()
        }
      }

      // 文件混合模式联动 - 启用完整配置
      if (isFileHybridMode.value) {
        if (!config.value.password) {
          generateRandomFilePassword()
        }
        if (!config.value.suo5_config?.auth) {
          generateRandomSuo5Password()
        }
        if (!config.value.suo5_config?.path) {
          generateRandomSuo5Path()
        }
        if (!config.value.output_filename) {
          generateRandomOutputFilename()
        }
      }
    }
  )

  // 脚本类型与注入类型的联动
  watch(
    () => config.value.script_type,
    (newType: ScriptType) => {
      if (isMemoryShellMode.value) {
        if (newType === 'jsp') {
          config.value.injection_type = 'tomcat_filter'
        } else if (newType === 'aspx') {
          config.value.injection_type = 'iis_httpmodule'
        } else {
          config.value.injection_type = undefined
        }
      }
    }
  )

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
      if (!config.value.password || config.value.password.trim() === '') {
        throw new Error('密码不能为空')
      }

      const result = await invoke<PayloadResult>('generate_payload_cmd', {
        config: config.value,
      })

      generatedResult.value = result
      
      history.value.unshift(result)
      
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

  // 客户端命令生成功能
  const generateClientCommand = (): string => {
    const commands: string[] = []
    
    if (isSuo5OnlyMode.value && config.value.suo5_config) {
      const { auth, path, timeout } = config.value.suo5_config
      commands.push(`# Suo5 客户端连接命令`)
      commands.push(`./suo5-client -auth "${auth}" -path "${path}" -timeout ${timeout}`)
    }
    
    if (generatedResult.value?.client_config) {
      const { key, iv, algorithm } = generatedResult.value.client_config
      commands.push(`\n# 客户端配置命令`)
      commands.push(`# Algorithm: ${algorithm}`)
      commands.push(`# Key: ${key}`)
      if (iv) {
        commands.push(`# IV: ${iv}`)
      }
    }
    
    if (config.value.self_destruct) {
      commands.push(`\n# 自毁模式已启用 - 载荷执行后将自动清除痕迹`)
    }
    
    return commands.join('\n')
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
  const clearHistory = async () => {
    try {
      await invoke('clear_payload_history')
      history.value = []
    } catch (err: any) {
      console.error('Failed to clear history:', err)
    }
  }

  // 复制代码到剪贴板
  const copyCode = async (code: string) => {
    try {
      await navigator.clipboard.writeText(code)
    } catch (err: any) {
      throw new Error(`复制失败：${err.message}`)
    }
  }

  // 模板管理
  const templates = ref<any[]>([])
  const isLoadingTemplates = ref(false)

  const loadTemplates = async () => {
    isLoadingTemplates.value = true
    try {
      templates.value = await invoke('get_payload_templates')
    } catch (err: any) {
      console.error('Failed to load templates:', err)
    } finally {
      isLoadingTemplates.value = false
    }
  }

  const addTemplate = async (template: any) => {
    try {
      const result = await invoke('add_payload_template', { template })
      templates.value.push(result)
      return result
    } catch (err: any) {
      throw new Error(`添加模板失败：${err.message}`)
    }
  }

  const updateTemplate = async (template: any) => {
    try {
      const result = await invoke('update_payload_template', { template })
      const index = templates.value.findIndex(t => t.name === template.name)
      if (index !== -1) {
        templates.value[index] = result
      }
      return result
    } catch (err: any) {
      throw new Error(`更新模板失败：${err.message}`)
    }
  }

  const deleteTemplate = async (name: string) => {
    try {
      await invoke('delete_payload_template', { name })
      templates.value = templates.value.filter(t => t.name !== name)
    } catch (err: any) {
      throw new Error(`删除模板失败：${err.message}`)
    }
  }

  const getTemplate = async (name: string) => {
    try {
      return await invoke('get_payload_template', { name })
    } catch (err: any) {
      throw new Error(`获取模板失败：${err.message}`)
    }
  }

  return {
    // 状态
    config,
    generatedResult,
    isGenerating,
    error,
    history,
    templates,
    isLoadingTemplates,
    
    // 计算属性
    isSimpleMode,
    isAdvancedMode,
    isFileBasedMode,
    isMemoryShellMode,
    isSuo5OnlyMode,
    isFileCommonMode,
    isFileProxyMode,
    isFileHybridMode,
    isShellMode,
    isFileMode,
    
    // 方法
    setMode,
    setScriptType,
    setFunctionType,
    setPassword,
    setEncodeType,
    setEncryptAlgo,
    setObfuscationLevel,
    setOutputFilename,
    setInjectionType,
    setSuo5Auth,
    setSuo5Path,
    setSuo5Timeout,
    setSelfDestruct,
    generateRandomSuo5Password,
    generateRandomSuo5Path,
    generateRandomFilePassword,
    generateRandomOutputFilename,
    generate,
    generatePreview,
    generateClientCommand,
    loadHistory,
    saveToFile,
    exportClientConfig,
    clearHistory,
    copyCode,
    loadTemplates,
    addTemplate,
    updateTemplate,
    deleteTemplate,
    getTemplate,
  }
})
