<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'

interface UserInfo {
  nickname: string
  avatar?: string
}

interface NavItem {
  title: string
  path: string
  badge?: number
}

interface Props {
  collapsed?: boolean
  userInfo?: UserInfo
  activeNav?: string
}

const props = withDefaults(defineProps<Props>(), {
  collapsed: false,
  userInfo: () => ({
    nickname: '管理员',
    avatar: ''
  }),
  activeNav: '首页'
})

const emit = defineEmits<{
  toggleCollapse: []
  navClick: [nav: string]
}>()

const router = useRouter()
const userStore = useUserStore()

const mainNavs = ref<NavItem[]>([
  { title: '首页', path: '/dashboard' },
  { title: '产品', path: '/product' },
  { title: '订单', path: '/order', badge: 3 },
  { title: '托管', path: '/hosting' },
  { title: '广告', path: '/ads' },
  { title: '数据', path: '/data' },
  { title: '采购', path: '/purchase' },
  { title: '仓库', path: '/warehouse' },
  { title: '物流', path: '/logistics' },
  { title: '授权', path: '/auth' },
  { title: '更多', path: '/more' }
])

const rightLinks = ref([
  { title: '联系客服', path: '/support' },
  { title: '帮助', path: '/help' },
  { title: '插件下载', path: '/plugins' },
  { title: '同步', path: '/sync' }
])

const handleToggle = () => {
  emit('toggleCollapse')
}

const handleNavClick = (nav: NavItem) => {
  emit('navClick', nav.title)
}

// 跳转到个人中心
const handleProfile = () => {
  router.push('/profile')
}

// 退出登录
const handleLogout = async () => {
  // 先清除本地认证信息
  userStore.clearAuth()
  // 跳转登录页
  router.push('/ms-auth-admin')
  layer.msg('退出成功', { icon: 1, time: 1000 })

  // 尝试调用后端退出接口（可选，失败也不影响）
  try {
    await userStore.logout()
  } catch (error) {
    // 忽略错误，错误由 request.ts 统一处理
  }
}
</script>

<template>
  <header class="lay-navbar">
    <!-- 左侧区域 -->
    <div class="navbar-left">
      <!-- 折叠按钮 -->
      <div class="collapse-btn" @click="handleToggle">
        <lay-icon type="layui-icon-spread-left" class="collapse-icon" />
      </div>
      
      <!-- Logo -->
      <div class="logo">
        <lay-icon type="layui-icon-face-smile" class="logo-icon" />
        <span class="logo-text">妙手ERP</span>
      </div>
      
      <!-- 主导航 -->
      <nav class="main-nav">
        <a
          v-for="nav in mainNavs"
          :key="nav.path"
          class="nav-link"
          :class="{ active: props.activeNav === nav.title }"
          @click="handleNavClick(nav)"
        >
          {{ nav.title }}
          <span v-if="nav.badge" class="nav-badge">{{ nav.badge }}</span>
        </a>
      </nav>
    </div>
    
    <!-- 右侧区域 -->
    <div class="navbar-right">
      <!-- 文字链接 -->
      <div class="right-links">
        <a v-for="link in rightLinks" :key="link.path" class="right-link">
          {{ link.title }}
        </a>
      </div>
      
      <!-- 订购按钮 -->
      <lay-button size="sm" class="order-btn">订购</lay-button>
      
      <!-- 用户头像 -->
      <lay-dropdown>
        <div class="user-info">
          <lay-avatar
            :src="props.userInfo.avatar || '/image/avatar.png'"
            :text="props.userInfo.nickname?.charAt(0)"
            size="sm"
            class="user-avatar"
          />
          <lay-icon type="layui-icon-down" class="dropdown-icon" />
        </div>
        <template #content>
          <lay-dropdown-menu>
            <lay-dropdown-menu-item @click="handleProfile">
              <lay-icon type="layui-icon-username" />
              个人中心
            </lay-dropdown-menu-item>
            <lay-dropdown-menu-item>
              <lay-icon type="layui-icon-set" />
              系统设置
            </lay-dropdown-menu-item>
            <lay-dropdown-menu-item divided @click="handleLogout">
              <lay-icon type="layui-icon-logout" />
              退出登录
            </lay-dropdown-menu-item>
          </lay-dropdown-menu>
        </template>
      </lay-dropdown>
    </div>
  </header>
</template>

<style scoped>
.lay-navbar {
  height: 48px;
  background-color: var(--erp-primary);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  position: sticky;
  top: 0;
  z-index: 100;
}

/* 左侧区域 */
.navbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
}

.collapse-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.collapse-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.collapse-icon {
  font-size: 16px;
  color: #fff;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon {
  font-size: 24px;
  color: #fff;
}

.logo-text {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

.main-nav {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: 16px;
}

.nav-link {
  position: relative;
  padding: 6px 12px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.nav-link:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.nav-link.active {
  background-color: rgba(255, 255, 255, 0.15);
  color: #fff;
  font-weight: 500;
}

.nav-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  font-size: 10px;
  line-height: 16px;
  text-align: center;
  background-color: var(--erp-red);
  color: #fff;
  border-radius: 8px;
}

/* 右侧区域 */
.navbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.right-links {
  display: flex;
  align-items: center;
  gap: 16px;
}

.right-link {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.85);
  text-decoration: none;
  cursor: pointer;
  transition: color 0.3s;
}

.right-link:hover {
  color: #fff;
}

.order-btn {
  background-color: var(--erp-orange) !important;
  border-color: var(--erp-orange) !important;
  font-weight: 500;
}

.order-btn:hover {
  background-color: var(--erp-orange-hover) !important;
  border-color: var(--erp-orange-hover) !important;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.user-avatar {
  flex-shrink: 0;
}

.dropdown-icon {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}
</style>