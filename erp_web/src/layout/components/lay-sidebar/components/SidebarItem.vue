<script setup lang="ts">
import { computed } from 'vue'
import SidebarLinkItem from './SidebarLinkItem.vue'

interface MenuItem {
  path: string
  name: string
  meta: {
    title: string
    icon?: string
    hidden?: boolean
  }
  children?: MenuItem[]
}

interface Props {
  item: MenuItem
  basePath: string
}

const props = defineProps<Props>()

const hasChildren = computed(() => {
  const children = props.item.children
  return children && children.length > 0 && !props.item.meta?.hidden
})

const visibleChildren = computed(() => {
  if (!props.item.children) return []
  return props.item.children.filter(child => !child.meta?.hidden)
})
</script>

<template>
  <SidebarLinkItem v-if="!hasChildren" :item="props.item" :base-path="props.basePath" />

  <lay-sub-menu v-else :title="props.item.meta.title">
    <template #icon>
      <lay-icon v-if="props.item.meta.icon" :type="props.item.meta.icon" />
    </template>
    <SidebarItem
      v-for="child in visibleChildren"
      :key="child.path"
      :item="child"
      :base-path="child.path"
    />
  </lay-sub-menu>
</template>
