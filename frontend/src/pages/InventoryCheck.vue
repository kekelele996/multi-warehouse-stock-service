<template>
  <div>
    <el-card>
      <el-form inline>
        <el-form-item label="仓库">
          <el-select v-model="warehouseId" style="width: 180px" @change="loadRecords">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品">
          <el-select v-model="productId" clearable filterable style="width: 200px" @change="loadRecords">
            <el-option v-for="p in products" :key="p.id" :label="`${p.name} (${p.sku})`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadRecords">查询</el-button>
      </el-form>
    </el-card>
    <el-card class="mt">
      <el-table :data="records" border>
        <el-table-column prop="product_id" label="商品ID" width="90" />
        <el-table-column prop="shelf_id" label="库位ID" width="90" />
        <el-table-column prop="batch_no" label="批次" />
        <el-table-column prop="quantity" label="账面数量" width="100" />
        <el-table-column label="实盘数量" width="140">
          <template #default="{ row }">
            <el-input-number v-model="actual[row.id]" :min="0" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="差异" width="100">
          <template #default="{ row }">
            <span :style="{ color: (actual[row.id] - row.quantity) !== 0 ? '#f56c6c' : '#67c23a' }">
              {{ (actual[row.id] ?? row.quantity) - row.quantity }}
            </span>
          </template>
        </el-table-column>
      </el-table>
      <EmptyState v-if="records.length === 0" text="暂无库存记录" />
      <RoleGuard :roles="['admin', 'warehouse_manager', 'operator']">
        <el-button type="primary" class="mt" :loading="saving" @click="saveCheck">保存盘点结果</el-button>
      </RoleGuard>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import EmptyState from '@/components/common/EmptyState.vue'
import RoleGuard from '@/components/common/RoleGuard.vue'
import { useWarehouseStore } from '@/stores/warehouseStore'
import { useProductStore } from '@/stores/productStore'
import { queryStock, adjustStock } from '@/api/stockRecord'
import type { StockRecord } from '@/types'

const warehouses = ref<any[]>([])
const products = ref<any[]>([])
const records = ref<StockRecord[]>([])
const actual = ref<Record<number, number>>({})
const warehouseId = ref<number>()
const productId = ref<number>()
const saving = ref(false)
const warehouseStore = useWarehouseStore()
const productStore = useProductStore()

async function loadRecords() {
  const res: any = await queryStock({ warehouse_id: warehouseId.value, product_id: productId.value, page: 1, page_size: 200 })
  records.value = res.data.list
  actual.value = {}
  records.value.forEach((r) => {
    actual.value[r.id] = r.quantity
  })
}

async function saveCheck() {
  saving.value = true
  try {
    for (const r of records.value) {
      await adjustStock({ product_id: r.product_id, warehouse_id: r.warehouse_id, shelf_id: r.shelf_id, actual_quantity: actual.value[r.id] })
    }
    ElMessage.success('盘点结果已保存')
    await loadRecords()
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  const wRes: any = await warehouseStore.fetchList({ page: 1, page_size: 100 })
  warehouses.value = warehouseStore.list
  if (warehouses.value.length > 0) warehouseId.value = warehouses.value[0].id
  await productStore.fetchList({ page: 1, page_size: 200 })
  products.value = productStore.list
  await loadRecords()
})
</script>

<style scoped>
.mt { margin-top: 16px; }
</style>
