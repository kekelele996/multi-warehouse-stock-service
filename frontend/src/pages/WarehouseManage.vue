<template>
  <div>
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">新建仓库</el-button>
    </div>
    <el-row :gutter="16">
      <el-col v-for="w in store.list" :key="w.id" :span="8" class="mb">
        <el-card>
          <template #header>
            <div class="card-head">
              <span>{{ w.name }}（{{ w.code }}）</span>
              <el-tag size="small" :type="w.status === 'active' ? 'success' : 'info'">{{ statusText(w.status) }}</el-tag>
            </div>
          </template>
          <p>地址：{{ w.address }}</p>
          <p>面积：{{ w.area_sqm }} m²</p>
          <p>电话：{{ w.phone || '-' }}</p>
          <el-button size="small" @click="openDetail(w)">查看库位</el-button>
          <RoleGuard :roles="['admin', 'warehouse_manager']">
            <el-button size="small" @click="toggleStatus(w)">{{ w.status === 'active' ? '停用' : '启用' }}</el-button>
          </RoleGuard>
        </el-card>
      </el-col>
    </el-row>
    <el-pagination class="mt" layout="total, prev, pager, next" :total="store.total" :page-size="pageSize" :current-page="page" @current-change="onPage" />

    <el-dialog v-model="createVisible" title="新建仓库" width="480px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="编码"><el-input v-model="form.code" /></el-form-item>
        <el-form-item label="地址"><el-input v-model="form.address" /></el-form-item>
        <el-form-item label="面积"><el-input-number v-model="form.area_sqm" :min="0" /></el-form-item>
        <el-form-item label="电话"><el-input v-model="form.phone" /></el-form-item>
        <el-button type="primary" @click="create">提交</el-button>
      </el-form>
    </el-dialog>

    <el-dialog v-model="detailVisible" :title="detail?.name" width="640px">
      <el-table :data="shelves" size="small" border>
        <el-table-column prop="shelf_no" label="库位号" />
        <el-table-column prop="layer_count" label="层数" />
        <el-table-column prop="column_count" label="列数" />
        <el-table-column prop="capacity" label="容量" />
        <el-table-column label="占用率">
          <template #default="{ row }">
            <OccupancyBar :used="stockCountByShelf(row.id)" :capacity="row.capacity" />
          </template>
        </el-table-column>
      </el-table>
      <RoleGuard :roles="['admin', 'warehouse_manager']">
        <div class="mt">
          <el-input v-model="shelfForm.shelf_no" placeholder="库位号" style="width: 140px" />
          <el-button type="primary" @click="addShelf">添加库位</el-button>
        </div>
      </RoleGuard>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import OccupancyBar from '@/components/common/OccupancyBar.vue'
import RoleGuard from '@/components/common/RoleGuard.vue'
import { useWarehouseStore } from '@/stores/warehouseStore'
import { useStockStore } from '@/stores/stockStore'
import { createWarehouse, changeWarehouseStatus, getWarehouse, createShelf } from '@/api/warehouse'
import type { Shelf, Warehouse } from '@/types'

const store = useWarehouseStore()
const stockStore = useStockStore()
const page = ref(1)
const pageSize = 12
const createVisible = ref(false)
const detailVisible = ref(false)
const detail = ref<Warehouse | null>(null)
const shelves = ref<Shelf[]>([])
const form = ref({ name: '', code: '', address: '', area_sqm: 0, phone: '' })
const shelfForm = ref({ shelf_no: '' })

async function load() {
  await store.fetchList({ page: page.value, page_size: pageSize })
}
function onPage(p: number) {
  page.value = p
  load()
}
function openCreate() {
  form.value = { name: '', code: '', address: '', area_sqm: 0, phone: '' }
  createVisible.value = true
}
async function create() {
  await createWarehouse(form.value)
  ElMessage.success('仓库创建成功')
  createVisible.value = false
  load()
}
async function openDetail(w: Warehouse) {
  const res: any = await getWarehouse(w.id)
  detail.value = res.data.warehouse
  shelves.value = res.data.shelves
  await stockStore.fetchList({ warehouse_id: w.id, page: 1, page_size: 200 })
  detailVisible.value = true
}
async function addShelf() {
  if (!detail.value) return
  await createShelf({ warehouse_id: detail.value.id, shelf_no: shelfForm.value.shelf_no, layer_count: 3, column_count: 4, capacity: 120 })
  ElMessage.success('库位已添加')
  shelfForm.value.shelf_no = ''
  openDetail(detail.value)
}
async function toggleStatus(w: Warehouse) {
  const next = w.status === 'active' ? 'inactive' : 'active'
  await changeWarehouseStatus(w.id, next)
  ElMessage.success('状态已更新')
  load()
}
function stockCountByShelf(shelfId: number): number {
  return stockStore.list.filter((r) => r.shelf_id === shelfId).reduce((s, r) => s + r.quantity, 0)
}
function statusText(s: string): string {
  if (s === 'active') return '启用'
  if (s === 'inactive') return '停用'
  return '维护中'
}
onMounted(load)
</script>

<style scoped>
.toolbar { margin-bottom: 16px; }
.mb { margin-bottom: 16px; }
.card-head { display: flex; justify-content: space-between; align-items: center; }
.mt { margin-top: 12px; }
</style>
