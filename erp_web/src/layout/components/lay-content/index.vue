<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  padding?: number | string
  background?: string
  showPadding?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  padding: 16,
  background: '#f3f4f6',
  showPadding: true
})

const contentStyle = computed(() => {
  const style: Record<string, string> = {
    backgroundColor: props.background
  }

  if (props.showPadding) {
    const paddingValue = typeof props.padding === 'number' ? `${props.padding}px` : props.padding
    style.padding = paddingValue
  }

  return style
})
</script>

<template>
  <main class="lay-content">
    <div class="content-wrapper" :style="contentStyle">
      <slot />
    </div>
  </main>
</template>

<style scoped>
.lay-content {
  flex: 1;
  min-height: calc(100vh - 48px);
  overflow: auto;
}

.content-wrapper {
  min-height: 100%;
  box-sizing: border-box;
}
</style>