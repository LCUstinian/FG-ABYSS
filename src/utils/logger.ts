/**
 * 环境感知的日志工具模块
 * 
 * 功能特性：
 * 1. 生产环境自动禁用调试日志
 * 2. 结构化日志输出
 * 3. 支持日志级别控制
 * 4. 防止敏感数据泄露
 */

const isDev = import.meta.env.DEV

export type LogLevel = 'debug' | 'info' | 'warn' | 'error'

export interface LoggerConfig {
  /** 日志前缀标识 */
  prefix: string
  /** 是否在生产环境启用 */
  enableInProd: boolean
  /** 启用的日志级别 */
  level: LogLevel
}

/**
 * 创建带前缀的日志工具
 * @param config 日志配置
 * @returns 日志工具对象
 */
export const createLogger = (config: LoggerConfig) => {
  const shouldLog = (level: LogLevel): boolean => {
    if (!isDev && !config.enableInProd) {
      return false
    }
    
    const levels: LogLevel[] = ['debug', 'info', 'warn', 'error']
    const currentIndex = levels.indexOf(config.level)
    const targetIndex = levels.indexOf(level)
    
    return targetIndex >= currentIndex
  }
  
  const formatArgs = (args: any[]): any[] => {
    return args.map(arg => {
      if (typeof arg === 'object' && arg !== null) {
        // 避免循环引用和过度展开
        try {
          return JSON.parse(JSON.stringify(arg))
        } catch {
          return '[Object]'
        }
      }
      return arg
    })
  }
  
  const debug = (...args: any[]) => {
    if (shouldLog('debug')) {
      console.debug(`[${config.prefix}]`, ...formatArgs(args))
    }
  }
  
  const log = (...args: any[]) => {
    if (shouldLog('info')) {
      console.log(`[${config.prefix}]`, ...formatArgs(args))
    }
  }
  
  const info = (...args: any[]) => {
    if (shouldLog('info')) {
      console.info(`[${config.prefix}]`, ...formatArgs(args))
    }
  }
  
  const warn = (...args: any[]) => {
    if (shouldLog('warn')) {
      console.warn(`[${config.prefix}]`, ...formatArgs(args))
    }
  }
  
  const error = (...args: any[]) => {
    if (shouldLog('error')) {
      console.error(`[${config.prefix}]`, ...formatArgs(args))
    }
  }
  
  return { debug, log, info, warn, error }
}

/**
 * 预定义的日志工具实例
 */
export const appLogger = createLogger({ prefix: 'App', enableInProd: false, level: 'info' })
export const apiLogger = createLogger({ prefix: 'API', enableInProd: false, level: 'warn' })
export const eventLogger = createLogger({ prefix: 'Event', enableInProd: false, level: 'info' })
export const componentLogger = createLogger({ prefix: 'Component', enableInProd: false, level: 'info' })

/**
 * 简易日志工具（默认配置）
 */
const logger = {
  debug: (...args: any[]) => isDev && console.debug(...args),
  log: (...args: any[]) => isDev && console.log(...args),
  info: (...args: any[]) => isDev && console.info(...args),
  warn: (...args: any[]) => isDev && console.warn(...args),
  error: (...args: any[]) => console.error(...args),
}

export default logger
