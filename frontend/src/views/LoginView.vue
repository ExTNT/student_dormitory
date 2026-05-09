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
  } catch {
    ElMessage.error('账号和密码错误');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <main class="login-page">
    <section class="login-copy">
      <span class="eyebrow">Dormitory Operations</span>
      <h1>宿舍管理系统</h1>
      <p>面向学生、宿管、维修、保洁与管理员的一站式住宿服务工作台。</p>
      <div class="metric-strip">
        <span>申请</span>
        <span>审批</span>
        <span>工单</span>
        <span>缴费</span>
      </div>
    </section>
    <section class="login-panel">
      <div class="panel-head">
        <h2>登录</h2>
        <span>REST API</span>
      </div>
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
  grid-template-columns: minmax(280px, 520px) minmax(360px, 430px);
  align-items: center;
  justify-content: center;
  gap: 58px;
  padding: 40px;
  background:
    radial-gradient(circle at 18% 22%, rgba(11, 93, 102, 0.22), transparent 26%),
    radial-gradient(circle at 82% 74%, rgba(201, 130, 43, 0.18), transparent 24%),
    linear-gradient(135deg, #edf4f6 0%, #f8faf8 48%, #eef5f0 100%);
}

.login-copy {
  position: relative;
  padding: 38px 0;
}

.login-copy::before {
  content: "";
  position: absolute;
  inset: 0 auto 0 -28px;
  width: 4px;
  border-radius: 8px;
  background: linear-gradient(180deg, var(--brand), var(--accent));
}

.eyebrow {
  display: inline-flex;
  margin-bottom: 20px;
  padding: 7px 10px;
  border: 1px solid rgba(11, 93, 102, 0.18);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.52);
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 800;
  text-transform: uppercase;
}

.login-copy h1 {
  margin: 0;
  font-size: clamp(42px, 7vw, 76px);
  line-height: 1.02;
  font-weight: 900;
}

.login-copy p {
  max-width: 430px;
  margin: 20px 0 0;
  color: #526371;
  font-size: 17px;
  line-height: 1.75;
}

.metric-strip {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 30px;
}

.metric-strip span {
  padding: 9px 13px;
  border: 1px solid rgba(11, 93, 102, 0.14);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.58);
  color: #26414b;
  font-weight: 750;
}

.login-panel {
  width: min(420px, calc(100vw - 32px));
  background: rgba(255, 255, 255, 0.84);
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 8px;
  padding: 28px;
  box-shadow: var(--shadow-lg);
  backdrop-filter: blur(20px);
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 22px;
}

.panel-head h2 {
  margin: 0;
  font-size: 26px;
}

.panel-head span {
  padding: 5px 9px;
  border-radius: 8px;
  background: var(--brand-soft);
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 800;
}

.login-panel :deep(.el-form-item__label) {
  font-weight: 750;
  color: #334451;
}

.login-button {
  width: 100%;
  height: 42px;
  margin-top: 4px;
}

@media (max-width: 860px) {
  .login-page {
    grid-template-columns: 1fr;
    gap: 22px;
    padding: 28px 18px;
  }

  .login-copy {
    padding: 10px 0 0 18px;
  }

  .login-panel {
    width: 100%;
  }
}
</style>
