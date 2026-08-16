<template>
  <el-container class="layout">
    <el-header class="header">
      <div class="brand" @click="router.push('/dashboard')">多仓库库存管理系统</div>
      <el-menu mode="horizontal" :default-active="active" router class="menu">
        <el-menu-item index="/dashboard">库存总览</el-menu-item>
        <el-menu-item index="/warehouses">仓库管理</el-menu-item>
        <el-menu-item index="/products">商品管理</el-menu-item>
        <el-menu-item index="/orders">出入库管理</el-menu-item>
        <el-menu-item index="/inventory-check">盘点管理</el-menu-item>
        <el-menu-item v-if="auth.isAdmin" index="/operation-logs">操作日志</el-menu-item>
        <el-menu-item index="/profile">个人中心</el-menu-item>
      </el-menu>
      <div class="user">
        <template v-if="auth.isLoggedIn">
          <span class="name">{{ auth.auth.user?.name || auth.auth.user?.phone }}</span>
          <el-button size="small" @click="logout">退出</el-button>
        </template>
        <template v-else>
          <el-button size="small" type="primary" @click="router.push('/login')">登录</el-button>
        </template>
      </div>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '@/hooks/useAuth'

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const active = computed(() => route.path)

function logout() {
  auth.auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout { min-height: 100vh; }
.header { display: flex; align-items: center; gap: 24px; border-bottom: 1px solid #e6e6e6; }
.brand { font-size: 20px; font-weight: 700; color: #409eff; cursor: pointer; }
.menu { flex: 1; border-bottom: none; }
.user { display: flex; align-items: center; gap: 8px; }
.name { color: #606266; }
</style>
