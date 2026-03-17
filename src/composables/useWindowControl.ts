/**
 * 窗口控制 Composable
 * 封装 Tauri 窗口操作逻辑
 */

import { ref, onMounted } from 'vue'
import { getCurrentWindow } from '@tauri-apps/api/window'

export function useWindowControl() {
  const isMaximized = ref(false)
  const appWindow = getCurrentWindow()

  /**
   * 最小化窗口
   */
  const minimizeWindow = async () => {
    try {
      await appWindow.minimize()
    } catch (error) {
      console.error('最小化窗口失败:', error)
    }
  }

  /**
   * 切换最大化状态
   */
  const toggleMaximize = async () => {
    try {
      if (isMaximized.value) {
        await appWindow.unmaximize()
        isMaximized.value = false
      } else {
        await appWindow.maximize()
        isMaximized.value = true
      }
    } catch (error) {
      console.error('切换最大化失败:', error)
    }
  }

  /**
   * 关闭窗口
   */
  const closeWindow = async () => {
    try {
      await appWindow.close()
    } catch (error) {
      console.error('关闭窗口失败:', error)
    }
  }

  /**
   * 检查窗口最大化状态
   */
  const checkMaximizeState = async () => {
    try {
      isMaximized.value = await appWindow.isMaximized()
    } catch (error) {
      console.error('检查最大化状态失败:', error)
    }
  }

  // 组件挂载时检查状态
  onMounted(() => {
    checkMaximizeState()
  })

  return {
    isMaximized,
    minimizeWindow,
    toggleMaximize,
    closeWindow,
    checkMaximizeState
  }
}
