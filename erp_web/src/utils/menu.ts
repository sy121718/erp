import router from '@/router'
import { useUserStore } from '@/stores/user'
import type { RouteRecordRaw } from 'vue-router'

export interface MenuItem {
  id: string
  title: string
  icon?: string
  iconColor?: string
  badge?: string
  badgeType?: 'new' | 'hot' | 'ai'
  children?: MenuItem[]
  path?: string
  order?: number
}

/**
 * 从路由配置生成菜单
 */
export function generateMenuFromRoutes(): MenuItem[] {
  const routes = router.options.routes
  const userStore = useUserStore()
  const menuMap = new Map<string, MenuItem>()
  const menuList: MenuItem[] = []

  console.log('🚀 开始生成菜单')
  console.log('👤 当前用户是否超管:', userStore.isAdmin)
  console.log('📦 用户信息:', userStore.adminInfo)

  // 遍历所有路由
  function processRoute(route: RouteRecordRaw, parentPath: string = '') {
    // 跳过隐藏的路由
    if (route.meta?.hidden) {
      console.log('⏭️ 跳过隐藏路由:', route.path)
      // 即使父路由隐藏，也要处理子路由
      if (route.children && route.children.length > 0) {
        route.children.forEach(child => {
          processRoute(child, route.path.startsWith('/') ? route.path : parentPath)
        })
      }
      return
    }

    // 跳过需要超管权限但当前用户不是超管的路由
    if (route.meta?.requiresAdmin && !userStore.isAdmin) {
      console.log('🔒 跳过需要超管权限的路由:', route.path)
      return
    }

    // 跳过 userTypes 不匹配的路由
    if (route.meta?.userTypes && userStore.userType) {
      if (!route.meta.userTypes.includes(userStore.userType as 'admin' | 'user')) {
        console.log('🔒 跳过用户类型不匹配的路由:', route.path, 'need:', route.meta.userTypes, 'current:', userStore.userType)
        return
      }
    }

    // 只处理有 title 的路由
    if (!route.meta?.title) {
      console.log('⏭️ 跳过无标题路由:', route.path)
      // 即使父路由无标题，也要处理子路由
      if (route.children && route.children.length > 0) {
        route.children.forEach(child => {
          processRoute(child, route.path.startsWith('/') ? route.path : parentPath)
        })
      }
      return
    }

    console.log('✅ 处理路由:', route.path, route.meta.title)

    const fullPath = route.path.startsWith('/') ? route.path : `${parentPath}/${route.path}`

    const menuItem: MenuItem = {
      id: route.name as string || route.path,
      title: route.meta.title,
      icon: route.meta.icon,
      path: fullPath,
      order: route.meta.order || 999
    }

    // 如果有父菜单ID，添加到父菜单下
    if (route.meta.parent) {
      const parentId = route.meta.parent
      if (!menuMap.has(parentId)) {
        // 创建父菜单
        menuMap.set(parentId, {
          id: parentId,
          title: getParentTitle(parentId),
          children: [],
          order: getParentOrder(parentId)
        })
        console.log('📁 创建父菜单:', parentId)
      }
      const parentMenu = menuMap.get(parentId)!
      if (!parentMenu.children) {
        parentMenu.children = []
      }
      parentMenu.children.push(menuItem)
      console.log('➕ 添加子菜单:', menuItem.title, '到', parentId)
    } else {
      // 顶级菜单
      menuMap.set(menuItem.id, menuItem)
      console.log('📌 添加顶级菜单:', menuItem.title)
    }

    // 处理子路由
    if (route.children && route.children.length > 0) {
      route.children.forEach(child => {
        processRoute(child, fullPath)
      })
    }
  }

  // 处理所有路由
  routes.forEach(route => processRoute(route))

  // 转换为数组并排序
  menuMap.forEach(menu => {
    // 如果是分组菜单，对其子菜单排序
    if (menu.children && menu.children.length > 0) {
      menu.children.sort((a, b) => (a.order || 999) - (b.order || 999))
    }
    menuList.push(menu)
  })

  // 顶级菜单排序
  menuList.sort((a, b) => (a.order || 999) - (b.order || 999))

  console.log('🎯 生成的菜单列表:', menuList)

  return menuList
}

/**
 * 获取父菜单标题
 */
function getParentTitle(parentId: string): string {
  const titleMap: Record<string, string> = {
    'system': '系统管理',
    'data': '数据管理',
    'product': '产品管理',
    'order': '订单管理'
  }
  return titleMap[parentId] || parentId
}

/**
 * 获取父菜单排序
 */
function getParentOrder(parentId: string): number {
  const orderMap: Record<string, number> = {
    'system': 1,
    'data': 2,
    'product': 3,
    'order': 4
  }
  return orderMap[parentId] || 999
}
