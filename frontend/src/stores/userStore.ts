import { defineStore } from 'pinia'
import { getMe, updateProfile } from '@/api/user'
import type { User } from '@/types'

export const useUserStore = defineStore('user', {
  state: () => ({ me: null as User | null }),
  actions: {
    async fetchMe() {
      const res: any = await getMe()
      this.me = res.data
      return res.data as User
    },
    async updateMe(data: { name?: string; avatar?: string }) {
      const res: any = await updateProfile(data)
      this.me = res.data
      return res.data as User
    },
  },
})
