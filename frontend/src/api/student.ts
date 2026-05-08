import { http } from './http';
import type { AvailableBed, Building, LifestyleSurvey, MyRequest, Room, Roommate } from '@/types';

export const studentApi = {
  me() {
    return http.get('/students/me').then((res) => res.data);
  },
  getSurvey() {
    return http.get<LifestyleSurvey>('/students/me/survey').then((res) => res.data);
  },
  saveSurvey(data: LifestyleSurvey) {
    return http.post<LifestyleSurvey>('/students/me/survey', data).then((res) => res.data);
  },
  createAllocation() {
    return http.post<{ id: number }>('/allocations').then((res) => res.data);
  },
  requests() {
    return http.get<MyRequest[]>('/students/me/requests').then((res) => res.data);
  },
  roommates() {
    return http.get<Roommate[]>('/students/me/roommates').then((res) => res.data);
  },
  createLeave(data: Record<string, unknown>) {
    return http.post<{ id: number }>('/leaves', data).then((res) => res.data);
  },
  createLateReturn(data: Record<string, unknown>) {
    return http.post<{ id: number }>('/late-returns', data).then((res) => res.data);
  },
  availableBeds(params?: { building_id?: number; floor?: number }) {
    return http.get<AvailableBed[]>('/beds/available', { params }).then((res) => res.data);
  },
  createRoomChange(data: Record<string, unknown>) {
    return http.post<{ id: number }>('/room-changes', data).then((res) => res.data);
  },
  createOffCampus(data: Record<string, unknown>) {
    return http.post<{ id: number }>('/off-campus', data).then((res) => res.data);
  },
  createRepair(data: { room_id: number; description: string }) {
    return http.post<{ id: number }>('/repairs', data).then((res) => res.data);
  },
  createCleaning(data: { building_id: number; location_desc: string }) {
    return http.post<{ id: number }>('/cleanings', data).then((res) => res.data);
  },
  buildings() {
    return http.get<Building[]>('/buildings').then((res) => res.data);
  },
  roomBalance(id: number) {
    return http.get<Room>(`/rooms/${id}/balance`).then((res) => res.data);
  },
  createPayment(data: { room_id: number; amount: number; payment_type: string }) {
    return http.post<{ id: number }>('/payments', data).then((res) => res.data);
  },
};
