<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const loading = ref(false);
const createdId = ref<number>();
const form = reactive({ description: '' });
const dormitoryText = computed(() => {
  const user = auth.user;
  if (!user?.has_bed) return '未分配宿舍';
  return [user.building_name, user.room_number, user.bed_label].filter(Boolean).join(' / ');
});

async function submit() {
  if (!auth.user?.room_id) return ElMessage.warning('当前账号未分配宿舍，无法提交维修申请');
  loading.value = true;
  try {
    const res = await studentApi.createRepair({ room_id: auth.user.room_id, description: form.description });
    createdId.value = res.id;
    ElMessage.success('维修申请已提交');
    form.description = '';
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="form-card">
    <h2>维修申请</h2>
    <el-alert title="房间信息由当前学生宿舍自动补齐，后端会校验当前学生是否入住该房间。" type="info" show-icon :closable="false" />
    <el-form :model="form" label-width="100px" style="margin-top: 16px">
      <el-form-item label="当前宿舍">
        <el-input :model-value="dormitoryText" disabled />
      </el-form-item>
      <el-form-item label="描述" required><el-input v-model="form.description" type="textarea" /></el-form-item>
      <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交</el-button></el-form-item>
    </el-form>
    <el-alert v-if="createdId" :title="`维修工单 ID：${createdId}`" type="success" show-icon />
  </section>
</template>
