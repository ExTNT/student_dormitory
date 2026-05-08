<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { studentApi } from '@/api/student';
import type { MyRequest } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<MyRequest[]>([]);
const loading = ref(false);
const sortedRows = computed(() => [...rows.value].sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at)));

async function fetchRows() {
  loading.value = true;
  try {
    rows.value = await studentApi.requests();
  } finally {
    loading.value = false;
  }
}

onMounted(fetchRows);
</script>

<template>
  <section class="page">
    <div class="toolbar">
      <h2>我的申请总览</h2>
      <el-button @click="fetchRows">刷新</el-button>
    </div>
    <el-table v-loading="loading" :data="sortedRows">
      <el-table-column prop="request_type" label="类型" width="150" />
      <el-table-column prop="request_id" label="申请 ID" width="110" />
      <el-table-column label="状态" width="120"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column prop="detail" label="详情" min-width="260" />
    </el-table>
  </section>
</template>
