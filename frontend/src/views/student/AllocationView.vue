<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { studentApi } from '@/api/student';

const createdId = ref<number>();
const loading = ref(false);

async function submit() {
  await ElMessageBox.confirm('确认提交新生分配申请？提交前请确保已完成生活习惯调查。', '确认申请', { type: 'warning' });
  loading.value = true;
  try {
    const res = await studentApi.createAllocation();
    createdId.value = res.id;
    ElMessage.success('申请已提交');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="form-card">
    <h2>新生分配申请</h2>
    <el-alert title="请先提交生活习惯调查，后端会根据生活习惯自动推荐床位。" type="info" show-icon :closable="false" />
    <div style="margin-top: 16px">
      <el-button type="primary" :loading="loading" @click="submit">提交申请</el-button>
    </div>
    <el-result v-if="createdId" icon="success" title="申请已创建" :sub-title="`申请 ID：${createdId}`" />
  </section>
</template>
