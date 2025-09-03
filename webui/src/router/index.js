import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../components/Login.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/login', 
			name: 'Login',
			component: Login
		},
		{
			path: '/', 
			name: 'Home',
			component: HomeView,
			beforeEnter: (to, from, next) => {
				// Check if user is logged in
				const userId = localStorage.getItem('userId');
				if (!userId) {
					next('/login');
				} else {
					next();
				}
			}
		},
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
