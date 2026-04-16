<template>
  <div class="titlebar" @dblclick="toggleMaximize">
    <div class="titlebar-drag-region">
      <div class="titlebar-icon">
        <n-icon size="20" color="#409eff">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path fill="currentColor" d="M12 2A10 10 0 0 0 2 12a10 10 0 0 0 10 10a10 10 0 0 0 10-10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 0-16a8 8 0 0 1 0 16z"/>
            <path fill="currentColor" d="M12 6a1 1 0 0 0-1 1v5a1 1 0 0 0 2 0V7a1 1 0 0 0-1-1zm0 8a1 1 0 1 0 0 2a1 1 0 0 0 0-2z"/>
          </svg>
        </n-icon>
      </div>
      <span class="titlebar-title">FG-ABYSS</span>
    </div>
    <div class="titlebar-controls">
      <button class="titlebar-button" @click="minimize" title="最小化">
        <n-icon size="14">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path fill="currentColor" d="M20 14H4v-2h16"/>
          </svg>
        </n-icon>
      </button>
      <button class="titlebar-button" @click="toggleMaximize" :title="isMaximized ? '还原' : '最大化'">
        <n-icon size="14">
          <svg v-if="!isMaximized" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path fill="currentColor" d="M18 4h-2V2H6v2H4v14h2v2h12v-2h2V4zm0 12H6V6h12v10z"/>
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path fill="currentColor" d="M16 4h2V2H8v2h2v2H6v10h2v2h10v-2h2V8h-4V4zm0 4h2v8H8V8h2V6h6v2z"/>
          </svg>
        </n-icon>
      </button>
      <button class="titlebar-button titlebar-close" @click="closeApp" title="关闭">
        <n-icon size="14">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path fill="currentColor" d="M19 6.41L17.59 5L12 10.59L6.41 5L5 6.41L10.59 12L5 17.59L6.41 19L12 13.41L17.59 19L19 17.59L13.41 12L19 6.41z"/>
          </svg>
        </n-icon>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCurrentWindow } from '@tauri-apps/api/window'

const appWindow = getCurrentWindow()
const isMaximized = ref(false)

const minimize = async () => {
  await appWindow.minimize()
}

const toggleMaximize = async () => {
  await appWindow.toggleMaximize()
  isMaximized.value = await appWindow.isMaximized()
}

const closeApp = async () => {
  await appWindow.close()
}

onMounted(async () => {
  isMaximized.value = await appWindow.isMaximized()
})
</script>

<style scoped>
.titlebar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 40px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  user-select: none;
  -webkit-app-region: drag;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 9999;
}

.titlebar-drag-region {
  display: flex;
  align-items: center;
  height: 100%;
  padding-left: 12px;
  flex: 1;
}

.titlebar-icon {
  margin-right: 8px;
  flex-shrink: 0;
}

.titlebar-title {
  font-size: 14px;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: 0.5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.titlebar-controls {
  display: flex;
  height: 100%;
  -webkit-app-region: no-drag;
  flex-shrink: 0;
}

.titlebar-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: 100%;
  border: none;
  background: transparent;
  color: #ffffff;
  cursor: pointer;
  transition: background-color 0.15s ease;
  padding: 0;
  margin: 0;
}

.titlebar-button:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.titlebar-button:active {
  background-color: rgba(255, 255, 255, 0.15);
}

.titlebar-close:hover {
  background-color: #e81123;
}

.titlebar-close:active {
  background-color: #c50b1a;
}
</style>
