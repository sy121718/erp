<script setup lang="ts">
interface Props {
  title: string
  icon?: string
  collapsed?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  icon: '',
  collapsed: false
})
</script>

<template>
  <div class="sidebar-logo" :class="{ collapsed: props.collapsed }">
    <div class="logo-content">
      <lay-icon
        v-if="props.icon"
        :type="props.icon"
        class="logo-icon"
      />
      <span v-if="!props.collapsed" class="logo-title">{{ props.title }}</span>
    </div>
  </div>
</template>

<style scoped>
.sidebar-logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--erp-bg-card);
  border-bottom: 1px solid var(--erp-border-color);
  flex-shrink: 0;
}

.logo-content {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
}

.logo-icon {
  font-size: 24px;
  color: var(--erp-primary);
}

.logo-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--erp-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar-logo.collapsed .logo-content {
  justify-content: center;
  padding: 0 8px;
}
</style>