import { App } from './App.vue.js'
import { router } from './router.js'

const app = Vue.createApp(App)
app.use(router)
app.mount('#app')
