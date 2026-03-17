import request from './request'

export interface RefreshTokenResult {
  access_token: string
  refresh_token: string
  expires_in: number
}

export const authService = {
  /**
   * 刷新Token
   */
  refreshToken: (refresh_token: string): Promise<RefreshTokenResult> => {
    return request.post('/admin/refresh-token', { refresh_token })
  },

  /**
   * 退出登录
   */
  logout: (): Promise<void> => {
    return request.post('/admin/logout')
  }
}
