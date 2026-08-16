import { defineStore } from 'pinia'
import { listOrders } from '@/api/stockOrder'
import type { StockOrder } from '@/types'

export const useOrderStore = defineStore('order', {
  state: () => ({ list: [] as StockOrder[], total: 0 }),
  actions: {
    async fetchList(params: Record<string, unknown> = {}) {
      const res: any = await listOrders(params)
      this.list = res.data.list
      this.total = res.data.total
    },
  },
})
