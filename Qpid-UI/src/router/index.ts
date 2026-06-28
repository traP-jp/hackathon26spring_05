import { createRouter, createWebHistory } from 'vue-router'

import Myself from '../components/Myself.vue'
import Likes from '../components/Likes.vue'
import Matching_Base from '../components/Matching_Base.vue'
import Login from '../components/Login.vue'
import { apiFetch } from '../api'


const routes = [
  { path: '/', component: Matching_Base },
  { path: '/me', component: Myself },
  { path: '/likes', component: Likes},
  { path: '/login', component: Login}

]

const router = createRouter({
  history: createWebHistory(),
  routes: routes as any
})

router.beforeEach(async (to) => {
  if (to.path === '/login') {
    return true
  }

  try {
    const response = await apiFetch('/api/me')
    if (response.status === 401) {
      return {
        path: '/login',
        query: { redirect: to.fullPath },
      }
    }
  } catch (error) {
    console.log('Authentication Error! :', error)
  }

  return true
})

export default router
