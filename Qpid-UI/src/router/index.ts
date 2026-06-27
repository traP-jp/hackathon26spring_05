import { createRouter, createWebHistory } from 'vue-router'

import Matching from '../components/Matching.vue'
import Myself from '../components/Myself.vue'
import Likes from '../components/Likes.vue'


const routes = [
  { path: '/', component: Matching },
  { path: '/me', component: Myself },
  { path: '/likes', component: Likes}

]

const router = createRouter({
  history: createWebHistory(),
  routes: routes as any
})

export default router