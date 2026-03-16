# Pinia Store 规范

## 目录结构

```
src/stores/
├── CLAUDE.md
├── index.ts        # Pinia 实例
├── user.ts         # 用户状态
├── app.ts          # 应用状态
└── ...
```

## Store 定义

```typescript
// src/stores/user.ts
import { defineStore } from 'pinia'
import { userService, type User } from '@/api'

interface UserState {
  userInfo: User | null
  token: string
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    userInfo: null,
    token: localStorage.getItem('token') || ''
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    nickname: (state) => state.userInfo?.nickname || '未登录'
  },

  actions: {
    async login(params: LoginParams) {
      const res = await userService.login(params)
      this.token = res.token
      this.userInfo = res.user
      localStorage.setItem('token', res.token)
    },

    logout() {
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
    },

    async getUserInfo() {
      const res = await userService.getInfo()
      this.userInfo = res
    }
  }
})
```

## 使用示例

```typescript
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// 读取状态
console.log(userStore.isLoggedIn)
console.log(userStore.nickname)

// 调用 action
await userStore.login({ username: 'admin', password: '123456' })
userStore.logout()
```

## 注意事项

1. Store 文件名使用 camelCase
2. 导出名称使用 use 前缀: `useUserStore`
3. 持久化数据使用 localStorage
4. 异步操作放在 actions 中