import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const ADMIN_LOGIN_PATH = '/ms-auth-admin'
const USER_LOGIN_PATH = '/login'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    hidden?: boolean
    icon?: string
    parent?: string
    order?: number
    requiresAdmin?: boolean
    userTypes?: ('admin' | 'user')[]
  }
}

const routes: RouteRecordRaw[] = [
  {
    path: ADMIN_LOGIN_PATH,
    name: 'AdminLogin',
    component: () => import('@/views/login-admin/index.vue'),
    meta: { title: '管理员登录', requiresAuth: false, hidden: true }
  },
  {
    path: USER_LOGIN_PATH,
    name: 'UserLogin',
    component: () => import('@/views/login-user/index.vue'),
    meta: { title: '用户登录', requiresAuth: false, hidden: true }
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
        meta: { title: '首页', icon: 'layui-icon-home', order: 0, hidden: true }
      },
      {
        path: 'admin',
        name: 'AdminList',
        component: () => import('@/views/admin/index.vue'),
        meta: { title: '管理员管理', parent: 'system', order: 1, requiresAdmin: true, userTypes: ['admin'] }
      },
      {
        path: 'user-manage',
        name: 'UserManage',
        component: () => import('@/views/user/index.vue'),
        meta: { title: '用户管理', parent: 'system', order: 2, userTypes: ['admin'] }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile-admin/index.vue'),
        meta: { title: '个人中心', hidden: true, userTypes: ['admin'] }
      },
      {
        path: 'collection',
        name: 'Collection',
        component: () => import('@/views/collection/index.vue'),
        meta: { title: '公用采集箱', hidden: true }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '页面不存在', hidden: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = `${to.meta.title || '妙手ERP'} - 妙手ERP`

  const userStore = useUserStore()
  const isLoggedIn = userStore.isLoggedIn
  const userType = userStore.userType

  if (to.meta.requiresAuth === false) {
    if (isLoggedIn && (to.name === 'AdminLogin' || to.name === 'UserLogin')) {
      next('/dashboard')
      return
    }
    next()
    return
  }

  if (!isLoggedIn) {
    next({ name: 'UserLogin' })
    return
  }

  if (to.meta.userTypes && userType && !to.meta.userTypes.includes(userType as 'admin' | 'user')) {
    next({ name: 'NotFound' })
    return
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next({ name: 'NotFound' })
    return
  }

  next()
})

export default router
export { ADMIN_LOGIN_PATH, USER_LOGIN_PATH }
