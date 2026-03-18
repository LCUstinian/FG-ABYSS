/**
 * Dashboard 数据管理 Composable
 * 提供首页所需的数据和操作方法
 */

import { ref, computed, onMounted, onUnmounted } from 'vue'
import { invoke } from '@/utils/tauri-mock-adapter'

export interface DashboardStats {
  totalProjects: number
  totalWebShells: number
  totalPayloads: number
  totalPlugins: number
  activeWebShells: number
  inactiveWebShells: number
}

export interface RecentActivity {
  id: string
  type: 'project' | 'webshell' | 'payload' | 'plugin'
  action: 'create' | 'update' | 'delete' | 'connect'
  description: string
  timestamp: string
  icon: string
  status: 'success' | 'info' | 'warning' | 'error'
}

export function useDashboard() {
  // 统计数据
  const stats = ref<DashboardStats>({
    totalProjects: 0,
    totalWebShells: 0,
    totalPayloads: 0,
    totalPlugins: 0,
    activeWebShells: 0,
    inactiveWebShells: 0
  })

  // 最近活动
  const activities = ref<RecentActivity[]>([])

  // 加载统计数据的函数
  const loadStats = async () => {
    try {
      // 从后端获取统计数据
      const [projects, webshells, payloads, plugins] = await Promise.all([
        invoke<any[]>('get_projects').catch(() => []),
        invoke<any[]>('get_all_webshells').catch(() => []),
        invoke<any[]>('get_generated_payloads').catch(() => []),
        invoke<any[]>('get_plugins').catch(() => [])
      ])

      stats.value = {
        totalProjects: projects.length,
        totalWebShells: webshells.length,
        totalPayloads: payloads.length,
        totalPlugins: plugins.length,
        activeWebShells: webshells.filter(w => w.status === 'active').length,
        inactiveWebShells: webshells.filter(w => w.status === 'inactive').length
      }
    } catch (error) {
      console.error('Failed to load dashboard stats:', error)
      // 使用模拟数据
      stats.value = {
        totalProjects: 24,
        totalWebShells: 156,
        totalPayloads: 89,
        totalPlugins: 12,
        activeWebShells: 94,
        inactiveWebShells: 62
      }
    }
  }

  // 加载最近活动
  const loadActivities = async () => {
    try {
      // 这里可以从后端获取最近活动记录
      // 目前使用模拟数据
      activities.value = [
        {
          id: '1',
          type: 'webshell',
          action: 'connect',
          description: '成功连接 WebShell',
          timestamp: '5 分钟前',
          icon: 'check',
          status: 'success'
        },
        {
          id: '2',
          type: 'project',
          action: 'create',
          description: '创建新项目 "测试项目"',
          timestamp: '10 分钟前',
          icon: 'plus',
          status: 'info'
        },
        {
          id: '3',
          type: 'webshell',
          action: 'update',
          description: '更新 WebShell 配置',
          timestamp: '15 分钟前',
          icon: 'edit',
          status: 'info'
        },
        {
          id: '4',
          type: 'plugin',
          action: 'update',
          description: '插件版本需要更新',
          timestamp: '1 小时前',
          icon: 'alert',
          status: 'warning'
        },
        {
          id: '5',
          type: 'payload',
          action: 'create',
          description: '下载新的 Payload 模板',
          timestamp: '3 小时前',
          icon: 'download',
          status: 'success'
        }
      ]
    } catch (error) {
      console.error('Failed to load activities:', error)
    }
  }

  // 刷新数据
  const refresh = async () => {
    await Promise.all([
      loadStats(),
      loadActivities()
    ])
  }

  // 生命周期钩子
  onMounted(() => {
    refresh()
  })

  return {
    stats,
    activities,
    refresh,
    loadStats,
    loadActivities
  }
}
