/**
 * 时间格式化工具函数
 * 统一处理所有时间字段的显示格式
 */

/**
 * 格式化时间为本地时间字符串
 * 格式：YYYY-MM-DD HH:mm:ss
 * 
 * @param time - 时间字符串或 Date 对象或时间戳
 * @returns 格式化后的时间字符串
 */
export function formatTime(time: string | Date | number): string {
  if (!time) return ''
  
  const date = normalizeTime(time)
  if (!isValidDate(date)) return ''
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

/**
 * 格式化时间为日期字符串
 * 格式：YYYY-MM-DD
 * 
 * @param time - 时间字符串或 Date 对象或时间戳
 * @returns 格式化后的日期字符串
 */
export function formatDate(time: string | Date | number): string {
  if (!time) return ''
  
  const date = normalizeTime(time)
  if (!isValidDate(date)) return ''
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  
  return `${year}-${month}-${day}`
}

/**
 * 格式化时间为相对时间（如：刚刚、5 分钟前、1 小时前等）
 * 
 * @param time - 时间字符串或 Date 对象或时间戳
 * @returns 相对时间字符串
 */
export function formatRelativeTime(time: string | Date | number): string {
  if (!time) return ''
  
  const date = normalizeTime(time)
  if (!isValidDate(date)) return ''
  
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const months = Math.floor(days / 30)
  const years = Math.floor(months / 12)
  
  if (seconds < 60) {
    return '刚刚'
  } else if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else if (days < 30) {
    return `${days}天前`
  } else if (months < 12) {
    return `${months}个月前`
  } else {
    return `${years}年前`
  }
}

/**
 * 格式化时间为详细时间字符串（包含毫秒）
 * 格式：YYYY-MM-DD HH:mm:ss.SSS
 * 
 * @param time - 时间字符串或 Date 对象或时间戳
 * @returns 格式化后的详细时间字符串
 */
export function formatTimeWithMs(time: string | Date | number): string {
  if (!time) return ''
  
  const date = normalizeTime(time)
  if (!isValidDate(date)) return ''
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  const milliseconds = String(date.getMilliseconds()).padStart(3, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}`
}

/**
 * 标准化时间为 Date 对象
 * 
 * @param time - 时间字符串或 Date 对象或时间戳
 * @returns Date 对象
 */
function normalizeTime(time: string | Date | number): Date {
  if (time instanceof Date) {
    return time
  }
  
  if (typeof time === 'number') {
    // 如果是时间戳（毫秒）
    return new Date(time)
  }
  
  if (typeof time === 'string') {
    // 处理 ISO 8601 格式（如：2024-01-15T10:30:00.000Z）
    if (time.includes('T')) {
      return new Date(time)
    }
    
    // 处理标准格式（如：2024-01-15 10:30:00）
    return new Date(time.replace(/-/g, '/'))
  }
  
  return new Date()
}

/**
 * 验证日期是否有效
 * 
 * @param date - Date 对象
 * @returns 是否有效
 */
function isValidDate(date: Date): boolean {
  return date instanceof Date && !isNaN(date.getTime())
}

/**
 * 格式化时间范围（如：2024-01-15 至 2024-01-20）
 * 
 * @param startTime - 开始时间
 * @param endTime - 结束时间
 * @returns 格式化后的时间范围字符串
 */
export function formatTimeRange(startTime: string | Date | number, endTime: string | Date | number): string {
  const start = formatTime(startTime)
  const end = formatTime(endTime)
  
  if (!start && !end) return ''
  if (!start) return `至 ${end}`
  if (!end) return `${start} 至今`
  
  return `${start} 至 ${end}`
}

// 默认导出（兼容性考虑）
export default {
  formatTime,
  formatDate,
  formatRelativeTime,
  formatTimeWithMs,
  formatTimeRange
}
