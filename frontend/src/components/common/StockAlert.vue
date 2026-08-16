<template>
  <el-tag :type="low ? 'danger' : 'success'" size="small">
    {{ low ? '库存预警' : '库存正常' }}
  </el-tag>
</template>

<script setup lang="ts">
defineProps<{ low: boolean }>()
</script>
