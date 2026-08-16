import { defineStore } from 'pinia'
import { listProducts, listCategories } from '@/api/product'
import type { Category, Product } from '@/types'

export const useProductStore = defineStore('product', {
  state: () => ({
    list: [] as Product[],
    total: 0,
    categories: [] as Category[],
  }),
  actions: {
    async fetchList(params: Record<string, unknown> = {}) {
      const res: any = await listProducts(params)
      this.list = res.data.list
      this.total = res.data.total
    },
    async fetchCategories() {
      const res: any = await listCategories()
      this.categories = res.data
    },
  },
})
