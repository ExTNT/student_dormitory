import { http } from './http';
import type { AllocationRequest, Role, User } from '@/types';

export const adminApi = {
  createUser(data: { username: string; password: string; role: Role; name: string; phone?: string }) {
    return http.post<User>('/users', data).then((res) => res.data);
  },
  allocations() {
    return http.get<AllocationRequest[]>('/allocations/pending').then((res) => res.data);
  },
  reviewAllocation(id: number, status: 'approved' | 'rejected') {
    return http.put(`/allocations/${id}/review`, { status });
  },
};
