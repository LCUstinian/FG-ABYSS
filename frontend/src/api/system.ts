/**
 * 系统状态 API
 * 调用后端 GetSystemStatus 方法获取真实系统状态
 */

// 导入自动生成的 Wails 绑定
import { GetSystemStatus } from '../../bindings/fg-abyss/internal/app/handlers/systemhandler'

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
    
    // 格式化内存使用信息
    let memoryUsage = 'N/A'
    if (status.memory && typeof status.memory === 'object') {
      const total = (status.memory.total as number) || 0
      const used = (status.memory.used as number) || 0
      memoryUsage = `${formatBytes(used)} / ${formatBytes(total)}`
    }
    
    // 获取进程 ID
    const processId = status.processID?.toString() || status.processId?.toString() || 'N/A'
    
    // 格式化 CPU 使用率
    let cpuUsage = 'N/A'
    if (status.cpuPercent !== undefined && status.cpuPercent !== null) {
      cpuUsage = `${(status.cpuPercent as number).toFixed(1)}%`
    } else if (status.processCpuPercent !== undefined && status.processCpuPercent !== null) {
      cpuUsage = `${(status.processCpuPercent as number).toFixed(1)}%`
    }
    
    // 格式化运行时间（转换为秒）
    let uptime = '0'
    if (status.uptime) {
      uptime = parseUptime(status.uptime as string)
    }
    
    return {
      memoryUsage,
      processId,
      cpuUsage,
      uptime
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

// 格式化字节数
function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// 解析运行时间字符串（如 "2h3m4.567s"）为秒数
function parseUptime(uptimeStr: string): string {
  try {
    // 匹配小时、分钟、秒
    const hoursMatch = uptimeStr.match(/(\d+)h/)
    const minutesMatch = uptimeStr.match(/(\d+)m/)
    const secondsMatch = uptimeStr.match(/(\d+\.?\d*)s/)
    
    const hours = hoursMatch ? parseInt(hoursMatch[1]) : 0
    const minutes = minutesMatch ? parseInt(minutesMatch[1]) : 0
    const seconds = secondsMatch ? parseFloat(secondsMatch[1]) : 0
    
    const totalSeconds = hours * 3600 + minutes * 60 + seconds
    return totalSeconds.toFixed(3)
  } catch {
    return '0'
  }
}
