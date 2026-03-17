<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import SidebarMenu from './components/SidebarMenu.vue'
import { generateMenuFromRoutes } from '@/utils/menu'

interface Props {
  collapsed?: boolean
  activeMenu?: string
}

const props = withDefaults(defineProps<Props>(), {
  collapsed: false,
  activeMenu: 'dashboard'
})

const router = useRouter()
const userStore = useUserStore()

// 从路由配置自动生成菜单
// 明确依赖 userStore.isAdmin，确保权限变化时重新计算
const menuList = computed(() => {
  // 这里明确依赖 isAdmin，当它变化时会触发重新计算
  const isAdmin = userStore.isAdmin
  console.log('🔄 computed 重新计算菜单，isAdmin:', isAdmin)
  return generateMenuFromRoutes()
})

// 监听 isAdmin 变化
watch(() => userStore.isAdmin, (newVal) => {
  console.log('👁️ isAdmin 变化:', newVal)
})

const emit = defineEmits<{
  menuClick: [menuId: string]
}>()

const handleMenuClick = (menuId: string) => {
  // 找到对应的路由并跳转
  const route = router.resolve({ name: menuId })
  if (route && route.path) {
    router.push(route.path)
  }
  emit('menuClick', menuId)
}
</script>

<template>
  <aside class="nav-vertical" :class="{ collapsed: props.collapsed }">
    <div class="nav-menu">
      <SidebarMenu
        :menu-list="menuList"
        :active-menu="props.activeMenu"
        @menu-click="handleMenuClick"
      />
    </div>
  </aside>
</template>

<style scoped>
.nav-vertical {
  width: 192px;
  height: calc(100vh - 48px);
  background-color: var(--erp-bg-card);
  border-right: 1px solid var(--erp-border-color);
  transition: width 0.3s;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.nav-vertical.collapsed {
  width: 64px;
}

.nav-menu {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>