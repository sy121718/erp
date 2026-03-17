/** 搜索表单 */
export interface SearchForm {
  keyword: string
}

/** 分页参数 */
export interface Pagination {
  page: number
  pageSize: number
  total: number
}

/** 表格列配置 */
export interface TableColumn {
  title: string
  key: string
  width?: number | string
  fixed?: string
  customSlot?: string
}

/** 管理员表单数据（新增/编辑） */
export interface AdminFormData {
  id?: number
  username: string
  password: string
  name: string
  email: string
  phone: string
}

/** 表单模式 */
export type FormMode = 'create' | 'edit'
