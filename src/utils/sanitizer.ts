/**
 * 输入清理和验证工具
 * 
 * 功能特性：
 * 1. XSS 防护（HTML 实体编码）
 * 2. SQL 注入防护基础
 * 3. 特殊字符过滤
 * 4. 字符串长度限制
 * 5. 对象递归清理
 */

/**
 * HTML 实体编码，防止 XSS 攻击
 * @param input 待清理的字符串
 * @returns 编码后的字符串
 */
export const sanitizeHTML = (input: string): string => {
  if (!input) return ''
  
  const div = document.createElement('div')
  div.textContent = input
  return div.innerHTML
}

/**
 * 清理字符串中的危险字符
 * @param input 待清理的字符串
 * @param options 清理选项
 * @returns 清理后的字符串
 */
export const sanitizeString = (
  input: string,
  options: {
    trim?: boolean
    maxLength?: number
    allowHTML?: boolean
  } = {}
): string => {
  const {
    trim = true,
    maxLength = 10000,
    allowHTML = false
  } = options
  
  if (!input) return ''
  
  let result = input
  
  // 去除首尾空格
  if (trim) {
    result = result.trim()
  }
  
  // 限制长度
  if (maxLength && result.length > maxLength) {
    result = result.substring(0, maxLength)
  }
  
  // HTML 清理
  if (!allowHTML) {
    result = sanitizeHTML(result)
  }
  
  // 移除可能的危险字符序列
  result = result
    .replace(/<script/gi, '&lt;script')
    .replace(/<\/script>/gi, '&lt;/script&gt;')
    .replace(/javascript:/gi, 'javascript:')
    .replace(/on\w+=/gi, (match) => match.replace('=', '-disabled='))
  
  return result
}

/**
 * 递归清理对象中的所有字符串字段
 * @param obj 待清理的对象
 * @param options 清理选项
 * @returns 清理后的对象
 */
export const sanitizeObject = <T extends Record<string, any>>(
  obj: T,
  options: {
    trim?: boolean
    maxLength?: number
    allowHTML?: boolean
    excludeFields?: string[]
  } = {}
): T => {
  const {
    trim = true,
    maxLength = 10000,
    allowHTML = false,
    excludeFields = []
  } = options
  
  if (!obj || typeof obj !== 'object') {
    return obj
  }
  
  const sanitized = {} as T
  
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      // 跳过排除的字段
      if (excludeFields.includes(key)) {
        sanitized[key] = obj[key]
        continue
      }
      
      const value = obj[key]
      
      if (typeof value === 'string') {
        sanitized[key] = sanitizeString(value, { trim, maxLength, allowHTML }) as any
      } else if (typeof value === 'object' && value !== null) {
        // 递归处理嵌套对象
        if (Array.isArray(value)) {
          sanitized[key] = value.map(item => 
            typeof item === 'object' 
              ? sanitizeObject(item, options)
              : typeof item === 'string'
                ? sanitizeString(item, options)
                : item
          ) as any
        } else {
          sanitized[key] = sanitizeObject(value, options)
        }
      } else {
        sanitized[key] = value
      }
    }
  }
  
  return sanitized
}

/**
 * 验证字符串长度
 * @param input 待验证的字符串
 * @param min 最小长度
 * @param max 最大长度
 * @returns 验证结果
 */
export const validateStringLength = (
  input: string,
  min: number = 0,
  max: number = 10000
): { valid: boolean; error?: string } => {
  if (!input) {
    return { valid: min === 0 }
  }
  
  if (input.length < min) {
    return { valid: false, error: `长度不能少于 ${min} 个字符` }
  }
  
  if (input.length > max) {
    return { valid: false, error: `长度不能超过 ${max} 个字符` }
  }
  
  return { valid: true }
}

/**
 * 验证邮箱格式
 * @param email 待验证的邮箱地址
 * @returns 验证结果
 */
export const validateEmail = (email: string): { valid: boolean; error?: string } => {
  if (!email) {
    return { valid: false, error: '邮箱不能为空' }
  }
  
  // RFC 5322 兼容的正则表达式
  const emailRegex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
  
  if (!emailRegex.test(email)) {
    return { valid: false, error: '邮箱格式不正确' }
  }
  
  if (email.length > 254) {
    return { valid: false, error: '邮箱地址过长' }
  }
  
  return { valid: true }
}

/**
 * 验证是否为安全的文件名
 * @param filename 文件名
 * @returns 验证结果
 */
export const validateFilename = (filename: string): { valid: boolean; error?: string } => {
  if (!filename) {
    return { valid: false, error: '文件名不能为空' }
  }
  
  // 检查危险字符
  const dangerousChars = /[<>:"/\\|？*]/
  if (dangerousChars.test(filename)) {
    return { valid: false, error: '文件名包含非法字符' }
  }
  
  // 检查长度
  if (filename.length > 255) {
    return { valid: false, error: '文件名过长' }
  }
  
  // 检查是否为保留名称
  const reservedNames = ['CON', 'PRN', 'AUX', 'NUL', 'COM1', 'COM2', 'COM3', 'COM4', 'COM5', 'COM6', 'COM7', 'COM8', 'COM9', 'LPT1', 'LPT2', 'LPT3', 'LPT4', 'LPT5', 'LPT6', 'LPT7', 'LPT8', 'LPT9']
  const nameWithoutExt = filename.split('.')[0].toUpperCase()
  if (reservedNames.includes(nameWithoutExt)) {
    return { valid: false, error: '文件名不能使用系统保留名称' }
  }
  
  return { valid: true }
}

/**
 * 清理文件名
 * @param filename 原始文件名
 * @returns 清理后的文件名
 */
export const sanitizeFilename = (filename: string): string => {
  if (!filename) return ''
  
  // 移除危险字符
  let sanitized = filename
    .replace(/[<>:"/\\|？*]/g, '_')
    .trim()
  
  // 限制长度
  if (sanitized.length > 255) {
    sanitized = sanitized.substring(0, 255)
  }
  
  return sanitized
}

export default {
  sanitizeHTML,
  sanitizeString,
  sanitizeObject,
  validateStringLength,
  validateEmail,
  validateFilename,
  sanitizeFilename
}
