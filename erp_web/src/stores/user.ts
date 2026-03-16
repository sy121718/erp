import { defineStore } from 'pinia'
import { authService, type LoginParams, type LoginResult } from '@/api/auth'

interface AdminInfo {
  id: number
  username: string
  nickname: string
  avatar?: string
}

interface UserState {
  accessToken: string
  refreshToken: string
  expiresIn: number
  adminInfo: AdminInfo | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    accessToken: localStorage.getItem('access_token') || '',
    refreshToken: localStorage.getItem('refresh_token') || '',
    expiresIn: 0,
    adminInfo: null
  }),

  getters: {
    isLoggedIn: (state) => !!state.accessToken,
    nickname: (state) => state.adminInfo?.nickname || '未登录',
    avatar: (state) => state.adminInfo?.avatar || ''
  },

  actions: {
    /**
     * 登录
     */
    async login(params: LoginParams) {
      const result: LoginResult = await authService.login(params)
      
      this.accessToken = result.access_token
      this.refreshToken = result.refresh_token
      this.expiresIn = result.expires_in
      this.adminInfo = result.admin

      // 持久化存储
      localStorage.setItem('access_token', result.access_token)
      localStorage.setItem('refresh_token', result.refresh_token)
      
      return result
    },

    /**
     * 刷新Token
     */
    async refreshAccessToken() {
      if (!this.refreshToken) {
        throw new Error('无刷新令牌')
      }

      const result = await authService.refreshToken(this.refreshToken)
      
      this.accessToken = result.access_token
      this.refreshToken = result.refresh_token
      this.expiresIn = result.expires_in

      localStorage.setItem('access_token', result.access_token)
      localStorage.setItem('refresh_token', result.refresh_token)
      
      return result
    },

    /**
     * 退出登录
     */
    async logout() {
      try {
        await authService.logout()
      } catch (error) {
        console.error('退出登录失败', error)
      } finally {
        this.clearAuth()
      }
    },

    /**
     * 清除认证信息
     */
    clearAuth() {
      this.accessToken = ''
      this.refreshToken = ''
      this.expiresIn = 0
      this.adminInfo = null

      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    },

    /**
     * 设置用户信息
     */
    setAdminInfo(info: AdminInfo) {
      this.adminInfo = info
    }
  }
})