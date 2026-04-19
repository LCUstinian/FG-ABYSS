# FG-ABYSS UI Design Spec

> **Version**: 1.8  
> **Date**: 2026-04-18  
> **Status**: Active  
> **Iteration rounds**: 9 (brainstorm R1 + 5× ui-ux-pro-max + 3× detail polish)  
> **Style**: Cyber Professional — modern geek security tool, dark-first, light theme equally polished

---

## 1. Design Language

### 1.1 Visual Direction

"Cyber Professional" — high-contrast near-black backgrounds, cold blue/cyan accent, monospace fonts for all technical data. Inspired by Caido, Warp Terminal, Raycast. Neither a traditional ugly security tool nor a generic SaaS product. Both dark and light themes are first-class citizens.

### 1.2 Color System

#### Dark Theme Tokens

| Token | Value | Usage |
|-------|-------|-------|
| `--bg-deep` | `#09090c` | Sidebar, titlebar |
| `--bg-base` | `#0d0e13` | Main content area |
| `--bg-elevated` | `#13141a` | Cards, panels |
| `--bg-hover` | `#191b23` | Row hover, active states |
| `--border` | `#1f2130` | Dividers, borders |
| `--text-1` | `#e2e4ed` | Primary text |
| `--text-2` | `#7a829a` | Secondary text |
| `--text-3` | `#646b85` | Disabled, muted (≥3:1 contrast on bg-base) |
| `--accent` | `#4f9cff` | Primary actions (default blue) |
| `--accent-bg` | `rgba(79,156,255,0.12)` | Selected rows, focus ring fill |

#### Light Theme Tokens

| Token | Value | Usage |
|-------|-------|-------|
| `--bg-deep` | `#eef0f6` | Sidebar, titlebar |
| `--bg-base` | `#f6f7fb` | Main content area |
| `--bg-elevated` | `#ffffff` | Cards, panels |
| `--bg-hover` | `#eeeff8` | Row hover |
| `--border` | `#dde0ec` | Dividers, borders |
| `--text-1` | `#181a28` | Primary text |
| `--text-2` | `#525870` | Secondary text |
| `--text-3` | `#808899` | Disabled, muted (≥3:1 contrast on bg-base) |
| `--accent` | `#2463eb` | Primary actions (default blue) |
| `--accent-bg` | `rgba(36,99,235,0.09)` | Selected rows, focus ring fill |

#### Accent Color Presets (6)

| Name | Dark | Light |
|------|------|-------|
| Blue (default) | `#4f9cff` | `#2463eb` |
| Cyan | `#22d3ee` | `#0891b2` |
| Purple | `#a78bfa` | `#7c3aed` |
| Pink | `#f472b6` | `#db2777` |
| Orange | `#fb923c` | `#ea580c` |
| Green | `#4ade80` | `#16a34a` |

#### Semantic Status Colors (theme-invariant)

| State | Color |
|-------|-------|
| Success / Active | `#4ade80` |
| Warning | `#fbbf24` |
| Error / Danger | `#f87171` |
| Info | `#60a5fa` |

#### WebShell Type Badge Colors

| Type | Dark bg | Dark text | Light bg | Light text |
|------|---------|-----------|----------|------------|
| PHP | `rgba(79,156,255,0.15)` | `#4f9cff` | `rgba(36,99,235,0.10)` | `#2463eb` |
| JSP | `rgba(251,146,60,0.15)` | `#fb923c` | `rgba(234,88,12,0.10)` | `#ea580c` |
| ASP | `rgba(167,139,250,0.15)` | `#a78bfa` | `rgba(124,58,237,0.10)` | `#7c3aed` |
| ASPX | `rgba(34,211,238,0.15)` | `#22d3ee` | `rgba(8,145,178,0.10)` | `#0891b2` |

### 1.3 Typography

- **UI Font**: `'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif`
- **Mono Font**: `'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace`

| Scale | Size | Weight | Usage |
|-------|------|--------|-------|
| xs | 12px | 400 | Timestamps, captions, statusbar |
| sm | 13px | 400 | Default UI text, table cells |
| md | 14px | 500 | Section labels, form labels |
| lg | 16px | 600 | Panel headings |
| xl | 20px | 700 | Stat card numbers |

**Mono font rules** (from Terminal CLI Monospace spec):
- Weight: always `400` — bold ruins monospace character
- Letter-spacing: `0` (monospaced font has built-in auto-spacing)
- Sizes: `12px / 13px / 14px` only — no 11px or 15px for mono fields
- Line-height: `1.4` for inline data, `1.2` for terminal/code blocks (information density)

**UI font line-heights**:
- Dense UI (table rows, form labels): `1.4`
- Readable content (descriptions, tooltips): `1.6`
- Headings: `1.25`

**Mono font is mandatory** for all technical data — see Section 5.4 for full list.

### 1.4 Spacing (4px grid — strictly enforced)

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Icon gaps, badge padding |
| sm | 8px | Within-component spacing |
| md | 16px | Card padding, section gaps |
| lg | 24px | Panel padding |
| xl | 32px | View-level spacing |
| xxl | 48px | Large section separators |

No 5px, 7px, 9px, or any non-4x multiple anywhere.

### 1.5 Border Radius

| Context | Value |
|---------|-------|
| Input fields, table cells | 4px |
| Cards, panels, buttons | 6px |
| Modal dialogs | 8px |
| Tags, badges | 4px |
| Status dots, avatars | 50% |

### 1.6 Visual Effects

- **Modal/Drawer backdrop**: `backdrop-filter: blur(16px)` + `rgba(9,9,12,0.88)` dark / `rgba(246,247,251,0.90)` light; modal scrim is a `<div>` with `position: fixed; inset: 0; z-index: 100` and these bg values — blur alone is insufficient over colorful content
- **Focus ring**: `box-shadow: 0 0 0 2px var(--accent-bg), 0 0 0 1px var(--accent)`
- **Dark mode card elevation**: `border: 1px solid rgba(255,255,255,0.08)` — hairline glass edge, more perceptually premium than `--border` (`#1f2130`) which blends into the dark base at low luminance
- **Light mode card elevation**: `box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05)` — shadow provides lift; border ensures definition on white-on-white surfaces
- **Code areas**: always dark (`var(--code-bg)` bg, `var(--code-border)` border) regardless of theme
- **Custom scrollbar**: 4px wide, `--border` thumb, hover `--text-3`, 2px radius

### 1.7 Semantic CSS Variables (theme-invariant, defined in `:root`)

These never change between dark/light themes:

```css
:root {
  --color-success: #4ade80;
  --color-warning: #fbbf24;
  --color-error:   #f87171;
  --color-info:    #60a5fa;
  --code-bg:       #0c0d11;
  --code-border:   #1a1c26;
  --wc-close-bg:   #ef4444;  /* window close button hover */
}
```

---

## 2. Theme System

### 2.1 Architecture

```
ThemeStore (Pinia)
  ├── mode: 'dark' | 'light' | 'system'
  ├── accentKey: AccentKey
  ├── resolvedMode: computed → 'dark' | 'light'
  ├── naiveTheme: computed → darkTheme | null
  └── themeOverrides: computed → GlobalThemeOverrides

App.vue
  └── <n-config-provider :theme :theme-overrides :locale>
        └── <n-message-provider>
              └── <n-dialog-provider>
                    └── <RouterView />

useThemeCssVars() composable
  └── watchEffect → document.documentElement CSS variables sync
```

### 2.2 ThemeStore

```typescript
// src/stores/theme.ts
import { defineStore } from 'pinia'
import { useOsTheme, darkTheme } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'

export type ThemeMode = 'dark' | 'light' | 'system'
export type AccentKey = 'blue' | 'cyan' | 'purple' | 'pink' | 'orange' | 'green'
export type FontSize = '12px' | '13px' | '14px' | '15px'

export const ACCENT_COLORS: Record<AccentKey, { dark: string; light: string }> = {
  blue:   { dark: '#4f9cff', light: '#2463eb' },
  cyan:   { dark: '#22d3ee', light: '#0891b2' },
  purple: { dark: '#a78bfa', light: '#7c3aed' },
  pink:   { dark: '#f472b6', light: '#db2777' },
  orange: { dark: '#fb923c', light: '#ea580c' },
  green:  { dark: '#4ade80', light: '#16a34a' },
}

export const useThemeStore = defineStore('theme', () => {
  const osTheme = useOsTheme()
  const mode = useLocalStorage<ThemeMode>('fg-theme-mode', 'dark')
  const accentKey = useLocalStorage<AccentKey>('fg-theme-accent', 'blue')
  const fontSize = useLocalStorage<FontSize>('fg-theme-fontsize', '13px')

  const resolvedMode = computed<'dark' | 'light'>(() =>
    mode.value === 'system' ? (osTheme.value ?? 'dark') : mode.value
  )
  const naiveTheme = computed(() =>
    resolvedMode.value === 'dark' ? darkTheme : null
  )
  const themeOverrides = computed(() =>
    buildOverrides(
      ACCENT_COLORS[accentKey.value][resolvedMode.value],
      resolvedMode.value === 'dark',
      fontSize.value,
    )
  )

  function setMode(m: ThemeMode) { mode.value = m }
  function setAccent(k: AccentKey) { accentKey.value = k }
  function setFontSize(s: FontSize) { fontSize.value = s }

  return { mode, accentKey, fontSize, resolvedMode, naiveTheme, themeOverrides, setMode, setAccent, setFontSize }
})
```

### 2.3 Naive UI themeOverrides Builder

```typescript
// src/theme/overrides.ts
import type { GlobalThemeOverrides } from 'naive-ui'

export function buildOverrides(accent: string, isDark: boolean, fontSize = '13px'): GlobalThemeOverrides {
  const bg0   = isDark ? '#09090c' : '#eef0f6'
  const bg1   = isDark ? '#0d0e13' : '#f6f7fb'
  const bg2   = isDark ? '#13141a' : '#ffffff'
  const bg3   = isDark ? '#191b23' : '#eeeff8'
  const bdr   = isDark ? '#1f2130' : '#dde0ec'
  const t1    = isDark ? '#e2e4ed' : '#181a28'
  const t2    = isDark ? '#7a829a' : '#525870'
  const t3    = isDark ? '#646b85' : '#808899'
  const abg   = isDark
    ? `${accent}1f`   // 12% opacity
    : `${accent}17`   // 9% opacity

  return {
    common: {
      bodyColor: bg1,
      primaryColor: accent,
      primaryColorHover: accent,
      primaryColorPressed: accent,
      primaryColorSuppl: accent,
      cardColor: bg2,
      popoverColor: bg2,
      modalColor: isDark ? 'rgba(9,9,12,0.88)' : 'rgba(246,247,251,0.90)',
      textColor1: t1,
      textColor2: t2,
      textColor3: t3,
      dividerColor: bdr,
      borderColor: bdr,
      inputColor: bg2,
      tableColor: bg1,
      tableHeaderColor: bg2,
      hoverColor: bg3,
      borderRadius: '6px',
      borderRadiusSmall: '4px',
      fontFamily: "'Inter', -apple-system, 'Segoe UI', sans-serif",
      fontFamilyMono: "'JetBrains Mono', 'Fira Code', monospace",
      fontSizeMedium: fontSize,
      fontSizeSmall: `${parseInt(fontSize) - 1}px`,
    },
    DataTable: {
      tdColor: bg1,
      tdColorHover: bg3,
      thColor: bg2,
      thTextColor: t2,
      tdTextColor: t1,
      borderColor: bdr,
      thFontWeight: '500',
    },
    Menu: {
      color: bg0,
      itemColorHover: bg3,
      itemColorActive: abg,
      itemColorActiveHover: abg,
      itemTextColorActive: accent,
      itemIconColorActive: accent,
      itemTextColorActiveHover: accent,
      itemIconColorActiveHover: accent,
    },
    Input: {
      color: bg2,
      colorFocus: bg2,
      border: `1px solid ${bdr}`,
      borderFocus: `1px solid ${accent}`,
      boxShadowFocus: `0 0 0 2px ${abg}`,
    },
    Button: {
      colorPrimary: accent,
      borderRadiusMedium: '6px',
      borderRadiusSmall: '4px',
    },
    Tag: { borderRadius: '4px' },
    Modal: {
      boxShadow: isDark
        ? '0 24px 48px rgba(0,0,0,0.6)'
        : '0 8px 32px rgba(0,0,0,0.14)',
    },
    Tabs: {
      colorSegment: bg2,
      tabColorSegment: bg3,
      tabTextColorActiveBar: accent,
      barColor: accent,
    },
    Slider: {
      fillColor: accent,
      dotColor: accent,
    },
  }
}
```

### 2.4 CSS Variables Sync

```typescript
// src/composables/useThemeCssVars.ts
import { watchEffect } from 'vue'
import { useThemeStore, ACCENT_COLORS } from '@/stores/theme'

const DARK_VARS: Record<string, string> = {
  '--bg-deep':     '#09090c',
  '--bg-base':     '#0d0e13',
  '--bg-elevated': '#13141a',
  '--bg-hover':    '#191b23',
  '--border':      '#1f2130',
  '--text-1':      '#e2e4ed',
  '--text-2':      '#7a829a',
  '--text-3':      '#646b85',
  '--font-mono':   "'JetBrains Mono', 'Fira Code', monospace",
}

const LIGHT_VARS: Record<string, string> = {
  '--bg-deep':     '#eef0f6',
  '--bg-base':     '#f6f7fb',
  '--bg-elevated': '#ffffff',
  '--bg-hover':    '#eeeff8',
  '--border':      '#dde0ec',
  '--text-1':      '#181a28',
  '--text-2':      '#525870',
  '--text-3':      '#808899',
  '--font-mono':   "'JetBrains Mono', 'Fira Code', monospace",
}

export function useThemeCssVars() {
  const store = useThemeStore()

  watchEffect(() => {
    const isDark = store.resolvedMode === 'dark'
    const accent = ACCENT_COLORS[store.accentKey][store.resolvedMode]
    const vars = isDark ? DARK_VARS : LIGHT_VARS
    const root = document.documentElement

    root.setAttribute('data-theme', store.resolvedMode)
    Object.entries(vars).forEach(([k, v]) => root.style.setProperty(k, v))
    root.style.setProperty('--accent', accent)
    root.style.setProperty('--accent-bg', isDark ? `${accent}1f` : `${accent}17`)
  })
}
// Called once in App.vue setup()
```

---

## 3. Layout Architecture

### 3.1 App Shell Grid

```
┌───────────────────────────────────────────┐  48px  titlebar
├──────────┬────────────────────────────────┤
│          │                                │
│ Sidebar  │        <RouterView />          │  flex-1
│ 200px    │                                │
├──────────┴────────────────────────────────┤  28px  statusbar
└───────────────────────────────────────────┘
```

```css
/* src/styles/app-shell.css */

/* Outer grid: rows only — sidebar width managed by flex below */
.app-shell {
  display: grid;
  grid-template-rows: 48px 1fr 28px;
  height: 100vh;
  overflow: hidden;
  background: var(--bg-base);
}

.titlebar-area  { grid-row: 1; }
.statusbar-area { grid-row: 3; }

/* Middle row: flex so sidebar width transition works correctly.
   CSS Grid cannot animate grid-template-columns on a CSS variable change. */
.app-main {
  grid-row: 2;
  display: flex;
  overflow: hidden;
}

.sidebar-panel {
  flex-shrink: 0;
  width: 200px;
  transition: width 220ms cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.sidebar-panel.is-collapsed {
  width: 64px;
}

.content-area {
  flex: 1;
  overflow: hidden;
}
```

### 3.2 CustomTitlebar

**Layout**:
```
[◈ logo] 渊渟  ABYSS         [☀/🌙]  [ZH/EN]  [─] [□] [✕]
←——————— data-tauri-drag-region ————————→ ←— no drag —→
```

**Specs**:
- Full 48px height, `background: var(--bg-deep)`, bottom 1px `var(--border)`
- Logo: 24×24px SVG, accent color
- App name "渊渟": 13px semi-bold, `--text-1`
- Subtitle "ABYSS": 11px, `--text-3`, 6px left margin
- Theme toggle + lang toggle: 32×32px icon buttons, `--bg-hover` on hover
- Theme toggle cycles: `dark → light → system`, tooltip shows current mode
- Window controls: 40×32px (min/max), 48×32px (close), no border-radius, flush to edge
  - Close hover: `#ef4444` bg, white icon
  - Min/Max hover: `var(--bg-hover)`
- `data-tauri-drag-region="false"` on the right controls group

```vue
<!-- src/components/layout/CustomTitlebar.vue -->
<div class="titlebar" data-tauri-drag-region>
  <div class="titlebar-left" data-tauri-drag-region>
    <AppLogo class="logo" />
    <span class="app-name">渊渟</span>
    <span class="app-sub">ABYSS</span>
  </div>
  <div class="titlebar-right" data-tauri-drag-region="false">
    <ThemeToggle />
    <LangToggle />
    <div class="win-controls">
      <button class="wc-min"   @click="minimize">&#x2212;</button>
      <button class="wc-max"   @click="toggleMaximize">&#x25A1;</button>
      <button class="wc-close" @click="close">&#x2715;</button>
    </div>
  </div>
</div>
```

### 3.3 Sidebar

**Layout** (200px expanded / 64px collapsed):
```
┌──────────────────┐
│ ▐ ◈  首页        │  ← active: 2px left accent border + --accent-bg
│   📁  项目       │  ← hover: --bg-hover
│   📦  载荷       │
│   🧩  插件       │
│   ⚙   设置       │
│                  │
│ [‹‹] 折叠        │  ← bottom toggle
└──────────────────┘
```

**Specs**:
- Background: `--bg-deep`, right border: 1px `--border`
- Nav item: height 40px, `padding: 0 12px`, `gap: 10px`, icon 18px, text 13px
- All items: `border-left: 2px solid transparent` (permanent 2px placeholder, prevents layout jitter on activation)
- Active: `border-left-color: var(--accent)`, `background: var(--accent-bg)`, text/icon `--accent`
- Hover: `background: var(--bg-hover)`
- Collapsed: icons only, centered; text hidden via `opacity: 0` + `width: 0` in 200ms
- Collapsed active: icon stays accent colored
- Collapsed tooltip: `n-tooltip` placement="right", delay 400ms
- Collapse button: bottom of sidebar, chevron icon rotates on collapse

```typescript
// Collapse state persisted
const sidebarCollapsed = useLocalStorage('fg-sidebar-collapsed', false)
```

### 3.4 StatusBar

**Layout**:
```
● 已连接  │  项目 5  │  Shell 12  │  载荷 25
```

**Specs**:
- Height: 28px, `background: var(--bg-deep)`, top 1px `var(--border)`
- Padding: `0 16px`, gap: 8px, font-size: 12px, color: `--text-2`
- Status dot: 8px circle
  - Active shells exist → `#4ade80` + CSS pulse animation (2.4s infinite)
  - No connections → `--text-3` static dot
- Numbers: `font-variant-numeric: tabular-nums`, mono font
- Vertical dividers: `--border` color, 12px height, `1px` wide

```css
@keyframes pulse-dot {
  0%, 100% { box-shadow: 0 0 0 0 rgba(74, 222, 128, 0.4) }
  50%       { box-shadow: 0 0 0 5px rgba(74, 222, 128, 0) }
}
.status-dot-active { animation: pulse-dot 2.4s ease-in-out infinite; }
```

### 3.5 Route Transitions

```vue
<!-- src/layouts/MainLayout.vue -->
<RouterView v-slot="{ Component }">
  <Transition name="fade" mode="out-in">
    <component :is="Component" :key="$route.fullPath" />
  </Transition>
</RouterView>
```

```css
.fade-enter-active, .fade-leave-active { transition: opacity 120ms ease }
.fade-enter-from,   .fade-leave-to     { opacity: 0 }
```

---

## 4. View Designs

### 4.1 HomeView (`/`)

```
┌────────────────────────────────────────────────────┐
│  欢迎回来                          今天 2026/04/18  │
├───────────┬───────────┬───────────┬────────────────┤
│ [活跃图标]  │ [项目图标] │ [载荷图标] │ [请求图标]    │
│ 活跃 Shell │ 项目       │ 载荷       │ 今日请求       │
│     3     │    5      │    25     │    142         │
├───────────┴───────────┴───────────┴────────────────┤
│  最近连接                                           │
│  ● shell-prod    192.168.1.1/shell.php  PHP  2m    │
│  ● shell-test    10.0.0.5/cmd.jsp       JSP  1h    │
│  ○ shell-dev     172.16.0.1/web.aspx    ASPX  3d   │
└────────────────────────────────────────────────────┘
```

**Stat cards**: equal-width, `--bg-elevated`, left 3px `--accent` border decoration  
**Card icons**: SVG line icons from `lucide-vue-next` — `Terminal` / `FolderOpen` / `Package` / `Zap`; 20px, `--accent` color  
**Numbers**: 20px semi-bold, mono, `font-variant-numeric: tabular-nums`  
**Recent list**: compact, status dot + name + URL (mono, `--text-2`) + type badge + relative time  
**Empty state (first launch)**: centered illustration + "添加第一个 WebShell" primary button

### 4.2 ProjectView (`/project`)

**Layout**: 240px project tree (left) + flex-1 shell table (right)

```
┌────────────┬──────────────────────────────────────────┐
│ + 新建项目  │ 🔍 搜索...  [PHP▾] [状态▾]  [+ 添加Shell]│
├────────────┤──────────────────────────────────────────┤
│ ▼ 全部      │ ☐ │ 状 │ 名称        │ URL              │
│   项目A  5  │───┼────┼─────────────┼──────────────────│
│   项目B  2  │   │ ●  │ shell-prod  │ https://xxx.com  │
│   项目C  1  │   │ ●  │ shell-test  │ http://10.0.0.5  │
│             │   │ ○  │ shell-dev   │ http://172...    │
│ ──────────  │──────────────────────────────────────────│
│ 🗑 回收站 3  │  < 1 2 3 >           共 8 条             │
└────────────┴──────────────────────────────────────────┘
```

**Project tree**:
- Background `--bg-deep`, right 1px `--border`
- Item: 40px height, 16px left padding, hover `--bg-hover`, active `--accent-bg`
- Shell count badge: right-aligned, `--text-3`, 12px
- "回收站": bottom of list, divider separator, click to filter table to soft-deleted shells

**Shell table columns**:
| Column | Width | Notes |
|--------|-------|-------|
| Checkbox | 40px | Batch select |
| Status dot | 40px | 🟢 active / ⚪ offline / 🔴 error |
| Name | 160px | `--text-1` |
| URL | flex | Mono font, `--text-2`, ellipsis + hover tooltip for full URL |
| Type | 80px | PHP/JSP/ASP/ASPX badge with type colors |
| Last connected | 120px | Relative time, `--text-2` |
| Tags | 120px | Max 2 tags visible + "+N" overflow |
| Actions | 80px | Hover-reveal icons (connect, edit, delete) |

**Row actions keyboard accessibility** (critical — hover-reveal is mouse-only):
- Row element has `tabIndex="0"`; `↑↓` arrow keys move between rows; `Tab` moves to next focusable region
- `Space` or `Enter` on focused row opens the context menu (same options as right-click `n-dropdown`)
- Hover-reveal action icons are **not** individually Tab-reachable — the context menu is the sole keyboard path
- Disabled/deleted rows: `tabIndex="-1"`, excluded from arrow navigation
- `aria-selected` reflects checkbox state; row `role="row"` inside `role="grid"` for screen readers

**Row states**:
- Selected: `--accent-bg` background
- Hover: `--bg-hover`
- Soft-deleted (recycle view): `--text-3`, italic, reduced opacity 0.6

**Row context menu** (`n-dropdown`):
```
🖥  连接控制台        ← --accent color
📋  复制 URL
✏️  编辑
🔄  测试连接
📤  导出配置
──────────────
🗑  移入回收站        ← #f87171 color
```

**Batch action bar** (slides up from bottom when N rows selected):
- Height 48px, `--bg-elevated`, top border, floating above statusbar
- Position: `position: absolute; bottom: 0; left: 0; right: 0` inside `.content-area` (which has `position: relative`)
- Content: "已选 N 个 │ 批量测试 │ 批量导出 │ 移入回收站"
- Slide-up transition: `translateY(48px → 0)`, 200ms ease-out; hidden with `translateY(48px)` + `pointer-events: none`

### 4.3 PayloadView (`/payload`)

**Layout**: 340px config panel (left) + flex-1 code preview (right)

```
┌──────────────────┬──────────────────────────────────┐
│ 语言              │ [复制] [下载] [查看混淆对比]       │
│ [PHP][JSP][ASP]  │ ┌──────────────────────────────┐ │
│                  │ │<?php                         │ │
│ 加密              │ │  $k="a8f3...";               │ │
│ ◉ AES-256-GCM    │ │  ...                         │ │
│ ○ XOR            │ │  (永远深色)                   │ │
│ ○ 无加密          │ └──────────────────────────────┘ │
│                  │  PHP 5.x-8.x │ AES-256 │ 1.4 KB  │
│ 密钥              │                                  │
│ [••••••••••] 👁  │                                  │
│                  │                                  │
│ 混淆等级          │                                  │
│ L0──●────────L5  │                                  │
│ ☑ 变量重命名      │                                  │
│ ☑ 字符串加密      │                                  │
│ ☐ 垃圾代码        │                                  │
│                  │                                  │
│ [    生成载荷   ] │                                  │
└──────────────────┴──────────────────────────────────┘
```

**Left config panel**:

All input groups have a **visible label above** (14px, `--text-2`, `font-weight: 500`) — never placeholder-only labels:

- **语言** label → `n-radio-group` button style, PHP / JSP / ASP / ASPX; selected button: `--accent` bg + white text
- **加密方式** label → `n-radio-group` vertical list; each option: radio + name + brief tag (e.g. "推荐", "兼容")
- **密钥** label → row: `n-input` (mono font, `type="password"`, 36px) + `EyeIcon` toggle + `CopyIcon` copy button
  - Copy button: `aria-label="复制密钥"`, copies plaintext key → icon flips to `CheckIcon` 1.5s then reverts
  - Password strength bar: 4-segment bar below input, colored red/orange/yellow/green by entropy score
- **混淆等级** label → row layout: `n-slider` (flex-1) + level badge `"L2"` (28px pill, mono, `--accent-bg`, `--accent`)
  - Slider steps 0–5 with tick marks; label updates live as handle moves
  - Below slider: checkboxes (☑ 变量重命名, ☑ 字符串加密, ☐ 垃圾代码); all disabled + opacity 0.38 at L0
- Generate button: bottom, full-width, primary, 40px height
  - **Loading state**: `:loading="loadingMap.generate"` — inline spinner, button disabled
  - **Success state**: button briefly flashes green (`#4ade80` bg) for 600ms, code panel updates
  - **Error state**: button gets red outline (`border: 1px solid #f87171`, bg reverts to default), `n-message` toast appears top-right with specific error text; outline clears after 3s or on next keystroke

**Right code panel**:
- Background always `#0c0d11`, border `1px solid #1a1c26`
- Syntax highlight: `shiki` with One Dark Pro theme
- Line numbers: `--text-3`
- Toolbar: copy (→ check icon 1.5s, no toast), download, obfuscation diff toggle
- Status bar: file type + encryption + size, right-aligned, 12px `--text-2`
- Obfuscation diff: two-column view, original left / obfuscated right

### 4.4 ConsoleView (standalone window)

Window label: `console-{webshell_id}`, `decorations: false`, min size 900×600px

```
┌──────────────────────────────────────────────────────┐  32px header
│ ◈ shell-prod  https://example.com/shell.php  🟢  [✕] │
├──────────────────────────────────────────────────────┤  40px tab bar
│ [📄 文件管理] [🗄️ 数据库] [⌨️ 终端] [🧩 plugin-1]   │
├──────────────────────────────────────────────────────┤
│                                                      │
│              Tab content area                        │  flex-1
│                                                      │
└──────────────────────────────────────────────────────┘
```

**Header**:
- Background `--bg-deep`, bottom 1px `--border`
- Shell name `--text-1`, URL `--text-2` mono, status dot with pulse if active
- Close button right-aligned, hover `#ef4444` (matches main window)

**Tab bar** (`n-tabs` type="card", custom styled):
- Inactive tab: `--bg-deep` bg, hover `--bg-hover`
- Active tab: `--bg-base` bg + 2px `--accent` bottom border
- Plugin tabs: dynamically appended, same style

**File Manager Tab**:
```
[/var/www/html] › html › uploads          [列表▾] [刷新]
──────────────────────────────────────────────────────
☐  名称 ↑          大小     权限      修改时间
   📁  uploads     —        drwxr-x   2024-03-01
   📄  index.php   4.2 KB   -rw-r--   2024-02-28
   📄  shell.php   1.1 KB   -rw-r--   2024-04-18
```
- Breadcrumb path: clickable segments, mono font
- File name: mono font
- Permissions: mono font, color-coded (writable = warning color)
- Row context menu: 查看 / 编辑 / 下载 / 重命名 / 删除
- Bottom drop zone: "拖拽文件到此上传", dashed border, activates on dragover

**Terminal Tab**:
- Full-height xterm.js instance
- Font: JetBrains Mono 14px, bg `#0c0d11`, cursor block in accent color
- Toolbar (top): clear screen / history toggle / copy selection
- Command history: Up/Down arrow navigation

**Database Tab**:
```
[MySQL ▾]  [host]  [port]  [user]  [password]  [连接]
──────────────────────────────────────────────────────
┌── SQL Editor (always dark) ──────────────────────┐
│ SELECT * FROM users LIMIT 100;                   │
│                                                  │
└──────────────────────────────────────────────────┘
[执行 F5]                              [导出 CSV]
──────────── Results ──────────────────────────────
id │ username │ email              │ created_at
1  │ admin    │ admin@example.com  │ 2024-01-01
```
- Connection inputs: mono font for host/port/credentials
- SQL editor: CodeMirror v6 (`@codemirror/lang-sql`), One Dark Pro theme, SQL syntax highlight
- Result table: sortable columns, row count display, `--text-2` for null values

### 4.5 PluginView (`/plugin`)

**Layout**: full-width grid of plugin cards + top toolbar

```
┌──────────────────────────────────────────────────────┐
│ 已安装 (3)    [全部▾]  [启用▾]              [+ 安装]  │
├──────────────────────────────────────────────────────┤
│ ┌──────────────┐  ┌──────────────┐  ┌─────────────┐ │
│ │ 文件管理      │  │ 数据库管理    │  │ 终端         │ │
│ │ v1.2.0  内置  │  │ v1.0.0  内置  │  │ v1.1.0  内置 │ │
│ │ 远程文件操作  │  │ SQL 查询执行  │  │ 命令执行      │ │
│ │ [●] 已启用    │  │ [●] 已启用    │  │ [●] 已启用   │ │
│ └──────────────┘  └──────────────┘  └─────────────┘ │
└──────────────────────────────────────────────────────┘
```

**Toolbar**:
- Filter tabs: 全部 / 已启用 / 已禁用 / 内置 / 第三方
- "+ 安装" button: opens Drawer for installing from file or URL

**Plugin card** (`--bg-elevated`, 6px radius, `padding: 16px`):
- Plugin name: 14px semi-bold, `--text-1`
- Version badge: `--text-3`, 12px mono
- Source badge: "内置" uses `--accent-bg` + `--accent`; "第三方" uses neutral
- Description: 13px, `--text-2`, max 2 lines with ellipsis
- Bottom: `n-switch` (enabled/disabled) left + "卸载" link right (only for third-party)
- Disabled state: card opacity 0.6, switch off

**Empty state**: "还没有安装插件 / 内置插件随应用分发" (no CTA — built-ins always present)

### 4.6 SettingsView (`/settings`)

**Layout**: 180px left submenu + flex-1 right panels

```
┌────────────┬──────────────────────────────────────┐
│ ▐ 外观      │  主题                                │
│   连接      │  ◉ 跟随系统  ○ 深色  ○ 浅色           │
│   安全      │                                      │
│   日志      │  强调色                               │
│   备份      │  🔵 🟦 💜 🩷 🟠 🟢 🟡               │
│   关于      │                                      │
│             │  字体大小   [13px ▾]                  │
│             │  语言       [简体中文 ▾]               │
└────────────┴──────────────────────────────────────┘
```

**Left submenu**: same style as main sidebar (active left accent border, `--bg-deep` bg)

**Right content**: `--bg-base` bg, content in `--bg-elevated` cards, `padding: 24px`, `gap: 16px`

**Appearance panel**:
- Theme mode: 3 card-style radio options (跟随系统/深色/浅色), each with small visual preview thumbnail
- Accent colors: 6 × 32px circles, selected state: white checkmark + 2px accent outline ring
- Font size dropdown: 12 / 13 / 14 / 15px options
- Language dropdown: 简体中文 / English

**Connection panel**:
- Proxy type: None / HTTP / SOCKS5 radio; selected type expands host+port inputs below (progressive disclosure, 200ms transition)
- Timeout: row — `n-slider` (flex-1, 1–60s, default 30s) + value badge `"30s"` (42px, mono, `--text-2`); badge updates live
- Retry count: row — label (flex-1) + `n-input-number` (80px, range 0–5, default 3)

**Security panel**:
- Master password: "修改主密码" button → slides Drawer from right (3-step: current → new → confirm)
- Key rotation: last rotated timestamp + "立即轮换" button
- Audit log toggle: `n-switch`

**About panel**:
- Version: mono font badge
- Links: GitHub, documentation
- License info: `--text-2`, small

---

## 5. Animations & Component Specs

### 5.1 Animation Principles

- Every animation conveys information (loading, state change, hierarchy), no decorative motion
- Duration budget: fast operations 100–150ms, page-level 200–250ms, nothing over 300ms
- **Exit is faster than enter**: leave duration = 60–70% of enter duration (feels responsive)
- Easing: enter `cubic-bezier(0.4, 0, 0.2, 1)` (ease-out), leave `cubic-bezier(0.4, 0, 1, 1)` (ease-in)
- **Spring physics** for modal/drawer enter: `cubic-bezier(0.34, 1.56, 0.64, 1)` — slight overshoot, natural bounce
- **`prefers-reduced-motion`**: all animations must be disabled or reduced when `prefers-reduced-motion: reduce` is active

```css
/* src/styles/animations.css — apply globally */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

- Status dot pulse: keep (functional — conveys live connection), but reduce to static glow when `prefers-reduced-motion`
- Only animate `opacity` and `transform` — never `width`, `height`, `top`, `left` (triggers layout reflow)

### 5.2 Animation Catalog

| Interaction | Enter | Leave | Easing |
|------------|-------|-------|--------|
| Page switch | Fade 120ms | Fade 80ms | ease-out / ease-in |
| Sidebar collapse | 220ms | 220ms | `cubic-bezier(0.4,0,0.2,1)` |
| Row hover | 80ms | 60ms | ease-out |
| Context menu open | Scale 0.95→1 + opacity 100ms | 70ms | ease-out |
| Modal enter | Scale 0.97→1 + opacity 180ms | 120ms | spring `cubic-bezier(0.34,1.56,0.64,1)` |
| Drawer enter | TranslateX 100%→0 + backdrop 200ms | 140ms | ease-out / ease-in |
| Copy success icon | copy→check + rotate(-10°→0) 80ms | stays 1.5s, auto-reverts | spring |
| Connect success | Status dot gray→green + scale 1→1.2→1 | — | 200ms spring |
| Connect fail shake | translateX ±4px × 3 | — | 200ms linear |
| Status dot pulse | Box-shadow 2.4s loop | static glow on reduced-motion | ease-in-out |
| Batch bar slide | TranslateY 48px→0 200ms | 140ms | ease-out |

### 5.3 Component Quick Reference

| Component | Height | Radius | Key Style |
|-----------|--------|--------|-----------|
| Button (primary, md) | 36px | 6px | `--accent` bg, white text; hover: dark `brightness(1.10)`, light `brightness(0.92)` |
| Button (default, md) | 36px | 6px | `--bg-elevated` bg, `--border` border; hover: `--bg-hover` |
| Button (danger) | 36px | 6px | Outline style default; hover: `rgba(248,113,113,0.12)` bg + `#f87171` border |
| Button (sm) | 28px | 4px | Same hover rules, smaller dimensions |
| Input | 36px | 6px | Focus: 2px `--accent-bg` + 1px `--accent` ring |
| Table row | 40px | — | Hover `--bg-hover`, selected `--accent-bg` |
| Sidebar nav item | 40px | — | Active: left 2px `--accent` + `--accent-bg` |
| Card | — | 6px | `--bg-elevated`, 1px `--border` |
| Modal | max-w 560px | 8px | backdrop-blur 16px |
| Drawer | w 380px | — | Slide from right |
| Tooltip | — | 4px | `--bg-elevated`, `padding: 6px 10px` |
| Tag/Badge | 20px | 4px | Type-specific colors (see Section 1.2) |
| Status dot | 8px | 50% | Active green pulse / offline gray / error red |
| Scrollbar | 4px wide | 2px | Thumb `--border`, hover `--text-3` |

### 5.4 Mandatory Mono Font Fields

Every instance of the following data types **must** render in `font-family: var(--font-mono)`:

- WebShell URL and IP addresses (tables, headers, status)
- Password / key fields (payload config, database connection)
- Session IDs, hash values, tokens
- File paths (breadcrumb navigation, file manager columns)
- Payload code preview (entire block)
- Terminal all output
- Database SQL editor input + query result values
- StatusBar counters
- Version strings (About panel)

### 5.5 Disabled State Spec

Disabled elements use consistent treatment across all components:
- Opacity: `0.38` (Material Design standard — enough to show state, low enough to signal inactivity)
- Cursor: `not-allowed`
- No hover/focus effects; `pointer-events: none`
- The `disabled` HTML attribute must be set (not just visual) for screen readers

In Naive UI: use `:disabled="true"` prop on all interactive components; Naive UI applies opacity automatically but verify it matches `0.38`.

### 5.6 Loading State Spec

| Situation | Pattern |
|-----------|---------|
| Table fetching data | `n-data-table` `:loading="loadingMap.list"` — shows built-in spinner overlay |
| Initial page load (>300ms) | Skeleton rows: gray shimmer blocks at expected row height (40px), 3–5 placeholder rows |
| Button async action | `:loading` prop on `n-button` — shows inline spinner, auto-disables |
| Panel/card fetching | `n-spin` wrapper with `size="small"`, `stroke-color="var(--accent)"` |

**Skeleton shimmer CSS**:
```css
@keyframes skeleton-shimmer {
  0%   { background-position: -400px 0 }
  100% { background-position: 400px 0 }
}
.skeleton-row {
  background: linear-gradient(90deg,
    var(--bg-elevated) 25%,
    var(--bg-hover) 50%,
    var(--bg-elevated) 75%);
  background-size: 800px 40px;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
  border-radius: 4px;
}
```

### 5.7 Cursor Rules

| Element | Cursor |
|---------|--------|
| All clickable buttons, links, rows | `cursor: pointer` |
| Disabled elements | `cursor: not-allowed` |
| Text inputs, textareas | `cursor: text` |
| Resizable panels, drag handles | `cursor: col-resize` / `cursor: grab` |
| Window controls (min/max/close) | `cursor: default` (system controls, not pointer) |
| Sidebar collapse button | `cursor: pointer` |
| Table column resize handle | `cursor: col-resize` |

### 5.8 Empty State Spec

All empty lists/tables use this layout:

```
         [SVG icon, 64px, --text-3 fill]

        还没有 [资源名]
      [One-sentence explanation]
         [+ Primary CTA Button]
```

- Icon: line-style SVG matching the feature, `--text-3` fill
- Title: `--text-1`, 14px, centered
- Subtitle: `--text-2`, 13px, centered (optional, omit for recycle bin)
- CTA: only where an action makes sense (no CTA on recycle bin empty state)

### 5.9 Accessibility Spec

**Icon-only buttons require `aria-label`**:

| Button | `aria-label` |
|--------|-------------|
| Minimize window | `"最小化"` |
| Maximize window | `"最大化"` |
| Close window | `"关闭"` |
| Theme toggle | `"切换主题: 深色"` (dynamic) |
| Language toggle | `"切换语言"` |
| Sidebar collapse | `"折叠侧边栏"` / `"展开侧边栏"` |
| Copy button | `"复制代码"` |
| Eye icon (password) | `"显示密码"` / `"隐藏密码"` |

**Color is not the only indicator** — status dots must include text fallback:
- `aria-label="状态: 活跃"` on green dot
- `aria-label="状态: 离线"` on gray dot
- `aria-label="状态: 错误"` on red dot
- StatusBar: `● 已连接` text label next to dot (already spec'd) ✅

**Keyboard navigation**:
- Sidebar nav items: fully `Tab`-navigable, `Enter` activates
- Table rows: `Tab` into first row, `↑↓` arrows navigate rows, `Enter` opens context menu
- Context menu: `↑↓` arrows, `Enter` selects, `Escape` closes
- Modal/Drawer: `Tab` traps focus inside; `Escape` closes; focus returns to trigger after close
- Skip link: `<a href="#main-content" class="skip-link">跳转到主要内容</a>` (visually hidden, visible on focus)

**Focus ring**: `box-shadow: 0 0 0 2px var(--accent-bg), 0 0 0 1px var(--accent)` — never `outline: none` without this replacement

**WCAG contrast targets** (verified):
- `--text-1` on `--bg-base`: ≥15:1 ✅ AAA
- `--text-2` on `--bg-base`: ≥5.2:1 ✅ AA
- `--text-3` on `--bg-base`: ≥3.1:1 ✅ (large text / decorative minimum)
- `--accent` on `--bg-base`: ≥7:1 ✅ AAA
- Light theme same or better ratios ✅

### 5.10 Global Consistency Rules

1. **Hover**: always `--bg-hover` / Naive `hoverColor` — no per-component custom hover colors
2. **Selection/Active**: always `--accent-bg` background + `--accent` text/icon

3. **Mono font**: use the mandatory list in 5.4 — not by feel, by checklist
4. **Code areas**: always dark (`#0c0d11`) regardless of current theme
5. **Spacing**: strictly 4px multiples — 5px, 7px, 9px values are bugs
6. **Border radius**: 4px (small elements), 6px (cards/buttons/inputs), 8px (modals only)
7. **Dividers**: 1px `--border` — no 2px dividers except the sidebar active indicator
8. **No hardcoded colors** in component styles — always use CSS variables or Naive UI tokens
9. **Loading states**: tables use `n-data-table` `:loading="loadingMap.list"` prop; panels use `n-spin` overlay; buttons use `:loading` prop. Use `LoadingMap` (`Record<string, boolean>`) per CLAUDE.md, not a single boolean
10. **Error display policy** (from CLAUDE.md): transient errors → `n-message` toast top-right; unrecoverable/destructive confirmations → `n-dialog` modal; form validation → inline under field
11. **`prefers-reduced-motion`**: all CSS animations and Vue transitions must check the global `@media (prefers-reduced-motion: reduce)` block — see Section 5.1
12. **`cursor: pointer`** on every interactive element; `cursor: not-allowed` on disabled; `cursor: default` on window controls
13. **No emoji as icons** — all iconography uses `lucide-vue-next` SVG components throughout

### 5.11 Light Theme Polish Checklist

Light mode is a first-class citizen, not an inversion. These rules prevent the common "light mode feels washed out" failure:

**Button hover behavior** (light mode is opposite of dark — darken, not lighten):
- Primary button hover: `filter: brightness(0.92)` — darkens `--accent` slightly for visible feedback
- Default button hover: `--bg-hover` (`#eeeff8`) — subtle blue-gray tint, NOT pure white
- Danger button hover: `background: rgba(248,113,113,0.12)` with `#f87171` border
- Never use `brightness(1.1)` on a light mode primary button — it bleaches to near-white

**Sidebar in light mode** — the `--bg-deep` (`#eef0f6`) vs `--bg-base` (`#f6f7fb`) gap is subtle (~7 lightness units). To maintain visual hierarchy:
- Sidebar gets explicit `border-right: 1px solid var(--border)` (`#dde0ec`) — always rendered, not just hover state
- Active nav item background: `--accent-bg` (`rgba(36,99,235,0.09)`) is slightly more saturated than neutral — visible contrast
- Sidebar font weight for active item: `font-weight: 600` in light mode (dark mode `500` is enough due to contrast)

**Code containers in light mode** — code panels (`--code-bg: #0c0d11`) create a "dark island" inside light UI:
- Apply `border-radius: 6px; box-shadow: 0 1px 4px rgba(0,0,0,0.12), 0 0 0 1px rgba(0,0,0,0.08)` to lift the island visually
- The code area header bar (toolbar with copy/download) gets `background: #16181f; border-bottom: 1px solid #1a1c26` — slightly lighter than code-bg for separation without breaking the dark island

**Modal scrim in light mode**:
- Backdrop: `background: rgba(246,247,251,0.90)` + `backdrop-filter: blur(16px)`
- The frosted-glass look in light mode requires the scrim to be a bit more opaque (`0.90`) vs dark (`0.88`) because light backgrounds bleed through blur more

**Card visual rhythm in light mode**:
- Cards use `box-shadow` for lift (from Section 1.6), NOT border alone — border alone makes cards invisible on `#ffffff` content area
- Project tree and Settings submenu: use `background: var(--bg-deep)` (`#eef0f6`) to distinguish from `--bg-base` content

### 5.12 Z-Index Layer Definition

Never use arbitrary z-index values. All layers must come from this table:

| Layer | Value | Elements |
|-------|-------|---------|
| Base | 0 | Normal flow content, table rows |
| Raised | 10 | Sticky table header, fixed sidebar |
| Dropdown | 20 | Context menus, select popups, tooltips, autocomplete |
| Overlay | 40 | Drawer/modal backdrop scrim |
| Modal | 100 | Modal dialogs, drawers, full-screen panels |
| Toast | 1000 | `n-message` toast notifications (always topmost) |

Set `n-config-provider` `:theme-overrides="{ common: { zIndexModal: 100, zIndexPopover: 20, zIndexMessage: 1000 } }"` so Naive UI components respect the same scale.

The batch action bar (Section 4.2) sits at the content layer (`z-index: 10`) — it's inside `.content-area`, not a portal, so it stacks naturally without fighting the Dropdown layer at 20.

---

## 6. File Structure

```
src/
├── styles/
│   ├── app-shell.css          # Grid layout, sidebar transition
│   ├── variables.css          # CSS custom property declarations
│   ├── scrollbar.css          # Custom scrollbar styles
│   └── animations.css         # @keyframes definitions
├── theme/
│   ├── overrides.ts           # buildOverrides() function
│   ├── tokens.ts              # DARK_VARS, LIGHT_VARS, ACCENT_COLORS
│   └── index.ts               # Re-exports
├── stores/
│   └── theme.ts               # useThemeStore
├── composables/
│   └── useThemeCssVars.ts     # CSS variable sync watchEffect
├── components/layout/
│   ├── CustomTitlebar.vue
│   ├── Sidebar.vue
│   ├── StatusBar.vue
│   └── AppShell.vue           # Assembles the grid
└── views/
    ├── HomeView.vue
    ├── ProjectView.vue
    ├── PayloadView.vue
    ├── PluginView.vue
    ├── ConsoleView.vue        # Used in console window
    └── settings/
        ├── SettingsView.vue
        ├── AppearancePanel.vue
        ├── ConnectionPanel.vue
        ├── SecurityPanel.vue
        ├── LoggingPanel.vue
        ├── BackupPanel.vue
        └── AboutPanel.vue
```
