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

export interface RefreshTokenResult {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface CaptchaResult {
  captcha_id: string
  captcha_image: string
}

export const authService = {
  /**
   * 管理员登录
   */
  login: (params: LoginParams): Promise<LoginResult> => {
    return request.post('/admin/login', params)
  },

  /**
   * 刷新Token
   */
  refreshToken: (refresh_token: string): Promise<RefreshTokenResult> => {
    return request.post('/admin/refresh-token', { refresh_token })
  },

  /**
   * 获取验证码
   */
  getCaptcha: (): Promise<CaptchaResult> => {
    return request.get('/admin/captcha')
  },

  /**
   * 退出登录
   */
  logout: (): Promise<void> => {
    return request.post('/admin/logout')
  }
}
