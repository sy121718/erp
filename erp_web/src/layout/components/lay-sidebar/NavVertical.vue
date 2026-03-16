<script setup lang="ts">
import { computed } from 'vue'
import SidebarMenu from './components/SidebarMenu.vue'

interface MenuItem {
  id: string
  title: string
  icon?: string
  iconColor?: string
  badge?: string
  badgeType?: 'new' | 'hot' | 'ai'
  children?: MenuChild[]
  active?: boolean
}

interface MenuChild {
  id: string
  title: string
  badge?: string
  badgeType?: 'new' | 'hot' | 'ai'
  active?: boolean
}

interface Props {
  collapsed?: boolean
  activeMenu?: string
}

const props = withDefaults(defineProps<Props>(), {
  collapsed: false,
  activeMenu: 'collection-public'
})

const menuList = computed<MenuItem[]>(() => [
  {
    id: 'general',
    title: '通用功能',
    children: [
      { id: 'collection-public', title: '产品采集', active: true },
      { id: 'collection-box', title: '公用采集箱', active: props.activeMenu === 'collection-box' },
      { id: 'ai-workbench', title: 'AI工作台', badge: 'NEW', badgeType: 'new' },
      { id: 'infringement', title: '侵权检测' },
      { id: 'cargo-center', title: '货盘中心', badge: 'AI', badgeType: 'ai' },
      { id: 'cross-border-hot', title: '跨境热卖' },
      { id: 'taobao-collection', title: '淘宝采集' },
      { id: 'cross-border-industry', title: '跨境产业带' },
      { id: 'local-cargo', title: '本土货盘' },
      { id: '1688-same', title: '1688搜同款' },
      { id: '1688-smart', title: '1688智能采集' },
      { id: 'selected-cargo', title: '妙手精选货盘', badge: 'HOT', badgeType: 'hot' },
      { id: 'ai-selection', title: 'AI选品', badge: 'AI', badgeType: 'ai' },
      { id: '1688-selection', title: '1688选品库' },
      { id: 'invite-supplier', title: '邀请供应商' }
    ]
  },
  {
    id: 'data-selection',
    title: '数据选品',
    children: [
      { id: 'shopee-hot', title: 'Shopee热销' },
      { id: 'tiktok-hot', title: 'TikTok热销' },
      { id: 'temu-hot', title: 'Temu热销' },
      { id: 'lazada-hot', title: 'Lazada热销' },
      { id: 'amazon-hot', title: 'Amazon热销' }
    ]
  },
  {
    id: 'quick-upload',
    title: '快速上货',
    children: [
      { id: 'collection-box-2', title: '采集箱' },
      { id: 'publish-record', title: '发布记录' }
    ]
  },
  {
    id: 'wildberries',
    title: 'Wildberries',
    icon: 'W',
    iconColor: '#9333ea'
  },
  {
    id: 'online-products',
    title: '在线产品',
    children: [
      { id: 'product-manage', title: '产品管理' }
    ]
  }
])

const emit = defineEmits<{
  menuClick: [menuId: string]
}>()

const handleMenuClick = (menuId: string) => {
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