// WASA Text Application - Version 2.2 - Cache Buster Edition
// Build timestamp: 2024-10-01-12:17:00-FORCE-REBUILD
import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'

import './assets/dashboard.css'
import './assets/main.css'

// Debug info for cache busting
console.log('WASA App initializing - Build 2.2');

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);

// Force new hash generation
const CACHE_BUSTER = 'v2.2-' + Date.now();
app.config.globalProperties.$cacheBuster = CACHE_BUSTER;
app.use(router)
app.mount('#app')
