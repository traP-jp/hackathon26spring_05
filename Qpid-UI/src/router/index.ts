import { createRouter, createWebHistory } from 'vue-router'

import Myself from '../components/Myself.vue'
import Likes from '../components/Likes.vue'
import Matching_Base from '@/components/Matching_Base.vue'


const routes = [
  { path: '/', component: Matching_Base },
  { path: '/me', component: Myself },
  { path: '/likes', component: Likes}

]

const router = createRouter({
  history: createWebHistory(),
  routes: routes as any
})

export default router