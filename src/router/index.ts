import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
  },
  {
    path: '/project',
    name: 'Project',
    component: () => import('@/views/ProjectView.vue'),
  },
  {
    path: '/payload',
    name: 'Payload',
    component: () => import('@/views/PayloadView.vue'),
  },
  {
    path: '/plugin',
    name: 'Plugin',
    component: () => import('@/views/PluginView.vue'),
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue'),
  },
  {
    path: '/console',
    name: 'Console',
    component: () => import('@/views/ConsoleView.vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
