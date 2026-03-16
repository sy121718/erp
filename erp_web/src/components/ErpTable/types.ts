/**
 * ErpTable 组件类型定义
 */

/** 表格列配置 */
export interface TableColumn {
  /** 标题 */
  title: string
  /** 字段名 */
  key?: string
  /** 宽度 */
  width?: string
  /** 对齐方式 */
  align?: 'left' | 'center' | 'right'
  /** 固定列 */
  fixed?: 'left' | 'right'
  /** 自定义插槽名称 */
  slot?: string
  /** 列类型 */
  type?: string
}

/** 分页配置 */
export interface TablePagination {
  /** 当前页 */
  page: number
  /** 每页条数 */
  limit: number
  /** 总条数 */
  total: number
}

/** 表格组件 Props */
export interface ErpTableProps {
  /** 列配置 */
  columns: TableColumn[]
  /** 数据源 */
  dataSource: any[]
  /** 加载状态 */
  loading?: boolean
  /** 分页配置 */
  pagination?: TablePagination
  /** 显示选择列 */
  showSelection?: boolean
  /** 显示序号列 */
  showIndex?: boolean
  /** 行唯一标识字段 */
  rowKey?: string
  /** 空数据提示 */
  emptyText?: string
  /** 表格高度 */
  height?: string | number
  /** 斑马纹 */
  stripe?: boolean
  /** 边框 */
  border?: boolean
}

/** 表格组件事件 */
export interface ErpTableEmits {
  /** 选中行变化 */
  (e: 'selectionChange', selectedKeys: (string | number)[]): void
  /** 页码变化 */
  (e: 'pageChange', page: number): void
  /** 每页条数变化 */
  (e: 'limitChange', limit: number): void
  /** 行点击 */
  (e: 'rowClick', row: any): void
}