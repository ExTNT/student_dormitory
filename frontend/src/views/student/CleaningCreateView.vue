<script setup lang="ts">
import { nextTick, onMounted, reactive, ref } from 'vue';
import { ElMessage, type UploadFile, type UploadInstance, type UploadProps } from 'element-plus';
import { studentApi } from '@/api/student';
import { attachmentApi } from '@/api/attachment';
import type { Building } from '@/types';
import AttachmentList from '@/components/AttachmentList.vue';

const buildings = ref<Building[]>([]);
const createdId = ref<number>();
const listRef = ref<InstanceType<typeof AttachmentList>>();
const uploadRef = ref<UploadInstance>();
const loading = ref(false);
const beforeFile = ref<File>();
const form = reactive({ building_id: undefined as number | undefined, location_desc: '' });

function selectBeforePhoto(uploadFile: UploadFile) {
  const file = uploadFile.raw;
  if (!file) return;
  const okType = ['image/jpeg', 'image/png'].includes(file.type);
  const okSize = file.size <= 5 * 1024 * 1024;
  if (!okType) {
    beforeFile.value = undefined;
    uploadRef.value?.clearFiles();
    ElMessage.warning('只允许上传 jpg/png 图片');
    return;
  }
  if (!okSize) {
    beforeFile.value = undefined;
    uploadRef.value?.clearFiles();
    ElMessage.warning('图片不能超过 5MB');
    return;
  }
  beforeFile.value = file;
}

function removeBeforePhoto() {
  beforeFile.value = undefined;
}

const onExceed: UploadProps['onExceed'] = () => {
  ElMessage.warning('保洁申请只能上传 1 张保洁前照片');
};

async function submit() {
  if (!form.building_id) return ElMessage.warning('请选择楼栋');
  if (!form.location_desc) return ElMessage.warning('请填写位置描述');
  if (!beforeFile.value) return ElMessage.warning('请先选择保洁前照片');
  loading.value = true;
  try {
    const res = await studentApi.createCleaning({ building_id: form.building_id, location_desc: form.location_desc });
    createdId.value = res.id;
    await attachmentApi.upload({
      file: beforeFile.value,
      owner_type: 'cleaning',
      owner_id: res.id,
      category: 'before',
    });
    ElMessage.success('保洁申请已提交，保洁前照片已上传');
    beforeFile.value = undefined;
    uploadRef.value?.clearFiles();
    form.location_desc = '';
    await nextTick();
    await listRef.value?.fetchList();
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
        <el-form-item label="保洁前照片" required>
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            accept="image/jpeg,image/png"
            :on-change="selectBeforePhoto"
            :on-remove="removeBeforePhoto"
            :on-exceed="onExceed"
          >
            <el-button>选择图片</el-button>
            <template #tip>
              <div class="muted">提交申请前必须选择 jpg/png 图片，最大 5MB。</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交</el-button></el-form-item>
      </el-form>
    </div>
    <div v-if="createdId" class="data-panel">
      <div class="toolbar">
        <strong>保洁前照片</strong>
      </div>
      <AttachmentList ref="listRef" owner-type="cleaning" :owner-id="createdId" category="before" />
    </div>
  </section>
</template>
