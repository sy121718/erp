<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { layer } from '@layui/layui-vue'
import { adminService } from '@/api/adminService'
import type { AdminInfo, UpdateAdminParams, ChangePasswordParams } from '@/api/adminService'

// 个人信息
const profile = ref<AdminInfo | null>(null)
const loading = ref(false)

// 编辑表单
const editForm = ref<UpdateAdminParams>({
  name: '',
  email: '',
  phone: ''
})

// 密码表单
const passwordForm = ref<ChangePasswordParams>({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

// Tab 切换
const currentTab = ref('profile')

// 获取个人信息
const fetchProfile = async () => {
  loading.value = true
  try {
    const data = await adminService.getProfile()
    profile.value = data
    // 填充编辑表单
    editForm.value = {
      name: data.name,
      email: data.email || '',
      phone: data.phone || ''
    }
  } catch (error) {
    // 错误由 request.ts 统一处理
  } finally {
    loading.value = false
  }
}

// 更新个人信息
const handleUpdateProfile = async () => {
  if (!editForm.value.name) {
    layer.msg('姓名不能为空', { icon: 0 })
    return
  }

  loading.value = true
  try {
    if (!profile.value) return
    const data = await adminService.updateProfile(profile.value.id, editForm.value)
    profile.value = data
    layer.msg('更新成功', { icon: 1 })
  } catch (error) {
    // 错误由 request.ts 统一处理
  } finally {
    loading.value = false
  }
}

// 修改密码
const handleChangePassword = async () => {
  if (!passwordForm.value.old_password) {
    layer.msg('请输入原密码', { icon: 0 })
    return
  }
  if (!passwordForm.value.new_password) {
    layer.msg('请输入新密码', { icon: 0 })
    return
  }
  if (passwordForm.value.new_password.length < 6) {
    layer.msg('新密码至少6位', { icon: 0 })
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    layer.msg('两次密码不一致', { icon: 0 })
    return
  }

  loading.value = true
  try {
    await adminService.changePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password
    })
    layer.msg('密码修改成功', { icon: 1 })
    // 清空密码表单
    passwordForm.value = {
      old_password: '',
      new_password: '',
      confirm_password: ''
    }
  } catch (error) {
    // 错误由 request.ts 统一处理
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProfile()
})
</script>

<template>
  <div class="profile-container">
    <lay-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h3>个人中心</h3>
        </div>
      </template>

      <lay-tab v-model="currentTab" type="brief">
        <!-- 个人信息 Tab -->
        <lay-tab-item title="个人信息" id="profile">
          <div class="profile-content" v-if="profile">
            <!-- 头像区域 -->
            <div class="avatar-section">
              <lay-avatar
                :src="profile.avatar || '/image/avatar.png'"
                :text="profile.name?.charAt(0)"
                size="lg"
              />
              <div class="user-info">
                <h2 class="user-name">{{ profile.name }}</h2>
                <p class="user-role">{{ profile.is_admin ? '超级管理员' : '管理员' }}</p>
              </div>
            </div>

            <!-- 信息列表 -->
            <div class="info-list">
              <div class="info-item">
                <span class="info-label">用户名</span>
                <span class="info-value">{{ profile.username }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">姓名</span>
                <span class="info-value">{{ profile.name }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">邮箱</span>
                <span class="info-value">{{ profile.email || '未设置' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">手机</span>
                <span class="info-value">{{ profile.phone || '未设置' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">状态</span>
                <lay-tag :type="profile.status === 1 ? 'normal' : 'danger'">
                  {{ profile.status === 1 ? '正常' : '禁用' }}
                </lay-tag>
              </div>
              <div class="info-item">
                <span class="info-label">创建时间</span>
                <span class="info-value">{{ profile.create_time }}</span>
              </div>
            </div>
          </div>
        </lay-tab-item>

        <!-- 编辑信息 Tab -->
        <lay-tab-item title="编辑信息" id="edit">
          <div class="edit-content">
            <lay-form class="edit-form" label-position="top">
              <lay-form-item label="姓名" required>
                <lay-input
                  v-model="editForm.name"
                  placeholder="请输入姓名"
                  :disabled="loading"
                />
              </lay-form-item>

              <lay-form-item label="邮箱">
                <lay-input
                  v-model="editForm.email"
                  placeholder="请输入邮箱"
                  :disabled="loading"
                />
              </lay-form-item>

              <lay-form-item label="手机">
                <lay-input
                  v-model="editForm.phone"
                  placeholder="请输入手机号"
                  :disabled="loading"
                />
              </lay-form-item>

              <lay-form-item>
                <lay-button
                  type="primary"
                  :loading="loading"
                  @click="handleUpdateProfile"
                >
                  保存修改
                </lay-button>
              </lay-form-item>
            </lay-form>
          </div>
        </lay-tab-item>

        <!-- 修改密码 Tab -->
        <lay-tab-item title="修改密码" id="password">
          <div class="password-content">
            <lay-form class="password-form" label-position="top">
              <lay-form-item label="原密码" required>
                <lay-input
                  v-model="passwordForm.old_password"
                  type="password"
                  placeholder="请输入原密码"
                  :disabled="loading"
                />
              </lay-form-item>

              <lay-form-item label="新密码" required>
                <lay-input
                  v-model="passwordForm.new_password"
                  type="password"
                  placeholder="请输入新密码（至少6位）"
                  :disabled="loading"
                />
              </lay-form-item>

              <lay-form-item label="确认密码" required>
                <lay-input
                  v-model="passwordForm.confirm_password"
                  type="password"
                  placeholder="请再次输入新密码"
                  :disabled="loading"
                />
              </lay-form-item>

              <lay-form-item>
                <lay-button
                  type="primary"
                  :loading="loading"
                  @click="handleChangePassword"
                >
                  确认修改
                </lay-button>
                <lay-button @click="passwordForm = { old_password: '', new_password: '', confirm_password: '' }">
                  重置
                </lay-button>
              </lay-form-item>
            </lay-form>
          </div>
        </lay-tab-item>
      </lay-tab>
    </lay-card>
  </div>
</template>

<style scoped>
.profile-container {
  padding: 16px;
  height: 100%;
}

.profile-card {
  height: 100%;
  min-height: calc(100vh - 80px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--erp-text-primary);
}

.profile-content {
  padding: 20px 0;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 30px;
  background: linear-gradient(135deg, var(--erp-primary) 0%, #04d4d3 100%);
  border-radius: 8px;
  margin-bottom: 30px;
}

.user-info {
  flex: 1;
}

.user-name {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #fff;
}

.user-role {
  margin: 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.9);
}

.info-list {
  background-color: var(--erp-bg-card);
  border: 1px solid var(--erp-border-color);
  border-radius: 8px;
  overflow: hidden;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--erp-border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  width: 120px;
  font-size: 14px;
  color: var(--erp-text-secondary);
  flex-shrink: 0;
}

.info-value {
  flex: 1;
  font-size: 14px;
  color: var(--erp-text-primary);
}

.edit-content,
.password-content {
  padding: 20px 0;
}

.edit-form,
.password-form {
  max-width: 600px;
}

:deep(.layui-form-item) {
  margin-bottom: 20px;
}

:deep(.layui-form-item:last-child) {
  margin-bottom: 0;
  margin-top: 30px;
}
</style>
