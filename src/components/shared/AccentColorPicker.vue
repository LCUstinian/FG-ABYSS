<template>
  <div class="accent-color-picker">
    <!-- 预设颜色选择区 -->
    <div class="preset-colors-section">
      <div class="section-title">
        <span class="title-icon">🎨</span>
        <span>{{ t('settings.presetColors') }}</span>
      </div>
      <div class="preset-colors-container">
        <button
          v-for="(color, index) in presetColors"
          :key="index"
          class="preset-color-btn"
          :class="{ active: selectedColor === color.value }"
          :style="{ backgroundColor: color.value }"
          @click="selectColor(color.value)"
          :title="color.name"
        >
          <span v-if="selectedColor === color.value" class="check-icon">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="white">
              <path fill-rule="evenodd" d="M13.707 5.293a1 1 0 010 1.414l-6 6a1 1 0 01-1.414 0l-3-3a1 1 0 011.414-1.414L7 10.586l5.293-5.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
            </svg>
          </span>
        </button>
      </div>
    </div>

    <!-- 自定义颜色选择器 -->
    <div class="custom-color-section">
      <div class="section-title">
        <span class="title-icon">⚙️</span>
        <span>{{ t('settings.customColor') }}</span>
      </div>
      
      <div class="custom-color-container">
        <!-- 颜色预览区 -->
        <div class="color-preview-section">
          <div class="color-preview" :style="{ backgroundColor: customColor }">
            <div class="preview-content">
              <span class="preview-label">{{ t('settings.preview') }}</span>
              <span class="color-code">{{ customColor.toUpperCase() }}</span>
            </div>
          </div>
        </div>

        <!-- 颜色选择器 -->
        <div class="color-picker-controls">
          <!-- 颜色滑块 -->
          <div class="slider-group">
            <div class="slider-item">
              <label class="slider-label">
                <span>{{ t('settings.hue') }}</span>
                <span class="slider-value">{{ hsl.h }}°</span>
              </label>
              <div class="hue-slider-container">
                <input
                  type="range"
                  min="0"
                  max="360"
                  v-model.number="hsl.h"
                  class="slider hue-slider"
                  :style="{ background: hueGradient }"
                />
              </div>
            </div>

            <div class="slider-item">
              <label class="slider-label">
                <span>{{ t('settings.saturation') }}</span>
                <span class="slider-value">{{ hsl.s }}%</span>
              </label>
              <input
                type="range"
                min="0"
                max="100"
                v-model.number="hsl.s"
                class="slider"
              />
            </div>

            <div class="slider-item">
              <label class="slider-label">
                <span>{{ t('settings.lightness') }}</span>
                <span class="slider-value">{{ hsl.l }}%</span>
              </label>
              <input
                type="range"
                min="0"
                max="100"
                v-model.number="hsl.l"
                class="slider"
              />
            </div>
          </div>

          <!-- 颜色输入框 -->
          <div class="input-group">
            <div class="input-item">
              <label class="input-label">HEX</label>
              <input
                type="text"
                v-model="hexInput"
                class="color-input hex-input"
                :class="{ error: hexError }"
                placeholder="#000000"
                @input="validateHexInput"
              />
            </div>

            <div class="input-item">
              <label class="input-label">RGB</label>
              <div class="rgb-inputs">
                <input
                  type="number"
                  v-model.number="rgb.r"
                  class="color-input rgb-input"
                  min="0"
                  max="255"
                  placeholder="0"
                  @input="updateFromRGB"
                />
                <input
                  type="number"
                  v-model.number="rgb.g"
                  class="color-input rgb-input"
                  min="0"
                  max="255"
                  placeholder="0"
                  @input="updateFromRGB"
                />
                <input
                  type="number"
                  v-model.number="rgb.b"
                  class="color-input rgb-input"
                  min="0"
                  max="255"
                  placeholder="0"
                  @input="updateFromRGB"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 应用按钮 -->
    <div class="action-buttons">
      <button
        class="btn-apply"
        :disabled="!canApply"
        @click="applyColor"
      >
        <span class="btn-icon">✓</span>
        <span>{{ t('settings.applyColor') }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface Props {
  modelValue?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '#3b82f6'
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  apply: [value: string]
}>()

// 预设颜色（10 种）
const presetColors = [
  { value: '#3b82f6', name: 'Blue' },
  { value: '#8b5cf6', name: 'Purple' },
  { value: '#ec4899', name: 'Pink' },
  { value: '#f59e0b', name: 'Amber' },
  { value: '#10b981', name: 'Emerald' },
  { value: '#06b6d4', name: 'Cyan' },
  { value: '#ef4444', name: 'Red' },
  { value: '#84cc16', name: 'Lime' },
  { value: '#f97316', name: 'Orange' },
  { value: '#14b8a6', name: 'Teal' },
]

// 当前选中的颜色
const selectedColor = ref(props.modelValue)
const customColor = ref(props.modelValue)

// HSL 值
const hsl = ref({ h: 210, s: 100, l: 50 })

// RGB 值
const rgb = ref({ r: 59, g: 130, b: 246 })

// HEX 输入
const hexInput = ref('#3b82f6')

// HEX 错误状态
const hexError = ref(false)

// 色相渐变
const hueGradient = computed(() => {
  const colors = [
    'hsl(0, 100%, 50%)',
    'hsl(60, 100%, 50%)',
    'hsl(120, 100%, 50%)',
    'hsl(180, 100%, 50%)',
    'hsl(240, 100%, 50%)',
    'hsl(300, 100%, 50%)',
    'hsl(360, 100%, 50%)'
  ]
  return `linear-gradient(to right, ${colors.join(', ')})`
})

// 是否可以应用
const canApply = computed(() => {
  return !hexError.value && customColor.value.trim() !== ''
})

// 选择预设颜色
const selectColor = (color: string) => {
  selectedColor.value = color
  customColor.value = color
  updateFromHex(color)
}

// 验证 HEX 输入
const validateHexInput = () => {
  const hex = hexInput.value.trim()
  const hexRegex = /^#[0-9A-Fa-f]{6}$/
  
  if (hexRegex.test(hex)) {
    hexError.value = false
    updateFromHex(hex)
  } else {
    hexError.value = true
  }
}

// 从 HEX 更新
const updateFromHex = (hex: string) => {
  const cleanHex = hex.replace('#', '')
  const r = parseInt(cleanHex.substring(0, 2), 16)
  const g = parseInt(cleanHex.substring(2, 4), 16)
  const b = parseInt(cleanHex.substring(4, 6), 16)
  
  rgb.value = { r, g, b }
  hsl.value = rgbToHsl(r, g, b)
  customColor.value = hex
}

// 从 RGB 更新
const updateFromRGB = () => {
  const { r, g, b } = rgb.value
  const hex = rgbToHex(r, g, b)
  
  hexInput.value = hex
  hsl.value = rgbToHsl(r, g, b)
  customColor.value = hex
}

// RGB 转 HEX
const rgbToHex = (r: number, g: number, b: number): string => {
  const toHex = (n: number) => {
    const clamped = Math.max(0, Math.min(255, n))
    const hex = clamped.toString(16)
    return hex.length === 1 ? '0' + hex : hex
  }
  return '#' + toHex(r) + toHex(g) + toHex(b)
}

// RGB 转 HSL
const rgbToHsl = (r: number, g: number, b: number) => {
  r /= 255
  g /= 255
  b /= 255
  
  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  let h = 0
  let s = 0
  const l = (max + min) / 2
  
  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    
    switch (max) {
      case r:
        h = ((g - b) / d + (g < b ? 6 : 0)) / 6
        break
      case g:
        h = ((b - r) / d + 2) / 6
        break
      case b:
        h = ((r - g) / d + 4) / 6
        break
    }
  }
  
  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    l: Math.round(l * 100)
  }
}

// HSL 转 RGB
const hslToRgb = (h: number, s: number, l: number) => {
  h /= 360
  s /= 100
  l /= 100
  
  let r: number, g: number, b: number
  
  if (s === 0) {
    r = g = b = l
  } else {
    const hue2rgb = (p: number, q: number, t: number) => {
      if (t < 0) t += 1
      if (t > 1) t -= 1
      if (t < 1/6) return p + (q - p) * 6 * t
      if (t < 1/2) return q
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6
      return p
    }
    
    const q = l < 0.5 ? l * (1 + s) : l + s - l * s
    const p = 2 * l - q
    r = hue2rgb(p, q, h + 1/3)
    g = hue2rgb(p, q, h)
    b = hue2rgb(p, q, h - 1/3)
  }
  
  return {
    r: Math.round(r * 255),
    g: Math.round(g * 255),
    b: Math.round(b * 255)
  }
}

// 监听 HSL 变化
watch(() => hsl.value, (newHsl) => {
  const rgbValue = hslToRgb(newHsl.h, newHsl.s, newHsl.l)
  rgb.value = rgbValue
  const hex = rgbToHex(rgbValue.r, rgbValue.g, rgbValue.b)
  hexInput.value = hex
  customColor.value = hex
}, { deep: true })

// 应用颜色
const applyColor = () => {
  emit('update:modelValue', customColor.value)
  emit('apply', customColor.value)
  selectedColor.value = customColor.value
}

// 监听外部值变化
watch(() => props.modelValue, (newValue) => {
  if (newValue !== customColor.value) {
    updateFromHex(newValue)
    selectedColor.value = newValue
  }
}, { immediate: true })
</script>

<style scoped>
.accent-color-picker {
  display: flex;
  flex-direction: column;
  gap: 28px;
  padding: 24px;
  background: var(--card-bg);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
}

.dark .accent-color-picker {
  border-color: var(--border-strong);
}

/* ===== 区域标题 ===== */
.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 16px;
}

.title-icon {
  font-size: 18px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

/* ===== 预设颜色区 ===== */
.preset-colors-container {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

.preset-color-btn {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  border: 3px solid var(--border-color);
  cursor: pointer;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  flex-shrink: 0;
}

.dark .preset-color-btn {
  border-color: var(--border-strong);
}

.preset-color-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.3), transparent);
  border-radius: 50%;
}

.preset-color-btn:hover {
  transform: scale(1.15);
  box-shadow: var(--shadow-md);
  border-color: var(--active-color);
}

.preset-color-btn.active {
  border-color: white;
  box-shadow: 0 0 0 4px var(--active-color-suppl), var(--shadow-md);
  animation: pulse-ring 2s infinite;
}

.check-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 1;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
  animation: checkmark-bounce 0.4s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

@keyframes pulse-ring {
  0%, 100% {
    box-shadow: 0 0 0 4px var(--active-color-suppl), var(--shadow-md);
  }
  50% {
    box-shadow: 0 0 0 8px var(--active-color-suppl), var(--shadow-md);
  }
}

@keyframes checkmark-bounce {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

/* ===== 自定义颜色区 ===== */
.custom-color-section {
  flex: 1;
}

.custom-color-container {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 20px;
  align-items: start;
}

/* 颜色预览 */
.color-preview-section {
  position: sticky;
  top: 0;
}

.color-preview {
  width: 100%;
  aspect-ratio: 1;
  border-radius: var(--border-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  border: 3px solid rgba(255, 255, 255, 0.15);
  transition: all var(--transition-normal);
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(255,255,255,0.1), rgba(255,255,255,0.05));
}

.dark .color-preview {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  border-color: rgba(255, 255, 255, 0.08);
}

.color-preview::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.15), transparent 40%, rgba(0,0,0,0.1) 60%, rgba(0,0,0,0.05));
  border-radius: var(--border-radius-lg);
}

.preview-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  position: relative;
  z-index: 1;
  padding: 12px;
}

.preview-label {
  font-size: 12px;
  color: white;
  opacity: 0.85;
  font-weight: 500;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.4);
  letter-spacing: 0.3px;
}

.color-code {
  font-size: 20px;
  color: white;
  font-weight: 700;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  text-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
  letter-spacing: 0.8px;
  background: rgba(0, 0, 0, 0.15);
  padding: 4px 12px;
  border-radius: 12px;
  backdrop-filter: blur(8px);
}

/* 颜色选择器控制区 */
.color-picker-controls {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

/* 滑块组 */
.slider-group {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.slider-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.slider-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.4px;
}

.slider-value {
  color: var(--text-color);
  font-weight: 700;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 11px;
  background: var(--content-bg);
  padding: 2px 8px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
}

.dark .slider-value {
  border-color: var(--border-strong);
}

.slider {
  width: 100%;
  height: 10px;
  border-radius: 5px;
  appearance: none;
  background: var(--content-bg);
  border: 1px solid var(--border-color);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.dark .slider {
  border-color: var(--border-strong);
}

.slider:hover {
  border-color: var(--active-color);
}

.slider::-webkit-slider-thumb {
  appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: white;
  border: 2px solid var(--active-color);
  cursor: ew-resize;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  transition: all var(--transition-fast);
}

.slider::-webkit-slider-thumb:hover {
  transform: scale(1.15);
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.3);
  border-color: var(--active-color);
}

.slider::-webkit-slider-thumb:active {
  transform: scale(1.1);
}

.slider::-moz-range-thumb {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: white;
  border: 2px solid var(--active-color);
  cursor: ew-resize;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  transition: all var(--transition-fast);
}

.slider::-moz-range-thumb:hover {
  transform: scale(1.15);
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.3);
}

.hue-slider {
  background: transparent;
  border: none;
}

/* 输入框组 */
.input-group {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 14px;
}

.input-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.input-label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.color-input {
  width: 100%;
  padding: 9px 12px;
  border-radius: var(--border-radius-md);
  border: 1.5px solid var(--border-color);
  background: var(--content-bg);
  color: var(--text-color);
  font-size: 13px;
  font-family: 'Monaco', 'Consolas', monospace;
  transition: all var(--transition-fast);
  text-align: center;
}

.dark .color-input {
  border-color: var(--border-strong);
}

.color-input:focus {
  outline: none;
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px var(--active-color-suppl);
  background: var(--card-bg);
}

.color-input.error {
  border-color: #ef4444;
  animation: shake 0.3s ease-in-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  50% { transform: translateX(4px); }
  75% { transform: translateX(-4px); }
}

.hex-input {
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.rgb-inputs {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.rgb-input {
  text-align: center;
}

/* ===== 应用按钮 ===== */
.action-buttons {
  display: flex;
  justify-content: center;
  padding-top: 8px;
  border-top: 1px solid var(--border-color);
}

.dark .action-buttons {
  border-top-color: var(--border-strong);
}

.btn-apply {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 14px 32px;
  background: linear-gradient(135deg, var(--active-color), var(--active-color-suppl));
  color: white;
  border: none;
  border-radius: var(--border-radius-md);
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
  box-shadow: var(--shadow-md);
}

.btn-apply:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.btn-apply:active:not(:disabled) {
  transform: translateY(0);
}

.btn-apply:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon {
  font-size: 18px;
  font-weight: 700;
}

/* ===== 响应式设计 ===== */
@media (max-width: 1024px) {
  .custom-color-container {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .color-preview-section {
    position: static;
  }
  
  .color-preview {
    max-width: 400px;
    margin: 0 auto;
    aspect-ratio: 2/1;
  }
  
  .input-group {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .accent-color-picker {
    padding: 20px;
  }
  
  .preset-colors-container {
    gap: 14px;
  }
  
  .preset-color-btn {
    width: 46px;
    height: 46px;
  }
  
  .color-preview {
    aspect-ratio: 3/2;
  }
  
  .slider {
    height: 9px;
  }
  
  .slider::-webkit-slider-thumb {
    width: 17px;
    height: 17px;
  }
  
  .color-input {
    padding: 8px 10px;
    font-size: 13px;
  }
  
  .btn-apply {
    padding: 12px 28px;
    font-size: 14px;
  }
}

@media (max-width: 480px) {
  .accent-color-picker {
    padding: 16px;
    gap: 20px;
  }
  
  .section-title {
    font-size: 14px;
  }
  
  .preset-colors-container {
    gap: 10px;
  }
  
  .preset-color-btn {
    width: 40px;
    height: 40px;
  }
  
  .color-preview {
    aspect-ratio: 4/3;
    border-radius: var(--border-radius-md);
  }
  
  .color-code {
    font-size: 18px;
    padding: 3px 10px;
  }
  
  .slider-item {
    gap: 5px;
  }
  
  .slider-label {
    font-size: 11px;
  }
  
  .slider-value {
    font-size: 10px;
    padding: 2px 6px;
  }
  
  .slider {
    height: 8px;
  }
  
  .slider::-webkit-slider-thumb {
    width: 16px;
    height: 16px;
  }
  
  .input-group {
    gap: 12px;
  }
  
  .input-label {
    font-size: 11px;
  }
  
  .color-input {
    padding: 7px 9px;
    font-size: 12px;
  }
  
  .rgb-inputs {
    gap: 6px;
  }
  
  .btn-apply {
    width: 100%;
    justify-content: center;
    padding: 12px 24px;
  }
}

/* ===== 超小屏幕优化 ===== */
@media (max-width: 360px) {
  .preset-color-btn {
    width: 36px;
    height: 36px;
  }
  
  .color-preview {
    aspect-ratio: 1/1;
  }
  
  .color-code {
    font-size: 16px;
    padding: 2px 8px;
  }
  
  .slider {
    height: 7px;
  }
  
  .slider::-webkit-slider-thumb {
    width: 15px;
    height: 15px;
  }
}
</style>
