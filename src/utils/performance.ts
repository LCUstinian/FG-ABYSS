/**
 * 性能监控工具
 */

// 性能指标记录
interface PerformanceMetric {
  name: string
  duration: number
  timestamp: number
}

class PerformanceMonitor {
  private metrics: PerformanceMetric[] = []
  private startTimeMap: Map<string, number> = new Map()
  private maxMetrics = 1000 // 最大存储的指标数量

  /**
   * 开始计时
   * @param name 操作名称
   */
  start(name: string): void {
    this.startTimeMap.set(name, performance.now())
  }

  /**
   * 结束计时并记录性能指标
   * @param name 操作名称
   * @returns 操作耗时（毫秒）
   */
  end(name: string): number {
    const startTime = this.startTimeMap.get(name)
    if (!startTime) {
      console.warn(`Performance monitor: No start time found for ${name}`)
      return 0
    }

    const duration = performance.now() - startTime
    const metric: PerformanceMetric = {
      name,
      duration,
      timestamp: Date.now()
    }

    this.metrics.push(metric)
    this.startTimeMap.delete(name)

    // 限制指标数量
    if (this.metrics.length > this.maxMetrics) {
      this.metrics.shift()
    }

    // 检查性能阈值
    if (duration > 200) {
      console.warn(`Performance warning: ${name} took ${duration.toFixed(2)}ms (threshold: 200ms)`)
    }

    return duration
  }

  /**
   * 获取性能指标
   * @param name 可选，操作名称
   * @returns 性能指标数组
   */
  getMetrics(name?: string): PerformanceMetric[] {
    if (name) {
      return this.metrics.filter(metric => metric.name === name)
    }
    return this.metrics
  }

  /**
   * 清除所有性能指标
   */
  clear(): void {
    this.metrics = []
    this.startTimeMap.clear()
  }

  /**
   * 生成性能报告
   * @returns 性能报告
   */
  generateReport(): string {
    const report: string[] = ['Performance Report:']
    
    // 按操作名称分组
    const groupedMetrics = this.metrics.reduce((acc, metric) => {
      if (!acc[metric.name]) {
        acc[metric.name] = []
      }
      acc[metric.name].push(metric.duration)
      return acc
    }, {} as { [key: string]: number[] })

    // 计算每个操作的统计数据
    Object.entries(groupedMetrics).forEach(([name, durations]) => {
      const avg = durations.reduce((sum, duration) => sum + duration, 0) / durations.length
      const min = Math.min(...durations)
      const max = Math.max(...durations)
      const p95 = this.calculatePercentile(durations, 95)

      report.push(`
${name}:`)
      report.push(`  Average: ${avg.toFixed(2)}ms`)
      report.push(`  Min: ${min.toFixed(2)}ms`)
      report.push(`  Max: ${max.toFixed(2)}ms`)
      report.push(`  95th percentile: ${p95.toFixed(2)}ms`)
      report.push(`  Count: ${durations.length}`)
    })

    return report.join('\n')
  }

  /**
   * 计算百分位数
   * @param values 值数组
   * @param percentile 百分位数（0-100）
   * @returns 计算结果
   */
  private calculatePercentile(values: number[], percentile: number): number {
    const sorted = values.sort((a, b) => a - b)
    const index = Math.ceil((percentile / 100) * sorted.length) - 1
    return sorted[Math.max(0, index)]
  }
}

// 导出单例实例
export const performanceMonitor = new PerformanceMonitor()

// 导出装饰器，用于监控函数性能
export function monitorPerformance(target: any, propertyKey: string, descriptor: PropertyDescriptor) {
  const originalMethod = descriptor.value
  
  descriptor.value = function(...args: any[]) {
    const methodName = `${target.constructor.name}.${propertyKey}`
    performanceMonitor.start(methodName)
    
    try {
      const result = originalMethod.apply(this, args)
      
      // 处理异步函数
      if (result instanceof Promise) {
        return result.then((res) => {
          performanceMonitor.end(methodName)
          return res
        })
      }
      
      performanceMonitor.end(methodName)
      return result
    } catch (error) {
      performanceMonitor.end(methodName)
      throw error
    }
  }
}
