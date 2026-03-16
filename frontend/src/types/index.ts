/**
 * 项目通用类型定义
 * 
 * 此文件包含项目中常用的 TypeScript 类型和接口定义
 * 用于替代 any 类型，提高类型安全性
 */

/**
 * 项目基本信息
 */
export interface Project {
  id: string
  name: string
  description?: string
  createdAt: string
  updatedAt: string
  deletedAt?: string
  status?: 'active' | 'inactive' | 'archived'
}

/**
 * WebShell 连接信息
 */
export interface WebShell {
  id: string
  projectId: string
  name: string
  url: string
  payload: string
  cryption: string
  encoding: string
  proxyType: string
  remark?: string
  status: 'active' | 'inactive'
  createdAt: string
  updatedAt: string
}

/**
 * 删除的项目（用于回收站）
 */
export interface DeletedProject extends Project {
  deletedAt: string
  deletedReason?: string
}

/**
 * 分页配置
 */
export interface PaginationConfig {
  page: number
  pageSize: number
  total: number
  showSizePicker: boolean
  pageSizes: number[]
  prefix?: (pagination: PaginationConfig) => string
}

/**
 * 表格列配置
 */
export interface TableColumn {
  key: string
  title: string
  width?: string | number
  minWidth?: string | number
  maxWidth?: string | number
  resizable?: boolean
  sortable?: boolean
  ellipsis?: boolean
  render?: (row: any, index: number) => any
}

/**
 * 系统状态信息
 */
export interface SystemStatus {
  memoryUsage: string
  processId: string
  cpuUsage: string
  uptime: string
}

/**
 * 表单验证结果
 */
export interface ValidationResult {
  valid: boolean
  error?: string
  field?: string
}

/**
 * API 响应基础结构
 */
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
  message?: string
  code?: number
}

/**
 * 下拉菜单选项
 */
export interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
}

/**
 * 按钮配置
 */
export interface ButtonConfig {
  label: string
  type?: 'primary' | 'success' | 'warning' | 'error' | 'info' | 'default'
  size?: 'small' | 'medium' | 'large'
  disabled?: boolean
  loading?: boolean
  icon?: string
  onClick?: () => void
}

/**
 * 模态框配置
 */
export interface ModalConfig {
  title: string
  show: boolean
  closable?: boolean
  maskClosable?: boolean
  preset?: 'dialog' | 'card'
  zIndex?: number
}

/**
 * 通知消息配置
 */
export interface NotificationConfig {
  type: 'success' | 'error' | 'warning' | 'info'
  content: string
  duration?: number
  showIcon?: boolean
}

/**
 * 主题配置
 */
export interface ThemeConfig {
  name: string
  isDark: boolean
  primaryColor?: string
  accentColor?: string
}

/**
 * 语言配置
 */
export interface LocaleConfig {
  code: string
  name: string
  flag?: string
}

/**
 * 用户偏好设置
 */
export interface UserPreferences {
  theme: 'light' | 'dark' | 'auto'
  language: string
  accentColor: string
  fontFamily?: string
  fontSize?: number
  sidebarCollapsed: boolean
}

/**
 * 事件定义
 */
export interface AppEvents {
  'project-change': (projectId: string) => void
  'webshell-create': (webshell: WebShell) => void
  'theme-change': (theme: ThemeConfig) => void
  'language-change': (locale: string) => void
}

/**
 * 工具函数类型
 */
export type AsyncCallback<T = void> = () => Promise<T>
export type SyncCallback<T = void> = () => T
export type Nullable<T> = T | null
export type Optional<T> = T | undefined
