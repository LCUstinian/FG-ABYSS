// 分页配置常量
const PAGINATION_CONFIG = {
  // 每行高度（像素）- 根据实际表格行高调整
  ROW_HEIGHT: 53,
  
  // 表格头部和底部占用的固定高度（像素）
  TABLE_HEADER_HEIGHT: 50,
  TABLE_FOOTER_HEIGHT: 70,
  
  // 其他 UI 元素占用的高度（搜索栏、按钮等）
  OTHER_UI_HEIGHT: 100,
  
  // 最小分页大小
  MIN_PAGE_SIZE: 5,
  
  // 最大分页大小
  MAX_PAGE_SIZE: 50,
  
  // 预设的分页大小选项
  PAGE_SIZE_OPTIONS: [5, 10, 20, 30, 50],
} as const

/**
 * 计算最优分页大小
 * @param dataCount 数据总条数
 * @param availableHeight 可用高度（像素）- 已经减去了所有固定元素高度的净可用高度
 * @returns 最优分页大小
 */
export function calculateOptimalPageSize(
  dataCount: number,
  availableHeight: number
): number {
  // availableHeight 已经是可用高度，直接计算理想行数
  const idealRows = Math.floor(availableHeight / PAGINATION_CONFIG.ROW_HEIGHT)
  
  // 确保行数在合理范围内
  let optimalSize = Math.max(
    PAGINATION_CONFIG.MIN_PAGE_SIZE,
    Math.min(idealRows, PAGINATION_CONFIG.MAX_PAGE_SIZE)
  )
  
  // 从预设选项中选择最接近的值
  const closestOption = PAGINATION_CONFIG.PAGE_SIZE_OPTIONS.reduce((prev, curr) => {
    return Math.abs(curr - optimalSize) < Math.abs(prev - optimalSize) ? curr : prev
  })
  
  // 如果数据量小于最优分页大小，则调整为数据量（但不小于最小值）
  if (dataCount > 0 && dataCount < closestOption) {
    optimalSize = Math.max(dataCount, PAGINATION_CONFIG.MIN_PAGE_SIZE)
  } else {
    optimalSize = closestOption
  }
  
  // 确保不超过最大值
  return Math.min(optimalSize, PAGINATION_CONFIG.MAX_PAGE_SIZE)
}

/**
 * 防抖函数
 * @param fn 需要防抖的函数
 * @param delay 延迟时间（毫秒）
 * @returns 防抖后的函数
 */
export function debounce<T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null
  
  return function (this: any, ...args: Parameters<T>) {
    if (timeoutId !== null) {
      clearTimeout(timeoutId)
    }
    
    timeoutId = setTimeout(() => {
      fn.apply(this, args)
      timeoutId = null
    }, delay)
  }
}

/**
 * 获取分页大小选项
 * @returns 分页大小选项数组
 */
export function getPageSizeOptions(): number[] {
  return [...PAGINATION_CONFIG.PAGE_SIZE_OPTIONS]
}

export { PAGINATION_CONFIG }
