<template>
  <div class="payload-view">
    <n-space vertical>
      <n-h1>{{ $t('payload.manager') }}</n-h1>
      <n-empty 
        v-if="payloadStore.configs.length === 0"
        description="暂无载荷配置，点击'创建载荷配置'按钮创建"
      >
        <template #extra>
          <n-button type="primary">
            {{ $t('payload.create') }}
          </n-button>
        </template>
      </n-empty>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { usePayloadStore } from '@/stores/payload'

const payloadStore = usePayloadStore()
</script>

<style scoped>
.payload-view {
  padding: 16px;
}
</style>
