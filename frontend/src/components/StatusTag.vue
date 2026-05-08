<script setup lang="ts">
import { statusTagType } from '@/utils/format';

defineProps<{ status: string }>();
</script>

<template>
  <el-tag :type="statusTagType[status] || 'info'">{{ status }}</el-tag>
</template>
