/**
 * 系统状态 API
 * 使用 Tauri Mock 适配层获取系统状态
 */

// 使用 Tauri Mock 适配层
import { invoke } from '@/utils/tauri-mock-adapter'
import type { SystemStatus as MockSystemStatus } from '@/utils/tauri-mock-adapter'

// 定义系统状态接口
export interface SystemStatus {
  memoryUsage: string
  processId: string
  cpuUsage: string
  uptime: string
}

/**
 * 获取系统状态信息
 * @returns 系统状态对象
 */
export async function getSystemStatus(): Promise<SystemStatus> {
  try {
    // 使用 Tauri Mock API 获取系统状态
    const status = await invoke<MockSystemStatus>('get_system_status')
    
    return {
      memoryUsage: status.memoryUsage || 'N/A',
      processId: status.processId || 'N/A',
      cpuUsage: status.cpuUsage || 'N/A',
      uptime: status.uptime || '0'
    }
  } catch (error) {
    console.error('Failed to fetch system status:', error)
    // 返回默认值
    return {
      memoryUsage: 'N/A',
      processId: 'N/A',
      cpuUsage: 'N/A',
      uptime: '0'
    }
  }
}
