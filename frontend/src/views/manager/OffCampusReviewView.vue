<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { managerApi } from '@/api/manager';
import type { PendingOffCampus } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import { displayBool, formatDateTime } from '@/utils/format';

const rows = ref<PendingOffCampus[]>([]);
async function fetchRows() { rows.value = await managerApi.offCampus(); }
async function review(id: number, status: 'approved' | 'rejected') { await managerApi.reviewOffCampus(id, status); ElMessage.success('审批成功'); await fetchRows(); }
onMounted(fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>待审批校外居住申请</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="rows">
      <el-table-column prop="id" label="ID" width="80" /><el-table-column prop="student_name" label="学生" />
      <el-table-column label="保留床位"><template #default="{ row }">{{ displayBool(row.retain_bed) }}</template></el-table-column>
      <el-table-column prop="destination" label="居住地址" /><el-table-column prop="reason" label="原因" min-width="220" />
      <el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column label="操作" width="150"><template #default="{ row }"><ReviewButtons @approve="review(row.id, 'approved')" @reject="review(row.id, 'rejected')" /></template></el-table-column>
    </el-table>
  </section>
</template>
