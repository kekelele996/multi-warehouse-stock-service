<template>
  <el-result v-if="hasError" icon="error" title="页面出错了" :sub-title="message">
    <template #extra>
      <el-button type="primary" @click="reload">刷新页面</el-button>
    </template>
  </el-result>
  <slot v-else />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onErrorCaptured } from 'vue'

const hasError = ref(false)
const message = ref('')

onErrorCaptured((err) => {
  hasError.value = true
  message.value = err.message
  return false
})

function reload() {
  window.location.reload()
}
</script>
