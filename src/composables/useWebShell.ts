/**
 * WebShell 管理 Composable
 * 封装 WebShell 相关的业务逻辑和状态管理
 */

import { ref, computed } from 'vue'
import { invoke, emitEvent } from '@/utils/tauri-mock-adapter'
import type { WebShell } from '@/types'

export interface WebShellFilters {
  page?: number
  pageSize?: number
  searchQuery?: string
  sortField?: string
  sortDirection?: 'asc' | 'desc'
}

export function useWebShell() {
  const webshells = ref<WebShell[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const filters = ref<WebShellFilters>({})

  /**
   * 获取 WebShell 列表
   */
  const fetchWebShells = async (projectId: string, filters?: WebShellFilters) => {
    loading.value = true
    error.value = null
    try {
      const webshellList = await invoke<WebShell[]>('get_webshells', {
        projectId,
        ...filters
      })
      webshells.value = webshellList
    } catch (err: any) {
      error.value = err.message || '获取 WebShell 列表失败'
      console.error('获取 WebShell 列表失败:', err)
    } finally {
      loading.value = false
    }
  }

  /**
   * 创建 WebShell
   */
  const createWebShell = async (webshellData: Partial<WebShell>) => {
    error.value = null
    try {
      const result = await invoke('create_webshell', webshellData)
      if (result.success) {
        await fetchWebShells(webshellData.projectId!)
        return result.webshell
      }
    } catch (err: any) {
      error.value = err.message || '创建 WebShell 失败'
      throw err
    }
  }

  /**
   * 删除 WebShell
   */
  const deleteWebShell = async (webshellId: string, projectId?: string) => {
    error.value = null
    try {
      await invoke('delete_webshell', { webshellId })
      if (projectId) {
        await fetchWebShells(projectId, filters.value)
      }
    } catch (err: any) {
      error.value = err.message || '删除 WebShell 失败'
      throw err
    }
  }

  /**
   * 恢复已删除的 WebShell
   */
  const recoverWebShell = async (webshellId: string, projectId?: string) => {
    error.value = null
    try {
      await invoke('recover_webshell', { webshellId })
      if (projectId) {
        await fetchWebShells(projectId, filters.value)
      }
    } catch (err: any) {
      error.value = err.message || '恢复 WebShell 失败'
      throw err
    }
  }

  /**
   * 连接 WebShell
   */
  const connectWebShell = async (webshell: WebShell) => {
    error.value = null
    try {
      const result = await invoke('connect_webshell', { id: webshell.id })
      if (result.success) {
        // 打开控制窗口
        await emitEvent('open-webshell-window', {
          id: webshell.id,
          name: webshell.remark || webshell.url,
          url: webshell.url
        })
      }
      return result
    } catch (err: any) {
      error.value = err.message || '连接 WebShell 失败'
      throw err
    }
  }

  /**
   * 更新筛选条件
   */
  const updateFilters = (newFilters: WebShellFilters) => {
    filters.value = { ...filters.value, ...newFilters }
  }

  /**
   * 清空筛选条件
   */
  const clearFilters = () => {
    filters.value = {}
  }

  return {
    webshells,
    loading,
    error,
    filters,
    fetchWebShells,
    createWebShell,
    deleteWebShell,
    recoverWebShell,
    connectWebShell,
    updateFilters,
    clearFilters
  }
}
