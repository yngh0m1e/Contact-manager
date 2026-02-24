import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Contacts from '../views/Contacts.vue'
import Profile from '../views/Profile.vue'

const routes = [
  { path: '/', component: Login },
  { path: '/contacts', component: Contacts, meta: { requiresAuth: true } },
  { path: '/profile', component: Profile, meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !localStorage.getItem('token')) {
    next('/')
  } else {
    next()
  }
})

export default router