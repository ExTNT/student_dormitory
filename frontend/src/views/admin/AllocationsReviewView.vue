<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { adminApi } from '@/api/admin';
import type { AllocationRequest } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<AllocationRequest[]>([]);
async function fetchRows() { rows.value = await adminApi.allocations(); }
async function review(id: number, status: 'approved' | 'rejected') { await adminApi.reviewAllocation(id, status); ElMessage.success('审批成功'); await fetchRows(); }
onMounted(fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>待审批新生分配</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="rows">
      <el-table-column prop="id" label="ID" width="80" /><el-table-column prop="student_id" label="学生 ID" />
      <el-table-column prop="recommended_room_id" label="推荐房间 ID" /><el-table-column prop="recommended_bed_id" label="推荐床位 ID" />
      <el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column prop="admin_id" label="管理员 ID" />
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column label="处理时间" width="190"><template #default="{ row }">{{ formatDateTime(row.resolved_at) }}</template></el-table-column>
      <el-table-column label="操作" width="150"><template #default="{ row }"><ReviewButtons @approve="review(row.id, 'approved')" @reject="review(row.id, 'rejected')" /></template></el-table-column>
    </el-table>
  </section>
</template>
