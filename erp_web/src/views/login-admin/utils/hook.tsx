import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'
import { adminService } from '@/api/adminService'
import type { LoginForm, LoginSubmitParams } from './types'

/** 登录页面组合式函数 */
export const useLogin = () => {
  const router = useRouter()
  const userStore = useUserStore()

  const formData = ref<LoginForm>({
    username: '',
    password: '',
    captcha_code: ''
  })

  const captchaId = ref('')
  const loading = ref(false)

  const fetchCaptcha = () => adminService.getCaptcha()

  const canSubmit = computed<boolean>(() => {
    return !!(
      formData.value.username &&
      formData.value.password &&
      formData.value.captcha_code
    )
  })

  const handleLogin = async (): Promise<void> => {
    if (!canSubmit.value) {
      layer.msg('请填写完整信息', { icon: 0 })
      return
    }

    loading.value = true
    try {
      const params: LoginSubmitParams = {
        username: formData.value.username,
        password: formData.value.password,
        captcha_id: captchaId.value,
        captcha_code: formData.value.captcha_code
      }

      await userStore.login(params)
      router.push('/')
    } catch {
      // 错误由 request.ts 统一处理，组件会自动刷新验证码
    } finally {
      loading.value = false
    }
  }

  return {
    formData,
    captchaId,
    loading,
    canSubmit,
    fetchCaptcha,
    handleLogin
  }
}
