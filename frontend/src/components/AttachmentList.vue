<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { attachmentApi } from '@/api/attachment';
import type { AttachmentMeta } from '@/types';
import AttachmentImage from './AttachmentImage.vue';
import { formatDateTime } from '@/utils/format';

const props = defineProps<{ ownerType: string; ownerId?: number; category?: string }>();
const list = ref<AttachmentMeta[]>([]);
const loading = ref(false);
const loaded = ref(false);

async function fetchList() {
  if (!props.ownerId) {
    list.value = [];
    loaded.value = true;
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
    loaded.value = true;
  }
}

defineExpose({ fetchList });
onMounted(fetchList);
watch(() => [props.ownerType, props.ownerId, props.category], fetchList);
</script>

<template>
  <div class="attachment-list">
    <div v-if="loading && !loaded" class="photo-empty">正在加载照片...</div>
    <div v-for="item in list" v-else :key="item.id" class="attachment-item">
      <AttachmentImage :id="item.id" />
      <div>
        <div>{{ item.file_name || `附件 ${item.id}` }}</div>
        <div class="muted">{{ item.category }} · {{ formatDateTime(item.uploaded_at) }}</div>
      </div>
    </div>
    <div v-if="loaded && !loading && !list.length" class="photo-empty">无照片</div>
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
  padding: 10px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.72);
}

.photo-empty {
  padding: 16px;
  border: 1px dashed var(--border-strong);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.54);
  color: var(--text-muted);
  text-align: center;
}
</style>
