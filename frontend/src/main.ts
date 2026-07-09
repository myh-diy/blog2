import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useCustomTheme } from './composables/useCustomTheme'
import './style.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// Apply saved custom theme before mount
const { apply } = useCustomTheme()
apply()

app.mount('#app')
