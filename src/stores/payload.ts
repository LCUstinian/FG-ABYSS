import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface PayloadConfig {
  id: string
  name: string
  url: string
  password: string
  payloadType: 'php' | 'jsp' | 'asp' | 'aspx'
  encryption: 'aes-256-gcm' | 'xor' | 'base64'
  obfuscationLevel: 1 | 2 | 3
  tags?: string[]
  group?: string
  createdAt: number
  updatedAt: number
}

export interface PayloadTemplate {
  id: string
  name: string
  type: 'php' | 'jsp' | 'asp' | 'aspx'
  code: string
  isBuiltIn: boolean
}

export interface PayloadHistory {
  id: string
  configId: string
  generatedCode: string
  generatedAt: number
}

export const usePayloadStore = defineStore('payload', () => {
  // State
  const configs = ref<PayloadConfig[]>([])
  const templates = ref<PayloadTemplate[]>([])
  const history = ref<PayloadHistory[]>([])
  const selectedConfig = ref<PayloadConfig | null>(null)
  const isGenerating = ref(false)
  
  // 搜索/过滤/排序
  const searchQuery = ref('')
  const filterType = ref<'all' | 'php' | 'jsp' | 'asp' | 'aspx'>('all')
  const filterTag = ref<string | null>(null)
  const sortBy = ref<'name' | 'createdAt' | 'updatedAt' | 'type'>('name')
  const sortOrder = ref<'asc' | 'desc'>('asc')
  const tags = ref<string[]>([])

  // Getters
  const configCount = computed(() => configs.value.length)
  const templateCount = computed(() => templates.value.length)
  const historyCount = computed(() => history.value.length)

  const filteredConfigs = computed(() => {
    let result = [...configs.value]
    
    // 搜索
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      result = result.filter(config =>
        config.name.toLowerCase().includes(query) ||
        config.url.toLowerCase().includes(query)
      )
    }
    
    // 类型过滤
    if (filterType.value !== 'all') {
      result = result.filter(config => config.payloadType === filterType.value)
    }
    
    // 标签过滤
    if (filterTag.value) {
      result = result.filter(config => config.tags?.includes(filterTag.value))
    }
    
    // 排序
    result.sort((a, b) => {
      const multiplier = sortOrder.value === 'asc' ? 1 : -1
      if (sortBy.value === 'name') {
        return a.name.localeCompare(b.name) * multiplier
      } else if (sortBy.value === 'type') {
        return a.payloadType.localeCompare(b.payloadType) * multiplier
      } else {
        return (a[sortBy.value] - b[sortBy.value]) * multiplier
      }
    })
    
    return result
  })

  // Actions
  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function setFilterType(type: 'all' | 'php' | 'jsp' | 'asp' | 'aspx') {
    filterType.value = type
  }

  function setFilterTag(tag: string | null) {
    filterTag.value = tag
  }

  function clearFilters() {
    searchQuery.value = ''
    filterType.value = 'all'
    filterTag.value = null
  }

  function setSortBy(sortByValue: 'name' | 'createdAt' | 'updatedAt' | 'type') {
    sortBy.value = sortByValue
  }

  function toggleSortOrder() {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  }

  return {
    // State
    configs,
    templates,
    history,
    selectedConfig,
    isGenerating,
    searchQuery,
    filterType,
    filterTag,
    sortBy,
    sortOrder,
    tags,
    // Getters
    configCount,
    templateCount,
    historyCount,
    filteredConfigs,
    // Actions
    setSearchQuery,
    setFilterType,
    setFilterTag,
    clearFilters,
    setSortBy,
    toggleSortOrder,
  }
})
