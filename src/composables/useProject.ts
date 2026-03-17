/**
 * 项目管理 Composable
 * 封装项目相关的业务逻辑和状态管理
 */

import { ref, computed } from 'vue'
import { invoke } from '@/utils/tauri-mock-adapter'
import type { Project } from '@/types'

export function useProject() {
  const projects = ref<Project[]>([])
  const selectedProject = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  /**
   * 获取项目列表
   */
  const fetchProjects = async () => {
    loading.value = true
    error.value = null
    try {
      const projectList = await invoke<Project[]>('get_projects')
      projects.value = projectList.sort((a, b) => {
        const timeA = new Date(a.createdAt).getTime()
        const timeB = new Date(b.createdAt).getTime()
        return timeA - timeB
      })
      
      // 如果没有选中项目，且项目列表不为空，选择第一个项目
      if (!selectedProject.value && projects.value.length > 0) {
        selectedProject.value = projects.value[0].id
      }
    } catch (err: any) {
      error.value = err.message || '获取项目列表失败'
      console.error('获取项目列表失败:', err)
    } finally {
      loading.value = false
    }
  }

  /**
   * 创建项目
   */
  const createProject = async (name: string, description?: string) => {
    error.value = null
    try {
      const result = await invoke('create_project', { name, description })
      if (result.success) {
        await fetchProjects()
        return result.project
      }
    } catch (err: any) {
      error.value = err.message || '创建项目失败'
      throw err
    }
  }

  /**
   * 删除项目
   */
  const deleteProject = async (projectId: string) => {
    error.value = null
    try {
      await invoke('delete_project', { projectId })
      await fetchProjects()
    } catch (err: any) {
      error.value = err.message || '删除项目失败'
      throw err
    }
  }

  /**
   * 恢复已删除的项目
   */
  const recoverProject = async (projectId: string) => {
    error.value = null
    try {
      await invoke('recover_project', { projectId })
      await fetchProjects()
    } catch (err: any) {
      error.value = err.message || '恢复项目失败'
      throw err
    }
  }

  /**
   * 选中项目
   */
  const selectProject = (projectId: string) => {
    selectedProject.value = projectId
  }

  /**
   * 获取当前选中的项目
   */
  const currentProject = computed(() => {
    return projects.value.find(p => p.id === selectedProject.value)
  })

  return {
    projects,
    selectedProject,
    loading,
    error,
    currentProject,
    fetchProjects,
    createProject,
    deleteProject,
    recoverProject,
    selectProject
  }
}
