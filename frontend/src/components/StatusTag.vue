<script setup lang="ts">
import { statusTagType } from '@/utils/format';

defineProps<{ status: string }>();
</script>

<template>
  <el-tag class="status-tag" :type="statusTagType[status] || 'info'" effect="light">{{ status }}</el-tag>
</template>

<style scoped>
.status-tag {
  min-width: 78px;
  justify-content: center;
  border-radius: 8px;
  font-weight: 800;
  text-transform: uppercase;
}
</style>
