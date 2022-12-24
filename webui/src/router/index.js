import {createRouter, createWebHashHistory} from 'vue-router'
import StreamView from '../views/StreamView.vue'
import LoginView from '../views/LoginView.vue'
import UserView from '../views/UserView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{name: 'login', path: '/', component: LoginView},
		{name: 'stream', path: '/stream', component: StreamView},
		{name: 'user', path: '/users/:username', component: UserView},
	]
})

export default router
