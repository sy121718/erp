import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'
import type { LoginForm } from './types'

/** 用户登录页面组合式函数 */
export const useLogin = () => {
  const router = useRouter()
  const userStore = useUserStore()

  const formData = ref<LoginForm>({
    username: '',
    password: ''
  })

  const loading = ref(false)

  const canSubmit = computed<boolean>(() => {
    return !!(formData.value.username && formData.value.password)
  })

  const handleLogin = async (): Promise<void> => {
    if (!canSubmit.value) {
      layer.msg('请填写完整信息', { icon: 0 })
      return
    }

    loading.value = true
    try {
      await userStore.userLogin({
        username: formData.value.username,
        password: formData.value.password
      })
      router.push('/')
    } catch {
      // 错误由 request.ts 统一处理
    } finally {
      loading.value = false
    }
  }

  return {
    formData,
    loading,
    canSubmit,
    handleLogin
  }
}
