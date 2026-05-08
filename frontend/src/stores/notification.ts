import { defineStore } from 'pinia';
import { http } from '@/api/http';
import type { Notification } from '@/types';

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [] as Notification[],
  }),
  actions: {
    async fetchNotifications() {
      this.notifications = await http.get<Notification[]>('/notifications').then((res) => res.data);
    },
    async markRead(id: number) {
      await http.put(`/notifications/${id}/read`);
      await this.fetchNotifications();
    },
  },
});
