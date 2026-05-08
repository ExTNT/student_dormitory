import { http } from './http';
import type { PendingLateReturn, PendingLeave, PendingOffCampus, PendingRoomChange } from '@/types';

export const managerApi = {
  leaves() {
    return http.get<PendingLeave[]>('/leaves/pending').then((res) => res.data);
  },
  reviewLeave(id: number, status: 'approved' | 'rejected') {
    return http.put(`/leaves/${id}/review`, { status });
  },
  lateReturns() {
    return http.get<PendingLateReturn[]>('/late-returns/pending').then((res) => res.data);
  },
  reviewLateReturn(id: number, status: 'approved' | 'rejected') {
    return http.put(`/late-returns/${id}/review`, { status });
  },
  roomChanges() {
    return http.get<PendingRoomChange[]>('/room-changes/pending').then((res) => res.data);
  },
  reviewRoomChange(id: number, status: 'approved' | 'rejected') {
    return http.put(`/room-changes/${id}/review`, { status });
  },
  offCampus() {
    return http.get<PendingOffCampus[]>('/off-campus/pending').then((res) => res.data);
  },
  reviewOffCampus(id: number, status: 'approved' | 'rejected') {
    return http.put(`/off-campus/${id}/review`, { status });
  },
};
