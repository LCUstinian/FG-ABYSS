import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface WebShell {
  id: string
  projectId: string
  name: string
  url: string
  password: string
  payloadType: 'php' | 'jsp' | 'asp' | 'aspx'
  status: 'online' | 'offline' | 'unknown'
  lastConnectedAt?: number
  createdAt: number
  updatedAt: number
}

export interface WebShellSession {
  webshellId: string
  isConnected: boolean
  lastActivity: number
  metadata?: Record<string, any>
}

export const useWebShellStore = defineStore('webshell', () => {
  // State
  const webshells = ref<WebShell[]>([])
  const sessions = ref<Record<string, WebShellSession>>({})
  const activeSession = ref<string | null>(null)
  const isConnecting = ref(false)

  // Getters
  const webshellCount = computed(() => webshells.value.length)
  const onlineCount = computed(() => 
    webshells.value.filter(w => w.status === 'online').length
  )
  const activeWebShell = computed(() => {
    if (!activeSession.value) return null
    return webshells.value.find(w => w.id === activeSession.value)
  })

  // Actions
  function setActiveSession(sessionId: string | null) {
    activeSession.value = sessionId
  }

  function updateSession(webshellId: string, session: Partial<WebShellSession>) {
    if (sessions.value[webshellId]) {
      sessions.value[webshellId] = {
        ...sessions.value[webshellId],
        ...session,
      }
    } else {
      sessions.value[webshellId] = {
        webshellId,
        isConnected: false,
        lastActivity: Date.now(),
        ...session,
      }
    }
  }

  return {
    // State
    webshells,
    sessions,
    activeSession,
    isConnecting,
    // Getters
    webshellCount,
    onlineCount,
    activeWebShell,
    // Actions
    setActiveSession,
    updateSession,
  }
})
