<script setup lang="ts">
import { computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { roleLabels } from '@/utils/format';

const auth = useAuthStore();

const links = computed(() => {
  switch (auth.user?.role) {
    case 'student':
      return [
        ['/student/requests', '我的申请'],
        ['/student/notifications', '通知'],
        ['/student/payment', '缴费'],
        ['/student/repair', '维修申请'],
        ['/student/cleaning', '保洁申请'],
      ];
    case 'repair_staff':
      return [['/repair/orders', '待处理维修工单']];
    case 'cleaning_staff':
      return [['/cleaning/orders', '待处理保洁工单']];
    case 'dormitory_manager':
      return [
        ['/manager/leaves', '待审批离校'],
        ['/manager/room-changes', '换寝审批'],
        ['/manager/summary', '统计'],
      ];
    case 'system_admin':
      return [
        ['/admin/allocations', '分配审批'],
        ['/admin/users', '用户创建'],
        ['/admin/summary', '统计'],
      ];
    default:
      return [];
  }
});
</script>

<template>
  <section class="page">
    <div class="data-panel">
      <h2>工作台</h2>
      <p>{{ auth.user?.name }} · {{ auth.user ? roleLabels[auth.user.role] : '' }}</p>
    </div>
    <el-row :gutter="16">
      <el-col v-for="[path, label] in links" :key="path" :xs="24" :sm="12" :md="8" :lg="6">
        <el-card shadow="never" class="quick-card" @click="$router.push(path)">
          <strong>{{ label }}</strong>
        </el-card>
      </el-col>
    </el-row>
  </section>
</template>

<style scoped>
h2 {
  margin: 0 0 8px;
}

.quick-card {
  margin-bottom: 16px;
  cursor: pointer;
  border-radius: 8px;
}
</style>
