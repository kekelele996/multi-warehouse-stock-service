import { defineStore } from 'pinia'
import { listWarehouses } from '@/api/warehouse'
import type { Warehouse } from '@/types'

export const useWarehouseStore = defineStore('warehouse', {
  state: () => ({ list: [] as Warehouse[], total: 0 }),
  actions: {
    async fetchList(params: { page?: number; page_size?: number } = {}) {
      const res: any = await listWarehouses(params)
      this.list = res.data.list
      this.total = res.data.total
    },
  },
})
