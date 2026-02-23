<template>
  <div class="tooltip-wrapper" ref="wrapperRef">
    <div class="tooltip-content">
      <slot></slot>
    </div>
    <Teleport to="body">
      <div class="tooltip-text" :class="{ 'dark': isDarkTheme }" :style="tooltipStyle">{{ text }}</div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

interface Props {
  text: string
}

const props = defineProps<Props>()

const isDarkTheme = ref(false)
const wrapperRef = ref<HTMLElement | null>(null)
const tooltipStyle = ref<Record<string, string>>({
  left: '0px',
  top: '0px',
  opacity: '0',
  visibility: 'hidden'
})

const updateTheme = () => {
  isDarkTheme.value = document.documentElement.classList.contains('dark')
}

const updateTooltipPosition = () => {
  if (!wrapperRef.value) return
  
  const rect = wrapperRef.value.getBoundingClientRect()
  tooltipStyle.value = {
    left: `${rect.left + rect.width / 2}px`,
    top: `${rect.bottom + 8}px`,
    opacity: '0',
    visibility: 'hidden'
  }
}

const showTooltip = () => {
  if (!wrapperRef.value) return
  
  updateTheme()
  
  const rect = wrapperRef.value.getBoundingClientRect()
  tooltipStyle.value = {
    left: `${rect.left + rect.width / 2}px`,
    top: `${rect.bottom + 8}px`,
    opacity: '1',
    visibility: 'visible'
  }
}

const hideTooltip = () => {
  tooltipStyle.value.opacity = '0'
  tooltipStyle.value.visibility = 'hidden'
}

onMounted(() => {
  updateTheme()
  updateTooltipPosition()
  window.addEventListener('storage', updateTheme)
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateTheme)
  
  if (wrapperRef.value) {
    wrapperRef.value.addEventListener('mouseenter', showTooltip)
    wrapperRef.value.addEventListener('mouseleave', hideTooltip)
  }
})

onUnmounted(() => {
  window.removeEventListener('storage', updateTheme)
  window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateTheme)
  
  if (wrapperRef.value) {
    wrapperRef.value.removeEventListener('mouseenter', showTooltip)
    wrapperRef.value.removeEventListener('mouseleave', hideTooltip)
  }
})
</script>

<style>
.tooltip-wrapper {
  position: relative;
  display: inline-block;
  z-index: 999999;
}

.tooltip-content {
  display: flex;
  align-items: center;
  justify-content: center;
}

.tooltip-text {
  position: fixed;
  transform: translateX(-50%);
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 8px;
  box-shadow: 
    0 4px 16px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  color: #1e293b;
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  pointer-events: none;
  transition: opacity 0.15s ease, transform 0.15s ease;
  z-index: 2147483647;
  overflow: visible;
}

.tooltip-text::before {
  content: '';
  position: absolute;
  top: -12px;
  left: 50%;
  transform: translateX(-50%);
  border: 6px solid transparent;
  border-bottom-color: rgba(255, 255, 255, 0.9);
  filter: drop-shadow(0 -2px 4px rgba(0, 0, 0, 0.05));
  z-index: 2147483647;
  transition: opacity 0.15s ease;
}

.tooltip-text.dark {
  background: rgba(30, 41, 59, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 
    0 4px 16px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.tooltip-text.dark::before {
  border-bottom-color: rgba(30, 41, 59, 0.95);
}
</style>
