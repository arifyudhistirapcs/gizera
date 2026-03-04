import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Vant, { Locale } from 'vant'
import idID from 'vant/es/locale/lang/id-ID'
import App from './App.vue'
import router from './router'
import 'vant/lib/index.css'
import './styles/horizon-mobile.css'

// Set Vant locale to Indonesian
Locale.use('id-ID', idID)

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Vant)

app.mount('#app')

// Register service worker
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js')
  })
}
