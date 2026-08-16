<template>
  <div>
    <el-row :gutter="16" class="mb">
      <el-col :span="6"><el-card><el-statistic title="仓库数量" :value="stats.warehouse_count" /></el-card></el-col>
      <el-col :span="6"><el-card><el-statistic title="今日出入库单" :value="stats.today_orders" /></el-card></el-col>
      <el-col :span="12"><el-card><el-statistic title="库存预警商品" :value="stats.low_stock?.length || 0" /></el-card></el-col>
    </el-row>
    <el-row :gutter="16">
      <el-col :span="12">
        <el-card>
          <div ref="barRef" class="chart" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card title="库存预警 TOP10">
          <el-table :data="stats.low_stock || []" size="small" max-height="320">
            <el-table-column prop="name" label="商品" />
            <el-table-column prop="sku" label="SKU" />
            <el-table-column prop="quantity" label="当前库存" />
            <el-table-column prop="min_stock" label="最低库存" />
            <el-table-column label="状态">
              <template #default="{ row }">
                <StockAlert :low="true" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, reactive, ref } from 'vue'
import * as echarts from 'echarts'
import { getDashboardStats } from '@/api/dashboard'
import StockAlert from '@/components/common/StockAlert.vue'

const stats = reactive<any>({ warehouse_count: 0, today_orders: 0, low_stock: [], warehouse_summary: [] })
const barRef = ref<HTMLElement>()

onMounted(async () => {
  const res: any = await getDashboardStats()
  Object.assign(stats, res.data)
  await nextTick()
  renderBar()
})

function renderBar() {
  if (!barRef.value) return
  const chart = echarts.init(barRef.value)
  chart.setOption({
    title: { text: '各仓库库存量' },
    tooltip: {},
    xAxis: { type: 'category', data: (stats.warehouse_summary || []).map((w: any) => `仓库#${w.warehouse_id}`) },
    yAxis: { type: 'value' },
    series: [{ name: '库存量', type: 'bar', data: (stats.warehouse_summary || []).map((w: any) => w.quantity) }],
  })
}
</script>

<style scoped>
.mb { margin-bottom: 16px; }
.chart { height: 320px; }
</style>
