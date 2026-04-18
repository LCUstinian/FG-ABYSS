<template>
  <div class="app-shell">
    <div class="titlebar-area">
      <CustomTitlebar />
    </div>

    <div class="app-main">
      <div class="sidebar-panel" :class="{ 'is-collapsed': collapsed }">
        <Sidebar />
      </div>
      <div class="content-area">
        <RouterView v-slot="{ Component }">
          <Transition name="fade" mode="out-in">
            <component :is="Component" :key="$route.fullPath" />
          </Transition>
        </RouterView>
      </div>
    </div>

    <div class="statusbar-area">
      <StatusBar />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import CustomTitlebar from './CustomTitlebar.vue'
import Sidebar from './Sidebar.vue'
import StatusBar from './StatusBar.vue'

const collapsed = useLocalStorage('fg-sidebar-collapsed', false)
</script>
