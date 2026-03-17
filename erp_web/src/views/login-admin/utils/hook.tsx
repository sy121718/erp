import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { layer } from '@layui/layui-vue'
import { useUserStore } from '@/stores/user'
import { adminService } from '@/api/adminService'
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
    captchaCode: ''
  })

  // 验证码图片
  const captchaImage = ref('')

  // 加载状态
  const loading = ref(false)

  // 验证码过期倒计时（1分钟 = 60秒）
  const countdown = ref(60)
  let countdownTimer: number | null = null

  // 开始倒计时
  const startCountdown = () => {
    // 清除之前的定时器
    if (countdownTimer) {
      clearInterval(countdownTimer)
    }

    countdown.value = 60 // 1分钟 = 60秒
    countdownTimer = window.setInterval(() => {
      countdown.value--
      if (countdown.value <= 0 && countdownTimer) {
        clearInterval(countdownTimer)
        countdownTimer = null
        // 验证码过期，自动刷新
        getCaptcha()
      }
    }, 1000)
  }

  // 绘制验证码图片
  const drawCaptcha = (code: string): string => {
    const canvas = document.createElement('canvas')
    const width = 120
    const height = 42
    canvas.width = width
    canvas.height = height
    const ctx = canvas.getContext('2d')

    if (!ctx) {
      return ''
    }

    // 背景色
    ctx.fillStyle = '#f3f4f6'
    ctx.fillRect(0, 0, width, height)

    // 绘制干扰线
    for (let i = 0; i < 5; i++) {
      ctx.strokeStyle = `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.5)`
      ctx.beginPath()
      ctx.moveTo(Math.random() * width, Math.random() * height)
      ctx.lineTo(Math.random() * width, Math.random() * height)
      ctx.stroke()
    }

    // 绘制干扰点
    for (let i = 0; i < 50; i++) {
      ctx.fillStyle = `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.5)`
      ctx.beginPath()
      ctx.arc(Math.random() * width, Math.random() * height, 1, 0, Math.PI * 2)
      ctx.fill()
    }

    // 绘制验证码文字
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

  // 获取验证码
  const getCaptcha = async (): Promise<void> => {
    try {
      const res = await adminService.getCaptcha()

      // 更新验证码数据
      captchaData.value = {
        captchaId: res.captcha_id,
        captchaCode: res.code
      }

      // 立即绘制验证码图片
      captchaImage.value = drawCaptcha(res.code)

      // 开始倒计时
      startCountdown()
    } catch (error) {
      // 错误由 request.ts 统一处理
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

      await userStore.login(params)
      // 登录成功立即跳转
      router.push('/')
    } catch (error) {
      // 错误由 request.ts 统一处理
      // 延迟刷新验证码
      setTimeout(() => {
        getCaptcha()
      }, 2000)
    } finally {
      loading.value = false
    }
  }

  // 清理定时器
  onUnmounted(() => {
    if (countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
  })

  // 初始化
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
