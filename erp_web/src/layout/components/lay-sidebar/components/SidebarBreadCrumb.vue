<script setup lang="ts">
interface BreadcrumbItem {
  title: string
  path: string
}

interface Props {
  items?: BreadcrumbItem[]
}

const props = withDefaults(defineProps<Props>(), {
  items: () => []
})
</script>

<template>
  <div class="sidebar-breadcrumb">
    <lay-breadcrumb v-if="props.items.length > 0">
      <lay-breadcrumb-item
        v-for="(item, index) in props.items"
        :key="item.path"
      >
        {{ item.title }}
      </lay-breadcrumb-item>
    </lay-breadcrumb>
    <span v-else class="breadcrumb-placeholder">首页</span>
  </div>
</template>

<style scoped>
.sidebar-breadcrumb {
  padding: 0 12px;
  background-color: transparent;
}

.breadcrumb-placeholder {
  font-size: 14px;
  color: #fff;
  font-weight: 500;
}

:deep(.layui-breadcrumb) {
  font-size: 14px;
}

:deep(.layui-breadcrumb-item) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.layui-breadcrumb-item:last-child) {
  color: #fff;
  font-weight: 500;
}

:deep(.layui-breadcrumb-separator) {
  color: rgba(255, 255, 255, 0.5);
}
</style>