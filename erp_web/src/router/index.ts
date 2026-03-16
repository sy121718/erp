import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

// Admin 登录路径（不使用常用路径）
const ADMIN_LOGIN_PATH = '/ms-auth-admin'

const routes: RouteRecordRaw[] = [
  {
    path: ADMIN_LOGIN_PATH,
    name: 'AdminLogin',
    component: () => import('@/views/login-admin/index.vue'),
    meta: { title: '管理员登录', requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'collection',
        name: 'Collection',
        component: () => import('@/views/collection/index.vue'),
        meta: { title: '公用采集箱' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '页面不存在' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || '妙手ERP'} - 妙手ERP`

  const userStore = useUserStore()
  const isLoggedIn = userStore.isLoggedIn

  // 不需要登录的页面
  if (to.meta.requiresAuth === false) {
    next()
    return
  }

  // 需要登录但未登录，跳转到 404（不暴露登录页）
  if (!isLoggedIn && to.path !== ADMIN_LOGIN_PATH) {
    next({ name: 'NotFound' })
    return
  }

  // 已登录但访问登录页
  if (isLoggedIn && to.path === ADMIN_LOGIN_PATH) {
    next('/dashboard')
    return
  }

  next()
})

export default router
export { ADMIN_LOGIN_PATH }