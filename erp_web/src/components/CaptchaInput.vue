<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue'

interface CaptchaResult {
  captcha_id: string
  display: string
  code: string
}

interface Props {
  /** 获取验证码的 API 函数 */
  fetchApi: () => Promise<CaptchaResult>
  /** 验证码输入值 (v-model:code) */
  code?: string
  /** 验证码 ID (v-model:captchaId) */
  captchaId?: string
  /** 倒计时秒数，默认 60 */
  countdown?: number
  /** 输入框 placeholder */
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  code: '',
  captchaId: '',
  countdown: 60,
  placeholder: '请输入验证码'
})

const emit = defineEmits<{
  'update:code': [value: string]
  'update:captchaId': [value: string]
}>()

const captchaImage = ref('')
const loading = ref(false)
const countdownValue = ref(0)
let countdownTimer: number | null = null

const startCountdown = () => {
  if (countdownTimer) clearInterval(countdownTimer)
  countdownValue.value = props.countdown
  countdownTimer = window.setInterval(() => {
    countdownValue.value--
    if (countdownValue.value <= 0 && countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
      refresh()
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
  const charWidth = Math.min(18, (width - 30) / code.length)
  const startX = 15

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

const refresh = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await props.fetchApi()
    emit('update:captchaId', res.captcha_id)
    captchaImage.value = drawCaptcha(res.display || res.code)
    startCountdown()
  } catch {
    // 错误由 request.ts 统一处理
  } finally {
    loading.value = false
  }
}

const handleInput = (val: string | number) => {
  emit('update:code', String(val))
}

watch(() => props.fetchApi, () => {
  refresh()
})

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
})

refresh()

defineExpose({ refresh })
</script>

<template>
  <div class="captcha-row">
    <lay-input
      :modelValue="code"
      @update:modelValue="handleInput"
      :placeholder="placeholder"
      size="lg"
      class="captcha-input"
    />
    <div class="captcha-wrapper">
      <div class="captcha-box" @click="refresh">
        <img v-if="captchaImage" :src="captchaImage" alt="验证码" class="captcha-img" />
        <span v-else class="captcha-placeholder">加载中...</span>
      </div>
      <div class="captcha-refresh" @click="refresh">
        点击刷新
      </div>
    </div>
  </div>
</template>

<style scoped>
.captcha-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.captcha-input {
  flex: 1;
}

.captcha-wrapper {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex-shrink: 0;
}

.captcha-box {
  width: 120px;
  height: 42px;
  background-color: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
}

.captcha-box:hover {
  border-color: var(--erp-primary);
}

.captcha-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  pointer-events: none;
}

.captcha-placeholder {
  font-size: 12px;
  color: var(--erp-text-tertiary);
}

.captcha-refresh {
  width: 120px;
  text-align: center;
  font-size: 12px;
  color: var(--erp-primary);
  cursor: pointer;
  user-select: none;
  padding: 4px 0;
  transition: all 0.2s;
}

.captcha-refresh:hover {
  color: var(--erp-primary-dark);
  text-decoration: underline;
}
</style>
