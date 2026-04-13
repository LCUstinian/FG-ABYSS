<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <n-message-provider>
      <n-dialog-provider>
        <AppLayout>
          <router-view />
        </AppLayout>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { darkTheme } from 'naive-ui'
import AppLayout from './components/AppLayout.vue'

const theme = ref(darkTheme)
</script>

<style>
#app {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

body {
  margin: 0;
  padding: 0;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}
</style>
