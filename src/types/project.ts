/**
 * 项目相关类型定义
 */

/**
 * 项目实体
 */
export interface Project {
  id: string
  name: string
  description?: string
  createdAt?: string
  updatedAt?: string
  deletedAt?: string
}

/**
 * 创建项目请求
 */
export interface CreateProjectRequest {
  name: string
  description?: string
}

/**
 * 项目统计信息
 */
export interface ProjectStats {
  totalWebShells: number
  activeWebShells: number
  totalPayloads: number
}
