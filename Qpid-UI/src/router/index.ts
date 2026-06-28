import { createRouter, createWebHistory } from 'vue-router'

import Myself from '../components/Myself.vue'
import Likes from '../components/Likes.vue'
import Matching_Base from '@/components/Matching_Base.vue'
import Login from '@/components/Login.vue'


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

router.beforeEach(async (to, from, next) => {
  if (to.path === '/login') {
    next();
    return;
  }

  try {
    const response = await fetch('https://qpid.trap.show/api/me');
    if (response.status === 401) {
      next('/login');
    } else {
      next();
    }
  } catch (error) {
    //next('/login');
    next();
  }
});

export default router