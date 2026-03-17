import { ref, reactive, onMounted } from 'vue'
import { layer } from '@layui/layui-vue'
import { adminService, type AdminInfo } from '@/api/adminService'
import type { SearchForm, Pagination, AdminFormData, FormMode } from './types'

/** 管理员列表组合式函数 */
export const useAdminList = () => {
  // 搜索表单
  const searchForm = reactive<SearchForm>({
    keyword: ''
  })

  // 分页
  const pagination = reactive<Pagination>({
    page: 1,
    pageSize: 10,
    total: 0
  })

  // 表格数据
  const tableData = ref<AdminInfo[]>([])
  const loading = ref(false)

  // 弹窗控制
  const formVisible = ref(false)
  const formMode = ref<FormMode>('create')
  const formLoading = ref(false)
  const formData = reactive<AdminFormData>({
    id: undefined,
    username: '',
    password: '',
    name: '',
    email: '',
    phone: '',
    status: undefined
  })

  // 表格列配置
  const columns = [
    { title: 'ID', key: 'id', width: 80, fixed: 'left' },
    { title: '用户名', key: 'username', width: 120 },
    { title: '姓名', key: 'name', width: 120 },
    { title: '邮箱', key: 'email', width: 180 },
    { title: '手机号', key: 'phone', width: 140 },
    { title: '状态', key: 'status', width: 100, customSlot: 'status' },
    { title: '最后登录时间', key: 'last_login_time', width: 180 },
    { title: '操作', key: 'action', width: 240, fixed: 'right', customSlot: 'action' }
  ]

  // 获取列表数据
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
    } catch (error) {
      // 错误由 request.ts 统一处理
    } finally {
      loading.value = false
    }
  }

  // 搜索
  const handleSearch = () => {
    pagination.page = 1
    fetchData()
  }

  // 重置搜索
  const handleReset = () => {
    searchForm.keyword = ''
    pagination.page = 1
    fetchData()
  }

  // 分页改变
  const handlePageChange = (page: number) => {
    pagination.page = page
    fetchData()
  }

  // 每页数量改变
  const handleSizeChange = (size: number) => {
    pagination.pageSize = size
    pagination.page = 1
    fetchData()
  }

  // 新增
  const handleAdd = () => {
    formMode.value = 'create'
    Object.assign(formData, {
      id: undefined,
      username: '',
      password: '',
      name: '',
      email: '',
      phone: '',
      status: undefined
    })
    formVisible.value = true
  }

  // 编辑
  const handleEdit = async (row: AdminInfo) => {
    formMode.value = 'edit'
    formLoading.value = true
    formVisible.value = true
    try {
      const detail = await adminService.getDetail(row.id)
      Object.assign(formData, {
        id: detail.id,
        username: detail.username,
        password: '',
        name: detail.name,
        email: detail.email || '',
        phone: detail.phone || '',
        status: detail.status
      })
    } catch (error) {
      formVisible.value = false
    } finally {
      formLoading.value = false
    }
  }

  // 提交表单
  const handleSubmit = async () => {
    // 表单验证
    if (!formData.username) {
      layer.msg('请输入用户名', { icon: 0 })
      return
    }
    if (formMode.value === 'create' && !formData.password) {
      layer.msg('请输入密码', { icon: 0 })
      return
    }
    if (!formData.name) {
      layer.msg('请输入姓名', { icon: 0 })
      return
    }

    formLoading.value = true
    try {
      if (formMode.value === 'create') {
        await adminService.create({
          username: formData.username,
          password: formData.password,
          name: formData.name,
          email: formData.email || undefined,
          phone: formData.phone || undefined
        })
        layer.msg('创建成功', { icon: 1 })
      } else if (formMode.value === 'edit' && formData.id) {
        await adminService.update(formData.id, {
          name: formData.name,
          email: formData.email || undefined,
          phone: formData.phone || undefined
        })
        layer.msg('更新成功', { icon: 1 })
      }
      formVisible.value = false
      fetchData()
    } catch (error) {
      // 错误由 request.ts 统一处理
    } finally {
      formLoading.value = false
    }
  }

  // 删除
  const handleDelete = (row: AdminInfo) => {
    layer.confirm(`确定要删除管理员【${row.username}】吗？`, {
      title: '确认删除',
      btn: [
        { text: '确定', callback: async () => {
          try {
            await adminService.delete(row.id)
            layer.msg('删除成功', { icon: 1 })
            fetchData()
          } catch (error) {
            // 错误由 request.ts 统一处理
          }
        }},
        { text: '取消', callback: () => true }
      ]
    })
  }

  // 禁用/解禁
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
    } catch (error) {
      // 错误由 request.ts 统一处理
    }
  }

  // 重置密码
  const handleResetPassword = (row: AdminInfo) => {
    layer.prompt({
      title: '重置密码',
      formType: 1,
      value: '',
      placeholder: '请输入新密码（6-50位）'
    }, async (value: string) => {
      if (value.length < 6 || value.length > 50) {
        layer.msg('密码长度为6-50位', { icon: 0 })
        return false
      }
      try {
        await adminService.resetPassword(row.id, { new_password: value })
        layer.msg('重置密码成功', { icon: 1 })
        return true
      } catch (error) {
        // 错误由 request.ts 统一处理
        return false
      }
    })
  }

  // 强制下线
  const handleForceLogout = (row: AdminInfo) => {
    layer.confirm(`确定要强制管理员【${row.username}】下线吗？`, {
      title: '确认操作',
      btn: [
        { text: '确定', callback: async () => {
          try {
            await adminService.forceLogout(row.id)
            layer.msg('强制下线成功', { icon: 1 })
          } catch (error) {
            // 错误由 request.ts 统一处理
          }
        }},
        { text: '取消', callback: () => true }
      ]
    })
  }

  // 初始化
  onMounted(() => {
    fetchData()
  })

  return {
    // 搜索
    searchForm,
    handleSearch,
    handleReset,

    // 分页
    pagination,
    handlePageChange,
    handleSizeChange,

    // 表格
    columns,
    tableData,
    loading,

    // 弹窗
    formVisible,
    formMode,
    formData,
    formLoading,

    // 操作
    handleAdd,
    handleEdit,
    handleSubmit,
    handleDelete,
    handleToggleStatus,
    handleResetPassword,
    handleForceLogout
  }
}
