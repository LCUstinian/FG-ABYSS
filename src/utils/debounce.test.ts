import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { debounce } from './debounce'

describe('debounce', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('应该在指定延迟后执行函数', () => {
    const mockFn = vi.fn()
    const debouncedFn = debounce(mockFn, 100)

    debouncedFn()
    expect(mockFn).not.toHaveBeenCalled()

    vi.advanceTimersByTime(99)
    expect(mockFn).not.toHaveBeenCalled()

    vi.advanceTimersByTime(1)
    expect(mockFn).toHaveBeenCalledTimes(1)
  })

  it('应该在连续调用时重置计时器', () => {
    const mockFn = vi.fn()
    const debouncedFn = debounce(mockFn, 100)

    debouncedFn()
    vi.advanceTimersByTime(50)
    debouncedFn()
    vi.advanceTimersByTime(50)
    expect(mockFn).not.toHaveBeenCalled()

    vi.advanceTimersByTime(50)
    expect(mockFn).toHaveBeenCalledTimes(1)
  })

  it('应该传递参数给原始函数', () => {
    const mockFn = vi.fn()
    const debouncedFn = debounce(mockFn, 100)

    debouncedFn('param1', 'param2')
    vi.advanceTimersByTime(100)

    expect(mockFn).toHaveBeenCalledWith('param1', 'param2')
  })

  it('应该使用默认延迟时间', () => {
    const mockFn = vi.fn()
    const debouncedFn = debounce(mockFn)

    debouncedFn()
    vi.advanceTimersByTime(300)
    expect(mockFn).toHaveBeenCalledTimes(1)
  })
})