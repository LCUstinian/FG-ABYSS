import { describe, it, expect } from 'vitest'
import { formatTime, formatDate, formatRelativeTime } from './formatTime'

describe('formatTime', () => {
  it('应该格式化日期时间字符串', () => {
    const date = new Date('2024-01-15T10:30:00Z')
    const result = formatTime(date)
    
    expect(result).toContain('2024')
    expect(result).toContain('01')
    expect(result).toContain('15')
  })

  it('应该处理当前时间', () => {
    const now = new Date()
    const result = formatTime(now)
    
    expect(result).toBeDefined()
    expect(typeof result).toBe('string')
  })

  it('应该处理空值', () => {
    expect(formatTime('')).toBe('')
    expect(formatTime(null as any)).toBe('')
    expect(formatTime(undefined as any)).toBe('')
  })

  it('应该格式化时间戳', () => {
    const timestamp = 1704067200000 // 2024-01-01 00:00:00
    const result = formatTime(timestamp)
    
    expect(result).toContain('2024')
    expect(result).toContain('01')
    expect(result).toContain('01')
  })
})

describe('formatDate', () => {
  it('应该只返回日期部分', () => {
    const date = new Date('2024-01-15T10:30:00Z')
    const result = formatDate(date)
    
    expect(result).toBe('2024-01-15')
  })

  it('应该处理空值', () => {
    expect(formatDate('')).toBe('')
  })
})

describe('formatRelativeTime', () => {
  it('应该处理刚刚的时间', () => {
    const now = new Date()
    const result = formatRelativeTime(now)
    
    expect(result).toBe('刚刚')
  })

  it('应该处理空值', () => {
    expect(formatRelativeTime('')).toBe('')
  })
})
