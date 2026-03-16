<script setup lang="ts">
import { computed } from 'vue'

interface MenuItem {
  path: string
  name: string
  meta: {
    title: string
    icon?: string
    hidden?: boolean
  }
}

interface Props {
  item: MenuItem
  basePath: string
}

const props = defineProps<Props>()

const menuPath = computed(() => {
  return props.item.path
})
</script>

<template>
  <lay-menu-item :to="menuPath">
    <template #icon>
      <lay-icon v-if="props.item.meta.icon" :type="props.item.meta.icon" />
    </template>
    <span>{{ props.item.meta.title }}</span>
  </lay-menu-item>
</template>
