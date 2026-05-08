<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';
import { toIso } from '@/utils/format';

const loading = ref(false);
const form = reactive({ type: 'normal', destination: '', emergency_contact: '', return_time: '', reason: '' });

async function submit() {
  loading.value = true;
  try {
    await studentApi.createLeave({ ...form, return_time: toIso(form.return_time) });
    ElMessage.success('申请已提交');
    Object.assign(form, { type: 'normal', destination: '', emergency_contact: '', return_time: '', reason: '' });
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="form-card">
    <h2>离校/节假日离校申请</h2>
    <el-form :model="form" label-width="120px">
      <el-form-item label="类型"><el-select v-model="form.type"><el-option label="普通离校" value="normal" /><el-option label="节假日离校" value="holiday" /></el-select></el-form-item>
      <el-form-item label="目的地" required><el-input v-model="form.destination" /></el-form-item>
      <el-form-item label="紧急联系人" required><el-input v-model="form.emergency_contact" /></el-form-item>
      <el-form-item label="返校时间" required><el-date-picker v-model="form.return_time" type="datetime" /></el-form-item>
      <el-form-item label="原因" required><el-input v-model="form.reason" type="textarea" /></el-form-item>
      <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交</el-button></el-form-item>
    </el-form>
  </section>
</template>
