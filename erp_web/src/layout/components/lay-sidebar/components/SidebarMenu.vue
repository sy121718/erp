<script setup lang="ts">
import { ref } from 'vue'

interface MenuChild {
  id: string
  title: string
  badge?: string
  badgeType?: 'new' | 'hot' | 'ai'
  active?: boolean
}

interface MenuItem {
  id: string
  title: string
  icon?: string
  iconColor?: string
  badge?: string
  badgeType?: 'new' | 'hot' | 'ai'
  children?: MenuChild[]
  expanded?: boolean
}

interface Props {
  menuList: MenuItem[]
  activeMenu?: string
}

const props = defineProps<Props>()

const expandedMenus = ref<Set<string>>(new Set(['general']))

const isExpanded = (menuId: string) => {
  return expandedMenus.value.has(menuId)
}

const toggleExpand = (menuId: string) => {
  if (expandedMenus.value.has(menuId)) {
    expandedMenus.value.delete(menuId)
  } else {
    expandedMenus.value.add(menuId)
  }
}

const emit = defineEmits<{
  menuClick: [menuId: string]
}>()

const handleMenuClick = (menuId: string) => {
  emit('menuClick', menuId)
}

const getBadgeClass = (type?: 'new' | 'hot' | 'ai') => {
  switch (type) {
    case 'new':
      return 'badge-new'
    case 'hot':
      return 'badge-hot'
    case 'ai':
      return 'badge-ai'
    default:
      return ''
  }
}
</script>

<template>
  <div class="sidebar-menu">
    <!-- 菜单项 -->
    <div
      v-for="menu in menuList"
      :key="menu.id"
      class="menu-item"
    >
      <!-- 有子菜单的项 -->
      <template v-if="menu.children && menu.children.length > 0">
        <!-- 父菜单标题 -->
        <div
          class="menu-title"
          :class="{ expanded: isExpanded(menu.id) }"
          @click="toggleExpand(menu.id)"
        >
          <span class="title-text">{{ menu.title }}</span>
          <lay-icon
            type="layui-icon-down"
            class="expand-icon"
            :class="{ expanded: isExpanded(menu.id) }"
          />
        </div>
        
        <!-- 子菜单 -->
        <div v-show="isExpanded(menu.id)" class="sub-menu">
          <!-- 两列布局的子菜单 -->
          <div class="sub-menu-grid">
            <div
              v-for="child in menu.children"
              :key="child.id"
              class="sub-menu-item"
              :class="{ active: child.active || props.activeMenu === child.id }"
              @click="handleMenuClick(child.id)"
            >
              <span class="item-text">{{ child.title }}</span>
              <span
                v-if="child.badge"
                class="item-badge"
                :class="getBadgeClass(child.badgeType)"
              >
                {{ child.badge }}
              </span>
            </div>
          </div>
        </div>
      </template>
      
      <!-- 无子菜单的项 (如 Wildberries) -->
      <template v-else>
        <div
          class="menu-single"
          :class="{ active: props.activeMenu === menu.id }"
          @click="handleMenuClick(menu.id)"
        >
          <span
            v-if="menu.icon"
            class="menu-icon"
            :style="{ color: menu.iconColor }"
          >
            {{ menu.icon }}
          </span>
          <span class="title-text">{{ menu.title }}</span>
          <span
            v-if="menu.badge"
            class="item-badge"
            :class="getBadgeClass(menu.badgeType)"
          >
            {{ menu.badge }}
          </span>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.sidebar-menu {
  padding: 8px 0;
}

/* 父菜单标题 */
.menu-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  font-size: 13px;
  font-weight: 500;
  color: var(--erp-text-primary);
  cursor: pointer;
  transition: background-color 0.3s;
}

.menu-title:hover {
  background-color: var(--erp-bg-base);
}

.expand-icon {
  font-size: 12px;
  color: var(--erp-text-tertiary);
  transition: transform 0.3s;
}

.expand-icon.expanded {
  transform: rotate(180deg);
}

/* 子菜单 */
.sub-menu {
  padding: 4px 0;
}

.sub-menu-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2px;
  padding: 0 8px;
}

.sub-menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  font-size: 12px;
  color: var(--erp-text-secondary);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
  overflow: hidden;
}

.sub-menu-item:hover {
  background-color: var(--erp-bg-base);
  color: var(--erp-primary);
}

.sub-menu-item.active {
  background-color: var(--erp-primary-light);
  color: var(--erp-primary);
}

.item-text {
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 单菜单项 */
.menu-single {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  font-size: 13px;
  color: var(--erp-text-primary);
  cursor: pointer;
  transition: all 0.3s;
}

.menu-single:hover {
  background-color: var(--erp-bg-base);
  color: var(--erp-primary);
}

.menu-single.active {
  background-color: var(--erp-primary-light);
  color: var(--erp-primary);
  border-left: 3px solid var(--erp-primary);
}

.menu-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  border-radius: 4px;
  background-color: rgba(147, 51, 234, 0.1);
}

/* 标签 */
.item-badge {
  font-size: 10px;
  padding: 1px 4px;
  border-radius: 3px;
  font-weight: 500;
  flex-shrink: 0;
}

.badge-new {
  background-color: var(--erp-red);
  color: #fff;
}

.badge-hot {
  background-color: var(--erp-red);
  color: #fff;
}

.badge-ai {
  background-color: var(--erp-blue);
  color: #fff;
}
</style>