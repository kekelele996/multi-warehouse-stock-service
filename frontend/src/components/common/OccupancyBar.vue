<template>
  <div>
    <el-progress :percentage="percentage" :color="color" />
    <span class="text">{{ percentage }}%</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ used: number; capacity: number }>()
const percentage = computed(() => {
  if (props.capacity <= 0) return 0
  return Math.min(100, Math.round((props.used / props.capacity) * 100))
})
const color = computed(() => {
  if (percentage.value >= 90) return '#f56c6c'
  if (percentage.value >= 70) return '#e6a23c'
  return '#67c23a'
})
</script>

<style scoped>
.text { margin-left: 8px; color: #909399; font-size: 12px; }
</style>
