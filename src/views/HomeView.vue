<template>
  <div class="home-container">
    <n-h1>{{ t('welcome') }}</n-h1>
    <n-p>FG-ABYSS - Advanced penetration testing framework</n-p>
    <n-button type="primary" @click="goToTerminal">
      {{ t('terminal') }}
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NH1, NP, NButton } from 'naive-ui'

const router = useRouter()
const { t } = useI18n()

function goToTerminal() {
  router.push('/terminal')
}
</script>

<style scoped>
.home-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 16px;
}
</style>
