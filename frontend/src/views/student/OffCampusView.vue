<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';

const loading = ref(false);
const form = reactive({ retain_bed: 0, reason: '', destination: '' });
async function submit() {
  loading.value = true;
  try {
    await studentApi.createOffCampus(form);
    ElMessage.success('申请已提交');
    Object.assign(form, { retain_bed: 0, reason: '', destination: '' });
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="form-card">
    <h2>校外居住申请</h2>
    <el-form :model="form" label-width="110px">
      <el-form-item label="保留床位"><el-switch v-model="form.retain_bed" :active-value="1" :inactive-value="0" /></el-form-item>
      <el-form-item label="原因" required><el-input v-model="form.reason" type="textarea" /></el-form-item>
      <el-form-item label="居住地址"><el-input v-model="form.destination" /></el-form-item>
      <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交</el-button></el-form-item>
    </el-form>
  </section>
</template>
