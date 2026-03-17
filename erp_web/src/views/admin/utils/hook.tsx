import { ref, reactive, onMounted } from 'vue'
import { layer } from '@layui/layui-vue'
import { adminService, type AdminInfo } from '@/api/adminService'
import type { SearchForm, Pagination } from './types'

/** 管理员列表页面状态管理 */
export const useAdminList = () => {
  const searchForm = reactive<SearchForm>({ keyword: '' })
  const pagination = reactive<Pagination>({ page: 1, pageSize: 10, total: 0 })
  const tableData = ref<AdminInfo[]>([])
  const loading = ref(false)

  // 弹窗状态
  const addVisible = ref(false)
  const editVisible = ref(false)
  const editAdminId = ref<number | undefined>(undefined)

  // ========== 数据获取 ==========

  const fetchData = async () => {
    loading.value = true
    try {
      const res = await adminService.getList({
        page: pagination.page,
        page_size: pagination.pageSize,
        keyword: searchForm.keyword
      })
      tableData.value = res.list
      pagination.total = res.total
    } finally {
      loading.value = false
    }
  }

  // ========== 搜索 ==========

  const handleSearch = () => {
    pagination.page = 1
    fetchData()
  }

  const handleReset = () => {
    searchForm.keyword = ''
    pagination.page = 1
    fetchData()
  }

  // ========== 分页 ==========

  const handlePageChange = (page: number) => {
    pagination.page = page
    fetchData()
  }

  const handleSizeChange = (size: number) => {
    pagination.pageSize = size
    pagination.page = 1
    fetchData()
  }

  // ========== 弹窗控制 ==========

  const handleAdd = () => {
    addVisible.value = true
  }

  const handleEdit = (row: AdminInfo) => {
    editAdminId.value = row.id
    editVisible.value = true
  }

  const handleFormSuccess = () => {
    fetchData()
  }

  // ========== 行操作 ==========

  const handleDelete = (row: AdminInfo) => {
    layer.confirm(`确定要删除管理员【${row.username}】吗？`, {
      title: '确认删除',
      btn: [
        { text: '确定', callback: async () => {
          try {
            await adminService.delete(row.id)
            layer.msg('删除成功', { icon: 1 })
            fetchData()
          } catch { /* 错误由 request.ts 统一处理 */ }
        }},
        { text: '取消', callback: () => true }
      ]
    })
  }

  const handleToggleStatus = async (row: AdminInfo) => {
    try {
      if (row.status === 1) {
        await adminService.ban(row.id)
        layer.msg('禁用成功', { icon: 1 })
      } else {
        await adminService.unban(row.id)
        layer.msg('解禁成功', { icon: 1 })
      }
      fetchData()
    } catch { /* 错误由 request.ts 统一处理 */ }
  }

  const handleResetPassword = (row: AdminInfo) => {
    layer.prompt({
      title: '重置密码',
      formType: 1,
      value: '',
      placeholder: '请输入新密码（6-50位）',
      yes: async (value: string) => {
        if (value.length < 6 || value.length > 50) {
          layer.msg('密码长度为6-50位', { icon: 0 })
          return false
        }
        try {
          await adminService.resetPassword(row.id, { new_password: value })
          layer.msg('重置密码成功', { icon: 1 })
          return true
        } catch {
          return false
        }
      }
    })
  }

  const handleForceLogout = (row: AdminInfo) => {
    layer.confirm(`确定要强制管理员【${row.username}】下线吗？`, {
      title: '确认操作',
      btn: [
        { text: '确定', callback: async () => {
          try {
            await adminService.forceLogout(row.id)
            layer.msg('强制下线成功', { icon: 1 })
          } catch { /* 错误由 request.ts 统一处理 */ }
        }},
        { text: '取消', callback: () => true }
      ]
    })
  }

  onMounted(() => fetchData())

  return {
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
  }
}
