<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { layer } from '@layui/layui-vue'
import { userService } from '@/api/userService'

interface Props {
  visible: boolean
  userId?: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'success': []
}>()

const loading = ref(false)
const formData = reactive({
  id: 0,
  username: '',
  name: '',
  email: '',
  phone: '',
  status: 1,
  points: 0
})

const loadDetail = async (id: number) => {
  loading.value = true
  try {
    const detail = await userService.adminGetDetail(id)
    Object.assign(formData, {
      id: detail.id,
      username: detail.username,
      name: detail.name,
      email: detail.email || '',
      phone: detail.phone || '',
      status: detail.status,
      points: detail.points
    })
  } catch {
    emit('update:visible', false)
  } finally {
    loading.value = false
  }
}

watch(() => props.visible, (val: boolean) => {
  if (val && props.userId) {
    loadDetail(props.userId)
  }
})

const handleSubmit = async () => {
  loading.value = true
  try {
    await userService.adminUpdate(formData.id, {
      name: formData.name || undefined,
      email: formData.email || undefined,
      phone: formData.phone || undefined
    })
    layer.msg('更新成功', { icon: 1 })
    emit('update:visible', false)
    emit('success')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <lay-layer
    :modelValue="visible"
    @update:modelValue="(val: boolean) => emit('update:visible', val)"
    title="编辑用户"
    :area="['480px', 'auto']"
  >
    <div class="dialog-form">
      <lay-form :model="formData">
        <lay-form-item label="用户名">
          <lay-input v-model="formData.username" disabled />
        </lay-form-item>
        <lay-form-item label="姓名">
          <lay-input v-model="formData.name" placeholder="请输入姓名" />
        </lay-form-item>
        <lay-form-item label="邮箱">
          <lay-input v-model="formData.email" placeholder="请输入邮箱" />
        </lay-form-item>
        <lay-form-item label="手机号">
          <lay-input v-model="formData.phone" placeholder="请输入手机号" />
        </lay-form-item>
        <lay-form-item label="状态">
          <lay-tag :type="formData.status === 1 ? 'normal' : 'danger'">
            {{ formData.status === 1 ? '正常' : '禁用' }}
          </lay-tag>
        </lay-form-item>
        <lay-form-item label="积分">
          <span>{{ formData.points }}</span>
        </lay-form-item>
      </lay-form>
      <div class="dialog-form__footer">
        <lay-button @click="emit('update:visible', false)">取消</lay-button>
        <lay-button type="primary" :loading="loading" @click="handleSubmit">确定</lay-button>
      </div>
    </div>
  </lay-layer>
</template>

<style scoped>
.dialog-form {
  padding: 20px 24px 0;
}
.dialog-form__footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 0;
  border-top: 1px solid var(--erp-border-color);
  margin-top: 8px;
}
</style>
