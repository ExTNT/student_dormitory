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

const visibleMenus = computed(() => {
  if (!auth.user) return [];
  const items = menus[auth.user.role];
  if (auth.user.role !== 'student') return items;
  return items.filter((item) => {
    if (item.path === '/student/survey') return !auth.user?.has_survey;
    if (item.path === '/student/allocation') return !auth.user?.has_bed;
    return true;
  });
});
const activeMenu = computed(() => visibleMenus.value.find((item) => item.path === route.path)?.label || '工作台');

async function logout() {
  await auth.logout();
  router.replace('/login');
}
</script>

<template>
  <el-container class="layout">
    <el-aside width="232px">
      <div class="brand">
        <span class="brand-mark">宿</span>
        <span>
          <strong>宿舍管理系统</strong>
          <small>Dormitory Ops</small>
        </span>
      </div>
      <el-menu router :default-active="route.path">
        <el-menu-item v-for="item in visibleMenus" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-title">
          <span>{{ activeMenu }}</span>
          <small>校园住宿服务与审批工作台</small>
        </div>
        <div class="userbar">
          <div class="user-chip">
            <span class="avatar">{{ auth.user?.name?.slice(0, 1) }}</span>
            <span class="user-meta">
              <strong>{{ auth.user?.name }}</strong>
              <small>{{ auth.user ? roleLabels[auth.user.role] : '' }}</small>
            </span>
          </div>
          <el-button :icon="SwitchButton" text class="logout-button" @click="logout">退出</el-button>
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
  background: transparent;
}

.brand {
  height: 70px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 18px 0 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.10);
  color: #eef8f7;
}

.brand-mark {
  display: grid;
  place-items: center;
  width: 38px;
  height: 38px;
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 8px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.22), rgba(255, 255, 255, 0.06));
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.16);
  color: #ffffff;
  font-weight: 850;
}

.brand strong,
.brand small {
  display: block;
}

.brand strong {
  font-size: 16px;
  letter-spacing: 0;
}

.brand small {
  margin-top: 2px;
  color: rgba(238, 248, 247, 0.62);
  font-size: 11px;
  text-transform: uppercase;
}

.el-aside {
  background:
    radial-gradient(circle at 24px 18px, rgba(255, 255, 255, 0.16), transparent 26px),
    linear-gradient(180deg, #0d343a 0%, #08272d 100%);
  border-right: 0;
  box-shadow: 14px 0 38px rgba(8, 39, 45, 0.16);
}

.el-menu {
  border-right: 0;
  padding: 12px;
  background: transparent;
}

.el-menu :deep(.el-menu-item) {
  height: 42px;
  margin: 4px 0;
  border-radius: 8px;
  color: rgba(238, 248, 247, 0.70);
}

.el-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
}

.el-menu :deep(.el-menu-item.is-active) {
  background: #f4fbfa;
  color: #0a515a;
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.18);
}

.el-menu :deep(.el-icon) {
  font-size: 17px;
}

.el-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 70px;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.78);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(18px);
}

.header-title span,
.header-title small {
  display: block;
}

.header-title span {
  font-weight: 800;
  font-size: 18px;
}

.header-title small {
  margin-top: 3px;
  color: var(--text-muted);
  font-size: 12px;
}

.userbar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 168px;
  padding: 7px 10px 7px 8px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.78);
}

.avatar {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  border-radius: 8px;
  background: var(--brand-soft);
  color: var(--brand-strong);
  font-weight: 850;
}

.user-meta strong,
.user-meta small {
  display: block;
  line-height: 1.2;
}

.user-meta strong {
  font-size: 14px;
}

.user-meta small {
  margin-top: 3px;
  color: var(--text-muted);
  font-size: 12px;
}

.logout-button {
  color: var(--text-muted);
}

.el-main {
  padding: 24px;
  overflow-x: hidden;
}

@media (max-width: 840px) {
  .el-aside {
    width: 78px !important;
  }

  .brand span:not(.brand-mark),
  .el-menu :deep(.el-menu-item span) {
    display: none;
  }

  .brand {
    justify-content: center;
    padding: 0;
  }

  .header-title small,
  .user-meta {
    display: none;
  }

  .user-chip {
    min-width: auto;
    padding: 6px;
  }
}
</style>
