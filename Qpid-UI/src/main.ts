//import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { startUpdateChecker } from './update-checker'


const app = createApp(App)

app.use(router as any) 
app.mount('#app')
startUpdateChecker()
