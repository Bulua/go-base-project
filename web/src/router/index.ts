import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/store/modules/auth'

const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/login/index.vue'),
    meta: { public: true, title: '登录' },
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layouts/default/index.vue'),
    redirect: '/dashboard/workbench',
    children: [],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { public: true, title: '页面不存在' },
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes,
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  if (!auth.isAuthenticated) {
    if (to.name === 'Login') return true
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.path === '/login') {
    return { path: '/' }
  }

  if (auth.menuRoutes.length === 0) {
    await auth.loadWorkspace()
    return { path: to.fullPath, replace: true }
  }

  return true
})

export default router
