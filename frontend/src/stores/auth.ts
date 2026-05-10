import { defineStore } from 'pinia';
import { authApi } from '@/api/auth';
import { tokenStorage } from '@/api/http';
import type { LoginRequest, User } from '@/types';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: tokenStorage.accessToken,
    refreshToken: tokenStorage.refreshToken,
    user: null as User | null,
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.accessToken),
  },
  actions: {
    setTokens(accessToken: string, refreshToken: string) {
      this.accessToken = accessToken;
      this.refreshToken = refreshToken;
      tokenStorage.accessToken = accessToken;
      tokenStorage.refreshToken = refreshToken;
    },
    async login(data: LoginRequest) {
      const tokens = await authApi.login(data);
      this.setTokens(tokens.access_token, tokens.refresh_token);
      await this.fetchMe();
    },
    async fetchMe() {
      this.user = await authApi.me();
      return this.user;
    },
    async updateMe(data: { phone?: string }) {
      this.user = await authApi.updateMe(data);
      return this.user;
    },
    async refresh() {
      this.accessToken = tokenStorage.accessToken;
      this.refreshToken = tokenStorage.refreshToken;
    },
    async logout() {
      const refreshToken = tokenStorage.refreshToken;
      try {
        if (refreshToken) await authApi.logout(refreshToken);
      } finally {
        this.accessToken = '';
        this.refreshToken = '';
        this.user = null;
        tokenStorage.clear();
      }
    },
  },
});
