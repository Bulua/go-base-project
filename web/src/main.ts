import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import { createPinia } from 'pinia'
import 'element-plus/dist/index.css'
import './assets/styles/main.css'
import App from './App.vue'
import router from './router'
import { vAuth } from './directives/vAuth'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.directive('auth', vAuth)
app.mount('#app')
