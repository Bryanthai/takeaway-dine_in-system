import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router' // Import your router instance

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router) // Use Vue Router
app.mount('#app')