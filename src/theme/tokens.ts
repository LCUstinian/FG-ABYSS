export type AccentKey = 'blue' | 'cyan' | 'purple' | 'pink' | 'orange' | 'green'

export const ACCENT_COLORS: Record<AccentKey, { dark: string; light: string }> = {
  blue:   { dark: '#4f9cff', light: '#2463eb' },
  cyan:   { dark: '#22d3ee', light: '#0891b2' },
  purple: { dark: '#a78bfa', light: '#7c3aed' },
  pink:   { dark: '#f472b6', light: '#db2777' },
  orange: { dark: '#fb923c', light: '#ea580c' },
  green:  { dark: '#4ade80', light: '#16a34a' },
}

export const DARK_VARS: Record<string, string> = {
  '--bg-deep':     '#09090c',
  '--bg-base':     '#0d0e13',
  '--bg-elevated': '#13141a',
  '--bg-hover':    '#191b23',
  '--border':      '#1f2130',
  '--text-1':      '#e2e4ed',
  '--text-2':      '#7a829a',
  '--text-3':      '#646b85',
}

export const LIGHT_VARS: Record<string, string> = {
  '--bg-deep':     '#eef0f6',
  '--bg-base':     '#f6f7fb',
  '--bg-elevated': '#ffffff',
  '--bg-hover':    '#eeeff8',
  '--border':      '#dde0ec',
  '--text-1':      '#181a28',
  '--text-2':      '#525870',
  '--text-3':      '#808899',
}
