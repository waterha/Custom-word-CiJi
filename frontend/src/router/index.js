import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/learn',
    name: 'Learn',
    component: () => import('../views/Learn.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/stats',
    name: 'LearningStats',
    component: () => import('../views/LearningStats.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/user-admin',
    name: 'UserAdmin',
    component: () => import('../views/UserAdmin.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/Admin.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true
    }
  },
  {
    path: '/monitor',
    name: 'Monitor',
    component: () => import('../views/Monitor.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const user = JSON.parse(localStorage.getItem('user'))

  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!token) {
      next({ path: '/login' })
    } else if (to.matched.some(record => record.meta.requiresAdmin)) {
      if (user && user.role === 'admin') {
        next()
      } else {
        next({ path: '/learn' })
      }
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
