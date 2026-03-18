import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import i18n from './i18n'
import naive from 'naive-ui'
import './styles/global.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(i18n)
app.use(naive)
app.mount('#app')
