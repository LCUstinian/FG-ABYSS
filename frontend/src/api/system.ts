/**
 * 系统状态 API
 * 调用后端 GetSystemStatus 方法获取真实系统状态
 */

// 导入自动生成的 Wails 绑定
import { GetSystemStatus } from '../../bindings/fg-abyss/app'
import type { SystemStatus as WailsSystemStatus } from '../../bindings/fg-abyss/models'

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
    // 使用自动生成的绑定函数调用后端
    const status = await GetSystemStatus()
    return {
      memoryUsage: status.memoryUsage,
      processId: status.processId,
      cpuUsage: status.cpuUsage,
      uptime: status.uptime
    }
  } catch (error) {
    console.error('Failed to fetch system status:', error)
    // 返回默认值
    return {
      memoryUsage: 'N/A',
      processId: 'N/A',
      cpuUsage: 'N/A',
      uptime: 'N/A'
    }
  }
}
