<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useAuthStore } from '@/stores/auth';
import { roleHome } from '@/router';

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();
const loading = ref(false);
const form = reactive({ username: '', password: '' });

async function submit() {
  loading.value = true;
  try {
    await auth.login(form);
    ElMessage.success('登录成功');
    router.replace((route.query.redirect as string) || roleHome[auth.user!.role]);
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <main class="login-page">
    <section class="login-panel">
      <h1>宿舍管理系统</h1>
      <el-form :model="form" label-position="top" @keyup.enter="submit">
        <el-form-item label="用户名" required>
          <el-input v-model="form.username" autocomplete="username" />
        </el-form-item>
        <el-form-item label="密码" required>
          <el-input v-model="form.password" type="password" autocomplete="current-password" show-password />
        </el-form-item>
        <el-button type="primary" :loading="loading" class="login-button" @click="submit">登录</el-button>
      </el-form>
      <p class="muted">测试账号密码均为 123456</p>
    </section>
  </main>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, #eef5ff 0%, #f7f9fc 45%, #eef8f3 100%);
}

.login-panel {
  width: min(420px, calc(100vw - 32px));
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 20px 50px rgb(31 41 55 / 10%);
}

h1 {
  margin: 0 0 24px;
  font-size: 26px;
}

.login-button {
  width: 100%;
}
</style>
