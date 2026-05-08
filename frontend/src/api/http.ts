import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { ElMessage } from 'element-plus';
import router from '@/router';
import type { LoginResponse } from '@/types';

export const API_BASE_URL = 'http://localhost:8080/api';

export const tokenStorage = {
  get accessToken() {
    return localStorage.getItem('access_token') || '';
  },
  set accessToken(value: string) {
    localStorage.setItem('access_token', value);
  },
  get refreshToken() {
    return localStorage.getItem('refresh_token') || '';
  },
  set refreshToken(value: string) {
    localStorage.setItem('refresh_token', value);
  },
  clear() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  },
};

interface RetryConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

let refreshPromise: Promise<LoginResponse> | null = null;

export const http = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
});

const refreshClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
});

http.interceptors.request.use((config) => {
  const token = tokenStorage.accessToken;
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

async function refreshTokens() {
  if (!refreshPromise) {
    refreshPromise = refreshClient
      .post<LoginResponse>('/auth/refresh', { refresh_token: tokenStorage.refreshToken })
      .then((res) => {
        tokenStorage.accessToken = res.data.access_token;
        tokenStorage.refreshToken = res.data.refresh_token;
        return res.data;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

function getBackendMessage(error: AxiosError<{ message?: string; error?: string }>) {
  return error.response?.data?.message || error.response?.data?.error || error.message || '请求失败';
}

http.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<{ message?: string; error?: string }>) => {
    const config = error.config as RetryConfig | undefined;
    if (error.response?.status === 401 && config && !config._retry && tokenStorage.refreshToken) {
      config._retry = true;
      try {
        const tokens = await refreshTokens();
        config.headers.Authorization = `Bearer ${tokens.access_token}`;
        return http(config);
      } catch (refreshError) {
        tokenStorage.clear();
        ElMessage.error('登录已过期，请重新登录');
        router.replace({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } });
        return Promise.reject(refreshError);
      }
    }
    if (error.response?.status !== 401) {
      ElMessage.error(getBackendMessage(error));
    }
    return Promise.reject(error);
  },
);
