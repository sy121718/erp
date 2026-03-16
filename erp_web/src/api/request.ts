import axios, { type AxiosError, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { layer } from '@layui/layui-vue'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000
})

// 请求拦截
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
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
    const { code, msg, data } = response.data
    if (code === 0) {
      return data
    }
    layer.msg(msg || '请求失败', { icon: 2 })
    return Promise.reject(new Error(msg || '请求失败'))
  },
  (error: AxiosError<{ msg?: string }>) => {
    const message = error.response?.data?.msg || error.message || '网络错误'
    layer.msg(message, { icon: 2 })
    return Promise.reject(error)
  }
)

export default request