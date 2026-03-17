import { defineStore } from 'pinia'
import { adminService, type LoginParams, type LoginResult } from '@/api/adminService'
import { userService, type UserLoginParams, type UserLoginResult, type UserInfo } from '@/api/userService'

export type UserType = 'admin' | 'user' | ''

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
  userType: UserType
  adminInfo: AdminInfo | null
  userInfo: UserInfo | null
  initialized: boolean
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    accessToken: localStorage.getItem('access_token') || '',
    refreshToken: localStorage.getItem('refresh_token') || '',
    expiresIn: 0,
    userType: (localStorage.getItem('user_type') as UserType) || '',
    adminInfo: null,
    userInfo: null,
    initialized: false
  }),

  getters: {
    isLoggedIn: (state) => !!state.accessToken,
    isAdminType: (state) => state.userType === 'admin',
    isUserType: (state) => state.userType === 'user',
    isAdmin: (state) => state.adminInfo?.is_admin || false,

    nickname: (state) => {
      if (state.userType === 'admin') return state.adminInfo?.name || '管理员'
      if (state.userType === 'user') return state.userInfo?.name || '用户'
      return '未登录'
    },

    avatar: (state) => {
      if (state.userType === 'admin') return state.adminInfo?.avatar || '/image/avatar.png'
      if (state.userType === 'user') return state.userInfo?.avatar || '/image/avatar.png'
      return '/image/avatar.png'
    }
  },

  actions: {
    /**
     * 管理员登录
     */
    async login(params: LoginParams) {
      const result: LoginResult = await adminService.login(params)

      this.accessToken = result.access_token
      this.refreshToken = result.refresh_token
      this.expiresIn = result.expires_in
      this.userType = 'admin'

      if (result.admin && !result.admin.avatar) {
        result.admin.avatar = '/image/avatar.png'
      }
      this.adminInfo = result.admin
      this.userInfo = null

      localStorage.setItem('access_token', result.access_token)
      localStorage.setItem('refresh_token', result.refresh_token)
      localStorage.setItem('user_type', 'admin')

      return result
    },

    /**
     * 用户登录
     */
    async userLogin(params: UserLoginParams) {
      const result: UserLoginResult = await userService.login(params)

      this.accessToken = result.access_token
      this.refreshToken = result.refresh_token
      this.expiresIn = result.expires_in
      this.userType = 'user'

      if (result.user && !result.user.avatar) {
        result.user.avatar = '/image/avatar.png'
      }
      this.userInfo = result.user
      this.adminInfo = null

      localStorage.setItem('access_token', result.access_token)
      localStorage.setItem('refresh_token', result.refresh_token)
      localStorage.setItem('user_type', 'user')

      return result
    },

    /**
     * 刷新Token
     */
    async refreshAccessToken() {
      if (!this.refreshToken) {
        throw new Error('无刷新令牌')
      }

      if (this.userType === 'user') {
        const result = await userService.refreshToken(this.refreshToken)
        this.accessToken = result.access_token
        this.refreshToken = result.refresh_token
        this.expiresIn = result.expires_in
        localStorage.setItem('access_token', result.access_token)
        localStorage.setItem('refresh_token', result.refresh_token)
        return result
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
      const wasUserType = this.userType
      this.clearAuth()

      try {
        if (wasUserType === 'user') {
          await userService.logout()
        } else {
          await adminService.logout()
        }
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
      this.userType = ''
      this.adminInfo = null
      this.userInfo = null
      this.initialized = false

      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user_type')
    },

    /**
     * 设置管理员信息
     */
    setAdminInfo(info: AdminInfo) {
      this.adminInfo = info
    },

    /**
     * 设置用户信息
     */
    setUserInfo(info: UserInfo) {
      this.userInfo = info
    },

    /**
     * 获取登录后的重定向路径
     */
    getLoginPath(): string {
      if (this.userType === 'admin') return '/ms-auth-admin'
      return '/login'
    },

    /**
     * 初始化用户信息
     */
    async initUserInfo() {
      if (this.initialized || !this.accessToken) {
        return
      }

      try {
        if (this.userType === 'user') {
          const info = await userService.getProfile()
          this.userInfo = info
        } else {
          const info = await adminService.getProfile()
          this.adminInfo = info
        }
        this.initialized = true
      } catch (error) {
        console.error('初始化用户信息失败:', error)
        this.clearAuth()
      }
    }
  }
})
