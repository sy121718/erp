<script setup lang="ts">
import { useLogin } from './utils/hook'

const {
  formData,
  captchaImage,
  loading,
  canSubmit,
  getCaptcha,
  handleLogin
} = useLogin()
</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h1 class="logo-title">妙手ERP-Copy</h1>
        <p class="logo-subtitle">用户登录</p>
      </div>

      <lay-form class="login-form">
        <lay-form-item>
          <lay-input
            v-model="formData.username"
            placeholder="请输入用户名"
            prefix-icon="layui-icon-username"
            size="lg"
          />
        </lay-form-item>

        <lay-form-item>
          <lay-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="layui-icon-password"
            size="lg"
          />
        </lay-form-item>

        <lay-form-item>
          <div class="captcha-row">
            <lay-input
              v-model="formData.captcha_code"
              placeholder="请输入验证码"
              size="lg"
              class="captcha-input"
            />
            <div class="captcha-wrapper">
              <div class="captcha-box" @click="getCaptcha">
                <img v-if="captchaImage" :src="captchaImage" alt="验证码" class="captcha-img" />
                <span v-else class="captcha-placeholder">暂无验证码</span>
              </div>
              <div class="captcha-refresh" @click="getCaptcha">
                点击刷新
              </div>
            </div>
          </div>
        </lay-form-item>

        <lay-form-item>
          <lay-button
            type="primary"
            size="lg"
            fluid
            :loading="loading"
            :disabled="!canSubmit"
            @click="handleLogin"
          >
            登 录
          </lay-button>
        </lay-form-item>
      </lay-form>

      <div class="login-footer">
        <span>© 2024 妙手ERP All Rights Reserved</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding-top: 15%;
  padding-right: 10%;
  background-image: url('/image/login.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.login-box {
  width: 400px;
  background-color: var(--erp-bg-card);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  padding: 48px 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--erp-text-primary);
  margin-bottom: 8px;
}

.logo-subtitle {
  font-size: 14px;
  color: var(--erp-text-tertiary);
}

.login-form {
  margin-bottom: 24px;
}

.captcha-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.captcha-input {
  flex: 1;
}

.captcha-wrapper {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex-shrink: 0;
}

.captcha-box {
  width: 120px;
  height: 42px;
  background-color: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
}

.captcha-box:hover {
  border-color: var(--erp-primary);
}

.captcha-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  pointer-events: none;
}

.captcha-placeholder {
  font-size: 12px;
  color: var(--erp-text-tertiary);
}

.captcha-refresh {
  width: 120px;
  text-align: center;
  font-size: 12px;
  color: var(--erp-primary);
  cursor: pointer;
  user-select: none;
  padding: 4px 0;
  transition: all 0.2s;
}

.captcha-refresh:hover {
  color: var(--erp-primary-dark);
  text-decoration: underline;
}

.login-footer {
  text-align: center;
  font-size: 12px;
  color: var(--erp-text-disabled);
}

:deep(.layui-input-lg) {
  height: 42px;
}

:deep(.layui-btn-lg) {
  height: 42px;
  font-size: 16px;
}

:deep(.layui-form-item) {
  margin-bottom: 20px;
}
</style>
