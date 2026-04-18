import type { GlobalThemeOverrides } from 'naive-ui'

export function buildOverrides(
  accent: string,
  isDark: boolean,
  fontSize = '13px',
): GlobalThemeOverrides {
  const bg0 = isDark ? '#09090c' : '#eef0f6'
  const bg1 = isDark ? '#0d0e13' : '#f6f7fb'
  const bg2 = isDark ? '#13141a' : '#ffffff'
  const bg3 = isDark ? '#191b23' : '#eeeff8'
  const bdr = isDark ? '#1f2130' : '#dde0ec'
  const t1  = isDark ? '#e2e4ed' : '#181a28'
  const t2  = isDark ? '#7a829a' : '#525870'
  const t3  = isDark ? '#646b85' : '#808899'
  // accent-bg: 12% opacity dark (0x1f ≈ 12%), 9% opacity light (0x17 ≈ 9%)
  const abg = isDark ? `${accent}1f` : `${accent}17`
  const sm  = `${parseInt(fontSize) - 1}px`

  return {
    common: {
      bodyColor:            bg1,
      primaryColor:         accent,
      primaryColorHover:    accent,
      primaryColorPressed:  accent,
      primaryColorSuppl:    accent,
      cardColor:            bg2,
      popoverColor:         bg2,
      modalColor:           isDark ? 'rgba(9,9,12,0.88)' : 'rgba(246,247,251,0.90)',
      textColor1:           t1,
      textColor2:           t2,
      textColor3:           t3,
      dividerColor:         bdr,
      borderColor:          bdr,
      inputColor:           bg2,
      tableColor:           bg1,
      tableHeaderColor:     bg2,
      hoverColor:           bg3,
      borderRadius:         '6px',
      borderRadiusSmall:    '4px',
      fontFamily:           "'Inter', -apple-system, 'Segoe UI', sans-serif",
      fontFamilyMono:       "'JetBrains Mono', 'Fira Code', monospace",
      fontSizeMedium:       fontSize,
      fontSizeSmall:        sm,
    },
    DataTable: {
      tdColor:      bg1,
      tdColorHover: bg3,
      thColor:      bg2,
      thTextColor:  t2,
      tdTextColor:  t1,
      borderColor:  bdr,
      thFontWeight: '500',
    },
    Menu: {
      color:                    bg0,
      itemColorHover:           bg3,
      itemColorActive:          abg,
      itemColorActiveHover:     abg,
      itemTextColorActive:      accent,
      itemIconColorActive:      accent,
      itemTextColorActiveHover: accent,
      itemIconColorActiveHover: accent,
    },
    Input: {
      color:          bg2,
      colorFocus:     bg2,
      border:         `1px solid ${bdr}`,
      borderFocus:    `1px solid ${accent}`,
      boxShadowFocus: `0 0 0 2px ${abg}`,
    },
    Button: {
      colorPrimary:       accent,
      borderRadiusMedium: '6px',
      borderRadiusSmall:  '4px',
    },
    Tag:    { borderRadius: '4px' },
    Modal:  {
      boxShadow: isDark
        ? '0 24px 48px rgba(0,0,0,0.6)'
        : '0 8px 32px rgba(0,0,0,0.14)',
    },
    Tabs: {
      colorSegment:           bg2,
      tabColorSegment:        bg3,
      tabTextColorActiveBar:  accent,
      barColor:               accent,
    },
    Slider: {
      fillColor: accent,
      dotColor:  accent,
    },
  }
}
