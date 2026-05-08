<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { adminApi } from '@/api/admin';
import type { Role } from '@/types';

const loading = ref(false);
const form = reactive({ username: '', password: '123456', role: 'student' as Role, name: '', phone: '' });
async function submit() {
  loading.value = true;
  try {
    await adminApi.createUser(form);
    ElMessage.success('用户已创建');
    Object.assign(form, { username: '', password: '123456', role: 'student', name: '', phone: '' });
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="form-card">
    <h2>创建用户</h2>
    <el-form :model="form" label-width="100px">
      <el-form-item label="用户名" required><el-input v-model="form.username" /></el-form-item>
      <el-form-item label="密码" required><el-input v-model="form.password" type="password" show-password /></el-form-item>
      <el-form-item label="角色" required><el-select v-model="form.role"><el-option label="学生" value="student" /><el-option label="维修人员" value="repair_staff" /><el-option label="保洁人员" value="cleaning_staff" /><el-option label="宿舍管理员" value="dormitory_manager" /><el-option label="系统管理员" value="system_admin" /></el-select></el-form-item>
      <el-form-item label="姓名" required><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="电话"><el-input v-model="form.phone" /></el-form-item>
      <el-form-item><el-button type="primary" :loading="loading" @click="submit">创建</el-button></el-form-item>
    </el-form>
  </section>
</template>
