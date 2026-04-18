import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/project',
      component: () => import('@/views/ProjectView.vue'),
    },
    {
      path: '/payload',
      component: () => import('@/views/PayloadView.vue'),
    },
    {
      path: '/plugin',
      component: () => import('@/views/PluginView.vue'),
    },
    {
      path: '/settings',
      component: () => import('@/views/settings/SettingsView.vue'),
    },
    {
      path: '/console',
      component: () => import('@/views/ConsoleView.vue'),
    },
  ],
})

export default router
