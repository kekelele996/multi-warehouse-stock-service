<template>
  <el-card title="操作日志">
    <el-table :data="list" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="operator_name" label="操作人" />
      <el-table-column prop="action" label="动作" />
      <el-table-column prop="entity_type" label="实体" />
      <el-table-column prop="entity_id" label="实体ID" />
      <el-table-column prop="ip" label="IP" />
      <el-table-column prop="created_at" label="时间" />
    </el-table>
    <el-pagination class="mt" layout="total, prev, pager, next" :total="total" :page-size="pageSize" :current-page="page" @current-change="onPage" />
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { listOperationLogs } from '@/api/operationLog'
import type { OperationLog } from '@/types'

const list = ref<OperationLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const res: any = await listOperationLogs({ page: page.value, page_size: pageSize })
    list.value = res.data.list
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}
function onPage(p: number) {
  page.value = p
  load()
}
onMounted(load)
</script>

<style scoped>
.mt { margin-top: 12px; }
</style>
