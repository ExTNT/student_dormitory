<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { House, User, Bell, Tickets, Tools, Brush, Setting, DataLine, Money, SwitchButton } from '@element-plus/icons-vue';
import { useAuthStore } from '@/stores/auth';
import { roleLabels } from '@/utils/format';
import type { Role } from '@/types';

const auth = useAuthStore();
const route = useRoute();
const router = useRouter();

const menus: Record<Role, { path: string; label: string; icon: unknown }[]> = {
  student: [
    { path: '/student/dashboard', label: '工作台', icon: House },
    { path: '/student/profile', label: '个人信息', icon: User },
    { path: '/student/survey', label: '生活习惯', icon: Tickets },
    { path: '/student/allocation', label: '分配申请', icon: Tickets },
    { path: '/student/requests', label: '我的申请', icon: Tickets },
    { path: '/student/roommates', label: '舍友信息', icon: User },
    { path: '/student/leave', label: '离校申请', icon: Tickets },
    { path: '/student/late-return', label: '晚归记录', icon: Tickets },
    { path: '/student/room-change', label: '换寝申请', icon: Tickets },
    { path: '/student/off-campus', label: '校外居住', icon: Tickets },
    { path: '/student/repair', label: '维修申请', icon: Tools },
    { path: '/student/cleaning', label: '保洁申请', icon: Brush },
    { path: '/student/payment', label: '水电缴费', icon: Money },
    { path: '/student/notifications', label: '通知', icon: Bell },
  ],
  repair_staff: [
    { path: '/repair/dashboard', label: '工作台', icon: House },
    { path: '/repair/orders', label: '维修工单', icon: Tools },
  ],
  cleaning_staff: [
    { path: '/cleaning/dashboard', label: '工作台', icon: House },
    { path: '/cleaning/orders', label: '保洁工单', icon: Brush },
  ],
  dormitory_manager: [
    { path: '/manager/dashboard', label: '工作台', icon: House },
    { path: '/manager/leaves', label: '离校审批', icon: Tickets },
    { path: '/manager/late-returns', label: '晚归审批', icon: Tickets },
    { path: '/manager/room-changes', label: '换寝审批', icon: Tickets },
    { path: '/manager/off-campus', label: '校外居住审批', icon: Tickets },
    { path: '/manager/repairs', label: '维修审核', icon: Tools },
    { path: '/manager/cleanings', label: '保洁审核', icon: Brush },
    { path: '/manager/summary', label: '楼栋统计', icon: DataLine },
    { path: '/manager/low-balance', label: '低余额房间', icon: Money },
  ],
  system_admin: [
    { path: '/admin/dashboard', label: '工作台', icon: House },
    { path: '/admin/users', label: '创建用户', icon: Setting },
    { path: '/admin/allocations', label: '分配审批', icon: Tickets },
    { path: '/admin/summary', label: '楼栋统计', icon: DataLine },
    { path: '/admin/low-balance', label: '低余额房间', icon: Money },
  ],
};

const visibleMenus = computed(() => (auth.user ? menus[auth.user.role] : []));

async function logout() {
  await auth.logout();
  router.replace('/login');
}
</script>

<template>
  <el-container class="layout">
    <el-aside width="232px">
      <div class="brand">宿舍管理系统</div>
      <el-menu router :default-active="route.path">
        <el-menu-item v-for="item in visibleMenus" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div></div>
        <div class="userbar">
          <span>{{ auth.user?.name }}</span>
          <el-tag>{{ auth.user ? roleLabels[auth.user.role] : '' }}</el-tag>
          <el-button :icon="SwitchButton" text @click="logout">退出</el-button>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout {
  min-height: 100vh;
}

.brand {
  height: 58px;
  display: flex;
  align-items: center;
  padding: 0 18px;
  font-weight: 700;
  border-bottom: 1px solid #e5e7eb;
}

.el-aside {
  background: #fff;
  border-right: 1px solid #e5e7eb;
}

.el-menu {
  border-right: 0;
}

.el-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
}

.userbar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.el-main {
  padding: 20px;
}
</style>
