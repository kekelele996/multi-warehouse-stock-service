<template>
  <div class="auth-wrap">
    <el-card class="auth-card">
      <h2>登录</h2>
      <el-form :model="form" label-width="70px" @keyup.enter="submit">
        <el-form-item label="手机号">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/authStore'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const form = reactive({ phone: '13800000002', password: 'User@123' })

async function submit() {
  loading.value = true
  try {
    await auth.login(form.phone, form.password)
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-wrap { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f0f2f5; }
.auth-card { width: 380px; }
</style>
