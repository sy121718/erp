<script setup lang="ts">
import { useUserList } from './utils/hook'
import UserTable from './components/UserTable.vue'
import UserEditForm from './components/UserEditForm.vue'

const {
  searchForm,
  handleSearch,
  handleReset,
  pagination,
  tableData,
  loading,
  editVisible,
  editUserId,
  handleEdit,
  handleFormSuccess,
  handlePageChange,
  handleSizeChange,
  handleToggleStatus,
  handleResetPassword
} = useUserList()
</script>

<template>
  <div class="user-container">
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
      <UserTable
        :data="tableData"
        :loading="loading"
        :pagination="pagination"
        @edit="handleEdit"
        @toggle-status="handleToggleStatus"
        @reset-password="handleResetPassword"
        @page-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </lay-card>

    <!-- 编辑弹窗 -->
    <UserEditForm
      v-model:visible="editVisible"
      :user-id="editUserId"
      @success="handleFormSuccess"
    />
  </div>
</template>

<style scoped>
.user-container {
  padding: 16px;
}
</style>
