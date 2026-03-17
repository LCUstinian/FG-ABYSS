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
  // 表格容器元素的选择器（用于直接获取可用高度）
  tableContainerSelector?: string
  // 表格区域的固定高度偏移量（像素）- 备选方案，如果不提供 selector 则使用此偏移量
  heightOffset?: number
}

export function useSmartPagination(options: UseSmartPaginationOptions) {
  const {
    total,
    onPageSizeChange,
    onPageChange,
    enableAuto = true,
    tableContainerSelector,
    heightOffset = 250 // 默认偏移量：工具栏 + 侧边栏 + 分页器等
  } = options

  // 当前页码
  const page = ref(1)
  // 当前分页大小 - 默认为 5
  const pageSize = ref(5)
  // 是否手动模式
  const isManualMode = ref(false)
  // 分页大小选项
  const pageSizeOptions = ref(getPageSizeOptions())
  
  // 如果禁用自动模式，强制设置为 5
  if (!enableAuto) {
    pageSize.value = 5
  }

  // 计算当前可视区域高度
  const viewportHeight = ref(window.innerHeight)
  
  // 计算表格区域的可用高度
  const getAvailableHeight = (): number => {
    // 优先使用表格容器的实际可视高度
    if (tableContainerSelector) {
      const container = document.querySelector(tableContainerSelector)
      if (container) {
        const rect = container.getBoundingClientRect()
        console.log('表格容器可视高度:', rect.height)
        return rect.height
      }
    }
    // 备选方案：使用窗口高度减去偏移量
    return window.innerHeight - heightOffset
  }

  // 计算最优分页大小
  const updatePageSize = () => {
    const totalCount = total()
    
    // 计算可用高度
    const availableHeight = getAvailableHeight()
    
    console.log('计算分页大小:', {
      totalCount,
      availableHeight,
      currentPageSize: pageSize.value,
      isManualMode: isManualMode.value
    })
    
    const newSize = calculateOptimalPageSize(totalCount, availableHeight)
    
    console.log('updatePageSize:', {
      newSize,
      currentPageSize: pageSize.value
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

  // 处理窗口大小调整
  const handleResize = () => {
    const oldHeight = viewportHeight.value
    viewportHeight.value = window.innerHeight
    console.log('窗口大小调整:', { oldHeight, newHeight: viewportHeight.value, isManualMode: isManualMode.value })
    if (!isManualMode.value) {
      console.log('自动模式下，准备更新分页大小')
      updatePageSize()
    } else {
      console.log('手动模式下，跳过自动更新')
    }
  }
  
  // 处理表格容器大小变化
  const handleContainerResize = () => {
    if (!isManualMode.value) {
      console.log('表格容器大小变化，重新计算分页')
      updatePageSize()
    }
  }
  
  // 防抖版本
  const debouncedHandleResize = debounce(handleResize, 300)
  const debouncedHandleContainerResize = debounce(handleContainerResize, 300)
  
  // 监听表格容器大小变化（使用 ResizeObserver）
  let resizeObserver: ResizeObserver | null = null
  if (tableContainerSelector && typeof window !== 'undefined' && window.ResizeObserver) {
    resizeObserver = new ResizeObserver(debouncedHandleContainerResize)
  }

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

  // 监听数据总数变化（仅在自动模式下且数据量大幅变化时更新分页大小）
  let previousTotal = 0
  watch(
    () => total(),
    (newTotal) => {
      console.log('watch total:', { newTotal, previousTotal, isManualMode: isManualMode.value, enableAuto })
      
      // 只在自动模式下才响应数据变化
      if (enableAuto && !isManualMode.value) {
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
    console.log('onMounted:', { enableAuto, currentTotal: total(), windowInnerHeight: window.innerHeight })
    
    // 监听表格容器大小变化
    if (tableContainerSelector && resizeObserver) {
      const container = document.querySelector(tableContainerSelector)
      if (container) {
        resizeObserver.observe(container)
        console.log('已启动监听表格容器大小变化:', tableContainerSelector)
      }
    }
    
    if (enableAuto) {
      // 使用可用高度计算初始分页大小
      const totalCount = total()
      const availableHeight = getAvailableHeight()
      const initialSize = calculateOptimalPageSize(totalCount, availableHeight)
      console.log('初始分页大小计算:', { 
        totalCount, 
        availableHeight, 
        initialSize, 
        currentPageSize: pageSize.value 
      })
      // 只在初始分页大小与当前值不同时才更新
      if (pageSize.value !== initialSize) {
        pageSize.value = initialSize
        console.log('设置初始分页大小:', initialSize)
      }
    } else {
      // 禁用自动模式时，使用默认分页大小 5
      pageSize.value = 5
      console.log('自动模式已禁用，使用默认分页大小：5')
    }
    window.addEventListener('resize', debouncedHandleResize)
  })

  // 卸载时移除监听
  onUnmounted(() => {
    window.removeEventListener('resize', debouncedHandleResize)
    // 停止监听表格容器
    if (resizeObserver) {
      resizeObserver.disconnect()
      console.log('已停止监听表格容器')
    }
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
