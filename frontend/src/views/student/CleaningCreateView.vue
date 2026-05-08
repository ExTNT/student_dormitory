<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';
import type { Building } from '@/types';
import ImageUploader from '@/components/ImageUploader.vue';
import AttachmentList from '@/components/AttachmentList.vue';

const buildings = ref<Building[]>([]);
const createdId = ref<number>();
const listRef = ref<InstanceType<typeof AttachmentList>>();
const loading = ref(false);
const form = reactive({ building_id: undefined as number | undefined, location_desc: '' });

async function submit() {
  if (!form.building_id) return ElMessage.warning('请选择楼栋');
  loading.value = true;
  try {
    const res = await studentApi.createCleaning({ building_id: form.building_id, location_desc: form.location_desc });
    createdId.value = res.id;
    ElMessage.success('保洁工单已创建，可上传保洁前照片');
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  buildings.value = await studentApi.buildings();
});
</script>

<template>
  <section class="page">
    <div class="form-card">
      <h2>保洁申请</h2>
      <el-form :model="form" label-width="110px">
        <el-form-item label="楼栋" required>
          <el-select v-model="form.building_id" filterable>
            <el-option v-for="b in buildings" :key="b.id" :label="b.name" :value="b.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="位置描述" required><el-input v-model="form.location_desc" /></el-form-item>
        <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交</el-button></el-form-item>
      </el-form>
    </div>
    <div v-if="createdId" class="data-panel">
      <div class="toolbar">
        <strong>保洁前照片</strong>
        <ImageUploader owner-type="cleaning" :owner-id="createdId" category="before" @success="listRef?.fetchList()" />
      </div>
      <AttachmentList ref="listRef" owner-type="cleaning" :owner-id="createdId" category="before" />
    </div>
  </section>
</template>
