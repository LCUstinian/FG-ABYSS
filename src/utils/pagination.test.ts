import { describe, it, expect } from 'vitest'
import { calculateOptimalPageSize, getPageSizeOptions, PAGINATION_CONFIG } from './pagination'

describe('pagination', () => {
  describe('calculateOptimalPageSize', () => {
    it('应该根据可用高度计算最优分页大小', () => {
      const result = calculateOptimalPageSize(100, 500)
      
      expect(result).toBeGreaterThanOrEqual(PAGINATION_CONFIG.MIN_PAGE_SIZE)
      expect(result).toBeLessThanOrEqual(PAGINATION_CONFIG.MAX_PAGE_SIZE)
    })

    it('应该处理数据量小于最优分页大小的情况', () => {
      const result = calculateOptimalPageSize(3, 500)
      
      expect(result).toBe(PAGINATION_CONFIG.MIN_PAGE_SIZE)
    })

    it('应该处理空数据的情况', () => {
      const result = calculateOptimalPageSize(0, 500)
      
      expect(result).toBeGreaterThanOrEqual(PAGINATION_CONFIG.MIN_PAGE_SIZE)
      expect(result).toBeLessThanOrEqual(PAGINATION_CONFIG.MAX_PAGE_SIZE)
    })

    it('应该返回预设选项中的值', () => {
      const result = calculateOptimalPageSize(100, 500)
      
      expect(PAGINATION_CONFIG.PAGE_SIZE_OPTIONS).toContain(result)
    })
  })

  describe('getPageSizeOptions', () => {
    it('应该返回分页大小选项数组', () => {
      const result = getPageSizeOptions()
      
      expect(Array.isArray(result)).toBe(true)
      expect(result).toEqual(PAGINATION_CONFIG.PAGE_SIZE_OPTIONS)
    })

    it('应该返回预设的分页大小选项', () => {
      const result = getPageSizeOptions()
      
      expect(result).toContain(5)
      expect(result).toContain(10)
      expect(result).toContain(20)
      expect(result).toContain(30)
      expect(result).toContain(50)
    })
  })

  describe('PAGINATION_CONFIG', () => {
    it('应该包含正确的配置值', () => {
      expect(PAGINATION_CONFIG.ROW_HEIGHT).toBe(53)
      expect(PAGINATION_CONFIG.TABLE_HEADER_HEIGHT).toBe(50)
      expect(PAGINATION_CONFIG.TABLE_FOOTER_HEIGHT).toBe(70)
      expect(PAGINATION_CONFIG.OTHER_UI_HEIGHT).toBe(100)
      expect(PAGINATION_CONFIG.MIN_PAGE_SIZE).toBe(5)
      expect(PAGINATION_CONFIG.MAX_PAGE_SIZE).toBe(50)
      expect(PAGINATION_CONFIG.PAGE_SIZE_OPTIONS).toEqual([5, 10, 20, 30, 50])
    })
  })
})