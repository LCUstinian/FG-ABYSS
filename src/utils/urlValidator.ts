/**
 * URL 验证和清理工具
 * 
 * 功能特性：
 * 1. URL 格式验证
 * 2. 协议白名单检查
 * 3. XSS 攻击检测
 * 4. 内网地址访问控制
 * 5. URL 清理和标准化
 */

export interface UrlValidationResult {
  /** 验证是否通过 */
  valid: boolean
  /** 错误信息（如果验证失败） */
  error?: string
  /** 清理后的 URL（如果验证通过） */
  sanitized?: string
}

export interface UrlValidationOptions {
  /** 允许的协议列表，默认 ['http:', 'https:'] */
  allowedProtocols?: string[]
  /** 是否允许 localhost，默认 false */
  allowLocalhost?: boolean
  /** 是否允许内网 IP 地址，默认 false */
  allowInternalIPs?: boolean
  /** 是否允许 IP 地址，默认 true */
  allowIPAddresses?: boolean
}

/**
 * 验证并清理 URL
 * @param url 待验证的 URL
 * @param options 验证选项
 * @returns 验证结果
 */
export const validateUrl = (
  url: string,
  options: UrlValidationOptions = {}
): UrlValidationResult => {
  const {
    allowedProtocols = ['http:', 'https:'],
    allowLocalhost = false,
    allowInternalIPs = false,
    allowIPAddresses = true
  } = options
  
  // 空值检查
  if (!url?.trim()) {
    return { valid: false, error: 'URL 不能为空' }
  }
  
  const trimmedUrl = url.trim()
  
  // XSS 和危险协议检查
  const dangerousPatterns = [
    /javascript:/i,
    /data:/i,
    /vbscript:/i,
    /<script/i,
    /<\/script>/i,
    /on\w+=/i,  // onclick=, onerror= 等
    /eval\(/i,
    /alert\(/i,
    /prompt\(/i,
    /confirm\(/i,
    /document\./i,
    /window\./i,
    /location\./i
  ]
  
  for (const pattern of dangerousPatterns) {
    if (pattern.test(trimmedUrl)) {
      return { valid: false, error: 'URL 包含不安全内容' }
    }
  }
  
  // URL 解析和验证
  try {
    const urlObj = new URL(trimmedUrl)
    
    // 协议白名单检查
    if (!allowedProtocols.includes(urlObj.protocol)) {
      return { 
        valid: false, 
        error: `不支持的协议：${urlObj.protocol}，允许的协议：${allowedProtocols.join(', ')}` 
      }
    }
    
    const hostname = urlObj.hostname.toLowerCase()
    
    // localhost 检查
    if (!allowLocalhost && (hostname === 'localhost' || hostname === '127.0.0.1')) {
      return { valid: false, error: '不允许访问本地地址' }
    }
    
    // IP 地址检查
    if (!allowIPAddresses && isIPAddress(hostname)) {
      return { valid: false, error: '不允许使用 IP 地址' }
    }
    
    // 内网 IP 检查
    if (!allowInternalIPs && isInternalIP(hostname)) {
      return { valid: false, error: '不允许访问内网地址' }
    }
    
    // 返回清理后的 URL
    return { 
      valid: true, 
      sanitized: urlObj.toString() 
    }
  } catch (error) {
    return { valid: false, error: 'URL 格式不正确' }
  }
}

/**
 * 检查是否为 IP 地址
 */
const isIPAddress = (hostname: string): boolean => {
  // IPv4
  const ipv4Pattern = /^(\d{1,3}\.){3}\d{1,3}$/
  if (ipv4Pattern.test(hostname)) {
    return true
  }
  
  // IPv6（简化检查）
  const ipv6Pattern = /^([0-9a-fA-F]{0,4}:){2,7}[0-9a-fA-F]{0,4}$/
  if (ipv6Pattern.test(hostname)) {
    return true
  }
  
  return false
}

/**
 * 检查是否为内网 IP 地址
 */
const isInternalIP = (hostname: string): boolean => {
  const internalIPPatterns = [
    /^192\.168\./,                    // 192.168.x.x
    /^10\./,                          // 10.x.x.x
    /^172\.(1[6-9]|2[0-9]|3[01])\./,  // 172.16.x.x - 172.31.x.x
    /^127\./,                         // 127.x.x.x (loopback)
    /^0\.0\.0\.0$/,                   // 0.0.0.0
    /^169\.254\./                     // 169.254.x.x (link-local)
  ]
  
  for (const pattern of internalIPPatterns) {
    if (pattern.test(hostname)) {
      return true
    }
  }
  
  return false
}

/**
 * 清理 URL（移除危险字符和编码）
 * @param url 待清理的 URL
 * @returns 清理后的 URL
 */
export const sanitizeUrl = (url: string): string => {
  try {
    const urlObj = new URL(url)
    
    // 清理查询参数中的危险字符
    urlObj.searchParams.forEach((value, key) => {
      // 移除可能的脚本标签
      const sanitized = value
        .replace(/<script/gi, '')
        .replace(/<\/script>/gi, '')
        .replace(/javascript:/gi, '')
        .replace(/on\w+=/gi, '')
      
      urlObj.searchParams.set(key, sanitized)
    })
    
    return urlObj.toString()
  } catch {
    return url
  }
}

/**
 * 批量验证 URL 列表
 * @param urls URL 列表
 * @param options 验证选项
 * @returns 验证结果数组
 */
export const validateUrls = (
  urls: string[],
  options: UrlValidationOptions = {}
): UrlValidationResult[] => {
  return urls.map(url => validateUrl(url, options))
}

export default {
  validateUrl,
  validateUrls,
  sanitizeUrl
}
