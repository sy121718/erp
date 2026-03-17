import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { layer } from '@layui/layui-vue'
import router, { ADMIN_LOGIN_PATH } from '@/router'
import { isWhiteListed, addSignToRequest } from '@/utils/sign'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000
})

// 请求拦截
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 添加 Token
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 添加签名（白名单路由不需要签名）
    if (!isWhiteListed(config.url || '')) {
      let signParams

      if (config.method?.toUpperCase() === 'GET') {
        // GET 请求：对 query 参数签名
        signParams = addSignToRequest(config.method, config.url || '', config.params)
      } else {
        // POST/PUT/DELETE 请求：对 body 参数签名
        signParams = addSignToRequest(config.method || 'POST', config.url || '', config.data)
      }

      config.headers['X-Timestamp'] = signParams.timestamp
      config.headers['X-Nonce'] = signParams.nonce
      config.headers['X-Sign'] = signParams.sign
    }

    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// 响应拦截
request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 检查响应头中是否有新的token
    const newToken = response.headers['x-new-access-token']
    if (newToken) {
      // 保存新的token
      localStorage.setItem('access_token', newToken)
    }

    const { code, msg, data } = response.data
    if (code === 0) {
      return data
    }

    // 业务错误码处理
    layer.msg(msg || '请求失败', { icon: 2 })
    return Promise.reject(new Error(msg || '请求失败'))
  },
  (error: AxiosError<{ msg?: string; code?: number }>) => {
    const { response } = error

    // 401 未授权 - Token 过期
    if (response?.status === 401) {
      // 清除本地认证信息
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')

      // 跳转到登录页
      router.push(ADMIN_LOGIN_PATH)
      layer.msg('登录已过期，请重新登录', { icon: 2 })

      return Promise.reject(error)
    }

    // 其他错误
    const msg = response?.data?.msg || error.message || '网络错误'
    layer.msg(msg, { icon: 2 })
    return Promise.reject(error)
  }
)

export default request
