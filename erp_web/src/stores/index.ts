import { createPinia } from 'pinia'

const pinia = createPinia()

export default pinia

// 导出各个 Store
export * from './user'