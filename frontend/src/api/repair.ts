import { http } from './http';
import type { PendingRepair } from '@/types';

export const repairApi = {
  pending() {
    return http.get<PendingRepair[]>('/repairs/pending').then((res) => res.data);
  },
  accept(id: number) {
    return http.put(`/repairs/${id}/accept`);
  },
  complete(id: number, repair_description: string) {
    return http.put(`/repairs/${id}/repair`, { repair_description });
  },
  review(id: number, data: { status: 'completed' | 'rejected'; comment?: string }) {
    return http.put(`/repairs/${id}/review`, data);
  },
};
