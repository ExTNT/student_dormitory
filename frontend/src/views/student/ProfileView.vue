<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { formatDateTime, roleLabels } from '@/utils/format';

const auth = useAuthStore();
</script>

<template>
  <section class="data-panel">
    <h2>个人信息</h2>
    <el-descriptions :column="1" border>
      <el-descriptions-item label="用户名">{{ auth.user?.username }}</el-descriptions-item>
      <el-descriptions-item label="姓名">{{ auth.user?.name }}</el-descriptions-item>
      <el-descriptions-item label="角色">{{ auth.user ? roleLabels[auth.user.role] : '' }}</el-descriptions-item>
      <el-descriptions-item label="电话">{{ auth.user?.phone || '-' }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ formatDateTime(auth.user?.created_at) }}</el-descriptions-item>
    </el-descriptions>
  </section>
</template>
