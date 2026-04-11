/**
 * 系统状态 API
 * 使用 Tauri Mock 适配层获取系统状态
 */

// 使用 Tauri Mock 适配层
import { invoke } from '@/utils/tauri-mock-adapter'
import type { SystemStatus as MockSystemStatus } from '@/utils/tauri-mock-adapter'
import { performanceMonitor } from '@/utils/performance'

// 定义系统状态接口
export interface SystemStatus {
  memoryUsage: string
  processId: string
  cpuUsage: string
  uptime: string
}

// 系统状态缓存
let systemStatusCache: SystemStatus | null = null
let lastFetchTime: number = 0
const CACHE_DURATION = 500 // 缓存持续时间（毫秒）

/**
 * 获取系统状态信息
 * @returns 系统状态对象
 */
export async function getSystemStatus(): Promise<SystemStatus> {
  performanceMonitor.start('getSystemStatus')
  try {
    // 检查缓存是否有效
    const now = Date.now()
    if (systemStatusCache && (now - lastFetchTime) < CACHE_DURATION) {
      performanceMonitor.end('getSystemStatus')
      return systemStatusCache
    }
    
    // 使用 Tauri Mock API 获取系统状态
    const status = await invoke<MockSystemStatus>('get_system_status')
    
    const systemStatus: SystemStatus = {
      memoryUsage: status.memoryUsage || 'N/A',
      processId: status.processId || 'N/A',
      cpuUsage: status.cpuUsage || 'N/A',
      uptime: status.uptime || '0'
    }
    
    // 更新缓存
    systemStatusCache = systemStatus
    lastFetchTime = now
    
    performanceMonitor.end('getSystemStatus')
    return systemStatus
  } catch (error) {
    console.error('Failed to fetch system status:', error)
    // 返回缓存值或默认值
    if (systemStatusCache) {
      performanceMonitor.end('getSystemStatus')
      return systemStatusCache
    }
    performanceMonitor.end('getSystemStatus')
    return {
      memoryUsage: 'N/A',
      processId: 'N/A',
      cpuUsage: 'N/A',
      uptime: '0'
    }
  }
}
