<template>
  <div>
    <div class="toolbar">
      <el-tabs v-model="typeTab" @tab-change="onTabChange">
        <el-tab-pane label="全部" name="" />
        <el-tab-pane v-for="opt in OrderTypeOptions" :key="opt.value" :label="opt.label" :name="opt.value" />
      </el-tabs>
      <RoleGuard :roles="['admin', 'warehouse_manager', 'operator']">
        <el-button type="primary" @click="openCreate">创建单据</el-button>
      </RoleGuard>
    </div>
    <el-table :data="store.list" v-loading="loading" border>
      <el-table-column prop="order_no" label="单号" width="170" />
      <el-table-column label="类型" width="90">
        <template #default="{ row }">{{ OrderTypeText[row.order_type] }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><StatusBadge :status="row.status" /></template>
      </el-table-column>
      <el-table-column prop="source_warehouse_id" label="源仓库" width="90" />
      <el-table-column prop="target_warehouse_id" label="目标仓库" width="90" />
      <el-table-column prop="remark" label="备注" />
      <el-table-column label="操作" width="300">
        <template #default="{ row }">
          <el-button size="small" @click="detail(row)">详情</el-button>
          <RoleGuard :roles="['admin', 'warehouse_manager', 'operator']">
            <el-button v-if="row.status === 'draft'" size="small" type="primary" @click="submit(row)">提交</el-button>
            <el-button v-if="row.status === 'submitted'" size="small" type="success" @click="approve(row)">审批</el-button>
            <el-button v-if="['submitted', 'processing'].includes(row.status)" size="small" @click="execute(row)">执行</el-button>
            <el-button v-if="['draft', 'submitted'].includes(row.status)" size="small" type="danger" @click="cancel(row)">取消</el-button>
          </RoleGuard>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="mt" layout="total, prev, pager, next" :total="store.total" :page-size="pageSize" :current-page="page" @current-change="onPage" />

    <el-dialog v-model="createVisible" title="创建出入库单" width="620px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="单据类型">
          <el-select v-model="form.order_type" style="width: 100%">
            <el-option v-for="opt in OrderTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="源仓库">
          <el-input-number v-model="form.source_warehouse_id" :min="0" />
        </el-form-item>
        <el-form-item label="目标仓库">
          <el-input-number v-model="form.target_warehouse_id" :min="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" />
        </el-form-item>
        <el-divider>明细行（商品ID / 库位ID / 数量）</el-divider>
        <div v-for="(item, idx) in form.items" :key="idx" class="item-row">
          <el-input-number v-model="item.product_id" :min="1" placeholder="商品ID" />
          <el-input-number v-model="item.shelf_id" :min="0" placeholder="库位ID" />
          <el-input-number v-model="item.quantity" :min="1" placeholder="数量" />
          <el-button type="danger" text @click="removeItem(idx)">删除</el-button>
        </div>
        <el-button size="small" @click="addItem">+ 添加明细</el-button>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="create">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" :title="detailOrder?.order_no" width="620px">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="类型">{{ detailOrder ? OrderTypeText[detailOrder.order_type] : '' }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ detailOrder ? OrderStatusText[detailOrder.status] : '' }}</el-descriptions-item>
        <el-descriptions-item label="源仓库">{{ detailOrder?.source_warehouse_id }}</el-descriptions-item>
        <el-descriptions-item label="目标仓库">{{ detailOrder?.target_warehouse_id }}</el-descriptions-item>
      </el-descriptions>
      <el-table :data="detailItems" border size="small" class="mt">
        <el-table-column prop="product_id" label="商品ID" />
        <el-table-column prop="shelf_id" label="库位ID" />
        <el-table-column prop="quantity" label="数量" />
        <el-table-column prop="actual_quantity" label="实收数量" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import StatusBadge from '@/components/common/StatusBadge.vue'
import RoleGuard from '@/components/common/RoleGuard.vue'
import { useOrderStore } from '@/stores/orderStore'
import { createOrder, submitOrder, approveOrder, executeOrder, cancelOrder, getOrder } from '@/api/stockOrder'
import { OrderStatusText, OrderTypeOptions, OrderTypeText } from '@/constants/order'
import { generateOrderNo } from '@/utils/orderNo'
import type { StockOrder, StockOrderItem } from '@/types'

const store = useOrderStore()
const loading = ref(false)
const page = ref(1)
const pageSize = 10
const typeTab = ref('')
const createVisible = ref(false)
const detailVisible = ref(false)
const detailOrder = ref<StockOrder | null>(null)
const detailItems = ref<StockOrderItem[]>([])
const form = ref({
  order_type: 'inbound',
  source_warehouse_id: 0,
  target_warehouse_id: 1,
  remark: '',
  items: [] as { product_id: number; shelf_id: number; quantity: number }[],
})

async function load() {
  loading.value = true
  try {
    await store.fetchList({ page: page.value, page_size: pageSize, order_type: typeTab.value || undefined })
  } finally {
    loading.value = false
  }
}
function onPage(p: number) {
  page.value = p
  load()
}
function onTabChange() {
  page.value = 1
  load()
}
function openCreate() {
  form.value = { order_type: 'inbound', source_warehouse_id: 0, target_warehouse_id: 1, remark: '', items: [{ product_id: 1, shelf_id: 1, quantity: 10 }] }
  createVisible.value = true
}
function addItem() {
  form.value.items.push({ product_id: 1, shelf_id: 1, quantity: 1 })
}
function removeItem(idx: number) {
  form.value.items.splice(idx, 1)
}
async function create() {
  await createOrder({ ...form.value, order_no: generateOrderNo(form.value.order_type) })
  ElMessage.success('单据创建成功')
  createVisible.value = false
  load()
}
async function submit(row: StockOrder) {
  await submitOrder(row.id)
  ElMessage.success('已提交')
  load()
}
async function approve(row: StockOrder) {
  await approveOrder(row.id)
  ElMessage.success('已审批')
  load()
}
async function execute(row: StockOrder) {
  await executeOrder(row.id)
  ElMessage.success('已执行')
  load()
}
async function cancel(row: StockOrder) {
  await cancelOrder(row.id)
  ElMessage.success('已取消')
  load()
}
async function detail(row: StockOrder) {
  const res: any = await getOrder(row.id)
  detailOrder.value = res.data.order
  detailItems.value = res.data.items
  detailVisible.value = true
}
onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; justify-content: space-between; align-items: flex-start; }
.mt { margin-top: 12px; }
.item-row { display: flex; gap: 8px; margin-bottom: 8px; align-items: center; }
</style>
