import { http } from './http';
import type { LoginRequest, LoginResponse, User } from '@/types';

export const authApi = {
  login(data: LoginRequest) {
    return http.post<LoginResponse>('/auth/login', data, { suppressErrorMessage: true } as any).then((res) => res.data);
  },
  logout(refreshToken: string) {
    return http.post('/auth/logout', { refresh_token: refreshToken });
  },
  me() {
    return http.get<User>('/students/me').then((res) => res.data);
  },
  updateMe(data: { phone?: string }) {
    return http.put<User>('/students/me', data).then((res) => res.data);
  },
};
