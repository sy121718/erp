<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'
import NavVertical from './components/lay-sidebar/NavVertical.vue'
import LayNavbar from './components/lay-navbar/index.vue'
import LayContent from './components/lay-content/index.vue'

const router = useRouter()
const userStore = useUserStore()

const collapsed = ref(false)

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
          await userStore.logout()
          router.push('/login')
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
      <NavVertical active-menu="collection-box" />
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