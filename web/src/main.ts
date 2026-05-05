import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import * as ElIcons from '@element-plus/icons-vue'
import { createPinia } from 'pinia'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './assets/styles/main.css'
import App from './App.vue'
import router from './router'
import { vAuth } from './directives/vAuth'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.directive('auth', vAuth)

for (const [name, component] of Object.entries(ElIcons)) {
  app.component(name, component)
}

app.mount('#app')
