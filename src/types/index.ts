/**
 * 类型统一导出
 * 
 * 此文件导出所有类型定义，方便统一导入
 */

// 通用类型
export * from './common'

// 业务领域类型
export * from './project'
export * from './webshell'

// 保留原有的向后兼容的类型定义（逐步迁移后可删除）
export interface DeletedProject extends Project {
  deletedAt: string
  deletedReason?: string
}

export interface PaginationConfig {
  page: number
  pageSize: number
  total: number
  showSizePicker: boolean
  pageSizes: number[]
  prefix?: (pagination: PaginationConfig) => string
}

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

export interface ValidationResult {
  valid: boolean
  error?: string
  field?: string
}

export interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
}

export interface ButtonConfig {
  label: string
  type?: 'primary' | 'success' | 'warning' | 'error' | 'info' | 'default'
  size?: 'small' | 'medium' | 'large'
  disabled?: boolean
  loading?: boolean
  icon?: string
  onClick?: () => void
}

export interface ModalConfig {
  title: string
  show: boolean
  closable?: boolean
  maskClosable?: boolean
  preset?: 'dialog' | 'card'
  zIndex?: number
}

export interface NotificationConfig {
  type: 'success' | 'error' | 'warning' | 'info'
  content: string
  duration?: number
  showIcon?: boolean
}

export interface ThemeConfig {
  name: string
  isDark: boolean
  primaryColor?: string
  accentColor?: string
}

export interface LocaleConfig {
  code: string
  name: string
  flag?: string
}

export interface UserPreferences {
  theme: 'light' | 'dark' | 'auto'
  language: string
  accentColor: string
  fontFamily?: string
  fontSize?: number
  sidebarCollapsed: boolean
}

export type AsyncCallback<T = void> = () => Promise<T>
export type SyncCallback<T = void> = () => T
export type Nullable<T> = T | null
export type Optional<T> = T | undefined
