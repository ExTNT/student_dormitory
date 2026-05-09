import { http } from './http';
import type { PendingCleaning } from '@/types';

export const cleaningApi = {
  list() {
    return http.get<PendingCleaning[]>('/cleanings').then((res) => res.data);
  },
  pending() {
    return http.get<PendingCleaning[]>('/cleanings/pending').then((res) => res.data);
  },
  accept(id: number) {
    return http.put(`/cleanings/${id}/accept`);
  },
  complete(id: number) {
    return http.put(`/cleanings/${id}/clean`);
  },
  review(id: number, data: { status: 'completed' | 'rejected'; comment?: string }) {
    return http.put(`/cleanings/${id}/review`, data);
  },
};
