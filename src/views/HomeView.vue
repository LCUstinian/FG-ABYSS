<template>
  <div class="home-view">
    <n-grid :cols="4" :x-gap="16" :y-gap="16">
      <n-grid-item>
        <n-card :title="$t('nav.payload')">
          <n-statistic :value="payloadStore.configCount">
            <template #prefix>
              <n-icon size="24">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                  <path fill="currentColor" d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" points="14 2 14 8 20 8"/>
                </svg>
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card :title="$t('nav.project')">
          <n-statistic :value="projectStore.projectCount">
            <template #prefix>
              <n-icon size="24">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                  <path fill="currentColor" d="M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8z"/>
                </svg>
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card title="WebShell">
          <n-statistic :value="webshellStore.onlineCount">
            <template #prefix>
              <n-icon size="24" color="#18a058">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                  <circle cx="12" cy="12" r="10" fill="currentColor"/>
                </svg>
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card :title="$t('nav.plugin')">
          <n-statistic :value="8">
            <template #prefix>
              <n-icon size="24">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                  <path fill="currentColor" d="M20.5 11H19V7c0-1.1-.9-2-2-2h-4V3.5C13 2.12 11.88 1 10.5 1S8 2.12 8 3.5V5H4c-1.1 0-2 .9-2 2v4H1.5C.12 11-1 12.12-1 13.5S.12 16 1.5 16H4v4c0 1.1.9 2 2 2h4v1.5c0 1.38 1.12 2.5 2.5 2.5s2.5-1.12 2.5-2.5V22h4c1.1 0 2-.9 2-2v-4h1.5c1.38 0 2.5-1.12 2.5-2.5s-1.12-2.5-2.5-2.5z"/>
                </svg>
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
    </n-grid>

    <n-grid :cols="1" :x-gap="16" :y-gap="16" style="margin-top: 16px;">
      <n-grid-item>
        <n-card title="快速操作">
          <n-space>
            <n-button type="primary">
              <template #icon>
                <n-icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/></svg></n-icon>
              </template>
              新建项目
            </n-button>
            <n-button type="success">
              <template #icon>
                <n-icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="currentColor" d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/></svg></n-icon>
              </template>
              创建载荷
            </n-button>
            <n-button type="info">
              <template #icon>
                <n-icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path fill="currentColor" d="M20.5 11H19V7c0-1.1-.9-2-2-2h-4V3.5C13 2.12 11.88 1 10.5 1S8 2.12 8 3.5V5H4c-1.1 0-2 .9-2 2v4H1.5C.12 11-1 12.12-1 13.5S.12 16 1.5 16H4v4c0 1.1.9 2 2 2h4v1.5c0 1.38 1.12 2.5 2.5 2.5s2.5-1.12 2.5-2.5V22h4c1.1 0 2-.9 2-2v-4h1.5c1.38 0 2.5-1.12 2.5-2.5s-1.12-2.5-2.5-2.5z"/></svg></n-icon>
              </template>
              安装插件
            </n-button>
          </n-space>
        </n-card>
      </n-grid-item>
    </n-grid>

    <n-grid :cols="1" :x-gap="16" :y-gap="16" style="margin-top: 16px;">
      <n-grid-item>
        <n-card title="最近活动">
          <n-timeline>
            <n-timeline-item type="success" title="创建了载荷" time="10:30" content="test.php" />
            <n-timeline-item type="info" title="启用了插件" time="10:25" content="Lua Runner" />
            <n-timeline-item type="warning" title="创建了项目" time="10:20" content="Project A" />
          </n-timeline>
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { usePayloadStore } from '@/stores/payload'
import { useProjectStore } from '@/stores/project'
import { useWebShellStore } from '@/stores/webshell'
import { useI18n } from 'vue-i18n'

const payloadStore = usePayloadStore()
const projectStore = useProjectStore()
const webshellStore = useWebShellStore()
const { t } = useI18n()
</script>

<style scoped>
.home-view {
  padding: 16px;
}
</style>
