import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import App from './app.vue'
import router from './router'
import i18n from './i18n'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)
app.use(naive)

app.mount('#app')

if (import.meta.env.DEV) {
  console.log('FG-ABYSS started in development mode')
}
