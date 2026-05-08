<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { repairApi } from '@/api/repair';
import type { PendingRepair } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<PendingRepair[]>([]);
const reviewable = computed(() => rows.value.filter((r) => r.status === 'repaired'));
async function fetchRows() { rows.value = await repairApi.pending(); }
async function review(id: number, status: 'completed' | 'rejected', comment?: string) { await repairApi.review(id, { status, comment }); ElMessage.success('审核成功'); await fetchRows(); }
onMounted(fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>维修审核</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="reviewable">
      <el-table-column prop="request_id" label="工单 ID" width="100" /><el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column prop="student_name" label="学生" /><el-table-column prop="room_number" label="房间" /><el-table-column prop="description" label="描述" min-width="220" />
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column prop="repair_staff_name" label="维修人员" />
      <el-table-column label="操作" width="150"><template #default="{ row }"><ReviewButtons with-comment approve-text="完成" @approve="(c) => review(row.request_id, 'completed', c)" @reject="(c) => review(row.request_id, 'rejected', c)" /></template></el-table-column>
    </el-table>
  </section>
</template>
