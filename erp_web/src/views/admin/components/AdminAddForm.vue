<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { layer } from '@layui/layui-vue'
import { adminService } from '@/api/adminService'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'success': []
}>()

const loading = ref(false)
const formData = reactive({
  username: '',
  password: '',
  name: '',
  email: '',
  phone: ''
})

watch(() => props.visible, (val: boolean) => {
  if (val) {
    Object.assign(formData, { username: '', password: '', name: '', email: '', phone: '' })
  }
})

const handleSubmit = async () => {
  if (!formData.username) return layer.msg('请输入用户名', { icon: 0 })
  if (!formData.password) return layer.msg('请输入密码', { icon: 0 })
  if (!formData.name) return layer.msg('请输入姓名', { icon: 0 })

  loading.value = true
  try {
    await adminService.create({
      username: formData.username,
      password: formData.password,
      name: formData.name,
      email: formData.email || undefined,
      phone: formData.phone || undefined
    })
    layer.msg('创建成功', { icon: 1 })
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
    title="新增管理员"
    :area="['480px', 'auto']"
  >
    <div class="dialog-form">
      <lay-form :model="formData">
        <lay-form-item label="用户名" required>
          <lay-input v-model="formData.username" placeholder="请输入用户名（3-50位）" />
        </lay-form-item>
        <lay-form-item label="密码" required>
          <lay-input v-model="formData.password" type="password" placeholder="请输入密码（6-50位）" />
        </lay-form-item>
        <lay-form-item label="姓名" required>
          <lay-input v-model="formData.name" placeholder="请输入姓名" />
        </lay-form-item>
        <lay-form-item label="邮箱">
          <lay-input v-model="formData.email" placeholder="请输入邮箱" />
        </lay-form-item>
        <lay-form-item label="手机号">
          <lay-input v-model="formData.phone" placeholder="请输入手机号" />
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
