<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { managerApi } from '@/api/manager';
import type { PendingLateReturn } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import { formatDate, formatDateTime } from '@/utils/format';

const rows = ref<PendingLateReturn[]>([]);
async function fetchRows() { rows.value = await managerApi.lateReturns(); }
async function review(id: number, status: 'approved' | 'rejected') { await managerApi.reviewLateReturn(id, status); ElMessage.success('审批成功'); await fetchRows(); }
onMounted(fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>待审批晚归记录</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="rows">
      <el-table-column prop="id" label="ID" width="80" /><el-table-column prop="student_name" label="学生" />
      <el-table-column label="晚归日期"><template #default="{ row }">{{ formatDate(row.return_date) }}</template></el-table-column>
      <el-table-column prop="reason" label="原因" min-width="220" /><el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column label="操作" width="150"><template #default="{ row }"><ReviewButtons @approve="review(row.id, 'approved')" @reject="review(row.id, 'rejected')" /></template></el-table-column>
    </el-table>
  </section>
</template>
