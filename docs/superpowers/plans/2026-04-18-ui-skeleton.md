# FG-ABYSS Frontend UI Skeleton Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Scaffold the complete Tauri v2 + Vue 3 frontend UI shell — theme system, layout, routing, and all view stubs — so that `pnpm tauri dev` renders a fully navigable app with correct visual design and working theme switching, using mock data (no Tauri commands wired yet).

**Architecture:** Feature-based layered Vue 3 SPA inside Tauri v2. Theme system lives in `src/theme/` (tokens + Naive UI overrides) and `src/stores/theme.ts` (Pinia). CSS custom properties are synced by `useThemeCssVars()`. Layout is a CSS Grid shell (3 rows: 48px titlebar / flex-1 main / 28px statusbar) with Flexbox middle row for sidebar collapse animation. Views are functional UI shells with hardcoded mock data — zero `invoke()` calls in this plan.

**Tech Stack:** Tauri v2, Vue 3.4 + TypeScript 5, Naive UI v2, Pinia v2, Vue Router 4, VueUse v11, vue-i18n v10, lucide-vue-next, Vitest v1

**Out of scope:** Rust backend features, Tauri command wiring, Shiki syntax highlighting, xterm.js terminal, CodeMirror SQL editor — all in subsequent plans.

**Design spec:** `docs/superpowers/specs/2026-04-18-fg-abyss-ui-design.md`

---

### Task 1: Project manifest files

**Files:**
- Create: `package.json`
- Create: `vite.config.ts`
- Create: `tsconfig.json`
- Create: `tsconfig.node.json`
- Create: `index.html`

- [ ] **Step 1: Create `package.json`**

```json
{
  "name": "fg-abyss",
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "preview": "vite preview",
    "tauri": "tauri",
    "test": "vitest run",
    "typecheck": "vue-tsc --noEmit"
  },
  "dependencies": {
    "@tauri-apps/api": "^2",
    "@vueuse/core": "^11",
    "naive-ui": "^2",
    "pinia": "^2",
    "vue": "^3.4",
    "vue-i18n": "^10",
    "vue-router": "^4",
    "lucide-vue-next": "^0.400.0"
  },
  "devDependencies": {
    "@tauri-apps/cli": "^2",
    "@vitejs/plugin-vue": "^5",
    "@vue/test-utils": "^2",
    "jsdom": "^24",
    "typescript": "^5",
    "vite": "^5",
    "vitest": "^1",
    "vue-tsc": "^2"
  }
}
```

- [ ] **Step 2: Create `vite.config.ts`**

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  clearScreen: false,
  server: {
    port: 1420,
    strictPort: true,
    watch: {
      ignored: ['**/src-tauri/**'],
    },
  },
  envPrefix: ['VITE_', 'TAURI_ENV_*'],
  build: {
    target: process.env.TAURI_ENV_PLATFORM === 'windows' ? 'chrome105' : 'safari13',
    minify: !process.env.TAURI_ENV_DEBUG ? 'esbuild' : false,
    sourcemap: !!process.env.TAURI_ENV_DEBUG,
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
  },
})
```

- [ ] **Step 3: Create `tsconfig.json`**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "moduleResolution": "bundler",
    "strict": true,
    "jsx": "preserve",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

- [ ] **Step 4: Create `tsconfig.node.json`**

```json
{
  "compilerOptions": {
    "composite": true,
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true,
    "module": "ESNext",
    "target": "ES2022"
  },
  "include": ["vite.config.ts"]
}
```

- [ ] **Step 5: Create `index.html`**

```html
<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>FG-ABYSS</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

- [ ] **Step 6: Commit**

```bash
git init
git add package.json vite.config.ts tsconfig.json tsconfig.node.json index.html
git commit -m "chore: add project manifest files"
```

---

### Task 2: Rust/Tauri backend skeleton

**Files:**
- Create: `src-tauri/Cargo.toml`
- Create: `src-tauri/build.rs`
- Create: `src-tauri/src/main.rs`
- Create: `src-tauri/src/lib.rs`
- Create: `src-tauri/tauri.conf.json`
- Create: `src-tauri/capabilities/default.json`

- [ ] **Step 1: Create `src-tauri/Cargo.toml`**

```toml
[package]
name = "fg-abyss"
version = "0.1.0"
edition = "2021"

[lib]
name = "fg_abyss_lib"
crate-type = ["staticlib", "cdylib", "rlib"]

[build-dependencies]
tauri-build = { version = "2", features = [] }

[dependencies]
tauri = { version = "2", features = [] }
serde = { version = "1", features = ["derive"] }
serde_json = "1"

[profile.release]
codegen-units = 1
lto = true
opt-level = "s"
panic = "abort"
strip = true
```

- [ ] **Step 2: Create `src-tauri/build.rs`**

```rust
fn main() {
    tauri_build::build()
}
```

- [ ] **Step 3: Create `src-tauri/src/main.rs`**

```rust
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
    fg_abyss_lib::run()
}
```

- [ ] **Step 4: Create `src-tauri/src/lib.rs`**

```rust
pub fn run() {
    tauri::Builder::default()
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

- [ ] **Step 5: Create `src-tauri/tauri.conf.json`**

```json
{
  "$schema": "https://schema.tauri.app/config/2",
  "productName": "fg-abyss",
  "version": "0.1.0",
  "identifier": "com.fg-abyss.app",
  "build": {
    "beforeDevCommand": "pnpm dev",
    "devUrl": "http://localhost:1420",
    "beforeBuildCommand": "pnpm build",
    "frontendDist": "../dist"
  },
  "app": {
    "windows": [
      {
        "title": "FG-ABYSS",
        "width": 1280,
        "height": 800,
        "minWidth": 1024,
        "minHeight": 768,
        "decorations": false,
        "resizable": true,
        "center": true
      }
    ],
    "security": {
      "csp": null
    }
  },
  "bundle": {
    "active": true,
    "targets": "all",
    "icon": []
  }
}
```

- [ ] **Step 6: Create `src-tauri/capabilities/default.json`**

```json
{
  "$schema": "../gen/schemas/desktop-schema.json",
  "identifier": "default",
  "description": "Default capability",
  "windows": ["main"],
  "permissions": [
    "core:default",
    "core:window:allow-minimize",
    "core:window:allow-toggle-maximize",
    "core:window:allow-close"
  ]
}
```

- [ ] **Step 7: Commit**

```bash
git add src-tauri/
git commit -m "chore: add Tauri v2 Rust skeleton"
```

---

### Task 3: Install dependencies + smoke test

**Files:** none (installs to node_modules)

- [ ] **Step 1: Install frontend dependencies**

Run: `pnpm install`

Expected: `node_modules/` created, no errors

- [ ] **Step 2: Create minimal `src/main.ts` to verify Vite starts**

```typescript
import { createApp } from 'vue'

const app = createApp({ template: '<div>loading...</div>' })
app.mount('#app')
```

- [ ] **Step 3: Verify Vite dev server starts**

Run: `pnpm dev`

Expected: `http://localhost:1420` serves the page with "loading..." text. Stop with Ctrl+C.

- [ ] **Step 4: Verify TypeScript compiles**

Run: `pnpm typecheck`

Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add src/main.ts
git commit -m "chore: verify Vite dev server and TypeScript setup"
```

---

### Task 4: CSS foundation

**Files:**
- Create: `src/styles/variables.css`
- Create: `src/styles/app-shell.css`
- Create: `src/styles/scrollbar.css`
- Create: `src/styles/animations.css`

- [ ] **Step 1: Create `src/styles/variables.css`**

```css
/* Theme-invariant semantic variables. Theme-variant vars are set by useThemeCssVars() at runtime. */
:root {
  --color-success:  #4ade80;
  --color-warning:  #fbbf24;
  --color-error:    #f87171;
  --color-info:     #60a5fa;
  --code-bg:        #0c0d11;
  --code-border:    #1a1c26;
  --wc-close-bg:    #ef4444;
  --font-mono:      'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
}

*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html, body, #app {
  height: 100%;
  overflow: hidden;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  font-size: 13px;
  line-height: 1.4;
  background: var(--bg-base);
  color: var(--text-1);
  -webkit-font-smoothing: antialiased;
}
```

- [ ] **Step 2: Create `src/styles/app-shell.css`**

```css
.app-shell {
  display: grid;
  grid-template-rows: 48px 1fr 28px;
  height: 100vh;
  overflow: hidden;
  background: var(--bg-base);
}

.titlebar-area  { grid-row: 1; }
.statusbar-area { grid-row: 3; }

/* Middle row uses Flexbox so sidebar width transition works.
   CSS Grid cannot animate grid-template-columns on a CSS variable. */
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
  position: relative; /* for batch action bar absolute positioning */
}
```

- [ ] **Step 3: Create `src/styles/scrollbar.css`**

```css
::-webkit-scrollbar {
  width: 4px;
  height: 4px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 2px;
}
::-webkit-scrollbar-thumb:hover {
  background: var(--text-3);
}
```

- [ ] **Step 4: Create `src/styles/animations.css`**

```css
/* Disable all animations for users who prefer reduced motion */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}

/* Page route fade */
.fade-enter-active { transition: opacity 120ms ease-out; }
.fade-leave-active { transition: opacity 80ms ease-in; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

/* Status dot pulse — reduced to static glow when prefers-reduced-motion */
@keyframes pulse-dot {
  0%, 100% { box-shadow: 0 0 0 0 rgba(74, 222, 128, 0.4); }
  50%       { box-shadow: 0 0 0 5px rgba(74, 222, 128, 0); }
}
.status-dot-active {
  animation: pulse-dot 2.4s ease-in-out infinite;
}

/* Skeleton shimmer for loading states */
@keyframes skeleton-shimmer {
  0%   { background-position: -400px 0; }
  100% { background-position:  400px 0; }
}
.skeleton-row {
  background: linear-gradient(
    90deg,
    var(--bg-elevated) 25%,
    var(--bg-hover)    50%,
    var(--bg-elevated) 75%
  );
  background-size: 800px 40px;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
  border-radius: 4px;
}

/* Batch action bar slide */
.batch-bar-enter-active { transition: transform 200ms ease-out; }
.batch-bar-leave-active { transition: transform 140ms ease-in; }
.batch-bar-enter-from, .batch-bar-leave-to { transform: translateY(48px); }
```

- [ ] **Step 5: Run typecheck**

Run: `pnpm typecheck`

Expected: PASS (CSS files don't affect TS check)

- [ ] **Step 6: Commit**

```bash
git add src/styles/
git commit -m "feat: add CSS foundation (variables, app-shell, scrollbar, animations)"
```

---

### Task 5: Theme tokens + buildOverrides

**Files:**
- Create: `src/theme/tokens.ts`
- Create: `src/theme/overrides.ts`
- Create: `src/theme/index.ts`
- Create: `src/theme/overrides.test.ts`

- [ ] **Step 1: Write the failing test**

```typescript
// src/theme/overrides.test.ts
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pnpm test`

Expected: FAIL with "Cannot find module './overrides'"

- [ ] **Step 3: Create `src/theme/tokens.ts`**

```typescript
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
```

- [ ] **Step 4: Create `src/theme/overrides.ts`**

```typescript
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
  // accent-bg: 12% opacity dark / 9% opacity light (hex alpha suffix)
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
      colorSegment:        bg2,
      tabColorSegment:     bg3,
      tabTextColorActiveBar: accent,
      barColor:            accent,
    },
    Slider: {
      fillColor: accent,
      dotColor:  accent,
    },
  }
}
```

- [ ] **Step 5: Create `src/theme/index.ts`**

```typescript
export { buildOverrides } from './overrides'
export { ACCENT_COLORS, DARK_VARS, LIGHT_VARS } from './tokens'
export type { AccentKey } from './tokens'
```

- [ ] **Step 6: Run test to verify it passes**

Run: `pnpm test`

Expected: 4 tests PASS

- [ ] **Step 7: Commit**

```bash
git add src/theme/
git commit -m "feat: add theme tokens and buildOverrides with tests"
```

---

### Task 6: ThemeStore + useThemeCssVars

**Files:**
- Create: `src/stores/theme.ts`
- Create: `src/composables/useThemeCssVars.ts`
- Create: `src/stores/theme.test.ts`

- [ ] **Step 1: Write the failing test**

```typescript
// src/stores/theme.test.ts
import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore } from './theme'

describe('useThemeStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('default mode is dark', () => {
    const store = useThemeStore()
    expect(store.mode).toBe('dark')
    expect(store.resolvedMode).toBe('dark')
  })

  it('setMode changes mode', () => {
    const store = useThemeStore()
    store.setMode('light')
    expect(store.mode).toBe('light')
    expect(store.resolvedMode).toBe('light')
  })

  it('setAccent changes accentKey', () => {
    const store = useThemeStore()
    store.setAccent('cyan')
    expect(store.accentKey).toBe('cyan')
  })

  it('naiveTheme is null for light mode', () => {
    const store = useThemeStore()
    store.setMode('light')
    expect(store.naiveTheme).toBeNull()
  })
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pnpm test`

Expected: FAIL with "Cannot find module './theme'"

- [ ] **Step 3: Create `src/stores/theme.ts`**

```typescript
import { defineStore } from 'pinia'
import { computed } from 'vue'
import { darkTheme, useOsTheme } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'
import { buildOverrides, ACCENT_COLORS } from '@/theme'
import type { AccentKey } from '@/theme'

export type ThemeMode = 'dark' | 'light' | 'system'
export type FontSize = '12px' | '13px' | '14px' | '15px'

export { ACCENT_COLORS }
export type { AccentKey }

export const useThemeStore = defineStore('theme', () => {
  const osTheme = useOsTheme()
  const mode      = useLocalStorage<ThemeMode>('fg-theme-mode',   'dark')
  const accentKey = useLocalStorage<AccentKey>('fg-theme-accent', 'blue')
  const fontSize  = useLocalStorage<FontSize>( 'fg-theme-fontsize', '13px')

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

  function setMode(m: ThemeMode)    { mode.value = m }
  function setAccent(k: AccentKey)  { accentKey.value = k }
  function setFontSize(s: FontSize) { fontSize.value = s }

  return {
    mode, accentKey, fontSize,
    resolvedMode, naiveTheme, themeOverrides,
    setMode, setAccent, setFontSize,
  }
})
```

- [ ] **Step 4: Create `src/composables/useThemeCssVars.ts`**

```typescript
import { watchEffect } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { ACCENT_COLORS, DARK_VARS, LIGHT_VARS } from '@/theme'

export function useThemeCssVars() {
  const store = useThemeStore()

  watchEffect(() => {
    const isDark = store.resolvedMode === 'dark'
    const accent = ACCENT_COLORS[store.accentKey][store.resolvedMode]
    const vars   = isDark ? DARK_VARS : LIGHT_VARS
    const root   = document.documentElement

    root.setAttribute('data-theme', store.resolvedMode)
    Object.entries(vars).forEach(([k, v]) => root.style.setProperty(k, v))
    root.style.setProperty('--accent', accent)
    // accent-bg: 12% opacity dark (0x1f ≈ 12%), 9% opacity light (0x17 ≈ 9%)
    root.style.setProperty('--accent-bg', isDark ? `${accent}1f` : `${accent}17`)
    root.style.setProperty('--font-mono', "'JetBrains Mono', 'Fira Code', monospace")
  })
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `pnpm test`

Expected: 4 tests PASS

- [ ] **Step 6: Typecheck**

Run: `pnpm typecheck`

Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add src/stores/ src/composables/
git commit -m "feat: add ThemeStore and useThemeCssVars composable"
```

---

### Task 7: Shared types + i18n + router

**Files:**
- Create: `src/shared/types/error.ts`
- Create: `src/shared/types/loading.ts`
- Create: `src/shared/utils/invoke.ts`
- Create: `src/i18n/index.ts`
- Create: `src/router/index.ts`

- [ ] **Step 1: Create `src/shared/types/error.ts`**

```typescript
export interface AppError {
  kind:
    | 'Database' | 'Http' | 'Crypto' | 'Connection' | 'InvalidResponse'
    | 'CircuitOpen' | 'NotFound' | 'InvalidInput' | 'Io' | 'Serialize'
    | 'Locked' | 'Plugin' | 'NeedsRedeploy' | 'MemShellExpired'
  message: string
}
```

- [ ] **Step 2: Create `src/shared/types/loading.ts`**

```typescript
// Use LoadingMap instead of a single boolean. Allows loading.list, loading.create,
// loading['delete-uuid'] to be independent per CLAUDE.md.
export type LoadingMap = Record<string, boolean>

export function useLoadingMap(): LoadingMap {
  return {}
}
```

- [ ] **Step 3: Create `src/shared/utils/invoke.ts`**

```typescript
import { invoke as tauriInvoke } from '@tauri-apps/api/core'
import type { AppError } from '@/shared/types/error'

export async function invoke<T>(cmd: string, args?: Record<string, unknown>): Promise<T> {
  try {
    return await tauriInvoke<T>(cmd, args)
  } catch (e) {
    throw e as AppError
  }
}
```

- [ ] **Step 4: Create `src/i18n/index.ts`**

```typescript
import { createI18n } from 'vue-i18n'

const messages = {
  'zh-CN': {
    nav: {
      home:    '首页',
      project: '项目',
      payload: '载荷',
      plugin:  '插件',
      settings: '设置',
    },
    app: {
      name: '渊渟',
      sub:  'ABYSS',
    },
    common: {
      loading:  '加载中...',
      empty:    '暂无数据',
      confirm:  '确认',
      cancel:   '取消',
      save:     '保存',
      delete:   '删除',
      edit:     '编辑',
      create:   '新建',
      search:   '搜索',
      refresh:  '刷新',
      copy:     '复制',
      download: '下载',
    },
  },
  'en-US': {
    nav: {
      home:    'Home',
      project: 'Projects',
      payload: 'Payload',
      plugin:  'Plugins',
      settings: 'Settings',
    },
    app: {
      name: 'FG',
      sub:  'ABYSS',
    },
    common: {
      loading:  'Loading...',
      empty:    'No data',
      confirm:  'Confirm',
      cancel:   'Cancel',
      save:     'Save',
      delete:   'Delete',
      edit:     'Edit',
      create:   'New',
      search:   'Search',
      refresh:  'Refresh',
      copy:     'Copy',
      download: 'Download',
    },
  },
}

export const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'en-US',
  messages,
})
```

- [ ] **Step 5: Create `src/router/index.ts`**

```typescript
import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/project',
      component: () => import('@/views/ProjectView.vue'),
    },
    {
      path: '/payload',
      component: () => import('@/views/PayloadView.vue'),
    },
    {
      path: '/plugin',
      component: () => import('@/views/PluginView.vue'),
    },
    {
      path: '/settings',
      component: () => import('@/views/settings/SettingsView.vue'),
    },
    {
      path: '/console',
      component: () => import('@/views/ConsoleView.vue'),
    },
  ],
})

export default router
```

- [ ] **Step 6: Typecheck**

Run: `pnpm typecheck`

Expected: PASS (view files don't exist yet, but lazy imports don't fail typecheck)

- [ ] **Step 7: Commit**

```bash
git add src/shared/ src/i18n/ src/router/
git commit -m "feat: add shared types, i18n, and router"
```

---

### Task 8: App.vue + main.ts

**Files:**
- Create: `src/App.vue`
- Modify: `src/main.ts`

- [ ] **Step 1: Create `src/App.vue`**

```vue
<template>
  <n-config-provider
    :theme="themeStore.naiveTheme"
    :theme-overrides="themeStore.themeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme-overrides-common="{ zIndexModal: 100, zIndexPopover: 20, zIndexMessage: 1000 }"
  >
    <n-message-provider placement="top-right">
      <n-dialog-provider>
        <AppShell />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { zhCN, dateZhCN } from 'naive-ui'
import { useThemeStore } from '@/stores/theme'
import { useThemeCssVars } from '@/composables/useThemeCssVars'
import AppShell from '@/components/layout/AppShell.vue'

const themeStore = useThemeStore()
useThemeCssVars()
</script>
```

- [ ] **Step 2: Rewrite `src/main.ts`**

```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { i18n } from '@/i18n'
import router from '@/router'
import App from './App.vue'

// Global styles — order matters: variables first, then layout, then effects
import './styles/variables.css'
import './styles/app-shell.css'
import './styles/scrollbar.css'
import './styles/animations.css'

const app = createApp(App)
app.use(createPinia())
app.use(i18n)
app.use(router)
app.mount('#app')
```

- [ ] **Step 3: Typecheck**

Run: `pnpm typecheck`

Expected: errors about missing `AppShell.vue` — that's expected at this stage

- [ ] **Step 4: Commit**

```bash
git add src/App.vue src/main.ts
git commit -m "feat: wire App.vue with n-config-provider and theme system"
```

---

### Task 9: Layout components — CustomTitlebar

**Files:**
- Create: `src/components/layout/CustomTitlebar.vue`

- [ ] **Step 1: Create `src/components/layout/CustomTitlebar.vue`**

```vue
<template>
  <div class="titlebar" data-tauri-drag-region>
    <div class="titlebar-left" data-tauri-drag-region>
      <div class="logo" aria-hidden="true">◈</div>
      <span class="app-name">{{ t('app.name') }}</span>
      <span class="app-sub">{{ t('app.sub') }}</span>
    </div>

    <div class="titlebar-right" data-tauri-drag-region="false">
      <!-- Theme toggle -->
      <button
        class="tb-icon-btn"
        :aria-label="`切换主题: ${modeLabel}`"
        :title="`切换主题: ${modeLabel}`"
        @click="cycleMode"
      >
        <Moon v-if="themeStore.mode === 'dark'" :size="16" />
        <Sun v-else-if="themeStore.mode === 'light'" :size="16" />
        <Monitor v-else :size="16" />
      </button>

      <!-- Language toggle -->
      <button
        class="tb-icon-btn"
        aria-label="切换语言"
        title="切换语言"
        @click="toggleLang"
      >
        <Languages :size="16" />
      </button>

      <!-- Window controls -->
      <div class="win-controls">
        <button class="wc-btn wc-min" aria-label="最小化" @click="minimize">
          <Minus :size="12" />
        </button>
        <button class="wc-btn wc-max" aria-label="最大化" @click="toggleMaximize">
          <Square :size="12" />
        </button>
        <button class="wc-btn wc-close" aria-label="关闭" @click="close">
          <X :size="12" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Moon, Sun, Monitor, Languages, Minus, Square, X } from 'lucide-vue-next'
import { getCurrentWindow } from '@tauri-apps/api/window'
import { useThemeStore } from '@/stores/theme'
import type { ThemeMode } from '@/stores/theme'

const { t, locale } = useI18n()
const themeStore = useThemeStore()
const win = getCurrentWindow()

const modeLabel = computed(() => ({
  dark: '深色', light: '浅色', system: '跟随系统',
}[themeStore.mode]))

const modeOrder: ThemeMode[] = ['dark', 'light', 'system']
function cycleMode() {
  const idx = modeOrder.indexOf(themeStore.mode)
  themeStore.setMode(modeOrder[(idx + 1) % 3])
}

function toggleLang() {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
}

const minimize       = () => win.minimize()
const toggleMaximize = () => win.toggleMaximize()
const close          = () => win.close()
</script>

<style scoped>
.titlebar {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-deep);
  border-bottom: 1px solid var(--border);
  user-select: none;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-left: 16px;
  height: 100%;
}

.logo {
  font-size: 20px;
  color: var(--accent);
  line-height: 1;
}

.app-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-1);
}

.app-sub {
  font-size: 11px;
  color: var(--text-3);
  margin-left: 2px;
  letter-spacing: 0.08em;
}

.titlebar-right {
  display: flex;
  align-items: center;
  height: 100%;
}

.tb-icon-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-2);
  cursor: pointer;
  border-radius: 4px;
  margin: 0 2px;
  transition: background 80ms ease-out, color 80ms ease-out;
}

.tb-icon-btn:hover {
  background: var(--bg-hover);
  color: var(--text-1);
}

.win-controls {
  display: flex;
  align-items: center;
  height: 100%;
  margin-left: 4px;
}

.wc-btn {
  width: 40px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-2);
  cursor: default;
  transition: background 80ms ease-out, color 80ms ease-out;
}

.wc-btn:hover {
  background: var(--bg-hover);
  color: var(--text-1);
}

.wc-close {
  width: 48px;
}

.wc-close:hover {
  background: var(--wc-close-bg);
  color: #fff;
}
</style>
```

- [ ] **Step 2: Typecheck**

Run: `pnpm typecheck`

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add src/components/layout/CustomTitlebar.vue
git commit -m "feat: add CustomTitlebar with theme toggle and window controls"
```

---

### Task 10: Layout components — Sidebar

**Files:**
- Create: `src/components/layout/Sidebar.vue`

- [ ] **Step 1: Create `src/components/layout/Sidebar.vue`**

```vue
<template>
  <nav
    class="sidebar"
    :class="{ 'is-collapsed': collapsed }"
    role="navigation"
    aria-label="主导航"
  >
    <div class="nav-items">
      <n-tooltip
        v-for="item in navItems"
        :key="item.path"
        placement="right"
        :delay="400"
        :disabled="!collapsed"
      >
        <template #trigger>
          <RouterLink
            :to="item.path"
            class="nav-item"
            :class="{ 'is-active': route.path === item.path }"
            :aria-label="collapsed ? item.label : undefined"
          >
            <component :is="item.icon" :size="18" class="nav-icon" />
            <span class="nav-label">{{ item.label }}</span>
          </RouterLink>
        </template>
        {{ item.label }}
      </n-tooltip>
    </div>

    <button
      class="collapse-btn"
      :aria-label="collapsed ? '展开侧边栏' : '折叠侧边栏'"
      @click="collapsed = !collapsed"
    >
      <ChevronLeft
        :size="16"
        class="collapse-icon"
        :class="{ 'is-rotated': collapsed }"
      />
      <span class="nav-label">折叠</span>
    </button>
  </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useLocalStorage } from '@vueuse/core'
import { Home, FolderOpen, Package, Puzzle, Settings } from 'lucide-vue-next'
import { ChevronLeft } from 'lucide-vue-next'

const route = useRoute()
const collapsed = useLocalStorage('fg-sidebar-collapsed', false)

const navItems = [
  { path: '/',         label: '首页', icon: Home       },
  { path: '/project',  label: '项目', icon: FolderOpen  },
  { path: '/payload',  label: '载荷', icon: Package     },
  { path: '/plugin',   label: '插件', icon: Puzzle      },
  { path: '/settings', label: '设置', icon: Settings    },
]
</script>

<style scoped>
.sidebar {
  height: 100%;
  background: var(--bg-deep);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.nav-items {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 12px;
  text-decoration: none;
  color: var(--text-2);
  border-left: 2px solid transparent; /* permanent — prevents 2px jitter on activation */
  cursor: pointer;
  transition: background 80ms ease-out, color 80ms ease-out;
  white-space: nowrap;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-1);
}

.nav-item.is-active {
  border-left-color: var(--accent);
  background: var(--accent-bg);
  color: var(--accent);
  font-weight: 600;
}

.nav-item.is-active .nav-icon {
  color: var(--accent);
}

.nav-icon {
  flex-shrink: 0;
}

.nav-label {
  font-size: 13px;
  overflow: hidden;
  white-space: nowrap;
  transition: opacity 200ms ease-out, width 200ms ease-out;
  width: auto;
  opacity: 1;
}

/* Collapsed: hide text */
.sidebar.is-collapsed .nav-label {
  opacity: 0;
  width: 0;
}

.sidebar.is-collapsed .nav-item {
  padding: 0;
  justify-content: center;
}

.collapse-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 12px;
  border: none;
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  border-top: 1px solid var(--border);
  width: 100%;
  transition: background 80ms ease-out, color 80ms ease-out;
  white-space: nowrap;
}

.collapse-btn:hover {
  background: var(--bg-hover);
  color: var(--text-2);
}

.collapse-icon {
  flex-shrink: 0;
  transition: transform 220ms cubic-bezier(0.4, 0, 0.2, 1);
}

.collapse-icon.is-rotated {
  transform: rotate(180deg);
}

.sidebar.is-collapsed .collapse-btn {
  padding: 0;
  justify-content: center;
}
</style>
```

- [ ] **Step 2: Typecheck**

Run: `pnpm typecheck`

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add src/components/layout/Sidebar.vue
git commit -m "feat: add collapsible Sidebar with keyboard navigation"
```

---

### Task 11: Layout components — StatusBar

**Files:**
- Create: `src/components/layout/StatusBar.vue`

- [ ] **Step 1: Create `src/components/layout/StatusBar.vue`**

```vue
<template>
  <div class="statusbar" role="status" aria-live="polite">
    <div class="status-left">
      <span
        class="status-dot"
        :class="{ 'status-dot-active': activeCount > 0 }"
        :aria-label="`状态: ${activeCount > 0 ? '活跃' : '离线'}`"
      />
      <span class="status-text">{{ activeCount > 0 ? '已连接' : '未连接' }}</span>
    </div>

    <div class="divider" aria-hidden="true" />
    <span class="stat-item">
      项目 <span class="stat-num">{{ projectCount }}</span>
    </span>

    <div class="divider" aria-hidden="true" />
    <span class="stat-item">
      Shell <span class="stat-num">{{ shellCount }}</span>
    </span>

    <div class="divider" aria-hidden="true" />
    <span class="stat-item">
      载荷 <span class="stat-num">{{ payloadCount }}</span>
    </span>
  </div>
</template>

<script setup lang="ts">
// Mock data — will be replaced by store values in feature plans
const activeCount  = 0
const projectCount = 0
const shellCount   = 0
const payloadCount = 0
</script>

<style scoped>
.statusbar {
  height: 28px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 16px;
  background: var(--bg-deep);
  border-top: 1px solid var(--border);
  font-size: 12px;
  color: var(--text-2);
}

.status-left {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-3);
  flex-shrink: 0;
}

.status-dot.status-dot-active {
  background: var(--color-success);
}

.status-text {
  font-size: 12px;
}

.divider {
  width: 1px;
  height: 12px;
  background: var(--border);
}

.stat-item {
  font-size: 12px;
  color: var(--text-2);
}

.stat-num {
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  color: var(--text-1);
  margin-left: 4px;
}
</style>
```

- [ ] **Step 2: Commit**

```bash
git add src/components/layout/StatusBar.vue
git commit -m "feat: add StatusBar with status dot and counts"
```

---

### Task 12: AppShell — first visual milestone

**Files:**
- Create: `src/components/layout/AppShell.vue`

- [ ] **Step 1: Create `src/components/layout/AppShell.vue`**

```vue
<template>
  <div class="app-shell">
    <div class="titlebar-area">
      <CustomTitlebar />
    </div>

    <div class="app-main">
      <div class="sidebar-panel" :class="{ 'is-collapsed': collapsed }">
        <Sidebar />
      </div>
      <div class="content-area">
        <RouterView v-slot="{ Component }">
          <Transition name="fade" mode="out-in">
            <component :is="Component" :key="$route.fullPath" />
          </Transition>
        </RouterView>
      </div>
    </div>

    <div class="statusbar-area">
      <StatusBar />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import CustomTitlebar from './CustomTitlebar.vue'
import Sidebar from './Sidebar.vue'
import StatusBar from './StatusBar.vue'

const collapsed = useLocalStorage('fg-sidebar-collapsed', false)
</script>
```

- [ ] **Step 2: Create stub view files so router can resolve them**

Create these 6 empty stub files:

`src/views/HomeView.vue`:
```vue
<template><div class="view-stub" style="padding:24px;color:var(--text-1)">Home</div></template>
```

`src/views/ProjectView.vue`:
```vue
<template><div class="view-stub" style="padding:24px;color:var(--text-1)">Projects</div></template>
```

`src/views/PayloadView.vue`:
```vue
<template><div class="view-stub" style="padding:24px;color:var(--text-1)">Payload</div></template>
```

`src/views/PluginView.vue`:
```vue
<template><div class="view-stub" style="padding:24px;color:var(--text-1)">Plugins</div></template>
```

`src/views/ConsoleView.vue`:
```vue
<template><div class="view-stub" style="padding:24px;color:var(--text-1)">Console</div></template>
```

`src/views/settings/SettingsView.vue`:
```vue
<template><div class="view-stub" style="padding:24px;color:var(--text-1)">Settings</div></template>
```

- [ ] **Step 3: Typecheck**

Run: `pnpm typecheck`

Expected: PASS

- [ ] **Step 4: Visual verification — run the app**

Run: `pnpm tauri dev`

Verify:
- App opens without chrome decorations (frameless)
- Dark theme applied — near-black background visible
- Titlebar shows "◈ 渊渟 ABYSS", moon icon, languages icon, window controls
- Sidebar shows 5 nav items: 首页/项目/载荷/插件/设置
- Clicking sidebar items navigates routes (stub content changes)
- Collapse button at sidebar bottom collapses sidebar to icon-only
- Theme toggle cycles dark → light → system → dark
- Light theme shows correct light colors
- StatusBar shows at bottom with dot + counts

Stop with Ctrl+C.

- [ ] **Step 5: Commit**

```bash
git add src/components/layout/AppShell.vue src/views/
git commit -m "feat: add AppShell with layout and stub views — first visual milestone"
```

---

### Task 13: HomeView

**Files:**
- Modify: `src/views/HomeView.vue`

- [ ] **Step 1: Replace HomeView stub with full implementation**

```vue
<template>
  <div class="home-view">
    <div class="home-header">
      <h1 class="home-title">欢迎回来</h1>
      <span class="home-date">{{ todayStr }}</span>
    </div>

    <!-- Stat cards -->
    <div class="stat-cards">
      <div class="stat-card" v-for="card in statCards" :key="card.label">
        <component :is="card.icon" :size="20" class="stat-icon" />
        <div class="stat-body">
          <div class="stat-num">{{ card.value }}</div>
          <div class="stat-label">{{ card.label }}</div>
        </div>
      </div>
    </div>

    <!-- Recent connections -->
    <div class="recent-section">
      <h2 class="section-title">最近连接</h2>
      <div v-if="recentShells.length === 0" class="empty-state">
        <Globe :size="64" class="empty-icon" />
        <div class="empty-title">还没有 WebShell</div>
        <div class="empty-sub">前往「项目」页面添加第一个 WebShell</div>
        <RouterLink to="/project">
          <n-button type="primary" size="medium">前往项目</n-button>
        </RouterLink>
      </div>
      <div v-else class="recent-list">
        <div
          v-for="shell in recentShells"
          :key="shell.id"
          class="recent-item"
        >
          <span
            class="recent-dot"
            :class="shell.active ? 'dot-active' : 'dot-offline'"
            :aria-label="`状态: ${shell.active ? '活跃' : '离线'}`"
          />
          <span class="recent-name">{{ shell.name }}</span>
          <span class="recent-url">{{ shell.url }}</span>
          <n-tag size="small" :color="typeColor(shell.type)">{{ shell.type }}</n-tag>
          <span class="recent-time">{{ shell.relativeTime }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Terminal, FolderOpen, Package, Zap, Globe } from 'lucide-vue-next'

const todayStr = new Date().toLocaleDateString('zh-CN', {
  year: 'numeric', month: '2-digit', day: '2-digit',
}).replace(/\//g, '/')

const statCards = [
  { icon: Terminal,   value: 0, label: '活跃 Shell' },
  { icon: FolderOpen, value: 0, label: '项目'       },
  { icon: Package,    value: 0, label: '载荷'       },
  { icon: Zap,        value: 0, label: '今日请求'   },
]

const recentShells: Array<{
  id: string; name: string; url: string
  type: string; active: boolean; relativeTime: string
}> = []

const TYPE_COLORS: Record<string, { color: string; textColor: string; borderColor: string }> = {
  PHP:  { color: 'rgba(79,156,255,0.15)',  textColor: '#4f9cff', borderColor: 'transparent' },
  JSP:  { color: 'rgba(251,146,60,0.15)',  textColor: '#fb923c', borderColor: 'transparent' },
  ASP:  { color: 'rgba(167,139,250,0.15)', textColor: '#a78bfa', borderColor: 'transparent' },
  ASPX: { color: 'rgba(34,211,238,0.15)',  textColor: '#22d3ee', borderColor: 'transparent' },
}

function typeColor(type: string) {
  return TYPE_COLORS[type] ?? TYPE_COLORS.PHP
}
</script>

<style scoped>
.home-view {
  height: 100%;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.home-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.home-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-1);
}

.home-date {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}

.stat-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-left: 3px solid var(--accent);
}

[data-theme="light"] .stat-card {
  border: none;
  border-left: 3px solid var(--accent);
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}

.stat-icon {
  color: var(--accent);
  flex-shrink: 0;
}

.stat-num {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-1);
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
}

.stat-label {
  font-size: 12px;
  color: var(--text-2);
  margin-top: 2px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
  margin-bottom: 12px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 48px 0;
  color: var(--text-3);
}

.empty-icon { color: var(--text-3); }

.empty-title {
  font-size: 14px;
  color: var(--text-1);
}

.empty-sub {
  font-size: 13px;
  color: var(--text-2);
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 40px;
  padding: 0 12px;
  border-radius: 6px;
  background: var(--bg-elevated);
}

.recent-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.dot-active  { background: var(--color-success); }
.dot-offline { background: var(--text-3); }

.recent-name {
  font-size: 13px;
  color: var(--text-1);
  min-width: 120px;
}

.recent-url {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-2);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.recent-time {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}
</style>
```

- [ ] **Step 2: Typecheck**

Run: `pnpm typecheck`

Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add src/views/HomeView.vue
git commit -m "feat: implement HomeView with stat cards and recent connections"
```

---

### Task 14: ProjectView

**Files:**
- Modify: `src/views/ProjectView.vue`

- [ ] **Step 1: Replace ProjectView stub**

```vue
<template>
  <div class="project-view">
    <!-- Project tree -->
    <div class="project-tree">
      <div class="tree-header">
        <n-button size="small" type="primary" ghost>
          <template #icon><Plus :size="14" /></template>
          新建项目
        </n-button>
      </div>
      <div class="tree-list">
        <div
          v-for="proj in projects"
          :key="proj.id"
          class="tree-item"
          :class="{ 'is-active': selectedProject === proj.id }"
          @click="selectedProject = proj.id"
        >
          <span class="tree-name">{{ proj.name }}</span>
          <span class="tree-count">{{ proj.shellCount }}</span>
        </div>
        <div
          class="tree-item tree-item--recycle"
          :class="{ 'is-active': selectedProject === '__recycle__' }"
          @click="selectedProject = '__recycle__'"
        >
          <Trash2 :size="14" />
          <span class="tree-name">回收站</span>
          <span class="tree-count">0</span>
        </div>
      </div>
    </div>

    <!-- Shell table -->
    <div class="shell-table-area">
      <div class="table-toolbar">
        <n-input
          v-model:value="searchText"
          placeholder="搜索..."
          clearable
          size="small"
          style="width: 220px"
        >
          <template #prefix><Search :size="14" /></template>
        </n-input>
        <n-select
          v-model:value="typeFilter"
          :options="typeOptions"
          placeholder="类型"
          clearable
          size="small"
          style="width: 100px"
        />
        <div style="flex:1" />
        <n-button type="primary" size="small">
          <template #icon><Plus :size="14" /></template>
          添加 Shell
        </n-button>
      </div>

      <n-data-table
        :columns="columns"
        :data="filteredShells"
        :row-key="(row: Shell) => row.id"
        :loading="false"
        size="small"
        class="shell-table"
      />

      <div class="table-footer">
        <n-pagination :page="1" :page-count="1" size="small" />
        <span class="total-count">共 {{ filteredShells.length }} 条</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import type { DataTableColumns } from 'naive-ui'
import { Plus, Search, Trash2 } from 'lucide-vue-next'

interface Project { id: string; name: string; shellCount: number }
interface Shell {
  id: string; name: string; url: string
  type: 'PHP' | 'JSP' | 'ASP' | 'ASPX'
  status: 'active' | 'offline' | 'error'
  lastConnected: string
  projectId: string
}

const projects = ref<Project[]>([
  { id: 'all', name: '全部', shellCount: 0 },
])

const selectedProject = ref('all')
const searchText = ref('')
const typeFilter = ref<string | null>(null)

const typeOptions = [
  { label: 'PHP',  value: 'PHP'  },
  { label: 'JSP',  value: 'JSP'  },
  { label: 'ASP',  value: 'ASP'  },
  { label: 'ASPX', value: 'ASPX' },
]

const shells = ref<Shell[]>([])

const filteredShells = computed(() => {
  let list = shells.value
  if (typeFilter.value) list = list.filter(s => s.type === typeFilter.value)
  if (searchText.value) {
    const q = searchText.value.toLowerCase()
    list = list.filter(s => s.name.toLowerCase().includes(q) || s.url.toLowerCase().includes(q))
  }
  return list
})

const STATUS_DOT: Record<string, string> = {
  active:  'var(--color-success)',
  offline: 'var(--text-3)',
  error:   'var(--color-error)',
}

const TYPE_COLORS: Record<string, { bg: string; text: string }> = {
  PHP:  { bg: 'rgba(79,156,255,0.15)',  text: '#4f9cff' },
  JSP:  { bg: 'rgba(251,146,60,0.15)',  text: '#fb923c' },
  ASP:  { bg: 'rgba(167,139,250,0.15)', text: '#a78bfa' },
  ASPX: { bg: 'rgba(34,211,238,0.15)',  text: '#22d3ee' },
}

const columns: DataTableColumns<Shell> = [
  {
    type: 'selection',
    width: 40,
  },
  {
    key: 'status',
    title: '状',
    width: 40,
    render: (row) => h('span', {
      style: {
        display: 'inline-block',
        width: '8px', height: '8px',
        borderRadius: '50%',
        background: STATUS_DOT[row.status],
      },
      'aria-label': `状态: ${row.status}`,
    }),
  },
  {
    key: 'name',
    title: '名称',
    width: 160,
    render: (row) => h('span', { style: { color: 'var(--text-1)' } }, row.name),
  },
  {
    key: 'url',
    title: 'URL',
    render: (row) => h('span', {
      style: { fontFamily: 'var(--font-mono)', fontSize: '12px', color: 'var(--text-2)' },
      title: row.url,
    }, row.url),
  },
  {
    key: 'type',
    title: '类型',
    width: 80,
    render: (row) => {
      const c = TYPE_COLORS[row.type] ?? TYPE_COLORS.PHP
      return h('span', {
        style: {
          padding: '2px 6px',
          borderRadius: '4px',
          background: c.bg,
          color: c.text,
          fontSize: '12px',
          fontFamily: 'var(--font-mono)',
        },
      }, row.type)
    },
  },
  {
    key: 'lastConnected',
    title: '最近连接',
    width: 120,
    render: (row) => h('span', { style: { color: 'var(--text-2)', fontSize: '12px' } }, row.lastConnected),
  },
]
</script>

<style scoped>
.project-view {
  height: 100%;
  display: flex;
  overflow: hidden;
}

.project-tree {
  width: 240px;
  flex-shrink: 0;
  background: var(--bg-deep);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.tree-header {
  padding: 12px;
  border-bottom: 1px solid var(--border);
}

.tree-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.tree-item {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  padding: 0 16px;
  cursor: pointer;
  color: var(--text-2);
  font-size: 13px;
  transition: background 80ms ease-out;
}

.tree-item:hover { background: var(--bg-hover); color: var(--text-1); }
.tree-item.is-active { background: var(--accent-bg); color: var(--accent); font-weight: 600; }

.tree-name { flex: 1; }
.tree-count { font-size: 12px; color: var(--text-3); font-family: var(--font-mono); }

.tree-item--recycle {
  border-top: 1px solid var(--border);
  color: var(--text-3);
  margin-top: 8px;
}

.shell-table-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.table-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
}

.shell-table {
  flex: 1;
  overflow: hidden;
}

.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-top: 1px solid var(--border);
}

.total-count {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}
</style>
```

- [ ] **Step 2: Typecheck + commit**

Run: `pnpm typecheck`

Expected: PASS

```bash
git add src/views/ProjectView.vue
git commit -m "feat: implement ProjectView with project tree and shell table"
```

---

### Task 15: PayloadView

**Files:**
- Modify: `src/views/PayloadView.vue`

- [ ] **Step 1: Replace PayloadView stub**

```vue
<template>
  <div class="payload-view">
    <!-- Left config panel -->
    <div class="config-panel">
      <div class="config-section">
        <label class="field-label">语言</label>
        <n-radio-group v-model:value="lang" name="lang" size="small">
          <n-radio-button value="PHP">PHP</n-radio-button>
          <n-radio-button value="JSP">JSP</n-radio-button>
          <n-radio-button value="ASP">ASP</n-radio-button>
          <n-radio-button value="ASPX">ASPX</n-radio-button>
        </n-radio-group>
      </div>

      <div class="config-section">
        <label class="field-label">加密方式</label>
        <n-radio-group v-model:value="encMode" name="enc">
          <div class="radio-list">
            <n-radio value="aes">AES-256-GCM <n-tag size="tiny" type="info">推荐</n-tag></n-radio>
            <n-radio value="xor">XOR <n-tag size="tiny">兼容</n-tag></n-radio>
            <n-radio value="none">无加密</n-radio>
          </div>
        </n-radio-group>
      </div>

      <div class="config-section">
        <label class="field-label">密钥</label>
        <div class="key-row">
          <n-input
            v-model:value="keyValue"
            :type="showKey ? 'text' : 'password'"
            placeholder="输入密钥"
            size="small"
            style="font-family: var(--font-mono); flex:1"
          />
          <n-button text size="small" :aria-label="showKey ? '隐藏密码' : '显示密码'" @click="showKey = !showKey">
            <EyeOff v-if="showKey" :size="14" />
            <Eye v-else :size="14" />
          </n-button>
          <n-button text size="small" aria-label="复制密钥" @click="copyKey">
            <Check v-if="keyCopied" :size="14" style="color: var(--color-success)" />
            <Copy v-else :size="14" />
          </n-button>
        </div>
        <!-- Password strength bar -->
        <div class="strength-bar">
          <div
            v-for="i in 4"
            :key="i"
            class="strength-seg"
            :class="{ 'seg-active': keyStrength >= i }"
            :style="{ background: strengthColor(keyStrength, i) }"
          />
        </div>
      </div>

      <div class="config-section">
        <label class="field-label">混淆等级</label>
        <div class="slider-row">
          <n-slider v-model:value="obfLevel" :min="0" :max="5" :step="1" style="flex:1" />
          <span class="level-badge">L{{ obfLevel }}</span>
        </div>
        <div class="obf-checks">
          <n-checkbox v-model:checked="obfVarRename"  :disabled="obfLevel === 0">变量重命名</n-checkbox>
          <n-checkbox v-model:checked="obfStrEncrypt" :disabled="obfLevel === 0">字符串加密</n-checkbox>
          <n-checkbox v-model:checked="obfJunk"       :disabled="obfLevel === 0">垃圾代码</n-checkbox>
        </div>
      </div>

      <n-button
        type="primary"
        block
        size="medium"
        style="margin-top: auto; height: 40px"
        :loading="generating"
        :style="generateBtnStyle"
        @click="generate"
      >
        生成载荷
      </n-button>
    </div>

    <!-- Right code panel -->
    <div class="code-panel">
      <div class="code-toolbar">
        <n-button text size="small" aria-label="复制代码" @click="copyCode">
          <Check v-if="codeCopied" :size="14" style="color: var(--color-success)" />
          <Copy v-else :size="14" />
        </n-button>
        <n-button text size="small" aria-label="下载">
          <Download :size="14" />
        </n-button>
      </div>
      <div class="code-area">
        <pre class="code-placeholder">{{ codeOutput || '// 点击「生成载荷」生成代码' }}</pre>
      </div>
      <div class="code-status">
        <span>{{ lang }} · {{ encLabel[encMode] }} · {{ codeOutput ? (codeOutput.length + ' B') : '--' }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Eye, EyeOff, Copy, Check, Download } from 'lucide-vue-next'

const lang     = ref<'PHP' | 'JSP' | 'ASP' | 'ASPX'>('PHP')
const encMode  = ref<'aes' | 'xor' | 'none'>('aes')
const keyValue = ref('')
const showKey  = ref(false)
const keyCopied = ref(false)
const obfLevel = ref(0)
const obfVarRename  = ref(true)
const obfStrEncrypt = ref(true)
const obfJunk       = ref(false)
const generating = ref(false)
const codeOutput = ref('')
const codeCopied = ref(false)
const generateError = ref(false)

const encLabel: Record<string, string> = {
  aes: 'AES-256', xor: 'XOR', none: '无加密',
}

const keyStrength = computed(() => {
  const k = keyValue.value
  if (k.length === 0) return 0
  if (k.length < 8)   return 1
  if (k.length < 16)  return 2
  const hasNum = /\d/.test(k)
  const hasSym = /[^a-zA-Z0-9]/.test(k)
  return hasNum && hasSym ? 4 : 3
})

function strengthColor(strength: number, seg: number): string {
  if (strength < seg) return 'var(--bg-hover)'
  if (strength === 1) return '#f87171'
  if (strength === 2) return '#fb923c'
  if (strength === 3) return '#fbbf24'
  return '#4ade80'
}

const generateBtnStyle = computed(() => ({
  ...(generateError.value ? { border: '1px solid #f87171' } : {}),
}))

async function copyKey() {
  if (!keyValue.value) return
  await navigator.clipboard.writeText(keyValue.value)
  keyCopied.value = true
  setTimeout(() => (keyCopied.value = false), 1500)
}

async function copyCode() {
  if (!codeOutput.value) return
  await navigator.clipboard.writeText(codeOutput.value)
  codeCopied.value = true
  setTimeout(() => (codeCopied.value = false), 1500)
}

async function generate() {
  generating.value = true
  generateError.value = false
  // Placeholder — actual Tauri invoke wired in payload feature plan
  await new Promise(r => setTimeout(r, 600))
  codeOutput.value = `<?php\n// Generated ${lang.value} payload\n// Key: ${keyValue.value || '(empty)'}\n// Enc: ${encMode.value}\n?>`
  generating.value = false
}
</script>

<style scoped>
.payload-view {
  height: 100%;
  display: flex;
  overflow: hidden;
}

.config-panel {
  width: 340px;
  flex-shrink: 0;
  border-right: 1px solid var(--border);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow-y: auto;
}

.config-section { display: flex; flex-direction: column; gap: 8px; }

.field-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-2);
}

.radio-list { display: flex; flex-direction: column; gap: 8px; }

.key-row { display: flex; align-items: center; gap: 4px; }

.strength-bar {
  display: flex;
  gap: 4px;
  height: 3px;
  margin-top: 6px;
}

.strength-seg {
  flex: 1;
  border-radius: 2px;
  background: var(--bg-hover);
  transition: background 200ms;
}

.slider-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.level-badge {
  min-width: 28px;
  text-align: center;
  background: var(--accent-bg);
  color: var(--accent);
  font-family: var(--font-mono);
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
}

.obf-checks { display: flex; flex-direction: column; gap: 6px; padding-top: 8px; }

.code-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--code-bg);
  border: 1px solid var(--code-border);
  margin: 16px;
  border-radius: 6px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.12), 0 0 0 1px rgba(0,0,0,0.08);
}

.code-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-bottom: 1px solid var(--code-border);
  background: #16181f;
}

.code-area {
  flex: 1;
  overflow: auto;
  padding: 16px;
}

.code-placeholder {
  font-family: var(--font-mono);
  font-size: 13px;
  color: #abb2bf;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.code-status {
  padding: 6px 12px;
  font-size: 12px;
  color: var(--text-2);
  border-top: 1px solid var(--code-border);
  text-align: right;
  font-family: var(--font-mono);
  background: #16181f;
}
</style>
```

- [ ] **Step 2: Typecheck + commit**

Run: `pnpm typecheck`

Expected: PASS

```bash
git add src/views/PayloadView.vue
git commit -m "feat: implement PayloadView with config panel and code preview"
```

---

### Task 16: PluginView

**Files:**
- Modify: `src/views/PluginView.vue`

- [ ] **Step 1: Replace PluginView stub**

```vue
<template>
  <div class="plugin-view">
    <div class="plugin-toolbar">
      <span class="plugin-count">已安装 ({{ filteredPlugins.length }})</span>
      <div class="filter-tabs">
        <button
          v-for="tab in filterTabs"
          :key="tab.value"
          class="filter-tab"
          :class="{ 'is-active': activeFilter === tab.value }"
          @click="activeFilter = tab.value"
        >{{ tab.label }}</button>
      </div>
      <div style="flex:1" />
      <n-button type="primary" size="small">
        <template #icon><Plus :size="14" /></template>
        安装
      </n-button>
    </div>

    <div class="plugin-grid" v-if="filteredPlugins.length > 0">
      <div
        v-for="plugin in filteredPlugins"
        :key="plugin.id"
        class="plugin-card"
        :class="{ 'is-disabled': !plugin.enabled }"
      >
        <div class="card-header">
          <span class="plugin-name">{{ plugin.name }}</span>
          <div class="card-badges">
            <span class="version-badge">{{ plugin.version }}</span>
            <n-tag v-if="plugin.builtin" size="tiny" :color="{ color: 'var(--accent-bg)', textColor: 'var(--accent)', borderColor: 'transparent' }">内置</n-tag>
            <n-tag v-else size="tiny">第三方</n-tag>
          </div>
        </div>
        <p class="plugin-desc">{{ plugin.description }}</p>
        <div class="card-footer">
          <div style="display:flex;align-items:center;gap:8px">
            <n-switch
              :value="plugin.enabled"
              size="small"
              @update:value="(v: boolean) => togglePlugin(plugin.id, v)"
            />
            <span class="switch-label">{{ plugin.enabled ? '已启用' : '已禁用' }}</span>
          </div>
          <n-button v-if="!plugin.builtin" text size="tiny" type="error">卸载</n-button>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <Puzzle :size="64" class="empty-icon" />
      <div class="empty-title">还没有安装插件</div>
      <div class="empty-sub">内置插件随应用分发</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Plus, Puzzle } from 'lucide-vue-next'

interface Plugin {
  id: string; name: string; version: string
  builtin: boolean; description: string; enabled: boolean
}

const plugins = ref<Plugin[]>([
  { id: 'file-manager', name: '文件管理', version: 'v1.2.0', builtin: true,  description: '远程文件系统浏览、上传、下载、编辑', enabled: true  },
  { id: 'db-manager',   name: '数据库管理', version: 'v1.0.0', builtin: true,  description: 'SQL 查询执行与结果导出', enabled: true  },
  { id: 'terminal',     name: '终端',   version: 'v1.1.0', builtin: true,  description: '远程命令执行与进程管理', enabled: true  },
])

const activeFilter = ref('all')

const filterTabs = [
  { label: '全部',   value: 'all'      },
  { label: '已启用', value: 'enabled'  },
  { label: '已禁用', value: 'disabled' },
  { label: '内置',   value: 'builtin'  },
  { label: '第三方', value: 'third'    },
]

const filteredPlugins = computed(() => {
  switch (activeFilter.value) {
    case 'enabled':  return plugins.value.filter(p => p.enabled)
    case 'disabled': return plugins.value.filter(p => !p.enabled)
    case 'builtin':  return plugins.value.filter(p => p.builtin)
    case 'third':    return plugins.value.filter(p => !p.builtin)
    default:         return plugins.value
  }
})

function togglePlugin(id: string, enabled: boolean) {
  const p = plugins.value.find(x => x.id === id)
  if (p) p.enabled = enabled
}
</script>

<style scoped>
.plugin-view {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 16px;
  overflow-y: auto;
}

.plugin-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.plugin-count {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
}

.filter-tabs {
  display: flex;
  gap: 4px;
}

.filter-tab {
  padding: 4px 12px;
  border-radius: 4px;
  border: none;
  background: transparent;
  color: var(--text-2);
  font-size: 13px;
  cursor: pointer;
  transition: background 80ms;
}

.filter-tab:hover { background: var(--bg-hover); color: var(--text-1); }
.filter-tab.is-active { background: var(--accent-bg); color: var(--accent); }

.plugin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.plugin-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: opacity 200ms;
}

[data-theme="light"] .plugin-card {
  border: none;
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}

.plugin-card.is-disabled { opacity: 0.6; }

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.plugin-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
}

.card-badges {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.version-badge {
  font-size: 12px;
  color: var(--text-3);
  font-family: var(--font-mono);
}

.plugin-desc {
  font-size: 13px;
  color: var(--text-2);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
}

.switch-label {
  font-size: 12px;
  color: var(--text-2);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 80px 0;
  color: var(--text-3);
}

.empty-icon { color: var(--text-3); }
.empty-title { font-size: 14px; color: var(--text-1); }
.empty-sub   { font-size: 13px; color: var(--text-2); }
</style>
```

- [ ] **Step 2: Typecheck + commit**

Run: `pnpm typecheck`

Expected: PASS

```bash
git add src/views/PluginView.vue
git commit -m "feat: implement PluginView with card grid and enable/disable toggle"
```

---

### Task 17: SettingsView

**Files:**
- Modify: `src/views/settings/SettingsView.vue`
- Create: `src/views/settings/AppearancePanel.vue`

- [ ] **Step 1: Create `src/views/settings/AppearancePanel.vue`**

```vue
<template>
  <div class="appearance-panel">
    <div class="settings-card">
      <h3 class="card-title">主题</h3>
      <div class="theme-cards">
        <button
          v-for="m in themeModes"
          :key="m.value"
          class="theme-card"
          :class="{ 'is-active': themeStore.mode === m.value }"
          @click="themeStore.setMode(m.value)"
        >
          <div class="theme-preview" :data-preview="m.value" />
          <span class="theme-label">{{ m.label }}</span>
        </button>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">强调色</h3>
      <div class="accent-circles">
        <button
          v-for="(colors, key) in ACCENT_COLORS"
          :key="key"
          class="accent-circle"
          :style="{ background: colors[themeStore.resolvedMode] }"
          :aria-label="`强调色: ${key}`"
          :class="{ 'is-selected': themeStore.accentKey === key }"
          @click="themeStore.setAccent(key as AccentKey)"
        >
          <Check v-if="themeStore.accentKey === key" :size="14" color="#fff" />
        </button>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">字体大小</h3>
      <n-select
        :value="themeStore.fontSize"
        :options="fontSizeOptions"
        size="small"
        style="width: 100px"
        @update:value="themeStore.setFontSize"
      />
    </div>

    <div class="settings-card">
      <h3 class="card-title">语言</h3>
      <n-select
        :value="locale"
        :options="langOptions"
        size="small"
        style="width: 140px"
        @update:value="(v: string) => (locale = v)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Check } from 'lucide-vue-next'
import { useThemeStore, ACCENT_COLORS } from '@/stores/theme'
import type { AccentKey } from '@/stores/theme'
import type { ThemeMode, FontSize } from '@/stores/theme'

const { locale } = useI18n()
const themeStore = useThemeStore()

const themeModes: Array<{ value: ThemeMode; label: string }> = [
  { value: 'system', label: '跟随系统' },
  { value: 'dark',   label: '深色'     },
  { value: 'light',  label: '浅色'     },
]

const fontSizeOptions: Array<{ label: string; value: FontSize }> = [
  { label: '12px', value: '12px' },
  { label: '13px', value: '13px' },
  { label: '14px', value: '14px' },
  { label: '15px', value: '15px' },
]

const langOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English',  value: 'en-US' },
]
</script>

<style scoped>
.appearance-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.settings-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

[data-theme="light"] .settings-card {
  border: none;
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-1);
}

.theme-cards {
  display: flex;
  gap: 12px;
}

.theme-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border: 2px solid var(--border);
  border-radius: 6px;
  background: transparent;
  cursor: pointer;
  transition: border-color 100ms;
}

.theme-card:hover { border-color: var(--text-2); }
.theme-card.is-active { border-color: var(--accent); }

.theme-preview {
  width: 64px;
  height: 40px;
  border-radius: 4px;
}

.theme-preview[data-preview="dark"]   { background: #0d0e13; }
.theme-preview[data-preview="light"]  { background: #f6f7fb; border: 1px solid #dde0ec; }
.theme-preview[data-preview="system"] {
  background: linear-gradient(135deg, #0d0e13 50%, #f6f7fb 50%);
}

.theme-label {
  font-size: 12px;
  color: var(--text-2);
}

.accent-circles {
  display: flex;
  gap: 8px;
}

.accent-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: box-shadow 100ms, transform 100ms;
}

.accent-circle.is-selected {
  box-shadow: 0 0 0 2px var(--bg-base), 0 0 0 4px var(--accent);
  transform: scale(1.05);
}
</style>
```

- [ ] **Step 2: Replace SettingsView stub**

```vue
<template>
  <div class="settings-view">
    <!-- Left submenu -->
    <div class="settings-nav">
      <button
        v-for="item in navItems"
        :key="item.key"
        class="settings-nav-item"
        :class="{ 'is-active': activeSection === item.key }"
        @click="activeSection = item.key"
      >
        <component :is="item.icon" :size="16" />
        <span>{{ item.label }}</span>
      </button>
    </div>

    <!-- Right panels -->
    <div class="settings-content">
      <h2 class="section-heading">{{ currentLabel }}</h2>
      <AppearancePanel v-if="activeSection === 'appearance'" />
      <div v-else class="placeholder-panel">
        <span style="color: var(--text-3); font-size: 13px">{{ currentLabel }} — 配置项待实现</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Palette, Wifi, ShieldCheck, ScrollText, Archive, Info } from 'lucide-vue-next'
import AppearancePanel from './AppearancePanel.vue'

const activeSection = ref('appearance')

const navItems = [
  { key: 'appearance', label: '外观',   icon: Palette      },
  { key: 'connection', label: '连接',   icon: Wifi         },
  { key: 'security',   label: '安全',   icon: ShieldCheck  },
  { key: 'logs',       label: '日志',   icon: ScrollText   },
  { key: 'backup',     label: '备份',   icon: Archive      },
  { key: 'about',      label: '关于',   icon: Info         },
]

const currentLabel = computed(() =>
  navItems.find(n => n.key === activeSection.value)?.label ?? ''
)
</script>

<style scoped>
.settings-view {
  height: 100%;
  display: flex;
  overflow: hidden;
}

.settings-nav {
  width: 180px;
  flex-shrink: 0;
  background: var(--bg-deep);
  border-right: 1px solid var(--border);
  padding: 8px 0;
  display: flex;
  flex-direction: column;
}

.settings-nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 16px;
  border: none;
  border-left: 2px solid transparent;
  background: transparent;
  color: var(--text-2);
  font-size: 13px;
  cursor: pointer;
  text-align: left;
  transition: background 80ms;
}

.settings-nav-item:hover { background: var(--bg-hover); color: var(--text-1); }
.settings-nav-item.is-active {
  border-left-color: var(--accent);
  background: var(--accent-bg);
  color: var(--accent);
  font-weight: 600;
}

.settings-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.section-heading {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-1);
  margin-bottom: 20px;
}

.placeholder-panel {
  padding: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
```

- [ ] **Step 3: Typecheck + commit**

Run: `pnpm typecheck`

Expected: PASS

```bash
git add src/views/settings/
git commit -m "feat: implement SettingsView with functional Appearance panel"
```

---

### Task 18: ConsoleView scaffold

**Files:**
- Modify: `src/views/ConsoleView.vue`

- [ ] **Step 1: Replace ConsoleView stub**

```vue
<template>
  <div class="console-view">
    <!-- Console header -->
    <div class="console-header">
      <div class="console-shell-info">
        <span class="console-icon" aria-hidden="true">◈</span>
        <span class="console-name">{{ shellName }}</span>
        <span class="console-url">{{ shellUrl }}</span>
        <span
          class="console-dot"
          :class="{ 'status-dot-active': isActive }"
          aria-label="状态: 活跃"
        />
      </div>
      <button class="console-close" aria-label="关闭" @click="closeWindow">
        <X :size="14" />
      </button>
    </div>

    <!-- Tab bar -->
    <n-tabs
      v-model:value="activeTab"
      type="card"
      size="small"
      class="console-tabs"
    >
      <n-tab-pane name="file" tab="文件管理">
        <div class="tab-content-placeholder">
          文件管理器将在 Console 功能计划中实现
        </div>
      </n-tab-pane>
      <n-tab-pane name="database" tab="数据库">
        <div class="tab-content-placeholder">
          数据库管理器将在 Console 功能计划中实现
        </div>
      </n-tab-pane>
      <n-tab-pane name="terminal" tab="终端">
        <div class="tab-content-placeholder">
          xterm.js 终端将在 Console 功能计划中实现
        </div>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { X } from 'lucide-vue-next'
import { getCurrentWindow } from '@tauri-apps/api/window'

// Console window reads webshell_id from URL query param (hash router)
// URL format: tauri://localhost/#/console?id={webshell_id}
const route = useRoute()
const webshellId = route.query.id as string | undefined

const shellName = webshellId ? `shell-${webshellId.slice(0, 8)}` : 'Unknown'
const shellUrl  = 'https://example.com/shell.php'
const isActive  = ref(true)
const activeTab = ref('terminal')

async function closeWindow() {
  await getCurrentWindow().close()
}
</script>

<style scoped>
.console-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-base);
}

.console-header {
  height: 32px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background: var(--bg-deep);
  border-bottom: 1px solid var(--border);
  user-select: none;
}

.console-shell-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.console-icon { color: var(--accent); }
.console-name { color: var(--text-1); font-weight: 600; }
.console-url  { color: var(--text-2); font-family: var(--font-mono); font-size: 12px; }

.console-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-success);
}

.console-close {
  width: 32px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--text-2);
  cursor: pointer;
  border-radius: 4px;
}

.console-close:hover { background: var(--wc-close-bg); color: #fff; }

.console-tabs {
  flex: 1;
  overflow: hidden;
}

.tab-content-placeholder {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-3);
  font-size: 13px;
  padding: 40px;
}
</style>
```

- [ ] **Step 2: Final typecheck**

Run: `pnpm typecheck`

Expected: PASS (0 errors)

Run: `pnpm test`

Expected: 8 tests PASS (4 from overrides + 4 from theme store)

- [ ] **Step 3: Final visual verification**

Run: `pnpm tauri dev`

Walk through each route and verify:
- `/` — HomeView: stat cards (all zeros), empty state with CTA button
- `/project` — ProjectView: project tree, empty shell table, toolbar
- `/payload` — PayloadView: config panel left, dark code area right; click "生成载荷" shows mock code
- `/plugin` — PluginView: 3 built-in plugin cards, enable/disable switch works
- `/settings` — SettingsView: click "外观" → theme cards, accent color circles, font size dropdown all work; switching theme and accent updates UI immediately
- Theme toggle in titlebar: cycles dark → light → system
- Sidebar collapse button: collapses to icons, tooltips appear on hover
- ConsoleView: accessible at `/#/console?id=test` — shows header + tab stubs

Stop with Ctrl+C.

- [ ] **Step 4: Commit**

```bash
git add src/views/ConsoleView.vue
git commit -m "feat: implement ConsoleView scaffold with tab structure"
```

---

## Self-Review

### Spec coverage check

| Spec requirement | Covered by task |
|-----------------|----------------|
| Theme tokens (dark/light, 6 accents) | Task 5 |
| buildOverrides (Naive UI overrides) | Task 5 |
| ThemeStore (mode/accent/fontSize, persisted) | Task 6 |
| useThemeCssVars CSS variable sync | Task 6 |
| App shell CSS grid 48/1fr/28 | Task 4, 12 |
| Sidebar collapse (flex transition) | Task 10, 12 |
| CustomTitlebar (drag, controls, toggles) | Task 9 |
| StatusBar (dot, counts) | Task 11 |
| Route transitions (fade 120ms/80ms) | Task 4, 12 |
| HomeView (stat cards, recent list, empty state) | Task 13 |
| ProjectView (project tree, shell table, batch-ready) | Task 14 |
| PayloadView (config panel, key copy, slider badge, generate states) | Task 15 |
| PluginView (card grid, filter tabs, enable/disable) | Task 16 |
| SettingsView (appearance panel fully functional) | Task 17 |
| ConsoleView scaffold (header, tabs) | Task 18 |
| CSS card elevation (dark hairline, light shadow) | Task 13-16 (applied per component) |
| Button hover states | Task 5 (overrides), individual components |
| prefers-reduced-motion block | Task 4 |
| ARIA labels on icon buttons | Task 9, 10, 11, 15 |
| Mono font on technical data | Task 11, 13, 14, 15 |
| Empty states (icon + text + CTA) | Task 13, 16 |
| i18n foundation | Task 7 |
| LoadingMap type | Task 7 |
| AppError type + invoke utility | Task 7 |

### No-placeholder check ✓

All steps contain complete code. No "TBD", "TODO", or "implement later" in any step.

### Type consistency check ✓

- `AccentKey` defined in `tokens.ts`, re-exported through `theme/index.ts`, imported in `stores/theme.ts` and `AppearancePanel.vue`
- `ThemeMode`, `FontSize` defined in `stores/theme.ts`, used in `AppearancePanel.vue`
- `Shell` interface defined locally in `ProjectView.vue` — will be moved to `features/webshell/types.ts` in feature plan
- `buildOverrides(accent, isDark, fontSize)` signature consistent across test file and consumer in `stores/theme.ts`

---

**Plan complete and saved to `docs/superpowers/plans/2026-04-18-ui-skeleton.md`. Two execution options:**

**1. Subagent-Driven (recommended)** — Fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** — Execute tasks in this session using executing-plans, batch with checkpoints

**Which approach?**
