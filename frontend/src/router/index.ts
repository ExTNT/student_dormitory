import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import type { Role } from '@/types';

const roleHome: Record<Role, string> = {
  student: '/student/dashboard',
  repair_staff: '/repair/dashboard',
  cleaning_staff: '/cleaning/dashboard',
  dormitory_manager: '/manager/dashboard',
  system_admin: '/admin/dashboard',
};

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: () => import('@/views/LoginView.vue'), meta: { public: true } },
  { path: '/403', component: () => import('@/views/ForbiddenView.vue'), meta: { public: true } },
  {
    path: '/',
    component: () => import('@/views/MainLayout.vue'),
    children: [
      { path: 'student/dashboard', component: () => import('@/views/DashboardView.vue'), meta: { roles: ['student'] } },
      { path: 'student/profile', component: () => import('@/views/student/ProfileView.vue'), meta: { roles: ['student'] } },
      { path: 'student/survey', component: () => import('@/views/student/SurveyView.vue'), meta: { roles: ['student'] } },
      { path: 'student/allocation', component: () => import('@/views/student/AllocationView.vue'), meta: { roles: ['student'] } },
      { path: 'student/requests', component: () => import('@/views/student/MyRequestsView.vue'), meta: { roles: ['student'] } },
      { path: 'student/roommates', component: () => import('@/views/student/RoommatesView.vue'), meta: { roles: ['student'] } },
      { path: 'student/leave', component: () => import('@/views/student/LeaveView.vue'), meta: { roles: ['student'] } },
      { path: 'student/late-return', component: () => import('@/views/student/LateReturnView.vue'), meta: { roles: ['student'] } },
      { path: 'student/room-change', component: () => import('@/views/student/RoomChangeView.vue'), meta: { roles: ['student'] } },
      { path: 'student/off-campus', component: () => import('@/views/student/OffCampusView.vue'), meta: { roles: ['student'] } },
      { path: 'student/repair', component: () => import('@/views/student/RepairCreateView.vue'), meta: { roles: ['student'] } },
      { path: 'student/cleaning', component: () => import('@/views/student/CleaningCreateView.vue'), meta: { roles: ['student'] } },
      { path: 'student/payment', component: () => import('@/views/student/PaymentView.vue'), meta: { roles: ['student'] } },
      { path: 'student/notifications', component: () => import('@/views/NotificationsView.vue'), meta: { roles: ['student'] } },
      { path: 'repair/dashboard', component: () => import('@/views/DashboardView.vue'), meta: { roles: ['repair_staff'] } },
      { path: 'repair/orders', component: () => import('@/views/repair/RepairOrdersView.vue'), meta: { roles: ['repair_staff'] } },
      { path: 'cleaning/dashboard', component: () => import('@/views/DashboardView.vue'), meta: { roles: ['cleaning_staff'] } },
      { path: 'cleaning/orders', component: () => import('@/views/cleaning/CleaningOrdersView.vue'), meta: { roles: ['cleaning_staff'] } },
      { path: 'manager/dashboard', component: () => import('@/views/DashboardView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/leaves', component: () => import('@/views/manager/LeavesReviewView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/late-returns', component: () => import('@/views/manager/LateReturnsReviewView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/room-changes', component: () => import('@/views/manager/RoomChangesReviewView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/off-campus', component: () => import('@/views/manager/OffCampusReviewView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/repairs', component: () => import('@/views/manager/RepairReviewView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/cleanings', component: () => import('@/views/manager/CleaningReviewView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/summary', component: () => import('@/views/dashboard/SummaryView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'manager/low-balance', component: () => import('@/views/dashboard/LowBalanceView.vue'), meta: { roles: ['dormitory_manager'] } },
      { path: 'admin/dashboard', component: () => import('@/views/DashboardView.vue'), meta: { roles: ['system_admin'] } },
      { path: 'admin/users', component: () => import('@/views/admin/UserCreateView.vue'), meta: { roles: ['system_admin'] } },
      { path: 'admin/allocations', component: () => import('@/views/admin/AllocationsReviewView.vue'), meta: { roles: ['system_admin'] } },
      { path: 'admin/summary', component: () => import('@/views/dashboard/SummaryView.vue'), meta: { roles: ['system_admin'] } },
      { path: 'admin/low-balance', component: () => import('@/views/dashboard/LowBalanceView.vue'), meta: { roles: ['system_admin'] } },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to) => {
  const auth = useAuthStore();
  if (to.meta.public) return true;
  if (!auth.accessToken) return { path: '/login', query: { redirect: to.fullPath } };
  if (!auth.user) {
    try {
      await auth.fetchMe();
    } catch {
      return { path: '/login', query: { redirect: to.fullPath } };
    }
  }
  const roles = to.meta.roles as Role[] | undefined;
  if (roles?.length && auth.user && !roles.includes(auth.user.role)) return '/403';
  return true;
});

export { roleHome };
export default router;
