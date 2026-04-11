import { describe, it, expect } from 'vitest'
import { sanitizeString, sanitizeHTML, sanitizeFilename } from './sanitizer'

describe('sanitizer', () => {
  describe('sanitizeString', () => {
    it('应该清理输入字符串', () => {
      const input = '<script>alert("XSS")</script>test'
      const result = sanitizeString(input)
      
      expect(result).not.toContain('<script>')
      expect(result).toContain('test')
    })

    it('应该处理空值', () => {
      expect(sanitizeString('')).toBe('')
      expect(sanitizeString(null as any)).toBe('')
      expect(sanitizeString(undefined as any)).toBe('')
    })
  })

  describe('sanitizeHTML', () => {
    it('应该清理HTML标签', () => {
      const html = '<div><p>Test</p><script>alert("XSS")</script></div>'
      const result = sanitizeHTML(html)
      
      expect(result).not.toContain('<div>')
      expect(result).not.toContain('<script>')
    })

    it('应该处理空值', () => {
      expect(sanitizeHTML('')).toBe('')
    })
  })

  describe('sanitizeFilename', () => {
    it('应该清理文件名中的危险字符', () => {
      const filename = 'file<name>with"dangerous"chars.txt'
      const result = sanitizeFilename(filename)
      
      expect(result).toBe('file_name_with_dangerous_chars.txt')
      expect(result).not.toContain('<')
      expect(result).not.toContain('>')
      expect(result).not.toContain('"')
    })

    it('应该处理空值', () => {
      expect(sanitizeFilename('')).toBe('')
    })
  })
})