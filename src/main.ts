import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { i18n } from '@/i18n'
import router from '@/router'
import App from './App.vue'

// Global styles — order matters: variables first, then layout, then effects
import './styles/variables.css'
import './styles/app-shell.css'
import './styles/scrollbar.css'
import './styles/animations.css'

const app = createApp(App)
app.use(createPinia())
app.use(i18n)
app.use(router)
app.mount('#app')
