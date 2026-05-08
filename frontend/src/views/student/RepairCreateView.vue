<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';

const loading = ref(false);
const createdId = ref<number>();
const form = reactive({ room_id: undefined as number | undefined, description: '' });
async function submit() {
  if (!form.room_id) return ElMessage.warning('请输入房间 ID');
  loading.value = true;
  try {
    const res = await studentApi.createRepair({ room_id: form.room_id, description: form.description });
    createdId.value = res.id;
    ElMessage.success('维修申请已提交');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="form-card">
    <h2>维修申请</h2>
    <el-alert title="后端会校验当前学生是否入住该房间。" type="info" show-icon :closable="false" />
    <el-form :model="form" label-width="100px" style="margin-top: 16px">
      <el-form-item label="房间 ID" required><el-input-number v-model="form.room_id" :min="1" /></el-form-item>
      <el-form-item label="描述" required><el-input v-model="form.description" type="textarea" /></el-form-item>
      <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交</el-button></el-form-item>
    </el-form>
    <el-alert v-if="createdId" :title="`维修工单 ID：${createdId}`" type="success" show-icon />
  </section>
</template>
