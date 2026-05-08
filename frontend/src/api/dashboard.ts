import { http } from './http';
import type { DormitorySummary, LowBalanceRoom } from '@/types';

export const dashboardApi = {
  summary() {
    return http.get<DormitorySummary[]>('/dashboard/summary').then((res) => res.data);
  },
  lowBalance() {
    return http.get<LowBalanceRoom[]>('/dashboard/low-balance').then((res) => res.data);
  },
};
