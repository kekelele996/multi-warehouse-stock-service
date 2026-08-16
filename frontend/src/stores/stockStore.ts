import { defineStore } from 'pinia'
import { queryStock } from '@/api/stockRecord'
import type { StockRecord } from '@/types'

export const useStockStore = defineStore('stock', {
  state: () => ({ list: [] as StockRecord[], total: 0 }),
  actions: {
    async fetchList(params: Record<string, unknown> = {}) {
      const res: any = await queryStock(params)
      this.list = res.data.list
      this.total = res.data.total
    },
  },
})
