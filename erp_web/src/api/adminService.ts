import request from './request'

// ========== 类型定义 ==========

/** 管理员信息 */
export interface AdminInfo {
  id: number
  username: string
  name: string
  avatar?: string
  email?: string
  phone?: string
  status: number // 0=禁用 1=正常
  is_admin: boolean
  last_login_time?: string
  create_time: string
}

/** 登录结果 */
export interface LoginResult {
  access_token: string
  refresh_token: string
  expires_in: number
  admin: AdminInfo
}

/** 刷新Token结果 */
export interface RefreshTokenResult {
  access_token: string
  refresh_token: string
  expires_in: number
}

/** 验证码结果 */
export interface CaptchaResult {
  captcha_id: string
  code: string
}

/** 登录参数 */
export interface LoginParams {
  username: string
  password: string
  captcha_id: string
  captcha_code: string
}

/** 管理员列表参数 */
export interface AdminListParams {
  page?: number
  page_size?: number
  keyword?: string
}

/** 管理员列表结果 */
export interface AdminListResult {
  list: AdminInfo[]
  total: number
  page: number
  page_size: number
}

/** 创建管理员参数 */
export interface CreateAdminParams {
  username: string
  password: string
  name: string
  email?: string
  phone?: string
}

/** 更新管理员参数 */
export interface UpdateAdminParams {
  name: string
  email?: string
  phone?: string
}

/** 修改密码参数 */
export interface ChangePasswordParams {
  old_password: string
  new_password: string
  confirm_password: string
}

/** 重置密码参数 */
export interface ResetPasswordParams {
  new_password: string
}

// ========== API 接口 ==========

export const adminService = {
  /**
   * 管理员登录
   */
  login: (params: LoginParams): Promise<LoginResult> => {
    return request.post('/admin/login', params)
  },

  /**
   * 获取验证码
   */
  getCaptcha: (): Promise<CaptchaResult> => {
    return request.get('/admin/captcha')
  },

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
  },

  /**
   * 获取当前管理员信息
   */
  getProfile: (): Promise<AdminInfo> => {
    return request.get('/admin/profile')
  },

  /**
   * 更新个人信息
   */
  updateProfile: (id: number, params: UpdateAdminParams): Promise<AdminInfo> => {
    return request.post(`/admin/update/${id}`, params)
  },

  /**
   * 修改密码
   */
  changePassword: (params: ChangePasswordParams): Promise<void> => {
    return request.post('/admin/password/change', params)
  },

  /**
   * 获取管理员列表
   */
  getList: (params: AdminListParams): Promise<AdminListResult> => {
    return request.get('/admin/list', { params })
  },

  /**
   * 获取管理员详情
   */
  getDetail: (id: number): Promise<AdminInfo> => {
    return request.get(`/admin/${id}`)
  },

  /**
   * 创建管理员
   */
  create: (params: CreateAdminParams): Promise<AdminInfo> => {
    return request.post('/admin/create', params)
  },

  /**
   * 更新管理员
   */
  update: (id: number, params: UpdateAdminParams): Promise<AdminInfo> => {
    return request.post(`/admin/update/${id}`, params)
  },

  /**
   * 删除管理员
   */
  delete: (id: number): Promise<void> => {
    return request.post(`/admin/delete/${id}`)
  },

  /**
   * 禁用管理员
   */
  ban: (id: number): Promise<void> => {
    return request.post(`/admin/ban/${id}`)
  },

  /**
   * 解禁管理员
   */
  unban: (id: number): Promise<void> => {
    return request.post(`/admin/unban/${id}`)
  },

  /**
   * 重置密码（仅超管）
   */
  resetPassword: (id: number, params: ResetPasswordParams): Promise<void> => {
    return request.post(`/admin/password/reset/${id}`, params)
  },

  /**
   * 强制下线（仅超管）
   */
  forceLogout: (id: number): Promise<void> => {
    return request.post(`/admin/force-logout/${id}`)
  }
}
