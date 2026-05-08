<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { cleaningApi } from '@/api/cleaning';
import type { PendingCleaning } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<PendingCleaning[]>([]);
const reviewable = computed(() => rows.value.filter((r) => r.status === 'cleaned'));
async function fetchRows() { rows.value = await cleaningApi.pending(); }
async function review(id: number, status: 'completed' | 'rejected', comment?: string) { await cleaningApi.review(id, { status, comment }); ElMessage.success('审核成功'); await fetchRows(); }
onMounted(fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>保洁审核</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="reviewable">
      <el-table-column prop="request_id" label="工单 ID" width="100" /><el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column prop="student_name" label="学生" /><el-table-column prop="building_name" label="楼栋" /><el-table-column prop="location_desc" label="位置" min-width="220" />
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column prop="cleaner_name" label="保洁人员" />
      <el-table-column label="操作" width="150"><template #default="{ row }"><ReviewButtons with-comment approve-text="完成" @approve="(c) => review(row.request_id, 'completed', c)" @reject="(c) => review(row.request_id, 'rejected', c)" /></template></el-table-column>
    </el-table>
  </section>
</template>
