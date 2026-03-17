import { defineStore } from 'pinia'
import { adminService, type LoginParams, type LoginResult } from '@/api/adminService'

interface AdminInfo {
  id: number
  username: string
  name: string
  avatar?: string
  email?: string
  phone?: string
  status: number
  is_admin: boolean
}

interface UserState {
  accessToken: string
  refreshToken: string
  expiresIn: number
  adminInfo: AdminInfo | null
  initialized: boolean // 是否已初始化
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    accessToken: localStorage.getItem('access_token') || '',
    refreshToken: localStorage.getItem('refresh_token') || '',
    expiresIn: 0,
    adminInfo: null,
    initialized: false
  }),

  getters: {
    isLoggedIn: (state) => !!state.accessToken,
    nickname: (state) => state.adminInfo?.name || '未登录',
    avatar: (state) => state.adminInfo?.avatar || '/image/avatar.png',
    isAdmin: (state) => state.adminInfo?.is_admin || false
  },

  actions: {
    /**
     * 登录
     */
    async login(params: LoginParams) {
      const result: LoginResult = await adminService.login(params)

      this.accessToken = result.access_token
      this.refreshToken = result.refresh_token
      this.expiresIn = result.expires_in

      // 设置默认头像
      if (result.admin && !result.admin.avatar) {
        result.admin.avatar = '/image/avatar.png'
      }
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

      const result = await adminService.refreshToken(this.refreshToken)

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
      // 先清除本地认证信息（避免 401 触发刷新 token）
      this.clearAuth()

      // 调用后端退出接口（在白名单中，不需要 token）
      try {
        await adminService.logout()
      } catch (error) {
        console.error('退出登录失败', error)
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
    },

    /**
     * 初始化用户信息
     * 页面刷新后，如果有 token，重新获取用户信息
     */
    async initUserInfo() {
      // 如果已经初始化或者没有 token，直接返回
      if (this.initialized || !this.accessToken) {
        return
      }

      try {
        const info = await adminService.getProfile()
        this.adminInfo = info
        this.initialized = true
        console.log('✅ 初始化用户信息成功:', info)
      } catch (error) {
        console.error('❌ 初始化用户信息失败:', error)
        // 获取失败，清除 token
        this.clearAuth()
      }
    }
  }
})
