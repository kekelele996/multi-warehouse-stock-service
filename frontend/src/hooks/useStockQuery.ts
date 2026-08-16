import { reactive, ref } from 'vue'
import { queryStock } from '@/api/stockRecord'
import type { StockRecord } from '@/types'

export function useStockQuery() {
  const loading = ref(false)
  const list = ref<StockRecord[]>([])
  const total = ref(0)
  const filters = reactive<{ product_id?: number; warehouse_id?: number; shelf_id?: number }>({})

  async function load(params: { page?: number; page_size?: number } = {}) {
    loading.value = true
    try {
      const res: any = await queryStock({ ...filters, ...params })
      list.value = res.data.list
      total.value = res.data.total
    } finally {
      loading.value = false
    }
  }

  return { loading, list, total, filters, load }
}
