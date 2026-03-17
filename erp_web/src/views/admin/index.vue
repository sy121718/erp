<script setup lang="ts">
import { useAdminList } from './utils/hook'

const {
  searchForm,
  handleSearch,
  handleReset,
  pagination,
  handlePageChange,
  handleSizeChange,
  columns,
  tableData,
  loading,
  formVisible,
  formMode,
  formData,
  formLoading,
  handleAdd,
  handleEdit,
  handleSubmit,
  handleDelete,
  handleToggleStatus,
  handleResetPassword,
  handleForceLogout
} = useAdminList()
</script>

<template>
  <div class="admin-container">
    <!-- 搜索区域 -->
    <lay-card class="mb-4">
      <lay-form :model="searchForm" :pane="true">
        <lay-row :space="16">
          <lay-col :md="8">
            <lay-form-item label="关键词">
              <lay-input
                v-model="searchForm.keyword"
                placeholder="用户名/姓名/邮箱/手机号"
                allow-clear
                @pressEnter="handleSearch"
              />
            </lay-form-item>
          </lay-col>
          <lay-col :md="8">
            <lay-form-item>
              <lay-button type="primary" @click="handleSearch">查询</lay-button>
              <lay-button @click="handleReset">重置</lay-button>
            </lay-form-item>
          </lay-col>
        </lay-row>
      </lay-form>
    </lay-card>

    <!-- 表格区域 -->
    <lay-card>
      <div class="flex justify-between mb-4">
        <lay-button type="primary" @click="handleAdd">新增管理员</lay-button>
      </div>

      <lay-table
        :columns="columns"
        :dataSource="tableData"
        :loading="loading"
        :default-toolbar="true"
        :page="{
          current: pagination.page,
          limit: pagination.pageSize,
          total: pagination.total
        }"
        @change="handlePageChange"
        @limitChange="handleSizeChange"
      >
        <!-- 状态 -->
        <template #status="{ row }">
          <lay-tag :type="row.status === 1 ? 'normal' : 'danger'">
            {{ row.status === 1 ? '正常' : '禁用' }}
          </lay-tag>
        </template>

        <!-- 是否超管 -->
        <template #is_admin="{ row }">
          <lay-tag :type="row.is_admin ? 'normal' : ''">
            {{ row.is_admin ? '是' : '否' }}
          </lay-tag>
        </template>

        <!-- 操作 -->
        <template #action="{ row }">
          <lay-button size="xs" type="primary" @click="handleEdit(row)">编辑</lay-button>
          <lay-button
            size="xs"
            :type="row.status === 1 ? 'danger' : 'success'"
            @click="handleToggleStatus(row)"
          >
            {{ row.status === 1 ? '禁用' : '解禁' }}
          </lay-button>
          <lay-button size="xs" type="warm" @click="handleResetPassword(row)">重置密码</lay-button>
          <lay-button size="xs" type="danger" @click="handleForceLogout(row)">强制下线</lay-button>
          <lay-button size="xs" type="danger" @click="handleDelete(row)">删除</lay-button>
        </template>
      </lay-table>
    </lay-card>

    <!-- 新增/编辑/查看弹窗 -->
    <lay-layer
      v-model="formVisible"
      :title="formMode === 'create' ? '新增管理员' : '编辑管理员'"
      :area="['600px', 'auto']"
      :btn="[
        { text: '取消', callback: () => formVisible = false },
        { text: '确定', callback: handleSubmit }
      ]"
    >
      <lay-form :model="formData" :pane="true" class="p-4">
        <lay-form-item label="用户名" required>
          <lay-input
            v-model="formData.username"
            placeholder="请输入用户名（3-50位）"
            :disabled="formMode !== 'create'"
            allow-clear
          />
        </lay-form-item>

        <lay-form-item v-if="formMode === 'create'" label="密码" required>
          <lay-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码（6-50位）"
            allow-clear
          />
        </lay-form-item>

        <lay-form-item label="姓名" required>
          <lay-input
            v-model="formData.name"
            placeholder="请输入姓名"
            allow-clear
          />
        </lay-form-item>

        <lay-form-item label="邮箱">
          <lay-input
            v-model="formData.email"
            placeholder="请输入邮箱"
            allow-clear
          />
        </lay-form-item>

        <lay-form-item label="手机号">
          <lay-input
            v-model="formData.phone"
            placeholder="请输入手机号"
            allow-clear
          />
        </lay-form-item>
      </lay-form>
    </lay-layer>
  </div>
</template>

<style scoped>
.admin-container {
  padding: 16px;
}

.mb-4 {
  margin-bottom: 16px;
}

.p-4 {
  padding: 16px;
}

.flex {
  display: flex;
}

.justify-between {
  justify-content: space-between;
}
</style>
