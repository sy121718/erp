import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'
import { authService } from '@/api/auth'
import type { LoginForm, CaptchaData, LoginSubmitParams } from './types'

/** 登录页面组合式函数 */
export const useLogin = () => {
  const router = useRouter()
  const userStore = useUserStore()

  // 表单数据
  const formData = ref<LoginForm>({
    username: '',
    password: '',
    captcha_code: ''
  })

  // 验证码
  const captchaData = ref<CaptchaData>({
    captchaId: '',
    captchaImage: ''
  })

  // 加载状态
  const loading = ref(false)

  // 获取验证码
  const getCaptcha = async (): Promise<void> => {
    try {
      const res = await authService.getCaptcha()
      captchaData.value = {
        captchaId: res.captcha_id,
        captchaImage: res.captcha_image
      }
    } catch (error) {
      console.error('获取验证码失败', error)
      layer.msg('获取验证码失败', { icon: 2 })
    }
  }

  // 表单验证
  const canSubmit = computed<boolean>(() => {
    return !!(
      formData.value.username &&
      formData.value.password &&
      formData.value.captcha_code
    )
  })

  // 登录
  const handleLogin = async (): Promise<void> => {
    if (!canSubmit.value) {
      layer.msg('请填写完整信息', { icon: 2 })
      return
    }

    loading.value = true
    try {
      const params: LoginSubmitParams = {
        username: formData.value.username,
        password: formData.value.password,
        captcha_id: captchaData.value.captchaId,
        captcha_code: formData.value.captcha_code
      }

      await userStore.login(params)
      layer.msg('登录成功', { icon: 1 })
      router.push('/')
    } catch (error) {
      layer.msg('登录失败', { icon: 2 })
      getCaptcha()
    } finally {
      loading.value = false
    }
  }

  // 初始化
  getCaptcha()

  return {
    formData,
    captchaData,
    loading,
    canSubmit,
    getCaptcha,
    handleLogin
  }
}
