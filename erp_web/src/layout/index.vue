<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'
import NavVertical from './components/lay-sidebar/NavVertical.vue'
import LayNavbar from './components/lay-navbar/index.vue'
import LayContent from './components/lay-content/index.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const collapsed = ref(false)

// 初始化用户信息
onMounted(async () => {
  await userStore.initUserInfo()
})

// 当前激活的菜单（从路由名称获取）
const activeMenu = computed(() => route.name as string || 'Dashboard')

const userInfo = computed(() => ({
  nickname: userStore.nickname,
  avatar: userStore.avatar
}))

const handleToggleCollapse = () => {
  collapsed.value = !collapsed.value
}

const handleLogout = async () => {
  layer.confirm('确定要退出登录吗？', {
    btn: [
      {
        text: '确定',
        callback: async (id: string | number) => {
          layer.close(id)
          const loginPath = userStore.getLoginPath()
          await userStore.logout()
          router.push(loginPath)
        }
      },
      {
        text: '取消',
        callback: (id: string | number) => {
          layer.close(id)
        }
      }
    ]
  })
}

const handleRefresh = () => {
  window.location.reload()
}

// 菜单点击事件
const handleMenuClick = (menuId: string) => {
  // 菜单点击会由 NavVertical 内部处理路由跳转
  // 这里可以添加额外的逻辑（如关闭折叠菜单等）
}
</script>

<template>
  <div class="app-container">
    <!-- 顶部导航 -->
    <LayNavbar
      :collapsed="collapsed"
      :user-info="userInfo"
      active-nav="首页"
      @toggle-collapse="handleToggleCollapse"
    />
    <div class="app-body">
      <!-- 左侧边栏 -->
      <NavVertical :active-menu="activeMenu" @menu-click="handleMenuClick" />
      <!-- 主内容区 -->
      <LayContent>
        <router-view />
      </LayContent>
    </div>
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--erp-bg-base);
}

.app-body {
  display: flex;
  flex: 1;
}
</style>