<template>
  <n-config-provider
    :theme="themeStore.naiveTheme"
    :theme-overrides="themeStore.themeOverrides"
    :locale="naiveLocale"
    :date-locale="naiveDateLocale"
  >
    <n-message-provider placement="top-right">
      <n-dialog-provider>
        <AppShell />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { zhCN, enUS, dateZhCN, dateEnUS } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { useThemeStore } from '@/stores/theme'
import { useThemeCssVars } from '@/composables/useThemeCssVars'
import AppShell from '@/components/layout/AppShell.vue'

const { locale } = useI18n()
const themeStore = useThemeStore()
useThemeCssVars()

const naiveLocale     = computed(() => locale.value === 'zh-CN' ? zhCN : enUS)
const naiveDateLocale = computed(() => locale.value === 'zh-CN' ? dateZhCN : dateEnUS)
</script>
