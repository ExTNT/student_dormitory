<script setup lang="ts">
import { ElMessage, type UploadRequestOptions } from 'element-plus';
import { attachmentApi } from '@/api/attachment';
import type { AttachmentMeta } from '@/types';

const props = defineProps<{
  ownerType: string;
  ownerId?: number;
  category: string;
  sortOrder?: number;
}>();

const emit = defineEmits<{ success: [meta: AttachmentMeta] }>();

async function upload(options: UploadRequestOptions) {
  if (!props.ownerId) {
    ElMessage.warning('请先创建工单');
    return;
  }
  const file = options.file;
  const meta = await attachmentApi.upload({
    file,
    owner_type: props.ownerType,
    owner_id: props.ownerId,
    category: props.category,
    sort_order: props.sortOrder,
  });
  ElMessage.success('上传成功');
  emit('success', meta);
  options.onSuccess(meta);
}

function beforeUpload(file: File) {
  const okType = ['image/jpeg', 'image/png'].includes(file.type);
  const okSize = file.size <= 5 * 1024 * 1024;
  if (!okType) ElMessage.warning('只允许上传 jpg/png 图片');
  if (!okSize) ElMessage.warning('图片不能超过 5MB');
  return okType && okSize;
}
</script>

<template>
  <el-upload :http-request="upload" :before-upload="beforeUpload" :show-file-list="false" accept="image/jpeg,image/png">
    <el-button type="primary" class="upload-button">上传图片</el-button>
  </el-upload>
</template>

<style scoped>
.upload-button {
  min-width: 96px;
}
</style>
