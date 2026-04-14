import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import type { Ref } from 'vue'

interface UseScrollLoadOptions {
  /** 距离底部多少像素时触发加载 */
  threshold?: number
  /** 是否立即加载 */
  immediate?: boolean
  /** 加载回调 */
  onLoad: () => Promise<void>
  /** 是否还有更多数据 */
  hasMore?: Ref<boolean>
  /** 是否正在加载 */
  loading?: Ref<boolean>
}

/**
 * 滚动加载 Hook
 * 用于实现无限滚动加载列表
 */
export function useScrollLoad(options: UseScrollLoadOptions) {
  const {
    threshold = 100,
    immediate = true,
    onLoad,
    hasMore,
    loading,
  } = options

  const containerRef = ref<HTMLElement | null>(null)
  const isLoading = ref(false)
  const isFinished = ref(false)

  const handleScroll = async () => {
    if (!containerRef.value || isLoading.value || isFinished.value) {
      return
    }

    if (hasMore && !hasMore.value) {
      isFinished.value = true
      return
    }

    const { scrollTop, scrollHeight, clientHeight } = containerRef.value

    // 距离底部小于阈值时触发加载
    if (scrollHeight - scrollTop - clientHeight < threshold) {
      isLoading.value = true
      if (loading) {
        loading.value = true
      }

      try {
        await onLoad()
      } finally {
        isLoading.value = false
        if (loading) {
          loading.value = false
        }
      }
    }
  }

  const refresh = async () => {
    isLoading.value = true
    if (loading) {
      loading.value = true
    }

    try {
      await onLoad()
    } finally {
      isLoading.value = false
      if (loading) {
        loading.value = false
      }
    }
  }

  onMounted(() => {
    if (immediate) {
      refresh()
    }

    if (containerRef.value) {
      containerRef.value.addEventListener('scroll', handleScroll)
    }
  })

  onBeforeUnmount(() => {
    if (containerRef.value) {
      containerRef.value.removeEventListener('scroll', handleScroll)
    }
  })

  return {
    containerRef,
    isLoading,
    isFinished,
    refresh,
  }
}

interface UseVirtualListOptions<T> {
  /** 列表项高度 */
  itemHeight: number
  /** 缓冲区大小（显示区域上下各显示多少项） */
  overscan?: number
}

interface UseVirtualListReturn<T> {
  /** 可见的列表项 */
  visibleItems: Ref<T[]>
  /** 列表总高度 */
  totalHeight: Ref<string>
  /** 偏移量 */
  offset: Ref<string>
  /** 滚动处理 */
  handleScroll: (e: Event) => void
  /** 容器引用 */
  containerRef: Ref<HTMLElement | null>
}

/**
 * 虚拟列表 Hook
 * 用于优化长列表性能，只渲染可见区域的元素
 */
export function useVirtualList<T>(items: Ref<T[]>, options: UseVirtualListOptions<T>): UseVirtualListReturn<T> {
  const { itemHeight, overscan = 5 } = options

  const containerRef = ref<HTMLElement | null>(null)
  const scrollTop = ref(0)
  const containerHeight = ref(0)
  const visibleItems = ref<T[]>([])

  const totalItems = items.value.length
  const totalHeightValue = totalItems * itemHeight

  const totalHeight = ref(`${totalHeightValue}px`)

  const updateVisibleItems = () => {
    if (!containerRef.value) return

    const startIndex = Math.max(0, Math.floor(scrollTop.value / itemHeight) - overscan)
    const visibleCount = Math.ceil(containerHeight.value / itemHeight)
    const endIndex = Math.min(totalItems, startIndex + visibleCount + overscan * 2)

    visibleItems.value = items.value.slice(startIndex, endIndex)
  }

  const handleScroll = (e: Event) => {
    const target = e.target as HTMLElement
    scrollTop.value = target.scrollTop
    containerHeight.value = target.clientHeight
    updateVisibleItems()
  }

  const refresh = () => {
    updateVisibleItems()
  }

  watch([items, containerRef], () => {
    refresh()
  }, { immediate: true })

  return {
    containerRef,
    visibleItems,
    totalHeight,
    offset: ref('0px'),
    handleScroll,
  }
}

interface UseDebounceOptions {
  /** 延迟时间（毫秒） */
  delay?: number
  /** 是否立即执行 */
  immediate?: boolean
}

/**
 * 防抖 Hook
 * 延迟执行函数，避免频繁触发
 */
export function useDebounce<T extends (...args: any[]) => any>(
  fn: T,
  options: UseDebounceOptions = {}
) {
  const { delay = 300, immediate = false } = options

  let timeoutId: ReturnType<typeof setTimeout> | null = null
  let lastResult: ReturnType<T> | null = null

  const debounced = (...args: Parameters<T>): Promise<ReturnType<T>> => {
    return new Promise((resolve) => {
      if (timeoutId) {
        clearTimeout(timeoutId)
      }

      if (immediate && !timeoutId) {
        lastResult = fn(...args)
        resolve(lastResult)
      }

      timeoutId = setTimeout(() => {
        lastResult = fn(...args)
        resolve(lastResult)
        timeoutId = null
      }, delay)
    })
  }

  const cancel = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
  }

  return {
    debounced,
    cancel,
    lastResult: () => lastResult,
  }
}

/**
 * 节流 Hook
 * 限制函数执行频率
 */
export function useThrottle<T extends (...args: any[]) => any>(
  fn: T,
  delay: number = 300
) {
  let lastCall = 0
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  const throttled = (...args: Parameters<T>): Promise<ReturnType<T>> => {
    return new Promise((resolve) => {
      const now = Date.now()
      const remaining = delay - (now - lastCall)

      if (remaining <= 0) {
        lastCall = now
        resolve(fn(...args))
      } else {
        if (timeoutId) {
          clearTimeout(timeoutId)
        }

        timeoutId = setTimeout(() => {
          lastCall = Date.now()
          timeoutId = null
          resolve(fn(...args))
        }, remaining)
      }
    })
  }

  const cancel = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
  }

  return {
    throttled,
    cancel,
  }
}
