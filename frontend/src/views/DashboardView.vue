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
      return [
        ['/repair/orders', '接收维修工单'],
        ['/repair/my-orders', '我的维修工单'],
      ];
    case 'cleaning_staff':
      return [
        ['/cleaning/orders', '接收保洁工单'],
        ['/cleaning/my-orders', '我的保洁工单'],
      ];
    case 'dormitory_manager':
      return [
        ['/manager/leaves', '待审批离校'],
        ['/manager/room-changes', '换寝审批'],
        ['/manager/repairs', '待审核维修'],
        ['/manager/cleanings', '待审核保洁'],
        ['/manager/repairs/all', '全部维修工单'],
        ['/manager/cleanings/all', '全部保洁工单'],
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
    <div class="dashboard-hero">
      <div>
        <span class="eyebrow">Workspace</span>
        <h2>工作台</h2>
        <p>{{ auth.user?.name }} · {{ auth.user ? roleLabels[auth.user.role] : '' }}</p>
      </div>
      <div class="hero-stamp">
        <strong>{{ links.length }}</strong>
        <span>快捷入口</span>
      </div>
    </div>
    <el-row :gutter="18">
      <el-col v-for="[path, label] in links" :key="path" :xs="24" :sm="12" :md="8" :lg="6">
        <el-card shadow="never" class="quick-card" @click="$router.push(path)">
          <span class="quick-index">0{{ links.findIndex((item) => item[0] === path) + 1 }}</span>
          <strong>{{ label }}</strong>
          <small>进入处理</small>
        </el-card>
      </el-col>
    </el-row>
  </section>
</template>

<style scoped>
.dashboard-hero {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 18px;
  min-height: 180px;
  padding: 28px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(11, 93, 102, 0.92), rgba(7, 70, 77, 0.88)),
    linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.16));
  color: #f5fbfb;
  box-shadow: var(--shadow-md);
  overflow: hidden;
  position: relative;
}

.dashboard-hero::after {
  content: "";
  position: absolute;
  right: -80px;
  top: -90px;
  width: 260px;
  height: 260px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 50%;
}

.eyebrow {
  color: rgba(255, 255, 255, 0.66);
  font-size: 12px;
  font-weight: 850;
  text-transform: uppercase;
}

.dashboard-hero h2 {
  margin: 8px 0;
  color: #ffffff;
  font-size: 34px;
}

.dashboard-hero p {
  margin: 0;
  color: rgba(245, 251, 251, 0.78);
}

.hero-stamp {
  position: relative;
  z-index: 1;
  display: grid;
  place-items: center;
  width: 116px;
  height: 116px;
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.10);
  backdrop-filter: blur(10px);
}

.hero-stamp strong,
.hero-stamp span {
  display: block;
}

.hero-stamp strong {
  font-size: 42px;
  line-height: 1;
}

.hero-stamp span {
  color: rgba(255, 255, 255, 0.74);
  font-size: 12px;
}

.quick-card {
  position: relative;
  margin-bottom: 16px;
  min-height: 126px;
  cursor: pointer;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: var(--shadow-sm);
  transition: transform 180ms ease, box-shadow 180ms ease, border-color 180ms ease;
}

.quick-card:hover {
  transform: translateY(-3px);
  border-color: rgba(11, 93, 102, 0.26);
  box-shadow: var(--shadow-md);
}

.quick-card strong,
.quick-card small,
.quick-index {
  display: block;
}

.quick-index {
  color: var(--accent);
  font-size: 12px;
  font-weight: 900;
}

.quick-card strong {
  margin-top: 18px;
  font-size: 18px;
}

.quick-card small {
  margin-top: 8px;
  color: var(--text-muted);
}
</style>
