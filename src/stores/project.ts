import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Project {
  id: string
  name: string
  description?: string
  webshellCount: number
  tags?: string[]
  group?: string
  color?: string
  createdAt: number
  updatedAt: number
}

export const useProjectStore = defineStore('project', () => {
  // State
  const projects = ref<Project[]>([])
  const selectedProject = ref<string | null>(null)
  const isLoading = ref(false)
  
  // 搜索/过滤/排序
  const searchQuery = ref('')
  const filterGroup = ref<string | null>(null)
  const sortBy = ref<'name' | 'createdAt' | 'updatedAt' | 'webshellCount'>('name')
  const sortOrder = ref<'asc' | 'desc'>('asc')
  
  // 分组管理
  const groups = ref<Array<{
    id: string
    name: string
    color: string
    order: number
  }>>([])

  // Getters
  const projectCount = computed(() => projects.value.length)
  
  const activeProject = computed(() => 
    projects.value.find(p => p.id === selectedProject.value)
  )

  const filteredProjects = computed(() => {
    let result = [...projects.value]
    
    // 搜索
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      result = result.filter(project =>
        project.name.toLowerCase().includes(query) ||
        project.description?.toLowerCase().includes(query)
      )
    }
    
    // 分组过滤
    if (filterGroup.value) {
      result = result.filter(project => project.group === filterGroup.value)
    }
    
    // 排序
    result.sort((a, b) => {
      const multiplier = sortOrder.value === 'asc' ? 1 : -1
      if (sortBy.value === 'name') {
        return a.name.localeCompare(b.name) * multiplier
      } else if (sortBy.value === 'webshellCount') {
        return (a.webshellCount - b.webshellCount) * multiplier
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

  function setFilterGroup(group: string | null) {
    filterGroup.value = group
  }

  function setSortBy(sortByValue: 'name' | 'createdAt' | 'updatedAt' | 'webshellCount') {
    sortBy.value = sortByValue
  }

  function toggleSortOrder() {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  }

  function clearFilters() {
    searchQuery.value = ''
    filterGroup.value = null
  }

  function selectProject(projectId: string | null) {
    selectedProject.value = projectId
  }

  return {
    // State
    projects,
    selectedProject,
    isLoading,
    searchQuery,
    filterGroup,
    sortBy,
    sortOrder,
    groups,
    // Getters
    projectCount,
    activeProject,
    filteredProjects,
    // Actions
    setSearchQuery,
    setFilterGroup,
    setSortBy,
    toggleSortOrder,
    clearFilters,
    selectProject,
  }
})
