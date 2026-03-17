import request from './request'

// ========== 类型定义 ==========

/** 用户信息 */
export interface UserInfo {
  id: number
  username: string
  name: string
  avatar?: string
  email?: string
  phone?: string
  status: number
  points: number
  total_points: number
  used_points: number
  expire_time?: string
  last_login_time?: string
  create_time: string
}

/** 用户登录参数 */
export interface UserLoginParams {
  username: string
  password: string
  captcha_id: string
  captcha_code: string
}

/** 验证码结果 */
export interface UserCaptchaResult {
  captcha_id: string
  code: string
}

/** 用户登录结果 */
export interface UserLoginResult {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

/** 用户注册参数 */
export interface UserRegisterParams {
  username: string
  password: string
  name?: string
  email?: string
  phone?: string
}

/** 用户更新参数 */
export interface UserUpdateParams {
  name?: string
  email?: string
  phone?: string
  avatar?: string
}

/** 用户修改密码参数 */
export interface UserChangePasswordParams {
  old_password: string
  new_password: string
}

/** 刷新Token结果 */
export interface UserRefreshTokenResult {
  access_token: string
  refresh_token: string
  expires_in: number
}

/** 用户列表参数（管理员侧） */
export interface UserListParams {
  page?: number
  page_size?: number
  keyword?: string
}

/** 用户列表结果 */
export interface UserListResult {
  list: UserInfo[]
  total: number
  page: number
  page_size: number
}

/** 管理员更新用户参数 */
export interface AdminUpdateUserParams {
  name?: string
  email?: string
  phone?: string
}

// ========== 用户侧 API ==========

export const userService = {
  getCaptcha: (): Promise<UserCaptchaResult> => {
    return request.get('/user/captcha')
  },

  login: (params: UserLoginParams): Promise<UserLoginResult> => {
    return request.post('/user/login', params)
  },

  register: (params: UserRegisterParams): Promise<UserInfo> => {
    return request.post('/user/register', params)
  },

  refreshToken: (refresh_token: string): Promise<UserRefreshTokenResult> => {
    return request.post('/user/refresh-token', { refresh_token })
  },

  logout: (): Promise<void> => {
    return request.post('/user/logout')
  },

  getProfile: (): Promise<UserInfo> => {
    return request.get('/user/profile')
  },

  updateProfile: (params: UserUpdateParams): Promise<UserInfo> => {
    return request.post('/user/update', params)
  },

  changePassword: (params: UserChangePasswordParams): Promise<void> => {
    return request.post('/user/password/change', params)
  },

  // ========== 管理员侧用户管理 API ==========

  adminGetList: (params: UserListParams): Promise<UserListResult> => {
    return request.get('/admin/user/list', { params })
  },

  adminGetDetail: (id: number): Promise<UserInfo> => {
    return request.get(`/admin/user/${id}`)
  },

  adminUpdate: (id: number, params: AdminUpdateUserParams): Promise<UserInfo> => {
    return request.post(`/admin/user/update/${id}`, params)
  },

  adminBan: (id: number): Promise<void> => {
    return request.post(`/admin/user/ban/${id}`)
  },

  adminUnban: (id: number): Promise<void> => {
    return request.post(`/admin/user/unban/${id}`)
  },

  adminResetPassword: (id: number): Promise<void> => {
    return request.post(`/admin/user/password/reset/${id}`)
  }
}
