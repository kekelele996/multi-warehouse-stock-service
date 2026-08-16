<template>
  <el-tag :type="tagType">{{ text }}</el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { OrderStatusText } from '@/constants/order'

const props = defineProps<{ status: string }>()

const text = computed(() => OrderStatusText[props.status] || props.status)
const tagType = computed(() => {
  if (['completed', 'active'].includes(props.status)) return 'success'
  if (['cancelled', 'inactive', 'failed'].includes(props.status)) return 'info'
  if (['submitted', 'pending'].includes(props.status)) return 'warning'
  if (['processing', 'maintenance'].includes(props.status)) return 'primary'
  return 'default'
})
</script>
