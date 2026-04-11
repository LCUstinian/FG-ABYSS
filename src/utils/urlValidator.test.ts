import { describe, it, expect } from 'vitest'
import { validateUrl } from './urlValidator'

describe('urlValidator', () => {
  it('应该验证有效的URL', () => {
    const validUrls = [
      'https://example.com',
      'http://example.com:8080',
      'https://subdomain.example.com/path',
      'https://example.com/path?query=value'
    ]
    
    validUrls.forEach(url => {
      const result = validateUrl(url)
      expect(result.valid).toBe(true)
      expect(result.sanitized).toBeDefined()
    })
  })

  it('应该拒绝无效的URL', () => {
    const invalidUrls = [
      '',
      'example.com',
      'http://',
      'https://',
      'ftp://example.com',
      'mailto:user@example.com'
    ]
    
    invalidUrls.forEach(url => {
      const result = validateUrl(url)
      expect(result.valid).toBe(false)
      expect(result.error).toBeDefined()
    })
  })

  it('应该处理空值', () => {
    const result1 = validateUrl('')
    expect(result1.valid).toBe(false)
    expect(result1.error).toBeDefined()

    const result2 = validateUrl(null as any)
    expect(result2.valid).toBe(false)
    expect(result2.error).toBeDefined()

    const result3 = validateUrl(undefined as any)
    expect(result3.valid).toBe(false)
    expect(result3.error).toBeDefined()
  })
})