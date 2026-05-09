<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage, type UploadInstance, type UploadProps, type UploadRequestOptions } from 'element-plus';
import { attachmentApi } from '@/api/attachment';
import type { AttachmentMeta } from '@/types';

const props = defineProps<{
  ownerType: string;
  ownerId?: number;
  category: string;
  sortOrder?: number;
  buttonText?: string;
  tip?: string;
  disabled?: boolean;
}>();

const emit = defineEmits<{ success: [meta: AttachmentMeta] }>();
const uploadRef = ref<UploadInstance>();
const uploading = ref(false);

async function upload(options: UploadRequestOptions) {
  if (!props.ownerId) {
    ElMessage.warning('请先创建工单');
    options.onError(new Error('missing owner id') as any);
    return;
  }
  uploading.value = true;
  try {
    const meta = await attachmentApi.upload({
      file: options.file,
      owner_type: props.ownerType,
      owner_id: props.ownerId,
      category: props.category,
      sort_order: props.sortOrder,
    });
    ElMessage.success('上传成功');
    emit('success', meta);
    options.onSuccess(meta);
    uploadRef.value?.clearFiles();
  } catch (error) {
    options.onError(error as any);
  } finally {
    uploading.value = false;
  }
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const okType = ['image/jpeg', 'image/png'].includes(file.type);
  const okSize = file.size <= 5 * 1024 * 1024;
  if (!okType) ElMessage.warning('只允许上传 jpg/png 图片');
  if (!okSize) ElMessage.warning('图片不能超过 5MB');
  return okType && okSize;
};

const onExceed: UploadProps['onExceed'] = () => {
  ElMessage.warning('一次只能上传 1 张图片');
};
</script>

<template>
  <el-upload
    ref="uploadRef"
    :http-request="upload"
    :before-upload="beforeUpload"
    :on-exceed="onExceed"
    :show-file-list="false"
    :limit="1"
    :disabled="disabled || uploading"
    accept="image/jpeg,image/png"
  >
    <el-button type="primary" class="upload-button" :loading="uploading" :disabled="disabled || uploading">
      {{ buttonText || '上传图片' }}
    </el-button>
    <template v-if="tip" #tip>
      <div class="upload-tip">{{ tip }}</div>
    </template>
  </el-upload>
</template>

<style scoped>
.upload-button {
  min-width: 96px;
}

.upload-tip {
  margin-top: 6px;
  color: var(--text-muted);
  font-size: 12px;
}
</style>
