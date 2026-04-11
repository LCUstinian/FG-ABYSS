import { ref, computed } from 'vue'
import type { Plugin, PluginActionRequest, PluginLoadResult, PluginSignatureResult } from '@/types/plugin'

// 模拟内置插件数据
const builtInPlugins: Plugin[] = [
  {
    id: 'command-execution',
    name: '命令执行',
    version: '1.0.0',
    description: '提供伪终端(PTY)模拟，提升交互体验',
    author: 'FG-ABYSS Team',
    enabled: true,
    signature: 'BUILT-IN',
    isBuiltIn: true
  },
  {
    id: 'file-management',
    name: '文件管理',
    version: '1.0.0',
    description: '支持断点续传、哈希校验、时间戳篡改',
    author: 'FG-ABYSS Team',
    enabled: true,
    signature: 'BUILT-IN',
    isBuiltIn: true
  },
  {
    id: 'database-management',
    name: '数据库管理',
    version: '1.0.0',
    description: '自适应识别 MySQL, MSSQL, Oracle 等环境',
    author: 'FG-ABYSS Team',
    enabled: true,
    signature: 'BUILT-IN',
    isBuiltIn: true
  }
]

export function usePlugin() {
  const plugins = ref<Plugin[]>([...builtInPlugins])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性：按状态分类插件
  const enabledPlugins = computed(() => plugins.value.filter(p => p.enabled))
  const disabledPlugins = computed(() => plugins.value.filter(p => !p.enabled))
  const builtInPluginsList = computed(() => plugins.value.filter(p => p.isBuiltIn))
  const customPlugins = computed(() => plugins.value.filter(p => !p.isBuiltIn))

  // 加载插件列表
  const loadPlugins = async () => {
    loading.value = true
    error.value = null
    try {
      // 模拟从本地加载插件
      // 实际实现中，这里应该调用后端API加载插件
      await new Promise(resolve => setTimeout(resolve, 500))
      // 保持内置插件，后续可以添加从本地目录加载的自定义插件
    } catch (err) {
      error.value = '加载插件失败'
      console.error('加载插件失败:', err)
    } finally {
      loading.value = false
    }
  }

  // 执行插件操作
  const performPluginAction = async (action: PluginActionRequest): Promise<PluginLoadResult> => {
    loading.value = true
    error.value = null
    try {
      const { pluginId, action: actionType, path } = action
      const pluginIndex = plugins.value.findIndex(p => p.id === pluginId)

      if (pluginIndex === -1) {
        return {
          success: false,
          message: '插件不存在'
        }
      }

      switch (actionType) {
        case 'enable':
          plugins.value[pluginIndex].enabled = true
          return {
            success: true,
            message: '插件已启用',
            plugin: plugins.value[pluginIndex]
          }
        case 'disable':
          plugins.value[pluginIndex].enabled = false
          return {
            success: true,
            message: '插件已禁用',
            plugin: plugins.value[pluginIndex]
          }
        case 'uninstall':
          if (plugins.value[pluginIndex].isBuiltIn) {
            return {
              success: false,
              message: '内置插件无法卸载'
            }
          }
          plugins.value.splice(pluginIndex, 1)
          return {
            success: true,
            message: '插件已卸载'
          }
        case 'install':
          if (!path) {
            return {
              success: false,
              message: '插件路径不能为空'
            }
          }
          // 模拟安装插件
          const newPlugin: Plugin = {
            id: `custom-${Date.now()}`,
            name: '自定义插件',
            version: '1.0.0',
            description: '自定义插件描述',
            author: 'Unknown',
            enabled: true,
            signature: 'UNVERIFIED',
            path,
            isBuiltIn: false
          }
          plugins.value.push(newPlugin)
          return {
            success: true,
            message: '插件已安装',
            plugin: newPlugin
          }
        default:
          return {
            success: false,
            message: '不支持的操作'
          }
      }
    } catch (err) {
      error.value = '操作插件失败'
      console.error('操作插件失败:', err)
      return {
        success: false,
        message: '操作插件失败'
      }
    } finally {
      loading.value = false
    }
  }

  // 验证插件签名
  const verifyPluginSignature = async (pluginPath: string): Promise<PluginSignatureResult> => {
    loading.value = true
    error.value = null
    try {
      // 模拟签名验证
      // 实际实现中，这里应该调用后端API进行签名验证
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      // 模拟验证结果
      const isValid = Math.random() > 0.3 // 70%概率验证通过
      return {
        valid: isValid,
        message: isValid ? '签名验证通过' : '签名验证失败',
        pluginId: `plugin-${Date.now()}`
      }
    } catch (err) {
      error.value = '验证插件签名失败'
      console.error('验证插件签名失败:', err)
      return {
        valid: false,
        message: '验证插件签名失败'
      }
    } finally {
      loading.value = false
    }
  }

  // 初始化加载插件
  loadPlugins()

  return {
    plugins,
    loading,
    error,
    enabledPlugins,
    disabledPlugins,
    builtInPluginsList,
    customPlugins,
    loadPlugins,
    performPluginAction,
    verifyPluginSignature
  }
}
