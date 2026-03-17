import CryptoJS from 'crypto-js'

/**
 * 签名配置
 */
const SIGN_CONFIG = {
  secret: 'erp-sign-secret-key-2024-secure',
  expireTime: 300 // 5分钟
}

/**
 * 白名单路由（不需要签名）
 */
const WHITE_LIST = [
  '/health',
  '/api/ping',
  '/api/admin/login',
  '/api/admin/captcha',
  '/api/admin/refresh-token',
  '/api/user/captcha',
  '/api/user/login',
  '/api/user/register',
  '/api/user/refresh-token'
]

/**
 * 检查路径是否在白名单中
 */
export function isWhiteListed(path: string): boolean {
  return WHITE_LIST.some(p => path === p || path.startsWith(p + '/'))
}

/**
 * 生成随机字符串
 */
function generateNonce(length: number = 16): string {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let result = ''

  // 使用 Web Crypto API 生成更安全的随机数
  if (window.crypto && window.crypto.getRandomValues) {
    const values = new Uint32Array(length)
    window.crypto.getRandomValues(values)
    for (let i = 0; i < length; i++) {
      result += chars.charAt(values[i] % chars.length)
    }
  } else {
    // 降级方案：使用 Math.random()
    for (let i = 0; i < length; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
  }

  return result
}

/**
 * 生成签名
 * @param params 业务参数
 * @param timestamp 时间戳（毫秒）
 * @param nonce 随机字符串
 * @returns 签名字符串
 */
function generateSign(params: Record<string, string>, timestamp: number, nonce: string): string {
  // 复制参数，避免修改原参数
  const allParams: Record<string, string> = { ...params }

  // 添加 timestamp 和 nonce
  allParams.timestamp = timestamp.toString()
  allParams.nonce = nonce

  // 获取所有 key 并排序
  const keys = Object.keys(allParams).sort()

  // 拼接字符串: key1=value1&key2=value2
  const signStr = keys.map(key => `${key}=${allParams[key]}`).join('&')

  // HMAC-SHA256 签名
  const signature = CryptoJS.HmacSHA256(signStr, SIGN_CONFIG.secret)
  return CryptoJS.enc.Hex.stringify(signature)
}

/**
 * 生成签名参数
 * @param params 业务参数
 * @returns 签名参数对象
 */
export function generateSignParams(params: Record<string, string> = {}): {
  timestamp: number
  nonce: string
  sign: string
} {
  const timestamp = Date.now()
  const nonce = generateNonce(16)
  const sign = generateSign(params, timestamp, nonce)

  return { timestamp, nonce, sign }
}

/**
 * 为请求添加签名
 * @param method HTTP 方法
 * @param url 请求 URL
 * @param params 请求参数（GET: query, POST: body）
 * @returns 签名参数
 */
export function addSignToRequest(
  method: string,
  url: string,
  params?: Record<string, any>
): { timestamp: number; nonce: string; sign: string } {
  // 将参数转换为字符串键值对
  const stringParams: Record<string, string> = {}

  if (params) {
    // 递归展平嵌套对象
    flattenObject(params, '', stringParams)
  }

  return generateSignParams(stringParams)
}

/**
 * 递归展平嵌套对象为字符串键值对
 * @param obj 原始对象
 * @param prefix 键前缀
 * @param result 结果对象
 */
function flattenObject(obj: any, prefix: string, result: Record<string, string>): void {
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const value = obj[key]
      const newKey = prefix ? `${prefix}.${key}` : key

      // 跳过 undefined 和 null
      if (value === undefined || value === null) {
        continue
      }

      // 如果是对象（非数组、非日期、非文件）
      if (typeof value === 'object' && value !== null && !Array.isArray(value) && !(value instanceof Date) && !(value instanceof File)) {
        flattenObject(value, newKey, result)
      } else {
        // 转换为字符串
        result[newKey] = String(value)
      }
    }
  }
}
