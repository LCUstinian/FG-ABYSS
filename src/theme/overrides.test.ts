import { describe, it, expect } from 'vitest'
import { buildOverrides } from './overrides'

describe('buildOverrides', () => {
  it('dark mode sets bodyColor to dark bg-base', () => {
    const o = buildOverrides('#4f9cff', true)
    expect(o.common?.bodyColor).toBe('#0d0e13')
  })

  it('light mode sets bodyColor to light bg-base', () => {
    const o = buildOverrides('#2463eb', false)
    expect(o.common?.bodyColor).toBe('#f6f7fb')
  })

  it('primaryColor equals accent', () => {
    const o = buildOverrides('#22d3ee', true)
    expect(o.common?.primaryColor).toBe('#22d3ee')
  })

  it('fontSize propagates to fontSizeMedium', () => {
    const o = buildOverrides('#4f9cff', true, '14px')
    expect(o.common?.fontSizeMedium).toBe('14px')
    expect(o.common?.fontSizeSmall).toBe('13px')
  })
})
