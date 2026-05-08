<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { attachmentApi } from '@/api/attachment';
import type { AttachmentMeta } from '@/types';
import AttachmentImage from './AttachmentImage.vue';
import { formatDateTime } from '@/utils/format';

const props = defineProps<{ ownerType: string; ownerId?: number; category?: string }>();
const list = ref<AttachmentMeta[]>([]);
const loading = ref(false);

async function fetchList() {
  if (!props.ownerId) {
    list.value = [];
    return;
  }
  loading.value = true;
  try {
    list.value = await attachmentApi.list({
      owner_type: props.ownerType,
      owner_id: props.ownerId,
      category: props.category,
    });
  } finally {
    loading.value = false;
  }
}

defineExpose({ fetchList });
onMounted(fetchList);
watch(() => [props.ownerType, props.ownerId, props.category], fetchList);
</script>

<template>
  <div v-loading="loading" class="attachment-list">
    <div v-for="item in list" :key="item.id" class="attachment-item">
      <AttachmentImage :id="item.id" />
      <div>
        <div>{{ item.file_name || `附件 ${item.id}` }}</div>
        <div class="muted">{{ item.category }} · {{ formatDateTime(item.uploaded_at) }}</div>
      </div>
    </div>
    <el-empty v-if="!loading && !list.length" description="暂无附件" :image-size="56" />
  </div>
</template>

<style scoped>
.attachment-list {
  display: grid;
  gap: 10px;
}

.attachment-item {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
