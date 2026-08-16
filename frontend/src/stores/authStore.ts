import { defineStore } from 'pinia'
import { login as apiLogin, getMe } from '@/api/user'
import type { User } from '@/types'

const TOKEN_KEY = 'warehouse_token'
const USER_KEY = 'warehouse_user'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || '',
    user: JSON.parse(localStorage.getItem(USER_KEY) || 'null') as User | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    role: (state) => state.user?.role || '',
  },
  actions: {
    async login(phone: string, password: string) {
      const res: any = await apiLogin({ phone, password })
      this.token = res.data.token
      this.user = res.data.user
      localStorage.setItem(TOKEN_KEY, this.token)
      localStorage.setItem(USER_KEY, JSON.stringify(this.user))
    },
    async fetchMe() {
      const res: any = await getMe()
      this.user = res.data
      localStorage.setItem(USER_KEY, JSON.stringify(this.user))
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem(USER_KEY)
    },
  },
})
