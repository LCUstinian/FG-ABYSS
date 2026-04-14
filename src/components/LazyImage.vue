<template>
  <div class="lazy-image" :style="{ height }">
    <div v-if="loading" class="lazy-image-placeholder">
      <Loading :text="loadingText" />
    </div>
    <img
      v-show="loaded"
      ref="imgRef"
      :src="computedSrc"
      :alt="alt"
      :style="{ objectFit }"
      @load="handleLoad"
      @error="handleError"
    />
    <div v-if="error && showRetry" class="lazy-image-error">
      <n-button size="small" @click="retry">
        <template #icon>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23 4 23 10 17 10"/>
            <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
          </svg>
        </template>
        重试
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { NButton } from 'naive-ui'
import Loading from './Loading.vue'

interface Props {
  src: string
  alt?: string
  height?: string
  objectFit?: 'contain' | 'cover' | 'fill' | 'none' | 'scale-down'
  loadingText?: string
  showRetry?: boolean
  threshold?: number
}

const props = withDefaults(defineProps<Props>(), {
  alt: '',
  height: 'auto',
  objectFit: 'cover',
  loadingText: '加载中...',
  showRetry: true,
  threshold: 0,
})

const emit = defineEmits<{
  load: []
  error: [error: Event]
}>()

const imgRef = ref<HTMLImageElement | null>(null)
const loading = ref(true)
const loaded = ref(false)
const error = ref(false)
const observer = ref<IntersectionObserver | null>(null)

const computedSrc = computed(() => {
  return error.value ? '' : props.src
})

const handleLoad = () => {
  loading.value = false
  loaded.value = true
  error.value = false
  emit('load')
}

const handleError = (e: Event) => {
  loading.value = false
  error.value = true
  emit('error', e)
}

const retry = () => {
  loading.value = true
  loaded.value = false
  error.value = false
  if (imgRef.value) {
    imgRef.value.src = props.src
  }
}

const observeImage = () => {
  if (!imgRef.value) return

  observer.value = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          // 图片进入视口，开始加载
          if (imgRef.value) {
            imgRef.value.src = props.src
          }
          // 停止观察
          observer.value?.disconnect()
        }
      })
    },
    {
      rootMargin: `${props.threshold}px`,
    }
  )

  observer.value.observe(imgRef.value)
}

onMounted(() => {
  observeImage()
})

onBeforeUnmount(() => {
  observer.value?.disconnect()
})
</script>

<style scoped>
.lazy-image {
  position: relative;
  display: inline-block;
  width: 100%;
  background-color: var(--n-border-color);
  border-radius: 4px;
  overflow: hidden;
}

.lazy-image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  min-height: 100px;
}

.lazy-image-error {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
}
</style>
