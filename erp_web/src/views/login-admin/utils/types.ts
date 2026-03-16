/** 登录表单数据 */
export interface LoginForm {
  username: string
  password: string
  captcha_code: string
}

/** 验证码数据 */
export interface CaptchaData {
  captchaId: string
  captchaCode: string
}

/** 登录参数 (提交给后端) */
export interface LoginSubmitParams {
  username: string
  password: string
  captcha_id: string
  captcha_code: string
}
