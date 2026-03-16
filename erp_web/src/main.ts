import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import pinia from './stores'
import LayuiVue from '@layui/layui-vue'
import '@layui/layui-vue/lib/index.css'
import './style.css'

const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(LayuiVue)

app.mount('#app')