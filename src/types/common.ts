/**
 * 通用类型定义
 */

/**
 * 分页参数
 */
export interface PaginationParams {
  page: number
  pageSize: number
}

/**
 * 分页结果
 */
export interface PaginationResult<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
}

/**
 * API 响应基础类型
 */
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
}

/**
 * 系统状态
 */
export interface SystemStatus {
  memoryUsage: string
  processId: string
  cpuUsage: string
  uptime: string
}
