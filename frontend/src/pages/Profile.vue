<template>
  <el-card style="max-width: 480px">
    <div class="head">
      <el-avatar :size="72" :src="avatar">{{ (me?.name || '?')[0] }}</el-avatar>
    </div>
    <el-form :model="form" label-width="80px" class="mt">
      <el-form-item label="姓名"><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="手机号"><el-input :model-value="me?.phone" disabled /></el-form-item>
      <el-form-item label="角色">{{ me ? UserRoleText[me.role] : '' }}</el-form-item>
      <el-button type="primary" @click="save">保存</el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/userStore'
import { UserRoleText } from '@/constants/user'
import type { User } from '@/types'

const store = useUserStore()
const me = ref<User | null>(null)
const avatar = ref('')
const form = reactive({ name: '' })

onMounted(async () => {
  me.value = await store.fetchMe()
  form.name = me.value?.name || ''
  avatar.value = me.value?.avatar || ''
})

async function save() {
  me.value = await store.updateMe({ ...form, avatar: avatar.value })
  ElMessage.success('保存成功')
}
</script>

<style scoped>
.head { text-align: center; }
.mt { margin-top: 16px; }
</style>
