# API 接口规范

## 重要规则

**创建或修改任何 API 接口前，必须按以下顺序确认接口路径：**

1. **首先**：使用 Apipost MCP 工具查询接口文档
2. **其次**：如果文档不存在或有误，才查看后端路由文件 `erp_server/internal/handler/{模块}/{模块}_router.go`

**查询命令示例：**
```
# 搜索接口
apipost_list search="关键词" show_path=true

# 查看接口详情
apipost_detail target_id="接口ID"
```

## 目录结构

```
src/api/
├── index.ts          # 统一导出
├── request.ts        # axios 封装
├── auth.ts           # 认证相关接口
├── userService.ts    # 用户相关接口
├── productService.ts # 产品相关接口
└── ...
```

## 接口定义

```typescript
// src/api/auth.ts
import request from './request'

export interface LoginParams {
  username: string
  password: string
  captcha_id: string
  captcha_code: string
}

export interface LoginResult {
  access_token: string
  refresh_token: string
  expires_in: number
  admin: {
    id: number
    username: string
    nickname: string
    avatar?: string
  }
}

export const authService = {
  // 登录
  login: (params: LoginParams) => request.post('/admin/login', params),

  // 获取验证码
  getCaptcha: () => request.get('/admin/captcha'),

  // 刷新Token
  refreshToken: (refresh_token: string) => request.post('/admin/refresh-token', { refresh_token }),

  // 退出登录
  logout: () => request.post('/admin/logout')
}
```

## Request 封装

```typescript
// src/api/request.ts
import axios, { type AxiosError, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { layer } from '@layui/layui-vue'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000
})

// 请求拦截 - 添加 Token
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

// 响应拦截 - 统一处理响应
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
    layer.msg(error.response?.data?.msg || error.message || '网络错误', { icon: 2 })
    return Promise.reject(error)
  }
)

export default request
```

## 使用示例

```typescript
import { authService } from '@/api/auth'

// 获取验证码
const captcha = await authService.getCaptcha()

// 登录
const result = await authService.login({
  username: 'admin',
  password: '123456',
  captcha_id: 'xxx',
  captcha_code: '1234'
})
```
