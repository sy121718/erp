<script setup lang="ts">
import { useLogin } from './utils/hook'

const { formData, loading, canSubmit, handleLogin } = useLogin()
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
            @keyup.enter="handleLogin"
          />
        </lay-form-item>

        <lay-form-item>
          <lay-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="layui-icon-password"
            size="lg"
            @keyup.enter="handleLogin"
          />
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
