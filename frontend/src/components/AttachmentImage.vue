<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue';
import { Picture } from '@element-plus/icons-vue';
import { attachmentApi } from '@/api/attachment';

const props = defineProps<{ id?: number; width?: string; height?: string }>();
const url = ref('');
const loading = ref(false);
const error = ref(false);

async function load() {
  if (!props.id) return;
  if (url.value) URL.revokeObjectURL(url.value);
  loading.value = true;
  error.value = false;
  try {
    const blob = await attachmentApi.blob(props.id);
    url.value = URL.createObjectURL(blob);
  } catch {
    error.value = true;
  } finally {
    loading.value = false;
  }
}

watch(() => props.id, load, { immediate: true });
onBeforeUnmount(() => {
  if (url.value) URL.revokeObjectURL(url.value);
});
</script>

<template>
  <div v-loading="loading" class="attachment-image" :style="{ width: width || '72px', height: height || '72px' }">
    <el-image v-if="url && !error" :src="url" fit="cover" :preview-src-list="[url]" />
    <el-empty v-else-if="error" description="加载失败" :image-size="36" />
    <el-icon v-else><Picture /></el-icon>
  </div>
</template>

<style scoped>
.attachment-image,
.attachment-image :deep(.el-image) {
  display: grid;
  place-items: center;
  overflow: hidden;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(11, 93, 102, 0.06), transparent),
    #fafafa;
  box-shadow: var(--shadow-sm);
}
</style>
