import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'
import { userService } from '@/api/userService'
import type { LoginForm, CaptchaData, LoginSubmitParams } from './types'

/** 用户登录页面组合式函数 */
export const useLogin = () => {
  const router = useRouter()
  const userStore = useUserStore()

  const formData = ref<LoginForm>({
    username: '',
    password: '',
    captcha_code: ''
  })

  const captchaData = ref<CaptchaData>({
    captchaId: '',
    captchaCode: ''
  })

  const captchaImage = ref('')
  const loading = ref(false)

  const countdown = ref(60)
  let countdownTimer: number | null = null

  const startCountdown = () => {
    if (countdownTimer) {
      clearInterval(countdownTimer)
    }

    countdown.value = 60
    countdownTimer = window.setInterval(() => {
      countdown.value--
      if (countdown.value <= 0 && countdownTimer) {
        clearInterval(countdownTimer)
        countdownTimer = null
        getCaptcha()
      }
    }, 1000)
  }

  const drawCaptcha = (code: string): string => {
    const canvas = document.createElement('canvas')
    const width = 120
    const height = 42
    canvas.width = width
    canvas.height = height
    const ctx = canvas.getContext('2d')

    if (!ctx) return ''

    ctx.fillStyle = '#f3f4f6'
    ctx.fillRect(0, 0, width, height)

    for (let i = 0; i < 5; i++) {
      ctx.strokeStyle = `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.5)`
      ctx.beginPath()
      ctx.moveTo(Math.random() * width, Math.random() * height)
      ctx.lineTo(Math.random() * width, Math.random() * height)
      ctx.stroke()
    }

    for (let i = 0; i < 50; i++) {
      ctx.fillStyle = `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.5)`
      ctx.beginPath()
      ctx.arc(Math.random() * width, Math.random() * height, 1, 0, Math.PI * 2)
      ctx.fill()
    }

    ctx.font = 'bold 24px Arial'
    ctx.textBaseline = 'middle'

    const colors = ['#03b3b2', '#ff5722', '#4caf50', '#2196f3']
    const startX = 15
    const charWidth = 18

    for (let i = 0; i < code.length; i++) {
      ctx.fillStyle = colors[Math.floor(Math.random() * colors.length)]
      const x = startX + i * charWidth
      const y = height / 2 + (Math.random() * 10 - 5)
      const angle = (Math.random() * 30 - 15) * Math.PI / 180

      ctx.save()
      ctx.translate(x, y)
      ctx.rotate(angle)
      ctx.fillText(code[i], 0, 0)
      ctx.restore()
    }

    return canvas.toDataURL()
  }

  const getCaptcha = async (): Promise<void> => {
    try {
      const res = await userService.getCaptcha()
      captchaData.value = {
        captchaId: res.captcha_id,
        captchaCode: res.code
      }
      captchaImage.value = drawCaptcha(res.code)
      startCountdown()
    } catch {
      // 错误由 request.ts 统一处理
    }
  }

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
        captcha_id: captchaData.value.captchaId,
        captcha_code: formData.value.captcha_code
      }

      await userStore.userLogin(params)
      router.push('/')
    } catch {
      setTimeout(() => {
        getCaptcha()
      }, 2000)
    } finally {
      loading.value = false
    }
  }

  onUnmounted(() => {
    if (countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
  })

  getCaptcha()

  return {
    formData,
    captchaData,
    captchaImage,
    loading,
    canSubmit,
    getCaptcha,
    handleLogin
  }
}
