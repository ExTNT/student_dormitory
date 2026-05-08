import type { Role } from '@/types';

export const roleLabels: Record<Role, string> = {
  student: '学生',
  repair_staff: '维修人员',
  cleaning_staff: '保洁人员',
  dormitory_manager: '宿舍管理员',
  system_admin: '系统管理员',
};

export const statusTagType: Record<string, 'success' | 'warning' | 'danger' | 'primary' | 'info'> = {
  pending: 'warning',
  approved: 'success',
  rejected: 'danger',
  accepted: 'primary',
  repaired: 'primary',
  cleaned: 'primary',
  completed: 'success',
  paid: 'success',
};

export function formatDateTime(value?: string | null) {
  if (!value) return '-';
  return new Date(value).toLocaleString();
}

export function formatDate(value?: string | null) {
  if (!value) return '-';
  return new Date(value).toLocaleDateString();
}

export function toIso(value?: string | Date | null) {
  if (!value) return '';
  return new Date(value).toISOString();
}

export function displayBool(value: number | boolean | undefined) {
  return value ? '是' : '否';
}
