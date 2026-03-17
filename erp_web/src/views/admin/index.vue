<script setup lang="ts">
import { useAdminList } from './utils/hook'
import AdminTable from './components/AdminTable.vue'
import AdminAddForm from './components/AdminAddForm.vue'
import AdminEditForm from './components/AdminEditForm.vue'

const {
  searchForm,
  handleSearch,
  handleReset,
  pagination,
  tableData,
  loading,
  addVisible,
  editVisible,
  editAdminId,
  handleAdd,
  handleEdit,
  handleFormSuccess,
  handlePageChange,
  handleSizeChange,
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

      <AdminTable
        :data="tableData"
        :loading="loading"
        :pagination="pagination"
        @edit="handleEdit"
        @delete="handleDelete"
        @toggle-status="handleToggleStatus"
        @reset-password="handleResetPassword"
        @force-logout="handleForceLogout"
        @page-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </lay-card>

    <!-- 新增弹窗 -->
    <AdminAddForm
      v-model:visible="addVisible"
      @success="handleFormSuccess"
    />

    <!-- 编辑弹窗 -->
    <AdminEditForm
      v-model:visible="editVisible"
      :admin-id="editAdminId"
      @success="handleFormSuccess"
    />
  </div>
</template>

<style scoped>
.admin-container {
  padding: 16px;
}
</style>
