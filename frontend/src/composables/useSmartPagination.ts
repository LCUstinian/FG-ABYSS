import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { calculateOptimalPageSize, debounce, getPageSizeOptions } from '@/utils/pagination'

interface UseSmartPaginationOptions {
  // 数据总数
  total: () => number
  // 分页大小变化时的回调
  onPageSizeChange?: (size: number) => void
  // 页码变化时的回调
  onPageChange?: (page: number) => void
  // 是否启用自动分页（默认启用）
  enableAuto?: boolean
}

export function useSmartPagination(options: UseSmartPaginationOptions) {
  const {
    total,
    onPageSizeChange,
    onPageChange,
    enableAuto = true
  } = options

  // 当前页码
  const page = ref(1)
  // 当前分页大小
  const pageSize = ref(5)
  // 是否手动模式
  const isManualMode = ref(false)
  // 分页大小选项
  const pageSizeOptions = ref(getPageSizeOptions())

  // 计算当前可视区域高度
  const viewportHeight = ref(window.innerHeight)

  // 计算最优分页大小
  const updatePageSize = () => {
    const totalCount = total()
    const newSize = calculateOptimalPageSize(totalCount, viewportHeight.value)
    
    console.log('updatePageSize:', {
      totalCount,
      viewportHeight: viewportHeight.value,
      newSize,
      currentPageSize: pageSize.value,
      isManualMode: isManualMode.value
    })
    
    // 只在自动模式下且分页大小真正变化时才更新
    if (!isManualMode.value && pageSize.value !== newSize) {
      console.log('自动更新分页大小:', pageSize.value, '->', newSize)
      pageSize.value = newSize
      // 使用 nextTick 避免在 watch 中立即触发
      setTimeout(() => {
        onPageSizeChange?.(newSize)
      }, 0)
    } else if (isManualMode.value) {
      console.log('手动模式，跳过更新')
    } else {
      console.log('分页大小未变化，跳过更新')
    }
  }

  // 防抖的窗口大小调整处理
  const handleResize = debounce(() => {
    const oldHeight = viewportHeight.value
    viewportHeight.value = window.innerHeight
    console.log('窗口大小调整:', { oldHeight, newHeight: viewportHeight.value, isManualMode: isManualMode.value })
    if (!isManualMode.value) {
      console.log('自动模式下，准备更新分页大小')
      updatePageSize()
    } else {
      console.log('手动模式下，跳过自动更新')
    }
  }, 300)

  // 切换到手动模式
  const setManualMode = (manual: boolean) => {
    isManualMode.value = manual
  }

  // 手动设置分页大小
  const setPageSize = (size: number) => {
    pageSize.value = size
    page.value = 1
    setManualMode(true)
    onPageSizeChange?.(size)
  }

  // 设置页码
  const setPage = (newPage: number) => {
    page.value = newPage
    onPageChange?.(newPage)
  }

  // 重置为自动模式
  const resetToAuto = () => {
    setManualMode(false)
    updatePageSize()
  }

  // 监听数据总数变化（仅在首次加载数据或数据量大幅变化时更新分页大小）
  let previousTotal = 0
  watch(
    () => total(),
    (newTotal) => {
      console.log('watch total:', { newTotal, previousTotal, isManualMode: isManualMode.value })
      
      if (!isManualMode.value) {
        // 当数据总量从 0 变为非 0 时（首次加载数据），更新分页大小
        if (previousTotal === 0 && newTotal > 0) {
          previousTotal = newTotal
          updatePageSize()
          return
        }
        
        // 当数据总量变化超过 50% 时，更新分页大小
        if (previousTotal > 0) {
          const changeRatio = Math.abs(newTotal - previousTotal) / previousTotal
          console.log('数据变化比例:', changeRatio)
          if (changeRatio > 0.5) {
            previousTotal = newTotal
            updatePageSize()
          } else {
            previousTotal = newTotal
          }
        }
      }
    }
  )

  // 注意：不在这里监听 pageSize 变化，因为 onPageSizeChange 回调已经处理了

  // 挂载时设置初始分页大小并监听窗口变化
  onMounted(() => {
    console.log('onMounted:', { enableAuto, currentTotal: total(), viewportHeight: viewportHeight.value })
    if (enableAuto) {
      // 挂载时先基于窗口高度计算一个初始分页大小（此时数据可能还未加载）
      const totalCount = total()
      const initialSize = calculateOptimalPageSize(totalCount, viewportHeight.value)
      console.log('初始分页大小计算:', { totalCount, viewportHeight: viewportHeight.value, initialSize, currentPageSize: pageSize.value })
      // 只在初始分页大小与当前值不同时才更新
      if (pageSize.value !== initialSize) {
        pageSize.value = initialSize
        console.log('设置初始分页大小:', initialSize)
      }
    }
    window.addEventListener('resize', handleResize)
  })

  // 卸载时移除监听
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
  })

  return {
    // 响应式状态
    page,
    pageSize,
    isManualMode,
    pageSizeOptions,
    
    // 方法
    setPageSize,
    setPage,
    setManualMode,
    resetToAuto,
    updatePageSize
  }
}
