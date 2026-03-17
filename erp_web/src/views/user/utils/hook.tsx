import { ref, reactive, onMounted } from 'vue'
import { layer } from '@layui/layui-vue'
import { userService, type UserInfo } from '@/api/userService'
import type { SearchForm, Pagination } from './types'

/** 用户管理页面状态管理 */
export const useUserList = () => {
  const searchForm = reactive<SearchForm>({ keyword: '' })
  const pagination = reactive<Pagination>({ page: 1, pageSize: 10, total: 0 })
  const tableData = ref<UserInfo[]>([])
  const loading = ref(false)

  const editVisible = ref(false)
  const editUserId = ref<number | undefined>(undefined)

  const fetchData = async () => {
    loading.value = true
    try {
      const res = await userService.adminGetList({
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

  const handleSearch = () => {
    pagination.page = 1
    fetchData()
  }

  const handleReset = () => {
    searchForm.keyword = ''
    pagination.page = 1
    fetchData()
  }

  const handlePageChange = (page: number) => {
    pagination.page = page
    fetchData()
  }

  const handleSizeChange = (size: number) => {
    pagination.pageSize = size
    pagination.page = 1
    fetchData()
  }

  const handleEdit = (row: UserInfo) => {
    editUserId.value = row.id
    editVisible.value = true
  }

  const handleFormSuccess = () => {
    fetchData()
  }

  const handleToggleStatus = async (row: UserInfo) => {
    try {
      if (row.status === 1) {
        await userService.adminBan(row.id)
        layer.msg('封禁成功', { icon: 1 })
      } else {
        await userService.adminUnban(row.id)
        layer.msg('解封成功', { icon: 1 })
      }
      fetchData()
    } catch { /* 错误由 request.ts 统一处理 */ }
  }

  const handleResetPassword = (row: UserInfo) => {
    layer.confirm(`确定要将用户【${row.username}】的密码重置为 u123456 吗？`, {
      title: '确认重置密码',
      btn: [
        { text: '确定', callback: async () => {
          try {
            await userService.adminResetPassword(row.id)
            layer.msg('重置密码成功，新密码为 u123456', { icon: 1 })
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
    editVisible,
    editUserId,
    handleEdit,
    handleFormSuccess,
    handlePageChange,
    handleSizeChange,
    handleToggleStatus,
    handleResetPassword
  }
}
