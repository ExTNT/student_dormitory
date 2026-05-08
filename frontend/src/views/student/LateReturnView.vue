<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';

const loading = ref(false);
const form = reactive({ return_date: '', reason: '' });
async function submit() {
  loading.value = true;
  try {
    await studentApi.createLateReturn(form);
    ElMessage.success('记录已提交');
    Object.assign(form, { return_date: '', reason: '' });
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="form-card">
    <h2>晚归记录</h2>
    <el-form :model="form" label-width="100px">
      <el-form-item label="晚归日期" required><el-date-picker v-model="form.return_date" value-format="YYYY-MM-DD" /></el-form-item>
      <el-form-item label="原因" required><el-input v-model="form.reason" type="textarea" /></el-form-item>
      <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交</el-button></el-form-item>
    </el-form>
  </section>
</template>
