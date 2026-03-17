/**
 * WebShell 相关类型定义
 */

/**
 * WebShell 实体
 */
export interface WebShell {
  id: string
  projectId: string
  url: string
  payload: string
  cryption: string
  encoding: string
  proxyType: string
  remark: string
  status: string
  createdAt?: string
  updatedAt?: string
}

/**
 * 创建 WebShell 请求
 */
export interface CreateWebShellRequest {
  projectId: string
  url: string
  payload?: string
  cryption?: string
  encoding?: string
  proxyType?: string
  remark?: string
  status?: string
}

/**
 * WebShell 连接结果
 */
export interface WebShellConnectionResult {
  success: boolean
  connectionId?: string
  message?: string
}
