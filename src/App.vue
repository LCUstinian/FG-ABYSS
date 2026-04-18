<template>
  <n-config-provider
    :theme="themeStore.naiveTheme"
    :theme-overrides="themeStore.themeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme-overrides-common="{ zIndexModal: 100, zIndexPopover: 20, zIndexMessage: 1000 }"
  >
    <n-message-provider placement="top-right">
      <n-dialog-provider>
        <AppShell />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { zhCN, dateZhCN } from 'naive-ui'
import { useThemeStore } from '@/stores/theme'
import { useThemeCssVars } from '@/composables/useThemeCssVars'
import AppShell from '@/components/layout/AppShell.vue'

const themeStore = useThemeStore()
useThemeCssVars()
</script>
