import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import 'virtual:uno.css'
import '@unocss/reset/tailwind.css'
import './assets/globals.css'
import App from './App.vue'
import SuitesView from './pages/SuitesView.vue'
import SuiteDetailView from './pages/SuiteDetailView.vue'
import RunDetailView from './pages/RunDetailView.vue'
import HelpView from './pages/HelpView.vue'
import SettingsView from './pages/SettingsView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/suites' },
    { path: '/suites', component: SuitesView },
    { path: '/suites/:id', component: SuiteDetailView },
    { path: '/runs/:id', component: RunDetailView },
    { path: '/help', component: HelpView },
    { path: '/settings', component: SettingsView },
  ],
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
