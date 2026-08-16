<template>
  <div>
    <div class="toolbar">
      <el-select v-model="categoryFilter" placeholder="分类" clearable style="width: 160px" @change="load">
        <el-option v-for="c in store.categories" :key="c.id" :label="c.name" :value="c.id" />
      </el-select>
      <el-input v-model="keyword" placeholder="搜索名称/SKU/条码" clearable style="width: 220px" @change="load" />
      <el-button type="primary" @click="openCreate">新建商品</el-button>
    </div>
    <el-table :data="store.list" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column prop="sku" label="SKU" width="120" />
      <el-table-column label="分类" width="120">
        <template #default="{ row }">{{ categoryName(row.category_id) }}</template>
      </el-table-column>
      <el-table-column prop="spec" label="规格" width="100" />
      <el-table-column prop="unit" label="单位" width="70" />
      <el-table-column prop="min_stock" label="最低库存" width="90" />
      <el-table-column prop="max_stock" label="最高库存" width="90" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" @click="edit(row)">编辑</el-button>
          <el-button size="small" @click="showStock(row)">库存分布</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="mt" layout="total, prev, pager, next" :total="store.total" :page-size="pageSize" :current-page="page" @current-change="onPage" />

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑商品' : '新建商品'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="SKU"><el-input v-model="form.sku" :disabled="!!editing" /></el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.category_id" style="width: 100%">
            <el-option v-for="c in store.categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="规格"><el-input v-model="form.spec" /></el-form-item>
        <el-form-item label="单位"><el-input v-model="form.unit" /></el-form-item>
        <el-form-item label="最低库存"><el-input-number v-model="form.min_stock" :min="0" /></el-form-item>
        <el-form-item label="最高库存"><el-input-number v-model="form.max_stock" :min="0" /></el-form-item>
        <el-button type="primary" @click="save">保存</el-button>
      </el-form>
    </el-dialog>

    <el-dialog v-model="stockVisible" title="各仓库库存分布" width="560px">
      <el-table :data="stockList" border size="small">
        <el-table-column prop="warehouse_id" label="仓库ID" width="90" />
        <el-table-column prop="shelf_id" label="库位ID" width="90" />
        <el-table-column prop="batch_no" label="批次" />
        <el-table-column prop="quantity" label="数量" width="90" />
        <el-table-column label="预警">
          <template #default="{ row }">
            <StockAlert :low="Number(row.quantity) < (editing?.min_stock || 0)" />
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import StockAlert from '@/components/common/StockAlert.vue'
import { useProductStore } from '@/stores/productStore'
import { createProduct, updateProduct } from '@/api/product'
import { queryStock } from '@/api/stockRecord'
import type { Product, StockRecord } from '@/types'

const store = useProductStore()
const loading = ref(false)
const page = ref(1)
const pageSize = 10
const categoryFilter = ref<number>()
const keyword = ref('')
const dialogVisible = ref(false)
const stockVisible = ref(false)
const editing = ref<Product | null>(null)
const stockList = ref<StockRecord[]>([])
const form = ref({ name: '', sku: '', category_id: 0, spec: '', unit: '件', min_stock: 0, max_stock: 0 })

async function load() {
  loading.value = true
  try {
    await store.fetchList({ page: page.value, page_size: pageSize, category_id: categoryFilter.value, keyword: keyword.value || undefined })
  } finally {
    loading.value = false
  }
}
function onPage(p: number) {
  page.value = p
  load()
}
function categoryName(id: number): string {
  return store.categories.find((c) => c.id === id)?.name || `#${id}`
}
function openCreate() {
  editing.value = null
  form.value = { name: '', sku: '', category_id: store.categories[0]?.id || 0, spec: '', unit: '件', min_stock: 0, max_stock: 0 }
  dialogVisible.value = true
}
function edit(row: Product) {
  editing.value = row
  form.value = { name: row.name, sku: row.sku, category_id: row.category_id, spec: row.spec, unit: row.unit, min_stock: row.min_stock, max_stock: row.max_stock }
  dialogVisible.value = true
}
async function save() {
  if (editing.value) {
    await updateProduct(editing.value.id, { ...form.value })
  } else {
    await createProduct({ ...form.value })
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  load()
}
async function showStock(row: Product) {
  const res: any = await queryStock({ product_id: row.id, page: 1, page_size: 200 })
  stockList.value = res.data.list
  editing.value = row
  stockVisible.value = true
}
onMounted(async () => {
  await store.fetchCategories()
  await load()
})
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.mt { margin-top: 12px; }
</style>
